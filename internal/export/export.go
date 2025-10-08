package export

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"github.com/willie68/gowillie68/pkg/fileutils"
	"github.com/willie68/osmltools/internal/export/geojsonexporter"
	"github.com/willie68/osmltools/internal/export/gpxexporter"
	"github.com/willie68/osmltools/internal/export/jsonexporter"
	"github.com/willie68/osmltools/internal/export/kmlexporter"
	"github.com/willie68/osmltools/internal/export/nmeaexporter"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/osml"
	"github.com/willie68/osmltools/internal/trackutils"
)

const (
	JSONFormat    = "JSON"
	NMEAFormat    = "NMEA"
	GPXFormat     = "GPX"
	KMLFormat     = "KML"
	KMZFormat     = "KMZ"
	GEOJSONFormat = "GEOJSON"
)

var (
	// ErrUnknownExporter error for unknown exporter
	ErrUnknownExporter = fmt.Errorf("unknown exporter")
	// SupportedFormats all supported export formats
	SupportedFormats = []string{JSONFormat, NMEAFormat, GPXFormat, KMLFormat, KMZFormat, GEOJSONFormat}
)

type formatExporter interface {
	ExportTrack(track model.TrackPoints, output io.Writer) error
}

type checkerSrv interface {
	AnalyseLoggerFile(fr *model.FileResult, lf string) ([]*model.LogLine, error)
	CorrectTimeStamp(ls []*model.LogLine) ([]*model.LogLine, bool, error)
}

type exporter struct {
	log    logging.Logger
	chk    checkerSrv
	exp    formatExporter
	tracks map[string]trackFileData
}

type trackFileData struct {
	Name  string
	Files []string
}

func provide(inj do.Injector) (*exporter, error) {
	return &exporter{
		log:    *logging.New().WithName("Exporter"),
		chk:    do.MustInvokeAs[checkerSrv](inj),
		tracks: make(map[string]trackFileData),
	}, nil
}

func Init(inj do.Injector) {
	do.Provide(inj, provide)
}

// Export get the exporter and execute it on the sd file set
func (e *exporter) Export(sdCardFolder, outputFolder string, files []string, format, name string) error {
	outTempl := filepath.Join(outputFolder, fmt.Sprintf("track_%%04d.%s", strings.ToLower(format)))
	e.log.Infof("exporter called: sd %s, out: %s, format: %s", sdCardFolder, outTempl, format)

	exp, err := e.checkExporter(format)
	if err != nil {
		return err
	}
	e.exp = exp

	fs, err := os.Stat(sdCardFolder)
	if err != nil {
		return err
	}

	if fs.IsDir() && (len(files) == 0) {
		files, err = osml.GetDataFiles(sdCardFolder)
		if err != nil {
			return err
		}
	}

	if !fs.IsDir() {
		files = append(files, filepath.Base(sdCardFolder))
		sdCardFolder = filepath.Dir(sdCardFolder)
	}

	e.log.Infof("Found %d files on sd card", len(files))

	ls, count, processedFiles, err := ReadLogFiles(files, sdCardFolder, e, outTempl, name)
	if err != nil {
		return err
	}

	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		return err
	}
	err = e.exportFile(ls, count, outTempl, name, processedFiles)
	if err != nil {
		return err
	}

	js, err := json.MarshalIndent(e.tracks, "", "  ")
	if err != nil {
		return err
	}
	of := filepath.Join(outputFolder, "tracks.json")
	err = os.WriteFile(of, js, os.ModePerm)
	return err
}

// Export get the exporter and execute it on the sd file set
func (e *exporter) ExportTrack(trackfile, outputfile, format string) error {
	e.log.Infof("track exporter called: track %s, out: %s, format: %s", trackfile, outputfile, format)

	exp, err := e.checkExporter(format)
	if err != nil {
		return err
	}
	e.exp = exp

	if !fileutils.FileExists(trackfile) {
		return fmt.Errorf("the track file %s does not exist", trackfile)
	}

	if model.IsOldTrackVersion(trackfile) {
		return fmt.Errorf("can't export an old track file %s", trackfile)
	}

	track, nmealines, err := trackutils.ReadTrackAndNmea(trackfile)
	if err != nil {
		return err
	}

	lls, err := model.ParseLines2LogLines(nmealines, false)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(outputfile), os.ModePerm)
	if err != nil {
		return err
	}

	err = e.exportTrackFile(lls, track.Name, outputfile)
	if err != nil {
		return err
	}

	return err
}

func ReadLogFiles(files []string, sdCardFolder string, e *exporter, outTempl string, name string) ([]*model.LogLine, int, []string, error) {
	ls := make([]*model.LogLine, 0)
	count := 0
	today := time.Time{}
	processedFiles := make([]string, 0)

	for _, file := range files {
		lf := filepath.Join(sdCardFolder, file)
		e.log.Infof("analysing file: %s", lf)
		lss, err := e.chk.AnalyseLoggerFile(nil, lf)
		if err != nil {
			return nil, 0, nil, err
		}
		lss, _, err = e.chk.CorrectTimeStamp(lss)
		if err != nil {
			return nil, 0, nil, err
		}
		if len(lss) > 0 {
			if today.IsZero() {
				today = lss[0].CorrectTimeStamp
			} else {
				if lss[0].CorrectTimeStamp.Sub(today).Hours() > 24 {
					e.exportFile(ls, count, outTempl, name, processedFiles)
					processedFiles = make([]string, 0)
					ls = make([]*model.LogLine, 0)
					processedFiles = append(processedFiles, lf)
					ls = append(ls, lss...)
					count++
				}
			}
			processedFiles = append(processedFiles, lf)
			ls = append(ls, lss...)
		}
	}

	sort.Slice(ls, func(i, j int) bool {
		return ls[i].CorrectTimeStamp.Before(ls[j].CorrectTimeStamp)
	})
	return ls, count, processedFiles, nil
}

func (e *exporter) exportFile(ls []*model.LogLine, count int, outTempl, name string, filelist []string) error {
	if len(ls) == 0 {
		return nil
	}
	if name == "" {
		name = fmt.Sprintf("Track %04d", count)
	}
	of := fmt.Sprintf(outTempl, count)
	tr := &model.TrackPoints{
		Name:     name,
		LogLines: ls,
	}
	tr, err := model.GetWaypoints(tr)
	if err != nil {
		return err
	}

	fn := filepath.Base(of)
	e.tracks[fn] = trackFileData{
		Name:  name,
		Files: filelist,
	}
	fs, err := os.Create(of)
	if err != nil {
		return err
	}
	defer fs.Close()

	e.log.Infof("exporting %d loglines to %s", len(ls), of)
	return e.exp.ExportTrack(*tr, fs)
}

func (e *exporter) exportTrackFile(ls []*model.LogLine, name, outputfile string) error {
	if len(ls) == 0 {
		return nil
	}
	tr := &model.TrackPoints{
		Name:     name,
		LogLines: ls,
	}
	tr, err := model.GetWaypoints(tr)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outputfile), os.ModePerm); err != nil {
		return err
	}

	fs, err := os.Create(outputfile)
	if err != nil {
		return err
	}
	defer fs.Close()

	e.log.Infof("exporting %d loglines to %s", len(ls), outputfile)
	return e.exp.ExportTrack(*tr, fs)
}

func (e *exporter) checkExporter(format string) (formatExporter, error) {
	switch format {
	case JSONFormat:
		return jsonexporter.New(), nil
	case NMEAFormat:
		return nmeaexporter.New(), nil
	case GPXFormat:
		return gpxexporter.New(), nil
	case KMLFormat:
		return kmlexporter.New(), nil
	case KMZFormat:
		return kmlexporter.New().WithCompressed(true), nil
	case GEOJSONFormat:
		return geojsonexporter.New(), nil
	}
	return nil, ErrUnknownExporter
}

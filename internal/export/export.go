package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"github.com/willie68/osmltools/internal/check"
	"github.com/willie68/osmltools/internal/export/gpxexporter"
	"github.com/willie68/osmltools/internal/export/kmlexporter"
	"github.com/willie68/osmltools/internal/export/nmeaexporter"
	"github.com/willie68/osmltools/internal/interfaces"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/osml"
)

const (
	NMEAFormat = "NMEA"
	GPXFormat  = "GPX"
	KMLFormat  = "KML"
	KMZFormat  = "KMZ"
)

var (
	// ErrUnknownExporter error for unknown exporter
	ErrUnknownExporter = fmt.Errorf("unknown exporter")
	// SupportedFormats all supported export formats
	SupportedFormats = []string{NMEAFormat, GPXFormat, KMLFormat, KMZFormat}
)

type Exporter struct {
	log    logging.Logger
	chk    check.Checker
	exp    interfaces.FormatExporter
	tracks map[string]trackData
}

type trackData struct {
	Name  string
	Files []string
}

func Init(inj do.Injector) {
	exp := Exporter{
		log:    *logging.New().WithName("Exporter"),
		chk:    do.MustInvoke[check.Checker](inj),
		tracks: make(map[string]trackData),
	}
	do.ProvideValue(inj, exp)
}

// Export get the exporter and execute it on the sd file set
func (e *Exporter) Export(sdCardFolder, outputFolder, format, name string) error {
	outTempl := filepath.Join(outputFolder, fmt.Sprintf("track_%%04d.%s", strings.ToLower(format)))
	e.log.Infof("exporter called: sd %s, out: %s, format: %s", sdCardFolder, outTempl, format)

	exp, err := e.checkExporter(format)
	if err != nil {
		return err
	}
	e.exp = exp

	files, err := osml.GetDataFiles(sdCardFolder)
	if err != nil {
		return err
	}
	e.log.Infof("Found %d files on sd card", len(files))

	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		return err
	}
	ls := make([]*model.LogLine, 0)
	count := 0
	today := time.Time{}
	processedFiles := make([]string, 0)

	for _, lf := range files {
		e.log.Infof("analysing file: %s", lf)
		lss, err := e.chk.AnalyseLoggerFile(nil, lf)
		if err != nil {
			return err
		}
		lss, _, err = e.chk.CorrectTimeStamp(lss)
		if err != nil {
			return err
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
	err = e.exportFile(ls, count, outTempl, name, processedFiles)
	if err != nil {
		return err
	}

	js, err := json.Marshal(e.tracks)
	if err != nil {
		return err
	}
	of := filepath.Join(outputFolder, "tracks.json")
	err = os.WriteFile(of, js, os.ModePerm)
	return err
}

func (e *Exporter) exportFile(ls []*model.LogLine, count int, outTempl, name string, filelist []string) error {
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].CorrectTimeStamp.Before(ls[j].CorrectTimeStamp)
	})
	if name == "" {
		name = fmt.Sprintf("Track %04d", count)
	}
	of := fmt.Sprintf(outTempl, count)
	tr := model.Track{
		Name:     name,
		LogLines: ls,
	}

	fn := filepath.Base(of)
	e.tracks[fn] = trackData{
		Name:  name,
		Files: filelist,
	}

	e.log.Infof("exporting %d loglines to %s", len(ls), of)
	return e.exp.ExportTrack(tr, of)
}

func (e *Exporter) checkExporter(format string) (interfaces.FormatExporter, error) {
	switch format {
	case NMEAFormat:
		return nmeaexporter.New(), nil
	case GPXFormat:
		return gpxexporter.New(), nil
	case KMLFormat:
		return kmlexporter.New(), nil
	case KMZFormat:
		return kmlexporter.New().WithCompressed(true), nil
	}
	return nil, ErrUnknownExporter
}

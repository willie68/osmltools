package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/samber/do/v2"
	"github.com/willie68/osmltools/internal/check"
	"github.com/willie68/osmltools/internal/export/geojsonexporter"
	"github.com/willie68/osmltools/internal/export/gpxexporter"
	"github.com/willie68/osmltools/internal/export/jsonexporter"
	"github.com/willie68/osmltools/internal/export/kmlexporter"
	"github.com/willie68/osmltools/internal/export/nmeaexporter"
	"github.com/willie68/osmltools/internal/interfaces"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/osml"
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

type Exporter struct {
	log    logging.Logger
	chk    check.Checker
	exp    interfaces.FormatExporter
	tracks map[string]trackFileData
}

type trackFileData struct {
	Name  string
	Files []string
}

func Init(inj do.Injector) {
	exp := Exporter{
		log:    *logging.New().WithName("Exporter"),
		chk:    do.MustInvoke[check.Checker](inj),
		tracks: make(map[string]trackFileData),
	}
	do.ProvideValue(inj, exp)
}

// Export get the exporter and execute it on the sd file set
func (e *Exporter) Export(sdCardFolder, outputFolder string, files []string, format, name string) error {
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

	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		return err
	}
	ls := make([]*model.LogLine, 0)
	count := 0
	today := time.Time{}
	processedFiles := make([]string, 0)

	for _, file := range files {
		lf := filepath.Join(sdCardFolder, file)
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

	js, err := json.MarshalIndent(e.tracks, "", "  ")
	if err != nil {
		return err
	}
	of := filepath.Join(outputFolder, "tracks.json")
	err = os.WriteFile(of, js, os.ModePerm)
	return err
}

func (e *Exporter) exportFile(ls []*model.LogLine, count int, outTempl, name string, filelist []string) error {
	if len(ls) == 0 {
		return nil
	}
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].CorrectTimeStamp.Before(ls[j].CorrectTimeStamp)
	})
	if name == "" {
		name = fmt.Sprintf("Track %04d", count)
	}
	of := fmt.Sprintf(outTempl, count)
	tr := &model.Track{
		Name:     name,
		LogLines: ls,
	}
	tr, err := e.GetWaypoints(tr)
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

func (e *Exporter) checkExporter(format string) (interfaces.FormatExporter, error) {
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

// GetWaypoints extracts the waypoints from the log lines of the track
func (e *Exporter) GetWaypoints(track *model.Track) (*model.Track, error) {
	e.log.Infof("extracting waypoints from %d log lines", len(track.LogLines))

	track.Waypoints = make([]*model.Waypoint, 0)

	for _, ll := range track.LogLines {
		if ll.NMEAMessage != nil {
			switch ll.NMEAMessage.Prefix() {
			case "GPRMC":
				rmc, ok := ll.NMEAMessage.(nmea.RMC)
				if ok && rmc.Validity == "A" { // only valid
					track.End = &model.Waypoint{
						Lat:   rmc.Latitude,
						Lon:   rmc.Longitude,
						Time:  ll.CorrectTimeStamp,
						Speed: rmc.Speed,
						Ele:   0.0,
					}
					track.Waypoints = append(track.Waypoints, track.End)
					if track.Start == nil {
						track.Start = track.End
					}
				}
			case "GPGGA":
				if track.End != nil {
					gga, ok := ll.NMEAMessage.(nmea.GGA)
					if ok {
						if track.End.Ele == 0.0 {
							track.End.Ele = gga.Altitude
						}
					}
				}
			case "SDDBT":
				if track.End != nil {
					dbt, ok := ll.NMEAMessage.(nmea.DBT)
					if ok {
						depth := dbt.DepthFeet * 0.3048 // convert feet to meters
						if track.End.Depth == 0.0 {
							track.End.Depth = depth
						}
					}
				}
			case "SDDPT":
				if track.End != nil {
					dpt, ok := ll.NMEAMessage.(nmea.DPT)
					if ok {
						if track.End.Depth == 0.0 {
							track.End.Depth = dpt.Depth
						}
					}
				}
			}
		}
	}
	if track.Start != nil {
		track.Start.Name = "Start"
	}
	if track.End != nil {
		track.End.Name = "End"
	}
	return track, nil
}

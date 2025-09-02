package gpxexporter

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/adrianmo/go-nmea"
	"github.com/twpayne/go-gpx"
	"github.com/willie68/osmltools/internal/interfaces"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

var _ interfaces.FormatExporter = &GPXExporter{}

type GPXExporter struct {
	log logging.Logger
}

func New() *GPXExporter {
	return &GPXExporter{
		log: *logging.New().WithName("GPXExporter"),
	}
}

func (e *GPXExporter) ExportTrack(track model.Track, outputfile string) error {
	e.log.Infof("exporting %d loglines to gpx file %s", len(track.LogLines), outputfile)

	g := NewGPX()

	wpts := make([]*gpx.WptType, 0)
	var lastwpt *gpx.WptType
	var startwpt *gpx.WptType

	for _, ll := range track.LogLines {
		if ll.NMEAMessage != nil {
			if ll.NMEAMessage.Prefix() == "GPRMC" {
				rmc, ok := ll.NMEAMessage.(nmea.RMC)
				if ok && rmc.Validity == "A" { // only valid
					lastwpt = &gpx.WptType{
						Lat:   rmc.Latitude,
						Lon:   rmc.Longitude,
						Time:  ll.CorrectTimeStamp,
						Speed: rmc.Speed,
						Ele:   0.0,
					}
					wpts = append(wpts, lastwpt)
					if startwpt == nil {
						startwpt = lastwpt
					}
				}
			}
			if lastwpt != nil {
				if ll.NMEAMessage.Prefix() == "GPGGA" {
					gga, ok := ll.NMEAMessage.(nmea.GGA)
					if ok {
						if lastwpt.Ele == 0.0 {
							lastwpt.Ele = gga.Altitude
						}
					}
				}

				if lastwpt.Extensions == nil {
					depth := 0.0
					if ll.NMEAMessage.Prefix() == "SDDBT" {
						dbt, ok := ll.NMEAMessage.(nmea.DBT)
						if ok {
							depth = dbt.DepthFeet * 0.3048 // convert feet to meters
						}
					}
					if depth == 0.0 && ll.NMEAMessage.Prefix() == "SDDPT" {
						dpt, ok := ll.NMEAMessage.(nmea.DPT)
						if ok {
							depth = dpt.Depth
						}
					}
					if depth != 0.0 {
						lastwpt.Extensions = &gpx.ExtensionsType{
							XML: []byte(fmt.Sprintf("<gpxx:Depth>%f</gpxx:Depth>", depth)),
						}
					}
				}
			}

		}
	}
	if startwpt != nil {
		startwpt.Name = "Start"
	}
	if lastwpt != nil {
		lastwpt.Name = "End"
	}
	g.Trk = []*gpx.TrkType{{
		Name: track.Name,
		TrkSeg: []*gpx.TrkSegType{&gpx.TrkSegType{
			TrkPt: wpts,
		}},
	}}
	g.Wpt = []*gpx.WptType{
		startwpt,
		lastwpt,
	}
	g.XMLAttrs = make(map[string]string)
	g.XMLAttrs["xmlns:gpxx"] = "http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd"

	fs, err := os.Create(outputfile)
	if err != nil {
		return err
	}
	defer fs.Close()

	if _, err := fmt.Fprint(fs, xml.Header); err != nil {
		return err
	}

	if err := g.WriteIndent(fs, "", "  "); err != nil {
		return err
	}

	e.log.Info("output file written")
	return nil
}

func NewGPX() *gpx.GPX {
	return &gpx.GPX{
		Version: "1.1",
		Creator: "ExpertGPS 1.1 - http://www.topografix.com",
	}
}

package kmlexporter

import (
	"fmt"
	"os"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/twpayne/go-kml/v3"
	"github.com/willie68/osmltools/internal/interfaces"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

var _ interfaces.FormatExporter = &KMLExporter{}

type KMLExporter struct {
	log logging.Logger
}

type waypoint struct {
	Name  string
	Lat   float64
	Lon   float64
	Time  time.Time
	Speed float64
	Ele   float64
	Depth float64
}

func New() *KMLExporter {
	return &KMLExporter{
		log: *logging.New().WithName("KMLExporter"),
	}
}

func (e *KMLExporter) ExportTrack(track model.Track, outputfile string) error {
	e.log.Infof("exporting %d loglines to kml file %s", len(track.LogLines), outputfile)

	wpts := make([]*waypoint, 0)
	var lastwpt *waypoint
	var startwpt *waypoint

	for _, ll := range track.LogLines {
		if ll.NMEAMessage != nil {
			if ll.NMEAMessage.Prefix() == "GPRMC" {
				rmc, ok := ll.NMEAMessage.(nmea.RMC)
				if ok && rmc.Validity == "A" { // only valid
					lastwpt = &waypoint{
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
				if depth != 0.0 && lastwpt.Depth == 0.0 {
					lastwpt.Depth = depth
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

	fs, err := os.Create(outputfile)
	if err != nil {
		return err
	}
	defer fs.Close()

	kos := make([]kml.Coordinate, 0)
	for _, wpt := range wpts {
		kos = append(kos, kml.Coordinate{
			Lon: wpt.Lon,
			Lat: wpt.Lat,
			Alt: wpt.Ele,
		})
	}

	gxkos := make([]kml.Element, 0)
	for _, wpt := range wpts {
		gxkos = append(gxkos, kml.GxCoord(kml.Coordinate{
			Lon: wpt.Lon,
			Lat: wpt.Lat,
			Alt: -wpt.Depth,
		}))
	}
	k := kml.KML(
		kml.Document(
			kml.Name(track.Name),
			kml.Description(fmt.Sprintf("Exported with osmltools - %d points", len(kos))),
			kml.Placemark(
				kml.Name(track.Name),
				kml.LineString(kml.Coordinates(kos...)),
			),
			kml.Placemark(
				kml.Name("Water depth profile"),
				kml.GxTrack(gxkos...),
			),
		),
	)
	if err := k.WriteIndent(fs, "", "  "); err != nil {
		return err
	}

	e.log.Info("output file written")
	return nil
}

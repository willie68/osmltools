package kmlexporter

import (
	"fmt"
	"io"

	"github.com/twpayne/go-kml/v3"
	"github.com/willie68/osmltools/internal/interfaces"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

var _ interfaces.FormatExporter = &KMLExporter{}

type KMLExporter struct {
	log        logging.Logger
	compressed bool
}

// New returns a new KMLExporter
func New() *KMLExporter {
	return &KMLExporter{
		log:        *logging.New().WithName("KMLExporter"),
		compressed: false,
	}
}

// WithCompressed sets if the output should be compressed to a kmz file
func (e *KMLExporter) WithCompressed(compressed bool) *KMLExporter {
	e.compressed = compressed
	return e
}

// ExportTrack exports the given track to a kml or kmz file
func (e *KMLExporter) ExportTrack(track model.TrackPoints, output io.Writer) error {
	kos := make([]kml.Coordinate, 0)
	for _, wpt := range track.Waypoints {
		kos = append(kos, kml.Coordinate{
			Lon: wpt.Lon,
			Lat: wpt.Lat,
			Alt: wpt.Ele,
		})
	}

	gxkos := make([]kml.Element, 0)
	for _, wpt := range track.Waypoints {
		gxkos = append(gxkos, kml.GxCoord(kml.Coordinate{
			Lon: wpt.Lon,
			Lat: wpt.Lat,
			Alt: -wpt.Depth,
		}))
	}

	kd := kml.KML(kml.Document(
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
	))

	if e.compressed {
		if err := kml.WriteKMZ(output, map[string]any{"doc.kml": kd}); err != nil {
			return err
		}
	} else {
		if err := kml.KML(kd).WriteIndent(output, "", "  "); err != nil {
			return err
		}
	}
	return nil
}

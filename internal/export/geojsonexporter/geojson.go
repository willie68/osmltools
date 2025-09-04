package geojsonexporter

import (
	"os"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/willie68/osmltools/internal/interfaces"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

var _ interfaces.FormatExporter = &GeoJSONExporter{}

type GeoJSONExporter struct {
	log logging.Logger
}

func New() *GeoJSONExporter {
	return &GeoJSONExporter{
		log: *logging.New().WithName("GeoJSONExporter"),
	}
}

func (e *GeoJSONExporter) ExportTrack(track model.Track, outputfile string) error {
	e.log.Infof("exporting %d loglines to geojson file %s", len(track.LogLines), outputfile)

	coords := make([]geom.Coord, 0)
	depths := make([]float64, 0)
	speeds := make([]float64, 0)
	times := make([]string, 0)
	// collect coordinates, depths, speeds and times
	for _, wpt := range track.Waypoints {
		coords = append(coords, geom.Coord{wpt.Lon, wpt.Lat, wpt.Ele})
		depths = append(depths, wpt.Depth)
		speeds = append(speeds, wpt.Speed)
		times = append(times, wpt.Time.Format("2006-01-02T15:04:05Z"))
	}

	ls := geom.NewLineString(geom.XYZ)
	ls.SetCoords(coords)
	tf := &geojson.Feature{
		Geometry: ls,
		Properties: map[string]any{
			"name":   track.Name,
			"depths": depths,
			"speeds": speeds,
			"times":  times,
		},
	}
	ts := &geojson.Feature{
		Geometry: geom.NewPointFlat(geom.XY, []float64{track.Start.Lon, track.Start.Lat}),
		Properties: map[string]any{
			"name": "start",
		},
	}
	te := &geojson.Feature{
		Geometry: geom.NewPointFlat(geom.XY, []float64{track.End.Lon, track.End.Lat}),
		Properties: map[string]any{
			"name": "end",
		},
	}

	fc := geojson.FeatureCollection{
		Features: []*geojson.Feature{
			tf,
			ts,
			te,
		},
	}

	rawJSON, err := fc.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(outputfile, rawJSON, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

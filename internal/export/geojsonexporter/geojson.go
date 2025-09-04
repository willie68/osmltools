package geojsonexporter

import (
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
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

	coords := make([]orb.Point, 0)
	depths := make([]float64, 0)
	speeds := make([]float64, 0)
	times := make([]string, 0)

	// collect coordinates, depths, speeds and times
	for _, wpt := range track.Waypoints {
		coords = append(coords, orb.Point{wpt.Lon, wpt.Lat})
		depths = append(depths, wpt.Depth)
		speeds = append(speeds, wpt.Speed)
		times = append(times, wpt.Time.Format("2006-01-02T15:04:05Z"))
	}

	fc := geojson.NewFeatureCollection()
	f := geojson.NewFeature(orb.LineString(coords))
	f.Properties["depth"] = depths
	f.Properties["speed"] = speeds
	f.Properties["time"] = times
	f.Properties["name"] = track.Name
	fc.Append(f)
	fs := geojson.NewFeature(orb.Point{track.Start.Lon, track.Start.Lat})
	fs.Properties["name"] = "start"
	fc.Append(fs)
	fe := geojson.NewFeature(orb.Point{track.End.Lon, track.End.Lat})
	fe.Properties["name"] = "end"
	fc.Append(fe)
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

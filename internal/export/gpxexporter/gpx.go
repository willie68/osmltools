package gpxexporter

import (
	"encoding/xml"
	"fmt"
	"os"

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

	g.Trk = []*gpx.TrkType{{
		Name: track.Name,
		TrkSeg: []*gpx.TrkSegType{{
			TrkPt: e.ConvertToWPTTypes(track.Waypoints),
		}},
	}}
	g.Wpt = []*gpx.WptType{
		e.ConvertToWPTType(track.Start),
		e.ConvertToWPTType(track.End),
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

func (e *GPXExporter) ConvertToWPTType(wpt *model.Waypoint) *gpx.WptType {
	if wpt == nil {
		return nil
	}

	gwpt := &gpx.WptType{
		Lat:   wpt.Lat,
		Lon:   wpt.Lon,
		Time:  wpt.Time,
		Ele:   wpt.Ele,
		Name:  wpt.Name,
		Speed: wpt.Speed,
	}
	if wpt.Depth != 0.0 {
		gwpt.Extensions = &gpx.ExtensionsType{
			XML: []byte(fmt.Sprintf("<gpxx:Depth>%f</gpxx:Depth>", wpt.Depth)),
		}
	}
	return gwpt
}

func (e *GPXExporter) ConvertToWPTTypes(wpts []*model.Waypoint) []*gpx.WptType {
	if wpts == nil {
		return nil
	}
	gwpts := make([]*gpx.WptType, 0)
	for _, wpt := range wpts {
		gwpts = append(gwpts, e.ConvertToWPTType(wpt))
	}
	return gwpts
}

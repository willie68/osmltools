package jsonexporter

import (
	"encoding/json"
	"io"

	"github.com/willie68/osmltools/internal/interfaces"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

var _ interfaces.FormatExporter = &JSONExporter{}

type JSONExporter struct {
	log        logging.Logger
	compressed bool
}

// New returns a new KMLExporter
func New() *JSONExporter {
	return &JSONExporter{
		log:        *logging.New().WithName("JSONExporter"),
		compressed: false,
	}
}

// ExportTrack exports the given track to a kml or kmz file
func (e *JSONExporter) ExportTrack(track model.Track, output io.Writer) error {
	js, err := json.Marshal(track.Waypoints)
	if err != nil {
		return err
	}
	_, err = output.Write(js)
	return err
}

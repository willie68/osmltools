package nmeaexporter

import (
	"fmt"
	"io"

	"github.com/willie68/osmltools/internal/interfaces"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

var _ interfaces.FormatExporter = &NMEAExporter{}

type NMEAExporter struct {
	log logging.Logger
}

func New() *NMEAExporter {
	return &NMEAExporter{
		log: *logging.New().WithName("NMEAExporter"),
	}
}

func (e *NMEAExporter) ExportTrack(track model.TrackPoints, output io.Writer) error {
	for _, ll := range track.LogLines {
		fmt.Fprintln(output, ll.NMEAString())
	}
	return nil
}

package nmeaexporter

import (
	"fmt"
	"os"

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

func (e *NMEAExporter) ExportTrack(track model.Track, outputfile string) error {
	e.log.Infof("exporting %d loglines to nmea file %s", len(track.LogLines), outputfile)

	fs, err := os.Create(outputfile)
	if err != nil {
		return err
	}
	defer fs.Close()
	for _, ll := range track.LogLines {
		fmt.Fprintln(fs, ll.NMEAString())
	}
	e.log.Info("output file written")
	return nil
}

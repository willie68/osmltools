package export

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"github.com/willie68/osmltools/internal/check"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/osml"
)

const (
	NMEAFormat = "NMEA"
	GPXFormat  = "GPX"
	KMLFormat  = "KML"
)

var (
	SupportedFormats = []string{NMEAFormat, GPXFormat, KMLFormat}
)

type Exporter struct {
	log logging.Logger
	chk check.Checker
}

func Init(inj do.Injector) {
	exp := Exporter{
		log: *logging.New().WithName("Exporter"),
		chk: do.MustInvoke[check.Checker](inj),
	}
	do.ProvideValue(inj, exp)
}

// Export get the exporter and execute it on the sd file set
func (e *Exporter) Export(sdCardFolder, outputFolder, format string) error {
	outTempl := filepath.Join(outputFolder, fmt.Sprintf("track_%%04d.%s", strings.ToLower(format)))
	e.log.Infof("exporter called: sd %s, out: %s, format: %s", sdCardFolder, outTempl, format)
	files, err := osml.GetDataFiles(sdCardFolder)
	if err != nil {
		return err
	}
	e.log.Infof("Found %d files on sd card", len(files))

	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		return err
	}
	ls := make([]*model.LogLine, 0)
	count := 0
	today := time.Time{}

	for _, lf := range files {
		e.log.Infof("analysing file: %s", lf)
		lss, err := e.chk.AnalyseLoggerFile(nil, lf)
		if err != nil {
			return err
		}
		lss, err = e.chk.CorrectTimeStamp(lss)
		if err != nil {
			return err
		}
		if len(lss) > 0 {
			if today.IsZero() {
				today = lss[0].CorrectTimeStamp
			} else {
				if lss[0].CorrectTimeStamp.Sub(today).Hours() > 24 {
					ls = append(ls, lss...)
					e.exportFile(ls, count, outTempl, format)
					ls = make([]*model.LogLine, 0)
					count++
				}
			}
		}
		ls = append(ls, lss...)
	}
	return e.exportFile(ls, count, outTempl, format)
}

func (e *Exporter) exportFile(ls []*model.LogLine, count int, outTempl string, format string) error {
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].CorrectTimeStamp.Before(ls[j].CorrectTimeStamp)
	})

	fs, err := os.Create(fmt.Sprintf(outTempl, count))
	if err != nil {
		return err
	}
	defer fs.Close()
	for _, ll := range ls {
		fmt.Fprintln(fs, ll.NMEAString())
	}
	e.log.Infof("output file written with %d sentences", len(ls))
	return nil
}

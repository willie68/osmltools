package export

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/samber/do/v2"
	utils "github.com/willie68/gowillie68/pkg"
	"github.com/willie68/osmltools/internal/check"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/osml"
)

const ()

var ()

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

func (e *Exporter) Export(sdCardFolder, outputFile, format string) error {
	e.log.Infof("exporter called: sd %s, out: %s, format: %s", sdCardFolder, outputFile, format)
	files, err := osml.GetDataFiles(sdCardFolder)
	if err != nil {
		return err
	}
	e.log.Infof("Found %d files on sd card", len(files))

	err = os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	if err != nil {
		return err
	}
	ls := make([]model.LogLine, 0)

	for _, lf := range files {
		lss, err := e.chk.AnalyseLoggerFile(nil, lf)
		if err != nil {
			return err
		}
		ls = append(ls, lss...)
	}

	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Timestamp.Before(ls[j].Timestamp)
	})
	if utils.FileExists(outputFile) {
		os.Remove(outputFile)
	}
	fs, err := os.OpenFile(outputFile, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fs.Close()
	for _, ll := range ls {
		fmt.Fprintf(fs, "%s\r\n", ll.String())
	}
	e.log.Infof("output file written with %d sentences", len(ls))
	return nil
}

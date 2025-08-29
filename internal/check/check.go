package check

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/samber/do/v2"
	utils "github.com/willie68/gowillie68/pkg"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/osml"
	"github.com/willie68/osmltools/internal/osmlnmea"
)

const (
	fmtDateOnly = "2006-01-02"
)

var (
	ErrOutputfileAlreadyExists = errors.New("the output file already exists")
)

type Checker struct {
	files       []string
	log         logging.Logger
	UnknownTags int
	ErrorTags   int
}

func Init(inj do.Injector) {
	chk := Checker{
		log: *logging.New().WithName("Checker"),
	}
	do.ProvideValue(inj, chk)
}

func (c *Checker) Check(sdCardFolder, outputFolder string, overwrite, report bool) error {
	c.log.Infof("check called: sd %s, out: %s", sdCardFolder, outputFolder)
	result := model.NewCheckResult()
	files, err := osml.GetDataFiles(sdCardFolder)
	if err != nil {
		return err
	}
	c.files = make([]string, 0)
	c.files = append(c.files, files...)
	c.log.Infof("Found %d files on sd card", len(files))

	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		return err
	}

	for _, lf := range files {
		err1 := c.checkFile(lf, result, outputFolder, overwrite)
		if err1 != nil {
			return err1
		}
	}
	if report {
		err = c.WriteResult(outputFolder, *result)
		if err != nil {
			return err
		}
	}
	c.log.Infof("all files parsed with %d errors and %d unknown tags", c.ErrorTags, c.UnknownTags)
	return nil
}

func (c *Checker) checkFile(loggerfile string, result *model.CheckResult, outputFolder string, overwrite bool) error {
	ofn := filepath.Base(loggerfile)
	fr := model.NewFileResult().WithOrigin(ofn)
	result.WithFileResult(ofn, fr)
	set := c.ErrorTags
	sut := c.UnknownTags
	c.log.Infof("start with file %s", loggerfile)
	ls, err := c.AnalyseLoggerFile(fr, loggerfile)
	if err != nil {
		return err
	}
	err = c.outputToFolder(fr, loggerfile, outputFolder, ls, overwrite)
	if err != nil {
		return err
	}
	c.log.Infof("file parsed with %d errors and %d unknown tags", c.ErrorTags-set, c.UnknownTags-sut)
	return nil
}

func (c *Checker) AnalyseLoggerFile(fr *model.FileResult, lf string) ([]model.LogLine, error) {
	f, err := os.Open(lf)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	ls := make([]model.LogLine, 0)
	count := 0
	// Loop through the file and read each line
	for scanner.Scan() {
		count++
		line := scanner.Text() // Get the line as a string
		ll, ok, err := model.ParseLogLine(line)
		if err != nil {
			if ok {
				c.UnknownTags++
				ls := fmt.Sprintf("warning unknown NMEA Tag in line %d: %s", count, line)
				c.log.Debug(ls)
				model.AddWarning(fr, ls)
			} else {
				c.ErrorTags++
				ls := fmt.Sprintf("error in line %d: %s: %v", count, line, err)
				c.log.Debug(ls)
				model.AddError(fr, ls)
			}
		}
		if ok {
			ls = append(ls, *ll)
		}
	}

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		c.log.Fatalf("error reading file: %v", err)
		return nil, err
	}

	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Timestamp.Before(ls[j].Timestamp)
	})
	return ls, nil
}

func (c *Checker) outputToFolder(fr *model.FileResult, lf, of string, ls []model.LogLine, overwrite bool) error {
	vesselID, ft := c.getFileInfo(ls)
	filedate := ft.Format(fmtDateOnly)
	ofn := fmt.Sprintf("%d-%s-%s.nmea", vesselID, utils.FileNameWithoutExtension(filepath.Base(lf)), filedate)
	off := filepath.Join(of, ofn)
	if !overwrite && utils.FileExists(off) {
		return errors.Join(ErrOutputfileAlreadyExists, fmt.Errorf("output file name: %s", off))
	}
	fr.WithFilename(ofn).WithVesselID(vesselID).WithCreated(ft)
	fo, err := os.OpenFile(off, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		err := fo.Close()
		if err != nil {
			c.log.Fatalf("error closing output file: %v", err)
		}
		os.Chtimes(off, time.Unix(0, 0), ft)
	}()

	for _, ll := range ls {
		fo.WriteString(fmt.Sprintf("%s;%s;", ll.Timestamp, ll.Channel))
		if ll.NMEAMessage != nil {
			fo.WriteString(ll.NMEAMessage.String())
		} else {
			fo.WriteString(ll.Unknown)
		}
		fo.WriteString("\r\n")
	}
	c.log.Infof("writing clean up file to %s", off)
	return nil
}

func (c *Checker) getFileInfo(ls []model.LogLine) (vesselID int64, creationDate time.Time) {
	creationDate = time.Now()
	vesselID = int64(0)
	dateFound := false
	vesselFound := false
	for _, ll := range ls {
		if ll.NMEAMessage != nil {
			tt, ok := c.processRMC(ll, dateFound)
			if ok {
				creationDate = tt
				dateFound = true
			}
			vid, ok := c.processCFG(ll, vesselFound)
			if ok {
				vesselID = vid
				vesselFound = true
			}
		}
		if vesselFound && dateFound {
			break
		}
	}
	return
}

func (c *Checker) processRMC(ll model.LogLine, found bool) (creationDate time.Time, ok bool) {
	ok = false
	if ll.NMEAMessage.DataType() == "RMC" {
		rmc, tok := ll.NMEAMessage.(nmea.RMC)
		if tok && !found {
			creationDate = nmea.DateTime(2000, rmc.Date, rmc.Time)
			ok = true
		}
	}
	return
}

func (c *Checker) processCFG(ll model.LogLine, found bool) (vesselID int64, ok bool) {
	ok = false
	if ll.NMEAMessage.DataType() == "OSMCFG" {
		cfg, tok := ll.NMEAMessage.(osmlnmea.OSMCFG)
		if tok && !found {
			vesselID = cfg.VesselID
			ok = true
		}
	}
	return
}

func (c *Checker) WriteResult(of string, res model.CheckResult) error {
	fn := filepath.Join(of, "report.json")
	return os.WriteFile(fn, []byte(res.String()), os.ModePerm)
}

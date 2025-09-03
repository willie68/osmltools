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

// Check checks the sd card folder and writes the cleaned up NMEA files to the output folder
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

// checkFile checks a single logger file and writes the cleaned up NMEA file to the output folder
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
	ls, err = c.CorrectTimeStamp(ls)
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

// AnalyseLoggerFile analyses a single logger file and returns the log lines found
func (c *Checker) AnalyseLoggerFile(fr *model.FileResult, lf string) ([]*model.LogLine, error) {
	f, err := os.Open(lf)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	ls := make([]*model.LogLine, 0)
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
			ls = append(ls, ll)
		}
	}

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		c.log.Fatalf("error reading file: %v", err)
		return nil, err
	}

	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Duration < ls[j].Duration
	})
	return ls, nil
}

// CorrectTimeStamp corrects the timestamp of the log lines. It searches for the first RMC sentence and uses its time as reference.
func (c *Checker) CorrectTimeStamp(ls []*model.LogLine) ([]*model.LogLine, error) {
	// get first RMC for getting the right time information
	td := time.Time{}
	found := false
	for _, ll := range ls {
		ts, ok := c.getRMCTime(ll, td)
		if ok {
			found = true
			ll.CorrectTimeStamp = ts
			td = ts.Add(-ll.Duration)
			break
		}
	}
	if found {
		c.log.Infof("reference time found: %s", td.String())
	} else {
		c.log.Infof("no reference time found")
	}

	// setting the timestamp right
	for _, ll := range ls {
		ts, ok := c.getRMCTime(ll, td)
		if ok {
			ll.CorrectTimeStamp = ts
			td = ts.Add(-ll.Duration)
		} else {
			ll.CorrectTimeStamp = td.Add(ll.Duration)
		}
	}

	return ls, nil
}

func (c *Checker) getRMCTime(ll *model.LogLine, ts time.Time) (time.Time, bool) {
	newTime := false
	if ll.NMEAMessage != nil {
		if ll.NMEAMessage.Prefix() == "GPRMC" {
			rmc, ok := ll.NMEAMessage.(nmea.RMC)
			if ok {
				ts = nmea.DateTime(0, rmc.Date, rmc.Time)
				newTime = true
			}
		}
	}
	return ts, newTime
}

func (c *Checker) outputToFolder(fr *model.FileResult, lf, of string, ls []*model.LogLine, overwrite bool) error {
	vesselID, ft := c.getFileInfo(ls)
	filedate := ft.Format(fmtDateOnly)
	ofn := fmt.Sprintf("%d-%s-%s.nmea", vesselID, utils.FileNameWithoutExtension(filepath.Base(lf)), filedate)
	off := filepath.Join(of, ofn)
	if !overwrite && utils.FileExists(off) {
		return errors.Join(ErrOutputfileAlreadyExists, fmt.Errorf("output file name: %s", off))
	}
	fr.WithFilename(ofn).WithVesselID(vesselID).WithCreated(ft)
	fo, err := os.Create(off)
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
		fo.WriteString(fmt.Sprintf("%s\r\n", ll.NMEAString()))
	}
	c.log.Infof("writing clean up file to %s", off)
	return nil
}

func (c *Checker) getFileInfo(ls []*model.LogLine) (vesselID int64, creationDate time.Time) {
	creationDate = time.Time{}
	vesselID = int64(0)
	dateFound := false
	vesselFound := false
	for _, ll := range ls {
		if ll.NMEAMessage != nil {
			tt, ok := c.processRMC(*ll, dateFound)
			if ok {
				creationDate = tt
				dateFound = true
			}
			vid, ok := c.processCFG(*ll, vesselFound)
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

// WriteResult writes the check result to a json file in the output folder
func (c *Checker) WriteResult(of string, res model.CheckResult) error {
	fn := filepath.Join(of, "report.json")
	return os.WriteFile(fn, []byte(res.String()), os.ModePerm)
}

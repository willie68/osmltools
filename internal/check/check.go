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
	"github.com/willie68/gowillie68/pkg/fileutils"
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
	ErrNotFound                = errors.New("the requested object was not found")
)

type checker struct {
	files []string
	log   logging.Logger
}

func provide(inj do.Injector) (*checker, error) {
	return &checker{
		log: *logging.New().WithName("Checker"),
	}, nil
}

func Init(inj do.Injector) {
	do.Provide(inj, provide)
}

// Check checks the sd card folder and writes the cleaned up NMEA files to the output folder
func (c *checker) Check(sdCardFolder, outputFolder string, overwrite, report bool) (*model.CheckResult, error) {
	c.log.Infof("check called: sd %s, out: %s", sdCardFolder, outputFolder)
	fs, err := os.Stat(sdCardFolder)
	if err != nil {
		return nil, osml.ErrWrongCardFolder
	}
	result := model.NewCheckResult()
	c.files = make([]string, 0)
	if fs.IsDir() {
		files, err := osml.GetDataFiles(sdCardFolder)
		if err != nil {
			return nil, err
		}
		c.files = append(c.files, files...)
	} else {
		// only a single file should be checked
		c.files = append(c.files, sdCardFolder)
	}
	c.log.Infof("Found %d files on sd card", len(c.files))

	if outputFolder != "" {
		err = os.MkdirAll(outputFolder, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	for _, lf := range c.files {
		err1 := c.checkFile(lf, result, outputFolder, overwrite)
		if err1 != nil {
			return nil, err1
		}
	}
	if report && outputFolder != "" {
		err = c.WriteResult(outputFolder, *result)
		if err != nil {
			return nil, err
		}
	}
	c.log.Infof("all files parsed with %d errors and %d unknown tags", result.ErrorTags, result.UnknownTags)
	result.Calc()
	return result, nil
}

// checkFile checks a single logger file and writes the cleaned up NMEA file to the output folder
func (c *checker) checkFile(loggerfile string, result *model.CheckResult, outputFolder string, overwrite bool) error {
	ofn := filepath.Base(loggerfile)
	fr := model.NewFileResult().WithOrigin(ofn)
	result.WithFileResult(ofn, fr)
	c.log.Infof("start with file %s", loggerfile)
	ls, err := c.AnalyseLoggerFile(fr, loggerfile)
	if err != nil {
		return err
	}
	ls, ok, err := c.CorrectTimeStamp(ls)
	if err != nil {
		return err
	}
	ver, err := c.GetVersion(ls)
	if err == nil {
		fr.Version = ver
	} else {
		fr.Version = "n.N."
	}
	if len(ls) == 0 {
		c.log.Infof("no valid nmea lines found in file %s", loggerfile)
		fr.AddErrors("I", fmt.Sprintf("no valid nmea lines found in file %s", loggerfile))
		return nil
	}
	fr.FirstTimestamp = ls[0].CorrectTimeStamp
	fr.LastTimestamt = ls[len(ls)-1].CorrectTimeStamp
	if !ok {
		c.log.Infof("no valid time stamp found in file %s", loggerfile)
		fr.AddErrors("I", fmt.Sprintf("no valid time stamp found in file %s", loggerfile))
	}
	if outputFolder != "" {
		err = c.outputToFolder(fr, loggerfile, outputFolder, ls, overwrite)
		if err != nil {
			return err
		}
	}
	result.ErrorTags += fr.ErrorTags
	result.UnknownTags += fr.UnknownTags
	c.log.Infof("file parsed with %d errors and %d unknown tags", fr.ErrorTags, fr.UnknownTags)
	return nil
}

func (c *checker) GetVersion(ls []*model.LogLine) (string, error) {
	for _, ll := range ls {
		if ll.NMEAMessage != nil {
			if ll.NMEAMessage.Prefix() == "POSMST" {
				st, ok := ll.NMEAMessage.(osmlnmea.OSMST)
				if ok {
					return st.Version, nil
				}
			}
		}
	}
	return "", ErrNotFound
}

// AnalyseLoggerFile analyses a single logger file and returns the log lines found
func (c *checker) AnalyseLoggerFile(fr *model.FileResult, lf string) ([]*model.LogLine, error) {
	fs, err := os.Stat(lf)
	if err != nil {
		return nil, err
	}
	if fr != nil {
		fr.Size = fs.Size()
	}
	f, err := os.Open(lf)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	ls := make([]*model.LogLine, 0)
	count := 0
	erTags := 0
	unTags := 0
	// Loop through the file and read each line
	for scanner.Scan() {
		count++
		line := scanner.Text() // Get the line as a string
		ll, ok, err := model.ParseLogLine(line)
		if err != nil {
			if ok {
				unTags++
				ls := fmt.Sprintf("warning unknown NMEA Tag in line %d: %s", count, line)
				c.log.Debug(ls)
				model.AddWarning(fr, ls)
			} else {
				erTags++
				ls := fmt.Sprintf("error in line %d: %s: %v", count, line, err)
				c.log.Debug(ls)
				ch := "I"
				if ll != nil {
					ch = ll.Channel
				}
				model.AddError(ch, fr, ls)
			}
		}
		if ok {
			ls = append(ls, ll)
		}
	}
	if fr != nil {
		fr.DatagramCount = count
		fr.ErrorTags += erTags
		fr.UnknownTags += unTags
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
func (c *checker) CorrectTimeStamp(ls []*model.LogLine) ([]*model.LogLine, bool, error) {
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

	return ls, found, nil
}

func (c *checker) getRMCTime(ll *model.LogLine, ts time.Time) (time.Time, bool) {
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

func (c *checker) outputToFolder(fr *model.FileResult, lf, of string, ls []*model.LogLine, overwrite bool) error {
	vesselID, ft := c.getFileInfo(ls)
	filedate := ft.Format(fmtDateOnly)
	ofn := fmt.Sprintf("%d-%s-%s.nmea", vesselID, fileutils.FileNameWithoutExtension(filepath.Base(lf)), filedate)
	off := filepath.Join(of, ofn)
	if !overwrite && fileutils.FileExists(off) {
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

func (c *checker) getFileInfo(ls []*model.LogLine) (vesselID int64, creationDate time.Time) {
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

func (c *checker) processRMC(ll model.LogLine, found bool) (creationDate time.Time, ok bool) {
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

func (c *checker) processCFG(ll model.LogLine, found bool) (vesselID int64, ok bool) {
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
func (c *checker) WriteResult(of string, res model.CheckResult) error {
	fn := filepath.Join(of, "report.json")
	return os.WriteFile(fn, []byte(res.JSON()), os.ModePerm)
}

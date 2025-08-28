package check

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/samber/do/v2"
	utils "github.com/willie68/gowillie68/pkg"
	"github.com/willie68/gowillie68/pkg/fileutils"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

var (
	ErrWrongCardFolder         = errors.New("sd card folder not exists")
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

func (c *Checker) Check(sdCardFolder, outputFolder string, overwrite bool) error {
	c.log.Infof("check called: sd %s, out: %s", sdCardFolder, outputFolder)

	files, err := c.getDataFiles(sdCardFolder)
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
		set := c.ErrorTags
		sut := c.UnknownTags
		c.log.Infof("start with file %s", lf)
		ls, err := c.analyseLoggerFile(lf, outputFolder)
		if err != nil {
			return err
		}
		err = c.outputToFolder(lf, outputFolder, ls)
		if err != nil {
			return err
		}
		c.log.Infof("file parsed with %d errors and %d unknown tags", c.ErrorTags-set, c.UnknownTags-sut)
	}

	c.log.Infof("all files parsed with %d errors and %d unknown tags", c.ErrorTags, c.UnknownTags)
	return nil
}

func (c *Checker) analyseLoggerFile(lf, of string) ([]model.LogLine, error) {
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
				c.log.Debugf("warning unknown NMEA Tag in line %d: %s", count, line)
			} else {
				c.ErrorTags++
				c.log.Errorf("error in line %d: %s: %v", count, line, err)
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

	slices.SortFunc(ls, func(a, b model.LogLine) int {
		return strings.Compare(strings.ToLower(a.Timestamp), strings.ToLower(b.Timestamp))
	})
	return ls, nil
}

func (c *Checker) outputToFolder(lf, of string, ls []model.LogLine) error {
	off := filepath.Join(of, utils.FileNameWithoutExtension(filepath.Base(lf))+".nmea")
	if utils.FileExists(off) {
		return errors.Join(ErrOutputfileAlreadyExists, fmt.Errorf("output file name: %s", off))
	}
	fo, err := os.OpenFile(off, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		err := fo.Close()
		if err != nil {
			c.log.Fatalf("error closing output file: %v", err)
		}
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

func (c *Checker) getDataFiles(sdCardFolder string) ([]string, error) {
	files := make([]string, 0)
	if !utils.FileExists(sdCardFolder) {
		return files, ErrWrongCardFolder
	}
	err := fileutils.GetFiles(sdCardFolder, "data", func(fileinfo fs.DirEntry) bool {
		files = append(files, filepath.Join(sdCardFolder, fileinfo.Name()))
		return true
	})
	return files, err
}

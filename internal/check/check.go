package check

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/samber/do"
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
	files []string
	log   logging.Logger
}

func init() {
	chk := Checker{
		log: *logging.New().WithName("Checker"),
	}
	do.ProvideValue(nil, chk)
}

func (c *Checker) Check(sdCardFolder, outputFolder string) error {
	c.log.Infof("check called: sd %s, out: %s", sdCardFolder, outputFolder)

	files, err := c.getDataFiles(sdCardFolder)
	if err != nil {
		return err
	}
	c.files = make([]string, 0)
	c.files = append(c.files, files...)

	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		return err
	}

	for _, lf := range files {
		c.log.Infof("start with file %s", lf)
		err := c.analyseLoggerFile(lf, outputFolder)
		if err != nil {
			return err
		}
	}

	c.log.Infof("Found %d files on sd card", len(files))
	return nil
}

func (c *Checker) analyseLoggerFile(lf, of string) error {
	f, err := os.Open(lf)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	count := 0
	ls := make([]model.LogLine, 0)
	// Loop through the file and read each line
	for scanner.Scan() {
		count++
		line := scanner.Text() // Get the line as a string
		ll, ok, err := model.ParseLogLine(line)
		if err != nil {
			if ok {
				c.log.Alertf("warning unknown NMEA Tag in line %d: %s", count, line)
			} else {
				c.log.Errorf("error in line %d: %s", count, line)
			}
		}
		if ok {
			ls = append(ls, *ll)
		}
	}

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		c.log.Fatalf("error reading file: %s", err)
		return err
	}

	defer f.Close()
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

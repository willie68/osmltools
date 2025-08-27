package check

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/samber/do"
	utils "github.com/willie68/gowillie68/pkg"
	"github.com/willie68/gowillie68/pkg/fileutils"
	"github.com/willie68/osmltools/internal/logging"
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

func (c *Checker) Check(sdCardFolder, outputFile string) error {
	c.log.Infof("check called: sd %s, out: %s", sdCardFolder, outputFile)

	files, err := c.getDataFiles(sdCardFolder)
	if err != nil {
		return nil
	}

	if utils.FileExists(outputFile) {
		return ErrOutputfileAlreadyExists
	}
	c.log.Infof("Found %d files on sd card", len(files))
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

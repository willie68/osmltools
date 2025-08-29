package osml

import (
	"errors"
	"io/fs"
	"path/filepath"

	utils "github.com/willie68/gowillie68/pkg"
	"github.com/willie68/gowillie68/pkg/fileutils"
)

var (
	ErrWrongCardFolder = errors.New("sd card folder not exists")
)

func GetDataFiles(sdCardFolder string) ([]string, error) {
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

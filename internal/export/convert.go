package export

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/willie68/osmltools/internal/model"
)

func (e *Exporter) Convert(sdCardFolder string, files []string) (*model.Track, error) {
	fs, err := os.Stat(sdCardFolder)
	if err != nil {
		return nil, err
	}
	if fs.IsDir() && (len(files) == 0) {
		return nil, errors.New("sd card file is not a file")
	}
	if !fs.IsDir() {
		files = append(files, fs.Name())
		sdCardFolder = filepath.Dir(sdCardFolder)
	}
	var ls []*model.LogLine
	for _, file := range files {
		file = strings.TrimSpace(file)
		fp := filepath.Join(sdCardFolder, file)
		if _, err := os.Stat(fp); err != nil {
			return nil, err
		}
		ll, err := e.chk.AnalyseLoggerFile(nil, fp)
		if err != nil {
			return nil, err
		}
		ls = append(ls, ll...)
	}

	tr := &model.Track{
		Name:     sdCardFolder,
		LogLines: ls,
	}

	if len(ls) > 0 {
		ls, _, err = e.chk.CorrectTimeStamp(ls)
		if err != nil {
			return nil, err
		}

		sort.Slice(ls, func(i, j int) bool {
			return ls[i].CorrectTimeStamp.Before(ls[j].CorrectTimeStamp)
		})
		tr, err = e.GetWaypoints(tr)
		if err != nil {
			return nil, err
		}
	}

	return tr, nil
}

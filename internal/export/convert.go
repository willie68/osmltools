package export

import (
	"errors"
	"os"
	"sort"

	"github.com/willie68/osmltools/internal/model"
)

func (e *Exporter) Convert(sdCardFile string) (*model.Track, error) {
	fs, err := os.Stat(sdCardFile)
	if err != nil {
		return nil, err
	}
	if fs.IsDir() {
		return nil, errors.New("sd card file is not a file")
	}

	ls, err := e.chk.AnalyseLoggerFile(nil, sdCardFile)
	if err != nil {
		return nil, err
	}

	tr := &model.Track{
		Name:     sdCardFile,
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

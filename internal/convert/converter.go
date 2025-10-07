package convert

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samber/do/v2"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/trackutils"
)

type checkerSrv interface {
	AnalyseLoggerFile(fr *model.FileResult, lf string) ([]*model.LogLine, error)
	CorrectTimeStamp(ls []*model.LogLine) ([]*model.LogLine, bool, error)
}

type converter struct {
	log logging.Logger
	chk checkerSrv
}

func Init(inj do.Injector) {
	exp := converter{
		log: *logging.New().WithName("Converter"),
		chk: do.MustInvokeAs[checkerSrv](inj),
	}
	do.ProvideValue(inj, &exp)
}

func (c *converter) Convert(sdCardFolder string, files []string, track string) (tps *model.TrackPoints, err error) {
	if track != "" {
		tps, err = c.TrackPoints(track)
		if err != nil {
			return nil, err
		}
	} else {
		tps, err = c.convertData(sdCardFolder, files)
		if err != nil {
			return nil, err
		}
	}
	return tps, nil
}

func (c *converter) TrackPoints(trackfile string) (*model.TrackPoints, error) {
	if model.IsOldTrackVersion(trackfile) {
		return c.OldTrackPoints(trackfile)
	}
	return c.NewTrackPoints(trackfile)
}

func (c *converter) NewTrackPoints(trackfile string) (*model.TrackPoints, error) {
	track, nmealines, err := trackutils.ReadTrackAndNmea(trackfile)

	lls, err := model.ParseLines2LogLines(nmealines, false)
	if err != nil {
		return nil, err
	}

	// merge nmea with logline list
	tps := &model.TrackPoints{
		Name:     track.Name,
		LogLines: lls,
	}

	tps, err = model.GetWaypoints(tps)
	if err != nil {
		return nil, err
	}
	return tps, nil
}

func (c *converter) OldTrackPoints(trackfile string) (*model.TrackPoints, error) {
	track, nmealines, err := trackutils.ReadOldTrackAndNmea(trackfile)

	if err != nil {
		return nil, err
	}
	lls := make([]*model.LogLine, 0, len(nmealines))

	for _, l := range nmealines {
		ll, ok, err := model.ParseNMEALogLine(l, true)
		if err != nil {
			c.log.Errorf("error parsing nmea line: %v", err)
		}
		if ok {
			lls = append(lls, ll)
		}
	}

	// merge nmea with logline list
	tps := &model.TrackPoints{
		Name:     track.Name,
		LogLines: lls,
	}

	tps, err = model.GetWaypoints(tps)
	if err != nil {
		return nil, err
	}
	return tps, nil
}

func (c *converter) convertData(sdCardFolder string, files []string) (*model.TrackPoints, error) {
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
		ll, err := c.chk.AnalyseLoggerFile(nil, fp)
		if err != nil {
			return nil, err
		}
		ls = append(ls, ll...)
	}

	tr := &model.TrackPoints{
		Name:     sdCardFolder,
		LogLines: ls,
	}

	if len(ls) > 0 {
		ls, _, err = c.chk.CorrectTimeStamp(ls)
		if err != nil {
			return nil, err
		}

		sort.Slice(ls, func(i, j int) bool {
			return ls[i].CorrectTimeStamp.Before(ls[j].CorrectTimeStamp)
		})
		tr, err = model.GetWaypoints(tr)
		if err != nil {
			return nil, err
		}
	}

	return tr, nil
}

package track

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"github.com/willie68/osmltools/internal/export/nmeaexporter"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/trackutils"
)

type checkerSrv interface {
	AnalyseLoggerFile(fr *model.FileResult, lf string) ([]*model.LogLine, error)
	CorrectTimeStamp(ls []*model.LogLine) ([]*model.LogLine, bool, error)
}

// Manager the track manager service
type manager struct {
	log *logging.Logger
	chk checkerSrv
}

// Init init this service and provide it to di
func Init(inj do.Injector) {
	do.Provide(inj, func(inj do.Injector) (*manager, error) {
		return &manager{
			log: logging.New().WithName("Trackmanager"),
			chk: do.MustInvokeAs[checkerSrv](inj),
		}, nil
	})
}

func (m *manager) NewTrack(sdCardFolder string, files []string, trackfile string, track model.Track) error {
	m.log.Infof("Creating new track file %s", trackfile)
	track.MapFile = trackutils.NMEAFile

	ll, err := m.ReadLogFiles(files, sdCardFolder)
	if err != nil {
		return err
	}
	m.log.Infof("found %d loglines for the track.", len(ll))
	tps := &model.TrackPoints{
		Name:     track.Name,
		LogLines: ll,
	}
	// Create the ZIP file
	if err := os.MkdirAll(filepath.Dir(trackfile), os.ModePerm); err != nil {
		return err
	}
	outFile, err := os.Create(trackfile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	jsf, err := zipWriter.Create(trackutils.NMEAFile)
	if err != nil {
		m.log.Errorf("Failed to add track.nmea: %v", err)
	}
	err = nmeaexporter.New().ExportTrack(*tps, jsf)
	if err != nil {
		m.log.Errorf("Failed to export nmea: %v", err)
	}

	track, err = m.copyFiles2Zip(sdCardFolder, files, zipWriter, track)
	if err != nil {
		m.log.Errorf("Failed to add files: %v", err)
		return err
	}

	err = m.createTrackJSON(zipWriter, track)
	if err != nil {
		m.log.Errorf("Failed to create JSON: %v", err)
		return err
	}

	return nil
}

func (m *manager) createTrackJSON(zipWriter *zip.Writer, track model.Track) error {
	jsf, err := zipWriter.Create(trackutils.JSONFile)
	if err != nil {
		m.log.Errorf("Failed to add track.json: %v", err)
		return err
	}
	err = json.NewEncoder(jsf).Encode(track)
	if err != nil {
		m.log.Errorf("error marshal track: %v", err)
		return err
	}
	return nil
}

func (m *manager) copyFiles2Zip(sdCardFolder string, files []string, zipWriter *zip.Writer, track model.Track) (model.Track, error) {
	for _, file := range files {
		sdf := filepath.Join(sdCardFolder, strings.TrimSpace(file))
		sd, err := addFileToZip(zipWriter, sdf)
		if err != nil {
			m.log.Errorf("Failed to add %s: %v", file, err)
			return track, err
		}
		track.Files = append(track.Files, *sd)
	}
	return track, nil
}

// addFileToZip adds a file to the given zip.Writer
func addFileToZip(zipWriter *zip.Writer, filename string) (*model.SourceData, error) {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fileToZip.Close()

	// Get file info to preserve file name
	info, err := fileToZip.Stat()
	if err != nil {
		return nil, err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return nil, err
	}
	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	wr := io.MultiWriter(writer, h)

	n, err := io.Copy(wr, fileToZip)
	sd := model.SourceData{
		FileName: filepath.Base(filename),
		Size:     n,
		Hash:     fmt.Sprintf("sha256:%s", hex.EncodeToString(h.Sum(nil))),
		Modified: info.ModTime(),
	}

	return &sd, err
}

func (m *manager) ReadLogFiles(files []string, sdCardFolder string) ([]*model.LogLine, error) {
	ls := make([]*model.LogLine, 0)
	today := time.Time{}

	for _, file := range files {
		sdf := filepath.Join(sdCardFolder, strings.TrimSpace(file))
		m.log.Infof("analysing file: %s", sdf)
		lss, err := m.chk.AnalyseLoggerFile(nil, sdf)
		if err != nil {
			return nil, err
		}
		lss, _, err = m.chk.CorrectTimeStamp(lss)
		if err != nil {
			return nil, err
		}
		if len(lss) > 0 {
			if today.IsZero() {
				today = lss[0].CorrectTimeStamp
			}
			ls = append(ls, lss...)
		}
	}

	sort.Slice(ls, func(i, j int) bool {
		return ls[i].CorrectTimeStamp.Before(ls[j].CorrectTimeStamp)
	})
	return ls, nil
}

package track

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"github.com/willie68/osmltools/internal/check"
	"github.com/willie68/osmltools/internal/export/nmeaexporter"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

const (
	mapFile  = "track.nmea"
	jsonFile = "track.json"
)

// Manager the track manager interface
type Manager interface {
	NewTrack(sdCardFolder string, files []string, trackfile string, track model.Track) error
	AddTrack(sdCardFolder string, files []string, trackfile string) error
	DeleteTrack(sdCardFolder string, trackfile string) error
	ListTrack(sdCardFolder string, trackfile string) (*model.Track, error)
}

// Manager the track manager service
type manager struct {
	log *logging.Logger
	chk check.Checker
}

var _ Manager = &manager{}

// Init init this service and provide it to di
func Init(inj do.Injector) {
	trm := manager{
		log: logging.New().WithName("Trackmanager"),
		chk: do.MustInvoke[check.Checker](inj),
	}
	do.ProvideValue(inj, &trm)
}

func (m *manager) NewTrack(sdCardFolder string, files []string, trackfile string, track model.Track) error {
	m.log.Infof("Creating new track file %s", trackfile)
	track.MapFile = mapFile

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
	outFile, err := os.Create(trackfile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	jsf, err := zipWriter.Create(mapFile)
	if err != nil {
		m.log.Errorf("Failed to add track.nmea: %v", err)
	}
	err = nmeaexporter.New().ExportTrack(*tps, jsf)
	if err != nil {
		m.log.Errorf("Failed to export nmea: %v", err)
	}

	for _, file := range files {
		sdf := filepath.Join(sdCardFolder, strings.TrimSpace(file))
		sd, err := addFileToZip(zipWriter, sdf)
		if err != nil {
			m.log.Errorf("Failed to add %s: %v", file, err)
		}
		track.Files = append(track.Files, *sd)
	}
	jsf, err = zipWriter.Create(jsonFile)
	if err != nil {
		m.log.Errorf("Failed to add track.json: %v", err)
	}
	err = json.NewEncoder(jsf).Encode(track)
	if err != nil {
		m.log.Errorf("error marshal track: %v", err)
	}

	return nil
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
		FileName: filename,
		Size:     n,
		Hash:     fmt.Sprintf("sha256:%s", hex.EncodeToString(h.Sum(nil))),
	}

	return &sd, err
}

func (m *manager) AddTrack(sdCardFolder string, files []string, trackfile string) error {
	m.log.Infof("Adding data to track file %s", trackfile)
	if m.IsOldVersion(trackfile) {
		return errors.New("can't add data to an old track file.")
	}
	return nil
}

func (m *manager) DeleteTrack(sdCardFolder string, trackfile string) error {
	m.log.Infof("Deleting track file %s", trackfile)
	return nil
}

func SplitMultiValueParam(value string) []string {
	return strings.FieldsFunc(value, func(r rune) bool {
		return r == ' ' || r == ',' || r == ';'
	})
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

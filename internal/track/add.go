package track

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/willie68/osmltools/internal/export/nmeaexporter"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/trackutils"
)

// AddTrack adds data from the given files to the given track file
func (m *manager) AddTrack(sdCardFolder string, files []string, trackfile string) error {
	m.log.Infof("Adding data to track file %s", trackfile)
	if model.IsOldTrackVersion(trackfile) {
		return errors.New("can't add data to an old track file")
	}

	// read sd files and build logline list
	ll, err := m.ReadLogFiles(files, sdCardFolder)
	if err != nil {
		return err
	}

	// open zip, load track and nmea file
	track, nmealines, err := trackutils.ReadTrackAndNmea(trackfile)
	if err != nil {
		return err
	}
	lls := make([]*model.LogLine, 0, len(nmealines))
	for _, l := range nmealines {
		ll, ok, err := model.ParseNMEALogLine(l, false)
		if err != nil {
			return err
		}
		if ok {
			lls = append(lls, ll)
		}
	}

	// merge nmea with logline list
	tps := model.TrackPoints{
		Name:     track.Name,
		LogLines: make([]*model.LogLine, 0, len(ll)+len(lls)),
	}
	tps.LogLines = append(tps.LogLines, lls...)
	tps.LogLines = append(tps.LogLines, ll...)

	fmt.Printf("lines:%d, track: %v", len(tps.LogLines), track)
	// update with new nmea file add source files to zip
	return m.openNewZipCopyContent(sdCardFolder, files, trackfile, tps, *track)
}

func (m *manager) openNewZipCopyContent(sdCardFolder string, files []string, trackfile string, tps model.TrackPoints, track model.Track) error {
	// Create a temporary file
	tmpFile, err := os.CreateTemp(filepath.Dir(trackfile), "updated-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name()) // Clean up if needed

	// Create ZIP writer to temp file
	zipWriter := zip.NewWriter(tmpFile)

	// Open original ZIP for reading
	err = m.copyOldFiles(trackfile, zipWriter)
	if err != nil {
		return err
	}

	// create new NMEA File
	jsf, err := zipWriter.Create(trackutils.NMEAFile)
	if err != nil {
		m.log.Errorf("Failed to add track.nmea: %v", err)
	}
	err = nmeaexporter.New().ExportTrack(tps, jsf)
	if err != nil {
		m.log.Errorf("Failed to export nmea: %v", err)
	}

	// copy new data files
	track, err = m.copyFiles2Zip(sdCardFolder, files, zipWriter, track)
	if err != nil {
		m.log.Errorf("Failed to add files: %v", err)
	}

	// create Track JSON

	err = m.createTrackJSON(zipWriter, track)
	if err != nil {
		m.log.Errorf("Failed to create JSON: %v", err)
		return err
	}

	// Finalize ZIP
	if err := zipWriter.Close(); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	// Replace original ZIP with temp file
	return os.Rename(tmpFile.Name(), trackfile)
}

func (m *manager) copyOldFiles(trackfile string, zipWriter *zip.Writer) error {
	r, err := zip.OpenReader(trackfile)
	if err != nil {
		return err
	}
	// copy old files
	for _, f := range r.File {
		var fw io.Writer
		// ignore nmea and track file
		if f.Name == trackutils.JSONFile || f.Name == trackutils.NMEAFile {
			continue
		}

		// Copy original file
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fw, err = zipWriter.CreateHeader(&f.FileHeader)
		if err != nil {
			return err
		}
		if _, err := io.Copy(fw, rc); err != nil {
			return err
		}
	}

	r.Close()
	return nil
}

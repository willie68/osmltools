package trackutils

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/willie68/osmltools/internal/model"
)

const (
	NMEAFile = "track.nmea"
	JSONFile = "track.json"
)

func ReadTrackAndNmea(tf string) (*model.Track, []string, error) {
	if _, err := os.Stat(tf); err != nil {
		return nil, nil, fmt.Errorf("error zip file %s does not exists: %v", tf, err)
	}

	r, err := zip.OpenReader(tf)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening zip file %s: %v", tf, err)
	}
	defer r.Close()
	var lines []string
	track := &model.Track{}
	for _, f := range r.File {
		if f.Name == JSONFile {
			track, err = Track(f)
			if err != nil {
				return nil, nil, err
			}
		}
		if f.Name == NMEAFile {
			lines, err = NMEA(f)
			if err != nil {
				return nil, nil, err
			}
		}

	}
	return track, lines, nil
}

func Track(f *zip.File) (track *model.Track, err error) {
	rc, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening json file %s: %v", f.Name, err)
	}
	defer rc.Close()

	err = json.NewDecoder(rc).Decode(&track)
	if err != nil {
		return nil, err
	}
	return
}

func NMEA(f *zip.File) (lines []string, err error) {
	rc, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening json file %s: %v", f.Name, err)
	}
	defer rc.Close()
	scanner := bufio.NewScanner(rc)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return
}

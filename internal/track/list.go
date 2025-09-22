package track

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"

	"github.com/willie68/osmltools/internal/model"
)

func (m *manager) ListTrack(sdCardFolder string, trackfile string) (*model.Track, error) {
	m.log.Infof("Listing track file %s", trackfile)
	if m.IsOldVersion(trackfile) {
		return m.ListOldTrack(trackfile)
	}
	return m.ListNewTrack(trackfile)
}

func (m *manager) ListNewTrack(tf string) (*model.Track, error) {
	if _, err := os.Stat(tf); err != nil {
		return nil, fmt.Errorf("error zip file %s does not exists: %v", tf, err)
	}

	r, err := zip.OpenReader(tf)
	if err != nil {
		return nil, fmt.Errorf("error opening zip file %s: %v", tf, err)
	}
	defer r.Close()
	track := model.Track{}
	for _, f := range r.File {
		if f.Name == jsonFile {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("error opening json file %s: %v", f.Name, err)
			}
			defer rc.Close()

			err = json.NewDecoder(rc).Decode(&track)
			if err != nil {
				return nil, err
			}
		}

	}
	return &track, nil
}

func (m *manager) ListOldTrack(tf string) (*model.Track, error) {
	if _, err := os.Stat(tf); err != nil {
		return nil, fmt.Errorf("error zip file %s does not exists: %v", tf, err)
	}

	r, err := zip.OpenReader(tf)
	if err != nil {
		return nil, fmt.Errorf("error opening zip file %s: %v", tf, err)
	}
	defer r.Close()
	track := model.Track{}
	for _, f := range r.File {
		if f.Name == "route.properties" {
			err := m.routeProps2track(&track, f)
			if err != nil {
				return nil, err
			}
		}
	}
	return &track, nil
}

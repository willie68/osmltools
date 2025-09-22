package track

import (
	"archive/zip"
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

func (m *manager) ListNewTrack(trackfile string) (*model.Track, error) {
	tr := model.Track{
		Name:        "Testtrack",
		Description: "This is a test track",
		VesselID:    1234,
		Files: []model.SourceData{
			{
				FileName: "data1.log",
				Size:     123456,
				Hash:     "abcdef1234567890",
			},
			{
				FileName: "data2.log",
				Size:     654321,
				Hash:     "098765fedcba",
			},
		},
		MapFile: mapFile,
	}
	return &tr, nil
	//return model.Track{}, nil
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

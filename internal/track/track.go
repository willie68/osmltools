package track

import (
	"archive/zip"
	"os"

	"github.com/samber/do/v2"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

// Manager the track manager interface
type Manager interface {
	NewTrack(sdCardFolder string, files []string, trackfile string, track model.Track) error
	AddTrack(sdCardFolder string, files []string, trackfile string) error
	DeleteTrack(sdCardFolder string, trackfile string) error
	ListTrack(sdCardFolder string, trackfile string) (model.Track, error)
}

// Manager the track manager service
type manager struct {
	log *logging.Logger
}

// Init init this service and provide it to di
func Init(inj do.Injector) {
	trm := manager{
		log: logging.New().WithName("Trackmanager"),
	}
	do.ProvideValue(inj, trm)
}

func (m *manager) NewTrack(sdCardFolder string, files []string, trackfile string, track model.Track) error {
	m.log.Infof("Creating new track file %s", trackfile)
	return nil
}

func (m *manager) AddTrack(sdCardFolder string, files []string, trackfile string) error {
	m.log.Infof("Adding data to track file %s", trackfile)
	return nil
}

func (m *manager) DeleteTrack(sdCardFolder string, trackfile string) error {
	m.log.Infof("Deleting track file %s", trackfile)
	return nil
}

func (m *manager) ListTrack(sdCardFolder string, trackfile string) (model.Track, error) {
	m.log.Infof("Listing track file %s", trackfile)
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
		MapFile: "track.nmea",
	}
	return tr, nil
	//return model.Track{}, nil
}

// IsOldVersion checks if the given track file is in the old format (contains route.properties)
func (m *manager) IsOldVersion(tf string) bool {
	if _, err := os.Stat(tf); err != nil {
		m.log.Errorf("error zip file %s does not exists: %v", tf, err)
		return false
	}

	r, err := zip.OpenReader(tf)
	if err != nil {
		m.log.Errorf("error opening zip file %s: %v", tf, err)
		return false
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "route.properties" {
			return true
		}
	}
	return false
}

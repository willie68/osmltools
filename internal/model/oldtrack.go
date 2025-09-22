package model

import (
	"archive/zip"
	"os"

	"github.com/willie68/osmltools/internal/logging"
)

var (
	log *logging.Logger = logging.New().WithName("model")
)

// IsOldTrackVersion checks if the given track file is in the old format (contains route.properties)
func IsOldTrackVersion(tf string) bool {
	if _, err := os.Stat(tf); err != nil {
		log.Errorf("error zip file %s does not exists: %v", tf, err)
		return false
	}

	r, err := zip.OpenReader(tf)
	if err != nil {
		log.Errorf("error opening zip file %s: %v", tf, err)
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

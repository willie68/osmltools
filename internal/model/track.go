package model

import (
	"encoding/json"
)

// Track track structure, containing metadata and the list of source data files and the map file for the ui
type Track struct {
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	VesselID    int32        `json:"vessel_id,omitempty"`
	Files       []SourceData `json:"files,omitempty"`
	MapFile     string       `json:"map_file,omitempty"`
}

// SourceData information about a source data file
type SourceData struct {
	FileName string `json:"file_name,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Hash     string `json:"hash,omitempty"`
}

// JSON return the json representation of the version
func (t *Track) JSON() (string, error) {
	js, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

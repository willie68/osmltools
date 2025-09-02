package interfaces

import "github.com/willie68/osmltools/internal/model"

type FormatExporter interface {
	ExportTrack(track model.Track, outputfile string) error
}

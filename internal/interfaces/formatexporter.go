package interfaces

import (
	"io"

	"github.com/willie68/osmltools/internal/model"
)

type FormatExporter interface {
	ExportTrack(track model.TrackPoints, output io.Writer) error
}

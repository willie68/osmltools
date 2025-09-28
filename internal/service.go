package internal

import (
	"github.com/samber/do/v2"
	"github.com/willie68/osmltools/internal/backup"
	"github.com/willie68/osmltools/internal/check"
	"github.com/willie68/osmltools/internal/config"
	"github.com/willie68/osmltools/internal/convert"
	"github.com/willie68/osmltools/internal/export"
	"github.com/willie68/osmltools/internal/track"
)

var (
	// Inj the injector used for the programm
	Inj = do.New()
)

// Init initialise all needed services
func Init() {
	check.Init(Inj)
	config.Init(Inj)
	export.Init(Inj)
	backup.Init(Inj)
	track.Init(Inj)
	convert.Init(Inj)
}

package internal

import (
	"github.com/willie68/osmltools/internal/check"
	"github.com/willie68/osmltools/internal/config"
	"github.com/willie68/osmltools/internal/export"
)

func Init() {
	check.Init(nil)
	config.Init(nil)
	export.Init(nil)
}

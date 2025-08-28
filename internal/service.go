package internal

import (
	"github.com/willie68/osmltools/internal/check"
	"github.com/willie68/osmltools/internal/config"
)

func Init() {
	check.Init(nil)
	config.Init(nil)
}

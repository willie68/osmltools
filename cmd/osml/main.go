package main

import (
	"os"

	"github.com/willie68/osmltools/cmd/osml/cmd"
	"github.com/willie68/osmltools/internal/config"
	"github.com/willie68/osmltools/internal/logging"

	// blank imports for every processor, so that the processor is registered to the di framework
	_ "github.com/willie68/osmltools/internal/check"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	ver := config.NewVersion().WithCommit(commit).WithDate(date).WithVersion(version)
	cmd.CmdVersion = *ver
	err := cmd.Execute()
	if err != nil {
		logging.Root.Errorf("error on command: %v", err)
		os.Exit(1)
	}
}

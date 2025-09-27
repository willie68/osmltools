package main

import (
	"os"

	"github.com/willie68/osmltools/cmd/osml/cmd"
	"github.com/willie68/osmltools/internal/logging"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		if cmd.JSONOutput {
			cmd.OutputErrorJSON(err)
		} else {
			logging.Root.Errorf("error on command: %v", err)
		}
		os.Exit(1)
	}
}

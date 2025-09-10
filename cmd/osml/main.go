package main

import (
	"fmt"
	"os"

	"github.com/willie68/osmltools/cmd/osml/cmd"
	"github.com/willie68/osmltools/internal/logging"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		if cmd.JSONOutput {
			fmt.Printf("{\"error\": \"%v\"}", err)
		} else {
			logging.Root.Errorf("error on command: %v", err)
		}
		os.Exit(1)
	}
}

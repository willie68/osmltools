//go:build windows

package sdformatter

import (
	"fmt"
	"os/exec"
	"strings"
)

func setLabel(devicePath, label string) error {
	// devicePath wie "E:"
	devicePath = strings.TrimSuffix(devicePath, "\\")
	cmd := exec.Command("cmd", "/C", "label", devicePath, label)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("label Fehler: %v, Ausgabe: %s", err, string(out))
	}
	return nil
}

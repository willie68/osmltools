//go:build linux

package sdformatter

import (
	"fmt"
	"os/exec"
)

func setLabel(devicePath, label string) error {
	// Versuche fatlabel
	cmd := exec.Command("fatlabel", devicePath, label)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("fatlabel Fehler: %v, Ausgabe: %s", err, string(out))
	}
	return nil
}
//go:build darwin

package sdformatter

import (
	"fmt"
	"os/exec"
)

func setLabel(devicePath, label string) error {
	// devicePath wie "/Volumes/SDCARD"
	cmd := exec.Command("diskutil", "rename", devicePath, label)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("diskutil Fehler: %v, Ausgabe: %s", err, string(out))
	}
	return nil
}
//go:build windows

package sdformatter

import (
	"fmt"
	"os/exec"
)

func formatFAT32(devicePath string) error {
	// Format using the format command
	// devicePath should be like "\\\\.\\E:"
	cmd := exec.Command("cmd", "/C", "format", devicePath, "/FS:FAT32", "/Q", "/Y")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to format SD card: %v, output: %s", err, string(out))
	}
	return nil
}
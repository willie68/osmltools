//go:build darwin

package sdformatter

import (
	"fmt"
	"os/exec"
)

func formatFAT32(devicePath string) error {
	// Format the disk using diskutil
	// devicePath should be like "/dev/disk2"
	cmd := exec.Command("diskutil", "eraseDisk", "FAT32", "SDCARD", devicePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to format SD card: %v, output: %s", err, string(out))
	}
	return nil
}
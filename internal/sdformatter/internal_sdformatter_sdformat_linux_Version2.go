//go:build linux

package sdformatter

import (
	"fmt"
	"os/exec"
)

func formatFAT32(devicePath string) error {
	// Unmount the device if mounted
	_ = exec.Command("umount", devicePath).Run()

	// Format as FAT32 using mkfs.vfat
	cmd := exec.Command("mkfs.vfat", "-F", "32", devicePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to format SD card: %v, output: %s", err, string(out))
	}
	return nil
}
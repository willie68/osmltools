package sdformatter

// FormatFAT32 formats the SD card at the given device path with FAT32.
// Device path examples:
//   Linux:   "/dev/sdb"
//   macOS:   "/dev/disk2"
//   Windows: "\\\\.\\E:"
func FormatFAT32(devicePath string) error {
	return formatFAT32(devicePath)
}
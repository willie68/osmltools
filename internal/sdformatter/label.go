package sdformatter

// SetLabel setzt das Volume Label des angegebenen Laufwerks plattformunabhängig.
// devicePath Beispiele:
//  - Linux: "/dev/sdb1"
//  - macOS: "/Volumes/SDCARD"
//  - Windows: "E:"
func SetLabel(devicePath, label string) error {
	return setLabel(devicePath, label)
}
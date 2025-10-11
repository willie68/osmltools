//go:build windows
// +build windows

package check

import (
	"time"

	"golang.org/x/sys/windows"
)

// modifyFileTime modifies the file time of the given file to the newTime for windows only
func (c *checker) modifyFileTime(path string, newTime time.Time) error {
	// This example is for Windows. For other OS, the implementation will differ.
	// You may need to use syscall or a third-party package for cross-platform support.
	// Here, we use the golang.org/x/sys/windows package for Windows.
	// Make sure to import "golang.org/x/sys/windows"
	// and run `go get golang.org/x/sys/windows` to install it.
	// Open file with FILE_WRITE_ATTRIBUTES
	pathUTF16, err := windows.UTF16PtrFromString(path)
	if err != nil {
		c.log.Fatalf("error converting path: %v", err)
		return err
	}

	handle, err := windows.CreateFile(
		pathUTF16,
		windows.FILE_WRITE_ATTRIBUTES,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		c.log.Fatalf("CreateFile failed: %v", err)
		return err
	}
	defer windows.CloseHandle(handle)

	// Set new time
	ft := windows.NsecToFiletime(newTime.UnixNano())

	err = windows.SetFileTime(handle, nil, &ft, &ft) // only setting LastWriteTime
	if err != nil {
		c.log.Fatalf("SetFileTime failed: %v", err)
		return err
	}
	return nil
}

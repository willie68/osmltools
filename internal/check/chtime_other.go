//go:build !windows
// +build !windows

package check

import (
	"errors"
	"time"
)

func (c *checker) modifyFileTime(_ string, _ time.Time) error {
	return errors.New("not implemented on this platform")
}

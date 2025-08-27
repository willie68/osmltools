package check

import (
	"github.com/samber/do"
	"github.com/willie68/osmltools/internal/logging"
)

type Checker struct {
	log logging.Logger
}

func init() {
	chk := Checker{
		log: *logging.New().WithName("Checker"),
	}
	do.ProvideValue(nil, chk)
}

func (c *Checker) Check(sdCardFolder, outputFile string) error {
	c.log.Infof("check called: sd %s, out: %s", sdCardFolder, outputFile)
	return nil
}

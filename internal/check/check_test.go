package check

import (
	"path/filepath"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	utils "github.com/willie68/gowillie68/pkg"
)

const testdata = "../../testdata"

func TestCheckSDInit(t *testing.T) {
	ast := assert.New(t)

	chk, err := do.Invoke[Checker](nil)
	ast.NoError(err)
	ast.NotNil(chk)
}

func TestCheckBasicCheck(t *testing.T) {
	ast := assert.New(t)

	chk, err := do.Invoke[Checker](nil)
	ast.NoError(err)
	ast.NotNil(chk)

	err = chk.Check(filepath.Join(testdata, "sdcard"), filepath.Join(testdata, "temp/track.nmea"))
	ast.NoError(err)
	ast.True(utils.FileExists(filepath.Join(testdata, "temp/track.nmea")))
}

func TestCheckWrongSDCardFolder(t *testing.T) {
	ast := assert.New(t)

	chk, err := do.Invoke[Checker](nil)
	ast.NoError(err)
	ast.NotNil(chk)

	err = chk.Check("./testdata/sdcard1", "./testdata/temp/track.nmea")
	ast.ErrorIs(ErrWrongCardFolder, err)
}

func TestCheckNMEAFIleAlreadyExists(t *testing.T) {
	ast := assert.New(t)

	chk, err := do.Invoke[Checker](nil)
	ast.NoError(err)
	ast.NotNil(chk)

	err = chk.Check(filepath.Join(testdata, "sdcard"), filepath.Join(testdata, "track.nmea"))
	ast.ErrorIs(ErrOutputfileAlreadyExists, err)
}

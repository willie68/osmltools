package check

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	utils "github.com/willie68/gowillie68/pkg"
	"github.com/willie68/osmltools/internal/logging"
)

const testdata = "../../testdata"

type CheckSuite struct {
	suite.Suite
	ast *assert.Assertions
	chk Checker
}

func TestCheckSuite(t *testing.T) {
	suite.Run(t, new(CheckSuite))
}

func (s *CheckSuite) SetupTest() {
	s.ast = assert.New(s.T())
	s.chk = Checker{
		log: *logging.New().WithName("testchecker").WithLevel(logging.Debug),
	}
}

func (s *CheckSuite) TestCheckBasicCheck() {
	of := filepath.Join(testdata, "temp")
	os.RemoveAll(of)
	os.MkdirAll(of, os.ModePerm)

	err := s.chk.Check(
		filepath.Join(testdata, "sdcard"),
		of,
		true,
	)
	s.ast.NoError(err)
	s.ast.Equal(7, s.chk.ErrorTags)
	s.ast.Equal(24421, s.chk.UnknownTags)
	s.ast.True(utils.FileExists(filepath.Join(of, "DATA001231.nmea")))
	s.ast.True(utils.FileExists(filepath.Join(of, "DATA001232.nmea")))
	s.ast.True(utils.FileExists(filepath.Join(of, "DATA001233.nmea")))
	s.ast.True(utils.FileExists(filepath.Join(of, "DATA001234.nmea")))
	s.ast.True(utils.FileExists(filepath.Join(of, "DATA001235.nmea")))
}

func (s *CheckSuite) TestCheckWrongSDCardFolder() {
	err := s.chk.Check("./testdata/sdcard1", "./testdata/temp/track.nmea", false)
	s.ast.ErrorIs(ErrWrongCardFolder, err)
}

func (s *CheckSuite) TestCheckNMEAFileAlreadyExists() {
	err := s.chk.Check(
		filepath.Join(testdata, "sdcard"),
		filepath.Join(testdata, "already"),
		false,
	)
	s.ast.ErrorIs(err, ErrOutputfileAlreadyExists)
}

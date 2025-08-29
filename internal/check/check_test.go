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
		log: *logging.New().WithName("testchecker").WithLevel(logging.Error),
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
		true,
	)
	s.ast.NoError(err)
	s.ast.Equal(8147, s.chk.ErrorTags)
	s.ast.Equal(16281, s.chk.UnknownTags)
	s.ast.True(utils.FileExists(filepath.Join(of, "597-DATA001231-2016-09-11.nmea")))
	s.ast.True(utils.FileExists(filepath.Join(of, "597-DATA001232-2016-09-11.nmea")))
	s.ast.True(utils.FileExists(filepath.Join(of, "597-DATA001233-2016-09-11.nmea")))
	s.ast.True(utils.FileExists(filepath.Join(of, "597-DATA001234-2016-09-11.nmea")))
	s.ast.True(utils.FileExists(filepath.Join(of, "597-DATA001235-2016-09-11.nmea")))
	s.ast.True(utils.FileExists(filepath.Join(of, "report.json")))
}

func (s *CheckSuite) TestCheckWrongSDCardFolder() {
	err := s.chk.Check("./testdata/sdcard1", "./testdata/temp/track.nmea", false, false)
	s.ast.ErrorIs(ErrWrongCardFolder, err)
}

func (s *CheckSuite) TestCheckNMEAFileAlreadyExists() {
	err := s.chk.Check(
		filepath.Join(testdata, "sdcard"),
		filepath.Join(testdata, "already"),
		false,
		false,
	)
	s.ast.ErrorIs(err, ErrOutputfileAlreadyExists)
}

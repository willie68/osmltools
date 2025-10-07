package check

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/willie68/gowillie68/pkg/fileutils"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/osml"
)

const testdata = "../../testdata"

type checkerSrv interface {
	Check(sdCardFolder, outputFolder string, overwrite, report bool) (*model.CheckResult, error)
}

type CheckSuite struct {
	suite.Suite
	ast *assert.Assertions
	chk checkerSrv
}

func TestCheckSuite(t *testing.T) {
	suite.Run(t, new(CheckSuite))
}

func (s *CheckSuite) SetupTest() {
	s.ast = assert.New(s.T())
	inj := do.New()
	Init(inj)
	s.chk = do.MustInvokeAs[checkerSrv](inj)
}

func (s *CheckSuite) TestCheckBasicCheck() {
	of := filepath.Join(testdata, "temp")
	os.RemoveAll(of)
	os.MkdirAll(of, os.ModePerm)

	res, err := s.chk.Check(
		filepath.Join(testdata, "sdcard"),
		of,
		true,
		true,
	)
	s.ast.NoError(err)
	s.ast.Equal(8147, res.ErrorTags)
	s.ast.Equal(16281, res.UnknownTags)
	s.ast.True(fileutils.FileExists(filepath.Join(of, "597-DATA001231-2016-09-11.nmea")))
	s.ast.True(fileutils.FileExists(filepath.Join(of, "597-DATA001232-2016-09-11.nmea")))
	s.ast.True(fileutils.FileExists(filepath.Join(of, "597-DATA001233-2016-09-11.nmea")))
	s.ast.True(fileutils.FileExists(filepath.Join(of, "597-DATA001234-2016-09-11.nmea")))
	s.ast.True(fileutils.FileExists(filepath.Join(of, "597-DATA001235-2016-09-11.nmea")))
	s.ast.True(fileutils.FileExists(filepath.Join(of, "report.json")))
}

func (s *CheckSuite) TestCheckWrongSDCardFolder() {
	_, err := s.chk.Check("./testdata/sdcard1", "./testdata/temp/track.nmea", false, false)
	s.ast.ErrorIs(osml.ErrWrongCardFolder, err)
}

func (s *CheckSuite) TestCheckNMEAFileAlreadyExists() {
	_, err := s.chk.Check(
		filepath.Join(testdata, "sdcard"),
		filepath.Join(testdata, "already"),
		false,
		false,
	)
	s.ast.ErrorIs(err, ErrOutputfileAlreadyExists)
}

func (s *CheckSuite) TestCheckEmptyFile() {
	res, err := s.chk.Check(
		filepath.Join(testdata, "empty", "DATA001231.DAT"),
		filepath.Join(testdata, "tmp"),
		false,
		false,
	)
	s.ast.NoError(err)
	s.ast.Equal(1, res.ErrorCount)
}

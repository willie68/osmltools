package backup

import (
	"os"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/suite"
)

const (
	testZIP = "../../testdata/bck/bck_20250913160522.zip"
)

type backupSrv interface {
	Backup(sdCardFolder, outputFolder string) (string, error)
	Restore(zipfile, sdCardFolder string) (string, error)
}

type BackupTestSuite struct {
	suite.Suite
	bck backupSrv
}

func TestUploadTestSuite(t *testing.T) {
	suite.Run(t, new(BackupTestSuite))
}

func (s *BackupTestSuite) SetupTest() {
	inj := do.New()
	Init(inj)
	s.bck = do.MustInvokeAs[backupSrv](inj)
	s.NotNil(s.bck)
}

func (s *BackupTestSuite) TestBackup() {
	filename, err := s.bck.Backup("../../testdata/sdCard", "../../testdata/bck/")
	s.NoError(err)
	s.NotEmpty(filename)
}

func (s *BackupTestSuite) TestRestore() {
	err := os.MkdirAll("../../testdata/rst", os.ModePerm)
	s.NoError(err)

	filename, err := s.bck.Restore(testZIP, "../../testdata/rst")
	s.NoError(err)
	s.Equal(testZIP, filename)
}

func (s *BackupTestSuite) TestRestore1() {
	err := os.MkdirAll("../../testdata/rst", os.ModePerm)
	s.NoError(err)
	zip := "../../testdata/bck/bck_20250913160205.zip"

	filename, err := s.bck.Restore(zip, "../../testdata/rst1")
	s.NoError(err)
	s.Equal(zip, filename)
}

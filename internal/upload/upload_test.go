package upload

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/suite"
)

type Manager interface {
	StoreCredentials(user, password string) error
	GetCredentials(user string) (string, error)
}

type UploadTestSuite struct {
	suite.Suite
	inj do.Injector
	upl Manager
}

func TestUploadTestSuite(t *testing.T) {
	suite.Run(t, new(UploadTestSuite))
}

func (s *UploadTestSuite) SetupTest() {
	s.inj = do.New()
	Init(s.inj)
	s.upl = do.MustInvokeAs[Manager](s.inj)
}

func (s *UploadTestSuite) TestStoreCread() {
	user := "Willie"
	password := "meinSuperDuperöäü?Password"

	s.upl.StoreCredentials(user, password)

	pwd, err := s.upl.GetCredentials(user)

	s.NoError(err)
	s.Equal(password, pwd)
}

func (s *UploadTestSuite) TestGetWrongUser() {
	pwd, err := s.upl.GetCredentials("wronguser")

	s.Error(err)
	s.Empty(pwd)
}

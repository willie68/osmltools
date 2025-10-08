package upload

import (
	"github.com/samber/do/v2"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/zalando/go-keyring"
)

type manager struct {
	log *logging.Logger
}

const (
	osmlService = "osml-upload"
)

// Init init this service and provide it to di
func Init(inj do.Injector) {
	do.Provide(inj, func(inj do.Injector) (*manager, error) {
		return &manager{
			log: logging.New().WithName("Upload"),
		}, nil
	})
}

func (m *manager) StoreCredentials(user, password string) error {
	err := keyring.Set(osmlService, user, password)
	if err != nil {
		return err
	}
	return nil
}

func (m *manager) GetCredentials(user string) (string, error) {
	secret, err := keyring.Get(osmlService, user)
	if err != nil {
		return "", err
	}
	return secret, nil
}

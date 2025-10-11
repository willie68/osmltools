package upload

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

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

type FileUpload struct {
	FilePath string
	VesselID int
	Username string
	Hash     string
}

func (fu *FileUpload) Upload(url string) error {
	if err := fu.CalculateHash(); err != nil {
		return err
	}

	exists, err := fu.CheckDuplicate(url)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("⚠️ Datei wurde bereits hochgeladen. Upload übersprungen.")
		return nil
	}

	// ... (Multipart-Upload wie zuvor)
	file, err := os.Open(fu.FilePath)
	if err != nil {
		return fmt.Errorf("Fehler beim Öffnen der Datei: %w", err)
	}
	defer file.Close()

	// Multipart-Form vorbereiten
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Datei anhängen
	part, err := writer.CreateFormFile("datei", filepath.Base(fu.FilePath))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}

	// Weitere Felder
	writer.WriteField("fahrzeugid", fmt.Sprintf("%d", fu.VesselID))
	writer.WriteField("benutzer", fu.Username)
	writer.Close()

	// HTTP-Request senden
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Antwort anzeigen
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Antwort vom Server:", string(respBody))

	return nil
}

func (fu *FileUpload) CheckDuplicate(serverURL string) (bool, error) {
	checkURL := fmt.Sprintf("%s?checkhash=%s", serverURL, fu.Hash)
	resp, err := http.Get(checkURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body) == "exists", nil
}

func (fu *FileUpload) CalculateHash() error {
	file, err := os.Open(fu.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	fu.Hash = hex.EncodeToString(hash.Sum(nil))
	return nil
}

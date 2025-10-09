package upload

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

type CheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Simulierter Server für Testfälle
func mockCheckHashServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hash := r.URL.Query().Get("checkhash")

		var resp CheckResponse
		switch hash {
		case "validbutnewhash":
			resp = CheckResponse{"notfound", "Datei nicht vorhanden"}
		case "existinghash":
			resp = CheckResponse{"exists", "Datei bereits vorhanden"}
		case "invalidhash":
			w.WriteHeader(http.StatusBadRequest)
			resp = CheckResponse{"error", "Ungültiger Hash"}
		default:
			resp = CheckResponse{"error", "Unbekannter Testfall"}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

// Testfunktion
func TestCheckDuplicate(t *testing.T) {
	server := mockCheckHashServer(t)
	defer server.Close()

	tests := []struct {
		name     string
		hash     string
		expected string
	}{
		{"Datei noch nicht vorhanden", "validbutnewhash", "notfound"},
		{"Datei vorhanden", "existinghash", "exists"},
		{"Hash ungültig", "invalidhash", "error"},
		{"Unbekannter Hash", "unknownhash", "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := server.URL + "?checkhash=" + tt.hash
			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("Fehler bei Anfrage: %v", err)
			}
			defer resp.Body.Close()

			var result CheckResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("Fehler beim Dekodieren: %v", err)
			}

			if result.Status != tt.expected {
				t.Errorf("Erwartet %s, erhalten %s", tt.expected, result.Status)
			}
		})
	}
}

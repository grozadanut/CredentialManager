package credential

import (
	"encoding/json"
	"moqui/runtime/component/CredentialManager/src/test/go/utils"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func GetBaseUrl() string {
	url := os.Getenv("MOQUI_URL")
	if url == "" {
		url = "http://moqui:8080"
	}
	return url
}

func NewRequest(t *testing.T, method, path, username, password string, params map[string]string) *utils.GenericValue {
	u, err := url.Parse(GetBaseUrl() + path)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}
	return utils.GenericValueOf(result)
}

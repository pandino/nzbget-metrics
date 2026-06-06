package nzbget

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func rpcServer(t *testing.T, result any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "testuser" || pass != "testpass" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		raw, _ := json.Marshal(result)
		resp := map[string]any{"result": json.RawMessage(raw), "error": nil}
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestClientStatus(t *testing.T) {
	want := StatusResult{
		DownloadRate:    1024000,
		RemainingSizeLo: 512,
		RemainingSizeHi: 0,
		ServerStandBy:   false,
	}
	srv := rpcServer(t, want)
	defer srv.Close()

	c := &Client{
		url:      srv.URL,
		username: "testuser",
		password: "testpass",
		http:     &http.Client{},
	}

	got, err := c.Status()
	if err != nil {
		t.Fatal(err)
	}
	if got.DownloadRate != want.DownloadRate {
		t.Errorf("DownloadRate: got %v, want %v", got.DownloadRate, want.DownloadRate)
	}
	if got.RemainingSizeLo != want.RemainingSizeLo {
		t.Errorf("RemainingSizeLo: got %v, want %v", got.RemainingSizeLo, want.RemainingSizeLo)
	}
}

func TestClientVersion(t *testing.T) {
	srv := rpcServer(t, "23.0")
	defer srv.Close()

	c := &Client{
		url:      srv.URL,
		username: "testuser",
		password: "testpass",
		http:     &http.Client{},
	}

	v, err := c.Version()
	if err != nil {
		t.Fatal(err)
	}
	if v != "23.0" {
		t.Errorf("Version: got %q, want %q", v, "23.0")
	}
	// Second call returns cached value.
	v2, _ := c.Version()
	if v2 != "23.0" {
		t.Errorf("Version (cached): got %q", v2)
	}
}

func TestClientQueuedCount(t *testing.T) {
	items := []map[string]any{{"NZBID": 1}, {"NZBID": 2}}
	srv := rpcServer(t, items)
	defer srv.Close()

	c := &Client{
		url:      srv.URL,
		username: "testuser",
		password: "testpass",
		http:     &http.Client{},
	}

	n, err := c.QueuedCount()
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Errorf("QueuedCount: got %d, want 2", n)
	}
}

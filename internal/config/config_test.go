package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	// Point at a nonexistent config file so only defaults apply.
	t.Setenv("NZBGET_CONFIG", "/nonexistent/nzbget.conf")
	clearEnv(t)

	cfg := Load()

	if cfg.Host != "127.0.0.1" {
		t.Errorf("Host: got %q, want %q", cfg.Host, "127.0.0.1")
	}
	if cfg.Port != 6789 {
		t.Errorf("Port: got %d, want %d", cfg.Port, 6789)
	}
	if cfg.Username != "nzbget" {
		t.Errorf("Username: got %q, want %q", cfg.Username, "nzbget")
	}
	if cfg.Password != "tegbzn6789" {
		t.Errorf("Password: got %q, want %q", cfg.Password, "tegbzn6789")
	}
	if cfg.ListenAddr != ":9452" {
		t.Errorf("ListenAddr: got %q, want %q", cfg.ListenAddr, ":9452")
	}
	if cfg.MetricsPath != "/metrics" {
		t.Errorf("MetricsPath: got %q, want %q", cfg.MetricsPath, "/metrics")
	}
}

func TestEnvOverrides(t *testing.T) {
	t.Setenv("NZBGET_CONFIG", "/nonexistent/nzbget.conf")
	t.Setenv("NZBGET_HOST", "10.0.0.1")
	t.Setenv("NZBGET_PORT", "7000")
	t.Setenv("NZBGET_USERNAME", "admin")
	t.Setenv("NZBGET_PASSWORD", "secret")
	t.Setenv("METRICS_PORT", "9999")
	t.Setenv("METRICS_PATH", "/prom")

	cfg := Load()

	if cfg.Host != "10.0.0.1" {
		t.Errorf("Host: got %q", cfg.Host)
	}
	if cfg.Port != 7000 {
		t.Errorf("Port: got %d", cfg.Port)
	}
	if cfg.Username != "admin" {
		t.Errorf("Username: got %q", cfg.Username)
	}
	if cfg.Password != "secret" {
		t.Errorf("Password: got %q", cfg.Password)
	}
	if cfg.ListenAddr != ":9999" {
		t.Errorf("ListenAddr: got %q", cfg.ListenAddr)
	}
	if cfg.MetricsPath != "/prom" {
		t.Errorf("MetricsPath: got %q", cfg.MetricsPath)
	}
}

func TestConfFileOverride(t *testing.T) {
	clearEnv(t)

	dir := t.TempDir()
	confFile := filepath.Join(dir, "nzbget.conf")
	err := os.WriteFile(confFile, []byte(`
# comment
ControlIP=192.168.1.100
ControlPort=8080
ControlUsername=myuser
ControlPassword=mypass
`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("NZBGET_CONFIG", confFile)

	cfg := Load()

	if cfg.Host != "192.168.1.100" {
		t.Errorf("Host: got %q", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("Port: got %d", cfg.Port)
	}
	if cfg.Username != "myuser" {
		t.Errorf("Username: got %q", cfg.Username)
	}
	if cfg.Password != "mypass" {
		t.Errorf("Password: got %q", cfg.Password)
	}
}

func TestEnvTakesPriorityOverConf(t *testing.T) {
	dir := t.TempDir()
	confFile := filepath.Join(dir, "nzbget.conf")
	os.WriteFile(confFile, []byte("ControlPort=8080\n"), 0600)
	t.Setenv("NZBGET_CONFIG", confFile)
	t.Setenv("NZBGET_PORT", "9000")

	cfg := Load()
	if cfg.Port != 9000 {
		t.Errorf("Port: got %d, want 9000", cfg.Port)
	}
}

func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"NZBGET_HOST", "NZBGET_PORT", "NZBGET_USERNAME", "NZBGET_PASSWORD", "METRICS_PORT", "METRICS_PATH"} {
		t.Setenv(k, "")
	}
}

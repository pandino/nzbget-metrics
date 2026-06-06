package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host        string
	Port        int
	Username    string
	Password    string
	ListenAddr  string
	MetricsPath string
}

func Load() *Config {
	confPath := envOrDefault("NZBGET_CONFIG", "/config/nzbget.conf")

	fileVals := parseConfFile(confPath)

	host := firstNonEmpty(os.Getenv("NZBGET_HOST"), fileVals["ControlIP"], "127.0.0.1")
	portStr := firstNonEmpty(os.Getenv("NZBGET_PORT"), fileVals["ControlPort"], "6789")
	username := firstNonEmpty(os.Getenv("NZBGET_USERNAME"), fileVals["ControlUsername"], "nzbget")
	password := firstNonEmpty(os.Getenv("NZBGET_PASSWORD"), fileVals["ControlPassword"], "tegbzn6789")
	metricsPort := firstNonEmpty(os.Getenv("METRICS_PORT"), "", "9452")
	metricsPath := firstNonEmpty(os.Getenv("METRICS_PATH"), "", "/metrics")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 6789
	}

	listenPort, err := strconv.Atoi(metricsPort)
	if err != nil {
		listenPort = 9452
	}

	return &Config{
		Host:        host,
		Port:        port,
		Username:    username,
		Password:    password,
		ListenAddr:  ":" + strconv.Itoa(listenPort),
		MetricsPath: metricsPath,
	}
}

func parseConfFile(path string) map[string]string {
	vals := make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		return vals
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		vals[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return vals
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

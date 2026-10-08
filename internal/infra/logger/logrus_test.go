package logger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
)

func TestRequestDataCannotForgeLogLines(t *testing.T) {
	previous := config.EnvConfig
	config.EnvConfig = &config.Config{App: config.App{Env: config.ReleaseMode}}
	t.Cleanup(func() { config.EnvConfig = previous })
	logger := InitLog()
	if logger.Level != logrus.InfoLevel {
		t.Fatalf("release log level = %s", logger.Level)
	}
	var output bytes.Buffer
	logger.SetOutput(&output)
	logger.WithField("path", "/login\nforged entry").Info("request completed")

	scanner := bufio.NewScanner(&output)
	if !scanner.Scan() {
		t.Fatal("missing log entry")
	}
	var entry map[string]any
	if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
		t.Fatalf("log is not JSON: %v", err)
	}
	if entry["path"] != "/login\nforged entry" || scanner.Scan() {
		t.Fatalf("request data altered log record boundary: %q", output.String())
	}
}

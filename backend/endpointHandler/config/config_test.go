package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	content := `
server:
  port: "8000"

database:
  host: "localhost"
  user: "postgres"
  password: "00000"
  dbname: "jira2"
  port: "5432"
  sslmode: "disable"

connector:
  baseURL: "http://localhost:8080/api/v1/connector"
`

	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.Server.Port != "8000" {
		t.Errorf("expected server.port=8000, got %s", cfg.Server.Port)
	}
	if cfg.Database.DBName != "jira2" {
		t.Errorf("expected database.dbname=jira2, got %s", cfg.Database.DBName)
	}
	if cfg.Connector.BaseURL != "http://localhost:8080/api/v1/connector" {
		t.Errorf("unexpected connector.baseURL: %s", cfg.Connector.BaseURL)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent-config-file.yaml")
	if err == nil {
		t.Errorf("expected error for missing config file, got nil")
	}
}


package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPrefersEnvOverFile(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv(envConfigDir, tempDir)
	t.Setenv("APIMUX_BASE_URL", "http://env.example")
	t.Setenv("APIMUX_API_KEY", "env-key")

	if err := os.WriteFile(filepath.Join(tempDir, "config.json"), []byte(`{"base_url":"http://file.example","api_key":"file-key"}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.BaseURL != "http://env.example" || cfg.APIKey != "env-key" {
		t.Fatalf("unexpected config: %#v", cfg)
	}
}

func TestSaveWritesConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv(envConfigDir, tempDir)

	if err := Save(Config{BaseURL: "http://127.0.0.1:8081", APIKey: "abc"}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.BaseURL != "http://127.0.0.1:8081" || cfg.APIKey != "abc" {
		t.Fatalf("unexpected config: %#v", cfg)
	}
}

func TestLoadDetailedTracksValueSources(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv(envConfigDir, tempDir)
	t.Setenv("APIMUX_BASE_URL", "http://env.example")

	if err := os.WriteFile(filepath.Join(tempDir, "config.json"), []byte(`{"base_url":"http://file.example","api_key":"file-key"}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	loaded, err := LoadDetailed()
	if err != nil {
		t.Fatalf("load detailed: %v", err)
	}

	if loaded.Config.BaseURL != "http://env.example" {
		t.Fatalf("unexpected base url: %s", loaded.Config.BaseURL)
	}
	if loaded.Config.APIKey != "file-key" {
		t.Fatalf("unexpected api key: %s", loaded.Config.APIKey)
	}
	if loaded.Sources["base_url"] != "env" {
		t.Fatalf("expected env source for base_url, got %q", loaded.Sources["base_url"])
	}
	if loaded.Sources["api_key"] != "file" {
		t.Fatalf("expected file source for api_key, got %q", loaded.Sources["api_key"])
	}
	if loaded.Path == "" {
		t.Fatal("expected config path")
	}
}

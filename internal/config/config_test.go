package config

import (
	"os"
	"path/filepath"
	"strings"
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

func TestPathDefaultsToHomeApimuxConfig(t *testing.T) {
	t.Setenv(envConfigDir, "")
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	path, err := Path()
	if err != nil {
		t.Fatalf("config path: %v", err)
	}

	expected := filepath.Join(homeDir, ".apimux", "config.json")
	if path != expected {
		t.Fatalf("expected %s, got %s", expected, path)
	}
}

func TestLoadFallsBackToLegacyUserConfigPath(t *testing.T) {
	t.Setenv(envConfigDir, "")
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(homeDir, ".config"))

	legacyBase, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("user config dir: %v", err)
	}
	legacyPath := filepath.Join(legacyBase, "apimux", "config.json")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("mkdir legacy config dir: %v", err)
	}
	if err := os.WriteFile(legacyPath, []byte(`{"base_url":"http://legacy.example","api_key":"legacy-key"}`), 0o600); err != nil {
		t.Fatalf("write legacy config: %v", err)
	}

	loaded, err := LoadDetailed()
	if err != nil {
		t.Fatalf("load detailed: %v", err)
	}

	expectedPath := filepath.Join(homeDir, ".apimux", "config.json")
	if loaded.Path != expectedPath {
		t.Fatalf("expected primary path %s, got %s", expectedPath, loaded.Path)
	}
	if loaded.LegacyPath != legacyPath {
		t.Fatalf("expected legacy path %s, got %s", legacyPath, loaded.LegacyPath)
	}
	if loaded.Config.BaseURL != "http://legacy.example" || loaded.Config.APIKey != "legacy-key" {
		t.Fatalf("unexpected config: %#v", loaded.Config)
	}
}

func TestSaveReadsLegacyButWritesPrimaryPath(t *testing.T) {
	t.Setenv(envConfigDir, "")
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(homeDir, ".config"))

	legacyBase, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("user config dir: %v", err)
	}
	legacyPath := filepath.Join(legacyBase, "apimux", "config.json")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("mkdir legacy config dir: %v", err)
	}
	if err := os.WriteFile(legacyPath, []byte(`{"base_url":"http://legacy.example","api_key":"legacy-key"}`), 0o600); err != nil {
		t.Fatalf("write legacy config: %v", err)
	}

	if err := Save(Config{APIKey: "new-key"}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	primaryPath := filepath.Join(homeDir, ".apimux", "config.json")
	body, err := os.ReadFile(primaryPath)
	if err != nil {
		t.Fatalf("read primary config: %v", err)
	}
	if !strings.Contains(string(body), `"base_url": "http://legacy.example"`) {
		t.Fatalf("expected legacy base URL preserved in primary config, got %s", string(body))
	}
	if !strings.Contains(string(body), `"api_key": "new-key"`) {
		t.Fatalf("expected new API key in primary config, got %s", string(body))
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

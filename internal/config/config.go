package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const envConfigDir = "APIMUX_CONFIG_DIR"

type Config struct {
	BaseURL string `json:"base_url,omitempty"`
	APIKey  string `json:"api_key,omitempty"`
}

type Loaded struct {
	Config     Config            `json:"config"`
	Path       string            `json:"path"`
	LegacyPath string            `json:"legacy_path,omitempty"`
	Sources    map[string]string `json:"sources"`
}

func Load() (Config, error) {
	loaded, err := LoadDetailed()
	if err != nil {
		return Config{}, err
	}
	return loaded.Config, nil
}

func LoadDetailed() (Loaded, error) {
	cfg := Config{
		BaseURL: "http://127.0.0.1:8081",
	}
	sources := map[string]string{
		"base_url": "default",
		"api_key":  "unset",
	}

	filePath, legacyPath, fileCfg, usedLegacy, err := readFile()
	if err != nil {
		return Loaded{}, err
	}
	if fileCfg.BaseURL != "" {
		cfg.BaseURL = fileCfg.BaseURL
		sources["base_url"] = "file"
	}
	if fileCfg.APIKey != "" {
		cfg.APIKey = fileCfg.APIKey
		sources["api_key"] = "file"
	}

	if value := strings.TrimSpace(os.Getenv("APIMUX_BASE_URL")); value != "" {
		cfg.BaseURL = value
		sources["base_url"] = "env"
	}
	if value := strings.TrimSpace(os.Getenv("APIMUX_API_KEY")); value != "" {
		cfg.APIKey = value
		sources["api_key"] = "env"
	}

	loaded := Loaded{
		Config:  cfg,
		Path:    filePath,
		Sources: sources,
	}
	if usedLegacy {
		loaded.LegacyPath = legacyPath
	}
	return loaded, nil
}

func Save(update Config) error {
	_, _, current, _, err := readFile()
	if err != nil {
		return err
	}
	if update.BaseURL != "" {
		current.BaseURL = strings.TrimSpace(update.BaseURL)
	}
	if update.APIKey != "" {
		current.APIKey = strings.TrimSpace(update.APIKey)
	}

	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	body, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return os.WriteFile(path, body, 0o600)
}

func Path() (string, error) {
	primaryPath, _, err := paths()
	return primaryPath, err
}

func readFile() (string, string, Config, bool, error) {
	primaryPath, legacyPath, err := paths()
	if err != nil {
		return "", "", Config{}, false, err
	}
	body, err := os.ReadFile(primaryPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", "", Config{}, false, err
		}
		if legacyPath == "" || legacyPath == primaryPath {
			return primaryPath, legacyPath, Config{}, false, nil
		}
		legacyBody, legacyErr := os.ReadFile(legacyPath)
		if legacyErr != nil {
			if errors.Is(legacyErr, os.ErrNotExist) {
				return primaryPath, legacyPath, Config{}, false, nil
			}
			return "", "", Config{}, false, legacyErr
		}
		var legacyCfg Config
		if err := json.Unmarshal(legacyBody, &legacyCfg); err != nil {
			return "", "", Config{}, false, err
		}
		return primaryPath, legacyPath, legacyCfg, true, nil
	}

	var cfg Config
	if err := json.Unmarshal(body, &cfg); err != nil {
		return "", "", Config{}, false, err
	}
	return primaryPath, legacyPath, cfg, false, nil
}

func paths() (string, string, error) {
	if dir := strings.TrimSpace(os.Getenv(envConfigDir)); dir != "" {
		path := filepath.Join(dir, "config.json")
		return path, "", nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	primaryPath := filepath.Join(home, ".apimux", "config.json")
	legacyBase, err := os.UserConfigDir()
	if err != nil {
		return primaryPath, "", nil
	}
	return primaryPath, filepath.Join(legacyBase, "apimux", "config.json"), nil
}

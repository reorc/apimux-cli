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
	Config  Config            `json:"config"`
	Path    string            `json:"path"`
	Sources map[string]string `json:"sources"`
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

	if fileCfg, err := readFile(); err != nil {
		return Loaded{}, err
	} else {
		if fileCfg.BaseURL != "" {
			cfg.BaseURL = fileCfg.BaseURL
			sources["base_url"] = "file"
		}
		if fileCfg.APIKey != "" {
			cfg.APIKey = fileCfg.APIKey
			sources["api_key"] = "file"
		}
	}

	if value := strings.TrimSpace(os.Getenv("APIMUX_BASE_URL")); value != "" {
		cfg.BaseURL = value
		sources["base_url"] = "env"
	}
	if value := strings.TrimSpace(os.Getenv("APIMUX_API_KEY")); value != "" {
		cfg.APIKey = value
		sources["api_key"] = "env"
	}

	configPath, err := path()
	if err != nil {
		return Loaded{}, err
	}

	return Loaded{
		Config:  cfg,
		Path:    configPath,
		Sources: sources,
	}, nil
}

func Save(update Config) error {
	current, err := readFile()
	if err != nil {
		return err
	}
	if update.BaseURL != "" {
		current.BaseURL = strings.TrimSpace(update.BaseURL)
	}
	if update.APIKey != "" {
		current.APIKey = strings.TrimSpace(update.APIKey)
	}

	path, err := path()
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
	return path()
}

func readFile() (Config, error) {
	path, err := path()
	if err != nil {
		return Config{}, err
	}
	body, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(body, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func path() (string, error) {
	if dir := strings.TrimSpace(os.Getenv(envConfigDir)); dir != "" {
		return filepath.Join(dir, "config.json"), nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "apimux", "config.json"), nil
}

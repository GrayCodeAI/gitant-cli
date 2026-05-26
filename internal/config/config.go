package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Settings persisted in ~/.gitant/config.json (like gh's config.yml).
type Settings struct {
	DaemonURL string `json:"daemon_url,omitempty"`
	WebURL    string `json:"web_url,omitempty"`
	UCANToken string `json:"ucan_token,omitempty"`
}

func path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gitant", "config.json"), nil
}

// Load reads settings from disk. Missing file returns zero values.
func Load() (Settings, error) {
	p, err := path()
	if err != nil {
		return Settings{}, err
	}
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return Settings{}, nil
	}
	if err != nil {
		return Settings{}, err
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return Settings{}, err
	}
	return s, nil
}

// Save writes settings to disk.
func Save(s Settings) error {
	p, err := path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

// Get returns a config value by key: daemon_url, web_url, ucan_token.
func Get(key string) (string, error) {
	s, err := Load()
	if err != nil {
		return "", err
	}
	switch key {
	case "daemon_url":
		return s.DaemonURL, nil
	case "web_url":
		return s.WebURL, nil
	case "ucan_token":
		return s.UCANToken, nil
	default:
		return "", os.ErrInvalid
	}
}

// Set updates one config key and saves.
func Set(key, value string) error {
	s, err := Load()
	if err != nil {
		return err
	}
	switch key {
	case "daemon_url":
		s.DaemonURL = value
	case "web_url":
		s.WebURL = value
	case "ucan_token":
		s.UCANToken = value
	default:
		return os.ErrInvalid
	}
	return Save(s)
}

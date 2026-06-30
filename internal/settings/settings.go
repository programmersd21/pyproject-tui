// Package settings stores persistent UI preferences.
package settings

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Config stores persistent application preferences.
type Config struct {
	Theme           string `toml:"theme"`
	AnimationsOn    bool   `toml:"animations"`
	UIDensity       string `toml:"ui_density"`   // "compact", "normal", "comfortable"
	BorderStyle     string `toml:"border_style"` // "rounded", "normal", "thick", "double"
	ShowLineNumbers bool   `toml:"show_line_numbers"`
}

const appDirName = "pyproject-tui"

// DefaultConfig returns the default settings.
func DefaultConfig() Config {
	return Config{
		Theme:           "tokyo-night",
		AnimationsOn:    true,
		UIDensity:       "normal",
		BorderStyle:     "rounded",
		ShowLineNumbers: false,
	}
}

// DefaultPath returns the OS-specific config file path for pyproject-tui.
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, appDirName, "config.toml"), nil
}

// Load reads a settings file.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	// Validate and set defaults for empty values
	if cfg.Theme == "" {
		cfg.Theme = "tokyo-night"
	}
	if cfg.UIDensity == "" {
		cfg.UIDensity = "normal"
	}
	if cfg.BorderStyle == "" {
		cfg.BorderStyle = "rounded"
	}

	return cfg, nil
}

// LoadDefault reads the config from the OS-specific default location.
func LoadDefault() (Config, string, error) {
	path, err := DefaultPath()
	if err != nil {
		return DefaultConfig(), "", err
	}
	cfg, err := Load(path)
	return cfg, path, err
}

// Save writes the settings file atomically.
func Save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// SaveDefault writes the config to the OS-specific default location.
func SaveDefault(cfg Config) (string, error) {
	path, err := DefaultPath()
	if err != nil {
		return "", err
	}
	return path, Save(path, cfg)
}

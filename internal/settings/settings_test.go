package settings

import (
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Theme == "" || cfg.UIDensity == "" || cfg.BorderStyle == "" {
		t.Error("Default config should have all required fields")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, ".test.toml")

	testCfg := Config{
		Theme:           "nord",
		AnimationsOn:    false,
		UIDensity:       "compact",
		BorderStyle:     "thick",
		ShowLineNumbers: true,
	}

	if err := Save(testPath, testCfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(testPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Theme != testCfg.Theme || loaded.UIDensity != testCfg.UIDensity {
		t.Error("Loaded config doesn't match saved")
	}
}

func TestLoadNonexistent(t *testing.T) {
	cfg, err := Load("/nonexistent/.settings.toml")
	if err != nil {
		t.Errorf("Should return default config, got error: %v", err)
	}
	if cfg.Theme == "" {
		t.Error("Should return valid default config")
	}
}

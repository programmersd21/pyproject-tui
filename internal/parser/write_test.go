package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteRoundTrip(t *testing.T) {
	p, err := Load("../../testdata/complete.toml")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	tmp := filepath.Join(t.TempDir(), "pyproject.toml")
	p.FilePath = tmp
	p.Project.Description = "Updated"
	if writeErr := Write(p); writeErr != nil {
		t.Fatalf("Write: %v", writeErr)
	}
	if _, statErr := os.Stat(tmp); statErr != nil {
		t.Fatalf("stat written file: %v", statErr)
	}
	reloaded, reloadErr := Load(tmp)
	if reloadErr != nil {
		t.Fatalf("reload: %v", reloadErr)
	}
	if reloaded.Project.Description != "Updated" {
		t.Fatalf("description = %q", reloaded.Project.Description)
	}
}

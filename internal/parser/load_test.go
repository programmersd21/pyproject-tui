package parser

import "testing"

func TestLoad(t *testing.T) {
	p, err := Load("../../testdata/complete.toml")
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if p.Project == nil || p.Project.Name != "my-package" {
		t.Fatalf("unexpected project: %#v", p.Project)
	}
	if p.BuildSystem == nil || p.BuildSystem.BuildBackend != "hatchling.build" {
		t.Fatalf("unexpected build system: %#v", p.BuildSystem)
	}
}

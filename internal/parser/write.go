// Package parser loads and writes pyproject.toml documents.
package parser

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Write rebuilds the pyproject.toml file from the typed structure and writes it atomically.
func Write(p *PyProject) error {
	if p == nil {
		return nil
	}
	root := deepCopyMap(p.Raw)
	if root == nil {
		root = map[string]any{}
	}
	if p.Project != nil {
		root["project"] = projectToMap(p.Project)
	}
	if p.BuildSystem != nil {
		root["build-system"] = buildSystemToMap(p.BuildSystem)
	}
	if len(p.Tool) > 0 {
		tool := map[string]any{}
		for k, v := range p.Tool {
			tool[k] = deepCopyAny(v)
		}
		root["tool"] = tool
	}

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	enc.SetIndentTables(true)
	if err := enc.Encode(root); err != nil {
		return err
	}

	tmp := p.FilePath + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, p.FilePath); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// projectToMap converts a typed project section into a TOML-ready map.
func projectToMap(p *ProjectSection) map[string]any {
	if p == nil {
		return nil
	}
	out := map[string]any{}
	if p.Name != "" {
		out["name"] = p.Name
	}
	if p.Version != "" {
		out["version"] = p.Version
	}
	if p.Description != "" {
		out["description"] = p.Description
	}
	if p.Readme != nil {
		out["readme"] = p.Readme
	}
	if p.RequiresPython != "" {
		out["requires-python"] = p.RequiresPython
	}
	if p.License != nil {
		out["license"] = p.License
	}
	if len(p.Authors) > 0 {
		out["authors"] = p.Authors
	}
	if len(p.Maintainers) > 0 {
		out["maintainers"] = p.Maintainers
	}
	if len(p.Keywords) > 0 {
		out["keywords"] = p.Keywords
	}
	if len(p.Classifiers) > 0 {
		out["classifiers"] = p.Classifiers
	}
	if len(p.URLs) > 0 {
		out["urls"] = p.URLs
	}
	if len(p.Dependencies) > 0 {
		out["dependencies"] = p.Dependencies
	}
	if len(p.OptionalDependencies) > 0 {
		out["optional-dependencies"] = p.OptionalDependencies
	}
	if len(p.Scripts) > 0 {
		out["scripts"] = p.Scripts
	}
	if len(p.EntryPoints) > 0 {
		out["entry-points"] = p.EntryPoints
	}
	if len(p.Dynamic) > 0 {
		out["dynamic"] = p.Dynamic
	}
	return out
}

// buildSystemToMap converts a typed build-system section into a TOML-ready map.
func buildSystemToMap(b *BuildSystemSection) map[string]any {
	if b == nil {
		return nil
	}
	out := map[string]any{}
	if len(b.Requires) > 0 {
		out["requires"] = b.Requires
	}
	if b.BuildBackend != "" {
		out["build-backend"] = b.BuildBackend
	}
	if len(b.BackendPath) > 0 {
		out["backend-path"] = b.BackendPath
	}
	return out
}

// WriteFile writes data to path atomically using a temp file and rename.
func WriteFile(path string, data []byte) error {
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

// EnsureDir makes parent directories for a path.
func EnsureDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0o755)
}

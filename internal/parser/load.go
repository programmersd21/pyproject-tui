// Package parser loads and writes pyproject.toml documents.
package parser

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// IsSyntaxError reports whether err is a TOML syntax error.
func IsSyntaxError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "toml")
}

// LoadRaw reads and parses the raw TOML document as a map.
func LoadRaw(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	raw := map[string]any{}
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// Load reads a pyproject.toml file from disk and parses it into a PyProject.
func Load(path string) (*PyProject, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := map[string]any{}
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	pp := &PyProject{
		Raw:      raw,
		FilePath: path,
	}
	if projectRaw, ok := raw["project"].(map[string]any); ok {
		pp.Project = parseProject(projectRaw)
	}
	if buildRaw, ok := raw["build-system"].(map[string]any); ok {
		pp.BuildSystem = parseBuildSystem(buildRaw)
	}
	if toolRaw, ok := raw["tool"].(map[string]any); ok && len(toolRaw) > 0 {
		pp.Tool = map[string]any{}
		for k, v := range toolRaw {
			pp.Tool[k] = v
		}
	}
	return pp, nil
}

func parseProject(m map[string]any) *ProjectSection {
	p := &ProjectSection{}
	if v, ok := m["name"].(string); ok {
		p.Name = v
	}
	if v, ok := m["version"].(string); ok {
		p.Version = v
	}
	if v, ok := m["description"].(string); ok {
		p.Description = v
	}
	p.Readme = m["readme"]
	if v, ok := m["requires-python"].(string); ok {
		p.RequiresPython = v
	}
	p.License = m["license"]
	p.Authors = parseAuthorArray(m["authors"])
	p.Maintainers = parseAuthorArray(m["maintainers"])
	p.Keywords = parseStringSlice(m["keywords"])
	p.Classifiers = parseStringSlice(m["classifiers"])
	p.URLs = parseStringMap(m["urls"])
	if deps, ok := m["dependencies"].([]any); ok {
		p.Dependencies = parseStringSlice(deps)
	}
	if depTable, ok := m["dependencies"].(map[string]any); ok {
		p.Dependencies = parseStringSlice(depTable["requires"])
	}
	p.OptionalDependencies = parseMapOfStringSlice(m["optional-dependencies"])
	p.Scripts = parseStringMap(m["scripts"])
	p.EntryPoints = parseMapOfStringMap(m["entry-points"])
	p.Dynamic = parseStringSlice(m["dynamic"])
	return p
}

func parseBuildSystem(m map[string]any) *BuildSystemSection {
	b := &BuildSystemSection{}
	b.Requires = parseStringSlice(m["requires"])
	if v, ok := m["build-backend"].(string); ok {
		b.BuildBackend = v
	}
	b.BackendPath = parseStringSlice(m["backend-path"])
	return b
}

func parseAuthorArray(v any) []AuthorEntry {
	items, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]AuthorEntry, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			a := AuthorEntry{}
			if v, ok := m["name"].(string); ok {
				a.Name = v
			}
			if v, ok := m["email"].(string); ok {
				a.Email = v
			}
			out = append(out, a)
		}
	}
	return out
}

func parseStringSlice(v any) []string {
	switch val := v.(type) {
	case []string:
		return append([]string(nil), val...)
	case []any:
		out := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	case string:
		return []string{val}
	default:
		return nil
	}
}

func parseStringMap(v any) map[string]string {
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	out := map[string]string{}
	for k, val := range m {
		if s, ok := val.(string); ok {
			out[k] = s
		}
	}
	return out
}

func parseMapOfStringSlice(v any) map[string][]string {
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	out := map[string][]string{}
	for k, val := range m {
		out[k] = parseStringSlice(val)
	}
	return out
}

func parseMapOfStringMap(v any) map[string]map[string]string {
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	out := map[string]map[string]string{}
	for k, val := range m {
		if inner, ok := val.(map[string]any); ok {
			out[k] = parseStringMap(inner)
		}
	}
	return out
}

// NewEmpty returns a minimal PyProject with defaults for a new file.
func NewEmpty(path string) *PyProject {
	return &PyProject{
		Project: &ProjectSection{
			Name:           "",
			Version:        "0.1.0",
			Description:    "",
			RequiresPython: ">=3.11",
		},
		BuildSystem: &BuildSystemSection{
			Requires:     []string{"hatchling"},
			BuildBackend: "hatchling.build",
		},
		Tool:     map[string]any{},
		Raw:      map[string]any{},
		FilePath: path,
	}
}

func deepCopyMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		dst[k] = deepCopyAny(src[k])
	}
	return dst
}

func deepCopyAny(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return deepCopyMap(val)
	case map[string]string:
		out := make(map[string]string, len(val))
		for k, vv := range val {
			out[k] = vv
		}
		return out
	case []any:
		out := make([]any, len(val))
		for i := range val {
			out[i] = deepCopyAny(val[i])
		}
		return out
	case []string:
		out := make([]string, len(val))
		copy(out, val)
		return out
	case []AuthorEntry:
		out := make([]AuthorEntry, len(val))
		copy(out, val)
		return out
	default:
		return v
	}
}

// Clone returns a deep copy of the PyProject.
func Clone(p *PyProject) *PyProject {
	if p == nil {
		return nil
	}
	cp := &PyProject{
		Project:     cloneProject(p.Project),
		BuildSystem: cloneBuildSystem(p.BuildSystem),
		Tool:        map[string]any{},
		Raw:         deepCopyMap(p.Raw),
		FilePath:    p.FilePath,
	}
	for k, v := range p.Tool {
		cp.Tool[k] = deepCopyAny(v)
	}
	return cp
}

func cloneProject(p *ProjectSection) *ProjectSection {
	if p == nil {
		return nil
	}
	cp := *p
	cp.Authors = append([]AuthorEntry(nil), p.Authors...)
	cp.Maintainers = append([]AuthorEntry(nil), p.Maintainers...)
	cp.Keywords = append([]string(nil), p.Keywords...)
	cp.Classifiers = append([]string(nil), p.Classifiers...)
	cp.Dependencies = append([]string(nil), p.Dependencies...)
	cp.Dynamic = append([]string(nil), p.Dynamic...)
	if p.URLs != nil {
		cp.URLs = map[string]string{}
		for k, v := range p.URLs {
			cp.URLs[k] = v
		}
	}
	if p.OptionalDependencies != nil {
		cp.OptionalDependencies = map[string][]string{}
		for k, v := range p.OptionalDependencies {
			cp.OptionalDependencies[k] = append([]string(nil), v...)
		}
	}
	if p.Scripts != nil {
		cp.Scripts = map[string]string{}
		for k, v := range p.Scripts {
			cp.Scripts[k] = v
		}
	}
	if p.EntryPoints != nil {
		cp.EntryPoints = map[string]map[string]string{}
		for k, v := range p.EntryPoints {
			inner := map[string]string{}
			for kk, vv := range v {
				inner[kk] = vv
			}
			cp.EntryPoints[k] = inner
		}
	}
	return &cp
}

func cloneBuildSystem(b *BuildSystemSection) *BuildSystemSection {
	if b == nil {
		return nil
	}
	cp := *b
	cp.Requires = append([]string(nil), b.Requires...)
	cp.BackendPath = append([]string(nil), b.BackendPath...)
	return &cp
}

// String returns a display string for a section ID.
func (id SectionID) String() string {
	switch id.Kind {
	case KindProject:
		return "[project]"
	case KindBuildSystem:
		return "[build-system]"
	case KindTool:
		return fmt.Sprintf("[tool.%s]", id.ToolKey)
	default:
		return "[unknown]"
	}
}

// ParseSectionTitle converts a sidebar title into a SectionID.
func ParseSectionTitle(title string) SectionID {
	switch {
	case title == "[project]":
		return SectionID{Kind: KindProject}
	case title == "[build-system]":
		return SectionID{Kind: KindBuildSystem}
	case strings.HasPrefix(title, "[tool.") && strings.HasSuffix(title, "]"):
		return SectionID{Kind: KindTool, ToolKey: strings.TrimSuffix(strings.TrimPrefix(title, "[tool."), "]")}
	default:
		return SectionID{}
	}
}

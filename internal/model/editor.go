// Package model contains Bubble Tea models for the pyproject-tui application.
package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/programmersd21/pyproject-tui/internal/parser"
	"github.com/programmersd21/pyproject-tui/internal/theme"
	"github.com/programmersd21/pyproject-tui/internal/ui"
)

// EditMode identifies the editor interaction mode.
type EditMode int

const (
	// EditModeNormal is the browsing mode.
	EditModeNormal EditMode = iota
	// EditModeEditValue edits an existing value.
	EditModeEditValue
	// EditModeAddKey prompts for a new key.
	EditModeAddKey
	// EditModeAddValue prompts for a new value.
	EditModeAddValue
)

// EditorModel implements the right pane.
type EditorModel struct {
	section         parser.SectionID
	fields          []parser.Field
	cursor          int
	editMode        EditMode
	inputKey        textinput.Model
	inputValue      textinput.Model
	viewport        viewport.Model
	width           int
	height          int
	showLineNumbers bool
	focused         bool
}

// NewEditor creates a new editor model.
func NewEditor() EditorModel {
	key := textinput.New()
	key.Placeholder = "key_name"
	key.Prompt = ""

	val := textinput.New()
	val.Placeholder = "value (e.g. \"my-app\", 123, [\"a\", \"b\"])"
	val.Prompt = ""

	return EditorModel{
		inputKey:   key,
		inputValue: val,
		viewport:   viewport.New(0, 0),
	}
}

// getFieldType determines the field type for rendering/editing.
func getFieldType(kind parser.SectionKind, key string, val any) parser.FieldType {
	switch kind {
	case parser.KindProject:
		switch key {
		case "name", "version", "description", "requires-python":
			return parser.FieldTypeString
		case "readme", "license", "entry-points":
			return parser.FieldTypeAny
		case "authors", "maintainers":
			return parser.FieldTypeAuthorArray
		case "keywords", "classifiers", "dependencies", "dynamic":
			return parser.FieldTypeStringArray
		case "urls", "scripts":
			return parser.FieldTypeStringMap
		case "optional-dependencies":
			return parser.FieldTypeStringMapOfArray
		}
	case parser.KindBuildSystem:
		switch key {
		case "requires", "backend-path":
			return parser.FieldTypeStringArray
		case "build-backend":
			return parser.FieldTypeString
		}
	}

	if val != nil {
		switch val.(type) {
		case bool:
			return parser.FieldTypeBool
		case int, int64:
			return parser.FieldTypeInt
		case []string, []any:
			return parser.FieldTypeStringArray
		case map[string]string:
			return parser.FieldTypeStringMap
		case map[string][]string:
			return parser.FieldTypeStringMapOfArray
		}
	}
	return parser.FieldTypeAny
}

// SetSection rebuilds fields for the given section.
//
//nolint:gocyclo // Field mapping requires exhaustive switch statements
func (m *EditorModel) SetSection(id parser.SectionID, p *parser.PyProject) {
	oldKey := ""
	if m.cursor >= 0 && m.cursor < len(m.fields) {
		oldKey = m.fields[m.cursor].Key
	}

	m.section = id
	m.fields = m.fields[:0]

	switch id.Kind {
	case parser.KindProject:
		m.fields = appendProjectFields(m.fields, p)
	case parser.KindBuildSystem:
		m.fields = appendBuildSystemFields(m.fields, p)
	case parser.KindTool:
		m.fields = appendToolFields(m.fields, id.ToolKey, p)
	}

	// Restore cursor position
	m.cursor = 0
	if oldKey != "" {
		for i, f := range m.fields {
			if f.Key == oldKey {
				m.cursor = i
				break
			}
		}
	}
	if m.cursor >= len(m.fields) {
		m.cursor = len(m.fields) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}

	m.rebuildViewport()
}

// SetFocused changes focus styling for the editor pane.
func (m *EditorModel) SetFocused(focused bool) {
	m.focused = focused
}

func appendProjectFields(fields []parser.Field, p *parser.PyProject) []parser.Field {
	if p == nil || p.Project == nil {
		return fields
	}
	pr := p.Project
	appendField := func(key string, val any) {
		fields = append(fields, parser.Field{Key: key, Value: val, Type: getFieldType(parser.KindProject, key, val)})
	}
	if pr.Name != "" {
		appendField("name", pr.Name)
	}
	if pr.Version != "" {
		appendField("version", pr.Version)
	}
	if pr.Description != "" {
		appendField("description", pr.Description)
	}
	if pr.Readme != nil {
		appendField("readme", pr.Readme)
	}
	if pr.RequiresPython != "" {
		appendField("requires-python", pr.RequiresPython)
	}
	if pr.License != nil {
		appendField("license", pr.License)
	}
	if len(pr.Authors) > 0 {
		appendField("authors", pr.Authors)
	}
	if len(pr.Maintainers) > 0 {
		appendField("maintainers", pr.Maintainers)
	}
	if len(pr.Keywords) > 0 {
		appendField("keywords", pr.Keywords)
	}
	if len(pr.Classifiers) > 0 {
		appendField("classifiers", pr.Classifiers)
	}
	if len(pr.URLs) > 0 {
		appendField("urls", pr.URLs)
	}
	if len(pr.Dependencies) > 0 {
		appendField("dependencies", pr.Dependencies)
	}
	if len(pr.OptionalDependencies) > 0 {
		appendField("optional-dependencies", pr.OptionalDependencies)
	}
	if len(pr.Scripts) > 0 {
		appendField("scripts", pr.Scripts)
	}
	if len(pr.EntryPoints) > 0 {
		appendField("entry-points", pr.EntryPoints)
	}
	if len(pr.Dynamic) > 0 {
		appendField("dynamic", pr.Dynamic)
	}
	return fields
}

func appendBuildSystemFields(fields []parser.Field, p *parser.PyProject) []parser.Field {
	if p == nil || p.BuildSystem == nil {
		return fields
	}
	b := p.BuildSystem
	if len(b.Requires) > 0 {
		fields = append(fields, parser.Field{Key: "requires", Value: b.Requires, Type: getFieldType(parser.KindBuildSystem, "requires", b.Requires)})
	}
	if b.BuildBackend != "" {
		fields = append(fields, parser.Field{Key: "build-backend", Value: b.BuildBackend, Type: getFieldType(parser.KindBuildSystem, "build-backend", b.BuildBackend)})
	}
	if len(b.BackendPath) > 0 {
		fields = append(fields, parser.Field{Key: "backend-path", Value: b.BackendPath, Type: getFieldType(parser.KindBuildSystem, "backend-path", b.BackendPath)})
	}
	return fields
}

func appendToolFields(fields []parser.Field, toolKey string, p *parser.PyProject) []parser.Field {
	if p == nil || p.Tool == nil {
		return fields
	}
	tool, ok := p.Tool[toolKey]
	if !ok {
		return fields
	}
	root, ok := tool.(map[string]any)
	if !ok {
		return fields
	}
	keys := make([]string, 0, len(root))
	for k := range root {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		val := root[k]
		fields = append(fields, parser.Field{Key: k, Value: val, Type: getFieldType(parser.KindTool, k, val)})
	}
	return fields
}

// Update updates the editor model.
func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if m.editMode != EditModeNormal {
			return m.updateEditMode(key)
		}
		return m.handleNormalKey(key)
	}
	return m, nil
}

func (m EditorModel) handleNormalKey(msg tea.KeyMsg) (EditorModel, tea.Cmd) {
	switch msg.String() {
	case "down", "j":
		if m.cursor < len(m.fields)-1 {
			m.cursor++
		}
		m.rebuildViewport()
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
		m.rebuildViewport()
	case "e", "enter":
		if len(m.fields) > 0 {
			m.editMode = EditModeEditValue
			m.inputValue.SetValue(m.currentValueString())
			m.inputValue.Focus()
		}
	case "a":
		m.editMode = EditModeAddKey
		m.inputKey.SetValue("")
		m.inputKey.Focus()
	}
	return m, nil
}

func (m EditorModel) updateEditMode(msg tea.KeyMsg) (EditorModel, tea.Cmd) {
	switch m.editMode {
	case EditModeEditValue:
		var cmd tea.Cmd
		m.inputValue, cmd = m.inputValue.Update(msg)
		switch msg.String() {
		case "esc":
			m.editMode = EditModeNormal
			return m, nil
		case "enter":
			m.editMode = EditModeNormal
			return m, func() tea.Msg {
				return FieldEditedMsg{Key: m.currentKey(), NewValue: m.inputValue.Value()}
			}
		}
		return m, cmd
	case EditModeAddKey:
		var cmd tea.Cmd
		m.inputKey, cmd = m.inputKey.Update(msg)
		switch msg.String() {
		case "esc":
			m.editMode = EditModeNormal
			return m, nil
		case "enter":
			keyVal := strings.TrimSpace(m.inputKey.Value())
			if keyVal == "" {
				m.editMode = EditModeNormal
				return m, nil
			}
			m.editMode = EditModeAddValue
			m.inputValue.SetValue("")
			m.inputValue.Focus()
			return m, nil
		}
		return m, cmd
	case EditModeAddValue:
		var cmd tea.Cmd
		m.inputValue, cmd = m.inputValue.Update(msg)
		switch msg.String() {
		case "esc":
			m.editMode = EditModeNormal
			return m, nil
		case "enter":
			m.editMode = EditModeNormal
			return m, func() tea.Msg {
				return FieldAddedMsg{Key: m.inputKey.Value(), Value: m.inputValue.Value()}
			}
		}
		return m, cmd
	}
	return m, nil
}

func (m EditorModel) currentValueString() string {
	if len(m.fields) == 0 {
		return ""
	}
	val := m.fields[m.cursor].Value
	switch v := val.(type) {
	case string:
		return v
	case []string:
		// Render as TOML array representation, e.g. ["a", "b"]
		parts := make([]string, len(v))
		for i, s := range v {
			parts[i] = fmt.Sprintf("%q", s)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]string:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, len(keys))
		for i, k := range keys {
			parts[i] = fmt.Sprintf("%s = %q", k, v[k])
		}
		return "{" + strings.Join(parts, ", ") + "}"
	}
	return ui.RenderAny(val, 0)
}

func (m EditorModel) currentKey() string {
	if len(m.fields) == 0 {
		return ""
	}
	return m.fields[m.cursor].Key
}

// SetShowLineNumbers toggles line number rendering.
func (m *EditorModel) SetShowLineNumbers(enabled bool) {
	m.showLineNumbers = enabled
	m.rebuildViewport()
}

func (m *EditorModel) rebuildViewport() {
	s := theme.Styles()
	if s == nil {
		m.viewport.SetContent("Loading...")
		return
	}

	lines := make([]string, 0, len(m.fields))
	for i, field := range m.fields {
		keyStr := field.Key
		if len(keyStr) > 22 {
			keyStr = ui.TruncateString(keyStr, 22)
		}

		prefix := "  "
		if m.showLineNumbers {
			prefix = fmt.Sprintf("%2d ", i+1)
		}
		var key string
		value := renderValue(field)
		if i == m.cursor {
			prefix = "→ "
			key = s.Cursor.Render(fmt.Sprintf("%-22s", keyStr))
		} else {
			key = s.Key.Render(fmt.Sprintf("%-22s", keyStr))
		}
		row := fmt.Sprintf("%s%s  %s", prefix, key, value)
		lines = append(lines, row)
	}
	if len(lines) == 0 {
		lines = append(lines, "  "+s.Value.Render("<no fields defined>"))
	}
	m.viewport.SetContent(strings.Join(lines, "\n"))
}

func renderValue(field parser.Field) string {
	s := theme.Styles()
	if s == nil {
		return fmt.Sprint(field.Value)
	}

	switch field.Type {
	case parser.FieldTypeString:
		return s.Value.Render(fmt.Sprintf("%q", field.Value))
	case parser.FieldTypeStringArray:
		if arr, ok := field.Value.([]string); ok {
			return s.ArrayValue.Render(ui.RenderStringSlice(arr, 5))
		}
	case parser.FieldTypeStringMap:
		if m, ok := field.Value.(map[string]string); ok {
			return s.ArrayValue.Render(ui.RenderStringMap(m, 96))
		}
	case parser.FieldTypeAuthorArray:
		if arr, ok := field.Value.([]parser.AuthorEntry); ok {
			parts := make([]string, 0, len(arr))
			for _, a := range arr {
				switch {
				case a.Name != "" && a.Email != "":
					parts = append(parts, fmt.Sprintf("{name = %q, email = %q}", a.Name, a.Email))
				case a.Name != "":
					parts = append(parts, fmt.Sprintf("{name = %q}", a.Name))
				case a.Email != "":
					parts = append(parts, fmt.Sprintf("{email = %q}", a.Email))
				}
			}
			return s.ArrayValue.Render("[" + strings.Join(parts, ", ") + "]")
		}
	case parser.FieldTypeStringMapOfArray:
		if m, ok := field.Value.(map[string][]string); ok {
			return s.ArrayValue.Render(fmt.Sprintf("{%d groups}", len(m)))
		}
	case parser.FieldTypeAny:
		return s.Value.Render(ui.RenderAny(field.Value, 0))
	case parser.FieldTypeBool, parser.FieldTypeInt:
		return s.Value.Render(fmt.Sprint(field.Value))
	}
	return s.Value.Render(ui.RenderAny(field.Value, 0))
}

// View renders the editor pane.
func (m EditorModel) View() string {
	s := theme.Styles()
	t := theme.Active()
	if s == nil || t == nil {
		return m.viewport.View()
	}

	vpHeight := m.height
	var inputView string

	if m.editMode == EditModeEditValue || m.editMode == EditModeAddKey || m.editMode == EditModeAddValue {
		vpHeight = m.height - 3
		if vpHeight < 1 {
			vpHeight = 1
		}

		var prompt, inputStr string
		switch m.editMode {
		case EditModeEditValue:
			prompt = s.Key.Render(fmt.Sprintf(" %s → ", m.currentKey()))
			inputStr = m.inputValue.View()
		case EditModeAddKey:
			prompt = s.Cursor.Render(" key → ")
			inputStr = m.inputKey.View()
		case EditModeAddValue:
			prompt = s.Key.Render(fmt.Sprintf(" %s → ", m.inputKey.Value()))
			inputStr = m.inputValue.View()
		}

		innerW := m.width - 4
		if innerW < 20 {
			innerW = 20
		}
		divider := s.Divider.Render(strings.Repeat("─", innerW))
		inputView = "\n" + divider + "\n" + s.Cursor.Render(prompt) + inputStr
	}

	m.viewport.Height = vpHeight
	m.viewport.Width = m.width - 4 // Account for Padding(0, 2) = 4 chars consumed

	// Auto-scroll viewport to keep cursor in view
	if len(m.fields) > 0 {
		if m.cursor < m.viewport.YOffset {
			m.viewport.YOffset = m.cursor
		} else if m.cursor >= m.viewport.YOffset+m.viewport.Height {
			m.viewport.YOffset = m.cursor - m.viewport.Height + 1
		}
	}

	body := m.viewport.View()
	content := lipgloss.NewStyle().
		Padding(0, 2).
		Width(m.width).
		Height(m.height).
		Render(body + inputView)

	if !m.focused {
		content = lipgloss.NewStyle().Faint(true).Render(content)
	}

	return content
}

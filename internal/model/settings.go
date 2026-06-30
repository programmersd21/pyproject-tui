// Package model contains Bubble Tea models for the pyproject-tui application.
package model

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/programmersd21/pyproject-tui/internal/settings"
	"github.com/programmersd21/pyproject-tui/internal/theme"
	"github.com/programmersd21/pyproject-tui/internal/ui"
)

// SettingsChangedMsg signals that settings were modified.
type SettingsChangedMsg struct {
	Config settings.Config
}

// SettingsModel implements the Settings page.
type SettingsModel struct {
	width   int
	height  int
	cursor  int
	options []settingsOption
	config  settings.Config // Own copy, not pointer
}

type settingsOption struct {
	category string
	label    string
	value    string
	choices  []string
	action   func(*SettingsModel)
}

// NewSettingsModel creates a new settings page model.
func NewSettingsModel(cfg settings.Config) SettingsModel {
	m := SettingsModel{
		config: cfg,
	}
	m.rebuildOptions()
	return m
}

// Config returns the current configuration.
func (m *SettingsModel) Config() settings.Config {
	return m.config
}

func (m *SettingsModel) rebuildOptions() {
	themeNames := theme.List()
	m.options = []settingsOption{
		{
			category: "Appearance",
			label:    "Theme",
			value:    m.config.Theme,
			choices:  themeNames,
			action:   (*SettingsModel).cycleTheme,
		},
		{
			category: "Appearance",
			label:    "Animations",
			value:    boolToString(m.config.AnimationsOn),
			choices:  []string{"On", "Off"},
			action:   (*SettingsModel).toggleAnimations,
		},
		{
			category: "Appearance",
			label:    "Line Numbers",
			value:    boolToString(m.config.ShowLineNumbers),
			choices:  []string{"On", "Off"},
			action:   (*SettingsModel).toggleLineNumbers,
		},
		{
			category: "Layout",
			label:    "UI Density",
			value:    m.config.UIDensity,
			choices:  []string{"compact", "normal", "comfortable"},
			action:   (*SettingsModel).cycleUIDensity,
		},
		{
			category: "Layout",
			label:    "Border Style",
			value:    m.config.BorderStyle,
			choices:  []string{"rounded", "normal", "thick", "double"},
			action:   (*SettingsModel).cycleBorderStyle,
		},
	}
}

func boolToString(b bool) string {
	if b {
		return "On"
	}
	return "Off"
}

func (m *SettingsModel) cycleTheme() {
	current := m.config.Theme
	themeNames := theme.List()
	for i, name := range themeNames {
		if name == current {
			m.config.Theme = themeNames[(i+1)%len(themeNames)]
			m.rebuildOptions()
			return
		}
	}
	if len(themeNames) > 0 {
		m.config.Theme = themeNames[0]
		m.rebuildOptions()
	}
}

func (m *SettingsModel) toggleAnimations() {
	m.config.AnimationsOn = !m.config.AnimationsOn
	m.rebuildOptions()
}

func (m *SettingsModel) toggleLineNumbers() {
	m.config.ShowLineNumbers = !m.config.ShowLineNumbers
	m.rebuildOptions()
}

func (m *SettingsModel) cycleUIDensity() {
	densities := []string{"compact", "normal", "comfortable"}
	for i, d := range densities {
		if d == m.config.UIDensity {
			m.config.UIDensity = densities[(i+1)%len(densities)]
			m.rebuildOptions()
			return
		}
	}
}

func (m *SettingsModel) cycleBorderStyle() {
	styles := []string{"rounded", "normal", "thick", "double"}
	for i, s := range styles {
		if s == m.config.BorderStyle {
			m.config.BorderStyle = styles[(i+1)%len(styles)]
			m.rebuildOptions()
			return
		}
	}
}

// Update handles settings page input.
func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.options) - 1
			}
		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.options) {
				m.cursor = 0
			}
		case "enter", "right", "l", " ":
			if m.cursor >= 0 && m.cursor < len(m.options) {
				m.options[m.cursor].action(&m)
				// Notify parent that settings changed
				return m, func() tea.Msg {
					return SettingsChangedMsg{Config: m.config}
				}
			}
		}
	}
	return m, nil
}

// View renders the settings page.
func (m SettingsModel) View() string {
	s := theme.Styles()
	if s == nil {
		return "Loading settings..."
	}
	t := theme.Active()

	var inner strings.Builder

	// Render options grouped by category
	currentCategory := ""
	for i, opt := range m.options {
		if opt.category != currentCategory {
			if currentCategory != "" {
				inner.WriteString("\n")
			}
			categoryStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(t.Accent)).
				Bold(true).
				MarginTop(1).
				MarginBottom(1)
			inner.WriteString(categoryStyle.Render("  " + categoryIcon(opt.category) + " " + opt.category))
			inner.WriteString("\n")
			currentCategory = opt.category
		}

		prefix := "    "
		labelStyle := s.Key
		valueStyle := s.Value

		if i == m.cursor {
			prefix = "  > "
			labelStyle = s.Cursor
			valueStyle = s.Cursor
		}

		label := labelStyle.Render(fmt.Sprintf("%-20s", opt.label))

		var valueDisplay string
		if opt.label == "Theme" {
			valueDisplay = m.renderThemePreview(opt.value, i == m.cursor)
		} else {
			valueDisplay = valueStyle.Render(opt.value)
		}

		fmt.Fprintf(&inner, "%s%s  %s\n", prefix, label, valueDisplay)
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.BorderFocused)).
		Padding(1, 3)

	boxW := m.width - 8
	if boxW > 60 {
		boxW = 60
	}
	boxStyle = boxStyle.Width(boxW)

	title := s.TitleBar.Render("  Settings")
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.TextDim))
	hint := hintStyle.Render("  [up/down] navigate   [enter] change   [esc] back   [s] save")

	content := title + "\n\n" + boxStyle.Render(inner.String()) + "\n" + hint
	return ui.CenterBox(content, m.width, ui.BodyHeight(m.height))
}

func categoryIcon(category string) string {
	switch category {
	case "Appearance":
		return "\u25C6"
	case "Layout":
		return "\u25CB"
	default:
		return "\u25AA"
	}
}

func (m SettingsModel) renderThemePreview(themeName string, selected bool) string {
	t, ok := theme.Get(themeName)
	if !ok {
		return themeName
	}

	// Render theme name
	nameStyle := theme.Styles().Value
	if selected {
		nameStyle = theme.Styles().Cursor
	}

	var sb strings.Builder
	sb.WriteString(nameStyle.Render(fmt.Sprintf("%-15s", t.DisplayName)))
	sb.WriteString("  ")

	// Render color preview circles
	for _, color := range t.PreviewColors {
		circle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Render("\u25CF")
		sb.WriteString(circle)
		sb.WriteString(" ")
	}

	return sb.String()
}

// SetSize updates the settings model dimensions.
func (m *SettingsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

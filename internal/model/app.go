// Package model contains Bubble Tea models for the pyproject-tui application.
package model

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
	"github.com/programmersd21/pyproject-tui/internal/parser"
	"github.com/programmersd21/pyproject-tui/internal/settings"
	"github.com/programmersd21/pyproject-tui/internal/theme"
	"github.com/programmersd21/pyproject-tui/internal/ui"
)

// Pane identifies the active focus area.
type Pane int

const (
	// PaneSidebar focuses the section list.
	PaneSidebar Pane = iota
	// PaneEditor focuses the field editor.
	PaneEditor
)

// ViewMode identifies the current view.
type ViewMode int

const (
	// ViewDashboard shows the main menu.
	ViewDashboard ViewMode = iota
	// ViewEditor shows the two-pane editor.
	ViewEditor
	// ViewSettings shows the settings page.
	ViewSettings
	// ViewHelp shows the help overlay.
	ViewHelp
)

// AppModel is the root Bubble Tea model.
type AppModel struct {
	pyproject    *parser.PyProject
	filePath     string
	dirty        bool
	width        int
	height       int
	sidebar      SidebarModel
	editor       EditorModel
	settingsView SettingsModel
	focus        Pane
	undoStack    []*parser.PyProject
	redoStack    []*parser.PyProject
	viewMode     ViewMode
	statusMsg    string
	statusIsErr  bool
	keys         KeyMap
	readOnly     bool
	settings     settings.Config
	settingsPath string
	version      string

	// Confirmation and prompts
	confirmActive bool
	confirmMsg    string
	confirmAction PendingAction
	promptingTool bool
	toolNameInput textinput.Model
	promptingOpen bool
	openPathInput textinput.Model

	// Dashboard animation
	dashboardCursor int
	sandField       *sandField
}

// NewAppModel constructs the application model.
func NewAppModel(p *parser.PyProject, readOnly bool) AppModel {
	cfg, cfgPath, _ := settings.LoadDefault()
	theme.Initialize()
	theme.SetDisplayConfig(theme.DisplayConfig{
		Density:         cfg.UIDensity,
		BorderStyle:     cfg.BorderStyle,
		Animations:      cfg.AnimationsOn,
		ShowLineNumbers: cfg.ShowLineNumbers,
	})
	ui.ApplyTheme(cfg.Theme)

	sidebar := NewSidebar(p)
	editor := NewEditor()
	if p != nil {
		section := parser.SectionID{Kind: parser.KindProject}
		if p.Project == nil && p.BuildSystem != nil {
			section = parser.SectionID{Kind: parser.KindBuildSystem}
		}
		editor.SetSection(section, p)
		sidebar.MoveToSection(section)
	}

	ti := textinput.New()
	ti.Placeholder = "e.g. poetry, ruff, black"
	ti.Prompt = " > "
	oi := textinput.New()
	oi.Placeholder = "/path/to/pyproject.toml"
	oi.Prompt = " > "

	m := AppModel{
		pyproject:     p,
		filePath:      p.FilePath,
		sidebar:       sidebar,
		editor:        editor,
		settingsView:  NewSettingsModel(cfg),
		focus:         PaneSidebar,
		keys:          DefaultKeyMap(),
		readOnly:      readOnly,
		settings:      cfg,
		settingsPath:  cfgPath,
		toolNameInput: ti,
		openPathInput: oi,
		viewMode:      ViewDashboard,
		sandField:     newSandField(80, 24),
	}
	m.sidebar.SetFocused(true)
	m.editor.SetFocused(false)
	return m
}

// SetStatus updates the status line.
func (m *AppModel) SetStatus(msg string, isErr bool) {
	m.statusMsg = msg
	m.statusIsErr = isErr
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*33, func(_ time.Time) tea.Msg {
		return TickMsg{}
	})
}

// SetVersion sets the application version for display in the footer.
func (m *AppModel) SetVersion(v string) {
	m.version = v
}

// Init implements tea.Model.
func (m AppModel) Init() tea.Cmd {
	return tickCmd()
}

// Update implements tea.Model.
//
//nolint:gocyclo // Bubble Tea Update functions are inherently complex state machines
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 1. Handle window sizing first
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = sizeMsg.Width
		m.height = sizeMsg.Height
		m.sidebar.SetSize(ui.SidebarWidth(sizeMsg.Width), ui.BodyHeight(sizeMsg.Height))
		m.editor.width = ui.EditorWidth(sizeMsg.Width)
		m.editor.height = ui.BodyHeight(sizeMsg.Height)
		m.editor.SetShowLineNumbers(m.settings.ShowLineNumbers)
		m.sidebar.SetFocused(m.focus == PaneSidebar)
		m.settingsView.SetSize(sizeMsg.Width, ui.BodyHeight(sizeMsg.Height))
		m.sandField.setSize(sizeMsg.Width, ui.BodyHeight(sizeMsg.Height))
		return m, nil
	}

	// 2. Handle animation ticks
	if _, ok := msg.(TickMsg); ok {
		return m, tickCmd()
	}

	// 2b. Handle open file prompt
	if m.promptingOpen {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "esc":
				m.promptingOpen = false
				return m, nil
			case "enter":
				path := strings.TrimSpace(m.openPathInput.Value())
				if path != "" {
					m.promptingOpen = false
					return m, func() tea.Msg {
						return OpenFileMsg{Path: path}
					}
				}
				m.promptingOpen = false
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.openPathInput, cmd = m.openPathInput.Update(msg)
		return m, cmd
	}

	// 3. Handle Settings view
	if m.viewMode == ViewSettings {
		if key, ok := msg.(tea.KeyMsg); ok && (key.String() == "esc" || key.String() == "s") {
			m.applySettings(m.settingsView.Config())
			if key.String() == "esc" {
				m.viewMode = ViewDashboard
			}
			return m, nil
		}
		var cmd tea.Cmd
		m.settingsView, cmd = m.settingsView.Update(msg)
		return m, cmd
	}

	// 4. Handle Help view
	if m.viewMode == ViewHelp {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "?", "esc":
				m.viewMode = ViewDashboard
			}
		}
		return m, nil
	}

	// 5. Handle active tool-name prompt input
	if m.promptingTool {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "esc":
				m.promptingTool = false
				return m, nil
			case "enter":
				name := strings.TrimSpace(m.toolNameInput.Value())
				if name != "" {
					m.pushUndo()
					if m.pyproject.Tool == nil {
						m.pyproject.Tool = map[string]any{}
					}
					m.pyproject.Tool[name] = map[string]any{}
					m.sidebar.UpdateSections(m.pyproject)
					newSec := parser.SectionID{Kind: parser.KindTool, ToolKey: name}
					m.sidebar.MoveToSection(newSec)
					m.editor.SetSection(newSec, m.pyproject)
					m.dirty = true
					m.statusMsg = fmt.Sprintf("Added [tool.%s]", name)
					m.statusIsErr = false
				}
				m.promptingTool = false
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.toolNameInput, cmd = m.toolNameInput.Update(msg)
		return m, cmd
	}

	// 6. Handle active confirmations
	if m.confirmActive {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "y", "enter":
				m.confirmActive = false
				switch m.confirmAction {
				case ActionQuitDirty:
					return m, tea.Quit
				case ActionDelete:
					m.pushUndo()
					m.applyDelete(m.editor.currentKey())
					m.editor.SetSection(m.editor.section, m.pyproject)
					m.dirty = true
					m.statusMsg = "Deleted field."
					m.statusIsErr = false
				case ActionDeleteSection:
					selected := m.sidebar.SelectedSection()
					if selected.Kind == parser.KindTool {
						m.pushUndo()
						delete(m.pyproject.Tool, selected.ToolKey)
						m.sidebar.UpdateSections(m.pyproject)
						firstSec := m.sidebar.SelectedSection()
						m.editor.SetSection(firstSec, m.pyproject)
						m.dirty = true
						m.statusMsg = fmt.Sprintf("Deleted [tool.%s]", selected.ToolKey)
						m.statusIsErr = false
					}
				}
				return m, nil
			case "n", "esc":
				m.confirmActive = false
				m.statusMsg = "Canceled"
				m.statusIsErr = false
				return m, nil
			}
		}
		return m, nil
	}

	// 7. Handle Dashboard view
	if m.viewMode == ViewDashboard {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "up", "k":
				m.dashboardCursor--
				if m.dashboardCursor < 0 {
					m.dashboardCursor = len(dashboardOptions) - 1
				}
				return m, nil
			case "down", "j":
				m.dashboardCursor++
				if m.dashboardCursor >= len(dashboardOptions) {
					m.dashboardCursor = 0
				}
				return m, nil
			case "enter":
				switch m.dashboardCursor {
				case 0: // Open Editor
					m.viewMode = ViewEditor
				case 1: // Open File
					m.promptingOpen = true
					m.openPathInput.SetValue("")
					m.openPathInput.Focus()
				case 2: // Open Config
					if m.settingsPath != "" {
						return m, func() tea.Msg {
							return NativeOpenFileMsg{Path: m.settingsPath}
						}
					}
					m.statusMsg = "Config file not found"
					m.statusIsErr = true
				case 3: // Settings
					m.viewMode = ViewSettings
				case 4: // Undo History
					m.statusMsg = fmt.Sprintf("Undo: %d edits  |  Redo: %d edits", len(m.undoStack), len(m.redoStack))
					m.statusIsErr = false
				case 5: // Help
					m.viewMode = ViewHelp
				case 6: // Quit
					if m.dirty {
						m.confirmActive = true
						m.confirmAction = ActionQuitDirty
						m.confirmMsg = "You have unsaved changes. Quit anyway? [y/N]"
					} else {
						return m, tea.Quit
					}
				}
				return m, nil
			case "e":
				m.viewMode = ViewEditor
				return m, nil
			case "o":
				m.promptingOpen = true
				m.openPathInput.SetValue("")
				m.openPathInput.Focus()
				return m, nil
			case "i":
				if m.settingsPath != "" {
					return m, func() tea.Msg {
						return NativeOpenFileMsg{Path: m.settingsPath}
					}
				}
				m.statusMsg = "Config file not found"
				m.statusIsErr = true
				return m, nil
			case "t":
				m.cycleTheme()
				return m, nil
			case "c":
				m.viewMode = ViewSettings
				return m, nil
			case "u":
				m.statusMsg = fmt.Sprintf("Undo: %d edits  |  Redo: %d edits", len(m.undoStack), len(m.redoStack))
				m.statusIsErr = false
				return m, nil
			case "h", "?":
				m.viewMode = ViewHelp
				return m, nil
			case "q", "ctrl+c":
				if m.dirty {
					m.confirmActive = true
					m.confirmAction = ActionQuitDirty
					m.confirmMsg = "You have unsaved changes. Quit anyway? [y/N]"
				} else {
					return m, tea.Quit
				}
				return m, nil
			case "esc":
				// Already on dashboard, do nothing
				return m, nil
			}
		}
		// Non-key messages (e.g. OpenFileMsg) fall through to section 9
	}

	// 8. Handle field inline text input inside the Editor
	if m.editor.editMode != EditModeNormal {
		var cmd tea.Cmd
		m.editor, cmd = m.editor.Update(msg)
		return m, cmd
	}

	// 9. Handle background message updates
	switch msg := msg.(type) {
	case OpenFileMsg:
		return m.handleOpenFile(msg.Path)
	case NativeOpenFileMsg:
		return m, openFileNative(msg.Path)
	case NativeOpenFileResultMsg:
		if msg.Err != nil {
			m.statusMsg = "Error opening file: " + msg.Err.Error()
			m.statusIsErr = true
		} else {
			m.statusMsg = "Opened config file in default editor."
			m.statusIsErr = false
		}
		return m, nil
	case SettingsChangedMsg:
		m.applySettings(msg.Config)
		return m, nil
	case FieldEditedMsg:
		m.pushUndo()
		m.applyEdit(msg.Key, parseInputValue(msg.NewValue.(string)))
		m.editor.SetSection(m.editor.section, m.pyproject)
		m.dirty = true
		m.statusMsg = "Modified field."
		m.statusIsErr = false
		return m, nil
	case FieldAddedMsg:
		m.pushUndo()
		m.applyAdd(msg.Key, parseInputValue(msg.Value.(string)))
		m.editor.SetSection(m.editor.section, m.pyproject)
		m.dirty = true
		m.statusMsg = "Added field."
		m.statusIsErr = false
		return m, nil
	}

	// 10. Handle normal keyboard shortcuts in Editor view
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "tab":
			if m.focus == PaneSidebar {
				m.focus = PaneEditor
			} else {
				m.focus = PaneSidebar
			}
			m.sidebar.SetFocused(m.focus == PaneSidebar)
			m.editor.SetFocused(m.focus == PaneEditor)
		case "shift+tab":
			if m.focus == PaneSidebar {
				m.focus = PaneEditor
			} else {
				m.focus = PaneSidebar
			}
			m.sidebar.SetFocused(m.focus == PaneSidebar)
			m.editor.SetFocused(m.focus == PaneEditor)
		case "left", "h":
			if m.focus == PaneEditor {
				m.focus = PaneSidebar
				m.sidebar.SetFocused(true)
				m.editor.SetFocused(false)
			}
		case "right", "l":
			if m.focus == PaneSidebar {
				m.focus = PaneEditor
				m.sidebar.SetFocused(false)
				m.editor.SetFocused(true)
			}
		case "?":
			m.viewMode = ViewHelp
		case "q", "ctrl+c":
			if m.dirty {
				m.confirmActive = true
				m.confirmAction = ActionQuitDirty
				m.confirmMsg = "You have unsaved changes. Quit anyway? [y/N]"
				return m, nil
			}
			return m, tea.Quit
		case "s":
			if m.readOnly {
				m.statusMsg = "Read-only: unable to save."
				m.statusIsErr = true
				break
			}
			err := parser.Write(m.pyproject)
			if err != nil {
				m.statusMsg = "Error: " + err.Error()
				m.statusIsErr = true
			} else {
				m.dirty = false
				m.statusMsg = "Saved successfully."
				m.statusIsErr = false
			}
		case "u":
			if len(m.undoStack) > 0 {
				m.redoStack = append(m.redoStack, parser.Clone(m.pyproject))
				last := m.undoStack[len(m.undoStack)-1]
				m.undoStack = m.undoStack[:len(m.undoStack)-1]
				m.pyproject = parser.Clone(last)
				m.sidebar.UpdateSections(m.pyproject)
				m.editor.SetSection(m.editor.section, m.pyproject)
				m.dirty = len(m.undoStack) > 0
				m.statusMsg = "Undo applied."
				m.statusIsErr = false
			} else {
				m.statusMsg = "No undo history."
				m.statusIsErr = false
			}
		case "r":
			if len(m.redoStack) > 0 {
				m.undoStack = append(m.undoStack, parser.Clone(m.pyproject))
				next := m.redoStack[len(m.redoStack)-1]
				m.redoStack = m.redoStack[:len(m.redoStack)-1]
				m.pyproject = parser.Clone(next)
				m.sidebar.UpdateSections(m.pyproject)
				m.editor.SetSection(m.editor.section, m.pyproject)
				m.dirty = true
				m.statusMsg = "Redo applied."
				m.statusIsErr = false
			} else {
				m.statusMsg = "No redo history."
				m.statusIsErr = false
			}
		case "t":
			m.cycleTheme()
		case "c":
			m.viewMode = ViewSettings
		case "esc":
			m.viewMode = ViewDashboard
			m.statusMsg = "Returned to dashboard"
			m.statusIsErr = false
		case "enter", "e":
			if m.focus == PaneSidebar {
				m.focus = PaneEditor
				m.sidebar.SetFocused(false)
				m.editor.SetFocused(true)
				selected := m.sidebar.SelectedSection()
				m.editor.SetSection(selected, m.pyproject)
			} else {
				m.editor, _ = m.editor.Update(key)
			}
		case "a":
			if m.focus == PaneEditor {
				m.editor, _ = m.editor.Update(key)
			} else {
				m.promptingTool = true
				m.toolNameInput.SetValue("")
				m.toolNameInput.Focus()
			}
		case "d":
			if m.focus == PaneEditor {
				if len(m.editor.fields) > 0 {
					m.confirmActive = true
					m.confirmAction = ActionDelete
					m.confirmMsg = fmt.Sprintf("Delete key %q? [y/N]", m.editor.currentKey())
				}
			} else {
				selected := m.sidebar.SelectedSection()
				if selected.Kind == parser.KindTool {
					m.confirmActive = true
					m.confirmAction = ActionDeleteSection
					m.confirmMsg = fmt.Sprintf("Delete [tool.%s] section and all its keys? [y/N]", selected.ToolKey)
				} else {
					m.statusMsg = "Cannot delete standard sections."
					m.statusIsErr = true
				}
			}
		case "up", "down", "j", "k":
			if m.focus == PaneEditor {
				m.editor, _ = m.editor.Update(key)
			} else {
				m.sidebar, _ = m.sidebar.Update(key)
				// Live-preview sections on selection move
				selected := m.sidebar.SelectedSection()
				m.editor.SetSection(selected, m.pyproject)
			}
		}
	}

	return m, nil
}

// View renders the application.
func (m AppModel) View() string {
	s := theme.Styles()
	if s == nil {
		// Fallback if theme system isn't initialized
		return "Initializing theme system..."
	}

	title := fmt.Sprintf("  pyproject-tui  -  %s", m.filePath)
	if m.dirty {
		title += " *"
	}
	activeTheme := theme.Active()
	if activeTheme != nil {
		title += fmt.Sprintf("   [%s]", activeTheme.DisplayName)
	}
	title = ui.TruncateString(title, m.width)
	top := s.TitleBar.Width(m.width).Render(title)

	status := m.statusMsg
	if status == "" {
		status = "Ready"
	}
	if m.statusIsErr {
		status = s.Error.Render(status)
	} else if m.dirty {
		status = s.Success.Render(status)
	}

	keyHints := "[esc] dashboard  [c] settings  [?] help  [q] quit"
	switch m.viewMode {
	case ViewEditor:
		focusIndicator := ""
		if m.focus == PaneSidebar {
			focusIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color(activeTheme.Accent)).Bold(true).Render(" Sidebar ") + "  "
		} else {
			focusIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color(activeTheme.Accent)).Bold(true).Render(" Editor ") + "  "
		}
		keyHints = focusIndicator + "[↑/↓] move  [tab] switch  [e] edit  [a] add  [d] delete  [s] save  [u] undo  [r] redo  [t] theme  [?] help  [esc] dashboard  [q] quit"
	case ViewSettings:
		keyHints = "[↑/↓] navigate  [enter/space/→] change  [s] save  [esc] back"
	case ViewHelp:
		keyHints = "[?/esc] close help"
	case ViewDashboard:
		keyHints = "[↑/↓/j/k] navigate  [enter] select  [e] editor  [o] open file  [i] config  [c] settings  [?] help  [q] quit"
	}
	if m.promptingOpen {
		keyHints = "[enter] open  [esc] cancel"
	}

	versionStr := ""
	if m.version != "" {
		versionStr = " " + s.Success.Render("v"+m.version)
	}
	bottom := s.StatusBar.Width(m.width).Render(status + "    " + keyHints + versionStr)
	horizontalLine := s.Divider.Render(strings.Repeat("─", m.width))

	// Route to appropriate view - overlays first, then view modes
	var body string
	switch {
	case m.promptingOpen:
		promptContent := " Open pyproject.toml \n\n Enter a file path:\n " + m.openPathInput.View()
		body = ui.CenterBox(s.InputBox.Width(50).Render(promptContent), m.width, ui.BodyHeight(m.height))
	case m.promptingTool:
		promptContent := " Add New [tool.*] Section \n\n Enter tool name (e.g. poetry, ruff):\n " + m.toolNameInput.View()
		body = ui.CenterBox(s.InputBox.Width(50).Render(promptContent), m.width, ui.BodyHeight(m.height))
	case m.confirmActive:
		body = ui.CenterBox(s.ConfirmBox.Render(m.confirmMsg), m.width, ui.BodyHeight(m.height))
	case m.viewMode == ViewDashboard:
		body = m.dashboardView()
	case m.viewMode == ViewSettings:
		body = ui.CenterBox(m.settingsView.View(), m.width, ui.BodyHeight(m.height))
	case m.viewMode == ViewHelp:
		body = ui.CenterBox(s.HelpBox.Render(helpContent()), m.width, ui.BodyHeight(m.height))
	default:
		dividerHeight := ui.BodyHeight(m.height)
		dividerColor := activeTheme.Border
		if m.focus == PaneEditor {
			dividerColor = activeTheme.BorderFocused
		}
		verticalDivider := lipgloss.NewStyle().Width(1).Foreground(lipgloss.Color(dividerColor)).Render(strings.Repeat("│\n", dividerHeight))
		body = lipgloss.JoinHorizontal(lipgloss.Top, m.sidebar.View(), verticalDivider, m.editor.View())
	}

	return top + "\n" + horizontalLine + "\n" + body + "\n" + horizontalLine + "\n" + bottom
}

func helpContent() string {
	s := theme.Styles()
	if s == nil {
		return "Help system loading..."
	}

	sections := []struct {
		title string
		items []struct {
			keys string
			desc string
		}
	}{
		{
			title: "Navigation",
			items: []struct {
				keys string
				desc string
			}{
				{"Tab / Shift+Tab", "Switch focus between sidebar and editor"},
				{"h / ←", "Move focus to sidebar"},
				{"l / →", "Move focus to editor"},
				{"j / ↓", "Move down"},
				{"k / ↑", "Move up"},
				{"Esc", "Return to dashboard"},
			},
		},
		{
			title: "Editing",
			items: []struct {
				keys string
				desc string
			}{
				{"Enter / e", "Edit selected field or open section"},
				{"a", "Add new field (editor) or tool section (sidebar)"},
				{"d", "Delete field or section"},
				{"Esc", "Cancel editing"},
			},
		},
		{
			title: "File & History",
			items: []struct {
				keys string
				desc string
			}{
				{"o", "Open a pyproject.toml file"},
				{"i", "Open pyproject-tui config directly"},
				{"s", "Save changes to disk"},
				{"u", "Undo last change (up to 50 edits)"},
				{"r", "Redo last undone change"},
			},
		},
		{
			title: "General",
			items: []struct {
				keys string
				desc string
			}{
				{"e", "Open editor"},
				{"c", "Open settings"},
				{"t", "Cycle through themes"},
				{"? / h", "Toggle help"},
				{"q / Ctrl+C", "Quit application"},
			},
		},
	}

	var lines []string
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Active().Accent)).
		Bold(true)

	keyStyle := s.Key
	descStyle := s.Value

	lines = append(lines, titleStyle.Render("Keyboard Shortcuts"))

	for _, section := range sections {
		sectionTitleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Active().AccentSoft)).
			Bold(true).
			MarginTop(1)
		lines = append(lines, sectionTitleStyle.Render(section.title))

		for _, item := range section.items {
			key := keyStyle.Render(fmt.Sprintf("  %-20s", item.keys))
			desc := descStyle.Render(item.desc)
			lines = append(lines, fmt.Sprintf("%s  %s", key, desc))
		}
	}

	return strings.Join(lines, "\n")
}

func (m *AppModel) cycleTheme() {
	m.applySettings(settings.Config{
		Theme:           ui.NextTheme(),
		AnimationsOn:    m.settings.AnimationsOn,
		UIDensity:       m.settings.UIDensity,
		BorderStyle:     m.settings.BorderStyle,
		ShowLineNumbers: m.settings.ShowLineNumbers,
	})
	m.sandField.refreshColors()
	m.statusMsg = "Theme: " + m.settings.Theme
	m.statusIsErr = false
}

func (m *AppModel) applySettings(cfg settings.Config) {
	m.settings = cfg
	m.settingsView = NewSettingsModel(cfg)
	ui.ApplyTheme(cfg.Theme)
	theme.SetDisplayConfig(theme.DisplayConfig{
		Density:         cfg.UIDensity,
		BorderStyle:     cfg.BorderStyle,
		Animations:      cfg.AnimationsOn,
		ShowLineNumbers: cfg.ShowLineNumbers,
	})
	m.editor.SetShowLineNumbers(cfg.ShowLineNumbers)
	if err := settings.Save(m.settingsPath, cfg); err != nil {
		m.statusMsg = "Error saving settings: " + err.Error()
		m.statusIsErr = true
		return
	}
	m.statusMsg = "Settings saved."
	m.statusIsErr = false
}

func (m *AppModel) handleOpenFile(path string) (tea.Model, tea.Cmd) {
	loaded, err := parser.Load(path)
	if err != nil {
		if parser.IsSyntaxError(err) {
			loaded = parser.NewEmpty(path)
			if raw, rawErr := parser.LoadRaw(path); rawErr == nil {
				loaded.Raw = raw
			}
			m.pyproject = loaded
			m.filePath = path
			m.sidebar = NewSidebar(loaded)
			m.editor.SetSection(parser.SectionID{Kind: parser.KindProject}, loaded)
			m.statusMsg = err.Error()
			m.statusIsErr = true
			m.viewMode = ViewEditor
			return m, nil
		}
		m.statusMsg = "Error opening file: " + err.Error()
		m.statusIsErr = true
		return m, nil
	}
	m.pyproject = loaded
	m.filePath = path
	m.sidebar = NewSidebar(loaded)
	section := parser.SectionID{Kind: parser.KindProject}
	if loaded.Project == nil && loaded.BuildSystem != nil {
		section = parser.SectionID{Kind: parser.KindBuildSystem}
	}
	m.sidebar.MoveToSection(section)
	m.editor.SetSection(section, loaded)
	m.viewMode = ViewEditor
	m.dirty = false
	m.statusMsg = "Opened " + path
	m.statusIsErr = false
	return m, nil
}

func (m *AppModel) pushUndo() {
	cp := parser.Clone(m.pyproject)
	if cp == nil {
		return
	}
	m.undoStack = append(m.undoStack, cp)
	if len(m.undoStack) > 50 {
		m.undoStack = m.undoStack[len(m.undoStack)-50:]
	}
	m.redoStack = nil
}

// type helpers for parsing input string to typed data structure
func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toStringSlice(v any) []string {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case []string:
		return val
	case []any:
		out := make([]string, 0, len(val))
		for _, x := range val {
			out = append(out, toString(x))
		}
		return out
	case string:
		parts := strings.Split(val, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
		return out
	default:
		return []string{toString(v)}
	}
}

func toStringMap(v any) map[string]string {
	if v == nil {
		return nil
	}
	out := map[string]string{}
	switch val := v.(type) {
	case map[string]string:
		return val
	case map[string]any:
		for k, x := range val {
			out[k] = toString(x)
		}
	}
	return out
}

func toAuthorArray(v any) []parser.AuthorEntry {
	if v == nil {
		return nil
	}
	if arr, ok := v.([]parser.AuthorEntry); ok {
		return arr
	}
	items, ok := v.([]any)
	if !ok {
		if m, ok := v.(map[string]any); ok {
			return []parser.AuthorEntry{parseAuthorMap(m)}
		}
		return nil
	}
	out := make([]parser.AuthorEntry, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			out = append(out, parseAuthorMap(m))
		}
	}
	return out
}

func parseAuthorMap(m map[string]any) parser.AuthorEntry {
	a := parser.AuthorEntry{}
	if name, ok := m["name"].(string); ok {
		a.Name = name
	}
	if email, ok := m["email"].(string); ok {
		a.Email = email
	}
	return a
}

func toMapOfStringSlice(v any) map[string][]string {
	if v == nil {
		return nil
	}
	if m, ok := v.(map[string][]string); ok {
		return m
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	out := map[string][]string{}
	for k, val := range m {
		out[k] = toStringSlice(val)
	}
	return out
}

func toMapOfStringMap(v any) map[string]map[string]string {
	if v == nil {
		return nil
	}
	if m, ok := v.(map[string]map[string]string); ok {
		return m
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	out := map[string]map[string]string{}
	for k, val := range m {
		out[k] = toStringMap(val)
	}
	return out
}

func parseInputValue(input string) any {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	// Attempt parsing TOML representation (e.g. array, integer, boolean)
	tomlStr := "val = " + input
	var m map[string]any
	if err := toml.Unmarshal([]byte(tomlStr), &m); err == nil {
		if val, exists := m["val"]; exists {
			return val
		}
	}

	// Fallback to literal string representation
	return input
}

func (m *AppModel) applyEdit(key string, value any) {
	if m.pyproject == nil {
		return
	}
	switch m.editor.section.Kind {
	case parser.KindProject:
		if m.pyproject.Project == nil {
			m.pyproject.Project = &parser.ProjectSection{}
		}
		switch key {
		case "name":
			m.pyproject.Project.Name = toString(value)
		case "version":
			m.pyproject.Project.Version = toString(value)
		case "description":
			m.pyproject.Project.Description = toString(value)
		case "readme":
			m.pyproject.Project.Readme = value
		case "requires-python":
			m.pyproject.Project.RequiresPython = toString(value)
		case "license":
			m.pyproject.Project.License = value
		case "authors":
			m.pyproject.Project.Authors = toAuthorArray(value)
		case "maintainers":
			m.pyproject.Project.Maintainers = toAuthorArray(value)
		case "keywords":
			m.pyproject.Project.Keywords = toStringSlice(value)
		case "classifiers":
			m.pyproject.Project.Classifiers = toStringSlice(value)
		case "urls":
			m.pyproject.Project.URLs = toStringMap(value)
		case "dependencies":
			m.pyproject.Project.Dependencies = toStringSlice(value)
		case "optional-dependencies":
			m.pyproject.Project.OptionalDependencies = toMapOfStringSlice(value)
		case "scripts":
			m.pyproject.Project.Scripts = toStringMap(value)
		case "entry-points":
			m.pyproject.Project.EntryPoints = toMapOfStringMap(value)
		case "dynamic":
			m.pyproject.Project.Dynamic = toStringSlice(value)
		}
	case parser.KindBuildSystem:
		if m.pyproject.BuildSystem == nil {
			m.pyproject.BuildSystem = &parser.BuildSystemSection{}
		}
		switch key {
		case "requires":
			m.pyproject.BuildSystem.Requires = toStringSlice(value)
		case "build-backend":
			m.pyproject.BuildSystem.BuildBackend = toString(value)
		case "backend-path":
			m.pyproject.BuildSystem.BackendPath = toStringSlice(value)
		}
	case parser.KindTool:
		if m.pyproject.Tool == nil {
			m.pyproject.Tool = map[string]any{}
		}
		root, _ := m.pyproject.Tool[m.editor.section.ToolKey].(map[string]any)
		if root == nil {
			root = map[string]any{}
		}
		root[key] = value
		m.pyproject.Tool[m.editor.section.ToolKey] = root
	}
}

func (m *AppModel) applyAdd(key string, value any) {
	m.applyEdit(key, value)
}

func (m *AppModel) applyDelete(key string) {
	if m.pyproject == nil {
		return
	}
	switch m.editor.section.Kind {
	case parser.KindProject:
		if m.pyproject.Project == nil {
			return
		}
		switch key {
		case "name":
			m.pyproject.Project.Name = ""
		case "version":
			m.pyproject.Project.Version = ""
		case "description":
			m.pyproject.Project.Description = ""
		case "readme":
			m.pyproject.Project.Readme = nil
		case "requires-python":
			m.pyproject.Project.RequiresPython = ""
		case "license":
			m.pyproject.Project.License = nil
		case "authors":
			m.pyproject.Project.Authors = nil
		case "maintainers":
			m.pyproject.Project.Maintainers = nil
		case "keywords":
			m.pyproject.Project.Keywords = nil
		case "classifiers":
			m.pyproject.Project.Classifiers = nil
		case "urls":
			m.pyproject.Project.URLs = nil
		case "dependencies":
			m.pyproject.Project.Dependencies = nil
		case "optional-dependencies":
			m.pyproject.Project.OptionalDependencies = nil
		case "scripts":
			m.pyproject.Project.Scripts = nil
		case "entry-points":
			m.pyproject.Project.EntryPoints = nil
		case "dynamic":
			m.pyproject.Project.Dynamic = nil
		}
	case parser.KindBuildSystem:
		if m.pyproject.BuildSystem == nil {
			return
		}
		switch key {
		case "requires":
			m.pyproject.BuildSystem.Requires = nil
		case "build-backend":
			m.pyproject.BuildSystem.BuildBackend = ""
		case "backend-path":
			m.pyproject.BuildSystem.BackendPath = nil
		}
	case parser.KindTool:
		if m.pyproject.Tool == nil {
			return
		}
		if root, ok := m.pyproject.Tool[m.editor.section.ToolKey].(map[string]any); ok {
			delete(root, key)
			m.pyproject.Tool[m.editor.section.ToolKey] = root
		}
	}
}

type dashboardOption struct {
	key   string
	label string
	desc  string
}

var dashboardOptions = []dashboardOption{
	{"e", "Open Editor", "Edit project metadata & dependencies"},
	{"o", "Open File", "Load a pyproject.toml from disk"},
	{"i", "Open Config", "Edit pyproject-tui settings file"},
	{"c", "Settings", "Configure themes and preferences"},
	{"u", "Undo History", "View and recover previous edits"},
	{"h", "Help Guide", "View keyboard shortcuts and usage"},
	{"q", "Quit", "Exit pyproject-tui"},
}

func (m AppModel) dashboardView() string {
	s := theme.Styles()
	if s == nil {
		return "Loading..."
	}

	logoLines := []string{
		`██████╗ ██╗   ██╗██████╗ ██████╗  ██████╗      ██╗███████╗ ██████╗████████╗   ████████╗██╗   ██╗██╗`,
		`██╔══██╗╚██╗ ██╔╝██╔══██╗██╔══██╗██╔═══██╗     ██║██╔════╝██╔════╝╚══██╔══╝   ╚══██╔══╝██║   ██║██║`,
		`██████╔╝ ╚████╔╝ ██████╔╝██████╔╝██║   ██║     ██║█████╗  ██║        ██║         ██║   ██║   ██║██║`,
		`██╔═══╝   ╚██╔╝  ██╔═══╝ ██╔══██╗██║   ██║██   ██║██╔══╝  ██║        ██║         ██║   ██║   ██║██║`,
		`██║        ██║   ██║     ██║  ██║╚██████╔╝╚██████╔╝███████╗╚██████╗  ██║         ██║   ╚██████╔╝██║`,
		`╚═╝        ╚═╝   ╚═╝     ╚═╝  ╚═╝ ╚═════╝  ╚═════╝ ╚══════╝ ╚═════╝  ╚═╝         ╚═╝    ╚═════╝ ╚═╝`,
	}
	minLogoWidth := 72
	if m.width < minLogoWidth {
		return ui.CenterBox(m.renderMenu(), m.width, ui.BodyHeight(m.height))
	}
	logo := m.sandField.render(logoLines)
	subtitle := s.StatusBar.Render("Keyboard-driven metadata & dependency manager")

	content := fmt.Sprintf("%s\n\n%s\n\n\n%s", logo, subtitle, m.renderMenu())
	return ui.CenterBox(content, m.width, ui.BodyHeight(m.height))
}

func (m AppModel) renderMenu() string {
	s := theme.Styles()
	if s == nil {
		return ""
	}
	var menuLines []string
	for i, opt := range dashboardOptions {
		prefix := "  "
		var keyStr, labelStr, descStr string
		if i == m.dashboardCursor {
			prefix = "→ "
			keyStr = s.Cursor.Render(opt.key)
			labelStr = s.Cursor.Render(fmt.Sprintf("%-15s", opt.label))
		} else {
			keyStr = s.MenuKey.Render(opt.key)
			labelStr = s.Value.Render(fmt.Sprintf("%-15s", opt.label))
		}
		descStr = s.MenuDesc.Render(opt.desc)
		row := fmt.Sprintf("%s[%s] %s  %s", prefix, keyStr, labelStr, descStr)
		menuLines = append(menuLines, row)
	}
	return strings.Join(menuLines, "\n")
}

func openFileNative(path string) tea.Cmd {
	return func() tea.Msg {
		var cmd string
		var args []string
		switch runtime.GOOS {
		case "windows":
			cmd = "cmd"
			args = []string{"/c", "start", "", path}
		case "darwin":
			cmd = "open"
			args = []string{path}
		default:
			cmd = "xdg-open"
			args = []string{path}
		}
		if err := exec.Command(cmd, args...).Start(); err != nil {
			return NativeOpenFileResultMsg{Err: err}
		}
		return NativeOpenFileResultMsg{}
	}
}

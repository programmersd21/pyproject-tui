// Package theme provides a semantic theme engine for consistent styling across the application.
package theme

import (
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// Theme defines semantic color roles for the entire application.
// Components consume semantic roles rather than hardcoded colors.
type Theme struct {
	// Identity
	Name        string
	DisplayName string

	// Semantic Colors
	Background       string // Main background
	Surface          string // Panel/card background
	SurfaceSecondary string // Secondary panels

	Text      string // Primary text
	TextMuted string // Secondary/muted text
	TextDim   string // Disabled/inactive text

	Accent       string // Primary accent (selection, focus)
	AccentSoft   string // Subtle accent highlights
	AccentBright string // Bright accent for emphasis

	Success string // Success states
	Warning string // Warning states
	Error   string // Error states

	Border        string // Default borders
	BorderFocused string // Focused element borders
	Selection     string // Selected items background

	Shadow string // Shadows and depth

	// Preview Colors (for theme selector)
	PreviewColors []string
}

// Registry holds all available themes and manages the active theme.
type Registry struct {
	mu      sync.RWMutex
	themes  map[string]*Theme
	active  *Theme
	styles  *StyleSet
	display DisplayConfig
	ready   bool
}

// DisplayConfig controls non-color presentation settings.
type DisplayConfig struct {
	Density         string
	BorderStyle     string
	Animations      bool
	ShowLineNumbers bool
}

var globalRegistry = &Registry{
	themes: make(map[string]*Theme),
	display: DisplayConfig{
		Density:     "normal",
		BorderStyle: "rounded",
		Animations:  true,
	},
}

var initOnce sync.Once

// Initialize registers built-in themes and sets the default active theme.
func Initialize() {
	initOnce.Do(func() {
		globalRegistry.mu.Lock()
		defer globalRegistry.mu.Unlock()
		globalRegistry.themes["tokyo-night"] = tokyoNight()
		globalRegistry.themes["catppuccin"] = catppuccin()
		globalRegistry.themes["nord"] = nord()
		globalRegistry.themes["gruvbox"] = gruvbox()
		globalRegistry.themes["rose-pine"] = rosePine()
		globalRegistry.themes["everforest"] = everforest()
		globalRegistry.themes["python"] = python()
		globalRegistry.themes["midnight"] = midnight()
		globalRegistry.themes["minimal"] = minimal()
		globalRegistry.themes["sage"] = sage()
		globalRegistry.active = globalRegistry.themes["tokyo-night"]
		globalRegistry.styles = generateStyles(globalRegistry.active, globalRegistry.display)
		globalRegistry.ready = true
	})
}

// Register adds a theme to the global registry.
func Register(t *Theme) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.themes[t.Name] = t
}

// SetActive activates a theme by name and regenerates all styles.
func SetActive(name string) {
	Initialize()
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	if t, ok := globalRegistry.themes[name]; ok {
		globalRegistry.active = t
		globalRegistry.styles = generateStyles(t, globalRegistry.display)
	}
}

// SetDisplayConfig updates display preferences and regenerates styles.
func SetDisplayConfig(cfg DisplayConfig) {
	Initialize()
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.display = cfg
	if globalRegistry.active != nil {
		globalRegistry.styles = generateStyles(globalRegistry.active, globalRegistry.display)
	}
}

// Active returns the currently active theme.
func Active() *Theme {
	Initialize()
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	if globalRegistry.active == nil {
		return nil
	}
	return globalRegistry.active
}

// Styles returns the generated styles for the active theme.
func Styles() *StyleSet {
	Initialize()
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	return globalRegistry.styles
}

// List returns all registered theme names.
func List() []string {
	Initialize()
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	names := make([]string, 0, len(globalRegistry.themes))
	for name := range globalRegistry.themes {
		names = append(names, name)
	}
	return names
}

// Get returns a theme by name.
func Get(name string) (*Theme, bool) {
	Initialize()
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	t, ok := globalRegistry.themes[name]
	return t, ok
}

// StyleSet contains all Lip Gloss styles generated from the active theme.
type StyleSet struct {
	// Layout
	TitleBar  lipgloss.Style
	StatusBar lipgloss.Style
	Sidebar   lipgloss.Style
	Editor    lipgloss.Style
	EditorBox lipgloss.Style
	Divider   lipgloss.Style

	// Sidebar
	SidebarItem     lipgloss.Style
	SidebarSelected lipgloss.Style
	SidebarFocused  lipgloss.Style

	// Editor
	Key        lipgloss.Style
	Value      lipgloss.Style
	ArrayValue lipgloss.Style
	Cursor     lipgloss.Style
	DirtyField lipgloss.Style

	// Dialogs
	HelpBox    lipgloss.Style
	ConfirmBox lipgloss.Style
	InputBox   lipgloss.Style

	// Status
	DirtyDot lipgloss.Style
	Error    lipgloss.Style
	Success  lipgloss.Style
	Warning  lipgloss.Style

	// Dashboard
	Logo       lipgloss.Style
	MenuOption lipgloss.Style
	MenuCursor lipgloss.Style
	MenuKey    lipgloss.Style
	MenuDesc   lipgloss.Style
}

// generateStyles creates Lip Gloss styles from a theme.
func generateStyles(t *Theme, display DisplayConfig) *StyleSet {
	sidebarPadding := 1
	editorPadding := 2
	boxPaddingY := 2
	boxPaddingX := 3
	switch display.Density {
	case "compact":
		sidebarPadding = 0
		editorPadding = 1
		boxPaddingY = 1
		boxPaddingX = 2
	case "comfortable":
		sidebarPadding = 2
		editorPadding = 3
		boxPaddingY = 3
		boxPaddingX = 4
	}
	borderStyle := lipgloss.NormalBorder()
	switch display.BorderStyle {
	case "thick":
		borderStyle = lipgloss.ThickBorder()
	case "double":
		borderStyle = lipgloss.DoubleBorder()
	case "rounded":
		borderStyle = lipgloss.RoundedBorder()
	case "normal":
		borderStyle = lipgloss.NormalBorder()
	}
	return &StyleSet{
		// Layout
		TitleBar:  lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text)).Bold(true),
		StatusBar: lipgloss.NewStyle().Foreground(lipgloss.Color(t.TextMuted)),
		Sidebar: lipgloss.NewStyle().
			Padding(0, sidebarPadding),
		Editor: lipgloss.NewStyle().Padding(0, editorPadding),
		EditorBox: lipgloss.NewStyle().
			Border(borderStyle).
			BorderForeground(lipgloss.Color(t.Border)).
			Padding(0, editorPadding),
		Divider: lipgloss.NewStyle().Foreground(lipgloss.Color(t.Border)),

		// Sidebar
		SidebarItem:     lipgloss.NewStyle().Foreground(lipgloss.Color(t.TextMuted)),
		SidebarSelected: lipgloss.NewStyle().Foreground(lipgloss.Color(t.AccentSoft)).Bold(true),
		SidebarFocused:  lipgloss.NewStyle().BorderForeground(lipgloss.Color(t.BorderFocused)),

		// Editor
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color(t.AccentSoft)).Bold(true),
		Value:      lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text)),
		ArrayValue: lipgloss.NewStyle().Foreground(lipgloss.Color(t.TextMuted)),
		Cursor:     lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Bold(true),
		DirtyField: lipgloss.NewStyle().Foreground(lipgloss.Color(t.Warning)).Underline(true),

		// Dialogs
		HelpBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.BorderFocused)).
			Padding(boxPaddingY, boxPaddingX),
		ConfirmBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.Error)).
			Padding(1, boxPaddingX),
		InputBox: lipgloss.NewStyle().
			Border(borderStyle).
			BorderForeground(lipgloss.Color(t.Accent)).
			Padding(1, editorPadding),

		// Status
		DirtyDot: lipgloss.NewStyle().Foreground(lipgloss.Color(t.Warning)),
		Error:    lipgloss.NewStyle().Foreground(lipgloss.Color(t.Error)),
		Success:  lipgloss.NewStyle().Foreground(lipgloss.Color(t.Success)),
		Warning:  lipgloss.NewStyle().Foreground(lipgloss.Color(t.Warning)),

		// Dashboard
		Logo:       lipgloss.NewStyle().Bold(true),
		MenuOption: lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text)),
		MenuCursor: lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Bold(true),
		MenuKey:    lipgloss.NewStyle().Foreground(lipgloss.Color(t.AccentSoft)),
		MenuDesc:   lipgloss.NewStyle().Foreground(lipgloss.Color(t.TextMuted)),
	}
}

// Package ui contains layout helpers and styles for pyproject-tui.
package ui

import "github.com/programmersd21/pyproject-tui/internal/theme"

// ApplyTheme sets the active theme by name and refreshes all styles.
func ApplyTheme(name string) {
	theme.Initialize()
	theme.SetActive(name)
	RefreshStyles()
}

// ThemeNames returns all available theme names.
func ThemeNames() []string {
	theme.Initialize()
	return theme.List()
}

// ActiveThemeName returns the name of the currently active theme.
func ActiveThemeName() string {
	theme.Initialize()
	t := theme.Active()
	if t == nil {
		return "tokyo-night"
	}
	return t.Name
}

// NextTheme cycles to the next theme in the list.
func NextTheme() string {
	theme.Initialize()
	names := theme.List()
	if len(names) == 0 {
		return "tokyo-night"
	}

	current := ActiveThemeName()
	for i, name := range names {
		if name == current {
			next := names[(i+1)%len(names)]
			return next
		}
	}

	return names[0]
}

// GetTheme returns a theme by name.
func GetTheme(name string) *theme.Theme {
	theme.Initialize()
	t, _ := theme.Get(name)
	return t
}

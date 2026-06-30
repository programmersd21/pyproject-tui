package theme

import "testing"

func TestThemeRegistry(t *testing.T) {
	Initialize()
	testTheme := &Theme{
		Name:        "test-theme",
		DisplayName: "Test",
		Accent:      "#FF0000",
	}
	Register(testTheme)

	retrieved, ok := Get("test-theme")
	if !ok || retrieved.Name != "test-theme" {
		t.Error("Failed to retrieve registered theme")
	}
}

func TestSetActive(t *testing.T) {
	Initialize()
	SetActive("tokyo-night")
	active := Active()
	if active == nil || active.Name != "tokyo-night" {
		t.Error("SetActive failed")
	}

	styles := Styles()
	if styles == nil {
		t.Error("Styles should not be nil after SetActive")
	}
}

func TestAllThemesRegistered(t *testing.T) {
	Initialize()
	expected := []string{"tokyo-night", "catppuccin", "nord", "gruvbox", "rose-pine", "everforest", "python", "midnight", "minimal", "sage"}

	for _, name := range expected {
		if _, ok := Get(name); !ok {
			t.Errorf("Expected theme '%s' not registered", name)
		}
	}
}

func TestThemePreviewColors(t *testing.T) {
	Initialize()
	theme, _ := Get("tokyo-night")
	if len(theme.PreviewColors) == 0 {
		t.Error("Theme should have preview colors")
	}

	for _, color := range theme.PreviewColors {
		if len(color) != 7 || color[0] != '#' {
			t.Errorf("Invalid preview color: %s", color)
		}
	}
}

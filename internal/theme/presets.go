// Package theme provides a semantic theme engine for consistent styling across the application.
package theme

func tokyoNight() *Theme {
	return &Theme{
		Name:          "tokyo-night",
		DisplayName:   "Tokyo Night",
		Text:          "#c0caf5",
		TextMuted:     "#9aa5ce",
		TextDim:       "#565f89",
		Accent:        "#7aa2f7",
		AccentSoft:    "#7dcfff",
		AccentBright:  "#2ac3de",
		Success:       "#9ece6a",
		Warning:       "#e0af68",
		Error:         "#f7768e",
		Border:        "#292e42",
		BorderFocused: "#7aa2f7",
		Selection:     "#364a82",
		PreviewColors: []string{"#7aa2f7", "#7dcfff", "#9ece6a", "#e0af68", "#f7768e", "#bb9af7"},
	}
}

func catppuccin() *Theme {
	return &Theme{
		Name:          "catppuccin",
		DisplayName:   "Catppuccin",
		Text:          "#cdd6f4",
		TextMuted:     "#bac2de",
		TextDim:       "#6c7086",
		Accent:        "#89b4fa",
		AccentSoft:    "#94e2d5",
		AccentBright:  "#74c7ec",
		Success:       "#a6e3a1",
		Warning:       "#f9e2af",
		Error:         "#f38ba8",
		Border:        "#45475a",
		BorderFocused: "#89b4fa",
		Selection:     "#585b70",
		PreviewColors: []string{"#89b4fa", "#94e2d5", "#a6e3a1", "#f9e2af", "#f38ba8", "#cba6f7"},
	}
}

func nord() *Theme {
	return &Theme{
		Name:          "nord",
		DisplayName:   "Nord",
		Text:          "#eceff4",
		TextMuted:     "#d8dee9",
		TextDim:       "#4c566a",
		Accent:        "#88c0d0",
		AccentSoft:    "#8fbcbb",
		AccentBright:  "#81a1c1",
		Success:       "#a3be8c",
		Warning:       "#ebcb8b",
		Error:         "#bf616a",
		Border:        "#434c5e",
		BorderFocused: "#88c0d0",
		Selection:     "#434c5e",
		PreviewColors: []string{"#88c0d0", "#81a1c1", "#a3be8c", "#ebcb8b", "#bf616a", "#b48ead"},
	}
}

func gruvbox() *Theme {
	return &Theme{
		Name:          "gruvbox",
		DisplayName:   "Gruvbox",
		Text:          "#ebdbb2",
		TextMuted:     "#d5c4a1",
		TextDim:       "#665c54",
		Accent:        "#83a598",
		AccentSoft:    "#8ec07c",
		AccentBright:  "#b8bb26",
		Success:       "#b8bb26",
		Warning:       "#fabd2f",
		Error:         "#fb4934",
		Border:        "#504945",
		BorderFocused: "#83a598",
		Selection:     "#504945",
		PreviewColors: []string{"#83a598", "#8ec07c", "#b8bb26", "#fabd2f", "#fb4934", "#d3869b"},
	}
}

func rosePine() *Theme {
	return &Theme{
		Name:          "rose-pine",
		DisplayName:   "Rose Pine",
		Text:          "#e0def4",
		TextMuted:     "#908caa",
		TextDim:       "#6e6a86",
		Accent:        "#c4a7e7",
		AccentSoft:    "#9ccfd8",
		AccentBright:  "#ebbcba",
		Success:       "#9ccfd8",
		Warning:       "#f6c177",
		Error:         "#eb6f92",
		Border:        "#403d52",
		BorderFocused: "#c4a7e7",
		Selection:     "#403d52",
		PreviewColors: []string{"#c4a7e7", "#9ccfd8", "#ebbcba", "#f6c177", "#eb6f92", "#31748f"},
	}
}

func everforest() *Theme {
	return &Theme{
		Name:          "everforest",
		DisplayName:   "Everforest",
		Text:          "#d3c6aa",
		TextMuted:     "#9da9a0",
		TextDim:       "#505a60",
		Accent:        "#7fbbb3",
		AccentSoft:    "#83c092",
		AccentBright:  "#a7c080",
		Success:       "#a7c080",
		Warning:       "#dbbc7f",
		Error:         "#e67e80",
		Border:        "#3d484d",
		BorderFocused: "#7fbbb3",
		Selection:     "#3d484d",
		PreviewColors: []string{"#7fbbb3", "#83c092", "#a7c080", "#dbbc7f", "#e67e80", "#d699b6"},
	}
}

func python() *Theme {
	return &Theme{
		Name:          "python",
		DisplayName:   "Python",
		Text:          "#F8F8F2",
		TextMuted:     "#B6C2D9",
		TextDim:       "#5E6A7D",
		Accent:        "#3776AB",
		AccentSoft:    "#FFD43B",
		AccentBright:  "#4B8BBE",
		Success:       "#A6E3A1",
		Warning:       "#FFD43B",
		Error:         "#FF6B6B",
		Border:        "#2A3554",
		BorderFocused: "#FFD43B",
		Selection:     "#203A5B",
		Shadow:        "#0F172A",
		PreviewColors: []string{"#3776AB", "#4B8BBE", "#FFD43B", "#E6B422", "#A6E3A1", "#FF6B6B"},
	}
}

func midnight() *Theme {
	return &Theme{
		Name:          "midnight",
		DisplayName:   "Midnight",
		Text:          "#d8dee9",
		TextMuted:     "#8892b0",
		TextDim:       "#3b4252",
		Accent:        "#a3be8c",
		AccentSoft:    "#81a1c1",
		AccentBright:  "#88c0d0",
		Success:       "#a3be8c",
		Warning:       "#d08770",
		Error:         "#bf616a",
		Border:        "#3b4252",
		BorderFocused: "#a3be8c",
		Selection:     "#2e3440",
		PreviewColors: []string{"#81a1c1", "#88c0d0", "#8fbcbb", "#a3be8c", "#d08770", "#bf616a"},
	}
}

func minimal() *Theme {
	return &Theme{
		Name:          "minimal",
		DisplayName:   "Minimal",
		Text:          "#e8e8e8",
		TextMuted:     "#707070",
		TextDim:       "#303030",
		Accent:        "#ffffff",
		AccentSoft:    "#a0a0a0",
		AccentBright:  "#e0e0e0",
		Success:       "#a0d0a0",
		Warning:       "#d8a15a",
		Error:         "#d46a6a",
		Border:        "#303030",
		BorderFocused: "#ffffff",
		Selection:     "#202020",
		PreviewColors: []string{"#707070", "#a0a0a0", "#e0e0e0", "#ffffff", "#d8a15a", "#d46a6a"},
	}
}

func sage() *Theme {
	return &Theme{
		Name:          "sage",
		DisplayName:   "Sage",
		Text:          "#ecefea",
		TextMuted:     "#6b7767",
		TextDim:       "#3d443b",
		Accent:        "#d6e4d0",
		AccentSoft:    "#a3b19b",
		AccentBright:  "#8f9e8b",
		Success:       "#a3b19b",
		Warning:       "#e2b07e",
		Error:         "#c87a7a",
		Border:        "#3d443b",
		BorderFocused: "#d6e4d0",
		Selection:     "#2c312a",
		PreviewColors: []string{"#6b7767", "#8f9e8b", "#a3b19b", "#d6e4d0", "#e2b07e", "#c87a7a"},
	}
}

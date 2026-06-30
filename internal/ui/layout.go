// Package ui contains layout helpers and styles for pyproject-tui.
package ui

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

// SidebarWidth returns responsive sidebar width based on terminal width.
func SidebarWidth(totalWidth int) int {
	// Responsive: 25 (small), 28 (medium), 32 (large)
	if totalWidth < 80 {
		return 25
	}
	if totalWidth < 120 {
		return 28
	}
	return 32
}

// EditorWidth returns the editor pane width given total terminal width.
func EditorWidth(totalWidth int) int {
	sidebar := SidebarWidth(totalWidth)
	editorW := totalWidth - sidebar - 1 // 1 char for vertical divider
	if editorW < 20 {
		editorW = 20
	}
	return editorW
}

// BodyHeight returns the height available for panes.
func BodyHeight(totalHeight int) int {
	// title=1, status=1, 2 borders=2
	reserved := 4
	h := totalHeight - reserved
	if h < 5 {
		h = 5
	}
	return h
}

// CenterBox returns a string centered in a box of the given dimensions using lipgloss.Place.
func CenterBox(content string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// TruncateString truncates s to maxLen runes, appending "..." if truncated.
func TruncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	if maxLen <= 1 {
		return "…"
	}
	return string(runes[:maxLen-1]) + "…"
}

// RenderStringMap renders a map[string]string as an inline TOML-like string.
func RenderStringMap(m map[string]string, maxLen int) string {
	if len(m) == 0 {
		return "{}"
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(m))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s = %q", k, m[k]))
	}
	s := "{" + strings.Join(parts, ", ") + "}"
	return TruncateString(s, maxLen)
}

// RenderStringSlice renders []string inline.
func RenderStringSlice(s []string, inlineMax int) string {
	if len(s) == 0 {
		return "[]"
	}
	if len(s) > inlineMax {
		return fmt.Sprintf("[%d items]", len(s))
	}
	quoted := make([]string, 0, len(s))
	for _, item := range s {
		quoted = append(quoted, fmt.Sprintf("%q", item))
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

// RenderAny renders any TOML value as a single-line string for display.
func RenderAny(v any, depth int) string {
	_ = depth
	switch val := v.(type) {
	case nil:
		return "null"
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int64:
		return fmt.Sprintf("%d", val)
	case int:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%g", val)
	case []any:
		parts := make([]string, 0, len(val))
		for _, item := range val {
			parts = append(parts, RenderAny(item, depth+1))
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case []string:
		parts := make([]string, 0, len(val))
		for _, item := range val {
			parts = append(parts, fmt.Sprintf("%q", item))
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]any:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(val))
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s = %s", k, RenderAny(val[k], depth+1)))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case map[string]string:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(val))
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s = %q", k, val[k]))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	default:
		return fmt.Sprint(v)
	}
}

// ReplaceRuneAt replaces the rune at position pos in s with ch.
// If pos is out of bounds, s is returned unchanged.
func ReplaceRuneAt(s string, pos int, ch rune) string {
	runes := []rune(s)
	if pos < 0 || pos >= len(runes) {
		return s
	}
	runes[pos] = ch
	return string(runes)
}

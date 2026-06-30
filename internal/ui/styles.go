// Package ui contains layout helpers and styles for pyproject-tui.
// This package now acts as a compatibility shim for the new theme system.
package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/programmersd21/pyproject-tui/internal/theme"
)

// Legacy color variables for backward compatibility
var (
	ColorPrimary   = lipgloss.Color("#88C0D0")
	ColorSecondary = lipgloss.Color("#81A1C1")
	ColorAccent    = lipgloss.Color("#A3BE8C")
	ColorOrange    = lipgloss.Color("#D08770")
	ColorRed       = lipgloss.Color("#BF616A")
	ColorText      = lipgloss.Color("#D8DEE9")
	ColorSubtext   = lipgloss.Color("#8892B0")
	ColorDim       = lipgloss.Color("#3B4252")
	ColorBg        = lipgloss.Color("#1E222B")
	ColorSelected  = lipgloss.Color("#2E3440")
)

// Style accessors - these now delegate to the theme system
var (
	TitleBarStyle        = theme.Styles().TitleBar
	StatusBarStyle       = theme.Styles().StatusBar
	SidebarStyle         = theme.Styles().Sidebar
	EditorStyle          = theme.Styles().Editor
	DividerStyle         = theme.Styles().Divider
	SidebarItemStyle     = theme.Styles().SidebarItem
	SidebarSelectedStyle = theme.Styles().SidebarSelected
	SidebarFocusedStyle  = theme.Styles().SidebarFocused
	KeyStyle             = theme.Styles().Key
	ValueStyle           = theme.Styles().Value
	ArrayValueStyle      = theme.Styles().ArrayValue
	CursorStyle          = theme.Styles().Cursor
	DirtyFieldStyle      = theme.Styles().DirtyField
	HelpBoxStyle         = theme.Styles().HelpBox
	ConfirmBoxStyle      = theme.Styles().ConfirmBox
	DirtyDotStyle        = theme.Styles().DirtyDot
	ErrorStyle           = theme.Styles().Error
	SuccessStyle         = theme.Styles().Success
)

// RefreshStyles updates all style references from the theme system.
// Call this after changing the active theme.
func RefreshStyles() {
	s := theme.Styles()
	if s == nil {
		return
	}

	TitleBarStyle = s.TitleBar
	StatusBarStyle = s.StatusBar
	SidebarStyle = s.Sidebar
	EditorStyle = s.Editor
	DividerStyle = s.Divider
	SidebarItemStyle = s.SidebarItem
	SidebarSelectedStyle = s.SidebarSelected
	SidebarFocusedStyle = s.SidebarFocused
	KeyStyle = s.Key
	ValueStyle = s.Value
	ArrayValueStyle = s.ArrayValue
	CursorStyle = s.Cursor
	DirtyFieldStyle = s.DirtyField
	HelpBoxStyle = s.HelpBox
	ConfirmBoxStyle = s.ConfirmBox
	DirtyDotStyle = s.DirtyDot
	ErrorStyle = s.Error
	SuccessStyle = s.Success
}

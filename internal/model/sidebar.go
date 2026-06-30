// Package model contains Bubble Tea models for the pyproject-tui application.
package model

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/programmersd21/pyproject-tui/internal/parser"
	"github.com/programmersd21/pyproject-tui/internal/theme"
)

type sidebarDelegate struct{}

func (d sidebarDelegate) Height() int                             { return 1 }
func (d sidebarDelegate) Spacing() int                            { return 0 }
func (d sidebarDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d sidebarDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	it := item.(SidebarItem)
	s := theme.Styles()
	if s == nil {
		_, _ = fmt.Fprint(w, "  "+it.title)
		return
	}

	style := s.SidebarItem
	prefix := "  "
	if index == m.Index() {
		style = s.SidebarSelected
		prefix = "→ "
	}
	if m.IsFiltered() {
		style = style.Faint(true)
	}
	_, _ = fmt.Fprint(w, style.Render(prefix+it.title))
}

// SidebarItem is a navigable section entry.
type SidebarItem struct {
	id    parser.SectionID
	title string
}

// FilterValue returns the searchable value for the item.
func (s SidebarItem) FilterValue() string { return s.title }

// SidebarModel implements the left pane.
type SidebarModel struct {
	list     list.Model
	sections []SidebarItem
	focused  bool
	width    int
	height   int
}

// NewSidebar creates a sidebar from the parsed file.
func NewSidebar(p *parser.PyProject) SidebarModel {
	items := make([]list.Item, 0, 8)
	if p != nil && p.Project != nil {
		items = append(items, SidebarItem{id: parser.SectionID{Kind: parser.KindProject}, title: "[project]"})
	}
	if p != nil && p.BuildSystem != nil {
		items = append(items, SidebarItem{id: parser.SectionID{Kind: parser.KindBuildSystem}, title: "[build-system]"})
	}
	if p != nil && len(p.Tool) > 0 {
		keys := make([]string, 0, len(p.Tool))
		for k := range p.Tool {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			items = append(items, SidebarItem{id: parser.SectionID{Kind: parser.KindTool, ToolKey: k}, title: "[tool." + k + "]"})
		}
	}
	if len(items) == 0 {
		items = append(items, SidebarItem{id: parser.SectionID{Kind: parser.KindProject}, title: "[project]"})
	}
	l := list.New(items, sidebarDelegate{}, 30, 10)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)

	if t := theme.Active(); t != nil {
		l.Styles.TitleBar = lipgloss.NewStyle()
		l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text))
		l.Styles.HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(t.TextDim))
		l.Styles.NoItems = lipgloss.NewStyle().Foreground(lipgloss.Color(t.TextDim))
		l.Styles.StatusBar = lipgloss.NewStyle().Foreground(lipgloss.Color(t.TextDim))
		l.Styles.StatusEmpty = lipgloss.NewStyle().Foreground(lipgloss.Color(t.TextDim))
		l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(t.TextDim))
		l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent))
		l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent))
		l.Styles.Spinner = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent))
	}

	return SidebarModel{list: l, sections: toSidebarItems(items), width: 30, height: 10}
}

func toSidebarItems(items []list.Item) []SidebarItem {
	out := make([]SidebarItem, 0, len(items))
	for _, item := range items {
		out = append(out, item.(SidebarItem))
	}
	return out
}

// UpdateSections updates the items and sections in the sidebar model.
func (m *SidebarModel) UpdateSections(p *parser.PyProject) {
	newSidebar := NewSidebar(p)
	m.sections = newSidebar.sections

	// Convert SidebarItem list to []list.Item
	listItems := make([]list.Item, len(newSidebar.sections))
	for i, item := range newSidebar.sections {
		listItems[i] = item
	}
	m.list.SetItems(listItems)
}

// Update updates the sidebar model.
func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the sidebar pane.
func (m SidebarModel) View() string {
	s := theme.Styles()
	if s == nil {
		return m.list.View()
	}

	content := m.list.View()
	if !m.focused {
		content = lipgloss.NewStyle().Faint(true).Render(content)
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(0, 1).
		Render(content)
}

// SelectedSection returns the currently selected section.
func (m SidebarModel) SelectedSection() parser.SectionID {
	if len(m.sections) == 0 {
		return parser.SectionID{Kind: parser.KindProject}
	}
	i := m.list.Index()
	if i < 0 || i >= len(m.sections) {
		return m.sections[0].id
	}
	return m.sections[i].id
}

// SetFocused changes focus styling.
func (m *SidebarModel) SetFocused(focused bool) {
	m.focused = focused
}

// SetSize updates the sidebar viewport size.
func (m *SidebarModel) SetSize(width, height int) {
	m.width = width
	// Account for Padding(0, 1) = 2 chars consumed by style
	contentWidth := width - 2
	if contentWidth < 10 {
		contentWidth = 10
	}
	m.height = height
	m.list.SetSize(contentWidth, height)
}

// MoveToSection selects the given section.
func (m *SidebarModel) MoveToSection(id parser.SectionID) {
	for i, sec := range m.sections {
		if sec.id == id {
			m.list.Select(i)
			return
		}
	}
}

func (m SidebarModel) String() string {
	return strings.TrimSpace(lipgloss.NewStyle().Render(m.View()))
}

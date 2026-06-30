package model

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines keybindings used by the application.
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Edit     key.Binding
	Add      key.Binding
	Delete   key.Binding
	Save     key.Binding
	Undo     key.Binding
	Help     key.Binding
	Quit     key.Binding
	Confirm  key.Binding
	Cancel   key.Binding
}

// DefaultKeyMap returns the default bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Left:     key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
		Right:    key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
		Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch pane")),
		ShiftTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "switch pane")),
		Edit:     key.NewBinding(key.WithKeys("e", "enter"), key.WithHelp("e/↵", "edit")),
		Add:      key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		Delete:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		Save:     key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "save")),
		Undo:     key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "undo")),
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Confirm:  key.NewBinding(key.WithKeys("enter", "y"), key.WithHelp("↵/y", "confirm")),
		Cancel:   key.NewBinding(key.WithKeys("esc", "n"), key.WithHelp("esc/n", "cancel")),
	}
}

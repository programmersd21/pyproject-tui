package model

import "github.com/programmersd21/pyproject-tui/internal/parser"

// TickMsg is sent to trigger animated gradient updates on the homepage.
type TickMsg struct{}

// SectionSelectedMsg is sent when a section is selected in the sidebar.
type SectionSelectedMsg struct {
	ID parser.SectionID
}

// FieldEditedMsg is sent when a field edit is committed.
type FieldEditedMsg struct {
	Key      string
	NewValue any
}

// FieldDeletedMsg is sent when a field is deleted.
type FieldDeletedMsg struct {
	Key string
}

// FieldAddedMsg is sent when a new key-value pair is added.
type FieldAddedMsg struct {
	Key   string
	Value any
}

// OpenFileMsg requests opening a different pyproject.toml file.
type OpenFileMsg struct {
	Path string
}

// SaveResultMsg is sent when save completes or errors.
type SaveResultMsg struct {
	Err error
}

// PendingAction identifies a confirmation dialog action.
type PendingAction int

const (
	// ActionDelete confirms deletion of a field.
	ActionDelete PendingAction = iota
	// ActionQuitDirty confirms quitting with unsaved changes.
	ActionQuitDirty
	// ActionDeleteSection confirms deletion of a tool section.
	ActionDeleteSection
)

// ConfirmMsg is the result of a confirmation dialog.
type ConfirmMsg struct {
	Confirmed bool
	Action    PendingAction
}

// NativeOpenFileMsg requests opening a file with the OS default application.
type NativeOpenFileMsg struct {
	Path string
}

// NativeOpenFileResultMsg is the result of attempting to open a file natively.
type NativeOpenFileResultMsg struct {
	Err error
}

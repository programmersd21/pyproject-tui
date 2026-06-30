package model

import (
	"testing"

	"github.com/programmersd21/pyproject-tui/internal/parser"
)

func TestParseInputValue(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"\"hello\"", "hello"},
		{"123", int64(123)},
		{"true", true},
		{"false", false},
		{"[\"a\", \"b\"]", []any{"a", "b"}},
		{"literal string", "literal string"},
	}

	for _, tt := range tests {
		got := parseInputValue(tt.input)
		switch gotVal := got.(type) {
		case []any:
			expVal, ok := tt.expected.([]any)
			if !ok || len(gotVal) != len(expVal) {
				t.Errorf("parseInputValue(%q) = %#v, want %#v", tt.input, got, tt.expected)
			} else {
				for i := range gotVal {
					if gotVal[i] != expVal[i] {
						t.Errorf("parseInputValue(%q)[%d] = %#v, want %#v", tt.input, i, gotVal[i], expVal[i])
					}
				}
			}
		default:
			if got != tt.expected {
				t.Errorf("parseInputValue(%q) = %#v, want %#v", tt.input, got, tt.expected)
			}
		}
	}
}

func TestAppModelWorkflow(t *testing.T) {
	pp := parser.NewEmpty("pyproject.toml")
	app := NewAppModel(pp, false)

	// Verify initial state
	if app.dirty {
		t.Error("new app model should not be dirty")
	}

	// 1. Edit a field
	app.applyEdit("name", "new-package")
	if pp.Project.Name != "new-package" {
		t.Errorf("expected project name to be 'new-package', got %q", pp.Project.Name)
	}

	// 2. Add a field
	app.applyAdd("description", "A description")
	if pp.Project.Description != "A description" {
		t.Errorf("expected project description to be 'A description', got %q", pp.Project.Description)
	}

	// 3. Delete a field
	app.applyDelete("description")
	if pp.Project.Description != "" {
		t.Errorf("expected project description to be empty, got %q", pp.Project.Description)
	}
}

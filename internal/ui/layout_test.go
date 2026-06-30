package ui

import "testing"

func TestEditorWidth(t *testing.T) {
	if got := EditorWidth(100); got != 71 {
		t.Fatalf("EditorWidth(100) = %d, want 71", got)
	}
}

func TestTruncateString(t *testing.T) {
	got := TruncateString("abcdef", 4)
	if got != "abc…" {
		t.Fatalf("TruncateString = %q, want %q", got, "abc…")
	}
}

func TestRenderStringSlice(t *testing.T) {
	got := RenderStringSlice([]string{"a", "b"}, 5)
	if got != `["a", "b"]` {
		t.Fatalf("RenderStringSlice = %q, want %q", got, `["a", "b"]`)
	}
}

package model

import "testing"

func TestStatusNextStatus_Cycle(t *testing.T) {
	tests := []struct {
		current  Status
		expected Status
	}{
		{StatusTodo, StatusInProgress},
		{StatusInProgress, StatusDone},
		{StatusDone, StatusCancelled},
		{StatusCancelled, StatusTodo},
	}

	for _, tt := range tests {
		got := tt.current.NextStatus()
		if got != tt.expected {
			t.Errorf("%s.NextStatus() = %s, want %s", tt.current, got, tt.expected)
		}
	}
}

func TestStatusSymbol(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusTodo, "[ ]"},
		{StatusInProgress, "[…]"},
		{StatusDone, "[✓]"},
		{StatusCancelled, "[x]"},
	}

	for _, tt := range tests {
		got := tt.status.Symbol()
		if got != tt.expected {
			t.Errorf("%s.Symbol() = %q, want %q", tt.status, got, tt.expected)
		}
	}
}

func TestStatusLabel(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusTodo, ""},
		{StatusInProgress, "@in_progress"},
		{StatusDone, "@done"},
		{StatusCancelled, "@cancelled"},
	}

	for _, tt := range tests {
		got := tt.status.Label()
		if got != tt.expected {
			t.Errorf("%s.Label() = %q, want %q", tt.status, got, tt.expected)
		}
	}
}

func TestStatusString(t *testing.T) {
	if StatusTodo.String() != "todo" {
		t.Errorf("expected 'todo', got %q", StatusTodo.String())
	}
}

func TestInvalidStatusDefaults(t *testing.T) {
	invalid := Status("unknown")
	if invalid.NextStatus() != StatusTodo {
		t.Error("invalid status should default to todo")
	}
	if invalid.Symbol() != "[ ]" {
		t.Error("invalid status symbol should default to [ ]")
	}
	if invalid.Label() != "" {
		t.Error("invalid status label should be empty")
	}
}

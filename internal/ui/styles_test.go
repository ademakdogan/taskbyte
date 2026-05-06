package ui

import (
	"testing"

	"github.com/adem/taskbyte/internal/config"
)

func TestNewStyles(t *testing.T) {
	cfg := config.DefaultConfig()
	styles := NewStyles(cfg)

	// Verify styles are created without panic
	_ = styles.TodoStyle.Render("test")
	_ = styles.InProgressStyle.Render("test")
	_ = styles.DoneStyle.Render("test")
	_ = styles.CancelledStyle.Render("test")
	_ = styles.DateHeader.Render("test")
	_ = styles.FocusedItem.Render("test")
	_ = styles.ErrorStyle.Render("test")
	_ = styles.HelpStyle.Render("test")
}

func TestGetColor_ValidColors(t *testing.T) {
	for name := range Colors {
		c := getColor(name)
		if c == "" {
			t.Errorf("expected color for %q, got empty", name)
		}
	}
}

func TestGetColor_Invalid(t *testing.T) {
	c := getColor("nonexistent")
	if c == "" {
		t.Error("invalid color should fallback to white")
	}
}

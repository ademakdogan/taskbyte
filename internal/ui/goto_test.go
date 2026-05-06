package ui

import (
	"testing"

	"github.com/adem/taskbyte/internal/config"
	"github.com/adem/taskbyte/internal/db"
)

func TestRenderGoto(t *testing.T) {
	cfg := config.DefaultConfig()
	styles := NewStyles(cfg)

	stats := []db.DateStats{
		{Date: "2026-03-31", Total: 5, Todo: 2, InProgress: 1, Done: 1, Cancelled: 1},
		{Date: "2026-04-01", Total: 3, Todo: 1, InProgress: 0, Done: 2, Cancelled: 0},
	}

	output := RenderGoto(stats, 0, "", cfg.DateFormat, 80, styles)

	if output == "" {
		t.Error("goto render should not be empty")
	}
}

func TestRenderGoto_Empty(t *testing.T) {
	cfg := config.DefaultConfig()
	styles := NewStyles(cfg)

	output := RenderGoto(nil, 0, "", cfg.DateFormat, 80, styles)

	if output == "" {
		t.Error("empty goto render should not be empty")
	}
}

func TestRenderGoto_WithInput(t *testing.T) {
	cfg := config.DefaultConfig()
	styles := NewStyles(cfg)

	stats := []db.DateStats{
		{Date: "2026-03-31", Total: 5, Todo: 2, InProgress: 1, Done: 1, Cancelled: 1},
	}

	output := RenderGoto(stats, 0, "31.03", cfg.DateFormat, 80, styles)

	if output == "" {
		t.Error("goto render with input should not be empty")
	}
}

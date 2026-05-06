package ui

import (
	"testing"

	"github.com/adem/taskbyte/internal/config"
	"github.com/adem/taskbyte/internal/db"
)

func TestAggregateStats(t *testing.T) {
	dateStats := []db.DateStats{
		{Date: "2026-03-31", Total: 5, Todo: 2, InProgress: 1, Done: 1, Cancelled: 1},
		{Date: "2026-04-01", Total: 3, Todo: 1, InProgress: 0, Done: 2, Cancelled: 0},
	}

	result := AggregateStats(dateStats)

	if result.Total != 8 {
		t.Errorf("expected 8 total, got %d", result.Total)
	}
	if result.Todo != 3 {
		t.Errorf("expected 3 todo, got %d", result.Todo)
	}
	if result.Done != 3 {
		t.Errorf("expected 3 done, got %d", result.Done)
	}
}

func TestAggregateStats_Empty(t *testing.T) {
	result := AggregateStats(nil)
	if result.Total != 0 {
		t.Errorf("expected 0 total for nil, got %d", result.Total)
	}
}

func TestRenderStats(t *testing.T) {
	cfg := config.DefaultConfig()
	styles := NewStyles(cfg)

	data := StatsData{Todo: 5, InProgress: 2, Done: 8, Cancelled: 1, Total: 16}
	output := RenderStats(data, "all", styles)

	if output == "" {
		t.Error("stats render should not be empty")
	}
}

func TestRenderStats_Empty(t *testing.T) {
	cfg := config.DefaultConfig()
	styles := NewStyles(cfg)

	data := StatsData{}
	output := RenderStats(data, "all", styles)

	if output == "" {
		t.Error("empty stats render should not be empty")
	}
}

package service

import (
	"testing"
	"time"

	"github.com/adem/taskbyte/internal/db"
)

func TestFilterStatsByRange_All(t *testing.T) {
	stats := []db.DateStats{
		{Date: "2020-01-01", Total: 1},
		{Date: "2026-03-31", Total: 2},
	}

	filtered := FilterStatsByRange(stats, "all")
	if len(filtered) != 2 {
		t.Errorf("expected 2 for 'all', got %d", len(filtered))
	}
}

func TestFilterStatsByRange_Default(t *testing.T) {
	stats := []db.DateStats{
		{Date: "2020-01-01", Total: 1},
	}

	filtered := FilterStatsByRange(stats, "unknown")
	if len(filtered) != 1 {
		t.Errorf("unknown range should return all, got %d", len(filtered))
	}
}

func TestFilterStatsByRange_Day(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	stats := []db.DateStats{
		{Date: yesterday, Total: 1},
		{Date: today, Total: 2},
	}

	filtered := FilterStatsByRange(stats, "day")
	if len(filtered) != 1 {
		t.Errorf("expected 1 for 'day', got %d", len(filtered))
	}
	if filtered[0].Date != today {
		t.Errorf("expected today's date, got %s", filtered[0].Date)
	}
}

func TestFilterStatsByRange_Month(t *testing.T) {
	now := time.Now()
	thisMonth := now.Format("2006-01-02")
	lastMonth := now.AddDate(0, -1, 0).Format("2006-01-02")

	stats := []db.DateStats{
		{Date: lastMonth, Total: 1},
		{Date: thisMonth, Total: 2},
	}

	filtered := FilterStatsByRange(stats, "month")
	if len(filtered) < 1 {
		t.Error("month filter should include at least this month's data")
	}
}

func TestFilterStatsByRange_InvalidDate(t *testing.T) {
	stats := []db.DateStats{
		{Date: "invalid-date", Total: 1},
		{Date: time.Now().Format("2006-01-02"), Total: 2},
	}

	filtered := FilterStatsByRange(stats, "day")
	// Invalid date should be skipped
	if len(filtered) != 1 {
		t.Errorf("expected 1 (invalid skipped), got %d", len(filtered))
	}
}

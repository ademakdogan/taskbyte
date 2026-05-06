package service

import (
	"time"

	"github.com/adem/taskbyte/internal/db"
)

// FilterStatsByRange filters date stats based on a time range.
func FilterStatsByRange(stats []db.DateStats, rangeType string) []db.DateStats {
	now := time.Now()
	var startDate time.Time

	switch rangeType {
	case "day":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startDate = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "all":
		return stats
	default:
		return stats
	}

	var filtered []db.DateStats
	for _, s := range stats {
		t, err := time.Parse("2006-01-02", s.Date)
		if err != nil {
			continue
		}
		if !t.Before(startDate) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

package ui

import (
"fmt"
"math"
"strings"

"github.com/adem/taskbyte/internal/db"
)

// StatsData holds aggregated statistics.
type StatsData struct {
	Todo       int
	InProgress int
	Done       int
	Cancelled  int
	Total      int
}

// AggregateStats calculates totals from date stats.
func AggregateStats(dateStats []db.DateStats) StatsData {
	var sd StatsData
	for _, ds := range dateStats {
		sd.Todo += ds.Todo
		sd.InProgress += ds.InProgress
		sd.Done += ds.Done
		sd.Cancelled += ds.Cancelled
		sd.Total += ds.Total
	}
	return sd
}

// RenderStats renders the statistics view with a simple pie chart.
func RenderStats(data StatsData, rangeLabel string, styles Styles) string {
	var s strings.Builder

	title := fmt.Sprintf("=================== STATS (%s) ===================", strings.ToUpper(rangeLabel))
	s.WriteString(styles.Title.Render(title) + "\n\n")

	if data.Total == 0 {
		s.WriteString(styles.Subtle.Render("  No tasks to display.\n"))
		return s.String()
	}

	// Bar chart
	barWidth := 40
	segments := []struct {
		label string
		count int
		block string
		style func(string) string
	}{
		{"Todo", data.Todo, "█", func(t string) string { return styles.TodoStyle.Render(t) }},
		{"In Progress", data.InProgress, "█", func(t string) string { return styles.InProgressStyle.Render(t) }},
		{"Done", data.Done, "█", func(t string) string { return styles.DoneStyle.Render(t) }},
		{"Cancelled", data.Cancelled, "█", func(t string) string { return styles.CancelledStyle.Render(t) }},
	}

	// Render bar
	s.WriteString("  ")
	for _, seg := range segments {
		if seg.count == 0 {
			continue
		}
		width := int(math.Round(float64(seg.count) / float64(data.Total) * float64(barWidth)))
		if width == 0 && seg.count > 0 {
			width = 1
		}
		s.WriteString(seg.style(strings.Repeat(seg.block, width)))
	}
	s.WriteString("\n\n")

	// Legend
	for _, seg := range segments {
		if seg.count == 0 {
			continue
		}
		pct := float64(seg.count) / float64(data.Total) * 100
		line := fmt.Sprintf("  %s %-14s %3d  (%5.1f%%)",
seg.style("●"),
seg.label,
seg.count,
pct,
)
		s.WriteString(line + "\n")
	}

	s.WriteString(fmt.Sprintf("\n  Total: %d tasks\n", data.Total))

	s.WriteString("\n" + styles.Subtle.Render("  "+strings.Repeat("=", 52)) + "\n")
	s.WriteString(styles.HelpStyle.Render("[Esc]: Back"))

	return s.String()
}

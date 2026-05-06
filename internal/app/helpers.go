package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// clearErrorMsg is sent to clear the error after a timeout.
type clearErrorMsg struct{}

// clearErrorAfter returns a command that clears the error after a delay.
func clearErrorAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

// formatBytes returns a human-readable byte size.
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return "< 1 KB"
	}
	return ""
}

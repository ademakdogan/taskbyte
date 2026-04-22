package ui

import tea "github.com/charmbracelet/bubbletea"

// Key checks if a message matches a specific key.
func Key(msg tea.KeyMsg, keys ...string) bool {
	for _, k := range keys {
		if msg.String() == k {
			return true
		}
	}
	return false
}

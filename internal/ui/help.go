package ui

import "strings"

// SlashCommandHelp returns a formatted help text for all slash commands.
func SlashCommandHelp() string {
	commands := []struct {
		cmd  string
		desc string
	}{
		{"/date <DD.MM.YYYY>", "Switch to a different date"},
		{"/sort <type>", "Sort tasks: date, progress, date-reverse, progress-reverse"},
		{"/export [path]", "Export current list as markdown"},
		{"/hide", "Hide completed and cancelled tasks"},
		{"/show", "Show hidden tasks"},
		{"/migrate [date]", "Move unfinished tasks from date to today"},
		{"/migrate-all", "Move all unfinished tasks to today"},
		{"/stats [range]", "Show statistics: day, week, month, all"},
		{"/settings", "Open settings page"},
	}

	var s strings.Builder
	s.WriteString("Available Commands:\n")
	for _, c := range commands {
		s.WriteString("  " + c.cmd + "\n")
		s.WriteString("    " + c.desc + "\n")
	}
	return s.String()
}

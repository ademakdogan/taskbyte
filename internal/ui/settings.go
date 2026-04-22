package ui

import (
"fmt"
"strings"

"github.com/adem/taskbyte/internal/config"
)

// SettingItem represents a single setting entry.
type SettingItem struct {
	Label    string
	Type     string // "bool", "select"
	Value    string
	Options  []string
	Section  string // for grouping
}

// BuildSettingsItems creates the list of settings from config.
func BuildSettingsItems(cfg config.Config) []SettingItem {
	return []SettingItem{
		{Label: "Auto-hide Completed Tasks", Type: "bool", Value: boolStr(cfg.AutoHideCompleted), Section: "General"},
		{Label: "Insert Mode Prompt History", Type: "bool", Value: boolStr(cfg.InsertPromptHistory), Section: "General"},
		{Label: "Default Date Format", Type: "select", Value: cfg.DateFormat, Options: config.ValidDateFormats(), Section: "General"},
		{Label: "Todo Task Color", Type: "select", Value: cfg.Theme.TodoColor, Options: config.ValidColors(), Section: "Theme & Colors"},
		{Label: "In Progress Task Color", Type: "select", Value: cfg.Theme.InProgressColor, Options: config.ValidColors(), Section: "Theme & Colors"},
		{Label: "Done Task Color", Type: "select", Value: cfg.Theme.DoneColor, Options: config.ValidColors(), Section: "Theme & Colors"},
		{Label: "Cancelled Task Color", Type: "select", Value: cfg.Theme.CancelledColor, Options: config.ValidColors(), Section: "Theme & Colors"},
	}
}

// ApplySettingsItems applies the settings items back to a config.
func ApplySettingsItems(items []SettingItem) config.Config {
	cfg := config.DefaultConfig()
	for _, item := range items {
		switch item.Label {
		case "Auto-hide Completed Tasks":
			cfg.AutoHideCompleted = item.Value == "True"
		case "Insert Mode Prompt History":
			cfg.InsertPromptHistory = item.Value == "True"
		case "Default Date Format":
			cfg.DateFormat = item.Value
		case "Todo Task Color":
			cfg.Theme.TodoColor = item.Value
		case "In Progress Task Color":
			cfg.Theme.InProgressColor = item.Value
		case "Done Task Color":
			cfg.Theme.DoneColor = item.Value
		case "Cancelled Task Color":
			cfg.Theme.CancelledColor = item.Value
		}
	}
	return cfg
}

// RenderSettings renders the settings page.
func RenderSettings(items []SettingItem, cursor int, dropdownOpen bool, dropdownCursor int, styles Styles) string {
	var s strings.Builder

	title := "=================== SETTINGS ==================="
	s.WriteString(styles.Title.Render(title) + "\n\n")

	currentSection := ""
	for i, item := range items {
		// Section header
		if item.Section != currentSection {
			currentSection = item.Section
			if i > 0 {
				s.WriteString("\n")
			}
			if currentSection == "Theme & Colors" {
				s.WriteString(styles.Subtle.Render("  --- Theme & Colors ---") + "\n")
			}
		}

		prefix := "  "
		if i == cursor {
			prefix = styles.FocusedItem.Render("> ")
		}

		label := fmt.Sprintf("%-30s", item.Label)
		value := fmt.Sprintf("[ %s ]", item.Value)

		if i == cursor {
			label = styles.Highlight.Render(label)
			value = styles.FocusedItem.Render(value)
		}

		s.WriteString(prefix + label + " : " + value + "\n")

		// Dropdown
		if dropdownOpen && i == cursor && item.Type == "select" {
			s.WriteString(renderDropdown(item.Options, item.Value, dropdownCursor, styles))
		}
	}

	s.WriteString("\n" + styles.Subtle.Render("  "+strings.Repeat("=", 48)) + "\n")

	var help string
	if dropdownOpen {
		help = "[↑/↓]: Select  |  [Enter]: Confirm  |  [Esc]: Cancel"
	} else {
		help = "[↑/↓]: Navigate  |  [Enter]: Change  |  [Esc]: Return"
	}
	s.WriteString(styles.HelpStyle.Render(help))

	return s.String()
}

func renderDropdown(options []string, current string, cursor int, styles Styles) string {
	var s strings.Builder

	maxLen := 0
	for _, opt := range options {
		if len(opt) > maxLen {
			maxLen = len(opt)
		}
	}
	boxWidth := maxLen + 4

	indent := strings.Repeat(" ", 34)
	s.WriteString(indent + "+" + strings.Repeat("-", boxWidth) + "+\n")

	for i, opt := range options {
		prefix := "| "
		suffix := " |"
		if i == cursor {
			prefix = "|>"
			suffix = "<|"
		}

		padded := fmt.Sprintf("%-*s", maxLen+2, opt)
		if i == cursor {
			padded = styles.Highlight.Render(padded)
		} else if opt == current {
			padded = styles.FocusedItem.Render(padded)
		}

		s.WriteString(indent + prefix + padded + suffix + "\n")
	}

	s.WriteString(indent + "+" + strings.Repeat("-", boxWidth) + "+\n")
	return s.String()
}

func boolStr(b bool) string {
	if b {
		return "True"
	}
	return "False"
}

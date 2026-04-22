package ui

import (
"github.com/charmbracelet/lipgloss"

"github.com/adem/taskbyte/internal/config"
)

// Colors maps config color names to lipgloss colors.
var Colors = map[string]lipgloss.Color{
	"white":       lipgloss.Color("#FFFFFF"),
	"orange":      lipgloss.Color("#FF8C00"),
	"red":         lipgloss.Color("#FF4444"),
	"dark_green":  lipgloss.Color("#2E8B57"),
	"light_green": lipgloss.Color("#90EE90"),
	"cyan":        lipgloss.Color("#00CED1"),
	"blue":        lipgloss.Color("#4169E1"),
	"gray":        lipgloss.Color("#808080"),
	"yellow":      lipgloss.Color("#FFD700"),
	"magenta":     lipgloss.Color("#DA70D6"),
}

// Styles holds all application styles derived from config.
type Styles struct {
	// Task status styles
	TodoStyle       lipgloss.Style
	InProgressStyle lipgloss.Style
	DoneStyle       lipgloss.Style
	CancelledStyle  lipgloss.Style

	// UI element styles
	DateHeader    lipgloss.Style
	InputBorder   lipgloss.Style
	FocusedItem   lipgloss.Style
	StatusBar     lipgloss.Style
	ModeIndicator lipgloss.Style
	Title         lipgloss.Style
	Subtle        lipgloss.Style
	Highlight     lipgloss.Style
	ErrorStyle    lipgloss.Style
	HelpStyle     lipgloss.Style
}

// NewStyles creates styles from the current config.
func NewStyles(cfg config.Config) Styles {
	todoColor := getColor(cfg.Theme.TodoColor)
	inProgressColor := getColor(cfg.Theme.InProgressColor)
	doneColor := getColor(cfg.Theme.DoneColor)
	cancelledColor := getColor(cfg.Theme.CancelledColor)

	return Styles{
		TodoStyle:       lipgloss.NewStyle().Foreground(todoColor),
		InProgressStyle: lipgloss.NewStyle().Foreground(inProgressColor),
		DoneStyle:       lipgloss.NewStyle().Foreground(doneColor).Strikethrough(true),
		CancelledStyle:  lipgloss.NewStyle().Foreground(cancelledColor).Strikethrough(true),

		DateHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00CED1")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#555555")).
			Padding(0, 1),

		InputBorder: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#555555")).
			Padding(0, 1),

		FocusedItem: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00CED1")),

		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#333333")).
			Foreground(lipgloss.Color("#AAAAAA")).
			Padding(0, 1),

		ModeIndicator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00CED1")),

		Subtle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555")),

		Highlight: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")),

		ErrorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")).
			Bold(true),

		HelpStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
	}
}

func getColor(name string) lipgloss.Color {
	if c, ok := Colors[name]; ok {
		return c
	}
	return lipgloss.Color("#FFFFFF")
}

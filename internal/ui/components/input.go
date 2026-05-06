package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// InputBox renders a bordered input box with a label.
type InputBox struct {
	Label     string
	Value     string
	ModeLabel string
	Width     int
	Style     lipgloss.Style
	Subtle    lipgloss.Style
}

// NewInputBox creates a new InputBox.
func NewInputBox(label, value, modeLabel string, width int, style, subtle lipgloss.Style) InputBox {
	return InputBox{
		Label:     label,
		Value:     value,
		ModeLabel: modeLabel,
		Width:     width,
		Style:     style,
		Subtle:    subtle,
	}
}

// Render returns the string representation of the input box.
func (ib InputBox) Render() string {
	boxWidth := max(ib.Width, 45)
	innerWidth := boxWidth - 4 // account for borders

	top := ib.Subtle.Render("┌─ " + ib.Label + " " + strings.Repeat("─", max(0, innerWidth-len(ib.Label)-2)) + "┐")

	inputLine := "│ › " + ib.Value + "█"
	padding := max(0, innerWidth-len(ib.Value)-3)
	inputLine += strings.Repeat(" ", padding) + " │"

	bottomPadding := max(0, innerWidth-len(ib.ModeLabel)-2)
	bottom := ib.Subtle.Render("└" + strings.Repeat("─", bottomPadding) + " " + ib.ModeLabel + " ─┘")

	return top + "\n" + inputLine + "\n" + bottom
}

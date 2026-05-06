package components

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestInputBox_Render(t *testing.T) {
	style := lipgloss.NewStyle()
	subtle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555"))

	ib := NewInputBox("Today", "Hello world", "—insert—", 45, style, subtle)
	output := ib.Render()

	if output == "" {
		t.Error("render should not be empty")
	}
}

func TestInputBox_EmptyValue(t *testing.T) {
	style := lipgloss.NewStyle()
	subtle := lipgloss.NewStyle()

	ib := NewInputBox("Today", "", "—insert—", 45, style, subtle)
	output := ib.Render()

	if output == "" {
		t.Error("render with empty value should not be empty")
	}
}

func TestInputBox_LongValue(t *testing.T) {
	style := lipgloss.NewStyle()
	subtle := lipgloss.NewStyle()

	longText := "This is a very long task description that exceeds the normal width"
	ib := NewInputBox("Today", longText, "—edit—", 45, style, subtle)
	output := ib.Render()

	if output == "" {
		t.Error("render with long value should not be empty")
	}
}

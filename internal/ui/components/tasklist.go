package components

import (
"fmt"
"strings"

"github.com/charmbracelet/lipgloss"

"github.com/adem/taskbyte/internal/model"
)

// TaskListRenderer handles the rendering of task lists.
type TaskListRenderer struct {
	TodoStyle       lipgloss.Style
	InProgressStyle lipgloss.Style
	DoneStyle       lipgloss.Style
	CancelledStyle  lipgloss.Style
	FocusedStyle    lipgloss.Style
	ErrorStyle      lipgloss.Style
}

// RenderTaskLine renders a single task line with appropriate styling.
func (r TaskListRenderer) RenderTaskLine(task model.Task, focused bool) string {
	symbol := task.Status.Symbol()
	label := task.Status.Label()
	style := r.styleFor(task.Status)

	var parts []string
	parts = append(parts, symbol)
	parts = append(parts, task.Text)

	if label != "" {
		parts = append(parts, label)
	}

	text := strings.Join(parts, " ")

	// Add status timestamp
	if task.StatusChangedAt != nil && task.Status != model.StatusTodo {
		ts := task.StatusChangedAt.Local().Format("02.01.2006 15:04")
		text += " - " + ts
	}

	if focused {
		symbolRendered := r.FocusedStyle.Render(symbol)
		textRendered := style.Render(task.Text)
		result := symbolRendered + " " + textRendered

		if label != "" {
			result += " " + style.Render(label)
		}
		if task.StatusChangedAt != nil && task.Status != model.StatusTodo {
			ts := task.StatusChangedAt.Local().Format("02.01.2006 15:04")
			result += " " + style.Render(fmt.Sprintf("- %s", ts))
		}
		return result
	}

	return style.Render(text)
}

func (r TaskListRenderer) styleFor(status model.Status) lipgloss.Style {
	switch status {
	case model.StatusTodo:
		return r.TodoStyle
	case model.StatusInProgress:
		return r.InProgressStyle
	case model.StatusDone:
		return r.DoneStyle
	case model.StatusCancelled:
		return r.CancelledStyle
	default:
		return r.TodoStyle
	}
}

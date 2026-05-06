package components

import (
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/adem/taskbyte/internal/model"
)

func newTestRenderer() TaskListRenderer {
	return TaskListRenderer{
		TodoStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")),
		InProgressStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8C00")),
		DoneStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#2E8B57")).Strikethrough(true),
		CancelledStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Strikethrough(true),
		FocusedStyle:    lipgloss.NewStyle().Bold(true),
		ErrorStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")),
	}
}

func TestRenderTaskLine_Todo(t *testing.T) {
	r := newTestRenderer()
	task := model.Task{Text: "Buy groceries", Status: model.StatusTodo}

	output := r.RenderTaskLine(task, false)
	if output == "" {
		t.Error("todo task render should not be empty")
	}
}

func TestRenderTaskLine_InProgress(t *testing.T) {
	r := newTestRenderer()
	now := time.Now()
	task := model.Task{Text: "Working on it", Status: model.StatusInProgress, StatusChangedAt: &now}

	output := r.RenderTaskLine(task, false)
	if output == "" {
		t.Error("in_progress task render should not be empty")
	}
}

func TestRenderTaskLine_Done_Focused(t *testing.T) {
	r := newTestRenderer()
	now := time.Now()
	task := model.Task{Text: "Completed", Status: model.StatusDone, StatusChangedAt: &now}

	output := r.RenderTaskLine(task, true)
	if output == "" {
		t.Error("focused done task render should not be empty")
	}
}

func TestRenderTaskLine_Cancelled(t *testing.T) {
	r := newTestRenderer()
	task := model.Task{Text: "Cancelled task", Status: model.StatusCancelled}

	output := r.RenderTaskLine(task, false)
	if output == "" {
		t.Error("cancelled task render should not be empty")
	}
}

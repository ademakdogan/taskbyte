package app

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/adem/taskbyte/internal/config"
	"github.com/adem/taskbyte/internal/db"
	"github.com/adem/taskbyte/internal/service"
)

func newTestApp(t *testing.T) Model {
	t.Helper()
	d, err := db.NewInMemory()
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	t.Cleanup(func() { d.Close() })

	repo := db.NewRepository(d)
	svc := service.NewTaskService(repo)
	cfg := config.DefaultConfig()

	m := New(svc, cfg)
	m.width = 80
	m.height = 24
	return m
}

func TestNew_DefaultMode(t *testing.T) {
	m := newTestApp(t)
	if m.mode != ModeViewer {
		t.Errorf("default mode should be Viewer, got %d", m.mode)
	}
}

func TestNew_DefaultDate(t *testing.T) {
	m := newTestApp(t)
	if m.currentDate != service.TodayString() {
		t.Errorf("default date should be today, got %s", m.currentDate)
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	m := newTestApp(t)

	msg := tea.WindowSizeMsg{Width: 100, Height: 40}
	newM, _ := m.Update(msg)
	model := newM.(Model)

	if model.width != 100 || model.height != 40 {
		t.Errorf("expected 100x40, got %dx%d", model.width, model.height)
	}
}

func TestUpdate_TasksLoaded(t *testing.T) {
	m := newTestApp(t)

	// Simulate tasks loaded message
	msg := tasksLoadedMsg{tasks: nil, err: nil}
	newM, _ := m.Update(msg)
	model := newM.(Model)

	if model.err != nil {
		t.Errorf("unexpected error: %v", model.err)
	}
}

func TestUpdate_QuitCtrlC(t *testing.T) {
	m := newTestApp(t)

	key := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := m.Update(key)

	if cmd == nil {
		t.Error("ctrl+c should produce a quit command")
	}
}

func TestUpdate_ViewerToInsert(t *testing.T) {
	m := newTestApp(t)

	key := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}
	newM, _ := m.Update(key)
	model := newM.(Model)

	if model.mode != ModeInsert {
		t.Errorf("pressing 'i' should switch to Insert, got %d", model.mode)
	}
}

func TestUpdate_ViewerToSearch(t *testing.T) {
	m := newTestApp(t)

	key := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	newM, _ := m.Update(key)
	model := newM.(Model)

	if model.mode != ModeSearch {
		t.Errorf("pressing 's' should switch to Search, got %d", model.mode)
	}
}

func TestUpdate_InsertEscToViewer(t *testing.T) {
	m := newTestApp(t)
	m.mode = ModeInsert

	key := tea.KeyMsg{Type: tea.KeyEsc}
	newM, _ := m.Update(key)
	model := newM.(Model)

	if model.mode != ModeViewer {
		t.Errorf("Esc from Insert should return to Viewer, got %d", model.mode)
	}
}

func TestUpdate_ViewerToGoto(t *testing.T) {
	m := newTestApp(t)

	key := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	newM, _ := m.Update(key)
	model := newM.(Model)

	if model.mode != ModeGoto {
		t.Errorf("pressing 'g' should switch to Goto, got %d", model.mode)
	}
}

func TestView_EmptyTasks(t *testing.T) {
	m := newTestApp(t)
	output := m.View()

	if output == "" {
		t.Error("view should not be empty")
	}
}

func TestView_WithTasks(t *testing.T) {
	m := newTestApp(t)
	m.svc.AddTask("Test task", service.TodayString())

	// Load tasks
	cmd := m.loadTasks()
	msg := cmd()
	newM, _ := m.Update(msg)
	model := newM.(Model)

	output := model.View()
	if output == "" {
		t.Error("view with tasks should not be empty")
	}
}

func TestClearErrorMsg(t *testing.T) {
	m := newTestApp(t)
	m.err = fmt.Errorf("test error")

	newM, _ := m.Update(clearErrorMsg{})
	model := newM.(Model)

	if model.err != nil {
		t.Error("error should be cleared")
	}
}

func TestGhostText(t *testing.T) {
	m := newTestApp(t)
	m.mode = ModeInsert

	m.inputValue = "/"
	ghost := m.getGhostText()
	if ghost == "" {
		t.Error("ghost text for '/' should not be empty")
	}

	m.inputValue = "/d"
	ghost = m.getGhostText()
	if ghost == "" {
		t.Error("ghost text for '/d' should not be empty")
	}

	m.inputValue = "hello"
	ghost = m.getGhostText()
	if ghost != "" {
		t.Error("ghost text for non-slash input should be empty")
	}
}

func TestDeleteConfirm(t *testing.T) {
	m := newTestApp(t)
	m.svc.AddTask("Test task", service.TodayString())

	// Load tasks
	cmd := m.loadTasks()
	msg := cmd()
	newM, _ := m.Update(msg)
	m = newM.(Model)

	// Press r once (first press)
	key := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newM, _ = m.Update(key)
	m = newM.(Model)

	if !m.deleteConfirm {
		t.Error("first 'r' should set deleteConfirm")
	}

	// Press r again (confirms delete)
	newM, _ = m.Update(key)
	m = newM.(Model)

	if m.deleteConfirm {
		t.Error("second 'r' should clear deleteConfirm after delete")
	}
}

package service

import (
	"testing"

	"github.com/adem/taskbyte/internal/db"
	"github.com/adem/taskbyte/internal/model"
)

func newTestService(t *testing.T) *TaskService {
	t.Helper()
	d, err := db.NewInMemory()
	if err != nil {
		t.Fatalf("failed to create in-memory db: %v", err)
	}
	t.Cleanup(func() { d.Close() })
	repo := db.NewRepository(d)
	return NewTaskService(repo)
}

func TestAddTask(t *testing.T) {
	svc := newTestService(t)

	id, err := svc.AddTask("Buy groceries", "2026-03-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id <= 0 {
		t.Errorf("expected positive id, got %d", id)
	}
}

func TestAddTask_EmptyText(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.AddTask("", "2026-03-31")
	if err == nil {
		t.Error("expected error for empty text")
	}
}

func TestAddTask_WhitespaceOnly(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.AddTask("   ", "2026-03-31")
	if err == nil {
		t.Error("expected error for whitespace-only text")
	}
}

func TestCycleStatus(t *testing.T) {
	svc := newTestService(t)

	id, _ := svc.AddTask("Test task", "2026-03-31")

	// Todo -> InProgress
	status, err := svc.CycleStatus(id)
	if err != nil {
		t.Fatalf("cycle failed: %v", err)
	}
	if status != model.StatusInProgress {
		t.Errorf("expected in_progress, got %s", status)
	}

	// InProgress -> Done
	status, _ = svc.CycleStatus(id)
	if status != model.StatusDone {
		t.Errorf("expected done, got %s", status)
	}

	// Done -> Cancelled
	status, _ = svc.CycleStatus(id)
	if status != model.StatusCancelled {
		t.Errorf("expected cancelled, got %s", status)
	}

	// Cancelled -> Todo
	status, _ = svc.CycleStatus(id)
	if status != model.StatusTodo {
		t.Errorf("expected todo, got %s", status)
	}
}

func TestSetStatus(t *testing.T) {
	svc := newTestService(t)

	id, _ := svc.AddTask("Test", "2026-03-31")

	err := svc.SetStatus(id, model.StatusDone)
	if err != nil {
		t.Fatalf("set status failed: %v", err)
	}

	task, _ := svc.GetTaskByID(id)
	if task.Status != model.StatusDone {
		t.Errorf("expected done, got %s", task.Status)
	}
}

func TestEditTask(t *testing.T) {
	svc := newTestService(t)

	id, _ := svc.AddTask("Original", "2026-03-31")

	err := svc.EditTask(id, "Updated")
	if err != nil {
		t.Fatalf("edit failed: %v", err)
	}

	task, _ := svc.GetTaskByID(id)
	if task.Text != "Updated" {
		t.Errorf("expected 'Updated', got %q", task.Text)
	}
}

func TestEditTask_EmptyText(t *testing.T) {
	svc := newTestService(t)

	id, _ := svc.AddTask("Original", "2026-03-31")

	err := svc.EditTask(id, "")
	if err == nil {
		t.Error("expected error for empty text")
	}
}

func TestDeleteTask(t *testing.T) {
	svc := newTestService(t)

	id, _ := svc.AddTask("Delete me", "2026-03-31")

	err := svc.DeleteTask(id)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	tasks, _ := svc.GetTasksForDate("2026-03-31")
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestSearchTasks(t *testing.T) {
	svc := newTestService(t)

	svc.AddTask("Buy groceries", "2026-03-31")
	svc.AddTask("Buy coffee", "2026-04-01")
	svc.AddTask("Clean house", "2026-03-31")

	results, err := svc.SearchTasks("Buy")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestSearchTasks_EmptyQuery(t *testing.T) {
	svc := newTestService(t)

	results, err := svc.SearchTasks("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results != nil {
		t.Errorf("expected nil for empty query, got %v", results)
	}
}

func TestMigrateTasks(t *testing.T) {
	svc := newTestService(t)

	svc.AddTask("Task 1", "2026-03-30")
	id2, _ := svc.AddTask("Task 2", "2026-03-30")
	svc.SetStatus(id2, model.StatusDone)

	count, err := svc.MigrateTasks("2026-03-30")
	if err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 migrated, got %d", count)
	}
}

func TestMigrateTasks_SameDay(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.MigrateTasks(TodayString())
	if err == nil {
		t.Error("expected error when migrating from today to today")
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		input    string
		format   string
		expected string
		hasError bool
	}{
		{"31.03.2026", "DD.MM.YYYY", "2026-03-31", false},
		{"03.31.2026", "MM.DD.YYYY", "2026-03-31", false},
		{"2026-03-31", "YYYY-MM-DD", "2026-03-31", false},
		{"2026.03.31", "YYYY.MM.DD", "2026-03-31", false},
		{"invalid", "DD.MM.YYYY", "", true},
	}

	for _, tt := range tests {
		got, err := FormatDate(tt.input, tt.format)
		if tt.hasError && err == nil {
			t.Errorf("FormatDate(%q, %q) expected error", tt.input, tt.format)
		}
		if !tt.hasError && got != tt.expected {
			t.Errorf("FormatDate(%q, %q) = %q, want %q", tt.input, tt.format, got, tt.expected)
		}
	}
}

func TestStorageToDisplay(t *testing.T) {
	got := StorageToDisplay("2026-03-31", "DD.MM.YYYY")
	if got != "31.03.2026" {
		t.Errorf("got %q, want 31.03.2026", got)
	}

	got = StorageToDisplay("2026-03-31", "YYYY-MM-DD")
	if got != "2026-03-31" {
		t.Errorf("got %q, want 2026-03-31", got)
	}
}

func TestIsToday(t *testing.T) {
	if !IsToday(TodayString()) {
		t.Error("today should be today")
	}
	if IsToday("2020-01-01") {
		t.Error("2020-01-01 should not be today")
	}
}

func TestValidateDateInput(t *testing.T) {
	if !ValidateDateInput("31.03.2026") {
		t.Error("31.03.2026 should be valid")
	}
	if ValidateDateInput("2026-03-31") {
		t.Error("2026-03-31 should not match DD.MM.YYYY pattern")
	}
	if ValidateDateInput("invalid") {
		t.Error("invalid should not be valid")
	}
}

func TestGetDateStats(t *testing.T) {
	svc := newTestService(t)

	svc.AddTask("T1", "2026-03-31")
	id2, _ := svc.AddTask("T2", "2026-03-31")
	svc.SetStatus(id2, model.StatusDone)

	stats, err := svc.GetDateStats()
	if err != nil {
		t.Fatalf("get stats failed: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 date group, got %d", len(stats))
	}
	if stats[0].Total != 2 {
		t.Errorf("expected 2 total, got %d", stats[0].Total)
	}
}

package service

import (
	"testing"

	"github.com/adem/taskbyte/internal/db"
	"github.com/adem/taskbyte/internal/model"
)

func TestMigrateAllTasks(t *testing.T) {
	d, err := db.NewInMemory()
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer d.Close()

	repo := db.NewRepository(d)
	svc := NewTaskService(repo)

	svc.AddTask("Old task 1", "2026-01-01")
	svc.AddTask("Old task 2", "2026-02-01")
	id3, _ := svc.AddTask("Done task", "2026-01-15")
	svc.SetStatus(id3, model.StatusDone)

	today := TodayString()
	svc.AddTask("Today task", today)

	count, err := svc.MigrateAllTasks()
	if err != nil {
		t.Fatalf("migrate all failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 migrated, got %d", count)
	}

	// Verify tasks are now on today's date
	todayTasks, _ := svc.GetTasksForDate(today)
	todoCount := 0
	for _, task := range todayTasks {
		if task.Status == model.StatusTodo {
			todoCount++
		}
	}
	if todoCount < 2 {
		t.Errorf("expected at least 2 todo tasks today, got %d", todoCount)
	}
}

func TestMigrateTasks_FromSpecificDate(t *testing.T) {
	d, err := db.NewInMemory()
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer d.Close()

	repo := db.NewRepository(d)
	svc := NewTaskService(repo)

	svc.AddTask("Task A", "2026-03-15")
	svc.AddTask("Task B", "2026-03-15")
	id3, _ := svc.AddTask("Task C", "2026-03-15")
	svc.SetStatus(id3, model.StatusCancelled)

	count, err := svc.MigrateTasks("2026-03-15")
	if err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 migrated (1 cancelled stays), got %d", count)
	}
}

func TestMigrateTasks_EmptySource(t *testing.T) {
	d, err := db.NewInMemory()
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer d.Close()

	repo := db.NewRepository(d)
	svc := NewTaskService(repo)

	count, err := svc.MigrateTasks("2026-01-01")
	if err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 migrated from empty date, got %d", count)
	}
}

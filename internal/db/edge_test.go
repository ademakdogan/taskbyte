package db

import (
"testing"

"github.com/adem/taskbyte/internal/model"
)

func TestGetByID_NotFound(t *testing.T) {
	_, repo := newTestDB(t)

	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

func TestGetAll(t *testing.T) {
	_, repo := newTestDB(t)

	repo.Create("Task 1", "2026-03-31")
	repo.Create("Task 2", "2026-04-01")
	repo.Create("Task 3", "2026-03-30")

	tasks, err := repo.GetAll()
	if err != nil {
		t.Fatalf("get all failed: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestGetAll_Empty(t *testing.T) {
	_, repo := newTestDB(t)

	tasks, err := repo.GetAll()
	if err != nil {
		t.Fatalf("get all empty failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0, got %d", len(tasks))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	_, repo := newTestDB(t)

	err := repo.Delete(999)
	if err != nil {
		t.Errorf("deleting non-existent should not error: %v", err)
	}
}

func TestMultipleStatusChanges(t *testing.T) {
	_, repo := newTestDB(t)

	id, _ := repo.Create("Multi status", "2026-03-31")

	statuses := []model.Status{
		model.StatusInProgress,
		model.StatusDone,
		model.StatusCancelled,
		model.StatusTodo,
		model.StatusInProgress,
		model.StatusDone,
	}

	for _, s := range statuses {
		err := repo.UpdateStatus(id, s)
		if err != nil {
			t.Fatalf("update to %s failed: %v", s, err)
		}
		task, _ := repo.GetByID(id)
		if task.Status != s {
			t.Errorf("expected %s, got %s", s, task.Status)
		}
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	_, repo := newTestDB(t)

	repo.Create("Buy GROCERIES", "2026-03-31")

	// SQLite LIKE is case-insensitive for ASCII
	tasks, _ := repo.Search("groceries")
	if len(tasks) != 1 {
		t.Errorf("expected 1 result for case-insensitive search, got %d", len(tasks))
	}
}

func TestCreateMultipleSameDate(t *testing.T) {
	_, repo := newTestDB(t)

	for i := 0; i < 100; i++ {
		_, err := repo.Create("Task", "2026-03-31")
		if err != nil {
			t.Fatalf("create %d failed: %v", i, err)
		}
	}

	tasks, _ := repo.GetByDate("2026-03-31")
	if len(tasks) != 100 {
		t.Errorf("expected 100 tasks, got %d", len(tasks))
	}
}

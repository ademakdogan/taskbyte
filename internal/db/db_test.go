package db

import (
	"testing"

	"github.com/adem/taskbyte/internal/model"
)

func newTestDB(t *testing.T) (*DB, *Repository) {
	t.Helper()
	d, err := NewInMemory()
	if err != nil {
		t.Fatalf("failed to create in-memory db: %v", err)
	}
	t.Cleanup(func() { d.Close() })
	return d, NewRepository(d)
}

func TestNewInMemory(t *testing.T) {
	d, err := NewInMemory()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer d.Close()

	// Verify tables exist
	var name string
	err = d.Conn().QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tasks'").Scan(&name)
	if err != nil {
		t.Fatalf("tasks table not found: %v", err)
	}
	if name != "tasks" {
		t.Errorf("expected 'tasks', got %q", name)
	}
}

func TestCreate(t *testing.T) {
	_, repo := newTestDB(t)

	id, err := repo.Create("Buy groceries", "2026-03-31")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if id <= 0 {
		t.Errorf("expected positive id, got %d", id)
	}
}

func TestGetByDate(t *testing.T) {
	_, repo := newTestDB(t)

	repo.Create("Task A", "2026-03-31")
	repo.Create("Task B", "2026-03-31")
	repo.Create("Task C", "2026-04-01")

	tasks, err := repo.GetByDate("2026-03-31")
	if err != nil {
		t.Fatalf("get by date failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestGetByDate_Empty(t *testing.T) {
	_, repo := newTestDB(t)

	tasks, err := repo.GetByDate("2026-01-01")
	if err != nil {
		t.Fatalf("get by date failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestUpdateStatus(t *testing.T) {
	_, repo := newTestDB(t)

	id, _ := repo.Create("Test task", "2026-03-31")

	err := repo.UpdateStatus(id, model.StatusInProgress)
	if err != nil {
		t.Fatalf("update status failed: %v", err)
	}

	task, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if task.Status != model.StatusInProgress {
		t.Errorf("expected in_progress, got %s", task.Status)
	}
	if task.StatusChangedAt == nil {
		t.Error("StatusChangedAt should be set")
	}
}

func TestUpdateText(t *testing.T) {
	_, repo := newTestDB(t)

	id, _ := repo.Create("Original text", "2026-03-31")

	err := repo.UpdateText(id, "Updated text")
	if err != nil {
		t.Fatalf("update text failed: %v", err)
	}

	task, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if task.Text != "Updated text" {
		t.Errorf("expected 'Updated text', got %q", task.Text)
	}
}

func TestDelete(t *testing.T) {
	_, repo := newTestDB(t)

	id, _ := repo.Create("To delete", "2026-03-31")

	err := repo.Delete(id)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	tasks, _ := repo.GetByDate("2026-03-31")
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks after delete, got %d", len(tasks))
	}
}

func TestSearch(t *testing.T) {
	_, repo := newTestDB(t)

	repo.Create("Buy groceries", "2026-03-31")
	repo.Create("Buy coffee", "2026-04-01")
	repo.Create("Clean house", "2026-03-31")

	tasks, err := repo.Search("Buy")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 results, got %d", len(tasks))
	}
}

func TestSearch_NoResults(t *testing.T) {
	_, repo := newTestDB(t)

	repo.Create("Buy groceries", "2026-03-31")

	tasks, err := repo.Search("nonexistent")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 results, got %d", len(tasks))
	}
}

func TestGetUnfinishedByDate(t *testing.T) {
	_, repo := newTestDB(t)

	id1, _ := repo.Create("Todo task", "2026-03-31")
	id2, _ := repo.Create("In progress task", "2026-03-31")
	id3, _ := repo.Create("Done task", "2026-03-31")

	repo.UpdateStatus(id2, model.StatusInProgress)
	repo.UpdateStatus(id3, model.StatusDone)

	_ = id1
	tasks, err := repo.GetUnfinishedByDate("2026-03-31")
	if err != nil {
		t.Fatalf("get unfinished failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 unfinished, got %d", len(tasks))
	}
}

func TestGetAllUnfinished(t *testing.T) {
	_, repo := newTestDB(t)

	repo.Create("Task 1", "2026-03-31")
	repo.Create("Task 2", "2026-04-01")
	id3, _ := repo.Create("Task 3", "2026-03-30")
	repo.UpdateStatus(id3, model.StatusDone)

	tasks, err := repo.GetAllUnfinished()
	if err != nil {
		t.Fatalf("get all unfinished failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 unfinished, got %d", len(tasks))
	}
}

func TestGetDateStats(t *testing.T) {
	_, repo := newTestDB(t)

	repo.Create("T1", "2026-03-31")
	id2, _ := repo.Create("T2", "2026-03-31")
	id3, _ := repo.Create("T3", "2026-03-31")
	repo.Create("T4", "2026-04-01")

	repo.UpdateStatus(id2, model.StatusDone)
	repo.UpdateStatus(id3, model.StatusCancelled)

	stats, err := repo.GetDateStats()
	if err != nil {
		t.Fatalf("get stats failed: %v", err)
	}
	if len(stats) != 2 {
		t.Fatalf("expected 2 date groups, got %d", len(stats))
	}

	march := stats[0]
	if march.Total != 3 || march.Todo != 1 || march.Done != 1 || march.Cancelled != 1 {
		t.Errorf("unexpected stats for march: %+v", march)
	}
}

func TestUpdateDate(t *testing.T) {
	_, repo := newTestDB(t)

	id, _ := repo.Create("Migrate me", "2026-03-30")

	err := repo.UpdateDate(id, "2026-03-31")
	if err != nil {
		t.Fatalf("update date failed: %v", err)
	}

	task, _ := repo.GetByID(id)
	if task.Date != "2026-03-31" {
		t.Errorf("expected 2026-03-31, got %s", task.Date)
	}
}

func TestStatusTransitionCycle(t *testing.T) {
	_, repo := newTestDB(t)

	id, _ := repo.Create("Cycle test", "2026-03-31")

	statuses := []model.Status{
		model.StatusInProgress,
		model.StatusDone,
		model.StatusCancelled,
		model.StatusTodo,
	}

	for _, expected := range statuses {
		repo.UpdateStatus(id, expected)
		task, _ := repo.GetByID(id)
		if task.Status != expected {
			t.Errorf("expected %s, got %s", expected, task.Status)
		}
	}
}

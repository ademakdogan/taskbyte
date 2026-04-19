package db

import (
"database/sql"
"fmt"
"time"

"github.com/adem/taskbyte/internal/model"
)

// Repository provides CRUD operations for tasks.
type Repository struct {
	db *DB
}

// NewRepository creates a new task repository.
func NewRepository(db *DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new task and returns its ID.
func (r *Repository) Create(text, date string) (int, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := r.db.conn.Exec(
"INSERT INTO tasks (text, status, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
text, string(model.StatusTodo), date, now, now,
)
	if err != nil {
		return 0, fmt.Errorf("create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}
	return int(id), nil
}

// GetByDate returns all tasks for a given date (YYYY-MM-DD).
func (r *Repository) GetByDate(date string) ([]model.Task, error) {
	rows, err := r.db.conn.Query(
"SELECT id, text, status, date, created_at, updated_at, status_changed_at FROM tasks WHERE date = ? ORDER BY id",
date,
)
	if err != nil {
		return nil, fmt.Errorf("query tasks by date: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// GetAll returns all tasks ordered by date descending then id.
func (r *Repository) GetAll() ([]model.Task, error) {
	rows, err := r.db.conn.Query(
"SELECT id, text, status, date, created_at, updated_at, status_changed_at FROM tasks ORDER BY date DESC, id",
)
	if err != nil {
		return nil, fmt.Errorf("query all tasks: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// Search finds tasks whose text contains the query string.
func (r *Repository) Search(query string) ([]model.Task, error) {
	rows, err := r.db.conn.Query(
"SELECT id, text, status, date, created_at, updated_at, status_changed_at FROM tasks WHERE text LIKE ? ORDER BY date DESC, id",
"%"+query+"%",
)
	if err != nil {
		return nil, fmt.Errorf("search tasks: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// UpdateText updates the text of a task.
func (r *Repository) UpdateText(id int, newText string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.conn.Exec(
"UPDATE tasks SET text = ?, updated_at = ? WHERE id = ?",
newText, now, id,
)
	if err != nil {
		return fmt.Errorf("update task text: %w", err)
	}
	return nil
}

// UpdateStatus updates the status of a task and records the timestamp.
func (r *Repository) UpdateStatus(id int, status model.Status) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.conn.Exec(
"UPDATE tasks SET status = ?, status_changed_at = ?, updated_at = ? WHERE id = ?",
string(status), now, now, id,
)
	if err != nil {
		return fmt.Errorf("update task status: %w", err)
	}
	return nil
}

// UpdateDate moves a task to a different date.
func (r *Repository) UpdateDate(id int, newDate string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.conn.Exec(
"UPDATE tasks SET date = ?, updated_at = ? WHERE id = ?",
newDate, now, id,
)
	if err != nil {
		return fmt.Errorf("update task date: %w", err)
	}
	return nil
}

// Delete removes a task by ID.
func (r *Repository) Delete(id int) error {
	_, err := r.db.conn.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

// GetUnfinishedByDate returns todo and in_progress tasks for a given date.
func (r *Repository) GetUnfinishedByDate(date string) ([]model.Task, error) {
	rows, err := r.db.conn.Query(
"SELECT id, text, status, date, created_at, updated_at, status_changed_at FROM tasks WHERE date = ? AND status IN (?, ?) ORDER BY id",
date, string(model.StatusTodo), string(model.StatusInProgress),
)
	if err != nil {
		return nil, fmt.Errorf("query unfinished tasks: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// GetAllUnfinished returns all todo and in_progress tasks across all dates.
func (r *Repository) GetAllUnfinished() ([]model.Task, error) {
	rows, err := r.db.conn.Query(
"SELECT id, text, status, date, created_at, updated_at, status_changed_at FROM tasks WHERE status IN (?, ?) ORDER BY date, id",
string(model.StatusTodo), string(model.StatusInProgress),
)
	if err != nil {
		return nil, fmt.Errorf("query all unfinished tasks: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// DateStats holds statistics for a single date.
type DateStats struct {
	Date       string
	Total      int
	Todo       int
	InProgress int
	Done       int
	Cancelled  int
}

// GetDateStats returns task statistics grouped by date.
func (r *Repository) GetDateStats() ([]DateStats, error) {
	rows, err := r.db.conn.Query(`
		SELECT
			date,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'todo' THEN 1 ELSE 0 END) as todo,
			SUM(CASE WHEN status = 'in_progress' THEN 1 ELSE 0 END) as in_progress,
			SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END) as done,
			SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) as cancelled
		FROM tasks
		GROUP BY date
		ORDER BY date
	`)
	if err != nil {
		return nil, fmt.Errorf("query date stats: %w", err)
	}
	defer rows.Close()

	var stats []DateStats
	for rows.Next() {
		var s DateStats
		if err := rows.Scan(&s.Date, &s.Total, &s.Todo, &s.InProgress, &s.Done, &s.Cancelled); err != nil {
			return nil, fmt.Errorf("scan date stats: %w", err)
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

// GetByID returns a single task by its ID.
func (r *Repository) GetByID(id int) (*model.Task, error) {
	row := r.db.conn.QueryRow(
"SELECT id, text, status, date, created_at, updated_at, status_changed_at FROM tasks WHERE id = ?",
id,
)

	t, err := scanTask(row)
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	return t, nil
}

func scanTasks(rows *sql.Rows) ([]model.Task, error) {
	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		var statusStr string
		var createdAt, updatedAt string
		var statusChangedAt sql.NullString

		if err := rows.Scan(&t.ID, &t.Text, &statusStr, &t.Date, &createdAt, &updatedAt, &statusChangedAt); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		t.Status = model.Status(statusStr)
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		if statusChangedAt.Valid {
			parsed, _ := time.Parse(time.RFC3339, statusChangedAt.String)
			t.StatusChangedAt = &parsed
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func scanTask(row *sql.Row) (*model.Task, error) {
	var t model.Task
	var statusStr string
	var createdAt, updatedAt string
	var statusChangedAt sql.NullString

	if err := row.Scan(&t.ID, &t.Text, &statusStr, &t.Date, &createdAt, &updatedAt, &statusChangedAt); err != nil {
		return nil, err
	}

	t.Status = model.Status(statusStr)
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	if statusChangedAt.Valid {
		parsed, _ := time.Parse(time.RFC3339, statusChangedAt.String)
		t.StatusChangedAt = &parsed
	}
	return &t, nil
}

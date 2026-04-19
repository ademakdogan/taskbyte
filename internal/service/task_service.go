package service

import (
"fmt"
"regexp"
"strings"
"time"

"github.com/adem/taskbyte/internal/db"
"github.com/adem/taskbyte/internal/model"
)

// TaskService provides business logic for task operations.
type TaskService struct {
	repo *db.Repository
}

// NewTaskService creates a new task service.
func NewTaskService(repo *db.Repository) *TaskService {
	return &TaskService{repo: repo}
}

// AddTask creates a new task for the given date.
func (s *TaskService) AddTask(text, date string) (int, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0, fmt.Errorf("task text cannot be empty")
	}
	return s.repo.Create(text, date)
}

// GetTasksForDate returns all tasks for a date.
func (s *TaskService) GetTasksForDate(date string) ([]model.Task, error) {
	return s.repo.GetByDate(date)
}

// CycleStatus advances the task to the next status in the cycle.
func (s *TaskService) CycleStatus(taskID int) (model.Status, error) {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return "", fmt.Errorf("get task: %w", err)
	}

	next := task.Status.NextStatus()
	if err := s.repo.UpdateStatus(taskID, next); err != nil {
		return "", err
	}
	return next, nil
}

// SetStatus sets a task to a specific status directly.
func (s *TaskService) SetStatus(taskID int, status model.Status) error {
	return s.repo.UpdateStatus(taskID, status)
}

// EditTask updates the text of a task.
func (s *TaskService) EditTask(taskID int, newText string) error {
	newText = strings.TrimSpace(newText)
	if newText == "" {
		return fmt.Errorf("task text cannot be empty")
	}
	return s.repo.UpdateText(taskID, newText)
}

// DeleteTask removes a task.
func (s *TaskService) DeleteTask(taskID int) error {
	return s.repo.Delete(taskID)
}

// SearchTasks finds tasks matching the query string.
func (s *TaskService) SearchTasks(query string) ([]model.Task, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}
	return s.repo.Search(query)
}

// MigrateTasks moves unfinished tasks from a source date to today.
func (s *TaskService) MigrateTasks(sourceDate string) (int, error) {
	today := TodayString()
	if sourceDate == today {
		return 0, fmt.Errorf("cannot migrate from today to today")
	}

	tasks, err := s.repo.GetUnfinishedByDate(sourceDate)
	if err != nil {
		return 0, err
	}

	for _, t := range tasks {
		if err := s.repo.UpdateDate(t.ID, today); err != nil {
			return 0, fmt.Errorf("migrate task %d: %w", t.ID, err)
		}
	}
	return len(tasks), nil
}

// MigrateAllTasks moves all unfinished tasks from all dates to today.
func (s *TaskService) MigrateAllTasks() (int, error) {
	today := TodayString()
	tasks, err := s.repo.GetAllUnfinished()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, t := range tasks {
		if t.Date == today {
			continue
		}
		if err := s.repo.UpdateDate(t.ID, today); err != nil {
			return 0, fmt.Errorf("migrate task %d: %w", t.ID, err)
		}
		count++
	}
	return count, nil
}

// GetDateStats returns statistics grouped by date.
func (s *TaskService) GetDateStats() ([]db.DateStats, error) {
	return s.repo.GetDateStats()
}

// GetTaskByID returns a task by its ID.
func (s *TaskService) GetTaskByID(id int) (*model.Task, error) {
	return s.repo.GetByID(id)
}

// TodayString returns today's date as YYYY-MM-DD.
func TodayString() string {
return time.Now().Format("2006-01-02")
}

// FormatDate converts a display date (DD.MM.YYYY) to storage format (YYYY-MM-DD).
func FormatDate(displayDate, format string) (string, error) {
var goFormat string
switch format {
case "DD.MM.YYYY":
goFormat = "02.01.2006"
case "MM.DD.YYYY":
goFormat = "01.02.2006"
case "YYYY-MM-DD":
goFormat = "2006-01-02"
case "YYYY.MM.DD":
goFormat = "2006.01.02"
default:
goFormat = "02.01.2006"
}

t, err := time.Parse(goFormat, displayDate)
if err != nil {
return "", fmt.Errorf("invalid date %q for format %q: %w", displayDate, format, err)
}
return t.Format("2006-01-02"), nil
}

// StorageToDisplay converts a storage date (YYYY-MM-DD) to the user's display format.
func StorageToDisplay(storageDate, format string) string {
	t, err := time.Parse("2006-01-02", storageDate)
	if err != nil {
		return storageDate
	}

	switch format {
	case "DD.MM.YYYY":
		return t.Format("02.01.2006")
	case "MM.DD.YYYY":
		return t.Format("01.02.2006")
	case "YYYY-MM-DD":
		return t.Format("2006-01-02")
	case "YYYY.MM.DD":
		return t.Format("2006.01.02")
	default:
		return t.Format("02.01.2006")
	}
}

// IsToday checks if a storage-format date is today.
func IsToday(storageDate string) bool {
	return storageDate == TodayString()
}

// YesterdayString returns yesterday's date as YYYY-MM-DD.
func YesterdayString() string {
return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
}

// ValidateDateInput checks if a string looks like a valid date.
var datePattern = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)

func ValidateDateInput(input string) bool {
return datePattern.MatchString(input)
}

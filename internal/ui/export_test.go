package ui

import (
"os"
"path/filepath"
"testing"

"github.com/adem/taskbyte/internal/model"
)

func TestExportToMarkdown(t *testing.T) {
	tmpDir := t.TempDir()

	tasks := []model.Task{
		{ID: 1, Text: "Buy groceries", Status: model.StatusTodo, Date: "2026-03-31"},
		{ID: 2, Text: "Clean house", Status: model.StatusDone, Date: "2026-03-31"},
	}

	path, err := ExportToMarkdown(tasks, "2026-03-31", "DD.MM.YYYY", tmpDir)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("export file not found: %v", err)
	}

	// Verify content
	data, _ := os.ReadFile(path)
	content := string(data)
	if content == "" {
		t.Error("exported file should not be empty")
	}

	// Verify extension
	ext := filepath.Ext(path)
	if ext != ".md" {
		t.Errorf("expected .md extension, got %q", ext)
	}
}

func TestExportToMarkdown_DefaultPath(t *testing.T) {
	tasks := []model.Task{
		{ID: 1, Text: "Test task", Status: model.StatusTodo, Date: "2026-03-31"},
	}

	path, err := ExportToMarkdown(tasks, "2026-03-31", "DD.MM.YYYY", "")
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Clean up
	defer os.Remove(path)

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("default path export file not found: %v", err)
	}
}

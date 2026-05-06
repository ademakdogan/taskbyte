package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adem/taskbyte/internal/model"
	"github.com/adem/taskbyte/internal/service"
)

// ExportToMarkdown exports tasks to a markdown file.
func ExportToMarkdown(tasks []model.Task, date, dateFormat, exportPath string) (string, error) {
	if exportPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		exportPath = home
	}

	// Ensure directory exists
	if err := os.MkdirAll(exportPath, 0o755); err != nil {
		return "", fmt.Errorf("create export dir: %w", err)
	}

	// Generate filename
	dateDisplay := service.StorageToDisplay(date, dateFormat)
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("todo_%s_%s.md", strings.ReplaceAll(dateDisplay, ".", "-"), timestamp)
	fullPath := filepath.Join(exportPath, filename)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Todo List - %s\n\n", dateDisplay))
	sb.WriteString(fmt.Sprintf("Exported: %s\n\n", time.Now().Format("02.01.2006 15:04:05")))

	for _, task := range tasks {
		checkbox := "- [ ]"
		switch task.Status {
		case model.StatusDone:
			checkbox = "- [x]"
		case model.StatusInProgress:
			checkbox = "- [~]"
		case model.StatusCancelled:
			checkbox = "- [!]"
		}

		line := fmt.Sprintf("%s %s", checkbox, task.Text)
		if task.Status.Label() != "" {
			line += " " + task.Status.Label()
		}
		if task.StatusChangedAt != nil {
			line += fmt.Sprintf(" (%s)", task.StatusChangedAt.Local().Format("02.01.2006 15:04"))
		}
		sb.WriteString(line + "\n")
	}

	if err := os.WriteFile(fullPath, []byte(sb.String()), 0o644); err != nil {
		return "", fmt.Errorf("write export file: %w", err)
	}

	return fullPath, nil
}

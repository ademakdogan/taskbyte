package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/adem/taskbyte/internal/app"
	"github.com/adem/taskbyte/internal/config"
	"github.com/adem/taskbyte/internal/db"
	"github.com/adem/taskbyte/internal/service"
)

func main() {
	// Ensure data directory exists
	_, err := config.EnsureDataDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Open database
	dbPath, err := config.DBPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving database path: %v\n", err)
		os.Exit(1)
	}

	database, err := db.New(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Create service
	repo := db.NewRepository(database)
	svc := service.NewTaskService(repo)

	// Run TUI
	model := app.New(svc, cfg)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupTestEnv(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	// Reset global state
	mu.Lock()
	loaded = false
	current = Config{}
	mu.Unlock()

	return tmpDir
}

func TestLoad_CreatesDefault(t *testing.T) {
	tmpDir := setupTestEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := DefaultConfig()
	if cfg != expected {
		t.Errorf("got %+v, want %+v", cfg, expected)
	}

	// Verify file was written
	path := filepath.Join(tmpDir, appName, "config.json")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("config file not created: %v", err)
	}
}

func TestLoad_ReadsExisting(t *testing.T) {
	tmpDir := setupTestEnv(t)

	// Create custom config
	dir := filepath.Join(tmpDir, appName)
	os.MkdirAll(dir, 0o755)

	custom := Config{
		AutoHideCompleted:   true,
		InsertPromptHistory: false,
		DateFormat:          "YYYY-MM-DD",
		Theme: ThemeConfig{
			TodoColor:       "cyan",
			InProgressColor: "yellow",
			DoneColor:       "light_green",
			CancelledColor:  "magenta",
		},
	}
	data, _ := json.MarshalIndent(custom, "", "  ")
	os.WriteFile(filepath.Join(dir, "config.json"), data, 0o644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg != custom {
		t.Errorf("got %+v, want %+v", cfg, custom)
	}
}

func TestLoad_CorruptJSON_ResetsToDefault(t *testing.T) {
	tmpDir := setupTestEnv(t)

	dir := filepath.Join(tmpDir, appName)
	os.MkdirAll(dir, 0o755)
	os.WriteFile(filepath.Join(dir, "config.json"), []byte("{invalid json"), 0o644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := DefaultConfig()
	if cfg != expected {
		t.Errorf("corrupt config should reset to default, got %+v", cfg)
	}
}

func TestSave_PersistsChanges(t *testing.T) {
	tmpDir := setupTestEnv(t)

	// Load defaults first
	_, err := Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	// Modify and save
	modified := DefaultConfig()
	modified.AutoHideCompleted = true
	modified.DateFormat = "YYYY-MM-DD"

	if err := Save(modified); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Reload from disk
	mu.Lock()
	loaded = false
	mu.Unlock()

	_ = tmpDir // keep env set
	cfg, err := Load()
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}

	if cfg.AutoHideCompleted != true {
		t.Error("AutoHideCompleted not persisted")
	}
	if cfg.DateFormat != "YYYY-MM-DD" {
		t.Errorf("DateFormat not persisted, got %q", cfg.DateFormat)
	}
}

func TestGet_BeforeLoad_ReturnsDefault(t *testing.T) {
	mu.Lock()
	loaded = false
	current = Config{}
	mu.Unlock()

	cfg := Get()
	expected := DefaultConfig()
	if cfg != expected {
		t.Errorf("Get before Load should return default, got %+v", cfg)
	}
}

func TestDefaultConfig_Values(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.AutoHideCompleted != false {
		t.Error("AutoHideCompleted should default to false")
	}
	if cfg.InsertPromptHistory != true {
		t.Error("InsertPromptHistory should default to true")
	}
	if cfg.DateFormat != "DD.MM.YYYY" {
		t.Errorf("DateFormat should default to DD.MM.YYYY, got %q", cfg.DateFormat)
	}
	if cfg.Theme.TodoColor != "white" {
		t.Errorf("TodoColor should default to white, got %q", cfg.Theme.TodoColor)
	}
	if cfg.Theme.InProgressColor != "orange" {
		t.Errorf("InProgressColor should default to orange, got %q", cfg.Theme.InProgressColor)
	}
	if cfg.Theme.DoneColor != "dark_green" {
		t.Errorf("DoneColor should default to dark_green, got %q", cfg.Theme.DoneColor)
	}
	if cfg.Theme.CancelledColor != "red" {
		t.Errorf("CancelledColor should default to red, got %q", cfg.Theme.CancelledColor)
	}
}

func TestValidDateFormats(t *testing.T) {
	formats := ValidDateFormats()
	if len(formats) < 3 {
		t.Errorf("expected at least 3 date formats, got %d", len(formats))
	}
}

func TestValidColors(t *testing.T) {
	colors := ValidColors()
	if len(colors) < 5 {
		t.Errorf("expected at least 5 colors, got %d", len(colors))
	}
}

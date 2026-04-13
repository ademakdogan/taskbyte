package config

import (
"os"
"path/filepath"
"runtime"
"testing"
)

func TestDataDir_Default(t *testing.T) {
	// Clear any XDG override
	t.Setenv("XDG_DATA_HOME", "")

	dir, err := DataDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	home, _ := os.UserHomeDir()
	var expected string
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			expected = filepath.Join(appData, appName)
		} else {
			expected = filepath.Join(home, "AppData", "Roaming", appName)
		}
	default:
		expected = filepath.Join(home, ".local", "share", appName)
	}

	if dir != expected {
		t.Errorf("got %q, want %q", dir, expected)
	}
}

func TestDataDir_XDGOverride(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	dir, err := DataDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(tmpDir, appName)
	if dir != expected {
		t.Errorf("got %q, want %q", dir, expected)
	}
}

func TestEnsureDataDir_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	dir, err := EnsureDataDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("expected directory, got file")
	}
}

func TestEnsureDataDir_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	dir1, err := EnsureDataDir()
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}

	dir2, err := EnsureDataDir()
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}

	if dir1 != dir2 {
		t.Errorf("paths differ: %q vs %q", dir1, dir2)
	}
}

func TestDBPath(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	dbPath, err := DBPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(tmpDir, appName, "taskbyte.db")
	if dbPath != expected {
		t.Errorf("got %q, want %q", dbPath, expected)
	}
}

func TestConfigPath(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	cfgPath, err := ConfigPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(tmpDir, appName, "config.json")
	if cfgPath != expected {
		t.Errorf("got %q, want %q", cfgPath, expected)
	}
}

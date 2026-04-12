package config

import (
"os"
"path/filepath"
"runtime"
)

const appName = "taskbyte"

// DataDir returns the XDG-compliant data directory for TaskByte.
// On Linux: $XDG_DATA_HOME/taskbyte or ~/.local/share/taskbyte
// On macOS: $XDG_DATA_HOME/taskbyte or ~/.local/share/taskbyte
// On Windows: %APPDATA%/taskbyte
func DataDir() (string, error) {
	var base string

	if env := os.Getenv("XDG_DATA_HOME"); env != "" {
		base = env
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		switch runtime.GOOS {
		case "windows":
			if appData := os.Getenv("APPDATA"); appData != "" {
				base = appData
			} else {
				base = filepath.Join(home, "AppData", "Roaming")
			}
		default: // linux, darwin
			base = filepath.Join(home, ".local", "share")
		}
	}

	dir := filepath.Join(base, appName)
	return dir, nil
}

// EnsureDataDir creates the data directory if it does not exist.
func EnsureDataDir() (string, error) {
	dir, err := DataDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	return dir, nil
}

// DBPath returns the full path to the SQLite database file.
func DBPath() (string, error) {
	dir, err := EnsureDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "taskbyte.db"), nil
}

// ConfigPath returns the full path to the config JSON file.
func ConfigPath() (string, error) {
	dir, err := EnsureDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

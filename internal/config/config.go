package config

import (
"encoding/json"
"errors"
"os"
"sync"
)

// ThemeConfig holds color settings for task statuses.
type ThemeConfig struct {
	TodoColor       string `json:"todo_color"`
	InProgressColor string `json:"in_progress_color"`
	DoneColor       string `json:"done_color"`
	CancelledColor  string `json:"cancelled_color"`
}

// Config holds all application settings.
type Config struct {
	AutoHideCompleted   bool        `json:"auto_hide_completed"`
	InsertPromptHistory bool        `json:"insert_prompt_history"`
	DateFormat          string      `json:"date_format"`
	Theme               ThemeConfig `json:"theme"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		AutoHideCompleted:   false,
		InsertPromptHistory: true,
		DateFormat:          "DD.MM.YYYY",
		Theme: ThemeConfig{
			TodoColor:       "white",
			InProgressColor: "orange",
			DoneColor:       "dark_green",
			CancelledColor:  "red",
		},
	}
}

// ValidDateFormats returns all supported date format strings.
func ValidDateFormats() []string {
	return []string{"DD.MM.YYYY", "MM.DD.YYYY", "YYYY-MM-DD", "YYYY.MM.DD"}
}

// ValidColors returns all supported color names.
func ValidColors() []string {
	return []string{"white", "orange", "red", "dark_green", "light_green", "cyan", "blue", "gray", "yellow", "magenta"}
}

var (
current Config
mu      sync.RWMutex
loaded  bool
)

// Load reads the config from disk, or creates a default one if it doesn't exist.
func Load() (Config, error) {
mu.Lock()
defer mu.Unlock()

path, err := ConfigPath()
if err != nil {
return Config{}, err
}

data, err := os.ReadFile(path)
if err != nil {
if errors.Is(err, os.ErrNotExist) {
cfg := DefaultConfig()
if writeErr := writeConfig(path, cfg); writeErr != nil {
return Config{}, writeErr
}
current = cfg
loaded = true
return cfg, nil
}
return Config{}, err
}

var cfg Config
if err := json.Unmarshal(data, &cfg); err != nil {
// If JSON is corrupt, reset to defaults
cfg = DefaultConfig()
if writeErr := writeConfig(path, cfg); writeErr != nil {
return Config{}, writeErr
}
}

current = cfg
loaded = true
return cfg, nil
}

// Save persists the given config to disk.
func Save(cfg Config) error {
mu.Lock()
defer mu.Unlock()

path, err := ConfigPath()
if err != nil {
return err
}

current = cfg
return writeConfig(path, cfg)
}

// Get returns the current in-memory config. Must call Load first.
func Get() Config {
mu.RLock()
defer mu.RUnlock()

if !loaded {
return DefaultConfig()
}
return current
}

func writeConfig(path string, cfg Config) error {
data, err := json.MarshalIndent(cfg, "", "  ")
if err != nil {
return err
}
return os.WriteFile(path, data, 0o644)
}

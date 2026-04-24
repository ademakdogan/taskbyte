package ui

import (
"testing"

"github.com/adem/taskbyte/internal/config"
)

func TestBuildSettingsItems(t *testing.T) {
	cfg := config.DefaultConfig()
	items := BuildSettingsItems(cfg)

	if len(items) != 7 {
		t.Errorf("expected 7 settings items, got %d", len(items))
	}

	// Verify first item
	if items[0].Label != "Auto-hide Completed Tasks" {
		t.Errorf("first item should be Auto-hide, got %q", items[0].Label)
	}
	if items[0].Value != "False" {
		t.Errorf("auto-hide should default to False, got %q", items[0].Value)
	}
}

func TestApplySettingsItems(t *testing.T) {
	cfg := config.DefaultConfig()
	items := BuildSettingsItems(cfg)

	// Modify some values
	items[0].Value = "True" // Auto-hide
	items[2].Value = "YYYY-MM-DD" // Date format
	items[5].Value = "light_green" // Done color

	newCfg := ApplySettingsItems(items)

	if !newCfg.AutoHideCompleted {
		t.Error("AutoHideCompleted should be true")
	}
	if newCfg.DateFormat != "YYYY-MM-DD" {
		t.Errorf("DateFormat should be YYYY-MM-DD, got %q", newCfg.DateFormat)
	}
	if newCfg.Theme.DoneColor != "light_green" {
		t.Errorf("DoneColor should be light_green, got %q", newCfg.Theme.DoneColor)
	}
}

func TestRenderSettings(t *testing.T) {
	cfg := config.DefaultConfig()
	items := BuildSettingsItems(cfg)
	styles := NewStyles(cfg)

	output := RenderSettings(items, 0, false, 0, styles)
	if output == "" {
		t.Error("settings render should not be empty")
	}
}

func TestRenderSettings_WithDropdown(t *testing.T) {
	cfg := config.DefaultConfig()
	items := BuildSettingsItems(cfg)
	styles := NewStyles(cfg)

	// Dropdown open on Done Task Color (index 5)
	output := RenderSettings(items, 5, true, 0, styles)
	if output == "" {
		t.Error("settings render with dropdown should not be empty")
	}
}

func TestBoolStr(t *testing.T) {
	if boolStr(true) != "True" {
		t.Error("true should return 'True'")
	}
	if boolStr(false) != "False" {
		t.Error("false should return 'False'")
	}
}

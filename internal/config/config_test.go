package config

import (
	"path/filepath"
	"testing"
)

func TestGetDefaultConfig(t *testing.T) {
	cfg := GetDefaultConfig()

	if cfg.Global.LogLevel != "info" {
		t.Errorf("expected loglevel 'info', got '%s'", cfg.Global.LogLevel)
	}
	if cfg.Global.NoColor != false {
		t.Errorf("expected nocolor false, got %v", cfg.Global.NoColor)
	}
	if cfg.Global.AutoUpdate != true {
		t.Errorf("expected autoupdate true, got %v", cfg.Global.AutoUpdate)
	}
	if cfg.Speedtest.Timeout != 30 {
		t.Errorf("expected timeout 30, got %d", cfg.Speedtest.Timeout)
	}
	if cfg.UI.HistoryLimit != 100 {
		t.Errorf("expected historylimit 100, got %d", cfg.UI.HistoryLimit)
	}
	if cfg.UI.SaveHistory != true {
		t.Errorf("expected savehistory true, got %v", cfg.UI.SaveHistory)
	}
}

func TestGet(t *testing.T) {
	cfg := Get()
	if cfg == nil {
		t.Error("expected config, got nil")
	}
}

func TestGetConfigDir(t *testing.T) {
	if err := Init("auraspeed"); err != nil {
		t.Fatalf("failed to init config: %v", err)
	}
	dir := GetConfigDir()
	if dir == "" {
		t.Error("expected non-empty config dir")
	}
}

func TestGetDataDir(t *testing.T) {
	if err := Init("auraspeed"); err != nil {
		t.Fatalf("failed to init config: %v", err)
	}
	dir := GetDataDir()
	if dir == "" {
		t.Error("expected non-empty data dir")
	}
}

func TestGetHistoryFile(t *testing.T) {
	if err := Init("auraspeed"); err != nil {
		t.Fatalf("failed to init config: %v", err)
	}
	file := GetHistoryFile()
	if file == "" {
		t.Error("expected non-empty history file path")
	}
	if filepath.Ext(file) != ".json" {
		t.Errorf("expected .json extension, got %s", filepath.Ext(file))
	}
}

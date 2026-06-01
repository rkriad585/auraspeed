package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestConfigAliases(t *testing.T) {
	cfg := GetDefaultConfig()

	if cfg.Aliases["st"] != "speedtest" {
		t.Errorf("expected alias 'st' -> 'speedtest', got '%s'", cfg.Aliases["st"])
	}
	if cfg.Aliases["si"] != "info" {
		t.Errorf("expected alias 'si' -> 'info', got '%s'", cfg.Aliases["si"])
	}
	if cfg.Aliases["net"] != "network" {
		t.Errorf("expected alias 'net' -> 'network', got '%s'", cfg.Aliases["net"])
	}
	if cfg.Aliases["hist"] != "history" {
		t.Errorf("expected alias 'hist' -> 'history', got '%s'", cfg.Aliases["hist"])
	}
}

func TestSpeedtestConfig(t *testing.T) {
	cfg := GetDefaultConfig()

	if cfg.Speedtest.DefaultServerID != 0 {
		t.Errorf("expected defaultserverid 0, got %d", cfg.Speedtest.DefaultServerID)
	}
	if cfg.Speedtest.ParallelDownloads != 4 {
		t.Errorf("expected paralleldownloads 4, got %d", cfg.Speedtest.ParallelDownloads)
	}
	if cfg.Speedtest.ParallelUploads != 2 {
		t.Errorf("expected paralleluploads 2, got %d", cfg.Speedtest.ParallelUploads)
	}
}

func TestUIConfig(t *testing.T) {
	cfg := GetDefaultConfig()

	if cfg.UI.Theme != "sunny-beach" {
		t.Errorf("expected theme 'sunny-beach', got '%s'", cfg.UI.Theme)
	}
	if cfg.UI.GraphHeight != 8 {
		t.Errorf("expected graphheight 8, got %d", cfg.UI.GraphHeight)
	}
	if cfg.UI.AutoRefresh != false {
		t.Errorf("expected autorefresh false, got %v", cfg.UI.AutoRefresh)
	}
	if cfg.UI.RefreshRate != 5 {
		t.Errorf("expected refreshrate 5, got %d", cfg.UI.RefreshRate)
	}
}

func TestConfigDirCreation(t *testing.T) {
	testAppName := "auraspeed-test-" + time.Now().Format("20060102150405")

	err := Init(testAppName)
	if err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	dir := GetConfigDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("config directory was not created: %s", dir)
	}

	dataDir := GetDataDir()
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Errorf("data directory was not created: %s", dataDir)
	}
}

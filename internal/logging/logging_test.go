package logging

import (
	"testing"

	"github.com/rs/zerolog"
)

func resetState() {
	globalLogger = nil
	globalNoColor = false
	currentLevel = zerolog.InfoLevel
	privacyMode = false
}

func TestNew(t *testing.T) {
	resetState()

	logger := New()
	if logger == nil {
		t.Error("expected logger, got nil")
	}
}

func TestSetLevel(t *testing.T) {
	resetState()
	New() // ensure initialized

	// Test valid level
	err := SetLevel("debug")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test invalid level
	err = SetLevel("invalid")
	if err == nil {
		t.Error("expected error for invalid level, got nil")
	}
}

func TestSetNoColor(t *testing.T) {
	resetState()
	New() // ensure initialized

	SetNoColor(true)
	if !globalNoColor {
		t.Error("expected globalNoColor to be true")
	}

	SetNoColor(false)
	if globalNoColor {
		t.Error("expected globalNoColor to be false")
	}
}

func TestGet(t *testing.T) {
	resetState()

	logger1 := Get()
	if logger1 == nil {
		t.Error("expected logger from Get(), got nil")
	}

	logger2 := Get()
	if logger1 != logger2 {
		t.Error("expected Get() to return the same instance (singleton)")
	}
}

package speedtest

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadHistoryEmpty(t *testing.T) {
	// Non-existent file
	results, err := loadHistoryForTest(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestLoadHistoryJSONArray(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "history.json")

	data := `[{"timestamp":"12:00:00","download":100.5,"upload":20.5,"ping":10,"isp":"Test"}]`
	if err := os.WriteFile(tmpFile, []byte(data), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	results, err := loadHistoryForTest(tmpFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Download != 100.5 {
		t.Errorf("expected download 100.5, got %f", results[0].Download)
	}
}

func TestLoadHistoryOldFormat(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "history.json")

	data := `{"timestamp":"12:00:00","download":100.5,"upload":20.5,"ping":10,"isp":"Test"}
{"timestamp":"12:01:00","download":200.5,"upload":30.5,"ping":15,"isp":"Test2"}`
	if err := os.WriteFile(tmpFile, []byte(data), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	results, err := loadHistoryForTest(tmpFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

// Helper to set history file path for testing
func loadHistoryForTest(path string) ([]TestResult, error) {
	// Read file and parse
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []TestResult{}, nil
		}
		return nil, err
	}
	defer file.Close()

	// Try JSON array first
	var results []TestResult
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&results); err == nil {
		return results, nil
	}

	// Fallback to JSON lines
	file.Seek(0, 0)
	results = []TestResult{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var res TestResult
		if err := json.Unmarshal([]byte(scanner.Text()), &res); err == nil {
			results = append(results, res)
		}
	}
	return results, scanner.Err()
}

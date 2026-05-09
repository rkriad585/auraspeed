package speedtest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"auraspeed/internal/config"
	"auraspeed/internal/logging"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TestResult represents a single speed test result.
type TestResult struct {
	Timestamp string  `json:"timestamp"`
	Download  float64 `json:"download"`
	Upload    float64 `json:"upload"`
	Ping      int64   `json:"ping"`
	ISP       string  `json:"isp"`
}

// LoadHistory loads test history from file with backward compatibility.
// It supports both the new JSON array format and the old JSON lines format.
// Returns an empty slice if the file doesn't exist.
func LoadHistory() ([]TestResult, error) {
	historyPath := config.GetHistoryFile()
	file, err := os.Open(historyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []TestResult{}, nil
		}
		return nil, err
	}
	defer file.Close()

	// Try to read as JSON array first
	var results []TestResult
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&results); err == nil {
		return results, nil
	}

	// Fallback: try to read as JSON lines (old format)
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var res TestResult
		if err := json.Unmarshal([]byte(scanner.Text()), &res); err == nil {
			results = append(results, res)
		}
	}
	return results, scanner.Err()
}

func saveToHistory(dl, ul float64, ping int64, isp string) {
	cfg := config.Get()
	if !cfg.UI.SaveHistory {
		return
	}
	res := TestResult{
		Timestamp: time.Now().Format("15:04:05"),
		Download:  dl / 1000000.0,
		Upload:    ul / 1000000.0,
		Ping:      ping,
		ISP:       isp,
	}

	historyPath := config.GetHistoryFile()

	// Load existing history
	results, err := LoadHistory()
	if err != nil {
		logger := logging.Get()
		logger.Errorf("Failed to load history: %v", err)
		results = []TestResult{}
	}

	// Append new result
	results = append(results, res)

	// Enforce history limit
	cfg = config.Get()
	limit := cfg.UI.HistoryLimit
	if limit <= 0 {
		limit = 100 // default
	}
	if len(results) > limit {
		results = results[len(results)-limit:]
	}

	// Write back as JSON array
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		logger := logging.Get()
		logger.Errorf("Failed to marshal history: %v", err)
		return
	}

	if err := os.WriteFile(historyPath, jsonData, 0644); err != nil {
		logger := logging.Get()
		logger.Errorf("Failed to write history file: %v", err)
	}
}

// showHistory displays test history in a TUI table modal.
func showHistory() {
	table := tview.NewTable().SetBorders(true)
	table.SetTitle(" TEST HISTORY ").SetBorder(true).SetBorderColor(tcell.ColorYellow)

	headers := []string{"Time", "Down", "Up", "Ping", "ISP"}
	for c, h := range headers {
		table.SetCell(0, c, tview.NewTableCell(h).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	}

	results, err := LoadHistory()
	if err == nil {
		for row, res := range results {
			table.SetCell(row+1, 0, tview.NewTableCell(res.Timestamp))
			table.SetCell(row+1, 1, tview.NewTableCell(fmt.Sprintf("%.1f", res.Download)))
			table.SetCell(row+1, 2, tview.NewTableCell(fmt.Sprintf("%.1f", res.Upload)))
			table.SetCell(row+1, 3, tview.NewTableCell(fmt.Sprintf("%dms", res.Ping)))
			table.SetCell(row+1, 4, tview.NewTableCell(res.ISP))
		}
	}

	modal := tview.NewGrid().SetColumns(0, 70, 0).SetRows(0, 12, 0).AddItem(table, 1, 1, 1, 1, 0, 0, true)
	pages.AddPage("history", modal, true, true)
}

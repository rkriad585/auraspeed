package speedtest

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	// downloadData stores the speed data points for graph rendering.
	downloadData []float64
	// graphMu protects downloadData from concurrent access.
	graphMu sync.Mutex
	// lastGraphUpdate tracks the last time the graph was redrawn.
	lastGraphUpdate time.Time
	// graphUpdateMu protects lastGraphUpdate from concurrent access.
	graphUpdateMu sync.Mutex
)

// refreshGraph redraws the speed graph with throttling (max once per 500ms).
// It takes a snapshot of downloadData to avoid holding locks during rendering.
func refreshGraph() {
	graphUpdateMu.Lock()
	if time.Since(lastGraphUpdate) < 500*time.Millisecond {
		graphUpdateMu.Unlock()
		return
	}
	lastGraphUpdate = time.Now()
	graphUpdateMu.Unlock()

	graphMu.Lock()
	if len(downloadData) < 2 {
		graphMu.Unlock()
		return
	}
	data := make([]float64, len(downloadData))
	copy(data, downloadData)
	graphMu.Unlock()

	if len(data) > 40 {
		data = data[len(data)-40:]
	}

	// Generate simple ASCII graph without escape sequences
	graph := generateSimpleGraph(data)
	app.QueueUpdateDraw(func() {
		graphBox.SetText(graph)
	})
}

// generateSimpleGraph creates a simple text-based graph
func generateSimpleGraph(data []float64) string {
	if len(data) == 0 {
		return "No data yet..."
	}

	maxVal := 0.0
	for _, v := range data {
		if v > maxVal {
			maxVal = v
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	height := 6
	width := 40
	graph := ""

	for h := height; h > 0; h-- {
		threshold := maxVal * float64(h) / float64(height)
		line := ""
		for _, v := range data {
			if len(line) >= width {
				break
			}
			if v >= threshold {
				line += "█"
			} else {
				line += " "
			}
		}
		graph += fmt.Sprintf("%6.1f |%s\n", maxVal*float64(h)/float64(height), line)
	}
	graph += "       +" + strings.Repeat("─", width)

	return graph
}

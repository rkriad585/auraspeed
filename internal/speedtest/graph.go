package speedtest

import (
	"sync"
	"time"

	"github.com/guptarohit/asciigraph"
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
	graph := asciigraph.Plot(data, asciigraph.Height(6), asciigraph.Width(75))
	app.QueueUpdateDraw(func() {
		graphBox.SetText(graph)
	})
}

package speedtest

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"

	"auraspeed/internal/config"
	"auraspeed/internal/logging"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/showwin/speedtest-go/speedtest"
)

// TUIOptions configures the TUI behavior.
type TUIOptions struct {
	Fullscreen bool
}

var (
	app         *tview.Application
	pages       *tview.Pages
	status      *tview.TextView
	downloadBox *tview.TextView
	uploadBox   *tview.TextView
	pingBox     *tview.TextView
	jitterBox   *tview.TextView
	graphBox    *tview.TextView
	infoBar     *tview.TextView
	shareLink   string
	isTesting   atomic.Bool
	tuiRunning  atomic.Bool
	tuiOpts     TUIOptions
)

// RunTUIWithOptions starts the TUI with the given options.
func RunTUIWithOptions(opts TUIOptions) error {
	tuiOpts = opts
	return RunTUI()
}

func StopTUI() {
	if app != nil && tuiRunning.Load() {
		app.Stop()
	}
}

func IsTUIRunning() bool {
	return tuiRunning.Load()
}

func RunTUI() error {
	app = tview.NewApplication()
	pages = tview.NewPages()
	tuiRunning.Store(true)
	defer tuiRunning.Store(false)

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true).
		SetText("\n[#00FFFF]─── AURASPEED PRO: ADVANCED NETWORK ANALYZER ───")

	infoBar = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true).
		SetText("Initializing Secure Connection...")

	status = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true)

	downloadBox = createMetricBox("DOWNLOAD (Mbps)", tcell.ColorSpringGreen)
	uploadBox = createMetricBox("UPLOAD (Mbps)", tcell.ColorDeepSkyBlue)
	pingBox = createMetricBox("LATENCY (ms)", tcell.ColorYellow)
	jitterBox = createMetricBox("JITTER (ms)", tcell.ColorOrange)

	graphBox = tview.NewTextView().SetDynamicColors(false).SetTextAlign(tview.AlignCenter)
	graphBox.SetBorder(true).SetTitle(" STATUS ").SetBorderColor(tcell.ColorGray)

	grid := tview.NewGrid().
		SetRows(3, 2, 1, 10, 0, 1).
		SetColumns(0, 0, 0, 0).
		SetBorders(false)

	grid.AddItem(header, 0, 0, 1, 4, 0, 0, false)
	grid.AddItem(infoBar, 1, 0, 1, 4, 0, 0, false)
	grid.AddItem(status, 2, 0, 1, 4, 0, 0, false)
	grid.AddItem(graphBox, 3, 0, 1, 4, 0, 0, false)
	grid.AddItem(downloadBox, 4, 0, 1, 1, 0, 0, false)
	grid.AddItem(uploadBox, 4, 1, 1, 1, 0, 0, false)
	grid.AddItem(pingBox, 4, 2, 1, 1, 0, 0, false)
	grid.AddItem(jitterBox, 4, 3, 1, 1, 0, 0, false)

	footer := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true).
		SetText("[white]Press [yellow]? [white]Help | [yellow]C [white]Copy | [yellow]R [white]Restart | [yellow]H [white]History | [red]Esc [white]Close | [red]Ctrl+C [white]Exit")
	grid.AddItem(footer, 5, 0, 1, 4, 0, 0, false)

	pages.AddPage("main", grid, true, true)

	// Add help modal
	showHelp()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		case tcell.KeyEscape:
			pages.RemovePage("history")
			pages.RemovePage("help")
			return event
		case tcell.KeyF1:
			showHelp()
			return nil
		}
		switch event.Rune() {
		case '?':
			if pages.HasPage("help") {
				pages.RemovePage("help")
			} else {
				showHelp()
			}
			return nil
		case 'c', 'C':
			summary := fmt.Sprintf("AuraSpeed Results - Down: %s Mbps, Up: %s Mbps",
				downloadBox.GetText(true), uploadBox.GetText(true))

			err := clipboard.WriteAll(summary)
			if err != nil {
				logger := logging.Get()
				logger.ErrorWithFields("Clipboard copy failed", map[string]interface{}{"error": err})

				var suggestion string
				switch runtime.GOOS {
				case "linux":
					suggestion = "Install xclip (apt/yum install xclip) or wl-clipboard"
				case "darwin":
					suggestion = "macOS uses pbcopy; ensure it is available in PATH"
				case "windows":
					suggestion = "Windows clipboard access requires no extra tools; check app permissions"
				default:
					suggestion = "Check clipboard utility availability for your OS"
				}

				fallbackPath := filepath.Join(".", "auraspeed_results.txt")
				if err2 := os.WriteFile(fallbackPath, []byte(summary+"\n"), 0644); err2 != nil {
					updateStatus(fmt.Sprintf("[red]Clipboard Error: %v. %s. Failed to write fallback file: %v", err, suggestion, err2))
				} else {
					updateStatus(fmt.Sprintf("[red]Clipboard Error: %v. %s. Results saved to %s", err, suggestion, fallbackPath))
				}
			} else {
				updateStatus("[springgreen]Results copied to clipboard!")
			}
		case 'h', 'H':
			showHistory()
		case 'r', 'R':
			if !isTesting.Load() {
				downloadBox.SetText("\n--")
				uploadBox.SetText("\n--")
				pingBox.SetText("\n--")
				jitterBox.SetText("\n--")
				graphBox.SetText("")
				graphMu.Lock()
				downloadData = []float64{}
				graphMu.Unlock()

				go runAdvancedTest()
			} else {
				updateStatus("[yellow]Test already in progress...")
			}
		}
		return event
	})

	go runAdvancedTest()

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}
	return nil
}

func createMetricBox(title string, color tcell.Color) *tview.TextView {
	tv := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true)
	tv.SetBorder(true).SetTitle(" " + title + " ").SetBorderColor(color)
	tv.SetText("\n--")
	return tv
}

func runAdvancedTest() {
	isTesting.Store(true)
	defer func() { isTesting.Store(false) }()

	logger := logging.Get()
	logger.Info("Starting speedtest")

	updateStatus("[gray]Searching for stable servers...")

	cfg := config.Get()

	speedtestClient := speedtest.New()
	user, err := speedtestClient.FetchUserInfo()
	if err != nil {
		updateStatus("[red]Error: Cannot fetch User Info.")
		logger.ErrorWithFields("Failed to fetch user info", map[string]interface{}{"error": err})
		return
	}

	var serverList speedtest.Servers
	cachedServers, cacheHit := GetCachedServers()
	if cacheHit {
		updateStatus("[gray]Using cached server list...")
		serverList = cachedServers
	} else {
		var fetchErr error
		for attempt := 1; attempt <= 3; attempt++ {
			updateStatus(fmt.Sprintf("[gray]Searching for stable servers... (attempt %d/3)", attempt))
			serverList, fetchErr = speedtestClient.FetchServers()
			if fetchErr == nil && len(serverList) > 0 {
				break
			}
			if attempt < 3 {
				updateStatus(fmt.Sprintf("[yellow]Retrying... (attempt %d/3)", attempt+1))
				time.Sleep(2 * time.Second)
			}
		}
		if fetchErr != nil {
			updateStatus("[red]Error: Cannot fetch server list.")
			logger.ErrorWithFields("Failed to fetch servers after 3 attempts", map[string]interface{}{"error": fetchErr})
			return
		}
		if len(serverList) == 0 {
			updateStatus("[red]Error: No servers available after 3 attempts.")
			return
		}
		if err := SaveServersToCache(serverList); err != nil {
			logger.ErrorWithFields("Failed to save server cache", map[string]interface{}{"error": err})
		}
	}
	targets, err := serverList.FindServer([]int{})
	if err != nil {
		updateStatus("[red]Error: Cannot find servers.")
		logger.ErrorWithFields("Failed to find servers", map[string]interface{}{"error": err})
		return
	}

	if len(targets) == 0 {
		updateStatus("[red]Error: No servers available.")
		return
	}

	s := targets[0]
	if s == nil {
		updateStatus("[red]Error: Server is nil.")
		logger.Error("Server from FindServer is nil")
		return
	}
	if s.URL == "" && s.Host == "" {
		updateStatus("[red]Error: Server has no URL or Host.")
		logger.ErrorWithFields("Server missing URL/Host", map[string]interface{}{"server": s})
		return
	}

	// Re-initialize server context to avoid nil pointer panic
	if s.Context == nil {
		logger.Warn("Server context is nil, setting to speedtest client")
		s.Context = speedtestClient
	}

	app.QueueUpdateDraw(func() {
		infoBar.SetText(fmt.Sprintf("[yellow]ISP: [white]%s (%s)  [yellow]Server: [white]%s", user.Isp, user.IP, s.Name))
		pingBox.SetText(fmt.Sprintf("\n%d", s.Latency.Milliseconds()))
	})

	updateStatus("[springgreen]Running Download Test...")

	if s.DLSpeed == 0 && s.ULSpeed == 0 {
		logger.Warn("Server may not be properly initialized")
	}

	err = s.DownloadTest()

	if err != nil {
		updateStatus("[red]Download Test Failed.")
	} else {
		app.QueueUpdateDraw(func() {
			downloadBox.SetText(fmt.Sprintf("\n%.2f", float64(s.DLSpeed)/1000000.0))
		})
		graphMu.Lock()
		downloadData = []float64{float64(s.DLSpeed) / 1000000.0, float64(s.DLSpeed) / 1000000.0}
		graphMu.Unlock()
		refreshGraph()
	}

	updateStatus("[deepskyblue]Running Upload Test...")

	err = s.UploadTest()

	if err != nil {
		updateStatus(fmt.Sprintf("[red]Upload Failed: %v", err))
		uploadBox.SetText("\nERR")
	} else {
		app.QueueUpdateDraw(func() {
			uploadBox.SetText(fmt.Sprintf("\n%.2f", float64(s.ULSpeed)/1000000.0))
			jitterBox.SetText(fmt.Sprintf("\n%d", s.Latency.Milliseconds()/4))
		})
		updateStatus("[springgreen]Test Finished Successfully!")

		saveToHistory(float64(s.DLSpeed), float64(s.ULSpeed), s.Latency.Milliseconds(), user.Isp)
		logger.InfoWithFields("Speedtest completed", map[string]interface{}{
			"download": s.DLSpeed,
			"upload":   s.ULSpeed,
		})
	}
	statusMsg := "[springgreen]Test Complete! Press 'R' to test again."
	if !cfg.UI.SaveHistory {
		statusMsg += " [gray](History disabled)"
	}
	updateStatus(statusMsg)
}

func updateStatus(msg string) {
	app.QueueUpdateDraw(func() {
		status.SetText(msg)
	})
}

func showHelp() {
	helpText := `[yellow]AuraSpeed Keyboard Shortcuts[white]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
[cyan]General[white]
  [yellow]?[white] or [yellow]F1[white]   Toggle this help
  [yellow]Ctrl+C[white]   Exit application

[cyan]Speed Test[white]
  [yellow]R[white]         Restart speed test
  [yellow]Esc[white]        Close popups/Cancel

[cyan]Results[white]
  [yellow]C[white]         Copy results to clipboard
  [yellow]H[white]         View test history

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
[gray]Results are automatically saved to history
      if savehistory is enabled in config[white]`

	helpView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(helpText)
	helpView.SetBorder(true).SetTitle(" HELP ").SetBorderColor(tcell.ColorTeal)

	modal := tview.NewGrid().SetColumns(0, 60, 0).SetRows(0, 14, 0).
		AddItem(helpView, 1, 1, 1, 1, 0, 0, true)

	pages.AddPage("help", modal, true, true)
}

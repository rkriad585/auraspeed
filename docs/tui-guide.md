# Interactive TUI Guide

The AuraSpeed TUI (Terminal User Interface) is built using `tview` and provides a real-time dashboard for network performance.

## Interface Layout

1.  **Header:** Displays current ISP and Public IP information.
2.  **Metrics Boxes:** Live displays for Download, Upload, Ping, and Jitter.
3.  **Graph Area:** A real-time ASCII graph showing throughput fluctuations.
4.  **Status Bar:** Displays current operation status and system messages.
5.  **Help Bar:** Quick reference for keyboard shortcuts.

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Space` | Start/Stop Speed Test |
| `C` | Copy results to clipboard |
| `H` | Open History overlay |
| `S` | Open Server Selection menu |
| `?` | Show Help modal |
| `Ctrl+C` / `Q` | Exit AuraSpeed |

## Features

### Real-time Graphing
The graph updates every 500ms during active tests, providing a visual representation of connection stability. It automatically scales based on the maximum speed detected during the session.

### Clipboard Integration
By pressing `C`, AuraSpeed formats your latest results (Download, Upload, Ping, ISP) into a clean string and copies it to your system clipboard for easy sharing.

### History Overlay
The history view within the TUI allows you to scroll through previous tests without leaving the interactive session.


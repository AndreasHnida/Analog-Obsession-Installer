package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// AppLogger writes timestamped messages to the in-app log panel.
// File logging is off by default; call SetFileLog(true) to enable it.
type AppLogger struct {
	rt     *widget.RichText
	scroll *container.Scroll
	file   *os.File
	mu     sync.Mutex
}

func newAppLogger(rt *widget.RichText, scroll *container.Scroll) *AppLogger {
	return &AppLogger{rt: rt, scroll: scroll}
}

// SetFileLog opens or closes the log file next to the executable.
func (l *AppLogger) SetFileLog(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if enabled {
		if l.file != nil {
			return
		}
		execPath, err := os.Executable()
		if err != nil {
			return
		}
		logPath := filepath.Join(filepath.Dir(execPath), "AOInstaller.log")
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return
		}
		l.file = f
		fmt.Fprintln(l.file, "─── Session started ───────────────────────────────")
	} else {
		if l.file == nil {
			return
		}
		fmt.Fprintln(l.file, "─── Session ended ─────────────────────────────────")
		l.file.Close()
		l.file = nil
	}
}

var monoStyle = fyne.TextStyle{Monospace: true}

func (l *AppLogger) Log(msg string) {
	line := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), msg)
	l.writeLine(line)

	var colorName fyne.ThemeColorName
	switch {
	case strings.Contains(msg, "✓"):
		colorName = colorNameLogSuccess
	case strings.Contains(msg, "✗"):
		colorName = colorNameLogError
	}

	fyne.Do(func() {
		// Newline between entries (not before the first).
		if len(l.rt.Segments) > 0 {
			l.rt.Segments = append(l.rt.Segments, &widget.TextSegment{
				Text:  "\n",
				Style: widget.RichTextStyle{Inline: true, TextStyle: monoStyle},
			})
		}
		l.rt.Segments = append(l.rt.Segments, &widget.TextSegment{
			Text: line,
			Style: widget.RichTextStyle{
				Inline:    true,
				ColorName: colorName,
				TextStyle: monoStyle,
			},
		})
		l.rt.Refresh()
		l.scroll.ScrollToBottom()
	})
}

func (l *AppLogger) writeLine(line string) {
	l.mu.Lock()
	if l.file != nil {
		fmt.Fprintln(l.file, line)
	}
	l.mu.Unlock()
}

func (l *AppLogger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		fmt.Fprintln(l.file, "─── Session ended ─────────────────────────────────")
		l.file.Close()
		l.file = nil
	}
}

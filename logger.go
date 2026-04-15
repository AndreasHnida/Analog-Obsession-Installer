package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// AppLogger writes timestamped messages to the in-app log panel and to a
// log file next to the executable.
type AppLogger struct {
	rt     *widget.RichText
	scroll *container.Scroll
	file   *os.File
}

func newAppLogger(rt *widget.RichText, scroll *container.Scroll) *AppLogger {
	l := &AppLogger{rt: rt, scroll: scroll}

	execPath, err := os.Executable()
	if err == nil {
		logPath := filepath.Join(filepath.Dir(execPath), "AOInstaller.log")
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err == nil {
			l.file = f
			l.writeLine("─── Session started ───────────────────────────────")
		}
	}
	return l
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
	if l.file != nil {
		fmt.Fprintln(l.file, line)
	}
}

func (l *AppLogger) Close() {
	if l.file != nil {
		l.writeLine("─── Session ended ─────────────────────────────────")
		l.file.Close()
	}
}

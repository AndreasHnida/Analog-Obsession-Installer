package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// AppLogger writes timestamped messages to the in-app log panel and to a
// log file next to the executable.
type AppLogger struct {
	entry  *widget.Entry
	scroll *container.Scroll
	file   *os.File
}

func newAppLogger(entry *widget.Entry, scroll *container.Scroll) *AppLogger {
	l := &AppLogger{entry: entry, scroll: scroll}

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

func (l *AppLogger) Log(msg string) {
	line := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), msg)

	cur := l.entry.Text
	if cur == "" {
		l.entry.SetText(line)
	} else {
		l.entry.SetText(cur + "\n" + line)
	}
	l.scroll.ScrollToBottom()

	l.writeLine(line)
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

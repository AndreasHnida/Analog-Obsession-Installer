package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icon.png
var iconBytes []byte

// appIcon is the application icon used in the window title bar and taskbar.
var appIcon = fyne.NewStaticResource("icon.png", iconBytes)

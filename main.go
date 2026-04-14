package main

import (
	"context"
	"fmt"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type pluginRow struct {
	plugin Plugin
	check  *widget.Check
	badge  *widget.Label
	obj    fyne.CanvasObject
}

func defaultVST3Path() string {
	switch runtime.GOOS {
	case "windows":
		return `C:\Program Files\Common Files\VST3`
	case "darwin":
		return "/Library/Audio/Plug-Ins/VST3"
	default:
		return "/usr/lib/vst3"
	}
}

func main() {
	a := app.NewWithID("com.analogobsession.aoinstaller")
	a.Settings().SetTheme(bwTheme{})
	a.SetIcon(appIcon)

	w := a.NewWindow("Analog Obsession Installer")
	w.Resize(fyne.NewSize(640, 820))
	w.CenterOnScreen()

	// ── Header ────────────────────────────────────────────────────────────────
	titleLabel := widget.NewLabelWithStyle(
		"Analog Obsession",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	subtitleLabel := widget.NewLabelWithStyle(
		"Free VST3 Plugin Installer  •  Windows",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	// ── Install path row ──────────────────────────────────────────────────────
	pathEntry := widget.NewEntry()
	pathEntry.SetText(defaultVST3Path())
	pathEntry.SetPlaceHolder("VST3 install directory...")

	browseBtn := widget.NewButton("Browse…", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			pathEntry.SetText(uri.Path())
		}, w)
	})
	pathRow := container.NewBorder(nil, nil, nil, browseBtn, pathEntry)

	// ── Plugin rows ───────────────────────────────────────────────────────────
	prows := make([]*pluginRow, len(Plugins))
	objs := make([]fyne.CanvasObject, len(Plugins))
	for i, p := range Plugins {
		check := widget.NewCheck(p.Name, nil)
		badge := widget.NewLabel("")
		row := container.NewBorder(nil, nil, nil, badge, check)
		prows[i] = &pluginRow{plugin: p, check: check, badge: badge, obj: row}
		objs[i] = row
	}

	pluginScroll := container.NewVScroll(container.NewVBox(objs...))
	pluginScroll.SetMinSize(fyne.NewSize(600, 360))

	// ── Status refresh ────────────────────────────────────────────────────────
	refreshStatus := func(vst3Dir string) {
		for _, r := range prows {
			if IsInstalled(r.plugin.Name, vst3Dir) {
				r.badge.SetText("✓ installed")
			} else {
				r.badge.SetText("")
			}
		}
	}
	pathEntry.OnChanged = func(s string) { refreshStatus(s) }

	// ── Selection buttons ─────────────────────────────────────────────────────
	selectAllBtn := widget.NewButton("Select All", func() {
		for _, r := range prows {
			r.check.SetChecked(true)
		}
	})
	deselectAllBtn := widget.NewButton("Deselect All", func() {
		for _, r := range prows {
			r.check.SetChecked(false)
		}
	})
	selectionRow := container.NewHBox(selectAllBtn, deselectAllBtn)

	// ── Log panel ─────────────────────────────────────────────────────────────
	logEntry := widget.NewMultiLineEntry()
	logEntry.Disable() // read-only; text colour set via theme.ColorNameDisabled
	logEntry.SetMinRowsVisible(6)

	logScroll := container.NewVScroll(logEntry)
	logScroll.SetMinSize(fyne.NewSize(600, 110))

	logger := newAppLogger(logEntry, logScroll)
	defer logger.Close()

	logger.Log("Ready — select plugins and press Install.")

	// ── Progress bar ──────────────────────────────────────────────────────────
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	// ── Install button ────────────────────────────────────────────────────────
	var installBtn *widget.Button
	var uninstallBtn *widget.Button
	var actionRow *fyne.Container
	var cancelInstall context.CancelFunc

	var installAction func()
	installAction = func() {
		selected := make([]Plugin, 0, len(Plugins))
		for _, r := range prows {
			if r.check.Checked {
				selected = append(selected, r.plugin)
			}
		}
		if len(selected) == 0 {
			dialog.ShowInformation("Nothing selected",
				"Please check at least one plugin before installing.", w)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancelInstall = cancel

		installBtn.SetText("Cancel Installation")
		installBtn.Importance = widget.DangerImportance
		installBtn.OnTapped = func() {
			cancelInstall()
			logger.Log("Cancelling — waiting for current download to stop…")
		}
		installBtn.Refresh()
		actionRow.Objects = []fyne.CanvasObject{installBtn}
		actionRow.Layout = layout.NewGridLayout(1)
		actionRow.Refresh()
		progressBar.Show()
		progressBar.SetValue(0)

		go func() {
			defer func() {
				cancel()
				installBtn.SetText("Install Selected")
				installBtn.Importance = widget.HighImportance
				installBtn.OnTapped = installAction
				installBtn.Refresh()
				actionRow.Objects = []fyne.CanvasObject{installBtn, uninstallBtn}
				actionRow.Layout = layout.NewGridLayout(2)
				actionRow.Refresh()
				refreshStatus(pathEntry.Text)
				progressBar.Hide()
			}()

			logger.Log(fmt.Sprintf("Starting installation of %d plugin(s)…", len(selected)))
			total := float64(len(selected))
			errCount := 0
			cancelled := false

			for i, p := range selected {
				if ctx.Err() != nil {
					cancelled = true
					break
				}

				if err := InstallPlugin(ctx, p, pathEntry.Text, func(msg string) {
					logger.Log(msg)
				}); err != nil {
					if ctx.Err() != nil {
						cancelled = true
						break
					}
					errCount++
					logger.Log(fmt.Sprintf("✗ %s — %v", p.Name, err))
				} else {
					logger.Log(fmt.Sprintf("✓ %s installed.", p.Name))
				}

				progressBar.SetValue(float64(i+1) / total)
			}

			switch {
			case cancelled:
				logger.Log("Installation cancelled.")
				progressBar.SetValue(0)
			case errCount == 0:
				logger.Log(fmt.Sprintf("✓ Done! %d plugin(s) installed successfully.", len(selected)))
			default:
				logger.Log(fmt.Sprintf("Done — %d error(s). See log above.", errCount))
			}
		}()
	}

	installBtn = widget.NewButton("Install Selected", installAction)
	installBtn.Importance = widget.HighImportance

	// ── Uninstall button ──────────────────────────────────────────────────────
	uninstallBtn = widget.NewButton("Uninstall Selected", func() {
		selected := make([]*pluginRow, 0)
		for _, r := range prows {
			if r.check.Checked && IsInstalled(r.plugin.Name, pathEntry.Text) {
				selected = append(selected, r)
			}
		}
		if len(selected) == 0 {
			dialog.ShowInformation("Nothing to uninstall",
				"No checked plugins are currently installed at the selected path.", w)
			return
		}

		dialog.ShowConfirm(
			"Uninstall plugins",
			fmt.Sprintf("Remove %d plugin(s) from disk?", len(selected)),
			func(ok bool) {
				if !ok {
					return
				}
				uninstallBtn.Disable()
				errCount := 0
				for _, r := range selected {
					if err := UninstallPlugin(r.plugin.Name, pathEntry.Text); err != nil {
						errCount++
						logger.Log(fmt.Sprintf("✗ %s — %v", r.plugin.Name, err))
					} else {
						logger.Log(fmt.Sprintf("✓ %s removed.", r.plugin.Name))
					}
				}
				refreshStatus(pathEntry.Text)
				if errCount == 0 {
					logger.Log(fmt.Sprintf("✓ Removed %d plugin(s).", len(selected)))
				} else {
					logger.Log(fmt.Sprintf("Done — %d error(s).", errCount))
				}
				uninstallBtn.Enable()
			}, w)
	})
	uninstallBtn.Importance = widget.DangerImportance

	// ── Root layout ───────────────────────────────────────────────────────────
	actionRow = container.NewGridWithColumns(2, installBtn, uninstallBtn)

	content := container.NewVBox(
		titleLabel,
		subtitleLabel,
		widget.NewSeparator(),
		pathRow,
		selectionRow,
		widget.NewSeparator(),
		pluginScroll,
		widget.NewSeparator(),
		progressBar,
		logScroll,
		actionRow,
	)

	w.SetContent(container.NewPadded(container.NewPadded(content)))

	// Initial scan — pre-select every plugin not yet installed.
	refreshStatus(pathEntry.Text)
	for _, r := range prows {
		if !IsInstalled(r.plugin.Name, pathEntry.Text) {
			r.check.SetChecked(true)
		}
	}

	w.ShowAndRun()
}

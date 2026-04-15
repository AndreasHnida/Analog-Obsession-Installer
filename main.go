package main

import (
	"context"
	"fmt"
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
	"path"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type pluginRow struct {
	plugin    Plugin
	check     *widget.Check
	vst3Badge *canvas.Text
	aaxBadge  *canvas.Text
	obj       fyne.CanvasObject
}

// infoBtn is a small icon button that fires onHover(desc) on mouse-in and
// onHover("") on mouse-out, and opens the plugin page in the browser on tap.
type infoBtn struct {
	widget.Button
	desc    string
	onHover func(string)
}

func newInfoBtn(desc, pageURL string, onHover func(string)) *infoBtn {
	b := &infoBtn{desc: desc, onHover: onHover}
	b.Button = widget.Button{
		Icon:       theme.InfoIcon(),
		Importance: widget.LowImportance,
		OnTapped: func() {
			if pageURL == "" {
				return
			}
			u, err := url.Parse(pageURL)
			if err != nil {
				return
			}
			_ = fyne.CurrentApp().OpenURL(u)
		},
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *infoBtn) MouseIn(e *desktop.MouseEvent) {
	b.Button.MouseIn(e)
	if b.onHover != nil {
		b.onHover(b.desc)
	}
}

func (b *infoBtn) MouseOut() {
	b.Button.MouseOut()
	if b.onHover != nil {
		b.onHover("")
	}
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

func defaultAAXPath() string {
	switch runtime.GOOS {
	case "windows":
		return `C:\Program Files\Common Files\Avid\Audio\Plug-Ins`
	case "darwin":
		return "/Library/Application Support/Avid/Audio/Plug-Ins"
	default:
		return ""
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
	titleLabel := canvas.NewText("Analog Obsession", colBlack)
	titleLabel.TextSize = 24
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter
	subtitleLabel := widget.NewLabelWithStyle(
		"Batch Plugin Installer",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	prefs := a.Preferences()

	// ── Install path rows ─────────────────────────────────────────────────────
	pathEntry := widget.NewEntry()
	pathEntry.SetText(defaultVST3Path())
	pathEntry.SetPlaceHolder("VST3 install directory...")
	if saved := prefs.String("vst3Dir"); saved != "" {
		pathEntry.SetText(saved)
	}

	aaxPathEntry := widget.NewEntry()
	aaxPathEntry.SetText(defaultAAXPath())
	aaxPathEntry.SetPlaceHolder("AAX install directory...")
	if saved := prefs.String("aaxDir"); saved != "" {
		aaxPathEntry.SetText(saved)
	}

	browseBtn := widget.NewButton("Browse…", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			pathEntry.SetText(uri.Path())
		}, w)
	})
	vst3PathLabel := widget.NewLabel("VST3")
	pathRow := container.NewBorder(nil, nil, vst3PathLabel, browseBtn, pathEntry)
	pathRow.Hide()

	aaxBrowseBtn := widget.NewButton("Browse…", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			aaxPathEntry.SetText(uri.Path())
		}, w)
	})
	aaxPathLabel := widget.NewLabel("AAX")
	aaxPathRow := container.NewBorder(nil, nil, aaxPathLabel, aaxBrowseBtn, aaxPathEntry)
	aaxPathRow.Hide()

	// ── Plugin rows ───────────────────────────────────────────────────────────
	hintLabel := canvas.NewText("", colGrey150)
	hintLabel.TextSize = 10
	hintLabel.Alignment = fyne.TextAlignCenter

	const imgW, imgH = float32(88), float32(44)

	// greyPx is a 1×1 grey image used as placeholder while device images load.
	greyPx := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	greyPx.SetNRGBA(0, 0, struct{ R, G, B, A uint8 }{220, 220, 220, 255})

	prows := make([]*pluginRow, len(Plugins))
	objs := make([]fyne.CanvasObject, len(Plugins))
	for i, p := range Plugins {
		check := widget.NewCheck("", nil)

		// name + description stacked vertically
		nameText := canvas.NewText(p.Name, colBlack)
		nameText.TextSize = 12

		descText := canvas.NewText(p.Desc, colGrey150)
		descText.TextSize = 10
		descText.TextStyle = fyne.TextStyle{Italic: true}

		textBox := container.NewVBox(
			container.New(layout.NewCustomPaddedLayout(4, 0, 0, 0), nameText),
			container.New(layout.NewCustomPaddedLayout(0, 4, 0, 0), descText),
		)

		// device image — starts as grey placeholder, loaded asynchronously
		imgCanvas := canvas.NewImageFromImage(greyPx)
		imgCanvas.FillMode = canvas.ImageFillStretch
		imgCanvas.SetMinSize(fyne.NewSize(imgW, imgH))

		if p.ImgURL != "" {
			if data, err := assetFS.ReadFile("assets/" + path.Base(p.ImgURL)); err == nil {
				if decoded, _, err := image.Decode(bytes.NewReader(data)); err == nil {
					imgCanvas.Image = decoded
					imgCanvas.FillMode = canvas.ImageFillContain
				}
			}
		}

		vst3Badge := canvas.NewText("vst3:", colGrey150)
		vst3Badge.TextSize = 9
		vst3Badge.Alignment = fyne.TextAlignTrailing

		aaxBadge := canvas.NewText("aax:", colGrey150)
		aaxBadge.TextSize = 9
		aaxBadge.Alignment = fyne.TextAlignTrailing

		info := newInfoBtn(p.Desc, p.PageURL, func(s string) {
			hintLabel.Text = s
			hintLabel.Refresh()
		})
		right := container.NewHBox(
			container.New(
				layout.NewCustomPaddedLayout(0, 0, 0, 8),
				container.New(layout.NewGridLayout(2), vst3Badge, aaxBadge),
			),
			info,
		)

		center := container.NewBorder(nil, nil, imgCanvas, nil, textBox)
		row := container.NewBorder(nil, nil, check, right, center)
		prows[i] = &pluginRow{plugin: p, check: check, vst3Badge: vst3Badge, aaxBadge: aaxBadge, obj: row}
		objs[i] = row
	}


	pluginScroll := container.NewVScroll(container.NewVBox(objs...))
	pluginScroll.SetMinSize(fyne.NewSize(600, 340))

	// ── Status refresh ────────────────────────────────────────────────────────
	countLabel := widget.NewLabel("")
	countLabel.Importance = widget.LowImportance

	refreshStatus := func(vst3Dir, aaxDir string) {
		installed := 0
		for _, r := range prows {
			vst3ok := IsInstalled(r.plugin, vst3Dir)
			aaxok := IsAAXInstalled(r.plugin, aaxDir)
			if vst3ok || aaxok {
				installed++
			}
			if vst3ok {
				r.vst3Badge.Text = "vst3: ✓"
				r.vst3Badge.Color = colGrey80
			} else {
				r.vst3Badge.Text = "vst3:"
				r.vst3Badge.Color = colGrey150
			}
			r.vst3Badge.Refresh()
			if aaxok {
				r.aaxBadge.Text = "aax: ✓"
				r.aaxBadge.Color = colGrey80
			} else {
				r.aaxBadge.Text = "aax:"
				r.aaxBadge.Color = colGrey150
			}
			r.aaxBadge.Refresh()
		}
		countLabel.SetText(fmt.Sprintf("%d / %d installed", installed, len(Plugins)))
	}
	pathEntry.OnChanged = func(s string) {
		prefs.SetString("vst3Dir", s)
		refreshStatus(s, aaxPathEntry.Text)
	}
	aaxPathEntry.OnChanged = func(s string) {
		prefs.SetString("aaxDir", s)
		refreshStatus(pathEntry.Text, s)
	}

	// ── Format checkboxes + selection buttons ────────────────────────────────
	var installBtn *widget.Button
	var vst3Check, aaxCheck *widget.Check

	updateInstallBtn := func() {
		// installBtn may not be set yet during initial setup; guard nil
		if installBtn == nil {
			return
		}
		anyPlugin := false
		for _, r := range prows {
			if r.check.Checked {
				anyPlugin = true
				break
			}
		}
		anyFormat := vst3Check.Checked || aaxCheck.Checked
		if anyPlugin && anyFormat {
			installBtn.Enable()
		} else {
			installBtn.Disable()
		}
	}

	// Wire plugin check callbacks now that updateInstallBtn is defined.
	for _, r := range prows {
		r := r
		r.check.OnChanged = func(_ bool) { updateInstallBtn() }
	}

	vst3Check = widget.NewCheck("VST3", func(_ bool) {
		updateInstallBtn()
	})
	vst3Check.SetChecked(true)

	aaxCheck = widget.NewCheck("AAX", func(checked bool) {
		if !checked {
			aaxPathRow.Hide()
		}
		updateInstallBtn()
	})

	var logger *AppLogger // initialized after log widgets are created below

	logFileCheck := widget.NewCheck("Write log file", func(checked bool) {
		prefs.SetBool("fileLog", checked)
		if logger != nil {
			logger.SetFileLog(checked)
		}
	})
	logFileCheck.SetChecked(prefs.Bool("fileLog"))

	settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		d := dialog.NewCustom("Settings", "Done",
			container.NewVBox(
				container.NewBorder(nil, nil, vst3PathLabel, browseBtn, pathEntry),
				container.NewBorder(nil, nil, aaxPathLabel, aaxBrowseBtn, aaxPathEntry),
				widget.NewSeparator(),
				logFileCheck,
			), w)
		d.Resize(fyne.NewSize(500, 160))
		d.Show()
	})
	settingsBtn.Importance = widget.LowImportance

	formatRow := container.NewHBox(vst3Check, aaxCheck, settingsBtn)

	selectAllBtn := widget.NewButton("Select All", func() {
		for _, r := range prows {
			r.check.SetChecked(true)
		}
		updateInstallBtn()
	})
	deselectAllBtn := widget.NewButton("Deselect All", func() {
		for _, r := range prows {
			r.check.SetChecked(false)
		}
		updateInstallBtn()
	})
	selectionRow := container.NewBorder(nil, nil, nil, formatRow,
		container.NewHBox(selectAllBtn, deselectAllBtn))

	// ── Search + count row ────────────────────────────────────────────────────
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Filter plugins…")
	searchEntry.OnChanged = func(q string) {
		q = strings.ToLower(q)
		for _, r := range prows {
			if q == "" || strings.Contains(strings.ToLower(r.plugin.Name), q) {
				r.obj.Show()
			} else {
				r.obj.Hide()
			}
		}
		pluginScroll.Refresh()
	}
	searchRow := container.NewBorder(nil, nil, nil, countLabel, searchEntry)

	// ── Log panel ─────────────────────────────────────────────────────────────
	logRT := widget.NewRichText()
	logRT.Wrapping = fyne.TextWrapWord

	logScroll := container.NewVScroll(logRT)
	logScroll.SetMinSize(fyne.NewSize(600, 110))

	logger = newAppLogger(logRT, logScroll)
	defer logger.Close()
	if prefs.Bool("fileLog") {
		logger.SetFileLog(true)
	}

	logger.Log("Ready — select plugins and press Install.")

	// ── Progress bar ──────────────────────────────────────────────────────────
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	// ── Install button ────────────────────────────────────────────────────────
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
				fyne.Do(func() {
					installBtn.SetText("Install Selected")
					installBtn.Importance = widget.HighImportance
					installBtn.OnTapped = installAction
					installBtn.Refresh()
					actionRow.Objects = []fyne.CanvasObject{installBtn, uninstallBtn}
					actionRow.Layout = layout.NewGridLayout(2)
					actionRow.Refresh()
					refreshStatus(pathEntry.Text, aaxPathEntry.Text)
					updateInstallBtn()
					progressBar.Hide()
				})
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

				if err := InstallPlugin(ctx, p, pathEntry.Text, aaxPathEntry.Text,
					vst3Check.Checked, aaxCheck.Checked, func(msg string) {
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

				v := float64(i+1) / total
				fyne.Do(func() { progressBar.SetValue(v) })
			}

			switch {
			case cancelled:
				logger.Log("Installation cancelled.")
				fyne.Do(func() { progressBar.SetValue(0) })
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
	type uninstallEntry struct {
		row    *pluginRow
		doVST3 bool
		doAAX  bool
	}
	uninstallBtn = widget.NewButton("Uninstall Selected", func() {
		selected := make([]uninstallEntry, 0)
		for _, r := range prows {
			if !r.check.Checked {
				continue
			}
			doVST3 := vst3Check.Checked && IsInstalled(r.plugin, pathEntry.Text)
			doAAX := aaxCheck.Checked && IsAAXInstalled(r.plugin, aaxPathEntry.Text)
			if doVST3 || doAAX {
				selected = append(selected, uninstallEntry{r, doVST3, doAAX})
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
				for _, e := range selected {
					if e.doVST3 {
						if err := UninstallPlugin(e.row.plugin, pathEntry.Text); err != nil {
							errCount++
							logger.Log(fmt.Sprintf("✗ %s VST3 — %v", e.row.plugin.Name, err))
						} else {
							logger.Log(fmt.Sprintf("✓ %s VST3 removed.", e.row.plugin.Name))
						}
					}
					if e.doAAX {
						if err := UninstallAAX(e.row.plugin, aaxPathEntry.Text); err != nil {
							errCount++
							logger.Log(fmt.Sprintf("✗ %s AAX — %v", e.row.plugin.Name, err))
						} else {
							logger.Log(fmt.Sprintf("✓ %s AAX removed.", e.row.plugin.Name))
						}
					}
				}
				refreshStatus(pathEntry.Text, aaxPathEntry.Text)
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

	versionLabel := widget.NewLabelWithStyle(
		"v1.0.0",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)
	versionLabel.Importance = widget.LowImportance

	content := container.NewVBox(
		titleLabel,
		subtitleLabel,
		widget.NewSeparator(),
		selectionRow,
		searchRow,
		widget.NewSeparator(),
		pluginScroll,
		hintLabel,
		widget.NewSeparator(),
		progressBar,
		logScroll,
		actionRow,
		versionLabel,
	)

	w.SetContent(container.NewPadded(container.NewPadded(content)))

	// Initial scan — pre-select every plugin not yet installed.
	refreshStatus(pathEntry.Text, aaxPathEntry.Text)
	for _, r := range prows {
		if !IsInstalled(r.plugin, pathEntry.Text) {
			r.check.SetChecked(true)
		}
	}
	updateInstallBtn()

	w.ShowAndRun()
}

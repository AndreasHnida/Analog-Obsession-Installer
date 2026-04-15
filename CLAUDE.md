# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Fan-made batch installer for Analog Obsession's free VST3/AAX plugins. Windows-targeted GUI app written in Go using the Fyne toolkit. Single `main` package, flat layout.

## Build

Cross-compile from Linux to Windows (CGO + MinGW required by Fyne):

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
  go build -ldflags="-H windowsgui -s -w" -o AOInstaller.exe .
```

Native dev requires Go 1.23+, `gcc-mingw-w64-x86-64`, `libgl1-mesa-dev`, `xorg-dev`.

`rsrc_windows_amd64.syso` embeds the Windows icon + UAC manifest (`app.manifest` requests admin elevation). Regenerate with `rsrc` if `app.rc`/`app.manifest`/`logo.ico` change.

## Test

```bash
go test ./...
go test -run TestInstallFromLocalZip -v
```

`installer_test.go` skips unless `EDComp_1.0.zip` sits in the repo root. Note: it currently references a stale symbol `copyVST3Files` (renamed to `copyPluginFiles`); the test won't compile until updated.

## Architecture

Five source files, one package:

- `main.go` — Fyne UI. Builds plugin rows, format checkboxes (VST3/AAX), search filter, install/uninstall/cancel button state machine, log panel, progress bar. Persists install paths via `fyne.Preferences`. Install runs in goroutine with `context.CancelFunc` for cancellation.
- `installer.go` — Pure logic, no UI. `InstallPlugin` = download zip → temp extract → walk for `.vst3`/`.aaxplugin` bundles → copy to dest. Has zip-slip guard. `IsInstalled`/`IsAAXInstalled`/`Uninstall*` use `findBundle` for fuzzy matching of versioned bundle names (`Name_1.0.vst3` etc.).
- `plugins.go` — Static `Plugins []Plugin` catalogue (name, zip URL, page URL, description, image URL, optional `BundleName` override). Sorted A→Z.
- `apptheme.go` — `bwTheme` pure black/white Fyne theme. Red kept only for danger/cancel.
- `logger.go`, `logo.go`, `assets.go` — UI log sink, embedded app icon, `embed.FS` for `assets/` device images.

Bundle naming: `Plugin.Bundle()` returns `BundleName` if set, else `Name`. Match logic in `findBundle` accepts `name`, `name_*`, `name.*`, `name-*` (case-insensitive) so version suffixes resolve.

Default install paths (`defaultVST3Path` / `defaultAAXPath`) branch by `runtime.GOOS` — Windows is the only fully supported target. AAX path empty on Linux.

## Adding a plugin

Append to `Plugins` in `plugins.go` (keep alphabetical). Drop matching device PNG into `assets/` named exactly `path.Base(p.ImgURL)` — `main.go` reads from `embed.FS` by basename.

package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// bwTheme is a pure black-and-white theme.
// Checkboxes, progress bars, and focused borders use dark grey.
// The danger (cancel) colour stays red — the only accent in the UI.
type bwTheme struct{}

var _ fyne.Theme = (*bwTheme)(nil)

const (
	colorNameLogSuccess fyne.ThemeColorName = "logSuccess"
	colorNameLogError   fyne.ThemeColorName = "logError"
)

var (
	colWhite      = color.White
	colBlack      = color.Black
	colNearBlack  = color.NRGBA{R: 30, G: 30, B: 30, A: 255}
	colGrey20     = color.NRGBA{R: 20, G: 20, B: 20, A: 255}
	colGrey80     = color.NRGBA{R: 80, G: 80, B: 80, A: 255}
	colGrey150    = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
	colGrey160    = color.NRGBA{R: 160, G: 160, B: 160, A: 255}
	colGrey180    = color.NRGBA{R: 180, G: 180, B: 180, A: 255}
	colGrey200    = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	colGrey210    = color.NRGBA{R: 210, G: 210, B: 210, A: 255}
	colGrey220    = color.NRGBA{R: 220, G: 220, B: 220, A: 255}
	colGrey240    = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	colGrey245    = color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	colShadow     = color.NRGBA{R: 0, G: 0, B: 0, A: 25}
	colBrightGreen = color.NRGBA{R: 0, G: 230, B: 80, A: 255}
	colBrightRed   = color.NRGBA{R: 255, G: 60, B: 60, A: 255}
)

func (bwTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	// ── Backgrounds ─────────────────────────────────────────────────────────
	case theme.ColorNameBackground:
		return colWhite
	case theme.ColorNameOverlayBackground:
		return colWhite
	case theme.ColorNameMenuBackground:
		return colWhite
	case theme.ColorNameHeaderBackground:
		return colGrey245
	case theme.ColorNameInputBackground:
		return colGrey245
	case theme.ColorNameScrollBarBackground:
		return colGrey240

	// ── Foregrounds ──────────────────────────────────────────────────────────
	case theme.ColorNameForeground:
		return colBlack
	case theme.ColorNamePlaceHolder:
		return colGrey150
	case theme.ColorNameForegroundOnPrimary:
		return colWhite // text on primary-coloured (dark) backgrounds
	case theme.ColorNameForegroundOnError:
		return colWhite
	case theme.ColorNameForegroundOnSuccess:
		return colWhite
	case theme.ColorNameForegroundOnWarning:
		return colBlack

	// ── Accent / interactive ─────────────────────────────────────────────────
	case theme.ColorNamePrimary:
		return colNearBlack // checkboxes, progress bar, HighImportance button
	case theme.ColorNameFocus:
		return color.Transparent
	case theme.ColorNameHover:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0} // no hover tint on buttons
	case theme.ColorNamePressed:
		return colGrey180
	case theme.ColorNameSelection:
		return colGrey210
	case theme.ColorNameHyperlink:
		return colGrey20

	// ── Buttons ──────────────────────────────────────────────────────────────
	case theme.ColorNameButton:
		return colGrey240
	case theme.ColorNameDisabledButton:
		return colGrey220
	case theme.ColorNameDisabled:
		return colGrey80 // dark enough to read in the log panel

	// ── Borders / decoration ─────────────────────────────────────────────────
	case theme.ColorNameInputBorder:
		return colGrey180
	case theme.ColorNameScrollBar:
		return colGrey160
	case theme.ColorNameSeparator:
		return colGrey200
	case theme.ColorNameShadow:
		return colShadow

	// ── Log panel colours ────────────────────────────────────────────────────
	case colorNameLogSuccess:
		return colBrightGreen
	case colorNameLogError:
		return colBrightRed

	// ── Status colours — keep red/amber/green for cancel button etc. ─────────
	case theme.ColorNameError, theme.ColorNameWarning, theme.ColorNameSuccess:
		return theme.LightTheme().Color(name, variant)

	default:
		return theme.LightTheme().Color(name, variant)
	}
}

func (bwTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.LightTheme().Font(style)
}

func (bwTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.LightTheme().Icon(name)
}

func (bwTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameText:
		return 12 // default 14 → -2
	case theme.SizeNameHeadingText:
		return 22 // default 24 → -2
	case theme.SizeNameSubHeadingText:
		return 16 // default 18 → -2
	case theme.SizeNameCaptionText:
		return 9
	default:
		return theme.LightTheme().Size(name)
	}
}

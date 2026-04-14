package main

import "fyne.io/fyne/v2"

// appIcon is the application icon used in the window title bar and taskbar.
// Fyne renders SVG natively so no rasterisation needed.
var appIcon = fyne.NewStaticResource("logo.svg", []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 200">
  <rect width="200" height="200" fill="white"/>
  <polyline points="8,164 48,24 88,164"
            fill="none" stroke="black" stroke-width="11"
            stroke-linejoin="miter" stroke-miterlimit="10"/>
  <line x1="26" y1="102" x2="70" y2="102"
        stroke="black" stroke-width="11" stroke-linecap="butt"/>
  <circle cx="152" cy="94" r="46"
          fill="none" stroke="black" stroke-width="11"/>
</svg>`))

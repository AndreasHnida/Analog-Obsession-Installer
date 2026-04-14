# Analog Obsession Installer

A fan-made batch installer for the free VST3 plugins by **Analog Obsession**.

> **Disclaimer** — This is an independent fan project with no affiliation to, endorsement by, or relationship with Analog Obsession or its creator. All plugins are the intellectual property of their respective author. This tool simply automates downloading and placing the files you could install manually yourself.

---

## What it does

- Lists all free Analog Obsession VST3 plugins in one window
- Shows which plugins are already installed
- Downloads and installs selected plugins in one click
- Supports cancellation mid-install and removal of installed plugins
- Writes a log file next to the executable for troubleshooting
- Requires no separate installer — single portable `.exe`

## Requirements

- Windows 10 or 11 (64-bit)
- Administrator privileges (the installer will prompt via UAC)
- Internet connection

## Usage

1. Download `AOInstaller.exe` from [Releases](../../releases)
2. Run it — UAC will ask for elevation
3. Select the plugins you want
4. Click **Install Selected**

The default install path is `C:\Program Files\Common Files\VST3`, which is standard for all major DAWs. You can change it via the Browse button.

## Credits

All plugins are created by **Analog Obsession** and distributed free of charge.  
Please support the creator directly:

| | |
|---|---|
| **Website** | [analogobsession.com](https://analogobsession.com) |
| **Patreon** | [patreon.com/analogobsession](https://www.patreon.com/analogobsession) |

If you find these plugins useful, consider becoming a patron. The plugins are free — the creator relies on community support to keep making them.

## Building from source

Requirements: Go 1.23+, `gcc-mingw-w64-x86-64`, `libgl1-mesa-dev`, `xorg-dev`

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
  go build -ldflags="-H windowsgui -s -w" -o AOInstaller.exe .
```

## License

MIT — see [LICENSE](LICENSE)  
Plugin files remain the property of their respective author and are subject to their original license terms.

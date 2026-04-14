package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// InstallPlugin downloads the plugin zip, extracts it to a temp directory,
// finds all .vst3 files/bundles, and copies them to vst3Dir.
// onStatus is called with human-readable progress messages.
// Cancelling ctx aborts the operation; any partial download is cleaned up.
func InstallPlugin(ctx context.Context, p Plugin, vst3Dir string, onStatus func(string)) error {
	// --- Download ---
	onStatus(fmt.Sprintf("[%s] Downloading...", p.Name))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.ZipURL, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status)
	}

	tmpZip, err := os.CreateTemp("", "aoinstaller-*.zip")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpZip.Name())

	if _, err = io.Copy(tmpZip, resp.Body); err != nil {
		tmpZip.Close()
		return fmt.Errorf("write zip: %w", err)
	}
	tmpZip.Close()

	if err := ctx.Err(); err != nil {
		return err
	}

	// --- Extract ---
	onStatus(fmt.Sprintf("[%s] Extracting...", p.Name))

	tmpDir, err := os.MkdirTemp("", "aoinstaller-extract-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := unzip(tmpZip.Name(), tmpDir); err != nil {
		return fmt.Errorf("extract: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	// --- Install ---
	onStatus(fmt.Sprintf("[%s] Installing to %s...", p.Name, vst3Dir))

	if err := os.MkdirAll(vst3Dir, 0o755); err != nil {
		return fmt.Errorf("create VST3 dir: %w", err)
	}

	return copyVST3Files(tmpDir, vst3Dir, p.Name)
}

// unzip extracts all entries from src into dst.
func unzip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(dst, filepath.FromSlash(f.Name))

		// Prevent zip-slip attacks.
		if !strings.HasPrefix(
			filepath.Clean(target)+string(os.PathSeparator),
			filepath.Clean(dst)+string(os.PathSeparator),
		) {
			return fmt.Errorf("illegal path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, f.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := writeZipEntry(f, target); err != nil {
			return err
		}
	}
	return nil
}

func writeZipEntry(f *zip.File, dst string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}

// copyVST3Files walks srcDir, finds every .vst3 entry (file or bundle folder),
// and copies it into dstDir. Returns an error if no .vst3 is found.
func copyVST3Files(srcDir, dstDir, pluginName string) error {
	var found bool

	err := fs.WalkDir(os.DirFS(srcDir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.EqualFold(filepath.Ext(path), ".vst3") {
			return nil
		}

		found = true
		abs := filepath.Join(srcDir, filepath.FromSlash(path))
		dst := filepath.Join(dstDir, filepath.Base(abs))

		if d.IsDir() {
			if err := copyDir(abs, dst); err != nil {
				return err
			}
			return fs.SkipDir // don't recurse inside the bundle
		}
		return copyFile(abs, dst)
	})
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("no .vst3 file found inside the zip for %q", pluginName)
	}
	return nil
}

// copyDir recursively copies a directory tree from src to dst.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip Windows shell metadata files — they may have hidden/system
		// attributes that block overwriting on reinstall.
		if strings.EqualFold(info.Name(), "desktop.ini") {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

// IsInstalled reports whether a plugin's .vst3 bundle exists in vst3Dir.
// Matches case-insensitively and tolerates version suffixes in the filename
// (e.g. "room041_2.0.vst3" matches plugin name "Room041").
func IsInstalled(name, vst3Dir string) bool {
	entries, err := os.ReadDir(vst3Dir)
	if err != nil {
		return false
	}
	lname := strings.ToLower(name)
	for _, e := range entries {
		en := strings.ToLower(e.Name())
		if !strings.HasSuffix(en, ".vst3") {
			continue
		}
		stem := strings.TrimSuffix(en, ".vst3")
		if stem == lname || strings.HasPrefix(stem, lname+"_") || strings.HasPrefix(stem, lname+".") {
			return true
		}
	}
	return false
}

// UninstallPlugin removes a plugin's .vst3 bundle from vst3Dir.
func UninstallPlugin(name, vst3Dir string) error {
	return os.RemoveAll(filepath.Join(vst3Dir, name+".vst3"))
}

// copyFile copies a single file from src to dst, preserving permissions.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

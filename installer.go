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
// and copies the requested format bundles to their respective directories.
// onStatus is called with human-readable progress messages.
// Cancelling ctx aborts the operation; any partial download is cleaned up.
func InstallPlugin(ctx context.Context, p Plugin, vst3Dir, aaxDir string, installVST3, installAAX bool, onStatus func(string)) error {
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

	if installVST3 {
		if err := os.MkdirAll(vst3Dir, 0o755); err != nil {
			return fmt.Errorf("create VST3 dir: %w", err)
		}
	}
	if installAAX && aaxDir != "" {
		if err := os.MkdirAll(aaxDir, 0o755); err != nil {
			return fmt.Errorf("create AAX dir: %w", err)
		}
	}

	return copyPluginFiles(tmpDir, vst3Dir, aaxDir, p.Name, installVST3, installAAX, onStatus)
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

// copyPluginFiles walks srcDir and copies .vst3 and/or .aaxplugin bundles
// into their respective destination directories.
func copyPluginFiles(srcDir, vst3Dir, aaxDir, pluginName string, installVST3, installAAX bool, onStatus func(string)) error {
	var foundVST3, foundAAX bool

	err := fs.WalkDir(os.DirFS(srcDir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		ext := strings.ToLower(filepath.Ext(path))

		if ext == ".vst3" && installVST3 {
			foundVST3 = true
			abs := filepath.Join(srcDir, filepath.FromSlash(path))
			dst := filepath.Join(vst3Dir, filepath.Base(abs))
			onStatus(fmt.Sprintf("[%s] VST3: %s → %s", pluginName, filepath.Base(abs), dst))
			if d.IsDir() {
				if err := copyDir(abs, dst); err != nil {
					return fmt.Errorf("copyDir %s: %w", filepath.Base(abs), err)
				}
				return fs.SkipDir
			}
			return copyFile(abs, dst)
		}

		if ext == ".aaxplugin" && installAAX {
			foundAAX = true
			abs := filepath.Join(srcDir, filepath.FromSlash(path))
			dst := filepath.Join(aaxDir, filepath.Base(abs))
			onStatus(fmt.Sprintf("[%s] AAX: %s → %s", pluginName, filepath.Base(abs), dst))
			if d.IsDir() {
				if err := copyDir(abs, dst); err != nil {
					return fmt.Errorf("copyDir %s: %w", filepath.Base(abs), err)
				}
				return fs.SkipDir
			}
			return copyFile(abs, dst)
		}

		return nil
	})
	if err != nil {
		return err
	}
	if installVST3 && !foundVST3 {
		return fmt.Errorf("no .vst3 found in zip for %q", pluginName)
	}
	if installAAX && !foundAAX {
		return fmt.Errorf("no .aaxplugin found in zip for %q", pluginName)
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

// findBundle scans dir for a bundle matching bundleName with the given
// extension, case-insensitively and tolerating common version suffixes.
// Returns the full path or "".
func findBundle(bundleName, dir, ext string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	lname := strings.ToLower(bundleName)
	for _, e := range entries {
		en := strings.ToLower(e.Name())
		if !strings.HasSuffix(en, ext) {
			continue
		}
		stem := strings.TrimSuffix(en, ext)
		if stem == lname ||
			strings.HasPrefix(stem, lname+"_") ||
			strings.HasPrefix(stem, lname+".") ||
			strings.HasPrefix(stem, lname+"-") {
			return filepath.Join(dir, e.Name())
		}
	}
	return ""
}

// IsInstalled reports whether a plugin's .vst3 bundle exists in vst3Dir.
func IsInstalled(p Plugin, vst3Dir string) bool {
	return findBundle(p.Bundle(), vst3Dir, ".vst3") != ""
}

// IsAAXInstalled reports whether a plugin's .aaxplugin bundle exists in aaxDir.
func IsAAXInstalled(p Plugin, aaxDir string) bool {
	return findBundle(p.Bundle(), aaxDir, ".aaxplugin") != ""
}

// UninstallPlugin removes a plugin's .vst3 bundle from vst3Dir.
func UninstallPlugin(p Plugin, vst3Dir string) error {
	path := findBundle(p.Bundle(), vst3Dir, ".vst3")
	if path == "" {
		return fmt.Errorf("%s not found in %s", p.Bundle(), vst3Dir)
	}
	return os.RemoveAll(path)
}

// UninstallAAX removes a plugin's .aaxplugin bundle from aaxDir.
func UninstallAAX(p Plugin, aaxDir string) error {
	path := findBundle(p.Bundle(), aaxDir, ".aaxplugin")
	if path == "" {
		return fmt.Errorf("%s not found in %s", p.Bundle(), aaxDir)
	}
	return os.RemoveAll(path)
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

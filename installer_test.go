package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallFromLocalZip(t *testing.T) {
	zipPath := "EDComp_1.0.zip"
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Skip("EDComp_1.0.zip not present, skipping")
	}

	outDir := t.TempDir()

	// Replicate what InstallPlugin does after download.
	tmpDir := t.TempDir()
	if err := unzip(zipPath, tmpDir); err != nil {
		t.Fatalf("unzip: %v", err)
	}

	if err := copyVST3Files(tmpDir, outDir, "EDComp"); err != nil {
		t.Fatalf("copyVST3Files: %v", err)
	}

	// Walk the output and report what was installed.
	t.Log("Files written to output dir:")
	filepath.Walk(outDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(outDir, path)
		if info.IsDir() {
			t.Logf("  DIR  %s/", rel)
		} else {
			t.Logf("  FILE %s  (%d bytes)", rel, info.Size())
		}
		return nil
	})

	// Assert the VST3 bundle folder exists at the top level.
	bundlePath := filepath.Join(outDir, "EDComp.vst3")
	if info, err := os.Stat(bundlePath); err != nil || !info.IsDir() {
		t.Errorf("expected EDComp.vst3 bundle folder at output root, got: %v", err)
	}

	// Assert the inner DLL exists.
	dllPath := filepath.Join(outDir, "EDComp.vst3", "Contents", "x86_64-win", "EDComp.vst3")
	if info, err := os.Stat(dllPath); err != nil || info.IsDir() {
		t.Errorf("expected inner DLL at %s, got: %v", dllPath, err)
	}

	// Assert AAX was NOT copied.
	aaxPath := filepath.Join(outDir, "EDComp.aaxplugin")
	if _, err := os.Stat(aaxPath); err == nil {
		t.Errorf("AAX bundle should NOT have been copied, but it was")
	}
}

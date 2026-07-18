package aria2

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckInstalledUsesConfiguredBinary(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "custom-aria2c")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := CheckInstalled(binary); err != nil {
		t.Fatalf("configured executable should be accepted: %v", err)
	}
}

func TestCheckInstalledRejectsMissingConfiguredBinary(t *testing.T) {
	if err := CheckInstalled(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("missing configured executable should be rejected")
	}
}

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultDownloadDir(t *testing.T) {
	cfg := Default()
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, "Downloads")
	if cfg.DownloadDir != want {
		t.Fatalf("got %q want %q", cfg.DownloadDir, want)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfg := Default()
	cfg.DownloadDir = "/tmp/custom"
	cfg.RPCPort = 6801

	if err := Save(cfg); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.DownloadDir != cfg.DownloadDir {
		t.Fatalf("download dir: got %q want %q", loaded.DownloadDir, cfg.DownloadDir)
	}
	if loaded.RPCPort != cfg.RPCPort {
		t.Fatalf("rpc port: got %d want %d", loaded.RPCPort, cfg.RPCPort)
	}

	path, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("config permissions: got %o want 600", info.Mode().Perm())
	}
}

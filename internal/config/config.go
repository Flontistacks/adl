package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDir  = ".config/adl"
	ConfigFile = "config.yaml"
)

type Config struct {
	DownloadDir string `yaml:"download_dir"`
	Aria2Path   string `yaml:"aria2_path"`
	RPCPort     int    `yaml:"rpc_port"`
}

func Default() Config {
	home, _ := os.UserHomeDir()
	downloads := filepath.Join(home, "Downloads")
	aria2, _ := lookPath("aria2c")
	return Config{
		DownloadDir: downloads,
		Aria2Path:   aria2,
		RPCPort:     6800,
	}
}

func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigDir, ConfigFile), nil
}

func Load() (Config, error) {
	cfg := Default()
	path, err := Path()
	if err != nil {
		return cfg, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if cfg.DownloadDir == "" {
		cfg.DownloadDir = Default().DownloadDir
	}
	if cfg.Aria2Path == "" {
		cfg.Aria2Path, _ = lookPath("aria2c")
	}
	if cfg.RPCPort == 0 {
		cfg.RPCPort = 6800
	}
	return cfg, nil
}

func Save(cfg Config) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func lookPath(name string) (string, error) {
	return execLookPath(name)
}

package aria2

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Daemon struct {
	cmd    *exec.Cmd
	client *Client
}

func sessionSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate rpc secret: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func Start(binary string, port int, downloadDir string) (*Daemon, error) {
	if binary == "" {
		return nil, fmt.Errorf("aria2c not found; install with: brew install aria2")
	}
	secret, err := sessionSecret()
	if err != nil {
		return nil, err
	}
	args := []string{
		"--enable-rpc",
		"--rpc-listen-all=false",
		fmt.Sprintf("--dir=%s", downloadDir),
		"--quiet=true",
		fmt.Sprintf("--rpc-listen-port=%d", port),
		"--rpc-secret=" + secret,
	}
	cmd := exec.Command(binary, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start aria2c: %w", err)
	}

	client := NewClient(port, secret)
	d := &Daemon{cmd: cmd, client: client}

	// Wait for RPC to become ready
	var lastErr error
	for i := 0; i < 30; i++ {
		if _, err := client.call("aria2.getVersion"); err == nil {
			return d, nil
		} else {
			lastErr = err
		}
		time.Sleep(100 * time.Millisecond)
	}
	_ = d.Stop()
	return nil, fmt.Errorf("aria2 RPC not ready: %v", lastErr)
}

func (d *Daemon) Client() *Client {
	return d.client
}

func (d *Daemon) Stop() error {
	if d.cmd == nil || d.cmd.Process == nil {
		return nil
	}
	_ = d.cmd.Process.Kill()
	_, _ = d.cmd.Process.Wait()
	return nil
}

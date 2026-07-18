package aria2

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Daemon struct {
	cmd      *exec.Cmd
	client   *Client
	stopOnce sync.Once
	stopErr  error
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
	d.stopOnce.Do(func() {
		d.stopErr = d.stop()
	})
	return d.stopErr
}

func (d *Daemon) stop() error {
	if d.cmd == nil || d.cmd.Process == nil {
		return nil
	}

	done := make(chan error, 1)
	go func() {
		done <- d.cmd.Wait()
	}()

	signalErr := d.cmd.Process.Signal(os.Interrupt)
	if errors.Is(signalErr, os.ErrProcessDone) {
		signalErr = nil
	}

	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()

	select {
	case <-done:
		if signalErr != nil {
			return fmt.Errorf("stop aria2c: %w", signalErr)
		}
		return nil
	case <-timer.C:
		if err := d.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return fmt.Errorf("kill aria2c after shutdown timeout: %w", err)
		}
		<-done
		return nil
	}
}

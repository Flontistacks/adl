package aria2

import (
	"fmt"
	"os/exec"
)

func CheckInstalled() error {
	_, err := exec.LookPath("aria2c")
	if err != nil {
		return fmt.Errorf("aria2c not found. Install with: brew install aria2")
	}
	return nil
}

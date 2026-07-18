package aria2

import (
	"fmt"
	"os/exec"
)

func CheckInstalled(binary string) error {
	if binary == "" {
		return fmt.Errorf("aria2c not found. Install with: brew install aria2")
	}
	_, err := exec.LookPath(binary)
	if err != nil {
		return fmt.Errorf("aria2c not found at %q. Install with: brew install aria2", binary)
	}
	return nil
}

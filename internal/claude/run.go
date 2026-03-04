package claude

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FindClaudeBin returns the path to the claude binary.
// Checks common locations: hint (if provided), PATH, ~/.local/bin/claude
func FindClaudeBin(hint string) string {
	if hint != "" {
		return hint
	}

	if path, err := exec.LookPath("claude"); err == nil {
		return path
	}

	home, err := os.UserHomeDir()
	if err == nil {
		local := filepath.Join(home, ".local", "bin", "claude")
		if _, err := os.Stat(local); err == nil {
			return local
		}
	}

	return "claude"
}

// Run executes `claude --print <prompt>` and returns the response text.
// It filters CLAUDECODE from the environment (prevents nested session errors).
// If the claude binary is not found, returns a clear error message.
func Run(claudeBin, prompt string) (string, error) {
	if claudeBin == "" {
		var err error
		claudeBin, err = exec.LookPath("claude")
		if err != nil {
			return "", fmt.Errorf("claude not found at %s — install from https://claude.ai/code", "claude")
		}
	}

	cmd := exec.Command(claudeBin, "--print", prompt)

	filtered := make([]string, 0, len(os.Environ()))
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "CLAUDECODE=") {
			filtered = append(filtered, e)
		}
	}
	cmd.Env = filtered

	out, err := cmd.Output()
	if err != nil {
		var execErr *exec.Error
		if errors.As(err, &execErr) && errors.Is(execErr.Err, exec.ErrNotFound) {
			return "", fmt.Errorf("claude not found at %s — install from https://claude.ai/code", claudeBin)
		}
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

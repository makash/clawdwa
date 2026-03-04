package systemd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"text/template"
)

// IsAvailable returns true if systemd is running on this system.
func IsAvailable() bool {
	if _, err := os.Stat("/run/systemd/systemd.pid"); err == nil {
		return true
	}
	cmd := exec.Command("systemctl", "is-system-running")
	out, err := cmd.Output()
	if err == nil {
		return true
	}
	output := strings.TrimSpace(string(out))
	return strings.Contains(output, "running") || strings.Contains(output, "degraded")
}

const unitTemplate = `[Unit]
Description=clawdwa WhatsApp Claude Bot
After=network.target

[Service]
Type=simple
ExecStart={{.BinaryPath}}
Restart=always
RestartSec=5
User={{.User}}
Environment=HOME={{.Home}}

[Install]
WantedBy=multi-user.target
`

type unitData struct {
	BinaryPath string
	User       string
	Home       string
}

// Install generates and installs a systemd unit file for clawdwa.
// Writes to /etc/systemd/system/clawdwa.service using sudo.
// Then runs: systemctl daemon-reload && systemctl enable --now clawdwa
func Install() error {
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	username := os.Getenv("USER")
	if username == "" {
		u, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}
		username = u.Username
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	tmpl, err := template.New("unit").Parse(unitTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse unit template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, unitData{
		BinaryPath: binaryPath,
		User:       username,
		Home:       homeDir,
	}); err != nil {
		return fmt.Errorf("failed to render unit template: %w", err)
	}
	unitContent := buf.String()

	fmt.Println("Writing /etc/systemd/system/clawdwa.service ...")
	cmd := exec.Command("sudo", "tee", "/etc/systemd/system/clawdwa.service")
	cmd.Stdin = strings.NewReader(unitContent)
	cmd.Stdout = io.Discard
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to write unit file: %w", err)
	}
	fmt.Println("Unit file written.")

	fmt.Println("Running systemctl daemon-reload ...")
	if err := exec.Command("sudo", "systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("failed to run daemon-reload: %w", err)
	}
	fmt.Println("daemon-reload complete.")

	fmt.Println("Enabling and starting clawdwa service ...")
	if err := exec.Command("sudo", "systemctl", "enable", "--now", "clawdwa").Run(); err != nil {
		return fmt.Errorf("failed to enable/start clawdwa service: %w", err)
	}
	fmt.Println("clawdwa service enabled and started.")

	return nil
}

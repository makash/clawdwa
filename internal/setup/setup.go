package setup

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/makash/clawdwa/internal/claude"
	"github.com/makash/clawdwa/internal/config"
	"github.com/makash/clawdwa/internal/systemd"
	"github.com/makash/clawdwa/internal/wa"
)

// Run runs the interactive setup wizard.
func Run() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== clawdwa setup ===")
	fmt.Println()

	// Step 1: Connect to WhatsApp (shows QR if not authenticated).
	fmt.Println("Connecting to WhatsApp...")
	client, err := wa.New()
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	defer client.Disconnect()

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	fmt.Println("✓ Connected to WhatsApp")
	fmt.Println()

	// Step 2: List and pick a group.
	fmt.Print("Search groups by name (or press Enter to list all): ")
	search, _ := reader.ReadString('\n')
	search = strings.ToLower(strings.TrimSpace(search))

	groups, err := client.GetGroups(ctx)
	if err != nil {
		return fmt.Errorf("get groups: %w", err)
	}

	var filtered []wa.Group
	for _, g := range groups {
		if search == "" || strings.Contains(strings.ToLower(g.Name), search) {
			filtered = append(filtered, g)
		}
	}

	if len(filtered) == 0 {
		return fmt.Errorf("no groups found matching %q — create a WhatsApp group first", search)
	}

	fmt.Println()
	fmt.Println("Matching groups:")
	for i, g := range filtered {
		fmt.Printf("  %d. %s\n", i+1, g.Name)
	}
	fmt.Println()

	fmt.Print("Enter group number: ")
	numStr, _ := reader.ReadString('\n')
	numStr = strings.TrimSpace(numStr)
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > len(filtered) {
		return fmt.Errorf("invalid selection %q", numStr)
	}

	selected := filtered[num-1]
	fmt.Printf("✓ Selected: %s (%s)\n", selected.Name, selected.JID)
	fmt.Println()

	// Step 3: Bot phone number (to filter self-messages).
	fmt.Print("Enter your WhatsApp number with country code (e.g. 919900000000): ")
	phone, _ := reader.ReadString('\n')
	phone = strings.TrimSpace(phone)
	botJID := phone + "@s.whatsapp.net"
	fmt.Printf("✓ Bot JID: %s\n", botJID)
	fmt.Println()

	// Step 4: Find claude binary.
	claudeBin := claude.FindClaudeBin("")
	fmt.Printf("✓ Claude binary: %s\n", claudeBin)
	fmt.Println()

	// Step 5: Write config.
	cfg := &config.Config{
		GroupJID:  selected.JID,
		GroupName: selected.Name,
		BotJID:    botJID,
		Prefixes:  []string{"!", "@claude"},
		ClaudeBin: claudeBin,
	}
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	fmt.Printf("✓ Config saved to %s\n", config.Path())
	fmt.Println()

	// Step 6: Optional systemd install.
	if systemd.IsAvailable() {
		fmt.Print("Install as systemd service? (requires sudo) [y/N]: ")
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(strings.ToLower(ans))
		if ans == "y" || ans == "yes" {
			if err := systemd.Install(); err != nil {
				fmt.Fprintf(os.Stderr, "systemd install failed: %v\n", err)
				fmt.Println("You can start the bot manually with: clawdwa")
			}
		}
	}

	fmt.Println()
	fmt.Println("=== Setup complete! ===")
	fmt.Println()
	fmt.Printf("Bot will respond to messages in '%s' starting with:\n", selected.Name)
	fmt.Println("  !        e.g. \"! what is recursion?\"")
	fmt.Println("  @claude  e.g. \"@claude write a bash script\"")
	return nil
}

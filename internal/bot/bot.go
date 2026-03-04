package bot

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"github.com/makash/clawdwa/internal/claude"
	"github.com/makash/clawdwa/internal/config"
	"github.com/makash/clawdwa/internal/wa"
)

// Run starts the bot. Blocks until interrupted.
func Run(cfg *config.Config) error {
	db, err := openBotDB()
	if err != nil {
		return fmt.Errorf("open bot db: %w", err)
	}
	defer db.Close()

	client, err := wa.New()
	if err != nil {
		return fmt.Errorf("create wa client: %w", err)
	}
	defer client.Disconnect()

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	fmt.Printf("✓ Bot running. Listening in group: %s\n", cfg.GroupName)
	fmt.Println("  Send '! your question' or '@claude your question' to trigger Claude.")
	fmt.Println("  Press Ctrl+C to stop.")

	client.AddMessageHandler(func(sender, chatJID, msgID, text string) {
		// Only handle messages from the configured group.
		if chatJID != cfg.GroupJID {
			return
		}

		// Skip bot's own messages.
		if strings.HasPrefix(sender, strings.TrimSuffix(cfg.BotJID, "@s.whatsapp.net")) {
			return
		}

		// Extract prompt from prefix.
		prompt, ok := extractPrompt(cfg.Prefixes, text)
		if !ok {
			return
		}

		// Dedup: skip already-processed messages.
		if seen, _ := isProcessed(db, msgID); seen {
			return
		}
		if err := markProcessed(db, msgID); err != nil {
			fmt.Fprintf(os.Stderr, "dedup error: %v\n", err)
			return
		}

		fmt.Printf("[%s] from %s: %s\n", time.Now().Format("15:04:05"), sender, prompt)

		// Call Claude.
		response, err := claude.Run(cfg.ClaudeBin, prompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "claude error: %v\n", err)
			response = fmt.Sprintf("Error calling Claude: %v", err)
		}

		// Send reply.
		if err := client.SendText(ctx, cfg.GroupJID, response); err != nil {
			fmt.Fprintf(os.Stderr, "send error: %v\n", err)
		}
	})

	// Wait for interrupt.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("\nStopping bot...")
	return nil
}

// Status shows whether the bot is running (via pidfile or process search).
func Status() {
	out, err := exec.Command("pgrep", "-f", "clawdwa").Output()
	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		fmt.Println("✗ Bot is not running")
		return
	}
	fmt.Printf("✓ Bot is running (PIDs: %s)\n", strings.TrimSpace(string(out)))
}

// Stop kills running clawdwa processes.
func Stop() {
	err := exec.Command("pkill", "-f", "clawdwa").Run()
	if err != nil {
		fmt.Println("Bot was not running")
	} else {
		fmt.Println("✓ Bot stopped")
	}
}

func extractPrompt(prefixes []string, text string) (string, bool) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(text, prefix+" ") {
			return strings.TrimPrefix(text, prefix+" "), true
		}
		if strings.HasPrefix(text, prefix) && len(text) > len(prefix) {
			rest := text[len(prefix):]
			if rest[0] != ' ' {
				return strings.TrimSpace(rest), true
			}
		}
	}
	return "", false
}

func openBotDB() (*sql.DB, error) {
	if err := os.MkdirAll(config.Dir(), 0700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", "file:"+config.BotDBPath()+"?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS processed_messages (
		id TEXT PRIMARY KEY,
		processed_at INTEGER NOT NULL
	)`)
	return db, err
}

func isProcessed(db *sql.DB, id string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM processed_messages WHERE id = ?", id).Scan(&count)
	return count > 0, err
}

func markProcessed(db *sql.DB, id string) error {
	_, err := db.Exec("INSERT OR IGNORE INTO processed_messages (id, processed_at) VALUES (?, ?)",
		id, time.Now().Unix())
	return err
}

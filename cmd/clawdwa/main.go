package main

import (
	"fmt"
	"os"

	"github.com/makash/clawdwa/internal/bot"
	"github.com/makash/clawdwa/internal/config"
	"github.com/makash/clawdwa/internal/setup"
	"github.com/makash/clawdwa/internal/wa"
)

var version = "dev"

func main() {
	sub := ""
	if len(os.Args) > 1 {
		sub = os.Args[1]
	}

	switch sub {
	case "setup":
		if err := setup.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "setup failed: %v\n", err)
			os.Exit(1)
		}
	case "status":
		bot.Status()
	case "stop":
		bot.Stop()
	case "send":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: clawdwa send <message>")
			os.Exit(1)
		}
		msg := os.Args[2]
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "no config found — run clawdwa setup first\n")
			os.Exit(1)
		}
		if err := wa.SendOnce(cfg, msg); err != nil {
			fmt.Fprintf(os.Stderr, "send failed: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println(version)
	case "help", "--help", "-h":
		printHelp()
	default:
		run()
	}
}

func run() {
	// First run: no config → setup first
	if !config.Exists() {
		fmt.Println("Welcome to clawdwa! Running setup...")
		fmt.Println()
		if err := setup.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "setup failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := bot.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "bot error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`clawdwa — WhatsApp ↔ Claude. In both directions.

Usage:
  clawdwa           Start the bot (runs setup on first run)
  clawdwa setup     Re-run setup (change group, phone number)
  clawdwa send MSG  Send a message to the configured group
  clawdwa status    Show bot status and recent log lines
  clawdwa stop      Stop the bot / systemd service
  clawdwa version   Show version

Group members trigger Claude with:
  ! your question
  @claude your question
`)
}

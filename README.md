# whatsapp-claudecode

Send WhatsApp messages to Claude Code and get responses — in any WhatsApp group.

## How it works

```
WhatsApp group member sends "! your question"
  → bot picks it up via whatsapp-cli
  → pipes it to Claude Code CLI
  → replies in the group
```

Triggers: messages starting with `!` or `@claude`

## Requirements

- Linux (tested on Ubuntu)
- Claude Code CLI (`claude`) — authenticated
- [mise](https://mise.run) for Go installation
- A WhatsApp account

## Install

```bash
git clone https://github.com/YOUR_USERNAME/whatsapp-claudecode
cd whatsapp-claudecode

# Install mise (if not already)
curl https://mise.run | sh
eval "$(~/.local/bin/mise activate bash)"

# Install whatsapp-cli
mise use --global go@latest
go install github.com/vicentereig/whatsapp-cli@latest
```

## Setup (run once)

```bash
bash setup.sh
```

This will:
1. Authenticate your WhatsApp account (QR code scan)
2. Let you pick which group to use
3. Ask for your phone number (to filter bot's own messages)
4. Write `store/config.sh` and initialize the database

## Run

```bash
bash bot.sh
```

That's it. The bot manages sync internally — no second terminal needed.

## Usage

In your configured WhatsApp group:

```
! what is a binary search tree?
@claude write a python script to rename files
```

## Change group

```bash
# Stop the bot first (Ctrl+C), then:
bash setup.sh
bash bot.sh
```

## Claude Code Skill

If you use Claude Code, install the skill for easy management:

```bash
cp whatsapp-bot.md ~/.claude/skills/whatsapp-bot.md
```

Then use `/whatsapp-bot start`, `/whatsapp-bot stop`, `/whatsapp-bot status`, `/whatsapp-bot change-group`.

## Files

```
bot.sh        — main bot (run this)
setup.sh      — first-time setup wizard
sync.sh       — standalone sync (optional, bot.sh handles this)
store/        — WhatsApp session + message DB (gitignored)
```

## Built with

- [whatsapp-cli](https://github.com/vicentereig/whatsapp-cli) — WhatsApp Web protocol
- [Claude Code](https://claude.ai/code) — AI responses

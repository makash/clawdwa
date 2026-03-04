# clawdwa

Use Claude AI directly from WhatsApp — no app, no signup. Just message a group.

## How it works

One person (the admin) sets this up on a Linux machine. Everyone else just sends messages in a WhatsApp group as usual.

```
Group member sends:  "! what is a binary search tree?"
Bot replies:         "A binary search tree is..."

Group member sends:  "@claude write a bash script to rename files"
Bot replies:         "Here's a bash script..."
```

## For group members

Nothing to install. Just join the WhatsApp group and send messages starting with `!` or `@claude`.

## For the admin (one-time setup)

**Requirements:**
- Linux machine (always-on, e.g. a VPS or home server)
- [Claude Code](https://claude.ai/code) installed and authenticated
- [mise](https://mise.run) for installing Go and whatsapp-cli

**Install:**

```bash
# Install mise
curl https://mise.run | sh
eval "$(~/.local/bin/mise activate bash)"

# Install Go and whatsapp-cli
mise use --global go@latest
go install github.com/vicentereig/whatsapp-cli@latest

# Clone and set up
git clone https://github.com/makash/clawdwa
cd clawdwa
bash setup.sh
```

`setup.sh` will:
1. Scan a QR code to link your WhatsApp account
2. Let you pick which group the bot listens to
3. Ask for your phone number (so the bot ignores its own messages)

**Run:**

```bash
bash bot.sh
```

That's it. Leave it running. The bot manages everything internally.

**Change the group later:**

```bash
# Stop the bot (Ctrl+C), then:
bash setup.sh
bash bot.sh
```

## Optional: Claude Code skill

If you use Claude Code, copy the skill for `/whatsapp-bot` shortcuts:

```bash
cp whatsapp-bot.md ~/.claude/skills/whatsapp-bot.md
```

Then use `/whatsapp-bot status`, `/whatsapp-bot stop`, etc.

## Built with

- [whatsapp-cli](https://github.com/vicentereig/whatsapp-cli) — WhatsApp Web protocol
- [Claude Code](https://claude.ai/code) — AI responses

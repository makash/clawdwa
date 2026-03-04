# clawdwa

WhatsApp ↔ Claude. In both directions.

```
Group member sends:  "! what is a binary search tree?"
Bot replies:         "A binary search tree is..."

You tell Claude:     "Deploy and notify the group when done."
Claude sends:        clawdwa send "✅ Deploy finished. All checks passed."
```

---

## For the admin

**Requirements:** [Claude Code](https://claude.ai/code) installed and authenticated on a Linux or macOS machine.

```bash
curl -fsSL https://github.com/makash/clawdwa/releases/latest/download/clawdwa_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o clawdwa
chmod +x clawdwa
./clawdwa
```

On first run, clawdwa will:
1. Show a QR code — scan with WhatsApp to link your account
2. Let you pick which group the bot listens to
3. Ask for your phone number (so the bot ignores its own messages)
4. Optionally install itself as a systemd service

### Subcommands

```
clawdwa           Start the bot (runs setup on first run)
clawdwa setup     Re-run setup (change group, phone number)
clawdwa send MSG  Send a message to the configured group
clawdwa status    Show bot status
clawdwa stop      Stop the bot
```

---

## For group members

Nothing to install. Just send messages in the WhatsApp group starting with `!` or `@claude`.

---

## Claude Code skill (admin only)

Install the skill so Claude can proactively message your group:

```bash
cp whatsapp-bot.md ~/.claude/skills/whatsapp-bot.md
```

Then in any Claude Code session:

```
Run the deployment, and send a WhatsApp when it's done.
```

Claude will call `clawdwa send "✅ Deployment complete"` automatically.

Also gives you `/whatsapp-bot status`, `/whatsapp-bot stop`, and `/whatsapp-bot send` shortcuts in Claude Code.

---

## Demo

**Setup**

![clawdwa setup](https://raw.githubusercontent.com/makash/clawdwa/master/demo/setup.gif)

**Running**

![clawdwa running](https://raw.githubusercontent.com/makash/clawdwa/master/demo/start.gif)

---

## Built with

- [whatsmeow](https://go.mau.fi/whatsmeow) — WhatsApp Web protocol in Go (MPL 2.0)
- [Claude Code](https://claude.ai/code) — AI responses

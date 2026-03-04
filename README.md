# clawdwa

WhatsApp ↔ Claude. In both directions.

## Setup

![clawdwa setup](https://raw.githubusercontent.com/makash/clawdwa/master/demo/setup.gif)

## Running

![clawdwa running](https://raw.githubusercontent.com/makash/clawdwa/master/demo/start.gif)

Send messages from WhatsApp to get Claude's help. Or tell Claude to notify your group when a task completes — from the terminal.

```
Group member sends:  "! what is a binary search tree?"
Bot replies:         "A binary search tree is..."

Claude sends:        clawdwa send "✅ Deploy finished. All checks passed."
Group receives:      "✅ Deploy finished. All checks passed."
```

## How it works

One person (the admin) sets this up on a Linux machine. Everyone else just sends messages in a WhatsApp group as usual.

## For group members

Nothing to install. Join the WhatsApp group and send messages starting with `!` or `@claude`.

## For the admin — Go binary (recommended)

Download the binary for your platform from [GitHub Releases](https://github.com/makash/clawdwa/releases), then:

```bash
chmod +x clawdwa
./clawdwa
```

That's it. On first run, clawdwa will:
1. Show a QR code to link your WhatsApp account
2. Let you pick which group the bot listens to
3. Ask for your phone number
4. Optionally install itself as a systemd service

**Requirements:** [Claude Code](https://claude.ai/code) installed and authenticated.

### Subcommands

```
clawdwa           Start the bot (runs setup on first run)
clawdwa setup     Re-run setup (change group, phone number)
clawdwa send MSG  Send a message to the configured group
clawdwa status    Show bot status
clawdwa stop      Stop the bot
```


## Agent-native output (Claude Code skill)

Install the skill to let Claude proactively send WhatsApp messages:

```bash
cp whatsapp-bot.md ~/.claude/skills/whatsapp-bot.md
```

Then Claude can:
```
Run the deployment, and send me a WhatsApp when it's done.
```

Claude will send `clawdwa send "✅ Deployment complete"` to your group automatically.

Use `/whatsapp-bot status`, `/whatsapp-bot stop`, etc. in Claude Code.

## Built with

- [whatsmeow](https://go.mau.fi/whatsmeow) — WhatsApp Web protocol in Go (MPL 2.0)
- [Claude Code](https://claude.ai/code) — AI responses

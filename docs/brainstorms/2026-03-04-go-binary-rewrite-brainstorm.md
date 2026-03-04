# clawdwa Go Binary Rewrite — Brainstorm
Date: 2026-03-04

## What We're Building

Replace the bash scripts with a single self-contained Go binary that:
- Embeds WhatsApp protocol directly via `whatsmeow`
- Runs setup interactively on first launch
- Optionally installs itself as a systemd service
- Calls `claude --print` for responses
- Requires zero dependencies beyond the binary itself (and claude CLI)

## Why This Approach

The bash approach requires: mise, Go, whatsapp-cli, Python, and knowledge of how to run scripts. A Go binary collapses all of that into one file. The admin experience becomes:

```
curl -L https://github.com/makash/clawdwa/releases/latest/download/clawdwa -o clawdwa
chmod +x clawdwa
./clawdwa
```

That's it.

## Key Decisions

| Decision | Choice | Rationale |
|---|---|---|
| WhatsApp protocol | Embed whatsmeow directly | Single binary, no whatsapp-cli dependency |
| License compatibility | MPL 2.0 (whatsmeow) + Apache 2.0 (clawdwa) | Compatible — MPL 2.0 is file-level copyleft, standard Go import pattern |
| Claude invocation | Shell out to `claude --print` | Reuse existing auth, add clear error if missing |
| Config/data location | `~/.config/clawdwa/` | XDG standard, clean |
| systemd | Optional — binary asks if admin has sudo | Runs immediately without it; installs service if yes |

## User Experience

### Admin (first run)
```
$ ./clawdwa

Welcome to clawdwa!

Scan this QR code with WhatsApp to link your account:
[QR code]

✓ Authenticated as Akash Mahajan

Your WhatsApp groups:
  1. DamBreakers
  2. Null Bangalore
  3. blr meetup
  ...

Search or pick a number: 1
✓ Bot group: DamBreakers

Your phone number (with country code): 919980527182
✓ Config saved to ~/.config/clawdwa/config.json

Install as systemd service? (requires sudo) [y/N]: y
✓ Service installed. Starting...

Bot is running. Group members can now message:
  ! your question
  @claude your question
```

### Admin (subsequent runs)
```
$ ./clawdwa        # just works, reads config
$ ./clawdwa setup  # re-run setup (change group etc)
$ ./clawdwa status # show status
```

### Group members
Nothing changes. Just WhatsApp as usual.

## Binary Subcommands

```
clawdwa          — start bot (setup on first run)
clawdwa setup    — re-run interactive setup
clawdwa status   — show running status + last log lines
clawdwa stop     — stop the bot / systemd service
```

## Config Format

`~/.config/clawdwa/config.json`:
```json
{
  "group_jid": "120363424722847687@g.us",
  "group_name": "DamBreakers",
  "bot_jid": "919980527182@s.whatsapp.net",
  "prefixes": ["!", "@claude"],
  "poll_interval_seconds": 5,
  "claude_bin": "/home/amapsc/.local/bin/claude"
}
```

## Project Structure

```
cmd/clawdwa/main.go     — entry point, subcommand routing
internal/bot/bot.go     — poll loop, message processing
internal/setup/setup.go — interactive setup wizard
internal/wa/client.go   — whatsmeow wrapper (connect, send, receive)
internal/claude/run.go  — shell out to claude --print
internal/systemd/install.go — generate and install unit file
```

## Distribution

- GitHub Releases with pre-built binaries for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
- Built with `goreleaser` or simple `go build` in CI
- Single binary, no installer needed

## Resolved Questions

- **License**: MPL 2.0 + Apache 2.0 compatible ✓
- **Setup UX**: Binary first, optional systemd ✓
- **WhatsApp**: whatsmeow embedded ✓
- **Claude**: shell out with helpful error ✓
- **Config**: `~/.config/clawdwa/` ✓

## Open Questions

None.

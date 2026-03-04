---
title: "feat: Rewrite clawdwa as a self-contained Go binary"
type: feat
status: active
date: 2026-03-04
---

# feat: Rewrite clawdwa as a self-contained Go binary

## Overview

Replace the current bash+whatsapp-cli setup with a single Go binary that embeds WhatsApp protocol directly via `whatsmeow`. Admin experience collapses to: download binary, run it.

**Key value propositions:**
- Zero dependencies beyond the binary (and `claude` CLI)
- Agent-native output: Claude can proactively SEND WhatsApp messages (e.g. "deploy and notify me on WhatsApp when done")
- Works as a Claude Code skill — send tasks to Claude from WhatsApp, get results back

## Project Structure

```
cmd/clawdwa/main.go           — entry point, subcommand routing
internal/bot/bot.go           — event loop, message processing
internal/setup/setup.go       — interactive setup wizard
internal/wa/client.go         — whatsmeow wrapper (connect, send, receive)
internal/claude/run.go        — shell out to claude --print
internal/systemd/install.go   — generate and install unit file
go.mod
go.sum
.goreleaser.yaml
```

Config: `~/.config/clawdwa/config.json`
WhatsApp state: `~/.config/clawdwa/wa.db` (whatsmeow SQLite store)
Bot dedup: `~/.config/clawdwa/bot.db`

## Acceptance Criteria

- [x] `./clawdwa` on first run: QR code → auth → group picker → phone number → optional systemd → bot running
- [x] `./clawdwa` on subsequent runs: reads config, starts bot immediately
- [x] `./clawdwa setup` — re-run wizard (change group etc.)
- [x] `./clawdwa status` — show running status + recent log lines
- [x] `./clawdwa stop` — stop bot / systemd service
- [x] `./clawdwa send "message"` — send a message to configured group (agent-native output)
- [x] Bot responds to `!` and `@claude` prefixes in configured group
- [x] Bot never replies to its own messages (sender JID filter)
- [x] Messages processed exactly once (SQLite dedup)
- [x] `claude` binary missing → clear error with install instructions
- [x] `goreleaser` builds linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
- [x] CGO_ENABLED=0 (pure Go, static binary) via modernc.org/sqlite
- [x] Update README to lead with agent-native value prop: "WhatsApp ↔ Claude. In both directions."

## Implementation Plan

### Phase 1: Scaffold

- [x] `go mod init github.com/makash/clawdwa`
- [x] Add dependencies: `go.mau.fi/whatsmeow`, `modernc.org/sqlite` (pure-Go, required for CGO_ENABLED=0), `github.com/mdp/qrterminal/v3`
- [x] `cmd/clawdwa/main.go`: parse os.Args[1] for subcommands, dispatch

**main.go skeleton:**
```go
func main() {
    sub := ""
    if len(os.Args) > 1 {
        sub = os.Args[1]
    }
    switch sub {
    case "setup":   setup.Run()
    case "status":  bot.Status()
    case "stop":    bot.Stop()
    case "send":    wa.Send(os.Args[2:])
    default:        run()   // if config.json missing → setup first, then start bot
    }
}
```

### Phase 2: WhatsApp Client (`internal/wa/client.go`)

- [x] Open whatsmeow SQLite store at `~/.config/clawdwa/wa.db`
- [x] Connect with reconnect on `events.Disconnected`
- [x] QR auth via `GetQRChannel` → render with `qrterminal`
- [x] `SendText(jid, text string) error` wrapper
- [x] `OnMessage(func(sender, text string))` callback registration

**Key whatsmeow patterns:**
```go
// Auth
qrChan, _ := client.GetQRChannel(ctx)
client.Connect()
for evt := range qrChan {
    if evt.Event == "code" {
        qrterminal.Generate(evt.Code, qrterminal.L, os.Stdout)
    }
}

// Receive
client.AddEventHandler(func(evt interface{}) {
    switch v := evt.(type) {
    case *events.Message:
        // process v.Message, v.Info.Sender, v.Info.Chat
    case *events.Disconnected:
        client.Connect()
    }
})

// Send
client.SendMessage(ctx, jid, &waProto.Message{
    Conversation: proto.String(text),
})
```

### Phase 3: Setup Wizard (`internal/setup/setup.go`)

- [x] Check `~/.config/clawdwa/config.json` — if exists, skip setup unless `setup` subcommand
- [x] Connect WA, show QR if not authenticated
- [x] Fetch joined groups: `client.GetJoinedGroups()`
- [x] Display numbered list with optional name filter (bufio.Scanner prompt)
- [x] Prompt for bot phone number (own JID) — format: `919900000000@s.whatsapp.net`
- [x] Write `config.json`
- [x] Ask: "Install as systemd service? (requires sudo) [y/N]" → call `systemd.Install()` if yes

**config.json:**
```json
{
  "group_jid": "120363424722847687@g.us",
  "group_name": "DamBreakers",
  "bot_jid": "919900000000@s.whatsapp.net",
  "prefixes": ["!", "@claude"],
  "claude_bin": "/home/amapsc/.local/bin/claude"
}
```

### Phase 4: Bot Loop (`internal/bot/bot.go`)

- [x] Open `~/.config/clawdwa/bot.db`, create `processed_messages(id TEXT PRIMARY KEY, processed_at INTEGER)` if not exists
- [x] Register `OnMessage` handler
- [x] Filter: only group messages where `chat.JID == config.GroupJID`
- [x] Filter: skip messages from `config.BotJID` (self-filter)
- [x] Check/insert into `processed_messages` (dedup)
- [x] Extract prompt: strip `! ` or `@claude ` prefix
- [x] Call `claude.Run(prompt)` → get response text
- [x] `wa.SendText(groupJID, response)`

### Phase 5: Claude Runner (`internal/claude/run.go`)

- [x] `exec.Command(config.ClaudeBin, "--print", prompt)`
- [x] Set `env -u CLAUDECODE` equivalent: filter out CLAUDECODE from `cmd.Env`
- [x] Capture stdout, return as string
- [x] If binary not found: print clear error with install instructions

```go
func Run(claudeBin, prompt string) (string, error) {
    cmd := exec.Command(claudeBin, "--print", prompt)
    cmd.Env = filterEnv(os.Environ(), "CLAUDECODE")
    out, err := cmd.Output()
    if err != nil {
        if errors.Is(err, exec.ErrNotFound) {
            return "", fmt.Errorf("claude not found at %s — install from https://claude.ai/code", claudeBin)
        }
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}
```

### Phase 6: systemd Install (`internal/systemd/install.go`)

- [x] Check if systemd available: `systemctl is-system-running` or check `/run/systemd/`
- [x] `IsAvailable() bool`
- [x] Generate unit file from template (embed binary path via `os.Executable()`)
- [x] Write to `/etc/systemd/system/clawdwa.service` (requires sudo — detect and use `sudo tee`)
- [x] `systemctl daemon-reload && systemctl enable --now clawdwa`

**Unit file template:**
```ini
[Unit]
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
```

### Phase 7: `send` Subcommand

- [x] `clawdwa send "message text"` — loads config, connects WA, sends to `config.GroupJID`, exits
- [x] Note: each `send` invocation reconnects (~2-3s WA handshake). This is acceptable for agent-native use (fire-and-forget notifications). Not a streaming path.
- [x] This is the agent-native output path — Claude Code can call this to notify users on WhatsApp

### Phase 8: Distribution

- [x] `.goreleaser.yaml` with:
  - `CGO_ENABLED=0` (requires modernc.org/sqlite for pure-Go SQLite)
  - Targets: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
  - `fetch-depth: 0` for goreleaser (needs git tags for version)
  - GitHub Releases with checksums

```yaml
# .goreleaser.yaml
builds:
  - env: [CGO_ENABLED=0]
    goos: [linux, darwin]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w -X main.version={{.Version}}
archives:
  - format: binary
    name_template: "clawdwa_{{ .Os }}_{{ .Arch }}"
checksum:
  name_template: "checksums.txt"
release:
  github:
    owner: makash
    name: clawdwa
```

## Dependencies

```
go.mau.fi/whatsmeow          — WhatsApp Web protocol (MPL 2.0)
modernc.org/sqlite           — Pure-Go SQLite (CGO_ENABLED=0)
github.com/mdp/qrterminal/v3 — QR code rendering in terminal
```

Note: `modernc.org/sqlite` instead of `mattn/go-sqlite3` enables `CGO_ENABLED=0` for static cross-compilation.

## License Note

- clawdwa: Apache 2.0
- whatsmeow: MPL 2.0 (file-level copyleft, compatible with Apache 2.0 via standard Go import)

## References

- Brainstorm: `docs/brainstorms/2026-03-04-go-binary-rewrite-brainstorm.md`
- whatsmeow: `go.mau.fi/whatsmeow` (godoc + example bots in repo)
- Current bash bot: `bot.sh`, `setup.sh`
- Current skill: `~/.claude/skills/whatsapp-bot.md`

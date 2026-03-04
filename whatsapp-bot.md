# whatsapp-bot

Manage the WhatsApp Claude bot at ~/whatsapp-claudecode.

## Subcommands

Detect the subcommand from the user's message or arguments:

- `setup`        — run the setup wizard
- `start`        — start the bot
- `stop`         — stop the bot
- `status`       — show bot status and recent log lines
- `send`         — send a message to the configured group
- `change-group` — stop bot, re-run setup to pick a new group, restart

---

## setup

Run the interactive setup wizard:

```bash
bash ~/whatsapp-claudecode/setup.sh
```

This authenticates WhatsApp, lets the admin choose a group, captures the phone number, writes `store/config.sh`, and initializes `store/bot.db`.

---

## start

Start the bot (manages sync internally — one command, no second terminal):

```bash
bash ~/whatsapp-claudecode/bot.sh
```

The bot runs in the foreground. Press Ctrl+C to stop.

---

## stop

Stop any running bot processes:

```bash
pkill -f "bash.*bot.sh" 2>/dev/null && echo "Bot stopped" || echo "Bot was not running"
```

---

## status

Show whether the bot is running and the last 20 log lines:

```bash
if pgrep -f "bash.*bot.sh" > /dev/null; then
  echo "✓ Bot is running (PID: $(pgrep -f 'bash.*bot.sh'))"
else
  echo "✗ Bot is not running"
fi
echo ""
echo "--- Last 20 log lines ---"
tail -20 ~/whatsapp-claudecode/bot.log 2>/dev/null || echo "(no log file yet)"
```

---

## send

Send a message directly to the configured WhatsApp group. Useful for Claude to proactively notify users.

Steps:
1. Source the config to get GROUP_JID and WA_BIN:
   ```bash
   source ~/whatsapp-claudecode/store/config.sh
   ```
2. Send the message:
   ```bash
   "$WA_BIN" --store "$STORE_DIR" send --to "$GROUP_JID" --message "<your message here>"
   ```

Example — Claude sending a proactive update:
```bash
source ~/whatsapp-claudecode/store/config.sh
"$WA_BIN" --store "$STORE_DIR" send --to "$GROUP_JID" --message "✅ Task complete: the deployment finished successfully."
```

---

## change-group

Stop the bot, re-run setup to pick a different group, then restart:

```bash
# 1. Stop the bot
pkill -f "bash.*bot.sh" 2>/dev/null && echo "Bot stopped" || echo "Bot was not running"

# 2. Re-run setup (overwrites store/config.sh with new group)
bash ~/whatsapp-claudecode/setup.sh

# 3. Restart the bot
bash ~/whatsapp-claudecode/bot.sh
```

The old group is simply replaced in `store/config.sh`. No data is lost — `store/bot.db` keeps all processed message IDs across group changes.

---

## Notes

- The bot only responds to messages starting with `!` or `@claude` in the configured group
- Bot's own messages are always ignored (no reply loops)
- Each message is processed independently (stateless)
- Config lives at `~/whatsapp-claudecode/store/config.sh`
- Logs at `~/whatsapp-claudecode/bot.log`

#!/usr/bin/env bash
# WhatsApp Claude Bot
# Usage: bash bot.sh
# Requires: setup.sh to have been run first

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/store/config.sh"

if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "ERROR: Config not found. Run setup first:"
  echo "  bash $SCRIPT_DIR/setup.sh"
  exit 1
fi

# shellcheck source=/dev/null
source "$CONFIG_FILE"

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

# Strip trigger prefix and return the clean prompt.
# Returns 1 if message is not a bot command.
extract_prompt() {
  local content="$1"
  local stripped
  # Support: "! ...", "!...", "@claude ..."
  if [[ "$content" =~ ^!\ (.+)$ ]]; then
    echo "${BASH_REMATCH[1]}"
  elif [[ "$content" =~ ^!([^[:space:]].+)$ ]]; then
    echo "${BASH_REMATCH[1]}"
  elif [[ "$content" =~ ^@claude\ (.+)$ ]]; then
    echo "${BASH_REMATCH[1]}"
  else
    return 1
  fi
}

is_processed() {
  python3 -c "
import sqlite3, sys
con = sqlite3.connect(sys.argv[1])
cur = con.execute('SELECT 1 FROM processed_messages WHERE id = ?', (sys.argv[2],))
print('yes' if cur.fetchone() else 'no')
con.close()
" "$BOT_DB" "$1"
}

mark_processed() {
  python3 -c "
import sqlite3, sys, time
con = sqlite3.connect(sys.argv[1])
con.execute('INSERT OR IGNORE INTO processed_messages (id, processed_at) VALUES (?, ?)',
            (sys.argv[2], int(time.time())))
con.commit()
con.close()
" "$BOT_DB" "$1"
}

process_message() {
  local id="$1" sender="$2" content="$3"

  # Skip own messages
  [[ "$sender" == "${BOT_JID%%@*}" ]] && return
  [[ "$sender" == "$BOT_JID" ]] && return

  # Skip if already processed
  [[ "$(is_processed "$id")" == "yes" ]] && return

  # Extract prompt (skip if not a bot command)
  local prompt
  prompt=$(extract_prompt "$content") || return

  mark_processed "$id"
  log "[$id] From $sender: $content"

  local response
  response=$(cd "$SCRIPT_DIR" && env -u CLAUDECODE "$CLAUDE_BIN" --print "$prompt" 2>&1) || {
    response="Sorry, I encountered an error processing your request."
  }

  log "[$id] Response: ${response:0:80}..."

  "$WA_BIN" --store "$STORE_DIR" send --to "$GROUP_JID" --message "$response"
  log "[$id] Sent to $GROUP_JID"
}

main() {
  log "Bot starting (group: $GROUP_JID, bot: $BOT_JID)"

  # Start sync in background; kill it when bot exits
  "$WA_BIN" --store "$STORE_DIR" sync &
  SYNC_PID=$!
  trap "log 'Stopping...'; kill $SYNC_PID 2>/dev/null; exit 0" INT TERM EXIT

  # Give sync a moment to connect
  sleep 3
  log "Polling every ${POLL_INTERVAL}s for messages prefixed with ! or @claude ..."

  while true; do
    while IFS= read -r json_line; do
      id=$(echo "$json_line"     | python3 -c "import sys,json; d=json.loads(sys.stdin.read()); print(d['id'])")
      sender=$(echo "$json_line" | python3 -c "import sys,json; d=json.loads(sys.stdin.read()); print(d['sender'])")
      content=$(echo "$json_line"| python3 -c "import sys,json; d=json.loads(sys.stdin.read()); print(d['content'])")
      process_message "$id" "$sender" "$content"
    done < <(python3 - "$STORE_DIR/messages.db" "$GROUP_JID" "$BOT_JID" << 'PYEOF'
import sqlite3, sys, json

db_path, group_jid, bot_jid = sys.argv[1], sys.argv[2], sys.argv[3]
bot_num = bot_jid.split("@")[0]

try:
    con = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    cur = con.execute("""
        SELECT id, sender, content FROM messages
        WHERE chat_jid = ?
          AND sender != ?
          AND content != ""
          AND (content LIKE '!%' OR content LIKE '@claude%')
        ORDER BY timestamp ASC
    """, (group_jid, bot_num))
    for row in cur.fetchall():
        print(json.dumps({"id": row[0], "sender": row[1], "content": row[2]}))
    con.close()
except Exception as e:
    print(f"DB error: {e}", file=sys.stderr)
PYEOF
)
    sleep "$POLL_INTERVAL"
  done
}

main

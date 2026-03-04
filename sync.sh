#!/usr/bin/env bash
# Standalone sync runner (optional — bot.sh manages sync automatically)
# Use this if you want sync running separately for debugging

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/store/config.sh"

if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "ERROR: Run setup.sh first"
  exit 1
fi

source "$CONFIG_FILE"
"$WA_BIN" --store "$STORE_DIR" sync

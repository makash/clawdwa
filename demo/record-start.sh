#!/usr/bin/env bash
# Record the bot-start demo.
# Usage: bash demo/record-start.sh

DEMO_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY="$DEMO_DIR/../clawdwa"

mkdir -p "$DEMO_DIR"

asciinema rec "$DEMO_DIR/start.cast" \
  --title "clawdwa start" \
  --command "expect $DEMO_DIR/start.exp $BINARY" \
  --overwrite

agg "$DEMO_DIR/start.cast" "$DEMO_DIR/start.gif" \
  --font-size 16 \
  --cols 80 \
  --rows 24

echo "Done: $DEMO_DIR/start.gif"

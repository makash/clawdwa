#!/usr/bin/env bash
# Record the setup demo.
# Usage: bash demo/record-setup.sh

DEMO_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY="$DEMO_DIR/../clawdwa"

mkdir -p "$DEMO_DIR"

asciinema rec "$DEMO_DIR/setup.cast" \
  --title "clawdwa setup" \
  --command "expect $DEMO_DIR/setup.exp $BINARY" \
  --overwrite

agg "$DEMO_DIR/setup.cast" "$DEMO_DIR/setup.gif" \
  --font-size 16 \
  --cols 80 \
  --rows 24

echo "Done: $DEMO_DIR/setup.gif"

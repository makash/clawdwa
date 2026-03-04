# WhatsApp Claude Bot — Architecture Brainstorm
Date: 2026-03-04

## What We're Building

A WhatsApp bot that lets a small trusted group query Claude Code via a dedicated WhatsApp group. Messages prefixed with `!` (or similar) trigger Claude; bot replies have no prefix, making loop prevention trivial. The bot runs on a Linux VM, polls the local WhatsApp SQLite DB, and sends responses back to the group.

## Why This Approach

- **Group-based membership** — add/remove people without touching code
- **Prefix trigger** — eliminates the reply-loop problem architecturally (no heuristics, no ID tracking)
- **SQLite direct read** — avoids the DB lock conflict between `sync` and `messages list`
- **Stateless per-message** — each `!` message is an independent Claude prompt; no session state to corrupt or lose

## Key Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Trigger mechanism | `!` prefix in a group | Loop-proof; group = easy membership control |
| Memory/context | None (stateless) | Less brittle; can add history later |
| DB access | Direct SQLite read | Avoids lock conflict with sync process |
| Deduplication | Message ID set in SQLite bot table | Bounded, fast, no flat-file growth |
| Users | Small trusted group | Group admin controls membership |

## Architecture

```
WhatsApp group (trusted members)
  → member sends "! what is X?"
  → whatsapp-cli sync writes to messages.db
  → bot polls messages.db every 5s
  → finds unprocessed messages WHERE content LIKE '!%'
  → strips prefix, sends to: claude --print "<prompt>"
  → sends reply back to group JID
  → marks message ID as processed in bot.db
```

## Bot DB Schema (SQLite)

```sql
CREATE TABLE processed_messages (
    id TEXT PRIMARY KEY,
    processed_at INTEGER NOT NULL
);
```

Simple, bounded (one row per processed message), queryable.

## Loop Prevention

Two-layer defense:
1. **Prefix filter** — only `!`-prefixed messages are processed; bot replies never start with `!`
2. **Processed ID table** — even if a message somehow passes the prefix check, it won't be processed twice

## Open Questions

None.

## Resolved Questions

- **Memory**: stateless per-message (less brittle)
- **Users**: small trusted group, managed via WA group membership
- **Trigger**: group + prefix (`!` or `@claude`)
- **Storage**: SQLite bot table for processed IDs
- **Prefix**: `!` or `@claude` — either triggers Claude, prompt stripped before sending
- **Reply target**: group (everyone sees Q&A)
- **Rate limiting**: none for now
- **Self-filter**: yes — bot's own JID is always ignored (belt-and-suspenders)

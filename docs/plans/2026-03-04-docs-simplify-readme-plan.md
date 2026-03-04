---
title: "docs: Simplify README for admin install and skill prominence"
type: docs
status: completed
date: 2026-03-04
---

# docs: Simplify README for admin install and skill prominence

Rewrite README.md so the admin can be running in under 2 minutes, group members know what to do in one line, and the Claude Code skill is treated as a core feature — not a footnote.

## Acceptance Criteria

- [x] README opens with tagline + concrete example (no scrolling needed to understand what it does)
- [x] Admin install is a single `curl` one-liner (auto-detects OS/arch) followed by `chmod +x && ./clawdwa`
- [x] "For group members" is one sentence — nothing to install, just `!` or `@claude`
- [x] Claude Code skill section is prominent, above the GIFs, clarifies it is admin-only
- [x] GIFs moved below the instructions (context before visuals)
- [x] Subcommands table retained
- [x] No legacy content or "(recommended)" remnants

## Proposed Structure

```
# clawdwa
tagline + 2-line example block

## For the admin
curl one-liner that auto-detects OS + arch, downloads correct binary from GitHub Releases
chmod +x clawdwa && ./clawdwa
First-run bullet list (QR, group picker, phone, optional systemd)
Requirement: Claude Code link

## tagline example block shows:
  Group member sends "! what is X?" → bot replies
  Claude sends clawdwa send "✅ done" → group receives

### Subcommands
table

## For group members
One sentence.

## Claude Code skill (admin only)
cp one-liner
What it enables: proactive sends + /whatsapp-bot commands

## Demo
setup.gif
start.gif

## Built with
```

## Context

- Brainstorm: `docs/brainstorms/2026-03-04-readme-simplification-brainstorm.md`
- Current README: `README.md`
- Skill file: `whatsapp-bot.md`
- Releases: https://github.com/makash/clawdwa/releases
- The skill is **admin-only** — requires `clawdwa` binary installed locally

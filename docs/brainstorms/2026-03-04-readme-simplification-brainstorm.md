# README Simplification — Brainstorm
Date: 2026-03-04

## What We're Building

A simplified README where:
- Admin can be up and running in under 2 minutes
- Group members know exactly what to do in one sentence
- The Claude Code skill is front and center (it's the killer feature)
- GIFs provide visual proof after the reader already understands the value

## Key Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Primary audience | Admin first | Admin sets everything up; group members learn from admin |
| Skill placement | In README, more prominent | Core admin feature, not a footnote |
| GIF placement | After setup instructions | Context before visuals |
| Subcommands | Keep | Useful reference, admin stays in README |
| curl one-liner | Add | Current README skips the download step |

## Proposed Structure

```
# clawdwa
WhatsApp ↔ Claude. In both directions.

[tagline examples]

## For the admin

curl one-liner download
chmod + run
What first run does (3 bullets)
Requirement: Claude Code

### Subcommands
table

## For group members
One sentence. Nothing to install.

## Claude Code skill
cp one-liner
What it unlocks (proactive sends, /whatsapp-bot commands)
Clarify: admin-only

## Demo
setup.gif
start.gif

## Built with
```

## Resolved Questions

- **Is the skill admin-only?** Yes — requires `clawdwa` binary installed locally. Group members just use WhatsApp.
- **Move skill section?** No — keep in README, make more prominent.
- **Keep subcommands?** Yes.

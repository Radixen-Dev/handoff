---
name: handoff
description: 'Store and retrieve knowledge packages between AI agent sessions using the handoff CLI. Use when: context window is filling up, user asks for a "handoff", "knowledge transfer", or "save context", starting a new session and checking for prior work, or ending a session with meaningful state to preserve.'
argument-hint: 'Optional: project name to scope the packages (e.g. my-app)'
---

# handoff — Knowledge Transfer

## When This Skill Applies

- **Session start**: Check for existing packages — `handoff list`
- **Context filling up**: Proactively offer to store a transfer
- **Explicit request**: "do a handoff", "save context", "knowledge transfer", "hand off to next agent"
- **Session end**: Offer to preserve state before wrapping up

---

## Installation Check

```bash
which handoff
```

If not found, install:

```bash
# Homebrew
brew tap Radixen-Dev/tap && brew install handoff

# Go
go install github.com/Radixen-Dev/handoff@latest
```

---

## Commands

### Store a package

Reads from stdin. Returns an 8-char hex ID. **Always report the ID to the user.**

```bash
echo "<markdown content>" | handoff store \
  --name "<name>" \
  --summary "<one-line description>" \
  --ttl <duration> \
  --project "<project-key>" \
  --tags "<tag1>,<tag2>"
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--name` | Yes | — | e.g. `auth-state`, `todo-api-progress` |
| `--summary` | No | — | One-line human-readable description |
| `--ttl` | No | `7d` | `Nh` = hours, `Nd` = days (e.g. `14d`) |
| `--project` | No | — | Repo/app name — use consistently for scoping |
| `--tags` | No | — | Comma-separated, e.g. `auth,api,decisions` |

### Retrieve a package

```bash
handoff retrieve <id>             # by exact ID
handoff retrieve --name <name>    # by name (most recent match)
```

### List packages

```bash
handoff list                      # all non-expired
handoff list --project <key>      # filter by project
```

Columns: `ID  NAME  PROJECT  TAGS  EXPIRES`

### Remove expired packages

```bash
handoff gc
```

Note: expired packages are also auto-purged on every operation — this is optional.

---

## TTL Reference

| Value | When to use |
|-------|-------------|
| `2h`–`1d` | Throwaway / short-lived context |
| `7d` | Default — typical working sessions |
| `14d` | Active multi-week projects |
| `30d` | Long-running or frequently revisited work |

---

## Environment

| Variable | Default | Description |
|----------|---------|-------------|
| `HANDOFF_DB` | `~/.handoff/handoff.db` | Override the database path |

---

## Knowledge Package Format

Structure stored content as markdown. Be thorough — the next agent has no other context.

```markdown
## Context
What the project is, the tech stack, the current goal.

## Key Decisions
- Decision: rationale and alternatives considered
- Decision: rationale and alternatives considered

## Current State
What is complete. What is in progress. What is blocked and why.

## Next Steps
1. First concrete action
2. Second concrete action
3. Third concrete action

## Warnings / Gotchas
- Known issues or fragile areas
- Things that look wrong but are intentional
- Setup steps easy to miss

## Files of Note
- `path/to/file` — what it does and why it matters
```

---

## Full Example

```bash
echo '## Context
REST API for a todo app. Stack: Go + Chi + PostgreSQL. Goal: task CRUD + JWT auth.

## Key Decisions
- PostgreSQL over SQLite for multi-user support
- JWT: 15-min access tokens + 7-day refresh tokens in httponly cookies

## Current State
Complete: auth, GET /tasks, POST /tasks.
In progress: PATCH /tasks/:id — validation not yet written.

## Next Steps
1. Finish PATCH /tasks/:id — add go-playground/validator
2. Implement DELETE /tasks/:id
3. Integration tests with testcontainers-go

## Warnings / Gotchas
- Run `make migrate` before testing — migrations in /migrations
- JWT_SECRET env var required at startup

## Files of Note
- `internal/auth/middleware.go` — JWT validation middleware
- `internal/task/handler.go` — all task HTTP handlers' \
  | handoff store \
    --name "todo-api-state" \
    --summary "Task CRUD in progress, auth complete" \
    --ttl 14d \
    --project "todo-api" \
    --tags "api,crud,auth"
```

Output: `a3f9c12e`

Tell the user: *"Stored as `a3f9c12e`. Retrieve in the next session with `handoff retrieve a3f9c12e` or `handoff retrieve --name todo-api-state`."*

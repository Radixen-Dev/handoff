# handoff

**Transfer knowledge between AI agent sessions.**

[![Release](https://img.shields.io/github/v/release/Dborasik/handoff)](https://github.com/Dborasik/handoff/releases)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

When an AI agent's context window fills up, accumulated project knowledge is lost. `handoff` solves this by letting agents **store** structured knowledge packages to a local SQLite database and **retrieve** them in a fresh session — no cloud, no config, no server.

```
Session A (context full)          Session B (fresh start)
        │                                  │
        ▼                                  ▼
  handoff store ──► ~/.handoff/handoff.db ◄── handoff retrieve
```

---

## Install

### Homebrew (recommended)

```bash
brew tap Dborasik/tap
brew install handoff
```

### Go Install

Requires Go 1.21+. No CGO, no system dependencies.

```bash
go install github.com/Dborasik/handoff@latest
```

### Build from Source

```bash
git clone https://github.com/Dborasik/handoff.git
cd handoff
go build -o handoff .
```

---

## Configuring Your Agent

All agent instruction files live in [`instructions/`](instructions/) and [`*/skills/handoff/`](.github/skills/handoff/) — organized, versioned, and curl-able. Drop whichever file your agent needs into your own project.

There are two approaches:

### Option A — Always-on instructions

The agent reads these on every session automatically.

| Agent | Add this file to your project | Source |
|-------|-------------------------------|--------|
| GitHub Copilot | `.github/copilot-instructions.md` | [`instructions/copilot-instructions.md`](instructions/copilot-instructions.md) |
| Claude Code | `CLAUDE.md` | [`instructions/CLAUDE.md`](instructions/CLAUDE.md) |
| OpenAI Codex | `AGENTS.md` | [`instructions/AGENTS.md`](instructions/AGENTS.md) |
| Cursor | `.cursor/rules/handoff.mdc` | [`instructions/cursor.mdc`](instructions/cursor.mdc) |

```bash
# GitHub Copilot
mkdir -p .github
curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/instructions/copilot-instructions.md \
  -o .github/copilot-instructions.md

# Claude Code
curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/instructions/CLAUDE.md -o CLAUDE.md

# OpenAI Codex
curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/instructions/AGENTS.md -o AGENTS.md

# Cursor
mkdir -p .cursor/rules
curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/instructions/cursor.mdc \
  -o .cursor/rules/handoff.mdc
```

### Option B — Skill file (on-demand)

A skill is loaded by the agent only when relevant — lighter weight, and works across all providers that support the `skills/<name>/SKILL.md` convention.

| Location in your project | Supported by |
|--------------------------|-------------|
| `.github/skills/handoff/SKILL.md` | GitHub Copilot |
| `.agents/skills/handoff/SKILL.md` | OpenAI Codex and other agents |
| `.claude/skills/handoff/SKILL.md` | Claude Code |

```bash
# GitHub Copilot
mkdir -p .github/skills/handoff
curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/.github/skills/handoff/SKILL.md \
  -o .github/skills/handoff/SKILL.md

# OpenAI Codex and other agents
mkdir -p .agents/skills/handoff
curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/.agents/skills/handoff/SKILL.md \
  -o .agents/skills/handoff/SKILL.md

# Claude Code
mkdir -p .claude/skills/handoff
curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/.claude/skills/handoff/SKILL.md \
  -o .claude/skills/handoff/SKILL.md
```

Once in place, the agent will automatically check for existing packages at the start of each session and offer to store context when things get long.

---

## Quick Start

**In Session A** — agent stores its context before the window fills:

```bash
echo "## Context
We're building a REST API for a todo app. Using Go + Chi router.

## Key Decisions
- Postgres over SQLite for multi-user support
- JWT auth, 15-min access tokens + 7-day refresh tokens
- Endpoints: POST /tasks, GET /tasks, PATCH /tasks/:id, DELETE /tasks/:id

## Current State
Auth is complete. Working on task CRUD." | handoff store --name "todo-api-state" --project "todo-api" --ttl 7d
```

Output:
```
a3f9c12e
```

**In Session B** — agent retrieves context instantly:

```bash
handoff retrieve a3f9c12e
# or by name:
handoff retrieve --name "todo-api-state"
```

---

## Commands

### `handoff store`

Reads content from stdin and stores it as a named knowledge package. Prints the package ID on success.

```bash
echo "<content>" | handoff store --name <name> [options]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--name` | *(required)* | Name for the package |
| `--summary` | | Short one-line summary |
| `--ttl` | `7d` | Time-to-live: `Nh` (hours) or `Nd` (days) |
| `--project` | | Project grouping key |
| `--tags` | | Comma-separated tags |

**Examples:**

```bash
# Minimal
echo "context..." | handoff store --name "my-notes"

# Full options
cat notes.md | handoff store \
  --name "auth-design" \
  --summary "JWT auth architecture notes" \
  --ttl 14d \
  --project "myapp" \
  --tags "auth,api,decisions"
```

---

### `handoff retrieve`

Retrieves a package and writes its content to stdout. Look up by ID or name.

```bash
handoff retrieve <id>
handoff retrieve --name <name>
```

When using `--name`, the most recently stored package with that name is returned.

**Examples:**

```bash
handoff retrieve a3f9c12e
handoff retrieve --name "auth-design"

# Pipe directly into a file or another tool
handoff retrieve --name "auth-design" > context.md
```

---

### `handoff list`

Lists all non-expired packages in a table.

```bash
handoff list
handoff list --project <key>
```

**Example output:**

```
ID        NAME           PROJECT   TAGS          EXPIRES
a3f9c12e  auth-design    myapp     auth,api      2026-06-09 14:30
b1d2e3f4  db-schema      myapp     db,postgres   2026-06-02 09:15
```

| Flag | Description |
|------|-------------|
| `--project` | Filter results to a specific project |

---

### `handoff gc`

Manually removes all expired packages. Expired packages are also automatically removed on every database operation, so running this is optional.

```bash
handoff gc
# Removed 3 expired package(s).
```

---

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `HANDOFF_DB` | `~/.handoff/handoff.db` | Path to the SQLite database file |

Override the database location:

```bash
export HANDOFF_DB=/tmp/session.db
```

---

## How It Works

- **Storage**: A single SQLite file at `~/.handoff/handoff.db`, created automatically on first use. Pure Go — no CGO, no system libraries required.
- **Package IDs**: 8-character hex strings (e.g. `a3f9c12e`), randomly generated at store time.
- **TTL**: Specified at creation as a duration (`7d`, `2h`, etc.) and stored as an absolute expiry timestamp. Every database operation deletes expired rows before running, so GC is automatic.
- **No daemon, no config file**: The only state is the SQLite file.

---

## Agent Workflow

This is the intended usage pattern for AI coding agents:

1. **Agent A's context is filling up.** The user says: *"Do a knowledge transfer before we lose context."*
2. **Agent A** composes a markdown summary of the current project state.
3. **Agent A** runs:
   ```bash
   echo '<summary>' | handoff store --name "project-state" --ttl 14d
   ```
   and reports the ID back to the user (e.g. `a3f9c12e`).
4. **User starts a new session** with Agent B.
5. **User says:** *"Retrieve knowledge package `a3f9c12e`"* (or by name).
6. **Agent B** runs:
   ```bash
   handoff retrieve a3f9c12e
   ```
   and reads the full context, resuming work immediately.

---

## Uninstall

```bash
brew uninstall handoff
```

This removes the binary. `handoff` also creates a data directory at `~/.handoff/` the first time it runs — Homebrew does not touch this. To remove it too:

```bash
rm -rf ~/.handoff
```

> **Note:** `~/.handoff/handoff.db` contains all your stored knowledge packages. Only delete it if you're sure you no longer need them. If you used a custom path via `HANDOFF_DB`, remove that file instead.

---

## License

MIT — see [LICENSE](LICENSE).

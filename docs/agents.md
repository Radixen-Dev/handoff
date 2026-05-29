# Agent Setup

For `handoff` to be useful, your agent needs to know it exists. This page explains how to wire it in — both the file to add and what the agent will do once it has those instructions.

There are two approaches. Choose based on how often you expect to use `handoff` in the project:

!!! info "Windows users"
    The `curl` commands below use standard flags that work with `curl.exe` on Windows 10+. In PowerShell, `curl` is an alias for `Invoke-WebRequest` — use `curl.exe` explicitly to invoke the real curl binary. The `mkdir -p` flag is not needed on Windows; use `mkdir` without it (PowerShell's `mkdir` creates intermediate directories and does not error if the directory already exists).

| | Option A — Always-on | Option B — Skill file |
|--|---|---|
| **When loaded** | Every session, automatically | Only when the agent judges it relevant |
| **Context cost** | Small constant overhead per session | Near zero when not in use |
| **Best for** | Projects where context limits are a regular problem | Projects where `handoff` is only occasionally needed |
| **Cursor support** | :material-check: | :material-close: Not available |

---

## Option A — Always-on instruction files

Add a single file to your project. The agent reads it at the start of every session and will automatically check for existing packages, offer proactive transfers, and respond to explicit requests.

| Agent | File to create | Source |
|-------|---------------|--------|
| Claude Code | `CLAUDE.md` | [`instructions/CLAUDE.md`](https://github.com/Radixen-Dev/handoff/blob/main/instructions/CLAUDE.md) |
| GitHub Copilot | `.github/copilot-instructions.md` | [`instructions/copilot-instructions.md`](https://github.com/Radixen-Dev/handoff/blob/main/instructions/copilot-instructions.md) |
| OpenAI Codex | `AGENTS.md` | [`instructions/AGENTS.md`](https://github.com/Radixen-Dev/handoff/blob/main/instructions/AGENTS.md) |
| Cursor | `.cursor/rules/handoff.mdc` | [`instructions/cursor.mdc`](https://github.com/Radixen-Dev/handoff/blob/main/instructions/cursor.mdc) |

Use `curl` to pull the file directly into your project:

=== "Claude Code"

    ```bash
    curl -fsSL \
      https://raw.githubusercontent.com/Radixen-Dev/handoff/main/instructions/CLAUDE.md \
      -o CLAUDE.md
    ```

=== "GitHub Copilot"

    ```bash
    mkdir -p .github
    curl -fsSL \
      https://raw.githubusercontent.com/Radixen-Dev/handoff/main/instructions/copilot-instructions.md \
      -o .github/copilot-instructions.md
    ```

=== "OpenAI Codex"

    ```bash
    curl -fsSL \
      https://raw.githubusercontent.com/Radixen-Dev/handoff/main/instructions/AGENTS.md \
      -o AGENTS.md
    ```

=== "Cursor"

    ```bash
    mkdir -p .cursor/rules
    curl -fsSL \
      https://raw.githubusercontent.com/Radixen-Dev/handoff/main/instructions/cursor.mdc \
      -o .cursor/rules/handoff.mdc
    ```

!!! tip "Already have a CLAUDE.md or AGENTS.md?"
    Don't replace your existing file — append the `handoff` instructions to it. The instructions are plain markdown and compose cleanly with any existing content.

---

## Option B — Skill files

A skill file is loaded by the agent only when it determines the skill is relevant — typically when you ask for a handoff or when the context is getting large. This approach costs nothing in sessions where `handoff` is not needed.

!!! info "Cursor not supported"
    Cursor does not use skill files. Use Option A (above) instead.

| Location in your project | Supported by |
|--------------------------|-------------|
| `.github/skills/handoff/SKILL.md` | GitHub Copilot |
| `.agents/skills/handoff/SKILL.md` | OpenAI Codex and other agents |
| `.claude/skills/handoff/SKILL.md` | Claude Code |

=== "Claude Code"

    ```bash
    mkdir -p .claude/skills/handoff
    curl -fsSL \
      https://raw.githubusercontent.com/Radixen-Dev/handoff/main/.claude/skills/handoff/SKILL.md \
      -o .claude/skills/handoff/SKILL.md
    ```

=== "GitHub Copilot"

    ```bash
    mkdir -p .github/skills/handoff
    curl -fsSL \
      https://raw.githubusercontent.com/Radixen-Dev/handoff/main/.github/skills/handoff/SKILL.md \
      -o .github/skills/handoff/SKILL.md
    ```

=== "OpenAI Codex / other"

    ```bash
    mkdir -p .agents/skills/handoff
    curl -fsSL \
      https://raw.githubusercontent.com/Radixen-Dev/handoff/main/.agents/skills/handoff/SKILL.md \
      -o .agents/skills/handoff/SKILL.md
    ```

---

## What the agent will do

Once either option is in place, the agent behaves as follows:

**At the start of every session**
:   The agent runs `handoff list` (or `handoff list --project <name>`) to check for existing packages. If any are relevant to the current work, it retrieves and reads them before proceeding.

**During long sessions**
:   When the conversation is getting large, the agent proactively offers to store a knowledge transfer:
    > *"Our context is getting large — would you like me to store a knowledge transfer for the next session?"*

**When explicitly asked**
:   Any of the following phrases trigger an immediate store:
    `"do a handoff"`, `"save context"`, `"knowledge transfer"`, `"store session state"`, `"hand off to next agent"`

**Before closing a session**
:   The agent offers to preserve state before wrapping up, so nothing is lost.

**If `handoff` is not installed**
:   The agent checks with `which handoff` before using it. If it is not found, the agent guides the user through installation rather than failing silently.

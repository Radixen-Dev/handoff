# handoff

> **Transfer knowledge between AI agent sessions — local SQLite, no cloud, no server, no configuration.**

Every AI coding agent operates inside a finite context window. When that window fills up, the session ends and everything accumulated during it — decisions made, architecture discussed, progress tracked, gotchas discovered — is gone. The next agent starts from zero.

`handoff` solves this. Before the window closes, the agent stores a structured markdown summary to a local SQLite database. The next agent retrieves it and resumes work with full context, instantly.

![handoff flow diagram — Session A stores to handoff.db, Session B retrieves from handoff.db](https://raw.githubusercontent.com/Dborasik/handoff/main/assets/flow.png)

---

## Quick start

**Session A** — store context before the window fills:

```bash
echo "## Current State
Auth is complete. Working on task CRUD.

## Next Steps
1. Finish PATCH /tasks/:id
2. Add pagination to GET /tasks" | handoff store --name "api-state" --ttl 7d
```

```text
a3f9c12e
```

**Session B** — fresh agent retrieves it immediately:

```bash
handoff retrieve a3f9c12e
# or by name:
handoff retrieve --name "api-state"
```

The full context is written to stdout. The agent reads it and picks up where the last session left off.

---

## How it fits together

<div class="grid cards" markdown>

-   :material-database-outline: **Local SQLite**

    ---

    All data lives in `~/.handoff/handoff.db` on your own machine. Nothing leaves your system. No accounts, no API keys, no internet connection required.

-   :material-timer-sand: **Auto-expiry**

    ---

    Every package has a TTL. Packages expire silently in the background — old context never accumulates. Default TTL is 7 days.

-   :material-package-variant: **Single binary**

    ---

    Pure Go, no CGO, no runtime dependencies. One executable that works the same on macOS, Linux, and Windows.

-   :material-robot-outline: **Works with all major agents**

    ---

    Ready-made instruction files for Claude Code, GitHub Copilot, Cursor, and OpenAI Codex. Drop one file into your project and the agent handles the rest.

</div>

---

## Supported agents

`handoff` ships instruction files for every major AI coding agent:

| Agent | Supported |
|-------|-----------|
| Claude Code | :material-check-circle: Yes |
| GitHub Copilot | :material-check-circle: Yes |
| Cursor | :material-check-circle: Yes |
| OpenAI Codex | :material-check-circle: Yes |
| Any agent with terminal access | :material-check-circle: Yes |

See [Agent Setup](agents.md) for installation instructions.

---

## Explore the docs

<div class="grid cards" markdown>

-   :material-download: **Install**

    ---

    Homebrew, `go install`, or build from source.

    [:octicons-arrow-right-24: Install](install.md)

-   :material-console: **Commands**

    ---

    Full reference for `store`, `retrieve`, `list`, and `gc`.

    [:octicons-arrow-right-24: Commands](commands.md)

-   :material-robot: **Agent Setup**

    ---

    Wire `handoff` into your agent's instructions in one `curl` command.

    [:octicons-arrow-right-24: Agent Setup](agents.md)

-   :material-source-branch: **Workflow**

    ---

    The recommended handoff pattern and knowledge package format.

    [:octicons-arrow-right-24: Workflow](workflow.md)

-   :material-cog-outline: **How It Works**

    ---

    Database schema, ID generation, TTL, garbage collection, and design principles.

    [:octicons-arrow-right-24: How It Works](internals.md)

</div>

# Workflow

A handoff involves two sessions. Session A stores context before its window is exhausted. Session B retrieves that context and resumes work immediately. This page covers how those two sessions coordinate, how to structure the content of a package, and a complete worked example.

---

## The two-session pattern

### Session A — storing context

!!! example "Trigger phrases"
    These phrases all mean the same thing to an agent with `handoff` instructions: *compose and store a knowledge package right now.*

    `"do a handoff"` · `"save context"` · `"knowledge transfer"` · `"store session state"` · `"hand off to next agent"`

**Step 1 — The agent composes a summary.**
The agent writes a structured markdown document covering the current project state: what was decided, what is done, what is in progress, what comes next, and any gotchas to be aware of. See [Package structure](#package-structure) below for the recommended format.

**Step 2 — The agent stores it.**

```bash
echo '<summary>' | handoff store \
  --name "project-state" \
  --summary "One-line description of current state" \
  --ttl 14d \
  --project "myapp"
```

**Step 3 — The agent reports the ID.**
The command prints an 8-character hex ID. The agent passes it back to the user:

> *"Stored as `a3f9c12e`. In the next session, retrieve it with `handoff retrieve a3f9c12e` or `handoff retrieve --name project-state`."*

---

### Session B — retrieving context

**Step 1 — Check what's available.**
With [always-on instructions](agents.md) in place, the agent does this automatically. It can also be triggered manually:

```bash
handoff list
handoff list --project myapp
```

**Step 2 — Retrieve the package.**

```bash
# By ID (exact, always returns the right package)
handoff retrieve a3f9c12e

# By name (returns the most recently stored package with that name)
handoff retrieve --name "project-state"
```

**Step 3 — Resume work.**
The agent reads the content and has full context. No explanation from the user required.

---

## Package structure

`handoff` does not enforce any content structure — it stores whatever markdown you pipe in. The format below is a recommended template. It is designed around one principle: **the agent reading this package in the next session has no other context**. Write for that reader.

```markdown
## Context
What this project is and what we are currently building.
Include the tech stack, the overall goal, and any background
a new agent needs to understand the project from scratch.

## Key Decisions
- Decision: rationale and any alternatives that were considered
- Decision: rationale and any alternatives that were considered

## Current State
What is fully complete.
What is currently in progress and how far along it is.
What is blocked and why.

## Next Steps
1. First concrete action to take (specific, not vague)
2. Second concrete action to take
3. Third concrete action to take

## Warnings / Gotchas
- Known issues or fragile areas
- Things that look wrong but are intentional
- Setup steps that are easy to miss

## Files of Note
- `path/to/file.go` — what it does and why it matters
- `path/to/other.ts` — what it does and why it matters
```

### Why each section matters

**Context**
:   A fresh agent knows nothing about your project. Don't assume it has read the README or explored the codebase. Describe the project as if briefing someone for the first time.

**Key Decisions**
:   The most overlooked section. Architectural choices made during the session — technology selections, API design, trade-offs accepted — exist only in the conversation. Without documenting them, the next agent may re-examine or silently reverse decisions that were already settled.

**Current State**
:   Distinguishes what is production-ready from what is half-done. The next agent needs to know exactly where the work boundary is before writing a single line.

**Next Steps**
:   Write specific, actionable items — not vague goals. *"Add input validation to `PATCH /tasks/:id` using `go-playground/validator`"* is useful. *"Fix the tasks endpoint"* is not.

**Warnings / Gotchas**
:   Traps that cost time to rediscover. Include anything non-obvious, environment-specific, or that caused confusion during the session. This section has a disproportionate impact on the quality of the next session.

**Files of Note**
:   The next agent will explore the codebase. Point it directly to the relevant files so it doesn't have to infer them from the directory structure.

---

## Full example

A complete, realistic knowledge package with all sections filled in:

=== "macOS / Linux (bash/zsh)"

    ```bash
    handoff store \
      --name "todo-api-state" \
      --summary "Task CRUD in progress, auth complete, tests pending" \
      --ttl 14d \
      --project "todo-api" \
      --tags "api,crud,auth,go" << 'EOF'
    ## Context
    Building a REST API for a multi-user todo application.
    Stack: Go + Chi router + PostgreSQL.
    Goal: full CRUD for tasks with JWT authentication.
    Repo: github.com/example/todo-api

    ## Key Decisions
    - PostgreSQL over SQLite: needed for concurrent multi-user access
    - JWT auth: 15-min access tokens + 7-day refresh tokens in HttpOnly cookies
    - Chi router over Gin: lighter footprint, closer to the standard library
    - No ORM: using raw database/sql with pgx driver for clarity and control

    ## Current State
    Complete:
      - User registration and login
      - JWT middleware (applied per-route, not globally)
      - GET /tasks and POST /tasks

    In progress:
      - PATCH /tasks/:id — handler exists, input validation not yet written

    Not started:
      - DELETE /tasks/:id
      - Pagination on GET /tasks
      - Integration tests

    Blocked: nothing currently blocked.

    ## Next Steps
    1. Finish PATCH /tasks/:id — add validation with go-playground/validator
    2. Implement DELETE /tasks/:id
    3. Write integration tests using testcontainers-go
    4. Add cursor-based pagination to GET /tasks

    ## Warnings / Gotchas
    - Always run `make migrate` before testing — migrations live in /migrations/
    - JWT_SECRET env var must be set or the server panics on startup (see main.go:42)
    - The test database runs on port 5433, not 5432, to avoid conflicts with local Postgres
    - Refresh token rotation logic is in internal/auth/refresh.go, not middleware.go

    ## Files of Note
    - `internal/auth/middleware.go` — JWT validation middleware, applied per-route
    - `internal/auth/refresh.go` — refresh token rotation and revocation logic
    - `internal/task/handler.go` — all task HTTP handlers
    - `migrations/` — SQL migration files, applied in ascending filename order
    EOF
    ```

=== "Windows (PowerShell)"

    ```powershell
    @'
    ## Context
    Building a REST API for a multi-user todo application.
    Stack: Go + Chi router + PostgreSQL.
    Goal: full CRUD for tasks with JWT authentication.
    Repo: github.com/example/todo-api

    ## Key Decisions
    - PostgreSQL over SQLite: needed for concurrent multi-user access
    - JWT auth: 15-min access tokens + 7-day refresh tokens in HttpOnly cookies
    - Chi router over Gin: lighter footprint, closer to the standard library
    - No ORM: using raw database/sql with pgx driver for clarity and control

    ## Current State
    Complete:
      - User registration and login
      - JWT middleware (applied per-route, not globally)
      - GET /tasks and POST /tasks

    In progress:
      - PATCH /tasks/:id — handler exists, input validation not yet written

    Not started:
      - DELETE /tasks/:id
      - Pagination on GET /tasks
      - Integration tests

    Blocked: nothing currently blocked.

    ## Next Steps
    1. Finish PATCH /tasks/:id — add validation with go-playground/validator
    2. Implement DELETE /tasks/:id
    3. Write integration tests using testcontainers-go
    4. Add cursor-based pagination to GET /tasks

    ## Warnings / Gotchas
    - Always run `make migrate` before testing — migrations live in /migrations/
    - JWT_SECRET env var must be set or the server panics on startup (see main.go:42)
    - The test database runs on port 5433, not 5432, to avoid conflicts with local Postgres
    - Refresh token rotation logic is in internal/auth/refresh.go, not middleware.go

    ## Files of Note
    - `internal/auth/middleware.go` — JWT validation middleware, applied per-route
    - `internal/auth/refresh.go` — refresh token rotation and revocation logic
    - `internal/task/handler.go` — all task HTTP handlers
    - `migrations/` — SQL migration files, applied in ascending filename order
    '@ | handoff store `
      --name "todo-api-state" `
      --summary "Task CRUD in progress, auth complete, tests pending" `
      --ttl 14d `
      --project "todo-api" `
      --tags "api,crud,auth,go"
    ```

The agent in Session B retrieves this package, reads it, and has everything needed to continue working on `PATCH /tasks/:id` without asking any clarifying questions.

---

## Tips

!!! tip "Store early, not just when full"
    Don't wait until the context window is completely exhausted. A package stored with 20% headroom is more coherent and easier for the next agent to act on than one written under pressure at the very limit.

!!! tip "Reuse names intentionally"
    Package names are not unique. You can store a new package under the same name as a previous one — `retrieve --name` will always return the most recent. Use this to maintain a single "current state" name per project that gets overwritten on each transfer, while keeping older snapshots accessible by ID.

!!! tip "Use --project consistently"
    Setting `--project` to the repository or app name on every store makes `handoff list --project <name>` immediately useful. Without it, all packages from all projects appear together in `list` output.

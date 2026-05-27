# How It Works

This page covers the internals of `handoff`: where data lives, how the database is structured, how IDs are generated, how expiry works, and the design principles the tool is built around.

---

## Storage

All data is stored in a single SQLite file. The location depends on your operating system:

=== "macOS / Linux"

    ```text
    ~/.handoff/handoff.db
    ```

=== "Windows"

    ```text
    %USERPROFILE%\.handoff\handoff.db
    ```

    For example: `C:\Users\Alice\.handoff\handoff.db`

The directory is created automatically the first time any `handoff` command runs. No setup step is required.

`handoff` uses [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) — a pure-Go SQLite driver that bundles the SQLite engine directly into the binary. There is no CGO dependency and no requirement to have SQLite installed on the host system.

---

## Database schema

```sql
CREATE TABLE IF NOT EXISTS packages (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    summary    TEXT DEFAULT '',
    content    TEXT NOT NULL,
    tags       TEXT DEFAULT '[]',
    project    TEXT DEFAULT '',
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_packages_name
    ON packages(name);

CREATE INDEX IF NOT EXISTS idx_packages_project
    ON packages(project);

CREATE INDEX IF NOT EXISTS idx_packages_expires_at
    ON packages(expires_at);
```

### Column reference

`id`
:   8-character lowercase hex string. Primary key. Generated at store time using `crypto/rand`.

`name`
:   The name given by the `--name` flag. **Not unique** — multiple packages can share the same name. The most recently created package with a given name is returned by `retrieve --name`.

`summary`
:   Optional single-line description from `--summary`. Defaults to an empty string. Not used for lookup.

`content`
:   The full package body, stored verbatim as the text piped into `handoff store`.

`tags`
:   Stored as a JSON array string, e.g. `["auth","api","decisions"]`. Parsed from the comma-separated `--tags` value at store time.

`project`
:   Optional project scoping key from `--project`. Defaults to an empty string.

`created_at` / `expires_at`
:   UTC datetimes. `expires_at` is computed at store time as `created_at + TTL duration`. Both are stored and queried in UTC.

### Indexes

Three indexes are maintained to keep operations fast regardless of database size:

| Index | Used by |
|-------|---------|
| `idx_packages_name` | `retrieve --name` lookups |
| `idx_packages_project` | `list --project` filtering |
| `idx_packages_expires_at` | Garbage collection (`DELETE WHERE expires_at < NOW()`) |

---

## Package IDs

IDs are generated using Go's [`crypto/rand`](https://pkg.go.dev/crypto/rand) package:

1. 4 bytes are read from the OS cryptographic random source
2. The bytes are hex-encoded, producing an 8-character lowercase string (e.g. `a3f9c12e`)

This gives 2³² (approximately 4.3 billion) possible IDs. Given that a single user's local database typically holds a handful of packages at any time, the probability of a collision is negligible.

IDs are assigned at store time and never change.

---

## TTL and garbage collection

### How TTL works

TTL is specified at store time as a duration string and immediately converted to an absolute `expires_at` UTC timestamp. It cannot be changed after the package is stored.

**Supported formats:**

| Format | Example | Duration |
|--------|---------|---------|
| `Nh` | `2h` | N hours |
| `Nd` | `14d` | N days |

N must be a positive integer. `0h`, `-1d`, and similar are rejected with an error.

### Lazy garbage collection

`handoff` runs no background process. Instead, it garbage-collects expired packages *lazily* — at the beginning of every database operation:

```sql
DELETE FROM packages WHERE expires_at < <current UTC time>
```

This runs on every call to `store`, `retrieve`, `list`, and `gc`. The result is that expired packages are silently removed during normal use, and you will never see an expired package in `list` output or receive one from `retrieve`.

### Explicit GC

```bash
handoff gc
# Removed 3 expired package(s).
```

Running `handoff gc` triggers the same deletion and reports how many rows were removed. Use it when you want to confirm cleanup has occurred or reclaim disk space immediately.

---

## Configuration

`handoff` has exactly one configuration surface: an environment variable for the database path.

| Variable | Default (macOS / Linux) | Default (Windows) |
|----------|------------------------|-------------------|
| `HANDOFF_DB` | `~/.handoff/handoff.db` | `%USERPROFILE%\.handoff\handoff.db` |

There is no config file. No other settings exist.

=== "macOS / Linux"

    ```bash
    # Use a temporary database for a throwaway session (inline)
    HANDOFF_DB=/tmp/scratch.db handoff store --name "temp" --ttl 2h

    # Set for the current shell session
    export HANDOFF_DB="$(pwd)/.handoff.db"

    # Use a shared path
    export HANDOFF_DB=/var/shared/handoff.db
    ```

=== "Windows (PowerShell)"

    ```powershell
    # Use a temporary database for a throwaway session
    $env:HANDOFF_DB = "C:\Temp\scratch.db"
    handoff store --name "temp" --ttl 2h

    # Set for the current shell session
    $env:HANDOFF_DB = "$PWD\.handoff.db"

    # Use a shared path
    $env:HANDOFF_DB = "C:\Shared\handoff.db"
    ```

=== "Windows (CMD)"

    ```cmd
    set HANDOFF_DB=C:\Temp\scratch.db
    handoff store --name "temp" --ttl 2h
    ```

---

## Design principles

`handoff` was deliberately kept minimal. Every design decision reflects one of these principles:

1. **Zero config.** Works out of the box with no setup. The first command you run creates the database automatically.
2. **Single binary.** One executable, no runtime dependencies, no install scripts.
3. **Cross-platform.** Pure Go, no CGO. Identical behaviour on macOS, Linux, and Windows for both amd64 and arm64.
4. **Stdin/stdout.** Content flows through pipes. This is the natural interface for agents working in a terminal, and it makes `handoff` composable with other tools.
5. **Ephemeral by default.** TTL ensures stored context eventually disappears. Old packages do not accumulate forever.
6. **Simple over clever.** No plugin system, no network layer, no authentication, no daemon. A single SQLite file is the entire data layer.

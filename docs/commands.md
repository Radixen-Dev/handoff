# Commands

`handoff` exposes five subcommands. All database operations share the same SQLite file at `~/.handoff/handoff.db` by default, configurable via [`HANDOFF_DB`](internals.md#configuration).

| Command | Description |
|---------|-------------|
| [`store`](#handoff-store) | Read content from stdin and save it as a named knowledge package |
| [`retrieve`](#handoff-retrieve) | Fetch a package by ID or name and write its content to stdout |
| [`list`](#handoff-list) | List all non-expired packages in a table |
| [`gc`](#handoff-gc) | Manually delete all expired packages |
| [`completion`](#handoff-completion) | Generate a shell completion script for bash, zsh, fish, or PowerShell |

---

## handoff store

Reads content from stdin and saves it as a named knowledge package. On success, prints an 8-character hex package ID to stdout.

```bash
echo "<content>" | handoff store --name <name> [flags]
```

### Flags

`--name` *(required)*
:   A short, descriptive name for the package. Used to retrieve it by name in a future session. Names are **not unique** — multiple packages can share the same name. `retrieve --name` always returns the most recently stored match, so you can store successive snapshots under the same name and always get the latest.

`--summary`
:   A single-line human-readable description of the package contents. Not used for lookup — only appears in `list` output.

`--ttl` *(default: `7d`)*
:   How long the package survives before it is automatically deleted. Format: a positive integer followed by `h` (hours) or `d` (days). For example: `2h`, `1d`, `7d`, `14d`, `30d`. The TTL is converted to an absolute expiry timestamp at store time and cannot be changed afterwards.

`--project`
:   An optional project scoping key. Use the same value consistently across a project to group related packages. Acts as a filter argument for `handoff list --project`.

`--tags`
:   Optional comma-separated labels. Leading and trailing whitespace is trimmed from each tag. Example: `auth,api,decisions`.

!!! info "Content comes from stdin"
    There is no flag for the package body. Content must be piped in. The command exits with an error if stdin is empty or contains only whitespace.

### Output

On success, the package ID is printed to stdout followed by a newline:

```text
a3f9c12e
```

Pass this ID back to the user so they can retrieve the package by ID in the next session.

### Examples

```bash
# Minimal — name and content only
echo "notes about current state" | handoff store --name "scratch"

# All options
cat notes.md | handoff store \
  --name "auth-design" \
  --summary "JWT auth architecture decisions" \
  --ttl 14d \
  --project "myapp" \
  --tags "auth,api,decisions"

# Multi-line content via heredoc
handoff store --name "sprint-state" --ttl 7d --project "myapp" << 'EOF'
## Current State
Feature X is complete. Feature Y is in progress.

## Next Steps
1. Finish Y
2. Write tests
EOF
```

### Errors

| Error message | Cause |
|---------------|-------|
| `Error: --name is required` | The `--name` flag was omitted |
| `Error: no content provided on stdin` | Stdin was empty or contained only whitespace |
| `Error: invalid --ttl: too short` | TTL string was fewer than 2 characters |
| `Error: invalid --ttl: invalid number: <value>` | The numeric part of the TTL could not be parsed |
| `Error: invalid --ttl: unsupported unit '<unit>' (use h or d)` | TTL unit was not `h` or `d` |
| `Error: invalid --ttl: must be positive` | TTL value was `0` or negative |

---

## handoff retrieve

Retrieves a package and writes its full content to stdout. Lookup is either by exact package ID or by name.

```bash
handoff retrieve <id>
handoff retrieve --name <name>
```

### Flags

`<id>` *(positional argument)*
:   The exact 8-character hex package ID returned by `store`. When a positional argument is provided, it takes priority over `--name` even if both are given.

`--name`
:   Retrieve by package name. If multiple packages share the same name, the most recently stored one is returned. Expired packages are never returned.

!!! note
    At least one of a positional ID or `--name` must be provided. Providing both is valid — the positional ID is used.

### Output

The full package content is written to stdout exactly as it was stored. No metadata is prepended or appended. Use shell redirection to capture it:

```bash
handoff retrieve --name "auth-design" > context.md
```

### Examples

```bash
# Retrieve by exact ID
handoff retrieve a3f9c12e

# Retrieve by name (returns the most recent if the name has been reused)
handoff retrieve --name "auth-design"

# Save directly to a file
handoff retrieve a3f9c12e > context.md
```

Pipe to clipboard:

=== "macOS"

    ```bash
    handoff retrieve --name "sprint-state" | pbcopy
    ```

=== "Linux"

    ```bash
    handoff retrieve --name "sprint-state" | xclip -selection clipboard
    # or
    handoff retrieve --name "sprint-state" | xsel --clipboard
    ```

=== "Windows (PowerShell)"

    ```powershell
    handoff retrieve --name "sprint-state" | clip
    ```

### Errors

| Error message | Cause |
|---------------|-------|
| `Error: provide a package ID as argument or use --name` | Neither a positional ID nor `--name` was provided |
| `Error: package not found` | No non-expired package matched the given ID or name |

---

## handoff list

Lists all non-expired packages as a tab-aligned table, ordered from most recently stored to oldest.

```bash
handoff list
handoff list --project <key>
```

### Flags

`--project`
:   Filter results to packages whose project key matches the given value exactly. If omitted, all non-expired packages are shown regardless of project.

### Output

```text
ID        NAME           PROJECT   TAGS          EXPIRES
a3f9c12e  auth-design    myapp     auth,api      2026-06-09 14:30
b1d2e3f4  db-schema      myapp     db,postgres   2026-06-02 09:15
```

| Column | Description |
|--------|-------------|
| `ID` | 8-character package ID |
| `NAME` | Package name as given at store time |
| `PROJECT` | Project scoping key, or blank if not set |
| `TAGS` | Comma-separated tags, or blank if not set |
| `EXPIRES` | Expiry timestamp in local time, formatted `YYYY-MM-DD HH:MM` |

Packages whose `expires_at` has already passed are automatically deleted before the listing is produced, so they never appear here.

!!! info "No packages"
    If the database is empty, or no packages match the given `--project` filter, `No packages found.` is printed to stderr and the command exits with status `0`.

### Examples

```bash
# List all packages across all projects
handoff list

# List only packages for a specific project
handoff list --project myapp
```

### Errors

| Error message | Cause |
|---------------|-------|
| `No packages found.` *(stderr, exit 0)* | The database is empty or no packages match the `--project` filter — not a failure |

No other errors occur during normal use. Database access failures (e.g. permission denied on `~/.handoff/`) produce a message to stderr and exit with status `1`.

---

## handoff gc

Deletes all expired packages from the database and prints the number removed.

```bash
handoff gc
```

Expired packages are already deleted automatically at the start of every database operation, so running `gc` explicitly is never required. Use it when you want to confirm how many stale packages exist or reclaim disk space on demand.

### Flags

This command takes no flags.

### Output

```text
Removed 3 expired package(s).
```

If no packages were expired at the time of the call, the output is:

```text
Removed 0 expired package(s).
```

### Examples

```bash
# Run a manual garbage collection
handoff gc

# Check how many packages would be cleaned up
handoff gc
```

### Errors

No user-facing errors under normal operation. Database access failures produce a message to stderr and exit with status `1`.

---

## handoff completion

Generates a shell autocompletion script and writes it to stdout. This is a standard Cobra-generated utility and is not specific to `handoff`'s functionality.

```bash
handoff completion <shell>
```

Supported shells: `bash`, `zsh`, `fish`, `powershell`.

### Flags

`<shell>` *(positional argument, required)*
:   The target shell. Must be one of: `bash`, `zsh`, `fish`, `powershell`.

`--no-descriptions`
:   Disable completion descriptions. Supported by `bash`, `fish`, and `zsh`.

### Output

The completion script is written to stdout. Redirect it to the appropriate file for your shell to load it on startup.

### Examples

=== "bash"

    ```bash
    # Load for current session only
    source <(handoff completion bash)

    # Install permanently
    handoff completion bash > /etc/bash_completion.d/handoff
    ```

=== "zsh"

    ```bash
    # Load for current session only
    source <(handoff completion zsh)

    # Install permanently
    handoff completion zsh > "${fpath[1]}/_handoff"
    ```

=== "fish"

    ```bash
    handoff completion fish | source

    # Install permanently
    handoff completion fish > ~/.config/fish/completions/handoff.fish
    ```

=== "PowerShell"

    ```powershell
    handoff completion powershell | Out-String | Invoke-Expression
    ```

### Errors

| Error message | Cause |
|---------------|-------|
| `Error: unknown command "<shell>" for "handoff completion"` | The shell argument was not one of the four supported values |

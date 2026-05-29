# Install

`handoff` ships as a single static binary with no runtime dependencies. Choose the method that suits your setup.

---

## Methods

=== "Homebrew"

    !!! info "macOS and Linux only"
        Homebrew is not available on Windows. Windows users should use [Go install](#go-install) or a [pre-built binary](#pre-built-binaries).

    ```bash
    brew tap Radixen-Dev/tap
    brew install handoff
    ```

    To upgrade to the latest release:

    ```bash
    brew upgrade handoff
    ```

=== "Go install"

    Requires Go 1.26 or later. Run `go version` to check. The binary is fully self-contained — no CGO, no system libraries.

    ```bash
    go install github.com/Radixen-Dev/handoff@latest
    ```

    The binary is placed in your `GOPATH/bin` directory. Make sure that directory is on your `PATH`:

    === "macOS / Linux"

        ```bash
        export PATH="$PATH:$(go env GOPATH)/bin"
        ```

        Add this line to `~/.zshrc` or `~/.bashrc` to make it permanent.

    === "Windows (PowerShell)"

        ```powershell
        # Add for current session
        $env:PATH += ";$(go env GOPATH)\bin"

        # Add permanently (user-level)
        [Environment]::SetEnvironmentVariable(
            "PATH",
            $env:PATH + ";$(go env GOPATH)\bin",
            "User"
        )
        ```

=== "Build from source"

    Requires Go 1.26 or later.

    ```bash
    git clone https://github.com/Radixen-Dev/handoff.git
    cd handoff
    ```

    === "macOS / Linux"

        ```bash
        go build -o handoff .
        mv handoff /usr/local/bin/
        ```

    === "Windows (PowerShell)"

        ```powershell
        go build -o handoff.exe .
        Move-Item handoff.exe "$env:GOPATH\bin\handoff.exe"
        ```

---

## Verify

After installing, confirm the binary is available:

```bash
handoff --help
```

You should see:

```text
A CLI tool for storing and retrieving knowledge packages across AI agent context windows.

Usage:
  handoff [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  gc          Remove all expired packages
  help        Help about any command
  list        List available knowledge packages
  retrieve    Retrieve a knowledge package by ID or name
  store       Store a knowledge package (reads content from stdin)

Flags:
  -h, --help   help for handoff

Use "handoff [command] --help" for more information about a command.
```

---

## Pre-built binaries

Every [GitHub release](https://github.com/Radixen-Dev/handoff/releases) includes pre-built archives. A `checksums.txt` file is provided for verification.

| Platform | Architecture | Archive |
|----------|-------------|---------|
| macOS | Apple Silicon (arm64) | `handoff_darwin_arm64.tar.gz` |
| macOS | Intel (amd64) | `handoff_darwin_amd64.tar.gz` |
| Linux | x86-64 (amd64) | `handoff_linux_amd64.tar.gz` |
| Linux | ARM64 | `handoff_linux_arm64.tar.gz` |
| Windows | x86-64 (amd64) | `handoff_windows_amd64.zip` |
| Windows | ARM64 | `handoff_windows_arm64.zip` |

=== "macOS"

    ```bash
    # Apple Silicon
    curl -L https://github.com/Radixen-Dev/handoff/releases/latest/download/handoff_darwin_arm64.tar.gz \
      | tar -xz
    mv handoff /usr/local/bin/

    # Intel
    curl -L https://github.com/Radixen-Dev/handoff/releases/latest/download/handoff_darwin_amd64.tar.gz \
      | tar -xz
    mv handoff /usr/local/bin/
    ```

=== "Linux"

    ```bash
    # x86-64
    curl -L https://github.com/Radixen-Dev/handoff/releases/latest/download/handoff_linux_amd64.tar.gz \
      | tar -xz
    mv handoff /usr/local/bin/

    # ARM64
    curl -L https://github.com/Radixen-Dev/handoff/releases/latest/download/handoff_linux_arm64.tar.gz \
      | tar -xz
    mv handoff /usr/local/bin/
    ```

=== "Windows (PowerShell)"

    ```powershell
    # x86-64
    Invoke-WebRequest `
      -Uri https://github.com/Radixen-Dev/handoff/releases/latest/download/handoff_windows_amd64.zip `
      -OutFile handoff.zip
    Expand-Archive handoff.zip -DestinationPath .
    Move-Item handoff.exe "$env:GOPATH\bin\handoff.exe"

    # ARM64
    Invoke-WebRequest `
      -Uri https://github.com/Radixen-Dev/handoff/releases/latest/download/handoff_windows_arm64.zip `
      -OutFile handoff.zip
    Expand-Archive handoff.zip -DestinationPath .
    Move-Item handoff.exe "$env:GOPATH\bin\handoff.exe"
    ```

    !!! tip
        Move `handoff.exe` to any directory that is already on your `%PATH%`. If you are unsure, `$env:GOPATH\bin` is a good choice if you have Go installed. Otherwise use `C:\Windows\System32\` (requires admin) or create a dedicated `bin` folder in your home directory and add it to your user PATH.

---

## Uninstall

### Remove the binary

=== "Homebrew (macOS / Linux)"

    ```bash
    brew uninstall handoff
    ```

=== "Manual (macOS / Linux)"

    ```bash
    rm $(which handoff)
    ```

=== "Windows (PowerShell)"

    ```powershell
    Remove-Item (Get-Command handoff).Source
    ```

### Remove stored packages

The binary uninstall does not touch your stored packages. The data directory must be removed separately.

=== "macOS / Linux"

    ```bash
    rm -rf ~/.handoff
    ```

=== "Windows (PowerShell)"

    ```powershell
    Remove-Item -Recurse -Force "$env:USERPROFILE\.handoff"
    ```

!!! warning "This permanently deletes all knowledge packages"
    The database file contains all packages stored with `handoff store`. Only delete it if you are certain you no longer need them. If you used a custom path via `HANDOFF_DB`, remove that file instead.

---
title: Installation
---

# Installation and Setup

This guide will help you install Load-Pulse and verify that it's working correctly.

## Prerequisites

Before installing Load-Pulse, ensure you have the following tools installed:

- **[Go](https://go.dev/dl/)** (version 1.21 or later) – Required for building and running Load-Pulse
- **[Docker](https://www.docker.com/get-started)** – Containerization engine used by Load-Pulse for running distributed load tests
- **[Git](https://git-scm.com/)** – Version control system (optional, but recommended)

## Installation Methods

### Option 1: go install (Recommended)

If you have Go installed and your `GOBIN`/`GOPATH/bin` is on your system `PATH`:

```bash
go install github.com/Naganathan05/Load-Pulse@latest
```

This installs Load-Pulse globally, giving you a `loadpulse` command that you can run from any directory.

**Note:** Make sure `$GOPATH/bin` or `$GOBIN` is in your `PATH`. You can verify this by running:

```bash
# On Linux/macOS
echo $PATH | grep -o '[^:]*go[^:]*'

# On Windows PowerShell
$env:PATH -split ';' | Select-String -Pattern 'go'
```

### Option 2: Clone and Build Locally

If you prefer to build from source or want to contribute:

```bash
git clone https://github.com/Naganathan05/Load-Pulse.git
cd Load-Pulse
go build -o loadpulse ./main.go
```

**On Windows PowerShell:**

```powershell
git clone https://github.com/Naganathan05/Load-Pulse.git
cd Load-Pulse
go build -o loadpulse.exe .\main.go
```

After building, you can either:
- Add the binary to your `PATH`
- Run it directly: `./loadpulse` (Linux/macOS) or `.\loadpulse.exe` (Windows)

### Option 3: Run Without Installation

You can also run Load-Pulse directly without installing:

```bash
git clone https://github.com/Naganathan05/Load-Pulse.git
cd Load-Pulse
go run ./main.go --help
```

**On Windows PowerShell:**

```powershell
git clone https://github.com/Naganathan05/Load-Pulse.git
cd Load-Pulse
go run .\main.go --help
```

## Verifying Installation

After installation, verify that Load-Pulse is working correctly:

1. **Check the version:**

   ```bash
   loadpulse version
   ```

   This should display the installed Load-Pulse version.

2. **View available commands:**

   ```bash
   loadpulse --help
   ```

   This should list all available commands: `init`, `validate`, `run`, `clean`, and `version`.

3. **Ensure Docker is running:**

   Load-Pulse requires Docker to be running. Verify Docker is available:

   ```bash
   docker --version
   docker ps
   ```

   If Docker is not running, start Docker Desktop (or your Docker daemon) before using Load-Pulse.

## Next Steps

Once Load-Pulse is installed and verified:

1. **[Create a test configuration](/commands#loadpulse-init)** using `loadpulse init`
2. **[Validate your configuration](/commands#loadpulse-validate)** using `loadpulse validate`
3. **[Run your first load test](/commands#loadpulse-run)** using `loadpulse run`

For detailed information about each command, see the [Commands documentation](/commands).


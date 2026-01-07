---
title: CLI Commands
---

# CLI Commands and Usage

This page documents all Load-Pulse CLI commands, their options, and example usage.

All examples assume `loadpulse` is installed and available on your `PATH`. If you're running from source, replace `loadpulse` with `go run ./main.go` (or `go run .\main.go` on Windows).

## Command Overview

Load-Pulse provides the following commands:

- [`loadpulse init`](#loadpulse-init) – Create a test configuration file interactively
- [`loadpulse validate`](#loadpulse-validate) – Validate a test configuration file
- [`loadpulse run`](#loadpulse-run) – Execute a load test
- [`loadpulse clean`](#loadpulse-clean) – Clean up Docker containers
- [`loadpulse version`](#loadpulse-version) – Display version information

## `loadpulse init`

Interactively create a `testConfig.json` configuration file.

### Usage

```bash
loadpulse init
```

### Description

The `init` command launches an interactive wizard that guides you through creating a test configuration file. It prompts you for:

- Output file path (defaults to `testConfig.json` in the current directory)
- Base host URL for your API
- Test duration (in seconds)
- Request definitions (method, endpoint, data, connections, rate, concurrency limit)

### Example

```bash
$ loadpulse init
Enter output JSON file path [testConfig.json]: 
This wizard will create a testConfig configuration file at: 
testConfig.json

# Follow the interactive prompts...
```

### When to Use

Use `loadpulse init` when:
- Setting up a new load test configuration
- You want a guided, interactive way to create `testConfig.json`
- You're unfamiliar with the `testConfig.json` structure

For manual configuration, see the [Config Reference](/config-reference) page.

### Output

The command creates a `testConfig.json` file (or the path you specified) in the current directory with your configuration.

---

## `loadpulse validate`

Check that a test configuration file is well-formed and complete.

### Usage

```bash
loadpulse validate [testConfig-file]
```

### Arguments

- `testConfig-file` (optional) – Path to the test configuration JSON file. Defaults to `testConfig.json` in the current directory if not specified.

### Description

The `validate` command checks:
- That the file exists and is readable
- That the JSON is valid and properly formatted
- That all required fields are present
- That field types and values are correct

### Examples

```bash
# Validate the default testConfig.json
loadpulse validate

# Validate a specific file
loadpulse validate path/to/testConfig.json

# Validate a file in a different directory
loadpulse validate ../my-tests/api-test.json
```

### When to Use

Use `loadpulse validate` when:
- Before running a long load test to catch configuration errors early
- After manually editing `testConfig.json`

### Output

On success:
```
[INFO]: Validating test configuration file: testConfig.json
[INFO]: testConfig configuration is valid: testConfig.json
```

On failure, the command will display a error messages.

---

## `loadpulse run`

Run a load test using the given configuration file.

### Usage

```bash
loadpulse run --config <path-to-testConfig.json>
```

### Options

- `--config`, `-c` – Path to the test configuration JSON file (default: `testConfig.json`)

### Description

The `run` command:
1. Validates that Docker is running
2. Reads the test configuration file
3. Spins up Docker containers (load tester, aggregator, Redis, RabbitMQ)
4. Executes the load test according to your configuration
5. Displays aggregated results from the aggregator
6. Automatically cleans up containers after completion

### Examples

```bash
# Run with default testConfig.json
loadpulse run

# Run with explicit config file
loadpulse run --config testConfig.json

# Run with config file using short flag
loadpulse run -c my-test.json

# Run with config in a different directory
loadpulse run --config ../tests/api-load-test.json
```

### Prerequisites

- Docker must be running (Docker Desktop or Docker daemon)
- A valid `testConfig.json` file (use `loadpulse validate` to check)

### When to Use

Use `loadpulse run` when:
- You're ready to execute a load test
- You've validated your configuration
- Your target API server is running and accessible

### Output

During execution, you'll see:
- Initialization messages
- Docker container startup status
- A spinner indicating test progress
- Final aggregated statistics for each endpoint

For details on interpreting results, see the [Results documentation](/results).

---

## `loadpulse clean`

Stop and clean up Docker containers created by Load-Pulse.

### Usage

```bash
loadpulse clean
```

### Description

The `clean` command:
- Stops all containers started by Load-Pulse
- Removes containers and associated volumes
- Cleans up resources from previous test runs

### Examples

```bash
# Clean up containers
loadpulse clean
```

### When to Use

Use `loadpulse clean` when:
- Containers from a previous run didn't shut down properly
- You want to free up Docker resources
- You're troubleshooting container-related issues

**Note:** The `loadpulse run` command automatically cleans up containers after a successful test run, so manual cleanup is usually not necessary.

### Output

```
[INFO]: Starting cleanup process ...
[INFO]: Container Cleanup Successfully Completed
```

If no containers are found:
```
[INFO]: No containers found to clean up
```

---

## `loadpulse version`

Print the installed Load-Pulse version.

### Usage

```bash
loadpulse version
```

### Description

Displays the version number of the installed Load-Pulse CLI.

### Example

```bash
$ loadpulse version
Load Pulse version: 1.0.0
```

### When to Use

Use `loadpulse version` to:
- Verify your installation
- Check which version you're running
- Report version information when seeking help

---

## Getting Help

For help with any command, use the `--help` flag:

```bash
loadpulse --help
loadpulse init --help
loadpulse run --help
# etc.
```

## Common Workflows

### First-Time Setup

```bash
# 1. Create a test configuration
loadpulse init

# 2. Validate it
loadpulse validate

# 3. Run the test
loadpulse run
```

### Running Multiple Tests

```bash
# Create and validate first test
loadpulse init
loadpulse validate

# Run first test
loadpulse run --config testConfig.json

# Create second test configuration
loadpulse init  # Specify different filename when prompted

# Run second test
loadpulse run --config testConfig-2.json
```

### Troubleshooting

```bash
# Check version
loadpulse version

# Validate configuration
loadpulse validate testConfig.json

# Clean up if containers are stuck
loadpulse clean

# Run test again
loadpulse run
```


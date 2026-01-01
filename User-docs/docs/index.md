# Load-Pulse

Load testing tool built in Go which works based on the Raft consensus
algorithm. Designed for scalable and distributed benchmarking using
Dockerized microservices.

## High-level architecture
![Load-Pulse high-level architecture](/images/Architecture_diagram.png)

## Installing the CLI

### Option 1: go install (recommended)

If you have Go installed and your `GOBIN`/`GOPATH/bin` is on `PATH`:

```bash
go install github.com/Naganathan05/Load-Pulse@latest
```

This gives you a global `loadpulse` command you can run from any
project directory.

### Option 2: clone and run locally

```bash
git clone https://github.com/Naganathan05/Load-Pulse.git
cd Load-Pulse
go run ./main.go --help
```

On Windows PowerShell from the repo root:

```powershell
go run .\main.go --help
```

## The test configuration file (testConfig.json)

Load-Pulse reads its test definition from a JSON file, called
`testConfig.json`. It describes which endpoints to hit and how.

High-level structure:

- `host` – Base URL for your API, for example:
  - `"http://localhost:8080/"`
  - `"http://host.docker.internal:8081/"`
- `duration` – How long the test should run (in seconds).
- `requests` – Array of request definitions. Each item has:
  - `method` – HTTP method, e.g. `"GET"`, `"POST"`.
  - `endpoint` – Path relative to `host`, e.g. `"api/admin/getAllDepartments"`.
  - `data` – Optional request body as a string (used mainly for `POST`/`PUT`).
  - `connections` – Number of HTTP connections to open for this request.
  - `rate` – Delay between requests in milliseconds.
  - `concurrencyLimit` – Maximum concurrent in‑flight requests for this endpoint.

This file can be manually created or using the interactive tool too.

## Core commands

All examples assume `loadpulse` is on your PATH.

### `loadpulse init`

Interactively create a `testConfig.json` file:

```bash
loadpulse init
```

Use this when you want to set up a new test configuration (host,
duration, endpoints, etc.) interactively.

### `loadpulse validate`

Check that a test configuration file is well‑formed and complete:

```bash
loadpulse validate path/to/testConfig.json
```

Use this before running a long test to catch mistakes in the config.

### `loadpulse run`

Run a load test using the given configuration file:

```bash
loadpulse run --config path/to/testConfig.json
```

Use this from your project repository to actually execute the load
test defined in your `testConfig.json`.

### `loadpulse clean`

Stop and clean up any containers created by previous runs:

```bash
loadpulse clean
```

### `loadpulse version`

Print the installed Load-Pulse version:

```bash
loadpulse version
```

## View load testing results

Once a test run completes, Load-Pulse prints a summary of results
from the aggregator. The key metrics include:

- Average response time
- Max response time
- Min response time
- Total requests sent
- Successful vs failed requests





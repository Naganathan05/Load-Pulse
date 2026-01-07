---
title: Configuration Reference
---

# testConfig.json Reference

This page provides a complete reference for the `testConfig.json` configuration file used by Load-Pulse.

## Overview

`testConfig.json` is a JSON file that defines your load test parameters, including the target server, test duration, and the specific HTTP requests to execute. You can create this file manually or use the [`loadpulse init`](/commands#loadpulse-init) command to generate it interactively.

## File Structure

```json
{
  "host": "string",
  "duration": number,
  "requests": [
    {
      "method": "string",
      "endpoint": "string",
      "data": "string",
      "connections": number,
      "rate": number,
      "concurrencyLimit": number
    }
  ]
}
```

## Top-Level Fields

### `host`

**Type:** `string`  
**Required:** Yes  
**Description:** The base URL for your API server.

**Examples:**
```json
"host": "http://localhost:8080/"
"host": "http://host.docker.internal:8081/"
"host": "https://api.example.com/"
```

**Important Notes:**
- Must include the protocol (`http://` or `https://`)
- Must include a trailing slash (`/`)
- For localhost servers accessed from Docker containers, use `host.docker.internal` instead of `localhost`
- The host is prepended to each request's `endpoint` field

### `duration`

**Type:** `number` (integer)  
**Required:** Yes  
**Description:** How long the load test should run, specified in **seconds**.

**Examples:**
```json
"duration": 10    // Run for 10 seconds
"duration": 60    // Run for 1 minute
"duration": 300   // Run for 5 minutes
```

**Considerations:**
- Longer durations provide more stable averages but take more time
- Shorter durations are useful for quick smoke tests
- The test will run for exactly this duration, then stop

### `requests`

**Type:** `array` of request objects  
**Required:** Yes  
**Description:** An array of HTTP request definitions. Each request will be executed according to its configuration.

**Minimum:** At least one request object is required.

## Request Object Fields

Each object in the `requests` array defines a single HTTP endpoint to test.

### `method`

**Type:** `string`  
**Required:** Yes  
**Description:** The HTTP method to use for this request.

**Valid Values:**
- `"GET"`
- `"POST"`

**Example:**
```json
"method": "GET"
"method": "POST"
```

### `endpoint`

**Type:** `string`  
**Required:** Yes  
**Description:** The API endpoint path, relative to the `host` URL.

**Examples:**
```json
"endpoint": "api/users"
"endpoint": "api/admin/getAllDepartments"
"endpoint": "v1/products/123"
"endpoint": "health"
```

**Important Notes:**
- Do **not** include a leading slash (it's added automatically)
- Do **not** include query parameters here (include them in the endpoint path if needed)
- The full URL will be: `{host}{endpoint}`

**Full URL Examples:**
- Host: `"http://localhost:8080/"` + Endpoint: `"api/users"` = `http://localhost:8080/api/users`
- Host: `"https://api.example.com/"` + Endpoint: `"v1/products"` = `https://api.example.com/v1/products`

### `data`

**Type:** `string`  
**Required:** No (can be empty string)  
**Description:** The HTTP request body, typically used for `POST`, `PUT`, or `PATCH` requests.

**Examples:**
```json
"data": ""                                    // Empty for GET requests
"data": "{\"name\":\"John\",\"age\":30}"     // JSON string
"data": "name=John&age=30"                    // Form-encoded data
```

**Important Notes:**
- Must be a string (even if it contains JSON)
- For JSON payloads, escape quotes: `"{\"key\":\"value\"}"`
- Leave empty (`""`) for GET, HEAD, DELETE requests that don't need a body
- The content type is not explicitly set; ensure your server can parse the format you provide

### `connections`

**Type:** `number` (integer)  
**Required:** Yes  
**Description:** The number of HTTP connections to open for this endpoint.

**Examples:**
```json
"connections": 1      // Single connection
"connections": 10     // 10 concurrent connections
"connections": 100    // 100 concurrent connections
```

**What it means:**
- Each connection can send multiple requests over time
- More connections = more concurrent load
- Higher values simulate more simultaneous users

**Considerations:**
- Start with lower values (1-10) and increase gradually
- Very high values may overwhelm your server or network
- Each connection consumes resources on both client and server

### `rate`

**Type:** `number` (integer)  
**Required:** Yes  
**Description:** Delay between requests in **milliseconds**.

**Examples:**
```json
"rate": 1        // 1ms delay = very high rate (~1000 req/s per connection)
"rate": 100      // 100ms delay = 10 requests per second per connection
"rate": 1000     // 1 second delay = 1 request per second per connection
```

**Request Rate Calculation:**
```
Requests per second = 1000 / rate × connections
```

**Examples:**
- `rate: 100, connections: 10` = ~100 requests/second
- `rate: 1000, connections: 5` = 5 requests/second
- `rate: 1, connections: 50` = ~50,000 requests/second (very high!)

**Considerations:**
- Lower `rate` values = higher request frequency
- `rate: 1` sends requests as fast as possible (be careful!)
- Balance `rate` and `connections` to achieve desired load
- Consider your server's capacity when setting these values

### `concurrencyLimit`

**Type:** `number` (integer)  
**Required:** Yes  
**Description:** Maximum number of concurrent in-flight requests for this endpoint.

**Examples:**
```json
"concurrencyLimit": 1      // One request at a time
"concurrencyLimit": 10     // Up to 10 concurrent requests
"concurrencyLimit": 100    // Up to 100 concurrent requests
```

**What it means:**
- Limits how many requests can be "in progress" simultaneously
- Prevents overwhelming the server with too many concurrent requests
- Workers will wait if the limit is reached before sending new requests

**Relationship with `connections`:**
- `connections` determines how many HTTP connections are opened
- `concurrencyLimit` caps how many requests can be active at once
- If `concurrencyLimit < connections`, some connections may be idle

**Considerations:**
- Lower values provide more controlled load
- Higher values simulate more aggressive load patterns
- Should typically be ≤ `connections` for optimal resource usage

## Complete Example

Here's a complete `testConfig.json` example:

```json
{
  "host": "http://host.docker.internal:8081/",
  "duration": 30,
  "requests": [
    {
      "method": "GET",
      "endpoint": "api/users",
      "data": "",
      "connections": 10,
      "rate": 100,
      "concurrencyLimit": 5
    },
    {
      "method": "POST",
      "endpoint": "api/users",
      "data": "{\"name\":\"John Doe\",\"email\":\"john@example.com\"}",
      "connections": 5,
      "rate": 200,
      "concurrencyLimit": 3
    },
    {
      "method": "GET",
      "endpoint": "api/admin/getAllDepartments",
      "data": "",
      "connections": 20,
      "rate": 50,
      "concurrencyLimit": 10
    }
  ]
}
```

**What this configuration does:**
- Tests three endpoints on `http://host.docker.internal:8081/`
- Runs for 30 seconds
- GET `/api/users`: 10 connections, 100ms delay, max 5 concurrent
- POST `/api/users`: 5 connections, 200ms delay, max 3 concurrent
- GET `/api/admin/getAllDepartments`: 20 connections, 50ms delay, max 10 concurrent

## Configuration Tips

### Starting Out

1. **Begin with low values:**
   ```json
   {
     "connections": 1,
     "rate": 1000,
     "concurrencyLimit": 1
   }
   ```

2. **Gradually increase load:**
   - First increase `connections`
   - Then decrease `rate` (increase frequency)
   - Adjust `concurrencyLimit` based on your needs

### Testing Different Scenarios

**Light Load:**
```json
"connections": 5,
"rate": 500,
"concurrencyLimit": 3
```

**Moderate Load:**
```json
"connections": 20,
"rate": 100,
"concurrencyLimit": 10
```

**Heavy Load:**
```json
"connections": 100,
"rate": 10,
"concurrencyLimit": 50
```

### Common Pitfalls

1. **Using `localhost` in Docker:** Use `host.docker.internal` instead
2. **Missing trailing slash:** Always include `/` at the end of `host`
3. **JSON in `data` field:** Remember to escape quotes: `"{\"key\":\"value\"}"`
4. **Too aggressive settings:** Start conservative and increase gradually
5. **Mismatched limits:** `concurrencyLimit` should typically be ≤ `connections`

## Validation

Before running a test, validate your configuration:

```bash
loadpulse validate testConfig.json
```

This will catch:
- Invalid JSON syntax
- Missing required fields
- Invalid field types or values
- Configuration errors

## Next Steps

- [Create a configuration interactively](/commands#loadpulse-init) using `loadpulse init`
- [Validate your configuration](/commands#loadpulse-validate) before running tests
- [Run your load test](/commands#loadpulse-run) with `loadpulse run`
- [Understand your results](/results) after the test completes


---
title: Test Results
---

# Understanding Load Test Results

This page explains how to interpret the results reported by Load-Pulse after a test run completes.

## Overview

When a load test completes, Load-Pulse aggregates statistics from all worker nodes and displays a summary for each endpoint tested. The results are printed to the console and include key performance metrics.

## Result Format

Results are displayed in a formatted output with colored text for better readability. Each endpoint gets its own statistics block:

```
[LOG]: Test completed for endpoint: /api/users
------------------------------ STATS --------------------------------
[STATS]: Total requests completed: 1000
[STATS]: Total Number of Error Requests: 5
[STATS]: Average response size: 1024.5 bytes
[STATS]: Min response time: 45ms
[STATS]: Max response time: 1250ms
[STATS]: Average response time: 234ms
---------------------------------------------------------------------
```

## Metrics Explained

### Total Requests Completed

**What it means:** The total number of HTTP requests that were sent during the test run for this endpoint.

**What to look for:**
- Compare this number to your expected request count based on duration, rate, and connections
- A significantly lower number might indicate the test ended early or requests were throttled

**Example:**
```
[STATS]: Total requests completed: 1000
```

### Total Number of Error Requests

**What it means:** The count of requests that failed (network errors, timeouts, HTTP error status codes, etc.).

**What to look for:**
- **Low error rate (0-1%):** Generally acceptable for most applications
- **Moderate error rate (1-5%):** May indicate occasional issues under load
- **High error rate (>5%):** Suggests the server is struggling under load or there are configuration issues

**Example:**
```
[STATS]: Total Number of Error Requests: 5
```

**Success rate calculation:**
```
Success Rate = (Total Requests - Error Requests) / Total Requests × 100%
```

### Average Response Size

**What it means:** The average size (in bytes) of successful HTTP response bodies received.

**What to look for:**
- Helps understand data transfer requirements
- Larger response sizes may impact network bandwidth and response times
- Useful for capacity planning

**Example:**
```
[STATS]: Average response size: 1024.5 bytes
```

### Min Response Time

**What it means:** The fastest response time recorded for any single request during the test.

**What to look for:**
- Represents the best-case performance
- Useful for understanding baseline performance without load
- Compare with average response time to see variance

**Example:**
```
[STATS]: Min response time: 45ms
```

### Max Response Time

**What it means:** The slowest response time recorded for any single request during the test.

**What to look for:**
- Represents worst-case performance
- Large gaps between min and max indicate inconsistent performance
- Helps identify performance bottlenecks or timeout issues

**Example:**
```
[STATS]: Max response time: 1250ms
```

### Average Response Time

**What it means:** The mean response time across all successful requests for this endpoint.

**What to look for:**
- **Primary performance indicator** – most important metric for understanding typical user experience
- Compare against your SLA or performance requirements
- Consider percentiles (if available) for more detailed analysis

**Example:**
```
[STATS]: Average response time: 234ms
```

## Interpreting Results

### Good Performance Indicators

**Low average response time** relative to your requirements  
**Low error rate** (< 1%)  
**Consistent response times** (small gap between min and max)  
**High request completion rate** matching expected load

### Warning Signs

**High error rate** (> 5%) – Server may be overloaded or misconfigured  
**Large variance** between min and max response times – Inconsistent performance  
**Very high max response time** – Possible timeouts or bottlenecks  
**Lower than expected request count** – Test may have been throttled or ended early

### Performance Analysis Tips

1. **Compare endpoints:** Different endpoints may have different performance characteristics
2. **Consider load:** Response times typically increase with higher concurrency
3. **Check error patterns:** Are errors concentrated at certain times or evenly distributed?
4. **Review configuration:** Ensure your test configuration matches real-world usage patterns

## Example Analysis

Given these results:

```
[STATS]: Total requests completed: 1000
[STATS]: Total Number of Error Requests: 2
[STATS]: Average response size: 2048 bytes
[STATS]: Min response time: 50ms
[STATS]: Max response time: 500ms
[STATS]: Average response time: 150ms
```

**Analysis:**
- **Success rate:** 99.8% (998/1000) – Excellent
- **Average response time:** 150ms – Good for most web APIs
- **Consistency:** Max (500ms) is 3.3× the average, which is reasonable
- **Response size:** 2KB average – Moderate, shouldn't cause issues
- **Overall:** This endpoint is performing well under the tested load

## Limitations

Current Load-Pulse results include aggregate statistics. For more detailed analysis, consider:

- **Percentiles** (p50, p95, p99) – Not currently reported, but would provide better insight into response time distribution
- **Request rate over time** – Current results show totals, not time-series data
- **Per-endpoint breakdowns** – Results are shown per endpoint, which is helpful for identifying problem areas

## Next Steps

After reviewing your results:

1. **If performance is good:** Consider testing with higher load or longer duration
2. **If errors are high:** Check server logs, reduce concurrency, or investigate server capacity
3. **If response times are high:** Optimize your API, check database queries, or scale your infrastructure
4. **For more detailed analysis:** Consider integrating Load-Pulse results with monitoring tools or exporting data for further analysis

For information on configuring your tests, see the [Config Reference](/config-reference) page.


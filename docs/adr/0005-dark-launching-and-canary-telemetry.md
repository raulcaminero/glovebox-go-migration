# ADR 0005: Shadow traffic validation, canary rollouts, and OpenTelemetry verification

## Status
Accepted

## Context
When replacing production Node endpoints with Go handlers (especially high-throughput query endpoints like policy filtering or expiration reports), verifying correctness, response equivalence, and latency improvements in a staging environment is insufficient. Real production traffic often exposes edge cases, legacy data quirks, and unexpected query params.

The platform requires a structured strategy for validating Go endpoint performance and behavior before full live cutover.

## Decision
Implement a 3-tier cutover verification workflow:
1. **Header-based Override & Local Verification**: Developers and QA force traffic to `go-service` or `legacy-nest` using `X-Force-Backend` request headers.
2. **Shadow Traffic / Dark Launching (Read Mirroring)**: For read-heavy endpoints (`GET /api/v1/policies/expiring`), the gateway duplicates incoming production requests, sending a background request to `go-service` without returning its response to the client. Responses and error rates are diffed asynchronously.
3. **OpenTelemetry Telemetry & Metric Verification**: Both `legacy-nest` and `go-service` emit standardized OpenTelemetry metrics (p95/p99 latency, HTTP status codes, DB query durations, heap memory usage) to Prometheus/Grafana. Cutover proceeds only after Go demonstrates a lower p99 latency and 0% error rate diff.

## Alternatives & Rollout Patterns Considered

### 1. Direct DNS Cutover without Shadowing
- **Approach**: Route 100% of live traffic to the new Go endpoint as soon as unit/integration tests pass in CI.
- **Why Rejected**: Uncovers production bugs directly in front of live insurance agents. If an edge-case SQL query panics or returns malformed JSON, end users experience immediate outage.

### 2. Manual Sampling / Log Diffing
- **Approach**: Manually inspect application logs post-deploy to verify error counts.
- **Why Rejected**: Reactive and prone to human oversight. Does not provide quantitative latency or memory distribution benchmarks.

### 3. Percentage-Based Feature Flags (Canary Rollouts)
- **Approach**: Route 1% → 10% → 50% → 100% of live user traffic based on Tenant / Agency ID (e.g. via LaunchDarkly).
- **Why Adopted as Next Step**: Canary rollouts are used immediately after Shadow Traffic verification succeeds, enabling a controlled blast radius for write operations (`POST`, `PATCH`).

## Real-World Industry Precedents
- **Cloudflare & AWS**: Use dark launching / shadow traffic to validate rewritten HTTP proxies and router modules against petabytes of live traffic before making the new binary primary.
- **Netflix**: Uses automated canary analysis (Kayenta) comparing canary vs baseline OpenTelemetry metrics before promoting new backend service versions.

## Consequences
+ **Zero Customer Impact during Testing**: Shadow traffic bugs or panics occur silently in background routines without returning errors to live clients.
+ **Empirical Verification**: Performance gains (e.g., Go handling 10x throughput with 1/5th memory footprint) are backed by production metric telemetry rather than synthetic benchmarks.
+ **Controlled Blast Radius**: Header overrides allow targeted internal agency testing before enabling global routing rules.
- **Temporary Double I/O on Shadowed Reads**: Shadowing doubles database read queries during validation windows (mitigated by read replica scaling or limited shadow sampling rates).

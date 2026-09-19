# ADR 0003: Strangler-fig migration over a big-bang cutover

## Status
Accepted

## Context
GloveBoxCRM is a production system serving live insurance agencies; it
can't go down or regress mid-migration. The JD explicitly requires
"maintain backward compatibility during rollout" and a phased roadmap
(6-month: ship one Go module to production; 12-month: 30-50% migrated).

A full rewrite-then-cutover approach risks a long-lived, high-stakes
release with no incremental validation in production.

## Decision
Put a thin reverse-proxy gateway in front of both the legacy NestJS
service and the new Go service. The gateway routes requests by path
prefix against a list of "migrated" routes; anything not yet migrated
falls through to the legacy service unchanged. As each module's Go
implementation is production-ready, its route prefix moves from "legacy"
to "migrated" — a config change or feature flag toggle, not a redeploy of
either backend.

This repo's `gateway/` implements this pattern: `/api/v1/policyholders`
and `/api/v1/policies` are routed to `go-service`, while any non-migrated
route continues hitting `legacy-nest` until its own migration lands.

## Alternatives & Migration Patterns Considered

### 1. Big-Bang Rewrite & Cutover ("The Cold Cut")
- **Approach**: Build the entire Go application in a side repository over 6–12 months, then schedule a maintenance window to flip DNS.
- **Why Rejected**: Concentrates 100% of risk into a single release. Bugs found post-cutover affect the entire platform with no instant fallback path. High risk of feature drift between legacy development and rewrite development.

### 2. Branch-by-Abstraction (Internal Application-Level Strangling)
- **Approach**: Create abstract TypeScript interfaces inside NestJS (e.g., `PolicyService`), implementing them via internal gRPC/HTTP calls to the new Go service while keeping NestJS as the outer API layer.
- **Why Rejected**: Keeps Node on the critical path for every single request, preventing Go from owning its network throughput, low-latency zero-alloc handlers, and independent infrastructure scaling.

### 3. Event-Driven / Async Queue Strangling
- **Approach**: Intercept background worker tasks (e.g., carrier policy syncs, webhook ingestion) using a message broker (RabbitMQ/NATS/Kafka) and process them in Go first, before migrating REST APIs.
- **Why Adopted as Complementary**: While excellent for background jobs, user-facing CRUD endpoints still require an API-level proxy. This will be used in Phase 2 for asynchronous carrier integrations.

### 4. Shadow Traffic / Dark Launching (Read Replaying)
- **Approach**: Proxy duplicates production read requests (`GET /api/v1/policies/expiring`) to both NestJS and Go simultaneously in the background, diffing latency and response payloads in telemetry without returning the Go response to the client.
- **Why Selected for Verification**: Allows zero-risk verification of Go performance under real production traffic loads prior to enabling live routing.

## Real-World Industry Precedents
- **Shopify & GitHub**: Used API-level Strangler Fig proxies (via NGINX/Envoy) to incrementally peel off monolithic endpoints into isolated services without breaking mobile or third-party API clients.
- **Etsy & Stripe**: Pioneered dark-launching / shadow traffic replay to validate new payment/inventory engine rewrites against millions of live HTTP requests before cutting over authority.

## Consequences
+ **Zero-Downtime Rollout**: Each module migrates independently; a failure in one module doesn't affect legacy endpoints.
+ **Instant Rollback**: Instant fallback via header flag (`X-Force-Backend: legacy`) or gateway prefix list update.
+ **Matches 6/12-Month KPIs**: Enables incremental milestone delivery (e.g. policy module in month 3, contacts in month 6) rather than delayed value delivery.
- **Infrastructure Overhead**: Adds a gateway proxy component to the network critical path (mitigated by keeping gateway logic lightweight and in Go).
- **Dual Runtime Coexistence**: Operations must manage two application runtimes and deployment pipelines during the strangling window.


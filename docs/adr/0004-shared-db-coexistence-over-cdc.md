# ADR 0004: Shared Postgres database during strangler migration vs CDC microservice decoupling

## Status
Accepted

## Context
During a strangler-fig migration from NestJS to Go, both applications must co-exist in production. A critical architectural decision is how data storage and state are managed while routes are being incrementally cut over from NestJS to Go. 

The primary challenge is preventing data divergence, distributed transaction failures, or dual-write race conditions while maintaining zero downtime.

## Decision
Adopt a **Shared Postgres Database with Logical Schema & Domain Boundaries** during Phase 1 & Phase 2 of the migration. 

Both `legacy-nest` (via TypeORM) and `go-service` (via `pgx`/`sqlc`) connect to the same underlying PostgreSQL instance. As modules migrate to Go authority:
1. **Shared Read/Write Tables during strangling**: Tables for migrated domain models (e.g., `policyholders`, `policies`) are accessed by Go for newly migrated routes while legacy NestJS routes fall back to the same tables.
2. **Explicit Foreign Key and Index Ownership**: Database migrations are decoupled from application startup and managed centrally via `goose` migration scripts.
3. **Eventual DB Isolation (Phase 3)**: When 100% of routes are migrated to Go, database access is fully restricted to `go-service`, completing the transition without requiring an expensive cross-database streaming setup during early phases.

## Alternatives & Database Migration Patterns Considered

### 1. Change Data Capture (CDC via Debezium & Kafka)
- **Approach**: Provision a separate Postgres database for `go-service`. Stream all writes from legacy Postgres to Go Postgres in real-time via WAL (Write-Ahead Log) replication using Debezium and Kafka/NATS.
- **Why Deferred**: Adds significant operational complexity (Kafka cluster, schema registry, replication lag monitoring, out-of-order event handling) for a mid-sized CRM platform before Phase 1 KPIs are even met. CDC is reserved for Phase 3 if microservice database isolation becomes a business requirement.

### 2. Dual Writing at Application Layer
- **Approach**: Modify NestJS HTTP handlers or TypeORM subscribers to write to both the legacy database and a new Go database synchronously or asynchronously via an Outbox pattern.
- **Why Rejected**: Synchronous dual-writes introduce double latency and fail when one write fails (dual-write anomaly). Asynchronous dual-writes require complex reconciliation background jobs to handle split-brain state.

### 3. Database Views and Triggers Layer
- **Approach**: Move legacy table structures to custom Postgres views/triggers while Go writes to normalized new tables.
- **Why Rejected**: Pushes application logic into database triggers, making debugging and migration rollback obscure and difficult to test in CI/CD pipelines.

## Real-World Industry Precedents
- **GitHub Monolith Migration**: Kept MySQL as the unified data layer while strangling Ruby on Rails controllers into Go services, deferring database sharding/decoupling until service boundaries had stabilized in production.
- **Figma Service Extraction**: Maintained shared Postgres schemas between monolithic TS services and high-performance Rust services during initial cutover phases to avoid data synchronization race conditions.

## Consequences
+ **Zero Replication Lag**: No data inconsistency between legacy and new endpoints during parallel operation.
+ **Instant Safety**: Rolling back an endpoint from Go to NestJS via gateway configuration requires zero data backfilling or DB sync.
+ **Reduced Infrastructure Cost & Complexity**: Avoids deploying Kafka, Debezium, and secondary database instances during initial migration phases.
- **Shared Database Connection Limits**: Both NestJS and Go pool connections to the same Postgres instance, requiring careful connection pool tuning (`pgxpool` size vs TypeORM pool size).
- **Strict Schema Discipline Required**: Schema changes must be strictly backward-compatible (e.g. non-breaking column additions) until NestJS is fully decommissioned.

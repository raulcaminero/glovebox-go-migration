# ADR 0002: Use pgx + sqlc instead of a Go ORM for the data access layer

## Status
Accepted

## Context
The Go rewrite requires a high-performance database access layer for Postgres. The team is accustomed to TypeORM-style ergonomics from NestJS (entities, decorators, `repo.find(...)`). However, query performance, predictability, and safety matter for GloveBox's search and reporting endpoints (e.g., *"policies expiring in N days"* across large policyholder datasets).

The job description explicitly calls out "database query optimization" as a core requirement.

## Decision
Use `pgx/v5` as the PostgreSQL driver and `sqlc` to generate type-safe Go code from hand-written SQL files, rather than using a runtime Go ORM (such as GORM or Ent).

Queries live as plain `.sql` files in `db/queries/`. Running `sqlc generate` produces typed structs and execution methods in `internal/repo/sqlcgen/`. Domain services interact strictly with repository interfaces (`domain.PolicyholderRepo`, `domain.PolicyRepo`), insulating business logic from database driver imports.

## Alternatives & Persistence Layer Stacks Considered

### 1. GORM (Traditional Go ORM)
- **Approach**: Use GORM with struct annotations and dynamic query builder methods (`db.Where(...).Find(&policies)`).
- **Pros**: Familiar ORM ergonomics for developers coming from TypeORM or Prisma; auto-generates schema migrations.
- **Cons**: Uses runtime reflection for row scanning, adding garbage collection pressure; obscures generated SQL, making N+1 query patterns and inefficient JOINs easy to accidentally introduce.

### 2. Standard `database/sql` + Manual Row Scanning
- **Approach**: Write queries using standard library `database/sql` and hand-write `rows.Scan(&p.ID, &p.Name, ...)` calls.
- **Pros**: Zero third-party code generator dependencies; full raw SQL control.
- **Cons**: Highly repetitive boilerplate code; every schema update requires manual maintenance of column indexes and scanning lines.

### 3. Ent (Facebook's Graph ORM for Go)
- **Approach**: Define schemas using Ent's Go DSL and generate graph-based query builders.
- **Pros**: Strongly-typed graph traversals and schema validation.
- **Cons**: Steeper learning curve; schema DSL code generation creates high abstraction overhead for standard CRM relational queries.

## Real-World Industry Precedents
- **GitHub & Stripe**: Prefer raw SQL query generation tools (`sqlc`) or lightweight mappers over full ORMs for core transactional paths to maintain 100% reviewable SQL in pull requests and prevent N+1 performance regressions.

## Consequences
+ **100% Reviewable SQL**: Every query executed in production is plain SQL visible in PR diffs — no ORM black-box generation.
+ **Compile-Time Safety**: Parameter types and result column shapes are verified by `sqlc` during compilation; column type changes break the build before hitting production.
+ **Zero Reflection Overhead**: `sqlc` generates direct struct assignment code, maximizing execution speed and minimizing GC overhead.
- **Codegen Step**: Schema changes require running `sqlc generate` and committing generated Go artifacts (documented in Makefile `generate` target).
- **Explicit Migration Tooling**: Schema migrations are managed explicitly via `goose` scripts rather than inferred automatically from ORM structs.

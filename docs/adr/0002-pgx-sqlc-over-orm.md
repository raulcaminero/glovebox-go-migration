# ADR 0002: Use pgx + sqlc instead of a Go ORM for the data access layer

## Status
Accepted

## Context
The Go rewrite needs a data access layer for Postgres. The team is used to
TypeORM-style ergonomics from the NestJS side (entities, decorators,
`repo.find(...)`). Query performance and predictability matter for
GloveBox's search/filter and reporting endpoints (e.g., "policies expiring
in N days" across large policyholder sets), and the JD calls out "database
query optimization" explicitly as a required skill.

## Decision
Use `pgx/v5` as the Postgres driver and `sqlc` to generate type-safe Go
code from hand-written SQL, instead of a Go ORM (e.g., GORM, ent).

Queries live as plain `.sql` files in `db/queries/`; `sqlc generate`
produces typed structs and methods in `internal/repo/sqlcgen/`, which the
Postgres adapter (`internal/repo/*_pg.go`) wraps behind the domain's repo
interfaces. The service layer never imports pgx or sqlc types directly.

## Alternatives Considered
- **GORM**: familiar ORM ergonomics for the team, but hides the generated
  SQL, makes N+1 query patterns easy to introduce accidentally, and
  reflection-based scanning has a real performance cost at scale.
- **`database/sql` + hand-written scanning**: full control, no codegen
  dependency, but too much repetitive boilerplate as the schema grows —
  every new query means hand-writing `Scan()` calls and struct mapping.
- **ent (Facebook's Go ORM)**: strong type safety and graph-style
  queries, but a steeper learning curve and more magic (code generation
  from a schema DSL rather than from SQL you write yourself) than the
  team needs for a CRM's relatively standard query patterns.

## Consequences
+ Every query's exact SQL is visible and reviewable in a PR — no ORM
  black box generating unexpected joins or N+1 patterns.
+ Compile-time safety on query parameters and result shapes without
  runtime reflection.
+ Because the domain layer only depends on repo *interfaces*
  (`domain.PolicyholderRepo`, `domain.PolicyRepo`), the service layer and
  its tests never import pgx/sqlc — swapping the persistence layer later
  wouldn't touch business logic.
- Schema changes require running `sqlc generate` and committing the
  output; this is an extra step compared to an ORM inferring everything
  from struct tags, and needs to be documented for new engineers
  (see Makefile `generate` target).
- Team gives up ORM conveniences like automatic migrations-from-entities;
  migrations are managed explicitly via `goose`.

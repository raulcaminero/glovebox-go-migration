# GloveBox Go Migration — a working slice

This repo simulates the migration described in the **Full Stack Software
Engineer (Go & AI Tooling Lead)** role: taking one module of a Node/NestJS
CRM and rewriting it in Go, safely, incrementally, and with the AI-assisted
development workflow the role asks for.

It's not a toy CRUD app. It's a small but complete example of the actual
job: an architectural decision, a migration mechanism that doesn't require
downtime, the documentation that makes the decision durable, and the AI
tooling guardrails that keep a team's output consistent.

## What's here

```
legacy-nest/     "Before" — a minimal NestJS service (Policyholder CRUD,
                 TypeORM, SQLite) standing in for GloveBoxCRM's current
                 Node backend.

go-service/      "After" — the same domain rewritten in Go: clean
                 layered architecture, Postgres via pgx+sqlc, JWT auth,
                 Policyholder + Policy + Notes entities, structured slog logging,
                 OpenTelemetry trace spans, and table-driven unit tests.

gateway/         A strangler-fig reverse proxy that routes requests to
                 legacy-nest or go-service by path prefix — the mechanism
                 that lets the migration ship module by module instead
                 of as one big cutover.

docs/adr/        Three Architecture Decision Records for the real
                 decisions this repo makes.

docs/MIGRATION.md   A phased rollout roadmap mirroring the role's own
                    6-/12-month success metrics.

docs/AI_WORKFLOW.md How AI coding assistants were used to build this,
                    with a concrete before/after example of a caught
                    mistake.

CLAUDE.md        Coding standards fed to AI assistants working in this
                 repo — the guardrails referenced in AI_WORKFLOW.md.

tools/ai-pr-review/  A small Go program that sends a PR diff to the
                     Claude API for a first-pass check against CLAUDE.md.
```

## Architecture

```
                 ┌─────────────┐
   client  ───▶  │   gateway   │
                 └──────┬──────┘
                        │  routes by path prefix
             ┌──────────┴───────────┐
             ▼                      ▼
     ┌───────────────┐      ┌───────────────┐
     │  legacy-nest   │      │   go-service   │
     │  (NestJS/Nest) │      │      (Go)      │
     │  not-yet-       │      │  policyholders│
     │  migrated       │      │  policies     │
     │  routes         │      │  & notes      │
     └───────────────┘      └───────┬────────┘
                                     ▼
                              ┌──────────────┐
                              │  Postgres    │
                              │ (pgx + sqlc) │
                              └──────────────┘
```

`go-service` internally follows a ports-and-adapters layout:

```
transport/http  →  service  →  domain (interfaces, no framework deps)
                                   ▲
                                   │ implements
                              repo (pgx/sqlc Postgres adapter)
```

`domain` never imports pgx, sqlc, or chi. The service layer only depends
on domain interfaces. This is what makes the service layer unit-testable
with zero database (see `internal/service/*_test.go`, which use in-memory
fake repos) and what would let a future storage change happen without
touching business logic.

## Why these decisions — the ADRs

- [`0001-go-over-node.md`](docs/adr/0001-go-over-node.md) — why rewrite in
  Go at all, and why incrementally rather than all at once.
- [`0002-pgx-sqlc-over-orm.md`](docs/adr/0002-pgx-sqlc-over-orm.md) — why
  hand-written SQL + generated types instead of an ORM.
- [`0003-strangler-fig-over-big-bang.md`](docs/adr/0003-strangler-fig-over-big-bang.md) —
  why the gateway pattern instead of a cutover release.

## The migration roadmap

[`docs/MIGRATION.md`](docs/MIGRATION.md) sketches how this pattern would
scale to the real GloveBoxCRM codebase, phased against the same 6-/12-month
milestones in the role description (ship one module → onboard junior
engineers → 30-50% migrated).

## AI-assisted workflow

[`docs/AI_WORKFLOW.md`](docs/AI_WORKFLOW.md) documents how Claude was used
to build parts of this repo, including a concrete example of a first draft
that violated the architecture boundary and how `CLAUDE.md`'s stated rules
caught it in review. `tools/ai-pr-review/` is a small working tool that
automates a first pass of that same check against a PR diff.

## Running & Testing

**Prerequisites:** Go 1.22+, Node 20+, Docker (for Postgres).

```bash
# 1. Start Postgres
docker compose up -d postgres

# 2. Go service
cd go-service
go mod tidy
export DATABASE_URL="postgres://glovebox:glovebox@localhost:5435/glovebox"
export JWT_SECRET="dev-secret-change-me"
make migrate-up   # requires goose: https://github.com/pressly/goose
make run          # listens on :8080

# 3. Legacy Nest service (separate terminal)
cd legacy-nest
npm install
npm run start     # listens on :3001

# 4. Gateway (separate terminal)
cd gateway
go run .          # listens on :8000, routes between the two above
```

Then, once you have a JWT (any HS256 token signed with `JWT_SECRET`, with
a `sub` claim):

```bash
curl -H "Authorization: Bearer $TOKEN" \
     -X POST http://localhost:8000/api/v1/policyholders \
     -d '{"full_name":"Jane Doe","email":"jane@example.com"}'

# Which of your clients have coverage expiring in the next 30 days?
curl -H "Authorization: Bearer $TOKEN" \
     "http://localhost:8000/api/v1/policies/expiring?within_days=30"

# Add a note for a policyholder
curl -H "Authorization: Bearer $TOKEN" \
     -X POST http://localhost:8000/api/v1/policyholders/$ID/notes \
     -d '{"author":"Agent Smith","body":"Client requested quote update"}'
```

Watch the `X-Routed-To` response header — it shows whether the gateway
sent that request to `go-service` (migrated) or `legacy-nest` (not yet).

**Unit Tests:**
- Go service: `cd go-service && go test -v ./...`
- Legacy Nest service: `cd legacy-nest && npm test`

## What I'd do differently at GloveBox's actual scale

- **Multi-tenancy from day one.** This demo has no `agency_id` on any
  table. A real CRM serving many independent agencies needs tenant
  isolation baked into the schema and every query from the start, not
  retrofitted.
- **Transactional outbox for carrier sync.** Real policy data comes from
  external carrier integrations; I'd add an outbox pattern so a partial
  failure mid-sync doesn't leave policyholder and policy data
  inconsistent.
- **Contract tests between gateway and both backends**, not just unit
  tests within each — the riskiest part of a strangler-fig migration is a
  silent contract drift between the "before" and "after" implementations.
- **Structured audit logging** on every write, given insurance's
  regulatory/compliance requirements (the JD calls this out directly).

## A note on scope

My production background is TypeScript/Node (10+ years, including eBay and
Thryv). I'm not presenting this as years of production Go — I'm presenting
it as the evidence that matters more for this role: that I know how to
design a migration correctly, enforce architecture boundaries, set up AI
tooling guardrails, write testable services, and document decisions durably.

The patterns here — ports-and-adapters layering, interface-driven repos,
table-driven unit tests with zero database, a phased rollout roadmap that
maps to the role's stated KPIs — are the same patterns I've applied in
production TypeScript systems at scale. Go is the implementation language;
the architecture instincts are language-agnostic.

For React/Next.js frontend work, see:
- **[CMHub](https://github.com/raulcaminero/cmhub)** — multi-tenant fiscal ERP
  (Next.js App Router, NestJS, Prisma, pgvector AI copilot, Neon.tech Postgres)
- **[Portfolio](https://github.com/raulcaminero/Portafolio)** — production Next.js
  site deployed on Vercel


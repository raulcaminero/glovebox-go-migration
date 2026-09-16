# ADR 0001: Rewrite the CRM backend in Go instead of staying on Node/NestJS

## Status
Accepted

## Context
GloveBoxCRM's backend runs on NestJS today. As the platform grows toward
higher-throughput needs (active policy monitoring, real-time carrier sync,
analytics on tens of thousands of service transactions), the team needs to
decide whether to keep optimizing the Node stack or invest in a unified
rewrite.

Node/NestJS is productive for CRUD and has a large ecosystem, but:
- Concurrency model (event loop + async/await) makes CPU-bound or
  highly concurrent I/O workloads (e.g., syncing hundreds of carrier
  connections in parallel) harder to reason about and tune than a
  goroutine-based model.
- Runtime type safety depends entirely on discipline (TypeScript is
  erased at compile time); a class of bugs that Go's static typing and
  explicit error handling catch earlier slip through in fast-moving Node
  codebases.
- Deployment artifacts are heavier (Node runtime + node_modules) versus a
  single static Go binary, which simplifies containerization and cold
  starts.

## Decision
Migrate the backend to Go incrementally, using a strangler-fig pattern
(see ADR 0003) rather than a big-bang rewrite. New modules are built in Go
from day one; existing modules migrate module-by-module as they're
touched or as their throughput/reliability needs justify it.

## Alternatives Considered
- **Stay on NestJS, invest in performance tuning**: lower short-term
  risk, but doesn't address the structural concurrency and type-safety
  gaps, and doesn't unify the codebase — the JD explicitly frames
  fragmentation as a top problem.
- **Rewrite in a different compiled language (Rust, Java)**: Rust's
  learning curve and slower iteration speed are a poor fit for a team
  that needs to onboard engineers quickly; Java carries more ceremony and
  a heavier runtime than the "lean single binary" goal here.

## Consequences
+ Simpler deployment (single static binary, small containers).
+ Better default handling of concurrent I/O (goroutines, `context.Context`
  cancellation) for carrier integrations and background jobs.
+ Compile-time safety and explicit error handling reduce a class of
  runtime bugs common in loosely-typed Node code.
- Team must build Go proficiency; short-term velocity dip during ramp-up.
- Two languages in production during the migration window increases
  operational surface area (mitigated by ADR 0003's routing approach).
- Fewer off-the-shelf libraries than the Node ecosystem for some
  integrations; more code may need to be written in-house.

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
to "migrated" — a config change, not a redeploy of either backend.

This repo's `gateway/` implements this pattern: `/api/v1/policyholders`
and `/api/v1/policies` are already routed to `go-service`, while any
other route (in a real system: billing, claims, auth) would continue
hitting `legacy-nest` until its own migration lands.

## Alternatives Considered
- **Big-bang rewrite and cutover**: lower long-term complexity (only one
  backend running at a time), but concentrates all migration risk into a
  single release; a bug found post-cutover affects the entire platform
  with no fallback.
- **Branch-by-abstraction inside the existing Node codebase** (build the
  Go pieces as separate services called *from* Node, no gateway): avoids
  a new routing layer, but leaves Node as the system of record
  indefinitely and doesn't give Go services a direct path to owning their
  own endpoints and scaling independently.
- **Dual writes / shadow traffic** (send every request to both backends,
  compare responses, but only serve the legacy response): higher
  confidence before cutover, but meaningfully more implementation
  complexity than this migration's timeline justifies for a first phase.

## Consequences
+ Each module migrates independently; a failed or delayed migration for
  one module doesn't block others.
+ Rollback is a one-line config change (move a prefix back to legacy)
  instead of a redeploy or revert.
+ Matches the JD's phased KPI structure directly — each completed module
  is a visible, demoable milestone rather than an all-or-nothing bet.
- The gateway itself is a new piece of infrastructure that needs to be
  reliable and low-latency, since it now sits on the critical path for
  every request.
- Running two backends in parallel during the migration window means
  double the operational surface (two sets of logs, two deploy
  pipelines, two languages) until the migration completes.

# Go Migration Roadmap

This document sketches how the pattern demonstrated in this repo would
scale to GloveBoxCRM's actual codebase. It mirrors the 6-/12-month
milestones from the role's success metrics so progress against them is
easy to track.

## Guiding principles

1. **No big-bang cutover.** Every module moves behind the strangler-fig
   gateway (ADR 0003) independently, validated in production before the
   next module starts.
2. **Backward-compatible contracts.** A migrated module's HTTP API shape
   stays identical to the legacy one unless a breaking change is
   explicitly planned and versioned — clients (mobile app, agency web
   portal) should not need to change during the migration.
3. **Tests before traffic.** A module doesn't get its route flipped to
   the Go service until it has parity test coverage against the legacy
   behavior it replaces.

## Phase 0 — Foundations (this repo)
- Establish the layered architecture (`domain` / `service` / `repo` /
  `transport`) and conventions other engineers will follow.
- Stand up the gateway and prove the routing/rollback mechanism works.
- Document architecture decisions as ADRs so the "why" survives staffing
  changes.
- Set up AI-assisted workflow guardrails (see `AI_WORKFLOW.md`) so the
  velocity gains from Claude/Copilot don't come at the cost of consistency.

## Phase 1 — First production module (0–6 months)
Target: **contacts/policyholders and core policy CRUD** — high traffic,
well-understood domain, low risk of hidden edge cases. This validates the
whole pipeline (schema, auth, testing, deploy, gateway routing) on a
module simple enough to de-risk the pattern itself.
- Ship policyholder + policy endpoints in Go, behind the gateway.
- Publish this migration roadmap and the architecture/ADR conventions for
  the team (deliverable, not just a plan).
- Begin onboarding the first junior engineer(s) into the Go codebase,
  pairing on a small, well-scoped ticket in the new service.

**Success looks like:** one module running in production in Go, serving
real traffic, with no client-visible regressions, and documented enough
that a new hire can extend it without a walkthrough from the migration
lead.

## Phase 2 — Expand coverage (6–12 months)
Target: **30–50% of the platform migrated**, prioritized by:
1. Modules with the clearest business value from Go's concurrency model
   (carrier sync jobs, active policy monitoring — currently polling-heavy
   in Node).
2. Modules blocking new feature work, so the team isn't building new
   features twice (once in old patterns, once in new).
3. Modules with the best existing test coverage, to keep migration risk low.
- 2-3 junior engineers actively contributing to the Go codebase, not just
  the migration lead.
- Establish measurable baselines: build times, p95 API latency, deploy
  frequency — tracked before/after each module's migration so the
  platform-uptime and developer-velocity KPIs are backed by numbers, not
  impressions.

## Phase 3 — Beyond 12 months (indicative, not committed)
- Revisit which legacy modules remain and whether full migration is worth
  it versus maintaining a stable minority in Node long-term — not every
  system needs to move, and forcing the last 10% can cost more than it's
  worth.
- Formalize the AI-assisted workflow guardrails into a written engineering
  standard as team size grows.

## Risk register (abbreviated)
| Risk | Mitigation |
|---|---|
| Migrated module regresses a client-visible behavior | Gateway rollback is a config change, not a redeploy |
| Junior engineers ramp slowly on Go | Pair on Phase 1 module; ADRs + `CLAUDE.md` reduce tribal knowledge dependency |
| Two backends in parallel increases on-call surface | Keep the legacy service in maintenance-only mode — no new features land there once a module's migration starts |

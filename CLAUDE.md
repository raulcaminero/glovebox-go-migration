# Instructions for AI coding assistants working in this repo

This file exists so Claude/Copilot output on this codebase follows the
same conventions a human reviewer would enforce in a PR — the "quality
guardrails" the AI-tooling side of this role is about.

## Architecture rules (do not violate)
- `internal/domain` has **zero imports** from `pgx`, `sqlc`, `chi`, or any
  framework. If a generated diff imports one of these into `domain/`,
  reject it — that's an architecture boundary violation, not a style nit.
- The service layer (`internal/service`) depends only on domain
  interfaces (`domain.PolicyholderRepo`, etc.), never on the concrete
  Postgres adapter or sqlc types directly.
- HTTP handlers (`internal/transport/http`) stay thin: decode request,
  call one service method, encode response. Business logic/validation
  belongs in `internal/service`, not in a handler.

## Go conventions
- Every exported function that can fail returns `error` as the last
  value; wrap errors with `fmt.Errorf("doing X: %w", err)` for
  traceability, never swallow them silently.
- All I/O-bound functions take `context.Context` as the first parameter
  and respect cancellation.
- New query logic goes in `db/queries/*.sql` + `sqlc generate` — do not
  hand-write ad-hoc SQL strings inside Go files outside `internal/repo`.
- New business rules get a table-driven test in the corresponding
  `*_test.go` using the existing fake-repo pattern (see
  `internal/service/policyholder_test.go`) — no new test infra without a
  reason.

## What NOT to do
- Don't introduce a new dependency (module) without flagging it in the PR
  description — this includes swapping in an ORM, a different JWT
  library, etc.
- Don't "helpfully" refactor unrelated files in the same diff as a
  feature change; keep AI-assisted diffs scoped to what was asked.
- Don't generate code that silently changes an existing HTTP response
  shape — this repo cares about backward compatibility during the
  migration (see docs/adr/0003-strangler-fig-over-big-bang.md).

## How this file gets used
See `docs/AI_WORKFLOW.md` for a concrete before/after example of an
AI-assisted port from the legacy NestJS service to Go, including what the
first draft got wrong and how this file's rules caught it in review.

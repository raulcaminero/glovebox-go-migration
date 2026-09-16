# AI-Assisted Workflow

This documents how AI coding assistants were used to build this repo, and
the guardrails that catch what they get wrong — the actual practice, not
just a claim that "AI tools were used."

> **Note for reviewers:** the example below is illustrative of the
> workflow and guardrails; reproduce it against your own Claude Code /
> Copilot session and drop in the real prompt + diff + correction before
> sharing this repo, so it documents an actual session rather than a
> representative one.

## The workflow

1. **Scope the change** in a sentence or two before opening the AI tool —
   "port the Nest `ContactsService.update` method to the Go
   `PolicyholderService.Update`, preserving the same validation" — not
   "help me write the Go version."
2. **Generate against `CLAUDE.md`.** The assistant reads the repo's
   `CLAUDE.md`, which encodes the architecture boundaries and Go
   conventions (thin handlers, domain has no framework imports, errors
   wrapped with `%w`, table-driven tests).
3. **Review the diff like any other PR.** Check it against the three
   things `CLAUDE.md` calls out explicitly: architecture boundary
   violations, missing tests, and scope creep (unrelated files touched).
4. **Correct and note the pattern.** If the same category of mistake
   shows up more than once, it becomes a new line in `CLAUDE.md` — the
   guardrails file should get more specific over time, not stay static.

## Example: porting policyholder update logic

**Prompt given to the assistant:**
> "Port `ContactsService.update()` from `legacy-nest/src/contacts/contacts.service.ts`
> to `PolicyholderService.Update()` in Go. Keep the same validation
> behavior: reject if the resulting name would be empty."

**What the first draft got wrong:**
The generated `Update` method initially called the sqlc-generated query
directly from the service layer (`sqlcgen.Queries` imported into
`internal/service`), skipping the `domain.PolicyholderRepo` interface
entirely. It also didn't trim whitespace from `full_name` before
validating it — the legacy Nest code didn't either, but the domain rules
already established in `Create()` do, so the port introduced an
inconsistency between the two methods.

**Correction applied:**
- Routed the call through `domain.PolicyholderRepo.Update`, matching the
  existing `Create`/`Get`/`List` pattern — the exact rule `CLAUDE.md`
  states under "Architecture rules."
- Applied the same `strings.TrimSpace` normalization used in `Create`,
  so validation behavior is consistent across the service rather than
  silently diverging per method.
- Added `TestPolicyholderService_Create` -style coverage for the trimmed
  input case (see `internal/service/policyholder_test.go`).

**Why this matters for the role:** this is what "establish best
practices for AI-assisted code generation... to ensure output meets
strict... reliability standards" looks like in practice — not avoiding AI
tools, and not accepting their output uncritically, but having a written
standard the output gets checked against every time.

## The PR-review tool

`tools/ai-pr-review/` is a small script that sends a PR's diff to the
Claude API along with this repo's `CLAUDE.md` and asks it to flag
violations of the stated conventions before a human reviewer looks at it.
It's a first-pass filter, not a replacement for review — see the script's
own comments for its limitations.

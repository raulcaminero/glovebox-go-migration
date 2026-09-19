---
name: ai-pr-review
description: >
  Reviews the current git diff (or a specified branch diff) against this
  repo's CLAUDE.md architecture conventions. Flags violations, cites the
  specific rule broken, and quotes the offending lines. Activate when the
  user asks to "review my diff", "check my changes", "run ai pr review",
  or "does my code follow the architecture rules".
---

# AI PR Review Skill

You are performing a first-pass code review of a git diff against the
architecture conventions defined in `CLAUDE.md` at the root of this repo.

## Steps

1. **Get the current diff** by running:
   ```
   git diff master
   ```
   If that produces no output (already committed), try:
   ```
   git diff master...HEAD
   ```
   If the user specified a particular file or branch, use that instead.

2. **Read `CLAUDE.md`** from the repo root to load the architecture rules.

3. **Review the diff** strictly against the rules in `CLAUDE.md`. Do NOT
   invent new style opinions. Only flag what `CLAUDE.md` explicitly calls out:
   - Architecture boundary violations (e.g. `domain/` importing pgx/sqlc/chi)
   - Service layer depending on concrete adapters instead of domain interfaces
   - HTTP handlers containing business logic instead of thin decode → call → encode
   - Errors swallowed silently instead of wrapped with `fmt.Errorf("...: %w", err)`
   - New business rules without a table-driven test in `*_test.go`
   - New dependencies added without flagging them
   - Unrelated files touched in the same diff (scope creep)

4. **Report your findings** in this format:

   ### ✅ Passing
   List what the diff gets right against the stated conventions.

   ### ⚠️ Violations
   For each violation:
   - Which rule from `CLAUDE.md` is broken (quote it)
   - The offending file + line from the diff
   - A concrete fix

   ### 💡 Suggestions (optional)
   Non-blocking observations — things not in `CLAUDE.md` but worth noting
   as a human reviewer would.

## Important

- If the diff is clean, say so clearly and briefly. Don't pad the output.
- This is a first-pass filter, not a replacement for human review.
- Never approve or reject the diff outright — surface findings, let the
  developer decide.

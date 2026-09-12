---
name: generate-tests
description: Generate or update behavior-focused code and API test cases after each code-change batch, tracing requirements and diffs to executable tests and the active task's verification record.
---

# Generate Tests After Code Changes

Read AGENTS.md, `docs/status.md`, `docs/testing.md`, and the active task's plan,
PRD, design and initial test cases. Execute this workflow after every code-change
batch and before verification/delivery. Scope is the changed behavior and its
affected callers; do not expand into an unrelated full-project test rewrite.

1. Inspect the actual diff (staged and unstaged as appropriate), changed contracts,
   affected call paths, and existing tests. Use CodeGraph first when indexed.
   Derive expected behavior from PRD/AC IDs and contracts independently of the
   new implementation; a changed implementation is not its own test oracle.
2. Update `docs/tasks/<task-id>/04-test-cases.md`: assign stable Case IDs, reference
   requirements or confirmed behavior, and state preconditions, input, action,
   expected result, cleanup and automation location. Include relevant success,
   boundary, invalid input, missing resource, dependency failure and auth paths.
   Choose cases that could catch a real regression; do not enumerate irrelevant
   combinations or assert only implementation details and exact prose.
3. Implement the missing regression tests. Reuse existing meaningful tests when
   they already cover the behavior and record that mapping. Use the project's
   native unit-test location; put HTTP API tests in root `tests/api/` using Python
   and the existing client/runner. Test API status, response contract and relevant
   state changes, not only HTTP 200. Keep test fixtures isolated and repeatable.
4. For UI changes update the task's `ui-cases.json` with real page semantics,
   step-level assertions and cleanup. If no UI is affected, record non-applicability
   and its reason. Browser cases require an actual browser run to count as passed.
5. Run the affected tests. For a bug fix demonstrate the case detects the old bug
   when practical; never weaken assertions to fit a failure. Separate fixture,
   environment and assertion failures. No test service or empty discovery is
   not a successful API run. Do not put real credentials or tokens in artifacts.
6. Update `05-verification.md` with cases added/reused, command, exit code, actual
   result, coverage limits and unresolved issues. Update docs/status.md with the
   next actionable step. Mark done only after the required checks really pass.

This skill is an agent workflow, not a shell test generator or a file-save hook.
Reading it without deriving cases and verifying the result is not execution.
If a batch only changes documentation or otherwise needs no new executable test,
record why and perform the relevant link/configuration checks instead of adding
tests that merely mirror text.

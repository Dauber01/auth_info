# Project knowledge protocol

The repository is the evidence store. The knowledge base is a navigable map of
verified facts and decisions, not a second copy of the source code.

## Ownership and loading

- Framework: portable workflow, discovery rules, adapters, and reusable skills.
- Project: docs contains architecture, business rules, commands, decisions,
  pitfalls, task records and test plans. Project-specific skills reference them.
  Framework installation never replaces existing project documents or tests.
- Harness: generated configuration; independent permissions and hook semantics.
- Machine: personal preferences, credentials, and local runtime state; not shared.

AGENTS.md loads short entry links. docs/README.md routes tasks to topic documents;
docs/status.md records every task and handoff memory. Skills describe repeatable
work and reference knowledge; agents select
a role and read the same project knowledge.

## Minimal topic record

Each topic should state:

- Scope: the component or workflow and when this document is useful.
- Status: verified facts, intended rules, or open questions, labelled separately.
- Evidence: repository-relative source links, relevant symbols, or test commands.
- Freshness: verification date or revision, and changes that require rechecking.
- Content: boundaries, invariants, relevant failure paths, and practical commands.

New projects start with docs/README.md, status.md, testing.md and an explicit
unverified onboarding record. Add topic knowledge as code is inspected. Each
logical task must have its prescribed plan, PRD, prototype, design, case set and
verification records before implementation; mark non-applicable UI explicitly.

## Maintenance

Update related docs in the same change as interfaces, commands, business
rules, configuration, or architecture. For bug fixes, retain the confirmed cause
and regression check when reusable; do not turn a temporary failure into a rule
to skip tests. Keep past observations dated and re-run before asserting they hold.

Use links to existing contracts and instructions. Keep one authoritative location
per fact; route readers there from skills and indexes. Record why a decision was
made and its constraints. Avoid raw logs, credentials, production data, entire
transcripts, and unsupported claims of coverage or operational readiness. After
each code-change batch execute generate-tests, update the task case mapping and
actual verification record, then refresh status memory.

Mechanical checks validate files and discovery, not factual accuracy. Source
inspection and suitable tests are required to establish that a claim is true.

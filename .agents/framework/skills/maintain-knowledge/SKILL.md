---
name: maintain-knowledge
description: Initialize or update a project's agent knowledge base from repository evidence when onboarding, changing architecture or workflows, or correcting stale project documentation.
---

# Maintain Project Knowledge

Read AGENTS.md and `.agents/project/index.md` from the target project root.
Apply the [knowledge protocol](references/protocol.md). Keep framework rules
separate from project facts. Never copy facts from the framework's original host
project into a different project.

1. Identify the requested scope and existing project knowledge. For a new project,
   inventory repository instructions, manifests, build files, CI, entrypoints,
   contracts, and tests. Use CodeGraph first if the project has an index.
2. Follow relevant entrypoints and boundaries through current source. Record only
   claims supported by file paths, symbols, tests, or explicit user decisions.
3. Write concise project documents covering verified commands, architecture,
   domain rules, and operational constraints. Use the protocol's minimal record
   structure. Existing normative rules remain rules until implementation verifies
   them; examples are not automatically universal conventions.
4. Update `.agents/project/index.md` with when to read each document. Keep a short
   `.agents/project/context.md` for automatic loading. Detailed facts belong in
   topic documents, so the main entry stays focused.
5. Validate sources and links. Run relevant local checks when authorized and
   practical; record the command, result, revision or date, and environment limits.
   Mark anything not verified as unknown. Do not invoke deployment, database
   mutation, or external service calls merely to complete documentation.
6. Synchronize the native entries and run the framework check. Summarize what
   knowledge changed and which gaps remain; never claim the whole project was
   verified after inspecting only a sample.

When a new business rule or recurring workflow deserves a skill, keep its body
in the project's canonical `.agents/skills/<name>/` directory, link to knowledge
instead of duplicating it, and follow the portable frontmatter contract in
AGENTS.md. Generalize into the framework only after removing project assumptions.

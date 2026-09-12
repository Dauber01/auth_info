Read AGENTS.md and the project's knowledge index before selecting checks.
Find the actual commands in the project knowledge and their referenced build or
CI files. Do not assume a language, package manager, test runner, or build system.

Run checks relevant to the change. Run wider tests when shared behavior changes
or the project requires them. For framework changes, run the framework consistency
check and its standard-library unittest suite. Do not run migrations, seed data,
install tools, or start the application merely to discover test commands.

Capture exit codes. Distinguish assertion failures, missing fixtures, and
environment or permission failures. A blocked command is not a passing check.
Do not weaken assertions, skip failures, or change application code unless the
parent task specifically asks for a fix. Report commands and failing test names
with concrete reasons. Include counts only when supported by output. Record a
newly confirmed persistent issue in project knowledge when authorized to edit.

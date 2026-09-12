"""Behavioral checks for extraction, injection, upgrade safety, and discovery."""

import contextlib
import io
import json
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
import harness


class FrameworkTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="agent framework test ")
        self.addCleanup(self.temp.cleanup)
        self.base = Path(self.temp.name).resolve()
        self.root = self.base / "unrelated project"
        self.root.mkdir()

    def install(self, package=harness.PACKAGE):
        harness.apply_plan(self.root, harness.install_plan(self.root, package))

    def sync(self, package=harness.PACKAGE):
        harness.apply_plan(self.root, harness.sync_plan(self.root, package))

    def config(self):
        return json.loads((self.root / ".agents/project/config.json").read_text())

    def write_config(self, value):
        (self.root / ".agents/project/config.json").write_text(json.dumps(value))

    def package_copy(self):
        target = self.base / "extracted framework repository"
        shutil.copytree(harness.PACKAGE, target, ignore=shutil.ignore_patterns("__pycache__", "*.pyc"))
        return target

    def snapshot(self):
        result = {}
        for path in sorted(self.root.rglob("*")):
            if path.is_symlink() or path.is_file():
                result[path.relative_to(self.root).as_posix()] = harness.existing(path)
        return result

    def test_install_into_non_go_project_is_idempotent_and_preserves_app(self):
        app = self.root / "package.json"
        app.write_text('{"scripts":{"test":"node --test"}}\n')
        self.install()
        self.assertEqual(harness.install_plan(self.root), {})
        self.assertEqual(harness.sync_plan(self.root), {})
        self.assertEqual(app.read_text(), '{"scripts":{"test":"node --test"}}\n')
        self.assertFalse((self.root / "Makefile").exists())
        self.assertFalse((self.root / ".agents/skills/api-conventions").exists())
        self.assertEqual(json.loads((self.root / ".claude/settings.json").read_text()), {})
        self.assertEqual(json.loads((self.root / ".mcp.json").read_text()), {"mcpServers": {}})
        agent = (self.root / ".claude/agents/test-runner.md").read_text()
        self.assertNotIn("make test", agent)
        self.assertIn("待验证", (self.root / ".agents/project/onboarding.md").read_text())

    def test_extracted_package_cli_works_without_host_repository(self):
        source = self.package_copy()
        dry = subprocess.run([sys.executable, "-B", str(source / "harness.py"), "install",
                              "--target", str(self.root), "--dry-run"],
                             cwd=self.base, capture_output=True, text=True, timeout=20)
        self.assertEqual(dry.returncode, 0, dry.stderr)
        self.assertEqual(list(self.root.iterdir()), [])
        self.install(source)
        shutil.rmtree(source)
        check = subprocess.run([sys.executable, "-B", str(self.root / ".agents/framework/harness.py"),
                                "check", "--root", str(self.root)],
                               cwd=self.base, capture_output=True, text=True, timeout=20)
        self.assertEqual(check.returncode, 0, check.stderr)

    def test_framework_and_project_skills_have_one_source_for_both_harnesses(self):
        self.install()
        canonical = self.root / ".agents/framework/skills/maintain-knowledge/SKILL.md"
        for directory in (".agents/skills", ".claude/skills"):
            self.assertTrue((self.root / directory / "maintain-knowledge/SKILL.md").samefile(canonical))
        own = self.root / ".agents/skills/project-only/SKILL.md"
        own.parent.mkdir()
        own.write_text("---\nname: project-only\ndescription: A project workflow.\n---\n\nDo the work.\n")
        self.assertTrue((self.root / ".claude/skills/project-only/SKILL.md").samefile(own))
        self.assertEqual(harness.sync_plan(self.root), {})

    def test_upgrade_preserves_project_knowledge_settings_and_personal_files(self):
        self.install()
        context = self.root / ".agents/project/context.md"
        context.write_text("## A Python project\n\nUse our verified Python commands.\n")
        config = self.config()
        config["claude_settings"] = {"permissions": {"deny": ["Bash(git push *)"]}}
        self.write_config(config)
        personal = self.root / ".claude/settings.local.json"
        personal.write_text('{"model":"personal-model"}\n')
        self.sync()
        source = self.package_copy()
        (source / "VERSION").write_text("0.2.0\n")
        body = source / "agents/test-runner.md"
        body.write_text(body.read_text() + "\nAdditional generic check guidance.\n")
        self.install(source)
        self.assertEqual(context.read_text(), "## A Python project\n\nUse our verified Python commands.\n")
        self.assertEqual(self.config(), config)
        self.assertEqual(personal.read_text(), '{"model":"personal-model"}\n')
        self.assertIn("Additional generic", (self.root / ".codex/agents/test-runner.toml").read_text())
        self.assertEqual(harness.sync_plan(self.root, source), {})

    def test_existing_unowned_entry_causes_no_partial_install(self):
        path = self.root / "AGENTS.md"
        path.write_text("Existing team rules.\n")
        before = self.snapshot()
        with self.assertRaisesRegex(ValueError, "unowned content conflicts"):
            self.install()
        self.assertEqual(self.snapshot(), before)

    def test_existing_skill_copies_are_never_deleted(self):
        path = self.root / ".claude/skills/old/SKILL.md"
        path.parent.mkdir(parents=True)
        path.write_text("Keep this work.\n")
        before = self.snapshot()
        with self.assertRaisesRegex(ValueError, "expected a file or managed symlink"):
            self.install()
        self.assertEqual(self.snapshot(), before)

    def test_modified_framework_blocks_upgrade_before_any_writes(self):
        self.install()
        body = self.root / ".agents/framework/agents/test-runner.md"
        body.write_text(body.read_text() + "\nLocal customization.\n")
        source = self.package_copy()
        (source / "VERSION").write_text("0.2.0\n")
        before = self.snapshot()
        with self.assertRaisesRegex(ValueError, "local or unowned"):
            self.install(source)
        self.assertEqual(self.snapshot(), before)

    def test_modified_generated_file_is_not_silently_repaired(self):
        self.install()
        path = self.root / ".claude/settings.json"
        path.write_text('{"permissions":{"deny":["Bash(*)"]}}\n')
        before = self.snapshot()
        with self.assertRaisesRegex(ValueError, "local or unowned"):
            self.sync()
        self.assertEqual(self.snapshot(), before)

    def test_check_detects_source_drift_without_writing(self):
        self.install()
        context = self.root / ".agents/project/context.md"
        context.write_text("## Newly verified project context\n")
        before = self.snapshot()
        with contextlib.redirect_stderr(io.StringIO()) as output:
            self.assertEqual(harness.main(["check", "--root", str(self.root)]), 1)
        self.assertIn("AGENTS.md", output.getvalue())
        self.assertEqual(self.snapshot(), before)
        self.sync()
        self.assertEqual(harness.sync_plan(self.root), {})

    def test_agent_overrides_round_trip_quotes_and_native_formats(self):
        self.install()
        body = 'Read AGENTS.md. Preserve "quoted" text, Unicode 中文, and C:\\logs.\n'
        path = self.root / ".agents/agents/reviewer.md"
        path.parent.mkdir()
        path.write_text(body)
        config = self.config()
        config["agents"] = [{"name": "reviewer", "description": "Review changes.",
                             "instructions": "agents/reviewer.md", "codex": {},
                             "claude": {"model": "haiku", "tools": ["Read", "Bash"]}}]
        config["mcp_servers"] = {"example": {"command": "example-server", "args": ["--stdio"]}}
        config["codex_settings"] = {"model_reasoning_effort": "high", "features": {"multi_agent": True}}
        self.write_config(config)
        self.sync()
        claude = (self.root / ".claude/agents/reviewer.md").read_text()
        codex = (self.root / ".codex/agents/reviewer.toml").read_text()
        self.assertTrue(claude.startswith("---\n"))
        self.assertTrue(claude.endswith(body))
        self.assertIn('model: "haiku"', claude)
        self.assertNotIn("model =", codex)
        value = next(line.split(" = ", 1)[1] for line in codex.splitlines()
                     if line.startswith("developer_instructions = "))
        self.assertEqual(json.loads(value), body)
        self.assertEqual(json.loads((self.root / ".mcp.json").read_text())["mcpServers"], config["mcp_servers"])

    def test_disabled_agents_and_skills_cleanup_preserves_unmanaged_files(self):
        self.install()
        own = self.root / ".claude/agents/personal.md"
        own.write_text("User-maintained role.\n")
        config = self.config()
        config["disabled_agents"] = ["log-analyzer"]
        config["skills"] = []
        self.write_config(config)
        self.sync()
        self.assertFalse((self.root / ".claude/agents/log-analyzer.md").exists())
        self.assertFalse((self.root / ".agents/skills/maintain-knowledge").is_symlink())
        self.assertTrue((self.root / ".agents/framework/skills/maintain-knowledge/SKILL.md").is_file())
        self.assertEqual(own.read_text(), "User-maintained role.\n")

    def test_modified_obsolete_agent_is_preserved(self):
        self.install()
        agent = self.root / ".codex/agents/log-analyzer.toml"
        agent.write_text(agent.read_text() + "# Local work.\n")
        config = self.config()
        config["disabled_agents"] = ["log-analyzer"]
        self.write_config(config)
        before = self.snapshot()
        with self.assertRaisesRegex(ValueError, "local or unowned"):
            self.sync()
        self.assertEqual(self.snapshot(), before)

    def test_hooks_stay_native_and_are_never_executed_by_sync(self):
        self.install()
        config = self.config()
        native = {"hooks": {"SessionStart": [{"hooks": [{"type": "command", "command": "echo fixture"}]}]}}
        config["codex_hooks"] = native
        config["claude_settings"] = native
        self.write_config(config)
        self.sync()
        self.assertEqual(json.loads((self.root / ".codex/hooks.json").read_text()), native)
        self.assertEqual(json.loads((self.root / ".claude/settings.json").read_text()), native)
        config["codex_hooks"] = {}
        self.write_config(config)
        self.sync()
        self.assertFalse((self.root / ".codex/hooks.json").exists())

    def test_source_traversal_and_parent_symlinks_are_rejected(self):
        self.install()
        config = self.config()
        config["agents"] = [{"name": "test-runner", "instructions": "agents/../../outside.md"}]
        self.write_config(config)
        with self.assertRaisesRegex(ValueError, "inside framework/agents"):
            self.sync()
        other = self.base / "outside"
        other.mkdir()
        target = self.base / "fresh"
        target.mkdir()
        (target / ".codex").symlink_to(other, target_is_directory=True)
        with self.assertRaisesRegex(ValueError, "directory symlink"):
            harness.install_plan(target)
        self.assertEqual(list(other.iterdir()), [])
        self.assertFalse((target / ".agents").exists())

    def test_invalid_skill_metadata_resources_and_project_links_fail(self):
        self.install()
        path = self.root / ".agents/skills/example/SKILL.md"
        path.parent.mkdir()
        valid = "---\nname: example\ndescription: A skill.\n---\n\n"
        for content, message in [
            (valid + "Read [source](missing.md).\n", "missing resource"),
            (valid.replace("name: example", "name: example\ncontext: fork"), "frontmatter"),
            (valid.replace("description: A skill.", "description: >\n  A skill."), "frontmatter"),
            (valid.replace("description: A skill.", "description: [array]"), "ambiguous YAML"),
            (valid.replace("description: A skill.", "description: yes"), "ambiguous YAML"),
            (valid + "```python\nprint(1)\n", "code fence"),
        ]:
            path.write_text(content)
            with self.assertRaisesRegex(ValueError, message):
                self.sync()
        path.write_text(valid)
        (self.root / ".agents/project/index.md").write_text("Read [topic](missing-topic.md).\n")
        with self.assertRaisesRegex(ValueError, "missing resource"):
            self.sync()

    def test_unknown_schema_and_unsafe_lock_paths_are_rejected(self):
        self.install()
        config = self.config()
        config["schema_version"] = 99
        self.write_config(config)
        with self.assertRaisesRegex(ValueError, "schema_version"):
            self.sync()
        config["schema_version"] = 1
        self.write_config(config)
        lock = self.root / harness.GENERATED_LOCK
        value = json.loads(lock.read_text())
        value["files"]["../../outside"] = "sha256:invalid"
        lock.write_text(json.dumps(value))
        with self.assertRaisesRegex(ValueError, "outside managed scope"):
            self.sync()

    def test_obsolete_package_file_removal_is_tracked(self):
        source = self.package_copy()
        old = source / "skills/financial-analyzing/reference/old.md"
        old.write_text("An obsolete framework reference.\n")
        self.install(source)
        installed = self.root / ".agents/framework/skills/financial-analyzing/reference/old.md"
        self.assertTrue(installed.is_file())
        old.unlink()
        self.install(source)
        self.assertFalse(installed.exists())


if __name__ == "__main__":
    unittest.main()

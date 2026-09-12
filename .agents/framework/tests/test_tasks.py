"""Regression cases for task records, handoff memory, and upgrade ownership."""

import json
import shutil
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
import harness
import tasks


class TaskWorkflowTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="task records ")
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name).resolve()
        harness.apply_plan(self.root, harness.install_plan(self.root))
        self.name = "20260912-project-task"

    def create(self):
        tasks.create(self.root, self.name, "项目任务")
        return self.root / "docs/tasks" / self.name

    def mark_no_ui(self, folder):
        (folder / "ui-cases.json").write_text(json.dumps({
            "schema_version": 1, "applicable": False, "reason": "API-only change", "cases": []}))

    def test_create_task_and_refuse_duplicate_without_losing_work(self):
        folder = self.create()
        self.assertEqual(set(p.name for p in folder.iterdir()), set(tasks.FILES))
        self.assertEqual(tasks.validate(self.root)[self.name][0], "planned")
        (folder / "01-prd.md").write_text("User-maintained acceptance criteria.\n")
        before = (self.root / "docs/status.md").read_bytes()
        with self.assertRaisesRegex(ValueError, "already exists"):
            self.create()
        self.assertEqual((folder / "01-prd.md").read_text(), "User-maintained acceptance criteria.\n")
        self.assertEqual((self.root / "docs/status.md").read_bytes(), before)
        self.assertEqual(harness.sync_plan(self.root), {})

    def test_start_requires_ui_applicability_and_completion_requires_verification(self):
        folder = self.create()
        with self.assertRaisesRegex(ValueError, "decide whether UI"):
            tasks.set_status(self.root, self.name, "active", "Start implementation")
        self.mark_no_ui(folder)
        tasks.set_status(self.root, self.name, "active", "No UI; implement API behavior")
        with self.assertRaisesRegex(ValueError, "passed verification"):
            tasks.set_status(self.root, self.name, "done", "Completed")
        (folder / "05-verification.md").write_text("verification_status: passed\n\nCommand and evidence recorded.\n")
        tasks.set_status(self.root, self.name, "done", "Tests passed; retain decisions")
        self.assertEqual(tasks.validate(self.root)[self.name], ("done", "项目任务", "Tests passed; retain decisions"))

    def test_missing_material_and_unregistered_task_are_detected(self):
        folder = self.create()
        (folder / "01-prd.md").unlink()
        with self.assertRaises(OSError):
            tasks.validate(self.root)
        (folder / "01-prd.md").write_text("Restored PRD.\n")
        status = self.root / "docs/status.md"
        original = status.read_text()
        status.write_text("\n".join(line for line in original.splitlines() if self.name not in line))
        with self.assertRaisesRegex(ValueError, "entries must match"):
            tasks.validate(self.root)

    def test_upgrade_preserves_documents_cases_status_and_api_tests(self):
        folder = self.create()
        self.mark_no_ui(folder)
        tasks.set_status(self.root, self.name, "blocked", "Test service prerequisite remains")
        own = self.root / "tests/api/test_project.py"
        own.write_text("# project-owned API scenarios\n")
        (folder / "01-prd.md").write_text("Project acceptance criteria.\n")
        preserved = {path: path.read_bytes() for path in [own, folder / "01-prd.md",
                                                          folder / "ui-cases.json", self.root / "docs/status.md"]}
        with tempfile.TemporaryDirectory(prefix="framework release ") as release:
            source = Path(release) / "source"
            shutil.copytree(harness.PACKAGE, source, ignore=shutil.ignore_patterns("__pycache__"))
            (source / "VERSION").write_text("0.3.0\n")
            (source / "templates/docs/status.md").write_text("This must not replace project memory.\n")
            harness.apply_plan(self.root, harness.install_plan(self.root, source))
        for path, original in preserved.items():
            self.assertEqual(path.read_bytes(), original)

    def test_invalid_task_id_or_symlink_cannot_write_outside_tasks(self):
        with self.assertRaisesRegex(ValueError, "task ID"):
            tasks.create(self.root, "../../outside", "Invalid")
        outside = self.root / "outside"
        outside.mkdir()
        (self.root / "docs/tasks" / self.name).symlink_to(outside, target_is_directory=True)
        with self.assertRaises(ValueError):
            tasks.create(self.root, self.name, "Existing link")
        self.assertEqual(list(outside.iterdir()), [])

    def test_browser_case_requires_observable_step_assertions(self):
        folder = self.create()
        case = {"id": "UI-001", "title": "Submit empty form", "role": "visitor",
                "start_url": "${WEB_BASE_URL}/login", "preconditions": ["fresh session"],
                "steps": [{"action": "click", "target": "Sign in button", "expected": "Required fields shown"}],
                "cleanup": ["close isolated session"]}
        path = folder / "ui-cases.json"
        data = {"schema_version": 1, "applicable": True, "reason": "Login page changed", "cases": [case]}
        path.write_text(json.dumps(data))
        tasks.set_status(self.root, self.name, "active", "UI case is ready")
        del case["steps"][0]["expected"]
        path.write_text(json.dumps(data))
        with self.assertRaisesRegex(ValueError, "expected result"):
            tasks.validate(self.root)

    def test_status_update_keeps_memory_outside_the_table(self):
        folder = self.create()
        self.mark_no_ui(folder)
        path = self.root / "docs/status.md"
        old_row = next(line for line in path.read_text().splitlines() if line.startswith("| ["))
        path.write_text(path.read_text() + "\nHistorical record:\n" + old_row + "\n")
        tasks.set_status(self.root, self.name, "active", "Implementation started")
        self.assertTrue(path.read_text().endswith(old_row + "\n"))


if __name__ == "__main__":
    unittest.main()

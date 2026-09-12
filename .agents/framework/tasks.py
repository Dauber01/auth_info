#!/usr/bin/env python3
"""Create task records and maintain docs/status.md using Python 3.9+ only."""

import argparse
import json
import re
import sys
from pathlib import Path

ID = re.compile(r"[0-9]{8}-[a-z0-9]+(?:-[a-z0-9]+)*")
STATES = {"planned", "active", "blocked", "review", "done"}
FILES = ("00-plan.md", "01-prd.md", "02-prototype.md", "03-design.md",
         "04-test-cases.md", "ui-cases.json", "05-verification.md")
START, END = "<!-- TASKS:START -->", "<!-- TASKS:END -->"
ROW = re.compile(r"^\| \[([^]]+)\]\(tasks/([^/]+)/00-plan\.md\) \| ([^|]+) \| ([^|]+) \| ([^|]*) \|$", re.M)


def status_rows(content):
    if content.count(START) != 1 or content.count(END) != 1 or content.index(START) > content.index(END):
        raise ValueError("docs/status.md requires one ordered TASKS marker pair")
    section = content.split(START)[1].split(END)[0]
    rows = {}
    for match in ROW.finditer(section):
        name, folder, state, title, memory = (part.strip() for part in match.groups())
        if not ID.fullmatch(name) or name != folder or name in rows or state not in STATES:
            raise ValueError("invalid or duplicate task entry in docs/status.md")
        rows[name] = (state, title, memory)
    for line in section.splitlines():
        if line.startswith("| [") and not ROW.fullmatch(line):
            raise ValueError("malformed task row in docs/status.md")
    return rows


def validate_ui(path, state):
    value = json.loads(path.read_text(encoding="utf-8"))
    if set(value) != {"schema_version", "applicable", "reason", "cases"} or value["schema_version"] != 1:
        raise ValueError(f"{path}: invalid UI case schema")
    applicable = value["applicable"]
    if applicable is not None and not isinstance(applicable, bool):
        raise ValueError(f"{path}: applicable must be true, false, or null")
    if not isinstance(value["reason"], str) or not isinstance(value["cases"], list):
        raise ValueError(f"{path}: reason must be text and cases must be a list")
    if applicable is None and state != "planned":
        raise ValueError(f"{path}: decide whether UI cases apply before starting implementation")
    if applicable is False and (not value["reason"].strip() or value["cases"]):
        raise ValueError(f"{path}: non-applicable UI requires a reason and empty cases")
    if applicable is True and not value["cases"]:
        raise ValueError(f"{path}: applicable UI requires executable cases")
    names = set()
    for case in value["cases"]:
        required = {"id", "title", "role", "start_url", "preconditions", "steps", "cleanup"}
        if not isinstance(case, dict) or set(case) != required:
            raise ValueError(f"{path}: UI cases require {sorted(required)}")
        for key in ("id", "title", "role", "start_url"):
            if not isinstance(case[key], str) or not case[key].strip():
                raise ValueError(f"{path}: {key} must be nonempty text")
        if case["id"] in names:
            raise ValueError(f"{path}: duplicate UI case ID")
        names.add(case["id"])
        for key in ("preconditions", "cleanup"):
            if not isinstance(case[key], list) or not all(isinstance(item, str) for item in case[key]):
                raise ValueError(f"{path}: {key} must be a text list")
        if not isinstance(case["steps"], list) or not case["steps"]:
            raise ValueError(f"{path}: UI cases require steps")
        for step in case["steps"]:
            if (not isinstance(step, dict) or set(step) != {"action", "target", "expected"} or
                    not all(isinstance(item, str) and item.strip() for item in step.values())):
                raise ValueError(f"{path}: each UI step needs action, target, and expected result")


def validate(root, seeded_status=None):
    from harness import read_source, safe_path
    content = seeded_status if seeded_status is not None else read_source(root, "docs/status.md")
    rows = status_rows(content)
    directory = safe_path(root, "docs/tasks")
    if directory.is_symlink():
        raise ValueError("docs/tasks must not be a symlink")
    folders = {p.name: p for p in directory.glob("*") if p.is_dir()}
    if set(rows) != set(folders):
        raise ValueError("task directories and docs/status.md entries must match")
    for name, folder in folders.items():
        if folder.is_symlink() or not ID.fullmatch(name):
            raise ValueError(f"invalid task directory: {name}")
        for filename in FILES:
            text = read_source(root, f"docs/tasks/{name}/{filename}")
            if not text.strip():
                raise ValueError(f"{name}/{filename}: empty task material")
        state = rows[name][0]
        validate_ui(folder / "ui-cases.json", state)
        if state == "done":
            verification = (folder / "05-verification.md").read_text(encoding="utf-8")
            if not re.search(r"^verification_status: passed$", verification, re.M):
                raise ValueError(f"{name}: done requires a passed verification record")
    return rows


def cell(value):
    if not value.strip() or any(char in value for char in "|\r\n"):
        raise ValueError("task title and memory must be nonempty single lines without a pipe")
    return value.strip()


def create(root, name, title):
    from harness import apply_plan, read_source, safe_path, text_entry
    if not ID.fullmatch(name):
        raise ValueError("task ID must be YYYYMMDD-lowercase-slug")
    title = cell(title)
    rows = validate(root)
    if name in rows or safe_path(root, f"docs/tasks/{name}").exists():
        raise ValueError(f"task already exists: {name}")
    template = Path(__file__).resolve().parent / "task-templates"
    changes = {}
    for filename in FILES:
        body = (template / filename).read_text(encoding="utf-8")
        changes[f"docs/tasks/{name}/{filename}"] = text_entry(
            body.replace("{{task_id}}", name).replace("{{title}}", title))
    status = read_source(root, "docs/status.md")
    row = f"| [{name}](tasks/{name}/00-plan.md) | planned | {title} | 已建任务材料；先完善计划、PRD、原型、设计与用例。 |\n"
    changes["docs/status.md"] = text_entry(status.replace(END, row + END))
    for relative in changes:
        path = safe_path(root, relative)
        if path.is_symlink():
            raise ValueError(f"refusing task file symlink: {relative}")
    apply_plan(root, changes)


def set_status(root, name, state, memory):
    from harness import apply_plan, read_source, text_entry
    if state not in STATES:
        raise ValueError("unknown task status")
    memory = cell(memory)
    rows = validate(root)
    if name not in rows:
        raise ValueError(f"task does not exist: {name}")
    content = read_source(root, "docs/status.md")
    def replace(match):
        if match[1] != name:
            return match[0]
        return f"| [{name}](tasks/{name}/00-plan.md) | {state} | {rows[name][1]} | {memory} |"
    before, section = content.split(START, 1)
    section, after = section.split(END, 1)
    updated = before + START + ROW.sub(replace, section) + END + after
    validate(root, updated)
    apply_plan(root, {"docs/status.md": text_entry(updated)})


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, default=Path.cwd())
    sub = parser.add_subparsers(dest="command", required=True)
    new = sub.add_parser("new")
    new.add_argument("--id", required=True)
    new.add_argument("--title", required=True)
    update = sub.add_parser("status")
    update.add_argument("--id", required=True)
    update.add_argument("--state", choices=sorted(STATES), required=True)
    update.add_argument("--memory", required=True)
    sub.add_parser("check")
    args = parser.parse_args()
    root = args.root.resolve()
    try:
        if args.command == "new":
            create(root, args.id, args.title)
        elif args.command == "status":
            set_status(root, args.id, args.state, args.memory)
        else:
            validate(root)
        print("Task records are consistent.")
        return 0
    except (OSError, ValueError, TypeError, KeyError) as error:
        print(f"tasks: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())

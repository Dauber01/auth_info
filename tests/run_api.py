#!/usr/bin/env python3
"""Build an isolated Go HTTP fixture and run Python API cases against it."""

import json
import os
import queue
import subprocess
import sys
import tempfile
import threading
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def run_cases():
    # This branch runs only in the child after the fixture address is known.
    sys.path.insert(0, str(ROOT))
    suite = unittest.defaultTestLoader.discover(str(ROOT / "tests/api"), top_level_dir=str(ROOT))
    if suite.countTestCases() == 0:
        raise RuntimeError("no API cases discovered")
    return 0 if unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful() else 1


def main():
    if sys.argv[1:] == ["--run-cases"]:
        return run_cases()
    if sys.argv[1:]:
        raise ValueError("use tests/run_api.py without arguments")
    with tempfile.TemporaryDirectory(prefix="auth-info-api-") as directory:
        binary = Path(directory) / ("server.exe" if os.name == "nt" else "server")
        subprocess.run(["go", "build", "-o", str(binary), "./tests/api/server"],
                       cwd=ROOT, check=True, timeout=120)
        process = subprocess.Popen([str(binary)], cwd=ROOT, stdout=subprocess.PIPE, text=True)
        try:
            ready = queue.Queue()
            threading.Thread(target=lambda: ready.put(process.stdout.readline()), daemon=True).start()
            try:
                line = ready.get(timeout=20)
            except queue.Empty as error:
                raise RuntimeError("API fixture did not become ready within 20 seconds") from error
            if not line:
                raise RuntimeError("API fixture exited before reporting its address")
            base_url = json.loads(line)["base_url"]
            env = {**os.environ, "API_BASE_URL": base_url, "PYTHONDONTWRITEBYTECODE": "1"}
            result = subprocess.run([sys.executable, "-B", str(Path(__file__).resolve()), "--run-cases"],
                                    cwd=ROOT, env=env, timeout=90)
            return result.returncode
        finally:
            if process.poll() is None:
                process.terminate()
                try:
                    process.wait(timeout=8)
                except subprocess.TimeoutExpired:
                    process.kill()
                    process.wait(timeout=5)
            process.stdout.close()


if __name__ == "__main__":
    try:
        sys.exit(main())
    except (OSError, ValueError, RuntimeError, subprocess.SubprocessError) as error:
        print(f"API test setup failed: {error}", file=sys.stderr)
        sys.exit(1)

import subprocess
import tempfile
import unittest
from pathlib import Path


CHECKER = Path(__file__).resolve().parents[1] / "check_file_lines.sh"


def _write_lines(root: Path, relative: str, count: int) -> None:
    path = root / relative
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text("x\n" * count, encoding="utf-8")


def _stage_all(root: Path) -> None:
    subprocess.run(["git", "init", "-q"], cwd=root, check=True)
    subprocess.run(["git", "add", "."], cwd=root, check=True)


class CheckFileLinesTest(unittest.TestCase):
    def test_python_warning_and_generated_trees_are_not_scanned(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            _write_lines(root, "warning.py", 301)
            for tree in ("vendor", "cache", "generated"):
                _write_lines(root, f"{tree}/ignored.py", 501)
            _write_lines(root, "custom/ignored.py", 501)
            (root / ".line-limit-ignore").write_text("custom/*\n", encoding="utf-8")
            _stage_all(root)

            completed = subprocess.run(
                ["bash", str(CHECKER), "--all"],
                cwd=root,
                capture_output=True,
                text=True,
                check=False,
            )

            self.assertEqual(completed.returncode, 0, completed.stdout + completed.stderr)
            self.assertIn("warning.py", completed.stdout)
            for ignored in ("vendor/ignored.py", "cache/ignored.py", "generated/ignored.py", "custom/ignored.py"):
                self.assertNotIn(ignored, completed.stdout)

    def test_unignored_python_over_hard_limit_blocks(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            _write_lines(root, "too_long.py", 501)
            _stage_all(root)

            completed = subprocess.run(
                ["bash", str(CHECKER), "--all"],
                cwd=root,
                capture_output=True,
                text=True,
                check=False,
            )

            self.assertEqual(completed.returncode, 1, completed.stdout + completed.stderr)
            self.assertIn("too_long.py", completed.stdout)
            self.assertIn("ERROR files", completed.stdout)


if __name__ == "__main__":
    unittest.main()

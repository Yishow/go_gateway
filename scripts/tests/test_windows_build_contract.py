"""Build-command contracts; native artifact/runtime acceptance is separate."""
from pathlib import Path
import subprocess
import unittest

ROOT = Path(__file__).resolve().parents[2]

class WindowsBuildContract(unittest.TestCase):
    def plan(self, mode, target):
        return subprocess.run(["make", "-n", "build", f"BUILD_MODE={mode}", f"TARGET_GOOS={target}"], cwd=ROOT, text=True, capture_output=True, check=False)

    def test_desktop_flags_and_artifact(self):
        result = self.plan("desktop", "windows")
        self.assertEqual(0, result.returncode, result.stderr)
        for required in ("-tags desktop", "-H=windowsgui", "bin/gateway-desktop.exe", "GOOS=windows"):
            self.assertIn(required, result.stdout)

    def test_console_and_nonwindows_preserved(self):
        for target in ("windows", "linux", "darwin"):
            result = self.plan("console", target)
            self.assertEqual(0, result.returncode, result.stderr)
            self.assertNotIn("-H=windowsgui", result.stdout)
            self.assertNotIn("-tags desktop", result.stdout)
            self.assertIn("bin/test-ui.exe", result.stdout)

    def test_invalid_platform_or_mode_fails(self):
        self.assertNotEqual(0, self.plan("desktop", "linux").returncode)
        self.assertNotEqual(0, self.plan("invalid", "windows").returncode)

    def test_frontend_must_precede_backend(self):
        script = (ROOT / "scripts/build.ps1").read_text()
        for required in ('[ValidateSet("console", "desktop")]', '"-tags", "desktop"', '"-H=windowsgui"', '"index.html"', "gateway-desktop.exe"):
            self.assertIn(required, script)
        self.assertLess(script.index("$frontendBuildExitCode -ne 0"), script.index("& go build"))

    def test_automatic_commit_metadata_marks_dirty_worktree(self):
        self.assertIn("+dirty", (ROOT / "Makefile").read_text())
        self.assertIn("+dirty", (ROOT / "scripts/build.ps1").read_text())

if __name__ == "__main__":
    unittest.main()

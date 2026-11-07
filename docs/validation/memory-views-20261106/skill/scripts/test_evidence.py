"""Offline CLI regression tests; no application, database or network access."""
import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

SCRIPT = Path(__file__).with_name('evidence.py')


class EvidenceCLI(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix='recording-skill-check-')
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.repo = self.root / 'repo'
        self.repo.mkdir()
        self.git('init', '-q')
        (self.repo / 'source.txt').write_text('candidate-v1\n')
        self.git('add', 'source.txt')
        self.git('-c', 'user.name=Skill Test', '-c', 'user.email=skill@example.invalid',
                 '-c', 'commit.gpgsign=false', '-c', 'core.hooksPath=/dev/null', 'commit', '-qm', 'fixture')
        self.run = self.root / 'evidence'
        result = self.cli('init', '--repo', str(self.repo), '--out', str(self.run))
        self.assertEqual(result.returncode, 0, result.stderr)

    def git(self, *args):
        return subprocess.run(['git', '-C', str(self.repo), *args], check=True, capture_output=True)

    def cli(self, *args):
        return subprocess.run([sys.executable, str(SCRIPT), *args], capture_output=True, text=True)

    def complete(self):
        checks = json.loads((self.run / 'checks.json').read_text())
        (self.run / 'observations.json').write_text('{"synthetic_cli_test":true}\n')
        for row in checks:
            row.update(status='PASS', detail='Synthetic CLI test only, not product acceptance.',
                       evidence=['observations.json'])
        self.save(checks)
        return checks

    def save(self, rows):
        (self.run / 'checks.json').write_text(json.dumps(rows))

    def verify(self):
        return self.cli('check', '--run', str(self.run))

    def test_fresh_checklist_has_29_pending_gates_and_no_databases(self):
        rows = json.loads((self.run / 'checks.json').read_text())
        self.assertEqual(len(rows), 29)
        self.assertTrue(all(row['status'] == 'NOT_RUN' for row in rows))
        self.assertEqual(self.verify().returncode, 1)
        self.assertFalse(list(self.root.rglob('*.db')))

    def test_complete_structure_passes_but_labels_its_limit(self):
        self.complete()
        result = self.verify()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn('agent must review', result.stdout)

    def test_missing_sql_readback_or_duplicate_check_cannot_pass(self):
        rows = self.complete()
        self.save([row for row in rows if row['id'] != 'recording/mc_3e/sqlite/sql-readback'])
        self.assertEqual(self.verify().returncode, 1)
        self.save([row for row in rows if row['id'] != 'source/mc_3e/device-create'])
        self.assertEqual(self.verify().returncode, 1)
        self.save(rows + [rows[0]])
        self.assertEqual(self.verify().returncode, 1)

    def test_blocked_gate_and_empty_attachment_cannot_pass(self):
        rows = self.complete()
        rows[0]['status'] = 'BLOCKED'
        self.save(rows)
        self.assertEqual(self.verify().returncode, 1)
        rows[0]['status'] = 'PASS'
        rows[0]['evidence'] = []
        self.save(rows)
        self.assertEqual(self.verify().returncode, 1)
        rows[0]['evidence'] = ['missing.json']
        self.save(rows)
        self.assertEqual(self.verify().returncode, 1)
        rows[0]['evidence'] = ['observations.json']
        (self.run / 'observations.json').write_text('')
        self.save(rows)
        self.assertEqual(self.verify().returncode, 1)

    def test_outside_path_and_symlink_cannot_be_used_as_evidence(self):
        rows = self.complete()
        external = self.root / 'outside.json'
        external.write_text('{}')
        rows[0]['evidence'] = ['../outside.json']
        self.save(rows)
        self.assertEqual(self.verify().returncode, 1)
        (self.run / 'linked.json').symlink_to(external)
        rows[0]['evidence'] = ['linked.json']
        self.save(rows)
        self.assertEqual(self.verify().returncode, 1)

    def test_source_drift_cannot_inherit_pass(self):
        self.complete()
        (self.repo / 'source.txt').write_text('candidate-v2\n')
        result = self.verify()
        self.assertEqual(result.returncode, 1)
        self.assertIn('candidate source changed', result.stderr)

    def test_source_drift_hidden_by_git_index_flag_cannot_inherit_pass(self):
        self.complete()
        self.git('update-index', '--assume-unchanged', 'source.txt')
        (self.repo / 'source.txt').write_text('candidate-v2\n')
        result = self.verify()
        self.assertEqual(result.returncode, 1)
        self.assertIn('candidate source changed', result.stderr)

    @unittest.skipIf(os.name == 'nt', 'POSIX executable permission is unavailable on Windows')
    def test_untracked_script_permission_change_cannot_inherit_pass(self):
        script = self.repo / 'launch.sh'
        script.write_text('#!/bin/sh\necho candidate\n')
        script.chmod(0o644)
        self.run = self.root / 'evidence-untracked'
        result = self.cli('init', '--repo', str(self.repo), '--out', str(self.run))
        self.assertEqual(result.returncode, 0, result.stderr)
        self.complete()
        script.chmod(0o755)
        result = self.verify()
        self.assertEqual(result.returncode, 1)
        self.assertIn('candidate source changed', result.stderr)
        self.run = self.root / 'evidence-untracked-executable'
        result = self.cli('init', '--repo', str(self.repo), '--out', str(self.run))
        self.assertEqual(result.returncode, 0, result.stderr)
        self.complete()
        script.chmod(0o700)
        result = self.verify()
        self.assertEqual(result.returncode, 1)
        self.assertIn('candidate source changed', result.stderr)

    @unittest.skipIf(os.name == 'nt', 'Directory symlinks require Windows developer privileges')
    def test_untracked_source_symlink_target_change_cannot_inherit_pass(self):
        first = self.root / 'first-source'
        second = self.root / 'second-source'
        first.mkdir()
        second.mkdir()
        (first / 'config.txt').write_text('candidate-v1\n')
        (second / 'config.txt').write_text('candidate-v2\n')
        link = self.repo / 'linked-source'
        link.symlink_to(first, target_is_directory=True)
        self.run = self.root / 'evidence-symlink'
        result = self.cli('init', '--repo', str(self.repo), '--out', str(self.run))
        self.assertEqual(result.returncode, 0, result.stderr)
        self.complete()
        link.unlink()
        link.symlink_to(second, target_is_directory=True)
        result = self.verify()
        self.assertEqual(result.returncode, 1)
        self.assertIn('candidate source changed', result.stderr)

    def test_scope_cannot_be_silently_shrunk(self):
        self.complete()
        manifest = json.loads((self.run / 'manifest.json').read_text())
        manifest['scope']['protocols'] = ['mc_3e']
        (self.run / 'manifest.json').write_text(json.dumps(manifest))
        self.assertEqual(self.verify().returncode, 1)

    def test_existing_run_is_preserved(self):
        original = (self.run / 'manifest.json').read_bytes()
        result = self.cli('init', '--repo', str(self.repo), '--out', str(self.run))
        self.assertEqual(result.returncode, 1)
        self.assertEqual((self.run / 'manifest.json').read_bytes(), original)


if __name__ == '__main__':
    unittest.main()

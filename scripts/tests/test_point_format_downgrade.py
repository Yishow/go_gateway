import hashlib
import json
from pathlib import Path
import sqlite3
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
SCRIPT = ROOT / "scripts/check_point_format_downgrade.py"


class PointFormatDowngradeTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.path = Path(self.tmp.name) / "settings.db"
        with sqlite3.connect(self.path) as db:
            db.execute("CREATE TABLE points(id TEXT PRIMARY KEY,data_format TEXT)")
            db.execute("INSERT INTO points VALUES('p','CDAB')")

    def check(self, path=None):
        proc = subprocess.run(["python3", str(SCRIPT), "--sqlite-file", str(path or self.path),
                               "--target-format-version", "0"], capture_output=True, text=True)
        return proc.returncode, json.loads(proc.stdout)

    def migrate(self):
        with sqlite3.connect(self.path) as db:
            db.executescript((ROOT / "internal/datalink/schema/migrations/048_point_format_intent_sqlite.up.sql").read_text())

    def test_legacy_and_additive_migration_preserve_intent(self):
        self.assertEqual(0, self.check()[0])
        self.migrate()
        self.assertEqual(0, self.check()[0])
        with sqlite3.connect(self.path) as db:
            self.assertEqual(("CDAB", None, None, 0), db.execute("SELECT data_format,format_mode,format_override,format_version FROM points").fetchone())

    def test_override_inherit_and_unknown_block_without_writes(self):
        self.migrate()
        for mode, override, version, category in [("override", "ABCD", 1, "override"),
                ("inherit", "", 1, "inherit"), ("legacy", "ABCD", 0, "unknown"),
                ("legacy", "", 2, "unknown")]:
            with self.subTest(mode=mode, version=version):
                with sqlite3.connect(self.path) as db:
                    db.execute("UPDATE points SET format_mode=?,format_override=?,format_version=?", (mode, override, version))
                before = self.path.read_bytes()
                code, result = self.check()
                self.assertEqual(2, code)
                self.assertFalse(result["compatible_for_point_format"])
                self.assertEqual(1, result["blockers"][category])
                self.assertEqual(before, self.path.read_bytes())
                self.assertNotIn("p", result.get("point_ids", []))

    def test_missing_or_partial_schema_fails_closed(self):
        missing = Path(self.tmp.name) / "missing.db"
        self.assertEqual(1, self.check(missing)[0])
        self.assertFalse(missing.exists())
        with sqlite3.connect(self.path) as db:
            db.execute("ALTER TABLE points ADD COLUMN format_mode TEXT")
        self.assertEqual(1, self.check()[0])

    def test_base_columns_and_complete_expand_are_required_even_if_empty(self):
        schemas = ["id TEXT", "data_format TEXT",
                   "id TEXT,data_format TEXT,revision_id TEXT",
                   "data_format TEXT,format_mode TEXT,format_override TEXT,format_version INTEGER,revision_id TEXT",
                   "id TEXT,format_mode TEXT,format_override TEXT,format_version INTEGER,revision_id TEXT"]
        for index, schema in enumerate(schemas):
            for populated in (False, True):
                with self.subTest(schema=schema, populated=populated):
                    path = Path(self.tmp.name) / f"partial-{index}-{populated}.db"
                    with sqlite3.connect(path) as db:
                        db.execute(f"CREATE TABLE points({schema})")
                        if populated:
                            db.execute("INSERT INTO points DEFAULT VALUES")
                    before = path.read_bytes()
                    code, result = self.check(path)
                    self.assertEqual(1, code)
                    self.assertFalse(result["compatible_for_point_format"])
                    self.assertEqual("compatibility_evidence_unavailable", result["error"])
                    self.assertEqual(before, path.read_bytes())

    def test_disposable_backup_restore_never_changes_historical_sql(self):
        history = Path(self.tmp.name) / "history.db"
        with sqlite3.connect(history) as db:
            db.execute("CREATE TABLE samples(id INTEGER,value TEXT,source_revision TEXT)")
            db.execute("INSERT INTO samples VALUES(1,'18446744073709551615','frozen-old-source')")
        historical_hash = hashlib.sha256(history.read_bytes()).hexdigest()
        backup = Path(self.tmp.name) / "before-expand.db"
        with sqlite3.connect(self.path) as source, sqlite3.connect(backup) as dest:
            source.backup(dest)
        self.migrate()
        with sqlite3.connect(self.path) as db:
            db.execute("UPDATE points SET format_mode='override',format_override='DCBA',format_version=1")
        self.assertEqual(2, self.check()[0])
        # Restore only into a NEW disposable path; never overwrite the live file
        # or accepted records. Keep the post-migration source intact as evidence.
        restored = Path(self.tmp.name) / "restored-copy.db"
        with sqlite3.connect(backup) as source, sqlite3.connect(restored) as dest:
            source.backup(dest)
        self.assertEqual(0, self.check(restored)[0])
        self.assertEqual(2, self.check()[0])
        with sqlite3.connect(restored) as db:
            self.assertEqual("ok", db.execute("PRAGMA integrity_check").fetchone()[0])
            self.assertEqual(("p", "CDAB"), db.execute("SELECT * FROM points").fetchone())
        self.assertEqual(historical_hash, hashlib.sha256(history.read_bytes()).hexdigest())


if __name__ == "__main__":
    unittest.main()

#!/usr/bin/env python3
"""Read-only point-format gate for an operator-managed legacy binary rollback.

This checks one compatibility dimension; it never deploys, restores or edits data.
"""
import argparse
import json
from pathlib import Path
import sqlite3
import sys


def inspect(path):
    path = Path(path).resolve(strict=True)
    with sqlite3.connect(path.as_uri() + "?mode=ro", uri=True) as db:
        db.execute("PRAGMA query_only=ON")
        db.execute("BEGIN")
        columns = {row[1] for row in db.execute("PRAGMA table_info(points)")}
        if not {"id", "data_format"}.issubset(columns):
            raise ValueError("configuration schema unavailable")
        intent = {"format_mode", "format_override", "format_version", "revision_id"}
        present = columns & intent
        if present and present != intent:
            raise ValueError("partial intent schema")
        count = db.execute("SELECT count(*) FROM points").fetchone()[0]
        blockers = {"override": 0, "inherit": 0, "unknown": 0}
        if present:
            for mode, override, version in db.execute(
                    "SELECT format_mode,format_override,format_version FROM points"):
                if mode in (None, "", "legacy") and version == 0 and override in (None, ""):
                    continue
                if mode == "override" and version == 1 and override in ("ABCD", "BADC", "CDAB", "DCBA"):
                    blockers["override"] += 1
                elif mode == "inherit" and version == 1 and override in (None, ""):
                    # Legacy readers do not resolve inherited rule/connection policy.
                    # A coincidentally equal data_format projection is not proof.
                    blockers["inherit"] += 1
                else:
                    blockers["unknown"] += 1
        return {"scope": "point-format-only", "target_format_version": 0,
                "compatible_for_point_format": not any(blockers.values()),
                "point_count": count, "blockers": blockers}


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--sqlite-file", required=True)
    parser.add_argument("--target-format-version", type=int, choices=[0], required=True)
    args = parser.parse_args(argv)
    try:
        result = inspect(args.sqlite_file)
    except (OSError, sqlite3.Error, ValueError):
        print(json.dumps({"scope": "point-format-only", "compatible_for_point_format": False,
                          "error": "compatibility_evidence_unavailable"}))
        return 1
    print(json.dumps(result, sort_keys=True))
    return 0 if result["compatible_for_point_format"] else 2


if __name__ == "__main__":
    sys.exit(main())

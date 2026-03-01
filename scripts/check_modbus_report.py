#!/usr/bin/env python3
import argparse
import json
import sys
from pathlib import Path


def main() -> int:
    ap = argparse.ArgumentParser(description="Gate Modbus loadtest report")
    ap.add_argument("report", help="report json path")
    ap.add_argument("--require-pass", action="store_true", default=True)
    ap.add_argument("--max-mismatch", type=int, default=0)
    ap.add_argument("--max-bad-rows", type=int, default=0)
    ap.add_argument("--min-rows", type=int, default=1)
    args = ap.parse_args()

    p = Path(args.report)
    if not p.exists():
        print(f"[FAIL] report not found: {p}", file=sys.stderr)
        return 2

    data = json.loads(p.read_text())
    total_rows = int(data.get("total_rows", 0))
    bad_rows = int(data.get("bad_rows", 0))
    mismatch = int(data.get("mismatch_count", 0))
    passed = bool(data.get("pass", False))
    reasons = data.get("fail_reasons", []) or []

    print("[REPORT]", p)
    print(
        f"  pass={passed} rows={total_rows} bad_rows={bad_rows} mismatch={mismatch} "+
        f"collected={data.get('runtime_stats',{}).get('collected_total')} "+
        f"write_ok={data.get('runtime_stats',{}).get('write_success_total')} "+
        f"write_err={data.get('runtime_stats',{}).get('write_error_total')}"
    )

    failures = []
    if args.require_pass and not passed:
        failures.append("report.pass != true")
    if mismatch > args.max_mismatch:
        failures.append(f"mismatch_count={mismatch} > {args.max_mismatch}")
    if bad_rows > args.max_bad_rows:
        failures.append(f"bad_rows={bad_rows} > {args.max_bad_rows}")
    if total_rows < args.min_rows:
        failures.append(f"total_rows={total_rows} < {args.min_rows}")

    if failures:
        print("[FAIL] gate check failed:")
        for f in failures:
            print(" -", f)
        if reasons:
            print("[FAIL] report fail_reasons:")
            for r in reasons:
                print(" -", r)
        return 1

    print("[PASS] gate check ok")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

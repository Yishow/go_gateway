#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Run deterministic Modbus load tests (smoke + 10m + 30m) with report gating.

Usage:
  scripts/run_modbus_soak.sh [options]

Options:
  --workdir PATH          Repo path (default: current dir)
  --report-dir PATH       Output report dir (default: /tmp/modbus_reports_<ts>)
  --tcp N                 TCP server count (default: 15)
  --udp N                 UDP server count (default: 15)
  --interval-ms N         Polling interval ms (default: 500)
  --value-base N          Base register value for deterministic dataset (default: 1000)
  --tag-key-prefix TEXT   Tag key prefix (default: bench/server)
  --smoke-duration D      Smoke duration (default: 60s)
  --final-duration D      Final duration (default: 10m)
  --soak-duration D       Soak duration (default: 30m)
  --skip-smoke            Skip smoke stage
  --skip-final            Skip final 10m stage
  --skip-soak             Skip soak 30m stage
  --notify-on-fail        Send openclaw system event on failure (default: on)
  --no-notify-on-fail     Disable failure notification
  --dry-run               Print commands only
  -h, --help              Show help

Example:
  scripts/run_modbus_soak.sh \
    --workdir /home/yishow/github/go_gateway \
    --report-dir /tmp/modbus_reports_latest \
    --notify-on-fail
USAGE
}

WORKDIR="$(pwd)"
TS="$(date +%Y%m%d_%H%M%S)"
REPORT_DIR="/tmp/modbus_reports_${TS}"
TCP=15
UDP=15
INTERVAL_MS=500
VALUE_BASE=1000
TAG_KEY_PREFIX="bench/server"
SMOKE_DURATION="60s"
FINAL_DURATION="10m"
SOAK_DURATION="30m"
RUN_SMOKE=1
RUN_FINAL=1
RUN_SOAK=1
NOTIFY_ON_FAIL=1
DRY_RUN=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --workdir) WORKDIR="$2"; shift 2 ;;
    --report-dir) REPORT_DIR="$2"; shift 2 ;;
    --tcp) TCP="$2"; shift 2 ;;
    --udp) UDP="$2"; shift 2 ;;
    --interval-ms) INTERVAL_MS="$2"; shift 2 ;;
    --value-base) VALUE_BASE="$2"; shift 2 ;;
    --tag-key-prefix) TAG_KEY_PREFIX="$2"; shift 2 ;;
    --smoke-duration) SMOKE_DURATION="$2"; shift 2 ;;
    --final-duration) FINAL_DURATION="$2"; shift 2 ;;
    --soak-duration) SOAK_DURATION="$2"; shift 2 ;;
    --skip-smoke) RUN_SMOKE=0; shift ;;
    --skip-final) RUN_FINAL=0; shift ;;
    --skip-soak) RUN_SOAK=0; shift ;;
    --notify-on-fail) NOTIFY_ON_FAIL=1; shift ;;
    --no-notify-on-fail) NOTIFY_ON_FAIL=0; shift ;;
    --dry-run) DRY_RUN=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "[ERROR] unknown arg: $1" >&2; usage; exit 1 ;;
  esac
done

if [[ ! -d "$WORKDIR" ]]; then
  echo "[ERROR] workdir not found: $WORKDIR" >&2
  exit 1
fi

mkdir -p "$REPORT_DIR"

notify_fail() {
  local msg="$1"
  echo "[ALERT] $msg" >&2
  if [[ "$NOTIFY_ON_FAIL" -eq 1 ]] && command -v openclaw >/dev/null 2>&1; then
    openclaw system event --text "$msg" --mode now >/dev/null 2>&1 || true
  fi
}

run_stage() {
  local stage="$1"
  local duration="$2"

  local report_file="$REPORT_DIR/modbus_${stage}_report.json"
  local db_file="$REPORT_DIR/modbus_${stage}.db"

  local -a cmd=(
    docker run --rm
    -v "$WORKDIR":/workspace
    -v /home/yishow/.cache/go-mod:/go/pkg/mod
    -v /home/yishow/.cache/go-build:/root/.cache/go-build
    -v "$REPORT_DIR":"$REPORT_DIR"
    -w /workspace
    -e GOMODCACHE=/go/pkg/mod
    -e GOPROXY=off
    golang:1.25
    /usr/local/go/bin/go run ./cmd/loadtest_modbus
    -duration "$duration"
    -tcp "$TCP"
    -udp "$UDP"
    -interval-ms "$INTERVAL_MS"
    -value-base "$VALUE_BASE"
    -tag-key-prefix "$TAG_KEY_PREFIX"
    -db "$db_file"
    -report "$report_file"
  )

  echo "[RUN] stage=$stage duration=$duration report=$report_file"
  if [[ "$DRY_RUN" -eq 1 ]]; then
    printf '[DRY-RUN] %q ' "${cmd[@]}"; echo
    return 0
  fi

  (
    cd "$WORKDIR"
    "${cmd[@]}"
  )

  "$WORKDIR/scripts/check_modbus_report.py" "$report_file" \
    --require-pass \
    --max-mismatch 0 \
    --max-bad-rows 0 \
    --min-rows 1
}

if [[ "$RUN_SMOKE" -eq 1 ]]; then
  if ! run_stage "smoke" "$SMOKE_DURATION"; then
    notify_fail "FAIL: modbus smoke stage failed (report: $REPORT_DIR/modbus_smoke_report.json)"
    exit 1
  fi
fi

if [[ "$RUN_FINAL" -eq 1 ]]; then
  if ! run_stage "final_10m" "$FINAL_DURATION"; then
    notify_fail "FAIL: modbus final_10m stage failed (report: $REPORT_DIR/modbus_final_10m_report.json)"
    exit 1
  fi
fi

if [[ "$RUN_SOAK" -eq 1 ]]; then
  if ! run_stage "soak_30m" "$SOAK_DURATION"; then
    notify_fail "FAIL: modbus soak_30m stage failed (report: $REPORT_DIR/modbus_soak_30m_report.json)"
    exit 1
  fi
fi

if [[ "$DRY_RUN" -eq 0 ]]; then
  python3 - <<PY
import json
from pathlib import Path
root=Path("$REPORT_DIR")
files=[p for p in [root/"modbus_smoke_report.json", root/"modbus_final_10m_report.json", root/"modbus_soak_30m_report.json"] if p.exists()]
print("[SUMMARY]", root)
for f in files:
    d=json.loads(f.read_text())
    rs=d.get("runtime_stats",{})
    print(f" - {f.name}: pass={d.get('pass')} rows={d.get('total_rows')} mismatch={d.get('mismatch_count')} bad_rows={d.get('bad_rows')} collected={rs.get('collected_total')} write_ok={rs.get('write_success_total')} write_err={rs.get('write_error_total')}")
PY
fi

echo "[DONE] reports at: $REPORT_DIR"

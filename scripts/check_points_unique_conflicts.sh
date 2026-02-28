#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Check points unique-key conflicts before migration up/down.

Usage:
  scripts/check_points_unique_conflicts.sh --db <postgres|sqlite> --stage <up|down> [options]

Options:
  --db TYPE             Database type: postgres or sqlite
  --stage STAGE         Check stage: up or down
  --dsn TEXT            PostgreSQL DSN (required when --db postgres)
  --sqlite-file PATH    SQLite DB file path (required when --db sqlite)
  --show-sql            Print SQL before execution
  -h, --help            Show help

Exit code:
  0: no conflict
  2: conflict found
USAGE
}

DB_TYPE=""
STAGE=""
PG_DSN=""
SQLITE_FILE=""
SHOW_SQL=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --db) DB_TYPE="$2"; shift 2 ;;
    --stage) STAGE="$2"; shift 2 ;;
    --dsn) PG_DSN="$2"; shift 2 ;;
    --sqlite-file) SQLITE_FILE="$2"; shift 2 ;;
    --show-sql) SHOW_SQL=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "[ERROR] unknown arg: $1" >&2; usage; exit 1 ;;
  esac
done

if [[ "$DB_TYPE" != "postgres" && "$DB_TYPE" != "sqlite" ]]; then
  echo "[ERROR] --db must be postgres or sqlite" >&2
  exit 1
fi
if [[ "$STAGE" != "up" && "$STAGE" != "down" ]]; then
  echo "[ERROR] --stage must be up or down" >&2
  exit 1
fi
if [[ "$DB_TYPE" == "postgres" && -z "$PG_DSN" ]]; then
  echo "[ERROR] --dsn is required for postgres" >&2
  exit 1
fi
if [[ "$DB_TYPE" == "sqlite" && -z "$SQLITE_FILE" ]]; then
  echo "[ERROR] --sqlite-file is required for sqlite" >&2
  exit 1
fi
if [[ "$DB_TYPE" == "sqlite" && ! -f "$SQLITE_FILE" ]]; then
  echo "[ERROR] sqlite file not found: $SQLITE_FILE" >&2
  exit 1
fi

run_postgres_query() {
  local query="$1"
  if ! command -v psql >/dev/null 2>&1; then
    echo "[ERROR] psql not found in PATH" >&2
    exit 1
  fi
  psql "$PG_DSN" -v ON_ERROR_STOP=1 -At -F '|' -c "$query"
}

run_sqlite_query() {
  local query="$1"
  QUERY_TEXT="$query" python3 - "$SQLITE_FILE" <<'PY'
import os
import sqlite3
import sys

db_path = sys.argv[1]
query = os.environ["QUERY_TEXT"]

conn = sqlite3.connect(db_path)
try:
    cur = conn.cursor()
    rows = cur.execute(query).fetchall()
    for row in rows:
        print("|".join("" if value is None else str(value) for value in row))
finally:
    conn.close()
PY
}

run_query() {
  local query="$1"
  if [[ "$SHOW_SQL" -eq 1 ]]; then
    echo "[SQL]"
    echo "$query"
  fi
  if [[ "$DB_TYPE" == "postgres" ]]; then
    run_postgres_query "$query"
    return
  fi
  run_sqlite_query "$query"
}

TARGET_UP_QUERY=$(cat <<'SQL'
SELECT device_id, address, COALESCE(function, '') AS function_key, COUNT(*) AS cnt
FROM points
GROUP BY device_id, address, COALESCE(function, '')
HAVING COUNT(*) > 1
ORDER BY cnt DESC, device_id, address
LIMIT 200;
SQL
)

NORMALIZED_ALIAS_QUERY=$(cat <<'SQL'
WITH normalized AS (
  SELECT
    device_id,
    UPPER(TRIM(address)) AS address_norm,
    CASE LOWER(TRIM(COALESCE(function, '')))
      WHEN '' THEN '03'
      WHEN '3' THEN '03'
      WHEN '03' THEN '03'
      WHEN 'fc03' THEN '03'
      WHEN 'holding' THEN '03'
      WHEN 'read_holding' THEN '03'
      WHEN '1' THEN '01'
      WHEN '01' THEN '01'
      WHEN 'fc01' THEN '01'
      WHEN 'coil' THEN '01'
      WHEN '2' THEN '02'
      WHEN '02' THEN '02'
      WHEN 'fc02' THEN '02'
      WHEN 'discrete' THEN '02'
      WHEN '4' THEN '04'
      WHEN '04' THEN '04'
      WHEN 'fc04' THEN '04'
      WHEN 'input' THEN '04'
      WHEN 'read_input' THEN '04'
      ELSE UPPER(TRIM(COALESCE(function, '')))
    END AS function_norm
  FROM points
)
SELECT device_id, address_norm, function_norm, COUNT(*) AS cnt
FROM normalized
GROUP BY device_id, address_norm, function_norm
HAVING COUNT(*) > 1
ORDER BY cnt DESC, device_id, address_norm
LIMIT 200;
SQL
)

TARGET_DOWN_QUERY=$(cat <<'SQL'
SELECT device_id, address, COUNT(*) AS cnt
FROM points
GROUP BY device_id, address
HAVING COUNT(*) > 1
ORDER BY cnt DESC, device_id, address
LIMIT 200;
SQL
)

has_conflict=0

run_check() {
  local name="$1"
  local query="$2"
  local output
  output="$(run_query "$query")"
  if [[ -n "$output" ]]; then
    has_conflict=1
    echo "[FAIL] $name conflict detected:"
    echo "$output"
    return
  fi
  echo "[PASS] $name"
}

echo "[INFO] db=$DB_TYPE stage=$STAGE"
if [[ "$STAGE" == "up" ]]; then
  run_check "target-unique(device_id,address,function)" "$TARGET_UP_QUERY"
  run_check "normalized-alias(device_id,address_norm,function_norm)" "$NORMALIZED_ALIAS_QUERY"
else
  run_check "rollback-unique(device_id,address)" "$TARGET_DOWN_QUERY"
fi

if [[ "$has_conflict" -eq 1 ]]; then
  echo "[BLOCK] conflict found, migration should not proceed"
  exit 2
fi

echo "[OK] no conflict found"

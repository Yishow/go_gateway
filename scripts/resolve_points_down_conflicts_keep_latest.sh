#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Resolve down-migration conflicts on points unique(device_id,address) by keeping latest row.

Usage:
  scripts/resolve_points_down_conflicts_keep_latest.sh --db <postgres|sqlite> [options]

Options:
  --db TYPE             postgres or sqlite
  --dsn TEXT            PostgreSQL DSN (required for postgres)
  --sqlite-file PATH    SQLite DB file (required for sqlite)
  --execute             Apply deletion (default: preview only)
  --show-sql            Print SQL used
  -h, --help            Show help
USAGE
}

DB_TYPE=""
PG_DSN=""
SQLITE_FILE=""
EXECUTE=0
SHOW_SQL=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --db) DB_TYPE="$2"; shift 2 ;;
    --dsn) PG_DSN="$2"; shift 2 ;;
    --sqlite-file) SQLITE_FILE="$2"; shift 2 ;;
    --execute) EXECUTE=1; shift ;;
    --show-sql) SHOW_SQL=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "[ERROR] unknown arg: $1" >&2; usage; exit 1 ;;
  esac
done

if [[ "$DB_TYPE" != "postgres" && "$DB_TYPE" != "sqlite" ]]; then
  echo "[ERROR] --db must be postgres or sqlite" >&2
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

PG_CONFLICT_QUERY=$(cat <<'SQL'
SELECT device_id, address, COUNT(*) AS cnt
FROM points
GROUP BY device_id, address
HAVING COUNT(*) > 1
ORDER BY cnt DESC, device_id, address
LIMIT 200;
SQL
)

PG_DELETE_SQL=$(cat <<'SQL'
WITH ranked AS (
  SELECT
    id,
    ROW_NUMBER() OVER (
      PARTITION BY device_id, address
      ORDER BY
        COALESCE(updated_at, created_at) DESC,
        created_at DESC,
        id DESC
    ) AS rn
  FROM points
)
DELETE FROM points p
USING ranked r
WHERE p.id = r.id
  AND r.rn > 1;
SQL
)

SQLITE_CONFLICT_QUERY=$(cat <<'SQL'
SELECT device_id, address, COUNT(*) AS cnt
FROM points
GROUP BY device_id, address
HAVING COUNT(*) > 1
ORDER BY cnt DESC, device_id, address
LIMIT 200;
SQL
)

SQLITE_DELETE_SQL=$(cat <<'SQL'
DELETE FROM points
WHERE id IN (
  SELECT id FROM (
    SELECT
      id,
      ROW_NUMBER() OVER (
        PARTITION BY device_id, address
        ORDER BY
          COALESCE(updated_at, created_at) DESC,
          created_at DESC,
          id DESC
      ) AS rn
    FROM points
  )
  WHERE rn > 1
);
SQL
)

if [[ "$SHOW_SQL" -eq 1 ]]; then
  echo "[SQL] conflict query"
  if [[ "$DB_TYPE" == "postgres" ]]; then
    echo "$PG_CONFLICT_QUERY"
  else
    echo "$SQLITE_CONFLICT_QUERY"
  fi
  echo "[SQL] delete query"
  if [[ "$DB_TYPE" == "postgres" ]]; then
    echo "$PG_DELETE_SQL"
  else
    echo "$SQLITE_DELETE_SQL"
  fi
fi

if [[ "$DB_TYPE" == "postgres" ]]; then
  if ! command -v psql >/dev/null 2>&1; then
    echo "[ERROR] psql not found in PATH" >&2
    exit 1
  fi

  conflicts=$(psql "$PG_DSN" -v ON_ERROR_STOP=1 -At -F '|' -c "$PG_CONFLICT_QUERY")
  if [[ -z "$conflicts" ]]; then
    echo "[OK] no down conflicts found"
    exit 0
  fi

  echo "[INFO] conflicts detected:" 
  echo "$conflicts"

  if [[ "$EXECUTE" -eq 0 ]]; then
    echo "[PREVIEW] no changes applied. re-run with --execute to keep latest rows."
    exit 0
  fi

  psql "$PG_DSN" -v ON_ERROR_STOP=1 -c "$PG_DELETE_SQL"
  echo "[DONE] conflicts resolved by keeping latest rows"
  exit 0
fi

python3 - "$SQLITE_FILE" "$EXECUTE" <<'PY'
import sqlite3
import sys

db_path = sys.argv[1]
execute = sys.argv[2] == "1"

conflict_query = """
SELECT device_id, address, COUNT(*) AS cnt
FROM points
GROUP BY device_id, address
HAVING COUNT(*) > 1
ORDER BY cnt DESC, device_id, address
LIMIT 200;
"""

delete_sql = """
DELETE FROM points
WHERE id IN (
  SELECT id FROM (
    SELECT
      id,
      ROW_NUMBER() OVER (
        PARTITION BY device_id, address
        ORDER BY
          COALESCE(updated_at, created_at) DESC,
          created_at DESC,
          id DESC
      ) AS rn
    FROM points
  )
  WHERE rn > 1
);
"""

conn = sqlite3.connect(db_path)
try:
    cur = conn.cursor()
    rows = cur.execute(conflict_query).fetchall()
    if not rows:
        print("[OK] no down conflicts found")
        raise SystemExit(0)

    print("[INFO] conflicts detected:")
    for r in rows:
        print("|".join(str(v) for v in r))

    if not execute:
        print("[PREVIEW] no changes applied. re-run with --execute to keep latest rows.")
        raise SystemExit(0)

    cur.executescript(delete_sql)
    conn.commit()
    print("[DONE] conflicts resolved by keeping latest rows")
finally:
    conn.close()
PY

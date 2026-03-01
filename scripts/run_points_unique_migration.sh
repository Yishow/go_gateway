#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Run points unique-key migration with safety precheck.

Usage:
  scripts/run_points_unique_migration.sh --db <postgres|sqlite> --action <up|down> [options]

Options:
  --db TYPE                   Database type: postgres or sqlite
  --action ACTION             Migration action: up or down
  --dsn TEXT                  PostgreSQL DSN (required when --db postgres)
  --sqlite-file PATH          SQLite DB file path (required when --db sqlite)
  --execute                   Actually execute migration SQL
  --confirm-down              Required when --action down and --execute
  --down-conflict-strategy S  down conflict strategy: manual|keep-latest (default: keep-latest)
  --show-sql                  Show precheck SQL
  -h, --help                  Show help

Examples:
  # Print command templates only
  scripts/run_points_unique_migration.sh --db postgres --action up --dsn "$DATABASE_URL"

  # Execute rollback with keep-latest strategy
  scripts/run_points_unique_migration.sh \
    --db sqlite --action down --sqlite-file ./data/gateway.db \
    --execute --confirm-down --down-conflict-strategy keep-latest
USAGE
}

DB_TYPE=""
ACTION=""
PG_DSN=""
SQLITE_FILE=""
EXECUTE=0
CONFIRM_DOWN=0
SHOW_SQL=0
DOWN_CONFLICT_STRATEGY="keep-latest"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --db) DB_TYPE="$2"; shift 2 ;;
    --action) ACTION="$2"; shift 2 ;;
    --dsn) PG_DSN="$2"; shift 2 ;;
    --sqlite-file) SQLITE_FILE="$2"; shift 2 ;;
    --execute) EXECUTE=1; shift ;;
    --confirm-down) CONFIRM_DOWN=1; shift ;;
    --down-conflict-strategy) DOWN_CONFLICT_STRATEGY="$2"; shift 2 ;;
    --show-sql) SHOW_SQL=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "[ERROR] unknown arg: $1" >&2; usage; exit 1 ;;
  esac
done

if [[ "$DB_TYPE" != "postgres" && "$DB_TYPE" != "sqlite" ]]; then
  echo "[ERROR] --db must be postgres or sqlite" >&2
  exit 1
fi
if [[ "$ACTION" != "up" && "$ACTION" != "down" ]]; then
  echo "[ERROR] --action must be up or down" >&2
  exit 1
fi
if [[ "$DOWN_CONFLICT_STRATEGY" != "manual" && "$DOWN_CONFLICT_STRATEGY" != "keep-latest" ]]; then
  echo "[ERROR] --down-conflict-strategy must be manual or keep-latest" >&2
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

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

if [[ "$DB_TYPE" == "postgres" ]]; then
  MIGRATION_FILE="$REPO_ROOT/internal/datalink/schema/migrations/003_align_points_unique_function.${ACTION}.sql"
else
  MIGRATION_FILE="$REPO_ROOT/internal/datalink/schema/migrations/003_align_points_unique_function_sqlite.${ACTION}.sql"
fi

if [[ ! -f "$MIGRATION_FILE" ]]; then
  echo "[ERROR] migration file not found: $MIGRATION_FILE" >&2
  exit 1
fi

PRECHECK_CMD=("$REPO_ROOT/scripts/check_points_unique_conflicts.sh" "--db" "$DB_TYPE" "--stage" "$ACTION")
if [[ "$DB_TYPE" == "postgres" ]]; then
  PRECHECK_CMD+=("--dsn" "$PG_DSN")
else
  PRECHECK_CMD+=("--sqlite-file" "$SQLITE_FILE")
fi
if [[ "$SHOW_SQL" -eq 1 ]]; then
  PRECHECK_CMD+=("--show-sql")
fi

RESOLVER_CMD=("$REPO_ROOT/scripts/resolve_points_down_conflicts_keep_latest.sh" "--db" "$DB_TYPE")
if [[ "$DB_TYPE" == "postgres" ]]; then
  RESOLVER_CMD+=("--dsn" "$PG_DSN")
else
  RESOLVER_CMD+=("--sqlite-file" "$SQLITE_FILE")
fi
if [[ "$SHOW_SQL" -eq 1 ]]; then
  RESOLVER_CMD+=("--show-sql")
fi

run_precheck() {
  set +e
  "${PRECHECK_CMD[@]}"
  local rc=$?
  set -e
  return "$rc"
}

echo "[INFO] migration file: $MIGRATION_FILE"
echo "[INFO] precheck stage=$ACTION ..."
if run_precheck; then
  echo "[OK] precheck passed"
else
  rc=$?
  if [[ "$ACTION" == "down" && "$rc" -eq 2 && "$DOWN_CONFLICT_STRATEGY" == "keep-latest" ]]; then
    if [[ "$EXECUTE" -eq 1 ]]; then
      echo "[WARN] down conflicts found, apply keep-latest strategy"
      "${RESOLVER_CMD[@]}" --execute
      echo "[INFO] re-running precheck after conflict resolution"
      run_precheck
    else
      echo "[WARN] down conflicts found. strategy=keep-latest, but preview mode only"
      "${RESOLVER_CMD[@]}"
      echo "[BLOCK] re-run with --execute to apply conflict resolution then migration"
      exit 2
    fi
  else
    echo "[BLOCK] precheck failed (rc=$rc). strategy=$DOWN_CONFLICT_STRATEGY"
    exit "$rc"
  fi
fi

if [[ "$EXECUTE" -eq 0 ]]; then
  echo "[TEMPLATE] preview mode (no migration execution)."
  if [[ "$DB_TYPE" == "postgres" ]]; then
    echo "psql \"\$DATABASE_URL\" -v ON_ERROR_STOP=1 -f \"$MIGRATION_FILE\""
  else
    cat <<EOF
python3 - <<'PY'
import sqlite3
db = sqlite3.connect("$SQLITE_FILE")
try:
    db.executescript(open("$MIGRATION_FILE", "r", encoding="utf-8").read())
    db.commit()
finally:
    db.close()
PY
EOF
  fi
  exit 0
fi

if [[ "$ACTION" == "down" && "$CONFIRM_DOWN" -ne 1 ]]; then
  echo "[ERROR] rollback requires --confirm-down" >&2
  exit 1
fi

echo "[INFO] executing migration..."
if [[ "$DB_TYPE" == "postgres" ]]; then
  if ! command -v psql >/dev/null 2>&1; then
    echo "[ERROR] psql not found in PATH" >&2
    exit 1
  fi
  psql "$PG_DSN" -v ON_ERROR_STOP=1 -f "$MIGRATION_FILE"
else
  python3 - "$SQLITE_FILE" "$MIGRATION_FILE" <<'PY'
import sqlite3
import sys

db_path, migration_path = sys.argv[1], sys.argv[2]
conn = sqlite3.connect(db_path)
try:
    with open(migration_path, "r", encoding="utf-8") as fh:
        sql = fh.read()
    conn.executescript(sql)
    conn.commit()
finally:
    conn.close()
PY
fi

echo "[DONE] migration applied successfully"

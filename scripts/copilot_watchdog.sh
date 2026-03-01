#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
用法：
  scripts/copilot_watchdog.sh \
    --workdir /path/to/repo \
    --prompt-file /path/to/prompt.md \
    [--model gpt-5.3-codex] \
    [--idle-minutes 8] \
    [--max-restarts 2] \
    [--poll-seconds 15] \
    [--log-file /tmp/copilot-task.log] \
    [--dry-run] \
    [-- <extra copilot args>]

範例：
  scripts/copilot_watchdog.sh \
    --workdir /home/yishow/github/go_gateway \
    --prompt-file /home/yishow/github/go_gateway/docs/templates/task.prompt.md \
    --idle-minutes 8 \
    --max-restarts 2 \
    --log-file /tmp/go_gateway_copilot.log
USAGE
}

WORKDIR="$(pwd)"
MODEL="gpt-5.3-codex"
PROMPT_FILE=""
PROMPT_TEXT=""
IDLE_MINUTES=8
MAX_RESTARTS=2
POLL_SECONDS=15
LOG_FILE="/tmp/copilot_watchdog_$(date +%Y%m%d_%H%M%S).log"
DRY_RUN=0

EXTRA_ARGS=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --workdir)
      WORKDIR="$2"; shift 2 ;;
    --model)
      MODEL="$2"; shift 2 ;;
    --prompt-file)
      PROMPT_FILE="$2"; shift 2 ;;
    --prompt)
      PROMPT_TEXT="$2"; shift 2 ;;
    --idle-minutes)
      IDLE_MINUTES="$2"; shift 2 ;;
    --max-restarts)
      MAX_RESTARTS="$2"; shift 2 ;;
    --poll-seconds)
      POLL_SECONDS="$2"; shift 2 ;;
    --log-file)
      LOG_FILE="$2"; shift 2 ;;
    --dry-run)
      DRY_RUN=1; shift ;;
    --help|-h)
      usage; exit 0 ;;
    --)
      shift
      EXTRA_ARGS=("$@")
      break ;;
    *)
      echo "[ERROR] 未知參數: $1" >&2
      usage
      exit 1 ;;
  esac
done

if [[ -n "$PROMPT_FILE" ]]; then
  if [[ ! -f "$PROMPT_FILE" ]]; then
    echo "[ERROR] prompt file 不存在: $PROMPT_FILE" >&2
    exit 1
  fi
  PROMPT_TEXT="$(cat "$PROMPT_FILE")"
fi

if [[ -z "$PROMPT_TEXT" ]]; then
  echo "[ERROR] 需要 --prompt-file 或 --prompt" >&2
  exit 1
fi

if ! command -v copilot >/dev/null 2>&1; then
  echo "[ERROR] 找不到 copilot CLI，請先安裝。" >&2
  exit 1
fi

if [[ ! -d "$WORKDIR" ]]; then
  echo "[ERROR] workdir 不存在: $WORKDIR" >&2
  exit 1
fi

mkdir -p "$(dirname "$LOG_FILE")"
touch "$LOG_FILE"

run_once() {
  local attempt="$1"
  local idle_limit=$((IDLE_MINUTES * 60))

  {
    echo "================================================================"
    echo "[INFO] attempt=$attempt started_at=$(date -Iseconds)"
    echo "[INFO] workdir=$WORKDIR"
    echo "[INFO] model=$MODEL idle_minutes=$IDLE_MINUTES poll_seconds=$POLL_SECONDS"
    echo "================================================================"
  } >>"$LOG_FILE"

  local -a cmd=(
    copilot
    --model "$MODEL"
    --allow-all
    --no-ask-user
    --stream off
    -p "$PROMPT_TEXT"
  )

  if [[ ${#EXTRA_ARGS[@]} -gt 0 ]]; then
    cmd+=("${EXTRA_ARGS[@]}")
  fi

  if [[ "$DRY_RUN" -eq 1 ]]; then
    echo "[DRY-RUN] ${cmd[*]}" | tee -a "$LOG_FILE"
    return 0
  fi

  (
    cd "$WORKDIR"
    "${cmd[@]}"
  ) >>"$LOG_FILE" 2>&1 &
  local pid=$!
  local last_size=0
  local last_change
  last_change="$(date +%s)"

  while kill -0 "$pid" >/dev/null 2>&1; do
    sleep "$POLL_SECONDS"

    local size
    size=$(wc -c <"$LOG_FILE" | tr -d ' ')
    if (( size > last_size )); then
      last_size=$size
      last_change="$(date +%s)"
    fi

    local now
    now="$(date +%s)"
    if (( now - last_change > idle_limit )); then
      echo "[WARN] no new log for ${IDLE_MINUTES}m, restarting process pid=$pid" | tee -a "$LOG_FILE"
      kill "$pid" >/dev/null 2>&1 || true
      sleep 2
      kill -9 "$pid" >/dev/null 2>&1 || true
      wait "$pid" >/dev/null 2>&1 || true
      return 124
    fi
  done

  if wait "$pid"; then
    echo "[INFO] attempt=$attempt finished successfully at $(date -Iseconds)" | tee -a "$LOG_FILE"
    return 0
  else
    local rc=$?
    echo "[ERROR] attempt=$attempt exited with code=$rc at $(date -Iseconds)" | tee -a "$LOG_FILE"
    return "$rc"
  fi
}

attempt=1
max_attempts=$((MAX_RESTARTS + 1))

while (( attempt <= max_attempts )); do
  if run_once "$attempt"; then
    echo "[DONE] success. log=$LOG_FILE"
    exit 0
  fi

  rc=$?
  if [[ "$rc" -eq 124 && "$attempt" -lt "$max_attempts" ]]; then
    echo "[INFO] watchdog restart triggered (attempt ${attempt}/${max_attempts})" | tee -a "$LOG_FILE"
    attempt=$((attempt + 1))
    continue
  fi

  echo "[FAIL] stopped at attempt ${attempt}/${max_attempts}, rc=$rc, log=$LOG_FILE" >&2
  exit "$rc"
done

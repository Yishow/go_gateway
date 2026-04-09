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
    [--poll-seconds 2] \
    [--max-tree-mb 8192] \
    [--max-proc-mb 4096] \
    [--kill-grace-seconds 1] \
    [--log-file /tmp/copilot-task.log] \
    [--dry-run] \
    [-- <extra copilot args>]

範例：
  scripts/copilot_watchdog.sh \
    --workdir /home/yishow/github/go_gateway \
    --prompt-file /home/yishow/github/go_gateway/docs/templates/task.prompt.md \
    --idle-minutes 8 \
    --max-tree-mb 8192 \
    --max-proc-mb 4096 \
    --max-restarts 2 \
    --log-file /tmp/go_gateway_copilot.log

說明：
  - 只監控這次 copilot session 底下的子孫程序。
  - 只會處理 chrome/chromium/node 類程序，不會 kill copilot 主程序。
  - --max-tree-mb 或 --max-proc-mb 設成 0 可停用對應門檻。
USAGE
}

WORKDIR="$(pwd)"
MODEL="gpt-5.3-codex"
PROMPT_FILE=""
PROMPT_TEXT=""
IDLE_MINUTES=8
MAX_RESTARTS=2
POLL_SECONDS=2
MAX_TREE_MB=8192
MAX_PROC_MB=4096
KILL_GRACE_SECONDS=1
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
    --max-tree-mb)
      MAX_TREE_MB="$2"; shift 2 ;;
    --max-proc-mb)
      MAX_PROC_MB="$2"; shift 2 ;;
    --kill-grace-seconds)
      KILL_GRACE_SECONDS="$2"; shift 2 ;;
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

mb_to_kb() {
  local value="${1:-0}"
  if [[ "$value" =~ ^[0-9]+$ ]]; then
    printf '%s\n' "$((value * 1024))"
  else
    printf '0\n'
  fi
}

rss_kb_to_mb() {
  local rss_kb="${1:-0}"
  if [[ "$rss_kb" -le 0 ]]; then
    printf '0\n'
  else
    printf '%s\n' "$(((rss_kb + 1023) / 1024))"
  fi
}

is_target_command() {
  local args="$1"
  local lower
  lower="$(printf '%s' "$args" | tr '[:upper:]' '[:lower:]')"
  [[ "$lower" == *"chrome"* || "$lower" == *"chromium"* || "$lower" == *"/node"* || "$lower" == node* || "$lower" == *" node "* ]]
}

get_descendants() {
  local root_pid="$1"
  local queue=("$root_pid")
  local idx=0
  local child

  while (( idx < ${#queue[@]} )); do
    local parent="${queue[$idx]}"
    idx=$((idx + 1))

    while IFS= read -r child; do
      [[ -z "$child" ]] && continue
      queue+=("$child")
      printf '%s\n' "$child"
    done < <(pgrep -P "$parent" 2>/dev/null || true)
  done
}

snapshot_target_processes() {
  local root_pid="$1"
  local descendants=()
  local child
  local pid_csv
  local line pid rss args

  while IFS= read -r child; do
    [[ -n "$child" ]] && descendants+=("$child")
  done < <(get_descendants "$root_pid")

  if (( ${#descendants[@]} == 0 )); then
    return 0
  fi

  pid_csv="$(IFS=,; printf '%s' "${descendants[*]}")"

  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    pid=""
    rss=""
    args=""
    read -r pid rss args <<<"$line"
    [[ -z "$pid" || -z "$rss" ]] && continue
    if is_target_command "$args"; then
      printf '%s\t%s\t%s\n' "$pid" "$rss" "$args"
    fi
  done < <(ps -o pid=,rss=,args= -p "$pid_csv" 2>/dev/null || true)
}

terminate_target_tree() {
  local target_pid="$1"
  local root_pid="$2"
  local reason="$3"
  local rss_kb="$4"
  local cmd="$5"
  local all_pids=()
  local child
  local i
  local survivors=()
  local rss_mb

  if [[ "$target_pid" == "$root_pid" ]]; then
    return 0
  fi

  if ! kill -0 "$target_pid" >/dev/null 2>&1; then
    return 0
  fi

  while IFS= read -r child; do
    [[ -n "$child" && "$child" != "$root_pid" ]] && all_pids+=("$child")
  done < <(get_descendants "$target_pid")
  all_pids+=("$target_pid")
  rss_mb="$(rss_kb_to_mb "$rss_kb")"

  {
    echo "[WARN] memory guard: reason=$reason pid=$target_pid rss_mb=$rss_mb grace_seconds=$KILL_GRACE_SECONDS"
    echo "[WARN] cmd=$cmd"
  } | tee -a "$LOG_FILE"

  for ((i=${#all_pids[@]}-1; i>=0; i--)); do
    kill "${all_pids[$i]}" >/dev/null 2>&1 || true
  done

  sleep "$KILL_GRACE_SECONDS"

  for pid in "${all_pids[@]}"; do
    if kill -0 "$pid" >/dev/null 2>&1; then
      survivors+=("$pid")
    fi
  done

  if (( ${#survivors[@]} > 0 )); then
    echo "[WARN] escalating to SIGKILL for pids=${survivors[*]}" | tee -a "$LOG_FILE"
    for pid in "${survivors[@]}"; do
      kill -9 "$pid" >/dev/null 2>&1 || true
    done
  fi
}

enforce_memory_limits() {
  local root_pid="$1"
  local max_tree_kb
  local max_proc_kb
  local snapshot=()
  local selected=()
  local selected_ids=" "
  local total_kb=0
  local line pid rss cmd
  local sorted_lines
  local projected_kb

  max_tree_kb="$(mb_to_kb "$MAX_TREE_MB")"
  max_proc_kb="$(mb_to_kb "$MAX_PROC_MB")"

  if (( max_tree_kb <= 0 && max_proc_kb <= 0 )); then
    return 0
  fi

  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    snapshot+=("$line")
    IFS=$'\t' read -r pid rss cmd <<<"$line"
    total_kb=$((total_kb + rss))
  done < <(snapshot_target_processes "$root_pid")

  if (( ${#snapshot[@]} == 0 )); then
    return 0
  fi

  sorted_lines="$(printf '%s\n' "${snapshot[@]}" | sort -t$'\t' -k2,2nr)"

  if (( max_proc_kb > 0 )); then
    while IFS=$'\t' read -r pid rss cmd; do
      [[ -z "$pid" ]] && continue
      if (( rss > max_proc_kb )); then
        selected+=("$pid"$'\t'"$rss"$'\t'"$cmd")
        selected_ids+=" $pid "
      fi
    done <<<"$sorted_lines"
  fi

  if (( max_tree_kb > 0 && total_kb > max_tree_kb )); then
    projected_kb="$total_kb"
    while IFS=$'\t' read -r pid rss cmd; do
      [[ -z "$pid" ]] && continue
      if [[ "$selected_ids" == *" $pid "* ]]; then
        projected_kb=$((projected_kb - rss))
        continue
      fi
      selected+=("$pid"$'\t'"$rss"$'\t'"$cmd")
      selected_ids+=" $pid "
      projected_kb=$((projected_kb - rss))
      if (( projected_kb <= max_tree_kb )); then
        break
      fi
    done <<<"$sorted_lines"
  fi

  if (( ${#selected[@]} == 0 )); then
    return 0
  fi

  echo "[WARN] memory guard: target_total_mb=$(rss_kb_to_mb "$total_kb") tree_limit_mb=$MAX_TREE_MB proc_limit_mb=$MAX_PROC_MB" | tee -a "$LOG_FILE"

  for line in "${selected[@]}"; do
    IFS=$'\t' read -r pid rss cmd <<<"$line"
    terminate_target_tree "$pid" "$root_pid" "rss-threshold" "$rss" "$cmd"
  done
}

run_once() {
  local attempt="$1"
  local idle_limit=$((IDLE_MINUTES * 60))

  {
    echo "================================================================"
    echo "[INFO] attempt=$attempt started_at=$(date -Iseconds)"
    echo "[INFO] workdir=$WORKDIR"
    echo "[INFO] model=$MODEL idle_minutes=$IDLE_MINUTES poll_seconds=$POLL_SECONDS"
    echo "[INFO] memory_guard tree_mb=$MAX_TREE_MB proc_mb=$MAX_PROC_MB kill_grace_seconds=$KILL_GRACE_SECONDS"
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

    enforce_memory_limits "$pid"

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

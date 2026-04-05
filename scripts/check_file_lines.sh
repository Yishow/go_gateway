#!/usr/bin/env bash
set -euo pipefail

WARN_LIMIT="${WARN_LIMIT:-300}"
HARD_LIMIT="${HARD_LIMIT:-500}"
BASE_SHA=""
CHECK_ALL=0

usage() {
  cat <<'EOF'
Usage:
  scripts/check_file_lines.sh [--base <sha>] [--warn <n>] [--hard <n>] [--all]

Policy:
  - > WARN (default 300): warning
  - > HARD (default 500): fail

Legacy guard:
  If a file was already > HARD at base commit, it is allowed only when line count
  does not increase (same or fewer lines).
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --base)
      BASE_SHA="${2:-}"
      shift 2
      ;;
    --warn)
      WARN_LIMIT="${2:-}"
      shift 2
      ;;
    --hard)
      HARD_LIMIT="${2:-}"
      shift 2
      ;;
    --all)
      CHECK_ALL=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 2
      ;;
  esac
done

if ! [[ "$WARN_LIMIT" =~ ^[0-9]+$ ]] || ! [[ "$HARD_LIMIT" =~ ^[0-9]+$ ]]; then
  echo "[line-limit] WARN/HARD must be integers." >&2
  exit 2
fi

if (( HARD_LIMIT <= WARN_LIMIT )); then
  echo "[line-limit] HARD limit must be greater than WARN limit." >&2
  exit 2
fi

IGNORE_FILE=".line-limit-ignore"
ignore_patterns=()
if [[ -f "$IGNORE_FILE" ]]; then
  while IFS= read -r line; do
    line="${line%%#*}"
    line="${line#"${line%%[![:space:]]*}"}"
    line="${line%"${line##*[![:space:]]}"}"
    [[ -z "$line" ]] && continue
    ignore_patterns+=("$line")
  done < "$IGNORE_FILE"
fi

is_ignored() {
  local path="$1"
  for pat in "${ignore_patterns[@]:-}"; do
    if [[ "$path" == $pat ]]; then
      return 0
    fi
  done
  return 1
}

is_text_candidate() {
  local path="$1"
  local base
  base="$(basename "$path")"
  case "$base" in
    Makefile|Dockerfile|AGENTS.md|CLAUDE.md|README.md)
      return 0
      ;;
  esac

  case "$path" in
    *.go|*.ts|*.tsx|*.js|*.jsx|*.css|*.scss|*.sh|*.bash|*.zsh|*.ps1|*.md|*.txt|*.yml|*.yaml|*.json|*.toml|*.ini|*.env|*.sql)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

base_is_valid=0
if [[ -n "$BASE_SHA" ]] && git cat-file -e "${BASE_SHA}^{commit}" 2>/dev/null; then
  base_is_valid=1
fi

comparison_base=""
if (( base_is_valid == 1 )); then
  comparison_base="$BASE_SHA"
elif git rev-parse --verify HEAD >/dev/null 2>&1; then
  comparison_base="HEAD"
fi

if (( CHECK_ALL == 1 )); then
  mapfile -t changed_files < <(git ls-files)
else
  if (( base_is_valid == 1 )); then
    mapfile -t changed_files < <(git diff --name-only --diff-filter=ACMR "${BASE_SHA}...HEAD")
  else
    tmp_paths="$(mktemp)"
    git diff --name-only --diff-filter=ACMR --cached >> "$tmp_paths" || true
    git diff --name-only --diff-filter=ACMR >> "$tmp_paths" || true
    git ls-files --others --exclude-standard >> "$tmp_paths" || true
    mapfile -t changed_files < <(sort -u "$tmp_paths")
    rm -f "$tmp_paths"
    if [[ ${#changed_files[@]} -eq 0 ]] && git rev-parse --verify HEAD~1 >/dev/null 2>&1; then
      mapfile -t changed_files < <(git diff --name-only --diff-filter=ACMR HEAD~1...HEAD)
    fi
  fi
fi

if [[ ${#changed_files[@]} -eq 0 ]]; then
  echo "[line-limit] No files to check."
  exit 0
fi

warn_items=()
legacy_items=()
fail_items=()
checked_count=0

for f in "${changed_files[@]}"; do
  [[ -f "$f" ]] || continue
  is_ignored "$f" && continue
  is_text_candidate "$f" || continue

  checked_count=$((checked_count + 1))
  new_lines="$(wc -l < "$f" | tr -d '[:space:]')"

  old_exists=0
  old_lines=0
  if [[ -n "$comparison_base" ]] && git cat-file -e "${comparison_base}:${f}" 2>/dev/null; then
    old_exists=1
    old_lines="$(git show "${comparison_base}:${f}" | wc -l | tr -d '[:space:]')"
  fi

  if (( new_lines > HARD_LIMIT )); then
    if (( old_exists == 1 )) && (( old_lines > HARD_LIMIT )) && (( new_lines <= old_lines )); then
      legacy_items+=("${new_lines} ${f} (legacy ${old_lines} -> ${new_lines})")
      continue
    fi
    fail_items+=("${new_lines} ${f}")
    continue
  fi

  if (( new_lines > WARN_LIMIT )); then
    warn_items+=("${new_lines} ${f}")
  fi
done

echo "[line-limit] checked files: ${checked_count}"
echo "[line-limit] limits: warn>${WARN_LIMIT}, fail>${HARD_LIMIT}"

if [[ ${#warn_items[@]} -gt 0 ]]; then
  echo
  echo "[line-limit] WARNING files (> ${WARN_LIMIT} lines):"
  printf '  %s\n' "${warn_items[@]}"
fi

if [[ ${#legacy_items[@]} -gt 0 ]]; then
  echo
  echo "[line-limit] LEGACY allowed (> ${HARD_LIMIT}, not increased):"
  printf '  %s\n' "${legacy_items[@]}"
fi

if [[ ${#fail_items[@]} -gt 0 ]]; then
  echo
  echo "[line-limit] ERROR files (> ${HARD_LIMIT} lines):"
  printf '  %s\n' "${fail_items[@]}"
  echo
  echo "[line-limit] Failed. Split files or reduce line count."
  exit 1
fi

echo
echo "[line-limit] Passed."

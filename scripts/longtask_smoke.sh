#!/usr/bin/env bash
set -euo pipefail

CTRL="/home/yishow/.openclaw/agents/code-pro/agent/workspace/scripts/longtask_progress_controller.sh"
RECON="/home/yishow/.openclaw/agents/code-pro/agent/workspace/scripts/longtask_reconcile.sh"

if [[ ! -x "$CTRL" ]]; then
  echo "[ERR] 找不到 controller: $CTRL"
  exit 1
fi

SID="smoke-$(date +%s)"
TASK="go-gateway-longtask-smoke"

echo "[SMOKE] start sid=$SID"
"$CTRL" start --name "$TASK" --session-id "$SID" --interval-min 5 --summary-min 30

echo "[SMOKE] status after start"
"$CTRL" status --session-id "$SID"

echo "[SMOKE] finish sid=$SID"
"$CTRL" finish --session-id "$SID"

echo "[SMOKE] reconcile --all"
"$RECON" --all

echo "[SMOKE] status after finish"
"$CTRL" status --session-id "$SID"

echo "[OK] longtask smoke passed"

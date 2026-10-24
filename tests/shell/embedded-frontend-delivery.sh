#!/usr/bin/env bash
set -euo pipefail

# Both launchers now delegate embedded asset delivery to this shared builder.
# Verify current assets, stale removal, preserved placeholder, all targets,
# bootstrap-only preparation, and fail-closed installation/build/copy behavior.
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
exec node --test "$REPO_ROOT/scripts/tests/build-single.test.mjs"

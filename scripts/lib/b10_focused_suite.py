"""Focused backend transaction evidence for the B10 acceptance harness."""

from __future__ import annotations

import hashlib
import json
import os
import subprocess
from pathlib import Path
from typing import Any


FOCUSED_MODBUSSHARE_TEST_REGEX = (
    "Test(RegisterGeometry_(StrideValidation)|"
    "Reconciler_(AtomicCommitSuccess|StaleRevisionCASConflict|"
    "RollbackUncertaintyReportsDirtyUnknown|"
    "NilRevisionStoreFailsClosedBeforeMutation|"
    "EmptyExpectedRevisionFailsClosedBeforeMutation|"
    "ForeignWorkspaceProjectionIsPreservedAndExcluded|"
    "RevisionStatusReadsDurableDirtyMarker|"
    "RevisionStatusReadsDirtyMarkerAfterStoreRestart|"
    "SQLProjectionFailureRestoresOldDesiredAndDirtyUnknownOnRollbackFailure)|"
    "SQLWorkspaceRevisionStore_(PersistsAndCASChecksAcrossInstances|"
    "CommitsDesiredMappingsWithRevision)|"
    "Service_(HydrationFailureStopsRunningListener|"
    "StatusForWorkspaceScopesMappingsAndSharesLifecycleState|"
    "ApplySettings_BindConflictReportsTypedFailureWithoutFallbackTo5020)|"
    "Reconciler_StagedTransactionIsRejectedAfterGlobalDisable|"
    "Reconciler_ExactReplayRejectsExternalWorkspaceRevisionDrift|"
    "Service_StartCASChecksRevisionBeforeSamePortIdempotency)"
)

FOCUSED_HANDLER_TEST_REGEX = (
    "Test(ModbusShareHandler_(StatusExposesSafeBindFailureDiagnostic|"
    "StatusShowsDurableDirtyRecoveryAfterRestart)|"
    "StudioV2WorkspaceBootstrapSharesDurableRevisionAndDirtyRecovery|"
    "SourceRuleHandler_ApplyTags_RejectsRevisionMismatch)"
)


def run_focused_modbusshare_suite(root: Path, cache_dir: Path) -> dict[str, Any]:
    """Run deterministic durable/runtime transaction evidence outside the EXE."""
    test_env = os.environ.copy()
    test_env["GOCACHE"] = str(cache_dir)
    commands = [
        ["go", "test", "./internal/datalink/modbusshare", "-run", FOCUSED_MODBUSSHARE_TEST_REGEX, "-count=1"],
        ["go", "test", "./internal/api/handlers", "-run", FOCUSED_HANDLER_TEST_REGEX, "-count=1"],
    ]
    outputs = []
    exit_code = 0
    for command in commands:
        completed = subprocess.run(command, cwd=root, env=test_env, capture_output=True, text=True, check=False)
        outputs.append(completed.stdout + completed.stderr)
        exit_code = exit_code or completed.returncode
    output = "\n".join(outputs)[-12000:]
    return {
        "command": " && ".join(" ".join(command) for command in commands),
        "test_names": [
            "TestRegisterGeometry_StrideValidation",
            "TestReconciler_AtomicCommitSuccess",
            "TestReconciler_StaleRevisionCASConflict",
            "TestReconciler_RollbackUncertaintyReportsDirtyUnknown",
            "TestReconciler_NilRevisionStoreFailsClosedBeforeMutation",
            "TestReconciler_EmptyExpectedRevisionFailsClosedBeforeMutation",
            "TestReconciler_ForeignWorkspaceProjectionIsPreservedAndExcluded",
            "TestReconciler_RevisionStatusReadsDurableDirtyMarker",
            "TestReconciler_RevisionStatusReadsDirtyMarkerAfterStoreRestart",
            "TestSQLWorkspaceRevisionStore_PersistsAndCASChecksAcrossInstances",
            "TestSQLWorkspaceRevisionStore_CommitsDesiredMappingsWithRevision",
            "TestReconciler_SQLProjectionFailureRestoresOldDesiredAndDirtyUnknownOnRollbackFailure",
            "TestService_HydrationFailureStopsRunningListener",
            "TestService_StatusForWorkspaceScopesMappingsAndSharesLifecycleState",
            "TestService_ApplySettings_BindConflictReportsTypedFailureWithoutFallbackTo5020",
            "TestReconciler_StagedTransactionIsRejectedAfterGlobalDisable",
            "TestReconciler_ExactReplayRejectsExternalWorkspaceRevisionDrift",
            "TestService_StartCASChecksRevisionBeforeSamePortIdempotency",
            "TestModbusShareHandler_StatusExposesSafeBindFailureDiagnostic",
            "TestModbusShareHandler_StatusShowsDurableDirtyRecoveryAfterRestart",
            "TestStudioV2WorkspaceBootstrapSharesDurableRevisionAndDirtyRecovery",
            "TestSourceRuleHandler_ApplyTags_RejectsRevisionMismatch",
        ],
        "exit_code": exit_code,
        "output_hash": _hash(output),
    }


def _hash(value: Any) -> str:
    encoded = json.dumps(value, sort_keys=True, separators=(",", ":"), ensure_ascii=False).encode()
    return hashlib.sha256(encoded).hexdigest()

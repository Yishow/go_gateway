"""HTTP-only B10 negative acceptance steps.

The helpers deliberately use the public HTTP surface.  They do not open the
gateway database or import internal Go test seams.
"""

from __future__ import annotations

import hashlib
import json
from pathlib import Path
from typing import Any

from lib.b10_acceptance_helpers import (
    HTTP,
    SourceModbus,
    candidate_scope_path,
    data,
    free_port,
    read_holding,
    start_gateway,
    stop_gateway,
    wait_gateway,
)
from lib.b10_focused_suite import run_focused_modbusshare_suite
from lib.b10_projection_assertions import PROJECTION_STABILITY_FIELDS, assert_projection_unchanged


def _hash(value: Any) -> str:
    encoded = json.dumps(value, sort_keys=True, separators=(",", ":"), ensure_ascii=False).encode()
    return hashlib.sha256(encoded).hexdigest()


def _response_summary(method: str, path: str, status: int, payload: dict[str, Any]) -> dict[str, Any]:
    error = payload.get("error") if isinstance(payload.get("error"), dict) else {}
    body_data = payload.get("data")
    counts: dict[str, int] = {}
    if isinstance(body_data, dict):
        for key in ("mappings", "desired_mappings", "removed_spans", "invalidated_spans", "diagnostics", "results"):
            value = body_data.get(key)
            if isinstance(value, list):
                counts[key] = len(value)
    elif isinstance(body_data, list):
        counts["data"] = len(body_data)
    return {
        "request": {"method": method, "path": path},
        "status": status,
        "code": error.get("code"),
        "success": payload.get("success"),
        "counts": counts,
        "body_hash": _hash(payload),
    }


def _call(http: HTTP, method: str, path: str, body: dict[str, Any] | None, expected: set[int]) -> tuple[int, dict[str, Any], dict[str, Any]]:
    status, payload = http.call(method, path, body, expected=expected)
    return status, payload, _response_summary(method, path, status, payload)


def _settings(http: HTTP) -> dict[str, Any]:
    response = http.call("GET", "/api/v1/datalink/settings", expected={200})[1]
    item = next(item for item in response.get("data", []) if item.get("key") == "modbus_share")
    return dict(item["value"])


def _projection_state(http: HTTP, share_port: int, read_value: bool = True) -> dict[str, Any]:
    status = data(http.call("GET", "/api/v1/datalink/modbus-share/status", expected={200})[1])
    boot = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
    snapshot: dict[str, Any] = {
        "mapping_count": status.get("mapping_count"),
        "workspace_revision": boot.get("workspace_revision"),
        "settings_revision": status.get("settings_revision"),
        "lifecycle_state": status.get("lifecycle_state"),
        "status_hash": _hash(status),
    }
    if read_value:
        try:
            snapshot["register_value"] = read_holding("127.0.0.1", share_port)[0]
        except OSError:
            snapshot["register_value"] = None
    if status.get("lifecycle_state") == "running":
        mappings = http.call("GET", "/api/v1/datalink/modbus-share/mappings", expected={200})[1].get("data", [])
        snapshot["mapping_count"] = len(mappings)
        snapshot["mapping_hash"] = _hash(mappings)
    return snapshot


def _step(result: dict[str, Any], name: str, requests: list[dict[str, Any]], responses: list[dict[str, Any]], before: dict[str, Any] | None = None, after: dict[str, Any] | None = None, **extra: Any) -> None:
    entry: dict[str, Any] = {"name": name, "status": "pass", "requests": requests, "responses": responses}
    if before is not None:
        entry["before"] = before
    if after is not None:
        entry["after"] = after
    entry.update(extra)
    result.setdefault("steps", []).append(entry)


def run_pre_collision_negatives(
    http: HTTP,
    binary: Path,
    temp_root: Path,
    share_port: int,
    source: SourceModbus,
    result: dict[str, Any],
    command_log: list[str],
) -> None:
    """Exercise disabled, double-submit, stale revision, and bind conflict."""
    before_disabled = _projection_state(http, share_port)
    settings = _settings(http)
    disabled = dict(settings)
    disabled["enabled"] = False
    _, disabled_payload, disabled_response = _call(
        http,
        "PUT",
        "/api/v1/datalink/settings/modbus_share",
        {"value": disabled, "expected_settings_revision": settings["settings_revision"]},
        {200},
    )
    disabled_value = data(disabled_payload).get("value", disabled)
    boot = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
    reconcile_body = {
        "workspace_id": boot["workspace_id"],
        "expected_workspace_revision": boot["workspace_revision"],
        "expected_settings_revision": settings["settings_revision"],
        "readiness_token": boot["readiness_token"],
        "canonical_plan_signature": boot.get("canonical_plan", {}).get("signature", ""),
    }
    disabled_responses = [disabled_response]
    disabled_requests = [{"method": "PUT", "path": "/api/v1/datalink/settings/modbus_share"}]
    for method, path, body, expected in (
        ("GET", "/api/v1/datalink/modbus-share/mappings", None, {422}),
        ("POST", "/api/v1/datalink/modbus-share/reconcile", reconcile_body, {422}),
        ("PUT", "/api/v1/datalink/modbus-share/mappings/b10-disabled-probe", {"register": 0}, {422}),
        ("POST", "/api/v1/datalink/studio-v2/workspace/activate", {
            "workspace_revision": boot["workspace_revision"],
            "settings_revision": disabled_value.get("settings_revision", ""),
            "readiness_token": boot["readiness_token"],
            "pending_saves": 0,
        }, {422}),
    ):
        status, payload, summary = _call(http, method, path, body, expected)
        code = summary.get("code")
        if method == "PUT" and path.endswith("b10-disabled-probe") and code == "modbus_share_projection_required":
            summary["gate_interpretation"] = "projection_required (direct mutation seam remains closed)"
        elif code != "modbus_share_disabled":
            raise RuntimeError(f"disabled {method} {path} returned unexpected typed code: {payload}")
        disabled_requests.append({"method": method, "path": path})
        disabled_responses.append(summary)
    try:
        read_holding("127.0.0.1", share_port)
    except OSError:
        listener_disabled = True
    else:
        listener_disabled = False
    if not listener_disabled:
        raise RuntimeError("disabled mode still accepted a Modbus connection")
    enabled = dict(disabled_value)
    enabled.update({"enabled": True, "bind_address": "127.0.0.1", "port": share_port, "slave_id": 1, "capacity_registers": 64})
    _, enabled_payload, enabled_response = _call(
        http,
        "PUT",
        "/api/v1/datalink/settings/modbus_share",
        {"value": enabled, "expected_settings_revision": disabled_value["settings_revision"]},
        {200},
    )
    after_disabled = _projection_state(http, share_port)
    if before_disabled.get("mapping_hash") != after_disabled.get("mapping_hash") or before_disabled.get("register_value") != after_disabled.get("register_value"):
        raise RuntimeError(f"disabled requests changed the projection: before={before_disabled} after={after_disabled}")
    _step(result, "disabled-list-reconcile-activation-no-mutation", disabled_requests + [{"method": "PUT", "path": "/api/v1/datalink/settings/modbus_share"}], disabled_responses + [enabled_response], before_disabled, after_disabled, listener_disabled=listener_disabled, projection_hash=after_disabled.get("mapping_hash"))

    boot = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
    body = {
        "workspace_id": boot["workspace_id"],
        "expected_workspace_revision": boot["workspace_revision"],
        "expected_settings_revision": boot.get("canonical_plan", {}).get("settings_revision", ""),
        "readiness_token": boot["readiness_token"],
        "canonical_plan_signature": boot.get("canonical_plan", {}).get("signature", ""),
    }
    before_double = _projection_state(http, share_port)
    first_status, first_payload, first_response = _call(http, "POST", "/api/v1/datalink/modbus-share/reconcile", body, {200})
    second_status, second_payload, second_response = _call(http, "POST", "/api/v1/datalink/modbus-share/reconcile", body, {200})
    first_data, second_data = data(first_payload), data(second_payload)
    if first_data.get("outcome") != second_data.get("outcome") or first_data.get("new_workspace_revision") != second_data.get("new_workspace_revision"):
        raise RuntimeError(f"exact double-submit diverged: first={first_payload} second={second_payload}")
    after_double = _projection_state(http, share_port)
    _step(result, "exact-double-submit-idempotent", [{"method": "POST", "path": "/api/v1/datalink/modbus-share/reconcile"}] * 2, [first_response, second_response], before_double, after_double, outcomes=[first_data.get("outcome"), second_data.get("outcome")], revisions=[first_data.get("new_workspace_revision"), second_data.get("new_workspace_revision")], response_hashes=[first_response["body_hash"], second_response["body_hash"]])

    stale_settings = _settings(http)
    changed_settings = dict(stale_settings)
    changed_settings["capacity_registers"] = 63
    _call(http, "PUT", "/api/v1/datalink/settings/modbus_share", {"value": changed_settings, "expected_settings_revision": stale_settings["settings_revision"]}, {200})
    stale_baseline = _projection_state(http, share_port)
    current_boot = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
    stale_body = {
        "workspace_id": current_boot["workspace_id"],
        "expected_workspace_revision": current_boot["workspace_revision"],
        "expected_settings_revision": stale_settings["settings_revision"],
        "readiness_token": current_boot["readiness_token"],
        "desired_mappings": current_boot.get("canonical_desired_mappings", []),
    }
    stale_status, stale_payload, stale_response = _call(http, "POST", "/api/v1/datalink/modbus-share/reconcile", stale_body, {409})
    if stale_response.get("code") != "modbus_share_revision_conflict":
        raise RuntimeError(f"stale revision returned unexpected code: {stale_payload}")
    stale_after = _projection_state(http, share_port)
    assert_projection_unchanged(stale_baseline, stale_after, "stale reconcile")
    _step(result, "different-desired-stale-revision-conflict", [{"method": "PUT", "path": "/api/v1/datalink/settings/modbus_share"}, {"method": "POST", "path": "/api/v1/datalink/modbus-share/reconcile"}], [stale_response], stale_baseline, stale_after, stale_status=stale_status, unchanged_fields=list(PROJECTION_STABILITY_FIELDS))
    restored_settings = _settings(http)
    restored_settings["capacity_registers"] = 64
    _call(http, "PUT", "/api/v1/datalink/settings/modbus_share", {"value": restored_settings, "expected_settings_revision": restored_settings["settings_revision"]}, {200})

    conflict_http_port = free_port()
    conflict_db = temp_root / "bind-conflict.db"
    conflict_log = temp_root / "bind-conflict.log"
    conflict_proc = start_gateway(binary, conflict_db, conflict_http_port, conflict_log, Path(__file__).resolve().parents[2])
    try:
        conflict_http = HTTP(f"http://127.0.0.1:{conflict_http_port}")
        wait_gateway(conflict_http)
        conflict_settings = _settings(conflict_http)
        conflict_settings.update({"enabled": True, "bind_address": "127.0.0.1", "port": share_port, "slave_id": 1, "capacity_registers": 64})
        _, conflict_payload, conflict_response = _call(conflict_http, "PUT", "/api/v1/datalink/settings/modbus_share", {"value": conflict_settings, "expected_settings_revision": conflict_settings["settings_revision"]}, {422})
        conflict_status = data(conflict_http.call("GET", "/api/v1/datalink/modbus-share/status", expected={200})[1])
        retained = read_holding("127.0.0.1", share_port)[0]
        if conflict_response.get("code") != "modbus_share_listener_bind_failed" or conflict_status.get("lifecycle_state") != "failed" or retained != source.value:
            raise RuntimeError(f"bind conflict did not preserve old projection: {conflict_payload}, status={conflict_status}, value={retained}")
        _step(result, "listener-bind-conflict-old-projection-retained", [{"method": "PUT", "path": "/api/v1/datalink/settings/modbus_share"}], [conflict_response, _response_summary("GET", "/api/v1/datalink/modbus-share/status", 200, {"success": True, "data": conflict_status})], before_double, _projection_state(http, share_port), listener_failure="modbus_share_listener_bind_failed", retained_value=retained, partial_failure_exe_fallback="real listener bind conflict; primary old projection retained")
    finally:
        try:
            stop_gateway(conflict_proc, conflict_log)
        except Exception:
            conflict_proc.kill()
            conflict_proc.wait(timeout=3)
        command_log.extend(conflict_http.commands if "conflict_http" in locals() else [])


def run_post_collision_negatives(http: HTTP, share_port: int, source: SourceModbus, result: dict[str, Any]) -> None:
    """Exercise capacity overflow and a foreign workspace request."""
    overflow = data(http.call("POST", "/api/v1/datalink/studio-v2/workspace/source-rules", {
        "id": "b10-rule-capacity",
        "device_id": "b10-device",
        "start_address": "40003",
        "count": 1,
        "data_type": "int32",
        "naming_prefix": "B10CAP",
        "enabled": True,
        "share_enabled": True,
        "share_start_register": 40064,
        "share_stride": 2,
    }, expected={201})[1])
    boot = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
    candidates_path = candidate_scope_path(
        f"/api/v1/datalink/source-rules/{overflow['id']}/candidates",
        boot["workspace_id"],
        boot["workspace_revision"],
        overflow["revision_id"],
    )
    candidates = data(http.call("GET", candidates_path, expected={200})[1])
    tag_candidates = candidates.get("tags", {}).get("candidates", [])
    if not tag_candidates:
        raise RuntimeError(f"capacity fixture has no tag candidates: {candidates}")
    http.call("POST", f"/api/v1/datalink/source-rules/{overflow['id']}/tags/apply", {
        "workspace_id": boot["workspace_id"],
        "expected_workspace_revision": boot["workspace_revision"],
        "revision_id": overflow["revision_id"],
        "candidate_ids": [candidate["id"] for candidate in tag_candidates],
    }, expected={200})
    candidates = data(http.call("GET", candidates_path, expected={200})[1])
    boot = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
    desired = boot.get("canonical_desired_mappings", [])
    if not any(mapping.get("source_rule_id") == overflow["id"] for mapping in desired):
        raise RuntimeError(f"capacity fixture did not reach canonical desired set: rule={overflow}, candidates={candidates}, boot={boot}")
    body = {"workspace_id": boot["workspace_id"], "expected_workspace_revision": boot["workspace_revision"], "expected_settings_revision": boot["canonical_plan"]["settings_revision"], "readiness_token": boot["readiness_token"], "canonical_plan_signature": boot["canonical_plan"]["signature"]}
    before = _projection_state(http, share_port)
    status, payload, response = _call(http, "POST", "/api/v1/datalink/modbus-share/reconcile", body, {422})
    if response.get("code") not in {"modbus_share_capacity_exceeded", "modbus_share_invalid_geometry"}:
        raise RuntimeError(f"capacity overflow returned unexpected code: {payload}")
    after = _projection_state(http, share_port)
    if before.get("mapping_hash") != after.get("mapping_hash") or before.get("register_value") != after.get("register_value"):
        raise RuntimeError(f"capacity failure changed old projection: before={before} after={after}")
    _step(result, "capacity-overflow-old-projection-retained", [{"method": "POST", "path": "/api/v1/datalink/modbus-share/reconcile"}], [response], before, after, typed_code=response.get("code"), candidate_count=len(desired), status_code=status)

    # Leave the capacity fixture out of the next collision fixture's canonical
    # plan; this cleanup is also public HTTP and does not touch the database.
    http.call("PUT", f"/api/v1/datalink/studio-v2/workspace/source-rules/{overflow['id']}", {"share_enabled": False}, expected={200})

    foreign_baseline = _projection_state(http, share_port)
    foreign_boot = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
    foreign_body = {"workspace_id": "b10-foreign-workspace", "expected_workspace_revision": foreign_boot["workspace_revision"], "expected_settings_revision": foreign_boot["canonical_plan"]["settings_revision"], "readiness_token": foreign_boot["readiness_token"], "desired_mappings": []}
    foreign_status, foreign_payload, foreign_response = _call(http, "POST", "/api/v1/datalink/modbus-share/reconcile", foreign_body, {403})
    if foreign_response.get("code") != "modbus_share_workspace_scope":
        raise RuntimeError(f"foreign workspace request returned unexpected code: {foreign_payload}")
    foreign_after = _projection_state(http, share_port)
    assert_projection_unchanged(foreign_baseline, foreign_after, "cross-workspace reconcile")
    _step(result, "cross-workspace-ownership-no-mutation", [{"method": "POST", "path": "/api/v1/datalink/modbus-share/reconcile"}], [foreign_response], foreign_baseline, foreign_after, typed_code=foreign_response.get("code"), status_code=foreign_status, unchanged_fields=list(PROJECTION_STABILITY_FIELDS))

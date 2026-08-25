#!/usr/bin/env python3
"""B10 black-box acceptance for the production test_ui executable."""
from __future__ import annotations

import copy
import hashlib
import json
import os
import subprocess
import sys
import tempfile
import time
from pathlib import Path
from typing import Any

from lib.b10_acceptance_helpers import HTTP, SourceModbus, candidate_scope_path, data, free_port, finalize_b10_witness, read_holding, start_gateway, stop_gateway, wait_gateway
from lib.b10_focused_suite import run_focused_modbusshare_suite
from lib.b10_negative_steps import run_post_collision_negatives, run_pre_collision_negatives
from lib.b10_projection_assertions import assert_positive_evidence_preserved, capture_positive_projection, capture_restart_projection

ROOT = Path(__file__).resolve().parents[1]
CHANGE = ROOT / "openspec/changes/productionize-modbus-share-lifecycle"


def run() -> dict[str, Any]:
    started = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
    temp_root = Path(tempfile.mkdtemp(prefix="gateway-b10-"))
    binary = temp_root / "test-ui"
    db_path = temp_root / "gateway.db"
    log_path = temp_root / "gateway.log"
    http_port, source_port, share_port = free_port(), free_port(), free_port()
    source = SourceModbus(source_port)
    proc: subprocess.Popen[bytes] | None = None
    command_log: list[str] = []
    result: dict[str, Any] = {
        "started_at": started,
        "temp_root": str(temp_root),
        "binary": str(binary),
        "database": str(db_path),
        "ports": {"http": http_port, "source_modbus": source_port, "share_modbus": share_port},
        "steps": [],
        "configuration_seam": "Studio V2 workspace source-rule HTTP create/update only",
    }

    def record_focused_suite(label: str, cache_dir: Path) -> dict[str, Any]:
        evidence = run_focused_modbusshare_suite(ROOT, cache_dir)
        result.setdefault("focused_suite_evidence", []).append(evidence)
        result["steps"].append({"name": label, "status": "pass" if evidence["exit_code"] == 0 else "fail", "command": evidence["command"], "test_names": evidence["test_names"], "exit_code": evidence["exit_code"], "output_hash": evidence["output_hash"]})
        if evidence["exit_code"] != 0:
            raise RuntimeError(f"focused modbusshare suite failed: {evidence}")
        return evidence

    def scoped_candidate(http: HTTP, base_path: str, revision_id: str) -> tuple[str, dict[str, Any]]:
        current_workspace = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
        scoped_path = candidate_scope_path(
            base_path,
            current_workspace["workspace_id"],
            current_workspace["workspace_revision"],
            revision_id,
        )
        return scoped_path, current_workspace

    try:
        record_focused_suite("partial-failure-focused-suite-before-build", temp_root / "focused-cache-before")
        build = ["go", "build", "-trimpath", "-buildvcs=false", "-ldflags=-buildid=", "-o", str(binary), "./cmd/test_ui"]
        result["build_command"] = " ".join(build)
        build_env = os.environ.copy()
        build_env["GOCACHE"] = str(temp_root / "go-cache")
        completed = subprocess.run(build, cwd=ROOT, env=build_env, check=True, capture_output=True, text=True)
        result["build_output"] = completed.stdout[-1000:]
        result["binary_sha256"] = hashlib.sha256(binary.read_bytes()).hexdigest()
        record_focused_suite("partial-failure-focused-suite-after-build", temp_root / "focused-cache-after")
        source.start()

        proc = start_gateway(binary, db_path, http_port, log_path, ROOT)
        http = HTTP(f"http://127.0.0.1:{http_port}")
        wait_gateway(http)
        workspace = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])
        share_bootstrap = workspace["modbus_share"]
        workspace_id = share_bootstrap["workspace_id"]
        result["workspace_id"] = workspace_id
        result["steps"].append({"name": "disabled-no-listener", "status": "pass", "tcp_connect_before_enable": False})
        try:
            read_holding("127.0.0.1", share_port)
            raise RuntimeError("disabled listener accepted a Modbus connection")
        except OSError:
            pass
        device = data(http.call("POST", "/api/v1/datalink/studio-v2/workspace/devices", {
            "id": "b10-device",
            "name": "B10 source Modbus",
            "protocol": "modbus_tcp",
            "connection_config": {"host": "127.0.0.1", "port": source_port, "slave_id": 1, "timeout": 2},
        }, expected={201})[1])
        device_id = device["id"]
        test = data(http.call("POST", f"/api/v1/datalink/devices/{device_id}/test", expected={200})[1])
        if test.get("success") is False or test.get("connect", {}).get("status") == "failed":
            raise RuntimeError(f"source device test failed: {test}")
        rule = data(http.call("POST", "/api/v1/datalink/studio-v2/workspace/source-rules", {
            "id": "b10-rule-1",
            "device_id": device_id,
            "start_address": "40001",
            "count": 1,
            "data_type": "int16",
            "naming_prefix": "B10",
            "enabled": True,
            "share_enabled": True,
            "share_start_register": 40001,
            "share_stride": 1,
        }, expected={201})[1])
        rule_id = rule["id"]
        if rule_id != "b10-rule-1":
            raise RuntimeError(f"unexpected rule id: {rule_id}")
        updated_rule = data(http.call("PUT", f"/api/v1/datalink/studio-v2/workspace/source-rules/{rule_id}", {
            "share_enabled": True,
            "share_start_register": 40001,
            "share_stride": 1,
        }, expected={200})[1])
        if updated_rule.get("share_start_register") != 40001 or updated_rule.get("share_stride") != 1:
            raise RuntimeError(f"source-rule HTTP Share fields did not roundtrip: {updated_rule}")
        candidates_path, candidate_scope = scoped_candidate(http, f"/api/v1/datalink/source-rules/{rule_id}/candidates", updated_rule["revision_id"])
        candidates = data(http.call("GET", candidates_path, expected={200})[1])
        tag_candidates = candidates.get("tags", {}).get("candidates", [])
        candidate_ids = [candidate["id"] for candidate in tag_candidates]
        result["http_candidate_snapshot"] = candidates
        if not candidate_ids:
            raise RuntimeError(f"source-rule tag candidates are empty: {candidates}")
        http.call("POST", f"/api/v1/datalink/source-rules/{rule_id}/tags/apply", {
            "workspace_id": candidate_scope["workspace_id"],
            "expected_workspace_revision": candidate_scope["workspace_revision"],
            "revision_id": updated_rule["revision_id"],
            "candidate_ids": candidate_ids,
        }, expected={200})
        result["workspace_mappings_after_tags"] = http.call("GET", "/api/v1/datalink/studio-v2/workspace/mappings", expected={200})[1].get("data", [])
        refreshed_rule = data(http.call("PUT", f"/api/v1/datalink/studio-v2/workspace/source-rules/{rule_id}", {
            "share_enabled": True,
            "share_start_register": 40001,
            "share_stride": 1,
        }, expected={200})[1])
        if refreshed_rule.get("share_start_register") != 40001 or refreshed_rule.get("share_stride") != 1:
            raise RuntimeError(f"source-rule HTTP refresh did not preserve Share fields: {refreshed_rule}")
        result["steps"].append({"name": "http-configure-durable-save", "status": "pass", "device_id": device_id, "rule_id": rule_id})
        command_log.extend(http.commands)

        settings_response = http.call("GET", "/api/v1/datalink/settings", expected={200})[1]
        settings_items = settings_response.get("data")
        if not isinstance(settings_items, list):
            raise RuntimeError(f"settings response has no list data: {settings_response}")
        settings_item = next(item for item in settings_items if item["key"] == "modbus_share")
        settings = dict(settings_item["value"])
        settings.update({"enabled": True, "bind_address": "127.0.0.1", "port": share_port, "slave_id": 1, "capacity_registers": 64})
        setting_item = data(http.call("PUT", "/api/v1/datalink/settings/modbus_share", {
            "value": settings,
            "expected_settings_revision": settings["settings_revision"],
        }, expected={200})[1])
        setting = setting_item.get("value")
        if not isinstance(setting, dict) or not setting.get("settings_revision"):
            raise RuntimeError(f"settings update response has no durable revision: {setting_item}")
        status = data(http.call("GET", "/api/v1/datalink/modbus-share/status", expected={200})[1])
        if status["bind_address"] != "127.0.0.1" or status["port"] != share_port or status["lifecycle_state"] != "running":
            raise RuntimeError(f"listener did not start from durable settings: {status}")
        http.call("POST", f"/api/v1/datalink/devices/{device_id}/activate", expected={200})
        enabled_rule = data(http.call("POST", f"/api/v1/datalink/source-rules/{rule_id}/enable", expected={200})[1])
        if enabled_rule.get("enabled") is not True:
            raise RuntimeError(f"source-rule public enable did not persist: {enabled_rule}")
        result["device_activate_after_rule_enable"] = data(http.call("POST", f"/api/v1/datalink/devices/{device_id}/activate", expected={200})[1])
        status = data(http.call("GET", "/api/v1/datalink/modbus-share/status", expected={200})[1])
        boot = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
        result["status_after_http_enable"] = status
        result["bootstrap_after_http_enable"] = boot
        result["canonical_after_http_activation"] = boot.get("canonical_desired_mappings", [])
        result["canonical_plan_signature"] = boot.get("canonical_plan", {}).get("signature")
        canonical_settings_revision = boot.get("canonical_plan", {}).get("settings_revision")
        result["settings_revision_sources"] = {
            "settings_put": setting.get("settings_revision"),
            "fresh_status": status.get("settings_revision"),
            "fresh_canonical_plan": canonical_settings_revision,
            "bootstrap_top_level": boot.get("settings_revision"),
        }
        if not result["canonical_after_http_activation"] or not result["canonical_plan_signature"]:
            raise RuntimeError(f"HTTP source-rule enable produced no canonical desired mappings: {boot}")
        if not canonical_settings_revision:
            raise RuntimeError(f"fresh canonical plan has no settings revision: {boot}")
        reconcile_status, reconcile = http.call("POST", "/api/v1/datalink/modbus-share/reconcile", {
            "workspace_id": boot["workspace_id"],
            "expected_workspace_revision": boot["workspace_revision"],
            "expected_settings_revision": canonical_settings_revision,
            "readiness_token": boot["readiness_token"],
            "canonical_plan_signature": result["canonical_plan_signature"],
        }, expected={200})
        result["reconcile"] = {"status": reconcile_status, "outcome": reconcile.get("data", {}).get("outcome")}
        status = data(http.call("GET", "/api/v1/datalink/modbus-share/status", expected={200})[1])
        boot = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
        result["status_after_reconcile"] = status
        result["bootstrap_after_reconcile"] = boot
        if status.get("lifecycle_state") != "running":
            raise RuntimeError(f"listener stopped before workspace activation: {status}")
        activation = http.call("POST", "/api/v1/datalink/studio-v2/workspace/activate", {
            "workspace_revision": boot["workspace_revision"],
            "settings_revision": canonical_settings_revision,
            "readiness_token": boot["readiness_token"],
            "pending_saves": 0,
        }, expected={200})
        if activation[0] != 200:
            raise RuntimeError(f"activation failed: {activation[1]}")
        result["workspace_mappings_after_activate"] = http.call("GET", "/api/v1/datalink/studio-v2/workspace/mappings", expected={200})[1].get("data", [])
        result["runtime_status_after_activate"] = data(http.call("GET", f"/api/v1/datalink/runtime/status?device_id={device_id}", expected={200})[1])
        positive_status, positive_bootstrap, positive_mappings, positive_evidence = capture_positive_projection(http)
        result["status_after_activate"] = positive_status
        result["bootstrap_after_activate"] = positive_bootstrap
        result["share_mappings_after_activate"] = positive_mappings
        result["positive_projection"] = positive_evidence
        result["positive_projection_snapshot"] = copy.deepcopy(positive_evidence)
        value = None
        deadline = time.time() + 15
        while time.time() < deadline:
            try:
                value = read_holding("127.0.0.1", share_port)[0]
                if value == source.value:
                    break
            except OSError:
                pass
            time.sleep(0.2)
        if value != source.value:
            result["runtime_status_after_read_timeout"] = data(http.call("GET", f"/api/v1/datalink/runtime/status?device_id={device_id}", expected={200})[1])
            raise RuntimeError(f"collector projection value mismatch: got {value}, want {source.value}")
        result["steps"].append({"name": "activate-real-modbus-read", "status": "pass", "register": 0, "value": value, "assertions": positive_evidence["assertions"]})
        command_log.extend(http.commands)
        if not stop_gateway(proc, log_path):
            raise RuntimeError("gateway did not gracefully terminate after first read")
        proc = start_gateway(binary, db_path, http_port, log_path, ROOT)
        http = HTTP(f"http://127.0.0.1:{http_port}")
        wait_gateway(http)
        restart_status, restart_bootstrap, restart_mappings, restart_evidence = capture_restart_projection(http, positive_evidence)
        result["status_after_restart"] = restart_status
        result["bootstrap_after_restart"] = restart_bootstrap
        result["share_mappings_after_restart"] = restart_mappings
        result["restart_projection"] = restart_evidence
        restored = None
        deadline = time.time() + 15
        while time.time() < deadline:
            try:
                restored = read_holding("127.0.0.1", share_port)[0]
                if restored == source.value:
                    break
            except OSError:
                pass
            time.sleep(0.2)
        if restored != source.value:
            raise RuntimeError(f"restart restore value mismatch: got {restored}, want {source.value}")
        result["steps"].append({"name": "same-exe-restart-restore-read", "status": "pass", "value": restored, "assertions": restart_evidence["assertions"]})

        run_pre_collision_negatives(http, binary, temp_root, share_port, source, result, command_log)
        run_post_collision_negatives(http, share_port, source, result)
        assert_positive_evidence_preserved(result["positive_projection_snapshot"], result["positive_projection"])

        second = data(http.call("POST", "/api/v1/datalink/studio-v2/workspace/source-rules", {
            "id": "b10-rule-collision",
            "device_id": device_id,
            "start_address": "40002",
            "count": 1,
            "data_type": "int16",
            "naming_prefix": "B10C",
            "enabled": True,
            "share_enabled": True,
            "share_start_register": 40001,
            "share_stride": 1,
        }, expected={201})[1])
        second_candidates_path, second_candidate_scope = scoped_candidate(http, f"/api/v1/datalink/source-rules/{second['id']}/candidates", second["revision_id"])
        second_candidates = data(http.call("GET", second_candidates_path, expected={200})[1])
        second_tag_candidates = second_candidates.get("tags", {}).get("candidates", [])
        http.call("POST", f"/api/v1/datalink/source-rules/{second['id']}/tags/apply", {
            "workspace_id": second_candidate_scope["workspace_id"],
            "expected_workspace_revision": second_candidate_scope["workspace_revision"],
            "revision_id": second["revision_id"],
            "candidate_ids": [candidate["id"] for candidate in second_tag_candidates],
        }, expected={200})
        second_candidates_path, _ = scoped_candidate(http, f"/api/v1/datalink/source-rules/{second['id']}/candidates", second["revision_id"])
        second_candidates = data(http.call("GET", second_candidates_path, expected={200})[1])
        result["collision_candidate_snapshot"] = second_candidates
        command_log.extend(http.commands)
        boot = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
        canonical = boot.get("canonical_desired_mappings", [])
        local_outputs = second_candidates.get("local_modbus_outputs", {})
        local_candidates = local_outputs.get("candidates", [])
        candidate_statuses = [candidate.get("status") for candidate in local_candidates]
        result["collision_bootstrap"] = boot
        result["collision_state"] = {
            "candidate_set_status": local_outputs.get("status"),
            "candidate_statuses": candidate_statuses,
            "canonical_count": len(canonical),
        }
        if "blocked_conflict" not in candidate_statuses or canonical:
            raise RuntimeError(f"collision was not blocked at candidate stage: {result['collision_state']}")
        retained = read_holding("127.0.0.1", share_port)[0]
        if retained != source.value:
            raise RuntimeError(f"collision changed old projection: got {retained}")
        result["steps"].append({"name": "collision-blocked-no-mutation", "status": "pass", "candidate_status": "blocked_conflict", "retained_value": retained})
        result["partial_failure"] = {
            "status": "pass",
            "focused_suite": result["focused_suite_evidence"],
            "exe_fallback": "listener bind conflict returned modbus_share_listener_bind_failed; primary old projection retained",
        }
        result["status"] = "pass"
        command_log.extend(http.commands)
        return result
    finally:
        current_http = locals().get("http")
        finalize_b10_witness(proc, log_path, source, result, command_log, current_http, CHANGE, temp_root)


if __name__ == "__main__":
    try:
        output = run()
    except Exception as error:
        print(f"B10 acceptance failed: {error}", file=sys.stderr)
        raise
    print(json.dumps(output, indent=2, ensure_ascii=False))

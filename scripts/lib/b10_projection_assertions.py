"""B10 assertions for the positive and restart Share projections."""

from __future__ import annotations

import copy
from typing import Any

from lib.b10_acceptance_helpers import HTTP, data, list_data


MAPPING_CONTRACT_FIELDS = (
    "workspace_id", "source_rule_id", "source_rule_revision", "tag_id", "mapping_id",
    "data_type", "share_start_register", "zero_based_register", "span_registers",
    "stride_registers", "capacity_registers", "tag_key", "display_name",
)

PROJECTION_STABILITY_FIELDS = (
    "mapping_hash", "register_value", "mapping_count", "workspace_revision", "settings_revision",
)


def _mapping_contract(mapping: dict[str, Any]) -> dict[str, Any]:
    missing = [field for field in MAPPING_CONTRACT_FIELDS if field not in mapping]
    if missing:
        raise RuntimeError(f"mapping is missing identity/content fields: {missing}; mapping={mapping}")
    return {field: mapping[field] for field in MAPPING_CONTRACT_FIELDS}


def _sorted_mapping_contracts(mappings: list[dict[str, Any]]) -> list[dict[str, Any]]:
    return sorted((_mapping_contract(mapping) for mapping in mappings), key=lambda item: (item["mapping_id"], item["tag_id"]))


def assert_projection_unchanged(before: dict[str, Any], after: dict[str, Any], label: str) -> None:
    """Require a negative request to leave every observable projection field unchanged."""
    missing = [
        field
        for field in PROJECTION_STABILITY_FIELDS
        if field not in before or field not in after
    ]
    if missing:
        raise RuntimeError(f"{label} stability snapshot is missing fields: {missing}")
    before_projection = {field: before.get(field) for field in PROJECTION_STABILITY_FIELDS}
    after_projection = {field: after.get(field) for field in PROJECTION_STABILITY_FIELDS}
    if before_projection != after_projection:
        raise RuntimeError(f"{label} changed the projection: before={before_projection} after={after_projection}")


def assert_positive_projection(canonical_plan: dict[str, Any], production_mappings: list[dict[str, Any]], status: dict[str, Any], label: str = "positive Share projection") -> dict[str, Any]:
    canonical_mappings = canonical_plan.get("desired_mappings")
    if not isinstance(canonical_mappings, list) or not canonical_mappings:
        raise RuntimeError(f"{label} canonical desired mapping set is empty")
    if not production_mappings:
        raise RuntimeError(f"{label} production Share list is empty")
    expected = _sorted_mapping_contracts(canonical_mappings)
    actual = _sorted_mapping_contracts(production_mappings)
    if actual != expected:
        raise RuntimeError(f"{label} identity/content mismatch: canonical={expected}, production={actual}")
    mapping_count = status.get("mapping_count")
    if mapping_count != len(production_mappings) or mapping_count != len(canonical_mappings):
        raise RuntimeError(f"{label} status mapping_count={mapping_count} does not match production list={len(production_mappings)} and canonical plan={len(canonical_mappings)}")
    capacities = {mapping["capacity_registers"] for mapping in production_mappings}
    if capacities != {status.get("capacity_registers")}:
        raise RuntimeError(f"{label} mapping capacity_registers={capacities} does not match status capacity_registers={status.get('capacity_registers')}")
    if status.get("lifecycle_state") != "running":
        raise RuntimeError(f"{label} status is not running: {status}")
    return {
        "mappings": [dict(mapping) for mapping in production_mappings],
        "canonical_mappings": [dict(mapping) for mapping in canonical_mappings],
        "status": {key: status.get(key) for key in ("lifecycle_state", "enabled", "bind_address", "port", "slave_id", "capacity_registers", "mapping_count") if key in status},
        "revisions": {"workspace_revision": status.get("workspace_revision", canonical_plan.get("workspace_revision")), "settings_revision": status.get("settings_revision", canonical_plan.get("settings_revision"))},
        "assertions": ["production Share list is non-empty", "production Share identity/content equals canonical plan", "mapping capacity_registers equals listener capacity", "status mapping_count equals canonical and production list counts"],
    }


def assert_restart_projection(positive_evidence: dict[str, Any], restart_plan: dict[str, Any], restart_status: dict[str, Any], restart_mappings: list[dict[str, Any]]) -> dict[str, Any]:
    if restart_status.get("lifecycle_state") != "running":
        raise RuntimeError(f"restart listener lifecycle_state is not running: {restart_status}")
    restart_evidence = assert_positive_projection(restart_plan, restart_mappings, restart_status, "restart Share projection")
    restart_evidence["restart_listener_status"] = {
        "lifecycle_state": restart_status.get("lifecycle_state"),
        "bind_address": restart_status.get("bind_address"),
        "port": restart_status.get("port"),
        "slave_id": restart_status.get("slave_id"),
    }
    for revision_name in ("workspace_revision", "settings_revision"):
        expected = positive_evidence.get("revisions", {}).get(revision_name)
        observed = restart_evidence["revisions"].get(revision_name)
        if expected and observed != expected:
            raise RuntimeError(f"restart {revision_name} mismatch: before={expected}, after={observed}")
    if _sorted_mapping_contracts(restart_mappings) != _sorted_mapping_contracts(positive_evidence["mappings"]):
        raise RuntimeError("restart production Share list content mismatch")
    for key, expected in positive_evidence.get("status", {}).items():
        if key != "mapping_count" and expected is not None and restart_evidence["status"].get(key) != expected:
            raise RuntimeError(f"restart status {key} mismatch: before={expected}, after={restart_evidence['status'].get(key)}")
    restart_evidence["assertions"].extend(["restart listener status payload lifecycle_state=running", "restart workspace/settings revisions equal positive projection", "restart production Share list content and mapping capacity_registers equal positive projection", "restart status equals positive projection"])
    return restart_evidence


def assert_positive_evidence_preserved(expected: dict[str, Any], observed: dict[str, Any]) -> None:
    if observed != expected:
        raise RuntimeError("positive Share evidence was overwritten by a negative/capacity step")


def capture_positive_projection(http: HTTP) -> tuple[dict[str, Any], dict[str, Any], list[dict[str, Any]], dict[str, Any]]:
    status = data(http.call("GET", "/api/v1/datalink/modbus-share/status", expected={200})[1])
    bootstrap = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
    mappings = list_data(http.call("GET", "/api/v1/datalink/modbus-share/mappings", expected={200})[1])
    evidence = assert_positive_projection(bootstrap["canonical_plan"], mappings, status)
    evidence["bootstrap"] = copy.deepcopy(bootstrap)
    return status, bootstrap, mappings, evidence


def capture_restart_projection(http: HTTP, positive_evidence: dict[str, Any]) -> tuple[dict[str, Any], dict[str, Any], list[dict[str, Any]], dict[str, Any]]:
    status = data(http.call("GET", "/api/v1/datalink/modbus-share/status", expected={200})[1])
    bootstrap = data(http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})[1])["modbus_share"]
    mappings = list_data(http.call("GET", "/api/v1/datalink/modbus-share/mappings", expected={200})[1])
    evidence = assert_restart_projection(positive_evidence, bootstrap["canonical_plan"], status, mappings)
    return status, bootstrap, mappings, evidence

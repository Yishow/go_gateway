import copy
import io
import tempfile
import unittest
from pathlib import Path
from unittest import mock

from lib import b10_acceptance_helpers
from lib.b10_projection_assertions import (
    assert_positive_projection,
    assert_restart_projection,
    assert_positive_evidence_preserved,
    assert_projection_unchanged,
)


def _fixtures():
    mapping = {
        "workspace_id": "ws-a",
        "source_rule_id": "rule-a",
        "source_rule_revision": "rule-rev-a",
        "tag_id": "tag-a",
        "mapping_id": "mapping-a",
        "data_type": "int16",
        "share_start_register": 40001,
        "zero_based_register": 0,
        "span_registers": 1,
        "stride_registers": 1,
        "capacity_registers": 64,
        "tag_key": "B10_40001",
        "display_name": "B10_40001",
    }
    plan = {
        "workspace_revision": "workspace-rev-a",
        "settings_revision": "settings-rev-a",
        "capacity_registers": 64,
        "desired_mappings": [mapping],
    }
    status = {
        "lifecycle_state": "running",
        "mapping_count": 1,
        "workspace_revision": "workspace-rev-a",
        "settings_revision": "settings-rev-a",
        "capacity_registers": 64,
    }
    return plan, [copy.deepcopy(mapping)], status


class B10AcceptanceHelpersTest(unittest.TestCase):
    class _Process:
        def __init__(self, poll_value=None, *, kill_error=None):
            self.poll_value = poll_value
            self.kill_error = kill_error
            self.kill_calls = 0
            self.wait_calls = 0
            self._b10_log = io.StringIO()

        def poll(self):
            return self.poll_value

        def send_signal(self, value):
            self.signal = value

        def kill(self):
            self.kill_calls += 1
            if self.kill_error is not None:
                raise self.kill_error

        def wait(self, timeout=None):
            self.wait_calls += 1
            return self.poll_value

    class _Source:
        def __init__(self, error=None):
            self.error = error
            self.stop_calls = 0

        def stop(self):
            self.stop_calls += 1
            if self.error is not None:
                raise self.error

    def _finalize(self, proc, source, result):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            log_path = root / "gateway.log"
            log_path.write_text("failure log", encoding="utf-8")
            change = root / "change"
            temp_root = root / "temp"
            temp_root.mkdir()
            b10_acceptance_helpers.finalize_b10_witness(
                proc,
                log_path,
                source,
                result,
                [],
                None,
                change,
                temp_root,
            )
            witnesses = list((change / "witness").glob("b10-exe-acceptance-*.json"))
            self.assertEqual(len(witnesses), 1)
            self.assertIn("original acceptance failure", witnesses[0].read_text(encoding="utf-8"))
            self.assertFalse(temp_root.exists())

    def test_finalize_does_not_kill_process_that_already_exited(self):
        proc = self._Process(poll_value=17)
        source = self._Source()
        result = {"status": "fail", "error": "original acceptance failure"}

        self._finalize(proc, source, result)

        self.assertEqual(proc.kill_calls, 0)
        self.assertEqual(source.stop_calls, 1)
        self.assertTrue(proc._b10_log.closed)

    def test_stop_gateway_returns_true_after_graceful_wait(self):
        proc = self._Process(poll_value=None)

        self.assertTrue(b10_acceptance_helpers.stop_gateway(proc, Path("gateway.log")))

        self.assertEqual(proc.wait_calls, 1)
        self.assertTrue(proc._b10_log.closed)

    def test_candidate_scope_path_includes_workspace_revision_and_rule_revision(self):
        path = b10_acceptance_helpers.candidate_scope_path(
            "/api/v1/datalink/source-rules/rule-a/candidates",
            "workspace-a",
            "workspace-rev-a",
            "rule-rev-a",
        )

        self.assertEqual(
            path,
            "/api/v1/datalink/source-rules/rule-a/candidates?workspace_id=workspace-a&expected_workspace_revision=workspace-rev-a&revision_id=rule-rev-a",
        )

    def test_finalize_ignores_kill_process_lookup_error_and_writes_witness(self):
        proc = self._Process(poll_value=None, kill_error=ProcessLookupError("gone"))
        source = self._Source()
        result = {"status": "fail", "error": "original acceptance failure"}

        with mock.patch.object(b10_acceptance_helpers, "stop_gateway", side_effect=RuntimeError("stop failed")):
            self._finalize(proc, source, result)

        self.assertEqual(proc.kill_calls, 1)
        self.assertEqual(source.stop_calls, 1)
        self.assertTrue(proc._b10_log.closed)

    def test_finalize_writes_witness_when_source_stop_raises(self):
        proc = self._Process(poll_value=17)
        source = self._Source(RuntimeError("source cleanup failed"))
        result = {"status": "fail", "error": "original acceptance failure"}

        self._finalize(proc, source, result)

        self.assertEqual(source.stop_calls, 1)
        self.assertEqual(proc.kill_calls, 0)
        self.assertTrue(proc._b10_log.closed)

    def test_positive_projection_requires_non_empty_canonical_and_production_list(self):
        plan, mappings, status = _fixtures()
        evidence = assert_positive_projection(plan, mappings, status)
        self.assertEqual(evidence["mappings"], mappings)
        with self.assertRaisesRegex(RuntimeError, "production Share list is empty"):
            assert_positive_projection(plan, [], status)

    def test_positive_projection_requires_identity_content_and_status_count(self):
        plan, mappings, status = _fixtures()
        mappings[0]["mapping_id"] = "wrong"
        with self.assertRaisesRegex(RuntimeError, "identity/content mismatch"):
            assert_positive_projection(plan, mappings, status)
        _, mappings, status = _fixtures()
        status["mapping_count"] = 0
        with self.assertRaisesRegex(RuntimeError, "mapping_count"):
            assert_positive_projection(plan, mappings, status)

    def test_restart_requires_same_revisions_list_and_status(self):
        plan, mappings, status = _fixtures()
        evidence = assert_positive_projection(plan, mappings, status)
        restart_status = copy.deepcopy(status)
        restart_plan = copy.deepcopy(plan)
        restart_mappings = copy.deepcopy(mappings)
        assert_restart_projection(evidence, restart_plan, restart_status, restart_mappings)
        restart_status["workspace_revision"] = "wrong"
        with self.assertRaisesRegex(RuntimeError, "workspace_revision"):
            assert_restart_projection(evidence, restart_plan, restart_status, restart_mappings)

    def test_restart_requires_listener_status_payload_running(self):
        plan, mappings, status = _fixtures()
        evidence = assert_positive_projection(plan, mappings, status)
        status["lifecycle_state"] = "failed"
        with self.assertRaisesRegex(RuntimeError, "restart listener lifecycle_state"):
            assert_restart_projection(evidence, copy.deepcopy(plan), status, copy.deepcopy(mappings))

    def test_restart_rejects_mapping_capacity_drift(self):
        plan, mappings, status = _fixtures()
        evidence = assert_positive_projection(plan, mappings, status)
        restart_mapping = copy.deepcopy(mappings[0])
        restart_mapping["capacity_registers"] = 32
        with self.assertRaisesRegex(RuntimeError, "identity/content mismatch"):
            assert_restart_projection(evidence, copy.deepcopy(plan), copy.deepcopy(status), [restart_mapping])

    def test_positive_projection_rejects_status_capacity_drift(self):
        plan, mappings, status = _fixtures()
        status["capacity_registers"] = 32
        with self.assertRaisesRegex(RuntimeError, "mapping capacity_registers"):
            assert_positive_projection(plan, mappings, status)

    def test_negative_projection_assertion_rejects_old_witness_mutations(self):
        baseline = {
            "mapping_hash": "mapping-a",
            "register_value": 0x1234,
            "mapping_count": 1,
            "workspace_revision": "workspace-rev-a",
            "settings_revision": "settings-rev-a",
        }
        for field, changed in (
            ("mapping_hash", "mapping-b"),
            ("register_value", 0x4321),
            ("mapping_count", 0),
            ("workspace_revision", "workspace-rev-b"),
            ("settings_revision", "settings-rev-b"),
        ):
            with self.subTest(field=field):
                observed = dict(baseline)
                observed[field] = changed
                with self.assertRaisesRegex(RuntimeError, "changed the projection"):
                    assert_projection_unchanged(baseline, observed, "regression witness")

    def test_negative_projection_assertion_requires_complete_snapshot(self):
        baseline = {"mapping_hash": "mapping-a"}
        with self.assertRaisesRegex(RuntimeError, "missing fields"):
            assert_projection_unchanged(baseline, baseline, "incomplete witness")

    def test_negative_steps_cannot_overwrite_positive_evidence(self):
        plan, mappings, status = _fixtures()
        evidence = assert_positive_projection(plan, mappings, status)
        preserved = copy.deepcopy(evidence)
        assert_positive_evidence_preserved(preserved, evidence)
        evidence["mappings"].clear()
        with self.assertRaisesRegex(RuntimeError, "positive Share evidence was overwritten"):
            assert_positive_evidence_preserved(preserved, evidence)


if __name__ == "__main__":
    unittest.main()

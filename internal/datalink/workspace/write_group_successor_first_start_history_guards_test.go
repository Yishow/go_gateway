package workspace

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Intentionally contradictory rows represent late arrivals or durable corruption;
// they are added after valid review/readiness, never removed to fake NoHistory.
func TestSuccessorFirstStartAtomicLateHistoryRejectsConfirmation(t *testing.T) {
	for _, kind := range []string{"version", "runtime", "sample", "outbox_receipt", "checkpoint", "bucket", "counter_baseline", "counter_proof", "revision_end", "fence", "handoff_source", "handoff_target", "first_start_pending", "first_start_unknown"} {
		t.Run(kind, func(t *testing.T) {
			h, id, req := firstStartAtomicFixture(t)
			firstStartSeedHistory(t, h, kind, req.Envelope.Request.Stage.Envelope.Request.WorkspaceID, id)
			before := firstStartAtomicSnapshot(t, h)
			result, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
			require.Error(t, err, "late %s must invalidate the original no-history proof", kind)
			require.Nil(t, result)
			require.Equal(t, before, firstStartAtomicSnapshot(t, h))
			require.Zero(t, h.countRows("SELECT COUNT(*) FROM write_group_versions WHERE group_id=?", req.Envelope.Request.Stage.Current.GroupID))
		})
	}
}

func TestSuccessorFirstStartAtomicNoHistoryIncludesOldGenerationAndForeignRows(t *testing.T) {
	for _, foreign := range []bool{false, true} {
		for _, kind := range []string{"version", "runtime", "sample", "first_start_pending"} {
			t.Run(fmt.Sprintf("foreign_%t/%s", foreign, kind), func(t *testing.T) {
				h, id, req := firstStartAtomicFixture(t)
				scope := req.Envelope.Request.Stage.Envelope.Request
				head := req.Envelope.Request.Stage.Current.GroupID
				// Model a previously abandoned never-applied generation. Both generations
				// remain genuine tombstones and no applied/runtime history is discarded.
				_, err := h.db.Exec(`UPDATE write_groups SET status='deleted' WHERE id=?`, head)
				require.NoError(t, err)
				_, err = h.db.Exec(`UPDATE write_group_lineages SET current_group_id=?,reserved_group_id=NULL,reservation_operation_id=NULL,reservation_intent_digest=NULL,revision=revision+1 WHERE lineage_id=?`, head, scope.LineageID)
				require.NoError(t, err)
				check := func() error {
					tx, err := h.db.BeginTx(t.Context(), nil)
					require.NoError(t, err)
					defer tx.Rollback()
					return h.service.repo.requireSuccessorNoHistoryInTx(t.Context(), tx, scope.WorkspaceID, scope.LineageID, head)
				}
				require.NoError(t, check(), "the clean old-generation baseline must be valid")
				workspaceID := scope.WorkspaceID
				if foreign {
					workspaceID = "foreign-workspace"
				}
				firstStartSeedHistory(t, h, kind, workspaceID, id)
				before := firstStartAtomicSnapshot(t, h)
				require.Error(t, check(), "old generation G1 evidence cannot disappear behind G2 or a foreign workspace label")
				require.Equal(t, before, firstStartAtomicSnapshot(t, h))
			})
		}
	}
}

func firstStartHistoryRow(t *testing.T, h *circuitProfileHarness, table string, fields map[string]any) {
	t.Helper()
	columns := make([]string, 0, len(fields))
	for key := range fields {
		columns = append(columns, key)
	}
	sort.Strings(columns)
	args := make([]any, len(columns))
	placeholders := make([]string, len(columns))
	for i, key := range columns {
		args[i], placeholders[i] = fields[key], "?"
	}
	_, err := h.db.Exec("INSERT INTO "+table+" ("+strings.Join(columns, ",")+") VALUES ("+strings.Join(placeholders, ",")+")", args...)
	require.NoError(t, err, "seed contradictory %s evidence", table)
}

func firstStartSeedHistory(t *testing.T, h *circuitProfileHarness, kind, workspaceID, id string) {
	t.Helper()
	const at = "2026-10-28T10:00:00Z"
	row := map[string]any{"group_id": id, "group_revision": "late-revision"}
	table := ""
	switch kind {
	case "version":
		table = "write_group_versions"
		row["effective_at"], row["payload"] = at, "{}"
	case "runtime":
		table = "wg_runtime_versions"
		row["workspace_id"], row["payload"] = workspaceID, "{}"
	case "sample":
		table = "wg_delivery_samples"
		row["workspace_id"], row["sample_id"], row["member_key"] = workspaceID, "late-sample", "member"
		row["observed_at"], row["bucket_start"], row["created_at"] = at, at, at
		row["payload"], row["payload_digest"] = "{}", "digest"
	case "outbox_receipt":
		table = "wg_delivery_outbox"
		row["workspace_id"], row["effect_key"], row["record_id"] = workspaceID, "late-effect", "record"
		for _, key := range []string{"bucket_start", "next_retry_at", "created_at", "updated_at"} {
			row[key] = at
		}
		for _, key := range []string{"partition_key", "destination_scope", "connector_id", "connector_revision", "table_name", "payload_digest"} {
			row[key] = "late"
		}
		row["payload"] = "{}"
		firstStartHistoryRow(t, h, "wg_delivery_receipts", map[string]any{"effect_key": "late-effect", "payload_digest": "late", "committed_at": at})
	case "checkpoint":
		table = "wg_delivery_checkpoints"
		row["next_close"], row["updated_at"] = at, at
	case "bucket":
		table = "wg_delivery_buckets"
		row["entity_key"], row["bucket_start"], row["kind"], row["created_at"] = "entity", at, "no_data", at
	case "counter_baseline":
		table = "wg_delivery_counter_baselines"
		row["workspace_id"], row["member_key"], row["binding_key"] = workspaceID, "member", "binding"
		row["window_end"], row["created_at"], row["updated_at"], row["payload"], row["payload_digest"] = at, at, at, "{}", "digest"
	case "counter_proof":
		table = "wg_counter_proof_associations"
		row["workspace_id"], row["sample_seq"], row["proof_version"], row["association_version"] = workspaceID, 1, 1, 1
		for _, key := range []string{"sample_id", "sample_digest", "member_key", "proof_member_key", "proof_digest", "runtime_digest", "evidence_kind"} {
			row[key] = "late"
		}
		row["accepted_at"], row["associated_at"] = at, at
	case "revision_end":
		table = "write_group_revision_ends"
		row["workspace_id"], row["ends_at"], row["cutover_id"] = workspaceID, at, "late-cutover"
	case "fence":
		firstStartHistoryRow(t, h, "write_group_cutovers", map[string]any{"id": "late-fence-cutover", "workspace_id": workspaceID, "kind": "forward", "source_group_id": "unrelated-source", "boundary": at, "review_digest": "digest", "detail": `{"ended":[{"group_id":"unrelated-source","revision":"R1"}],"started":[{"group_id":"unrelated-target","revision":"R2"}]}`, "intent": "{}"})
		table = "write_group_ownership_fences"
		row["workspace_id"], row["fence_revision"], row["direction"], row["profile_kind"] = workspaceID, "late-fence-cutover", "ended", ""
		for _, key := range []string{"device_id", "point_id", "tag_id", "connector_id", "connector_revision", "database_name", "table_schema", "table_name"} {
			row[key] = "late"
		}
	case "handoff_source", "handoff_target":
		table = "wg_delivery_sample_handoffs"
		row = map[string]any{"transfer_id": "late-handoff", "workspace_id": workspaceID, "source_group_id": "other-source", "source_revision": "R1", "source_seq": 1, "target_group_id": "other-target", "target_revision": "R2", "origin_seq": 1, "state": "pending", "record": "{}"}
		if kind == "handoff_source" {
			row["source_group_id"] = id
		} else {
			row["target_group_id"] = id
		}
	case "first_start_pending", "first_start_unknown":
		table = "managed_schema_operations"
		status := "pending"
		if kind == "first_start_unknown" {
			status = "unknown"
		}
		row = map[string]any{"operation_id": "unresolved-first-start", "token": "unresolved-first-start-token", "workspace_id": workspaceID, "scope_key": "unresolved-first-start-scope", "owner": "owner", "action": "successor_first_start", "status": status, "created_at": at, "updated_at": at, "detail": fmt.Sprintf(`{"scope":{"predecessor_group_id":%q}}`, id), "payload_digest": "digest"}
	default:
		t.Fatalf("unhandled history kind %s", kind)
	}
	firstStartHistoryRow(t, h, table, row)
}

package workspace

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// This is the configuration boundary, not a UI copy endpoint or runtime Apply
// test. Seed one strictly encoded historical baseline for a real saved group;
// independent Create must neither copy it nor rewrite any of its stored bytes.
// The runtime durable tests separately exercise producing/consuming baselines.
func TestMixedIndependentCreatePreservesOldBaselineWithoutInheritance(t *testing.T) {
	ctx := t.Context()
	h := newCircuitProfileHarness(t, ":memory:")
	measurementID := seedMeasurement(t, h, h.primary(), "counter")
	_, err := h.db.Exec(`UPDATE measurement_definitions SET unit='kWh',series_epoch='epoch-1',numeric_encoding='float32',counter_policy='{}' WHERE id=?`, measurementID)
	require.NoError(t, err)
	draft := mixedDraft(h, WriteGroupMemberCalculationDelta, 15, "original_counter")
	draft.Members[0].MeasurementID = &measurementID
	original, err := h.service.Create(ctx, h.mutation(ctx, draft))
	require.NoError(t, err)
	old := original.Group
	definition, err := measurement.NewSQLRepository(h.db).GetByID(ctx, measurementID)
	require.NoError(t, err)
	proof, err := unscaledCounterBindingProof(old, old.Members[0], schema.DataTypeFloat32, definition)
	require.NoError(t, err)
	binding := aggregation.CounterBinding{
		WorkspaceID: old.WorkspaceID, GroupID: old.ID, ProfileRevision: old.Revision,
		MemberKey: proof.Member.SampleMemberKey(), MeasurementID: proof.MeasurementID,
		DefinitionRevision: proof.DefinitionRevision, SourceBindingRevision: proof.SourceBindingRevision,
		SeriesEpoch: proof.SeriesEpoch, Kind: aggregation.CounterKindMonotonic, Unit: proof.Unit,
	}
	end := time.Date(2026, 10, 1, 0, 0, 15, 0, time.UTC)
	value, err := measurement.NewFloat64(123.5)
	require.NoError(t, err)
	checkpoint := aggregation.CounterCheckpoint{Binding: binding, WindowEnd: end,
		Reading: aggregation.CounterReading{SampleID: "original-record-reading", ObservedAt: end.Add(-time.Second), Value: value}}
	payload, err := aggregation.EncodeCounterCheckpoint(checkpoint)
	require.NoError(t, err)
	decoded, err := aggregation.DecodeCounterCheckpoint(payload)
	require.NoError(t, err)
	require.Equal(t, checkpoint, decoded)
	bindingKey, err := binding.Key()
	require.NoError(t, err)
	digest := fmt.Sprintf("%x", sha256.Sum256(payload))
	stamp := end.Format("2006-01-02T15:04:05.000000000Z")
	_, err = h.db.Exec(`INSERT INTO wg_delivery_counter_baselines
  (workspace_id,group_id,group_revision,member_key,binding_key,window_end,payload,payload_digest,created_at,updated_at)
  VALUES (?,?,?,?,?,?,?,?,?,?)`, old.WorkspaceID, old.ID, old.Revision, binding.MemberKey, bindingKey, stamp, string(payload), digest, stamp, stamp)
	require.NoError(t, err)
	readOriginal := func() [10]string {
		var row [10]string
		require.NoError(t, h.db.QueryRow(`SELECT workspace_id,group_id,group_revision,member_key,binding_key,window_end,payload,payload_digest,created_at,updated_at
   FROM wg_delivery_counter_baselines WHERE group_id=? AND group_revision=?`, old.ID, old.Revision).
			Scan(&row[0], &row[1], &row[2], &row[3], &row[4], &row[5], &row[6], &row[7], &row[8], &row[9]))
		return row
	}
	before := readOriginal()
	lineage := func(id string) string {
		var value string
		require.NoError(t, h.db.QueryRow(`SELECT lineage_id FROM write_group_profile_keys WHERE group_id=? AND profile_kind='mixed' AND canonical_role='record'`, id).Scan(&value))
		require.NotEmpty(t, value)
		return value
	}
	oldLineage := lineage(old.ID)
	// Copying selected source metadata does not imply copying historical state.
	freshDraft := cloneWriteGroup(draft)
	freshDraft.RecordID = uuid.NewString()
	freshDraft.Destination.TableName = "independent_counter"
	created, err := h.service.Create(ctx, h.mutation(ctx, freshDraft))
	require.NoError(t, err)
	fresh := created.Group
	require.NotEqual(t, old.ID, fresh.ID)
	require.NotEqual(t, old.RecordID, fresh.RecordID)
	require.NotEqual(t, oldLineage, lineage(fresh.ID))
	require.Equal(t, old.Members[0].DeviceID, fresh.Members[0].DeviceID)
	require.Equal(t, old.Members[0].PointID, fresh.Members[0].PointID)
	require.Equal(t, old.Members[0].TagID, fresh.Members[0].TagID)
	require.Equal(t, old.Members[0].MeasurementID, fresh.Members[0].MeasurementID)
	require.Equal(t, before, readOriginal(), "all ten columns, including exact payload/digest, must survive independent creation unchanged")
	var count int
	require.NoError(t, h.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_counter_baselines WHERE group_id=?`, fresh.ID).Scan(&count))
	require.Zero(t, count, "no revision of the new record may inherit the old baseline")
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM wg_delivery_counter_baselines`))
	got, err := h.service.Get(ctx, old.ID)
	require.NoError(t, err)
	require.Equal(t, old, got.Group)
	require.Empty(t, fresh.AppliedRevision, "creating a draft is not activating a runtime")
	require.Zero(t, h.countRows(`SELECT COUNT(*) FROM wg_runtime_versions`))
	require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_cutovers`))
}

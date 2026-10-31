package workspace

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func mixedDraft(h *circuitProfileHarness, method WriteGroupMemberCalculation, seconds int64, table string) *WriteGroup {
	group := h.draft(WriteGroupProfileInstantaneous, "", "Mixed", h.primary())
	group.RecordID = uuid.NewString()
	group.RowPolicy.IntervalSeconds = 0
	group.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1}
	group.Destination.TableName = table
	group.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
	group.Members[0].CircuitName = ""
	group.Members[0].Calculation = method
	group.Members[0].Cadence = &WriteGroupMemberCadence{Interval: seconds, Unit: WriteGroupIntervalUnitSecond, Alignment: WriteGroupAlignmentUTCEpoch}
	return group
}

func TestMixedDraftCRUDRecordIdentityAndRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mixed.db")
	h := newCircuitProfileHarness(t, path)
	draft := mixedDraft(h, WriteGroupMemberCalculationLatest, 5, "record_a")
	request := h.mutation(t.Context(), draft)
	saved, err := h.service.Create(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, draft.RecordID, saved.Group.RecordID)
	replay, err := h.service.Create(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, saved.Group.ID, replay.Group.ID)
	second, err := h.service.Create(t.Context(), h.mutation(t.Context(), mixedDraft(h, WriteGroupMemberCalculationMean, 600, "record_b")))
	require.NoError(t, err)
	require.NotEqual(t, saved.Group.ID, second.Group.ID)
	conflicting := mixedDraft(h, WriteGroupMemberCalculationLatest, 15, "record_a")
	_, err = h.service.Create(t.Context(), h.mutation(t.Context(), conflicting))
	require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
	modified := cloneWriteGroup(saved.Group)
	modified.Members[0].Cadence.Interval = 15
	update := h.mutation(t.Context(), modified)
	update.ExpectedGroupRevision = saved.Group.Revision
	changed, err := h.service.Update(t.Context(), saved.Group.ID, update)
	require.NoError(t, err)
	_, err = h.service.Update(t.Context(), saved.Group.ID, update)
	require.Error(t, err)
	require.NoError(t, h.db.Close())
	db := openWorkspaceTestDB(t, path)
	restarted := newCircuitProfileHarnessOn(t, db, h)
	got, err := restarted.service.Get(t.Context(), saved.Group.ID)
	require.NoError(t, err)
	require.Equal(t, changed.Group, got.Group)
	require.Equal(t, int64(15), got.Group.Members[0].Cadence.Interval)
	require.Zero(t, got.Group.RowPolicy.IntervalSeconds)
	var policy string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT row_policy FROM write_groups WHERE id=?`, got.Group.ID).Scan(&policy))
	require.NotContains(t, policy, "interval_seconds")
	wire, err := json.Marshal(got.Group)
	require.NoError(t, err)
	require.NotContains(t, string(wire), "interval_seconds")
	_, err = NewSQLWriteGroupRepository(db).Get(t.Context(), "foreign-workspace", got.Group.ID)
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
}

func TestMixedDraftCounterConfirmationAndForeignIdentity(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	latest := mixedDraft(h, WriteGroupMemberCalculationLatest, 5, "counter_latest")
	_, err := h.service.Create(t.Context(), h.mutation(t.Context(), latest))
	require.NoError(t, err)
	delta := mixedDraft(h, WriteGroupMemberCalculationDelta, 600, "counter_delta")
	_, err = h.service.Create(t.Context(), h.mutation(t.Context(), delta))
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	id := seedMeasurement(t, h, h.primary(), "counter")
	delta.Members[0].MeasurementID = &id
	_, err = h.service.Create(t.Context(), h.mutation(t.Context(), delta))
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	_, err = h.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET unit='kWh',series_epoch='epoch-1',numeric_encoding='float32',counter_policy='{}' WHERE id=?`, id)
	require.NoError(t, err)
	saved, err := h.service.Create(t.Context(), h.mutation(t.Context(), delta))
	require.NoError(t, err)
	require.Equal(t, WriteGroupMemberCalculationDelta, saved.Group.Members[0].Calculation)
	require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_group_lineages`))
	foreign := cloneWriteGroup(saved.Group)
	foreign.RecordID = latest.RecordID
	request := h.mutation(t.Context(), foreign)
	request.ExpectedGroupRevision = saved.Group.Revision
	_, err = h.service.Update(t.Context(), saved.Group.ID, request)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	request = h.mutation(t.Context(), delta)
	request.WorkspaceID = "foreign-workspace"
	_, err = h.service.Create(t.Context(), request)
	require.Error(t, err)
	require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_groups`))
}

func TestMixedDraftCorruptCadenceAndReplayChangesRejected(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	draft := mixedDraft(h, WriteGroupMemberCalculationLatest, 5, "one")
	saved, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.NoError(t, err)
	for _, change := range []func(*WriteGroup){
		func(g *WriteGroup) { g.Members[0].Calculation = WriteGroupMemberCalculationMean },
		func(g *WriteGroup) { g.Members[0].Cadence.Interval = 15 },
		func(g *WriteGroup) { g.Members[0].SourceRevision = "stale" },
		func(g *WriteGroup) { g.Members[0].MappingRevision = "stale" },
	} {
		candidate := cloneWriteGroup(draft)
		change(candidate)
		_, err := h.service.Create(t.Context(), h.mutation(t.Context(), candidate))
		require.Error(t, err)
	}
	for _, raw := range []any{nil, `{}`, `{"interval":1,"unit":"hour","alignment":"utc_epoch"}`, `{"interval":5,"unit":"second","alignment":"utc_epoch","unknown":1}`} {
		_, err := h.db.ExecContext(t.Context(), `UPDATE write_group_members SET cadence=? WHERE group_id=?`, raw, saved.Group.ID)
		require.NoError(t, err)
		_, err = h.service.Get(t.Context(), saved.Group.ID)
		require.Error(t, err)
	}
}

func TestMixedDraftConcurrentCreateConverges(t *testing.T) {
	h := newCircuitProfileHarness(t, filepath.Join(t.TempDir(), "concurrent.db"))
	draft := mixedDraft(h, WriteGroupMemberCalculationLatest, 15, "shared_intent")
	request := h.mutation(t.Context(), draft)
	type outcome struct {
		result *WriteGroupSaveResult
		err    error
	}
	results := make(chan outcome, 4)
	for range 4 {
		go func() {
			svc := NewWriteGroupService(NewService(NewSQLRepository(h.db)), NewSQLWriteGroupRepository(h.db))
			result, err := svc.Create(t.Context(), request)
			results <- outcome{result, err}
		}()
	}
	ids := map[string]bool{}
	for range 4 {
		out := <-results
		require.NoError(t, out.err)
		ids[out.result.Group.ID] = true
	}
	require.Len(t, ids, 1)
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_profile_keys`))
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_lineages`))
}

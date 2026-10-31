package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedDraftDeleteRetainsCanonicalIdentity(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	draft := mixedDraft(h, WriteGroupMemberCalculationLatest, 5, "retained")
	created, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.NoError(t, err)
	request := h.mutation(t.Context(), created.Group)
	request.ExpectedGroupRevision = created.Group.Revision
	removed, err := h.service.Delete(t.Context(), created.Group.ID, request)
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDeleted, removed.Group.Status)
	require.Equal(t, draft.RecordID, removed.Group.RecordID)
	_, err = h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_profile_keys`))
}

func TestMixedDraftLongSharesOnlyExistingStorageScope(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	first := mixedDraft(h, WriteGroupMemberCalculationLatest, 5, "")
	first.LayoutVersion = WriteGroupLayoutManagedLongV1
	a, err := h.service.Create(t.Context(), h.mutation(t.Context(), first))
	require.NoError(t, err)
	second := mixedDraft(h, WriteGroupMemberCalculationMean, 60, "")
	second.LayoutVersion = WriteGroupLayoutManagedLongV1
	b, err := h.service.Create(t.Context(), h.mutation(t.Context(), second))
	require.NoError(t, err)
	require.Equal(t, a.Group.Destination.TableName, b.Group.Destination.TableName)
	require.NotEqual(t, a.Group.RecordID, b.Group.RecordID)
	collision := mixedDraft(h, WriteGroupMemberCalculationLatest, 15, a.Group.Destination.TableName)
	_, err = h.service.Create(t.Context(), h.mutation(t.Context(), collision))
	require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
}

func TestMixedDraftPreservesLegacyTargetAndSourceGuards(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	mixed := mixedDraft(h, WriteGroupMemberCalculationLatest, 5, "owned_table")
	_, err := h.service.Create(t.Context(), h.mutation(t.Context(), mixed))
	require.NoError(t, err)
	legacy := h.draft(WriteGroupProfileInstantaneous, "legacy-circuit", "Legacy", h.primary())
	legacy.Destination.TableName = "owned_table"
	_, err = h.ensure(legacy)
	require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
	legacy.Destination.TableName = "legacy_table"
	saved, err := h.ensure(legacy)
	require.NoError(t, err)
	copied := cloneWriteGroup(legacy)
	copied.Members[0].CircuitID = "different-circuit"
	copied.Destination.TableName = "different_table"
	_, err = h.ensure(copied)
	require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
	stale := mixedDraft(h, WriteGroupMemberCalculationLatest, 15, "stale")
	stale.Members[0].SourceRevision = "stale-source-revision"
	_, err = h.service.Create(t.Context(), h.mutation(t.Context(), stale))
	require.ErrorIs(t, err, ErrWriteGroupSourceRevisionConflict)
	require.NotEmpty(t, saved.Group.ID)
}

func TestMixedDraftDoesNotSilentlyConvertCustomStorage(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	draft := mixedDraft(h, WriteGroupMemberCalculationLatest, 5, "custom_table")
	draft.Destination.StorageStrategy = WriteGroupStorageStrategyCustom
	_, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_groups`))
}

func TestMixedMeanRejectsCounterSemanticsAndStaleType(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	id := seedConfirmedCounter(t, h, h.primary())
	draft := mixedDraft(h, WriteGroupMemberCalculationMean, 15, "mean_counter")
	draft.Members[0].MeasurementID = &id
	_, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	draft.Members[0].MeasurementID = nil
	draft.Members[0].ValueRole = WriteGroupValueRoleState
	_, err = h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	draft.Members[0].ValueRole = ""
	draft.Members[0].SourceValueType = "bool"
	_, err = h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_groups`))
}

func TestMixedCanonicalReadRejectsStoredGroupWindow(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	draft := mixedDraft(h, WriteGroupMemberCalculationLatest, 5, "no_group_window")
	created, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.NoError(t, err)
	_, err = h.db.ExecContext(t.Context(), `UPDATE write_groups SET row_policy=json_set(row_policy,'$.interval_seconds',60) WHERE id=?`, created.Group.ID)
	require.NoError(t, err)
	_, err = h.service.Get(t.Context(), created.Group.ID)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
}

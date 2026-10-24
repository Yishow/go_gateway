package workspace

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func flag(value bool) *bool { return &value }

// TestLongStorageRecreateMissingPolicy pins how the table-recovery policy is
// chosen and kept: a new managed long group saves it explicitly on, it can be
// explicitly turned off, it never exists on other layouts, an edit that does
// not mention it keeps what was saved, and nothing saved before the policy
// existed ever reads as authorized.
func TestLongStorageRecreateMissingPolicy(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	power := h.primary()
	energy := h.seedSource(power.deviceID, "energy")
	extra := h.seedSource(power.deviceID, "extra")

	t.Run("a new managed long group saves the policy explicitly on", func(t *testing.T) {
		created := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "main", power))
		require.NotNil(t, created.WritePolicy.RecreateMissing, "the default is saved as a value, not left implicit")
		require.True(t, *created.WritePolicy.RecreateMissing)
		saved, err := h.service.Get(t.Context(), created.ID)
		require.NoError(t, err)
		require.Equal(t, flag(true), saved.Group.WritePolicy.RecreateMissing, "and it survives a reload")
	})

	t.Run("an explicit off is kept", func(t *testing.T) {
		draft := h.longDraft(WriteGroupProfileCumulative, "off", energy)
		draft.WritePolicy.RecreateMissing = flag(false)
		created := h.mustEnsure(draft)
		require.Equal(t, flag(false), created.WritePolicy.RecreateMissing)
	})

	t.Run("a resend of the same intent is not a conflict, an opposite explicit policy is", func(t *testing.T) {
		again := h.longDraft(WriteGroupProfileInstantaneous, "main", power)
		first := h.mustEnsure(again)
		require.Equal(t, flag(true), first.WritePolicy.RecreateMissing, "a resend that omits the policy resolves to the saved group")

		opposite := h.longDraft(WriteGroupProfileInstantaneous, "main", power)
		opposite.WritePolicy.RecreateMissing = flag(false)
		_, err := h.ensure(opposite)
		require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict, "a different explicit policy is a different intent")
	})

	t.Run("only managed long groups carry a policy", func(t *testing.T) {
		wide := newBasicManagedDraft(h.fixture)
		wide.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
		created, err := h.service.Create(t.Context(), h.mutation(t.Context(), wide))
		require.NoError(t, err)
		require.Nil(t, created.Group.WritePolicy.RecreateMissing, "a wide group is not given a recovery policy")

		withPolicy := newBasicManagedDraft(h.fixture)
		withPolicy.Name = "wide with policy"
		withPolicy.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
		withPolicy.WritePolicy.RecreateMissing = flag(true)
		_, err = h.service.Create(t.Context(), h.mutation(t.Context(), withPolicy))
		require.ErrorIs(t, err, ErrWriteGroupValidation, "the policy only governs managed long storage")
	})

	t.Run("an edit keeps the policy it does not mention, and changing it needs a new confirmation", func(t *testing.T) {
		group := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "edit", extra))
		saved, err := h.service.Get(t.Context(), group.ID)
		require.NoError(t, err)
		proof := strings.Repeat("a", 64)
		confirmed, err := h.service.ConfirmManagedSchema(t.Context(), group.ID, h.mutationFor(saved), "op-1", proof)
		require.NoError(t, err)
		require.Equal(t, "op-1", confirmed.Group.Destination.SchemaRevision)

		omit := cloneWriteGroup(confirmed.Group)
		omit.WritePolicy.RecreateMissing = nil
		omit.Name = "renamed"
		kept, err := h.service.Update(t.Context(), group.ID, h.updateMutation(omit, confirmed))
		require.NoError(t, err)
		require.Equal(t, flag(true), kept.Group.WritePolicy.RecreateMissing, "an edit that leaves the policy out keeps the saved one")
		require.Equal(t, "op-1", kept.Group.Destination.SchemaRevision, "and the group's confirmation stands")

		flip := cloneWriteGroup(kept.Group)
		flip.WritePolicy.RecreateMissing = flag(false)
		changed, err := h.service.Update(t.Context(), group.ID, h.updateMutation(flip, kept))
		require.NoError(t, err)
		require.Equal(t, flag(false), changed.Group.WritePolicy.RecreateMissing)
		require.Empty(t, changed.Group.Destination.SchemaRevision, "a changed policy is shown and confirmed again")
		require.Empty(t, changed.Group.Destination.SchemaDigest)
		require.NotEqual(t, kept.Group.Revision, changed.Group.Revision)
	})

	t.Run("a group saved before the policy existed never reads as authorized", func(t *testing.T) {
		legacy := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "legacy", h.seedSource(power.deviceID, "legacy")))
		raw, err := json.Marshal(WriteGroupWritePolicy{Mode: "append", DedupeCapability: "receipt"})
		require.NoError(t, err)
		require.NotContains(t, string(raw), "recreate_missing", "the legacy payload is exactly what an older build stored")
		_, err = h.db.ExecContext(t.Context(), `UPDATE write_groups SET write_policy = ? WHERE id = ?`, string(raw), legacy.ID)
		require.NoError(t, err)

		loaded, err := h.service.Get(t.Context(), legacy.ID)
		require.NoError(t, err)
		require.Nil(t, loaded.Group.WritePolicy.RecreateMissing, "an absent policy is no authority")

		edit := cloneWriteGroup(loaded.Group)
		edit.Name = "legacy renamed"
		updated, err := h.service.Update(t.Context(), legacy.ID, h.updateMutation(edit, loaded))
		require.NoError(t, err)
		require.Nil(t, updated.Group.WritePolicy.RecreateMissing, "an edit does not quietly grant what was never confirmed")
	})
}

func (h *circuitProfileHarness) mutationFor(saved *WriteGroupSaveResult) WriteGroupMutation {
	return WriteGroupMutation{
		WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision,
		ExpectedGroupRevision: saved.Group.Revision, ExpectedConnectorRevision: saved.Group.Destination.ConnectorRevision,
	}
}

func (h *circuitProfileHarness) updateMutation(edit *WriteGroup, current *WriteGroupSaveResult) WriteGroupMutation {
	mutation := h.mutationFor(current)
	mutation.Group = edit
	return mutation
}

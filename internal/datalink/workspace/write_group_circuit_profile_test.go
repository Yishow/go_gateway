package workspace

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCircuitProfileCanonicalIdentity(t *testing.T) {
	t.Run("power and energy on one circuit keep independent stable groups across reload", func(t *testing.T) {
		dbPath := filepath.Join(t.TempDir(), "circuit-identity.db")
		h := newCircuitProfileHarness(t, dbPath)
		power := h.primary()
		energy := h.seedSource(power.deviceID, "energy")

		instantaneous := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-main", "Main feeder", power))
		cumulative := h.mustEnsure(h.draft(WriteGroupProfileCumulative, "circuit-main", "Main feeder", energy))

		require.NotEqual(t, instantaneous.ID, cumulative.ID)
		require.NotEqual(t, instantaneous.Destination.TableName, cumulative.Destination.TableName)
		require.Equal(t, WriteGroupProfileInstantaneous, instantaneous.RecordingProfile.Kind)
		require.Equal(t, WriteGroupProfileCumulative, cumulative.RecordingProfile.Kind)
		require.Equal(t, 15, instantaneous.RecordingProfile.Interval)
		require.Equal(t, WriteGroupIntervalUnitMinute, cumulative.RecordingProfile.Unit)
		require.Equal(t, WriteGroupLayoutManagedWideV1, instantaneous.EffectiveLayoutVersion())
		require.Equal(t, "circuit-main", instantaneous.Members[0].CircuitID)
		require.Equal(t, "Main feeder", instantaneous.Members[0].CircuitName)

		require.NoError(t, h.db.Close())
		restarted := newCircuitProfileHarnessOn(t, openWorkspaceTestDB(t, dbPath), h)
		listed, err := restarted.service.List(t.Context())
		require.NoError(t, err)
		require.Len(t, listed.Groups, 2)
		byID := map[string]*WriteGroup{}
		for _, group := range listed.Groups {
			byID[group.ID] = group
		}
		require.Equal(t, WriteGroupProfileInstantaneous, byID[instantaneous.ID].RecordingProfile.Kind)
		require.Equal(t, "circuit-main", byID[instantaneous.ID].Members[0].CircuitID)
		require.Equal(t, WriteGroupProfileCumulative, byID[cumulative.ID].RecordingProfile.Kind)
		require.Equal(t, "circuit-main", byID[cumulative.ID].Members[0].CircuitID)

		replayed := restarted.mustEnsure(restarted.draft(WriteGroupProfileInstantaneous, "circuit-main", "Main feeder", power))
		require.Equal(t, instantaneous.ID, replayed.ID)
	})

	t.Run("resent create reuses the canonical group without a second writer", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-main", "Main feeder", h.primary())
		stale := h.mutation(t.Context(), draft)

		first, err := h.service.EnsureCircuitProfile(t.Context(), stale)
		require.NoError(t, err)
		// The resend keeps the pre-create workspace revision, as a lost response would.
		second, err := h.service.EnsureCircuitProfile(t.Context(), stale)
		require.NoError(t, err)

		require.Equal(t, first.Group.ID, second.Group.ID)
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_profile_keys`))
	})

	t.Run("renaming the circuit keeps group identity and later creates reuse it", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		created := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-main", "Main feeder", h.primary()))

		renamed := cloneWriteGroup(created)
		renamed.Name = "Renamed feeder instantaneous"
		renamed.Members[0].CircuitName = "Renamed feeder"
		mutation := h.mutation(t.Context(), renamed)
		mutation.ExpectedGroupRevision = created.Revision
		updated, err := h.service.Update(t.Context(), created.ID, mutation)
		require.NoError(t, err)
		require.Equal(t, created.ID, updated.Group.ID)
		require.Equal(t, "circuit-main", updated.Group.Members[0].CircuitID)
		require.Equal(t, "Renamed feeder", updated.Group.Members[0].CircuitName)

		again := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-main", "Renamed feeder", h.primary()))
		require.Equal(t, created.ID, again.ID)
		stillOld := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-main", "Main feeder", h.primary()))
		require.Equal(t, created.ID, stillOld.ID)
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
	})

	t.Run("one device with two circuits keeps memberships separate", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		first := h.primary()
		second := h.seedSource(first.deviceID, "line-two")
		unselected := h.seedSource(first.deviceID, "line-three")

		a := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-a", "Line A", first))
		b := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-b", "Line B", second))

		require.NotEqual(t, a.ID, b.ID)
		require.Equal(t, "circuit-a", a.Members[0].CircuitID)
		require.Equal(t, "circuit-b", b.Members[0].CircuitID)
		require.Equal(t, 0, h.countRows(`SELECT COUNT(*) FROM write_group_members WHERE tag_id = ?`, unselected.tagID))
	})

	t.Run("unassigned raw selection is keyed by device and never by a circuit", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		raw := h.primary()
		other := h.seedSource(raw.deviceID, "other")

		unassigned := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "", "", raw))
		replayed := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "", "", raw))
		require.Equal(t, unassigned.ID, replayed.ID)
		require.Empty(t, unassigned.Members[0].CircuitID)
		require.Equal(t, 1, h.countRows(
			`SELECT COUNT(*) FROM write_group_profile_keys WHERE circuit_id = '' AND selection_scope = ?`, raw.deviceID))

		// A circuit whose opaque ID equals the device ID is a different key.
		lookalike := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, raw.deviceID, "Looks like the device", other))
		require.NotEqual(t, unassigned.ID, lookalike.ID)
		require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_group_profile_keys`))
	})

	t.Run("a resent create with different sources is a conflict, not a duplicate", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		first := h.primary()
		second := h.seedSource(first.deviceID, "second")
		h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-main", "Main", first))

		_, err := h.ensure(h.draft(WriteGroupProfileInstantaneous, "circuit-main", "Main", first, second))
		require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
	})
}

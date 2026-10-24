package workspace

import (
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCircuitProfileRejectsInvalidIdentity(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	first := h.primary()
	second := h.seedSource(first.deviceID, "second")
	otherDevice := seedAdditionalWriteGroupDevice(t.Context(), t, h.db, h.fixture.workspaceID, h.fixture.connectorID)
	foreign := circuitSource{deviceID: otherDevice.deviceID, pointID: otherDevice.pointID, tagID: otherDevice.tagID}

	build := func(mutate func(*WriteGroup), sources ...circuitSource) *WriteGroup {
		group := h.draft(WriteGroupProfileInstantaneous, "circuit-a", "A", sources...)
		mutate(group)
		return group
	}
	cases := map[string]*WriteGroup{
		"members from different circuits": build(func(g *WriteGroup) { g.Members[1].CircuitID = "circuit-b" }, first, second),
		"one member without a circuit":    build(func(g *WriteGroup) { g.Members[1].CircuitID, g.Members[1].CircuitName = "", "" }, first, second),
		"circuit name without an id":      build(func(g *WriteGroup) { g.Members[0].CircuitID = "" }, first),
		"cumulative without a circuit": build(func(g *WriteGroup) {
			g.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileCumulative, Interval: 10, Unit: WriteGroupIntervalUnitMinute}
			g.Members[0].CircuitID, g.Members[0].CircuitName = "", ""
		}, first),
		"unassigned across two devices": build(func(g *WriteGroup) {
			for i := range g.Members {
				g.Members[i].CircuitID, g.Members[i].CircuitName = "", ""
			}
		}, first, foreign),
		"unknown kind":             build(func(g *WriteGroup) { g.RecordingProfile.Kind = "average" }, first),
		"instantaneous in minutes": build(func(g *WriteGroup) { g.RecordingProfile.Unit = WriteGroupIntervalUnitMinute }, first),
		"cumulative in seconds": build(func(g *WriteGroup) {
			g.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileCumulative, Interval: 10, Unit: WriteGroupIntervalUnitSecond}
		}, first),
		"zero interval":               build(func(g *WriteGroup) { g.RecordingProfile.Interval = 0 }, first),
		"negative interval":           build(func(g *WriteGroup) { g.RecordingProfile.Interval = -15 }, first),
		"unsupported layout":          build(func(g *WriteGroup) { g.LayoutVersion = "managed_future_v9" }, first),
		"control character in id":     build(func(g *WriteGroup) { g.Members[0].CircuitID = "circuit\x00a" }, first),
		"oversized circuit id":        build(func(g *WriteGroup) { g.Members[0].CircuitID = strings.Repeat("c", maxCircuitIdentityLength+1) }, first),
		"no recording profile at all": build(func(g *WriteGroup) { g.RecordingProfile = nil }, first),
	}
	for name, group := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := h.ensure(group)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_groups`))
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_profile_keys`))
		})
	}

	t.Run("plain create cannot smuggle circuit identity into a non-profile or custom group", func(t *testing.T) {
		noProfile := h.draft(WriteGroupProfileInstantaneous, "circuit-a", "A", first)
		noProfile.RecordingProfile = nil
		_, err := h.service.Create(t.Context(), h.mutation(t.Context(), noProfile))
		require.ErrorIs(t, err, ErrWriteGroupValidation)

		custom := newBasicWriteGroup(h.fixture)
		custom.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond}
		_, err = h.service.Create(t.Context(), h.mutation(t.Context(), custom))
		require.ErrorIs(t, err, ErrWriteGroupValidation)
		require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_groups`))
	})
}

func TestCircuitProfileDuplicateOwnership(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	source := h.primary()
	owner := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-a", "A", source))

	t.Run("the same profile kind in another circuit is blocked", func(t *testing.T) {
		_, err := h.ensure(h.draft(WriteGroupProfileInstantaneous, "circuit-b", "B", source))
		require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
	})

	t.Run("an unassigned selection cannot shadow a circuit-owned source", func(t *testing.T) {
		_, err := h.ensure(h.draft(WriteGroupProfileInstantaneous, "", "", source))
		require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
	})

	t.Run("the other profile kind keeps a distinct output identity", func(t *testing.T) {
		cumulative := h.mustEnsure(h.draft(WriteGroupProfileCumulative, "circuit-a", "A", source))
		require.NotEqual(t, owner.ID, cumulative.ID)
		require.NotEqual(t, owner.Destination.TableName, cumulative.Destination.TableName)
	})

	t.Run("another destination is a separate ownership scope", func(t *testing.T) {
		connectorID := uuid.NewString()
		_, err := h.db.ExecContext(t.Context(), `
			INSERT INTO database_connectors (id, name, kind, connection_config, identity_revision, status, enabled)
			VALUES (?, 'Second target', 'sqlite', '{"database":"other-db","schema":"other-schema"}', 'connector-revision-1', 'ready', 1)
		`, connectorID)
		require.NoError(t, err)
		elsewhere := h.draft(WriteGroupProfileInstantaneous, "circuit-b", "B", source)
		elsewhere.Destination.ConnectorID = connectorID
		created := h.mustEnsure(elsewhere)
		require.NotEqual(t, owner.ID, created.ID)
	})

	t.Run("a tombstoned owner releases the source", func(t *testing.T) {
		deleted, err := h.service.Delete(t.Context(), owner.ID, WriteGroupMutation{
			WorkspaceID: h.fixture.workspaceID, ExpectedWorkspaceRevision: h.mutation(t.Context(), nil).ExpectedWorkspaceRevision,
			ExpectedGroupRevision: owner.Revision, ExpectedConnectorRevision: "connector-revision-1",
		})
		require.NoError(t, err)
		require.Equal(t, WriteGroupStatusDeleted, deleted.Group.Status)

		// The create-once key is retained, so the retired circuit is not reusable.
		_, err = h.ensure(h.draft(WriteGroupProfileInstantaneous, "circuit-a", "A", source))
		require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
		moved := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-c", "C", source))
		require.NotEqual(t, owner.ID, moved.ID)
	})
}

func TestCircuitProfileIdentityIsFixedAfterCreate(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	created := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-a", "A", h.primary()))

	update := func(mutate func(*WriteGroup)) (*WriteGroupSaveResult, error) {
		next := cloneWriteGroup(created)
		mutate(next)
		mutation := h.mutation(t.Context(), next)
		mutation.ExpectedGroupRevision = created.Revision
		return h.service.Update(t.Context(), created.ID, mutation)
	}

	rejected := map[string]func(*WriteGroup){
		"moving the group to another circuit": func(g *WriteGroup) { g.Members[0].CircuitID = "circuit-b" },
		"unassigning the group":               func(g *WriteGroup) { g.Members[0].CircuitID, g.Members[0].CircuitName = "", "" },
		"changing the profile kind": func(g *WriteGroup) {
			g.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileCumulative, Interval: 10, Unit: WriteGroupIntervalUnitMinute}
		},
		"dropping the profile": func(g *WriteGroup) {
			g.RecordingProfile = nil
			g.Members[0].CircuitID, g.Members[0].CircuitName = "", ""
		},
		"switching layout": func(g *WriteGroup) { g.LayoutVersion = "managed_long_v1" },
		"unknown layout":   func(g *WriteGroup) { g.LayoutVersion = "unknown" },
	}
	for name, mutate := range rejected {
		t.Run(name, func(t *testing.T) {
			_, err := update(mutate)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
		})
	}

	t.Run("an interval edit keeps the key and persists", func(t *testing.T) {
		updated, err := update(func(g *WriteGroup) { g.RecordingProfile.Interval = 30 })
		require.NoError(t, err)
		require.Equal(t, created.ID, updated.Group.ID)
		reloaded, err := h.service.Get(t.Context(), created.ID)
		require.NoError(t, err)
		require.Equal(t, 30, reloaded.Group.RecordingProfile.Interval)
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_profile_keys`))
	})

	t.Run("a legacy group cannot be given a profile by update", func(t *testing.T) {
		legacy := newBasicManagedDraft(h.fixture)
		legacy.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
		extra := h.seedSource(h.fixture.deviceID, "legacy-source")
		legacy.Members = []WriteGroupMember{{DeviceID: extra.deviceID, PointID: extra.pointID, TagID: extra.tagID, Required: true}}
		made, err := h.service.Create(t.Context(), h.mutation(t.Context(), legacy))
		require.NoError(t, err)

		promoted := cloneWriteGroup(made.Group)
		promoted.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond}
		promoted.Members[0].CircuitID = "circuit-z"
		mutation := h.mutation(t.Context(), promoted)
		mutation.ExpectedGroupRevision = made.Group.Revision
		_, err = h.service.Update(t.Context(), made.Group.ID, mutation)
		require.ErrorIs(t, err, ErrWriteGroupValidation)
		require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_profile_keys WHERE group_id = ?`, made.Group.ID))
	})

	t.Run("plain create registers the key instead of bypassing it", func(t *testing.T) {
		other := h.seedSource(h.fixture.deviceID, "bypass")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-bypass", "Bypass", other)
		draft.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
		first, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
		require.NoError(t, err)
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_profile_keys WHERE group_id = ?`, first.Group.ID))

		_, err = h.service.Create(t.Context(), h.mutation(t.Context(), draft))
		require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups WHERE id = ? OR name = ?`, first.Group.ID, draft.Name))
	})
}

func TestCircuitProfileConcurrentEnsureConverges(t *testing.T) {
	h := newCircuitProfileHarness(t, t.TempDir()+"/concurrent-circuit.db")
	draft := h.draft(WriteGroupProfileInstantaneous, "circuit-a", "A", h.primary())
	mutation := h.mutation(t.Context(), draft)

	const callers = 4
	type outcome struct {
		result *WriteGroupSaveResult
		err    error
	}
	results := make(chan outcome, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Go(func() {
			result, err := h.service.EnsureCircuitProfile(t.Context(), mutation)
			results <- outcome{result: result, err: err}
		})
	}
	wg.Wait()
	close(results)

	ids := map[string]struct{}{}
	for got := range results {
		require.NoError(t, got.err)
		ids[got.result.Group.ID] = struct{}{}
	}
	require.Len(t, ids, 1)
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
}

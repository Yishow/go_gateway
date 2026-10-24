package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func (h *circuitProfileHarness) longDraft(kind WriteGroupProfileKind, circuitID string, sources ...circuitSource) *WriteGroup {
	group := h.draft(kind, circuitID, "Circuit "+circuitID, sources...)
	group.LayoutVersion = WriteGroupLayoutManagedLongV1
	return group
}

func TestLongStorageScopeIdentity(t *testing.T) {
	base := longStorageScope{
		WorkspaceID: "ws", ConnectorID: "conn", Database: "db", TableSchema: "main",
		Mode: WriteGroupStorageModePlain, Layout: WriteGroupLayoutManagedLongV1,
	}
	require.Regexp(t, `^gw_long_[0-9a-f]{16}$`, base.tableName())
	require.Regexp(t, `^_gw_owner_[0-9a-f]{24}$`, base.ownerColumn())
	require.Equal(t, base.tableName(), base.tableName(), "deterministic")

	variants := map[string]longStorageScope{}
	for name, mutate := range map[string]func(*longStorageScope){
		"workspace": func(s *longStorageScope) { s.WorkspaceID = "other" },
		"connector": func(s *longStorageScope) { s.ConnectorID = "other" },
		"database":  func(s *longStorageScope) { s.Database = "other" },
		"schema":    func(s *longStorageScope) { s.TableSchema = "other" },
		"mode":      func(s *longStorageScope) { s.Mode = "timescale" },
		"layout":    func(s *longStorageScope) { s.Layout = "managed_long_v2" },
	} {
		scope := base
		mutate(&scope)
		variants[name] = scope
	}
	seenTables, seenOwners := map[string]string{base.tableName(): "base"}, map[string]string{base.ownerColumn(): "base"}
	for name, scope := range variants {
		require.NotContains(t, seenTables, scope.tableName(), "%s must own another table", name)
		require.NotContains(t, seenOwners, scope.ownerColumn(), "%s must own another marker", name)
		seenTables[scope.tableName()], seenOwners[scope.ownerColumn()] = name, name
	}
	// A field boundary is part of the identity: moving text between fields changes it.
	shifted := base
	shifted.Database, shifted.TableSchema = "dbmain", ""
	require.NotEqual(t, base.tableName(), shifted.tableName())
}

func TestLongStorageGroupsShareOneScopeTable(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	power := h.primary()
	energy := h.seedSource(power.deviceID, "energy")

	first := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "main", power))
	second := h.mustEnsure(h.longDraft(WriteGroupProfileCumulative, "main", energy))

	require.Equal(t, WriteGroupLayoutManagedLongV1, first.EffectiveLayoutVersion())
	require.Equal(t, WriteGroupStorageModePlain, first.StorageMode, "plain is explicit once saved")
	require.Equal(t, longStorageScopeOf(first).tableName(), first.Destination.TableName)
	require.Equal(t, first.Destination.TableName, second.Destination.TableName)
	require.NotEqual(t, first.ID, second.ID)
	require.NotEqual(t, managedGroupTableName(first.ID), first.Destination.TableName, "the table is not named by the group")
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_storage_scopes`), "one frozen scope for both profiles")
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_storage_scopes
		WHERE table_name = ? AND owner_column = ? AND generation = 1 AND length(layout_digest) = 64 AND storage_mode = 'plain'`,
		first.Destination.TableName, longStorageScopeOf(first).ownerColumn()))

	t.Run("a resent create reuses the group and its shared table", func(t *testing.T) {
		again := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "main", power))
		require.Equal(t, first.ID, again.ID)
		require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_groups`))
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_storage_scopes`))
	})

	t.Run("a caller cannot choose the shared table", func(t *testing.T) {
		other := h.seedSource(power.deviceID, "other")
		draft := h.longDraft(WriteGroupProfileInstantaneous, "elsewhere", other)
		draft.Destination.TableName = "attacker_table"
		_, err := h.ensure(draft)
		require.ErrorIs(t, err, ErrWriteGroupValidation, "a foreign table name is refused, not silently replaced")
		require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_groups`))

		draft.Destination.TableName = ""
		created := h.mustEnsure(draft)
		require.Equal(t, first.Destination.TableName, created.Destination.TableName)
	})

	t.Run("an update cannot move a group onto another table", func(t *testing.T) {
		saved, err := h.service.Get(t.Context(), first.ID)
		require.NoError(t, err)
		edit := cloneWriteGroup(saved.Group)
		edit.Destination.TableName = "attacker_table"
		mutation := h.mutation(t.Context(), edit)
		mutation.ExpectedGroupRevision = saved.Group.Revision
		_, err = h.service.Update(t.Context(), first.ID, mutation)
		require.ErrorIs(t, err, ErrWriteGroupValidation)

		edit.Destination.TableName = "" // the editor may omit it; the server keeps its own
		mutation = h.mutation(t.Context(), edit)
		mutation.ExpectedGroupRevision = saved.Group.Revision
		updated, err := h.service.Update(t.Context(), first.ID, mutation)
		require.NoError(t, err)
		require.Equal(t, first.Destination.TableName, updated.Group.Destination.TableName)
	})

	t.Run("wide groups create no storage scope", func(t *testing.T) {
		wide := newBasicManagedDraft(h.fixture)
		wide.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
		before := h.countRows(`SELECT COUNT(*) FROM write_group_storage_scopes`)
		created, err := h.service.Create(t.Context(), h.mutation(t.Context(), wide))
		require.NoError(t, err)
		require.Equal(t, WriteGroupLayoutManagedWideV1, created.Group.EffectiveLayoutVersion())
		require.Empty(t, created.Group.StorageMode)
		require.Equal(t, before, h.countRows(`SELECT COUNT(*) FROM write_group_storage_scopes`))
	})
}

func TestLongStorageRejectsInvalidShapes(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	source := h.primary()
	cases := map[string]func(*WriteGroup){
		"timescale is not available yet":     func(g *WriteGroup) { g.StorageMode = "timescale" },
		"an unknown storage mode":            func(g *WriteGroup) { g.StorageMode = "columnar" },
		"a long group needs a profile":       func(g *WriteGroup) { g.RecordingProfile = nil },
		"a long group needs a managed table": func(g *WriteGroup) { g.Destination.StorageStrategy = WriteGroupStorageStrategyCustom },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			draft := h.longDraft(WriteGroupProfileInstantaneous, "main", source)
			mutate(draft)
			_, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_groups`))
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_storage_scopes`))
		})
	}

	t.Run("a wide group cannot carry a storage mode", func(t *testing.T) {
		wide := h.draft(WriteGroupProfileInstantaneous, "main", "Main", source)
		wide.StorageMode = WriteGroupStorageModePlain
		_, err := h.ensure(wide)
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})

	t.Run("the storage mode is fixed after create", func(t *testing.T) {
		created := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "main", source))
		edit := cloneWriteGroup(created)
		edit.StorageMode = "timescale"
		mutation := h.mutation(t.Context(), edit)
		mutation.ExpectedGroupRevision = created.Revision
		_, err := h.service.Update(t.Context(), created.ID, mutation)
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})
}

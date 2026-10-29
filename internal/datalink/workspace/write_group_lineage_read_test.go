package workspace

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanonicalLineageCorruptionCannotFallThroughToCreate(t *testing.T) {
	cases := map[string]string{
		"missing_generation":        "DELETE FROM write_group_generations",
		"missing_lineage":           "DELETE FROM write_group_lineages",
		"missing_group":             "DELETE FROM write_groups",
		"null_key_reference":        "UPDATE KEY_TABLE SET lineage_id=NULL",
		"foreign_key_reference":     "UPDATE KEY_TABLE SET lineage_id='foreign-lineage'",
		"foreign_key_workspace":     "UPDATE KEY_TABLE SET workspace_id='foreign-workspace'",
		"foreign_lineage_workspace": "UPDATE write_group_lineages SET workspace_id='foreign-workspace'",
		"foreign_group_workspace":   "UPDATE write_groups SET workspace_id='foreign-workspace'",
		"missing_current":           "UPDATE write_group_lineages SET current_group_id='missing-head'",
		"empty_reference":           "UPDATE KEY_TABLE SET lineage_id=' '",
		"partial_reservation":       "UPDATE write_group_lineages SET reservation_operation_id='O2'",
		"missing_reserved":          "UPDATE write_group_lineages SET reserved_group_id='G2',reservation_operation_id='O2',reservation_intent_digest='intent-2'",
		"bad_genesis":               "UPDATE write_group_generations SET generation=2,predecessor_group_id='G0'",
	}
	for _, kind := range []string{"basic", "profile"} {
		for name, corruption := range cases {
			t.Run(kind+"/"+name, func(t *testing.T) {
				h := newCircuitProfileHarness(t, ":memory:")
				h.db.SetMaxOpenConns(1)
				mutation := h.mutation(t.Context(), lineageTestDraft(h, kind))
				created, err := ensureLineageKind(h, kind, mutation)
				require.NoError(t, err)
				seedLineageReadFixture(t, h, created.Group)
				current, err := h.workspace.GetOrCreate(t.Context())
				require.NoError(t, err)
				mutation.ExpectedWorkspaceRevision = current.DatabaseSetupRevision
				_, err = h.db.Exec("PRAGMA foreign_keys=OFF; PRAGMA ignore_check_constraints=ON")
				require.NoError(t, err)
				if name == "foreign_key_workspace" {
					// Keep the requested alias present and poison the other family.
					other := "write_group_profile_keys"
					if kind == "profile" {
						other = "write_group_basic_keys"
						_, err = h.db.Exec(`INSERT INTO write_group_basic_keys(workspace_id,device_id,canonical_role,group_id,lineage_id) SELECT workspace_id,'other-device','basic',current_group_id,lineage_id FROM write_group_lineages`)
					} else {
						_, err = h.db.Exec(`INSERT INTO write_group_profile_keys(workspace_id,circuit_id,profile_kind,canonical_role,group_id,lineage_id) SELECT workspace_id,'circuit-1','instantaneous','circuit',current_group_id,lineage_id FROM write_group_lineages`)
					}
					require.NoError(t, err)
					_, err = h.db.Exec(replaceLineageKeyTable(corruption, other))
				} else {
					_, err = h.db.Exec(replaceLineageKeyTable(corruption, lineageKeyTables(kind)[0]))
				}
				require.NoError(t, err)
				before := lineageReadSnapshot(t, h)
				_, err = ensureLineageKind(h, kind, mutation)
				require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
				require.ErrorContains(t, err, "write-group lineage is corrupt")
				require.False(t, errors.Is(err, ErrWriteGroupNotFound), "corruption must not invite create")
				require.False(t, errors.Is(err, ErrNotFound), "corruption must not invite create")
				require.Equal(t, before, lineageReadSnapshot(t, h), "read/replay must never repair metadata")
			})
		}
	}
}

func TestCanonicalLineageCurrentAndHistoricalIdentity(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	mutation := h.mutation(t.Context(), lineageTestDraft(h, "dual"))
	first, err := ensureLineageKind(h, "dual", mutation)
	require.NoError(t, err)
	lineage := seedLineageReadFixture(t, h, first.Group)
	genesis := cloneWriteGroup(first.Group)
	current := seedLineageGeneration(t, h, genesis, lineage, "G2", 2)
	seedLineageGeneration(t, h, current, lineage, "G3", 3)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET current_group_id='G2',reserved_group_id='G3',reservation_operation_id='O3',reservation_intent_digest='intent-3',revision=7 WHERE lineage_id=?`, lineage)
	require.NoError(t, err)
	before := lineageReadSnapshot(t, h)
	key, err := circuitProfileKeyOf(genesis)
	require.NoError(t, err)
	profile, err := h.service.repo.findProfileGroup(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, "G2", profile.ID)
	basic, err := h.service.repo.findBasicManaged(t.Context(), genesis.WorkspaceID, h.fixture.deviceID)
	require.NoError(t, err)
	require.Equal(t, "G2", basic.ID)
	for _, id := range []string{genesis.ID, "G2", "G3"} {
		got, err := h.service.repo.Get(t.Context(), genesis.WorkspaceID, id)
		require.NoError(t, err)
		require.Equal(t, h.fixture.deviceID, got.BasicManagedDeviceID)
		stored, found, err := h.service.repo.loadProfileKey(t.Context(), h.db, genesis.WorkspaceID, id)
		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, key, stored)
	}
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestCanonicalLineageOnlyAbsentKeysAreNotFound(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	_, err := h.service.repo.findBasicManaged(t.Context(), h.fixture.workspaceID, "missing-device")
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	key := circuitProfileKey{WorkspaceID: h.fixture.workspaceID, CircuitID: "missing-circuit", Kind: WriteGroupProfileInstantaneous, Role: "circuit"}
	_, err = h.service.repo.findProfileGroup(t.Context(), key)
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	group := newBasicWriteGroup(h.fixture)
	require.NoError(t, h.service.repo.Create(t.Context(), group))
	_, found, err := h.service.repo.loadProfileKey(t.Context(), h.db, group.WorkspaceID, group.ID)
	require.NoError(t, err)
	require.False(t, found)
	managed := newBasicWriteGroup(h.fixture)
	managed.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
	require.NoError(t, h.service.repo.Create(t.Context(), managed))
	loaded, err := h.service.repo.Get(t.Context(), managed.WorkspaceID, managed.ID)
	require.NoError(t, err)
	require.Empty(t, loaded.BasicManagedDeviceID)
	require.Zero(t, h.countRows("SELECT COUNT(*) FROM write_group_lineages"))
}

func seedLineageReadFixture(t *testing.T, h *circuitProfileHarness, group *WriteGroup) string {
	t.Helper()
	var lineage string
	err := h.db.QueryRow("SELECT lineage_id FROM write_group_generations WHERE group_id=?", group.ID).Scan(&lineage)
	if errors.Is(err, sql.ErrNoRows) {
		lineage = "L1"
		_, err = h.db.Exec(`INSERT INTO write_group_lineages(lineage_id,workspace_id,current_group_id) VALUES(?,?,?)`, lineage, group.WorkspaceID, group.ID)
		require.NoError(t, err)
		_, err = h.db.Exec(`INSERT INTO write_group_generations(group_id,lineage_id,generation) VALUES(?,?,1)`, group.ID, lineage)
		require.NoError(t, err)
		for _, table := range lineageKeyTables("dual") {
			_, err = h.db.Exec("UPDATE "+table+" SET lineage_id=? WHERE group_id=?", lineage, group.ID)
			require.NoError(t, err)
		}
	} else {
		require.NoError(t, err)
	}
	return lineage
}

func seedLineageGeneration(t *testing.T, h *circuitProfileHarness, predecessor *WriteGroup, lineage, id string, generation int) *WriteGroup {
	t.Helper()
	group := cloneWriteGroup(predecessor)
	group.ID, group.Revision = id, "R"+id
	group.BasicManagedDeviceID = ""
	require.NoError(t, h.service.repo.insertGroup(t.Context(), h.db, group))
	_, err := h.db.Exec(`INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id) VALUES(?,?,?,?)`, id, lineage, generation, predecessor.ID)
	require.NoError(t, err)
	return group
}

func replaceLineageKeyTable(query, table string) string {
	return strings.ReplaceAll(query, "KEY_TABLE", table)
}

func lineageReadSnapshot(t *testing.T, h *circuitProfileHarness) map[string]string {
	t.Helper()
	state := lifecycleReadState(t, h.db)
	for _, table := range []string{"write_group_lineages", "write_group_generations", "write_group_basic_keys", "write_group_profile_keys"} {
		rows, err := h.db.Query("SELECT * FROM " + table + " ORDER BY rowid")
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		all := [][]any{}
		for rows.Next() {
			values, pointers := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			require.NoError(t, rows.Scan(pointers...))
			all = append(all, values)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		raw, err := json.Marshal(all)
		require.NoError(t, err)
		state[table] = string(raw)
	}
	return state
}

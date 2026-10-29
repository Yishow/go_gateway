package workspace

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func TestCanonicalLineageInitializerRefusesRepair(t *testing.T) {
	for _, corruption := range []string{"DELETE FROM write_group_generations", "DELETE FROM write_group_lineages", "DELETE FROM write_group_generations; DELETE FROM write_group_lineages", "UPDATE write_group_basic_keys SET lineage_id=NULL"} {
		t.Run(corruption, func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			h.db.SetMaxOpenConns(1)
			first, err := ensureLineageKind(h, "basic", h.mutation(t.Context(), lineageTestDraft(h, "basic")))
			require.NoError(t, err)
			_, err = h.db.Exec("PRAGMA foreign_keys=OFF; " + corruption)
			require.NoError(t, err)
			before := lineageReadSnapshot(t, h)
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			tx, err := h.db.BeginTx(ctx, nil)
			require.NoError(t, err)
			_, err = h.service.repo.initializeCanonicalLineageInTx(ctx, tx, first.Group.WorkspaceID, first.Group.ID)
			require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
			require.NotErrorIs(t, err, ErrWriteGroupNotFound)
			require.NotErrorIs(t, err, ErrNotFound)
			require.NoError(t, tx.Commit(), "a refused initializer has made no repair even if caller commits")
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestCanonicalLineageMissingSchemaFailsClosed(t *testing.T) {
	for _, kind := range []string{"basic", "profile"} {
		for _, table := range []string{"write_group_lineages", "write_group_generations"} {
			t.Run(kind+"/"+table, func(t *testing.T) {
				h := newCircuitProfileHarness(t, ":memory:")
				h.db.SetMaxOpenConns(1)
				mutation := h.mutation(t.Context(), lineageTestDraft(h, kind))
				_, err := ensureLineageKind(h, kind, mutation)
				require.NoError(t, err)
				_, err = h.db.Exec("PRAGMA foreign_keys=OFF; DROP TABLE " + table)
				require.NoError(t, err)
				before := lifecycleReadState(t, h.db)
				_, err = ensureLineageKind(h, kind, mutation)
				require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
				require.NotErrorIs(t, err, ErrWriteGroupNotFound)
				require.NotErrorIs(t, err, ErrNotFound)
				require.Equal(t, before, lifecycleReadState(t, h.db))
			})
		}
	}
}

func TestCanonicalLineageGenerationGuardsAndCallerSnapshot(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	first, err := ensureLineageKind(h, "dual", h.mutation(t.Context(), lineageTestDraft(h, "dual")))
	require.NoError(t, err)
	lineage := requireLineageGenesis(t, h.db, first.Group)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	tx, err := h.db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()
	// The sole connection is held by tx; a nested repository DB read deadlocks.
	group := cloneWriteGroup(first.Group)
	group.ID, group.Revision = "G2", "R2"
	require.NoError(t, h.service.repo.insertGroup(ctx, tx, group))
	_, err = tx.ExecContext(ctx, `INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id) VALUES('G2',?,2,?)`, lineage, first.Group.ID)
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, `UPDATE write_group_lineages SET current_group_id='G2',revision=5 WHERE lineage_id=?`, lineage)
	require.NoError(t, err)
	metadata, found, err := h.service.repo.loadWriteGroupLineage(ctx, tx, first.Group.WorkspaceID, "G2")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, first.Group.ID, metadata.GenesisGroupID)
	require.Equal(t, "G2", metadata.CurrentGroupID)
	require.EqualValues(t, 2, metadata.Generation)
	require.EqualValues(t, 5, metadata.Revision)
	require.Equal(t, sql.NullString{String: first.Group.ID, Valid: true}, metadata.PredecessorGroupID)
	_, err = h.service.repo.initializeCanonicalLineageInTx(ctx, tx, first.Group.WorkspaceID, first.Group.ID)
	require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt, "initializer may not reinterpret a later lineage as fresh")
	_, _, err = h.service.repo.loadProfileKey(ctx, tx, "foreign-workspace", "G2")
	require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
	_, err = tx.ExecContext(ctx, `UPDATE write_group_generations SET predecessor_group_id='missing-predecessor' WHERE group_id='G2'`)
	require.NoError(t, err)
	_, _, err = h.service.repo.loadWriteGroupLineage(ctx, tx, first.Group.WorkspaceID, "G2")
	require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
	require.NoError(t, tx.Rollback())
	require.Equal(t, lineage, requireLineageGenesis(t, h.db, first.Group))
}

func TestCanonicalLineageDistinctGenerationEffectsPreserveReceipts(t *testing.T) {
	db, service, old := lifecyclePreservationFixture(t)
	seedLifecycleAcceptedData(t, db, old)
	before := lifecyclePreservedData(t, db)
	oldConfig, err := service.repo.Get(t.Context(), old.WorkspaceID, "G1")
	require.NoError(t, err)
	// Seed future metadata only: no live successor API is enabled by task 2.2.
	_, err = db.Exec(`INSERT INTO write_group_lineages(lineage_id,workspace_id,current_group_id) VALUES('L1',?,'G1')`, old.WorkspaceID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_generations(group_id,lineage_id,generation) VALUES('G1','L1',1)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_basic_keys(workspace_id,device_id,canonical_role,group_id,lineage_id) VALUES(?,?,'basic','G1','L1')`, old.WorkspaceID, old.Members[0].DeviceID)
	require.NoError(t, err)
	newGroup := cloneWriteGroup(old)
	newGroup.ID, newGroup.Revision, newGroup.AppliedRevision, newGroup.Status = "G2", "R2", "", WriteGroupStatusDraft
	require.NoError(t, service.repo.insertGroup(t.Context(), db, newGroup))
	_, err = db.Exec(`INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id) VALUES('G2','L1',2,'G1')`)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE write_group_lineages SET current_group_id='G2',revision=2 WHERE lineage_id='L1'`)
	require.NoError(t, err)
	when := time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)
	for _, profile := range []bool{false, true} {
		var effects []string
		for _, id := range []string{"G1", "G2"} {
			var record string
			if profile {
				record, err = snapshot.ProfileRecordID(old.WorkspaceID, id, "R1", "member-1", "instantaneous", when)
			} else {
				record, err = snapshot.RecordID(old.WorkspaceID, id, "R1", "", when)
			}
			require.NoError(t, err)
			effect, err := snapshot.EffectKey("original-destination", record)
			require.NoError(t, err)
			effects = append(effects, effect)
		}
		require.NotEqual(t, effects[0], effects[1], "group ID alone separates generation effect/receipt identities")
	}
	require.Equal(t, before, lifecyclePreservedData(t, db), "versions, accepted payloads, outbox and receipts are immutable")
	reloaded, err := service.repo.Get(t.Context(), old.WorkspaceID, "G1")
	require.NoError(t, err)
	require.Equal(t, oldConfig, reloaded)
}

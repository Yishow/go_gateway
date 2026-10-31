package workspace

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func mixedActiveIdentityPair(t *testing.T) (*circuitProfileHarness, *WriteGroup, *WriteGroup) {
	t.Helper()
	now := time.Date(2026, 10, 31, 12, 0, 7, 0, time.UTC)
	h, active := mixedHistoryFixture(t, &now)
	draft := mixedDraft(h, WriteGroupMemberCalculationMean, 60, "independent_mixed")
	draft.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
	draft.Destination.Database = active.Destination.Database
	created, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.NoError(t, err)
	return h, active, created.Group
}

func TestMixedActiveIdentityLoadsIndependentCanonicalRecords(t *testing.T) {
	h, active, candidate := mixedActiveIdentityPair(t)
	require.NotEqual(t, active.RecordID, candidate.RecordID)
	require.NotEqual(t, active.Destination.TableName, candidate.Destination.TableName)
	require.Equal(t, memberRefOf(active.Members[0]), memberRefOf(candidate.Members[0]))
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	first, err := h.service.repo.mixedRecordIdentityInTx(t.Context(), tx, active)
	require.NoError(t, err)
	second, err := h.service.repo.mixedRecordIdentityInTx(t.Context(), tx, candidate)
	require.NoError(t, err)
	require.NotEqual(t, first.lineageID, second.lineageID)
	overlap, err := h.service.repo.activeOverlapInTx(t.Context(), tx, candidate)
	require.NoError(t, err)
	require.Empty(t, overlap)
}

func TestMixedActiveIdentityRejectsUnprovenExceptions(t *testing.T) {
	cases := []string{"same record", "same lineage", "foreign workspace", "missing candidate key", "missing active key", "stale applied revision", "foreign applied revision", "foreign snapshot scope", "same target"}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			h, active, candidate := mixedActiveIdentityPair(t)
			// Mutations deliberately corrupt authoritative stored evidence, or submit
			// a foreign candidate. The real transactional loader must not grant reuse.
			switch name {
			case "same record":
				candidate.RecordID = active.RecordID
			case "foreign workspace":
				candidate.WorkspaceID = "foreign-workspace"
			case "same target":
				candidate.Destination.TableName = active.Destination.TableName
			case "missing candidate key", "missing active key":
				id := candidate.ID
				if name == "missing active key" {
					id = active.ID
				}
				_, err := h.db.Exec(`DELETE FROM write_group_profile_keys WHERE group_id=?`, id)
				require.NoError(t, err)
			case "same lineage":
				_, err := h.db.Exec(`UPDATE write_group_generations SET lineage_id=(SELECT lineage_id FROM write_group_generations WHERE group_id=?), generation=2, predecessor_group_id=? WHERE group_id=?`, active.ID, active.ID, candidate.ID)
				require.NoError(t, err)
			case "stale applied revision":
				_, err := h.db.Exec(`UPDATE write_groups SET applied_revision='absent-applied-revision' WHERE id=?`, active.ID)
				require.NoError(t, err)
			case "foreign applied revision", "foreign snapshot scope":
				saved := cloneWriteGroup(active)
				if name == "foreign applied revision" {
					saved.AppliedRevision = "foreign-applied-revision"
				} else {
					saved.WorkspaceID = "foreign-workspace"
				}
				payload, err := json.Marshal(saved)
				require.NoError(t, err)
				_, err = h.db.Exec(`UPDATE write_group_versions SET payload=? WHERE group_id=? AND group_revision=?`, string(payload), active.ID, active.AppliedRevision)
				require.NoError(t, err)
			}
			tx, err := h.db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			defer tx.Rollback()
			overlap, err := h.service.repo.activeOverlapInTx(t.Context(), tx, candidate)
			require.True(t, err != nil || len(overlap) > 0, "%s must not be granted independent-source reuse", name)
			if name == "stale applied revision" {
				require.ErrorIs(t, err, ErrWriteGroupAppliedRevisionUnavailable)
			}
			if name == "same target" {
				require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
			}
		})
	}
}

func TestMixedActiveIdentityDoesNotExemptLegacySameKind(t *testing.T) {
	h := newCircuitProfileHarness(t, filepath.Join(t.TempDir(), "workspace.db"))
	// The same real schema-ledger helper also confirms the legacy managed-wide
	// shape; this exercises its applied snapshot loader rather than pure identity.
	target := filepath.Join(t.TempDir(), "legacy-target.db")
	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	_, err = db.Exec(`PRAGMA user_version=0`)
	require.NoError(t, err)
	require.NoError(t, db.Close())
	setSQLiteConnectorConfig(t, h.db, h.fixture.connectorID, target)
	draft := h.draft(WriteGroupProfileInstantaneous, "", "", h.primary())
	draft.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
	draft.Destination.Database = target
	draft.Members[0].ValueRole = WriteGroupValueRoleGauge
	active := h.mustEnsure(draft)
	active = mixedHistoryConfirm(t, h, active)
	applied, err := h.service.Apply(t.Context(), active.ID, mixedHistoryRequest(t, h, active))
	require.NoError(t, err)
	candidate := cloneWriteGroup(applied.Group)
	candidate.ID = uuid.NewString()
	candidate.Destination.TableName = "another_legacy_target"
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	identity, err := h.service.repo.mixedRecordIdentityInTx(t.Context(), tx, candidate)
	require.NoError(t, err)
	require.Nil(t, identity)
	overlap, err := h.service.repo.activeOverlapInTx(t.Context(), tx, candidate)
	require.NoError(t, err)
	require.Contains(t, overlap, memberRefOf(active.Members[0]))
}

package workspace

import (
	"encoding/json"
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func recreateRepositorySnapshot(t *testing.T, h *circuitProfileHarness) map[string]string {
	t.Helper()
	state := lineageReadSnapshot(t, h)
	for _, table := range []string{"write_group_members", "write_group_storage_scopes", "write_group_ownership_fences"} {
		rows, err := h.db.Query("SELECT * FROM " + table + " ORDER BY rowid")
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		var all [][]any
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

func TestRecreateRepositoryRejectsChangedAuthorityWithoutEffects(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*testing.T, *circuitProfileHarness, *writeGroupRecreateReservationRequest)
		want   error
	}{
		{"stale_lineage", func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.ExpectedLineageRevision++
		}, ErrWriteGroupRevisionConflict},
		{"foreign_lineage", func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.LineageID = "foreign"
		}, ErrWriteGroupRevisionConflict},
		{"stale_group", func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.ExpectedPredecessorRevision = "stale"
		}, ErrWriteGroupRevisionConflict},
		{"stale_saved_connector", func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.ExpectedPredecessorConnectorRevision = "stale"
		}, ErrWriteGroupConnectorRevisionConflict},
		{"stale_target_connector", func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.ExpectedTargetConnectorRevision = "stale"
		}, ErrWriteGroupConnectorRevisionConflict},
		{"foreign_workspace", func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.WorkspaceID = "foreign"
			r.Group.WorkspaceID = "foreign"
		}, ErrWriteGroupLineageCorrupt},
		{"revision_overflow", func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			_, err := h.db.Exec("UPDATE write_group_lineages SET revision=?", int64(math.MaxInt64))
			require.NoError(t, err)
			r.ExpectedLineageRevision = math.MaxInt64
		}, ErrWriteGroupLineageCorrupt},
		{"future_generation", func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			seedLineageGeneration(t, h, r.Group, r.LineageID, "unexplained-G2", 2)
		}, ErrWriteGroupLineageCorrupt},
		{"historical_predecessor", func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			seedLineageGeneration(t, h, r.Group, r.LineageID, "current-G2", 2)
			_, err := h.db.Exec("UPDATE write_group_lineages SET current_group_id='current-G2'")
			require.NoError(t, err)
		}, ErrWriteGroupRevisionConflict},
	}
	for _, status := range []WriteGroupStatus{WriteGroupStatusDraft, WriteGroupStatusReady, WriteGroupStatusRunning, WriteGroupStatusDisabled} {
		cases = append(cases, struct {
			name   string
			mutate func(*testing.T, *circuitProfileHarness, *writeGroupRecreateReservationRequest)
			want   error
		}{string(status), func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			_, err := h.db.Exec("UPDATE write_groups SET status=?", status)
			require.NoError(t, err)
		}, ErrWriteGroupLifecycleBlocked})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			tc.mutate(t, h, &req)
			before := recreateRepositorySnapshot(t, h)
			_, err := reserveRepositoryTest(t, h, req)
			require.ErrorIs(t, err, tc.want)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func TestRecreateRepositoryCorruptionNeverRepairs(t *testing.T) {
	for name, query := range map[string]string{
		"missing_generation":    "DELETE FROM write_group_generations",
		"missing_lineage":       "DELETE FROM write_group_lineages",
		"missing_head":          "UPDATE write_group_lineages SET current_group_id='missing'",
		"missing_profile_alias": "DELETE FROM write_group_profile_keys",
		"foreign_alias":         "UPDATE write_group_profile_keys SET lineage_id='foreign'",
		"partial_reservation":   "UPDATE write_group_lineages SET reservation_operation_id='O2'",
		"generation_overflow":   "UPDATE write_group_generations SET generation=9223372036854775807",
		"basic_role_changed":    "UPDATE write_group_basic_keys SET canonical_role='other'",
		"basic_device_changed":  "UPDATE write_group_basic_keys SET device_id='other'",
	} {
		t.Run(name, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "dual")
			_, err := h.db.Exec("PRAGMA foreign_keys=OFF; PRAGMA ignore_check_constraints=ON")
			require.NoError(t, err)
			_, err = h.db.Exec(query)
			require.NoError(t, err)
			before := recreateRepositorySnapshot(t, h)
			_, err = reserveRepositoryTest(t, h, req)
			require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
			require.False(t, errors.Is(err, ErrWriteGroupNotFound))
			require.False(t, errors.Is(err, ErrNotFound))
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func TestRecreateRepositoryTransactionFailuresRollbackEveryArtifact(t *testing.T) {
	for _, point := range []struct{ name, table, event string }{
		{"group", "write_groups", "INSERT"}, {"member", "write_group_members", "INSERT"},
		{"generation", "write_group_generations", "INSERT"}, {"reservation", "write_group_lineages", "UPDATE"},
		{"storage", "write_group_storage_scopes", "INSERT"},
	} {
		t.Run(point.name, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "dual")
			if point.name == "storage" {
				recreateRepositoryNewTarget(t, h, &req)
			}
			_, err := h.db.Exec("CREATE TRIGGER abort_recreate AFTER " + point.event + " ON " + point.table + " BEGIN SELECT RAISE(ABORT,'injected successor failure'); END")
			require.NoError(t, err)
			before := recreateRepositorySnapshot(t, h)
			_, err = reserveRepositoryTest(t, h, req)
			require.ErrorContains(t, err, "injected successor failure")
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
	t.Run("reservation_CAS_loser", func(t *testing.T) {
		h, req := recreateRepositoryFixture(t, "profile")
		_, err := h.db.Exec(`CREATE TRIGGER compete_recreate AFTER INSERT ON write_group_generations WHEN NEW.generation=2
   BEGIN UPDATE write_group_lineages SET revision=revision+1 WHERE lineage_id=NEW.lineage_id; END`)
		require.NoError(t, err)
		before := recreateRepositorySnapshot(t, h)
		_, err = reserveRepositoryTest(t, h, req)
		require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
		require.Equal(t, before, recreateRepositorySnapshot(t, h))
	})
}

func recreateRepositoryNewTarget(t *testing.T, h *circuitProfileHarness, req *writeGroupRecreateReservationRequest) {
	t.Helper()
	_, err := h.db.Exec(`INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled)
 VALUES('recreate-target','recreate target','sqlite','{"database":"new-db","schema":"main"}','target-rev','ready',1)`)
	require.NoError(t, err)
	req.Group.Destination.ConnectorID, req.Group.Destination.ConnectorRevision = "recreate-target", "target-rev"
	req.ExpectedTargetConnectorRevision = "target-rev"
	req.Group.Destination.Database, req.Group.Destination.TableSchema, req.Group.Destination.TableName = "", "", ""
}

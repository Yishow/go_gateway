package workspace

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
)

func seedRecreateEngineEvidence(t *testing.T, h *circuitProfileHarness, req writeGroupRecreateReservationRequest) *WriteGroupAppliedSnapshot {
	t.Helper()
	group := cloneWriteGroup(req.Group)
	group.AppliedRevision = group.Revision
	group.Status = WriteGroupStatusReady
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, h.service.repo.insertVersionInTx(t.Context(), tx, group, time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)))
	require.NoError(t, tx.Commit())
	snapshot := &WriteGroupAppliedSnapshot{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision,
		AppliedRevision: group.Revision, Group: group, RuntimeLayout: &WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialect("sqlite")}}
	_, err = h.db.Exec("UPDATE write_groups SET applied_revision=? WHERE id=?", group.Revision, group.ID)
	require.NoError(t, err)
	saveRecreateEngineEvidence(t, h, req, snapshot)
	return snapshot
}

func saveRecreateEngineEvidence(t *testing.T, h *circuitProfileHarness, req writeGroupRecreateReservationRequest, snapshot *WriteGroupAppliedSnapshot) {
	t.Helper()
	payload, err := json.Marshal(snapshot)
	require.NoError(t, err)
	_, err = h.db.Exec(`INSERT OR REPLACE INTO wg_runtime_versions(workspace_id,group_id,group_revision,payload) VALUES(?,?,?,?)`, req.WorkspaceID, req.PredecessorGroupID, req.ExpectedPredecessorRevision, string(payload))
	require.NoError(t, err)
}

func TestRecreateRepositoryFrozenEngineSurvivesRemovedSourceConnector(t *testing.T) {
	h, req := recreateRepositoryFixture(t, "profile")
	seedRecreateEngineEvidence(t, h, req)
	recreateRepositoryNewTarget(t, h, &req)
	_, err := h.db.Exec("DELETE FROM database_connectors WHERE id=?", h.fixture.connectorID)
	require.NoError(t, err)
	result, err := reserveRepositoryTest(t, h, req)
	require.NoError(t, err)
	require.Equal(t, "recreate-target", result.Group.Destination.ConnectorID)
	require.Empty(t, result.Group.AppliedRevision)
}

func TestRecreateRepositoryInvalidFrozenEngineCannotUseLiveFallback(t *testing.T) {
	for name, mutate := range map[string]func(*WriteGroupAppliedSnapshot){
		"workspace":          func(s *WriteGroupAppliedSnapshot) { s.WorkspaceID = "foreign" },
		"group":              func(s *WriteGroupAppliedSnapshot) { s.GroupID = "foreign" },
		"revision":           func(s *WriteGroupAppliedSnapshot) { s.GroupRevision = "foreign" },
		"applied":            func(s *WriteGroupAppliedSnapshot) { s.AppliedRevision = "foreign" },
		"nested_workspace":   func(s *WriteGroupAppliedSnapshot) { s.Group.WorkspaceID = "foreign" },
		"nested_group":       func(s *WriteGroupAppliedSnapshot) { s.Group.ID = "foreign" },
		"nested_revision":    func(s *WriteGroupAppliedSnapshot) { s.Group.Revision = "foreign" },
		"nested_applied":     func(s *WriteGroupAppliedSnapshot) { s.Group.AppliedRevision = "foreign" },
		"missing_group":      func(s *WriteGroupAppliedSnapshot) { s.Group = nil },
		"missing_layout":     func(s *WriteGroupAppliedSnapshot) { s.RuntimeLayout = nil },
		"unknown_engine":     func(s *WriteGroupAppliedSnapshot) { s.RuntimeLayout.Dialect = "unknown" },
		"ambiguous_mysql":    func(s *WriteGroupAppliedSnapshot) { s.RuntimeLayout.Dialect = "mysql" },
		"connector":          func(s *WriteGroupAppliedSnapshot) { s.Group.Destination.ConnectorID = "foreign" },
		"connector_revision": func(s *WriteGroupAppliedSnapshot) { s.Group.Destination.ConnectorRevision = "foreign" },
		"database":           func(s *WriteGroupAppliedSnapshot) { s.Group.Destination.Database = "foreign" },
		"table":              func(s *WriteGroupAppliedSnapshot) { s.Group.Destination.TableName = "foreign" },
	} {
		t.Run(name, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			snapshot := seedRecreateEngineEvidence(t, h, req)
			mutate(snapshot)
			saveRecreateEngineEvidence(t, h, req, snapshot)
			before := recreateRepositorySnapshot(t, h)
			_, err := reserveRepositoryTest(t, h, req)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			var unavailable *WriteGroupTransitionError
			require.ErrorAs(t, err, &unavailable)
			require.Equal(t, "source-engine-evidence-unavailable", unavailable.ReasonCode)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func TestRecreateRepositoryEngineEvidenceStorageScopeAndLegacyFallback(t *testing.T) {
	for _, mode := range []string{"legacy_pinned", "missing_both", "invalid_json", "foreign_storage_scope", "duplicate_storage_scope"} {
		t.Run(mode, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			seedRecreateEngineEvidence(t, h, req)
			var err error
			switch mode {
			case "legacy_pinned", "missing_both":
				_, err = h.db.Exec("DELETE FROM wg_runtime_versions")
				require.NoError(t, err)
				if mode == "missing_both" {
					recreateRepositoryNewTarget(t, h, &req)
					_, err = h.db.Exec("DELETE FROM database_connectors WHERE id=?", h.fixture.connectorID)
				}
			case "invalid_json":
				_, err = h.db.Exec("UPDATE wg_runtime_versions SET payload='invalid'")
			case "foreign_storage_scope":
				_, err = h.db.Exec("UPDATE wg_runtime_versions SET workspace_id='foreign'")
			case "duplicate_storage_scope":
				_, err = h.db.Exec(`INSERT INTO wg_runtime_versions SELECT 'foreign',group_id,group_revision,payload FROM wg_runtime_versions`)
			}
			require.NoError(t, err)
			before := recreateRepositorySnapshot(t, h)
			_, err = reserveRepositoryTest(t, h, req)
			if mode == "legacy_pinned" {
				require.NoError(t, err)
				return
			}
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

package workspace

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLifecycleSuccessorSourceFailsClosed(t *testing.T) {
	for _, mode := range []string{"stale-connector", "missing-connector", "unknown-owner", "guard-error", "stale-source", "unknown-engine", "foreign-runtime", "missing-version", "foreign-version", "active-basic", "draft-profile", "reserved"} {
		t.Run(mode, func(t *testing.T) {
			h, id, req := successorActiveWideFixture(t)
			want := "source_definition_unavailable"
			var err error
			switch mode {
			case "stale-connector":
				_, err = h.db.Exec(`UPDATE database_connectors SET identity_revision='C2'`)
				want = "connector_identity_stale"
			case "missing-connector":
				_, err = h.db.Exec(`DELETE FROM database_connectors`)
				want = "connector_unavailable"
			case "unknown-owner":
				h.service.WithBacklogOwnershipGuard(nil)
				want = "ownership_unavailable"
			case "guard-error":
				h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return errors.New("private-dsn SQL secret") }))
				want = "ownership_unavailable"
			case "stale-source":
				_, err = h.db.Exec(`UPDATE tags SET data_type='float64'`)
			case "unknown-engine":
				_, err = h.db.Exec(`UPDATE wg_runtime_versions SET payload=json_set(payload,'$.runtime_layout.dialect','unknown')`)
			case "foreign-runtime":
				_, err = h.db.Exec(`UPDATE wg_runtime_versions SET workspace_id='foreign'`)
			case "missing-version":
				_, err = h.db.Exec(`DELETE FROM write_group_versions`)
			case "foreign-version":
				_, err = h.db.Exec(`UPDATE write_group_versions SET payload=json_set(payload,'$.workspace_id','foreign')`)
			case "active-basic":
				var r writeGroupRecreateReservationRequest
				h, r = recreateBasicFixture(t, "sqlite")
				id = r.PredecessorGroupID
				_, err = h.db.Exec(`UPDATE write_groups SET status='ready' WHERE id=?`, id)
				want = "active_basic_unsupported"
			case "draft-profile":
				_, err = h.db.Exec(`UPDATE write_groups SET status='draft' WHERE id=?`, id)
				want = "source_draft_not_active"
			case "reserved":
				review, e := h.service.SuccessorReview(t.Context(), id, req)
				require.NoError(t, e)
				_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
				want = "successor_reserved"
			}
			require.NoError(t, err)
			before := recreateRepositorySnapshot(t, h)
			view, err := h.service.Lifecycle(t.Context(), id)
			if mode == "missing-version" || mode == "foreign-version" {
				require.Error(t, err)
				require.Nil(t, view)
			} else {
				require.NoError(t, err)
				require.Nil(t, view.SuccessorSeed)
				for _, a := range view.Actions[3:] {
					require.False(t, a.Allowed)
					require.Equal(t, want, a.ReasonCode)
					require.True(t, a.RequiresReview)
				}
				if mode == "stale-connector" || mode == "missing-connector" {
					require.True(t, view.Actions[1].Allowed, "safe stop remains available")
				}
				if mode == "reserved" {
					reserved, e := h.service.Lifecycle(t.Context(), view.History.ReservedGroupID)
					require.NoError(t, e)
					for _, a := range reserved.Actions {
						require.False(t, a.Allowed)
					}
					require.Nil(t, reserved.SuccessorSeed)
				}
			}
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func TestLifecycleSuccessorLongEngineFactsRemainLocal(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb", "sqlserver"} {
		t.Run(family, func(t *testing.T) {
			var h *circuitProfileHarness
			var id string
			if family == "sqlserver" {
				h, id, _, _ = successorSQLServerUnitFixture(t, WriteGroupStorageModePlain)
			} else {
				h, id, _, _ = successorMySQLShapeFixture(t, family, WriteGroupStorageModePlain)
			}
			probe := &reservedProbeCounter{}
			h.service.WithManagedTableInspector(probe)
			before := recreateRepositorySnapshot(t, h)
			view, err := h.service.Lifecycle(t.Context(), id)
			require.NoError(t, err)
			require.True(t, view.Actions[4].Allowed, "%+v", view.Actions[4])
			require.NotNil(t, view.SuccessorSeed)
			require.Zero(t, probe.calls)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func TestLifecycleSuccessorHistoricalAndLegacyRemainBlocked(t *testing.T) {
	h, req := recreateBasicFixture(t, "sqlite")
	old, err := h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
	require.NoError(t, err)
	seedLineageGeneration(t, h, old, req.LineageID, "new-head", 2)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET current_group_id='new-head',revision=2 WHERE lineage_id=?`, req.LineageID)
	require.NoError(t, err)
	view, err := h.service.Lifecycle(t.Context(), old.ID)
	require.NoError(t, err)
	require.Nil(t, view.SuccessorSeed)
	for _, action := range view.Actions[3:] {
		require.False(t, action.Allowed)
		require.Equal(t, "historical_generation", action.ReasonCode)
	}
	_, service, legacy := lifecycleReadFixture(t, false)
	view, err = service.Lifecycle(t.Context(), legacy.ID)
	require.NoError(t, err)
	require.Nil(t, view.SuccessorSeed)
	for _, action := range view.Actions[3:] {
		require.False(t, action.Allowed)
		require.Equal(t, "canonical_required", action.ReasonCode)
	}
}

func TestLifecycleSuccessorNeverAppliedMySQLOnlyOffersReview(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			h, id, _, probe := successorMySQLNeverAppliedFixture(t, family)
			before := recreateRepositorySnapshot(t, h)
			view, err := h.service.Lifecycle(t.Context(), id)
			require.NoError(t, err)
			require.True(t, view.Actions[3].Allowed, "%+v", view.Actions[3])
			require.False(t, view.Actions[4].Allowed)
			require.NotNil(t, view.SuccessorSeed)
			require.Equal(t, "saved", view.SuccessorSeed.BasisKind)
			require.Empty(t, probe.calls, "server family/capability remains a later reserve-time check")
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

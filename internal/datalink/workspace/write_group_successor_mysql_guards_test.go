package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSuccessorMySQLFamilyModeAndDeletedMatrix(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
			for _, deleted := range []bool{false, true} {
				t.Run(family+"/"+string(mode)+map[bool]string{false: "/active", true: "/deleted"}[deleted], func(t *testing.T) {
					h, id, req, p := successorMySQLShapeFixture(t, family, mode)
					if deleted {
						_, err := h.db.Exec(`UPDATE write_groups SET status='deleted' WHERE id=?`, id)
						require.NoError(t, err)
					}
					review, err := h.service.SuccessorReview(t.Context(), id, req)
					require.NoError(t, err)
					require.Empty(t, review.Blockers)
					require.Len(t, review.RequiredChecks, 1)
					require.Zero(t, p.calls)
					saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
					require.NoError(t, err)
					require.EqualValues(t, 2, saved.Generation)
					require.Equal(t, 1, p.calls)
					require.Empty(t, saved.Group.AppliedRevision)
					require.Equal(t, WriteGroupStatusDraft, saved.Group.Status)
				})
			}
		}
	}
}

func TestSuccessorMySQLPreflightFailuresLeaveNoReservation(t *testing.T) {
	for _, failure := range []string{"unknown", "offline", "nil", "database", "revision", "capability", "capability-revision", "native", "workspace-during", "target-revision-during", "target-database-during", "source-during", "backlog-during"} {
		t.Run(failure, func(t *testing.T) {
			mode := WriteGroupStorageModePlain
			if failure == "native" {
				mode = WriteGroupStorageModeNativeMonthly
			}
			h, id, req, p := successorMySQLShapeFixture(t, "mysql", mode)
			review, err := h.service.SuccessorReview(t.Context(), id, req)
			require.NoError(t, err)
			require.Empty(t, review.Blockers)
			switch failure {
			case "unknown":
				p.capability.ServerFamily = "unknown"
			case "offline":
				p.err = errors.New("offline")
			case "nil":
				p.capability = nil
			case "database":
				p.capability.Database = "foreign"
			case "revision":
				p.capability.ConnectorRevision = "foreign"
			case "capability":
				p.capability.SupportsManagedLong = false
			case "capability-revision":
				p.capability.CapabilityRevision = ""
			case "native":
				p.capability.SupportsNativePartition = false
			case "workspace-during":
				p.during = func() {
					_, err := h.service.workspaceSvc.UpdateDatabaseSetup(t.Context(), req.ExpectedWorkspaceRevision, func(context.Context, *sql.Tx, *Record) error { return nil })
					require.NoError(t, err)
				}
			case "target-revision-during":
				p.during = func() {
					_, err := h.db.Exec(`UPDATE database_connectors SET identity_revision='changed' WHERE id=?`, req.Group.Destination.ConnectorID)
					require.NoError(t, err)
				}
			case "target-database-during":
				p.during = func() {
					_, err := h.db.Exec(`UPDATE database_connectors SET connection_config='{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/foreign?tls=true"}' WHERE id=?`, req.Group.Destination.ConnectorID)
					require.NoError(t, err)
				}
			case "source-during":
				p.during = func() {
					_, err := h.db.Exec(`UPDATE write_groups SET revision='changed' WHERE id=?`, id)
					require.NoError(t, err)
				}
			case "backlog-during":
				p.during = func() { h.service.WithBacklogOwnershipGuard(nil) }
			}
			_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
			require.Error(t, err)
			require.Equal(t, 1, p.calls)
			require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
			require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_generations`))
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_lineages WHERE reserved_group_id IS NOT NULL`))
		})
	}
}

func TestSuccessorMySQLHardBlockerNeverCallsProbe(t *testing.T) {
	h, id, req, p := successorMySQLFixture(t)
	req.Group.Members = nil
	review, err := h.service.SuccessorReview(t.Context(), id, req)
	if err == nil {
		require.NotEmpty(t, review.Blockers)
		require.Empty(t, review.RequiredChecks)
		_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
	}
	require.Error(t, err)
	require.Zero(t, p.calls)
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
}

func TestSuccessorMySQLPrivateRecreateStillRefuses(t *testing.T) {
	h, id, req, _ := successorMySQLFixture(t)
	_, err := h.db.Exec(`UPDATE write_groups SET status='deleted' WHERE id=?`, id)
	require.NoError(t, err)
	private := writeGroupRecreateReservationRequest{WorkspaceID: req.WorkspaceID, ExpectedWorkspaceRevision: req.ExpectedWorkspaceRevision, LineageID: req.LineageID, ExpectedLineageRevision: req.ExpectedLineageRevision, PredecessorGroupID: id, ExpectedPredecessorRevision: req.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision, ExpectedTargetConnectorRevision: req.ExpectedTargetConnectorRevision, OperationID: req.OperationID, Group: req.Group}
	_, err = reserveRepositoryTest(t, h, private)
	require.Error(t, err)
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
}

func TestSuccessorMySQLCorruptSourceNeverFallsBack(t *testing.T) {
	h, id, req, p := successorMySQLFixture(t)
	_, err := h.db.Exec(`UPDATE wg_runtime_versions SET payload='invalid' WHERE group_id=?`, id)
	require.NoError(t, err)
	review, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.NotEmpty(t, review.Blockers)
	require.Empty(t, review.RequiredChecks)
	_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
	require.Error(t, err)
	require.Zero(t, p.calls)
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
}

func TestSuccessorMySQLCallerFamilyCannotSupplyAuthority(t *testing.T) {
	h, id, req, p := successorMySQLFixture(t)
	review, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	request := WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}
	raw, err := json.Marshal(request)
	require.NoError(t, err)
	var forged map[string]any
	require.NoError(t, json.Unmarshal(raw, &forged))
	forged["server_family"] = "mysql"
	forged["capability"] = p.capability
	forged["required_checks"] = []any{}
	forged["layout"] = map[string]any{"server_family": "mysql"}
	raw, err = json.Marshal(forged)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &request))
	p.capability.ServerFamily = "unknown"
	_, err = h.service.ReserveSuccessor(t.Context(), id, request)
	require.Error(t, err)
	require.Equal(t, 1, p.calls)
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
}

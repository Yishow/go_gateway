package workspace

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

type reservedProbeCounter struct{ calls int }

func (p *reservedProbeCounter) InspectTableAtRevision(context.Context, string, string, string, string) (*dbtarget.TableInspection, error) {
	p.calls++
	return nil, nil
}

func reservedGuardFixture(t *testing.T) (*circuitProfileHarness, *WriteGroup) {
	t.Helper()
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	original := h.mustEnsure(lineageTestDraft(h, "profile"))
	lineage := requireLineageGenesis(t, h.db, original)
	successor := seedLineageGeneration(t, h, original, lineage, "reserved-G2", 2)
	_, err := h.db.Exec(`UPDATE write_groups SET status='deleted' WHERE id=?`, original.ID)
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET reserved_group_id=?,reservation_operation_id='O2',reservation_intent_digest='intent-O2',revision=2 WHERE lineage_id=?`, successor.ID, lineage)
	require.NoError(t, err)
	// A staged proof must never open ordinary schema/apply paths for this draft.
	_, err = h.db.Exec(`UPDATE write_groups SET destination_schema_revision='schema-op',destination_schema_digest=? WHERE id=?`, strings.Repeat("d", 64), successor.ID)
	require.NoError(t, err)
	result, err := h.service.Get(t.Context(), successor.ID)
	require.NoError(t, err)
	return h, result.Group
}

func requireSuccessorReserved(t *testing.T, err error) {
	t.Helper()
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	var typed *WriteGroupLifecycleActionError
	require.ErrorAs(t, err, &typed)
	require.Equal(t, "successor_reserved", typed.ReasonCode)
	require.Equal(t, "write group lifecycle action is unavailable", err.Error())
}

func TestReservedSuccessorOrdinaryEntrypointsBlocked(t *testing.T) {
	for _, action := range []string{"update", "rename", "delete", "disable", "apply", "preview", "confirm", "confirm_replay", "direct_apply"} {
		t.Run(action, func(t *testing.T) {
			h, group := reservedGuardFixture(t)
			before := lineageReadSnapshot(t, h)
			probe := &reservedProbeCounter{}
			h.service.WithManagedTableInspector(probe)
			mutation := lifecycleReadMutation(t, h.service, group)
			var err error
			switch action {
			case "update", "rename":
				mutation.Group = cloneWriteGroup(group)
				if action == "rename" {
					mutation.Group.Name = "renamed"
				}
				_, err = h.service.Update(t.Context(), group.ID, mutation)
			case "delete":
				_, err = h.service.Delete(t.Context(), group.ID, mutation)
			case "disable":
				_, err = h.service.Disable(t.Context(), group.ID, mutation)
			case "apply":
				_, err = h.service.Apply(t.Context(), group.ID, mutation)
			case "preview":
				_, err = h.service.ManagedSchemaScope(t.Context(), group.ID, mutation)
			case "confirm":
				_, err = h.service.ConfirmManagedSchema(t.Context(), group.ID, mutation, "new-schema-op", strings.Repeat("e", 64))
			case "confirm_replay":
				_, err = h.service.ConfirmManagedSchema(t.Context(), group.ID, mutation, "schema-op", strings.Repeat("d", 64))
			case "direct_apply":
				_, err = h.workspace.UpdateDatabaseSetup(t.Context(), mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
					_, inner := h.service.applyGroupInTx(ctx, tx, record, group.ID, mutation, &WriteGroupReadiness{Ready: true}, nextAlignedCutover, true)
					return inner
				})
			}
			requireSuccessorReserved(t, err)
			require.Zero(t, probe.calls)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestReservedSuccessorReadModelStaysReadableNotReady(t *testing.T) {
	h, group := reservedGuardFixture(t)
	before := lineageReadSnapshot(t, h)
	probe := &reservedProbeCounter{}
	h.service.WithManagedTableInspector(probe)
	got, err := h.service.Get(t.Context(), group.ID)
	require.NoError(t, err)
	require.Equal(t, group.ID, got.Group.ID)
	listed, err := h.service.List(t.Context())
	require.NoError(t, err)
	require.Len(t, listed.Groups, 2)
	ready, err := h.service.Readiness(t.Context(), group.ID)
	require.NoError(t, err)
	require.False(t, ready.Ready)
	require.False(t, ready.ConfigReady)
	require.False(t, ready.SchemaReady)
	require.Contains(t, issueCodesOf(ready), "successor_reserved")
	lifecycle, err := h.service.Lifecycle(t.Context(), group.ID)
	require.NoError(t, err)
	for _, action := range lifecycle.Actions {
		require.False(t, action.Allowed)
		require.Equal(t, "successor_reserved", action.ReasonCode)
	}
	intake, err := h.service.IntakeEligibility(t.Context(), group.ID)
	require.NoError(t, err)
	require.False(t, intake.Accepting)
	require.Zero(t, probe.calls)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestReservedSuccessorApplyBlocksOtherwiseReadySchema(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	_, err := h.db.Exec(`UPDATE database_connectors SET kind='postgres',connection_config='{"database":"persisted-db","schema":"persisted-schema"}'`)
	require.NoError(t, err)
	original := h.mustEnsure(lineageTestDraft(h, "profile"))
	lineage := requireLineageGenesis(t, h.db, original)
	successor := seedLineageGeneration(t, h, original, lineage, "ready-G2", 2)
	_, err = h.db.Exec(`UPDATE write_groups SET status='deleted' WHERE id=?`, original.ID)
	require.NoError(t, err)
	saveManagedReadinessProof(t, h.db, successor.ID)
	successor = lifecycleReadGroup(t, h.service, successor.ID)
	layout, err := managedSchemaLayout(successor, map[string]schema.DataType{h.fixture.tagID: schema.DataTypeFloat32}, "postgres")
	require.NoError(t, err)
	probes := 0
	inspector := &counterApplyInspector{layout: layout, mutate: func() { probes++ }}
	h.service.WithManagedTableInspector(inspector)
	ready, err := h.service.Readiness(t.Context(), successor.ID)
	require.NoError(t, err)
	require.True(t, ready.Ready, "all non-reservation readiness facts pass")
	require.Equal(t, 1, probes)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET reserved_group_id=?,reservation_operation_id='O2',reservation_intent_digest='intent-O2',revision=2`, successor.ID)
	require.NoError(t, err)
	before := lineageReadSnapshot(t, h)
	inspector.mutate = func() { probes++ }
	_, err = h.service.Apply(t.Context(), successor.ID, lifecycleReadMutation(t, h.service, successor))
	requireSuccessorReserved(t, err)
	require.Equal(t, 1, probes)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestReservedSuccessorRecordingStartScopeBlocked(t *testing.T) {
	h, group := reservedGuardFixture(t)
	mutation := lifecycleReadMutation(t, h.service, group)
	start := NewRecordingStartService(h.workspace, h.service, nil, nil)
	request := RecordingStartRequest{RequestID: "start-reserved", WorkspaceID: group.WorkspaceID, ExpectedWorkspaceRevision: mutation.ExpectedWorkspaceRevision,
		DeviceIDs: []string{h.fixture.deviceID}, Groups: []RecordingStartGroupIntent{{GroupID: group.ID, ExpectedGroupRevision: group.Revision, ExpectedConnectorRevision: group.Destination.ConnectorRevision}}}
	before := lineageReadSnapshot(t, h)
	progress := initialRecordingStart(request, "start-digest")
	err := start.validateRecordingStart(t.Context(), progress, true)
	requireSuccessorReserved(t, err)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestReservedSuccessorOldCutoverTargetActivationBlocked(t *testing.T) {
	h, group := reservedGuardFixture(t)
	mutation := lifecycleReadMutation(t, h.service, group)
	before := lineageReadSnapshot(t, h)
	_, err := h.workspace.UpdateDatabaseSetup(t.Context(), mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		plan := &cutoverPlan{req: CutoverRequest{ExpectedWorkspaceRevision: mutation.ExpectedWorkspaceRevision, Boundary: time.Now().Add(time.Hour)}, targets: []*WriteGroup{group}, readiness: map[string]*WriteGroupReadiness{group.ID: {Ready: true}}}
		return h.service.startTargetsAt(ctx, tx, record, plan, &WriteGroupCutover{})
	})
	requireSuccessorReserved(t, err)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestReservedSuccessorPublicRecordingStartCannotPrepareOrActivate(t *testing.T) {
	h, group := reservedGuardFixture(t)
	mutation := lifecycleReadMutation(t, h.service, group)
	probe := &reservedProbeCounter{}
	h.service.WithManagedTableInspector(probe)
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(h.db))
	start := NewRecordingStartService(h.workspace, h.service, nil, ledger)
	request := RecordingStartRequest{RequestID: "start-reserved", WorkspaceID: group.WorkspaceID, ExpectedWorkspaceRevision: mutation.ExpectedWorkspaceRevision,
		DeviceIDs: []string{h.fixture.deviceID}, Groups: []RecordingStartGroupIntent{{GroupID: group.ID, ExpectedGroupRevision: group.Revision, ExpectedConnectorRevision: group.Destination.ConnectorRevision}}}
	before := lineageReadSnapshot(t, h)
	operation, err := start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationFailed, operation.Status)
	require.Equal(t, "successor_reserved", operation.Reason)
	require.Equal(t, "review_selection", operation.NextAction)
	require.False(t, operation.Groups[0].Saved)
	require.False(t, operation.Groups[0].Ready)
	require.False(t, operation.Groups[0].Applied)
	require.False(t, operation.Devices[0].Activated)
	ledgerBefore := lifecyclePreservedData(t, h.db)
	replayed, err := start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, operation, replayed, "exact blocked replay must retain the reservation reason and non-retry guidance")
	require.Equal(t, ledgerBefore, lifecyclePreservedData(t, h.db), "blocked replay must not replace its operation result")
	require.Zero(t, probe.calls)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestReservedSuccessorDeleteKeepsReservationReasonAfterTargetChange(t *testing.T) {
	for _, state := range []string{"edited", "removed"} {
		t.Run(state, func(t *testing.T) {
			h, group := reservedGuardFixture(t)
			mutation := lifecycleReadMutation(t, h.service, group)
			var err error
			if state == "edited" {
				_, err = h.db.Exec(`UPDATE database_connectors SET identity_revision='edited-target' WHERE id=?`, group.Destination.ConnectorID)
			} else {
				_, err = h.db.Exec(`DELETE FROM database_connectors WHERE id=?`, group.Destination.ConnectorID)
			}
			require.NoError(t, err)
			before := lineageReadSnapshot(t, h)
			_, err = h.service.Delete(t.Context(), group.ID, mutation)
			requireSuccessorReserved(t, err)
			var action *WriteGroupLifecycleActionError
			require.ErrorAs(t, err, &action)
			require.Equal(t, "review_lifecycle", action.NextAction)
			require.Equal(t, before, lineageReadSnapshot(t, h))
			mutation.ExpectedConnectorRevision = "stale-saved-cas"
			_, err = h.service.Delete(t.Context(), group.ID, mutation)
			require.ErrorIs(t, err, ErrWriteGroupRevisionConflict, "saved connector CAS retains precedence over reservation check")
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

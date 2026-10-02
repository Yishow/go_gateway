package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/delivery"

	"github.com/stretchr/testify/require"
)

type writeGroupLifecycleInspector struct{}

func (writeGroupLifecycleInspector) InspectTable(context.Context, string, string, string) (*dbtarget.TableInspection, error) {
	return &dbtarget.TableInspection{
		Status: dbtarget.TableInspectionExists,
		Schema: "persisted-schema",
		Table:  "raw_values",
		Columns: []dbtarget.ColumnInfo{{
			Name:     "temperature",
			DataType: "REAL",
		}},
	}, nil
}

func TestGroupLifecycleDisableUsesCASAndStopsIntake(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `
		UPDATE database_connectors SET kind = ?, connection_config = ? WHERE id = ?
	`, "postgres", `{"database":"persisted-db","schema":"persisted-schema"}`, fixture.connectorID)
	require.NoError(t, err)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)

	beforeDisable, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	disabled, err := service.Disable(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeDisable.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDisabled, disabled.Group.Status)
	require.Equal(t, created.Group.AppliedRevision, disabled.Group.AppliedRevision)
	require.Equal(t, created.Group.Members, disabled.Group.Members)
	require.False(t, mustIntakeEligibility(t, service, created.Group.ID).Accepting)

	workspaceBeforeEdit, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	edited := cloneWriteGroup(disabled.Group)
	edited.Name = "disabled draft edit"
	edited.RowPolicy.IntervalSeconds = 15
	draft, err := service.Update(ctx, disabled.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeEdit.DatabaseSetupRevision,
		ExpectedGroupRevision: disabled.Group.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: edited,
	})
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDisabled, draft.Group.Status)
	require.False(t, mustIntakeEligibility(t, service, draft.Group.ID).Accepting)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	workspaceBeforeApply, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	reenabled, err := service.Apply(ctx, draft.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeApply.DatabaseSetupRevision,
		ExpectedGroupRevision: draft.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusReady, reenabled.Group.Status)
	require.False(t, mustIntakeEligibility(t, service, reenabled.Group.ID).Accepting)

	_, err = service.Disable(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeDisable.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
}

func TestGroupLifecycleApplySwitchesAtNextBucketWithoutChangingStableID(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)
	service.WithClock(func() time.Time { return now })

	beforeApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	first, err := service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeApply.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.Equal(t, created.Group.ID, first.Group.ID)
	require.Equal(t, created.Group.Revision, first.Group.Revision)
	require.Equal(t, first.Group.Revision, first.Group.AppliedRevision)
	effectiveAt := now.Truncate(15 * time.Second).Add(15 * time.Second)
	_, err = service.ResolveAppliedAt(ctx, first.Group.ID, now)
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	oldSnapshot, err := service.ResolveAppliedAt(ctx, first.Group.ID, effectiveAt)
	require.NoError(t, err)
	require.Equal(t, first.Group.Revision, oldSnapshot.GroupRevision)

	workspaceBeforeRename, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	rename := cloneWriteGroup(first.Group)
	rename.Name = "renamed"
	draft, err := service.Update(ctx, first.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeRename.DatabaseSetupRevision,
		ExpectedGroupRevision: first.Group.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: rename,
	})
	require.NoError(t, err)
	require.Equal(t, first.Group.ID, draft.Group.ID)
	require.NotEqual(t, first.Group.Revision, draft.Group.Revision)
	require.Equal(t, first.Group.AppliedRevision, draft.Group.AppliedRevision)
	require.Equal(t, WriteGroupStatusDraft, draft.Group.Status)

	now = effectiveAt.Add(2 * time.Second)
	workspaceBeforeRenameApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	renameApplied, err := service.Apply(ctx, draft.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeRenameApply.DatabaseSetupRevision,
		ExpectedGroupRevision: draft.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.Equal(t, first.Group.ID, renameApplied.Group.ID)
	require.Equal(t, draft.Group.Revision, renameApplied.Group.Revision)
	require.Equal(t, first.Group.AppliedRevision, renameApplied.Group.AppliedRevision)
	require.Equal(t, WriteGroupStatusReady, renameApplied.Group.Status)
	require.True(t, mustIntakeEligibility(t, service, renameApplied.Group.ID).Accepting)
	unchangedRuntime, err := service.ResolveAppliedAt(ctx, renameApplied.Group.ID, now)
	require.NoError(t, err)
	require.Equal(t, first.Group.Revision, unchangedRuntime.GroupRevision)
	require.Equal(t, first.Group.Name, unchangedRuntime.Group.Name)

	workspaceBeforePolicyEdit, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	policyEdit := cloneWriteGroup(renameApplied.Group)
	policyEdit.RowPolicy.IntervalSeconds = 30
	policyDraft, err := service.Update(ctx, renameApplied.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforePolicyEdit.DatabaseSetupRevision,
		ExpectedGroupRevision: renameApplied.Group.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: policyEdit,
	})
	require.NoError(t, err)
	require.Equal(t, first.Group.AppliedRevision, policyDraft.Group.AppliedRevision)
	workspaceBeforeSecondApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	second, err := service.Apply(ctx, policyDraft.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeSecondApply.DatabaseSetupRevision,
		ExpectedGroupRevision: policyDraft.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.Equal(t, first.Group.ID, second.Group.ID)
	require.Equal(t, policyDraft.Group.Revision, second.Group.Revision)
	require.NotEqual(t, first.Group.AppliedRevision, second.Group.AppliedRevision)
	require.Equal(t, second.Group.Revision, second.Group.AppliedRevision)

	secondEffectiveAt := now.Truncate(30 * time.Second).Add(30 * time.Second)
	oldAtCutover, err := service.ResolveAppliedAt(ctx, second.Group.ID, effectiveAt.Add(1*time.Second))
	require.NoError(t, err)
	require.Equal(t, first.Group.Revision, oldAtCutover.GroupRevision)
	newAtCutover, err := service.ResolveAppliedAt(ctx, second.Group.ID, secondEffectiveAt)
	require.NoError(t, err)
	require.Equal(t, second.Group.Revision, newAtCutover.GroupRevision)
	require.Equal(t, "renamed", newAtCutover.Group.Name)
	require.Equal(t, 30, newAtCutover.Group.RowPolicy.IntervalSeconds)
	require.True(t, mustIntakeEligibility(t, service, second.Group.ID).Accepting)
}

func TestGroupLifecycleDeleteAppliedRequiresOwnershipGuard(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE write_groups SET applied_revision = ? WHERE id = ?`, "applied-1", created.Group.ID)
	require.NoError(t, err)
	beforeDelete, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = service.Delete(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeDelete.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
}

func TestGroupLifecycleApplyRollsBackVersionAndGroupWhenWorkspaceSaveFails(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	beforeWorkspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	beforeGroup, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `CREATE TRIGGER fail_write_group_lifecycle_workspace_save
		BEFORE UPDATE ON system_settings
		WHEN NEW.key = 'studio_v2_workspace'
		BEGIN SELECT RAISE(ABORT, 'lifecycle-save-failed'); END`)
	require.NoError(t, err)

	_, err = service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeWorkspace.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.Error(t, err)
	afterWorkspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	afterGroup, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, beforeWorkspace, afterWorkspace)
	require.Equal(t, beforeGroup, afterGroup)
	var versionCount int
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, created.Group.ID,
	).Scan(&versionCount))
	require.Zero(t, versionCount)
}

func TestGroupLifecycleDeleteGuardFailureLeavesTombstoneAndBacklogUntouched(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE write_groups SET applied_revision = ? WHERE id = ?`, "applied-1", created.Group.ID)
	require.NoError(t, err)
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error {
		return errors.New("legacy backlog ownership unavailable")
	}))
	beforeWorkspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = service.Delete(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeWorkspace.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	after, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDraft, after.Group.Status)
	require.Equal(t, "applied-1", after.Group.AppliedRevision)
}

func TestGroupLifecycleDeleteWithOutboxPreservesPayloadAndReceipt(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)
	service.WithClock(func() time.Time { return now })
	beforeApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	firstApplied, err := service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeApply.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	now = time.Date(2026, 10, 2, 12, 0, 16, 0, time.UTC)

	workspaceBeforeEdit, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	edited := cloneWriteGroup(firstApplied.Group)
	edited.RowPolicy.IntervalSeconds = 30
	draft, err := service.Update(ctx, firstApplied.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeEdit.DatabaseSetupRevision,
		ExpectedGroupRevision: firstApplied.Group.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: edited,
	})
	require.NoError(t, err)
	workspaceBeforeSecondApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	secondApplied, err := service.Apply(ctx, draft.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeSecondApply.DatabaseSetupRevision,
		ExpectedGroupRevision: draft.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	now = time.Date(2026, 10, 2, 12, 0, 31, 0, time.UTC)
	require.True(t, mustIntakeEligibility(t, service, secondApplied.Group.ID).Accepting)

	payload := []byte(`{"group_id":"` + secondApplied.Group.ID + `","group_revision":"` + firstApplied.Group.Revision + `","connector_revision":"connector-revision-1","schema":"persisted-schema","table":"raw_values","record_id":"record-1","temperature":21.5}`)
	outbox := delivery.NewSQLOutbox(db)
	receipts := delivery.NewSQLReceiptLedger(db)
	require.NoError(t, outbox.Enqueue(&delivery.OutboxItem{
		ID: "outbox-1", DestinationID: secondApplied.Group.ID, DestinationRevision: secondApplied.Group.Destination.ConnectorRevision,
		PlanRevision: firstApplied.Group.Revision,
		RecordID:     "record-1", CalculationRevision: 1, Table: "raw_values", Payload: payload,
		ObservedAt: time.Date(2026, 10, 2, 11, 59, 0, 0, time.UTC),
	}))
	service.WithBacklogOwnershipGuard(sqliteWriteGroupBacklogGuard{})

	beforeDelete, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	deleted, err := service.Delete(ctx, secondApplied.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeDelete.DatabaseSetupRevision,
		ExpectedGroupRevision: secondApplied.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDeleted, deleted.Group.Status)
	require.Equal(t, secondApplied.Group.AppliedRevision, deleted.Group.AppliedRevision)

	var persistedPayload []byte
	require.NoError(t, db.QueryRowContext(ctx, `SELECT payload FROM gw_delivery_outbox WHERE id = ?`, "outbox-1").Scan(&persistedPayload))
	require.Equal(t, payload, persistedPayload)

	var sentPayload []byte
	var sentMu sync.Mutex
	sender := deliverySenderFunc(func(_ context.Context, item *delivery.OutboxItem) error {
		sentMu.Lock()
		defer sentMu.Unlock()
		sentPayload = append([]byte(nil), item.Payload...)
		return nil
	})
	worker := delivery.NewDeliveryWorker(secondApplied.Group.ID, outbox, receipts, sender, delivery.WorkerConfig{FlushInterval: time.Millisecond})
	workerCtx, cancel := context.WithCancel(ctx)
	cancel()
	done := make(chan struct{})
	go func() {
		worker.Start(workerCtx)
		close(done)
	}()
	<-done
	worker.Stop()
	sentMu.Lock()
	require.Equal(t, payload, sentPayload)
	sentMu.Unlock()
	hasReceipt, err := receipts.HasReceipt(secondApplied.Group.ID, "record-1", 1)
	require.NoError(t, err)
	require.True(t, hasReceipt)
	reloaded, err := service.Get(ctx, secondApplied.Group.ID)
	require.NoError(t, err)
	require.Equal(t, deleted.Group, reloaded.Group)
}

type deliverySenderFunc func(context.Context, *delivery.OutboxItem) error

func (f deliverySenderFunc) Send(ctx context.Context, item *delivery.OutboxItem) error {
	return f(ctx, item)
}

type sqliteWriteGroupBacklogGuard struct{}

type writeGroupBacklogGuardFunc func(context.Context, *sql.Tx, *WriteGroup) error

func (f writeGroupBacklogGuardFunc) CheckWriteGroupBacklog(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
	return f(ctx, tx, group)
}

func (sqliteWriteGroupBacklogGuard) CheckWriteGroupBacklog(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
	rows, err := tx.QueryContext(ctx, `SELECT destination_revision, plan_revision, payload FROM gw_delivery_outbox WHERE destination_id = ?`, group.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var destinationRevision, planRevision string
		var payload []byte
		if err := rows.Scan(&destinationRevision, &planRevision, &payload); err != nil {
			return err
		}
		if destinationRevision != group.Destination.ConnectorRevision || planRevision == "" || len(payload) == 0 {
			return errors.New("backlog ownership is not immutable")
		}
		var versionPayload string
		if err := tx.QueryRowContext(ctx, `
			SELECT payload FROM write_group_versions WHERE group_id = ? AND group_revision = ?
		`, group.ID, planRevision).Scan(&versionPayload); err != nil {
			return err
		}
		var versionGroup WriteGroup
		if err := json.Unmarshal([]byte(versionPayload), &versionGroup); err != nil {
			return err
		}
		var envelope struct {
			GroupID           string `json:"group_id"`
			GroupRevision     string `json:"group_revision"`
			ConnectorRevision string `json:"connector_revision"`
			Schema            string `json:"schema"`
			Table             string `json:"table"`
		}
		if err := json.Unmarshal(payload, &envelope); err != nil {
			return err
		}
		if envelope.GroupID != versionGroup.ID || envelope.GroupRevision != versionGroup.Revision ||
			envelope.ConnectorRevision != versionGroup.Destination.ConnectorRevision ||
			envelope.Schema != versionGroup.Destination.TableSchema || envelope.Table != versionGroup.Destination.TableName {
			return errors.New("backlog payload owner does not match saved version")
		}
	}
	return rows.Err()
}

func mustIntakeEligibility(t *testing.T, service *WriteGroupService, id string) *WriteGroupIntakeEligibility {
	t.Helper()
	eligibility, err := service.IntakeEligibility(t.Context(), id)
	require.NoError(t, err)
	return eligibility
}

func TestGroupLifecycleDisableStillWorksAfterTheDestinationEndpointWasEdited(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)

	// An operator edits the connector: its identity revision moves on, so the
	// group's saved destination revision is now stale.
	_, err = db.ExecContext(ctx, `UPDATE database_connectors SET identity_revision = 'connector-revision-2' WHERE id = ?`, fixture.connectorID)
	require.NoError(t, err)

	before, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	disabled, err := service.Disable(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: before.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err, "stopping intake must not depend on the endpoint still matching the saved revision")
	require.Equal(t, WriteGroupStatusDisabled, disabled.Group.Status)
	require.Equal(t, "connector-revision-1", disabled.Group.Destination.ConnectorRevision, "the group keeps the destination identity it was accepted for")
	require.False(t, mustIntakeEligibility(t, service, created.Group.ID).Accepting)

	// A stale expected revision is still a conflict: only the live-connector
	// requirement is dropped, not compare-and-swap on the group's own revision.
	other, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = service.Disable(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: other.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-2",
	})
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
}

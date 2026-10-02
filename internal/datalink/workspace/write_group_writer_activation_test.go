package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/delivery"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type writerOwnershipActivationProbe struct {
	calls              int
	request            WriterOwnershipActivationRequest
	versionInTx        bool
	versionPayloadInTx string
	effectiveInTx      time.Time
}

type writerOwnershipActivationFunc func(context.Context, *sql.Tx, WriterOwnershipActivationRequest) error

func (f writerOwnershipActivationFunc) ActivateInTx(
	ctx context.Context,
	tx *sql.Tx,
	request WriterOwnershipActivationRequest,
) error {
	return f(ctx, tx, request)
}

func (p *writerOwnershipActivationProbe) ActivateInTx(
	ctx context.Context,
	tx *sql.Tx,
	request WriterOwnershipActivationRequest,
) error {
	p.calls++
	p.request = request
	var count int
	if err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM write_group_versions
		WHERE group_id = ? AND group_revision = ?
	`, request.GroupID, request.AppliedRevision).Scan(&count); err != nil {
		return err
	}
	p.versionInTx = count == 1
	return tx.QueryRowContext(ctx, `
		SELECT payload, effective_at
		FROM write_group_versions
		WHERE group_id = ? AND group_revision = ?
	`, request.GroupID, request.AppliedRevision).Scan(&p.versionPayloadInTx, &p.effectiveInTx)
}

func newFixedWriterActivationGroup(ctx context.Context, t *testing.T, db *sql.DB) (*WriteGroupService, writeGroupFixture, *WriteGroupSaveResult) {
	t.Helper()
	fixture := seedWriteGroupFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `
		UPDATE database_connectors
		SET kind = ?, connection_config = ?
		WHERE id = ?
	`, schema.DatabaseConnectorKindPostgres, `{"database":"persisted-db","schema":"persisted-schema"}`, fixture.connectorID)
	require.NoError(t, err)
	workspaceSvc := NewService(NewSQLRepository(db))
	repo := NewSQLWriteGroupRepository(db)
	ids := []string{"group-A", "group-revision-1"}
	repo.newID = func() string {
		if len(ids) == 0 {
			return "unexpected-test-id"
		}
		id := ids[0]
		ids = ids[1:]
		return id
	}
	service := NewWriteGroupService(workspaceSvc, repo)
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	group := newBasicWriteGroup(fixture)
	group.RowPolicy.IntervalSeconds = 15
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: group,
	})
	require.NoError(t, err)
	return service, fixture, created
}

func TestWriteGroupWriterOwnershipActivationBarrierSeesPreparedVersion(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newFixedWriterActivationGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)
	service.WithClock(func() time.Time { return now })
	ownerKey := "test_writer_ownership_owner"
	_, err := db.ExecContext(ctx, `
		INSERT INTO system_settings (key, value, updated_at) VALUES (?, ?, ?)
	`, ownerKey, "legacy", now)
	require.NoError(t, err)
	var request WriterOwnershipActivationRequest
	var versionInTx bool
	var versionPayloadInTx string
	var effectiveInTx time.Time
	var ownerBefore string
	service.WithWriterOwnershipActivationBarrier(writerOwnershipActivationFunc(func(
		ctx context.Context,
		tx *sql.Tx,
		activation WriterOwnershipActivationRequest,
	) error {
		request = activation
		if err := tx.QueryRowContext(ctx, `
			SELECT payload, effective_at
			FROM write_group_versions WHERE group_id = ? AND group_revision = ?
		`, activation.GroupID, activation.AppliedRevision).Scan(&versionPayloadInTx, &effectiveInTx); err != nil {
			return err
		}
		versionInTx = true
		if err := tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = ?`, ownerKey).Scan(&ownerBefore); err != nil {
			return err
		}
		_, err := tx.ExecContext(ctx, `
			UPDATE system_settings SET value = ?, updated_at = ? WHERE key = ?
		`, `{"owner":"canonical","effective_at":"2026-10-02T12:00:15Z"}`, activation.EffectiveAt, ownerKey)
		return err
	}))

	workspaceBeforeApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	applied, err := service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID:               fixture.workspaceID,
		ExpectedWorkspaceRevision: workspaceBeforeApply.DatabaseSetupRevision,
		ExpectedGroupRevision:     created.Group.Revision,
		ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)

	require.Equal(t, "group-A", created.Group.ID)
	require.Equal(t, "group-revision-1", created.Group.Revision)
	require.True(t, versionInTx)
	require.Equal(t, "legacy", ownerBefore)
	require.Equal(t, time.Date(2026, 10, 2, 12, 0, 15, 0, time.UTC), request.EffectiveAt)
	require.Equal(t, request.EffectiveAt, effectiveInTx.UTC())
	var snapshot WriteGroup
	require.NoError(t, json.Unmarshal([]byte(versionPayloadInTx), &snapshot))
	require.Equal(t, *applied.Group, snapshot)
	require.Equal(t, WriterOwnershipActivationRequest{
		WorkspaceID:                   fixture.workspaceID,
		GroupID:                       created.Group.ID,
		ExpectedDatabaseSetupRevision: workspaceBeforeApply.DatabaseSetupRevision,
		ExpectedGroupRevision:         created.Group.Revision,
		ExpectedConnectorRevision:     "connector-revision-1",
		PreviousAppliedRevision:       "",
		AppliedRevision:               created.Group.Revision,
		EffectiveAt:                   time.Date(2026, 10, 2, 12, 0, 15, 0, time.UTC),
	}, request)
	var ownerAfter string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = ?`, ownerKey).Scan(&ownerAfter))
	require.Equal(t, `{"owner":"canonical","effective_at":"2026-10-02T12:00:15Z"}`, ownerAfter)
}

func TestWriteGroupWriterOwnershipActivationBarrierErrorsRollBackPreparedState(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	ownerKey := "test_writer_ownership_owner"
	barrierErr := errors.New("activation barrier rejected")
	service.WithWriterOwnershipActivationBarrier(writerOwnershipActivationFunc(func(
		ctx context.Context,
		tx *sql.Tx,
		_ WriterOwnershipActivationRequest,
	) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO system_settings (key, value, updated_at)
			VALUES (?, ?, ?)
			ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at
		`, ownerKey, "canonical", time.Now().UTC())
		if err != nil {
			return err
		}
		return barrierErr
	}))

	beforeWorkspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	beforeGroup, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	_, err = service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeWorkspace.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.ErrorIs(t, err, barrierErr)

	afterWorkspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	afterGroup, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, beforeWorkspace, afterWorkspace)
	require.Equal(t, beforeGroup, afterGroup)
	var versionCount, ownerCount int
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, created.Group.ID,
	).Scan(&versionCount))
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM system_settings WHERE key = ?`, ownerKey,
	).Scan(&ownerCount))
	require.Zero(t, versionCount)
	require.Zero(t, ownerCount)
}

func TestWriteGroupWriterOwnershipActivationBarrierPreservesUnavailableCause(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	barrierErr := errors.New("activation unavailable cause")
	service.WithWriterOwnershipActivationBarrier(writerOwnershipActivationFunc(func(
		context.Context,
		*sql.Tx,
		WriterOwnershipActivationRequest,
	) error {
		return errors.Join(ErrDatabaseSetupUnavailable, barrierErr)
	}))
	workspaceBefore, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBefore.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.ErrorIs(t, err, ErrWriteGroupServiceUnavailable)
	require.ErrorIs(t, err, ErrDatabaseSetupUnavailable)
	require.ErrorIs(t, err, barrierErr)
}

func TestWriteGroupWriterOwnershipActivationBarrierDoesNotRunForRename(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)
	service.WithClock(func() time.Time { return now })
	probe := &writerOwnershipActivationProbe{}
	service.WithWriterOwnershipActivationBarrier(probe)

	beforeFirstApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	first, err := service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeFirstApply.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	effectiveAt := now.Truncate(15 * time.Second).Add(15 * time.Second)

	beforeRename, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	rename := cloneWriteGroup(first.Group)
	rename.Name = "renamed"
	draft, err := service.Update(ctx, first.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeRename.DatabaseSetupRevision,
		ExpectedGroupRevision: first.Group.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: rename,
	})
	require.NoError(t, err)
	now = effectiveAt.Add(2 * time.Second)
	beforeRenameApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	renamed, err := service.Apply(ctx, first.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeRenameApply.DatabaseSetupRevision,
		ExpectedGroupRevision: draft.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)

	require.Equal(t, 1, probe.calls)
	require.Equal(t, first.Group.AppliedRevision, renamed.Group.AppliedRevision)
	require.Equal(t, first.Group.AppliedRevision, probe.request.AppliedRevision)
}

func TestWriteGroupWriterOwnershipActivationBarrierPreservesLegacyAndAcceptedRecords(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)
	service.WithClock(func() time.Time { return now })
	probe := &writerOwnershipActivationProbe{}
	service.WithWriterOwnershipActivationBarrier(probe)

	var mappingEnabled bool
	var mappingTransform string
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT enabled, transform_pipeline FROM mappings WHERE point_id = ? AND tag_id = ?
	`, fixture.pointID, fixture.tagID).Scan(&mappingEnabled, &mappingTransform))
	workspaceBeforeApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	first, err := service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeApply.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)

	payload := []byte(`{"group_id":"` + first.Group.ID + `","group_revision":"` + first.Group.Revision + `","value":21.5}`)
	outbox := delivery.NewSQLOutbox(db)
	receipts := delivery.NewSQLReceiptLedger(db)
	observedAt := time.Date(2026, 10, 2, 12, 0, 6, 0, time.UTC)
	require.NoError(t, outbox.Enqueue(&delivery.OutboxItem{
		ID: "activation-outbox-1", DestinationID: first.Group.ID,
		DestinationRevision: first.Group.Destination.ConnectorRevision, PlanRevision: first.Group.Revision,
		RecordID: "record-1", CalculationRevision: 1, Table: first.Group.Destination.TableName,
		Payload: payload, ObservedAt: observedAt,
	}))
	require.NoError(t, receipts.SaveReceipt(&delivery.Receipt{
		DestinationID: first.Group.ID, RecordID: "record-1", CalculationRevision: 1,
		Table: first.Group.Destination.TableName, DeliveredAt: observedAt,
	}))

	now = time.Date(2026, 10, 2, 12, 0, 16, 0, time.UTC)
	beforeEdit, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	edited := cloneWriteGroup(first.Group)
	edited.RowPolicy.IntervalSeconds = 30
	draft, err := service.Update(ctx, first.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeEdit.DatabaseSetupRevision,
		ExpectedGroupRevision: first.Group.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: edited,
	})
	require.NoError(t, err)
	beforeSecondApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	second, err := service.Apply(ctx, first.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeSecondApply.DatabaseSetupRevision,
		ExpectedGroupRevision: draft.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.NotEqual(t, first.Group.AppliedRevision, second.Group.AppliedRevision)

	var persistedPayload []byte
	var persistedPlanRevision, persistedDestinationRevision string
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT payload, plan_revision, destination_revision
		FROM gw_delivery_outbox WHERE id = ?
	`, "activation-outbox-1").Scan(&persistedPayload, &persistedPlanRevision, &persistedDestinationRevision))
	require.Equal(t, payload, persistedPayload)
	require.Equal(t, first.Group.Revision, persistedPlanRevision)
	require.Equal(t, first.Group.Destination.ConnectorRevision, persistedDestinationRevision)
	hasReceipt, err := receipts.HasReceipt(first.Group.ID, "record-1", 1)
	require.NoError(t, err)
	require.True(t, hasReceipt)
	var afterEnabled bool
	var afterTransform string
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT enabled, transform_pipeline FROM mappings WHERE point_id = ? AND tag_id = ?
	`, fixture.pointID, fixture.tagID).Scan(&afterEnabled, &afterTransform))
	require.Equal(t, mappingEnabled, afterEnabled)
	require.Equal(t, mappingTransform, afterTransform)
	require.Equal(t, 2, probe.calls)
	require.Equal(t, first.Group.AppliedRevision, probe.request.PreviousAppliedRevision)
	require.Equal(t, draft.Group.Revision, probe.request.AppliedRevision)
	require.Equal(t, beforeSecondApply.DatabaseSetupRevision, probe.request.ExpectedDatabaseSetupRevision)
	require.Equal(t, draft.Group.Revision, probe.request.ExpectedGroupRevision)
	require.Equal(t, "connector-revision-1", probe.request.ExpectedConnectorRevision)
	require.Equal(t, time.Date(2026, 10, 2, 12, 0, 30, 0, time.UTC), probe.request.EffectiveAt)
}

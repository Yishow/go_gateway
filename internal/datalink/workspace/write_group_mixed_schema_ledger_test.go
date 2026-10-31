package workspace

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"path/filepath"
	"testing"
)

func TestMixedSchemaLedgerFreshConfirmationCannotActivateOrAuthorizeChangedMethod(t *testing.T) {
	for _, layout := range []WriteGroupLayoutVersion{WriteGroupLayoutManagedWideProfileV1, WriteGroupLayoutManagedLongV1} {
		for _, change := range []string{"method", "cadence", "required", "measurement identity", "counter policy", "proof lost during inspection"} {
			t.Run(string(layout)+"/"+change, func(t *testing.T) { mixedSchemaLedgerCase(t, layout, change) })
		}
	}
}
func TestMixedSharedLongStillRequiresFreshSchemaProof(t *testing.T) {
	mixedSchemaLedgerCase(t, WriteGroupLayoutManagedLongV1, "shared legacy long")
}
func mixedSchemaLedgerCase(t *testing.T, layout WriteGroupLayoutVersion, change string) {
	ctx := t.Context()
	h := newCircuitProfileHarness(t, filepath.Join(t.TempDir(), "workspace.db"))
	target := filepath.Join(t.TempDir(), "target.db")
	external, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	t.Cleanup(func() { _ = external.Close() })
	_, err = external.Exec(`PRAGMA user_version=0`)
	require.NoError(t, err)
	setSQLiteConnectorConfig(t, h.db, h.fixture.connectorID, target)
	var legacy *WriteGroup
	if change == "shared legacy long" {
		old := h.longDraft(WriteGroupProfileInstantaneous, "main", h.primary())
		old.Destination.Database = target
		legacy = h.mustEnsure(old)
	}
	draft := mixedDraft(h, WriteGroupMemberCalculationMean, 15, "mixed_results")
	draft.LayoutVersion = layout
	if layout == WriteGroupLayoutManagedLongV1 {
		draft.StorageMode = WriteGroupStorageModePlain
		draft.Destination.TableName = ""
	}
	draft.Destination.Database = target
	if change == "counter policy" || change == "measurement identity" {
		draft.Members[0].Calculation = WriteGroupMemberCalculationDelta
		_, err = h.db.Exec(`UPDATE points SET data_type='uint64' WHERE id=?`, h.primary().pointID)
		require.NoError(t, err)
		_, err = h.db.Exec(`UPDATE tags SET data_type='uint64' WHERE id=?`, h.primary().tagID)
		require.NoError(t, err)
		id := seedMeasurement(t, h, h.primary(), "counter")
		draft.Members[0].MeasurementID = &id
		_, err = h.db.Exec(`UPDATE measurement_definitions SET unit='kWh',series_epoch='epoch-1',numeric_encoding='uint64',counter_policy='{}' WHERE id=?`, id)
		require.NoError(t, err)
	}

	saved, err := h.service.Create(ctx, h.mutation(ctx, draft))
	require.NoError(t, err)
	if layout == WriteGroupLayoutManagedWideProfileV1 {
		require.Empty(t, saved.Group.RowPolicy.BucketStartColumn)
		require.Equal(t, "due_at", saved.Group.RowPolicy.DueAtColumn)
	}
	if legacy != nil {
		require.Equal(t, legacy.Destination.TableName, saved.Group.Destination.TableName)
		require.NotEqual(t, legacy.ID, saved.Group.ID)
		require.Empty(t, saved.Group.Destination.SchemaRevision)
	}
	req := h.mutation(ctx, saved.Group)
	req.ExpectedGroupRevision = saved.Group.Revision
	scope, err := h.service.ManagedSchemaScope(ctx, saved.Group.ID, req)
	require.NoError(t, err)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(h.db))
	inspector := dbtarget.NewManagedTableInspector(targets, h.db)
	h.service.WithManagedTableInspector(inspector)
	inspect := func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		found, err := inspector.InspectTableAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, table)
		if err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		if table == scope.GroupLayout.TableName && layout == WriteGroupLayoutManagedLongV1 && found.Status == dbtarget.TableInspectionExists {
			indexes, e := inspector.InspectTableIndexesAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, table)
			if e != nil {
				return recordingplan.TargetTableInspection{}, e
			}
			found.Indexes = indexes
			found.IndexesInspected = true
		}
		return recordingPlanTableInspection(found), nil
	}
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(h.db))
	token, err := ledger.PrepareSchemaPreview(ctx, scope, inspect)
	require.NoError(t, err)
	_, err = h.service.ConfirmManagedSchemaWithStorage(ctx, saved.Group.ID, req, token.OperationID, token.Digest, &StorageEvidence{Layout: *scope.GroupLayout, CreatedTable: true})
	require.Error(t, err, "a preview alone is not a succeeded schema operation")
	unconfirmed, err := h.service.Get(ctx, saved.Group.ID)
	require.NoError(t, err)
	require.Empty(t, unconfirmed.Group.Destination.SchemaRevision)
	require.Equal(t, saved.WorkspaceRevision, unconfirmed.WorkspaceRevision)

	before, err := h.service.Readiness(ctx, saved.Group.ID)
	require.NoError(t, err)
	require.True(t, before.ConfigReady)
	require.False(t, before.SchemaReady)
	require.False(t, before.Ready)
	verified, err := ledger.ValidatePreviewTokenForApply(ctx, token.Token, scope)
	require.NoError(t, err)
	op, _, err := ledger.ApplySchemaPreview(ctx, verified, recordingplan.SchemaApplyTarget{Inspect: inspect, Execute: func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
		done, err := inspector.ExecuteSchemaStatements(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, statements)
		if done == nil {
			return recordingplan.SchemaExecution{}, err
		}
		return recordingplan.SchemaExecution{Committed: done.Committed, RolledBack: done.RolledBack}, err
	}})
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, op.Status)
	confirmed, err := h.service.ConfirmManagedSchemaWithStorage(ctx, saved.Group.ID, req, op.OperationID, op.VerifiedDigest, &StorageEvidence{Layout: *scope.GroupLayout, CreatedTable: true})
	require.NoError(t, err)
	if change == "proof lost during inspection" {
		initial, err := h.service.readinessSnapshot(ctx, saved.Group.ID, nil)
		require.NoError(t, err)
		require.NoError(t, initial.schemaProofErr)
		raced, err := h.service.evaluateReadinessSnapshot(ctx, saved.Group.ID, nil, initial, func(ctx context.Context) (*writeGroupReadinessSnapshot, error) {
			require.NoError(t, recordingplan.NewSQLRepository(h.db).DeletePreviewToken(ctx, token.Token))
			post, err := h.service.readinessSnapshot(ctx, saved.Group.ID, nil)
			require.NoError(t, err)
			require.Error(t, post.schemaProofErr)
			return post, err
		})
		require.NoError(t, err)
		require.False(t, raced.SchemaReady, "post-inspection proof loss must fail closed")
		require.False(t, raced.Ready)
		return
	}
	ready, err := h.service.Readiness(ctx, saved.Group.ID)
	require.NoError(t, err)
	require.True(t, ready.ConfigReady)
	require.True(t, ready.SchemaReady)
	require.False(t, ready.Ready)
	require.True(t, recordingReductionUnsupported(confirmed.Group, nil))
	apply := h.mutation(ctx, confirmed.Group)
	apply.ExpectedGroupRevision = confirmed.Group.Revision
	_, err = h.service.Apply(ctx, confirmed.Group.ID, apply)
	require.Error(t, err, "schema proof is not mixed activation")

	replay, err := h.service.ConfirmManagedSchemaWithStorage(ctx, saved.Group.ID, req, op.OperationID, op.VerifiedDigest, &StorageEvidence{Layout: *scope.GroupLayout, CreatedTable: true})
	require.NoError(t, err)
	require.Equal(t, confirmed.WorkspaceRevision, replay.WorkspaceRevision)
	require.Equal(t, confirmed.Group, replay.Group)
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM managed_schema_operations`))

	if legacy != nil {
		old, err := h.service.Get(ctx, legacy.ID)
		require.NoError(t, err)
		require.Equal(t, legacy, old.Group, "new mixed confirmation must not change legacy profile or consent")
		oldRequest := h.mutation(ctx, old.Group)
		oldRequest.ExpectedGroupRevision = old.Group.Revision
		oldScope, err := h.service.ManagedSchemaScope(ctx, old.Group.ID, oldRequest)
		require.NoError(t, err)
		_, err = ledger.ValidatePreviewTokenForApply(ctx, token.Token, oldScope)
		require.ErrorIs(t, err, recordingplan.ErrPreviewTokenStale, "same physical table never shares record consent")
		return
	}
	if change == "counter policy" {
		_, err = h.db.Exec(`UPDATE measurement_definitions SET definition_revision='definition-2',counter_policy='{"rollover_modulus":"65536","max_rate_per_second":"10"}' WHERE id=?`, *confirmed.Group.Members[0].MeasurementID)
		require.NoError(t, err)
		fresh := h.mutation(ctx, confirmed.Group)
		fresh.ExpectedGroupRevision = confirmed.Group.Revision
		current, err := h.service.ManagedSchemaScope(ctx, confirmed.Group.ID, fresh)
		require.NoError(t, err)
		require.NotEqual(t, scope.SourceDigest, current.SourceDigest, "same measurement ID with changed immutable facts must invalidate consent")
		stale, err := h.service.Readiness(ctx, confirmed.Group.ID)
		require.NoError(t, err)
		require.False(t, stale.SchemaReady)
		_, err = h.service.ConfirmManagedSchemaWithStorage(ctx, confirmed.Group.ID, fresh, op.OperationID, op.VerifiedDigest, &StorageEvidence{Layout: *scope.GroupLayout, CreatedTable: true})
		require.Error(t, err)
		token2, err := ledger.PrepareSchemaPreview(ctx, current, inspect)
		require.NoError(t, err)
		require.Empty(t, token2.Statements)
		op2, _, err := ledger.ApplySchemaPreview(ctx, token2, recordingplan.SchemaApplyTarget{Inspect: inspect, Execute: func(context.Context, []string) (recordingplan.SchemaExecution, error) {
			t.Fatal("compatible fresh proof must not repeat DDL")
			return recordingplan.SchemaExecution{}, nil
		}})
		require.NoError(t, err)
		_, err = h.service.ConfirmManagedSchemaWithStorage(ctx, confirmed.Group.ID, fresh, op2.OperationID, op2.VerifiedDigest, &StorageEvidence{Layout: *current.GroupLayout, CreatedTable: false})
		require.NoError(t, err)
		refreshedReady, err := h.service.Readiness(ctx, confirmed.Group.ID)
		require.NoError(t, err)
		require.True(t, refreshedReady.SchemaReady)
		require.False(t, refreshedReady.Ready)
		require.NotEqual(t, ready.SchemaDigest, refreshedReady.SchemaDigest)
		return
	}
	changed := cloneWriteGroup(confirmed.Group)
	switch change {
	case "method":
		changed.Members[0].Calculation = WriteGroupMemberCalculationLatest
	case "cadence":
		changed.Members[0].Cadence.Interval = 60
	case "required":
		changed.Members[0].Required = false
	case "measurement identity":
		id := seedMeasurement(t, h, h.primary(), "counter")
		_, err = h.db.Exec(`UPDATE measurement_definitions SET unit='kWh',series_epoch='epoch-1',numeric_encoding='uint64',counter_policy='{}' WHERE id=?`, id)
		require.NoError(t, err)
		changed.Members[0].MeasurementID = &id

	}
	update := h.mutation(ctx, changed)
	update.ExpectedGroupRevision = changed.Revision
	updated, err := h.service.Update(ctx, changed.ID, update)
	require.NoError(t, err)
	require.Empty(t, updated.Group.Destination.SchemaRevision)
	require.Empty(t, updated.Group.Destination.SchemaDigest)
	fresh := h.mutation(ctx, updated.Group)
	fresh.ExpectedGroupRevision = updated.Group.Revision
	current, err := h.service.ManagedSchemaScope(ctx, updated.Group.ID, fresh)
	require.NoError(t, err)
	_, err = ledger.ValidatePreviewTokenForApply(ctx, token.Token, current)
	require.ErrorIs(t, err, recordingplan.ErrPreviewTokenStale)
	// A new CAS envelope must not launder an old successful operation into new consent.
	_, err = h.service.ConfirmManagedSchemaWithStorage(ctx, updated.Group.ID, fresh, op.OperationID, op.VerifiedDigest, &StorageEvidence{Layout: *scope.GroupLayout, CreatedTable: true})
	require.Error(t, err, "old ledger proof must not confirm changed calculation with fresh CAS")
	got, err := h.service.Get(ctx, updated.Group.ID)
	require.NoError(t, err)
	require.Empty(t, got.Group.Destination.SchemaRevision)
	require.Equal(t, 0, h.countRows(`SELECT COUNT(*) FROM wg_runtime_versions`))
}

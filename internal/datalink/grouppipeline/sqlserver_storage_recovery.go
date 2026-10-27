package grouppipeline

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

func (r *StorageRecovery) sqlServerRestoreTarget(target groupdelivery.Target, schemaName string) recordingplan.StorageRestoreTarget {
	if r.sqlServerOps != nil {
		return sqlServerStorageRestoreTargetWithOps(target, schemaName, *r.sqlServerOps)
	}
	return sqlServerStorageRestoreTarget(target, schemaName)
}

func sqlServerCapabilityMatchesLayout(capability *dbtarget.SQLServerRecordingCapability, layout *recordingplan.GroupSchemaLayout) bool {
	if capability == nil || layout == nil || layout.SQLServer == nil {
		return false
	}
	e := layout.SQLServer
	return capability.Database == e.Database && capability.Schema == e.Schema && capability.ServerVersion == e.ServerVersion &&
		capability.ServerEdition == e.ServerEdition && capability.EngineEdition == e.EngineEdition &&
		capability.CompatibilityLevel == e.CompatibilityLevel && capability.CapabilityRevision == e.CapabilityRevision
}

func sqlServerRestoreRequest(snap *workspace.WriteGroupAppliedSnapshot) recordingplan.StorageRestoreRequest {
	group, proof := snap.Group, snap.RuntimeLayout.SQLServer
	table, owner, _ := workspace.LongStorageIdentity(group)
	return recordingplan.StorageRestoreRequest{
		Scope: workspace.StorageScopeOf(group), Dialect: string(dbtarget.SQLDialectSQLServer), ConnectorRevision: group.Destination.ConnectorRevision,
		TableName: table, OwnerColumn: owner, GroupID: group.ID,
		SQLServer: &recordingplan.SQLServerStorageRestoreEvidence{
			Group: sqlServerRestoreTargetInspection(proof.Group), Receipt: sqlServerRestoreTargetInspection(proof.Receipt),
		},
	}
}

func (r *StorageRecovery) sqlServerConfirmedFallbackLayout(ctx context.Context, req recordingplan.StorageRestoreRequest) (*recordingplan.GroupSchemaLayout, *recordingplan.StorageLineage, error) {
	lineage, err := recordingplan.ReadStorageLineage(ctx, r.store.DB(), req.Scope)
	if err != nil {
		return nil, nil, err
	}
	native := req.Scope.StorageMode == string(workspace.WriteGroupStorageModeNativeMonthly)
	if lineage == nil || !lineage.Confirmed() || lineage.StorageScope != req.Scope || lineage.Generation < 1 || lineage.Layout == nil || lineage.Layout.SQLServer == nil ||
		lineage.Layout.MySQL != nil || lineage.Layout.Partition != nil || lineage.ConnectorRevision != req.ConnectorRevision ||
		lineage.TableName != req.TableName || lineage.OwnerColumn != req.OwnerColumn ||
		(req.Scope.StorageMode != string(workspace.WriteGroupStorageModePlain) && !native) || req.Scope.LayoutVersion != string(workspace.WriteGroupLayoutManagedLongV1) {
		return nil, lineage, recordingplan.ErrPreviewTokenStale
	}
	fixed, err := recordingplan.ManagedSQLServerLongGroupLayout(req.TableName, req.OwnerColumn, *lineage.Layout.SQLServer)
	if err != nil || !reflect.DeepEqual(fixed, *lineage.Layout) || (fixed.SQLServer.NativeMonthly != nil) != native {
		return nil, lineage, recordingplan.ErrInvalidGroupSchemaLayout
	}
	encoded, err := json.Marshal(fixed)
	if err != nil {
		return nil, lineage, err
	}
	digest := sha256.Sum256(encoded)
	if lineage.LayoutDigest != hex.EncodeToString(digest[:]) {
		return nil, lineage, recordingplan.ErrPreviewTokenStale
	}
	return &fixed, lineage, nil
}

func (r *StorageRecovery) restoreSQLServer(ctx context.Context, item groupdelivery.OutboxItem, target groupdelivery.Target, cause error, snap *workspace.WriteGroupAppliedSnapshot) (groupdelivery.RestoreDecision, error) {
	blocked := groupdelivery.RestoreDecision{Action: groupdelivery.RestoreBlock, Code: recordingplan.StorageCodeRepairBlocked}
	if r.store == nil || r.ledger == nil || target.DB == nil || target.SQLServer == nil || target.SQLServer == snap.RuntimeLayout.SQLServer ||
		!sqlServerFrozenScopeMatches(item, snap) || !sqlServerCapabilityMatchesDestination(target.SQLServerCapability, item.Destination) ||
		!sqlServerCapabilityMatchesProof(target.SQLServerCapability, snap.RuntimeLayout.SQLServer) || !target.SQLServerCapability.SupportsMissingRepair {
		return blocked, nil
	}
	req := sqlServerRestoreRequest(snap)
	req.Deny = r.denial(ctx, snap, item, req.TableName, req.Dialect)
	layout, lineage, err := r.sqlServerConfirmedFallbackLayout(ctx, req)
	if err != nil || lineage.ConfirmedDigest != snap.Group.Destination.SchemaDigest || !r.sqlServerGroupConfirmed(ctx, snap, lineage) {
		return blocked, nil
	}
	req.FallbackLayout = layout
	release := r.lock(req.Scope)
	defer release()
	result, err := r.ledger.RestoreManagedStorage(ctx, req, r.sqlServerRestoreTarget(target, req.Scope.TableSchema))
	if err != nil {
		return groupdelivery.RestoreDecision{}, err
	}
	switch result.Outcome {
	case recordingplan.StorageRestored, recordingplan.StorageNotMissing:
		if result.Outcome == recordingplan.StorageNotMissing && !dbtarget.IsMissingRelationError(cause) && !errors.Is(cause, dbtarget.ErrSQLServerManagedPhysicalMismatch) {
			return groupdelivery.RestoreDecision{}, nil
		}
		current, err := r.sqlServerCurrentProof(ctx, snap, target, false)
		if err != nil {
			return blocked, nil
		}
		// The resolver allocated this attempt's proof independently of the saved
		// runtime snapshot. The sender's retry closure shares only this pointer.
		*target.SQLServer = *current
		return groupdelivery.RestoreDecision{Action: groupdelivery.RestoreRetryNow}, nil
	case recordingplan.StorageInProgress:
		return groupdelivery.RestoreDecision{Action: groupdelivery.RestoreWait}, nil
	default:
		blocked.Code = result.Code
		return blocked, nil
	}
}

func (r *StorageRecovery) sqlServerCurrentProof(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot, target groupdelivery.Target, allowMissing bool) (*dbtarget.SQLServerManagedRowProof, error) {
	if r == nil || r.store == nil || r.ledger == nil || snap == nil || snap.Group == nil || snap.RuntimeLayout == nil ||
		workspace.ValidateSQLServerRuntimeLayout(snap.Group, snap.RuntimeLayout) != nil || target.SQLServerCapability == nil {
		return nil, recordingplan.ErrTargetInspectionUnconfirmed
	}
	req := sqlServerRestoreRequest(snap)
	expectedCapabilityRevision := snap.RuntimeLayout.SQLServer.Group.Evidence.CapabilityRevision
	capabilityRevisionChanged := target.SQLServerCapability.CapabilityRevision != expectedCapabilityRevision
	if capabilityRevisionChanged && (snap.Group.StorageMode != workspace.WriteGroupStorageModeNativeMonthly ||
		!dbtarget.SQLServerNativeDeliveryCapabilityRevisionMatches(target.SQLServerCapability, expectedCapabilityRevision)) {
		return nil, recordingplan.ErrTargetInspectionUnconfirmed
	}
	_, lineage, err := r.sqlServerConfirmedFallbackLayout(ctx, req)
	if err != nil {
		return nil, err
	}
	if lineage.ConfirmedDigest != snap.Group.Destination.SchemaDigest || !r.sqlServerGroupConfirmed(ctx, snap, lineage) {
		return nil, recordingplan.ErrPreviewTokenStale
	}
	if lineage.Generation == 1 && snap.Group.StorageMode != workspace.WriteGroupStorageModeNativeMonthly {
		if err := r.ledger.VerifySQLServerStorageGeneration(ctx, req, req.SQLServer.Group, req.SQLServer.Receipt); err != nil {
			return nil, err
		}
		return dbtarget.CloneSQLServerManagedRowProof(snap.RuntimeLayout.SQLServer), nil
	}
	ops := sqlServerStorageRestoreTargetOps{Inspect: dbtarget.InspectOpenSQLServerManagedTable}
	if r.sqlServerOps != nil {
		ops = *r.sqlServerOps
	}
	if ops.Inspect == nil {
		return nil, recordingplan.ErrTargetInspectionUnconfirmed
	}
	group, err := ops.Inspect(ctx, target.DB, target.SQLServerCapability, req.TableName)
	if err != nil {
		return nil, mapSQLServerStorageRestoreError(err)
	}
	receipt, err := ops.Inspect(ctx, target.DB, target.SQLServerCapability, dbtarget.EffectReceiptTable)
	if err != nil {
		return nil, mapSQLServerStorageRestoreError(err)
	}
	if allowMissing && (group != nil && group.Inspection != nil && group.Inspection.Status == dbtarget.TableInspectionMissing ||
		receipt != nil && receipt.Inspection != nil && receipt.Inspection.Status == dbtarget.TableInspectionMissing) {
		// A proved missing object reaches the normal insert error -> confirmed
		// repair path; it does not authorize or synthesize a new physical proof.
		return dbtarget.CloneSQLServerManagedRowProof(snap.RuntimeLayout.SQLServer), nil
	}
	verifiedGroup := group
	verifiedReceipt := receipt
	if capabilityRevisionChanged {
		verifiedGroup, verifiedReceipt, err = normalizeSQLServerDeliveryInspectionsForVerification(group, receipt, expectedCapabilityRevision)
		if err != nil {
			return nil, err
		}
	}
	groupState := sqlServerRestoreTargetInspection(verifiedGroup)
	receiptState := sqlServerRestoreTargetInspection(verifiedReceipt)
	if err := r.ledger.VerifySQLServerStorageGeneration(ctx, req, groupState, receiptState); err != nil {
		return nil, err
	}
	return dbtarget.CloneSQLServerManagedRowProof(&dbtarget.SQLServerManagedRowProof{
		ConnectorID: req.Scope.ConnectorID, ConnectorRevision: req.ConnectorRevision, Group: group, Receipt: receipt,
	}), nil
}

func cloneSQLServerInspectionForDelivery(inspection *dbtarget.SQLServerManagedInspection) *dbtarget.SQLServerManagedInspection {
	if inspection == nil {
		return nil
	}
	return dbtarget.CloneSQLServerManagedRowProof(&dbtarget.SQLServerManagedRowProof{Group: inspection}).Group
}

func normalizeSQLServerDeliveryInspectionsForVerification(group, receipt *dbtarget.SQLServerManagedInspection, expectedRevision string) (normalizedGroup, normalizedReceipt *dbtarget.SQLServerManagedInspection, err error) {
	if group == nil || group.Evidence == nil || receipt == nil || receipt.Evidence == nil || expectedRevision == "" {
		return nil, nil, recordingplan.ErrTargetInspectionUnconfirmed
	}
	groupClone := cloneSQLServerInspectionForDelivery(group)
	receiptClone := cloneSQLServerInspectionForDelivery(receipt)
	if groupClone == nil || groupClone.Evidence == nil || receiptClone == nil || receiptClone.Evidence == nil {
		return nil, nil, recordingplan.ErrTargetInspectionUnconfirmed
	}
	groupClone.Evidence.CapabilityRevision = expectedRevision
	receiptClone.Evidence.CapabilityRevision = expectedRevision
	return groupClone, receiptClone, nil
}

func (r *StorageRecovery) sqlServerGroupConfirmed(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot, lineage *recordingplan.StorageLineage) bool {
	if snap == nil || snap.Group == nil || lineage == nil {
		return false
	}
	if lineage.ConfirmedOperationID != "" && lineage.ConfirmedOperationID == snap.Group.Destination.SchemaRevision {
		return true
	}
	return r.groupConfirmed(ctx, snap)
}

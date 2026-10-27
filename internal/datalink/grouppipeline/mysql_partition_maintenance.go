package grouppipeline

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// MaintainMySQLPartitions runs one bounded maintenance pass over applied native
// scopes. It shares the existing pipeline lifetime and operation ledger, while
// its asynchronous invocation leaves sample-window closure and delivery running.
func (p *Pipeline) MaintainMySQLPartitions(ctx context.Context) error {
	if p.recovery == nil || p.deps.Groups == nil || p.deps.Destinations == nil || !p.partitionMu.TryLock() {
		return nil
	}
	defer p.partitionMu.Unlock()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	listed, err := p.deps.Groups.List(ctx)
	if err != nil {
		return err
	}
	seen := map[recordingplan.StorageScope]bool{}
	var failures []error
	for _, group := range listed.Groups {
		if err := ctx.Err(); err != nil {
			return errors.Join(append(failures, err)...)
		}
		if group == nil || group.AppliedRevision == "" ||
			group.StorageMode != workspace.WriteGroupStorageModeNativeMonthly && group.Revision == group.AppliedRevision {
			continue
		}
		snapshot, err := p.deps.Groups.ResolveAppliedAt(ctx, group.ID, p.config.Now())
		if err != nil || snapshot == nil {
			failures = append(failures, errors.New("native partition maintenance: applied snapshot unavailable"))
			continue
		}
		frozen, err := p.loadRuntimeSnapshot(ctx, groupKey(snapshot))
		if err != nil || frozen == nil {
			failures = append(failures, errors.New("native partition maintenance: frozen proof unavailable"))
			continue
		}
		if frozen.Group.StorageMode != workspace.WriteGroupStorageModeNativeMonthly || frozen.RuntimeLayout.Dialect != dbtarget.SQLDialectMySQL {
			continue
		}
		scope := workspace.StorageScopeOf(frozen.Group)
		if seen[scope] {
			continue
		}
		seen[scope] = true
		if err := p.maintainMySQLPartitionScope(ctx, frozen); err != nil {
			failures = append(failures, err)
		}
	}
	return errors.Join(failures...)
}

func (p *Pipeline) maintainMySQLPartitionScope(ctx context.Context, frozen *workspace.WriteGroupAppliedSnapshot) (err error) {
	group := frozen.Group
	opened, err := p.deps.Destinations.OpenDestination(ctx, group.Destination.ConnectorID, group.Destination.ConnectorRevision)
	if err != nil {
		return err
	}
	if opened == nil {
		return dbtarget.ErrMySQLManagedServerMismatch
	}
	if opened.Close != nil {
		defer func() { err = errors.Join(err, opened.Close()) }()
	}
	if opened.Kind != schema.DatabaseConnectorKindMySQL || opened.DB == nil ||
		!mysqlCapabilityMatchesProof(opened.MySQL, frozen.RuntimeLayout.MySQL) {
		return dbtarget.ErrMySQLManagedServerMismatch
	}
	target := recordingplan.MySQLPartitionMaintenanceTarget{
		Inspect: func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
			if table != group.Destination.TableName && table != dbtarget.EffectReceiptTable {
				return recordingplan.TargetTableInspection{}, recordingplan.ErrTargetInspectionUnconfirmed
			}
			physical, err := dbtarget.InspectOpenMySQLManagedTable(ctx, opened.DB, opened.MySQL, group.Destination.Database, table)
			return mysqlRestoreTargetInspection(physical), mapMySQLStorageRestoreError(err)
		},
		Reorganize: func(ctx context.Context, before, after *recordingplan.MySQLPartitionEvidence) error {
			proof := mysqlProofAtPartitions(frozen.RuntimeLayout.MySQL, before)
			return mapMySQLStorageRestoreError(dbtarget.ReorganizeOpenMySQLFuturePartition(ctx, opened.DB, opened.MySQL,
				group.Destination.Database, proof.Group, mysqlPartitionEvidenceForDB(after)))
		},
	}
	result, err := p.recovery.ledger.MaintainMySQLPartitions(ctx, recordingplan.MySQLPartitionMaintenanceRequest{
		Scope: workspace.StorageScopeOf(group), GroupID: group.ID,
		ConnectorRevision: group.Destination.ConnectorRevision, At: p.config.Now(),
	}, target)
	if err != nil {
		return err
	}
	if result.Code != "" {
		return fmt.Errorf("%s: %s", result.Code, result.Reason)
	}
	return nil
}

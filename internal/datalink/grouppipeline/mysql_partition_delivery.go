package grouppipeline

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"modernc.org/sqlite"
)

var errMySQLPartitionProofReadTransient = errors.New("mysql partition proof read temporarily unavailable")

// mysqlDeliveryProof keeps the accepted payload and immutable runtime descriptor
// intact while binding the SQL transaction to the ledger-authorized current
// partition catalog. Every other physical fact still comes from that descriptor.
func (p *Pipeline) mysqlDeliveryProof(ctx context.Context, frozen *workspace.WriteGroupAppliedSnapshot) (*dbtarget.MySQLManagedRowProof, error) {
	proof := frozen.RuntimeLayout.MySQL
	if frozen.Group.StorageMode != workspace.WriteGroupStorageModeNativeMonthly {
		return proof, nil
	}
	if p.deps.Store == nil {
		return nil, recordingplan.ErrMySQLPartitionMaintenanceStateMismatch
	}
	scope := workspace.StorageScopeOf(frozen.Group)
	lineage, err := recordingplan.ReadStorageLineage(ctx, p.deps.Store.DB(), scope)
	if err != nil {
		return nil, mysqlPartitionProofReadError(err)
	}
	if !lineage.Confirmed() || lineage.ConnectorRevision != frozen.Group.Destination.ConnectorRevision ||
		lineage.ConfirmedDigest != frozen.Group.Destination.SchemaDigest {
		return nil, recordingplan.ErrMySQLPartitionMaintenanceStateMismatch
	}
	state, err := recordingplan.ReadMySQLPartitionMaintenanceState(ctx, p.deps.Store.DB(), scope)
	if err != nil {
		return nil, mysqlPartitionProofReadError(err)
	}
	actual := mysqlRestoreTargetInspection(proof.Group)
	if state != nil {
		actual.MySQL.Partitions = state.Current
	}
	if _, err := recordingplan.NormalizeMySQLPartitionCatalog(ctx, p.deps.Store.DB(), scope, *lineage.Layout, actual); err != nil {
		if errors.Is(err, recordingplan.ErrIncompatibleExistingTable) {
			return nil, dbtarget.ErrMySQLManagedPartitionChanged
		}
		return nil, mysqlPartitionProofReadError(err)
	}
	return mysqlProofAtPartitions(proof, actual.MySQL.Partitions), nil
}

func mysqlPartitionProofReadError(err error) error {
	var sqliteErr *sqlite.Error
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) ||
		errors.As(err, &sqliteErr) && dbtarget.ClassifyInsertError(sqliteErr) == dbtarget.InsertErrorTransient {
		return fmt.Errorf("%w: %w", errMySQLPartitionProofReadTransient, err)
	}
	return err
}

func mysqlProofAtPartitions(proof *dbtarget.MySQLManagedRowProof, current *recordingplan.MySQLPartitionEvidence) *dbtarget.MySQLManagedRowProof {
	result := *proof
	group := *proof.Group
	evidence := *proof.Group.Evidence
	evidence.Partitions = mysqlPartitionEvidenceForDB(current)
	group.Evidence = &evidence
	result.Group = &group
	return &result
}

func mysqlPartitionEvidenceForDB(source *recordingplan.MySQLPartitionEvidence) *dbtarget.MySQLPartitionEvidence {
	if source == nil {
		return nil
	}
	result := &dbtarget.MySQLPartitionEvidence{Method: source.Method, Expression: source.Expression,
		Partitions: make([]dbtarget.MySQLPartitionBoundary, len(source.Partitions))}
	for i, partition := range source.Partitions {
		result.Partitions[i] = dbtarget.MySQLPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return result
}

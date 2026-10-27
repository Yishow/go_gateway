package grouppipeline

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"modernc.org/sqlite"
)

var errSQLServerProofReadTransient = errors.New("sqlserver current physical proof is temporarily unavailable")

// sqlServerDeliveryProof leaves the accepted snapshot unchanged while binding
// this delivery attempt to the current ledger-authorized physical generation.
func (p *Pipeline) sqlServerDeliveryProof(ctx context.Context, frozen *workspace.WriteGroupAppliedSnapshot, opened *dbtarget.OpenedDestination) (*dbtarget.SQLServerManagedRowProof, error) {
	if p.recovery == nil || opened == nil {
		return nil, groupdelivery.ErrTargetBlocked
	}
	proof, err := p.recovery.sqlServerCurrentProof(ctx, frozen, groupdelivery.Target{
		DB: opened.DB, Kind: opened.Kind, SQLServerCapability: opened.SQLServer,
	}, true)
	var sqliteErr *sqlite.Error
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, dbtarget.ErrSQLServerManagedTableUnconfirmed) ||
		errors.As(err, &sqliteErr) && dbtarget.ClassifyInsertError(sqliteErr) == dbtarget.InsertErrorTransient {
		return nil, fmt.Errorf("%w: %w", errSQLServerProofReadTransient, err)
	}
	return proof, err
}

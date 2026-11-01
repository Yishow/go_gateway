package workspace

import (
	"context"
	"database/sql"
	"errors"

	"go-gateway/internal/datalink/sqltx"
)

// AcceptAcquisitionSourceGraph revalidates original SQL evidence while pinning
// source writers. The callback must publish only through nonblocking owned CAS;
// no I/O or blocking lifecycle/metadata lock may run inside it.
func (s *Service) AcceptAcquisitionSourceGraph(ctx context.Context, graph AcquisitionSourceGraph, stage func() error) error {
	return s.AcceptAcquisitionSourceGraphs(ctx, []AcquisitionSourceGraph{graph}, stage)
}

// AcceptAcquisitionSourceGraphs accepts a complete subset of device graphs in
// one original SQL reservation. Unknown devices must remain unknown at runtime;
// an empty set is not evidence and must use the ordinary no-gap install path.
func (s *Service) AcceptAcquisitionSourceGraphs(ctx context.Context, graphs []AcquisitionSourceGraph, stage func() error) error {
	if s == nil {
		return ErrAcquisitionSourceGraph
	}
	acceptor, ok := s.repo.(interface {
		AcceptAcquisitionSourceGraphs(context.Context, []AcquisitionSourceGraph, func() error) error
	})
	if !ok {
		return ErrAcquisitionSourceGraph
	}
	return acceptor.AcceptAcquisitionSourceGraphs(ctx, graphs, stage)
}

// AcceptAcquisitionSourceGraph uses the already published SQLite zero-row writer
// reservation pattern without borrowing any point-format operation authority.
// Any commit/cleanup failure forbids success; the caller must revoke only the
// exact geometry receipt produced by its callback, never a newer publication.
func (r *SQLRepository) AcceptAcquisitionSourceGraph(ctx context.Context, graph AcquisitionSourceGraph, stage func() error) error {
	return r.AcceptAcquisitionSourceGraphs(ctx, []AcquisitionSourceGraph{graph}, stage)
}

// AcceptAcquisitionSourceGraphs checks every original witness before invoking
// stage exactly once. Duplicate devices or foreign owners/workspaces reject the
// entire set; it cannot partially publish the graphs that happened to pass.
func (r *SQLRepository) AcceptAcquisitionSourceGraphs(ctx context.Context, graphs []AcquisitionSourceGraph, stage func() error) (resultErr error) {
	if r == nil || r.db == nil || stage == nil || len(graphs) == 0 {
		return ErrAcquisitionSourceGraph
	}
	devices := make(map[string]bool, len(graphs))
	workspaceID := ""
	for _, graph := range graphs {
		if !graph.Complete() || graph.data.owner != r || graph.data.workspaceID == "" || graph.data.witness == "" || graph.data.semantic == "" || graph.data.original == "" || graph.data.input.Device.ID == "" || devices[graph.data.input.Device.ID] {
			return ErrAcquisitionSourceGraph
		}
		if workspaceID != "" && graph.data.workspaceID != workspaceID {
			return ErrAcquisitionSourceGraph
		}
		workspaceID = graph.data.workspaceID
		devices[graph.data.input.Device.ID] = true
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	owned, err := sqltx.Begin(ctx, r.db, nil)
	if err != nil {
		return acquisitionGraphTransactionError(err)
	}
	defer func() { //nolint:contextcheck // sqltx.Rollback intentionally uses independent bounded cleanup after caller cancellation.
		if err := owned.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			resultErr = acquisitionGraphTransactionError(errors.Join(resultErr, err))
		}
	}()
	var sqliteVersion string
	if err := owned.Tx.QueryRowContext(ctx, `SELECT sqlite_version()`).Scan(&sqliteVersion); err != nil || sqliteVersion == "" {
		return ErrAcquisitionSourceGraph
	}
	if _, err := owned.Tx.ExecContext(ctx, `UPDATE system_settings SET value=value WHERE 0`); err != nil {
		return acquisitionGraphTransactionError(err)
	}
	for _, graph := range graphs {
		if err := r.validateAcquisitionSourceGraphTx(ctx, owned.Tx, workspaceID, graph.data.input, graph.data.witness); err != nil {
			return err
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := stage(); err != nil {
		return err
	}
	if err := owned.Commit(); err != nil { //nolint:contextcheck // sqltx.Commit owns independent bounded cleanup; its terminal API has no context parameter.
		return acquisitionGraphTransactionError(err)
	}
	return nil
}

package workspace

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/sqltx"
)

// ErrAcquisitionSourceGraph is a closed failure, never a raw SQL diagnostic.
var ErrAcquisitionSourceGraph = errors.New("acquisition source evidence unavailable or changed")

// ReadAcquisitionSourceGraph compares the original source projection against one
// canonical SQL snapshot. It does not refresh/rebase projection or install state.
func (r *SQLRepository) ReadAcquisitionSourceGraph(ctx context.Context, projection *RuntimeProjection, deviceID string) (out AcquisitionSourceGraph, resultErr error) {
	input, err := projectionAcquisitionRows(projection, deviceID)
	if err != nil || r == nil || r.db == nil {
		return AcquisitionSourceGraph{}, ErrAcquisitionSourceGraph
	}
	owned, err := sqltx.Begin(ctx, r.db, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return AcquisitionSourceGraph{}, ErrAcquisitionSourceGraph
	}
	defer func() { //nolint:contextcheck // sqltx.Rollback intentionally uses independent bounded cleanup after caller cancellation.
		if err := owned.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			out = AcquisitionSourceGraph{}
			resultErr = acquisitionGraphTransactionError(err)
		}
	}()
	witness, err := sourcerule.AcquisitionSourceDigest(input, true)
	if err != nil || r.validateAcquisitionSourceGraphTx(ctx, owned.Tx, projection.WorkspaceID, input, witness) != nil {
		return AcquisitionSourceGraph{reason: "source-snapshot-changed"}, ErrAcquisitionSourceGraph
	}
	excluded, reason := sourcerule.AcquisitionSourceExclusions(input)
	if reason != "" {
		return AcquisitionSourceGraph{reason: reason}, ErrAcquisitionSourceGraph
	}
	semantic, err := sourcerule.AcquisitionSourceDigest(input, false)
	if err != nil {
		return AcquisitionSourceGraph{}, ErrAcquisitionSourceGraph
	}
	original, err := acquisitionProjectionDigest(input)
	if err != nil {
		return AcquisitionSourceGraph{}, ErrAcquisitionSourceGraph
	}
	ruleSemantics := map[string]string{}
	for _, rule := range input.Rules {
		part := sourcerule.AcquisitionSourceRows{Device: input.Device, Rules: []*schema.SourceRule{rule}}
		ids := map[string]bool{}
		for _, link := range input.Links {
			if link.RuleID == rule.ID {
				part.Links = append(part.Links, link)
				ids[link.PointID] = true
			}
		}
		for _, point := range input.Points {
			if ids[point.ID] {
				part.Points = append(part.Points, point)
			}
		}
		digest, err := sourcerule.AcquisitionSourceDigest(part, false)
		if err != nil {
			return AcquisitionSourceGraph{}, ErrAcquisitionSourceGraph
		}
		ruleSemantics[rule.ID] = digest
	}
	out.data = &acquisitionSourceGraphData{owner: r, workspaceID: projection.WorkspaceID, witness: witness, semantic: semantic, original: original, input: input, excluded: excluded, ruleSemantics: ruleSemantics}
	return out, nil
}

func (r *SQLRepository) validateAcquisitionSourceGraphTx(ctx context.Context, tx *sql.Tx, workspaceID string, input sourcerule.AcquisitionSourceRows, witness string) error {
	record, err := readWorkspaceRecord(ctx, tx)
	if err != nil || record.ID != workspaceID || !slices.Contains(record.OrderedDeviceIDs, input.Device.ID) {
		return ErrAcquisitionSourceGraph
	}
	current, err := sourcerule.ReadAcquisitionSourceRowsTx(ctx, tx, input.Device.ID)
	if err != nil {
		return ErrAcquisitionSourceGraph
	}
	actual, err := sourcerule.AcquisitionSourceDigest(current, true)
	if err != nil || actual != witness {
		return ErrAcquisitionSourceGraph
	}
	return nil
}

type acquisitionGraphStorageError struct{ cause error }

func (e *acquisitionGraphStorageError) Error() string { return ErrAcquisitionSourceGraph.Error() }
func (e *acquisitionGraphStorageError) Unwrap() []error {
	return []error{ErrAcquisitionSourceGraph, e.cause}
}
func acquisitionGraphTransactionError(cause error) error {
	return &acquisitionGraphStorageError{cause: cause}
}

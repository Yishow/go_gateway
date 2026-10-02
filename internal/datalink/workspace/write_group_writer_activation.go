package workspace

import (
	"context"
	"database/sql"
	"time"
)

// WriterOwnershipActivationRequest identifies the local lifecycle transition
// prepared by Apply. EffectiveAt is always normalized to UTC by Apply.
type WriterOwnershipActivationRequest struct {
	WorkspaceID                   string    `json:"workspace_id"`
	GroupID                       string    `json:"group_id"`
	ExpectedDatabaseSetupRevision string    `json:"expected_database_setup_revision"`
	ExpectedGroupRevision         string    `json:"expected_group_revision"`
	ExpectedConnectorRevision     string    `json:"expected_connector_revision"`
	PreviousAppliedRevision       string    `json:"previous_applied_revision"`
	AppliedRevision               string    `json:"applied_revision"`
	EffectiveAt                   time.Time `json:"effective_at"`
}

// WriterOwnershipActivationBarrier is the transaction-scoped seam for the
// later writer-owner consumer. Implementations must not commit, begin another
// transaction, start or stop runtime, or modify Share state.
type WriterOwnershipActivationBarrier interface {
	ActivateInTx(context.Context, *sql.Tx, WriterOwnershipActivationRequest) error
}

// WithWriterOwnershipActivationBarrier installs the optional preparation
// barrier. A nil barrier keeps Apply domain-only and preserves the current
// writer path until a production consumer is explicitly connected.
func (s *WriteGroupService) WithWriterOwnershipActivationBarrier(
	barrier WriterOwnershipActivationBarrier,
) *WriteGroupService {
	if s != nil {
		s.writerOwnershipActivation = barrier
	}
	return s
}

func (s *WriteGroupService) activateWriterOwnershipInTx(
	ctx context.Context,
	tx *sql.Tx,
	request WriterOwnershipActivationRequest,
) error {
	if s == nil || s.writerOwnershipActivation == nil {
		return nil
	}
	return s.writerOwnershipActivation.ActivateInTx(ctx, tx, request)
}

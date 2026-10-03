package recordingplan

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
)

// RecordingStartAction separates a start intent from schema and test-write effects.
const RecordingStartAction = "recording_start"

type recordingStartResumer interface {
	ResumeRecordingStartOperation(context.Context, string, string, string) (*SchemaOperation, error)
}

func recordingStartOperationID(workspaceID, requestID string) string {
	identity := sha256.Sum256([]byte(workspaceID + "\x00" + requestID))
	return "start-" + hex.EncodeToString(identity[:])
}

// FindRecordingStart looks up a recorded request without touching live resources.
// Its scope and digest must match even when the result is already complete.
func (s *Service) FindRecordingStart(ctx context.Context, workspaceID, requestID, digest string) (*SchemaOperation, error) {
	if s == nil || s.repo == nil {
		return nil, ErrSchemaOperationResult
	}
	op, err := s.GetSchemaOperation(ctx, workspaceID, recordingStartOperationID(workspaceID, requestID))
	if err != nil {
		return nil, err
	}
	if op == nil || op.Action != RecordingStartAction || op.PayloadDigest != digest {
		return nil, ErrSchemaOperationMismatch
	}
	return op, nil
}

// ClaimRecordingStart records a request in the existing operation ledger.
// Request identity is namespaced by workspace; it is not a schema preview token.
func (s *Service) ClaimRecordingStart(ctx context.Context, workspaceID, requestID, digest, detail string) (*SchemaOperation, ClaimOutcome, error) {
	if s == nil || s.repo == nil || strings.TrimSpace(workspaceID) == "" ||
		strings.TrimSpace(requestID) == "" || len(requestID) > 128 || len(digest) != 64 {
		return nil, "", ErrSchemaOperationResult
	}
	if _, err := hex.DecodeString(digest); err != nil {
		return nil, "", ErrSchemaOperationResult
	}
	id := recordingStartOperationID(workspaceID, requestID)
	owner, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create recording start owner: %w", err)
	}
	now := time.Now().UTC()
	op, outcome, err := s.repo.ClaimSchemaOperation(ctx, &SchemaOperation{
		OperationID: id, Token: RecordingStartAction + ":" + id,
		WorkspaceID: workspaceID, ScopeKey: RecordingStartAction + ":" + workspaceID,
		Owner: owner, Action: RecordingStartAction, Status: SchemaOperationRunning,
		PayloadDigest: digest, Detail: detail, CreatedAt: now, UpdatedAt: now,
	})
	if err != nil || op == nil {
		return op, outcome, err
	}
	if outcome != ClaimScopeBusy && (op.Action != RecordingStartAction || op.PayloadDigest != digest) {
		return nil, "", ErrSchemaOperationMismatch
	}
	return op, outcome, nil
}

// ResumeRecordingStart grants an execution owner only after the application
// has revalidated the persisted intent and its recorded progress.
func (s *Service) ResumeRecordingStart(ctx context.Context, workspaceID, operationID, digest string) (*SchemaOperation, ClaimOutcome, error) {
	op, err := s.repo.GetSchemaOperation(ctx, workspaceID, operationID)
	if err != nil {
		return nil, "", err
	}
	if op.Action != RecordingStartAction || op.PayloadDigest != digest {
		return nil, "", ErrSchemaOperationMismatch
	}
	if op.Status == SchemaOperationSucceeded {
		return withoutOwner(*op), ClaimCompleted, nil
	}
	owner, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create recording start resume owner: %w", err)
	}
	var acquired *SchemaOperation
	if op.Status.IsActive() {
		repository, ok := s.repo.(operationExecutionRepository)
		if !ok {
			return nil, "", ErrSchemaOperationResult
		}
		acquired, err = repository.TakeOverSchemaOperation(ctx, operationID, owner, time.Now().UTC().Add(-SchemaOperationLease))
	} else {
		repository, ok := s.repo.(recordingStartResumer)
		if !ok {
			return nil, "", ErrSchemaOperationResult
		}
		acquired, err = repository.ResumeRecordingStartOperation(ctx, operationID, owner, digest)
	}
	if err != nil {
		return nil, "", err
	}
	if acquired != nil {
		return acquired, ClaimAcquired, nil
	}
	current, err := s.GetSchemaOperation(ctx, workspaceID, operationID)
	if err != nil {
		return nil, "", err
	}
	return current, claimOutcomeFor(current), nil
}

// SaveRecordingStartProgress renews the owner's lease while recording progress.
func (s *Service) SaveRecordingStartProgress(ctx context.Context, operationID, owner, detail string) error {
	repository, ok := s.repo.(operationExecutionRepository)
	if !ok {
		return ErrSchemaOperationResult
	}
	return repository.SaveSchemaOperationProgress(ctx, operationID, owner, detail)
}

// SaveRecordingStartProgressInTx atomically records a local group save/apply
// result with its workspace revision. The transaction is the configuration
// store's transaction; no external database participates.
func (s *Service) SaveRecordingStartProgressInTx(ctx context.Context, tx *sql.Tx, operationID, owner, detail string) error {
	if tx == nil || owner == "" {
		return ErrSchemaOperationNotOwned
	}
	result, err := tx.ExecContext(ctx, adaptPlaceholders(`UPDATE managed_schema_operations
		SET detail=$1,updated_at=$2 WHERE operation_id=$3 AND owner=$4
		AND action=$5 AND status IN ('pending','running')`),
		detail, time.Now().UTC(), operationID, owner, RecordingStartAction)
	if err != nil {
		return fmt.Errorf("record recording start transaction progress: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read recording start transaction progress: %w", err)
	}
	if affected != 1 {
		return ErrSchemaOperationNotOwned
	}
	return nil
}

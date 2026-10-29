package recordingplan

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// ErrTestWriteRecoveryUnavailable means the ledger cannot establish a complete,
// trustworthy recovery view. It never means that no write occurred.
var ErrTestWriteRecoveryUnavailable = errors.New("test write recovery unavailable")

const testWriteRecoveryLimit = 1000

type testWriteRecoveryRecord struct {
	operation SchemaOperation
	token     *SchemaPreviewToken
}

type testWriteRecoveryRepository interface {
	readTestWriteRecovery(context.Context, string) ([]testWriteRecoveryRecord, error)
}

// TestWriteRecoveryList is a complete workspace-scoped ledger observation, not
// authorization to repeat a write. Failed observations return no partial list.
type TestWriteRecoveryList struct {
	WorkspaceID string                       `json:"workspace_id"`
	Complete    bool                         `json:"complete"`
	Operations  []TestWriteRecoveryOperation `json:"operations"`
}

// TestWriteRecoveryOperation contains only immutable locators and validated
// outcomes. Tokens, execution proofs and destination details remain server-side.
type TestWriteRecoveryOperation struct {
	OperationID       string                `json:"operation_id"`
	WorkspaceID       string                `json:"workspace_id"`
	WorkspaceRevision string                `json:"workspace_revision"`
	GroupID           string                `json:"group_id"`
	GroupRevision     string                `json:"group_revision"`
	ConnectorID       string                `json:"connector_id"`
	ConnectorRevision string                `json:"connector_revision"`
	TargetFingerprint string                `json:"target_fingerprint"`
	Action            string                `json:"action"`
	Status            SchemaOperationStatus `json:"status"`
	WriteOutcome      string                `json:"write_outcome,omitempty"`
	CleanupStatus     string                `json:"cleanup_status,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	CompletedAt       *time.Time            `json:"completed_at,omitempty"`
}

// ListUnresolvedTestWrites recovers operations independently of browser state or
// mutable groups. It never claims operations, opens destinations or repairs data.
func (s *Service) ListUnresolvedTestWrites(ctx context.Context, workspaceID string) (*TestWriteRecoveryList, error) {
	reader, ok := s.repo.(testWriteRecoveryRepository)
	if !ok || strings.TrimSpace(workspaceID) == "" {
		return nil, ErrTestWriteRecoveryUnavailable
	}
	records, err := reader.readTestWriteRecovery(ctx, workspaceID)
	if err != nil || len(records) > testWriteRecoveryLimit {
		return nil, ErrTestWriteRecoveryUnavailable
	}
	result := &TestWriteRecoveryList{WorkspaceID: workspaceID, Complete: true, Operations: make([]TestWriteRecoveryOperation, 0)}
	for _, record := range records {
		op, token := record.operation, record.token
		if !validRecoveryBinding(workspaceID, op, token) {
			return nil, ErrTestWriteRecoveryUnavailable
		}
		if recoveryResolved(op) {
			continue
		}
		// Hash the saved target and connector revision, never an editable current target.
		target, _ := json.Marshal([]string{token.WorkspaceID, token.ConnectorID, token.ConnectorRevision, token.Dialect, token.Database, token.Schema, token.TablePrefix})
		sum := sha256.Sum256(target)
		result.Operations = append(result.Operations, TestWriteRecoveryOperation{
			OperationID: op.OperationID, WorkspaceID: workspaceID, WorkspaceRevision: token.WorkspaceRevision,
			GroupID: token.PlanID, GroupRevision: token.PlanRevision, ConnectorID: token.ConnectorID, ConnectorRevision: token.ConnectorRevision,
			TargetFingerprint: hex.EncodeToString(sum[:]), Action: TestWriteAction, Status: op.Status,
			WriteOutcome: op.WriteOutcome, CleanupStatus: op.CleanupStatus, CreatedAt: op.CreatedAt, UpdatedAt: op.UpdatedAt, CompletedAt: op.CompletedAt,
		})
	}
	return result, nil
}

func validRecoveryBinding(workspaceID string, op SchemaOperation, token *SchemaPreviewToken) bool {
	if !validRecoveryShape(op) || token == nil || op.WorkspaceID != workspaceID || token.WorkspaceID != workspaceID || op.Action != TestWriteAction || token.Action != TestWriteAction ||
		op.OperationID == "" || token.OperationID != op.OperationID || token.Token != op.Token || op.ScopeKey != token.ScopeKey() ||
		token.PlanID == "" || token.PlanRevision == "" || token.WorkspaceRevision == "" || token.ConnectorID == "" || token.ConnectorRevision == "" || token.TablePrefix == "" ||
		op.CreatedAt.IsZero() || op.UpdatedAt.IsZero() {
		return false
	}
	if recoveryCleanupClaim(op) && token.Dialect != "sqlserver" {
		return false
	}
	switch token.Dialect {
	case "sqlite", "postgres", "mysql", "sqlserver":
	default:
		return false
	}
	digest, err := previewDigest(token)
	if err != nil || subtle.ConstantTimeCompare([]byte(digest), []byte(token.Digest)) != 1 ||
		len(token.Statements) != 1 || !strings.HasPrefix(token.Statements[0], testWriteContentPrefix) || TestWriteContentDigest(token) == "" || op.PayloadDigest != TestWriteContentDigest(token) {
		return false
	}
	switch op.Status {
	case SchemaOperationPending, SchemaOperationRunning, SchemaOperationSucceeded, SchemaOperationPartial, SchemaOperationFailed, SchemaOperationUnknown:
	default:
		return false
	}
	switch op.WriteOutcome {
	case "", "written_verified", "written_unverified", "failed", "unknown":
	default:
		return false
	}
	switch op.CleanupStatus {
	case "", "not_attempted", "cleaned", "failed", "unknown":
	default:
		return false
	}
	return true
}

// validRecoveryShape matches the existing operation contract; contradictory
// facts are an unavailable observation, not evidence that the operation vanished.
func validRecoveryShape(op SchemaOperation) bool {
	if op.OperationID == "" || op.PayloadDigest == "" || op.CreatedAt.IsZero() || op.UpdatedAt.IsZero() || op.UpdatedAt.Before(op.CreatedAt) {
		return false
	}
	if op.Status.IsActive() {
		return (op.WriteOutcome == "" && op.CleanupStatus == "" && op.CompletedAt == nil) || recoveryCleanupClaim(op)
	}
	if op.CompletedAt == nil || op.CompletedAt.IsZero() || op.CompletedAt.Before(op.CreatedAt) || op.CleanupStatus == "" {
		return false
	}
	switch op.WriteOutcome {
	case "written_verified":
		return op.Status == SchemaOperationSucceeded && op.CleanupStatus != "not_attempted"
	case "written_unverified":
		return op.Status == SchemaOperationPartial
	case "failed":
		return op.Status == SchemaOperationFailed && op.CleanupStatus == "not_attempted"
	case "unknown":
		return op.Status == SchemaOperationUnknown
	default:
		return false
	}
}

// A committed but unverified write with confirmed cleanup has no pending work;
// its readback failure remains available through the ordinary operation GET.
func recoveryResolved(op SchemaOperation) bool {
	return validRecoveryShape(op) && ((op.Status == SchemaOperationSucceeded && op.WriteOutcome == "written_verified" && op.CleanupStatus == "cleaned") ||
		(op.Status == SchemaOperationPartial && op.WriteOutcome == "written_unverified" && op.CleanupStatus == "cleaned") ||
		(op.Status == SchemaOperationFailed && op.WriteOutcome == "failed" && op.CleanupStatus == "not_attempted"))
}

// recoveryCleanupClaim recognizes only the existing cleanup-only reclaim shape.
// Its previous completion/result are retained evidence, not a new insert result.
func recoveryCleanupClaim(op SchemaOperation) bool {
	return op.Status == SchemaOperationRunning && op.WriteOutcome == "written_verified" &&
		(op.CleanupStatus == "failed" || op.CleanupStatus == "unknown") && op.CompletedAt != nil &&
		!op.CompletedAt.IsZero() && !op.CompletedAt.Before(op.CreatedAt) && !op.UpdatedAt.Before(*op.CompletedAt) && strings.TrimSpace(op.Detail) != ""
}

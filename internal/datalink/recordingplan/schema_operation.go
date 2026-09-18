package recordingplan

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
)

// SchemaOperationStatus is the durable state of one confirmed schema change.
type SchemaOperationStatus string

const (
	SchemaOperationPending   SchemaOperationStatus = "pending"
	SchemaOperationRunning   SchemaOperationStatus = "running"
	SchemaOperationSucceeded SchemaOperationStatus = "succeeded"
	SchemaOperationPartial   SchemaOperationStatus = "partial"
	SchemaOperationFailed    SchemaOperationStatus = "failed"
	SchemaOperationUnknown   SchemaOperationStatus = "unknown"
)

// IsActive reports whether the operation may still be executing and so holds its scope.
func (s SchemaOperationStatus) IsActive() bool {
	return s == SchemaOperationPending || s == SchemaOperationRunning
}

// ClaimOutcome says how a confirmation relates to the operation ledger.
type ClaimOutcome string

const (
	// ClaimAcquired means this caller now owns execution of the operation.
	ClaimAcquired ClaimOutcome = "acquired"
	// ClaimInProgress means the same operation is still pending or running.
	ClaimInProgress ClaimOutcome = "in_progress"
	// ClaimCompleted means the same operation already ended; its result is retained.
	ClaimCompleted ClaimOutcome = "completed"
	// ClaimScopeBusy means another operation holds the same target scope.
	ClaimScopeBusy ClaimOutcome = "scope_busy"
	// ClaimNone means no operation exists yet for the confirmation.
	ClaimNone ClaimOutcome = "none"
)

var (
	ErrSchemaOperationNotFound = errors.New("schema operation not found")
	ErrSchemaOperationMismatch = errors.New("schema operation does not belong to the preview token")
	ErrSchemaOperationNotOwned = errors.New("schema operation is not owned by this execution")
	ErrSchemaOperationResult   = errors.New("schema operation result is invalid")
)

// SchemaOperation is the ledger entry for one preview confirmation. The token,
// scope and owner stay server-side.
type SchemaOperation struct {
	OperationID        string                `json:"operation_id"`
	Token              string                `json:"-"`
	WorkspaceID        string                `json:"-"`
	ScopeKey           string                `json:"-"`
	Owner              string                `json:"-"`
	Action             string                `json:"action"`
	Status             SchemaOperationStatus `json:"status"`
	ExecutedStatements int                   `json:"executed_statements"`
	VerifiedDigest     string                `json:"verified_digest,omitempty"`
	Reason             string                `json:"reason,omitempty"`
	NextAction         string                `json:"next_action,omitempty"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
	CompletedAt        *time.Time            `json:"completed_at,omitempty"`
}

// SchemaOperationResult is the terminal outcome the owning execution records.
type SchemaOperationResult struct {
	Status             SchemaOperationStatus
	ExecutedStatements int
	VerifiedDigest     string
	Reason             string
	NextAction         string
}

// ClaimSchemaApply obtains execution rights for a validated preview. The claim
// is keyed by the operation issued with the token, so a repeated or concurrent
// confirmation gets the existing operation instead of a second execution.
func (s *Service) ClaimSchemaApply(ctx context.Context, token *SchemaPreviewToken) (*SchemaOperation, ClaimOutcome, error) {
	if token == nil || isLegacyPreviewToken(token) {
		return nil, "", ErrPreviewTokenLegacy
	}
	owner, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create operation owner: %w", err)
	}
	now := time.Now().UTC()
	return s.repo.ClaimSchemaOperation(ctx, &SchemaOperation{
		OperationID: token.OperationID, Token: token.Token, WorkspaceID: token.WorkspaceID,
		ScopeKey: schemaOperationScopeKey(token), Owner: "claim-" + owner, Action: token.Action,
		Status: SchemaOperationRunning, CreatedAt: now, UpdatedAt: now,
	})
}

// ResolveSchemaApplyReplay reports, without executing or recording anything,
// whether an apply request repeats the operation issued with its token or meets
// another operation that holds the same scope.
func (s *Service) ResolveSchemaApplyReplay(ctx context.Context, workspaceID, tokenValue, operationID string) (*SchemaOperation, ClaimOutcome, error) {
	token, err := s.PreviewTokenForWorkspace(ctx, workspaceID, tokenValue)
	if err != nil {
		return nil, "", err
	}
	if isLegacyPreviewToken(token) {
		return nil, "", ErrPreviewTokenLegacy
	}
	if strings.TrimSpace(operationID) != token.OperationID {
		return nil, "", ErrSchemaOperationMismatch
	}
	op, err := s.GetSchemaOperation(ctx, workspaceID, token.OperationID)
	if err == nil {
		return op, claimOutcomeFor(op), nil
	}
	if !errors.Is(err, ErrSchemaOperationNotFound) {
		return nil, "", err
	}
	busy, err := s.repo.FindActiveSchemaOperation(ctx, workspaceID, schemaOperationScopeKey(token))
	switch {
	case err == nil:
		return withoutOwner(*busy), ClaimScopeBusy, nil
	case errors.Is(err, ErrSchemaOperationNotFound):
		return nil, ClaimNone, nil
	default:
		return nil, "", err
	}
}

// PreviewTokenForWorkspace returns a stored preview only to the workspace that
// issued it; a token of another workspace looks unknown.
func (s *Service) PreviewTokenForWorkspace(ctx context.Context, workspaceID, tokenValue string) (*SchemaPreviewToken, error) {
	token, err := s.repo.GetPreviewToken(ctx, strings.TrimSpace(tokenValue))
	if err != nil {
		return nil, err
	}
	if token.WorkspaceID != workspaceID {
		return nil, fmt.Errorf("%w: %s", ErrPreviewTokenNotFound, tokenValue)
	}
	return token, nil
}

// GetSchemaOperation returns an operation only within its owning workspace;
// the execution owner is never handed to readers.
func (s *Service) GetSchemaOperation(ctx context.Context, workspaceID, operationID string) (*SchemaOperation, error) {
	op, err := s.repo.GetSchemaOperation(ctx, workspaceID, strings.TrimSpace(operationID))
	if err != nil {
		return nil, err
	}
	return withoutOwner(*op), nil
}

// FinishSchemaOperation records the terminal result; only the owner that
// acquired the operation may do so, and only while it is still active.
func (s *Service) FinishSchemaOperation(ctx context.Context, operationID, owner string, result SchemaOperationResult) (*SchemaOperation, error) {
	if !isTerminalSchemaOperationStatus(result.Status) || result.ExecutedStatements < 0 {
		return nil, ErrSchemaOperationResult
	}
	return s.repo.FinishSchemaOperation(ctx, strings.TrimSpace(operationID), owner, result)
}

func isTerminalSchemaOperationStatus(status SchemaOperationStatus) bool {
	switch status {
	case SchemaOperationSucceeded, SchemaOperationPartial, SchemaOperationFailed, SchemaOperationUnknown:
		return true
	default:
		return false
	}
}

func claimOutcomeFor(op *SchemaOperation) ClaimOutcome {
	if op.Status.IsActive() {
		return ClaimInProgress
	}
	return ClaimCompleted
}

// withoutOwner returns a copy that cannot be used to finish the operation.
func withoutOwner(op SchemaOperation) *SchemaOperation {
	op.Owner = ""
	return &op
}

// SchemaOperationLease bounds how long an active operation may hold its scope
// without a terminal result before a retry may resolve it from target
// evidence. It matches the preview validity window by design; it is not a
// measured execution time.
const SchemaOperationLease = 10 * time.Minute

// schemaOperationScopeKey identifies the target tables two operations must not
// change at the same time.
func schemaOperationScopeKey(token *SchemaPreviewToken) string {
	return token.ScopeKey()
}

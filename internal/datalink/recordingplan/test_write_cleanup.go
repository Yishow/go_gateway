package recordingplan

import (
	"context"
	"crypto/subtle"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/common"
)

const (
	writeOutcomeVerified = "written_verified"
	cleanupStatusCleaned = "cleaned"
)

// ClaimTestWriteCleanup atomically acquires execution rights to retry only the
// cleanup of an operation that previously verified write and readback.
func (s *Service) ClaimTestWriteCleanup(ctx context.Context, workspaceID, tokenValue, operationID string, currentScope TestWriteScope, expectedContentDigest string) (*SchemaOperation, ClaimOutcome, error) {
	token, err := s.testWriteTokenForWorkspace(ctx, workspaceID, tokenValue)
	if err != nil {
		return nil, "", err
	}
	if strings.TrimSpace(operationID) != token.OperationID {
		return nil, "", ErrSchemaOperationMismatch
	}
	digest, err := previewDigest(token)
	if err != nil {
		return nil, "", err
	}
	if subtle.ConstantTimeCompare([]byte(digest), []byte(token.Digest)) != 1 {
		return nil, "", fmt.Errorf("%w: preview content changed", ErrPreviewTokenStale)
	}
	tokenContentDigest := TestWriteContentDigest(token)
	if tokenContentDigest == "" || tokenContentDigest != expectedContentDigest {
		return nil, "", fmt.Errorf("%w: test content changed", ErrPreviewTokenStale)
	}
	op, err := s.repo.GetSchemaOperation(ctx, workspaceID, token.OperationID)
	if err != nil {
		return nil, "", err
	}
	if op.Action != TestWriteAction {
		return nil, "", ErrPreviewTokenKind
	}
	if op.PayloadDigest != tokenContentDigest {
		return nil, "", fmt.Errorf("%w: operation payload digest does not match preview token", ErrSchemaOperationMismatch)
	}
	if op.WriteOutcome != writeOutcomeVerified {
		return nil, "", fmt.Errorf("%w: cannot repair cleanup when write is not verified", ErrSchemaOperationResult)
	}
	if op.CleanupStatus == cleanupStatusCleaned {
		return op, ClaimCompleted, nil
	}
	if op.Status.IsActive() {
		return op, ClaimInProgress, nil
	}
	savedScope := TestWriteScope{
		WorkspaceID:       token.WorkspaceID,
		WorkspaceRevision: token.WorkspaceRevision,
		GroupID:           token.PlanID,
		GroupRevision:     token.PlanRevision,
		ConnectorID:       token.ConnectorID,
		ConnectorRevision: token.ConnectorRevision,
		Dialect:           token.Dialect,
		Database:          token.Database,
		Schema:            token.Schema,
		Table:             token.TablePrefix,
	}.normalized()
	currentScope = currentScope.normalized()
	if savedScope != currentScope {
		return nil, "", fmt.Errorf("%w: test write scope changed", ErrPreviewTokenStale)
	}
	if strings.TrimSpace(op.Detail) == "" {
		return nil, "", fmt.Errorf("%w: missing original test write proof", ErrSchemaOperationResult)
	}
	reclaimer, ok := s.repo.(TestWriteCleanupRepository)
	if !ok {
		return nil, "", fmt.Errorf("%w: ledger does not support cleanup reclaim", ErrSchemaOperationResult)
	}
	ownerID, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create cleanup owner: %w", err)
	}
	return reclaimer.ReclaimTestWriteCleanupOperation(ctx, op.OperationID, workspaceID, token.Token, token.ScopeKey(), op.Detail, "claim-"+ownerID)
}

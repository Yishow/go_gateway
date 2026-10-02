package recordingplan

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
)

const (
	// TestWriteAction marks a token that may only confirm one explicit test write.
	TestWriteAction = "test_write"
	// TestWriteLease bounds how long an active test write may hold its table
	// without progress before a retry may adopt it and reconcile from target
	// evidence. It exceeds the sum of the write, readback and cleanup bounds.
	TestWriteLease = 3 * time.Minute

	testWriteContentPrefix = "content-digest:"
)

// ErrPreviewTokenKind means a token was issued for a different action.
var ErrPreviewTokenKind = errors.New("preview token is for a different action")

// TestWriteScope is the server-resolved scope a test-write preview is bound to.
type TestWriteScope struct {
	WorkspaceID       string
	WorkspaceRevision string
	GroupID           string
	GroupRevision     string
	ConnectorID       string
	ConnectorRevision string
	Dialect           string
	Database          string
	Schema            string
	Table             string
}

func (s TestWriteScope) normalized() TestWriteScope {
	s.Dialect = normalizePreviewDialect(s.Dialect)
	s.Table = strings.TrimSpace(s.Table)
	return s
}

// PrepareTestWritePreview persists a revision- and content-bound preview in
// the shared token ledger. It never touches the target. contentDigest binds
// the exact rows the confirmation will write.
func (s *Service) PrepareTestWritePreview(ctx context.Context, scope TestWriteScope, contentDigest string) (*SchemaPreviewToken, error) {
	scope = scope.normalized()
	if strings.TrimSpace(contentDigest) == "" || scope.Table == "" {
		return nil, fmt.Errorf("%w: test write preview needs a table and content", ErrPreviewTokenStale)
	}
	tokenID, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("create test write token: %w", err)
	}
	operationID, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("create test write operation: %w", err)
	}
	now := time.Now().UTC().Truncate(time.Second)
	token := &SchemaPreviewToken{
		Token: "tok-" + tokenID, OperationID: "op-" + operationID, Action: TestWriteAction,
		WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision,
		PlanID: scope.GroupID, PlanRevision: scope.GroupRevision,
		ConnectorID: scope.ConnectorID, ConnectorRevision: scope.ConnectorRevision,
		Dialect: scope.Dialect, Database: scope.Database, Schema: scope.Schema, TablePrefix: scope.Table,
		Statements: []string{testWriteContentPrefix + contentDigest},
		Tables:     []SchemaPreviewTable{{Name: scope.Table, Action: TestWriteAction}},
		ExpiresAt:  now.Add(schemaPreviewTTL), CreatedAt: now,
	}
	if token.Digest, err = previewDigest(token); err != nil {
		return nil, err
	}
	if err := s.repo.SavePreviewToken(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to save test write token: %w", err)
	}
	return token, nil
}

// testWriteTokenForWorkspace returns the stored token only to its workspace and
// only when it was issued for a test write.
func (s *Service) testWriteTokenForWorkspace(ctx context.Context, workspaceID, tokenValue string) (*SchemaPreviewToken, error) {
	token, err := s.PreviewTokenForWorkspace(ctx, workspaceID, tokenValue)
	if err != nil {
		return nil, err
	}
	if token.Action != TestWriteAction {
		return nil, ErrPreviewTokenKind
	}
	return token, nil
}

// TestWriteContentDigest returns the content digest a test-write token binds.
func TestWriteContentDigest(token *SchemaPreviewToken) string {
	if token == nil || len(token.Statements) != 1 {
		return ""
	}
	return strings.TrimPrefix(token.Statements[0], testWriteContentPrefix)
}

// ValidateTestWriteToken re-checks a stored preview against the current
// server-resolved scope and recomputed content before its first claim. An
// expired or stale token never writes; replays of an existing operation go
// through ResolveTestWriteReplay first and do not need a live token.
func (s *Service) ValidateTestWriteToken(ctx context.Context, workspaceID, tokenValue string, current TestWriteScope, currentContentDigest string) (*SchemaPreviewToken, error) {
	current = current.normalized()
	token, err := s.testWriteTokenForWorkspace(ctx, workspaceID, tokenValue)
	if err != nil {
		return nil, err
	}
	if isLegacyPreviewToken(&SchemaPreviewToken{
		OperationID: token.OperationID, WorkspaceRevision: token.WorkspaceRevision, PlanRevision: token.PlanRevision,
		ConnectorID: token.ConnectorID, ConnectorRevision: token.ConnectorRevision, Dialect: token.Dialect,
		TablePrefix: token.TablePrefix, Digest: token.Digest, Action: SchemaApplyAction,
	}) {
		return nil, ErrPreviewTokenLegacy
	}
	digest, err := previewDigest(token)
	if err != nil {
		return nil, err
	}
	if subtle.ConstantTimeCompare([]byte(digest), []byte(token.Digest)) != 1 {
		return nil, fmt.Errorf("%w: preview content changed", ErrPreviewTokenStale)
	}
	if token.IsExpired() {
		return nil, ErrPreviewTokenExpired
	}
	for _, check := range []struct{ field, stored, current string }{
		{"workspace revision", token.WorkspaceRevision, current.WorkspaceRevision},
		{"group", token.PlanID, current.GroupID},
		{"group revision", token.PlanRevision, current.GroupRevision},
		{"connector", token.ConnectorID, current.ConnectorID},
		{"connector revision", token.ConnectorRevision, current.ConnectorRevision},
		{"dialect", token.Dialect, current.Dialect},
		{"database", token.Database, current.Database},
		{"schema", token.Schema, current.Schema},
		{"table", token.TablePrefix, current.Table},
	} {
		if check.stored != check.current {
			return nil, fmt.Errorf("%w: %s changed", ErrPreviewTokenStale, check.field)
		}
	}
	if TestWriteContentDigest(token) != currentContentDigest {
		return nil, fmt.Errorf("%w: test content changed", ErrPreviewTokenStale)
	}
	return token, nil
}

// ResolveTestWriteReplay reports, without executing or recording anything,
// whether a confirmation repeats an operation already in the ledger.
func (s *Service) ResolveTestWriteReplay(ctx context.Context, workspaceID, tokenValue, operationID string) (*SchemaOperation, ClaimOutcome, error) {
	token, err := s.testWriteTokenForWorkspace(ctx, workspaceID, tokenValue)
	if err != nil {
		return nil, "", err
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
	return nil, ClaimNone, nil
}

// ClaimTestWrite obtains execution rights once for a validated preview. The
// claim is keyed by the operation issued with the token and shares the target
// scope with other test writes of the same table.
func (s *Service) ClaimTestWrite(ctx context.Context, token *SchemaPreviewToken) (*SchemaOperation, ClaimOutcome, error) {
	if token == nil || token.Action != TestWriteAction {
		return nil, "", ErrPreviewTokenKind
	}
	owner, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create operation owner: %w", err)
	}
	now := time.Now().UTC()
	return s.repo.ClaimSchemaOperation(ctx, &SchemaOperation{
		OperationID: token.OperationID, Token: token.Token, WorkspaceID: token.WorkspaceID,
		ScopeKey: token.ScopeKey(), Owner: "claim-" + owner, Action: TestWriteAction,
		Status: SchemaOperationRunning, PayloadDigest: TestWriteContentDigest(token), CreatedAt: now, UpdatedAt: now,
	})
}

// TakeOverTestWrite adopts a running test write whose lease ran out so a
// restarted gateway can reconcile it from target evidence. It returns nil when
// the lease is still live, the operation ended, or the ledger cannot adopt.
func (s *Service) TakeOverTestWrite(ctx context.Context, operationID string) (*SchemaOperation, error) {
	takeover, ok := s.repo.(operationExecutionRepository)
	if !ok {
		return nil, nil
	}
	owner, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("create operation owner: %w", err)
	}
	return takeover.TakeOverSchemaOperation(ctx, strings.TrimSpace(operationID), "claim-"+owner, time.Now().UTC().Add(-TestWriteLease))
}

// SaveTestWriteProgress records which step an owned test write reached, before
// that step's side effect, and renews the lease. A restarted gateway uses it to
// tell a write that never started from one that must not repeat.
func (s *Service) SaveTestWriteProgress(ctx context.Context, operationID, owner, detail string) error {
	progress, ok := s.repo.(operationExecutionRepository)
	if !ok {
		return fmt.Errorf("%w: ledger cannot record progress", ErrSchemaOperationResult)
	}
	return progress.SaveSchemaOperationProgress(ctx, strings.TrimSpace(operationID), owner, detail)
}

package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"
)

// SchemaExecution reports what a target adapter left applied.
type SchemaExecution struct {
	// Committed counts statements that remain applied on the target.
	Committed int
	// RolledBack is true when the adapter undid every statement of the batch.
	RolledBack bool
}

// SchemaExecutor runs confirmed statements on the saved target.
type SchemaExecutor func(ctx context.Context, statements []string) (SchemaExecution, error)

// SchemaApplyTarget is the saved target an apply inspects and changes.
type SchemaApplyTarget struct {
	Inspect TargetInspector
	Execute SchemaExecutor
}

var (
	ErrSchemaTargetChanged           = errors.New("target schema changed after preview")
	ErrSchemaAdapterUnverified       = errors.New("managed schema execution is not verified for this adapter")
	ErrSchemaOperationUnacknowledged = errors.New("schema operation result could not be recorded")

	errSchemaExecutorMissing = errors.New("no schema executor for the target")
)

// Reasons recorded on finished schema operations.
const (
	SchemaReasonRolledBack       = "statement_failed_rolled_back"
	SchemaReasonPartiallyApplied = "statement_failed_partially_applied"
	SchemaReasonVerificationGap  = "verification_mismatch"
	SchemaReasonUnverifiable     = "verification_unavailable"
	SchemaReasonPermissionDenied = "permission_denied"
	// SchemaReasonRecovered marks a terminal result resolved from the actual
	// target after the owning execution lost its acknowledgement.
	SchemaReasonRecovered = "recovered_from_target_evidence"

	schemaNextActionRepair = "preview again to create the missing tables"
	schemaNextActionCheck  = "check the target tables, then query the operation or preview again"
)

// ManagedSchemaExecutionVerified reports adapters whose managed DDL execution
// and post-execution verification have been tested.
func ManagedSchemaExecutionVerified(dialect string) bool {
	switch normalizePreviewDialect(dialect) {
	case dialectSQLite, dialectPostgres:
		return true
	default:
		return false
	}
}

// ApplySchemaPreview executes a validated preview at most once. A confirmation
// whose operation is already recorded gets that operation back. Otherwise the
// target must still match the preview before execution rights are claimed, and
// the recorded result reflects the schema actually found afterwards. An active
// operation whose claim lease has run out is resolved from the target instead
// of answering 202 forever.
func (s *Service) ApplySchemaPreview(ctx context.Context, token *SchemaPreviewToken, target SchemaApplyTarget) (*SchemaOperation, ClaimOutcome, error) {
	if token == nil || isLegacyPreviewToken(token) {
		return nil, "", ErrPreviewTokenLegacy
	}
	existing, err := s.GetSchemaOperation(ctx, token.WorkspaceID, token.OperationID)
	if err == nil {
		if existing.Status.IsActive() && schemaOperationLeaseExpired(existing) {
			return s.resolveStaleSchemaOperation(ctx, token, target, existing)
		}
		return existing, claimOutcomeFor(existing), nil
	}
	if !errors.Is(err, ErrSchemaOperationNotFound) {
		return nil, "", err
	}
	if !ManagedSchemaExecutionVerified(token.Dialect) {
		return nil, "", ErrSchemaAdapterUnverified
	}
	if token.IsExpired() {
		return nil, "", ErrPreviewTokenExpired
	}
	pendingBefore, err := confirmPreviewMatchesTarget(ctx, token, target.Inspect)
	if err != nil {
		if errors.Is(err, ErrPreviewTokenStale) {
			if existing, getErr := s.GetSchemaOperation(ctx, token.WorkspaceID, token.OperationID); getErr == nil {
				return existing, claimOutcomeFor(existing), nil
			}
		}
		return nil, "", err
	}

	op, outcome, err := s.ClaimSchemaApply(ctx, token)
	if err != nil || outcome != ClaimAcquired {
		return op, outcome, err
	}
	execution, execErr := runSchemaStatements(ctx, target.Execute, token.Statements)
	result := verifyAppliedSchema(ctx, token.TablePrefix, target.Inspect, pendingBefore, execution, execErr)
	finished, err := s.repo.FinishSchemaOperation(ctx, op.OperationID, op.Owner, result)
	if err != nil {
		return withoutOwner(*op), ClaimAcquired, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	return withoutOwner(*finished), ClaimAcquired, nil
}

// confirmPreviewMatchesTarget recomputes the preview from the target as it is
// now and returns how many tables are pending; any difference makes it stale.
func confirmPreviewMatchesTarget(ctx context.Context, token *SchemaPreviewToken, inspect TargetInspector) (int, error) {
	if inspect == nil {
		return 0, fmt.Errorf("%w: no target inspector", ErrTargetInspectionUnconfirmed)
	}
	statements, err := GenerateManagedSchemaDDL(token.Dialect, token.TablePrefix)
	if err != nil {
		return 0, err
	}
	pending, _, err := inspectManagedTables(ctx, token.TablePrefix, inspect)
	if err != nil {
		return 0, err
	}
	if !slices.Equal(statementsForTables(statements, token.TablePrefix, pending), token.Statements) {
		return 0, fmt.Errorf("%w: %w", ErrPreviewTokenStale, ErrSchemaTargetChanged)
	}
	return len(pending), nil
}

func runSchemaStatements(ctx context.Context, execute SchemaExecutor, statements []string) (SchemaExecution, error) {
	if len(statements) == 0 {
		return SchemaExecution{}, nil
	}
	if execute == nil {
		return SchemaExecution{RolledBack: true}, errSchemaExecutorMissing
	}
	return execute(ctx, statements)
}

func schemaOperationLeaseExpired(op *SchemaOperation) bool {
	return time.Since(op.UpdatedAt) > SchemaOperationLease
}

// resolveStaleSchemaOperation closes an active operation whose claim lease has
// run out. It never re-executes anything: the recorded result is derived from
// the schema actually found on the target, which is the evidence the ledger
// keeps the scope for.
func (s *Service) resolveStaleSchemaOperation(ctx context.Context, token *SchemaPreviewToken, target SchemaApplyTarget, stale *SchemaOperation) (*SchemaOperation, ClaimOutcome, error) {
	result := recoveryResultFromTarget(ctx, token.TablePrefix, target.Inspect)
	cutoff := time.Now().UTC().Add(-SchemaOperationLease)
	finished, err := s.repo.FinishStaleSchemaOperation(ctx, stale.OperationID, result, cutoff)
	if err != nil {
		return withoutOwner(*stale), "", fmt.Errorf("resolve stale schema operation %s: %w", stale.OperationID, err)
	}
	if finished == nil {
		// The owning execution finished the operation first; its result wins.
		current, err := s.GetSchemaOperation(ctx, token.WorkspaceID, token.OperationID)
		if err != nil {
			return nil, "", err
		}
		return current, claimOutcomeFor(current), nil
	}
	return finished, ClaimCompleted, nil
}

// recoveryResultFromTarget decides the honest outcome for a lost execution:
// a target that carries every managed table proves success, everything else
// stays unknown because what happened cannot be proven.
func recoveryResultFromTarget(ctx context.Context, prefix string, inspect TargetInspector) SchemaOperationResult {
	pending, tables, err := inspectManagedTables(ctx, prefix, inspect)
	switch {
	case err != nil:
		return SchemaOperationResult{Status: SchemaOperationUnknown, Reason: SchemaReasonUnverifiable, NextAction: schemaNextActionCheck}
	case len(pending) == 0:
		return SchemaOperationResult{Status: SchemaOperationSucceeded, VerifiedDigest: schemaTablesDigest(tables), Reason: SchemaReasonRecovered}
	default:
		return SchemaOperationResult{Status: SchemaOperationUnknown, Reason: SchemaReasonUnverifiable, NextAction: schemaNextActionRepair}
	}
}

// verifyAppliedSchema decides the outcome from the tables actually present
// after execution; an adapter's own report never counts as success on its own.
func verifyAppliedSchema(ctx context.Context, prefix string, inspect TargetInspector, pendingBefore int, execution SchemaExecution, execErr error) SchemaOperationResult {
	result := SchemaOperationResult{ExecutedStatements: execution.Committed}
	pending, tables, err := inspectManagedTables(ctx, prefix, inspect)
	switch {
	case err != nil:
		result.Status, result.Reason, result.NextAction = SchemaOperationUnknown, SchemaReasonUnverifiable, schemaNextActionCheck
	case len(pending) == 0:
		result.Status, result.VerifiedDigest = SchemaOperationSucceeded, schemaTablesDigest(tables)
	case execErr != nil && errors.Is(execErr, ErrTargetPermissionDenied) && len(pending) == pendingBefore:
		result.Status, result.Reason = SchemaOperationFailed, SchemaReasonPermissionDenied
		result.NextAction = "grant permission to create tables in the target schema, then preview again"
	case execErr != nil && execution.RolledBack && len(pending) == pendingBefore:
		result.Status, result.Reason, result.NextAction = SchemaOperationFailed, SchemaReasonRolledBack, "fix the reported cause, then preview again"
	case execErr != nil:
		result.Status, result.Reason, result.NextAction = SchemaOperationPartial, SchemaReasonPartiallyApplied, schemaNextActionRepair
	default:
		result.Status, result.Reason, result.NextAction = SchemaOperationPartial, SchemaReasonVerificationGap, schemaNextActionRepair
	}
	return result
}

func schemaTablesDigest(tables []SchemaPreviewTable) string {
	payload, err := json.Marshal(tables)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

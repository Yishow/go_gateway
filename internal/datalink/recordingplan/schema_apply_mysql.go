package recordingplan

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
)

const (
	mysqlSchemaNextActionRetry      = "grant target permission if needed, then retry the same confirmed preview"
	mysqlSchemaNextActionNewPreview = "grant target permission, then prepare a new schema preview"
	mysqlSchemaNextActionInspect    = "inspect the target catalog; do not rebuild or replace an incompatible object"
)

// applyMySQLSchemaPreview runs only the canonical MySQL group layout. The
// generic MySQL managed-schema gate remains closed for every other path.
func (s *Service) applyMySQLSchemaPreview(ctx context.Context, supplied *SchemaPreviewToken, target SchemaApplyTarget) (*SchemaOperation, ClaimOutcome, error) {
	if isLegacyPreviewToken(supplied) {
		return nil, "", ErrPreviewTokenLegacy
	}
	token, err := s.persistedGroupPreviewToken(ctx, supplied)
	if err != nil {
		return nil, "", err
	}
	if token.GroupLayout == nil || token.GroupLayout.MySQL == nil || token.Dialect != dialectMySQL ||
		token.Action != SchemaApplyAction || !mysqlSchemaApplyBaselineSupported(token.GroupLayout) {
		return nil, "", ErrSchemaAdapterUnverified
	}
	if err := validateMySQLPreviewScope(token.Database, token.Schema, token.Dialect); err != nil {
		return nil, "", err
	}
	if target.Inspect == nil {
		return nil, "", fmt.Errorf("%w: no target inspector", ErrTargetInspectionUnconfirmed)
	}
	if len(token.Statements) > 0 && target.Execute == nil {
		return nil, "", errSchemaExecutorMissing
	}
	manifest, err := newMySQLSchemaApplyManifest(token)
	if err != nil {
		return nil, "", err
	}

	existing, err := s.repo.GetSchemaOperation(ctx, token.WorkspaceID, token.OperationID)
	if err == nil {
		if existing.WorkspaceID != token.WorkspaceID || existing.Token != token.Token || existing.Action != token.Action ||
			existing.ScopeKey != schemaOperationScopeKey(token) {
			return nil, "", ErrSchemaOperationMismatch
		}
		switch existing.Status {
		case SchemaOperationPending, SchemaOperationRunning:
			if !schemaOperationLeaseExpired(existing) {
				return withoutOwner(*existing), ClaimInProgress, nil
			}
			return s.takeOverMySQLSchemaApply(ctx, token, target, existing)
		case SchemaOperationPartial, SchemaOperationUnknown:
			return s.reclaimMySQLSchemaApply(ctx, token, target, existing)
		default:
			return withoutOwner(*existing), ClaimCompleted, nil
		}
	}
	if !errors.Is(err, ErrSchemaOperationNotFound) {
		return nil, "", err
	}
	if token.IsExpired() {
		return nil, "", ErrPreviewTokenExpired
	}

	if _, err := confirmPreviewMatchesTarget(ctx, token, target.Inspect); err != nil {
		return nil, "", err
	}
	if err := verifyMySQLApplyStepCatalog(ctx, token, target.Inspect, manifest); err != nil {
		return nil, "", err
	}
	detail, err := encodeMySQLSchemaApplyManifest(manifest)
	if err != nil {
		return nil, "", err
	}
	op, outcome, err := s.claimMySQLSchemaApply(ctx, token, detail)
	if err != nil || outcome != ClaimAcquired {
		return op, outcome, err
	}
	return s.runMySQLSchemaApply(ctx, token, target, op, manifest)
}

func (s *Service) claimMySQLSchemaApply(ctx context.Context, token *SchemaPreviewToken, detail string) (*SchemaOperation, ClaimOutcome, error) {
	owner, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create MySQL schema operation owner: %w", err)
	}
	now := time.Now().UTC()
	return s.repo.ClaimSchemaOperation(ctx, &SchemaOperation{
		OperationID: token.OperationID, Token: token.Token, WorkspaceID: token.WorkspaceID,
		ScopeKey: schemaOperationScopeKey(token), Owner: "claim-" + owner, Action: token.Action,
		Status: SchemaOperationRunning, Detail: detail, CreatedAt: now, UpdatedAt: now,
	})
}

func (s *Service) takeOverMySQLSchemaApply(ctx context.Context, token *SchemaPreviewToken, target SchemaApplyTarget, stale *SchemaOperation) (*SchemaOperation, ClaimOutcome, error) {
	executionRepo, ok := s.repo.(operationExecutionRepository)
	if !ok {
		return withoutOwner(*stale), ClaimInProgress, nil
	}
	ownerID, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create MySQL schema operation owner: %w", err)
	}
	owner := "claim-" + ownerID
	taken, err := executionRepo.TakeOverSchemaOperation(ctx, stale.OperationID, owner, time.Now().UTC().Add(-SchemaOperationLease))
	if err != nil {
		return withoutOwner(*stale), "", fmt.Errorf("take over MySQL schema operation: %w", err)
	}
	if taken == nil {
		current, getErr := s.repo.GetSchemaOperation(ctx, token.WorkspaceID, token.OperationID)
		if getErr != nil {
			return nil, "", getErr
		}
		return withoutOwner(*current), claimOutcomeFor(current), nil
	}
	if taken.Owner != owner || taken.Token != token.Token || taken.Action != token.Action ||
		taken.ScopeKey != schemaOperationScopeKey(token) {
		return withoutOwner(*taken), ClaimAcquired, ErrSchemaOperationMismatch
	}
	manifest, err := mysqlSchemaApplyManifestFromDetail(taken.Detail, token)
	if err != nil {
		return s.finishMySQLSchemaApply(ctx, taken, mysqlSchemaApplyManifest{}, SchemaOperationResult{
			Status: SchemaOperationUnknown, Reason: SchemaReasonUnverifiable, NextAction: mysqlSchemaNextActionInspect,
			Detail: taken.Detail,
		})
	}
	if err := verifyMySQLApplyStepCatalog(ctx, token, target.Inspect, manifest); err != nil {
		status, reason := mysqlSchemaPreflightFailure(err)
		return s.finishMySQLSchemaApply(ctx, taken, manifest, SchemaOperationResult{
			Status: status, ExecutedStatements: verifiedMySQLSchemaStepCount(manifest),
			Reason: reason, NextAction: mysqlSchemaNextActionInspect,
		})
	}
	return s.runMySQLSchemaApply(ctx, token, target, taken, manifest)
}

func (s *Service) reclaimMySQLSchemaApply(ctx context.Context, token *SchemaPreviewToken, target SchemaApplyTarget, terminal *SchemaOperation) (*SchemaOperation, ClaimOutcome, error) {
	manifest, err := mysqlSchemaApplyManifestFromDetail(terminal.Detail, token)
	if err != nil {
		return withoutOwner(*terminal), ClaimCompleted, err
	}
	if err := verifyMySQLApplyStepCatalog(ctx, token, target.Inspect, manifest); err != nil {
		return withoutOwner(*terminal), ClaimCompleted, err
	}
	resumeRepo, ok := s.repo.(mysqlSchemaOperationResumeRepository)
	if !ok {
		return withoutOwner(*terminal), ClaimCompleted, ErrSchemaAdapterUnverified
	}
	ownerID, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create MySQL schema operation owner: %w", err)
	}
	owner := "claim-" + ownerID
	reclaimed, outcome, err := resumeRepo.ReclaimMySQLSchemaApplyOperation(
		ctx, terminal.OperationID, token.WorkspaceID, token.Token, schemaOperationScopeKey(token), terminal.Detail, owner,
	)
	if err != nil || outcome != ClaimAcquired {
		return reclaimed, outcome, err
	}
	if reclaimed.Owner != owner || reclaimed.Status != SchemaOperationRunning {
		return withoutOwner(*reclaimed), ClaimAcquired, ErrSchemaOperationNotOwned
	}
	manifest, err = mysqlSchemaApplyManifestFromDetail(reclaimed.Detail, token)
	if err != nil {
		return s.finishMySQLSchemaApply(ctx, reclaimed, mysqlSchemaApplyManifest{}, SchemaOperationResult{
			Status: SchemaOperationUnknown, Reason: SchemaReasonUnverifiable, NextAction: mysqlSchemaNextActionInspect,
			Detail: reclaimed.Detail,
		})
	}
	return s.runMySQLSchemaApply(ctx, token, target, reclaimed, manifest)
}

func verifyMySQLApplyStepCatalog(ctx context.Context, token *SchemaPreviewToken, inspect TargetInspector, manifest mysqlSchemaApplyManifest) error {
	found, err := inspectMySQLGroupTables(ctx, *token.GroupLayout, token.Database, inspect)
	if err != nil {
		return err
	}
	pendingStatements, err := groupSchemaStatements(token.Dialect, *token.GroupLayout, found.tables)
	if err != nil {
		return err
	}
	planned := make(map[string]struct{}, len(token.Statements))
	for _, statement := range token.Statements {
		planned[statement] = struct{}{}
	}
	for _, statement := range pendingStatements {
		if _, ok := planned[statement]; !ok {
			return fmt.Errorf("%w: MySQL target needs DDL outside the confirmed preview", ErrSchemaTargetChanged)
		}
	}
	for _, step := range manifest.Steps {
		evidence, err := inspectMySQLApplyStep(ctx, token, *token.GroupLayout, mysqlApplyStepRef{
			Kind: step.Kind, Name: step.Name, SQLHash: step.SQLHash,
		}, inspect)
		if err != nil {
			return err
		}
		if evidence.State == mysqlStepStateMissing {
			if step.Status == mysqlSchemaStepVerified {
				return fmt.Errorf("%w: a verified MySQL DDL object is missing", ErrSchemaTargetChanged)
			}
			continue
		}
		if evidence.State != mysqlStepStateExact {
			return ErrTargetInspectionUnconfirmed
		}
		if step.Status == mysqlSchemaStepVerified && step.Fingerprint != evidence.Fingerprint {
			return fmt.Errorf("%w: MySQL DDL fingerprint changed", ErrSchemaTargetChanged)
		}
	}
	return nil
}

func (s *Service) runMySQLSchemaApply(ctx context.Context, token *SchemaPreviewToken, target SchemaApplyTarget, op *SchemaOperation, manifest mysqlSchemaApplyManifest) (*SchemaOperation, ClaimOutcome, error) {
	verified := 0
	for index := range manifest.Steps {
		step := &manifest.Steps[index]
		evidence, err := inspectMySQLApplyStep(ctx, token, *token.GroupLayout, mysqlApplyStepRef{
			Kind: step.Kind, Name: step.Name, SQLHash: step.SQLHash,
		}, target.Inspect)
		if err != nil {
			return s.finishMySQLCatalogFailure(ctx, op, manifest, step, err, verified)
		}
		if evidence.State == mysqlStepStateExact {
			if step.Status == mysqlSchemaStepVerified && step.Fingerprint != evidence.Fingerprint {
				return s.finishMySQLCatalogFailure(ctx, op, manifest, step, ErrSchemaTargetChanged, verified)
			}
			step.Status, step.Fingerprint, step.FailureCode = mysqlSchemaStepVerified, evidence.Fingerprint, ""
			if err := s.saveMySQLSchemaApplyProgress(ctx, op, manifest); err != nil {
				return withoutOwner(*op), ClaimAcquired, err
			}
			verified++
			continue
		}
		if evidence.State != mysqlStepStateMissing {
			return s.finishMySQLCatalogFailure(ctx, op, manifest, step, ErrTargetInspectionUnconfirmed, verified)
		}
		if step.Status == mysqlSchemaStepVerified {
			return s.finishMySQLCatalogFailure(ctx, op, manifest, step, ErrSchemaTargetChanged, verified)
		}

		step.Status, step.Fingerprint, step.FailureCode = mysqlSchemaStepStarted, "", ""
		if err := s.saveMySQLSchemaApplyProgress(ctx, op, manifest); err != nil {
			return withoutOwner(*op), ClaimAcquired, err
		}
		if err := s.requireMySQLSchemaOperationOwner(ctx, op); err != nil {
			return withoutOwner(*op), ClaimAcquired, err
		}
		statement := token.Statements[index]
		_, execErr := runSchemaStatements(ctx, target.Execute, []string{statement})

		after, inspectErr := inspectMySQLApplyStep(ctx, token, *token.GroupLayout, mysqlApplyStepRef{
			Kind: step.Kind, Name: step.Name, SQLHash: step.SQLHash,
		}, target.Inspect)
		if inspectErr != nil {
			return s.finishMySQLSchemaApply(ctx, op, manifest, SchemaOperationResult{
				Status: SchemaOperationUnknown, ExecutedStatements: verified,
				Reason: SchemaReasonUnverifiable, NextAction: mysqlSchemaNextActionInspect,
			})
		}
		if after.State == mysqlStepStateExact {
			step.Status, step.Fingerprint, step.FailureCode = mysqlSchemaStepVerified, after.Fingerprint, ""
			if err := s.saveMySQLSchemaApplyProgress(ctx, op, manifest); err != nil {
				return withoutOwner(*op), ClaimAcquired, err
			}
			verified++
			continue
		}
		if after.State != mysqlStepStateMissing {
			return s.finishMySQLCatalogFailure(ctx, op, manifest, step, ErrTargetInspectionUnconfirmed, verified)
		}

		step.Status = mysqlSchemaStepFailed
		step.FailureCode = mysqlSchemaDDLFailureCode(execErr)
		if err := s.saveMySQLSchemaApplyProgress(ctx, op, manifest); err != nil {
			return withoutOwner(*op), ClaimAcquired, err
		}
		status := SchemaOperationFailed
		if verified > 0 {
			status = SchemaOperationPartial
		}
		reason, next := SchemaReasonPartiallyApplied, mysqlSchemaNextActionInspect
		if errors.Is(execErr, ErrTargetPermissionDenied) {
			reason, next = SchemaReasonPermissionDenied, mysqlSchemaNextActionNewPreview
			if status == SchemaOperationPartial {
				next = mysqlSchemaNextActionRetry
			}
		}
		if execErr == nil {
			reason = SchemaReasonVerificationGap
		}
		return s.finishMySQLSchemaApply(ctx, op, manifest, SchemaOperationResult{
			Status: status, ExecutedStatements: verified, Reason: reason, NextAction: next,
		})
	}

	found, err := inspectMySQLGroupTables(ctx, *token.GroupLayout, token.Database, target.Inspect)
	if err != nil {
		return s.finishMySQLSchemaApply(ctx, op, manifest, SchemaOperationResult{
			Status: SchemaOperationUnknown, ExecutedStatements: verified,
			Reason: SchemaReasonUnverifiable, NextAction: mysqlSchemaNextActionInspect,
		})
	}
	if found.pendingGroup || found.pendingReceipt {
		return s.finishMySQLSchemaApply(ctx, op, manifest, SchemaOperationResult{
			Status: SchemaOperationPartial, ExecutedStatements: verified,
			Reason: SchemaReasonVerificationGap, NextAction: mysqlSchemaNextActionInspect,
		})
	}
	return s.finishMySQLSchemaApply(ctx, op, manifest, SchemaOperationResult{
		Status: SchemaOperationSucceeded, ExecutedStatements: verified,
		VerifiedDigest: schemaTablesDigest(found.tables),
	})
}

func (s *Service) finishMySQLCatalogFailure(ctx context.Context, op *SchemaOperation, manifest mysqlSchemaApplyManifest, step *mysqlSchemaApplyStep, cause error, verified int) (*SchemaOperation, ClaimOutcome, error) {
	if errors.Is(cause, ErrTargetInspectionUnconfirmed) || errors.Is(cause, ErrTargetPermissionDenied) {
		return s.finishMySQLSchemaApply(ctx, op, manifest, SchemaOperationResult{
			Status: SchemaOperationUnknown, ExecutedStatements: verified,
			Reason: SchemaReasonUnverifiable, NextAction: mysqlSchemaNextActionInspect,
		})
	}
	if step != nil {
		step.Status, step.Fingerprint, step.FailureCode = mysqlSchemaStepFailed, "", "catalog_mismatch"
	}
	return s.finishMySQLSchemaApply(ctx, op, manifest, SchemaOperationResult{
		Status: SchemaOperationFailed, ExecutedStatements: verified,
		Reason: SchemaReasonVerificationGap, NextAction: mysqlSchemaNextActionInspect,
	})
}

func mysqlSchemaPreflightFailure(err error) (status SchemaOperationStatus, reason string) {
	if errors.Is(err, ErrTargetInspectionUnconfirmed) || errors.Is(err, ErrTargetPermissionDenied) {
		return SchemaOperationUnknown, SchemaReasonUnverifiable
	}
	return SchemaOperationFailed, SchemaReasonVerificationGap
}

func verifiedMySQLSchemaStepCount(manifest mysqlSchemaApplyManifest) int {
	verified := 0
	for _, step := range manifest.Steps {
		if step.Status == mysqlSchemaStepVerified {
			verified++
		}
	}
	return verified
}

func mysqlSchemaDDLFailureCode(err error) string {
	if errors.Is(err, ErrTargetPermissionDenied) {
		return "permission_denied"
	}
	return "statement_failed"
}

func (s *Service) saveMySQLSchemaApplyProgress(ctx context.Context, op *SchemaOperation, manifest mysqlSchemaApplyManifest) error {
	progressRepo, ok := s.repo.(operationExecutionRepository)
	if !ok {
		return fmt.Errorf("%w: MySQL schema progress is unsupported", ErrSchemaOperationUnacknowledged)
	}
	detail, err := encodeMySQLSchemaApplyManifest(manifest)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	if err := progressRepo.SaveSchemaOperationProgress(ctx, op.OperationID, op.Owner, detail); err != nil {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	op.Detail = detail
	return nil
}

func (s *Service) requireMySQLSchemaOperationOwner(ctx context.Context, op *SchemaOperation) error {
	current, err := s.repo.GetSchemaOperation(ctx, op.WorkspaceID, op.OperationID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	if current.Owner != op.Owner || !current.Status.IsActive() {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, ErrSchemaOperationNotOwned)
	}
	return nil
}

func (s *Service) finishMySQLSchemaApply(ctx context.Context, op *SchemaOperation, manifest mysqlSchemaApplyManifest, result SchemaOperationResult) (*SchemaOperation, ClaimOutcome, error) {
	if result.Detail == "" {
		detail, err := encodeMySQLSchemaApplyManifest(manifest)
		if err != nil {
			return withoutOwner(*op), ClaimAcquired, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
		}
		result.Detail = detail
	}
	finished, err := s.repo.FinishSchemaOperation(ctx, op.OperationID, op.Owner, result)
	if err != nil {
		return withoutOwner(*op), ClaimAcquired, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	return withoutOwner(*finished), ClaimAcquired, nil
}

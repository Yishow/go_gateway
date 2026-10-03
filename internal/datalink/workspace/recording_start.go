package workspace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"go-gateway/internal/datalink/recordingplan"
)

// RecordingStartService coordinates only saved setup, group Apply and scoped
// activation. It shares the schema/test-write operation ledger and performs no DDL.
type RecordingStartService struct {
	workspace     *Service
	groups        *WriteGroupService
	activate      recordingStartActivator
	ledger        *recordingplan.Service
	barrier       func(context.Context, RecordingStartRequest) error
	resumeBarrier func(context.Context, RecordingStartRequest) error
}

// NewRecordingStartService uses the already wired configuration services.
func NewRecordingStartService(ws *Service, groups *WriteGroupService, activation recordingStartActivator, ledger *recordingplan.Service) *RecordingStartService {
	return &RecordingStartService{workspace: ws, groups: groups, activate: activation, ledger: ledger}
}

// WithActivationBarrier retains the existing Share hydration/revision guards.
func (s *RecordingStartService) WithActivationBarrier(barrier, resumeBarrier func(context.Context, RecordingStartRequest) error) *RecordingStartService {
	s.barrier = barrier
	s.resumeBarrier = resumeBarrier
	return s
}

// Get reads recorded progress without saving, applying, activating or inspecting SQL.
func (s *RecordingStartService) Get(ctx context.Context, id string) (*RecordingStartOperation, error) {
	if s == nil || s.workspace == nil || s.ledger == nil {
		return nil, ErrWriteGroupServiceUnavailable
	}
	record, err := s.workspace.repo.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, recordingplan.ErrSchemaOperationNotFound
	}
	if err != nil {
		return nil, err
	}
	op, err := s.ledger.GetSchemaOperation(ctx, record.ID, id)
	if err != nil {
		return nil, err
	}
	p, err := decodeRecordingStart(op)
	if err != nil {
		return nil, err
	}
	return recordingStartView(op, p), nil
}

// Start returns or resumes the same request's recorded scope. Local save/Apply
// results and their ledger checkpoints commit in one configuration transaction.
func (s *RecordingStartService) Start(ctx context.Context, request RecordingStartRequest) (*RecordingStartOperation, error) {
	if s == nil || s.workspace == nil || s.groups == nil || s.ledger == nil {
		return nil, ErrWriteGroupServiceUnavailable
	}
	request, digest, err := normalizeRecordingStart(request)
	if err != nil {
		return nil, err
	}
	record, err := s.workspace.GetOrCreate(ctx)
	if err != nil {
		return nil, err
	}
	if record.ID != request.WorkspaceID {
		return nil, ErrRecordingStartInvalid
	}
	previous, err := s.ledger.FindRecordingStart(ctx, request.WorkspaceID, request.RequestID, digest)
	if err != nil && !errors.Is(err, recordingplan.ErrSchemaOperationNotFound) {
		return nil, err
	}
	if previous != nil && previous.Status == recordingplan.SchemaOperationSucceeded {
		p, err := decodeRecordingStart(previous)
		if err != nil {
			return nil, err
		}
		return recordingStartView(previous, p), nil
	}
	if err := s.validateRecordingStartSelection(ctx, request); err != nil {
		return nil, err
	}
	p := initialRecordingStart(request, digest)
	var initialValidationErr error
	if previous == nil {
		// Capture source/schema identity in the claim itself, before the first
		// effect or checkpoint can be interrupted. A failed capture has no
		// effects and remains explicitly recoverable as an initial intent.
		initialValidationErr = s.validateRecordingStart(ctx, p, true)
		if initialValidationErr != nil {
			p = initialRecordingStart(request, digest)
		}
	}
	encoded, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("encode recording start: %w", err)
	}
	op, outcome, err := s.ledger.ClaimRecordingStart(ctx, request.WorkspaceID, request.RequestID, digest, string(encoded))
	if err != nil {
		return nil, err
	}
	if outcome == recordingplan.ClaimScopeBusy {
		return nil, ErrRecordingStartBusy
	}
	if outcome != recordingplan.ClaimAcquired {
		initialValidationErr = nil // A concurrent claimant owns the persisted intent.
		p, err = decodeRecordingStart(op)
		if err != nil {
			return nil, err
		}
		if op.Status == recordingplan.SchemaOperationSucceeded {
			return recordingStartView(op, p), nil
		}
		if err := s.validateRecordingStart(ctx, p, recordingStartNeedsInitialization(p)); err != nil {
			view := recordingStartView(op, p)
			view.Reason, view.NextAction = "stale_intent", "revalidate"
			return view, nil //nolint:nilerr // A stale replay returns persisted progress and an explicit revalidation action.
		}
		op, outcome, err = s.ledger.ResumeRecordingStart(ctx, request.WorkspaceID, op.OperationID, digest)
		if err != nil {
			return nil, err
		}
		if outcome != recordingplan.ClaimAcquired {
			return recordingStartView(op, p), nil
		}
		// Takeover returns the latest fenced checkpoint, which can be newer
		// than the progress read before the lease acquisition.
		p, err = decodeRecordingStart(op)
		if err != nil {
			return nil, err
		}
	}
	if initialValidationErr != nil {
		return s.finishRecordingStart(ctx, op, p, initialValidationErr)
	}
	if err := s.validateRecordingStart(ctx, p, true); err != nil {
		return s.finishRecordingStart(ctx, op, p, err)
	}
	if err := s.runRecordingStart(ctx, op, p); err != nil {
		return s.finishRecordingStart(ctx, op, p, err)
	}
	p.View.Stage = "complete"
	return s.finishRecordingStart(ctx, op, p, nil)
}

func decodeRecordingStart(op *recordingplan.SchemaOperation) (*recordingStartProgress, error) {
	if op == nil || op.Action != recordingplan.RecordingStartAction {
		return nil, recordingplan.ErrSchemaOperationNotFound
	}
	var p recordingStartProgress
	if err := json.Unmarshal([]byte(op.Detail), &p); err != nil || p.Version != 1 ||
		p.Request.WorkspaceID != op.WorkspaceID || p.View.IntentDigest != op.PayloadDigest ||
		len(p.Scopes) != len(p.Request.Groups) || len(p.View.Groups) != len(p.Request.Groups) ||
		len(p.DeviceDigests) != len(p.Request.DeviceIDs) ||
		len(p.View.Devices) != len(p.Request.DeviceIDs) {
		return nil, recordingplan.ErrSchemaOperationResult
	}
	_, digest, err := normalizeRecordingStart(p.Request)
	if err != nil || digest != op.PayloadDigest || p.View.WorkspaceID != p.Request.WorkspaceID ||
		!slices.Equal(p.View.DeviceIDs, p.Request.DeviceIDs) {
		return nil, recordingplan.ErrSchemaOperationResult
	}
	for i, group := range p.Request.Groups {
		if p.View.Groups[i].GroupID != group.GroupID || p.View.Groups[i].GroupRevision == "" {
			return nil, recordingplan.ErrSchemaOperationResult
		}
	}
	for i, deviceID := range p.Request.DeviceIDs {
		if p.View.Devices[i].DeviceID != deviceID {
			return nil, recordingplan.ErrSchemaOperationResult
		}
	}
	return &p, nil
}

func recordingStartView(op *recordingplan.SchemaOperation, p *recordingStartProgress) *RecordingStartOperation {
	view := p.View
	view.OperationID, view.Action, view.Status = op.OperationID, op.Action, op.Status
	view.SetupRevision = p.WorkspaceRevision
	view.Reason, view.NextAction = op.Reason, op.NextAction
	view.CreatedAt, view.UpdatedAt = op.CreatedAt, op.UpdatedAt
	return &view
}

func (s *RecordingStartService) saveRecordingStartProgress(ctx context.Context, op *recordingplan.SchemaOperation, p *recordingStartProgress) error {
	encoded, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode recording start progress: %w", err)
	}
	return s.ledger.SaveRecordingStartProgress(ctx, op.OperationID, op.Owner, string(encoded))
}

func (s *RecordingStartService) finishRecordingStart(ctx context.Context, op *recordingplan.SchemaOperation, p *recordingStartProgress, cause error) (*RecordingStartOperation, error) {
	settle, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	if cause != nil {
		// A failed local transaction may have called its checkpoint before
		// rollback. Only the ledger's committed progress is a result fact.
		persisted, err := s.ledger.GetSchemaOperation(settle, p.Request.WorkspaceID, op.OperationID)
		if err != nil {
			return nil, err
		}
		p, err = decodeRecordingStart(persisted)
		if err != nil {
			return nil, err
		}
	}
	status, reason, action := recordingplan.SchemaOperationSucceeded, "", ""
	if cause != nil {
		status, reason, action = recordingplan.SchemaOperationFailed, "start_failed", "retry_same_request"
		for _, group := range p.View.Groups {
			if group.Applied {
				status = recordingplan.SchemaOperationPartial
			}
		}
		for _, device := range p.View.Devices {
			if device.Activated {
				status = recordingplan.SchemaOperationPartial
			}
		}
		switch {
		case errors.Is(cause, ErrWriteGroupApplyNotReady):
			reason, action = "preparation_required", "prepare_schema"
		case errors.Is(cause, ErrSetupRevisionConflict), errors.Is(cause, ErrWriteGroupRevisionConflict):
			reason, action = "stale_intent", "revalidate"
		case errors.Is(cause, ErrRecordingStartInvalid), errors.Is(cause, ErrReadinessScopeInvalid):
			reason, action = "invalid_scope", "review_selection"
		case errors.Is(cause, ErrRecordingStartShareNotReady):
			reason, action = "share_not_ready", "refresh_share"
		case errors.Is(cause, ErrReadinessUnavailable), errors.Is(cause, ErrWriteGroupServiceUnavailable):
			reason, action = "service_unavailable", "retry_same_request"
		case errors.Is(cause, recordingplan.ErrSchemaOperationNotOwned):
			return s.Get(settle, op.OperationID)
		}
		var blocked *ReadinessBlockedError
		if errors.As(cause, &blocked) {
			reason, action = "scope_not_ready", "review_device"
		}
	}
	encoded, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("encode recording start result: %w", err)
	}
	completed, err := s.ledger.FinishSchemaOperation(settle, op.OperationID, op.Owner, recordingplan.SchemaOperationResult{
		Status: status, Reason: reason, NextAction: action, PayloadDigest: p.View.IntentDigest, Detail: string(encoded),
	})
	if err != nil {
		return nil, err
	}
	return recordingStartView(completed, p), nil
}

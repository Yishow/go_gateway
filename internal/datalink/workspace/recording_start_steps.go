package workspace

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"

	"go-gateway/internal/datalink/recordingplan"
)

func (s *RecordingStartService) runRecordingStart(ctx context.Context, op *recordingplan.SchemaOperation, p *recordingStartProgress) error {
	if err := s.saveRecordingStartProgress(ctx, op, p); err != nil {
		return err
	}
	for i, intent := range p.Request.Groups {
		if p.View.Groups[i].Saved {
			continue
		}
		p.View.Stage = recordingStartSaveStage
		if err := s.validateRecordingStart(ctx, p, false); err != nil {
			return err
		}
		if intent.Draft != nil {
			_, err := s.groups.updateWithCheckpoint(ctx, intent.GroupID, s.recordingStartMutation(p, i, intent.Draft), s.recordingStartCheckpoint(op, p, i, false))
			if err != nil {
				return err
			}
		} else {
			p.View.Groups[i].Saved = true
			if err := s.saveRecordingStartProgress(ctx, op, p); err != nil {
				return err
			}
		}
	}
	p.View.Stage = "readiness"
	if err := s.saveRecordingStartProgress(ctx, op, p); err != nil {
		return err
	}
	if err := s.validateRecordingStart(ctx, p, false); err != nil {
		return err
	}
	if err := s.checkRecordingStartBarrier(ctx, op, p); err != nil {
		return err
	}
	for i, intent := range p.Request.Groups {
		current, err := s.groups.Get(ctx, intent.GroupID)
		if err != nil {
			return err
		}
		if current.Group.Destination.SchemaRevision == "" || current.Group.Destination.SchemaDigest == "" {
			return ErrWriteGroupApplyNotReady
		}
		ready, err := s.groups.Readiness(ctx, intent.GroupID)
		if err != nil {
			return err
		}
		if !ready.Ready {
			return ErrWriteGroupApplyNotReady
		}
		p.View.Groups[i].Ready = true
	}
	ids := recordingStartGroupIDs(p)
	ready, err := s.workspace.ReadinessScope(ctx, p.Request.DeviceIDs, ids)
	if err != nil {
		return err
	}
	if ready == nil || !ready.Ready {
		return &ReadinessBlockedError{Operation: "recording_start", Summary: ready}
	}
	if err := s.saveRecordingStartProgress(ctx, op, p); err != nil {
		return err
	}
	for i, intent := range p.Request.Groups {
		if p.View.Groups[i].Applied {
			continue
		}
		p.View.Stage = "apply"
		if err := s.saveRecordingStartProgress(ctx, op, p); err != nil {
			return err
		}
		if err := s.validateRecordingStart(ctx, p, false); err != nil {
			return err
		}
		_, err := s.groups.applyWithCheckpoint(ctx, intent.GroupID, s.recordingStartMutation(p, i, nil), s.recordingStartCheckpoint(op, p, i, true))
		if err != nil {
			return err
		}
	}
	p.View.Stage = "activation"
	if err := s.saveRecordingStartProgress(ctx, op, p); err != nil {
		return err
	}
	if err := s.validateRecordingStart(ctx, p, false); err != nil {
		return err
	}
	if err := s.checkRecordingStartBarrier(ctx, op, p); err != nil {
		return err
	}
	if s.activate == nil {
		return ErrWriteGroupServiceUnavailable
	}
	response, err := s.activate.ActivateScope(ctx, p.Request.DeviceIDs, ids)
	if err != nil {
		return err
	}
	if response == nil || response.WorkspaceID != p.Request.WorkspaceID {
		return ErrRecordingStartInvalid
	}
	for _, result := range response.Results {
		i := slices.Index(p.Request.DeviceIDs, result.DeviceID)
		if i < 0 {
			return ErrRecordingStartInvalid
		}
		if result.Status == ActivationResultStatusSuccess {
			p.View.Devices[i].Activated, p.View.Devices[i].Reason = true, ""
		} else {
			p.View.Devices[i].Reason = "activation_failed"
		}
		if err := s.saveRecordingStartProgress(ctx, op, p); err != nil {
			return err
		}
	}
	for _, device := range p.View.Devices {
		if !device.Activated {
			return fmt.Errorf("selected device activation did not complete")
		}
	}
	return nil
}

func (s *RecordingStartService) checkRecordingStartBarrier(ctx context.Context, op *recordingplan.SchemaOperation, p *recordingStartProgress) error {
	if s.barrier == nil {
		return nil
	}
	check := s.barrier
	if p.BarrierVerified && s.resumeBarrier != nil {
		check = s.resumeBarrier
	}
	if err := check(ctx, p.Request); err != nil {
		return err
	}
	if !p.BarrierVerified {
		p.BarrierVerified = true
		return s.saveRecordingStartProgress(ctx, op, p)
	}
	return nil
}

func recordingStartGroupIDs(p *recordingStartProgress) []string {
	ids := make([]string, len(p.Request.Groups))
	for i, group := range p.Request.Groups {
		ids[i] = group.GroupID
	}
	return ids
}

func (s *RecordingStartService) recordingStartMutation(p *recordingStartProgress, i int, draft *WriteGroup) WriteGroupMutation {
	return WriteGroupMutation{
		WorkspaceID: p.Request.WorkspaceID, ExpectedWorkspaceRevision: p.WorkspaceRevision,
		ExpectedGroupRevision:     p.View.Groups[i].GroupRevision,
		ExpectedConnectorRevision: p.Request.Groups[i].ExpectedConnectorRevision, Group: draft,
	}
}

func (s *RecordingStartService) recordingStartCheckpoint(op *recordingplan.SchemaOperation, p *recordingStartProgress, i int, applied bool) func(context.Context, *sql.Tx, *Record) error {
	return func(ctx context.Context, tx *sql.Tx, record *Record) error {
		group, err := s.groups.repo.GetInTx(ctx, tx, p.Request.WorkspaceID, p.Request.Groups[i].GroupID)
		if err != nil {
			return err
		}
		types, err := loadWriteGroupTagTypes(ctx, tx, s.groups.repo.query, group)
		if err != nil {
			return err
		}
		encoded, err := json.Marshal([]any{group.Members, types})
		if err != nil {
			return fmt.Errorf("encode recording start source scope: %w", err)
		}
		digest := sha256.Sum256(encoded)
		p.Scopes[i] = recordingStartGroupScope{
			SourceDigest: hex.EncodeToString(digest[:]), SchemaRevision: group.Destination.SchemaRevision,
			SchemaDigest: group.Destination.SchemaDigest,
		}
		p.View.Groups[i].GroupRevision, p.View.Groups[i].Saved = group.Revision, true
		if applied {
			p.View.Groups[i].Applied, p.View.Groups[i].AppliedRevision = true, group.AppliedRevision
		}
		p.WorkspaceRevision = record.DatabaseSetupRevision
		encoded, err = json.Marshal(p)
		if err != nil {
			return fmt.Errorf("encode recording start transaction progress: %w", err)
		}
		return s.ledger.SaveRecordingStartProgressInTx(ctx, tx, op.OperationID, op.Owner, string(encoded))
	}
}

package workspace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
	"strings"
)

const recordingStartSaveStage = "save"

func normalizeRecordingStart(request RecordingStartRequest) (RecordingStartRequest, string, error) {
	if strings.TrimSpace(request.RequestID) == "" || len(request.RequestID) > 128 ||
		strings.TrimSpace(request.WorkspaceID) == "" || len(request.DeviceIDs) == 0 ||
		len(request.DeviceIDs) > 64 || len(request.Groups) > 64 {
		return request, "", ErrRecordingStartInvalid
	}
	encoded, err := json.Marshal(request)
	if err != nil || json.Unmarshal(encoded, &request) != nil {
		return request, "", ErrRecordingStartInvalid
	}
	slices.Sort(request.DeviceIDs)
	for i, id := range request.DeviceIDs {
		if id == "" || (i > 0 && id == request.DeviceIDs[i-1]) {
			return request, "", ErrRecordingStartInvalid
		}
	}
	slices.SortFunc(request.Groups, func(left, right RecordingStartGroupIntent) int { return strings.Compare(left.GroupID, right.GroupID) })
	for i, group := range request.Groups {
		if group.GroupID == "" || group.ExpectedGroupRevision == "" || group.ExpectedConnectorRevision == "" ||
			(i > 0 && group.GroupID == request.Groups[i-1].GroupID) {
			return request, "", ErrRecordingStartInvalid
		}
		if group.Draft != nil {
			if group.Draft.WorkspaceID != "" && group.Draft.WorkspaceID != request.WorkspaceID {
				return request, "", ErrRecordingStartInvalid
			}
			for _, member := range group.Draft.Members {
				if !slices.Contains(request.DeviceIDs, member.DeviceID) {
					return request, "", ErrRecordingStartInvalid
				}
			}
		}
	}
	encoded, err = json.Marshal(request)
	if err != nil {
		return request, "", ErrRecordingStartInvalid
	}
	digest := sha256.Sum256(encoded)
	return request, hex.EncodeToString(digest[:]), nil
}

func initialRecordingStart(request RecordingStartRequest, digest string) *recordingStartProgress {
	p := &recordingStartProgress{
		Version: 1, Request: request, WorkspaceRevision: request.ExpectedWorkspaceRevision,
		Scopes:        make([]recordingStartGroupScope, len(request.Groups)),
		DeviceDigests: make([]string, len(request.DeviceIDs)),
		View: RecordingStartOperation{
			WorkspaceID: request.WorkspaceID, DeviceIDs: slices.Clone(request.DeviceIDs), IntentDigest: digest,
			Stage: recordingStartSaveStage, Groups: []RecordingStartGroupProgress{}, Devices: []RecordingStartDeviceProgress{},
		},
	}
	for _, group := range request.Groups {
		p.View.Groups = append(p.View.Groups, RecordingStartGroupProgress{GroupID: group.GroupID, GroupRevision: group.ExpectedGroupRevision})
	}
	for _, id := range request.DeviceIDs {
		p.View.Devices = append(p.View.Devices, RecordingStartDeviceProgress{DeviceID: id})
	}
	return p
}

// Only a claim with no recorded preparation or effects may initialize again.
func recordingStartNeedsInitialization(p *recordingStartProgress) bool {
	if p.BarrierVerified || p.View.Stage != recordingStartSaveStage || p.WorkspaceRevision != p.Request.ExpectedWorkspaceRevision {
		return false
	}
	for _, progress := range p.View.Groups {
		if progress.Saved || progress.Ready || progress.Applied || progress.AppliedRevision != "" {
			return false
		}
	}
	for _, progress := range p.View.Devices {
		if progress.Activated {
			return false
		}
	}
	for _, digest := range p.DeviceDigests {
		if digest != "" {
			return false
		}
	}
	for _, scope := range p.Scopes {
		if scope != (recordingStartGroupScope{}) {
			return false
		}
	}
	return true
}

func (s *RecordingStartService) validateRecordingStartSelection(ctx context.Context, request RecordingStartRequest) error {
	record, err := s.workspace.GetOrCreate(ctx)
	if err != nil {
		return err
	}
	if record.ID != request.WorkspaceID {
		return ErrRecordingStartInvalid
	}
	for _, id := range request.DeviceIDs {
		if !slices.Contains(record.OrderedDeviceIDs, id) {
			return ErrRecordingStartInvalid
		}
	}
	for _, intent := range request.Groups {
		current, err := s.groups.Get(ctx, intent.GroupID)
		if err != nil {
			return err
		}
		if current.Group.WorkspaceID != request.WorkspaceID || current.Group.Status == WriteGroupStatusDeleted {
			return ErrRecordingStartInvalid
		}
		for _, member := range current.Group.Members {
			if !slices.Contains(request.DeviceIDs, member.DeviceID) {
				return ErrRecordingStartInvalid
			}
		}
	}
	return nil
}

func (s *RecordingStartService) validateRecordingStart(ctx context.Context, p *recordingStartProgress, initialize bool) error {
	if err := s.validateRecordingStartSelection(ctx, p.Request); err != nil {
		return err
	}
	record, err := s.workspace.GetOrCreate(ctx)
	if err != nil {
		return err
	}
	if record.DatabaseSetupRevision != p.WorkspaceRevision {
		return ErrSetupRevisionConflict
	}
	for i, id := range p.Request.DeviceIDs {
		var protocol, config string
		err := s.groups.repo.db.QueryRowContext(ctx, s.groups.repo.query(`SELECT protocol,connection_config FROM devices WHERE id=$1`), id).Scan(&protocol, &config)
		if err != nil {
			return err
		}
		digest := sha256.Sum256([]byte(protocol + "\x00" + config))
		value := hex.EncodeToString(digest[:])
		if initialize && p.DeviceDigests[i] == "" {
			p.DeviceDigests[i] = value
		} else if p.DeviceDigests[i] != value {
			return ErrWriteGroupRevisionConflict
		}
	}
	for i, intent := range p.Request.Groups {
		current, err := s.groups.Get(ctx, intent.GroupID)
		if err != nil {
			return err
		}
		group := current.Group
		if group.Revision != p.View.Groups[i].GroupRevision || group.Destination.ConnectorRevision != intent.ExpectedConnectorRevision {
			return ErrWriteGroupRevisionConflict
		}
		if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged {
			return ErrRecordingStartInvalid
		}
		scope, err := s.groups.ManagedSchemaScope(ctx, group.ID, WriteGroupMutation{
			WorkspaceID: p.Request.WorkspaceID, ExpectedWorkspaceRevision: p.WorkspaceRevision,
			ExpectedGroupRevision: group.Revision, ExpectedConnectorRevision: intent.ExpectedConnectorRevision,
		})
		if err != nil {
			return err
		}
		now := recordingStartGroupScope{SourceDigest: scope.SourceDigest, SchemaRevision: scope.SchemaRevision, SchemaDigest: scope.SchemaDigest}
		if initialize && p.Scopes[i].SourceDigest == "" {
			p.Scopes[i] = now
		} else if p.Scopes[i] != now {
			return ErrWriteGroupRevisionConflict
		}
		if p.View.Groups[i].Applied && group.AppliedRevision != p.View.Groups[i].AppliedRevision {
			return ErrWriteGroupRevisionConflict
		}
	}
	return nil
}

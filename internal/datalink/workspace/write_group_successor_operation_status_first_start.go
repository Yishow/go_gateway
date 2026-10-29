package workspace

import (
	"context"
	"database/sql"

	"go-gateway/internal/datalink/recordingplan"
)

func (s *WriteGroupService) validateSuccessorStatusFirstStart(ctx context.Context, tx *sql.Tx, confirmation *WriteGroupSuccessorActivationConfirmation, started *WriteGroupSuccessorFirstStartResult, reserved *WriteGroupSuccessorPreparation, binding *successorReserveSchemaBinding) error {
	e := confirmation.Envelope
	if started.AppliedRevision != reserved.Group.Revision || e.Request.Stage.Current.LineageRevision != reserved.LineageRevision ||
		started.LineageRevision != reserved.LineageRevision+1 || started.Generation != reserved.Generation ||
		e.Request.Stage.Current.TargetConnectorRevision != reserved.Group.Destination.ConnectorRevision {
		return ErrWriteGroupLineageCorrupt
	}
	for _, digest := range []string{confirmation.ReadinessDigest, e.TargetDefinitionDigest, e.CounterProofDigest, e.SchemaDigest} {
		if !successorStatusDigestValid(digest) {
			return ErrWriteGroupLineageCorrupt
		}
	}
	token, err := recordingplan.SucceededSchemaTokenInTx(ctx, tx, binding.WorkspaceID, e.SchemaRevision, e.SchemaDigest)
	if err != nil {
		return ErrWriteGroupLineageCorrupt
	}
	if err := successorStatusTokenScope(token, reserved, binding); err != nil {
		return err
	}
	target, err := s.repo.VersionByRevisionInTx(ctx, tx, binding.WorkspaceID, binding.GroupID, started.AppliedRevision)
	if err != nil || !target.EffectiveAt.Equal(started.Boundary) || target.Group.Destination.SchemaRevision != e.SchemaRevision || target.Group.Destination.SchemaDigest != e.SchemaDigest {
		return ErrWriteGroupLineageCorrupt
	}
	actual := e.Request.Stage.Envelope
	actual.Request.Group = cloneWriteGroup(target.Group)
	if target.Group.isWideProfile() || target.Group.RecordingProfile == nil {
		actual.Request.Group.Destination.TableName = ""
	}
	_, _, digest, err := normalizeSuccessorEnvelope(actual)
	if err != nil || digest != reserved.IntentDigest {
		return ErrWriteGroupLineageCorrupt
	}
	return nil
}

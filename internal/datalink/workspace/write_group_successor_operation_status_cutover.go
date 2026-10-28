package workspace

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/recordingplan"
)

func (s *WriteGroupService) successorStatusCutoverInTx(ctx context.Context, tx *sql.Tx, original *WriteGroupSuccessorEnvelope, reserved *WriteGroupSuccessorPreparation, binding *successorReserveSchemaBinding) (*WriteGroupCutover, error) {
	events, err := s.repo.listCutoversInTx(ctx, tx, binding.WorkspaceID, binding.PredecessorGroupID)
	if err != nil {
		return nil, err
	}
	var result *WriteGroupCutover
	for _, event := range events {
		if !isSuccessorActivationIntent(event.intent) {
			continue
		}
		var stored struct {
			Version            string                                    `json:"version"`
			PredecessorGroupID string                                    `json:"predecessor_group_id"`
			Confirmation       WriteGroupSuccessorActivationConfirmation `json:"confirmation"`
		}
		if json.Unmarshal([]byte(event.intent), &stored) != nil || stored.Version != successorActivationIntentVersion {
			return nil, ErrWriteGroupLineageCorrupt
		}
		if stored.Confirmation.Envelope.Request.Stage.Envelope.Request.OperationID != original.Request.OperationID {
			continue
		}
		confirmation, intent, err := successorActivationIntent(stored.PredecessorGroupID, stored.Confirmation)
		if err != nil || intent != event.intent || !sameSuccessorJSON(confirmation.Envelope.Request.Stage.Envelope, *original) || result != nil {
			return nil, ErrWriteGroupLineageCorrupt
		}
		if err := s.validateSuccessorStatusCutover(ctx, tx, event, confirmation, reserved, binding); err != nil {
			return nil, err
		}
		result = event
	}
	return result, nil
}

func (s *WriteGroupService) validateSuccessorStatusCutover(ctx context.Context, tx *sql.Tx, event *WriteGroupCutover, confirmation WriteGroupSuccessorActivationConfirmation, reserved *WriteGroupSuccessorPreparation, binding *successorReserveSchemaBinding) error {
	e := confirmation.Envelope
	cas := e.Request.Stage.Current
	if successorActivationKind(e) != successorActivationCutover || event.ID == "" || event.WorkspaceID != binding.WorkspaceID || event.SourceGroupID != binding.PredecessorGroupID ||
		event.Kind != WriteGroupCutoverForward || event.PriorCutoverID != "" || event.ReviewDigest != confirmation.ReviewDigest ||
		event.OwnershipFenceRevision != event.ID || !event.Boundary.Equal(e.Request.Boundary) ||
		len(event.Ended) != 1 || len(event.Started) != 1 || event.Ended[0] != (CutoverEnd{GroupID: binding.PredecessorGroupID, Revision: e.SourceAppliedRevision}) ||
		event.Started[0] != (CutoverStart{GroupID: binding.GroupID, Revision: reserved.Group.Revision}) ||
		e.SourceAppliedRevision == "" || e.SourceAppliedRevision != e.Request.Stage.Envelope.SourceAppliedRevision ||
		cas.GroupID != binding.GroupID || cas.GroupRevision != reserved.Group.Revision || cas.LineageRevision != reserved.LineageRevision ||
		cas.TargetConnectorRevision != reserved.Group.Destination.ConnectorRevision {
		return ErrWriteGroupLineageCorrupt
	}
	for _, digest := range []string{confirmation.ReadinessDigest, e.SourceDefinitionDigest, e.TargetDefinitionDigest, e.CounterProofDigest, e.SchemaDigest} {
		if !successorStatusDigestValid(digest) {
			return ErrWriteGroupLineageCorrupt
		}
	}
	wanted := make([]CutoverMemberRef, 0, len(reserved.Group.Members))
	for _, member := range reserved.Group.Members {
		wanted = append(wanted, memberRefOf(member))
	}
	sortMemberRefs(wanted)
	actual := append([]CutoverMemberRef(nil), event.Moved...)
	sortMemberRefs(actual)
	if !slices.Equal(wanted, actual) {
		return ErrWriteGroupLineageCorrupt
	}
	token, err := recordingplan.SucceededSchemaTokenInTx(ctx, tx, binding.WorkspaceID, e.SchemaRevision, e.SchemaDigest)
	if err != nil {
		return ErrWriteGroupLineageCorrupt
	}
	if err := successorStatusTokenScope(token, reserved, binding); err != nil {
		return err
	}
	source, err := s.repo.VersionByRevisionInTx(ctx, tx, binding.WorkspaceID, binding.PredecessorGroupID, e.SourceAppliedRevision)
	if err != nil {
		return ErrWriteGroupLineageCorrupt
	}
	digest, err := successorDomainDigest("successor-activation-source/v1", source)
	if err != nil || digest != e.SourceDefinitionDigest || !event.Boundary.After(source.EffectiveAt) {
		return ErrWriteGroupLineageCorrupt
	}
	target, err := s.repo.VersionByRevisionInTx(ctx, tx, binding.WorkspaceID, binding.GroupID, reserved.Group.Revision)
	if err != nil || !target.EffectiveAt.Equal(event.Boundary) || target.Group.Destination.SchemaRevision != e.SchemaRevision || target.Group.Destination.SchemaDigest != e.SchemaDigest {
		return ErrWriteGroupLineageCorrupt
	}
	actualEnvelope := e.Request.Stage.Envelope
	actualEnvelope.Request.Group = cloneWriteGroup(target.Group)
	if target.Group.isWideProfile() || target.Group.RecordingProfile == nil {
		actualEnvelope.Request.Group.Destination.TableName = ""
	}
	_, _, actualDigest, err := normalizeSuccessorEnvelope(actualEnvelope)
	if err != nil || actualDigest != reserved.IntentDigest {
		return ErrWriteGroupLineageCorrupt
	}
	var ends, cutoverID string
	if err := tx.QueryRowContext(ctx, s.repo.query(`SELECT ends_at,cutover_id FROM write_group_revision_ends WHERE workspace_id=$1 AND group_id=$2 AND group_revision=$3`), binding.WorkspaceID, binding.PredecessorGroupID, e.SourceAppliedRevision).Scan(&ends, &cutoverID); err != nil {
		return ErrWriteGroupLineageCorrupt
	}
	boundary, err := time.Parse(cutoverTimeLayout, ends)
	if err != nil || !boundary.Equal(event.Boundary) || cutoverID != event.ID {
		return ErrWriteGroupLineageCorrupt
	}
	return s.successorStatusFenceInTx(ctx, tx, event, source.Group, target.Group)
}

func successorStatusDigestValid(value string) bool {
	decoded, err := hex.DecodeString(value)
	return err == nil && len(decoded) == 32 && strings.ToLower(value) == value
}

func (s *WriteGroupService) successorStatusFenceInTx(ctx context.Context, tx *sql.Tx, event *WriteGroupCutover, source, target *WriteGroup) error {
	var total int
	if err := tx.QueryRowContext(ctx, s.repo.query(`SELECT COUNT(*) FROM write_group_ownership_fences WHERE fence_revision=$1`), event.ID).Scan(&total); err != nil {
		return err
	}
	if total != len(source.Members)+len(target.Members) {
		return ErrWriteGroupLineageCorrupt
	}
	for _, side := range []struct {
		direction string
		group     *WriteGroup
	}{{"ended", source}, {"started", target}} {
		g := side.group
		for _, member := range g.Members {
			var count int
			err := tx.QueryRowContext(ctx, s.repo.query(`SELECT COUNT(*) FROM write_group_ownership_fences WHERE fence_revision=$1 AND workspace_id=$2 AND direction=$3 AND group_id=$4 AND group_revision=$5 AND device_id=$6 AND point_id=$7 AND tag_id=$8 AND profile_kind=$9 AND connector_id=$10 AND connector_revision=$11 AND database_name=$12 AND table_schema=$13 AND table_name=$14`), event.ID, event.WorkspaceID, side.direction, g.ID, g.Revision, member.DeviceID, member.PointID, member.TagID, profileKindOf(g), g.Destination.ConnectorID, g.Destination.ConnectorRevision, g.Destination.Database, g.Destination.TableSchema, g.Destination.TableName).Scan(&count)
			if err != nil {
				return err
			}
			if count != 1 {
				return ErrWriteGroupLineageCorrupt
			}
		}
	}
	return nil
}

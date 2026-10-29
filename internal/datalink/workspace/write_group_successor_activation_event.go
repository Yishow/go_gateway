package workspace

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
)

// successorActivationIntent stores the complete frozen request, not just a
// digest or the mutable canonical head. It remains usable after reservation
// release, restart, and a lost commit acknowledgement.
func successorActivationIntent(id string, req WriteGroupSuccessorActivationConfirmation) (WriteGroupSuccessorActivationConfirmation, string, error) {
	id = strings.TrimSpace(id)
	encoded, err := json.Marshal(req)
	if err != nil {
		return req, "", ErrWriteGroupValidation
	}
	var frozen WriteGroupSuccessorActivationConfirmation
	if json.Unmarshal(encoded, &frozen) != nil {
		return req, "", ErrWriteGroupValidation
	}
	e, review, _, err := normalizeSuccessorEnvelope(frozen.Envelope.Request.Stage.Envelope)
	if err != nil {
		return req, "", err
	}
	if id == "" || id != e.PredecessorGroupID || review != frozen.Envelope.Request.Stage.ReviewDigest || frozen.ReadinessDigest == "" {
		return req, "", ErrWriteGroupValidation
	}
	frozen.Envelope.Request.Stage.Envelope = e
	frozen.Envelope.Request.Boundary = frozen.Envelope.Request.Boundary.UTC()
	normalized, digest, err := normalizeSuccessorActivationEnvelope(frozen.Envelope)
	if err != nil {
		return req, "", err
	}
	frozen.Envelope = normalized
	if digest != frozen.ReviewDigest {
		return req, "", writeGroupRevisionConflict("successor activation review is stale")
	}
	encoded, err = json.Marshal(struct {
		Version            string                                    `json:"version"`
		PredecessorGroupID string                                    `json:"predecessor_group_id"`
		Confirmation       WriteGroupSuccessorActivationConfirmation `json:"confirmation"`
	}{successorActivationIntentVersion, id, frozen})
	if err != nil {
		return req, "", ErrWriteGroupValidation
	}
	return frozen, string(encoded), nil
}

func (s *WriteGroupService) recordedSuccessorActivation(ctx context.Context, req WriteGroupSuccessorActivationConfirmation, intent string) (*WriteGroupCutover, error) {
	c, stored, err := s.repo.cutoverByDigestInTx(ctx, s.repo.db, req.Envelope.Request.Stage.Envelope.Request.WorkspaceID, req.ReviewDigest)
	if err != nil || c == nil {
		return c, err
	}
	if stored != intent || !isSuccessorActivationIntent(stored) {
		return nil, writeGroupRevisionConflict("successor activation belongs to another intent")
	}
	if c.Kind != WriteGroupCutoverForward || c.SourceGroupID != req.Envelope.Request.Stage.Envelope.PredecessorGroupID || !c.Boundary.Equal(req.Envelope.Request.Boundary) || len(c.Ended) != 1 || len(c.Started) != 1 || c.Ended[0].GroupID != c.SourceGroupID || c.Ended[0].Revision != req.Envelope.SourceAppliedRevision || c.Started[0].GroupID != req.Envelope.Request.Stage.Current.GroupID || c.Started[0].Revision != req.Envelope.Request.Stage.Current.GroupRevision || c.OwnershipFenceRevision != c.ID {
		return nil, ErrWriteGroupLineageCorrupt
	}
	group := req.Envelope.Request.Stage.Envelope.Request.Group
	if group == nil || c.PriorCutoverID != "" {
		return nil, ErrWriteGroupLineageCorrupt
	}
	wanted := make([]CutoverMemberRef, 0, len(group.Members))
	for _, member := range group.Members {
		wanted = append(wanted, memberRefOf(member))
	}
	sortMemberRefs(wanted)
	actual := append([]CutoverMemberRef(nil), c.Moved...)
	sortMemberRefs(actual)
	if !slices.Equal(wanted, actual) {
		return nil, ErrWriteGroupLineageCorrupt
	}
	return c, nil
}

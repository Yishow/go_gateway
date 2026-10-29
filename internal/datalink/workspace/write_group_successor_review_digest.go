package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
)

const successorReviewVersion = "successor-review-v1"

func normalizeSuccessorReviewRequest(id string, req WriteGroupSuccessorReviewRequest) (WriteGroupSuccessorReviewRequest, error) {
	if req.Boundary.IsZero() || req.Boundary.Year() < 1 || req.Boundary.Year() > 9999 {
		return req, ErrWriteGroupValidation
	}
	normalized, _, err := normalizeRecreateReservationRequest(writeGroupRecreateReservationRequest{
		WorkspaceID: req.WorkspaceID, ExpectedWorkspaceRevision: req.ExpectedWorkspaceRevision, LineageID: req.LineageID,
		ExpectedLineageRevision: req.ExpectedLineageRevision, PredecessorGroupID: id,
		ExpectedPredecessorRevision: req.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision,
		ExpectedTargetConnectorRevision: req.ExpectedTargetConnectorRevision, OperationID: req.OperationID, Group: req.Group,
	})
	if err != nil {
		return req, err
	}
	if (normalized.Group.isWideProfile() || normalized.Group.RecordingProfile == nil && normalized.Group.EffectiveLayoutVersion() == WriteGroupLayoutManagedWideV1) && normalized.Group.Destination.TableName != "" {
		return req, transitionRejection("wide-table-identity-change")
	}
	return WriteGroupSuccessorReviewRequest{WorkspaceID: normalized.WorkspaceID, ExpectedWorkspaceRevision: normalized.ExpectedWorkspaceRevision,
		LineageID: normalized.LineageID, ExpectedLineageRevision: normalized.ExpectedLineageRevision, ExpectedPredecessorRevision: normalized.ExpectedPredecessorRevision,
		ExpectedPredecessorConnectorRevision: normalized.ExpectedPredecessorConnectorRevision, ExpectedTargetConnectorRevision: normalized.ExpectedTargetConnectorRevision,
		OperationID: normalized.OperationID, Boundary: req.Boundary.UTC(), Group: normalized.Group}, nil
}

// normalizeSuccessorEnvelope is independent of live state, so exact retry can
// prove its original identity after legitimate schema/workspace CAS advances.
func normalizeSuccessorEnvelope(env WriteGroupSuccessorEnvelope) (WriteGroupSuccessorEnvelope, string, string, error) {
	env.PredecessorGroupID = strings.TrimSpace(env.PredecessorGroupID)
	if env.Version != successorReviewVersion || env.Generation < 1 || env.PredecessorGroupID == "" || len(env.SourceDefinitionDigest) != 64 {
		return env, "", "", ErrWriteGroupValidation
	}
	decoded, err := hex.DecodeString(env.SourceDefinitionDigest)
	if err != nil || len(decoded) != 32 || strings.ToLower(env.SourceDefinitionDigest) != env.SourceDefinitionDigest {
		return env, "", "", ErrWriteGroupValidation
	}
	env.Request, err = normalizeSuccessorReviewRequest(env.PredecessorGroupID, env.Request)
	if err != nil {
		return env, "", "", err
	}
	review, err := successorDomainDigest("go-gateway/successor-review/v1", env)
	if err != nil {
		return env, "", "", err
	}
	reservation, err := successorDomainDigest("go-gateway/successor-reservation/v1", struct {
		Envelope     WriteGroupSuccessorEnvelope
		ReviewDigest string
	}{env, review})
	return env, review, reservation, err
}

func successorDomainDigest(domain string, value any) (string, error) {
	encoded, err := json.Marshal(struct {
		Domain string `json:"domain"`
		Value  any    `json:"value"`
	}{domain, value})
	if err != nil {
		return "", ErrWriteGroupValidation
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}

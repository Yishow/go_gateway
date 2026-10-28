package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// Cutover kinds.
const (
	// WriteGroupCutoverForward moves a wide or custom group's sources to long
	// profile groups at one exact instant.
	WriteGroupCutoverForward = "forward"
	// WriteGroupCutoverRollback returns those sources to the group that held
	// them, again at one exact, new instant. It never touches history.
	WriteGroupCutoverRollback = "rollback"
)

var (
	// ErrCutoverNotApplicable means the reviewed difference cannot be applied;
	// the issues say why and nothing was changed.
	ErrCutoverNotApplicable = errors.New("write group cutover is not applicable")
	// ErrWriteGroupCutoverRequired means a long group would take over sources an
	// active group still owns; only a reviewed cutover may move them.
	ErrWriteGroupCutoverRequired = errors.New("sources owned by an active group need a reviewed cutover")
)

// Stable issue codes of a cutover review. They name facts, never values.
const (
	CutoverIssueSourceUnavailable    = "source_unavailable"
	CutoverIssueSourceNotApplied     = "source_not_applied"
	CutoverIssueSourceDraftPending   = "source_has_unapplied_draft"
	CutoverIssueSourceEnded          = "source_already_ended"
	CutoverIssueSourceIsLong         = "source_is_long"
	CutoverIssueTargetNotLong        = "target_not_long"
	CutoverIssueTargetActive         = "target_already_active"
	CutoverIssueTargetNotReady       = "target_not_ready"
	CutoverIssueTargetOutsideSource  = "target_member_not_in_source"
	CutoverIssueTargetsOverlap       = "targets_overlap"
	CutoverIssueUnselectedSources    = "unselected_sources"
	CutoverIssueBoundaryNotFuture    = "boundary_not_future"
	CutoverIssueBoundaryBeforeActive = "boundary_before_active_revision"
	CutoverIssueNoTargets            = "no_targets"
	CutoverIssueNotRollbackable      = "cutover_not_rollbackable"
	// The remaining wide revision (the saved draft that keeps the sources that
	// are not moved) is not what the cutover can carry out.
	CutoverIssueRemainingMismatch      = "remaining_draft_mismatch"
	CutoverIssueRemainingNotReady      = "remaining_draft_not_ready"
	CutoverIssueRemainingChangesPolicy = "remaining_draft_changes_policy"
	// CutoverIssueRemainingOmitsColumn: the remaining wide revision would stop
	// writing a column the destination declares NOT NULL without a default, so its
	// first insert would fail.
	CutoverIssueRemainingOmitsColumn = "remaining_draft_omits_required_column"
	// CutoverIssueRestoredNotReady: the group a rollback would bring back is not
	// ready against the real sources and destination now.
	CutoverIssueRestoredNotReady = "restored_group_not_ready"
)

// CutoverMemberRef names one source by its stable identity.
type CutoverMemberRef struct {
	DeviceID string `json:"device_id"`
	PointID  string `json:"point_id"`
	TagID    string `json:"tag_id"`
}

// CutoverTargetRef names a long group and the revision the reviewer saw.
type CutoverTargetRef struct {
	GroupID          string `json:"group_id"`
	ExpectedRevision string `json:"expected_revision"`
}

// CutoverRequest is the shared shape of a review and its apply. Apply must
// carry the digest the review returned, so what was reviewed is what runs.
type CutoverRequest struct {
	Kind                      string             `json:"kind"`
	WorkspaceID               string             `json:"workspace_id"`
	ExpectedWorkspaceRevision string             `json:"expected_workspace_revision"`
	SourceGroupID             string             `json:"source_group_id"`
	ExpectedSourceRevision    string             `json:"expected_source_revision"`
	ExpectedConnectorRevision string             `json:"expected_connector_revision"`
	Targets                   []CutoverTargetRef `json:"targets,omitempty"`
	// CutoverID names the forward cutover a rollback reverses.
	CutoverID string `json:"cutover_id,omitempty"`
	// Boundary is the exact instant the sources change hands. Zero asks the
	// review to propose one; apply requires the reviewed value.
	Boundary     time.Time `json:"boundary"`
	ReviewDigest string    `json:"review_digest,omitempty"`
}

// CutoverGroupView is one side of the reviewed difference.
type CutoverGroupView struct {
	GroupID       string `json:"group_id"`
	Name          string `json:"name"`
	Revision      string `json:"revision"`
	Layout        string `json:"layout"`
	StorageMode   string `json:"storage_mode,omitempty"`
	ProfileKind   string `json:"profile_kind,omitempty"`
	IntervalSecs  int    `json:"interval_seconds"`
	ConnectorID   string `json:"connector_id"`
	TableSchema   string `json:"table_schema"`
	TableName     string `json:"table_name"`
	MemberCount   int    `json:"member_count"`
	SchemaConfirm bool   `json:"schema_confirmed"`
}

// CutoverReview is the read-only difference an operator confirms.
type CutoverReview struct {
	Kind       string             `json:"kind"`
	Boundary   time.Time          `json:"boundary"`
	Digest     string             `json:"digest"`
	Applicable bool               `json:"applicable"`
	Issues     []string           `json:"issues"`
	Source     CutoverGroupView   `json:"source"`
	Targets    []CutoverGroupView `json:"targets"`
	Moved      []CutoverMemberRef `json:"moved"`
	// Unselected are the sources of the applied group that no target takes.
	Unselected []CutoverMemberRef `json:"unselected"`
	// Remaining are the sources of the saved draft that keeps recording in the
	// wide group; it must name exactly the unselected ones.
	Remaining []CutoverMemberRef `json:"remaining"`
}

// CutoverEnd is a revision that stops owning its sources at the boundary.
type CutoverEnd struct {
	GroupID  string `json:"group_id"`
	Revision string `json:"revision"`
}

// CutoverStart is a revision that starts owning sources at the boundary.
type CutoverStart struct {
	GroupID  string `json:"group_id"`
	Revision string `json:"revision"`
}

// WriteGroupCutover is the immutable record of one cutover.
type WriteGroupCutover struct {
	ID             string             `json:"id"`
	WorkspaceID    string             `json:"workspace_id"`
	Kind           string             `json:"kind"`
	SourceGroupID  string             `json:"source_group_id"`
	Boundary       time.Time          `json:"boundary"`
	PriorCutoverID string             `json:"prior_cutover_id,omitempty"`
	ReviewDigest   string             `json:"review_digest"`
	Ended          []CutoverEnd       `json:"ended"`
	Started        []CutoverStart     `json:"started"`
	Moved          []CutoverMemberRef `json:"moved"`
	CreatedAt      time.Time          `json:"created_at"`

	// OwnershipFenceRevision is an opaque persisted CAS identity. Empty means
	// legacy/unknown ownership evidence, not permission to accept or close.
	OwnershipFenceRevision string `json:"ownership_fence_revision,omitempty"`
	// intent is the whole request this cutover answered; it is what a replay is
	// compared with and is not part of the public record.
	intent string
}

// cutoverLead is how far ahead a proposed boundary lies, so a reviewer has time
// to confirm and apply before the instant passes.
const cutoverLead = time.Minute

func memberRefOf(m WriteGroupMember) CutoverMemberRef {
	return CutoverMemberRef{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID}
}

func sortMemberRefs(refs []CutoverMemberRef) {
	slices.SortFunc(refs, func(a, b CutoverMemberRef) int {
		return strings.Compare(a.DeviceID+"\x00"+a.PointID+"\x00"+a.TagID, b.DeviceID+"\x00"+b.PointID+"\x00"+b.TagID)
	})
}

func cutoverGroupView(group *WriteGroup, confirmed bool) CutoverGroupView {
	view := CutoverGroupView{
		GroupID: group.ID, Name: group.Name, Revision: group.Revision, Layout: string(group.EffectiveLayoutVersion()),
		IntervalSecs: group.RowPolicy.IntervalSeconds, ConnectorID: group.Destination.ConnectorID,
		TableSchema: group.Destination.TableSchema, TableName: group.Destination.TableName,
		MemberCount: len(group.Members), SchemaConfirm: confirmed,
	}
	if group.isLongLayout() {
		view.StorageMode = string(group.StorageMode)
		if group.RecordingProfile != nil {
			view.ProfileKind = string(group.RecordingProfile.Kind)
		}
	}
	return view
}

// cutoverDigest binds a review to everything it was computed from: the
// workspace revision, the revisions and schema confirmations of every group it
// read, the sources that move and stay, and the instant. Changing any of them
// makes the review stale.
func cutoverDigest(req CutoverRequest, moved, remaining []CutoverMemberRef, facts []string, boundary time.Time) (string, error) {
	targets := slices.Clone(req.Targets)
	slices.SortFunc(targets, func(a, b CutoverTargetRef) int { return strings.Compare(a.GroupID, b.GroupID) })
	facts = slices.Clone(facts)
	slices.Sort(facts)
	encoded, err := json.Marshal(struct {
		Kind, Workspace, WorkspaceRevision, Source, SourceRevision, Connector, Cutover string
		Targets                                                                        []CutoverTargetRef
		Moved, Remaining                                                               []CutoverMemberRef
		Facts                                                                          []string
		Boundary                                                                       string
	}{req.Kind, req.WorkspaceID, req.ExpectedWorkspaceRevision, req.SourceGroupID, req.ExpectedSourceRevision, req.ExpectedConnectorRevision, req.CutoverID,
		targets, moved, remaining, facts, boundary.UTC().Format(time.RFC3339Nano)})
	if err != nil {
		return "", fmt.Errorf("encode cutover review: %w", err)
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}

// groupFact names the state of one group a review read: its revisions and the
// schema confirmation that makes its layout trusted.
func groupFact(role string, group *WriteGroup) string {
	return fmt.Sprintf("%s|%s|rev=%s|applied=%s|connector=%s|schema=%s/%s", role, group.ID, group.Revision, group.AppliedRevision,
		group.Destination.ConnectorRevision, group.Destination.SchemaRevision, group.Destination.SchemaDigest)
}

// cutoverIntent is the whole request in one comparable form. A replay is
// recognized only when the intent is identical, never by the digest alone.
func cutoverIntent(req CutoverRequest) (string, error) {
	targets := slices.Clone(req.Targets)
	slices.SortFunc(targets, func(a, b CutoverTargetRef) int { return strings.Compare(a.GroupID, b.GroupID) })
	req.Targets = targets
	req.Boundary = req.Boundary.UTC().Truncate(time.Microsecond)
	encoded, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("encode cutover intent: %w", err)
	}
	return string(encoded), nil
}

func normalizeCutoverRequest(req CutoverRequest) (CutoverRequest, error) {
	req.WorkspaceID = strings.TrimSpace(req.WorkspaceID)
	req.ExpectedWorkspaceRevision = strings.TrimSpace(req.ExpectedWorkspaceRevision)
	req.SourceGroupID = strings.TrimSpace(req.SourceGroupID)
	req.ExpectedSourceRevision = strings.TrimSpace(req.ExpectedSourceRevision)
	req.ExpectedConnectorRevision = strings.TrimSpace(req.ExpectedConnectorRevision)
	req.CutoverID = strings.TrimSpace(req.CutoverID)
	if req.Kind == "" {
		req.Kind = WriteGroupCutoverForward
	}
	if req.WorkspaceID == "" || req.SourceGroupID == "" || req.ExpectedSourceRevision == "" || req.ExpectedConnectorRevision == "" ||
		(req.Kind != WriteGroupCutoverForward && req.Kind != WriteGroupCutoverRollback) ||
		(req.Kind == WriteGroupCutoverRollback && req.CutoverID == "") {
		return CutoverRequest{}, fmt.Errorf("write group cutover request: %w", ErrWriteGroupValidation)
	}
	req.Targets = slices.Clone(req.Targets)
	for i := range req.Targets {
		req.Targets[i].GroupID = strings.TrimSpace(req.Targets[i].GroupID)
		req.Targets[i].ExpectedRevision = strings.TrimSpace(req.Targets[i].ExpectedRevision)
		if req.Targets[i].GroupID == "" || req.Targets[i].ExpectedRevision == "" {
			return CutoverRequest{}, fmt.Errorf("write group cutover target: %w", ErrWriteGroupValidation)
		}
	}
	if !req.Boundary.IsZero() {
		req.Boundary = req.Boundary.UTC().Truncate(time.Microsecond)
	}
	return req, nil
}

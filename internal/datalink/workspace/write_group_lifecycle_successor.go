package workspace

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

// WriteGroupSuccessorSeed is sanitized candidate meaning, not review, schema or
// mutation authority. Its opaque digest is informational; review always reloads
// its own immutable source and independently validates all caller input and CAS.
type WriteGroupSuccessorSeed struct {
	BasisKind     string      `json:"basis_kind"`
	BasisRevision string      `json:"basis_revision"`
	BasisDigest   string      `json:"basis_digest"`
	Group         *WriteGroup `json:"group"`
}

func (s *WriteGroupService) lifecycleSuccessorSeed(ctx context.Context, tx *sql.Tx, group *WriteGroup, history *WriteGroupLifecycleHistory, actions []WriteGroupLifecycleAction) *WriteGroupSuccessorSeed {
	block := func(reason, next string) *WriteGroupSuccessorSeed {
		for i := 3; i < len(actions); i++ {
			actions[i].Allowed, actions[i].ReasonCode, actions[i].NextAction = false, reason, next
		}
		return nil
	}
	if history == nil || !history.Canonical {
		return block("canonical_required", "review_lifecycle")
	}
	if history.ReservedGroupID != "" {
		return block("successor_reserved", "review_lifecycle")
	}
	if history.CurrentGroupID != group.ID {
		return block("historical_generation", "review_lifecycle")
	}
	if group.Status != WriteGroupStatusDeleted {
		if group.RecordingProfile == nil {
			return block("active_basic_unsupported", "review_lifecycle")
		}
		if group.Status == WriteGroupStatusDraft && group.AppliedRevision != "" {
			return block("source_draft_not_active", "review_lifecycle")
		}
		if group.AppliedRevision == "" {
			return block("source_definition_unavailable", "review_lifecycle")
		}
	}
	if err := validateLiveWriteGroupConnector(ctx, tx, s.repo, group.Destination.ConnectorID, group.Destination.ConnectorRevision); err != nil {
		return block(lifecycleSuccessorReason(err), "review_connector")
	}
	if group.AppliedRevision == "" {
		// A missing pointer alone is not proof that the saved draft was never applied.
		if err := s.repo.requireSuccessorNoHistoryInTx(ctx, tx, group.WorkspaceID, history.LineageID, group.ID); err != nil {
			return block("source_definition_unavailable", "review_lifecycle")
		}
	} else if s.backlogGuard == nil || s.backlogGuard.CheckWriteGroupBacklog(ctx, tx, group) != nil {
		return block("ownership_unavailable", "inspect_ownership")
	}
	source, err := s.repo.recreateStructuralSourceInTx(ctx, tx, group)
	if err != nil || source == nil {
		return block("source_definition_unavailable", "review_lifecycle")
	}
	if err := validateLiveWriteGroupConnector(ctx, tx, s.repo, source.Destination.ConnectorID, source.Destination.ConnectorRevision); err != nil {
		return block(lifecycleSuccessorReason(err), "review_connector")
	}
	candidate := lifecycleSuccessorCandidate(source, group.BasicManagedDeviceID)
	if err := s.lifecycleSuccessorSourceInTx(ctx, tx, group, source, cloneWriteGroup(candidate)); err != nil {
		return block(lifecycleSuccessorReason(err), "review_lifecycle")
	}
	digest, err := successorDomainDigest("go-gateway/lifecycle-successor-basis/v1", source)
	if err != nil {
		return block("source_definition_unavailable", "review_lifecycle")
	}
	seed := &WriteGroupSuccessorSeed{BasisKind: "applied", BasisRevision: group.AppliedRevision, BasisDigest: digest, Group: candidate}
	if group.AppliedRevision == "" {
		seed.BasisKind, seed.BasisRevision = "saved", group.Revision
	}
	if group.Status == WriteGroupStatusDeleted {
		actions[3].Allowed, actions[3].ReasonCode, actions[3].NextAction = true, "", "review_lifecycle"
		actions[4].ReasonCode, actions[4].NextAction = "recording_removed", "recreate"
	} else {
		actions[3].ReasonCode, actions[3].NextAction = "recording_not_removed", "review_lifecycle"
		actions[4].Allowed, actions[4].ReasonCode, actions[4].NextAction = true, "", "review_lifecycle"
	}
	return seed
}

// Validate only that the canonical source can enter the supported review flow.
// Checking source against itself is intentional: no target has been selected,
// and this must not confuse destination-unchanged with source ineligibility.
func (s *WriteGroupService) lifecycleSuccessorSourceInTx(ctx context.Context, tx *sql.Tx, predecessor, source, candidate *WriteGroup) error {
	if err := s.repo.validateSuccessorLocalDefinitionInTx(ctx, tx, predecessor, candidate); err != nil {
		return err
	}
	if source.RecordingProfile == nil {
		return s.repo.validateRecreatedBasicInTx(ctx, tx, predecessor, candidate)
	}
	if _, err := s.repo.recreatedProfileIdentityInTx(ctx, tx, predecessor, candidate); err != nil {
		return err
	}
	engine, err := s.repo.recreateSourceEngineInTx(ctx, tx, source)
	if err != nil && predecessor.AppliedRevision == "" {
		// As in the existing local review, a never-applied MySQL driver can open
		// review, but source and target family both require fresh reserve-time proof.
		var kind string
		if readErr := tx.QueryRowContext(ctx, s.repo.query(`SELECT kind FROM database_connectors WHERE id=$1`), source.Destination.ConnectorID).Scan(&kind); readErr == nil && strings.ToLower(strings.TrimSpace(kind)) == writeGroupConnectorKindMySQL {
			engine, err = writeGroupConnectorKindMySQL, nil
		}
	}
	if err != nil {
		return err
	}
	tags, err := loadWriteGroupTagTypes(ctx, tx, s.repo.query, source)
	if err != nil {
		return err
	}
	return ValidateWriteGroupTransitionStructure(source, candidate, engine, engine, tags)
}

func lifecycleSuccessorCandidate(source *WriteGroup, basicDeviceID string) *WriteGroup {
	candidate := cloneWriteGroup(source)
	candidate.ID, candidate.Revision, candidate.AppliedRevision = "", "", ""
	candidate.BasicManagedDeviceID = basicDeviceID
	candidate.Status = WriteGroupStatusDraft
	candidate.CreatedAt, candidate.UpdatedAt = time.Time{}, time.Time{}
	candidate.Migration = WriteGroupMigration{}
	candidate.Destination.SchemaRevision, candidate.Destination.SchemaDigest = "", ""
	if candidate.isWideProfile() || candidate.RecordingProfile == nil {
		candidate.Destination.TableName = ""
	}
	return candidate
}

func lifecycleSuccessorReason(err error) string {
	if errors.Is(err, ErrWriteGroupConnectorRevisionConflict) {
		return "connector_identity_stale"
	}
	if errors.Is(err, ErrWriteGroupConnectorNotFound) {
		return "connector_unavailable"
	}
	if errors.Is(err, ErrWriteGroupCircuitProfileConflict) {
		return "ownership_unavailable"
	}
	// Never reflect raw diagnostics, driver errors or arbitrary private evidence.
	return "source_definition_unavailable"
}

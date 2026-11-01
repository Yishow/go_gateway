package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
)

func (r *SQLWriteGroupRepository) validateDraftUpdate(ctx context.Context, runner writeGroupSQLRunner, group, existing *WriteGroup) error {
	if err := r.validateDefinitionForDraftUpdate(ctx, runner, group, make(map[string]*measurement.MeasurementDefinition), existing); err != nil {
		return err
	}
	return r.checkProfileSourceOwnership(ctx, runner, group)
}

func savedDraftMember(group, existing *WriteGroup, member *WriteGroupMember) *WriteGroupMember {
	if existing == nil || existing.Status == WriteGroupStatusDeleted || group.ID != existing.ID || group.WorkspaceID != existing.WorkspaceID {
		return nil
	}
	for i := range existing.Members {
		prior := &existing.Members[i]
		if member.DeviceID == prior.DeviceID && member.PointID == prior.PointID && member.TagID == prior.TagID {
			return prior
		}
	}
	return nil
}

const writeGroupMappingSelect = `SELECT id, point_id, tag_id, transform_pipeline, status,
	COALESCE(rule_candidate_id, ''), COALESCE(proposed_signature, ''),
	COALESCE(last_applied_signature, ''), COALESCE(blocking_reason, ''), enabled
	FROM mappings WHERE point_id = $1 AND tag_id = $2 AND enabled `

func scanWriteGroupMapping(row interface{ Scan(...any) error }) (mappingRevisionSource, error) {
	var mapping mappingRevisionSource
	err := row.Scan(&mapping.ID, &mapping.PointID, &mapping.TagID, &mapping.TransformPipeline, &mapping.Status,
		&mapping.RuleCandidateID, &mapping.ProposedSignature, &mapping.LastAppliedSignature, &mapping.BlockingReason, &mapping.Enabled)
	return mapping, err
}

func (r *SQLWriteGroupRepository) readWriteGroupMapping(ctx context.Context, runner writeGroupSQLRunner, member, prior *WriteGroupMember) (mappingRevisionSource, error) {
	activeQuery := writeGroupMappingSelect + `= ` + r.trueLiteral() + ` ORDER BY id ASC LIMIT 1`
	mapping, err := scanWriteGroupMapping(runner.QueryRowContext(ctx, r.query(activeQuery), member.PointID, member.TagID))
	if !errors.Is(err, sql.ErrNoRows) || prior == nil {
		return mapping, err
	}
	// A stopped source remains draft-editable only when its complete mapping
	// identity is still the one already saved by the server. Client IDs or an
	// arbitrary surviving point/tag pair are not continuity evidence.
	rows, err := runner.QueryContext(ctx, r.query(writeGroupMappingSelect+`<> `+r.trueLiteral()+` ORDER BY id ASC`), member.PointID, member.TagID)
	if err != nil {
		return mappingRevisionSource{}, err
	}
	var found *mappingRevisionSource
	hasRows := false
	for rows.Next() {
		hasRows = true
		candidate, scanErr := scanWriteGroupMapping(rows)
		if scanErr != nil {
			return mappingRevisionSource{}, errors.Join(scanErr, rows.Close())
		}
		matches, matchErr := stoppedMappingMatchesSaved(candidate, prior.MappingRevision)
		if matchErr != nil {
			return mappingRevisionSource{}, errors.Join(matchErr, rows.Close())
		}
		if matches {
			if found != nil {
				return mappingRevisionSource{}, errors.Join(draftSourceRevisionConflict("saved mapping identity is ambiguous"), rows.Close())
			}
			found = &candidate
		}
	}
	rowErr := errors.Join(rows.Err(), rows.Close())
	if rowErr != nil {
		return mappingRevisionSource{}, rowErr
	}
	if !hasRows {
		return mappingRevisionSource{}, sql.ErrNoRows
	}
	if found == nil {
		return mappingRevisionSource{}, draftSourceRevisionConflict("saved mapping identity is unavailable")
	}
	if found.RuleCandidateID != "" {
		if err := r.validateStoppedRuleMapping(ctx, runner, member, *found); err != nil {
			return mappingRevisionSource{}, err
		}
	}
	return *found, nil
}

func (r *SQLWriteGroupRepository) validateStoppedRuleMapping(ctx context.Context, runner writeGroupSQLRunner, member *WriteGroupMember, mapping mappingRevisionSource) error {
	// Count all owners, not only matching links: an extra conflicting owner is
	// not made safe by the presence of one valid-looking link.
	var count int
	err := runner.QueryRowContext(ctx, r.query(`SELECT COUNT(*) FROM source_rule_links WHERE point_id=$1`), member.PointID).Scan(&count)
	if err != nil {
		return err
	}
	if count != 1 {
		return draftSourceRevisionConflict("source rule ownership is ambiguous")
	}
	var ruleID, address string
	err = runner.QueryRowContext(ctx, r.query(`SELECT r.id,l.address FROM source_rule_links l
		JOIN source_rules r ON r.id=l.rule_id JOIN points p ON p.id=l.point_id
		WHERE l.point_id=$1 AND l.tag_id=$2 AND l.mapping_id=$3
		AND r.device_id=$4 AND p.device_id=r.device_id AND l.address=p.address`), member.PointID, member.TagID, mapping.ID, member.DeviceID).Scan(&ruleID, &address)
	if errors.Is(err, sql.ErrNoRows) {
		return draftSourceRevisionConflict("saved rule mapping relation changed")
	}
	if err != nil {
		return err
	}
	candidateID, err := sourcerule.MappingCandidateID(ruleID, address)
	if err != nil {
		return err
	}
	if candidateID != mapping.RuleCandidateID {
		return draftSourceRevisionConflict("saved rule mapping candidate changed")
	}
	return nil
}

func stoppedMappingMatchesSaved(mapping mappingRevisionSource, savedRevision string) (bool, error) {
	if mapping.Enabled || savedRevision == "" || mapping.BlockingReason != "" ||
		mapping.Status != string(schema.MappingStatusDraft) && mapping.Status != string(schema.MappingStatusActive) {
		return false, nil
	}
	// Stop may change only enabled and active -> draft. All other hashed fields
	// (including mapping ID, pipeline and candidate signatures) must still match.
	states := []mappingRevisionSource{mapping}
	mapping.Enabled = true
	states = append(states, mapping)
	if mapping.Status == string(schema.MappingStatusDraft) {
		mapping.Status = string(schema.MappingStatusActive)
		states = append(states, mapping)
	}
	for _, state := range states {
		revision, err := state.revision()
		if err != nil || revision == savedRevision {
			return err == nil, err
		}
	}
	return false, nil
}

func draftSourceRevisionConflict(reason string) error {
	return fmt.Errorf("%w: %w: %s", ErrWriteGroupSourceRevisionConflict, ErrSetupRevisionConflict, reason)
}

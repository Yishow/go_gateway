package workspace

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// WriteGroupMigrationSourceKindRecordingPlan marks a group whose provenance is
// a legacy recording plan.
const WriteGroupMigrationSourceKindRecordingPlan = "recording-plan"

// ErrRecordingPlanGroupUnresolved means a legacy plan does not map to exactly
// one canonical group. Callers must refuse rather than guess one.
var ErrRecordingPlanGroupUnresolved = errors.New("recording plan does not resolve to exactly one write group")

// ResolveRecordingPlanGroup returns the canonical group a legacy recording plan
// stands for, using only recorded migration provenance. No provenance, or more
// than one group claiming the plan, is an error: a plan is never matched by
// guessing the first group, a name or a table.
func (s *WriteGroupService) ResolveRecordingPlanGroup(ctx context.Context, planID string) (string, error) {
	planID = strings.TrimSpace(planID)
	if planID == "" {
		return "", fmt.Errorf("%w: plan id is required", ErrRecordingPlanGroupUnresolved)
	}
	listed, err := s.List(ctx)
	if err != nil {
		return "", err
	}
	var matches []string
	for _, group := range listed.Groups {
		if group == nil || group.Status == WriteGroupStatusDeleted ||
			strings.TrimSpace(group.Migration.SourceKind) != WriteGroupMigrationSourceKindRecordingPlan {
			continue
		}
		for _, id := range group.Migration.SourceIDs {
			if strings.TrimSpace(id) == planID {
				matches = append(matches, group.ID)
				break
			}
		}
	}
	if len(matches) != 1 {
		return "", fmt.Errorf("%w: %d candidates", ErrRecordingPlanGroupUnresolved, len(matches))
	}
	return matches[0], nil
}

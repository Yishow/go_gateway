package workspace

import (
	"context"
	"fmt"
)

type readinessWriteGroupService interface {
	List(context.Context) (*WriteGroupListResult, error)
	Readiness(context.Context, string) (*WriteGroupReadiness, error)
}

// WithWriteGroupReadiness adds canonical database groups to setup readiness.
func (s *Service) WithWriteGroupReadiness(groups readinessWriteGroupService) *Service {
	s.readinessWriteGroups = groups
	return s
}

type writeGroupReadinessCoverage struct {
	groupIDs map[string]bool
	members  map[writeGroupReadinessMemberKey]bool
	pointIDs map[string]bool
}

type writeGroupReadinessMemberKey struct {
	connectorID, pointID, tagID string
}

func (s *Service) canonicalWriteGroupReadiness(ctx context.Context, record *Record) (writeGroupReadinessCoverage, []ReadinessIssue, error) {
	coverage := writeGroupReadinessCoverage{
		groupIDs: map[string]bool{}, members: map[writeGroupReadinessMemberKey]bool{},
	}
	if s.readinessWriteGroups == nil {
		return coverage, nil, nil
	}
	listed, err := s.readinessWriteGroups.List(ctx)
	if err != nil {
		return coverage, nil, fmt.Errorf("list canonical groups for readiness: %w", err)
	}
	if listed == nil {
		return coverage, nil, ErrReadinessUnavailable
	}
	if listed.WorkspaceID != record.ID || listed.WorkspaceRevision != record.DatabaseSetupRevision {
		return coverage, nil, ErrSetupRevisionConflict
	}
	issues := make([]ReadinessIssue, 0)
	for _, group := range listed.Groups {
		if group == nil || group.WorkspaceID != record.ID {
			return coverage, nil, ErrReadinessUnavailable
		}
		coverage.groupIDs[group.ID] = true
		for _, member := range group.Members {
			coverage.members[writeGroupReadinessMemberKey{
				connectorID: group.Destination.ConnectorID, pointID: member.PointID, tagID: member.TagID,
			}] = true
		}
		if group.Status == WriteGroupStatusDeleted || group.Status == WriteGroupStatusDisabled {
			continue
		}
		readiness, err := s.readinessWriteGroups.Readiness(ctx, group.ID)
		if err != nil {
			return coverage, nil, fmt.Errorf("evaluate canonical group readiness: %w", err)
		}
		if readiness == nil {
			return coverage, nil, ErrReadinessUnavailable
		}
		if readiness.WorkspaceID != record.ID || readiness.WorkspaceRevision != record.DatabaseSetupRevision ||
			readiness.GroupID != group.ID || readiness.GroupRevision != group.Revision {
			return coverage, nil, ErrSetupRevisionConflict
		}
		issues = append(issues, readiness.Issues...)
		if group.AppliedRevision == "" {
			issues = append(issues, writeGroupReadinessIssue(
				"write-group-apply-required", "saved write group has not been applied", group.ID,
			))
		}
	}
	return coverage, issues, nil
}

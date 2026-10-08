package workspace

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ErrReadinessScopeInvalid reports a scope that is invalid for the persisted workspace.
var ErrReadinessScopeInvalid = errors.New("workspace readiness scope invalid")

type readinessScopeSelection struct {
	record    *Record
	deviceIDs []string
	groupIDs  []string
	groups    map[string]*WriteGroup
}

// ReadinessScope evaluates only the persisted device and group scope supplied by the caller.
func (s *Service) ReadinessScope(ctx context.Context, deviceIDs, groupIDs []string) (*ReadinessSummary, error) {
	_, summary, err := s.readinessScope(ctx, deviceIDs, groupIDs)
	return summary, err
}

func (s *Service) readinessScope(
	ctx context.Context,
	deviceIDs, groupIDs []string,
) (*readinessScopeSelection, *ReadinessSummary, error) {
	return s.readinessScopeWithApplyRequirement(ctx, deviceIDs, groupIDs, false)
}

func (s *Service) activationReadinessScope(
	ctx context.Context,
	deviceIDs, groupIDs []string,
) (*readinessScopeSelection, *ReadinessSummary, error) {
	return s.readinessScopeWithApplyRequirement(ctx, deviceIDs, groupIDs, true)
}

func (s *Service) readinessScopeWithApplyRequirement(
	ctx context.Context,
	deviceIDs, groupIDs []string,
	requireApplied bool,
) (*readinessScopeSelection, *ReadinessSummary, error) {
	selection, err := s.resolveReadinessScope(ctx, deviceIDs, groupIDs)
	if err != nil {
		return nil, nil, err
	}

	issues := make([]ReadinessIssue, 0, len(selection.deviceIDs)+len(selection.groupIDs))
	for _, deviceID := range selection.deviceIDs {
		readiness, err := s.readinessDevices.CheckReadiness(ctx, deviceID)
		if err != nil {
			return nil, nil, fmt.Errorf("evaluate scoped device readiness %s: %w", deviceID, err)
		}
		if readiness != nil && strings.TrimSpace(readiness.DeviceID) != "" && readiness.DeviceID != deviceID {
			return nil, nil, scopeRevisionError("device readiness belongs to a different persisted device")
		}
		if issue, ok := deviceReadinessIssue(deviceID, readiness); ok {
			issues = append(issues, issue)
		}
	}

	if len(selection.groupIDs) > 0 {
		groupIssues, err := s.readinessScopeGroupIssues(ctx, selection, requireApplied)
		if err != nil {
			return nil, nil, err
		}
		issues = append(issues, groupIssues...)
	}

	coverage := writeGroupReadinessCoverage{
		groupIDs: map[string]bool{}, members: map[writeGroupReadinessMemberKey]bool{},
	}
	if len(selection.groupIDs) > 0 {
		coverage.pointIDs = make(map[string]bool)
		for _, groupID := range selection.groupIDs {
			for _, member := range selection.groups[groupID].Members {
				if pointID := strings.TrimSpace(member.PointID); pointID != "" {
					coverage.pointIDs[pointID] = true
				}
			}
		}
	}

	// A scope without canonical groups is still allowed for Local Modbus Share.
	// Passing an empty connector deliberately keeps tag/mapping checks while
	// preventing the legacy database-target reader from becoming a prerequisite.
	downstreamIssues, err := s.downstreamReadinessIssues(ctx, selection.deviceIDs, "", coverage)
	if err != nil {
		return nil, nil, err
	}
	issues = append(issues, downstreamIssues...)

	return selection, summarizeScopedReadiness(issues), nil
}

func (s *Service) resolveReadinessScope(ctx context.Context, deviceIDs, groupIDs []string) (*readinessScopeSelection, error) {
	if s == nil || s.readinessDevices == nil {
		return nil, ErrReadinessUnavailable
	}
	record, err := s.GetOrCreate(ctx)
	if err != nil {
		return nil, err
	}
	normalizedDevices, err := normalizeReadinessScopeIDs(deviceIDs, "device")
	if err != nil {
		return nil, err
	}
	normalizedGroups, err := normalizeReadinessScopeIDs(groupIDs, "group")
	if err != nil {
		return nil, err
	}

	ownedDevices := make(map[string]bool, len(record.OrderedDeviceIDs))
	for _, deviceID := range record.OrderedDeviceIDs {
		ownedDevices[strings.TrimSpace(deviceID)] = true
	}
	selectedDevices := slices.Clone(normalizedDevices)
	selectedDeviceSet := make(map[string]bool, len(selectedDevices))
	for _, deviceID := range selectedDevices {
		if !ownedDevices[deviceID] {
			return nil, scopeInvalidError("device %s is not owned by the persisted workspace", deviceID)
		}
		selectedDeviceSet[deviceID] = true
	}

	selection := &readinessScopeSelection{
		record:    record,
		deviceIDs: selectedDevices,
		groupIDs:  normalizedGroups,
		groups:    make(map[string]*WriteGroup, len(normalizedGroups)),
	}
	if len(normalizedGroups) > 0 {
		if s.readinessWriteGroups == nil {
			return nil, ErrReadinessUnavailable
		}
		listed, err := s.readinessWriteGroups.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("list scoped canonical groups: %w", err)
		}
		if listed == nil {
			return nil, ErrReadinessUnavailable
		}
		if listed.WorkspaceID != record.ID || listed.WorkspaceRevision != record.DatabaseSetupRevision {
			return nil, scopeRevisionError("canonical group snapshot is stale")
		}
		available := make(map[string]*WriteGroup, len(listed.Groups))
		for _, group := range listed.Groups {
			if group == nil || strings.TrimSpace(group.ID) == "" {
				continue
			}
			available[group.ID] = group
		}
		for _, groupID := range normalizedGroups {
			group, ok := available[groupID]
			if !ok || group.WorkspaceID != record.ID {
				return nil, scopeInvalidError("group %s is not owned by the persisted workspace", groupID)
			}
			if group.Status == WriteGroupStatusDeleted || group.Status == WriteGroupStatusDisabled {
				return nil, scopeInvalidError("group %s is not active", groupID)
			}
			if strings.TrimSpace(group.Revision) == "" {
				return nil, scopeRevisionError("group %s has no persisted revision", groupID)
			}
			selection.groups[groupID] = group
			matchedSelectedDevice := false
			for _, member := range group.Members {
				memberDeviceID := strings.TrimSpace(member.DeviceID)
				if memberDeviceID == "" || !ownedDevices[memberDeviceID] {
					return nil, scopeInvalidError("group %s contains an unowned device", groupID)
				}
				if len(normalizedDevices) > 0 {
					if !selectedDeviceSet[memberDeviceID] {
						return nil, scopeInvalidError("group %s is outside the selected device scope", groupID)
					}
					matchedSelectedDevice = true
					continue
				}
				if !selectedDeviceSet[memberDeviceID] {
					selectedDevices = append(selectedDevices, memberDeviceID)
					selectedDeviceSet[memberDeviceID] = true
				}
				matchedSelectedDevice = true
			}
			if len(normalizedDevices) > 0 && !matchedSelectedDevice {
				return nil, scopeInvalidError("group %s has no selected member device", groupID)
			}
		}
	}
	selection.deviceIDs = selectedDevices
	if len(selection.deviceIDs) == 0 {
		return nil, scopeInvalidError("scope must select at least one persisted device")
	}
	for _, deviceID := range selection.deviceIDs {
		deviceRecord, err := s.readinessDevices.GetByID(ctx, deviceID)
		if err != nil {
			return nil, scopeInvalidError("read persisted device %s: %v", deviceID, err)
		}
		if deviceRecord == nil || strings.TrimSpace(deviceRecord.ID) != deviceID {
			return nil, scopeInvalidError("device %s is no longer persisted", deviceID)
		}
	}
	return selection, nil
}

func (s *Service) readinessScopeGroupIssues(ctx context.Context, selection *readinessScopeSelection, requireApplied bool) ([]ReadinessIssue, error) {
	issues := make([]ReadinessIssue, 0, len(selection.groupIDs))
	for _, groupID := range selection.groupIDs {
		group := selection.groups[groupID]
		readiness, err := s.readinessWriteGroups.Readiness(ctx, groupID)
		if err != nil {
			return nil, fmt.Errorf("evaluate scoped group readiness %s: %w", groupID, err)
		}
		if readiness == nil {
			return nil, ErrReadinessUnavailable
		}
		if readiness.WorkspaceID != selection.record.ID ||
			readiness.WorkspaceRevision != selection.record.DatabaseSetupRevision ||
			readiness.GroupID != groupID || readiness.GroupRevision != group.Revision ||
			readiness.AppliedRevision != group.AppliedRevision {
			return nil, scopeRevisionError("group %s readiness snapshot is stale", groupID)
		}
		issues = append(issues, readiness.Issues...)
		if !readiness.Ready && len(readiness.Issues) == 0 {
			issues = append(issues, writeGroupReadinessIssue(
				"write-group-not-ready", "selected write group is not ready", groupID,
			))
		}
		if requireApplied && group.AppliedRevision == "" && !hasScopedReadinessIssue(issues, "write-group-apply-required", groupID) {
			issues = append(issues, writeGroupReadinessIssue(
				"write-group-apply-required", "saved write group has not been applied", groupID,
			))
		}
	}
	return issues, nil
}

func normalizeReadinessScopeIDs(ids []string, kind string) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	result := make([]string, 0, len(ids))
	seen := make(map[string]bool, len(ids))
	for _, rawID := range ids {
		id := strings.TrimSpace(rawID)
		if id == "" {
			return nil, scopeInvalidError("%s scope contains an empty id", kind)
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		result = append(result, id)
	}
	return result, nil
}

func summarizeScopedReadiness(issues []ReadinessIssue) *ReadinessSummary {
	summary := &ReadinessSummary{Issues: issues}
	for _, issue := range issues {
		switch issue.Severity {
		case ReadinessSeverityBlocking:
			summary.BlockingCount++
		case ReadinessSeverityWarning:
			summary.WarningCount++
		}
	}
	summary.Ready = summary.BlockingCount == 0
	return summary
}

func hasScopedReadinessIssue(issues []ReadinessIssue, code, scope string) bool {
	return slices.ContainsFunc(issues, func(issue ReadinessIssue) bool {
		return issue.Code == code && issue.Scope == scope
	})
}

func scopeInvalidError(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrReadinessScopeInvalid, fmt.Sprintf(format, args...))
}

func scopeRevisionError(format string, args ...any) error {
	return fmt.Errorf("%w: %w: %s", ErrReadinessScopeInvalid, ErrSetupRevisionConflict, fmt.Sprintf(format, args...))
}

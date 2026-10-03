package modbusshare

import (
	"context"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// CheckDesiredProjectionForScope verifies selected mappings against the
// already-proven runtime projection without changing Share state. The desired
// slice is scoped by the caller; mappings outside it are deliberately ignored
// so a selected start cannot rewrite or remove another device's projection.
func (s *Service) CheckDesiredProjectionForScope(ctx context.Context, workspaceID string, desired []DesiredMapping) error {
	if s == nil {
		return NewError(ErrCodeProjectionRequired, "modbus Share service is unavailable", true)
	}
	workspaceID = strings.TrimSpace(workspaceID)
	if workspaceID == "" {
		return NewError(ErrCodeWorkspaceScope, "workspace id is required", false)
	}

	current, err := s.ListMappingsForWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}
	byTag := make(map[string]TagMirrorMapping, len(current))
	for _, mapping := range current {
		byTag[mapping.TagID] = mapping
	}

	seen := make(map[string]struct{}, len(desired))
	for _, candidate := range desired {
		if candidate.WorkspaceID != "" && candidate.WorkspaceID != workspaceID {
			return scopedProjectionMismatch("desired Share mapping belongs to another workspace")
		}
		if strings.TrimSpace(candidate.TagID) == "" ||
			strings.TrimSpace(candidate.SourceRuleID) == "" ||
			strings.TrimSpace(candidate.SourceRuleRevision) == "" {
			return scopedProjectionMismatch("desired Share mapping has incomplete persisted ownership")
		}
		if _, duplicate := seen[candidate.TagID]; duplicate {
			return scopedProjectionMismatch("desired Share projection contains duplicate tag %q", candidate.TagID)
		}
		seen[candidate.TagID] = struct{}{}

		mapping, ok := byTag[candidate.TagID]
		if !ok {
			return scopedProjectionMismatch("selected Share mapping %q is not aligned", candidate.TagID)
		}
		if mappingProjectionChanged(mapping, desiredProjectionMapping(workspaceID, candidate)) {
			return scopedProjectionMismatch("selected Share mapping %q is not aligned", candidate.TagID)
		}
	}
	return nil
}

// CheckDesiredProjectionForScopeWithSourceRules verifies the exact projection
// for selected devices while retaining valid mappings owned by other devices.
func (s *Service) CheckDesiredProjectionForScopeWithSourceRules(ctx context.Context, workspaceID string, selectedDeviceIDs []string, desired []DesiredMapping, sourceRules []*schema.SourceRule) error {
	if s == nil {
		return NewError(ErrCodeProjectionRequired, "modbus Share service is unavailable", true)
	}
	workspaceID = strings.TrimSpace(workspaceID)
	if workspaceID == "" {
		return NewError(ErrCodeWorkspaceScope, "workspace id is required", false)
	}
	selected := make(map[string]struct{}, len(selectedDeviceIDs))
	for _, rawID := range selectedDeviceIDs {
		deviceID := strings.TrimSpace(rawID)
		if deviceID != "" {
			selected[deviceID] = struct{}{}
		}
	}
	if len(selected) == 0 {
		return scopedProjectionMismatch("selected Share scope has no device identity")
	}
	ruleDevices := make(map[string]string, len(sourceRules))
	for _, rule := range sourceRules {
		if rule == nil || strings.TrimSpace(rule.ID) == "" || strings.TrimSpace(rule.DeviceID) == "" {
			continue
		}
		ruleDevices[rule.ID] = strings.TrimSpace(rule.DeviceID)
	}

	s.mu.RLock()
	rawMappings := make([]TagMirrorMapping, 0, len(s.mappings))
	desiredOwnership := s.desiredOwnership
	ownership := s.ownership
	for _, mapping := range s.mappings {
		if mapping.WorkspaceID == workspaceID {
			rawMappings = append(rawMappings, mapping)
		}
	}
	s.mu.RUnlock()

	selectedCurrent := make([]TagMirrorMapping, 0, len(rawMappings))
	for _, mapping := range rawMappings {
		ruleDeviceID, ok := ruleDevices[mapping.SourceRuleID]
		if !ok {
			return scopedProjectionMismatch("Share mapping %q has no persisted source-rule ownership", mapping.TagID)
		}
		if _, ok := selected[ruleDeviceID]; !ok {
			continue
		}
		if err := validateScopedMappingOwnership(ctx, mapping, workspaceID, desiredOwnership, ownership); err != nil {
			return err
		}
		selectedCurrent = append(selectedCurrent, mapping)
	}

	byTag := make(map[string]TagMirrorMapping, len(selectedCurrent))
	for _, mapping := range selectedCurrent {
		byTag[mapping.TagID] = mapping
	}
	seen := make(map[string]struct{}, len(desired))
	for _, candidate := range desired {
		if candidate.WorkspaceID != "" && candidate.WorkspaceID != workspaceID {
			return scopedProjectionMismatch("desired Share mapping belongs to another workspace")
		}
		if strings.TrimSpace(candidate.TagID) == "" ||
			strings.TrimSpace(candidate.SourceRuleID) == "" ||
			strings.TrimSpace(candidate.SourceRuleRevision) == "" {
			return scopedProjectionMismatch("desired Share mapping has incomplete persisted ownership")
		}
		if _, duplicate := seen[candidate.TagID]; duplicate {
			return scopedProjectionMismatch("desired Share projection contains duplicate tag %q", candidate.TagID)
		}
		seen[candidate.TagID] = struct{}{}

		mapping, ok := byTag[candidate.TagID]
		if !ok {
			return scopedProjectionMismatch("selected Share mapping %q is not aligned", candidate.TagID)
		}
		if mappingProjectionChanged(mapping, desiredProjectionMapping(workspaceID, candidate)) {
			return scopedProjectionMismatch("selected Share mapping %q is not aligned", candidate.TagID)
		}
	}
	if len(seen) != len(byTag) {
		return scopedProjectionMismatch("selected Share projection contains stale or extra mappings")
	}
	return nil
}

func validateScopedMappingOwnership(ctx context.Context, mapping TagMirrorMapping, workspaceID string, desiredChecker DesiredOwnershipChecker, checker OwnershipChecker) error {
	if strings.TrimSpace(mapping.SourceRuleID) == "" || strings.TrimSpace(mapping.SourceRuleRevision) == "" {
		return scopedProjectionMismatch("Share mapping %q has incomplete persisted ownership", mapping.TagID)
	}
	desired := desiredMappingFromMirror(mapping)
	if desiredChecker != nil {
		if err := desiredChecker(ctx, desired); err != nil {
			return scopedProjectionMismatch("Share mapping %q ownership is not proven", mapping.TagID)
		}
		return nil
	}
	if checker != nil && checker(ctx, workspaceID, mapping.TagID) {
		return nil
	}
	return NewError(ErrCodeWorkspaceScope, "durable workspace ownership is not configured", true)
}

func desiredProjectionMapping(workspaceID string, candidate DesiredMapping) TagMirrorMapping {
	span := candidate.SpanRegisters
	if span <= 0 {
		span = DataTypeSpan(candidate.DataType)
	}
	return TagMirrorMapping{
		WorkspaceID:        workspaceID,
		SourceRuleID:       candidate.SourceRuleID,
		SourceRuleRevision: candidate.SourceRuleRevision,
		TagID:              candidate.TagID,
		MappingID:          candidate.MappingID,
		Register:           candidate.ZeroBasedRegister,
		ShareStartRegister: candidate.ShareStartRegister,
		ZeroBasedRegister:  candidate.ZeroBasedRegister,
		SpanRegisters:      span,
		StrideRegisters:    candidate.StrideRegisters,
		CapacityRegisters:  candidate.CapacityRegisters,
		DataType:           candidate.DataType,
		TagKey:             candidate.TagKey,
		DisplayName:        candidate.DisplayName,
	}
}

func scopedProjectionMismatch(format string, args ...any) error {
	return NewError(
		ErrCodeProjectionRequired,
		fmt.Sprintf(format, args...),
		true,
	)
}

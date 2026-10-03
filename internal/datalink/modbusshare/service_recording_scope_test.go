package modbusshare

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go-gateway/internal/datalink/schema"
)

func TestService_CheckDesiredProjectionForScopeKeepsUnselectedMapping(t *testing.T) {
	ctx := t.Context()
	service := NewService(nil, 4096)
	selected := recordingScopeMirror("ws-1", "tag-a", "rule-a", "rev-a", 10)
	unselected := recordingScopeMirror("ws-1", "tag-b", "rule-b", "rev-b", 20)
	service.ReplaceMappings(map[string]TagMirrorMapping{
		selected.TagID:   selected,
		unselected.TagID: unselected,
	})
	service.SetDesiredMappingOwnershipChecker(func(context.Context, DesiredMapping) error { return nil })

	before := service.ListMappings()
	memoryBefore := service.MemorySnapshot()
	require.NoError(t, service.CheckDesiredProjectionForScope(ctx, "ws-1", []DesiredMapping{
		recordingScopeDesired(selected),
	}))

	assertMappingSetEqual(t, before, service.ListMappings())
	require.Equal(t, memoryBefore, service.MemorySnapshot())
}

func TestService_CheckDesiredProjectionForScopeRejectsSelectedMismatch(t *testing.T) {
	service := NewService(nil, 4096)
	current := recordingScopeMirror("ws-1", "tag-a", "rule-a", "rev-a", 10)
	service.ReplaceMappings(map[string]TagMirrorMapping{current.TagID: current})
	service.SetDesiredMappingOwnershipChecker(func(context.Context, DesiredMapping) error { return nil })

	before := service.ListMappings()
	desired := recordingScopeDesired(current)
	desired.ZeroBasedRegister = 11
	desired.ShareStartRegister = ZeroBasedToHuman(desired.ZeroBasedRegister)
	require.Error(t, service.CheckDesiredProjectionForScope(t.Context(), "ws-1", []DesiredMapping{desired}))

	assertMappingSetEqual(t, before, service.ListMappings())
}

func TestService_CheckDesiredProjectionForScopeRejectsSourceOwnershipDrift(t *testing.T) {
	service := NewService(nil, 4096)
	current := recordingScopeMirror("ws-1", "tag-a", "rule-a", "rev-a", 10)
	service.ReplaceMappings(map[string]TagMirrorMapping{current.TagID: current})
	service.SetDesiredMappingOwnershipChecker(func(context.Context, DesiredMapping) error { return nil })

	desired := recordingScopeDesired(current)
	desired.SourceRuleRevision = "rev-b"
	require.Error(t, service.CheckDesiredProjectionForScope(t.Context(), "ws-1", []DesiredMapping{desired}))
}

func TestService_CheckDesiredProjectionForScopeRejectsUnprovenRuntimeMapping(t *testing.T) {
	service := NewService(nil, 4096)
	current := recordingScopeMirror("ws-2", "tag-a", "rule-a", "rev-a", 10)
	service.ReplaceMappings(map[string]TagMirrorMapping{current.TagID: current})

	err := service.CheckDesiredProjectionForScope(t.Context(), "ws-2", []DesiredMapping{recordingScopeDesired(current)})
	require.Error(t, err)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeWorkspaceScope, shareErr.Code)
}

func TestService_CheckDesiredProjectionForScopeWithSourceRulesRejectsStaleSelectedMapping(t *testing.T) {
	service := NewService(nil, 4096)
	selected := recordingScopeMirror("ws-1", "tag-a", "rule-a", "rev-a", 10)
	unselected := recordingScopeMirror("ws-1", "tag-b", "rule-b", "rev-b", 20)
	service.ReplaceMappings(map[string]TagMirrorMapping{selected.TagID: selected, unselected.TagID: unselected})
	service.SetDesiredMappingOwnershipChecker(func(context.Context, DesiredMapping) error { return nil })

	err := service.CheckDesiredProjectionForScopeWithSourceRules(
		t.Context(), "ws-1", []string{"device-a"}, nil,
		[]*schema.SourceRule{{ID: "rule-a", DeviceID: "device-a"}, {ID: "rule-b", DeviceID: "device-b"}},
	)
	require.Error(t, err)
}

func TestService_CheckDesiredProjectionForScopeWithSourceRulesPreservesUnselectedMapping(t *testing.T) {
	service := NewService(nil, 4096)
	selected := recordingScopeMirror("ws-1", "tag-a", "rule-a", "rev-a", 10)
	unselected := recordingScopeMirror("ws-1", "tag-b", "rule-b", "rev-b", 20)
	service.ReplaceMappings(map[string]TagMirrorMapping{selected.TagID: selected, unselected.TagID: unselected})
	service.SetDesiredMappingOwnershipChecker(func(context.Context, DesiredMapping) error { return nil })

	require.NoError(t, service.CheckDesiredProjectionForScopeWithSourceRules(
		t.Context(), "ws-1", []string{"device-a"}, []DesiredMapping{recordingScopeDesired(selected)},
		[]*schema.SourceRule{{ID: "rule-a", DeviceID: "device-a"}, {ID: "rule-b", DeviceID: "device-b"}},
	))
}

func recordingScopeMirror(workspaceID, tagID, ruleID, ruleRevision string, register uint16) TagMirrorMapping {
	return TagMirrorMapping{
		WorkspaceID: workspaceID, SourceRuleID: ruleID, SourceRuleRevision: ruleRevision,
		TagID: tagID, MappingID: "mapping-" + tagID, Register: register,
		ShareStartRegister: ZeroBasedToHuman(register), ZeroBasedRegister: register,
		SpanRegisters: 1, StrideRegisters: 1, CapacityRegisters: 4096,
		DataType: schema.DataTypeInt16, TagKey: "key-" + tagID, DisplayName: "Display " + tagID,
	}
}

func recordingScopeDesired(mapping TagMirrorMapping) DesiredMapping {
	return DesiredMapping{
		WorkspaceID: mapping.WorkspaceID, SourceRuleID: mapping.SourceRuleID, SourceRuleRevision: mapping.SourceRuleRevision,
		TagID: mapping.TagID, MappingID: mapping.MappingID, DataType: mapping.DataType,
		ShareStartRegister: mapping.ShareStartRegister, ZeroBasedRegister: mapping.ZeroBasedRegister,
		SpanRegisters: mapping.SpanRegisters, StrideRegisters: mapping.StrideRegisters,
		CapacityRegisters: mapping.CapacityRegisters, TagKey: mapping.TagKey, DisplayName: mapping.DisplayName,
	}
}

func assertMappingSetEqual(t *testing.T, expected, actual []TagMirrorMapping) {
	t.Helper()
	require.ElementsMatch(t, expected, actual)
}

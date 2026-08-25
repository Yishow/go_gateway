package modbusshare

import (
	"context"
	"fmt"
)

// ClearSpan clears complete holding-register words for an obsolete mapping.
func (s *Service) ClearSpan(start uint16, count int) error {
	if count <= 0 {
		return nil
	}
	return s.bank.WriteSlice(int(start)*2, make([]byte, count*2))
}

// MemorySnapshot returns a copy used to restore an atomic projection on error.
func (s *Service) MemorySnapshot() []byte { return s.bank.Dump() }

// RestoreMemory restores a previously captured projection snapshot.
func (s *Service) RestoreMemory(snapshot []byte) error {
	s.mu.RLock()
	hook := s.restoreMemoryHook
	s.mu.RUnlock()
	if hook != nil {
		return hook(snapshot)
	}
	if len(snapshot) != s.bank.Size() {
		return fmt.Errorf("invalid memory snapshot size %d", len(snapshot))
	}
	return s.bank.WriteSlice(0, snapshot)
}

// FailClosed stops the listener and makes the Share hydration gate fail
// closed. It is used when rollback cannot establish a trustworthy final state.
func (s *Service) FailClosed(state HydrationState) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.server != nil {
		if err := s.server.Stop(); err != nil {
			s.hydrationErr = fmt.Errorf("stop modbus share listener: %w", err)
		}
	}
	state.State = HydrationStateFailed
	state.Readiness = false
	s.hydration = state
	s.state = HydrationStateFailed
}

// SetRestoreMemoryHook is a narrow failure-injection seam for rollback tests.
func (s *Service) SetRestoreMemoryHook(hook func([]byte) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.restoreMemoryHook = hook
}

// ListMappings lists all mappings.
func (s *Service) ListMappings() []TagMirrorMapping {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]TagMirrorMapping, 0, len(s.mappings))
	for _, m := range s.mappings {
		out = append(out, m)
	}
	return out
}

// ListMappingsByWorkspace returns the in-memory projection for one workspace.
// Ownership validation remains the responsibility of API-facing callers.
func (s *Service) ListMappingsByWorkspace(workspaceID string) []TagMirrorMapping {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]TagMirrorMapping, 0)
	for _, mapping := range s.mappings {
		if mapping.WorkspaceID == workspaceID {
			out = append(out, mapping)
		}
	}
	return out
}

// ListMappingsForWorkspace returns only mappings proven to belong to the
// hydrated workspace. Unknown/orphan runtime rows are withheld rather than
// leaking another workspace's state to the caller.
func (s *Service) ListMappingsForWorkspace(ctx context.Context, workspaceID string) ([]TagMirrorMapping, error) {
	if workspaceID == "" {
		return nil, NewError(ErrCodeWorkspaceScope, "workspace id is required", false)
	}
	s.mu.RLock()
	mappings := make([]TagMirrorMapping, 0, len(s.mappings))
	checker := s.ownership
	desiredChecker := s.desiredOwnership
	for _, mapping := range s.mappings {
		if mapping.WorkspaceID == workspaceID && mapping.SourceRuleID != "" && mapping.SourceRuleRevision != "" {
			mappings = append(mappings, mapping)
		}
	}
	s.mu.RUnlock()
	if checker == nil && desiredChecker == nil {
		return nil, NewError(ErrCodeWorkspaceScope, "durable workspace ownership is not configured", true)
	}
	owned := mappings[:0]
	for _, mapping := range mappings {
		desired := DesiredMapping{WorkspaceID: mapping.WorkspaceID, SourceRuleID: mapping.SourceRuleID, SourceRuleRevision: mapping.SourceRuleRevision, TagID: mapping.TagID, MappingID: mapping.MappingID, DataType: mapping.DataType, ShareStartRegister: mapping.ShareStartRegister, ZeroBasedRegister: mapping.ZeroBasedRegister, SpanRegisters: mapping.SpanRegisters, StrideRegisters: mapping.StrideRegisters, CapacityRegisters: mapping.CapacityRegisters, TagKey: mapping.TagKey, DisplayName: mapping.DisplayName}
		if (desiredChecker != nil && desiredChecker(ctx, desired) == nil) || (desiredChecker == nil && checker(ctx, workspaceID, mapping.TagID)) {
			owned = append(owned, mapping)
		}
	}
	return owned, nil
}

// HasMapping checks whether a tag mirror mapping exists.
func (s *Service) HasMapping(tagID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.mappings[tagID]
	return ok
}

// ReplaceMappings replaces the internal mappings map atomically and returns the previous mappings map.
func (s *Service) ReplaceMappings(newMappings map[string]TagMirrorMapping) map[string]TagMirrorMapping {
	s.mu.Lock()
	defer s.mu.Unlock()
	old := s.mappings
	s.mappings = make(map[string]TagMirrorMapping, len(newMappings))
	for k, v := range newMappings {
		s.mappings[k] = v
	}
	return old
}

// ReplaceMappingsForWorkspace swaps only one workspace's projection and
// preserves mappings owned by every other workspace.
func (s *Service) ReplaceMappingsForWorkspace(workspaceID string, newMappings map[string]TagMirrorMapping) map[string]TagMirrorMapping {
	s.mu.Lock()
	defer s.mu.Unlock()
	old := s.mappings
	next := make(map[string]TagMirrorMapping, len(old)+len(newMappings))
	for tagID, mapping := range old {
		if mapping.WorkspaceID != workspaceID {
			next[tagID] = mapping
		}
	}
	for tagID, mapping := range newMappings {
		next[tagID] = mapping
	}
	s.mappings = next
	return old
}

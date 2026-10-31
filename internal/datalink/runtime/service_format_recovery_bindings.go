package runtime

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"slices"
)

// These opaque values bind one attempt to the existing target mapping install.
// They are ephemeral values, never a second registry or externally injectable DTO.
type FormatRecoveryBindingProof struct {
	pointID           string
	exists            bool
	installed, values []mappingBinding
}
type PreparedFormatRecoveryBindings struct {
	pointID string
	profile connector.ReadProfile
	values  []mappingBinding
}
type FormatRecoveryBindingReceipt struct {
	pointID   string
	lease     uint64
	installed []mappingBinding
}

func (s *Service) CaptureFormatRecoveryBindings(pointID string) (FormatRecoveryBindingProof, bool) {
	if s == nil || pointID == "" {
		return FormatRecoveryBindingProof{}, false
	}
	s.mappingMu.RLock()
	defer s.mappingMu.RUnlock()
	values, exists := s.mappingIndex[pointID]
	return FormatRecoveryBindingProof{pointID: pointID, exists: exists, installed: values, values: slices.Clone(values)}, true
}

// PrepareFormatRecoveryBindings is pure: records come from the caller's one SQL
// snapshot, never a provider lookup while preparing or publishing runtime state.
func PrepareFormatRecoveryBindings(workspaceID string, p schema.Point, profile connector.ReadProfile, device schema.Device, mappings []*schema.Mapping, tags map[string]*schema.Tag) (PreparedFormatRecoveryBindings, error) {
	if workspaceID == "" || connector.ValidatePointReadProfile(&p, profile) != nil || device.ID != p.DeviceID || !profile.MatchesConnection(device.Protocol, device.ConnectionConfig) {
		return PreparedFormatRecoveryBindings{}, ErrTypedSampleRejected
	}
	out := PreparedFormatRecoveryBindings{pointID: p.ID, profile: profile, values: make([]mappingBinding, 0, len(mappings))}
	ids := map[string]bool{}
	for _, m := range mappings {
		if m == nil || m.ID == "" || m.PointID != p.ID || ids[m.ID] {
			return PreparedFormatRecoveryBindings{}, ErrTypedSampleRejected
		}
		ids[m.ID] = true
		if !m.Enabled {
			continue
		}
		binding, err := mappingBindingForRecords(workspaceID, m, &device, &p, tags[m.TagID], profile)
		if err != nil {
			return PreparedFormatRecoveryBindings{}, err
		}
		out.values = append(out.values, binding)
	}
	return out, nil
}

func sameRecoveryBindings(a, b []mappingBinding) bool {
	return slices.Equal(a, b)
}
func sameRecoveryBindingInstall(current []mappingBinding, proof FormatRecoveryBindingProof, exists bool) bool {
	if exists != proof.exists || len(current) != len(proof.installed) || !sameRecoveryBindings(current, proof.values) {
		return false
	}
	return len(current) == 0 || &current[0] == &proof.installed[0]
}

// ActivateFormatRecoveryBindings publishes target bindings before making the
// owned staged lease pollable. Lock order: mappingMu → pointMetaMu → scheduler.
// No database/provider call occurs under these locks. Ingestor captures binding
// and lease acceptance under mappingMu's read side before doing any SQL I/O.
func (s *Service) ActivateFormatRecoveryBindings(id string, lease uint64, profile connector.ReadProfile, proof FormatRecoveryBindingProof, prepared PreparedFormatRecoveryBindings) (FormatRecoveryBindingReceipt, bool) {
	if s == nil || s.scheduler == nil || lease == 0 || proof.pointID != id || prepared.pointID != id || prepared.profile != profile {
		return FormatRecoveryBindingReceipt{}, false
	}
	s.mappingMu.Lock()
	defer s.mappingMu.Unlock()
	current, exists := s.mappingIndex[id]
	if !sameRecoveryBindingInstall(current, proof, exists) {
		return FormatRecoveryBindingReceipt{}, false
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	next := slices.Clone(prepared.values)
	s.mappingIndex[id] = next
	if !s.scheduler.ActivatePointFormatRecovery(id, lease, profile) {
		if exists {
			s.mappingIndex[id] = current
		} else {
			delete(s.mappingIndex, id)
		}
		return FormatRecoveryBindingReceipt{}, false
	}
	return FormatRecoveryBindingReceipt{pointID: id, lease: lease, installed: next}, true
}

func (s *Service) MatchesFormatRecoveryBindings(prepared PreparedFormatRecoveryBindings) bool {
	if s == nil || prepared.pointID == "" {
		return false
	}
	s.mappingMu.RLock()
	defer s.mappingMu.RUnlock()
	return sameRecoveryBindings(s.mappingIndex[prepared.pointID], prepared.values)
}

// SuspendFormatRecovery removes only this attempt's published binding and lease.
// A concurrent foreign binding is retained even if our own lease still exists.
func (s *Service) SuspendFormatRecovery(id string, lease uint64, receipt FormatRecoveryBindingReceipt) bool {
	if s == nil || s.scheduler == nil || lease == 0 {
		return false
	}
	s.mappingMu.Lock()
	defer s.mappingMu.Unlock()
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	if !s.scheduler.SuspendPointLease(id, lease) {
		return false
	}
	if receipt.pointID == id && receipt.lease == lease {
		current := s.mappingIndex[id]
		if len(current) == len(receipt.installed) && (len(current) == 0 || &current[0] == &receipt.installed[0]) {
			delete(s.mappingIndex, id)
		}
	}
	return true
}

// MatchesFormatRecoveryBoundLease atomically checks the replay's binding and slot.
func (s *Service) MatchesFormatRecoveryBoundLease(id string, lease uint64, prepared PreparedFormatRecoveryBindings) bool {
	if s == nil || s.scheduler == nil || id != prepared.pointID {
		return false
	}
	s.mappingMu.RLock()
	defer s.mappingMu.RUnlock()
	return sameRecoveryBindings(s.mappingIndex[id], prepared.values) && s.scheduler.MatchesFormatRecoveryLease(id, lease, prepared.profile)
}

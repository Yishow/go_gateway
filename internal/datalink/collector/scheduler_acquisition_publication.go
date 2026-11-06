package collector

import (
	"maps"
	"slices"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

// AcquisitionGeometryDefinition is derived source data, never an authorization
// flag. The runtime must validate its original opaque SQL graph under the same
// writer reservation as PublishAcquisitionGeometry. No HTTP/config input may
// call this publication seam or manufacture installation receipts.
type AcquisitionGeometryDefinition struct {
	WorkspaceID    string
	SourceRevision string
	Revision       string
	Excluded       []optimizer.NativeRange
	Rules          map[string]string
}

// AcquisitionGeometryUpdate journals actual successful install/removal receipts
// for one synchronous mutation. It must not be shared across goroutines.
type AcquisitionGeometryUpdate struct {
	expectedPublication *acquisitionGeometry
	expectedEpoch       uint64
	requireEpoch        bool
	owner               *Scheduler
	mutation            *acquisitionGeometryMutation
	candidate           *acquisitionGeometry
	scopeRule           string
	autoInterval        bool
	sealed              bool
	deviceApplied       bool
}

// PrepareAcquisitionGeometry copies the original full device slot proof. It
// does not read current slots to bless changes performed by another operation.
// Empty scopeRule is reserved for the existing full-device source install.
func (s *Scheduler) PrepareAcquisitionGeometry(mutation AcquisitionGeometryMutation, device schema.Device, definition AcquisitionGeometryDefinition, scopeRule string) *AcquisitionGeometryUpdate {
	token := mutation.token
	if s == nil || token == nil || token.owner != s || token.deviceID != device.ID {
		return nil
	}
	// A scoped mutation cannot promote an unverified foreign source set. The
	// ordinary owned install remains legal, but complete geometry stays unknown.
	if scopeRule != "" && !token.before.authorized {
		definition = AcquisitionGeometryDefinition{}
	}
	cfg := deviceConfiguration(&device)
	return &AcquisitionGeometryUpdate{expectedPublication: token.marker, owner: s, mutation: token, scopeRule: scopeRule, autoInterval: device.AcquisitionAutoInterval, candidate: &acquisitionGeometry{workspaceID: definition.WorkspaceID, sourceRevision: definition.SourceRevision, revision: definition.Revision, excluded: slices.Clone(definition.Excluded), rules: maps.Clone(definition.Rules), device: cfg, slots: maps.Clone(token.before.slots)}}
}

// Installed records the exact successful return lease with its supplied frozen
// Point/profile. A later Capture/Get result is not an acceptable substitute.
func (u *AcquisitionGeometryUpdate) Installed(point *schema.Point, profile connector.ReadProfile, lease uint64) bool {
	if u == nil || u.sealed || point == nil || point.DeviceID != u.candidate.device.ID || lease == 0 || connector.ValidatePointReadProfile(point, profile) != nil {
		return false
	}
	u.candidate.slots[point.ID] = installedPointInfo(point, profile, u.candidate.device, true, lease, u.owner.installationScope)
	return true
}

// Removed records a successful owned removal, not inferred absence.
func (u *AcquisitionGeometryUpdate) Removed(pointID string) bool {
	if u == nil || u.sealed || pointID == "" {
		return false
	}
	delete(u.candidate.slots, pointID)
	return true
}

// AcquisitionGeometryPublication is immutable prepared data, not SQL authority.
// Only the original runtime SQL acceptance callback may consume it.
type AcquisitionGeometryPublication struct{ update *AcquisitionGeometryUpdate }

// Seal freezes this operation's receipt journal before SQL acceptance begins.
func (u *AcquisitionGeometryUpdate) Seal() AcquisitionGeometryPublication {
	if u == nil || u.sealed {
		return AcquisitionGeometryPublication{}
	}
	u.sealed = true
	// A caller may have value-copied the builder before Seal. Freeze storage,
	// not merely its flag, so that copy cannot mutate published slot evidence.
	frozen := *u
	candidate := *u.candidate
	candidate.slots = maps.Clone(u.candidate.slots)
	candidate.rules = maps.Clone(u.candidate.rules)
	candidate.excluded = slices.Clone(u.candidate.excluded)
	frozen.candidate = &candidate
	return AcquisitionGeometryPublication{update: &frozen}
}

// PublishAcquisitionGeometry atomically publishes derived geometry with the
// caller's already validated mapping/scope assignments. It uses TryLock because
// its caller holds a SQL writer reservation. The callback must not fail, block,
// acquire locks, perform I/O, invoke providers, or call scheduler methods.
func (s *Scheduler) PublishAcquisitionGeometry(publication AcquisitionGeometryPublication, assign func()) (AcquisitionGeometryReceipt, bool) {
	var receipt [1]AcquisitionGeometryReceipt
	ok := s.PublishAcquisitionGeometries([]AcquisitionGeometryPublication{publication}, receipt[:], assign)
	return receipt[0], ok
}

// PublishAcquisitionGeometries validates the entire prepared set under one
// TryLock before any assignment. The caller preallocates receipts outside its
// SQL callback. Blank semantic revisions stay unknown; they never grant gaps.
func (s *Scheduler) PublishAcquisitionGeometries(publications []AcquisitionGeometryPublication, receipts []AcquisitionGeometryReceipt, assign func()) bool {
	if s == nil || assign == nil || len(publications) == 0 || len(receipts) != len(publications) || !s.mu.TryLock() {
		return false
	}
	defer s.mu.Unlock()
	for i, publication := range publications {
		u := publication.update
		if u == nil || u.owner != s || !u.sealed || !s.acquisitionPublicationCurrentLocked(u) {
			return false
		}
		for _, previous := range publications[:i] {
			if previous.update.mutation.deviceID == u.mutation.deviceID {
				return false
			}
		}
	}
	assign()
	for i, publication := range publications {
		u := publication.update
		s.acquisitionGeometry[u.mutation.deviceID].publication = u.candidate
		receipts[i] = AcquisitionGeometryReceipt{ownedMutation: u.mutation, owner: s, deviceID: u.mutation.deviceID, device: u.candidate.device, present: true, publication: u.candidate, slots: u.candidate.slots}
	}
	return true
}

func (s *Scheduler) acquisitionPublicationCurrentLocked(u *AcquisitionGeometryUpdate) bool {
	state := s.acquisitionGeometry[u.mutation.deviceID]
	if state == nil || state.publication != u.expectedPublication || len(state.pending) != 1 || (u.requireEpoch && s.pointObservationLease != u.expectedEpoch) {
		return false
	}
	_, owned := state.pending[u.mutation]
	return owned && s.deviceConfigs[u.mutation.deviceID] == u.candidate.device && s.acquisitionSlotsCurrentLocked(u.mutation.deviceID, u.candidate.slots) && acquisitionForeignRulesUnchanged(u)
}

func acquisitionForeignRulesUnchanged(u *AcquisitionGeometryUpdate) bool {
	if u.scopeRule == "" || u.candidate.revision == "" {
		return true
	}
	before := u.mutation.before.publication
	if before == nil || before.revision == "" || !u.mutation.before.present || u.mutation.before.device != u.candidate.device {
		return false
	}
	for id, revision := range before.rules {
		if id != u.scopeRule && u.candidate.rules[id] != revision {
			return false
		}
	}
	for id, revision := range u.candidate.rules {
		if id != u.scopeRule && before.rules[id] != revision {
			return false
		}
	}
	return true
}

// RevokeAcquisitionGeometry fails closed only the exact pointer published by
// the failed SQL acceptance. A newer/foreign publication must remain untouched.
func (s *Scheduler) RevokeAcquisitionGeometry(receipt AcquisitionGeometryReceipt) bool {
	if s == nil || receipt.owner != s || receipt.publication == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.acquisitionGeometry[receipt.deviceID]
	if state == nil || state.publication != receipt.publication {
		return false
	}
	state.publication = &acquisitionGeometry{}
	return true
}

// ExpectedLease reads this operation's original or last owned receipt, never a
// fresh scheduler lookup. Absence is an explicit zero lease in the captured set.
func (u *AcquisitionGeometryUpdate) ExpectedLease(pointID string) (uint64, bool) {
	if u == nil || u.sealed || pointID == "" {
		return 0, false
	}
	return u.candidate.slots[pointID].observationLease, true
}

// ApplyAcquisitionDevice rejects a configuration/slot change since the original
// pre-provider capture before the old projection can overwrite that change.
// Unknown source evidence may use this guard but can never publish geometry.
func (s *Scheduler) ApplyAcquisitionDevice(u *AcquisitionGeometryUpdate) bool {
	if s == nil || u == nil || u.owner != s || u.sealed {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.acquisitionGeometry[u.mutation.deviceID]
	current, present := s.deviceConfigs[u.mutation.deviceID]
	before := u.mutation.before
	expected, expectedPresent := before.device, before.present
	if u.deviceApplied {
		// An identical duplicate in an ordinary legacy projection may repeat
		// this install. Journal only our successful write, never a fresh getter.
		expected, expectedPresent = u.candidate.device, true
	}
	if state == nil || state.publication != u.mutation.marker || len(state.pending) != 1 || present != expectedPresent || current != expected || !s.acquisitionSlotsMatchLocked(u.mutation.deviceID, u.candidate.slots, false) {
		return false
	}
	if _, owned := state.pending[u.mutation]; !owned {
		return false
	}
	s.addDeviceLocked(u.candidate.device)
	s.setDeviceAutoIntervalLocked(u.mutation.deviceID, u.autoInterval)
	u.deviceApplied = true
	return true
}

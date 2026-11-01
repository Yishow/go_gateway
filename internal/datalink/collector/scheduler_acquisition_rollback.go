package collector

import "go-gateway/internal/datalink/schema"

// PrepareAcquisitionGeometryRollback uses the original immutable slot/profile
// capture and exact owning publication, never a fresh lookup. The caller must
// reaccept the original SQL witness and supply successful compensation leases.
func (s *Scheduler) PrepareAcquisitionGeometryRollback(mutation AcquisitionGeometryMutation, receipt AcquisitionGeometryReceipt, device schema.Device, definition AcquisitionGeometryDefinition, epoch *uint64) *AcquisitionGeometryUpdate {
	if !s.CanRestoreAcquisitionGeometry(mutation) || mutation.token.before.publication.revision != definition.Revision {
		return nil
	}
	update := s.PrepareAcquisitionGeometry(mutation, device, definition, "")
	if update == nil {
		return nil
	}
	if receipt.owner != nil {
		if receipt.owner != s || receipt.deviceID != device.ID || receipt.publication == nil || receipt.ownedMutation != mutation.token {
			return nil
		}
		update.expectedPublication = receipt.publication
	}
	if epoch != nil {
		update.requireEpoch = true
		update.expectedEpoch = *epoch
	}
	return update
}

// Restored replaces only a captured slot's lease using the source operation's
// successful compensation receipt. Frozen profile and all source fields remain
// original. Exact final slot comparison rejects a mismatching restoration.
func (u *AcquisitionGeometryUpdate) Restored(id, revision string, lease uint64) bool {
	if u == nil || u.sealed || id == "" {
		return false
	}
	old, present := u.candidate.slots[id]
	if lease == 0 {
		return !present
	}
	if !present || old.configRevision != revision {
		return false
	}
	old.observationLease = lease
	u.candidate.slots[id] = old
	return true
}

// CanRestoreAcquisitionGeometry examines only the original pre-mutation
// authorization. An old opaque SQL graph alone cannot revive unknown metadata.
func (s *Scheduler) CanRestoreAcquisitionGeometry(mutation AcquisitionGeometryMutation) bool {
	return s != nil && mutation.token != nil && mutation.token.owner == s && mutation.token.before.authorized
}

// Published reports an owned forward publication, not a later capture.
func (r AcquisitionGeometryReceipt) Published() bool {
	return r.owner != nil && r.ownedMutation != nil && r.publication != nil
}

package sourcerule

import "go-gateway/internal/datalink/schema"

type stubRuntimeSync struct {
	leases    map[string]uint64
	revisions map[string]string
	nextLease uint64
	upserted  []string
	removed   []string
}

func (s *stubRuntimeSync) UpsertPoint(pointRecord *schema.Point) {
	if pointRecord == nil {
		return
	}
	s.upserted = append(s.upserted, pointRecord.ID)
	if s.leases == nil {
		s.leases = map[string]uint64{}
		s.revisions = map[string]string{}
	}
	if !pointRecord.Enabled {
		delete(s.leases, pointRecord.ID)
		delete(s.revisions, pointRecord.ID)
		return
	}
	s.nextLease++
	s.leases[pointRecord.ID] = s.nextLease
	s.revisions[pointRecord.ID] = pointRecord.RevisionID
}

func (s *stubRuntimeSync) RemovePoint(pointID string) {
	s.removed = append(s.removed, pointID)
}

func (s *stubRuntimeSync) CapturePointLease(id, revision string) uint64 {
	if s.revisions[id] != revision {
		return 0
	}
	return s.leases[id]
}
func (s *stubRuntimeSync) InstallPointForMutation(p *schema.Point) uint64 {
	s.UpsertPoint(p)
	return s.leases[p.ID]
}
func (s *stubRuntimeSync) SuspendPointLease(id string, lease uint64) bool {
	if s.leases[id] != lease {
		return false
	}
	delete(s.leases, id)
	delete(s.revisions, id)
	s.removed = append(s.removed, id)
	return true
}
func (s *stubRuntimeSync) RestorePointLease(p *schema.Point, lease uint64) bool {
	if lease == 0 || s.leases[p.ID] != lease {
		return false
	}
	s.UpsertPoint(p)
	return true
}

func (s *stubRuntimeSync) InstallPointForMutationAtLease(p *schema.Point, expected uint64) (uint64, bool) {
	if s.leases[p.ID] != expected {
		return 0, false
	}
	s.UpsertPoint(p)
	return s.leases[p.ID], true
}
func (s *stubRuntimeSync) RemovePointOwned(id string, expected uint64) (uint64, bool) {
	if s.leases[id] != expected {
		return 0, false
	}
	delete(s.leases, id)
	delete(s.revisions, id)
	s.removed = append(s.removed, id)
	return s.nextLease, true
}
func (s *stubRuntimeSync) RestoreRemovedPointOwned(p *schema.Point, removedLease, epoch uint64) (uint64, bool) {
	if removedLease == 0 || s.leases[p.ID] != 0 || s.nextLease != epoch {
		return 0, false
	}
	s.UpsertPoint(p)
	return s.leases[p.ID], true
}
func (s *stubRuntimeSync) RestorePointLeaseReceipt(p *schema.Point, lease uint64) (uint64, bool) {
	if !s.RestorePointLease(p, lease) {
		return 0, false
	}
	return s.leases[p.ID], true
}

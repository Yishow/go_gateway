package collector

import (
	"cmp"
	"slices"

	"go-gateway/internal/datalink/optimizer"
)

// acquisitionRequestHasGap uses only the canonical optimizer's component
// slices, not a new address parser. Overlapping active spans still cover cells.
func acquisitionRequestHasGap(request optimizer.AcquisitionRequest) bool {
	if request.Kind != "contiguous" {
		return false
	}
	points := slices.Clone(request.Points)
	slices.SortFunc(points, func(a, b optimizer.PointSlice) int { return cmp.Compare(a.Offset, b.Offset) })
	end := 0
	for _, point := range points {
		if point.Offset > end {
			return true
		}
		end = max(end, point.Offset+point.Length)
	}
	return end < request.Length
}

func (s *Scheduler) admitScheduledWire(snapshot scheduledSnapshot, points []pointInfo, gap bool) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.scheduledCurrentLocked(snapshot, points); err != nil {
		return err
	}
	// Only a read of otherwise unselected cells needs geometry authorization.
	// Quarantine never disables exact contiguous or active-only random requests.
	if gap && (snapshot.geometry == nil || s.currentAcquisitionGeometryLocked(snapshot.device.ID) != snapshot.geometry) {
		return errAcquisitionSuperseded
	}
	return nil
}

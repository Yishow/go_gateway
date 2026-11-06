package point

import (
	"context"

	"go-gateway/internal/datalink/schema"
)

func (r *MemoryRepository) updateMaterializedRevisions(_ context.Context, updates []RevisionRestore) error {
	if len(updates) == 0 {
		return ErrFormatInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	seen := make(map[string]bool, len(updates))
	deviceID := ""
	for _, update := range updates {
		p := update.Point
		if p == nil || p.ID == "" || update.ExpectedRevision == "" || seen[p.ID] {
			return ErrFormatInvalid
		}
		seen[p.ID] = true
		if deviceID == "" {
			deviceID = p.DeviceID
		} else if deviceID != p.DeviceID {
			return ErrFormatInvalid
		}
		old := r.points[p.ID]
		if old == nil || old.RevisionID != update.ExpectedRevision || old.DeviceID != p.DeviceID {
			return ErrFormatConflict
		}
		if err := validateMemoryOverrideUpdate(old, p); err != nil {
			return err
		}
	}
	occupied := map[string]bool{}
	key := func(p *schema.Point) string { return p.DeviceID + "\x00" + p.Address + "\x00" + p.Function }
	for id, p := range r.points {
		if !seen[id] {
			occupied[key(p)] = true
		}
	}
	for _, update := range updates {
		p := update.Point
		if occupied[key(p)] {
			return ErrFormatConflict
		}
		occupied[key(p)] = true
	}
	for _, update := range updates {
		old, next := r.points[update.Point.ID], clonePoint(update.Point)
		next.LastValue, next.LastReadAt, next.LastError = old.LastValue, old.LastReadAt, old.LastError
		r.points[next.ID] = next
	}
	return nil
}

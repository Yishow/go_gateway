package point

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrReadResultSuperseded = errors.New("point read result superseded")
var ErrReadResultCASUnavailable = errors.New("conditional point read result persistence unavailable")

// ReadResultRevisionWriter is a telemetry-only CAS, not a configuration edit.
// The expected identity must come from the captured read, never a later Get.
type ReadResultRevisionWriter interface {
	UpdateReadResultAtRevision(context.Context, string, string, string, any, string) error
}

func (s *Service) UpdateReadResultAtRevision(ctx context.Context, id, deviceID, revision string, value any, message string) error {
	writer, ok := s.repo.(ReadResultRevisionWriter)
	if !ok {
		return ErrReadResultCASUnavailable
	}
	return writer.UpdateReadResultAtRevision(ctx, id, deviceID, revision, value, message)
}
func (r *SQLRepository) UpdateReadResultAtRevision(ctx context.Context, id, deviceID, revision string, value any, message string) error {
	var lastValue *string
	if message == "" {
		text := fmt.Sprintf("%v", value)
		lastValue = &text
	}
	now := time.Now()
	result, err := r.db.ExecContext(ctx, `UPDATE points SET last_value=?,last_read_at=?,last_error=?,updated_at=? WHERE id=? AND device_id=? AND revision_id=?`, lastValue, now, message, now, id, deviceID, revision)
	if err != nil {
		return fmt.Errorf("conditional point telemetry update failed: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrReadResultSuperseded
	}
	return nil
}
func (r *MemoryRepository) UpdateReadResultAtRevision(_ context.Context, id, deviceID, revision string, value any, message string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.points[id]
	if !ok || p.DeviceID != deviceID || p.RevisionID != revision {
		return ErrReadResultSuperseded
	}
	now := time.Now()
	p.LastReadAt = &now
	p.LastError = message
	p.UpdatedAt = now
	if message == "" {
		if encoded, err := json.Marshal(value); err == nil {
			text := string(encoded)
			p.LastValue = &text
		}
	}
	return nil
}

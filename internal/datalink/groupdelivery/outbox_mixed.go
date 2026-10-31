package groupdelivery

import (
	"database/sql"
	"encoding/json"
	"go-gateway/internal/datalink/snapshot"
	"time"
)

func validateOutboxDispatch(item *OutboxItem, bucket string, due sql.NullString) error {
	if !due.Valid {
		var wire rowPayload
		if json.Unmarshal(item.Payload, &wire) == nil && (wire.Version == 2 || wire.DueAt != "") {
			return mixedInvalid("missing mixed dispatch column")
		}
		if bucket == "" {
			return mixedInvalid("outbox missing dispatch")
		}
		return nil
	}
	at, err := time.Parse(time.RFC3339Nano, due.String)
	if err != nil || at.IsZero() {
		return mixedInvalid("outbox due at")
	}
	if _, err = snapshot.CheckedWindowAdd(at, 0, "outbox dispatch"); err != nil {
		return err
	}
	if handoffHash(string(item.Payload)) != item.PayloadDigest {
		return mixedInvalid("outbox payload digest")
	}
	row, err := DecodeRowPayload(item.Payload)
	if err != nil {
		return err
	}
	if !row.DueAt.Equal(at) || !row.BucketStart.Equal(item.BucketStart) || row.RecordID != item.RecordID || row.EffectKey != item.EffectKey {
		return mixedInvalid("outbox scope mismatch")
	}
	var expected string
	if bucket == "" {
		expected, err = snapshot.MixedDispatchRecordID(item.Key.WorkspaceID, item.Key.GroupID, item.Key.GroupRevision, at)
	} else {
		expected, err = snapshot.ProfileRecordID(item.Key.WorkspaceID, item.Key.GroupID, item.Key.GroupRevision, item.EntityKey, "mixed", item.BucketStart)
	}
	if err != nil || expected != item.RecordID {
		return mixedInvalid("outbox record identity")
	}
	item.DueAt = at
	return nil
}

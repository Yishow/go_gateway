package snapshot

import (
	"go-gateway/internal/datalink/measurement"
	"time"
)

// MixedDispatchRecordID is the existing mixed-wide dispatch identity. Both the
// assembler and frozen encoder use this one authority; it is not a common
// member window and intentionally differs from the legacy bucket RecordID.
func MixedDispatchRecordID(workspaceID, groupID, groupRevision string, dueAt time.Time) (string, error) {
	return measurement.HashStringTuple("write-group-mixed-dispatch-v1", workspaceID, groupID, groupRevision, dueAt.UTC().Format(time.RFC3339Nano))
}

package dbtarget

import (
	"time"

	"go-gateway/internal/datalink/schema"
)

const (
	skipMappingNotCurrent = "mapping_not_current"
	skipWriteGroupOwned   = "write_group_owned"
	skipMappingCheckError = "mapping_check_failed"
	skipColumnReplaced    = "group_column_replaced"
	skipSnapshotReplaced  = "snapshot_superseded"
	skipIntervalFlushed   = "interval_already_flushed"
)

func (w *Writer) recordMappingSkip(connector *schema.DatabaseConnector, mapping *schema.DatabaseTargetMapping, observedAt time.Time, reason string) {
	if connector == nil || mapping == nil {
		return
	}
	attemptedAt := w.now().UTC()
	recordWriteHistory(connector.ID, WriteHistoryRecord{
		ObservedAt: observedAt, AttemptedAt: &attemptedAt, Status: "skipped", MappingID: mapping.ID,
		GroupKey: cloneOptionalString(mapping.GroupKey), TableName: mapping.TableName,
		EffectiveIntervalSeconds: effectiveWriteIntervalSeconds(mapping, connector), SkipReason: reason,
	})
}

// sameSnapshotBinding uses the same persisted target identity checked before
// delivery. A change must not inherit the old target's pending deadline.
func sameSnapshotBinding(a *schema.DatabaseTargetMapping, ac *schema.DatabaseConnector, b *schema.DatabaseTargetMapping, bc *schema.DatabaseConnector) bool {
	if a == nil || b == nil || !sameConnectorIdentity(ac, bc) {
		return false
	}
	return a.ID == b.ID && a.TagID == b.TagID && a.ConnectorID == b.ConnectorID &&
		a.TableSchema == b.TableSchema && a.TableName == b.TableName && a.ColumnName == b.ColumnName &&
		a.WriteMode == b.WriteMode && derefOptionalString(a.TimestampColumn) == derefOptionalString(b.TimestampColumn) &&
		derefOptionalString(a.GroupKey) == derefOptionalString(b.GroupKey) &&
		effectiveWriteIntervalSeconds(a, ac) == effectiveWriteIntervalSeconds(b, bc)
}

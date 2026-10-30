package workspace

import (
	"context"
	"fmt"
	"strings"
)

// checkMixedTargetOwnership permits source reuse, never a second independent
// managed-wide writer for the same configured destination. Physical catalog
// ownership/readiness remains mandatory before any writer can activate.
func (r *SQLWriteGroupRepository) checkMixedTargetOwnership(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup) error {
	if isMixedWriteGroup(group) {
		if _, err := mixedRecordProfileKey(group); err != nil {
			return err
		}
	}
	rows, err := runner.QueryContext(ctx, r.query(`SELECT id, destination_connector_id, destination_database,
 destination_table_schema, destination_table_name, COALESCE(layout_version,''), COALESCE(profile_kind,''), COALESCE(storage_mode,'')
 FROM write_groups WHERE workspace_id=$1 AND id<>$2`), group.WorkspaceID, group.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id, connector, database, namespace, table, layout, kind, mode string
		if err := rows.Scan(&id, &connector, &database, &namespace, &table, &layout, &kind, &mode); err != nil {
			return err
		}
		if !isMixedWriteGroup(group) && kind != string(WriteGroupProfileMixed) {
			continue
		}
		if !strings.EqualFold(namespace, group.Destination.TableSchema) || !strings.EqualFold(table, group.Destination.TableName) {
			continue
		}
		if connector != group.Destination.ConnectorID && database != "" && group.Destination.Database != "" && database != group.Destination.Database {
			continue
		}
		// The already-supported long table partitions rows by group/member identity.
		// Only the exact connector/database/schema/storage-mode scope may share it.
		if group.isLongLayout() && layout == string(WriteGroupLayoutManagedLongV1) && connector == group.Destination.ConnectorID && database == group.Destination.Database && mode == string(group.StorageMode) {
			continue
		}
		return fmt.Errorf("%w: destination table already belongs to group %s", ErrWriteGroupCircuitProfileConflict, id)
	}
	return rows.Err()
}

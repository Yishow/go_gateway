package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

func (s *WriteGroupService) purgeUnsharedStorageScopesInTx(ctx context.Context, tx *sql.Tx, deleted []*WriteGroup) error {
	candidates := make(map[longStorageScope]bool)
	surviving := make(map[longStorageScope]bool)
	deletedIDs := make(map[string]bool, len(deleted))
	for _, group := range deleted {
		deletedIDs[group.ID] = true
		if group.isLongLayout() {
			candidates[longStorageScopeOf(group)] = true
		}
	}
	remaining, err := s.repo.ListInTx(ctx, tx, deleted[0].WorkspaceID)
	if err != nil {
		return fmt.Errorf("read surviving recording storage: %w", err)
	}
	for _, group := range remaining {
		if !deletedIDs[group.ID] && group.isLongLayout() {
			surviving[longStorageScopeOf(group)] = true
		}
	}
	// Collect scopes from the deleted lineage's historical versions as well as
	// its current draft. A surviving version can still use an old scope even
	// after the current draft moved elsewhere. Unreadable versions fail closed.
	for _, table := range []string{"write_group_versions", "wg_runtime_versions"} {
		query := `SELECT v.group_id,v.payload FROM ` + table + ` v JOIN write_groups w ON w.id=v.group_id WHERE w.workspace_id=$1`
		rows, err := tx.QueryContext(ctx, s.repo.query(query), deleted[0].WorkspaceID)
		if err != nil {
			return fmt.Errorf("read surviving recording versions: %w", err)
		}
		for rows.Next() {
			var id, payload string
			if err := rows.Scan(&id, &payload); err != nil {
				rows.Close()
				return fmt.Errorf("scan surviving recording version: %w", err)
			}
			var group WriteGroup
			if table == "wg_runtime_versions" {
				var snapshot WriteGroupAppliedSnapshot
				if json.Unmarshal([]byte(payload), &snapshot) != nil || snapshot.Group == nil {
					rows.Close()
					return nil
				}
				group = *snapshot.Group
			} else if json.Unmarshal([]byte(payload), &group) != nil {
				rows.Close()
				return nil
			}
			if group.ID != id || group.WorkspaceID != deleted[0].WorkspaceID {
				rows.Close()
				return nil
			}
			if group.isLongLayout() {
				scope := longStorageScopeOf(&group)
				if deletedIDs[id] {
					candidates[scope] = true
				} else {
					surviving[scope] = true
				}
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return fmt.Errorf("iterate surviving recording versions: %w", err)
		}
		if err := rows.Close(); err != nil {
			return err
		}
	}
	for scope := range candidates {
		if surviving[scope] {
			continue
		}
		for _, table := range []string{"write_group_mysql_partition_maintenance", "write_group_sqlserver_partition_maintenance", "write_group_storage_scopes"} {
			query := `DELETE FROM ` + table + ` WHERE workspace_id=$1 AND connector_id=$2 AND database_name=$3 AND table_schema=$4 AND storage_mode=$5 AND layout_version=$6`
			if _, err := tx.ExecContext(ctx, s.repo.query(query), scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.Mode, scope.Layout); err != nil {
				return fmt.Errorf("delete unshared recording storage from %s: %w", table, err)
			}
		}
	}
	return nil
}

package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// migrationSourceOwnershipIssues prevents two migration adapters from
// claiming one legacy target mapping. An existing map for the exact source
// identity is the only replay exemption; tombstoned groups retain ownership.
func migrationSourceOwnershipIssues(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	workspaceID string,
	sourceKind string,
	sourceID string,
	migration WriteGroupMigration,
) ([]WriteGroupMigrationFinding, error) {
	rows, err := tx.QueryContext(ctx, repo.query(`
		SELECT id, migration
		FROM write_groups
		WHERE workspace_id = $1
		ORDER BY id ASC
	`), workspaceID)
	if err != nil {
		return nil, fmt.Errorf("read migration source ownership: %w", err)
	}
	defer rows.Close()
	mapGroupID, mapFound, err := repo.getMigrationMapInTx(ctx, tx, workspaceID, sourceKind, sourceID)
	if err != nil {
		return nil, err
	}
	issues := make([]WriteGroupMigrationFinding, 0)
	for rows.Next() {
		var groupID, rawMigration string
		if err := rows.Scan(&groupID, &rawMigration); err != nil {
			return nil, fmt.Errorf("scan migration source ownership: %w", err)
		}
		var existing WriteGroupMigration
		if strings.TrimSpace(rawMigration) != "" {
			if err := json.Unmarshal([]byte(rawMigration), &existing); err != nil {
				return nil, fmt.Errorf("decode migration source ownership: %w", ErrWriteGroupValidation)
			}
		}
		for _, source := range migration.SourceIDs {
			if !slices.Contains(existing.SourceIDs, source) {
				continue
			}
			if mapFound && strings.TrimSpace(existing.SourceKind) == strings.TrimSpace(sourceKind) && strings.TrimSpace(mapGroupID.GroupID) == strings.TrimSpace(groupID) {
				continue
			}
			issues = append(issues, WriteGroupMigrationFinding{
				Code:    "migration-source-already-owned",
				Message: fmt.Sprintf("legacy target mapping %s is already owned by canonical migration source %s", source, strings.TrimSpace(existing.SourceKind)),
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate migration source ownership: %w", err)
	}
	return issues, nil
}

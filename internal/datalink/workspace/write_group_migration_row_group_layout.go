package workspace

import (
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

func rowGroupMigrationDifferences() []WriteGroupMigrationFinding {
	return []WriteGroupMigrationFinding{
		{Code: "legacy-arrival-order-vs-max-observed", Message: "the snapshot selects the maximum observed time in each bucket; arrival order is not the selection rule"},
		{Code: "snapshot-observed-time-tie-sample-id", Message: "equal observed times use the greatest stable sample_id as the tie breaker"},
		{Code: "snapshot-utc-year-one-epoch-alignment", Message: "UTC bucket alignment uses Unix epoch boundaries instead of the legacy year-one origin"},
		{Code: "snapshot-late-sample-does-not-reopen", Message: "samples arriving after the lateness boundary do not reopen a closed snapshot bucket"},
		{Code: "snapshot-missing-or-stale-skips-row", Message: "a missing or stale required member skips the snapshot row under the saved freshness policy"},
		{Code: "legacy-timestamp-intent-not-bound", Message: "legacy timestamp intent is retained for review and is not bound to a new SQL column"},
	}
}

func validateRowGroupTargetLayout(
	rowGroup DatabaseRowGroup,
	connector writeGroupMigrationConnector,
	members []writeGroupRowGroupMember,
	addIssue func(string, string),
) {
	if len(members) == 0 {
		return
	}
	first := members[0].target
	interval := migrationPreviewInterval(first.writeGroupMigrationTarget, connector)
	sharedColumns := make(map[string]int)
	entityColumns := make(map[string]struct{}, len(members))
	firstTimestamp := migrationPreviewOptionalString(first.timestampColumn)
	if rowGroup.ConnectorID != "" && rowGroup.ConnectorID != first.connectorID {
		addIssue("row-group-target-scope-mismatch", "the saved row group connector does not match its target mappings")
	}
	if rowGroup.TableSchema != "" && rowGroup.TableSchema != first.tableSchema {
		addIssue("row-group-target-scope-mismatch", "the saved row group schema does not match its target mappings")
	}
	if rowGroup.TableName != "" && rowGroup.TableName != first.tableName {
		addIssue("row-group-target-scope-mismatch", "the saved row group table does not match its target mappings")
	}
	if strings.TrimSpace(rowGroup.ConnectorID) == "" || strings.TrimSpace(rowGroup.TableSchema) == "" || strings.TrimSpace(rowGroup.TableName) == "" {
		addIssue("row-group-scope-incomplete", "the saved row group has an incomplete connector, schema or table scope")
	}
	if _, err := migrationPreviewDestination(first.writeGroupMigrationTarget, connector); err != nil {
		addIssue("connector-scope-invalid", "the saved connector scope cannot be interpreted safely")
	} else if destination, destinationErr := migrationPreviewDestination(first.writeGroupMigrationTarget, connector); destinationErr == nil && destination.TableSchema != first.tableSchema {
		addIssue("connector-scope-mismatch", "the saved connector scope cannot preserve the row-group destination schema")
	}
	for _, member := range members {
		target := member.target
		if strings.TrimSpace(target.tableSchema) == "" || strings.TrimSpace(target.tableName) == "" || strings.TrimSpace(target.columnName) == "" {
			addIssue("row-group-target-scope-incomplete", "a row-group target mapping has an incomplete table or column scope")
		}
		if target.groupKey == nil || strings.TrimSpace(*target.groupKey) == "" {
			addIssue("row-group-group-key-missing", "every row-group member requires a non-empty persisted GroupKey")
		}
		if target.connectorID != first.connectorID || target.tableSchema != first.tableSchema || target.tableName != first.tableName {
			addIssue("row-group-target-scope-mismatch", "row-group targets do not share one connector, schema and table scope")
		}
		if strings.TrimSpace(target.writeMode) != string(schema.DatabaseWriteModeInsert) {
			addIssue("upsert-semantic-difference", "row-group target mappings must use insert; upsert is not representable by a basic snapshot")
		}
		if target.writeIntervalSeconds != nil && *target.writeIntervalSeconds < 0 {
			addIssue("negative-write-interval", "the row-group target mapping has a negative write interval")
		}
		if migrationPreviewInterval(target.writeGroupMigrationTarget, connector) != interval {
			addIssue("row-group-interval-mismatch", "row-group target mappings do not share one effective snapshot interval")
		}
		if migrationPreviewOptionalString(target.timestampColumn) != firstTimestamp {
			addIssue("row-group-timestamp-mismatch", "row-group target mappings use different legacy timestamp intent and cannot share one snapshot layout")
		}
		sharedColumns[target.columnName]++
		groupKey := ""
		if target.groupKey != nil {
			groupKey = strings.TrimSpace(*target.groupKey)
		}
		entityColumnKey := groupKey + "\x00" + target.columnName
		if _, exists := entityColumns[entityColumnKey]; exists {
			addIssue("row-group-entity-column-collision", "two row-group members use the same entity key and destination column")
		}
		entityColumns[entityColumnKey] = struct{}{}
	}
	for column, count := range sharedColumns {
		if count < 2 {
			continue
		}
		if len(rowGroup.GroupKeyColumns) == 0 {
			addIssue("row-group-group-key-columns-missing", "shared-column row groups require persisted group-key metadata")
		}
		if len(rowGroup.UniqueKeyColumns) == 0 {
			addIssue("row-group-unique-key-columns-missing", fmt.Sprintf("shared destination column %s requires persisted unique-key metadata", column))
		}
	}
}

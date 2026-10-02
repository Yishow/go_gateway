package workspace

import (
	"encoding/json"
	"fmt"
	"strings"
)

const defaultSingleMappingSnapshotIntervalSeconds = 15

func migrationPreviewSourceDigest(
	source *writeGroupMigrationSource,
	target writeGroupMigrationTarget,
	connector writeGroupMigrationConnector,
) (string, error) {
	mappingRevision, err := source.mapping.revision()
	if err != nil {
		return "", fmt.Errorf("derive migration mapping revision: %w", err)
	}
	sourceRevision, err := source.source.revision()
	if err != nil {
		return "", fmt.Errorf("derive migration source revision: %w", err)
	}
	connectorConfigRevision, err := hashWriteGroupRevision(connector.connectionConfig)
	if err != nil {
		return "", fmt.Errorf("derive migration connector revision: %w", err)
	}
	return hashWriteGroupRevision(
		"single-mapping-preview",
		source.mapping.ID,
		mappingRevision,
		sourceRevision,
		target.id,
		target.tagID,
		target.connectorID,
		target.tableSchema,
		target.tableName,
		target.columnName,
		target.writeMode,
		migrationPreviewOptionalString(target.timestampColumn),
		migrationPreviewOptionalString(target.groupKey),
		migrationPreviewOptionalInt(target.writeIntervalSeconds),
		fmt.Sprint(target.enabled),
		connector.id,
		connector.kind,
		connector.identityRevision,
		connectorConfigRevision,
		fmt.Sprint(connector.defaultWriteIntervalSecs),
	)
}

func migrationPreviewIntent(
	source *writeGroupMigrationSource,
	target writeGroupMigrationTarget,
	connector writeGroupMigrationConnector,
) *WriteGroupMigrationIntent {
	intent := &WriteGroupMigrationIntent{
		SourceID:             target.id,
		DeviceID:             source.source.DeviceID,
		PointID:              source.source.PointID,
		TagID:                target.tagID,
		ConnectorID:          connector.id,
		ConnectorRevision:    connector.identityRevision,
		TableSchema:          target.tableSchema,
		TableName:            target.tableName,
		ColumnName:           target.columnName,
		WriteMode:            target.writeMode,
		TimestampColumn:      cloneOptionalString(target.timestampColumn),
		GroupKey:             cloneOptionalString(target.groupKey),
		WriteIntervalSeconds: cloneOptionalInt(target.writeIntervalSeconds),
		IntervalSource:       migrationPreviewIntervalSource(target, connector),
		Enabled:              target.enabled,
	}
	// Invalid connector scope is already a blocking preview issue. Preserve
	// the resolved database for valid scopes without exposing its credentials.
	if destination, err := migrationPreviewDestination(target, connector); err == nil {
		intent.Database = destination.Database
	}
	return intent
}

func migrationPreviewDifferences(
	target writeGroupMigrationTarget,
	connector writeGroupMigrationConnector,
) []WriteGroupMigrationFinding {
	differences := []WriteGroupMigrationFinding{
		{
			Code:    "legacy-every-sample-vs-periodic-snapshot",
			Message: "legacy single-mapping writes are every-sample writes and are not equivalent to a basic periodic snapshot",
		},
		{
			Code:    "snapshot-bucket-max-observed-tie-sample-id",
			Message: "each snapshot bucket selects the sample with the maximum observed time; equal observed times use the greatest stable sample_id",
		},
		{
			Code:    "snapshot-late-sample-does-not-reopen",
			Message: "samples arriving after the bucket lateness boundary are marked late and do not reopen or rewrite a closed snapshot",
		},
		{
			Code:    "snapshot-missing-or-stale-skips-row",
			Message: "a missing or stale required member causes the snapshot row to be skipped under the saved skip_row and freshness policy",
		},
	}
	if target.timestampColumn != nil && strings.TrimSpace(*target.timestampColumn) != "" {
		differences = append(differences, WriteGroupMigrationFinding{
			Code:    "legacy-timestamp-column-ignored",
			Message: "legacy INSERT ignores the configured timestamp column; the original intent is retained, while this candidate has no timestamp binding and review does not activate a writer",
		})
	}
	switch {
	case target.writeIntervalSeconds != nil && *target.writeIntervalSeconds > 0:
		differences = append(differences, WriteGroupMigrationFinding{
			Code:    "legacy-interval-override",
			Message: "the target mapping supplies an interval override that must be reviewed against the basic snapshot interval",
		})
	case target.writeIntervalSeconds != nil:
		effectiveInterval := migrationPreviewInterval(target, connector)
		source := migrationPreviewIntervalSource(target, connector)
		sourceLabel := strings.ReplaceAll(source, "-", " ")
		differences = append(differences, WriteGroupMigrationFinding{
			Code:    "legacy-interval-nonpositive-fallback",
			Message: fmt.Sprintf("the legacy target mapping stores non-positive interval %d; the candidate ignores it and uses the %s interval of %d seconds", *target.writeIntervalSeconds, sourceLabel, effectiveInterval),
		})
	case connector.defaultWriteIntervalSecs > 0:
		differences = append(differences, WriteGroupMigrationFinding{
			Code:    "connector-interval-fallback",
			Message: "the candidate interval comes from the saved connector default because the target mapping has no override",
		})
	default:
		differences = append(differences, WriteGroupMigrationFinding{
			Code:    "snapshot-interval-defaulted",
			Message: fmt.Sprintf("the candidate uses the design default snapshot interval of %d seconds because no positive legacy or connector interval is saved", defaultSingleMappingSnapshotIntervalSeconds),
		})
	}
	return differences
}

func migrationPreviewCandidate(
	record *Record,
	source *writeGroupMigrationSource,
	target writeGroupMigrationTarget,
	connector writeGroupMigrationConnector,
	sourceDigest string,
	sourceRevision string,
	mappingRevision string,
) *WriteGroup {
	destination, err := migrationPreviewDestination(target, connector)
	if err != nil {
		return nil
	}
	name := strings.TrimSpace(source.source.TagDisplayName)
	if name == "" {
		name = strings.TrimSpace(source.source.TagKey)
	}
	if name == "" {
		name = target.tableName
	}
	if name == "" {
		name = "legacy mapping " + target.id
	}
	interval := migrationPreviewInterval(target, connector)
	return &WriteGroup{
		WorkspaceID: record.ID,
		Name:        name,
		Status:      WriteGroupStatusDraft,
		Members: []WriteGroupMember{{
			DeviceID:        source.source.DeviceID,
			PointID:         source.source.PointID,
			TagID:           source.source.TagID,
			SourceRevision:  sourceRevision,
			MappingRevision: mappingRevision,
			TargetColumn:    target.columnName,
			Required:        true,
			MaxAgeSeconds:   intPointer(interval),
		}},
		Destination: destination,
		RowPolicy: WriteGroupRowPolicy{
			IntervalSeconds:        interval,
			AllowedLatenessSeconds: 0,
			IncompletePolicy:       writeGroupIncompletePolicySkipRow,
		},
		WritePolicy: WriteGroupWritePolicy{Mode: writeGroupWriteModeAppend},
		Migration: WriteGroupMigration{
			SourceKind:     "legacy-single-mapping",
			SourceIDs:      []string{target.id},
			SourceRevision: sourceDigest,
			AdapterVersion: writeGroupMigrationPreviewAdapterVersion,
			ReviewResult:   WriteGroupMigrationPreviewStatusNeedsReview,
		},
	}
}

func migrationPreviewDigest(preview *WriteGroupMigrationPreview) (string, error) {
	copyPreview := *preview
	copyPreview.ReviewDigest = ""
	payload, err := json.Marshal(copyPreview)
	if err != nil {
		return "", err
	}
	return hashWriteGroupRevision(string(payload))
}

func migrationPreviewIntervalSource(target writeGroupMigrationTarget, connector writeGroupMigrationConnector) string {
	if target.writeIntervalSeconds != nil && *target.writeIntervalSeconds > 0 {
		return "legacy-target-override"
	}
	if connector.defaultWriteIntervalSecs > 0 {
		return "connector-default"
	}
	return "design-default"
}

func migrationPreviewInterval(target writeGroupMigrationTarget, connector writeGroupMigrationConnector) int {
	if target.writeIntervalSeconds != nil && *target.writeIntervalSeconds > 0 {
		return *target.writeIntervalSeconds
	}
	if connector.defaultWriteIntervalSecs > 0 {
		return connector.defaultWriteIntervalSecs
	}
	return defaultSingleMappingSnapshotIntervalSeconds
}

func intPointer(value int) *int {
	return &value
}

func migrationPreviewOptionalString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func migrationPreviewOptionalInt(value *int) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(*value)
}

func cloneOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneOptionalInt(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func migrationPreviewDestination(
	target writeGroupMigrationTarget,
	connector writeGroupMigrationConnector,
) (WriteGroupDestination, error) {
	destination := WriteGroupDestination{
		ConnectorID:       connector.id,
		ConnectorRevision: connector.identityRevision,
		TableSchema:       target.tableSchema,
		TableName:         target.tableName,
		StorageStrategy:   WriteGroupStorageStrategyCustom,
	}
	if err := resolveWriteGroupDestinationScope(connector.kind, connector.connectionConfig, &destination); err != nil {
		return WriteGroupDestination{}, err
	}
	return destination, nil
}

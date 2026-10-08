package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/common"
)

const writeGroupRowGroupMigrationAdapterVersion = "row-group-v1"

type writeGroupRowGroupTarget struct {
	writeGroupMigrationTarget
	pointID string
}

type writeGroupRowGroupMember struct {
	pointID string
	target  writeGroupRowGroupTarget
	source  *writeGroupMigrationSource
}

func (s *WriteGroupService) previewRowGroupMigration(
	ctx context.Context,
	workspaceID string,
	sourceIDs []string,
) (result *WriteGroupMigrationPreview, err error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	workspaceID = strings.TrimSpace(workspaceID)
	if workspaceID == "" {
		return nil, fmt.Errorf("row-group migration preview workspace id: %w", ErrWriteGroupValidation)
	}
	ids, err := normalizeMigrationPreviewIDs(sourceIDs)
	if err != nil {
		return nil, err
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("begin row-group migration preview", err)
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback row-group migration preview: %w", rollbackErr))
		}
	}()

	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("preview row-group migration workspace")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read row-group migration workspace", err)
	}
	if record.ID != workspaceID || strings.TrimSpace(record.DatabaseConnectorID) == "" {
		return nil, writeGroupNotFound("preview row-group migration workspace")
	}
	connector, err := readMigrationPreviewConnector(ctx, setupTx.SQLTx(), s.repo, record.DatabaseConnectorID)
	if err != nil {
		return nil, err
	}
	return buildRowGroupMigrationPreview(ctx, setupTx.SQLTx(), s.repo, record, connector, ids)
}

func buildRowGroupMigrationPreview(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	record *Record,
	connector writeGroupMigrationConnector,
	ids []string,
) (*WriteGroupMigrationPreview, error) {
	preview := &WriteGroupMigrationPreview{
		WorkspaceID:       record.ID,
		WorkspaceRevision: record.DatabaseSetupRevision,
		ConnectorRevision: connector.identityRevision,
		AdapterVersion:    writeGroupRowGroupMigrationAdapterVersion,
		Items:             make([]WriteGroupMigrationPreviewItem, 0, len(ids)),
	}
	for _, id := range ids {
		item, err := readRowGroupMigrationPreviewItem(ctx, tx, repo, record, connector, id)
		if err != nil {
			return nil, err
		}
		preview.Items = append(preview.Items, item)
	}
	blockOverlappingRowGroupSources(preview.Items)
	digest, err := migrationPreviewDigest(preview)
	if err != nil {
		return nil, fmt.Errorf("derive row-group migration preview digest: %w", err)
	}
	preview.ReviewDigest = digest
	return preview, nil
}

func blockOverlappingRowGroupSources(items []WriteGroupMigrationPreviewItem) {
	owners := make(map[string][]int)
	for index := range items {
		item := items[index]
		if item.Status != WriteGroupMigrationPreviewStatusNeedsReview || item.CandidateGroup == nil {
			continue
		}
		for _, sourceID := range item.CandidateGroup.Migration.SourceIDs {
			sourceID = strings.TrimSpace(sourceID)
			if sourceID != "" {
				owners[sourceID] = append(owners[sourceID], index)
			}
		}
	}
	overlaps := make([]string, 0)
	for sourceID, indexes := range owners {
		if len(indexes) > 1 {
			overlaps = append(overlaps, sourceID)
		}
	}
	slices.Sort(overlaps)
	for _, sourceID := range overlaps {
		for _, index := range owners[sourceID] {
			item := &items[index]
			item.Issues = append(item.Issues, WriteGroupMigrationFinding{
				Code:    "row-group-source-overlap",
				Message: fmt.Sprintf("legacy target mapping %s is selected by more than one row group in this review batch", sourceID),
			})
			item.SourceRevision = ""
			item.CandidateGroup = nil
			item.Status = WriteGroupMigrationPreviewStatusBlocked
		}
	}
}

func readRowGroupMigrationPreviewItem(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	record *Record,
	connector writeGroupMigrationConnector,
	rowGroupID string,
) (WriteGroupMigrationPreviewItem, error) {
	rowGroup, found := findDatabaseRowGroup(record.DatabaseRowGroups, rowGroupID)
	if !found {
		return WriteGroupMigrationPreviewItem{}, writeGroupNotFound("preview row-group migration source")
	}
	item := WriteGroupMigrationPreviewItem{
		SourceID:    rowGroupID,
		Differences: rowGroupMigrationDifferences(),
		Issues:      make([]WriteGroupMigrationFinding, 0),
		BeforeRowGroupIntent: &WriteGroupRowGroupMigrationIntent{
			SourceID:          rowGroupID,
			ConnectorID:       connector.id,
			ConnectorRevision: connector.identityRevision,
			RowGroup:          cloneDatabaseRowGroupValue(rowGroup),
			Members:           make([]WriteGroupMigrationIntent, 0, len(rowGroup.MemberPointIDs)),
		},
	}
	addIssue := func(code, message string) {
		item.Issues = append(item.Issues, WriteGroupMigrationFinding{Code: code, Message: message})
	}
	if rowGroup.ConnectorID != "" && rowGroup.ConnectorID != connector.id {
		addIssue("row-group-connector-mismatch", "the saved row group belongs to a different connector")
	}
	if connector.kind != writeGroupConnectorKindSQLite && connector.kind != writeGroupConnectorKindPostgres {
		addIssue("connector-kind-out-of-batch", "the saved connector kind is outside this preview adapter")
	}
	if len(rowGroup.MemberPointIDs) == 0 {
		addIssue("row-group-members-missing", "the saved row group has no member points")
	}

	members := make([]writeGroupRowGroupMember, 0, len(rowGroup.MemberPointIDs))
	seenPointIDs := make(map[string]struct{}, len(rowGroup.MemberPointIDs))
	for _, pointID := range rowGroup.MemberPointIDs {
		pointID = strings.TrimSpace(pointID)
		if pointID == "" {
			addIssue("row-group-member-missing", "the saved row group contains an empty member point")
			continue
		}
		if _, exists := seenPointIDs[pointID]; exists {
			addIssue("row-group-member-ambiguous", "the saved row group contains a duplicate member point")
			continue
		}
		seenPointIDs[pointID] = struct{}{}
		targets, err := readRowGroupTargetsForPoint(ctx, tx, repo, connector.id, pointID)
		if err != nil {
			return WriteGroupMigrationPreviewItem{}, err
		}
		enabledTargets := make([]writeGroupRowGroupTarget, 0, len(targets))
		for _, target := range targets {
			if target.enabled {
				enabledTargets = append(enabledTargets, target)
			}
		}
		if len(enabledTargets) == 0 {
			allSources, sourceErr := readMigrationPreviewSourcesForPoint(ctx, tx, repo, pointID)
			if sourceErr != nil {
				return WriteGroupMigrationPreviewItem{}, sourceErr
			}
			if _, _, sourceErr := selectRowGroupMigrationSource(record, allSources, pointID); sourceErr != nil {
				return WriteGroupMigrationPreviewItem{}, sourceErr
			}
			addIssue("row-group-target-missing-or-disabled", "the row-group member has no enabled target mapping in the saved connector scope")
			continue
		}
		if len(enabledTargets) > 1 {
			addIssue("row-group-multiple-target-mappings", "the row-group member resolves to more than one enabled target mapping")
			continue
		}
		target := enabledTargets[0]
		allTagTargets, err := readMigrationPreviewTargetsForTag(ctx, tx, repo, target.tagID)
		if err != nil {
			return WriteGroupMigrationPreviewItem{}, err
		}
		for _, otherTarget := range allTagTargets {
			if !otherTarget.enabled || otherTarget.connectorID == connector.id {
				continue
			}
			addIssue("row-group-multiple-output-connectors", "the selected tag has another enabled target mapping in a different connector")
			break
		}
		allSources, err := readMigrationPreviewSources(ctx, tx, repo, target.tagID)
		if err != nil {
			return WriteGroupMigrationPreviewItem{}, err
		}
		source, sourceIssues, sourceErr := selectRowGroupMigrationSource(record, allSources, pointID)
		if sourceErr != nil {
			return WriteGroupMigrationPreviewItem{}, sourceErr
		}
		item.Issues = append(item.Issues, sourceIssues...)
		member := writeGroupRowGroupMember{pointID: pointID, target: target, source: source}
		members = append(members, member)
		if source != nil {
			intent := migrationPreviewIntent(source, target.writeGroupMigrationTarget, connector)
			intent.Database = destinationDatabaseForRowGroup(target, connector)
			item.BeforeRowGroupIntent.Members = append(item.BeforeRowGroupIntent.Members, *intent)
		} else {
			item.BeforeRowGroupIntent.Members = append(item.BeforeRowGroupIntent.Members, WriteGroupMigrationIntent{
				SourceID:             target.id,
				PointID:              pointID,
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
				IntervalSource:       migrationPreviewIntervalSource(target.writeGroupMigrationTarget, connector),
				Enabled:              target.enabled,
			})
		}
	}

	if len(members) == 0 && len(rowGroup.MemberPointIDs) > 0 {
		addIssue("row-group-members-unresolved", "none of the saved row-group members has a reviewable target mapping")
	}
	validateRowGroupTargetLayout(rowGroup, connector, members, addIssue)
	if len(item.Issues) == 0 {
		candidate, sourceRevision, err := rowGroupMigrationCandidate(record, rowGroup, connector, members)
		if err != nil {
			return WriteGroupMigrationPreviewItem{}, err
		}
		ownershipIssues, ownershipErr := migrationSourceOwnershipIssues(ctx, tx, repo, record.ID, candidate.Migration.SourceKind, rowGroupID, candidate.Migration)
		if ownershipErr != nil {
			return WriteGroupMigrationPreviewItem{}, ownershipErr
		}
		item.Issues = append(item.Issues, ownershipIssues...)
		if len(item.Issues) == 0 {
			item.SourceRevision = sourceRevision
			item.CandidateGroup = candidate
		}
	}
	if len(item.Issues) == 0 {
		item.Status = WriteGroupMigrationPreviewStatusNeedsReview
	} else {
		item.Status = WriteGroupMigrationPreviewStatusBlocked
	}
	return item, nil
}

func findDatabaseRowGroup(groups []DatabaseRowGroup, id string) (DatabaseRowGroup, bool) {
	for _, group := range groups {
		if group.ID == id {
			return group, true
		}
	}
	return DatabaseRowGroup{}, false
}

func cloneDatabaseRowGroupValue(group DatabaseRowGroup) DatabaseRowGroup {
	group.MemberPointIDs = slices.Clone(group.MemberPointIDs)
	group.GroupKeyColumns = slices.Clone(group.GroupKeyColumns)
	group.UniqueKeyColumns = slices.Clone(group.UniqueKeyColumns)
	return group
}

func destinationDatabaseForRowGroup(target writeGroupRowGroupTarget, connector writeGroupMigrationConnector) string {
	destination, err := migrationPreviewDestination(target.writeGroupMigrationTarget, connector)
	if err != nil {
		return ""
	}
	return destination.Database
}

func rowGroupMigrationCandidate(
	record *Record,
	rowGroup DatabaseRowGroup,
	connector writeGroupMigrationConnector,
	members []writeGroupRowGroupMember,
) (*WriteGroup, string, error) {
	if len(members) == 0 {
		return nil, "", fmt.Errorf("row-group migration candidate has no members: %w", ErrWriteGroupValidation)
	}
	first := members[0].target
	destination, err := migrationPreviewDestination(first.writeGroupMigrationTarget, connector)
	if err != nil {
		return nil, "", err
	}
	interval := migrationPreviewInterval(first.writeGroupMigrationTarget, connector)
	candidate := &WriteGroup{
		WorkspaceID: record.ID,
		Name:        "row group " + rowGroup.ID,
		Status:      WriteGroupStatusDraft,
		Members:     make([]WriteGroupMember, 0, len(members)),
		Destination: destination,
		RowPolicy: WriteGroupRowPolicy{
			IntervalSeconds:        interval,
			AllowedLatenessSeconds: 0,
			IncompletePolicy:       writeGroupIncompletePolicySkipRow,
			GroupKeyColumns:        slices.Clone(rowGroup.GroupKeyColumns),
			UniqueKeyColumns:       slices.Clone(rowGroup.UniqueKeyColumns),
		},
		WritePolicy: WriteGroupWritePolicy{Mode: writeGroupWriteModeAppend},
		Migration: WriteGroupMigration{
			SourceKind:          writeGroupMigrationSourceKindLegacyRowGroup,
			LegacyRowGroupID:    rowGroup.ID,
			SourceIDs:           make([]string, 0, len(members)),
			TargetMappingPoints: make(map[string]string, len(members)),
			AdapterVersion:      writeGroupRowGroupMigrationAdapterVersion,
			ReviewResult:        WriteGroupMigrationPreviewStatusNeedsReview,
		},
	}
	digestParts := []string{"row-group-preview", rowGroup.ID, rowGroup.ConnectorID, rowGroup.TableSchema, rowGroup.TableName, connector.id, connector.identityRevision, connector.kind, connector.connectionConfig, fmt.Sprint(connector.defaultWriteIntervalSecs)}
	for _, pointID := range rowGroup.MemberPointIDs {
		digestParts = append(digestParts, "member-point", pointID)
	}
	for _, column := range rowGroup.GroupKeyColumns {
		digestParts = append(digestParts, "group-key-column", column)
	}
	for _, column := range rowGroup.UniqueKeyColumns {
		digestParts = append(digestParts, "unique-key-column", column)
	}
	refParts := make([]string, 0, len(record.DatabaseTargetRefs))
	for _, ref := range record.DatabaseTargetRefs {
		if ref.RowGroupID != rowGroup.ID {
			continue
		}
		refParts = append(refParts, ref.PointID+"\x00"+ref.RowGroupID)
	}
	slices.Sort(refParts)
	for _, ref := range refParts {
		digestParts = append(digestParts, "target-ref", ref)
	}
	for _, member := range members {
		target := member.target
		source := member.source
		if source == nil {
			return nil, "", fmt.Errorf("row-group migration candidate source missing: %w", ErrWriteGroupValidation)
		}
		mappingRevision, err := source.mapping.revision()
		if err != nil {
			return nil, "", fmt.Errorf("derive row-group mapping revision: %w", err)
		}
		sourceRevision, err := source.source.revision()
		if err != nil {
			return nil, "", fmt.Errorf("derive row-group source revision: %w", err)
		}
		digestParts = append(digestParts, target.id, target.pointID, target.tagID, target.connectorID, target.tableSchema, target.tableName, target.columnName, target.writeMode, migrationPreviewOptionalString(target.timestampColumn), migrationPreviewOptionalString(target.groupKey), migrationPreviewOptionalInt(target.writeIntervalSeconds), fmt.Sprint(target.enabled), mappingRevision, sourceRevision)
		memberRecord := WriteGroupMember{
			DeviceID:        source.source.DeviceID,
			PointID:         source.source.PointID,
			TagID:           source.source.TagID,
			EntityKey:       migrationPreviewOptionalString(target.groupKey),
			SourceRevision:  sourceRevision,
			MappingRevision: mappingRevision,
			TargetColumn:    target.columnName,
			Required:        true,
			MaxAgeSeconds:   common.Ptr(interval),
		}
		candidate.Members = append(candidate.Members, memberRecord)
		candidate.Migration.SourceIDs = append(candidate.Migration.SourceIDs, target.id)
		candidate.Migration.TargetMappingPoints[target.id] = target.pointID
	}
	slices.Sort(candidate.Migration.SourceIDs)
	sourceDigest, err := hashWriteGroupRevision(digestParts...)
	if err != nil {
		return nil, "", fmt.Errorf("derive row-group migration source digest: %w", err)
	}
	candidate.Migration.SourceRevision = sourceDigest
	return candidate, sourceDigest, nil
}

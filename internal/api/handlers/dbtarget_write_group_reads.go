package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

var errLegacyMappingReadConflict = errors.New("canonical group requires the write-group read interface")

const (
	legacyReadSourceKindSingleMapping = "legacy-single-mapping"
	legacyReadSourceKindRowGroup      = "legacy-row-group"
)

// writeGroupMigrationReader is used only for operator HTTP reads. Runtime
// continues reading the original mapping service until validated activation.
type writeGroupMigrationReader interface {
	List(context.Context) (*workspace.WriteGroupListResult, error)
}

// writeGroupLegacyRowGroupGuard is deliberately narrower than the full write
// group service so read-only handler test doubles remain valid.
type writeGroupLegacyRowGroupGuard interface {
	PreflightLegacyRowGroupReplacement(context.Context, string, []workspace.DatabaseRowGroup) error
	CheckLegacyRowGroupReplacementInTx(context.Context, *sql.Tx, *workspace.Record, string, []workspace.DatabaseRowGroup) error
}

type databaseTargetMappingReadResponse struct {
	*schema.DatabaseTargetMapping
	CanonicalGroup *workspace.WriteGroup         `json:"canonical_group,omitempty"`
	LegacyIntent   *schema.DatabaseTargetMapping `json:"legacy_intent,omitempty"`
	PointID        string                        `json:"-"`
	RowGroupID     string                        `json:"-"`
}

// WithWriteGroups enables canonical projections for migrated operator reads.
func (h *DatabaseTargetHandler) WithWriteGroups(groups *workspace.WriteGroupService) *DatabaseTargetHandler {
	if groups != nil {
		h.writeGroups = groups
		if h.connectorSvc != nil {
			h.connectorSvc.SetDeleteGuard(groups)
		}
		if h.mappingSvc != nil {
			dbtarget.WithLegacyWriteCoordinator(h.mappingSvc, groups)
		}
	}
	return h
}

// WithWriteGroups enables canonical projections on legacy workspace reads.
func (h *StudioV2WorkspaceDatabaseHandler) WithWriteGroups(groups *workspace.WriteGroupService) *StudioV2WorkspaceDatabaseHandler {
	if groups != nil {
		h.writeGroups = groups
		if h.mappingSvc != nil {
			dbtarget.WithLegacyWriteCoordinator(h.mappingSvc, groups)
		}
	}
	return h
}

func migrationReadGroups(ctx context.Context, reader writeGroupMigrationReader) (map[string]*workspace.WriteGroup, error) {
	groups := make(map[string]*workspace.WriteGroup)
	if reader == nil {
		return groups, nil
	}
	result, err := reader.List(ctx)
	if errors.Is(err, workspace.ErrNotFound) {
		return groups, nil
	}
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, workspace.ErrWriteGroupServiceUnavailable
	}
	for _, group := range result.Groups {
		if group == nil || group.WorkspaceID != result.WorkspaceID {
			return nil, workspace.ErrWriteGroupServiceUnavailable
		}
		if group.Migration.SourceKind != legacyReadSourceKindSingleMapping && group.Migration.SourceKind != legacyReadSourceKindRowGroup {
			continue
		}
		if len(group.Migration.SourceIDs) == 0 {
			return nil, fmt.Errorf("legacy mapping provenance is missing: %w", errLegacyMappingReadConflict)
		}
		for _, id := range group.Migration.SourceIDs {
			if id == "" || groups[id] != nil {
				return nil, fmt.Errorf("legacy mapping has ambiguous canonical ownership: %w", errLegacyMappingReadConflict)
			}
			groups[id] = group
		}
	}
	return groups, nil
}

func projectMigrationRead(row *schema.DatabaseTargetMapping, groups map[string]*workspace.WriteGroup) (databaseTargetMappingReadResponse, error) {
	if row == nil {
		return databaseTargetMappingReadResponse{}, workspace.ErrWriteGroupServiceUnavailable
	}
	group := groups[row.ID]
	if group == nil {
		return databaseTargetMappingReadResponse{DatabaseTargetMapping: row}, nil
	}
	member, rowGroupID, err := migrationReadMember(row, group)
	if err != nil {
		return databaseTargetMappingReadResponse{}, err
	}
	projected := *row
	projected.TagID = member.TagID
	projected.ConnectorID = group.Destination.ConnectorID
	projected.TableSchema = group.Destination.TableSchema
	projected.TableName = group.Destination.TableName
	projected.ColumnName = member.TargetColumn
	projected.WriteMode = schema.DatabaseWriteModeInsert
	// The legacy timestamp field is historical intent, not a snapshot binding.
	projected.TimestampColumn = nil
	projected.GroupKey = nil
	if member.EntityKey != "" {
		entityKey := member.EntityKey
		projected.GroupKey = &entityKey
	}
	interval := group.RowPolicy.IntervalSeconds
	projected.WriteIntervalSeconds = &interval
	projected.Enabled = group.Status != workspace.WriteGroupStatusDisabled && group.Status != workspace.WriteGroupStatusDeleted
	projected.UpdatedAt = group.UpdatedAt
	return databaseTargetMappingReadResponse{
		DatabaseTargetMapping: &projected, CanonicalGroup: group, LegacyIntent: row,
		PointID: member.PointID, RowGroupID: rowGroupID,
	}, nil
}

func migrationReadMember(row *schema.DatabaseTargetMapping, group *workspace.WriteGroup) (*workspace.WriteGroupMember, string, error) {
	conflict := func() (*workspace.WriteGroupMember, string, error) {
		return nil, "", fmt.Errorf("canonical group cannot be represented by the legacy mapping interface: %w", errLegacyMappingReadConflict)
	}
	if group.Destination.ConnectorID != row.ConnectorID || group.WritePolicy.Mode != "append" {
		return conflict()
	}
	if group.Migration.SourceKind == legacyReadSourceKindSingleMapping {
		if len(group.Members) != 1 || len(group.Migration.SourceIDs) != 1 {
			return conflict()
		}
		return &group.Members[0], group.ID, nil
	}
	if group.Migration.SourceKind != legacyReadSourceKindRowGroup || group.Migration.LegacyRowGroupID == "" ||
		len(group.Migration.TargetMappingPoints) != len(group.Members) || len(group.Migration.SourceIDs) != len(group.Members) ||
		group.RowPolicy.EntityKeyColumn != "" || group.RowPolicy.ValueColumn != "" ||
		group.RowPolicy.QualityColumn != "" || group.RowPolicy.ProvenanceColumn != "" {
		return conflict()
	}
	seenPoints := make(map[string]bool, len(group.Members))
	for _, sourceID := range group.Migration.SourceIDs {
		pointID := group.Migration.TargetMappingPoints[sourceID]
		if pointID == "" || seenPoints[pointID] {
			return conflict()
		}
		seenPoints[pointID] = true
		matches := 0
		for _, member := range group.Members {
			if member.PointID == pointID && member.EntityKey != "" {
				matches++
			}
		}
		if matches != 1 {
			return conflict()
		}
	}
	pointID := group.Migration.TargetMappingPoints[row.ID]
	for i := range group.Members {
		if group.Members[i].PointID == pointID {
			return &group.Members[i], group.Migration.LegacyRowGroupID, nil
		}
	}
	return conflict()
}

func validateMigrationReadRows(rows []*schema.DatabaseTargetMapping, groups map[string]*workspace.WriteGroup, filter dbtarget.TargetMappingListFilter) error {
	seen := make(map[string]bool, len(rows))
	for _, row := range rows {
		if row == nil {
			return workspace.ErrWriteGroupServiceUnavailable
		}
		seen[row.ID] = true
	}
	for id, group := range groups {
		if seen[id] || !missingMigrationTargetMatchesFilter(id, group, filter) {
			continue
		}
		return fmt.Errorf("canonical source target is missing from the legacy representation: %w", errLegacyMappingReadConflict)
	}
	return nil
}

func missingMigrationTargetMatchesFilter(id string, group *workspace.WriteGroup, filter dbtarget.TargetMappingListFilter) bool {
	row := &schema.DatabaseTargetMapping{
		ConnectorID: group.Destination.ConnectorID,
		Enabled:     group.Status != workspace.WriteGroupStatusDisabled && group.Status != workspace.WriteGroupStatusDeleted,
	}
	for _, member := range group.Members {
		if (group.Migration.SourceKind == legacyReadSourceKindSingleMapping && len(group.Members) == 1) ||
			member.PointID == group.Migration.TargetMappingPoints[id] {
			row.TagID = member.TagID
			break
		}
	}
	if row.TagID == "" {
		return migrationReadCandidateMatchesFilter(row, group, filter)
	}
	return migrationReadMatchesFilter(row, filter)
}

func renderMigrationReadError(c *gin.Context, err error) {
	if errors.Is(err, errLegacyMappingReadConflict) {
		renderWriteGroupTyped(c, http.StatusConflict, "WRITE_GROUP_LEGACY_READ_CONFLICT", "saved group requires the canonical write-group interface", false, "open_write_groups")
		return
	}
	renderWriteGroupError(c, err)
}

func migrationReadMatchesFilter(row *schema.DatabaseTargetMapping, filter dbtarget.TargetMappingListFilter) bool {
	return (filter.ConnectorID == nil || row.ConnectorID == *filter.ConnectorID) &&
		(filter.TagID == nil || row.TagID == *filter.TagID) &&
		(filter.Enabled == nil || row.Enabled == *filter.Enabled)
}

func migrationReadCandidateMatchesFilter(row *schema.DatabaseTargetMapping, group *workspace.WriteGroup, filter dbtarget.TargetMappingListFilter) bool {
	if group == nil {
		return migrationReadMatchesFilter(row, filter)
	}
	enabled := group.Status != workspace.WriteGroupStatusDisabled && group.Status != workspace.WriteGroupStatusDeleted
	connectorMatches := filter.ConnectorID == nil || group.Destination.ConnectorID == *filter.ConnectorID
	if group.Destination.ConnectorID != row.ConnectorID && filter.ConnectorID != nil {
		connectorMatches = connectorMatches || row.ConnectorID == *filter.ConnectorID
	}
	return connectorMatches &&
		(filter.Enabled == nil || enabled == *filter.Enabled) &&
		(filter.TagID == nil || slices.ContainsFunc(group.Members, func(member workspace.WriteGroupMember) bool {
			return member.TagID == *filter.TagID
		}))
}

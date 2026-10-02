package workspace

import (
	"context"

	"go-gateway/internal/datalink/recordingplan"
)

const (
	WriteGroupMigrationPreviewStatusNeedsReview    = "needs_review"
	WriteGroupMigrationPreviewStatusBlocked        = "blocked"
	writeGroupMigrationPreviewAdapterVersion       = "single-mapping-v1"
	writeGroupIncompletePolicySkipRow              = "skip_row"
	writeGroupWriteModeAppend                      = "append"
	writeGroupMigrationSourceKindLegacyRowGroup    = "legacy-row-group"
	writeGroupRecordingPlanMigrationAdapterVersion = "recording-plan-v1"
)

// WriteGroupMigrationPreview is a read-only compatibility review for legacy
// database target mappings. It never creates a canonical group or probes the
// external destination.
type WriteGroupMigrationPreview struct {
	WorkspaceID       string                           `json:"workspace_id"`
	WorkspaceRevision string                           `json:"workspace_revision"`
	ConnectorRevision string                           `json:"connector_revision"`
	AdapterVersion    string                           `json:"adapter_version"`
	ReviewDigest      string                           `json:"review_digest"`
	Items             []WriteGroupMigrationPreviewItem `json:"items"`
}

// WriteGroupMigrationPreviewItem describes one legacy mapping and its safe
// basic-group candidate, if the mapping is representable for review.
type WriteGroupMigrationPreviewItem struct {
	SourceID                  string                                  `json:"source_id"`
	SourceRevision            string                                  `json:"source_revision"`
	Status                    string                                  `json:"status"`
	Differences               []WriteGroupMigrationFinding            `json:"differences"`
	Issues                    []WriteGroupMigrationFinding            `json:"issues"`
	BeforeIntent              *WriteGroupMigrationIntent              `json:"before_intent,omitempty"`
	BeforeRowGroupIntent      *WriteGroupRowGroupMigrationIntent      `json:"before_row_group_intent,omitempty"`
	BeforeRecordingPlanIntent *WriteGroupRecordingPlanMigrationIntent `json:"before_recording_plan_intent,omitempty"`
	RepairAction              string                                  `json:"repair_action,omitempty"`
	CandidateGroup            *WriteGroup                             `json:"candidate_group,omitempty"`
}

// WriteGroupMigrationFinding is an explicit semantic difference or blocking
// reason; it deliberately excludes connector credentials and DSNs.
type WriteGroupMigrationFinding struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteGroupMigrationIntent captures the persisted legacy write intent
// without copying secrets from the connector configuration.
type WriteGroupMigrationIntent struct {
	SourceID             string  `json:"source_id"`
	DeviceID             string  `json:"device_id"`
	PointID              string  `json:"point_id"`
	TagID                string  `json:"tag_id"`
	ConnectorID          string  `json:"connector_id"`
	ConnectorRevision    string  `json:"connector_revision"`
	Database             string  `json:"database,omitempty"`
	TableSchema          string  `json:"table_schema"`
	TableName            string  `json:"table_name"`
	ColumnName           string  `json:"column_name"`
	WriteMode            string  `json:"write_mode"`
	TimestampColumn      *string `json:"timestamp_column,omitempty"`
	GroupKey             *string `json:"group_key,omitempty"`
	WriteIntervalSeconds *int    `json:"write_interval_seconds,omitempty"`
	IntervalSource       string  `json:"interval_source,omitempty"`
	Enabled              bool    `json:"enabled"`
}

// WriteGroupRowGroupMigrationIntent preserves the complete workspace row
// group and each target mapping intent before a reviewed conversion.
type WriteGroupRowGroupMigrationIntent struct {
	SourceID          string                      `json:"source_id"`
	ConnectorID       string                      `json:"connector_id"`
	ConnectorRevision string                      `json:"connector_revision"`
	RowGroup          DatabaseRowGroup            `json:"row_group"`
	Members           []WriteGroupMigrationIntent `json:"members"`
}

// WriteGroupRecordingPlanMigrationSource preserves the source-chain snapshot
// referenced by a recording plan without exposing device connection details.
type WriteGroupRecordingPlanMigrationSource struct {
	MeasurementID         string `json:"measurement_id"`
	DeviceID              string `json:"device_id"`
	PointID               string `json:"point_id"`
	TagID                 string `json:"tag_id"`
	DefinitionRevision    string `json:"definition_revision"`
	SourceBindingRevision string `json:"source_binding_revision"`
	SeriesEpoch           string `json:"series_epoch"`
	SourceRevision        string `json:"source_revision"`
	MappingRevision       string `json:"mapping_revision,omitempty"`
	Status                string `json:"status"`
}

// WriteGroupRecordingPlanMigrationIntent preserves the complete recording
// plan and its referenced source snapshots before a blocked migration review.
type WriteGroupRecordingPlanMigrationIntent struct {
	SourceID string                                   `json:"source_id"`
	Plan     recordingplan.RecordingPlan              `json:"plan"`
	Sources  []WriteGroupRecordingPlanMigrationSource `json:"sources"`
}

// PreviewSingleMappingMigration is the read-only migration seam. It reads the
// local workspace, legacy target mappings, and source chain in one transaction;
// it does not save a group or probe the external destination.
func (s *WriteGroupService) PreviewSingleMappingMigration(
	ctx context.Context,
	workspaceID string,
	sourceIDs []string,
) (*WriteGroupMigrationPreview, error) {
	return s.previewSingleMappingMigration(ctx, workspaceID, sourceIDs)
}

// PreviewRowGroupMigration is the read-only seam for workspace row-group
// compatibility review. It never changes canonical groups or legacy targets.
func (s *WriteGroupService) PreviewRowGroupMigration(
	ctx context.Context,
	workspaceID string,
	sourceIDs []string,
) (*WriteGroupMigrationPreview, error) {
	return s.previewRowGroupMigration(ctx, workspaceID, sourceIDs)
}

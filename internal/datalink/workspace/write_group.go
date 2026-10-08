package workspace

import (
	"maps"
	"slices"
	"time"

	"go-gateway/internal/datalink/common"
)

// WriteGroupStatus is the local lifecycle of a persisted Step 4 group.
type WriteGroupStatus string

const (
	WriteGroupStatusDraft    WriteGroupStatus = "draft"
	WriteGroupStatusReady    WriteGroupStatus = "ready"
	WriteGroupStatusRunning  WriteGroupStatus = "running"
	WriteGroupStatusDisabled WriteGroupStatus = "disabled"
	WriteGroupStatusDeleted  WriteGroupStatus = "deleted"
)

// WriteGroupStorageStrategy identifies the layout strategy of one group.
type WriteGroupStorageStrategy string

const (
	WriteGroupStorageStrategyManaged WriteGroupStorageStrategy = "managed"
	WriteGroupStorageStrategyCustom  WriteGroupStorageStrategy = "custom"
)

// WriteGroup is the persisted authority for one database output group.
type WriteGroup struct {
	ID              string `json:"id"`
	WorkspaceID     string `json:"workspace_id"`
	Revision        string `json:"revision"`
	AppliedRevision string `json:"applied_revision"`
	// BasicManagedDeviceID is read from the create-once key, never caller-assigned.
	BasicManagedDeviceID string                `json:"basic_managed_device_id,omitempty"`
	Name                 string                `json:"name"`
	Status               WriteGroupStatus      `json:"status"`
	Members              []WriteGroupMember    `json:"members"`
	Destination          WriteGroupDestination `json:"destination"`
	RowPolicy            WriteGroupRowPolicy   `json:"row_policy"`
	WritePolicy          WriteGroupWritePolicy `json:"write_policy"`
	Migration            WriteGroupMigration   `json:"migration"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

// WriteGroupMember stores persisted source and mapping identities.
type WriteGroupMember struct {
	DeviceID string `json:"device_id"`
	PointID  string `json:"point_id"`
	TagID    string `json:"tag_id"`
	// EntityKey is the raw legacy group key used only to partition buffered
	// rows. It is deliberately not interpreted as a destination SQL value.
	EntityKey       string  `json:"entity_key,omitempty"`
	SourceRevision  string  `json:"source_revision,omitempty"`
	MappingRevision string  `json:"mapping_revision,omitempty"`
	MeasurementID   *string `json:"measurement_id,omitempty"`
	TargetColumn    string  `json:"target_column"`
	Required        bool    `json:"required"`
	MaxAgeSeconds   *int    `json:"max_age_seconds,omitempty"`
}

// WriteGroupDestination binds one group to a saved connector identity.
type WriteGroupDestination struct {
	ConnectorID       string                    `json:"connector_id"`
	ConnectorRevision string                    `json:"connector_revision"`
	Database          string                    `json:"database"`
	TableSchema       string                    `json:"table_schema"`
	TableName         string                    `json:"table_name"`
	StorageStrategy   WriteGroupStorageStrategy `json:"storage_strategy"`
	SchemaRevision    string                    `json:"schema_revision,omitempty"`
	SchemaDigest      string                    `json:"schema_digest,omitempty"`
}

// WriteGroupRowPolicy reserves the basic snapshot row contract for later tasks.
type WriteGroupRowPolicy struct {
	IntervalSeconds        int      `json:"interval_seconds"`
	AllowedLatenessSeconds int      `json:"allowed_lateness_seconds"`
	IncompletePolicy       string   `json:"incomplete_policy,omitempty"`
	EntityKeyColumn        string   `json:"entity_key_column,omitempty"`
	GroupKeyColumns        []string `json:"group_key_columns,omitempty"`
	UniqueKeyColumns       []string `json:"unique_key_columns,omitempty"`
	ValueColumn            string   `json:"value_column,omitempty"`
	QualityColumn          string   `json:"quality_column,omitempty"`
	ProvenanceColumn       string   `json:"provenance_column,omitempty"`
	RecordKeyColumn        string   `json:"record_key_column,omitempty"`
	BucketStartColumn      string   `json:"bucket_start_column,omitempty"`
	GroupIDColumn          string   `json:"group_id_column,omitempty"`
	DeviceIDColumn         string   `json:"device_id_column,omitempty"`
}

// WriteGroupWritePolicy reserves append/latest and dedupe semantics.
type WriteGroupWritePolicy struct {
	Mode             string `json:"mode,omitempty"`
	DedupeCapability string `json:"dedupe_capability,omitempty"`
}

// WriteGroupMigration records compatibility provenance without activating it.
type WriteGroupMigration struct {
	SourceKind          string            `json:"source_kind,omitempty"`
	SourceIDs           []string          `json:"source_ids,omitempty"`
	LegacyRowGroupID    string            `json:"legacy_row_group_id,omitempty"`
	TargetMappingPoints map[string]string `json:"target_mapping_points,omitempty"`
	SourceRevision      string            `json:"source_revision,omitempty"`
	AdapterVersion      string            `json:"adapter_version,omitempty"`
	ReviewResult        string            `json:"review_result,omitempty"`
}

func cloneWriteGroup(group *WriteGroup) *WriteGroup {
	if group == nil {
		return nil
	}
	cloned := *group
	cloned.Members = slices.Clone(group.Members)
	for i, member := range group.Members {
		if member.MeasurementID != nil {
			cloned.Members[i].MeasurementID = common.Ptr(*member.MeasurementID)
		}
		if member.MaxAgeSeconds != nil {
			cloned.Members[i].MaxAgeSeconds = common.Ptr(*member.MaxAgeSeconds)
		}
	}
	cloned.Migration.SourceIDs = slices.Clone(group.Migration.SourceIDs)
	cloned.Migration.TargetMappingPoints = maps.Clone(group.Migration.TargetMappingPoints)
	cloned.RowPolicy.GroupKeyColumns = append([]string(nil), group.RowPolicy.GroupKeyColumns...)
	cloned.RowPolicy.UniqueKeyColumns = append([]string(nil), group.RowPolicy.UniqueKeyColumns...)
	return &cloned
}

package schema

import "time"

// SourceRule 來源規則，代表使用者在 Step 2 建立的持久化來源規則。
type SourceRule struct {
	ID               string    `json:"id" db:"id"`
	DeviceID         string    `json:"device_id" db:"device_id"`
	StartAddress     string    `json:"start_address" db:"start_address"`
	Count            int       `json:"count" db:"count"`
	DataType         DataType  `json:"data_type" db:"data_type"`
	NamingPrefix     string    `json:"naming_prefix" db:"naming_prefix"`
	Enabled          bool      `json:"enabled" db:"enabled"`
	Locked           bool      `json:"locked" db:"locked"`
	Origin           string    `json:"origin" db:"origin"`
	TemplateName     string    `json:"template_name,omitempty" db:"template_name"`
	SkippedAddresses string    `json:"skipped_addresses,omitempty" db:"skipped_addresses"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	RevisionID       string    `json:"revision_id" db:"revision_id"`
	// TargetDataType 目標資料型態（可空）。當與 DataType 不同時，自動建立 cast 步驟。
	// NULL 表示與 Point 讀取型別相同（即 DataType）。
	TargetDataType *DataType `json:"target_data_type,omitempty" db:"target_data_type"`
	// ScaleMultiplier 縮放倍率（可空）。與 ScaleOffset 構成線性縮放公式：value * multiplier + offset。
	// NULL 表示不進行縮放。
	ScaleMultiplier *float64 `json:"scale_multiplier,omitempty" db:"scale_multiplier"`
	// ScaleOffset 偏移量（可空）。與 ScaleMultiplier 構成線性縮放。
	// NULL 表示偏移量為 0（若 ScaleMultiplier 非 NULL）。
	ScaleOffset *float64 `json:"scale_offset,omitempty" db:"scale_offset"`
	// DataFormat 多暫存器／浮點解碼字節序（ABCD、BADC、CDAB、DCBA）；空字串表示使用設備連線預設。
	DataFormat string `json:"data_format,omitempty" db:"data_format"`
	// ShareEnabled controls whether this rule is exposed through the local Modbus Share output.
	ShareEnabled bool `json:"share_enabled" db:"share_enabled"`
	// ShareStartRegister is the optional local Modbus holding-register start address.
	ShareStartRegister *int `json:"share_start_register" db:"share_start_register"`
	// ShareStride is the optional number of local Modbus registers allocated per point.
	ShareStride *int `json:"share_stride" db:"share_stride"`
}

// SourceRuleLink 記錄來源規則衍生出的下游實體關聯。
type SourceRuleLink struct {
	ID        string    `json:"id" db:"id"`
	RuleID    string    `json:"rule_id" db:"rule_id"`
	Address   string    `json:"address" db:"address"`
	PointID   string    `json:"point_id" db:"point_id"`
	TagID     *string   `json:"tag_id,omitempty" db:"tag_id"`
	MappingID *string   `json:"mapping_id,omitempty" db:"mapping_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type SourceRuleCandidateType string

const (
	SourceRuleCandidateTypeTags               SourceRuleCandidateType = "tags"
	SourceRuleCandidateTypeMappings           SourceRuleCandidateType = "mappings"
	SourceRuleCandidateTypeDatabaseOutputs    SourceRuleCandidateType = "database_outputs"
	SourceRuleCandidateTypeLocalModbusOutputs SourceRuleCandidateType = "local_modbus_outputs"
)

type SourceRuleCandidateStatus string

const (
	SourceRuleCandidateStatusReady    SourceRuleCandidateStatus = "ready"
	SourceRuleCandidateStatusBlocked  SourceRuleCandidateStatus = "blocked"
	SourceRuleCandidateStatusDeferred SourceRuleCandidateStatus = "deferred"
)

// SourceRuleCandidateKind identifies one downstream object kind within a candidate type.
type SourceRuleCandidateKind string

const (
	// SourceRuleCandidateKindTag identifies a rule-derived tag candidate and its pending mapping intent.
	SourceRuleCandidateKindTag               SourceRuleCandidateKind = "tag"
	SourceRuleCandidateKindMapping           SourceRuleCandidateKind = "mapping"
	SourceRuleCandidateKindDatabaseOutput    SourceRuleCandidateKind = "database_output"
	SourceRuleCandidateKindLocalModbusOutput SourceRuleCandidateKind = "local_modbus_output"
)

// SourceRuleCandidateScopeField records one canonical binding-scope discriminator.
type SourceRuleCandidateScopeField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SourceRuleCandidateIdentity is the canonical identity contract for one rule-derived candidate.
type SourceRuleCandidateIdentity struct {
	SourceRuleID           string                          `json:"source_rule_id"`
	CandidateType          SourceRuleCandidateType         `json:"candidate_type"`
	CandidateKind          SourceRuleCandidateKind         `json:"candidate_kind"`
	DerivedFromRuleAddress string                          `json:"derived_from_rule_address"`
	TargetBindingScope     []SourceRuleCandidateScopeField `json:"target_binding_scope,omitempty"`
}

// SourceRuleCandidateSnapshot persists one candidate set for a rule revision.
type SourceRuleCandidateSnapshot struct {
	SourceRuleID  string                    `json:"source_rule_id" db:"source_rule_id"`
	RevisionID    string                    `json:"revision_id" db:"revision_id"`
	CandidateType SourceRuleCandidateType   `json:"candidate_type" db:"candidate_type"`
	Payload       string                    `json:"payload" db:"payload"`
	Status        SourceRuleCandidateStatus `json:"status" db:"status"`
	Reason        string                    `json:"reason,omitempty" db:"reason"`
	GeneratedAt   time.Time                 `json:"generated_at" db:"generated_at"`
}

// SourceRuleTagCandidate describes one rule-derived tag candidate snapshot entry.
type SourceRuleTagCandidate struct {
	ID                string                      `json:"id"`
	Identity          SourceRuleCandidateIdentity `json:"identity"`
	ProposedSignature string                      `json:"proposed_signature"`
	Address           string                      `json:"address"`
	PointID           string                      `json:"point_id"`
	TagID             *string                     `json:"tag_id,omitempty"`
	MappingID         *string                     `json:"mapping_id,omitempty"`
	TagKey            string                      `json:"tag_key,omitempty"`
	DisplayName       string                      `json:"display_name,omitempty"`
	DataType          DataType                    `json:"data_type"`
	TransformPipeline []TransformStep             `json:"transform_pipeline,omitempty"`
}

// SourceRuleDatabaseOutputCandidate describes one rule-scoped database output candidate.
type SourceRuleDatabaseOutputCandidate struct {
	ID                   string                      `json:"id"`
	Identity             SourceRuleCandidateIdentity `json:"identity"`
	ProposedSignature    string                      `json:"proposed_signature"`
	Address              string                      `json:"address"`
	PointID              string                      `json:"point_id"`
	TagID                *string                     `json:"tag_id,omitempty"`
	TagKey               string                      `json:"tag_key,omitempty"`
	DisplayName          string                      `json:"display_name,omitempty"`
	DataType             DataType                    `json:"data_type"`
	MappingID            *string                     `json:"mapping_id,omitempty"`
	ConnectorID          string                      `json:"connector_id,omitempty"`
	TableSchema          string                      `json:"table_schema,omitempty"`
	TableName            string                      `json:"table_name,omitempty"`
	ColumnName           string                      `json:"column_name,omitempty"`
	GroupKey             *string                     `json:"group_key"`
	WriteMode            DatabaseWriteMode           `json:"write_mode,omitempty"`
	TimestampColumn      *string                     `json:"timestamp_column,omitempty"`
	WriteIntervalSeconds *int                        `json:"write_interval_seconds"`
	Status               SourceRuleOutputStatus      `json:"status"`
	BlockingReason       string                      `json:"blocking_reason,omitempty"`
}

// SourceRuleOutputStatus identifies review/apply readiness of a rule-owned output candidate.
type SourceRuleOutputStatus string

const (
	SourceRuleOutputStatusReady     SourceRuleOutputStatus = "ready"
	SourceRuleOutputStatusBlocked   SourceRuleOutputStatus = "blocked"
	SourceRuleOutputStatusOutOfSync SourceRuleOutputStatus = "out_of_sync"
)

// SourceRuleLocalModbusOutputStatus identifies review/apply readiness of one rule-owned local Modbus candidate.
type SourceRuleLocalModbusOutputStatus string

const (
	SourceRuleLocalModbusOutputStatusDeferred        SourceRuleLocalModbusOutputStatus = "deferred"
	SourceRuleLocalModbusOutputStatusReady           SourceRuleLocalModbusOutputStatus = "ready"
	SourceRuleLocalModbusOutputStatusOutOfSync       SourceRuleLocalModbusOutputStatus = "out_of_sync"
	SourceRuleLocalModbusOutputStatusBlockedConflict SourceRuleLocalModbusOutputStatus = "blocked_conflict"
)

// SourceRuleLocalModbusOutputCandidate describes one rule-scoped local Modbus output candidate.
type SourceRuleLocalModbusOutputCandidate struct {
	ID                string                            `json:"id"`
	Identity          SourceRuleCandidateIdentity       `json:"identity"`
	ProposedSignature string                            `json:"proposed_signature"`
	Address           string                            `json:"address"`
	PointID           string                            `json:"point_id"`
	TagID             *string                           `json:"tag_id,omitempty"`
	TagKey            string                            `json:"tag_key,omitempty"`
	DisplayName       string                            `json:"display_name,omitempty"`
	DataType          DataType                          `json:"data_type"`
	Register          *uint16                           `json:"register,omitempty"`
	RegisterCount     int                               `json:"register_count"`
	Status            SourceRuleLocalModbusOutputStatus `json:"status"`
	BlockingReason    string                            `json:"blocking_reason,omitempty"`
	UpdatedAt         *time.Time                        `json:"updated_at,omitempty"`
}

type SourceRuleTagReviewDecisionAction string

const (
	SourceRuleTagReviewDecisionActionRename   SourceRuleTagReviewDecisionAction = "rename"
	SourceRuleTagReviewDecisionActionSkip     SourceRuleTagReviewDecisionAction = "skip"
	SourceRuleTagReviewDecisionActionOverride SourceRuleTagReviewDecisionAction = "override"
)

// SourceRuleTagReviewDecision persists one rule-scoped review choice for a canonical tag candidate identity.
type SourceRuleTagReviewDecision struct {
	SourceRuleID    string                            `json:"source_rule_id" db:"source_rule_id"`
	CandidateID     string                            `json:"candidate_id" db:"candidate_id"`
	Action          SourceRuleTagReviewDecisionAction `json:"action" db:"decision_type"`
	TagKey          string                            `json:"tag_key,omitempty" db:"tag_key"`
	OverrideTagID   *string                           `json:"override_tag_id,omitempty" db:"override_tag_id"`
	Stale           bool                              `json:"stale" db:"stale"`
	StaleRevisionID string                            `json:"stale_revision_id,omitempty" db:"stale_revision_id"`
	StaleAt         *time.Time                        `json:"stale_at,omitempty" db:"stale_at"`
	CreatedAt       time.Time                         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time                         `json:"updated_at" db:"updated_at"`
}

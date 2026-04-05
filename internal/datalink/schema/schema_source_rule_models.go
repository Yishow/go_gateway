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
	SourceRuleCandidateTypeDatabaseOutputs    SourceRuleCandidateType = "database_outputs"
	SourceRuleCandidateTypeLocalModbusOutputs SourceRuleCandidateType = "local_modbus_outputs"
)

type SourceRuleCandidateStatus string

const (
	SourceRuleCandidateStatusReady    SourceRuleCandidateStatus = "ready"
	SourceRuleCandidateStatusBlocked  SourceRuleCandidateStatus = "blocked"
	SourceRuleCandidateStatusDeferred SourceRuleCandidateStatus = "deferred"
)

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
	Address     string   `json:"address"`
	PointID     string   `json:"point_id"`
	TagID       *string  `json:"tag_id,omitempty"`
	MappingID   *string  `json:"mapping_id,omitempty"`
	TagKey      string   `json:"tag_key,omitempty"`
	DisplayName string   `json:"display_name,omitempty"`
	DataType    DataType `json:"data_type"`
}

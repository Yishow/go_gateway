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

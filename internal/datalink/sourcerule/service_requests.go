package sourcerule

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

type CreateRuleRequest struct {
	ID               string          `json:"id,omitempty"`
	DeviceID         string          `json:"device_id"`
	StartAddress     string          `json:"start_address"`
	Count            int             `json:"count"`
	DataType         schema.DataType `json:"data_type"`
	NamingPrefix     string          `json:"naming_prefix"`
	Enabled          bool            `json:"enabled"`
	Locked           bool            `json:"locked"`
	Origin           string          `json:"origin,omitempty"`
	TemplateName     string          `json:"template_name,omitempty"`
	SkippedAddresses []string        `json:"skipped_addresses,omitempty"`
	// TargetDataType 目標資料型態（可空）。預設與 DataType 相同。
	TargetDataType *schema.DataType `json:"target_data_type,omitempty"`
	// ScaleMultiplier 縮放倍率（可空）。
	ScaleMultiplier *float64 `json:"scale_multiplier,omitempty"`
	// ScaleOffset 偏移量（可空）。
	ScaleOffset *float64 `json:"scale_offset,omitempty"`
	// DataFormat 字節序格式（可空）。空字串表示使用設備連線預設。
	DataFormat string `json:"data_format,omitempty"`
	// ShareEnabled controls whether this rule is exposed through the local Modbus Share output.
	ShareEnabled bool `json:"share_enabled"`
	// ShareStartRegister is the optional local Modbus holding-register start address.
	ShareStartRegister *int `json:"share_start_register"`
	// ShareStride is the optional number of local Modbus registers allocated per point.
	ShareStride *int `json:"share_stride"`
	// LayoutMode is homogeneous (default) or mixed. A mixed rule derives
	// start_address, count and data_type from its first item.
	LayoutMode string `json:"layout_mode,omitempty" enums:"homogeneous,mixed"`
	// MixedItems lists one point per item; each item has its own address and data_type.
	MixedItems []measurement.MixedItem `json:"mixed_items,omitempty"`
}

type UpdateRuleRequest struct {
	StartAddress     *string          `json:"start_address,omitempty"`
	Count            *int             `json:"count,omitempty"`
	DataType         *schema.DataType `json:"data_type,omitempty"`
	NamingPrefix     *string          `json:"naming_prefix,omitempty"`
	Enabled          *bool            `json:"enabled,omitempty"`
	Locked           *bool            `json:"locked,omitempty"`
	TemplateName     *string          `json:"template_name,omitempty"`
	SkippedAddresses *[]string        `json:"skipped_addresses,omitempty"`
	// TargetDataType 目標資料型態（可空）。
	TargetDataType    *schema.DataType `json:"target_data_type,omitempty"`
	TargetDataTypeSet bool             `json:"-"`
	// ScaleMultiplier 縮放倍率（可空）。
	ScaleMultiplier    *float64 `json:"scale_multiplier,omitempty"`
	ScaleMultiplierSet bool     `json:"-"`
	// ScaleOffset 偏移量（可空）。
	ScaleOffset    *float64 `json:"scale_offset,omitempty"`
	ScaleOffsetSet bool     `json:"-"`
	// DataFormat 字節序格式（可空）；傳 null 並搭配 DataFormatSet 可清空為連線預設。
	DataFormat    *string `json:"data_format,omitempty"`
	DataFormatSet bool    `json:"-"`
	// ShareEnabled controls whether this rule is exposed through the local Modbus Share output.
	ShareEnabled *bool `json:"share_enabled,omitempty"`
	// ShareStartRegister is the optional local Modbus holding-register start address.
	ShareStartRegister    *int `json:"share_start_register,omitempty"`
	ShareStartRegisterSet bool `json:"-"`
	// ShareStride is the optional number of local Modbus registers allocated per point.
	ShareStride    *int `json:"share_stride,omitempty"`
	ShareStrideSet bool `json:"-"`
	// LayoutMode switches the rule layout; omitted keeps the stored layout.
	LayoutMode *string `json:"layout_mode,omitempty" enums:"homogeneous,mixed"`
	// MixedItems replaces the complete item list of a mixed rule; omitted keeps the stored items.
	MixedItems *[]measurement.MixedItem `json:"mixed_items,omitempty"`
}

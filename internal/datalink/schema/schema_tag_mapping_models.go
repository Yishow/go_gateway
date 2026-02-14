package schema

import "time"

// =============================================================================
// 標籤模型 (Tag)
// =============================================================================

// TagStatus 標籤生命週期狀態
type TagStatus string

const (
	// TagStatusDraft 草稿狀態
	TagStatusDraft TagStatus = "draft"
	// TagStatusActive 啟用狀態
	TagStatusActive TagStatus = "active"
	// TagStatusRetired 退役狀態 (不允許新映射)
	TagStatusRetired TagStatus = "retired"
)

// Tag 全域標籤實體
type Tag struct {
	// ID 主鍵，UUID 格式
	ID string `json:"id" db:"id"`

	// Key 標籤鍵 (唯一，大小寫不敏感)
	// 格式: ASCII 字母、數字、底線、連字號、點、斜線
	// 長度: 1-128 字元
	Key string `json:"key" db:"key"`

	// KeyLower 標籤鍵小寫版本 (用於唯一性檢查)
	KeyLower string `json:"-" db:"key_lower"`

	// DisplayName 顯示名稱
	DisplayName string `json:"display_name" db:"display_name"`

	// Description 標籤描述
	Description string `json:"description,omitempty" db:"description"`

	// Unit 單位 (如 ℃, kWh, rpm)
	Unit string `json:"unit,omitempty" db:"unit"`

	// DataType 資料型別
	DataType DataType `json:"data_type" db:"data_type"`

	// Status 標籤狀態
	Status TagStatus `json:"status" db:"status"`

	// Labels 標籤屬性 (JSON 物件，用於分類和篩選)
	Labels string `json:"labels,omitempty" db:"labels"`

	// CreatedAt 建立時間
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// UpdatedAt 更新時間
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// =============================================================================
// 映射模型 (Mapping)
// =============================================================================

// Mapping 點位到標籤的映射關係
type Mapping struct {
	// ID 主鍵，UUID 格式
	ID string `json:"id" db:"id"`

	// PointID 來源點位 ID (外鍵)
	PointID string `json:"point_id" db:"point_id"`

	// TagID 目標標籤 ID (外鍵)
	TagID string `json:"tag_id" db:"tag_id"`

	// TransformPipeline 轉換管線步驟 (JSON 陣列)
	// 詳見 TransformStep 類型定義
	TransformPipeline string `json:"transform_pipeline" db:"transform_pipeline"`

	// Enabled 是否啟用
	Enabled bool `json:"enabled" db:"enabled"`

	// CreatedAt 建立時間
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// UpdatedAt 更新時間
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// =============================================================================
// 轉換步驟模型 (TransformStep)
// =============================================================================

// TransformType 轉換步驟類型
type TransformType string

const (
	// TransformDecode 解碼 (如 BCD, Gray code)
	TransformDecode TransformType = "decode"
	// TransformCast 型別轉換
	TransformCast TransformType = "cast"
	// TransformScale 線性縮放 (value * scale + offset)
	TransformScale TransformType = "scale"
	// TransformLookup 查表替換
	TransformLookup TransformType = "lookup"
	// TransformConditional 條件判斷
	TransformConditional TransformType = "conditional"
	// TransformFormula 公式表達式
	TransformFormula TransformType = "formula"
)

// TransformStep 單一轉換步驟的配置
type TransformStep struct {
	// Type 轉換類型
	Type TransformType `json:"type"`

	// Order 執行順序 (從小到大)
	Order int `json:"order"`

	// Params 轉換參數 (依類型不同)
	Params map[string]interface{} `json:"params"`
}

// TransformParamsScale 線性縮放參數
type TransformParamsScale struct {
	// Scale 縮放係數
	Scale float64 `json:"scale"`
	// Offset 偏移量
	Offset float64 `json:"offset"`
}

// TransformParamsCast 型別轉換參數
type TransformParamsCast struct {
	// TargetType 目標型別
	TargetType DataType `json:"target_type"`
}

// TransformParamsLookup 查表參數
type TransformParamsLookup struct {
	// Table 對照表 (key-value 映射)
	Table map[string]interface{} `json:"table"`
	// Default 無匹配時的預設值
	Default interface{} `json:"default,omitempty"`
}

// TransformParamsConditional 條件判斷參數
type TransformParamsConditional struct {
	// Condition 條件表達式
	Condition string `json:"condition"`
	// TrueValue 條件為真時的值或表達式
	TrueValue interface{} `json:"true_value"`
	// FalseValue 條件為假時的值或表達式
	FalseValue interface{} `json:"false_value"`
}

// TransformParamsFormula 公式表達式參數
type TransformParamsFormula struct {
	// Expression 表達式字串
	Expression string `json:"expression"`
}

// TransformParamsDecode 解碼參數
type TransformParamsDecode struct {
	// Format 解碼格式 (bcd, gray, hex_string, etc.)
	Format string `json:"format"`
}

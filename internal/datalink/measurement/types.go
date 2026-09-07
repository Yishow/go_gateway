package measurement

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

// SemanticKind 定義量測項目的物理與記錄語意類型。
type SemanticKind string

const (
	// SemanticKindGauge 瞬時狀態值（如溫度、壓力、電壓、電流、功率、轉速）。
	SemanticKindGauge SemanticKind = "gauge"
	// SemanticKindRate 瞬時速率（如每小時流量 m3/h、每秒件數 pcs/s）。
	SemanticKindRate SemanticKind = "rate"
	// SemanticKindCounter 單調遞增累積量（如累積電量 kWh、累積用水 m3、累計生產件數）。
	SemanticKindCounter SemanticKind = "counter"
	// SemanticKindSignedCounter 可正可負的淨累積量（如淨注入電量、進出流量淨值）。
	SemanticKindSignedCounter SemanticKind = "signed_counter"
	// SemanticKindDelta 區間已計算差值（必須有明確區間或事件身份）。
	SemanticKindDelta SemanticKind = "delta"
	// SemanticKindState 設備運行狀態（如 0:停機, 1:運轉, 2:待機, 3:故障）。
	SemanticKindState SemanticKind = "state"
	// SemanticKindEvent 警報或離散事件（如 bitmask 警報旗標、異常觸發）。
	SemanticKindEvent SemanticKind = "event"
	// SemanticKindText 文字型資訊（如工單號碼、批次識別碼）。
	SemanticKindText SemanticKind = "text"
)

// IsValid 檢查 SemanticKind 是否合法。
func (k SemanticKind) IsValid() bool {
	switch k {
	case SemanticKindGauge, SemanticKindRate, SemanticKindCounter,
		SemanticKindSignedCounter, SemanticKindDelta, SemanticKindState,
		SemanticKindEvent, SemanticKindText:
		return true
	default:
		return false
	}
}

// CounterPolicy 定義計量累積值的邊界政策。
type CounterPolicy struct {
	ResetThreshold *float64 `json:"reset_threshold,omitempty"`
	MaxRollOver    *float64 `json:"max_roll_over,omitempty"`
	AllowNegative  bool     `json:"allow_negative"`
}

// MeasurementDefinition 記錄逐項量測的完整物理意義、語意、單位與身份。
type MeasurementDefinition struct {
	ID                    string            `json:"id" db:"id"`
	WorkspaceID           string            `json:"workspace_id" db:"workspace_id"`
	DeviceID              string            `json:"device_id" db:"device_id"`
	PointID               string            `json:"point_id" db:"point_id"`
	TagID                 *string           `json:"tag_id,omitempty" db:"tag_id"`
	EquipmentID           string            `json:"equipment_id" db:"equipment_id"`
	DefinitionRevision    string            `json:"definition_revision" db:"definition_revision"`
	SourceBindingRevision string            `json:"source_binding_revision" db:"source_binding_revision"`
	SeriesEpoch           string            `json:"series_epoch" db:"series_epoch"`
	Name                  string            `json:"name" db:"name"`
	Quantity              string            `json:"quantity" db:"quantity"`
	Unit                  string            `json:"unit,omitempty" db:"unit"`
	SemanticKind          SemanticKind      `json:"semantic_kind" db:"semantic_kind"`
	NumericEncoding       string            `json:"numeric_encoding,omitempty" db:"numeric_encoding"`
	CounterPolicy         *CounterPolicy    `json:"counter_policy,omitempty" db:"counter_policy"`
	StateMap              map[string]string `json:"state_map,omitempty" db:"state_map"`
	BitmaskLabels         map[int]string    `json:"bitmask_labels,omitempty" db:"bitmask_labels"`
	CreatedAt             time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at" db:"updated_at"`
}

// Validate 檢查 MeasurementDefinition 欄位合法性。
func (m *MeasurementDefinition) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return fmt.Errorf("measurement_id cannot be empty")
	}
	if strings.TrimSpace(m.DeviceID) == "" {
		return fmt.Errorf("device_id cannot be empty")
	}
	if strings.TrimSpace(m.PointID) == "" {
		return fmt.Errorf("point_id cannot be empty")
	}
	if !m.SemanticKind.IsValid() {
		return fmt.Errorf("invalid semantic_kind: %s", m.SemanticKind)
	}
	if strings.TrimSpace(m.SeriesEpoch) == "" {
		m.SeriesEpoch = "epoch-1"
	}
	if strings.TrimSpace(m.EquipmentID) == "" {
		m.EquipmentID = m.DeviceID
	}
	return nil
}

// SampleEnvelope 採集樣本的統一信封契約。
type SampleEnvelope struct {
	SampleID              string             `json:"sample_id"`
	WorkspaceID           string             `json:"workspace_id,omitempty"`
	MeasurementID         string             `json:"measurement_id"`
	SeriesEpoch           string             `json:"series_epoch"`
	DefinitionRevision    string             `json:"definition_revision"`
	SourceBindingRevision string             `json:"source_binding_revision"`
	AcquisitionID         string             `json:"acquisition_id"`
	SourceSequence        *int64             `json:"source_sequence,omitempty"`
	ObservedAt            time.Time          `json:"observed_at"`
	ReceivedAt            time.Time          `json:"received_at"`
	TimeOrigin            string             `json:"time_origin"`
	ValueType             string             `json:"value_type"`
	Value                 interface{}        `json:"value"`
	RawValue              interface{}        `json:"raw_value,omitempty"`
	Quality               schema.QualityFlag `json:"quality"`
	QualityReason         string             `json:"quality_reason,omitempty"`
}

// MarshalJSON 實作安全的 SampleEnvelope JSON 序列化，保護大整數精度。
func (s SampleEnvelope) MarshalJSON() ([]byte, error) {
	type Alias SampleEnvelope
	aux := struct {
		Alias
		ObservedAt string  `json:"observed_at"`
		ReceivedAt string  `json:"received_at"`
		ValueStr   *string `json:"value_str,omitempty"`
	}{
		Alias:      Alias(s),
		ObservedAt: s.ObservedAt.UTC().Format(time.RFC3339Nano),
		ReceivedAt: s.ReceivedAt.UTC().Format(time.RFC3339Nano),
	}

	// 大於 JS Number.MAX_SAFE_INTEGER (9007199254740991) 的整數提供 string 格式保護
	switch v := s.Value.(type) {
	case uint64:
		str := strconv.FormatUint(v, 10)
		aux.ValueStr = &str
		if v > 9007199254740991 {
			aux.Value = str
		}
	case int64:
		str := strconv.FormatInt(v, 10)
		aux.ValueStr = &str
		if v > 9007199254740991 || v < -9007199254740991 {
			aux.Value = str
		}
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			aux.Quality = schema.QualityBad
			aux.QualityReason = "invalid float value"
			aux.Value = nil
		}
	}

	return json.Marshal(aux)
}

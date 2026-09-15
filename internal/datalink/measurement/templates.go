package measurement

import (
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// MeasurementTemplate 描述版本化的量測佈局與語意範本。
type MeasurementTemplate struct { //nolint:revive // Preserve the exported Go name and its existing callers during lint maintenance.
	ID          string      `json:"id"`
	Version     string      `json:"version"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	LayoutMode  LayoutMode  `json:"layout_mode"`
	Items       []MixedItem `json:"items"`
}

const (
	quantityVoltage = "voltage"
	quantityCurrent = "current"
)

// BuiltinTemplates 內建標準範本庫。
var BuiltinTemplates = map[string]MeasurementTemplate{
	"three_phase_power_meter_v1": {
		ID:          "three_phase_power_meter_v1",
		Version:     "1.0.0",
		Name:        "三相多功能電表範本",
		Description: "標準三相電表：電壓(Va, Vb, Vc)、電流(Ia, Ib, Ic)、有效功率(kW)與累積電量(kWh)",
		LayoutMode:  LayoutModeMixed,
		Items: []MixedItem{
			{ItemID: "va", Name: "Phase A Voltage", RegisterOffset: 0, DataType: schema.DataTypeFloat32, Quantity: quantityVoltage, Unit: "V", SemanticKind: SemanticKindGauge},
			{ItemID: "vb", Name: "Phase B Voltage", RegisterOffset: 2, DataType: schema.DataTypeFloat32, Quantity: quantityVoltage, Unit: "V", SemanticKind: SemanticKindGauge},
			{ItemID: "vc", Name: "Phase C Voltage", RegisterOffset: 4, DataType: schema.DataTypeFloat32, Quantity: quantityVoltage, Unit: "V", SemanticKind: SemanticKindGauge},
			{ItemID: "ia", Name: "Phase A Current", RegisterOffset: 6, DataType: schema.DataTypeFloat32, Quantity: quantityCurrent, Unit: "A", SemanticKind: SemanticKindGauge},
			{ItemID: "ib", Name: "Phase B Current", RegisterOffset: 8, DataType: schema.DataTypeFloat32, Quantity: quantityCurrent, Unit: "A", SemanticKind: SemanticKindGauge},
			{ItemID: "ic", Name: "Phase C Current", RegisterOffset: 10, DataType: schema.DataTypeFloat32, Quantity: quantityCurrent, Unit: "A", SemanticKind: SemanticKindGauge},
			{ItemID: "kw", Name: "Active Power", RegisterOffset: 12, DataType: schema.DataTypeFloat32, Quantity: "power", Unit: "kW", SemanticKind: SemanticKindGauge},
			{ItemID: "kwh", Name: "Active Energy", RegisterOffset: 14, DataType: schema.DataTypeUint32, Quantity: "energy", Unit: "kWh", SemanticKind: SemanticKindCounter, DataFormat: "ABCD"},
		},
	},
	"mixed_equipment_sensors_v1": {
		ID:          "mixed_equipment_sensors_v1",
		Version:     "1.0.0",
		Name:        "三菱 PLC 混合感測器範本 (D0~D7)",
		Description: "連續暫存器包含溫度、壓力、瞬時流量、跨 register 累積流量、狀態、警報旗標與累計產量",
		LayoutMode:  LayoutModeMixed,
		Items: []MixedItem{
			{ItemID: "temp", Name: "溫度", Address: "D0", RegisterOffset: 0, DataType: schema.DataTypeInt16, Quantity: "temperature", Unit: "°C", SemanticKind: SemanticKindGauge, ScaleMultiplier: floatPtr(0.1)},
			{ItemID: "press", Name: "壓力", Address: "D1", RegisterOffset: 1, DataType: schema.DataTypeUint16, Quantity: "pressure", Unit: "bar", SemanticKind: SemanticKindGauge, ScaleMultiplier: floatPtr(0.01)},
			{ItemID: "flow_rate", Name: "瞬時流量", Address: "D2", RegisterOffset: 2, DataType: schema.DataTypeUint16, Quantity: "flow_rate", Unit: "m3/h", SemanticKind: SemanticKindRate},
			{ItemID: "total_flow", Name: "累積流量", Address: "D3", RegisterOffset: 3, DataType: schema.DataTypeUint32, Quantity: "volume", Unit: "m3", SemanticKind: SemanticKindCounter, DataFormat: "CDAB"},
			{ItemID: "status", Name: "運轉狀態", Address: "D5", RegisterOffset: 5, DataType: schema.DataTypeUint16, Quantity: "state", SemanticKind: SemanticKindState, StateMap: map[string]string{"0": "停止", "1": "運轉", "2": "待機", "3": "異常"}},
			{ItemID: "alarms", Name: "警報旗標", Address: "D6", RegisterOffset: 6, DataType: schema.DataTypeUint16, Quantity: "alarm", SemanticKind: SemanticKindEvent, BitmaskLabels: map[int]string{0: "過熱", 1: "過載", 2: "急停", 3: "感測器斷線"}},
			{ItemID: "prod_count", Name: "累計產量", Address: "D7", RegisterOffset: 7, DataType: schema.DataTypeUint16, Quantity: "production_count", Unit: "pcs", SemanticKind: SemanticKindCounter},
		},
	},
}

func floatPtr(v float64) *float64 {
	return &v
}

// TemplateApplyPreview 描述批次套用範本時各設備的預覽結果。
type TemplateApplyPreview struct {
	DeviceID      string                  `json:"device_id"`
	DeviceName    string                  `json:"device_name"`
	TemplateID    string                  `json:"template_id"`
	ProposedItems []MixedItem             `json:"proposed_items"`
	Definitions   []MeasurementDefinition `json:"definitions"`
	Confirmed     bool                    `json:"confirmed"`
	NeedsReview   bool                    `json:"needs_review"`
	ReviewNotes   []string                `json:"review_notes,omitempty"`
}

// GenerateTemplatePreview 為指定設備產生範本套用預覽，提示待確認項目。
func GenerateTemplatePreview(templateID, deviceID, deviceName, baseAddress string) (*TemplateApplyPreview, error) {
	tmpl, exists := BuiltinTemplates[templateID]
	if !exists {
		return nil, fmt.Errorf("template %s not found", templateID)
	}

	preview := &TemplateApplyPreview{
		DeviceID:      deviceID,
		DeviceName:    deviceName,
		TemplateID:    templateID,
		ProposedItems: make([]MixedItem, len(tmpl.Items)),
		Definitions:   make([]MeasurementDefinition, len(tmpl.Items)),
		Confirmed:     false,
		NeedsReview:   true,
		ReviewNotes:   make([]string, 0),
	}

	for i, itm := range tmpl.Items {
		itemCopy := itm
		if strings.TrimSpace(baseAddress) != "" && itemCopy.Address == "" {
			itemCopy.Address = fmt.Sprintf("%s+%d", baseAddress, itemCopy.RegisterOffset)
		}
		preview.ProposedItems[i] = itemCopy

		preview.Definitions[i] = MeasurementDefinition{
			ID:                    fmt.Sprintf("meas-%s-%s", deviceID, itemCopy.ItemID),
			WorkspaceID:           "default",
			DeviceID:              deviceID,
			PointID:               fmt.Sprintf("pt-%s-%s", deviceID, itemCopy.ItemID),
			EquipmentID:           deviceID,
			DefinitionRevision:    initialDefinitionRevision,
			SourceBindingRevision: initialDefinitionRevision,
			SeriesEpoch:           initialSeriesEpoch,
			Name:                  fmt.Sprintf("%s - %s", deviceName, itemCopy.Name),
			Quantity:              itemCopy.Quantity,
			Unit:                  itemCopy.Unit,
			SemanticKind:          itemCopy.SemanticKind,
			CounterPolicy:         itemCopy.CounterPolicy,
			StateMap:              itemCopy.StateMap,
			BitmaskLabels:         itemCopy.BitmaskLabels,
		}
	}

	preview.ReviewNotes = append(preview.ReviewNotes, "請確認各回路接線、暫存器起始位址與 CT/PT 換算倍率無誤後再啟用")
	return preview, nil
}

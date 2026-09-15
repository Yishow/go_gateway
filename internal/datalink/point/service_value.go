package point

import (
	"context"
	"encoding/json"

	"go-gateway/internal/datalink/schema"
)

// PointDetail 點位詳細資訊 (包含設備和群組資訊)
type PointDetail struct { //nolint:revive // Preserve the exported Go name and its existing callers during lint maintenance.
	*schema.Point
	DeviceName       string `json:"device_name,omitempty"`
	DeviceProtocol   string `json:"device_protocol,omitempty"`
	PollingGroupName string `json:"polling_group_name,omitempty"`
}

// ToJSON 將最後讀取值轉換為 JSON
func (s *Service) GetLastValue(ctx context.Context, pointID string) (interface{}, error) {
	point, err := s.repo.GetByID(ctx, pointID)
	if err != nil {
		return nil, err
	}

	if point.LastValue == nil {
		return nil, nil
	}

	var value interface{}
	if err := json.Unmarshal([]byte(*point.LastValue), &value); err != nil {
		// 非 JSON 格式，直接返回字串
		return *point.LastValue, nil //nolint:nilerr // Stored values also support plain text; parsing failure selects that representation.
	}

	return value, nil
}

package sourcerule

import (
	"context"
	"fmt"
	"go-gateway/internal/datalink/point"
	"strings"
)

func (s *Service) validatePointAddresses(ctx context.Context, deviceID string, keepPointIDs map[string]struct{}, addresses []string) error {
	filter := point.ListFilter{DeviceID: &deviceID}
	points, err := s.pointSvc.List(ctx, filter)
	if err != nil {
		return fmt.Errorf("取得既有點位失敗: %w", err)
	}

	existing := make(map[string]string, len(points))
	for _, item := range points {
		if _, ok := keepPointIDs[item.ID]; ok {
			continue
		}
		existing[strings.ToUpper(strings.TrimSpace(item.Address))] = item.ID
	}

	for _, address := range addresses {
		if existing[strings.ToUpper(strings.TrimSpace(address))] != "" {
			return fmt.Errorf("位址 %s 已存在既有點位", address)
		}
	}
	return nil
}

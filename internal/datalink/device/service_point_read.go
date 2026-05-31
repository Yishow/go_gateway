package device

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) ReadPointNow(
	ctx context.Context,
	point *schema.Point,
) (connector.ReadResult, error) {
	if point == nil {
		return connector.ReadResult{}, fmt.Errorf("point is required")
	}

	deviceRecord, err := s.repo.GetByID(ctx, point.DeviceID)
	if err != nil {
		return connector.ReadResult{}, fmt.Errorf("取得設備失敗: %w", err)
	}

	conn, err := s.connMgr.GetOrCreate(
		ctx,
		deviceRecord.ID,
		deviceRecord.Protocol,
		deviceRecord.ConnectionConfig,
	)
	if err != nil {
		return connector.ReadResult{}, fmt.Errorf("建立設備連線失敗: %w", err)
	}

	req := connector.ReadRequest{
		Address:    point.Address,
		Function:   point.Function,
		DataType:   point.DataType,
		DataFormat: point.DataFormat,
	}
	req.Count = schema.RegisterCountForDataType(point.DataType)

	result, err := conn.Read(ctx, req)
	if err != nil {
		return result, err
	}
	return result, nil
}

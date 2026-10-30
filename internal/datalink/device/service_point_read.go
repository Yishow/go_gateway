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
	profiles ...connector.ReadProfile,
) (connector.ReadResult, error) {
	if point == nil {
		return connector.ReadResult{}, fmt.Errorf("point is required")
	}

	var profile connector.ReadProfile
	if len(profiles) > 0 {
		profile = profiles[0]
	}
	if err := connector.ValidatePointReadProfile(point, profile); err != nil {
		return connector.ReadResult{Quality: schema.QualityBad, Error: err.Error()}, err
	}
	deviceRecord, err := s.repo.GetByID(ctx, point.DeviceID)
	if err != nil {
		return connector.ReadResult{}, fmt.Errorf("取得設備失敗: %w", err)
	}

	if !profile.MatchesConnection(deviceRecord.Protocol, deviceRecord.ConnectionConfig) {
		return connector.ReadResult{Quality: schema.QualityBad, Error: connector.ErrReadProfile.Error()}, connector.ErrReadProfile
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
		Profile:    profile,
	}
	count, countErr := connector.PointReadCountForAddress(conn.ProtocolType, point.Address, point.DataType)
	if countErr != nil {
		return connector.ReadResult{Quality: schema.QualityBad, Error: countErr.Error()}, countErr
	}
	req.Count = count

	result, err := conn.Read(ctx, req)
	if err != nil {
		return result, err
	}
	return result, nil
}

package point

import "context"

// UpdateReadResult 更新點位讀取結果
func (s *Service) UpdateReadResult(ctx context.Context, pointID string, value interface{}, errMsg string) error {
	return s.repo.UpdateReadResult(ctx, pointID, value, errMsg)
}

// BatchUpdateReadResult 批次更新讀取結果
func (s *Service) BatchUpdateReadResult(ctx context.Context, results []ReadResultUpdate) error {
	return s.repo.BatchUpdateReadResult(ctx, results)
}

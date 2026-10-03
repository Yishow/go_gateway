//go:build !f_write_group_fixture

package runtime

import (
	"context"
	"time"

	"go-gateway/internal/datalink/collector"
)

func (s *Service) consumeCollectedValue(runtimeCtx context.Context, cv collector.CollectedValue) {
	s.collectedTotal.Add(1)
	ctx, cancel := context.WithTimeout(runtimeCtx, 5*time.Second)
	s.handleCollectedValue(ctx, cv)
	cancel()
}

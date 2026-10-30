package collector

import (
	"context"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
	"time"
)

// acquisitionReadContext is shared by scheduled and manual reads. It only
// enables protocols whose bounded-transport contract has been installed.
func (s *Scheduler) acquisitionReadContext(ctx context.Context, conn *connector.ManagedConnection) (context.Context, context.CancelFunc) {
	if bounded, ok := conn.Protocol.(interface{ AcquisitionReadTimeout() time.Duration }); ok && conn.ProtocolType == schema.ProtocolMC3E {
		ctx, cancel := context.WithTimeout(ctx, bounded.AcquisitionReadTimeout())
		return wirebudget.With(ctx, s.config.DefaultRetryCount+1, s.config.DefaultRetryDelay), cancel
	}
	return ctx, func() {}
}

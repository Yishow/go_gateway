package main

import (
	"context"
	"time"

	"go-gateway/internal/datalink/modbusshare"
	datalinkruntime "go-gateway/internal/datalink/runtime"
)

// modbusShareTargetWriter keeps the production collector path connected to
// the authoritative Share projection. Diagnostic sync endpoints are not part
// of the runtime value-delivery path.
type modbusShareTargetWriter struct {
	share *modbusshare.Service
}

func (w modbusShareTargetWriter) targetActive(tagID string) bool {
	if w.share == nil || !w.share.Settings().Enabled {
		return false
	}
	status := w.share.Status()
	return status.ConfiguredEnabled && status.Enabled && w.share.HasMapping(tagID)
}

func (w modbusShareTargetWriter) WriteTagValue(ctx context.Context, tagID string, value any, _ time.Time) error {
	if w.share == nil || !w.share.HasMapping(tagID) {
		return nil
	}
	hydration, err := w.share.CheckHydration(ctx)
	if err != nil {
		return datalinkruntime.NewDeliveryError("modbus_share", "hydration", err)
	}
	if hydration.State != modbusshare.HydrationStateReady || !hydration.Readiness {
		return datalinkruntime.NewDeliveryError("modbus_share", "hydration", modbusshare.NewError(
			modbusshare.ErrCodeHydrationRequired,
			"workspace hydration is not ready",
			true,
		))
	}
	if !w.share.Settings().Enabled {
		return datalinkruntime.NewDeliveryError("modbus_share", "listener_disabled", modbusshare.NewError(modbusshare.ErrCodeDisabled, "Modbus Share is disabled", false))
	}
	if err := w.share.WriteProjectedTagValue(ctx, tagID, value); err != nil {
		return datalinkruntime.NewDeliveryError("modbus_share", "modbus_write", err)
	}
	return nil
}

type fanoutTargetWriter struct {
	targets []datalinkruntime.TargetWriter
}

func newProductionTargetWriter(database datalinkruntime.TargetWriter, share *modbusshare.Service) datalinkruntime.TargetWriter {
	return fanoutTargetWriter{targets: []datalinkruntime.TargetWriter{database, modbusShareTargetWriter{share: share}}}
}

func (w fanoutTargetWriter) WriteTagValue(ctx context.Context, tagID string, value any, observedAt time.Time) error {
	outcomes := w.WriteTagValueOutcomes(ctx, tagID, value, observedAt)
	var firstErr error
	for _, outcome := range outcomes {
		if outcome.Err != nil && firstErr == nil {
			firstErr = outcome.Err
		}
	}
	return firstErr
}

// WriteTagValueOutcomes attempts every production target and preserves each
// target result for runtime diagnostics. A failure from one target never
// prevents the other target from receiving the value.
func (w fanoutTargetWriter) WriteTagValueOutcomes(ctx context.Context, tagID string, value any, observedAt time.Time) []datalinkruntime.TargetDeliveryOutcome {
	results := make([]datalinkruntime.TargetDeliveryOutcome, 0, len(w.targets))
	for index, target := range w.targets {
		if target == nil {
			continue
		}
		if active, ok := target.(interface{ targetActive(string) bool }); ok && !active.targetActive(tagID) {
			continue
		}
		targetName := datalinkruntime.TargetDatabase
		if index == 1 {
			targetName = datalinkruntime.TargetModbusShare
		}
		results = append(results, datalinkruntime.TargetDeliveryOutcome{Target: targetName, Err: target.WriteTagValue(ctx, tagID, value, observedAt)})
	}
	return results
}

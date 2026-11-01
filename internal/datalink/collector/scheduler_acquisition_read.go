package collector

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/internal/protocol/modbus"
	"go-gateway/internal/protocol/wirebudget"
)

type acquisitionTrace struct{ requestID, planFingerprint, planReason string }

// readScheduledRequest owns exactly one original timeout and actual-send budget.
// Adapter reconnects and scheduler retries share both; neither can mint attempts.
func (s *Scheduler) readScheduledRequest(ctx context.Context, conn *connector.ManagedConnection, scan scheduledPlan, plan optimizer.AcquisitionPlan, index int, acquisitionID, phase string) ([]connector.PlannedPointResult, error) {
	request, ok := plan.Request(index)
	if !ok {
		return nil, connector.ErrPlannedRead
	}
	read := connector.PlannedReadRequest{Plan: plan, RequestIndex: index, AcquisitionID: acquisitionID, RequestID: fmt.Sprintf("%s/%s/%d", acquisitionID, phase, index)}
	points := make([]pointInfo, 0, len(request.Points))
	for _, point := range request.Points {
		read.Bindings = append(read.Bindings, scan.bindings[point.Point.ID])
		points = append(points, scan.points[point.Point.ID])
	}
	readCtx, cancel := s.acquisitionReadContext(ctx, conn)
	defer cancel()
	// This short captured-install check is the actual-send admission point. It
	// takes only scheduler.mu, never SQL, mapping locks or observer callbacks.
	gap := acquisitionRequestHasGap(request)
	readCtx = wirebudget.WithAdmission(readCtx, func() error { return s.admitScheduledWire(scan.snapshot, points, gap) })
	var results []connector.PlannedPointResult
	var err error
	for attempt := 0; attempt <= s.config.DefaultRetryCount; attempt++ {
		if attempt > 0 {
			if err = waitScheduledRetry(readCtx, s.config.DefaultRetryDelay, attempt); err != nil {
				break
			}
		}
		if err = readCtx.Err(); err != nil {
			break
		}
		results, err = conn.ReadPlanned(readCtx, read)
		if err == nil && plannedResultsGood(results) {
			break
		}
		if explicitCommandRejection(err) || errors.Is(err, errAcquisitionSuperseded) || errors.Is(err, wirebudget.ErrExhausted) {
			break
		}
		if budget := wirebudget.From(readCtx); budget != nil && budget.Remaining() == 0 {
			break
		}
	}
	if readCtx.Err() != nil {
		err = readCtx.Err()
	}
	return validateScheduledResults(read, results, err, s.now())
}

func waitScheduledRetry(ctx context.Context, delay time.Duration, attempt int) error {
	// Same delay*2^attempt contract as the scalar path, without shift overflow.
	if delay <= 0 {
		return ctx.Err()
	}
	for range attempt {
		if delay >= 15*time.Second {
			delay = 30 * time.Second
			break
		}
		delay *= 2
	}
	if delay < 0 {
		delay = 0
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func explicitCommandRejection(err error) bool {
	return mcprotocol.IsCommandRejected(err) || fatek.IsCommandRejected(err) || modbus.IsCommandRejected(err)
}

func plannedResultsGood(results []connector.PlannedPointResult) bool {
	if len(results) == 0 {
		return false
	}
	for _, result := range results {
		if result.Quality != schema.QualityGood || result.Error != "" {
			return false
		}
	}
	return true
}

// validateScheduledResults rejects incomplete or foreign scatter evidence as one
// failed request. It never borrows a different request's completion timestamp.
func validateScheduledResults(read connector.PlannedReadRequest, results []connector.PlannedPointResult, readErr error, completedAt time.Time) ([]connector.PlannedPointResult, error) {
	request, bindings, err := read.Snapshot()
	if err != nil {
		return nil, err
	}
	valid := len(results) == len(bindings)
	byID := make(map[string]connector.PlannedPointResult, len(results))
	var receivedAt time.Time
	var token connector.AcquisitionConnectionToken
	for _, result := range results {
		if _, exists := byID[result.Binding.Point.ID]; exists {
			valid = false
		}
		byID[result.Binding.Point.ID] = result
		if result.AcquisitionID != read.AcquisitionID || result.RequestID != read.RequestID || result.PlanFingerprint != read.Plan.Fingerprint() || result.Group != request.Group || result.ReceivedAt.IsZero() {
			valid = false
		}
		if receivedAt.IsZero() {
			receivedAt, token = result.ReceivedAt, result.ConnectionToken
		} else if !receivedAt.Equal(result.ReceivedAt) || token != result.ConnectionToken {
			valid = false
		}
	}
	for _, binding := range bindings {
		if result, ok := byID[binding.Point.ID]; !ok || result.Binding != binding {
			valid = false
		}
	}
	if !valid {
		// Keep a transport rejection only when the adapter supplied a complete
		// matching request result. Malformed evidence cannot authorize fallback.
		if readErr == nil || explicitCommandRejection(readErr) {
			readErr = connector.ErrPlannedRead
		}
		receivedAt = completedAt
		token = connector.AcquisitionConnectionToken{}
	}
	ordered := make([]connector.PlannedPointResult, 0, len(bindings))
	for _, binding := range bindings {
		result := byID[binding.Point.ID]
		if !valid {
			result = connector.PlannedPointResult{Binding: binding, Group: request.Group, PlanFingerprint: read.Plan.Fingerprint(), AcquisitionID: read.AcquisitionID, RequestID: read.RequestID, ReceivedAt: receivedAt, ConnectionToken: token}
		}
		if readErr != nil {
			result.ReadResult = connector.ReadResult{Quality: schema.QualityBad, Error: readErr.Error()}
		}
		ordered = append(ordered, result)
	}
	return ordered, readErr
}

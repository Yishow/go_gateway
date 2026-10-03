//go:build f_write_group_fixture

package runtime

import (
	"context"
	"errors"
	"sync/atomic"

	"go-gateway/internal/datalink/collector"
)

// AcceptFixtureCollectedValue routes a manually polled value through the same
// mapping and typed sink path used by the runtime consumer loop. It exists only
// in the acceptance binary; normal builds cannot call this control seam.
func (s *Service) AcceptFixtureCollectedValue(ctx context.Context, value collector.CollectedValue) error {
	if s == nil {
		return errors.New("fixture runtime unavailable")
	}
	reporter := &fixtureSampleErrorReporter{}
	ctx = withFixtureSampleErrorReporter(ctx, reporter)
	s.handleCollectedValue(ctx, value)
	if reporter.failed.Load() {
		return errors.New("fixture sample delivery failed")
	}
	return nil
}

type fixtureSampleErrorContextKey struct{}

type fixtureSampleErrorReporter struct {
	failed atomic.Bool
}

func withFixtureSampleErrorReporter(ctx context.Context, reporter *fixtureSampleErrorReporter) context.Context {
	return context.WithValue(ctx, fixtureSampleErrorContextKey{}, reporter)
}

// ReportFixtureSampleError lets the tagged fixture sink surface a bounded
// capture refusal without exposing its raw error through the runtime path.
func ReportFixtureSampleError(ctx context.Context) {
	if ctx == nil {
		return
	}
	reporter, ok := ctx.Value(fixtureSampleErrorContextKey{}).(*fixtureSampleErrorReporter)
	if ok && reporter != nil {
		reporter.failed.Store(true)
	}
}

// PauseFixture stops the fixture scheduler before manual polling. The tagged
// pause gate is added separately so this first boundary remains testable.
func (s *Service) PauseFixture(ctx context.Context) error {
	if s == nil || s.scheduler == nil {
		return errors.New("fixture runtime unavailable")
	}
	if err := s.scheduler.StopContext(ctx); err != nil {
		return err
	}
	if err := pauseFixtureConsumer(ctx, s); err != nil {
		return err
	}
	drainFixtureValues(s.scheduler.ValueChannel())
	return nil
}

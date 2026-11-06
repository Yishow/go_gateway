package device

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type acquisitionReadKey struct{}

type acquisitionReadGate struct {
	first    atomic.Bool
	waitAuto bool
}

type acquisitionReadBarrier struct {
	*SQLRepository
	reads    chan struct{}
	release  chan struct{}
	autoDone chan struct{}
}

func (r *acquisitionReadBarrier) GetByID(ctx context.Context, id string) (*schema.Device, error) {
	device, err := r.SQLRepository.GetByID(ctx, id)
	gate, _ := ctx.Value(acquisitionReadKey{}).(*acquisitionReadGate)
	if err != nil || gate == nil || !gate.first.CompareAndSwap(false, true) {
		return device, err
	}
	r.reads <- struct{}{}
	select {
	case <-r.release:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if gate.waitAuto {
		select {
		case <-r.autoDone:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return device, nil
}

func TestConcurrentPartialAcquisitionSettingsPreserveBothUpdates(t *testing.T) {
	db := setupTestDB(t)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	t.Cleanup(func() { _ = db.Close() })
	repo := NewSQLRepository(db)
	require.NoError(t, repo.Create(t.Context(), acquisitionTestDevice("dev")))
	barrier := &acquisitionReadBarrier{SQLRepository: repo, reads: make(chan struct{}, 2), release: make(chan struct{}), autoDone: make(chan struct{})}
	svc := NewService(barrier, nil)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	autoCtx := context.WithValue(ctx, acquisitionReadKey{}, &acquisitionReadGate{})
	frozenCtx := context.WithValue(ctx, acquisitionReadKey{}, &acquisitionReadGate{waitAuto: true})
	auto, frozen := true, 600000
	type result struct {
		frozen bool
		device *schema.Device
		err    error
	}
	results := make(chan result, 2)
	var group sync.WaitGroup
	group.Go(func() {
		dev, err := svc.Update(autoCtx, "dev", UpdateDeviceRequest{AcquisitionAutoInterval: &auto})
		close(barrier.autoDone)
		results <- result{device: dev, err: err}
	})
	group.Go(func() {
		dev, err := svc.Update(frozenCtx, "dev", UpdateDeviceRequest{FrozenThresholdMs: &frozen})
		results <- result{frozen: true, device: dev, err: err}
	})
	for range 2 {
		select {
		case <-barrier.reads:
		case <-ctx.Done():
			t.Fatal("concurrent requests did not reach their initial reads")
		}
	}
	close(barrier.release)
	group.Wait()
	for range 2 {
		result := <-results
		require.NoError(t, result.err)
		if result.frozen {
			require.True(t, result.device.AcquisitionAutoInterval, "the later successful response must include the previously saved auto setting")
		}
	}
	got, err := svc.GetByID(t.Context(), "dev")
	require.NoError(t, err)
	require.True(t, got.AcquisitionAutoInterval, "the successful auto update must survive the frozen update")
	require.Equal(t, 600000, got.FrozenThresholdMs)
}

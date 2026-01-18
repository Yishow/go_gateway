package collector

import (
	"testing"
	"time"
)

func TestEmitValueDoesNotBlock(t *testing.T) {
	config := DefaultSchedulerConfig()
	config.ValueBufferSize = 1
	scheduler := NewScheduler(config, nil)

	scheduler.valueChan <- CollectedValue{}

	done := make(chan struct{})
	go func() {
		scheduler.emitValue(CollectedValue{})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("emitValue blocked with full buffer")
	}
}

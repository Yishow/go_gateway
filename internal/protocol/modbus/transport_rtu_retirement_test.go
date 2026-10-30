package modbus

import (
	"context"
	"errors"
	"go.bug.st/serial"
	"testing"
)

type cancelDuringClosePort struct {
	serial.Port
	cancel context.CancelFunc
}

func (p *cancelDuringClosePort) Close() error { p.cancel(); return nil }
func TestRTUCancelDuringRetirementDoesNotOpenReplacement(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	tr := NewRTUTransport("/nonexistent-canceled-rtu-port", 9600, 8, 1, "N", 0)
	tr.port = &cancelDuringClosePort{cancel: cancel}
	err := tr.ConnectContext(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("after retiring old handle, canceled connect opened next port instead: %v", err)
	}
}

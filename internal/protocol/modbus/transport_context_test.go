package modbus

import (
	"context"
	"testing"
)

func TestTransportPreCanceledRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	transports := []Transport{NewTCPTransport("127.0.0.1", 1), NewUDPTransport("127.0.0.1", 1), NewRTUTransport("missing", 0, 0, 0, "", 0)}
	for _, tr := range transports {
		bounded, ok := tr.(interface {
			SendReceiveContext(context.Context, []byte) ([]byte, error)
		})
		if !ok {
			t.Errorf("%T lacks context send", tr)
			continue
		}
		if _, err := bounded.SendReceiveContext(ctx, []byte{1}); err != context.Canceled {
			t.Errorf("%T: %v", tr, err)
		}
	}
}

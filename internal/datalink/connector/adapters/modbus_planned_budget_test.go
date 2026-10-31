package adapters

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"testing"
	"time"
)

func TestModbusScalarRetryDefaultsRemainTwoAttempts(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(kind), func(t *testing.T) {
			peer := &modbusPlannedPeer{failure: io.EOF}
			c := modbusPlannedProtocol(kind, peer)
			_, err := c.Read(t.Context(), connector.ReadRequest{Address: "HR0", DataType: schema.DataTypeUint16})
			require.ErrorIs(t, err, io.EOF)
			require.Len(t, peer.frames, 2)
		})
	}
}
func TestModbusPlannedRetryBudgetAndCancellation(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		for _, tc := range []struct {
			name     string
			limit    int
			failure  error
			cancel   bool
			attempts int
		}{
			{"default", 0, io.EOF, false, 2}, {"shared", 3, io.EOF, false, 3}, {"no-retry-text", 4, errors.New("unexpected eof data"), false, 1}, {"canceled", 4, io.EOF, true, 0},
		} {
			t.Run(string(kind)+tc.name, func(t *testing.T) {
				req := modbusPlannedRequest(t, kind, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", "HR0", schema.DataTypeUint16, ""), modbusPlannedPoint("b", "HR1", schema.DataTypeUint16, ""))
				peer := &modbusPlannedPeer{failure: tc.failure}
				reader := requireModbusPlannedReader(t, modbusPlannedProtocol(kind, peer))
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				if tc.limit > 0 {
					ctx = wirebudget.With(ctx, tc.limit, 0)
				}
				if tc.cancel {
					cancel()
				}
				got, err := reader.ReadPlanned(ctx, req)
				require.Error(t, err)
				require.Len(t, peer.frames, tc.attempts)
				if tc.limit > 0 {
					require.Equal(t, tc.attempts, wirebudget.From(ctx).Used())
				}
				require.Len(t, got, 2)
				for _, r := range got {
					require.Equal(t, schema.QualityBad, r.Quality)
					require.Nil(t, r.Value)
				}
			})
		}
	}
}
func TestModbusPlannedDeadlineAndBackoffDoNotSendAgain(t *testing.T) {
	for _, waiting := range []bool{false, true} {
		t.Run(map[bool]string{false: "backoff", true: "read"}[waiting], func(t *testing.T) {
			req := modbusPlannedRequest(t, schema.ProtocolModbusTCP, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", "HR0", schema.DataTypeUint16, ""))
			peer := &modbusPlannedPeer{failure: io.EOF, wait: waiting}
			reader := requireModbusPlannedReader(t, modbusPlannedProtocol(schema.ProtocolModbusTCP, peer))
			ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
			defer cancel()
			ctx = wirebudget.With(ctx, 3, time.Second)
			_, err := reader.ReadPlanned(ctx, req)
			require.ErrorIs(t, err, context.DeadlineExceeded)
			require.Len(t, peer.frames, 1)
			require.Equal(t, 1, wirebudget.From(ctx).Used())
		})
	}
}

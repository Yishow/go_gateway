package adapters

import (
	"context"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// ReadPlanned executes only the selected immutable request through the existing
// client, transport, and retry budget. It never invents a gap or scalar fallback.
func (c *ModbusTCPConnector) ReadPlanned(ctx context.Context, req connector.PlannedReadRequest) ([]connector.PlannedPointResult, error) {
	return c.readPlanned(ctx, req, schema.ProtocolModbusTCP, "tcp", c.config.SlaveID)
}

// ReadPlanned keeps RTU framing and lifecycle in the existing delegated client.
func (c *ModbusRTUConnector) ReadPlanned(ctx context.Context, req connector.PlannedReadRequest) ([]connector.PlannedPointResult, error) {
	defer c.afterOperationContext(ctx)
	return c.readProxy().readPlanned(ctx, req, schema.ProtocolModbusRTU, "rtu", c.config.SlaveID)
}

// ReadPlanned preserves the existing UDP MBAP wire contract, not RTU-over-UDP.
func (c *ModbusUDPConnector) ReadPlanned(ctx context.Context, req connector.PlannedReadRequest) ([]connector.PlannedPointResult, error) {
	defer c.afterOperationContext(ctx)
	return c.readProxy().readPlanned(ctx, req, schema.ProtocolModbusUDP, "tcp", c.config.SlaveID)
}

func (c *ModbusTCPConnector) readPlanned(ctx context.Context, req connector.PlannedReadRequest, protocol schema.ProtocolType, frame string, station byte) ([]connector.PlannedPointResult, error) {
	request, bindings, err := req.Snapshot()
	if err == nil {
		err = validateModbusPlannedRequest(request, bindings, protocol, frame, station)
	}
	if err != nil {
		return modbusPlannedResults(req, request, bindings, nil, nil, time.Now(), err), err
	}
	ctx, cancel := c.readContext(ctx, 2)
	defer cancel()
	if err = c.ensureConnectionContext(ctx); err != nil {
		return modbusPlannedResults(req, request, bindings, nil, nil, time.Now(), err), err
	}
	defer c.afterOperationContext(ctx)
	var words []uint16
	var bits []bool
	err = c.withReadConnectionRetry(ctx, func() error {
		words, bits = nil, nil
		client := c.client.WithContext(ctx)
		function := modbusPlannedFunction(request.Group.Function)
		var readErr error
		// Validation bounded these values to the protocol's uint16 address space.
		address, count := uint16(request.Start), uint16(request.Length) // #nosec G115
		if function < 3 {
			bits, readErr = client.ReadBitsExact(function, address, count)
		} else {
			words, readErr = client.ReadWordsExact(function, address, count)
		}
		return readErr
	})
	completed := time.Now()
	return modbusPlannedResults(req, request, bindings, words, bits, completed, err), err
}

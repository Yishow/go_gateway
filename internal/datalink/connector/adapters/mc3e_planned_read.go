package adapters

import (
	"context"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/protocol/mcprotocol"
)

// ReadPlanned executes one immutable optimizer request, with no implicit scalar
// fallback. All wire attempts use the existing deadline/reconnect/send budget.
func (c *MC3EConnector) ReadPlanned(ctx context.Context, req connector.PlannedReadRequest) ([]connector.PlannedPointResult, error) {
	request, bindings, err := req.Snapshot()
	if err == nil {
		err = validateMCPlannedRequest(request, bindings)
	}
	if err != nil {
		return mcPlannedFailure(req, request, bindings, err), err
	}
	ctx, cancel := c.readContext(ctx)
	defer cancel()
	if err = c.ensureConnectionContext(ctx); err != nil {
		return mcPlannedFailure(req, request, bindings, err), err
	}
	defer c.afterOperation()
	var words []int
	var bits []bool
	err = c.withConnectionRetry(ctx, func() error {
		words, bits = nil, nil
		client := c.client.WithContext(ctx)
		var readErr error
		switch {
		case request.Kind == "random":
			items := make([]mcprotocol.RandomReadItem, len(request.Units))
			for i, u := range request.Units {
				items[i] = mcprotocol.RandomReadItem{Device: u.Area, Addr: u.Address}
			}
			words, readErr = client.RandomReadExact(items)
		case request.Points[0].NativeBits == 1:
			bits, readErr = client.BatchReadBitExact(request.Group.Area, request.Start, request.Length)
		default:
			words, readErr = client.BatchReadWordExact(request.Group.Area, request.Start, request.Length)
		}
		return readErr
	})
	completed := time.Now()
	if err != nil {
		return mcPlannedResults(req, request, bindings, nil, nil, completed, err), err
	}
	return mcPlannedResults(req, request, bindings, words, bits, completed, nil), nil
}
func validateMCPlannedRequest(r optimizer.AcquisitionRequest, bindings []connector.PlannedPointBinding) error {
	if (r.Group.Identity.Protocol != "mc" && r.Group.Identity.Protocol != "mitsubishi") || r.Group.Identity.Frame != "3e-binary" || (r.Kind != "contiguous" && r.Kind != "random") {
		return connector.ErrPlannedRead
	}
	for _, b := range bindings {
		if err := validateMCReadProfile(connector.ReadRequest{Address: b.Point.Address, Function: b.Point.Function, DataType: b.Point.DataType, Profile: b.Profile}); err != nil {
			return err
		}
	}
	return nil
}
func mcPlannedFailure(req connector.PlannedReadRequest, r optimizer.AcquisitionRequest, bindings []connector.PlannedPointBinding, err error) []connector.PlannedPointResult {
	return mcPlannedResults(req, r, bindings, nil, nil, time.Now(), err)
}

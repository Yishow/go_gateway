package adapters

import (
	"context"
	"strconv"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/protocol/fatek"
)

// ReadPlanned executes one frozen optimizer request through the existing FATEK
// client/transport budget. It never creates a fallback or another retry owner.
func (c *FatekConnector) ReadPlanned(ctx context.Context, req connector.PlannedReadRequest) ([]connector.PlannedPointResult, error) {
	request, bindings, err := req.Snapshot()
	if err == nil {
		err = c.validatePlannedRequest(request, bindings)
	}
	if err != nil {
		return fatekPlannedResults(req, request, bindings, fatekPlannedRaw{}, time.Now(), err), err
	}
	ctx, cancel := c.readContext(ctx, 1)
	defer cancel()
	if err = c.ensureConnectionContext(ctx); err != nil {
		return fatekPlannedResults(req, request, bindings, fatekPlannedRaw{}, time.Now(), err), err
	}
	defer c.afterOperationContext(ctx)
	client := c.client.WithContext(ctx)
	raw := fatekPlannedRaw{}
	switch {
	case request.Kind == "random":
		items := make([]fatek.RandomReadItem, len(request.Units))
		for i, u := range request.Units {
			items[i] = fatek.RandomReadItem{Symbol: u.Area, Addr: u.Address}
		}
		raw.random, err = client.ReadRandom(items)
	case request.Points[0].NativeBits == 1:
		raw.bits, err = client.ReadStatus(request.Group.Area, request.Start, request.Length)
	default:
		raw.words, err = client.ReadRegisters(request.Group.Area, request.Start, request.Length)
	}
	completed := time.Now()
	return fatekPlannedResults(req, request, bindings, raw, completed, err), err
}
func (c *FatekConnector) validatePlannedRequest(r optimizer.AcquisitionRequest, bindings []connector.PlannedPointBinding) error {
	station := int(c.config.StationNo)
	if station == 0 {
		station = fatek.DefaultStation
	}
	if r.Group.Identity.Protocol != "fatek" || r.Group.Identity.Frame != "ascii" || r.Group.Identity.Station != strconv.Itoa(station) || (r.Kind != "contiguous" && r.Kind != "random") {
		return connector.ErrPlannedRead
	}
	for _, p := range r.Points {
		if err := validateFatekPlannedSlice(r, p); err != nil {
			return err
		}
	}
	for _, b := range bindings {
		if err := validateFatekReadProfile(connector.ReadRequest{Address: b.Point.Address, Function: b.Point.Function, DataType: b.Point.DataType, Profile: b.Profile}); err != nil {
			return err
		}
	}
	return nil
}

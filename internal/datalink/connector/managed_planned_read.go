package connector

import (
	"context"
	"sync/atomic"

	"go-gateway/internal/protocol/wirebudget"
)

// ReadPlanned forwards one immutable request through the same managed gate as
// Read. A foreign device/config cannot borrow a current connection. Retirement
// and generation replacement remain solely owned by ConnectionManager.
func (mc *ManagedConnection) ReadPlanned(ctx context.Context, req PlannedReadRequest) ([]PlannedPointResult, error) {
	if err := mc.LockContext(ctx); err != nil {
		return nil, err
	}
	defer mc.Unlock()
	request, bindings, err := req.Snapshot()
	if err != nil {
		return nil, err
	}
	if request.Group.Identity.DeviceID != mc.DeviceID || request.Group.Identity.ConnectionRevision != ReadConnectionRevision(mc.Config) {
		return nil, ErrPlannedRead
	}
	for _, b := range bindings {
		if b.Profile.DeviceID() != mc.DeviceID || !b.Profile.MatchesConnection(mc.ProtocolType, mc.Config) {
			return nil, ErrReadProfile
		}
	}
	reader, ok := mc.Protocol.(PlannedReader)
	if !ok {
		return nil, ErrPlannedRead
	}
	var lastAttempt atomic.Pointer[AcquisitionConnectionToken]
	ctx = wirebudget.WithSendObserver(ctx, func() {
		token := mc.acquisitionTokenLocked()
		lastAttempt.Store(&token)
	})
	results, err := reader.ReadPlanned(ctx, req)
	for i := range results {
		// Adapter output cannot claim a connection token without a real send.
		results[i].ConnectionToken = AcquisitionConnectionToken{}
		if token := lastAttempt.Load(); token != nil {
			results[i].ConnectionToken = *token
		}
	}
	if err != nil {
		mc.LastError = err
	}
	return results, err
}

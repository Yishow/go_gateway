package connector

import (
	"context"
	"go-gateway/internal/datalink/schema"
	"time"
)

// Lock 鎖定連線以進行操作
func (mc *ManagedConnection) Lock() {
	mc.mu.Lock()
	mc.inUse = true
}

// LockContext acquires the same serialization gate, but canceled callers do not
// wait behind another request's I/O and never enter the protocol afterward.
func (mc *ManagedConnection) LockContext(ctx context.Context) error {
	if err := mc.mu.LockContext(ctx); err != nil {
		return err
	}
	mc.inUse = true
	return nil
}

// Unlock 解鎖連線
func (mc *ManagedConnection) Unlock() {
	mc.inUse = false
	mc.LastUsed = time.Now()
	mc.mu.Unlock()
}

// Read 透過連線讀取資料 (自動鎖定)
func (mc *ManagedConnection) Read(ctx context.Context, req ReadRequest) (ReadResult, error) {
	if err := mc.LockContext(ctx); err != nil {
		return ReadResult{Quality: schema.QualityBad, Error: err.Error()}, err
	}
	defer mc.Unlock()

	if err := ValidateReadProfile(req); err != nil {
		return ReadResult{Quality: schema.QualityBad, Error: err.Error()}, err
	}
	if !req.Profile.MatchesConnection(mc.ProtocolType, mc.Config) || (req.Profile.Present() && req.Profile.DeviceID() != mc.DeviceID) {
		return ReadResult{Quality: schema.QualityBad, Error: ErrReadProfile.Error()}, ErrReadProfile
	}
	// Only protocols whose production decoder consumes the profile may read
	// modern intent.
	if req.Profile.Present() && req.Profile.Resolution().Mode != schema.PointFormatLegacy {
		switch mc.ProtocolType {
		case schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP, schema.ProtocolMC3E, schema.ProtocolFatekFBs:
		default:
			return ReadResult{Quality: schema.QualityBad, Error: ErrReadProfile.Error()}, ErrReadProfile
		}
	}
	result, err := mc.Protocol.Read(ctx, req)
	if err != nil {
		mc.LastError = err
	}
	return result, err
}

// Write 透過連線寫入資料 (自動鎖定)
func (mc *ManagedConnection) Write(ctx context.Context, req WriteRequest) error {
	if err := mc.LockContext(ctx); err != nil {
		return err
	}
	defer mc.Unlock()

	err := mc.Protocol.Write(ctx, req)
	if err != nil {
		mc.LastError = err
	}
	return err
}

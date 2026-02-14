package connector

import (
	"context"
	"time"
)

// Lock 鎖定連線以進行操作
func (mc *ManagedConnection) Lock() {
	mc.mu.Lock()
	mc.inUse = true
}

// Unlock 解鎖連線
func (mc *ManagedConnection) Unlock() {
	mc.inUse = false
	mc.LastUsed = time.Now()
	mc.mu.Unlock()
}

// Read 透過連線讀取資料 (自動鎖定)
func (mc *ManagedConnection) Read(ctx context.Context, req ReadRequest) (ReadResult, error) {
	mc.Lock()
	defer mc.Unlock()

	result, err := mc.Protocol.Read(ctx, req)
	if err != nil {
		mc.LastError = err
	}
	return result, err
}

// Write 透過連線寫入資料 (自動鎖定)
func (mc *ManagedConnection) Write(ctx context.Context, req WriteRequest) error {
	mc.Lock()
	defer mc.Unlock()

	err := mc.Protocol.Write(ctx, req)
	if err != nil {
		mc.LastError = err
	}
	return err
}

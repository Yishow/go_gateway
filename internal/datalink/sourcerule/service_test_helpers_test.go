package sourcerule

import (
	"context"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
)

func seedActivationReadyDevice(ctx context.Context, repo *device.MemoryRepository, id string) (*schema.Device, error) {
	record, err := seedDeviceWithStatus(ctx, repo, id, schema.DeviceStatusActive)
	if err != nil {
		return nil, err
	}
	success := true
	record.LastTestSuccess = &success
	record.LastTestError = ""
	if err := repo.Update(ctx, record); err != nil {
		return nil, err
	}
	return repo.GetByID(ctx, id)
}

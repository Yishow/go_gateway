package point

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestServiceCreate_AllowsMCHexWordAndBitAddressShapes(t *testing.T) {
	svc := NewService(NewMemoryRepository(), nil)
	for _, address := range []string{"WA", "W1A", "BFF", "B100"} {
		t.Run(address, func(t *testing.T) {
			created, err := svc.Create(t.Context(), CreatePointRequest{
				DeviceID: "device-mc-hex", Name: "MC_" + address, Address: address,
				DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly,
			})
			require.NoError(t, err)
			require.Equal(t, address, created.Address)
		})
	}
}

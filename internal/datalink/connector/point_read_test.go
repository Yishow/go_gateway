package connector

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestPointReadCountForAddress(t *testing.T) {
	tests := []struct {
		name     string
		protocol schema.ProtocolType
		address  string
		dataType schema.DataType
		want     int
		wantErr  bool
	}{
		{name: "FATEK 16-bit float32", protocol: schema.ProtocolFatekFBs, address: "D0", dataType: schema.DataTypeFloat32, want: 2},
		{name: "FATEK 16-bit uint64", protocol: schema.ProtocolFatekFBs, address: "R0", dataType: schema.DataTypeUint64, want: 4},
		{name: "FATEK 32-bit float32", protocol: schema.ProtocolFatekFBs, address: "DR0", dataType: schema.DataTypeFloat32, want: 1},
		{name: "FATEK 32-bit uint64", protocol: schema.ProtocolFatekFBs, address: "DD0", dataType: schema.DataTypeUint64, want: 2},
		{name: "FATEK scalar int16", protocol: schema.ProtocolFatekFBs, address: "D0", dataType: schema.DataTypeInt16, want: 1},
		{name: "FATEK scalar bool", protocol: schema.ProtocolFatekFBs, address: "X0", dataType: schema.DataTypeBool, want: 1},
		{name: "MC typed count", protocol: schema.ProtocolMC3E, address: "D3", dataType: schema.DataTypeUint64, want: 1},
		{name: "Modbus word count", protocol: schema.ProtocolModbusTCP, address: "40001", dataType: schema.DataTypeUint64, want: 4},
		{name: "FATEK unsupported symbol", protocol: schema.ProtocolFatekFBs, address: "RT0", dataType: schema.DataTypeFloat32, wantErr: true},
		{name: "FATEK address out of bounds", protocol: schema.ProtocolFatekFBs, address: "D100000", dataType: schema.DataTypeInt16, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PointReadCountForAddress(tt.protocol, tt.address, tt.dataType)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

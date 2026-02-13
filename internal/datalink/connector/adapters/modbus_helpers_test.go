package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseModbusAddress(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		function    string
		wantAddr    uint16
		wantFunc    string
		expectError bool
	}{
		{name: "holding register traditional", address: "40001", wantAddr: 0, wantFunc: "03"},
		{name: "input register traditional", address: "30001", wantAddr: 0, wantFunc: "04"},
		{name: "discrete input traditional", address: "10001", wantAddr: 0, wantFunc: "02"},
		{name: "coil traditional", address: "1", wantAddr: 0, wantFunc: "01"},
		{name: "holding register alias", address: "HR100", wantAddr: 100, wantFunc: "03"},
		{name: "input register alias", address: "IR7", wantAddr: 7, wantFunc: "04"},
		{name: "coil alias", address: "CO8", wantAddr: 8, wantFunc: "01"},
		{name: "direct address with function", address: "60000", function: "16", wantAddr: 60000, wantFunc: "16"},
		{name: "invalid address", address: "bad", expectError: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			addr, fn, err := parseModbusAddress(tt.address, tt.function)
			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantAddr, addr)
			assert.Equal(t, tt.wantFunc, fn)
		})
	}
}

func TestToUint16Slice(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		want        []uint16
		expectError bool
	}{
		{name: "uint16 slice", input: []uint16{1, 2}, want: []uint16{1, 2}},
		{name: "int slice", input: []int{3, 4}, want: []uint16{3, 4}},
		{name: "interface slice", input: []interface{}{float64(5), int(6)}, want: []uint16{5, 6}},
		{name: "unsupported type", input: "x", expectError: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := toUint16Slice(tt.input)
			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

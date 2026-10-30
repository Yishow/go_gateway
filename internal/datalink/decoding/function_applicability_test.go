package decoding

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestPointFormatApplicabilityMatchesModbusFunctionAddressPrecedence(t *testing.T) {
	for _, tc := range []struct {
		name, address, function string
		ok                      bool
	}{
		{"raw03", "100", "03", true}, {"raw04", "0", "04", true}, {"raw01", "50000", "01", false}, {"raw02", "50000", "02", false},
		{"holding-alias", "HR100", "01", true}, {"input-alias", "IR100", "01", true}, {"coil-alias", "CO100", "03", false},
		{"traditional-holding", "40001", "01", true}, {"traditional-input", "30001", "01", true}, {"traditional-bit", "10001", "03", false},
		{"short-no-function", "100", "", false}, {"raw-no-function", "50000", "", true}, {"invalid-function", "50000", "99", false},
		{"raw-span-overflow", "65535", "03", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := schema.Point{Address: tc.address, Function: tc.function, DataType: schema.DataTypeInt32, FormatMode: schema.PointFormatOverride, FormatOverride: "ABCD"}
			_, err := ResolvePoint(p, FormatContext{Protocol: schema.ProtocolModbusTCP})
			if tc.ok {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, ErrInapplicableFormat)
			}
		})
	}
}

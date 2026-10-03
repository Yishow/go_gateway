package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroupSQLTypesCompatiblePreservesNumericPrecisionAndScale(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		actual   string
		want     bool
	}{
		{name: "bounded exact", expected: "NUMERIC(20,0)", actual: "numeric(20,0)", want: true},
		{name: "unbounded exact custom codec", expected: "NUMERIC(20,0)", actual: "NUMERIC", want: true},
		{name: "fractional scale", expected: "NUMERIC(20,0)", actual: "NUMERIC(20,2)", want: false},
		{name: "narrow precision", expected: "NUMERIC(20,0)", actual: "NUMERIC(18,0)", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, groupSQLTypesCompatible(tt.expected, tt.actual))
		})
	}
}

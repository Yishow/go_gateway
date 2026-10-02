package workspace

import (
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupReadinessUsesConservativeColumnTypeCompatibility(t *testing.T) {
	tests := []struct {
		name       string
		sourceType schema.DataType
		columnType string
		compatible bool
	}{
		{name: "float does not target integer", sourceType: schema.DataTypeFloat32, columnType: "INTEGER"},
		{name: "integer does not target real", sourceType: schema.DataTypeInt64, columnType: "REAL"},
		{name: "string does not target uuid", sourceType: schema.DataTypeString, columnType: "UUID"},
		{name: "string does not target json", sourceType: schema.DataTypeString, columnType: "JSON"},
		{name: "integer targets numeric", sourceType: schema.DataTypeInt64, columnType: "NUMERIC(20,6)", compatible: true},
		{name: "float targets decimal", sourceType: schema.DataTypeFloat32, columnType: "DECIMAL(20,6)", compatible: true},
		{name: "string targets text", sourceType: schema.DataTypeString, columnType: "TEXT", compatible: true},
		{name: "bool targets integer", sourceType: schema.DataTypeBool, columnType: "INTEGER", compatible: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			compatible, known := writeGroupColumnTypeCompatible(test.sourceType, dbtarget.ColumnInfo{DataType: test.columnType})
			require.True(t, known)
			require.Equal(t, test.compatible, compatible)
		})
	}
}

package dbtarget

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildWriteStatement_MySQLUpsertUsesDuplicateKeyUpdate(t *testing.T) {
	mapping := &schema.DatabaseTargetMapping{
		TableSchema:     "gateway_metrics",
		TableName:       "sensor_values",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtr("ts"),
	}

	query, args, err := buildWriteStatement(
		schema.DatabaseConnectorKindMySQL,
		mapping,
		42.5,
		time.Date(2026, 3, 16, 12, 45, 0, 0, time.UTC),
	)
	require.NoError(t, err)
	assert.Equal(t, "INSERT INTO `gateway_metrics`.`sensor_values` (`value`, `ts`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `value` = VALUES(`value`)", query)
	assert.Len(t, args, 2)
}

func TestBuildGroupedWriteStatement_MySQLUpsertUsesDuplicateKeyUpdate(t *testing.T) {
	query, args, err := buildGroupedWriteStatement(
		schema.DatabaseConnectorKindMySQL,
		groupedWriteKey{
			SchemaName:      "gateway_metrics",
			TableName:       "meter_rows",
			WriteMode:       schema.DatabaseWriteModeUpsert,
			TimestampColumn: "ts",
			BucketStart:     time.Date(2026, 3, 16, 12, 45, 0, 0, time.UTC),
		},
		map[string]any{"kw": 7.25, "a1": 42.5},
	)
	require.NoError(t, err)
	assert.Equal(t, "INSERT INTO `gateway_metrics`.`meter_rows` (`ts`, `a1`, `kw`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `a1` = VALUES(`a1`), `kw` = VALUES(`kw`)", query)
	assert.Len(t, args, 3)
}

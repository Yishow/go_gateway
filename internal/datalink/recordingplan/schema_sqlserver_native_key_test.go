package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerNativePartitionRequiresTimeInUniqueKey(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*GroupSchemaColumn)
	}{
		{"partition column missing from primary key", func(column *GroupSchemaColumn) { column.PrimaryKey = false }},
		{"partition column has incompatible type", func(column *GroupSchemaColumn) { column.SQLType = "BIGINT" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			layout := sqlServerNativeMonthlyLayoutForTest(t)
			for i := range layout.Columns {
				if layout.Columns[i].Name == "bucket_start" {
					test.mutate(&layout.Columns[i])
				}
			}
			_, err := GenerateGroupSchemaDDL(dialectSQLServer, layout)
			require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
		})
	}
}

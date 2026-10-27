package runtime

import (
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestNativeMonthlyBoundaryAllowsSQLServer(t *testing.T) {
	f, _, _ := sqlServerScaledIntakeFixture(t, false, "")
	f.config.Group.StorageMode = workspace.WriteGroupStorageModeNativeMonthly
	boundary, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err, "SQL Server native monthly must activate under managed_long_v1")
	require.NotNil(t, boundary)
}

func TestNativeMonthlyBoundaryRejectsOtherDialects(t *testing.T) {
	for _, dialect := range []dbtarget.SQLDialect{dbtarget.SQLDialectSQLite, dbtarget.SQLDialectPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			f := newLongFixture(t)
			f.config.Group.StorageMode = workspace.WriteGroupStorageModeNativeMonthly
			f.config.Dialect = dialect
			_, err := NewGroupBoundaryContext(t.Context(), f.config)
			require.Error(t, err, "a direct boundary caller cannot reinterpret native monthly as %s", dialect)
		})
	}
}

package sourcerule

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
)

func TestDerivedPointNoopSynchronizationDoesNotWriteConfiguration(t *testing.T) {
	db := setupSQLRepoDB(t)
	t.Cleanup(func() { _ = db.Close() })
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(t.Context(), dr, "d"))
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	svc := NewService(NewSQLRepository(db), device.NewService(dr, nil), ps, nil)
	rule, err := svc.Create(t.Context(), CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeUint32, NamingPrefix: "stable", Enabled: true})
	require.NoError(t, err)
	links, err := svc.ListLinks(t.Context(), rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)
	before, err := ps.GetByID(t.Context(), links[0].PointID)
	require.NoError(t, err)
	_, err = db.Exec(`CREATE TRIGGER reject_noop_point_write BEFORE UPDATE ON points BEGIN SELECT RAISE(ABORT,'no-op must not write'); END`)
	require.NoError(t, err)
	require.NoError(t, svc.SyncDerivedPointState(t.Context()))
	after, err := ps.GetByID(t.Context(), before.ID)
	require.NoError(t, err)
	require.Equal(t, before, after, "including exact revision and timestamps")
}

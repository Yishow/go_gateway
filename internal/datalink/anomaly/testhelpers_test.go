package anomaly

import (
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func openAnomalyTestDB(t *testing.T, deviceIDs ...string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+t.Name()+"?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	repo := device.NewSQLRepository(db)
	for _, id := range deviceIDs {
		require.NoError(t, repo.Create(t.Context(), &schema.Device{ID: id, Name: id, Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusDraft, ConnectionConfig: "{}"}))
	}
	return db
}

type testClock struct{ now time.Time }

func (c *testClock) Now() time.Time { return c.now }

func (c *testClock) Advance(d time.Duration) { c.now = c.now.Add(d) }

var anomalyEpoch = time.Date(2026, 11, 5, 8, 0, 0, 0, time.UTC)

package datalink

import (
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// This is configuration persistence evidence, not adapter/collector wiring.
func TestPointFormatIntentRoundtripDefaultChangeAndReset(t *testing.T) {
	path := filepath.Join(t.TempDir(), "intent.sqlite")
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	require.NoError(t, NewMigrator().Migrate(db))
	_, err = db.ExecContext(t.Context(), `INSERT INTO devices(id,name,protocol,status,connection_config) VALUES('device','device','modbus_tcp','active','{}')`)
	require.NoError(t, err)
	repo := point.NewSQLRepository(db)
	service := point.NewService(repo, point.NewMemoryPollingGroupRepository())
	a, err := service.Create(t.Context(), point.CreatePointRequest{DeviceID: "device", Name: "A", Address: "HR100", DataType: schema.DataTypeUint32, DataFormat: "ABCD"})
	require.NoError(t, err)
	b, err := service.Create(t.Context(), point.CreatePointRequest{DeviceID: "device", Name: "B", Address: "HR102", DataType: schema.DataTypeUint32, DataFormat: "ABCD"})
	require.NoError(t, err)
	override, format := schema.PointFormatOverride, "ABCD"
	candidate, err := decoding.WithIntent(*a, decoding.IntentEdit{Mode: &override, ConfiguredFormat: &format})
	require.NoError(t, err)
	require.NoError(t, repo.Update(t.Context(), &candidate))
	defaults := decoding.FormatContext{Protocol: schema.ProtocolModbusTCP, RuleOwned: true, RuleFormat: "ABCD"}
	assertEffective := func(id, format string, mode schema.PointFormatMode) {
		t.Helper()
		stored, err := repo.GetByID(t.Context(), id)
		require.NoError(t, err)
		got, err := decoding.ResolvePoint(*stored, defaults)
		require.NoError(t, err)
		require.Equal(t, format, got.EffectiveFormat)
		require.Equal(t, mode, got.Mode)
	}
	assertEffective(a.ID, "ABCD", override)
	assertEffective(b.ID, "ABCD", schema.PointFormatInherit)
	name := "A renamed"
	_, err = service.Update(t.Context(), a.ID, point.UpdatePointRequest{Name: &name})
	require.NoError(t, err)
	require.NoError(t, db.Close())
	db, err = sql.Open("sqlite", path)
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	require.NoError(t, NewMigrator().Migrate(db))
	repo = point.NewSQLRepository(db)
	defaults.RuleFormat = "CDAB"
	assertEffective(a.ID, "ABCD", override)
	assertEffective(b.ID, "CDAB", schema.PointFormatInherit)
	stored, err := repo.GetByID(t.Context(), a.ID)
	require.NoError(t, err)
	require.Equal(t, "ABCD", stored.DataFormat)
	require.Equal(t, name, stored.Name)
	inherit := schema.PointFormatInherit
	reset, err := decoding.WithIntent(*stored, decoding.IntentEdit{Mode: &inherit})
	require.NoError(t, err)
	require.NoError(t, repo.Update(t.Context(), &reset))
	assertEffective(a.ID, "CDAB", inherit)
	assertEffective(b.ID, "CDAB", inherit)
	require.NoError(t, NewMigrator().Migrate(db))
	assertEffective(a.ID, "CDAB", inherit)
}

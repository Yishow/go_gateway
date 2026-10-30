package dbtarget

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
	"testing"
)

func TestMixedLayoutsFreezeCallerOwnedMetadata(t *testing.T) {
	s := mixedWideSpec(SQLDialectSQLite)
	mod := measurement.NewUint64(65536)
	rate := measurement.NewUint64(100)
	s.Members[2].Counter.RolloverModulus = &mod
	s.Members[2].Counter.MaxRatePerSecond = &rate
	wide, issues := NewGroupRowLayout(s)
	require.Empty(t, issues)
	s.Members[0].Calculation = snapshot.CalculationLatest
	s.Members[2].Counter.Unit = "changed"
	mod = measurement.NewUint64(16)
	row, err := wide.EncodeRow(mixedEncodeFixture(t, snapshot.ModeMixedWide)[0])
	require.NoError(t, err)
	var entries []map[string]any
	require.NoError(t, json.Unmarshal([]byte(cellMap(row)["provenance"].(string)), &entries))
	c := entries[2]["counter"].(map[string]any)
	require.Equal(t, "kWh", c["unit"])
	require.Equal(t, "65536", c["rollover_modulus"])
	ls := mixedLongSpec(t, SQLDialectSQLite)
	long, issues := NewLongRowLayout(ls)
	require.Empty(t, issues)
	ls.Members[2].Counter.Unit = "changed"
	row, err = long.EncodeRow(mixedEncodeFixture(t, snapshot.ModeMixedLong)[2])
	require.NoError(t, err)
	require.Equal(t, "kWh", cellMap(row)["unit"])
}

package dbtarget

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/snapshot"
	"testing"
)

func TestMixedProvenanceSeparatesRawResultAndMemberWindows(t *testing.T) {
	wide, issues := NewGroupRowLayout(mixedWideSpec(SQLDialectSQLite))
	require.Empty(t, issues)
	row, err := wide.EncodeRow(mixedEncodeFixture(t, snapshot.ModeMixedWide)[0])
	require.NoError(t, err)
	var entries []map[string]any
	require.NoError(t, json.Unmarshal([]byte(cellMap(row)["provenance"].(string)), &entries))
	require.Len(t, entries, 3)
	for i, e := range entries {
		require.Equal(t, MixedResultShapeV1, e["result_shape"])
		require.Equal(t, []string{"mean", "mean", "delta"}[i], e["calculation"])
		require.Equal(t, mixedEncodeAt(120).Format(longTimeText), e["due_at"])
		require.Equal(t, "utc_epoch", e["alignment"])
		require.NotNil(t, e["raw"])
		require.NotNil(t, e["window"])
	}
	w := entries[0]["window"].(map[string]any)
	require.Equal(t, mixedEncodeAt(105).Format(longTimeText), w["aligned_start"])
	w = entries[2]["window"].(map[string]any)
	require.Equal(t, mixedEncodeAt(60).Format(longTimeText), w["aligned_start"])
	counter := entries[2]["counter"].(map[string]any)
	require.Equal(t, "8", counter["delta"])
	require.Equal(t, mixedEncodeAt(59).Format(longTimeText), counter["delta_from"])
	require.Equal(t, mixedEncodeAt(119).Format(longTimeText), counter["delta_to"])
	raw := entries[2]["raw"].(map[string]any)
	require.Contains(t, raw, "value")
	require.Equal(t, "1008", raw["value"])
	long, issues := NewLongRowLayout(mixedLongSpec(t, SQLDialectSQLite))
	require.Empty(t, issues)
	row, err = long.EncodeRow(mixedEncodeFixture(t, snapshot.ModeMixedLong)[2])
	require.NoError(t, err)
	cells := cellMap(row)
	require.Equal(t, "8", cells["value_decimal"])
	require.Nil(t, cells["delta_value"], "mixed result slot is authoritative; old legacy delta column must not imply another selected value")
	var e map[string]any
	require.NoError(t, json.Unmarshal([]byte(cells["provenance"].(string)), &e))
	require.Equal(t, "delta", e["calculation"])
	require.Nil(t, e["profile"], "no fake common profile interval")
	require.Equal(t, "8", e["counter"].(map[string]any)["delta"])
}

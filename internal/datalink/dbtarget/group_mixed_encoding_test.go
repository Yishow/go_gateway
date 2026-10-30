package dbtarget

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/snapshot"
	"testing"
)

func TestMixedWideEncoderStoresMeanMeanDeltaNotRaw(t *testing.T) {
	l, issues := NewGroupRowLayout(mixedWideSpec(SQLDialectSQLite))
	require.Empty(t, issues)
	out := mixedEncodeFixture(t, snapshot.ModeMixedWide)[0]
	row, err := l.EncodeRow(out)
	require.NoError(t, err)
	cells := cellMap(row)
	require.Equal(t, "22", cells["a"])
	require.Equal(t, "12", cells["b"])
	require.Equal(t, "8", cells["c"])
	require.Equal(t, mixedEncodeAt(120), row.DueAt)
	require.True(t, row.BucketStart.IsZero())
	require.Equal(t, mixedEncodeAt(120).Format("2006-01-02T15:04:05.999999999Z07:00"), cells["due_at"])
}
func TestMixedLongEncoderStoresSelectedResultSlots(t *testing.T) {
	l, issues := NewLongRowLayout(mixedLongSpec(t, SQLDialectSQLite))
	require.Empty(t, issues)
	rows := mixedEncodeFixture(t, snapshot.ModeMixedLong)
	require.Len(t, rows, 3)
	for i, out := range rows {
		row, err := l.EncodeRow(out)
		require.NoError(t, err)
		cells := cellMap(row)
		require.Equal(t, []string{"22", "12", "8"}[i], cells["value_decimal"])
		require.Nil(t, cells["value_int"])
		require.Equal(t, "mixed", cells["profile_kind"])
	}
}

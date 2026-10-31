package snapshot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMixedDispatchIdentityPreservesPublishedBytes(t *testing.T) {
	// Known pre-extraction v1 tuple hash, independently fixed as a literal.
	at := time.Date(2026, 10, 30, 12, 2, 0, 0, time.UTC)
	got, err := MixedDispatchRecordID("w", "g", "r", at)
	require.NoError(t, err)
	require.Equal(t, "a64e4192d18006f2c86358448b0504ebf379b7596cdc595f5c43cfb78a9ceb80", got)
	offset := at.In(time.FixedZone("+8", 8*3600))
	same, err := MixedDispatchRecordID("w", "g", "r", offset)
	require.NoError(t, err)
	require.Equal(t, got, same)
}

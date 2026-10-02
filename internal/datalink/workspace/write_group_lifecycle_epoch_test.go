package workspace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGroupLifecycleCutoverUsesUnixEpochBuckets(t *testing.T) {
	for _, test := range []struct {
		name string
		now  time.Time
		want time.Time
	}{
		{name: "inside", now: time.Unix(1, 900_000_000), want: time.Unix(7, 0)},
		{name: "on_boundary", now: time.Unix(7, 0), want: time.Unix(14, 0)},
		{name: "before_epoch", now: time.Unix(-1, 500_000_000), want: time.Unix(0, 0)},
		{name: "local_zone", now: time.Unix(1, 0).In(time.FixedZone("local", 8*60*60)), want: time.Unix(7, 0)},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := nextWriteGroupCutover(test.now, 0, 7)
			require.NoError(t, err)
			require.Equal(t, test.want.UTC(), got)
		})
	}
	got, err := nextWriteGroupCutover(time.Unix(15, 0), 7, 11)
	require.NoError(t, err)
	require.Equal(t, time.Unix(77, 0).UTC(), got)
}

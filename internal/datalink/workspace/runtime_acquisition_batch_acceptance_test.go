package workspace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func acquisitionTwoGraphs(t *testing.T) (acquisitionSQLFixture, []AcquisitionSourceGraph) {
	t.Helper()
	f := newAcquisitionSQLFixture(t)
	_, err := f.db.ExecContext(t.Context(), `INSERT INTO devices(id,name,description,protocol,status) VALUES('e','e','','modbus_tcp','active'); UPDATE system_settings SET value='{"id":"w","ordered_device_ids":["d","e"]}' WHERE key='studio_v2_workspace'`)
	require.NoError(t, err)
	projection, d := f.graph(t)
	e, ok := projection.AcquisitionSourceGraph("e")
	require.True(t, ok, e.Reason())
	return f, []AcquisitionSourceGraph{d, e}
}

func TestRuntimeAcquisitionBatchAcceptsOnceAndRejectsAnyStaleMember(t *testing.T) {
	for _, stale := range []bool{false, true} {
		t.Run(map[bool]string{false: "current", true: "second-stale"}[stale], func(t *testing.T) {
			f, graphs := acquisitionTwoGraphs(t)
			if stale {
				_, err := f.db.ExecContext(t.Context(), `UPDATE devices SET connection_config='{"timeout":2}' WHERE id='e'`)
				require.NoError(t, err)
			}
			called := 0
			err := f.service.AcceptAcquisitionSourceGraphs(t.Context(), graphs, func() error { called++; return nil })
			if stale {
				require.ErrorIs(t, err, ErrAcquisitionSourceGraph)
				require.Zero(t, called)
			} else {
				require.NoError(t, err)
				require.Equal(t, 1, called)
			}
		})
	}
}

func TestRuntimeAcquisitionBatchRejectsEmptyDuplicateForeignAndMalformed(t *testing.T) {
	for _, kind := range []string{"empty", "duplicate", "zero", "foreign-owner", "foreign-workspace", "malformed"} {
		t.Run(kind, func(t *testing.T) {
			f, graphs := acquisitionTwoGraphs(t)
			switch kind {
			case "empty":
				graphs = nil
			case "duplicate":
				graphs[1] = graphs[0]
			case "zero":
				graphs[1] = AcquisitionSourceGraph{}
			default:
				copied := *graphs[1].data
				graphs[1].data = &copied
				switch kind {
				case "foreign-owner":
					copied.owner = NewSQLRepository(f.db)
				case "foreign-workspace":
					copied.workspaceID = "foreign"
				case "malformed":
					copied.witness = ""
				}
			}
			require.ErrorIs(t, f.service.AcceptAcquisitionSourceGraphs(t.Context(), graphs, func() error { t.Fatal("invalid batch must not publish any member"); return nil }), ErrAcquisitionSourceGraph)
		})
	}
}

func TestRuntimeAcquisitionBatchCancellationAfterStageRequiresOwnedCleanup(t *testing.T) {
	f, graphs := acquisitionTwoGraphs(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	calls := 0
	err := f.service.AcceptAcquisitionSourceGraphs(ctx, graphs, func() error { calls++; cancel(); return nil })
	require.Equal(t, 1, calls)
	require.Error(t, err)
}

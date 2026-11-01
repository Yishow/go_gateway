package workspace

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeAcquisitionAcceptanceRechecksOriginalSourceWitness(t *testing.T) {
	for _, tc := range []struct{ name, query string }{
		{"device-config", `UPDATE devices SET connection_config='{"timeout":9}' WHERE id='d'`},
		{"device-status", `UPDATE devices SET status='disabled' WHERE id='d'`},
		{"point-revision", `UPDATE points SET revision_id='new-installation' WHERE id='a'`},
		{"disabled-point", `UPDATE points SET enabled=0 WHERE id='a'`},
		{"skip-content", `UPDATE source_rules SET skipped_addresses='["40002","40004"]'`},
		{"link-receipt", `UPDATE source_rule_links SET id='different-receipt' WHERE id='la'`},
		{"new-standalone-disabled", `INSERT INTO points(id,device_id,name,address,data_type,enabled) VALUES('new','d','new','40002','uint16',0)`},
		{"new-skipped-only-rule", `INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,skipped_addresses) VALUES('new','d','40004',1,'uint16','n','["40004"]')`},
		{"membership", `UPDATE system_settings SET value='{"id":"w","ordered_device_ids":[]}' WHERE key='studio_v2_workspace'`},
		{"workspace-identity", `UPDATE system_settings SET value='{"id":"foreign","ordered_device_ids":["d"]}' WHERE key='studio_v2_workspace'`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newAcquisitionSQLFixture(t)
			_, graph := f.graph(t)
			_, err := f.db.ExecContext(t.Context(), tc.query)
			require.NoError(t, err)
			called := false
			err = f.service.AcceptAcquisitionSourceGraph(t.Context(), graph, func() error { called = true; return nil })
			require.ErrorIs(t, err, ErrAcquisitionSourceGraph)
			require.False(t, called, "stale graph must not rebase during acceptance")
		})
	}
}

func TestRuntimeAcquisitionAcceptancePinsSQLWritersDuringStage(t *testing.T) {
	f := newAcquisitionSQLFixture(t)
	_, graph := f.graph(t)
	other, err := sql.Open("sqlite", f.path+"?_pragma=busy_timeout(1)")
	require.NoError(t, err)
	defer other.Close()
	called := false
	err = f.service.AcceptAcquisitionSourceGraph(t.Context(), graph, func() error {
		called = true
		_, blocked := other.ExecContext(t.Context(), `UPDATE points SET enabled=0 WHERE id='a'`)
		require.Error(t, blocked, "zero-row reservation must pin the canonical source writer")
		return nil
	})
	require.NoError(t, err)
	require.True(t, called)
	_, err = other.ExecContext(t.Context(), `UPDATE points SET enabled=0 WHERE id='a'`)
	require.NoError(t, err, "reservation must be released after acceptance")
	require.ErrorIs(t, f.service.AcceptAcquisitionSourceGraph(t.Context(), graph, func() error { t.Fatal("stale stage"); return nil }), ErrAcquisitionSourceGraph)
}

func TestRuntimeAcquisitionAcceptanceFailureNeverReportsSuccess(t *testing.T) {
	f := newAcquisitionSQLFixture(t)
	_, graph := f.graph(t)
	stageFailure := errors.New("owned CAS refused")
	require.ErrorIs(t, f.service.AcceptAcquisitionSourceGraph(t.Context(), graph, func() error { return stageFailure }), stageFailure)
	require.ErrorIs(t, f.service.AcceptAcquisitionSourceGraph(t.Context(), AcquisitionSourceGraph{}, func() error { t.Fatal("zero graph stage"); return nil }), ErrAcquisitionSourceGraph)
	require.ErrorIs(t, NewSQLRepository(f.db).AcceptAcquisitionSourceGraph(t.Context(), graph, func() error { t.Fatal("foreign owner stage"); return nil }), ErrAcquisitionSourceGraph)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	called := false
	err := f.service.AcceptAcquisitionSourceGraph(ctx, graph, func() error { called = true; cancel(); return nil })
	require.True(t, called)
	require.Error(t, err, "cancellation after staged publication must force caller-owned cleanup")
}

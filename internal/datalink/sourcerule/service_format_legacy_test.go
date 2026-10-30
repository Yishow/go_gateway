package sourcerule

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
)

func TestRuleFormatLegacyPolicySurvivesSQLRestart(t *testing.T) {
	for _, tc := range []struct {
		protocol schema.ProtocolType
		address  string
	}{
		{schema.ProtocolModbusTCP, "40001"}, {schema.ProtocolMC3E, "D0"}, {schema.ProtocolFatekFBs, "D0"},
	} {
		t.Run(string(tc.protocol), func(t *testing.T) {
			ctx := t.Context()
			db := setupSQLRepoDB(t)
			t.Cleanup(func() { _ = db.Close() })
			dr := device.NewSQLRepository(db)
			require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "legacy-device"))
			d, err := dr.GetByID(ctx, "legacy-device")
			require.NoError(t, err)
			d.Protocol = tc.protocol
			ok := true
			d.LastTestSuccess = &ok
			require.NoError(t, dr.Update(ctx, d))
			require.NoError(t, dr.UpdateTestResult(ctx, d.ID, true, ""))
			ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
			svc := NewService(NewSQLRepository(db), device.NewService(dr, nil), ps, nil)
			r, err := svc.Create(ctx, CreateRuleRequest{DeviceID: d.ID, StartAddress: tc.address, Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "LEGACY", DataFormat: "ABCD"})
			require.NoError(t, err)
			links, err := svc.ListLinks(ctx, r.ID)
			require.NoError(t, err)
			p, err := ps.GetByID(ctx, links[0].PointID)
			require.NoError(t, err)
			p.FormatMode = schema.PointFormatLegacy
			p.FormatVersion = 0
			p.FormatOverride = ""
			p.DataFormat = "BADC"
			require.NoError(t, ps.Restore(ctx, p))
			defaults := decoding.FormatContext{Protocol: tc.protocol, RuleOwned: true, RuleFormat: "ABCD", ConnectionFormat: "CDAB"}
			before, err := decoding.ResolvePoint(*p, defaults)
			require.NoError(t, err)
			_, err = svc.Update(ctx, r.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("NEW"), DataFormat: common.Ptr("DCBA")})
			require.NoError(t, err)
			require.NoError(t, svc.Enable(ctx, r.ID))
			require.NoError(t, svc.Disable(ctx, r.ID))
			restartedPS := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
			restarted := NewService(NewSQLRepository(db), device.NewService(device.NewSQLRepository(db), nil), restartedPS, nil)
			require.NoError(t, restarted.SyncRuleDerivedState(ctx, r.ID))
			got, err := restartedPS.GetByID(ctx, p.ID)
			require.NoError(t, err)
			require.Equal(t, p.FormatMode, got.FormatMode)
			require.Equal(t, p.FormatVersion, got.FormatVersion)
			require.Equal(t, p.FormatOverride, got.FormatOverride)
			require.Equal(t, p.DataFormat, got.DataFormat)
			defaults.RuleFormat = "DCBA"
			after, err := decoding.ResolvePoint(*got, defaults)
			require.NoError(t, err)
			require.Equal(t, before.Policy, after.Policy)
			require.Equal(t, before.EffectiveFormat, after.EffectiveFormat)
			require.Equal(t, before.NativeCompatibility, after.NativeCompatibility)
			require.NotEqual(t, p.RevisionID, got.RevisionID)
		})
	}
}

func TestRuleFormatRollbackRestoresExactIntentRevision(t *testing.T) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	t.Cleanup(func() { _ = db.Close() })
	ps := point.NewService(point.NewSQLRepository(db), nil)
	p, err := ps.Create(ctx, point.CreatePointRequest{DeviceID: "d", Name: "before", Address: "40001", DataType: schema.DataTypeInt32, DataFormat: "ABCD"})
	require.NoError(t, err)
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "ABCD"
	require.NoError(t, ps.Restore(ctx, p))
	applied, err := ps.Update(ctx, p.ID, point.UpdatePointRequest{Name: common.Ptr("after")})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE points SET last_error=? WHERE id=?`, "new-telemetry", p.ID)
	require.NoError(t, err)
	svc := NewService(NewSQLRepository(db), nil, ps, nil)
	require.NoError(t, svc.rollbackUpdatedPoints(ctx, []pointUpdatePlan{{point: p, appliedRevision: applied.RevisionID}}))
	got, err := ps.GetByID(ctx, p.ID)
	require.NoError(t, err)
	require.Equal(t, p.RevisionID, got.RevisionID)
	require.Equal(t, p.FormatMode, got.FormatMode)
	require.Equal(t, p.FormatOverride, got.FormatOverride)
	require.Equal(t, p.FormatVersion, got.FormatVersion)
	require.Equal(t, p.DataFormat, got.DataFormat)
	require.Equal(t, p.Name, got.Name)
	require.Equal(t, "new-telemetry", got.LastError)
}

func TestRuleFormatRollbackPreservesConcurrentOrdinaryEdit(t *testing.T) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	t.Cleanup(func() { _ = db.Close() })
	ps := point.NewService(point.NewSQLRepository(db), nil)
	before, err := ps.Create(ctx, point.CreatePointRequest{DeviceID: "d", Name: "before", Address: "40001", DataType: schema.DataTypeInt32})
	require.NoError(t, err)
	applied, err := ps.Update(ctx, before.ID, point.UpdatePointRequest{Name: common.Ptr("rule-name")})
	require.NoError(t, err)
	accepted, err := ps.Update(ctx, before.ID, point.UpdatePointRequest{Description: common.Ptr("operator-accepted")})
	require.NoError(t, err)
	syncer := &stubRuntimeSync{}
	svc := NewService(NewSQLRepository(db), nil, ps, syncer)
	err = svc.rollbackUpdatedPoints(ctx, []pointUpdatePlan{{point: before, appliedRevision: applied.RevisionID}})
	require.ErrorIs(t, err, point.ErrFormatConflict)
	got, err := ps.GetByID(ctx, before.ID)
	require.NoError(t, err)
	require.Equal(t, accepted.RevisionID, got.RevisionID)
	require.Equal(t, "operator-accepted", got.Description)
	require.Empty(t, syncer.upserted)
}

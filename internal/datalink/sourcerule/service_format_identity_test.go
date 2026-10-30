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

func TestRuleFormatIdentitySQLLifecycle(t *testing.T) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	t.Cleanup(func() { _ = db.Close() })
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "format-device"))
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	newService := func() *Service { return NewService(NewSQLRepository(db), device.NewService(dr, nil), ps, nil) }
	svc := newService()
	rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "format-device", StartAddress: "40001", Count: 3, DataType: schema.DataTypeInt32, NamingPrefix: "SRC", DataFormat: "ABCD"})
	require.NoError(t, err)
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 3)
	original := map[string]*schema.Point{}
	for _, l := range links {
		p, e := ps.GetByID(ctx, l.PointID)
		require.NoError(t, e)
		p.FormatMode = schema.PointFormatOverride
		p.FormatOverride = "ABCD"
		p.DataFormat = "ABCD"
		require.NoError(t, ps.Restore(ctx, p))
		original[l.Address] = p
	}
	assertRetained := func() {
		t.Helper()
		for _, before := range original {
			got, e := ps.GetByID(ctx, before.ID)
			require.NoError(t, e)
			require.Equal(t, before.FormatMode, got.FormatMode)
			require.Equal(t, before.FormatOverride, got.FormatOverride)
			require.Equal(t, before.FormatVersion, got.FormatVersion)
			require.Equal(t, before.DataFormat, got.DataFormat)
		}
	}
	for _, req := range []UpdateRuleRequest{
		{NamingPrefix: common.Ptr("RENAMED")},
		{DataFormat: common.Ptr("CDAB")},
		{Count: common.Ptr(4)},
		{Enabled: common.Ptr(false)},
	} {
		_, err = svc.Update(ctx, rule.ID, req)
		require.NoError(t, err)
		assertRetained()
	}
	require.NoError(t, dr.UpdateTestResult(ctx, "format-device", true, ""))
	require.NoError(t, svc.Enable(ctx, rule.ID))
	assertRetained()
	require.NoError(t, svc.Disable(ctx, rule.ID))
	assertRetained()
	svc = newService()
	require.NoError(t, svc.SyncRuleDerivedState(ctx, rule.ID))
	assertRetained()
	// Newly added identities inherit the revised rule default, without gaining
	// configured intent even though their compatibility DataFormat is populated.
	currentLinks, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	for _, l := range currentLinks {
		if _, existed := original[l.Address]; existed {
			continue
		}
		p, e := ps.GetByID(ctx, l.PointID)
		require.NoError(t, e)
		require.Equal(t, schema.PointFormatInherit, p.FormatMode)
		require.Empty(t, p.FormatOverride)
		resolved, e := decoding.ResolvePoint(*p, decoding.FormatContext{Protocol: schema.ProtocolModbusTCP, RuleOwned: true, RuleFormat: "CDAB"})
		require.NoError(t, e)
		require.Equal(t, "CDAB", resolved.EffectiveFormat)
		require.Equal(t, decoding.SourceRule, resolved.Source)
	}

	// A shifted range retains actual overlapping point IDs, not row-index intent.
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{StartAddress: common.Ptr("40003")})
	require.NoError(t, err)
	links, err = svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	for _, l := range links {
		p, e := ps.GetByID(ctx, l.PointID)
		require.NoError(t, e)
		if before, ok := original[l.Address]; ok {
			require.Equal(t, before.ID, p.ID)
			require.Equal(t, schema.PointFormatOverride, p.FormatMode)
		} else {
			require.Equal(t, schema.PointFormatInherit, p.FormatMode)
			require.Empty(t, p.FormatOverride)
		}
	}
	// Skip deletes identity. Unskip must not resurrect its override by address.
	victim := links[0]
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{SkippedAddresses: common.Ptr([]string{victim.Address})})
	require.NoError(t, err)
	_, err = ps.GetByID(ctx, victim.PointID)
	require.ErrorIs(t, err, point.ErrPointNotFound)
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{SkippedAddresses: common.Ptr([]string{})})
	require.NoError(t, err)
	links, err = svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	for _, l := range links {
		if l.Address == victim.Address {
			require.NotEqual(t, victim.PointID, l.PointID)
			p, e := ps.GetByID(ctx, l.PointID)
			require.NoError(t, e)
			require.Equal(t, schema.PointFormatInherit, p.FormatMode)
			require.Empty(t, p.FormatOverride)
		}
	}
	// A dangling rule link must not donate intent when its point was manually deleted.
	victim = links[0]
	deleted, err := ps.GetByID(ctx, victim.PointID)
	require.NoError(t, err)
	deleted.FormatMode = schema.PointFormatOverride
	deleted.FormatOverride = "DCBA"
	deleted.DataFormat = "DCBA"
	require.NoError(t, ps.Restore(ctx, deleted))
	require.NoError(t, ps.Delete(ctx, deleted.ID))
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("REBUILT")})
	require.NoError(t, err)
	links, err = svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	found := false
	for _, l := range links {
		if l.Address != victim.Address {
			continue
		}
		found = true
		require.NotEqual(t, victim.PointID, l.PointID)
		p, e := ps.GetByID(ctx, l.PointID)
		require.NoError(t, e)
		require.Equal(t, schema.PointFormatInherit, p.FormatMode)
		require.Empty(t, p.FormatOverride)
	}
	require.True(t, found)

}

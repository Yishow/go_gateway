package sourcerule

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestRuleFormatIncompatibleTypeRejectedBeforeWrites(t *testing.T) {
	for _, dtype := range []schema.DataType{schema.DataTypeBool, schema.DataTypeInt16, schema.DataTypeString} {
		t.Run(string(dtype), func(t *testing.T) {
			ctx := t.Context()
			db := setupSQLRepoDB(t)
			db.SetMaxOpenConns(1) // Keep SQLite total_changes on one connection.
			t.Cleanup(func() { _ = db.Close() })
			dr := device.NewSQLRepository(db)
			require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "d"))
			ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
			svc := NewService(NewSQLRepository(db), device.NewService(dr, nil), ps, nil)
			rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 2, DataType: schema.DataTypeInt32, NamingPrefix: "SRC"})
			require.NoError(t, err)
			links, err := svc.ListLinks(ctx, rule.ID)
			require.NoError(t, err)
			p, err := ps.GetByID(ctx, links[1].PointID)
			require.NoError(t, err)
			p.FormatMode = schema.PointFormatOverride
			p.FormatOverride = "CDAB"
			p.DataFormat = "CDAB"
			require.NoError(t, ps.Restore(ctx, p))
			// Restore updates persisted timestamps through a clone; compare the
			// actual pre-operation row, including its timestamp, not the input.
			p, err = ps.GetByID(ctx, p.ID)
			require.NoError(t, err)
			rule, err = svc.GetByID(ctx, rule.ID)
			require.NoError(t, err)
			candidates, err := svc.ListCandidateSnapshots(ctx, rule.ID)
			require.NoError(t, err)
			// This connection-local write spy includes rolled-back changes and
			// trigger writes, so compensating writes cannot hide behind equality.
			var writesBefore, writesAfter int64
			require.NoError(t, db.QueryRowContext(ctx, `SELECT total_changes()`).Scan(&writesBefore))
			_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{DataType: &dtype, NamingPrefix: common.Ptr("CHANGED")})
			require.ErrorIs(t, err, ErrValidation)
			require.Contains(t, err.Error(), "reset")
			require.NoError(t, db.QueryRowContext(ctx, `SELECT total_changes()`).Scan(&writesAfter))
			require.Equal(t, writesBefore, writesAfter, "preflight rejection must perform no SQL writes")
			got, err := ps.GetByID(ctx, p.ID)
			require.NoError(t, err)
			require.Equal(t, p, got)
			gotRule, err := svc.GetByID(ctx, rule.ID)
			require.NoError(t, err)
			require.Equal(t, rule, gotRule)
			gotLinks, err := svc.ListLinks(ctx, rule.ID)
			require.NoError(t, err)
			require.Equal(t, links, gotLinks)
			gotCandidates, err := svc.ListCandidateSnapshots(ctx, rule.ID)
			require.NoError(t, err)
			require.Equal(t, candidates, gotCandidates)
		})
	}
}

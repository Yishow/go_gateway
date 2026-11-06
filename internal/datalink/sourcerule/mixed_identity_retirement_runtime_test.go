package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestMixedRetirementRestoresOwnedRuntimeAfterFailure(t *testing.T) {
	for _, failure := range []string{"links", "runtime"} {
		t.Run(failure, func(t *testing.T) {
			db := setupSQLRepoDB(t)
			t.Cleanup(func() { _ = db.Close() })
			db.SetMaxOpenConns(1)
			_, err := db.ExecContext(t.Context(), `PRAGMA foreign_keys=ON`)
			require.NoError(t, err)
			dr := device.NewSQLRepository(db)
			success := true
			require.NoError(t, dr.Create(t.Context(), &schema.Device{ID: "d", Name: "MC", Protocol: schema.ProtocolMC3E, Status: schema.DeviceStatusActive, LastTestSuccess: &success, ConnectionConfig: `{"host":"127.0.0.1","port":5000,"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1"}}`, CreatedAt: time.Now(), UpdatedAt: time.Now()}))
			ps := point.NewService(point.NewSQLRepository(db), nil)
			runtime := &stubRuntimeSync{}
			svc := NewService(NewSQLRepository(db), device.NewService(dr, nil), ps, runtime)
			rule, err := svc.Create(t.Context(), CreateRuleRequest{ID: "replace", DeviceID: "d", NamingPrefix: "M", Enabled: true, LayoutMode: "mixed", MixedItems: []measurement.MixedItem{mixedItem("a", "D100", schema.DataTypeUint16), mixedItem("b", "D200", schema.DataTypeUint16)}})
			require.NoError(t, err)
			links, err := svc.ListLinks(t.Context(), rule.ID)
			require.NoError(t, err)
			require.Len(t, links, 2)
			byAddress := map[string]*schema.SourceRuleLink{}
			for _, link := range links {
				byAddress[link.Address] = link
				require.NotZero(t, runtime.leases[link.PointID])
			}
			if failure == "links" {
				query := fmt.Sprintf(`CREATE TRIGGER reject_move BEFORE INSERT ON source_rule_links WHEN NEW.address='D200' AND NEW.point_id='%s' BEGIN SELECT RAISE(ABORT, 'injected link failure'); END`, byAddress["D100"].PointID)
				_, err = db.ExecContext(t.Context(), query)
				require.NoError(t, err)
			} else {
				svc.SetShareRuntimeReconciler(ShareRuntimeReconcilerFunc(func(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
					return RuntimeReconcileOutcome{}, errors.New("injected runtime failure")
				}))
			}
			items := []measurement.MixedItem{mixedItem("a", "D200", schema.DataTypeUint16)}
			_, _, err = svc.UpdateWithRuntimeReconcile(t.Context(), rule.ID, UpdateRuleRequest{MixedItems: &items})
			require.Error(t, err)
			var state *modbusshare.Error
			if errors.As(err, &state) {
				require.NotEqual(t, modbusshare.ErrCodeDirtyUnknown, state.Code, err)
			}
			after, err := svc.ListLinks(t.Context(), rule.ID)
			require.NoError(t, err)
			require.Len(t, after, 2)
			require.Len(t, runtime.leases, 2)
			for _, link := range after {
				original := byAddress[link.Address]
				require.NotNil(t, original)
				require.Equal(t, original.PointID, link.PointID)
				p, err := ps.GetByID(t.Context(), link.PointID)
				require.NoError(t, err)
				require.Equal(t, link.Address, p.Address)
				require.NotZero(t, runtime.leases[p.ID])
				require.Equal(t, p.RevisionID, runtime.revisions[p.ID])
			}
		})
	}
}

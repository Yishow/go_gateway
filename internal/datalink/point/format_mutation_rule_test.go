package point

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

type formatRuleTestGuard struct{ calls int }

func (g *formatRuleTestGuard) AcquireRuleMutation(ctx context.Context, _ string) (context.Context, func()) {
	g.calls++
	return ctx, func() {}
}
func TestFormatMutationRuleIdentityAndRevision(t *testing.T) {
	for _, kind := range []string{"valid", "wrong-rule", "wrong-device", "wrong-address", "missing-rule", "stale-rule", "missing-guard", "mixed-address"} {
		t.Run(kind, func(t *testing.T) {
			db, svc, req := formatFixture(t)
			_, err := db.Exec(`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,data_format) VALUES('rule',?,'40001',1,'uint32','p','rule-old','ABCD')`, req.DeviceID)
			require.NoError(t, err)
			_, err = db.Exec(`INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('link','rule','40001',?)`, req.PointID)
			require.NoError(t, err)
			req.RuleID = "rule"
			req.ExpectedRuleRevision = "rule-old"
			guard := &formatRuleTestGuard{}
			svc.SetFormatRuleGuard(guard)
			want := ErrFormatScope
			switch kind {
			case "valid":
				want = nil
			case "wrong-rule":
				req.RuleID = "other"
			case "wrong-device":
				_, err = db.Exec(`UPDATE source_rules SET device_id='other'`)
				require.NoError(t, err)
			case "wrong-address":
				_, err = db.Exec(`UPDATE source_rule_links SET address='40002'`)
				require.NoError(t, err)
			case "missing-rule":
				req.RuleID = ""
				req.ExpectedRuleRevision = ""
			case "stale-rule":
				req.ExpectedRuleRevision = "stale"
				want = ErrFormatConflict
			case "mixed-address":
				address := "40002"
				req.Patch = &UpdatePointRequest{Address: &address}
				want = ErrFormatInvalid
			case "missing-guard":
				svc.SetFormatRuleGuard(nil)
				want = ErrFormatUnavailable
			}
			got, err := svc.MutateFormat(t.Context(), req)
			if want == nil {
				require.NoError(t, err)
				require.Equal(t, "rule-old", got.Resolution.RuleRevision)
				require.Equal(t, 1, guard.calls)
			} else {
				require.ErrorIs(t, err, want)
				require.Nil(t, got)
			}
		})
	}
}

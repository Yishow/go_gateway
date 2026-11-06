package point

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestMixedRuleItemFormatIsUsedWhenResettingPointIntent(t *testing.T) {
	db, svc, req := formatFixture(t)
	_, err := db.ExecContext(t.Context(), `INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,layout_mode,mixed_items) VALUES('rule',?,'40001',1,'uint32','p','rule-old','mixed',?)`, req.DeviceID, `[{"item_id":"value","address":"40001","data_type":"uint32","data_format":"DCBA"}]`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('link','rule','40001',?)`, req.PointID)
	require.NoError(t, err)
	req.RuleID = "rule"
	req.ExpectedRuleRevision = "rule-old"
	req.FormatMode = schema.PointFormatInherit
	req.ConfiguredFormat = ""
	svc.SetFormatRuleGuard(&formatRuleTestGuard{})
	result, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "DCBA", result.Resolution.EffectiveFormat)
}

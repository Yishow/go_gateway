package workspace

import (
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

const mixedGraphItems = `[{"item_id":"d39","name":"","address":"D39","register_offset":0,"data_type":"int16","semantic_kind":"gauge"},{"item_id":"d100","name":"","address":"D100","register_offset":0,"data_type":"int32","semantic_kind":"gauge"}]`

func newMixedAcquisitionFixture(t *testing.T, profile string) *Service {
	t.Helper()
	db := openWorkspaceTestDB(t, filepath.Join(t.TempDir(), "acquisition.db"))
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	config := `{"host":"127.0.0.1","port":5000,"network_no":0,"pc_no":255,"acquisition":{"version":1,"profile":"` + profile + `"}}`
	for _, query := range []string{
		`INSERT INTO devices(id,name,description,protocol,status,connection_config) VALUES('d','d','','mc_3e','active','` + config + `')`,
		`INSERT INTO points(id,device_id,name,address,data_type,revision_id) VALUES('h0','d','h0','D0','int16','p0'),('h1','d','h1','D1','int16','p1'),('m39','d','m39','D39','int16','p2'),('m100','d','m100','D100','int32','p3')`,
		`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,skipped_addresses) VALUES('range','d','D0',2,'int16','h','rev-range','[]')`,
		`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,skipped_addresses,layout_mode,mixed_items) VALUES('mixed','d','D39',2,'int16','m','rev-mixed','[]','mixed','` + mixedGraphItems + `')`,
		`INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('l0','range','D0','h0'),('l1','range','D1','h1'),('l39','mixed','D39','m39'),('l100','mixed','D100','m100')`,
		`INSERT INTO system_settings(key,value) VALUES('studio_v2_workspace','{"id":"w","ordered_device_ids":["d"]}')`,
	} {
		_, err := db.ExecContext(t.Context(), query)
		require.NoError(t, err, query)
	}
	return NewService(NewSQLRepository(db)).WithRuntimeProjectionServices(
		device.NewSQLRepository(db), sourcerule.NewService(sourcerule.NewSQLRepository(db), nil, nil, nil),
		point.NewSQLRepository(db), mapping.NewSQLRepository(db), tag.NewSQLRepository(db), nil, nil,
		point.NewMemoryPollingGroupRepository(),
	)
}

func TestRuntimeAcquisitionSourceGraphAcceptsMixedRuleOnRandomProfile(t *testing.T) {
	service := newMixedAcquisitionFixture(t, "mitsubishi-q-iqr-3e-binary-v1")
	projection, err := service.RuntimeProjection(t.Context())
	require.NoError(t, err)
	graph, ok := projection.AcquisitionSourceGraph("d")
	require.True(t, ok, graph.Reason())
	require.Empty(t, graph.Excluded())
	revision, ok := graph.RuleSemanticRevision("mixed")
	require.True(t, ok)
	require.Len(t, revision, 64)
	require.Len(t, graph.Points(), 4)
	require.NoError(t, service.AcceptAcquisitionSourceGraph(t.Context(), graph, func() error { return nil }))
}

func TestRuntimeAcquisitionSourceGraphExcludesMixedRuleOnSafeProfile(t *testing.T) {
	service := newMixedAcquisitionFixture(t, "safe-v1")
	projection, err := service.RuntimeProjection(t.Context())
	require.NoError(t, err)
	graph, ok := projection.AcquisitionSourceGraph("d")
	require.False(t, ok)
	require.False(t, graph.Complete())
	require.Equal(t, "source-layout-unsupported", graph.Reason())
}

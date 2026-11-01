package workspace

import (
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

type acquisitionSQLFixture struct {
	db      *sql.DB
	path    string
	service *Service
	repo    *SQLRepository
}

func newAcquisitionSQLFixture(t *testing.T) acquisitionSQLFixture {
	t.Helper()
	path := filepath.Join(t.TempDir(), "acquisition.db")
	db := openWorkspaceTestDB(t, path)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	for _, query := range []string{
		`INSERT INTO devices(id,name,description,protocol,status) VALUES('d','d','','modbus_tcp','active')`,
		`INSERT INTO points(id,device_id,name,address,data_type,revision_id) VALUES('a','d','a','40001','uint16','pa'),('b','d','b','40003','uint16','pb')`,
		`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,skipped_addresses) VALUES('r','d','40001',3,'uint16','p','rule-original','["40002"]')`,
		`INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('la','r','40001','a'),('lb','r','40003','b')`,
		`INSERT INTO system_settings(key,value) VALUES('studio_v2_workspace','{"id":"w","ordered_device_ids":["d"]}')`,
	} {
		_, err := db.ExecContext(t.Context(), query)
		require.NoError(t, err)
	}
	repo := NewSQLRepository(db)
	service := NewService(repo).WithRuntimeProjectionServices(
		device.NewSQLRepository(db), sourcerule.NewService(sourcerule.NewSQLRepository(db), nil, nil, nil),
		point.NewSQLRepository(db), mapping.NewSQLRepository(db), tag.NewSQLRepository(db), nil, nil,
		point.NewMemoryPollingGroupRepository(),
	)
	return acquisitionSQLFixture{db: db, path: path, service: service, repo: repo}
}

func (f acquisitionSQLFixture) graph(t *testing.T) (*RuntimeProjection, AcquisitionSourceGraph) {
	t.Helper()
	projection, err := f.service.RuntimeProjection(t.Context())
	require.NoError(t, err)
	graph, ok := projection.AcquisitionSourceGraph("d")
	require.True(t, ok, graph.Reason())
	return projection, graph
}

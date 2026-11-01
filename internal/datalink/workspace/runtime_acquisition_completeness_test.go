package workspace

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type acquisitionTruncatedPoints struct {
	repo      *point.SQLRepository
	seenLimit int
}

func (r *acquisitionTruncatedPoints) List(ctx context.Context, filter point.ListFilter) ([]*schema.Point, error) {
	r.seenLimit = filter.Limit
	rows, err := r.repo.List(ctx, filter)
	if len(rows) > 1 {
		rows = rows[:1]
	}
	return rows, err
}

func TestRuntimeAcquisitionSourceGraphCannotCertifyTruncatedPointReader(t *testing.T) {
	f := newAcquisitionSQLFixture(t)
	reader := &acquisitionTruncatedPoints{repo: point.NewSQLRepository(f.db)}
	f.service.projectionPoints = reader
	projection, err := f.service.RuntimeProjection(t.Context())
	require.NoError(t, err)
	require.Equal(t, 100000, reader.seenLimit, "exercise original bounded reader, not a replacement producer")
	require.Len(t, projection.Points, 1)
	graph, ok := projection.AcquisitionSourceGraph("d")
	require.False(t, ok, "same-tx unbounded enumeration must detect omitted linked rows")
	require.False(t, graph.Complete())
}

func TestRuntimeAcquisitionRuleSemanticRevisionIsForeignSafe(t *testing.T) {
	f := newAcquisitionSQLFixture(t)
	for _, query := range []string{
		`INSERT INTO points(id,device_id,name,address,data_type,revision_id) VALUES('c','d','c','40005','uint16','pc')`,
		`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id) VALUES('foreign','d','40005',1,'uint16','p','foreign-before')`,
		`INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('lc','foreign','40005','c')`,
	} {
		_, err := f.db.ExecContext(t.Context(), query)
		require.NoError(t, err)
	}
	_, before := f.graph(t)
	ownedBefore, ok := before.RuleSemanticRevision("r")
	require.True(t, ok)
	_, err := f.db.ExecContext(t.Context(), `UPDATE source_rules SET enabled=0,revision_id='foreign-after' WHERE id='foreign'; UPDATE points SET enabled=0 WHERE id='c'`)
	require.NoError(t, err)
	_, after := f.graph(t)
	ownedAfter, ok := after.RuleSemanticRevision("r")
	require.True(t, ok)
	require.Equal(t, ownedBefore, ownedAfter, "foreign disable must not demand a new source lease for unchanged r")
	require.NotEqual(t, before.SemanticRevision(), after.SemanticRevision(), "device-wide gap geometry must change")
	require.ErrorIs(t, f.service.AcceptAcquisitionSourceGraph(t.Context(), before, func() error { t.Fatal("foreign source drift accepted"); return nil }), ErrAcquisitionSourceGraph)
}

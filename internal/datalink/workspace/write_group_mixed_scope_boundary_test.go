package workspace

import (
	"context"
	"go-gateway/internal/datalink/dbtarget"
	"testing"

	"github.com/stretchr/testify/require"
)

type mixedForbiddenProbe struct{ calls int }

func (p *mixedForbiddenProbe) InspectTableAtRevision(context.Context, string, string, string, string) (*dbtarget.TableInspection, error) {
	p.calls++
	return &dbtarget.TableInspection{Status: "missing"}, nil
}
func (p *mixedForbiddenProbe) ProbeMySQLRecordingAtRevision(_ context.Context, id, rev string) (*dbtarget.MySQLRecordingCapability, error) {
	p.calls++
	return &dbtarget.MySQLRecordingCapability{ConnectorID: id, ConnectorRevision: rev, Database: "recordings", ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-mysql", Connected: true, MetadataStatus: "verified", TransactionalEngine: true}, nil
}
func (p *mixedForbiddenProbe) ProbeSQLServerRecordingAtRevision(_ context.Context, id, rev string) (*dbtarget.SQLServerRecordingCapability, error) {
	p.calls++
	return sqlServerTestCapability(id, rev), nil
}

func TestMixedPublicScopeRejectsUnsupportedDialectBeforeProbe(t *testing.T) {
	for _, kind := range []string{"mysql", "sqlserver"} {
		t.Run(kind, func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			var id string
			if kind == "mysql" {
				id = seedMySQLConnector(t, h.db, `{"dsn":"app:secret@tcp(127.0.0.1:3306)/recordings?tls=true"}`)
			} else {
				id = seedSQLServerConnector(t, h.db, `{"host":"sql.example","database":"recordings","schema":"dbo","username":"u","password":"p"}`)
			}
			g := mixedDraft(h, WriteGroupMemberCalculationMean, 15, "")
			g.LayoutVersion = WriteGroupLayoutManagedLongV1
			g.StorageMode = WriteGroupStorageModePlain
			g.Destination.ConnectorID = id
			g.Destination.ConnectorRevision = "connector-revision-" + kind
			g.Destination.Database = ""
			g.Destination.TableSchema = ""
			req := h.mutation(t.Context(), g)
			req.ExpectedConnectorRevision = g.Destination.ConnectorRevision
			saved, err := h.service.Create(t.Context(), req)
			require.NoError(t, err)
			p := &mixedForbiddenProbe{}
			h.service.WithManagedTableInspector(p)
			req = h.mutation(t.Context(), saved.Group)
			req.ExpectedGroupRevision = saved.Group.Revision
			req.ExpectedConnectorRevision = saved.Group.Destination.ConnectorRevision
			_, err = h.service.ManagedSchemaScope(t.Context(), saved.Group.ID, req)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			require.Zero(t, p.calls, "unsupported mixed dialect must not probe or inspect the target")
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM managed_schema_operations`))
		})
	}
}

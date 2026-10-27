package grouptestwrite

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type sqlServerTokenDestinations struct {
	sqlServerTestWriteDestinations
	revision string
	opens    int
}

func (d *sqlServerTokenDestinations) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	connector, err := d.sqlServerTestWriteDestinations.GetByID(ctx, id)
	if connector != nil && d.revision != "" {
		connector.IdentityRevision = d.revision
	}
	return connector, err
}

func (d *sqlServerTokenDestinations) OpenDestination(ctx context.Context, id, revision string) (*dbtarget.OpenedDestination, error) {
	d.opens++
	return d.sqlServerTestWriteDestinations.OpenDestination(ctx, id, revision)
}

// These are regressions of the actual shared token ledger with a SQL Server
// plan. The HTTP bridge's error-mapping fixture does not itself prove staleness.
func TestSQLServerTokenLedgerRejectsChangedScopeBeforeDestinationOpen(t *testing.T) {
	for _, change := range []string{"group", "workspace", "connector", "operation", "foreign workspace"} {
		t.Run(change, func(t *testing.T) {
			svc, p := sqlServerExecutionService(t)
			destination := &sqlServerTokenDestinations{}
			svc.deps.Destinations = destination
			ws := Workspace{ID: p.group.WorkspaceID, Revision: testRevision}
			preview, err := svc.Preview(t.Context(), ws, p.group.ID)
			require.NoError(t, err)
			req := Confirmation{Token: preview.Token, OperationID: preview.OperationID, GroupID: p.group.ID}
			var expected error
			switch change {
			case "group":
				p.group.Revision = "changed-group"
				expected = recordingplan.ErrPreviewTokenStale
			case "workspace":
				ws.Revision = "changed-workspace"
				expected = recordingplan.ErrPreviewTokenStale
			case "connector":
				destination.revision = "changed-connector"
				expected = ErrDestinationChanged
			case "operation":
				req.OperationID = "other-operation"
				expected = recordingplan.ErrSchemaOperationMismatch
			case "foreign workspace":
				ws.ID = "other-workspace"
				expected = recordingplan.ErrPreviewTokenNotFound
			}
			_, err = svc.Confirm(t.Context(), ws, req)
			require.ErrorIs(t, err, expected)
			require.Zero(t, destination.opens, "a stale or foreign token must not reach SQL Server")
			_, err = svc.deps.Ledger.GetSchemaOperation(t.Context(), p.group.WorkspaceID, preview.OperationID)
			require.ErrorIs(t, err, recordingplan.ErrSchemaOperationNotFound, "refusal must precede claiming any write")
		})
	}
}

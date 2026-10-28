package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"
)

// Delegates to the actual SQL Server inspector unless deliberately offline.
// A spy is observability/fault injection, never manufactured real-server proof.
type sqlServerSuccessorProbeSpy struct {
	*dbtarget.ManagedTableInspector
	calls      int
	offline    bool
	afterProbe func()
}

func (s *sqlServerSuccessorProbeSpy) ProbeSQLServerRecordingAtRevision(ctx context.Context, id, revision string) (*dbtarget.SQLServerRecordingCapability, error) {
	s.calls++
	if s.offline {
		return nil, errors.New("fixture target offline")
	}
	cap, err := s.ManagedTableInspector.ProbeSQLServerRecordingAtRevision(ctx, id, revision)
	if s.afterProbe != nil {
		s.afterProbe()
	}
	return cap, err
}
func (s *sqlServerSuccessorProbeSpy) InspectTableAtRevision(ctx context.Context, id, revision, namespace, table string) (*dbtarget.TableInspection, error) {
	s.calls++
	if s.offline {
		return nil, errors.New("fixture target offline")
	}
	return s.ManagedTableInspector.InspectTableAtRevision(ctx, id, revision, namespace, table)
}
func (s *sqlServerSuccessorProbeSpy) InspectSQLServerManagedTableAtRevision(ctx context.Context, id, revision, database, namespace, capability, table string) (*dbtarget.SQLServerManagedInspection, error) {
	s.calls++
	if s.offline {
		return nil, errors.New("fixture target offline")
	}
	return s.ManagedTableInspector.InspectSQLServerManagedTableAtRevision(ctx, id, revision, database, namespace, capability, table)
}
func (s *sqlServerSuccessorProbeSpy) ExecuteSQLServerManagedStatementsAtRevision(ctx context.Context, id, revision, database, namespace, capability string, native bool, statements []string) (*dbtarget.SchemaStatementExecution, error) {
	s.calls++
	if s.offline {
		return nil, errors.New("fixture target offline")
	}
	return s.ManagedTableInspector.ExecuteSQLServerManagedStatementsAtRevision(ctx, id, revision, database, namespace, capability, native, statements)
}

func TestSQLServerSuccessorReadinessRechecksReservedCASAfterActualProbe(t *testing.T) {
	sqlServerSuccessorFixtures(t) // Missing/unsafe real fixtures must mark the entrypoint NOT_RUN.
	for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
		t.Run(string(mode), func(t *testing.T) {
			r := newSQLServerSuccessorRig(t, mode, false, true)
			f := r.activationAPI()
			review := f.postActivation(t, f.activationRequest(t))
			var after map[string][]string
			spy := &sqlServerSuccessorProbeSpy{ManagedTableInspector: r.f.groups.ManagedTableInspector(r.targets)}
			spy.afterProbe = func() {
				if after != nil {
					return
				}
				_, err := r.f.db.Exec(`UPDATE write_groups SET revision='changed-during-probe' WHERE id=?`, r.prepared.Group.ID)
				require.NoError(t, err)
				after = sqlServerSuccessorLocalState(t, r.f.db)
			}
			r.f.groups.WithManagedTableInspector(spy)
			external := r.targetState(t)
			response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", map[string]any{"action": "activation-readiness", "envelope": review["envelope"], "review_digest": review["review_digest"]})
			require.Positive(t, spy.calls)
			require.NotNil(t, after)
			if response.Code == http.StatusOK {
				ready := decodeJSONBody(t, response)["data"].(map[string]any)
				require.Equal(t, false, ready["readiness"].(map[string]any)["ready"])
				require.Empty(t, ready["readiness_digest"])
			} else {
				require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
			}
			require.Equal(t, after, sqlServerSuccessorLocalState(t, r.f.db))
			require.Equal(t, external, r.targetState(t), "a stale scoped reload cannot mutate SQL")
		})
	}
}

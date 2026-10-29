package api

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"testing"
)

type mysqlSuccessorProbeSpy struct {
	*dbtarget.ManagedTableInspector
	calls      int
	offline    bool
	afterProbe func()
}

func (s *mysqlSuccessorProbeSpy) ProbeMySQLRecordingAtRevision(ctx context.Context, id, revision string) (*dbtarget.MySQLRecordingCapability, error) {
	s.calls++
	if s.offline {
		return nil, errors.New("fixture target offline")
	}
	cap, err := s.ManagedTableInspector.ProbeMySQLRecordingAtRevision(ctx, id, revision)
	if s.afterProbe != nil {
		s.afterProbe()
	}
	return cap, err
}

func (s *mysqlSuccessorProbeSpy) InspectTableAtRevision(ctx context.Context, id, revision, namespace, table string) (*dbtarget.TableInspection, error) {
	s.calls++
	if s.offline {
		return nil, errors.New("fixture target offline")
	}
	return s.ManagedTableInspector.InspectTableAtRevision(ctx, id, revision, namespace, table)
}
func (s *mysqlSuccessorProbeSpy) InspectMySQLManagedTableAtRevision(ctx context.Context, id, revision, database, capability, table string) (*dbtarget.MySQLManagedInspection, error) {
	s.calls++
	if s.offline {
		return nil, errors.New("fixture target offline")
	}
	return s.ManagedTableInspector.InspectMySQLManagedTableAtRevision(ctx, id, revision, database, capability, table)
}
func (s *mysqlSuccessorProbeSpy) ExecuteMySQLManagedStatementAtRevision(ctx context.Context, id, revision, database, capability, statement string) error {
	s.calls++
	if s.offline {
		return errors.New("fixture target offline")
	}
	return s.ManagedTableInspector.ExecuteMySQLManagedStatementAtRevision(ctx, id, revision, database, capability, statement)
}

func TestMySQLSuccessorLocalReviewsAndPendingReplayDoNotProbe(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
				t.Run(string(mode), func(t *testing.T) {
					r := newMySQLSuccessorRig(t, family, family, mode, false, false)
					spy := &mysqlSuccessorProbeSpy{ManagedTableInspector: r.f.groups.ManagedTableInspector(r.targets), offline: true}
					r.f.groups.WithManagedTableInspector(spy)
					before, external := successorAPIAllState(t, r.f.db), r.targetState(t)
					review := r.review(t)
					require.Len(t, review.RequiredChecks, 1)
					require.Equal(t, "target-family-check-required", review.RequiredChecks[0].ReasonCode)
					require.Zero(t, spy.calls, "local review must not open target even when offline")
					require.Equal(t, before, successorAPIAllState(t, r.f.db))
					require.Equal(t, external, r.targetState(t))
					spy.offline = false
					prepared := r.reserve(t, review)
					require.Positive(t, spy.calls)
					calls := spy.calls
					spy.offline = true
					before = successorAPIAllState(t, r.f.db)
					replay := r.reserve(t, review)
					require.Equal(t, prepared.Group.ID, replay.Group.ID)
					require.Equal(t, calls, spy.calls, "exact pending replay must not probe offline target")
					require.Equal(t, before, successorAPIAllState(t, r.f.db))
					require.Equal(t, external, r.targetState(t))
				})
			}
		})
	}
}

func TestMySQLSuccessorProbeFailureAndStaleConnectorCannotReserve(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			for _, fault := range []string{"offline", "revision", "database"} {
				t.Run(fault, func(t *testing.T) {
					r := newMySQLSuccessorRig(t, family, family, "plain", false, false)
					review := r.review(t)
					spy := &mysqlSuccessorProbeSpy{ManagedTableInspector: r.f.groups.ManagedTableInspector(r.targets), offline: fault == "offline"}
					var after map[string][]string
					spy.afterProbe = func() {
						if fault == "revision" {
							_, err := r.f.db.Exec(`UPDATE database_connectors SET identity_revision='changed-during-probe' WHERE id='connector-B'`)
							require.NoError(t, err)
						}
						if fault == "database" {
							_, err := r.f.db.Exec(`UPDATE database_connectors SET connection_config=replace(connection_config,?,?) WHERE id='connector-B'`, r.targetNamespace, r.f.namespace)
							require.NoError(t, err)
						}
						after = successorAPIAllState(t, r.f.db)
					}
					r.f.groups.WithManagedTableInspector(spy)
					before, external := successorAPIAllState(t, r.f.db), r.targetState(t)
					response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", map[string]any{"action": "reserve", "envelope": review.Envelope, "review_digest": review.ReviewDigest})
					require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
					require.Less(t, response.Code, 500, response.Body.String())
					require.Equal(t, 1, spy.calls)
					if after != nil {
						before = after
					}
					require.Equal(t, before, successorAPIAllState(t, r.f.db))
					require.Equal(t, external, r.targetState(t))
				})
			}
		})
	}
}

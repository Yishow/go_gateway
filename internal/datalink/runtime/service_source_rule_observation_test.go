package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type sourceObservationReader struct {
	record             *workspace.Record
	projection         *workspace.RuntimeProjection
	reads, projections int
	beforeRead         func(int)
	failure            error
	panicRead          bool
}

func (r *sourceObservationReader) ReadSnapshot(context.Context) (*workspace.Record, error) {
	r.reads++
	if r.panicRead {
		panic("password_fixture")
	}
	if r.beforeRead != nil {
		r.beforeRead(r.reads)
	}
	if r.failure != nil {
		return nil, r.failure
	}
	if r.record == nil {
		return nil, nil
	}
	v := *r.record
	v.OrderedDeviceIDs = append([]string(nil), r.record.OrderedDeviceIDs...)
	return &v, nil
}
func (r *sourceObservationReader) RuntimeProjection(context.Context) (*workspace.RuntimeProjection, error) {
	r.projections++
	return r.projection, nil
}

func newSourceObservationRuntime(t *testing.T) (*Service, *sourceObservationReader, sourcerule.RuntimeReconcileRequest, *[]sourcerule.RuntimeSourceObservation) {
	t.Helper()
	r := &sourceObservationReader{record: &workspace.Record{ID: "W", DatabaseSetupRevision: "observed-W2", OrderedDeviceIDs: []string{"D"}}, projection: &workspace.RuntimeProjection{WorkspaceID: "W", Version: "input-fingerprint-not-W2", DeviceIDs: []string{"D"}, Devices: []*schema.Device{{ID: "D", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`}}, Rules: []*schema.SourceRule{{ID: "R", DeviceID: "D", RevisionID: "source-R3"}}}}
	seen := []sourcerule.RuntimeSourceObservation{}
	s := &Service{scheduler: collector.NewScheduler(collector.DefaultSchedulerConfig(), nil), workspace: r, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{}, sourceRuleObserver: func(_ context.Context, o sourcerule.RuntimeSourceObservation) { seen = append(seen, o) }}
	s.running.Store(true)
	req := sourcerule.RuntimeReconcileRequest{Operation: sourcerule.RuntimeReconcileOperationDisable, Scope: sourcerule.RuntimeReconcileScope{RuleID: "R", DeviceID: "D"}, RuleRevision: "source-R3"}
	return s, r, req, &seen
}

func TestSourceObservationUsesObservedRevisionAndWaitsForFinalization(t *testing.T) {
	s, r, req, seen := newSourceObservationRuntime(t)
	out := s.ReconcileSourceRule(t.Context(), req)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, out.Status)
	require.Empty(t, *seen)
	require.NotNil(t, out.Observation)
	require.True(t, out.Observation.Verified)
	require.Equal(t, "observed-W2", out.Observation.Scope.WorkspaceRevision)
	require.NotEqual(t, r.projection.Version, out.Observation.Scope.WorkspaceRevision)
	s.FinalizeSourceRuleObservation(t.Context(), out, nil)
	require.Len(t, *seen, 1)
	require.Equal(t, sourcerule.SourceObservationApplied, (*seen)[0].Result)
	require.NotEmpty(t, (*seen)[0].ReferenceID)
	require.Equal(t, 1, r.projections, "logging must not call bootstrap-capable projection again")
	require.Equal(t, 3, r.reads)
	raw, err := json.Marshal(out)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "Observation")
	require.NotContains(t, string(raw), "observed-W2")
	raw, err = json.Marshal(req)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "source-R3")
}

func TestSourceObservationNilObserverDoesNotAddReads(t *testing.T) {
	s, r, req, _ := newSourceObservationRuntime(t)
	s.sourceRuleObserver = nil
	out := s.ReconcileSourceRule(t.Context(), req)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, out.Status)
	require.Nil(t, out.Observation)
	require.Zero(t, r.reads)
	require.Equal(t, 1, r.projections)
}

func TestSourceObservationMissingForeignOrChangedProofIsUnavailable(t *testing.T) {
	cases := map[string]func(*Service, *sourceObservationReader, *sourcerule.RuntimeReconcileRequest){
		"empty observed revision": func(_ *Service, r *sourceObservationReader, _ *sourcerule.RuntimeReconcileRequest) {
			r.record.DatabaseSetupRevision = ""
		},
		"missing snapshot": func(_ *Service, r *sourceObservationReader, _ *sourcerule.RuntimeReconcileRequest) { r.record = nil },
		"foreign workspace": func(_ *Service, r *sourceObservationReader, _ *sourcerule.RuntimeReconcileRequest) {
			r.record.ID = "foreign"
		},
		"foreign device": func(_ *Service, r *sourceObservationReader, _ *sourcerule.RuntimeReconcileRequest) {
			r.record.OrderedDeviceIDs = []string{"other"}
		},
		"changed revision": func(_ *Service, r *sourceObservationReader, _ *sourcerule.RuntimeReconcileRequest) {
			r.beforeRead = func(n int) {
				if n == 2 {
					r.record.DatabaseSetupRevision = "changed"
				}
			}
		},
		"changed membership": func(_ *Service, r *sourceObservationReader, _ *sourcerule.RuntimeReconcileRequest) {
			r.beforeRead = func(n int) {
				if n == 2 {
					r.record.OrderedDeviceIDs = []string{"D", "another"}
				}
			}
		},
		"foreign rule": func(_ *Service, r *sourceObservationReader, _ *sourcerule.RuntimeReconcileRequest) {
			r.projection.Rules[0].DeviceID = "foreign"
		},
		"wrong revision": func(_ *Service, _ *sourceObservationReader, q *sourcerule.RuntimeReconcileRequest) {
			q.RuleRevision = "wrong"
		},
		"missing revision": func(_ *Service, _ *sourceObservationReader, q *sourcerule.RuntimeReconcileRequest) {
			q.RuleRevision = ""
		},
		"deleted workspace": func(_ *Service, r *sourceObservationReader, _ *sourcerule.RuntimeReconcileRequest) {
			r.beforeRead = func(n int) {
				if n == 2 {
					r.record = nil
				}
			}
		},
		"provider error": func(_ *Service, r *sourceObservationReader, _ *sourcerule.RuntimeReconcileRequest) {
			r.failure = errors.New("password_fixture")
		},
		"provider panic": func(_ *Service, r *sourceObservationReader, _ *sourcerule.RuntimeReconcileRequest) {
			r.panicRead = true
		},
		"replacement after delete": func(_ *Service, _ *sourceObservationReader, q *sourcerule.RuntimeReconcileRequest) {
			q.Operation = sourcerule.RuntimeReconcileOperationDelete
		},
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			s, r, req, seen := newSourceObservationRuntime(t)
			change(s, r, &req)
			out := s.ReconcileSourceRule(t.Context(), req)
			require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, out.Status)
			s.FinalizeSourceRuleObservation(t.Context(), out, nil)
			require.Len(t, *seen, 1)
			require.Equal(t, sourcerule.SourceObservationUnavailable, (*seen)[0].Result)
			require.Empty(t, (*seen)[0].RuleID)
			require.Equal(t, 1, r.projections)
		})
	}
}

func TestSourceObservationAlignedNoOpAndDeferredAreNotApplied(t *testing.T) {
	for _, name := range []string{"zero-point create", "inactive", "not running", "missing reader"} {
		t.Run(name, func(t *testing.T) {
			s, r, req, seen := newSourceObservationRuntime(t)
			switch name {
			case "zero-point create":
				req.Operation = sourcerule.RuntimeReconcileOperationCreate
			case "inactive":
				r.projection.Devices[0].Status = schema.DeviceStatusDraft
			case "not running":
				s.running.Store(false)
			case "missing reader":
				s.workspace = workspaceProjectionStub{projection: r.projection}
			}
			out := s.ReconcileSourceRule(t.Context(), req)
			if name == "zero-point create" {
				require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, out.Status)
			}
			s.FinalizeSourceRuleObservation(t.Context(), out, nil)
			require.Len(t, *seen, 1)
			require.Equal(t, sourcerule.SourceObservationUnavailable, (*seen)[0].Result)
		})
	}
}

func TestSourceObservationRejectsAttemptWithoutBorrowingInstalledSuccess(t *testing.T) {
	s, r, req, seen := newSourceObservationRuntime(t)
	r.projection.Points = []*schema.Point{{ID: "P", DeviceID: "D"}}
	r.projection.Mappings = []*schema.Mapping{{ID: "M", PointID: "P", TagID: "missing", Enabled: true, TransformPipeline: "[]"}}
	out := s.ReconcileSourceRule(t.Context(), req)
	require.Equal(t, sourcerule.RuntimeReconcileStatusStale, out.Status)
	require.True(t, out.Observation.Verified)
	require.False(t, out.Observation.Applied)
	s.FinalizeSourceRuleObservation(t.Context(), out, &modbusshare.Error{Code: modbusshare.ErrCodeReconcileFailed})
	require.Len(t, *seen, 1)
	require.Equal(t, sourcerule.SourceObservationRejected, (*seen)[0].Result)
}

func TestSourceObservationFinalizationRechecksScopeAndInstalledInput(t *testing.T) {
	for _, name := range []string{"workspace revision", "workspace deletion", "device runtime revision", "dirty unknown"} {
		t.Run(name, func(t *testing.T) {
			s, r, req, seen := newSourceObservationRuntime(t)
			out := s.ReconcileSourceRule(t.Context(), req)
			require.True(t, out.Observation.Verified)
			var finalErr error
			switch name {
			case "workspace revision":
				r.record.DatabaseSetupRevision = "later"
			case "workspace deletion":
				r.record = nil
			case "device runtime revision":
				s.markDeviceProjectionAligned("D", &workspace.RuntimeProjection{Version: "later"})
			case "dirty unknown":
				finalErr = &modbusshare.Error{Code: modbusshare.ErrCodeDirtyUnknown}
			}
			s.FinalizeSourceRuleObservation(t.Context(), out, finalErr)
			require.Len(t, *seen, 1)
			require.Equal(t, sourcerule.SourceObservationUnavailable, (*seen)[0].Result)
		})
	}
}

func TestSourceObservationDeleteUsesOriginalRevisionWithAbsentCurrentRule(t *testing.T) {
	s, r, req, seen := newSourceObservationRuntime(t)
	r.projection.Rules = nil
	req.Operation = sourcerule.RuntimeReconcileOperationDelete
	out := s.ReconcileSourceRule(t.Context(), req)
	s.FinalizeSourceRuleObservation(t.Context(), out, nil)
	require.Len(t, *seen, 1)
	require.Equal(t, sourcerule.SourceObservationApplied, (*seen)[0].Result)
	require.Equal(t, "source-R3", (*seen)[0].RuleRevision)
}

func TestSourceObservationForeignFinalOutcomeNeverBorrowsProof(t *testing.T) {
	s, _, req, seen := newSourceObservationRuntime(t)
	out := s.ReconcileSourceRule(t.Context(), req)
	out.Scope.RuleID = "foreign"
	s.FinalizeSourceRuleObservation(t.Context(), out, &modbusshare.Error{Code: modbusshare.ErrCodeReconcileFailed})
	require.Len(t, *seen, 1)
	require.Equal(t, sourcerule.SourceObservationUnavailable, (*seen)[0].Result)
}

func TestSourceObservationCallbackPanicPreservesActualOutcome(t *testing.T) {
	s, _, req, _ := newSourceObservationRuntime(t)
	s.sourceRuleObserver = func(context.Context, sourcerule.RuntimeSourceObservation) { panic("private-source-error") }
	out := s.ReconcileSourceRule(t.Context(), req)
	require.NotPanics(t, func() { s.FinalizeSourceRuleObservation(t.Context(), out, nil) })
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, out.Status)
	require.True(t, s.sourceInputStillInstalled("D", "input-fingerprint-not-W2"))
}

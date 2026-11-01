package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type acquisitionMutationRecorder struct {
	*reconcileRuntimeSync
	begin           func(context.Context, RuntimeReconcileScope) context.Context
	finish          func()
	begun, finished int
	active          bool
}

func (r *acquisitionMutationRecorder) BeginSourceRuleAcquisitionMutation(ctx context.Context, scope RuntimeReconcileScope) (nextContext context.Context, finalize func()) {
	r.begun++
	r.active = true
	if r.begin != nil {
		ctx = r.begin(ctx, scope)
	}
	return ctx, func() {
		r.finished++
		r.active = false
		if r.finish != nil {
			r.finish()
		}
	}
}

func (r *acquisitionMutationRecorder) ReconcileSourceRule(ctx context.Context, request RuntimeReconcileRequest) RuntimeReconcileOutcome {
	if !r.active || ctx.Err() != nil {
		return RuntimeReconcileOutcome{Status: RuntimeReconcileStatusFailed}
	}
	return r.reconcileRuntimeSync.ReconcileSourceRule(ctx, request)
}

func TestAcquisitionMutationLifecycleCoversEveryPublicSourceMutation(t *testing.T) {
	for _, operation := range []string{"create", "update", "disable", "enable", "delete", "preflight", "canceled", "compensation"} {
		t.Run(operation, func(t *testing.T) {
			f := newRuntimeReconcileRuleFixture(t, true)
			r := &acquisitionMutationRecorder{reconcileRuntimeSync: f.runtimeSync}
			f.svc.runtimeSync = r
			r.begin = func(ctx context.Context, scope RuntimeReconcileScope) context.Context {
				require.Equal(t, f.deviceID, scope.DeviceID)
				before, err := f.svc.GetByID(ctx, f.rule.ID)
				require.NoError(t, err)
				require.Equal(t, f.rule.Enabled, before.Enabled, "quarantine begins before any durable write")
				require.Equal(t, f.rule.NamingPrefix, before.NamingPrefix)
				if operation == "canceled" {
					canceled, cancel := context.WithCancel(ctx)
					cancel()
					return canceled
				}
				return ctx
			}
			if operation == "compensation" {
				f.svc.SetShareRuntimeReconciler(failingShareReconciler{})
				r.finish = func() {
					restored, err := f.svc.GetByID(f.ctx, f.rule.ID)
					require.NoError(t, err)
					require.Equal(t, f.rule.NamingPrefix, restored.NamingPrefix, "finalize must follow compensation")
				}
			}
			var err error
			switch operation {
			case "create":
				_, _, err = f.svc.CreateWithRuntimeReconcile(f.ctx, CreateRuleRequest{ID: "new-rule", DeviceID: f.deviceID, StartAddress: "40002", Count: 1, DataType: schema.DataTypeInt16, Enabled: true})
			case "disable":
				_, err = f.svc.DisableWithRuntimeReconcile(f.ctx, f.rule.ID)
			case "enable":
				_, err = f.svc.EnableWithRuntimeReconcile(f.ctx, f.rule.ID)
			case "delete":
				_, err = f.svc.DeleteWithRuntimeReconcile(f.ctx, f.rule.ID)
			case "preflight":
				invalid := 0
				_, _, err = f.svc.UpdateWithRuntimeReconcile(f.ctx, f.rule.ID, UpdateRuleRequest{Count: &invalid})
			default:
				prefix := "new-prefix"
				_, _, err = f.svc.UpdateWithRuntimeReconcile(f.ctx, f.rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})
			}
			if operation == "preflight" || operation == "canceled" || operation == "compensation" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, 1, r.begun)
			require.Equal(t, 1, r.finished, "all returns must finalize the mutation lifecycle")
			require.False(t, r.active)
		})
	}
}

func TestAcquisitionMutationLifecycleAlsoCoversDirectSourceServices(t *testing.T) {
	for _, operation := range []string{"create", "update", "delete", "disable"} {
		t.Run(operation, func(t *testing.T) {
			f := newRuntimeReconcileRuleFixture(t, true)
			recorder := &acquisitionMutationRecorder{reconcileRuntimeSync: f.runtimeSync}
			f.svc.runtimeSync = recorder
			var err error
			switch operation {
			case "create":
				_, err = f.svc.Create(f.ctx, CreateRuleRequest{ID: "direct", DeviceID: f.deviceID, StartAddress: "40002", Count: 1, DataType: schema.DataTypeInt16})
			case "update":
				prefix := "changed"
				_, err = f.svc.Update(f.ctx, f.rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})
			case "delete":
				err = f.svc.Delete(f.ctx, f.rule.ID)
			case "disable":
				err = f.svc.setEnabled(f.ctx, f.rule.ID, false)
			}
			require.NoError(t, err)
			require.Equal(t, 1, recorder.begun)
			require.Equal(t, 1, recorder.finished)
			require.False(t, recorder.active)
		})
	}
}

func TestAcquisitionMutationLinksIgnoreOutputSelectorsButInvalidateSourceOwnership(t *testing.T) {
	f := newRuntimeReconcileRuleFixture(t, true)
	recorder := &acquisitionMutationRecorder{reconcileRuntimeSync: f.runtimeSync}
	f.svc.runtimeSync = recorder
	links, err := f.svc.ListLinks(f.ctx, f.rule.ID)
	require.NoError(t, err)
	require.NotEmpty(t, links)
	outputOnly := cloneSourceRuleLinks(links)
	tagID := "new-output-tag"
	outputOnly[0].TagID = &tagID
	require.NoError(t, f.svc.ReplaceLinks(f.ctx, f.rule.ID, outputOnly))
	require.Zero(t, recorder.begun, "output selectors must not quarantine acquisition")
	require.NoError(t, f.svc.ReplaceLinks(f.ctx, f.rule.ID, nil))
	require.Equal(t, 1, recorder.begun)
	require.Equal(t, 1, recorder.finished)
}

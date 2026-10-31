package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/sourcerule"
)

// newPointFormatRecovery uses the existing receipt, rule mutex, Share workspace
// serialization and runtime installation authority. It never starts recording.
func newPointFormatRecovery(services gatewayServices, reconciler *modbusshare.Reconciler) handlers.PointFormatRecovery {
	return func(ctx context.Context, req point.FormatRecoveryRequest) (out handlers.FormatRecoveryResult, resultErr error) {
		if services.point == nil || services.runtime == nil || services.sourceRule == nil || services.modbusShare == nil || reconciler == nil {
			return out, point.ErrFormatUnavailable
		}
		ctx, release := services.sourceRule.AcquireRuleMutation(ctx, req.RuleID)
		defer release()
		resultErr = reconciler.WithFormatRecovery(ctx, req.WorkspaceID, func(ctx context.Context) error {
			var err error
			out, err = recoverPointFormat(ctx, services, reconciler, req)
			return err
		})
		return
	}
}

func recoverPointFormat(ctx context.Context, services gatewayServices, reconciler *modbusshare.Reconciler, req point.FormatRecoveryRequest) (out handlers.FormatRecoveryResult, resultErr error) {
	runtime, share := services.runtime, services.modbusShare
	proof, ok := runtime.CaptureOrdinaryPoint(req.PointID)
	if !ok {
		return out, point.ErrFormatUnavailable
	}
	bindingProof, ok := runtime.CaptureFormatRecoveryBindings(req.PointID)
	if !ok {
		return out, point.ErrFormatUnavailable
	}
	settings := share.Settings()
	var source *sourcerule.ShareRecoverySnapshot
	evidence, err := services.point.ReadFormatRecovery(ctx, req, func(ctx context.Context, tx *sql.Tx) error {
		var e error
		source, e = sourcerule.ReadShareRecoverySnapshot(ctx, tx, req.WorkspaceID, settings)
		return e
	})
	if err != nil {
		return out, err
	}
	target := evidence.TargetInputs
	prepared, err := datalinkruntime.PrepareFormatRecoveryBindings(req.WorkspaceID, evidence.Point, evidence.Profile, target.Device, target.Mappings, target.Tags)
	if err != nil {
		return out, err
	}
	fence, err := share.CaptureFormatRecoveryFence(req.WorkspaceID, req.ExpectedWorkspaceRevision, req.ExpectedSettingsRevision)
	if err != nil {
		return out, err
	}
	// Completion binds the first accepted full input set, not merely the target
	// receipt. A later current snapshot may validate itself yet represent a new plan.
	key := req.WorkspaceID + "/" + req.PointID + "/" + req.RequestID + "/" + evidence.ReceiptFingerprint + "/" + source.WitnessDigest + "/" + target.WitnessDigest
	if !evidence.RuntimeRequired {
		key += fmt.Sprintf("/inactive:%d", proof.Epoch)
	}
	validate := func(ctx context.Context, tx *sql.Tx) error {
		return sourcerule.ValidateShareRecoverySnapshot(ctx, tx, req.WorkspaceID, settings, source.WitnessDigest)
	}
	validLease := func(lease uint64) bool {
		if !evidence.RuntimeRequired {
			return lease == 0 && runtime.MatchesFormatRecoveryAbsence(proof)
		}
		return runtime.MatchesFormatRecoveryBoundLease(req.PointID, lease, prepared)
	}
	// Successful replay is read-only, including its runtime lease and token.
	if !fence.Pending() {
		err = services.point.AcceptFormatRecovery(ctx, req, evidence, func() error {
			state, ok := share.ReplayFormatRecovery(key, validLease)
			if !ok {
				return point.ErrFormatConflict
			}
			out = formatRecoveryResult(state, evidence.RuntimeRequired, settings.Enabled, "replayed")
			return nil
		}, validate)
		return out, err
	}
	if settings.Enabled {
		_, err = reconciler.ReconcileFormatRecovery(ctx, modbusshare.ReconcileRequest{WorkspaceID: req.WorkspaceID, ExpectedWorkspaceRevision: req.ExpectedWorkspaceRevision, ExpectedSettingsRevision: req.ExpectedSettingsRevision, DesiredMappings: source.Desired}, fence, evidence.WorkspaceWitness)
		if err != nil {
			return out, err
		}
	}
	var lease uint64
	if evidence.RuntimeRequired {
		lease, ok = runtime.InstallPreparedFormatRecovery(&evidence.Point, proof, evidence.Profile)
		if !ok {
			return out, point.ErrFormatUnavailable
		}
	} else if !runtime.MatchesFormatRecoveryAbsence(proof) {
		return out, point.ErrFormatUnavailable
	}
	var published modbusshare.HydrationState
	var bindingReceipt datalinkruntime.FormatRecoveryBindingReceipt
	defer func() {
		if resultErr == nil {
			return
		}
		if published.ReadinessToken != "" {
			share.RevokeFormatRecoveryPublication(key, published.ReadinessToken, lease)
		}
		if lease != 0 {
			runtime.SuspendFormatRecovery(req.PointID, lease, bindingReceipt)
		}
	}()
	resultErr = services.point.AcceptFormatRecovery(ctx, req, evidence, func() error {
		var e error
		published, e = share.PublishFormatRecovery(ctx, fence, key, lease, func() bool {
			if !evidence.RuntimeRequired {
				return runtime.MatchesFormatRecoveryAbsence(proof)
			}
			var accepted bool
			bindingReceipt, accepted = runtime.ActivateFormatRecoveryBindings(req.PointID, lease, evidence.Profile, bindingProof, prepared)
			return accepted
		})
		return e
	}, validate)
	if resultErr != nil {
		return out, errors.Join(point.ErrFormatUnavailable, resultErr)
	}
	return formatRecoveryResult(published, evidence.RuntimeRequired, settings.Enabled, "recovered"), nil
}

func formatRecoveryResult(state modbusshare.HydrationState, runtimeRequired, shareEnabled bool, outcome string) handlers.FormatRecoveryResult {
	runtimeState := "not_started"
	if runtimeRequired {
		runtimeState = "installed"
	}
	return handlers.FormatRecoveryResult{Outcome: outcome, RuntimeState: runtimeState, ShareApplied: shareEnabled, Hydration: state}
}

package grouppipeline

import (
	"context"
	"errors"

	"go-gateway/internal/datalink/groupdelivery"
)

func validObservedQuota(q groupdelivery.QuotaStatus, key groupdelivery.GroupKey) bool {
	if !q.Configured || q.MaxBytes <= 0 || q.UsedBytes < 0 || q.MaxBytes > 9007199254740991 || q.UsedBytes > 9007199254740991 {
		return false
	}
	if q.Scope == groupdelivery.QuotaScopeGroup {
		if q.GroupID != key.GroupID {
			return false
		}
	} else if q.Scope != groupdelivery.QuotaScopeGlobal || q.GroupID != "" {
		return false
	}
	switch q.State {
	case groupdelivery.QuotaHard:
		return q.IntakeRefused
	case groupdelivery.QuotaOK, groupdelivery.QuotaWarning:
		return !q.IntakeRefused
	}
	return false
}

// The actual wrapped refusal, not a later most-constrained-scope read, supplies
// quota bytes. Disk-full's best-effort max=0 remains unavailable, never capacity.
func (p *Pipeline) recordQuotaRefusal(ctx context.Context, m *managed, err error) {
	if !p.quotaEnabled() {
		return
	}
	var refusal *groupdelivery.QuotaError
	if !errors.As(err, &refusal) {
		return
	}
	d := intakeDescriptor{key: groupdelivery.GroupKey{WorkspaceID: m.workspaceID, GroupID: m.groupID, GroupRevision: m.revision}, kind: IntakeQuotaRefused, managed: m, until: m.boundary.ReplayState().Until}
	if refusal.Reason == groupdelivery.ReasonQuotaHardLimit {
		d.quota = groupdelivery.QuotaStatus{Configured: true, State: groupdelivery.QuotaHard, Scope: refusal.Scope, GroupID: refusal.GroupID, UsedBytes: refusal.UsedBytes, MaxBytes: refusal.MaxBytes, IntakeRefused: true}
	}
	p.recordIntakeObservation(ctx, d)
}

// Reconcile is the ONLY recovery trigger. The broker owns the existing 15s
// scheduling and exact retained episode token; no hint means no metadata reads.
// This runs after reconcileMu has unlocked but before the shared batch flush.
func (p *Pipeline) observeQuotaRecovery(ctx context.Context) {
	if !p.quotaEnabled() || ctx.Err() != nil || p.deps.IntakeQuotaObservation.TakeQuotaRecoveryHint == nil {
		return
	}
	defer func() {
		if recover() != nil {
			p.recordIntakeObservation(ctx, intakeDescriptor{kind: IntakeQuotaRecovered})
		}
	}()
	hint, ok := p.deps.IntakeQuotaObservation.TakeQuotaRecoveryHint(ctx)
	if !ok {
		return
	}
	if (hint.Scope != groupdelivery.QuotaScopeGlobal && hint.Scope != groupdelivery.QuotaScopeGroup) || hint.Observe == nil || !safeOwnerID(hint.WorkspaceRevision) || !safeOwnerID(hint.ConnectorID) || !safeOwnerID(hint.ConnectorRevision) || !safeOwnerID(hint.ReferenceID) {
		p.recordIntakeObservation(ctx, intakeDescriptor{kind: IntakeQuotaRecovered})
		return
	}
	p.recordIntakeObservation(ctx, intakeDescriptor{key: hint.Key, kind: IntakeQuotaRecovered, hint: &hint})
}

// One shared queue and one loss owner preserve the original admission bound.
func (p *Pipeline) observationBatchPolicy() (limit int, loss func(uint64), enabled bool) {
	owner, intake := p.deps.OwnerObservation, p.deps.IntakeQuotaObservation
	if owner.enabled() {
		limit, loss, enabled = owner.RecordLimit, owner.OnDrop, true
	}
	if intake.Observe != nil && (intake.IntakeEnabled || intake.QuotaEnabled) {
		if !enabled {
			return intake.RecordLimit, intake.OnDrop, true
		}
		if intake.RecordLimit < limit {
			limit = intake.RecordLimit
		}
		if loss == nil {
			loss = intake.OnDrop
		}
	}
	return
}

func (p *Pipeline) flushIntakeObservation(ctx context.Context, b *ownerObservationBatch, d intakeDescriptor) {
	observation, emit := p.verifyIntakeObservation(ctx, d)
	if !emit {
		return
	}
	callback := p.deps.IntakeQuotaObservation.Observe
	if observation.Available && d.kind == IntakeQuotaRecovered && d.hint != nil {
		callback = d.hint.Observe
	}
	defer func() {
		if recover() != nil {
			b.drop()
		}
	}()
	callback(observation)
}

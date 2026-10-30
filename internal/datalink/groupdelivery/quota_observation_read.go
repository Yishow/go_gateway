package groupdelivery

import (
	"context"
	"fmt"
)

// QuotaStatusForScope measures only the requested quota scope. Group usage
// includes every revision of that group; global usage includes every group.
// The key identifies the affected group, not a revision-specific byte budget or
// proof of current ownership. Callers must verify their frozen identity separately.
// A missing policy for this exact scope returns an unconfigured status, never
// the other scope's measurement. This read does not update the intake usage cache.
func (s *Store) QuotaStatusForScope(ctx context.Context, key GroupKey, scope QuotaScope) (QuotaStatus, error) {
	if !key.Valid() {
		return QuotaStatus{}, ErrInvalidGroupKey
	}
	if scope != QuotaScopeGroup && scope != QuotaScopeGlobal {
		return QuotaStatus{}, fmt.Errorf("read quota scope: %w", ErrInvalidQuota)
	}
	if s == nil || s.db == nil || ctx == nil {
		return QuotaStatus{}, fmt.Errorf("read quota scope: store or context unavailable")
	}
	if err := ctx.Err(); err != nil {
		return QuotaStatus{}, err
	}
	if s.quota == nil {
		return QuotaStatus{}, nil
	}
	groupID, limit := "", s.quota.GlobalMaxBytes
	if scope == QuotaScopeGroup {
		groupID, limit = key.GroupID, s.quota.GroupMaxBytes
	}
	if limit == 0 {
		return QuotaStatus{}, nil
	}
	if limit < 0 {
		return QuotaStatus{}, fmt.Errorf("read quota scope: %w", ErrInvalidQuota)
	}
	used, err := measureUsage(ctx, s.db, groupID)
	if err != nil {
		return QuotaStatus{}, err
	}
	status := QuotaStatus{Configured: true, Scope: scope, GroupID: groupID,
		UsedBytes: used, MaxBytes: limit, State: s.quota.scopeState(used, limit)}
	switch status.State {
	case QuotaHard:
		status.IntakeRefused = true
		status.LossRiskNotice = "new readings for this scope are not recorded until space is available; already accepted data is kept"
	case QuotaWarning:
		status.LossRiskNotice = "capacity is running low; intake will be refused at the hard limit"
	}
	return status, nil
}

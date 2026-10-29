package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
)

func TestRuntimeDeliveryQuotaCopiesMeasuredScopeWithoutInference(t *testing.T) {
	for _, tc := range []struct {
		name, scope, state string
		used, maximum      int64
		refused            bool
	}{
		{"group-zero", "group", "ok", 0, 100, false},
		{"global-warning", "global", "warning", 83, 100, false},
		{"group-overshoot", "group", "hard_limit", 123, 100, true},
		{"global-safe-bound", "global", "hard_limit", 9007199254740991, 100, true},
		// Custom producer thresholds, not a browser/projection-derived percentage.
		{"producer-state", "group", "warning", 1, 100, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeDeliveryProjectionFixture()
			raw.Quota = groupdelivery.QuotaStatusView{Configured: true, Scope: groupdelivery.QuotaScope(tc.scope), State: tc.state,
				UsedBytes: tc.used, MaxBytes: tc.maximum, IntakeRefused: tc.refused, LossRiskNotice: "password=private-quota-secret"}
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.NotNil(t, got.Current.Quota)
			require.Equal(t, &runtimeview.QuotaEvidence{Scope: tc.scope, State: tc.state, UsedBytes: tc.used, MaxBytes: tc.maximum, IntakeRefused: tc.refused}, got.Current.Quota)
			require.Equal(t, "A", got.Current.Binding.EntityRevision)
			require.Nil(t, got.Current.Intake, "quota refusal is not exact revision intake evidence")
			require.Zero(t, got.Current.Stages.SQLCommitted)
			require.NotContains(t, string(runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)), "private-quota-secret")
			got.Current.Quota.UsedBytes = 7
			require.Equal(t, tc.used, raw.Quota.UsedBytes, "detached projection must not mutate the authority")
		})
	}
}

func TestRuntimeDeliveryQuotaUnconfiguredCannotInventScope(t *testing.T) {
	for _, state := range []string{"", "unconfigured"} {
		t.Run(state, func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeDeliveryProjectionFixture()
			raw.Quota = groupdelivery.QuotaStatusView{State: state}
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.Nil(t, got.Current.Quota)
			require.Contains(t, string(runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)), `"quota":null`)
		})
	}
}

func TestRuntimeDeliveryQuotaRejectsMalformedObservations(t *testing.T) {
	base := groupdelivery.QuotaStatusView{Configured: true, Scope: groupdelivery.QuotaScopeGroup, State: groupdelivery.QuotaOK, UsedBytes: 0, MaxBytes: 100}
	for name, change := range map[string]func(*groupdelivery.QuotaStatusView){
		"empty-scope":                   func(q *groupdelivery.QuotaStatusView) { q.Scope = "" },
		"foreign-scope-label":           func(q *groupdelivery.QuotaStatusView) { q.Scope = "device-A" },
		"unsafe-scope":                  func(q *groupdelivery.QuotaStatusView) { q.Scope = "db.private:5432" },
		"unknown-state":                 func(q *groupdelivery.QuotaStatusView) { q.State = "private-driver-error" },
		"configured-unconfigured-state": func(q *groupdelivery.QuotaStatusView) { q.State = "unconfigured" },
		"negative-used":                 func(q *groupdelivery.QuotaStatusView) { q.UsedBytes = -1 },
		"unsafe-used":                   func(q *groupdelivery.QuotaStatusView) { q.UsedBytes = 9007199254740992 },
		"zero-maximum":                  func(q *groupdelivery.QuotaStatusView) { q.MaxBytes = 0 },
		"negative-maximum":              func(q *groupdelivery.QuotaStatusView) { q.MaxBytes = -1 },
		"unsafe-maximum":                func(q *groupdelivery.QuotaStatusView) { q.MaxBytes = 9007199254740992 },
		"ok-refused":                    func(q *groupdelivery.QuotaStatusView) { q.IntakeRefused = true },
		"warning-refused":               func(q *groupdelivery.QuotaStatusView) { q.State = groupdelivery.QuotaWarning; q.IntakeRefused = true },
		"hard-not-refused":              func(q *groupdelivery.QuotaStatusView) { q.State = groupdelivery.QuotaHard },
		"unconfigured-with-scope": func(q *groupdelivery.QuotaStatusView) {
			*q = groupdelivery.QuotaStatusView{State: "unconfigured", Scope: groupdelivery.QuotaScopeGroup}
		},
		"unconfigured-with-used": func(q *groupdelivery.QuotaStatusView) {
			*q = groupdelivery.QuotaStatusView{State: "unconfigured", UsedBytes: 1}
		},
		"unconfigured-with-limit": func(q *groupdelivery.QuotaStatusView) {
			*q = groupdelivery.QuotaStatusView{State: "unconfigured", MaxBytes: 100}
		},
		"unconfigured-with-refusal": func(q *groupdelivery.QuotaStatusView) {
			*q = groupdelivery.QuotaStatusView{State: "unconfigured", IntakeRefused: true}
		},
		"unconfigured-with-state": func(q *groupdelivery.QuotaStatusView) { *q = groupdelivery.QuotaStatusView{State: "warning"} },
	} {
		t.Run(name, func(t *testing.T) {
			saved, frozen, scope, raw, _ := runtimeDeliveryProjectionFixture()
			raw.Quota = base
			change(&raw.Quota)
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.Nil(t, got)
			require.True(t, runtimeview.IsContractError(err, runtimeview.ReasonUnavailable), "%v", err)
		})
	}
}

func TestRuntimeDeliveryQuotaMissingEffectiveProofDoesNotInventCurrent(t *testing.T) {
	saved, _, scope, raw, at := runtimeDeliveryProjectionFixture()
	raw.RevisionStages = nil
	raw.Quota = groupdelivery.QuotaStatusView{Configured: true, Scope: groupdelivery.QuotaScopeGlobal, State: groupdelivery.QuotaWarning, UsedBytes: 80, MaxBytes: 100}
	got, err := projectRuntimeRuleDelivery(raw, saved, nil, scope, nil)
	require.NoError(t, err)
	require.Nil(t, got.Current, "a quota observation cannot establish current revision or zero stage counts")
	runtimeDeliveryMarshal(t, saved, nil, scope, got, at)
}

func TestRuntimeDeliveryQuotaCannotBorrowEditableStorageLineage(t *testing.T) {
	saved, frozen, scope, raw, at := runtimeDeliveryProjectionFixture()
	raw.Quota = groupdelivery.QuotaStatusView{Configured: true, Scope: groupdelivery.QuotaScopeGroup, State: groupdelivery.QuotaWarning, UsedBytes: 80, MaxBytes: 100}
	raw.Storage = &grouppipeline.StorageView{Generation: 9, HistoryMayBeLost: true, Reconstructions: 4, RepairState: "blocked", RepairReason: "/private/storage-secret.sqlite"}
	got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
	require.NoError(t, err)
	require.NotNil(t, got.Current.Quota)
	require.Nil(t, got.Current.Storage, "editable-group storage lacks the required frozen destination/revision lineage proof")
	wire := string(runtimeDeliveryMarshal(t, saved, frozen, scope, got, at))
	require.Contains(t, wire, `"storage":null`)
	require.NotContains(t, wire, "storage-secret")
}

func TestRuntimeDeliveryQuotaDoesNotReplaceExistingRevisionProof(t *testing.T) {
	saved, frozen, scope, raw, _ := runtimeDeliveryProjectionFixture()
	raw.Quota = groupdelivery.QuotaStatusView{Configured: true, Scope: groupdelivery.QuotaScopeGlobal, State: groupdelivery.QuotaOK, UsedBytes: 0, MaxBytes: 100}
	raw.RevisionStages.GroupRevision = "B"
	got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
	require.Nil(t, got)
	require.True(t, runtimeview.IsContractError(err, runtimeview.ReasonStale))
}

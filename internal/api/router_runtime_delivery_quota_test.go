package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/snapshot"
)

func runtimeQuotaConfigure(t *testing.T, r *runtimeDeliveryRig, config groupdelivery.QuotaConfig) {
	t.Helper()
	store, err := r.store.WithQuota(config)
	require.NoError(t, err)
	r.store = store
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{Groups: r.f.groups, Store: store,
		Destinations: dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))},
		grouppipeline.Config{Now: func() time.Time { return r.clock }})
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	r.router = NewRouter(&DatalinkServices{Workspace: r.f.workspace, WriteGroups: r.f.groups, WriteGroupDelivery: r.pipe,
		Device: device.NewService(device.NewSQLRepository(r.f.db), nil)})
}

func runtimeQuotaMeasured(t *testing.T, r *runtimeDeliveryRig, groupID string) int64 {
	t.Helper()
	value, err := r.store.Usage(t.Context(), groupID)
	require.NoError(t, err)
	return value
}

func TestRuntimeDeliveryQuotaSQLiteGroupIncludesJournalAndHistoricalOutbox(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.accept(t, "quota-current-sample")
	groupID := r.f.group["id"].(string)
	journal := runtimeQuotaMeasured(t, r, groupID)
	require.Positive(t, journal)
	r.historicalRows(t, 2)
	used := runtimeQuotaMeasured(t, r, groupID)
	require.Greater(t, used, journal)
	var storedJournal, storedOutbox int64
	require.NoError(t, r.f.db.QueryRow(`SELECT COALESCE(SUM(LENGTH(CAST(payload AS BLOB))),0) FROM wg_delivery_samples WHERE consumed=0 AND group_id=?`, groupID).Scan(&storedJournal))
	require.NoError(t, r.f.db.QueryRow(`SELECT COALESCE(SUM(LENGTH(CAST(payload AS BLOB))),0) FROM wg_delivery_outbox WHERE state!='sql_committed' AND group_id=?`, groupID).Scan(&storedOutbox))
	require.Equal(t, storedJournal+storedOutbox, used)
	runtimeQuotaConfigure(t, r, groupdelivery.QuotaConfig{GroupMaxBytes: used * 2, WarningRatio: 0.4, HardRatio: 0.9})
	g := runtimeTimingReadOnly(t, r)
	require.NotNil(t, g.Delivery.Current.Quota, "production scoped route must not discard actual quota")
	require.Equal(t, &runtimeview.QuotaEvidence{Scope: "group", State: "warning", UsedBytes: used, MaxBytes: used * 2}, g.Delivery.Current.Quota)
	require.Equal(t, r.f.group["applied_revision"], g.Delivery.Current.Binding.EntityRevision)
	require.Zero(t, g.Delivery.Current.Stages.Queued, "quota bytes include retained revisions, not just current output")
	require.Len(t, g.Delivery.Backlog, 1)
	require.Equal(t, "rev7", g.Delivery.Backlog[0].Binding.EntityRevision)
	require.EqualValues(t, 2, g.Delivery.Backlog[0].Pending)
	require.Nil(t, g.Delivery.Current.Storage)
}

func runtimeQuotaOtherGroupRow(t *testing.T, r *runtimeDeliveryRig) {
	t.Helper()
	start := r.clock.Add(-time.Minute)
	key := groupdelivery.GroupKey{WorkspaceID: r.f.record.ID, GroupID: "other-quota-group", GroupRevision: "other-rev"}
	outcome := snapshot.Outcome{Kind: snapshot.OutcomeRow, RecordID: "other-quota-record", EffectKey: "other-quota-effect", BucketStart: start, BucketEnd: start.Add(15 * time.Second)}
	row := dbtarget.EncodedRow{RecordID: outcome.RecordID, EffectKey: outcome.EffectKey, BucketStart: start, Cells: []dbtarget.EncodedCell{{Column: "value", Value: float64(123)}}}
	require.NoError(t, r.store.CommitClosure(t.Context(), key, groupdelivery.Closure{
		Destination: groupdelivery.FrozenDestination{Scope: "other-scope", ConnectorID: "other-C", ConnectorRevision: "1", Database: "/private/quota-other.sqlite", TableName: "other_table", DedupeCapability: "none"},
		NextClose:   start.Add(15 * time.Second), Buckets: []groupdelivery.ClosedBucket{{Outcome: outcome, Row: &row}},
	}))
}

func TestRuntimeDeliveryQuotaSQLiteGlobalDoesNotAllocateOtherGroupUsage(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.historicalRows(t, 1)
	own := runtimeQuotaMeasured(t, r, r.f.group["id"].(string))
	runtimeQuotaOtherGroupRow(t, r)
	global := runtimeQuotaMeasured(t, r, "")
	require.Greater(t, global, own)
	runtimeQuotaConfigure(t, r, groupdelivery.QuotaConfig{GroupMaxBytes: global * 10, GlobalMaxBytes: 1})
	g := runtimeTimingReadOnly(t, r)
	require.NotNil(t, g.Delivery.Current.Quota)
	require.Equal(t, &runtimeview.QuotaEvidence{Scope: "global", State: "hard_limit", UsedBytes: global, MaxBytes: 1, IntakeRefused: true}, g.Delivery.Current.Quota)
	require.Greater(t, g.Delivery.Current.Quota.UsedBytes, g.Delivery.Current.Quota.MaxBytes, "existing accepted data may exceed configured capacity; never clamp")
	require.Len(t, g.Delivery.Backlog, 1, "other-group identity/data is not exposed through global quota")
	require.EqualValues(t, 1, g.Delivery.Backlog[0].Pending)
	require.Zero(t, g.Delivery.Current.Stages.SQLCommitted)
}

func TestRuntimeDeliveryQuotaSQLiteUnconfiguredAndObservedZeroStayDistinct(t *testing.T) {
	for _, configured := range []bool{false, true} {
		t.Run(map[bool]string{false: "unconfigured", true: "configured-zero"}[configured], func(t *testing.T) {
			r := newRuntimeDeliveryRig(t)
			if configured {
				runtimeQuotaConfigure(t, r, groupdelivery.QuotaConfig{GroupMaxBytes: 100})
			}
			g := runtimeTimingReadOnly(t, r)
			if !configured {
				require.Nil(t, g.Delivery.Current.Quota, "empty producer scope cannot be fabricated as group")
				return
			}
			require.Equal(t, &runtimeview.QuotaEvidence{Scope: "group", State: "ok", UsedBytes: 0, MaxBytes: 100}, g.Delivery.Current.Quota)
		})
	}
}

func TestRuntimeDeliveryQuotaSQLiteSharedMembersKeepWholeGroupBytes(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.sharedGroup(t)
	r.accept(t, "quota-two-members")
	used := runtimeQuotaMeasured(t, r, r.f.group["id"].(string))
	require.Positive(t, used)
	var samples int
	require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed=0`).Scan(&samples))
	require.Equal(t, 2, samples)
	runtimeQuotaConfigure(t, r, groupdelivery.QuotaConfig{GroupMaxBytes: used * 2, WarningRatio: 0.4, HardRatio: 0.9})
	g := runtimeTimingReadOnly(t, r)
	require.True(t, g.Delivery.Current.Shared)
	require.NotNil(t, g.Delivery.Current.Quota)
	require.Equal(t, "group", g.Delivery.Current.Quota.Scope)
	require.Equal(t, used, g.Delivery.Current.Quota.UsedBytes, "do not divide group quota bytes between selected device/Tag members")
}

func TestRuntimeDeliveryQuotaSQLiteMissingFrozenProofKeepsBacklogOnly(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.historicalRows(t, 1)
	runtimeQuotaConfigure(t, r, groupdelivery.QuotaConfig{GroupMaxBytes: 1})
	_, err := r.f.db.Exec(`DELETE FROM wg_runtime_versions`)
	require.NoError(t, err)
	g := runtimeTimingReadOnly(t, r)
	require.Nil(t, g.Effective)
	require.Nil(t, g.Delivery.Current, "quota must not create a synthetic current binding or zero stages")
	require.Len(t, g.Delivery.Backlog, 1)
	require.Equal(t, "old_table", g.Delivery.Backlog[0].Destination.Table)
}

func TestRuntimeDeliveryQuotaMalformedProviderIsTyped503(t *testing.T) {
	for name, q := range map[string]groupdelivery.QuotaStatusView{
		"unknown-scope":              {Configured: true, Scope: "db.private:5432", State: "ok", MaxBytes: 100},
		"unknown-state":              {Configured: true, Scope: groupdelivery.QuotaScopeGlobal, State: "private-driver-error", MaxBytes: 100},
		"negative-used":              {Configured: true, Scope: groupdelivery.QuotaScopeGroup, State: "ok", UsedBytes: -1, MaxBytes: 100},
		"unsafe-count":               {Configured: true, Scope: groupdelivery.QuotaScopeGroup, State: "ok", MaxBytes: 9007199254740992},
		"contradictory-refusal":      {Configured: true, Scope: groupdelivery.QuotaScopeGroup, State: "hard_limit", UsedBytes: 99, MaxBytes: 100},
		"unconfigured-with-capacity": {State: "unconfigured", MaxBytes: 100},
	} {
		t.Run(name, func(t *testing.T) {
			r := newRuntimeDeliveryRig(t)
			q.LossRiskNotice = "password=private-quota-secret"
			r.useReader(&runtimeDeliveryFault{Pipeline: r.pipe, change: func(v *grouppipeline.DeliveryView) { v.Quota = q }})
			before := lifecycleAPIState(t, r.f.db)
			response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
			require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), "unavailable")
			for _, secret := range []string{"private-quota-secret", "private-driver-error", "db.private", r.f.file, `"rule_evidence"`} {
				require.NotContains(t, response.Body.String(), secret)
			}
			require.Equal(t, before, lifecycleAPIState(t, r.f.db))
		})
	}
}

func TestRuntimeDeliveryQuotaSQLiteDoesNotBorrowEditableStorage(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	runtimeQuotaConfigure(t, r, groupdelivery.QuotaConfig{GroupMaxBytes: 100})
	r.useReader(&runtimeDeliveryFault{Pipeline: r.pipe, change: func(v *grouppipeline.DeliveryView) {
		v.Storage = &grouppipeline.StorageView{Generation: 7, HistoryMayBeLost: true, Reconstructions: 3, RepairState: "blocked", RepairReason: "/private/unproven-storage.sqlite"}
		v.Quota.LossRiskNotice = "password=private-quota-secret"
	}})
	g := runtimeTimingReadOnly(t, r)
	require.NotNil(t, g.Delivery.Current.Quota)
	require.Nil(t, g.Delivery.Current.Storage, "raw editable lineage is not frozen destination proof")
}

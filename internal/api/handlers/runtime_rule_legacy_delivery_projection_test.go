package handlers

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/runtimeview"
)

func legacyDeliverySafeFixture() *grouppipeline.DeliveryView {
	at := time.Date(2026, 10, 29, 1, 2, 3, 456789000, time.UTC)
	receipt := &groupdelivery.CommittedEffect{GroupRevision: "rev7", ConnectorRevision: "7", RecordID: "record1", EffectKey: "effect1", PayloadDigest: strings.Repeat("a", 64), CommittedAt: at}
	record := &groupdelivery.SourceRecord{Available: true, Delivery: groupdelivery.SourceDeliveryCommitted, VerifiedSourceRecord: &groupdelivery.VerifiedSourceRecord{GroupRevision: "rev7", ConnectorRevision: "7", RecordID: "record1", EffectKey: "effect1", PayloadDigest: receipt.PayloadDigest, ProfileKind: "instantaneous", WindowStart: at.Add(-15 * time.Second), WindowEnd: at, Status: groupdelivery.SourceStatusPartial, Reason: "quality-bad", RawPresent: true, CommittedAt: &at, Aggregation: &groupdelivery.SourceAggregation{Method: measurement.MeanAggregationMethod, Version: measurement.MeanAggregationVersion, Basis: "float64", SourceType: "float64", ResultType: "float64", ValidCount: 2, RejectedCount: 1, Coverage: "partial", Reason: "bad"}}}
	return &grouppipeline.DeliveryView{GroupID: "G", Intake: grouppipeline.IntakeView{State: "blocked", Reason: "snapshot-unavailable"}, Stages: grouppipeline.StagesView{Collecting: 1, Queued: 2, Retrying: 3, Blocked: 4, Quarantined: 5, Unknown: 6, SQLCommitted: 7, Skipped: 8}, LastSQLCommittedAt: &at, LastSQLCommittedEffect: receipt, RevisionStages: &grouppipeline.RevisionStagesView{GroupRevision: "rev7", Stages: grouppipeline.StagesView{SQLCommitted: 7}}, OldestPendingSeconds: 125.5, NoDataBuckets: 9, SkippedBuckets: 10,
		Backlog:                  []grouppipeline.BacklogView{{GroupRevision: "rev6", ConnectorID: "C1", ConnectorRevision: "7", TableSchema: "main", TableName: "old_table", Pending: 3, ErrorCodes: []string{"target-blocked"}}},
		Attention:                []groupdelivery.AttentionItem{{EffectKey: "effect1", State: "blocked", StateRevision: 11, PayloadDigest: receipt.PayloadDigest, GroupRevision: "rev7", ConnectorID: "C1", ConnectorRevision: "7", TableSchema: "main", TableName: "old_table", ErrorCode: "target-blocked"}},
		CounterReplayQuarantines: []grouppipeline.CounterReplayQuarantine{{WorkspaceID: "W", GroupID: "G", GroupRevision: "rev5", PendingSamples: 12, Reason: string(runtimeview.ReasonCounterMissing), Action: grouppipeline.CounterReplayRepairAction}},
		RecentBucketIssues:       []groupdelivery.BucketIssue{{GroupRevision: "rev7", BucketStart: at, Kind: "skipped", Causes: []string{"bad"}}},
		Quota:                    groupdelivery.QuotaStatusView{Configured: true, State: "warning", Scope: "global", UsedBytes: 12, MaxBytes: 15, LossRiskNotice: "capacity is running low; intake will be refused at the hard limit"},
		Storage:                  &grouppipeline.StorageView{Generation: 2, HistoryMayBeLost: true, Reconstructions: 1, LastReconstructedAt: &at, RepairState: "blocked", RepairReason: recordingplan.StorageReasonReceiptsMissing, Events: []recordingplan.StorageEvent{{Generation: 2, GroupID: "G", Kind: "repair_blocked", Code: recordingplan.StorageCodeRepairBlocked, Reason: recordingplan.StorageReasonReceiptsMissing, OperationID: "op1", Occurrences: 1, OccurredAt: at}}},
		SourceEvidence:           &grouppipeline.SourceEvidenceView{State: grouppipeline.SourceEvidenceAvailable, GroupRevision: "rev7", ConnectorID: "C1", ConnectorRevision: "7", Sources: []grouppipeline.SourceEvidenceItem{{CircuitID: "circuit1", DeviceID: "device1", PointID: "point1", TagID: "tag1", Latest: record, LastSQLCommitted: record}}}}
}
func TestRuntimeLegacyDeliveryProjectionPreservesProofAndDetaches(t *testing.T) {
	raw := legacyDeliverySafeFixture()
	before, err := json.Marshal(raw)
	require.NoError(t, err)
	got, err := projectLegacyRuntimeDelivery(raw, "G")
	require.NoError(t, err)
	after, err := json.Marshal(got)
	require.NoError(t, err)
	require.JSONEq(t, string(before), string(after))
	got.Backlog[0].ErrorCodes[0] = "changed"
	got.Attention[0].ErrorCode = "changed"
	got.Storage.Events[0].Reason = "changed"
	got.SourceEvidence.Sources[0].Latest.Aggregation.Reason = "changed"
	*got.LastSQLCommittedAt = time.Time{}
	unchanged, err := json.Marshal(raw)
	require.NoError(t, err)
	require.Equal(t, string(before), string(unchanged))
}

// Walk every populated string leaf in this fixture of the actual wire DTO.
// Future nil optional objects still require explicit fixture/audit expansion.
func legacyDeliveryStrings(value reflect.Value, path string, result map[string]reflect.Value) {
	if value.Type() == reflect.TypeFor[time.Time]() {
		return
	}
	switch value.Kind() {
	case reflect.Pointer:
		if !value.IsNil() {
			legacyDeliveryStrings(value.Elem(), path, result)
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			if value.Type().Field(i).IsExported() {
				legacyDeliveryStrings(value.Field(i), path+"."+value.Type().Field(i).Name, result)
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			legacyDeliveryStrings(value.Index(i), path+"[]", result)
		}
	case reflect.String:
		result[path] = value
	}
}
func TestRuntimeLegacyDeliveryProjectionAuditsEveryStringBranch(t *testing.T) {
	fields := map[string]reflect.Value{}
	legacyDeliveryStrings(reflect.ValueOf(legacyDeliverySafeFixture()), "", fields)
	require.Greater(t, len(fields), 60)
	for path := range fields {
		t.Run(path, func(t *testing.T) {
			raw := legacyDeliverySafeFixture()
			values := map[string]reflect.Value{}
			legacyDeliveryStrings(reflect.ValueOf(raw), "", values)
			const hostile = "password=fixture-secret /private/raw.sql [::1]:502"
			values[path].SetString(hostile)
			before, _ := json.Marshal(raw)
			got, err := projectLegacyRuntimeDelivery(raw, "G")
			after, _ := json.Marshal(raw)
			require.Equal(t, string(before), string(after))
			diagnostic := strings.HasSuffix(path, ".Reason") || strings.HasSuffix(path, ".ErrorCodes[]") || strings.HasSuffix(path, ".ErrorCode") || strings.HasSuffix(path, ".Causes[]") || strings.HasSuffix(path, ".LossRiskNotice") || strings.HasSuffix(path, ".RepairReason") || strings.HasSuffix(path, ".Events[].Code")
			if strings.Contains(path, "CounterReplayQuarantines") {
				diagnostic = false
			}
			if diagnostic {
				require.NoError(t, err)
				encoded, e := json.Marshal(got)
				require.NoError(t, e)
				require.NotContains(t, string(encoded), "fixture-secret")
				require.Contains(t, string(encoded), "unavailable")
			} else {
				require.Error(t, err)
				require.True(t, runtimeview.IsContractError(err, runtimeview.ReasonUnavailable))
				require.Nil(t, got)
				require.NotContains(t, err.Error(), "fixture-secret")
			}
		})
	}
}
func TestRuntimeLegacyDeliveryProjectionKnownReasonCompatibility(t *testing.T) {
	for _, value := range append(append([]string{}, legacyDeliveryReasons...), "recording-ownership-conflict", "counter-replay-proof-conflict") {
		require.Equal(t, value, legacyDeliveryReason(value))
	}
	for _, value := range legacyDeliverySourceReasons {
		require.Equal(t, value, legacyDeliverySourceReason(value))
	}
	for _, value := range []string{"password=secret", "future-safe-looking-code", "/private/file"} {
		require.Equal(t, "unavailable", legacyDeliveryReason(value))
		require.Equal(t, "reason_unavailable", legacyDeliverySourceReason(value))
	}
	require.Empty(t, legacyDeliveryReason(""))
	require.Empty(t, legacyDeliverySourceReason(""))
}
func TestRuntimeLegacyDeliveryProjectionUnavailableAndEmptyCompatibility(t *testing.T) {
	for _, raw := range []*grouppipeline.DeliveryView{{GroupID: "G"}, {GroupID: "G", SourceEvidence: &grouppipeline.SourceEvidenceView{State: "unavailable", GroupRevision: "rev7", UnavailableReason: "applied_snapshot_unavailable", Sources: []grouppipeline.SourceEvidenceItem{}}}} {
		got, err := projectLegacyRuntimeDelivery(raw, "G")
		require.NoError(t, err)
		require.Equal(t, raw, got)
	}
	for _, raw := range []*grouppipeline.DeliveryView{nil, {GroupID: "other"}} {
		got, err := projectLegacyRuntimeDelivery(raw, "G")
		require.Error(t, err)
		require.Nil(t, got)
	}
	source := legacyDeliverySafeFixture()
	source.SourceEvidence.State = "unavailable"
	source.SourceEvidence.UnavailableReason = "password=secret"
	source.SourceEvidence.Sources = nil
	got, err := projectLegacyRuntimeDelivery(source, "G")
	require.Error(t, err)
	require.Nil(t, got)
	for _, cause := range []groupdelivery.SourceUnavailableCause{groupdelivery.SourceCausePayloadInvalid, groupdelivery.SourceCauseDigestMismatch, groupdelivery.SourceCauseIdentityMismatch, groupdelivery.SourceCauseStatusUnknown, groupdelivery.SourceCauseWindowInvalid} {
		require.True(t, legacyDeliverySourceRecordSafe(&groupdelivery.SourceRecord{Cause: cause}))
	}
}

func TestRuntimeLegacyDeliveryProjectionTypedUnavailableAndBoundaryCodes(t *testing.T) {
	for _, reason := range []string{"group-missing", "sink-missing", "group-not-applied", "member-identity-incomplete", "tag-type-unsupported", "layout-blocked", "destination-scope-invalid", "restore-failed", "dedupe-unsupported"} {
		t.Run(reason, func(t *testing.T) { require.Equal(t, reason, legacyDeliveryReason(reason)) })
	}
	for _, reason := range []string{"not_long_layout", "applied_snapshot_unavailable", "revision_mismatch"} {
		t.Run(reason, func(t *testing.T) {
			raw := &grouppipeline.DeliveryView{GroupID: "G", SourceEvidence: &grouppipeline.SourceEvidenceView{State: "unavailable", GroupRevision: "rev7", UnavailableReason: reason, Sources: []grouppipeline.SourceEvidenceItem{}}}
			got, err := projectLegacyRuntimeDelivery(raw, "G")
			require.NoError(t, err)
			require.Equal(t, reason, got.SourceEvidence.UnavailableReason)
		})
	}
	t.Run("unknown-enum", func(t *testing.T) {
		raw := &grouppipeline.DeliveryView{GroupID: "G", SourceEvidence: &grouppipeline.SourceEvidenceView{State: "unavailable", GroupRevision: "rev7", UnavailableReason: "password=secret"}}
		got, err := projectLegacyRuntimeDelivery(raw, "G")
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestRuntimeLegacyDeliveryProjectionPreservesStorageLineageReasons(t *testing.T) {
	// The storage-lineage producer emits these literals on successful recreation
	// and reconstruction; they are diagnostics, not new permission to repair.
	for _, reason := range []string{"table-created-again", "table-rebuilt"} {
		t.Run(reason, func(t *testing.T) {
			raw := legacyDeliverySafeFixture()
			raw.Storage.Events[0].Reason = reason
			got, err := projectLegacyRuntimeDelivery(raw, "G")
			require.NoError(t, err)
			require.Equal(t, reason, got.Storage.Events[0].Reason)
			require.Equal(t, reason, raw.Storage.Events[0].Reason)
		})
	}
}

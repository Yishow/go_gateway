package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func runtimeSourceDeliveryFixture() (*workspace.WriteGroup, *workspace.WriteGroupAppliedSnapshot, runtimeview.Scope, *grouppipeline.DeliveryView, time.Time) {
	saved, frozen, scope, raw, at := runtimeDeliveryProjectionFixture()
	frozen.Group.LayoutVersion = workspace.WriteGroupLayoutManagedLongV1
	frozen.Group.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond, SampleReduction: workspace.WriteGroupSampleReductionMean}
	frozen.Group.Members[0].CircuitID = "circuit-A"
	frozen.Group.Members[0].ValueRole = workspace.WriteGroupValueRole("gauge")
	frozen.Group.Members[0].SourceValueType = schema.DataTypeFloat64
	frozen.Group.Members[0].ResultValueType = measurement.ExactFloat64
	record := &groupdelivery.SourceRecord{Available: true, Delivery: groupdelivery.SourceDeliveryQueued, VerifiedSourceRecord: &groupdelivery.VerifiedSourceRecord{
		GroupRevision: "A", ConnectorRevision: "7", RecordID: "source-record-2", EffectKey: "source-effect-2", PayloadDigest: "source-digest-2", ProfileKind: "instantaneous", WindowStart: at.Add(-15 * time.Second), WindowEnd: at, Status: groupdelivery.SourceStatusPartial, RawPresent: true,
		Aggregation: &groupdelivery.SourceAggregation{Method: measurement.MeanAggregationMethod, Version: measurement.MeanAggregationVersion, Basis: string(measurement.MeanBasisFloat64), SourceType: "float64", ResultType: string(measurement.ExactFloat64), ValidCount: 1, Coverage: "partial"},
	}}
	raw.SourceEvidence = &grouppipeline.SourceEvidenceView{State: grouppipeline.SourceEvidenceAvailable, GroupRevision: "A", ConnectorID: "C1", ConnectorRevision: "7", Sources: []grouppipeline.SourceEvidenceItem{{DeviceID: "device-A", PointID: "point-A", TagID: "tag-A", CircuitID: "circuit-A", Latest: record}}}
	return saved, frozen, scope, raw, at
}

func TestRuntimeSourceDeliveryMeanPartialAndUnknownCoverageNeverComplete(t *testing.T) {
	for _, coverage := range []string{"partial", "unknown", "missing"} {
		t.Run(coverage, func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
			record := raw.SourceEvidence.Sources[0].Latest
			record.Aggregation.Coverage = coverage
			if coverage == "missing" {
				record.Status = groupdelivery.SourceStatusMissing
				record.RawPresent = false
				record.Aggregation.ValidCount = 0
			}
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.Equal(t, "available", got.Current.SourceState)
			require.Len(t, got.Current.Sources, 1)
			latest := got.Current.Sources[0].Latest
			require.Equal(t, coverage, latest.State)
			if coverage == "missing" {
				require.Nil(t, latest.Coverage)
			} else {
				require.Equal(t, coverage, *latest.Coverage)
			}
			require.Equal(t, "2026-10-01T12:00:00Z", *latest.WindowStart)
			require.Equal(t, "2026-10-01T12:00:15Z", *latest.WindowEnd)
			require.Nil(t, latest.Receipt)
			require.Nil(t, latest.DeltaKnown)
			require.Nil(t, got.TestWrite)
			runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)
		})
	}
}

func TestRuntimeSourceDeliveryDeltaUnknownNeverZero(t *testing.T) {
	for _, known := range []bool{false, true} {
		t.Run(map[bool]string{false: "unknown", true: "known"}[known], func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
			frozen.Group.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileCumulative, Interval: 10, Unit: workspace.WriteGroupIntervalUnitMinute}
			frozen.Group.RowPolicy.IntervalSeconds = 600
			record := raw.SourceEvidence.Sources[0].Latest
			record.ProfileKind, record.Status, record.Aggregation, record.DeltaKnown = "cumulative", groupdelivery.SourceStatusOK, nil, &known
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.Len(t, got.Current.Sources, 1)
			latest := got.Current.Sources[0].Latest
			require.Equal(t, known, *latest.DeltaKnown)
			require.Equal(t, map[bool]string{false: "unknown", true: "available"}[known], latest.State)
			wire := string(runtimeDeliveryMarshal(t, saved, frozen, scope, got, at))
			require.NotContains(t, wire, "delta_value")
			require.NotContains(t, wire, "raw_present")
		})
	}
}

func TestRuntimeSourceDeliveryLatestQueuedSeparateFromEarlierCommittedWindow(t *testing.T) {
	saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
	committed := *raw.SourceEvidence.Sources[0].Latest.VerifiedSourceRecord
	committed.RecordID, committed.EffectKey, committed.PayloadDigest = "source-record-1", "source-effect-1", "source-digest-1"
	committed.WindowStart, committed.WindowEnd, committed.CommittedAt = at.Add(-30*time.Second), at.Add(-15*time.Second), &at
	raw.SourceEvidence.Sources[0].LastSQLCommitted = &groupdelivery.SourceRecord{Available: true, Delivery: groupdelivery.SourceDeliveryCommitted, VerifiedSourceRecord: &committed}
	raw.RevisionStages.Stages.SQLCommitted = 10
	// Group's latest receipt may belong to a different source/window; never allocate it to each Tag.
	raw.LastSQLCommittedEffect = &groupdelivery.CommittedEffect{GroupRevision: "A", ConnectorRevision: "7", RecordID: "group-latest", EffectKey: "group-effect", PayloadDigest: "group-digest", CommittedAt: at}
	got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
	require.NoError(t, err)
	require.Len(t, got.Current.Sources, 1)
	source := got.Current.Sources[0]
	require.Nil(t, source.Latest.Receipt)
	require.Equal(t, "2026-10-01T12:00:00Z", *source.Latest.WindowStart)
	require.Equal(t, "2026-10-01T11:59:45Z", *source.LastCommitted.WindowStart)
	require.Equal(t, "source-record-1", source.LastCommitted.Receipt.RecordID)
	require.Equal(t, "group-latest", got.Current.LastCommit.RecordID)
	require.Equal(t, "runtime", source.LastCommitted.Receipt.Provenance)
	require.Nil(t, got.TestWrite)
	runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)
}

func TestRuntimeSourceDeliveryMissingRecordNoZeroWindowOrReceipt(t *testing.T) {
	for _, mode := range []string{"no-record", "unavailable", "missing-verified"} {
		t.Run(mode, func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
			switch mode {
			case "no-record":
				raw.SourceEvidence.Sources[0].Latest = nil
			case "unavailable":
				raw.SourceEvidence.Sources[0].Latest.Available = false
			case "missing-verified":
				raw.SourceEvidence.Sources[0].Latest.VerifiedSourceRecord = nil
			}
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.Len(t, got.Current.Sources, 1)
			latest := got.Current.Sources[0].Latest
			if mode == "no-record" {
				require.Nil(t, latest)
			} else {
				require.Equal(t, &runtimeview.SourceResult{State: "unavailable"}, latest)
			}
			require.Nil(t, got.Current.Sources[0].LastCommitted)
			runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)
		})
	}
}

func TestRuntimeSourceDeliveryLegacyCustomMissingFrozenProofUnavailable(t *testing.T) {
	for _, mode := range []string{"nil", "no-snapshot", "pipeline-now-B"} {
		t.Run(mode, func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
			raw.SourceEvidence = nil
			if mode != "nil" {
				raw.SourceEvidence = &grouppipeline.SourceEvidenceView{State: grouppipeline.SourceEvidenceUnavailable, GroupRevision: "A", UnavailableReason: grouppipeline.SourceEvidenceNoSnapshot}
			}
			if mode == "pipeline-now-B" {
				raw.SourceEvidence.UnavailableReason = grouppipeline.SourceEvidenceRevisionMismatch
			}
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.Equal(t, "unavailable", got.Current.SourceState)
			require.Empty(t, got.Current.Sources)
			require.Equal(t, "A", got.Current.Binding.EntityRevision)
			require.EqualValues(t, 3, got.Current.Stages.Queued)
			runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)
		})
	}
}

func TestRuntimeSourceDeliverySameGroupOtherDeviceNeverGetsSelectedSourceSuccess(t *testing.T) {
	saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
	frozen.Group.Members = append(frozen.Group.Members, workspace.WriteGroupMember{DeviceID: "device-B", PointID: "point-B", TagID: "tag-B", CircuitID: "circuit-B", TargetColumn: "other"})
	raw.SourceEvidence.Sources = append(raw.SourceEvidence.Sources, grouppipeline.SourceEvidenceItem{DeviceID: "device-B", PointID: "point-B", TagID: "tag-B", CircuitID: "circuit-B"})
	saved.Members = []workspace.WriteGroupMember{{DeviceID: "device-C", PointID: "point-C", TagID: "tag-C", TargetColumn: "draft"}}
	before, err := json.Marshal(raw)
	require.NoError(t, err)
	got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
	require.NoError(t, err)
	require.True(t, got.Current.Shared)
	require.Len(t, got.Current.Sources, 2)
	require.Equal(t, "device-B", got.Current.Sources[1].DeviceID)
	require.Nil(t, got.Current.Sources[1].Latest)
	require.Nil(t, got.Current.Sources[1].LastCommitted)
	after, err := json.Marshal(raw)
	require.NoError(t, err)
	require.Equal(t, string(before), string(after))
	runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)
}

func TestRuntimeSourceDeliveryForeignDuplicateTupleAndUnsafeProofFailClosed(t *testing.T) {
	cases := map[string]struct {
		mutate func(*grouppipeline.SourceEvidenceView)
		reason runtimeview.Reason
	}{
		"unknown-envelope":          {func(v *grouppipeline.SourceEvidenceView) { v.State = "future" }, runtimeview.ReasonUnavailable},
		"mixed-envelope":            {func(v *grouppipeline.SourceEvidenceView) { v.UnavailableReason = "revision_mismatch" }, runtimeview.ReasonUnavailable},
		"group-revision":            {func(v *grouppipeline.SourceEvidenceView) { v.GroupRevision = "B" }, runtimeview.ReasonStale},
		"connector":                 {func(v *grouppipeline.SourceEvidenceView) { v.ConnectorID = "C2" }, runtimeview.ReasonStale},
		"connector-revision":        {func(v *grouppipeline.SourceEvidenceView) { v.ConnectorRevision = "8" }, runtimeview.ReasonStale},
		"foreign-device":            {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].DeviceID = "device-B" }, runtimeview.ReasonForeign},
		"foreign-point":             {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].PointID = "point-B" }, runtimeview.ReasonForeign},
		"foreign-tag":               {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].TagID = "tag-B" }, runtimeview.ReasonForeign},
		"foreign-circuit":           {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].CircuitID = "circuit-B" }, runtimeview.ReasonForeign},
		"duplicate-tuple":           {func(v *grouppipeline.SourceEvidenceView) { v.Sources = append(v.Sources, v.Sources[0]) }, runtimeview.ReasonUnavailable},
		"unsafe-device":             {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].DeviceID = "http://private-host" }, runtimeview.ReasonUnavailable},
		"unsafe-circuit":            {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].CircuitID = "/private/config" }, runtimeview.ReasonUnavailable},
		"unsafe-record":             {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.RecordID = "postgres://secret" }, runtimeview.ReasonUnavailable},
		"unsafe-effect":             {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.EffectKey = "/private/secret" }, runtimeview.ReasonUnavailable},
		"unsafe-digest":             {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.PayloadDigest = "password=secret" }, runtimeview.ReasonUnavailable},
		"record-revision":           {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.GroupRevision = "B" }, runtimeview.ReasonStale},
		"record-connector-revision": {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.ConnectorRevision = "8" }, runtimeview.ReasonStale},
		"wrong-profile":             {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.ProfileKind = "cumulative" }, runtimeview.ReasonUnavailable},
		"unknown-status":            {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Status = "future" }, runtimeview.ReasonUnavailable},
		"unknown-delivery":          {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Delivery = "future" }, runtimeview.ReasonUnavailable},
		"available-with-cause": {func(v *grouppipeline.SourceEvidenceView) {
			v.Sources[0].Latest.Cause = groupdelivery.SourceCauseDigestMismatch
		}, runtimeview.ReasonUnavailable},
		"zero-start": {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.WindowStart = time.Time{} }, runtimeview.ReasonUnavailable},
		"empty-window": {func(v *grouppipeline.SourceEvidenceView) {
			v.Sources[0].Latest.WindowEnd = v.Sources[0].Latest.WindowStart
		}, runtimeview.ReasonUnavailable},
		"reversed-window": {func(v *grouppipeline.SourceEvidenceView) {
			v.Sources[0].Latest.WindowEnd = v.Sources[0].Latest.WindowStart.Add(-time.Second)
		}, runtimeview.ReasonUnavailable},
		"complete-unproved":    {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation.Coverage = "complete" }, runtimeview.ReasonUnavailable},
		"unknown-coverage":     {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation.Coverage = "future" }, runtimeview.ReasonUnavailable},
		"unknown-method":       {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation.Method = "future" }, runtimeview.ReasonUnavailable},
		"unknown-version":      {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation.Version = "future" }, runtimeview.ReasonUnavailable},
		"unknown-basis":        {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation.Basis = "future" }, runtimeview.ReasonUnavailable},
		"wrong-source-type":    {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation.SourceType = "int16" }, runtimeview.ReasonUnavailable},
		"wrong-result-type":    {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation.ResultType = "decimal" }, runtimeview.ReasonUnavailable},
		"negative-count":       {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation.ValidCount = -1 }, runtimeview.ReasonUnavailable},
		"negative-rejected":    {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation.RejectedCount = -1 }, runtimeview.ReasonUnavailable},
		"missing-with-reading": {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation.Coverage = "missing" }, runtimeview.ReasonUnavailable},
		"missing-mean":         {func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].Latest.Aggregation = nil }, runtimeview.ReasonUnavailable},
		"instantaneous-delta":  {func(v *grouppipeline.SourceEvidenceView) { b := false; v.Sources[0].Latest.DeltaKnown = &b }, runtimeview.ReasonUnavailable},
		"queued-with-commit-time": {func(v *grouppipeline.SourceEvidenceView) {
			at := v.Sources[0].Latest.WindowEnd
			v.Sources[0].Latest.CommittedAt = &at
		}, runtimeview.ReasonUnavailable},
		"committed-without-proof": {func(v *grouppipeline.SourceEvidenceView) {
			v.Sources[0].Latest.Delivery = groupdelivery.SourceDeliveryCommitted
		}, runtimeview.ReasonUnavailable},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			saved, frozen, scope, raw, _ := runtimeSourceDeliveryFixture()
			tc.mutate(raw.SourceEvidence)
			_, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.True(t, runtimeview.IsContractError(err, tc.reason), "expected %s, got %v", tc.reason, err)
		})
	}
}

func TestRuntimeSourceDeliveryReceiptMismatchNeverBorrowsGroupCommit(t *testing.T) {
	for _, mode := range []string{"unavailable", "missing-verified", "queued", "no-time", "zero-time", "missing-group-proof", "wrong-group", "wrong-connector"} {
		t.Run(mode, func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
			record := raw.SourceEvidence.Sources[0].Latest
			record.Delivery, record.CommittedAt = groupdelivery.SourceDeliveryCommitted, &at
			raw.SourceEvidence.Sources[0].LastSQLCommitted = record
			raw.SourceEvidence.Sources[0].Latest = nil
			raw.RevisionStages.Stages.SQLCommitted = 1
			raw.LastSQLCommittedEffect = &groupdelivery.CommittedEffect{GroupRevision: "A", ConnectorRevision: "7", RecordID: "group-record", EffectKey: "group-effect", PayloadDigest: "group-digest", CommittedAt: at}
			want := runtimeview.ReasonUnavailable
			switch mode {
			case "unavailable":
				record.Available = false
			case "missing-verified":
				record.VerifiedSourceRecord = nil
			case "queued":
				record.Delivery = groupdelivery.SourceDeliveryQueued
				record.CommittedAt = nil
			case "no-time":
				record.CommittedAt = nil
			case "zero-time":
				record.CommittedAt = &time.Time{}
			case "missing-group-proof":
				raw.LastSQLCommittedEffect = nil
				raw.RevisionStages.Stages.SQLCommitted = 0
			case "wrong-group":
				record.GroupRevision = "B"
				want = runtimeview.ReasonStale
			case "wrong-connector":
				record.ConnectorRevision = "8"
				want = runtimeview.ReasonStale
			}
			_, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.True(t, runtimeview.IsContractError(err, want), "%v", err)
		})
	}
}

func TestRuntimeSourceDeliveryLegacyLayoutsIgnoreMisleadingSourceEnvelope(t *testing.T) {
	for _, mode := range []string{"legacy-wide", "custom", "no-profile"} {
		t.Run(mode, func(t *testing.T) {
			saved, frozen, scope, raw, _ := runtimeSourceDeliveryFixture()
			switch mode {
			case "legacy-wide":
				frozen.Group.LayoutVersion = workspace.WriteGroupLayoutManagedWideV1
			case "custom":
				frozen.Group.LayoutVersion = "custom"
			case "no-profile":
				frozen.Group.RecordingProfile = nil
			}
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.Equal(t, "unavailable", got.Current.SourceState)
			require.Empty(t, got.Current.Sources)
		})
	}
}

func TestRuntimeSourceDeliveryVerifiedQualityStatesKeepActualReceipts(t *testing.T) {
	statuses := map[groupdelivery.SourceStatus]string{groupdelivery.SourceStatusOK: "available", groupdelivery.SourceStatusPartial: "partial", groupdelivery.SourceStatusMissing: "missing", groupdelivery.SourceStatusBaseline: "unknown", groupdelivery.SourceStatusBad: "unknown", groupdelivery.SourceStatusStale: "unknown", groupdelivery.SourceStatusInvalid: "unknown", groupdelivery.SourceStatusUncertain: "unknown", groupdelivery.SourceStatusGap: "unknown"}
	for status, want := range statuses {
		t.Run(string(status), func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
			frozen.Group.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionLatest
			record := raw.SourceEvidence.Sources[0].Latest
			record.Aggregation, record.Status = nil, status
			record.RawPresent = status == groupdelivery.SourceStatusOK || status == groupdelivery.SourceStatusPartial || status == groupdelivery.SourceStatusBaseline
			record.Delivery, record.CommittedAt = groupdelivery.SourceDeliveryCommitted, &at
			raw.SourceEvidence.Sources[0].LastSQLCommitted = record
			raw.RevisionStages.Stages.SQLCommitted = 1
			raw.LastSQLCommittedEffect = &groupdelivery.CommittedEffect{GroupRevision: "A", ConnectorRevision: "7", RecordID: "group-record", EffectKey: "group-effect", PayloadDigest: "group-digest", CommittedAt: at}
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.Equal(t, want, got.Current.Sources[0].Latest.State)
			require.NotNil(t, got.Current.Sources[0].Latest.Receipt)
			require.NotNil(t, got.Current.Sources[0].LastCommitted.Receipt)
			runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)
		})
	}
}

func TestRuntimeSourceDeliveryNonCommitStagesNeverBorrowReceipt(t *testing.T) {
	for _, stage := range []groupdelivery.SourceDelivery{groupdelivery.SourceDeliveryQueued, groupdelivery.SourceDeliveryRetrying, groupdelivery.SourceDeliveryBlocked, groupdelivery.SourceDeliveryQuarantined, groupdelivery.SourceDeliveryUnknown, groupdelivery.SourceDeliverySkipped} {
		t.Run(string(stage), func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
			raw.SourceEvidence.Sources[0].Latest.Delivery = stage
			raw.RevisionStages.Stages.SQLCommitted = 1
			raw.LastSQLCommittedEffect = &groupdelivery.CommittedEffect{GroupRevision: "A", ConnectorRevision: "7", RecordID: "group-record", EffectKey: "group-effect", PayloadDigest: "group-digest", CommittedAt: at}
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.Nil(t, got.Current.Sources[0].Latest.Receipt)
			require.Nil(t, got.Current.Sources[0].LastCommitted)
		})
	}
}

func TestRuntimeSourceDeliveryContradictoryRecordProofFailsClosed(t *testing.T) {
	for _, mode := range []string{"missing-with-raw", "ok-without-raw", "cumulative-without-delta", "cumulative-with-mean", "baseline-with-delta", "instantaneous-unknown-profile", "mean-without-gauge", "aggregation-on-latest", "year-outside-rfc3339"} {
		t.Run(mode, func(t *testing.T) {
			saved, frozen, scope, raw, _ := runtimeSourceDeliveryFixture()
			record := raw.SourceEvidence.Sources[0].Latest
			switch mode {
			case "missing-with-raw":
				record.Status = groupdelivery.SourceStatusMissing
			case "ok-without-raw":
				record.Status = groupdelivery.SourceStatusOK
				record.RawPresent = false
			case "cumulative-without-delta", "cumulative-with-mean", "baseline-with-delta":
				frozen.Group.RecordingProfile.Kind = workspace.WriteGroupProfileCumulative
				record.ProfileKind = "cumulative"
				if mode != "cumulative-with-mean" {
					record.Aggregation = nil
				}
				if mode != "cumulative-without-delta" {
					known := true
					record.DeltaKnown = &known
				}
				if mode == "baseline-with-delta" {
					record.Status = groupdelivery.SourceStatusBaseline
				}
			case "instantaneous-unknown-profile":
				frozen.Group.RecordingProfile.Kind = "future"
				record.ProfileKind = "future"
			case "mean-without-gauge":
				frozen.Group.Members[0].ValueRole = "state"
			case "aggregation-on-latest":
				frozen.Group.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionLatest
			case "year-outside-rfc3339":
				record.WindowStart = time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)
				record.WindowEnd = record.WindowStart.Add(time.Second)
			}
			_, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.True(t, runtimeview.IsContractError(err, runtimeview.ReasonUnavailable), "%v", err)
		})
	}
}

func TestRuntimeSourceDeliverySafeValueFreeDetachedProjection(t *testing.T) {
	saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
	record := raw.SourceEvidence.Sources[0].Latest
	record.Reason = "postgres://secret-host/password"
	record.Aggregation.Reason = "/private/credentials"
	// Exact-decimal is also a producer-supported mean basis, not a coerced float.
	record.Aggregation.Basis, record.Aggregation.ResultType = string(measurement.MeanBasisExact), string(measurement.ExactDecimal)
	frozen.Group.Members[0].ResultValueType = measurement.ExactDecimal
	got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
	require.NoError(t, err)
	wire := string(runtimeDeliveryMarshal(t, saved, frozen, scope, got, at))
	for _, forbidden := range []string{"secret-host", "credentials", "raw_present", "valid_count", "rejected_count", "aggregation", "value_decimal", "value_float", "provenance_json"} {
		require.NotContains(t, wire, forbidden)
	}
	*got.Current.Sources[0].Latest.Coverage = "unknown"
	require.Equal(t, "partial", record.Aggregation.Coverage)
	known := false
	frozen.Group.RecordingProfile.Kind = workspace.WriteGroupProfileCumulative
	record.ProfileKind, record.Aggregation, record.DeltaKnown = "cumulative", nil, &known
	got, err = projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
	require.NoError(t, err)
	*got.Current.Sources[0].Latest.DeltaKnown = true
	require.False(t, *record.DeltaKnown)
}

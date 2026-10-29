package handlers

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
)

func TestRuntimeDeliveryTimingCopiesPositiveObservedAgeWithoutRounding(t *testing.T) {
	for name, age := range map[string]float64{"fraction": 75.25, "subsecond": 0.125, "smallest-positive": math.SmallestNonzeroFloat64, "maximum": 9007199254740991} {
		t.Run(name, func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeDeliveryProjectionFixture()
			raw.OldestPendingSeconds = age
			// The age and counters are separate observations, not an atomic equality invariant.
			raw.RevisionStages.Stages = grouppipeline.StagesView{}
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.NotNil(t, got.OldestPendingSeconds)
			require.Equal(t, age, *got.OldestPendingSeconds)
			require.Equal(t, "group-all-revisions", got.OldestPendingScope)
			runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)
			// The output owns the nullable value; projection must not lend a pointer into its authority.
			*got.OldestPendingSeconds = 12
			require.Equal(t, age, raw.OldestPendingSeconds)
		})
	}
}

func TestRuntimeDeliveryTimingDisabledAndMissingEffectiveProofRetainHistoricalAge(t *testing.T) {
	for _, missingProof := range []bool{false, true} {
		t.Run(map[bool]string{false: "disabled", true: "missing-effective-proof"}[missingProof], func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeDeliveryProjectionFixture()
			saved.Status = workspace.WriteGroupStatusDisabled
			raw.OldestPendingSeconds = 75.25
			raw.Backlog = []grouppipeline.BacklogView{{GroupRevision: "old", ConnectorID: "C1", ConnectorRevision: "6", TableName: "old_table", Pending: 3}}
			// These independently observed counters do not prove or disprove the group-wide age.
			raw.RevisionStages.Stages = grouppipeline.StagesView{}
			if missingProof {
				frozen = nil
				raw.RevisionStages = nil
			}
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.Equal(t, 75.25, *got.OldestPendingSeconds)
			require.Equal(t, "group-all-revisions", got.OldestPendingScope)
			require.EqualValues(t, 3, got.Backlog[0].Pending)
			require.Equal(t, "6", got.Backlog[0].Destination.ConnectorRevision)
			if missingProof {
				require.Nil(t, got.Current)
			} else {
				require.NotNil(t, got.Current)
			}
			runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)
		})
	}
}

func TestRuntimeDeliveryTimingZeroAgePresenceAmbiguityStaysUnavailable(t *testing.T) {
	for _, age := range []float64{0, math.Copysign(0, -1)} {
		t.Run(map[bool]string{false: "zero", true: "negative-zero"}[math.Signbit(age)], func(t *testing.T) {
			saved, frozen, scope, raw, at := runtimeDeliveryProjectionFixture()
			raw.OldestPendingSeconds = age
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.NoError(t, err)
			require.Nil(t, got.OldestPendingSeconds)
			require.Contains(t, string(runtimeDeliveryMarshal(t, saved, frozen, scope, got, at)), `"oldest_pending_seconds":null`)
		})
	}
}

func TestRuntimeDeliveryTimingInvalidAgeFailsClosedWithoutClamp(t *testing.T) {
	for name, age := range map[string]float64{"negative": -0.125, "nan": math.NaN(), "positive-infinity": math.Inf(1), "negative-infinity": math.Inf(-1), "above-bound": math.Nextafter(9007199254740991, math.Inf(1)), "maximum-float": math.MaxFloat64} {
		t.Run(name, func(t *testing.T) {
			saved, frozen, scope, raw, _ := runtimeDeliveryProjectionFixture()
			raw.OldestPendingSeconds = age
			got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
			require.Nil(t, got)
			require.True(t, runtimeview.IsContractError(err, runtimeview.ReasonUnavailable), "%v", err)
		})
	}
}

func TestRuntimeDeliveryTimingNoSchedulingEvidenceDoesNotDeriveFrom15SecondInterval(t *testing.T) {
	saved, frozen, scope, raw, _ := runtimeDeliveryProjectionFixture()
	// A configured 15-second interval cannot turn a 10:00:01 observation into a 10:00:15 retry.
	observed := time.Date(2026, 10, 1, 10, 0, 1, 0, time.UTC)
	frozen.EffectiveAt = observed.Add(-time.Minute)
	require.Equal(t, 15, frozen.Group.RowPolicy.IntervalSeconds)
	raw.OldestPendingSeconds = 75.25
	got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
	require.NoError(t, err)
	require.Nil(t, got.Current.NextDueAt)
	require.Nil(t, got.Current.NextRetryAt)
	require.Nil(t, got.Current.AttemptDurationMS)
	wire := string(runtimeDeliveryMarshal(t, saved, frozen, scope, got, observed))
	for _, field := range []string{"next_due_at", "next_retry_at", "attempt_duration_ms"} {
		require.Contains(t, wire, `"`+field+`":null`)
	}
	require.NotContains(t, wire, "10:00:15Z")
}

func TestRuntimeDeliveryTimingSourceWindowAndCommittedTimesRetainExactExistingProof(t *testing.T) {
	saved, frozen, scope, raw, at := runtimeSourceDeliveryFixture()
	raw.OldestPendingSeconds = 75.25
	record := raw.SourceEvidence.Sources[0].Latest
	committed := at.Add(2250 * time.Millisecond)
	record.Delivery, record.CommittedAt = groupdelivery.SourceDeliveryCommitted, &committed
	raw.SourceEvidence.Sources[0].LastSQLCommitted = record
	raw.RevisionStages.Stages.SQLCommitted = 1
	raw.LastSQLCommittedEffect = &groupdelivery.CommittedEffect{GroupRevision: "A", ConnectorRevision: "7", RecordID: record.RecordID, EffectKey: record.EffectKey, PayloadDigest: record.PayloadDigest, CommittedAt: committed}
	got, err := projectRuntimeRuleDelivery(raw, saved, frozen, scope, nil)
	require.NoError(t, err)
	require.Equal(t, "2026-10-01T12:00:00Z", *got.Current.Sources[0].Latest.WindowStart)
	require.Equal(t, "2026-10-01T12:00:15Z", *got.Current.Sources[0].Latest.WindowEnd)
	require.Equal(t, "2026-10-01T12:00:17.25Z", got.Current.Sources[0].LastCommitted.Receipt.CommittedAt)
	require.Equal(t, "2026-10-01T12:00:17.25Z", got.Current.LastCommit.CommittedAt)
	require.Equal(t, 75.25, *got.OldestPendingSeconds)
	runtimeDeliveryMarshal(t, saved, frozen, scope, got, committed)
}

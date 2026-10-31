package handlers

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
	"strings"
	"testing"
	"time"
)

func TestRuntimeMixedRecordingFactsHaveNoCommonWindow(t *testing.T) {
	_, frozen, scope, _, _ := runtimeSourceDeliveryFixture()
	g := frozen.Group
	g.RecordID = "03760df3-f1e1-4c52-9adf-43785e54be31"
	g.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileMixed, ResultShape: workspace.WriteGroupResultShapeMemberCalculationV1}
	g.RowPolicy.IntervalSeconds = 0
	for i := range g.Members {
		g.Members[i].SourceRevision = "source-1"
		g.Members[i].MappingRevision = "mapping-1"
		g.Members[i].SourceValueType = schema.DataTypeFloat64
		g.Members[i].ResultValueType = measurement.ExactFloat64
		g.Members[i].Calculation = workspace.WriteGroupMemberCalculationMean
		g.Members[i].Cadence = &workspace.WriteGroupMemberCadence{Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond, Alignment: workspace.WriteGroupAlignmentUTCEpoch}
	}
	facts, err := runtimeRecordingFacts(g, frozen, scope)
	require.NoError(t, err)
	require.Nil(t, facts.RecordingProfile, "mixed is not a legacy profile with fake interval0")
	require.Nil(t, facts.IntervalSeconds)
	raw, err := json.Marshal(facts)
	require.NoError(t, err)
	var dto map[string]any
	require.NoError(t, json.Unmarshal(raw, &dto))
	mixed, ok := dto["mixed_recording"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, g.RecordID, mixed["record_id"])
	require.Equal(t, "member-calculation-v1", mixed["result_shape"])
	member := dto["members"].([]any)[0].(map[string]any)
	intent, ok := member["calculation"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "mean", intent["method"])
	require.Equal(t, float64(15), intent["interval"])
}

func TestRuntimeMixedSourceResultDoesNotTreatRawAsSelectedValue(t *testing.T) {
	for _, kind := range []string{"mean", "baseline", "not_due"} {
		t.Run(kind, func(t *testing.T) {
			_, frozen, _, raw, at := runtimeSourceDeliveryFixture()
			g := frozen.Group
			g.LayoutVersion = workspace.WriteGroupLayoutManagedWideProfileV1
			g.RecordID = "03760df3-f1e1-4c52-9adf-43785e54be31"
			g.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileMixed, ResultShape: workspace.WriteGroupResultShapeMemberCalculationV1}
			member := g.Members[0]
			member.Calculation = workspace.WriteGroupMemberCalculationMean
			member.Cadence = &workspace.WriteGroupMemberCadence{Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond, Alignment: workspace.WriteGroupAlignmentUTCEpoch}
			record := raw.SourceEvidence.Sources[0].Latest
			r := record.VerifiedSourceRecord
			id, err := snapshot.MixedDispatchRecordID(frozen.WorkspaceID, g.ID, frozen.AppliedRevision, at)
			require.NoError(t, err)
			r.RecordID = id
			r.ProfileKind = "mixed"
			start, end, next := at.Add(-15*time.Second), at, at.Add(15*time.Second)
			c := &groupdelivery.SourceCalculation{ResultShape: "member-calculation-v1", Method: "mean", IntervalSeconds: 15, Alignment: "utc_epoch", DueAt: at, NextDue: &next, OutcomeDigest: strings.Repeat("a", 64), OutcomeKind: "row", Status: "ok", ResultPresent: true, RowEmitted: true, RawPresent: true, SourceType: "float64", ResultType: "float64", AlignedStart: &start, AlignedEnd: &end, ActualStart: &start, ActualEnd: &end}
			r.Calculation = c
			r.WindowStart = start
			r.WindowEnd = end
			if kind == "baseline" {
				member.Calculation = workspace.WriteGroupMemberCalculationDelta
				c.Method = "delta"
				c.Status = "missing"
				c.ResultPresent = false
				c.RowEmitted = false
				c.OutcomeKind = "skipped"
				record.Delivery = groupdelivery.SourceDeliveryWindowSkipped
				r.EffectKey = ""
				r.PayloadDigest = ""
				r.Aggregation = nil
				r.Status = groupdelivery.SourceStatusBaseline
				known := false
				r.DeltaKnown = &known
			}
			if kind == "not_due" {
				// A different member dispatches inside this member's open window.
				c.DueAt = at.Add(-5 * time.Second)
				c.NextDue = &at
				r.RecordID, err = snapshot.MixedDispatchRecordID(frozen.WorkspaceID, g.ID, frozen.AppliedRevision, c.DueAt)
				require.NoError(t, err)
				c.Status = "not_due"
				c.ResultPresent = false
				c.RawPresent = false
				c.AlignedStart = nil
				c.AlignedEnd = nil
				c.ActualStart = nil
				c.ActualEnd = nil
				r.WindowStart = time.Time{}
				r.WindowEnd = time.Time{}
				r.RawPresent = false
				r.Aggregation = nil
				r.Status = groupdelivery.SourceStatus("not_due")
				record.Delivery = groupdelivery.SourceDeliveryCommitted
				r.CommittedAt = &at
			}
			result, err := runtimeSourceDeliveryResult(record, frozen, member, true)
			require.NoError(t, err)
			require.Equal(t, c.ResultPresent, result.Calculation.ResultPresent)
			require.Equal(t, c.RawPresent, result.Calculation.RawPresent)
			if kind == "baseline" {
				require.Equal(t, "unknown", result.State)
				require.Nil(t, result.Receipt)
			}
			if kind == "not_due" {
				require.Equal(t, "not_due", result.State)
				require.Nil(t, result.WindowStart)
				require.NotNil(t, result.Receipt)
				c.AlignedStart = &start
				_, err = runtimeSourceDeliveryResult(record, frozen, member, true)
				require.Error(t, err, "not_due cannot invent an aligned window")
				c.AlignedStart = nil
				c.NextDue = &c.DueAt
				_, err = runtimeSourceDeliveryResult(record, frozen, member, true)
				require.Error(t, err, "next due must be in the future when present")
				c.NextDue = nil
				_, err = runtimeSourceDeliveryResult(record, frozen, member, true)
				require.NoError(t, err, "ended member has no next due")
			}
			c.Method = "latest"
			_, err = runtimeSourceDeliveryResult(record, frozen, member, true)
			require.Error(t, err, "foreign calculation intent must not become selected-result evidence")
		})
	}
}

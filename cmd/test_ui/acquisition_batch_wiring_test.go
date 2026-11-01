package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

// The decorator observes the real ingestor->SampleSink boundary and still
// delegates to the production pipeline. Capturing is not a durable-write ACK.
type acquisitionWiringCapture struct {
	acceptanceCompletionObserver func(collector.AcquisitionCompletionObservation) // Installed before Start; bounded test-only collection.
	groupFixtureHooks
	pipe                     *grouppipeline.Pipeline
	samples                  chan measurement.SampleEnvelope
	observations             chan collector.AcquisitionObservation
	acceptanceInterest       *acceptanceInterestCapture             // Performance fixtures only.
	acceptanceSampleObserver func(measurement.SampleEnvelope)       // Optional test-only, installed before Start.
	acceptanceObserver       func(collector.AcquisitionObservation) // Optional test-only, installed before Start.
}

// This test wrapper installs its receiver before Start, after service wiring.
// It therefore enables detached copies even when a particular fixture leaves
// the receiver nil. Ordinary production wiring has no such interface and is nil.
func (f *acquisitionWiringCapture) acquisitionCompletionObserver() func(collector.AcquisitionCompletionObservation) {
	return func(observation collector.AcquisitionCompletionObservation) {
		if f.acceptanceCompletionObserver != nil {
			f.acceptanceCompletionObserver(observation)
		}
	}
}

func (f *acquisitionWiringCapture) acquisitionObserver(next func(collector.AcquisitionObservation)) func(collector.AcquisitionObservation) {
	return func(observation collector.AcquisitionObservation) {
		if next != nil {
			next(observation)
		}
		if f.acceptanceObserver != nil {
			f.acceptanceObserver(observation)
		}
		select {
		case f.observations <- observation:
		default:
		}
	}
}

// Historical wiring fixtures intentionally capture all ingestor offers. Only
// performance fixtures opt into the real pipeline's interest decision.
func (f *acquisitionWiringCapture) WantsSample(deviceID, pointID, tagID string) bool {
	if f.acceptanceInterest == nil {
		return true
	}
	return f.acceptanceInterest.wants(f.pipe, deviceID, pointID, tagID)
}

func (f *acquisitionWiringCapture) sampleSink(p *grouppipeline.Pipeline) datalinkruntime.SampleSink {
	f.pipe = p
	return f
}
func (f *acquisitionWiringCapture) AcceptSample(ctx context.Context, sample measurement.SampleEnvelope) error {
	err := f.pipe.AcceptSample(ctx, sample)
	if f.acceptanceSampleObserver != nil {
		f.acceptanceSampleObserver(sample)
	}
	select {
	case f.samples <- sample:
	case <-ctx.Done():
		return ctx.Err()
	}
	return err
}

func TestProductionScheduledAcquisitionBatchesContiguousPoints(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs, schema.ProtocolModbusTCP} {
		t.Run(string(protocol), func(t *testing.T) {
			plc, config := wordWireServer(t, protocol)
			plc.set(0, []uint16{11, 22, 33})
			f := newScheduledProductionFixture(t, protocol, config)
			services, capture := f.services, f.capture
			address := "D0"
			if protocol == schema.ProtocolModbusTCP {
				address = "40001"
			}
			rule, err := services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{DeviceID: "batch-device", StartAddress: address, Count: 3, DataType: schema.DataTypeInt16, NamingPrefix: "batch", Enabled: true})
			require.NoError(t, err)
			links, err := services.sourceRule.ListLinks(t.Context(), rule.ID)
			require.NoError(t, err)
			require.Len(t, links, 3)
			applyAcquisitionTags(t, services, rule)
			plc.mu.Lock()
			before := plc.calls
			plc.mu.Unlock()
			require.NoError(t, services.runtime.Start(t.Context()))
			var samples []measurement.SampleEnvelope
			for range 3 {
				select {
				case sample := <-capture.samples:
					samples = append(samples, sample)
				case <-time.After(5 * time.Second):
					t.Fatal("scheduled production ingestor did not reach the SampleSink")
				}
			}
			require.NoError(t, services.scheduler.Stop())
			plc.mu.Lock()
			calls := plc.calls - before
			plc.mu.Unlock()
			require.Equal(t, 1, calls, "one installed group scan must use one contiguous protocol read")
			requestIDs, pointIDs := map[string]bool{}, map[string]bool{}
			var planFingerprint string
			for range 3 {
				select {
				case observation := <-capture.observations:
					require.NotEmpty(t, observation.RequestID)
					require.NotEmpty(t, observation.PlanFingerprint)
					require.Equal(t, samples[0].AcquisitionID, observation.AcquisitionID)
					require.False(t, pointIDs[observation.PointID], "one final observation per Point")
					pointIDs[observation.PointID], requestIDs[observation.RequestID] = true, true
					if planFingerprint == "" {
						planFingerprint = observation.PlanFingerprint
					}
					require.Equal(t, planFingerprint, observation.PlanFingerprint)
				default:
					t.Fatal("production final acquisition observation missing")
				}
			}
			require.Len(t, requestIDs, 1)
			ids := map[string]bool{}
			values := map[any]bool{}
			for _, sample := range samples {
				require.Equal(t, schema.QualityGood, sample.Quality)
				require.Equal(t, samples[0].AcquisitionID, sample.AcquisitionID)
				require.Equal(t, samples[0].ReceivedAt, sample.ReceivedAt)
				require.Equal(t, sample.ReceivedAt, sample.ObservedAt)
				require.NotEmpty(t, sample.ConfigFingerprint)
				require.False(t, ids[sample.SampleID], "distinct points retain distinct stable sample identities")
				ids[sample.SampleID], values[sample.Value] = true, true
			}
			require.Equal(t, map[any]bool{int16(11): true, int16(22): true, int16(33): true}, values)
		})
	}
}

func applyAcquisitionTags(t *testing.T, services gatewayServices, rule *schema.SourceRule) {
	t.Helper()
	snapshots, err := services.sourceRule.ListCandidateSnapshots(t.Context(), rule.ID)
	require.NoError(t, err)
	var ids []string
	for _, snapshot := range snapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeTags {
			continue
		}
		var payload struct {
			Candidates []schema.SourceRuleTagCandidate `json:"candidates"`
		}
		require.NoError(t, json.Unmarshal([]byte(snapshot.Payload), &payload))
		for _, candidate := range payload.Candidates {
			ids = append(ids, candidate.ID)
		}
	}
	require.NotEmpty(t, ids)
	hydration, err := services.modbusShare.CheckHydration(t.Context())
	require.NoError(t, err)
	applied, err := services.sourceRule.ApplyTagCandidates(t.Context(), rule.ID, sourcerule.ApplyTagCandidatesRequest{WorkspaceID: hydration.WorkspaceID, ExpectedWorkspaceRevision: hydration.WorkspaceRevision, RevisionID: rule.RevisionID, CandidateIDs: ids})
	require.NoError(t, err)
	require.Len(t, applied.Results, len(ids))
	for _, result := range applied.Results {
		require.Equal(t, "applied", result.Status, result.Error)
	}
}

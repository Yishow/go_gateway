package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestProductionScheduledAcquisitionRecordingSubsetNeverNarrowsSource(t *testing.T) {
	words := map[int]uint16{}
	for address := range 21 {
		words[address] = uint16(address + 100)
	}
	peer, config := newAcquisitionWirePeer(t, schema.ProtocolMC3E, words)
	f := newScheduledProductionFixture(t, schema.ProtocolMC3E, config)
	rule, _ := f.createSource(t, "D0", 21, schema.DataTypeInt16, "source-all")
	applyAcquisitionTags(t, f.services, rule)
	links, err := f.services.sourceRule.ListLinks(t.Context(), rule.ID)
	require.NoError(t, err)
	members := make([]workspace.WriteGroupMember, 0, len(links))
	for _, link := range links {
		if link.Address != "D0" && link.Address != "D10" && link.Address != "D20" {
			continue
		}
		require.NotNil(t, link.TagID)
		members = append(members, workspace.WriteGroupMember{DeviceID: "batch-device", PointID: link.PointID, TagID: *link.TagID, TargetColumn: fmt.Sprintf("recorded_%d", len(members)), Required: true})
	}
	require.Len(t, members, 3)
	// The draft is real persisted recording membership, not acquisition input.
	// No external target is opened or applied, and no durability claim is made.
	destination, err := json.Marshal(map[string]string{"dsn": filepath.Join(t.TempDir(), "unopened.db")})
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('record-target','Target','sqlite',?,'target-v1','ready',1)`, string(destination))
	require.NoError(t, err)
	record, err := f.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	created, err := f.services.writeGroups.Create(t.Context(), workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: "target-v1", Group: &workspace.WriteGroup{
		WorkspaceID: record.ID, Name: "Only three recorded members", Members: members,
		Destination: workspace.WriteGroupDestination{ConnectorID: "record-target", ConnectorRevision: "target-v1", TableSchema: "main", TableName: "recorded", StorageStrategy: workspace.WriteGroupStorageStrategyCustom},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 15},
	}})
	require.NoError(t, err)
	require.Len(t, created.Group.Members, 3)
	// Explicit manual intent may read an installed source not selected by this
	// recording draft. It returns only that Point and never feeds the scheduled
	// channel; the following full 21-source rounds must remain unchanged.
	unrecordedPoint := ""
	for _, link := range links {
		if link.Address == "D1" {
			unrecordedPoint = link.PointID
		}
	}
	require.NotEmpty(t, unrecordedPoint)
	beforeManual := len(peer.log())
	manual := f.services.scheduler.PollNowContext(t.Context(), []string{unrecordedPoint})
	require.Len(t, manual, 1)
	require.Equal(t, unrecordedPoint, manual[0].PointID)
	require.Equal(t, schema.QualityGood, manual[0].Quality, manual[0].Error)
	require.Equal(t, int16(101), manual[0].Value)
	require.Equal(t, []acquisitionWireRequest{{Kind: "contiguous", Addresses: []int{1}}}, peer.log()[beforeManual:])
	require.Empty(t, f.services.scheduler.ValueChannel())
	require.Empty(t, f.capture.samples)
	require.Len(t, f.capture.observations, 1)
	require.Equal(t, unrecordedPoint, (<-f.capture.observations).PointID)
	baseline := len(peer.log())
	require.NoError(t, f.services.runtime.Start(t.Context()))
	samples := f.samples(t, 21)
	firstPlan, firstRequest := "", ""
	for range 21 {
		observation := <-f.capture.observations
		if firstPlan == "" {
			firstPlan, firstRequest = observation.PlanFingerprint, observation.RequestID
		}
		require.Equal(t, firstPlan, observation.PlanFingerprint)
		require.NotEmpty(t, firstRequest)
		require.Equal(t, firstRequest, observation.RequestID)
		require.Equal(t, samples[0].AcquisitionID, observation.AcquisitionID)
	}
	changed := *created.Group
	changed.Members = append([]workspace.WriteGroupMember(nil), changed.Members[:1]...)
	updated, err := f.services.writeGroups.Update(t.Context(), changed.ID, workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: created.WorkspaceRevision, ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "target-v1", Group: &changed,
	})
	require.NoError(t, err)
	require.Len(t, updated.Group.Members, 1)
	// SQL mutation can span several real ticks under race instrumentation.
	// Exclude two complete post-update rounds (running plus one pending),
	// then select a scan necessarily captured after the mutation.
	// Completion alone cannot prove that a previously captured scan is new.
	second := scheduledSamplesAfterRound(t, f, time.Now())
	require.NoError(t, f.services.scheduler.Stop())
	requests := peer.log()[baseline:]
	require.GreaterOrEqual(t, len(requests), 4)
	for _, request := range requests {
		require.Equal(t, "contiguous", request.Kind)
		require.Len(t, request.Addresses, 21, "source-rule membership remains all 21 points before and after recording selection changes")
	}
	seen := map[string]bool{}
	for _, sample := range samples {
		require.False(t, seen[sample.PointID])
		seen[sample.PointID] = true
		require.Equal(t, schema.QualityGood, sample.Quality)
		require.Equal(t, samples[0].AcquisitionID, sample.AcquisitionID)
		require.Equal(t, samples[0].ReceivedAt, sample.ReceivedAt)
	}
	require.Len(t, seen, 21)
	secondSeen := map[string]bool{}
	for _, sample := range second {
		require.True(t, seen[sample.PointID])
		require.False(t, secondSeen[sample.PointID])
		secondSeen[sample.PointID] = true
		require.Equal(t, schema.QualityGood, sample.Quality)
		require.NotEqual(t, samples[0].AcquisitionID, sample.AcquisitionID)
		require.Equal(t, second[0].AcquisitionID, sample.AcquisitionID)
	}
	requireScheduledScanObservation(t, f, second[0].AcquisitionID, firstPlan)
	stats, exists := f.services.scheduler.GetDeviceBreakerStats("batch-device")
	require.True(t, exists)
	window := collector.DefaultSchedulerConfig().BreakerConfig.WindowSize
	require.Equal(t, min(21, window), stats.TrackerStats.TotalRequests, "historical Point-result health fills its bounded rolling window, rather than counting the one wire exchange")
	require.Equal(t, min(21, window), stats.TrackerStats.SuccessCount)
}

// At the mutation cutoff this one device/group tuple can have one running scan
// and at most one captured pending scan. Excluding two complete distinct rounds
// proves the selected third scan was captured after the cutoff, even with due
// coalescing. Completion alone cannot prove when either earlier scan was captured.
func scheduledSamplesAfterRound(t *testing.T, f scheduledProductionFixture, completedAfter time.Time) []measurement.SampleEnvelope {
	t.Helper()
	const count = 21
	// The real SQLite ingestor is deliberately retained. Under race it can
	// deliver a FIFO backlog long after wire completion; this bounds only the
	// fixture observation wait, never a scan/request deadline or retry budget.
	deadline := time.NewTimer(30 * time.Second)
	defer deadline.Stop()
	barrierID := ""
	barrierIDs := make([]string, 0, 2)
	barrierPoints := make(map[string]bool, count)
	selected := make([]measurement.SampleEnvelope, 0, count)
	discarded := 0
	var lastCompletion time.Time
	for len(selected) < count {
		select {
		case sample := <-f.capture.samples:
			lastCompletion = sample.ReceivedAt
			if !sample.ReceivedAt.After(completedAfter) {
				discarded++
				continue
			}
			if len(barrierIDs) < 2 {
				if len(barrierPoints) == 0 {
					barrierID = sample.AcquisitionID
					require.NotContains(t, barrierIDs, barrierID)
				}
				require.Equal(t, barrierID, sample.AcquisitionID, "post-update barrier is one complete source scan")
				require.False(t, barrierPoints[sample.PointID])
				require.Equal(t, schema.QualityGood, sample.Quality)
				barrierPoints[sample.PointID] = true
				if len(barrierPoints) == count {
					barrierIDs = append(barrierIDs, barrierID)
					clear(barrierPoints)
				}
				continue
			}
			require.NotContains(t, barrierIDs, sample.AcquisitionID)
			if len(selected) > 0 {
				require.Equal(t, selected[0].AcquisitionID, sample.AcquisitionID, "after scan retains every source Point")
			}
			selected = append(selected, sample)
		case <-deadline.C:
			t.Fatalf("post-update source scan timeout: complete-barriers=%d/2 barrier=%d/%d selected=%d/%d discarded=%d cutoff=%s last-completion=%s rejected=%d runtime=%+v", len(barrierIDs), len(barrierPoints), count, len(selected), count, discarded, completedAfter, lastCompletion, f.services.runtime.RejectedObservationCount(), f.services.runtime.Snapshot())
		}
	}
	t.Logf("post-update source barrier: discarded=%d, elapsed=%s, final gateway-completion-to-sink delay=%s", discarded, time.Since(completedAfter), time.Since(selected[0].ReceivedAt))
	return selected
}

func requireScheduledScanObservation(t *testing.T, f scheduledProductionFixture, acquisitionID, planFingerprint string) {
	t.Helper()
	const count = 21
	deadline := time.NewTimer(7 * time.Second)
	defer deadline.Stop()
	requestID := ""
	for seen := 0; seen < count; {
		select {
		case observation := <-f.capture.observations:
			if observation.AcquisitionID != acquisitionID {
				continue
			}
			if requestID == "" {
				requestID = observation.RequestID
			}
			require.NotEmpty(t, requestID)
			require.Equal(t, requestID, observation.RequestID, "one final request supplies this entire source scan")
			require.Equal(t, planFingerprint, observation.PlanFingerprint, "record selection is not acquisition configuration identity")
			seen++
		case <-deadline.C:
			t.Fatalf("after-scan observation timeout: %d/%d", seen, count)
		}
	}
}

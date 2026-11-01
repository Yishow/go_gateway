package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

type acceptanceScheduledDevice struct {
	id, sourceHead string
	peer           *acquisitionWirePeer
	metrics        *acquisitionAcceptanceWireMetrics
	want           map[string]int16
}

func newAcceptanceScheduledLoad(t *testing.T, devices, points int, rtt time.Duration, captureAll ...bool) (scheduledProductionFixture, []acceptanceScheduledDevice, string) {
	t.Helper()
	load := make([]acceptanceScheduledDevice, devices)
	var fixture scheduledProductionFixture
	for i := range devices {
		memory := map[int]uint16{}
		for address := range max(points, 21) {
			memory[address] = uint16(address + 1)
		}
		peer, config := newAcquisitionWirePeer(t, schema.ProtocolMC3E, memory)
		config = withAcquisitionDeclaration(config, declaredAcquisitionProfile(schema.ProtocolMC3E))
		id := "batch-device"
		if i == 0 {
			fixture = newScheduledProductionFixture(t, schema.ProtocolMC3E, config)
		} else {
			id = fmt.Sprintf("load-device-%03d", i)
			var fields map[string]any
			require.NoError(t, json.Unmarshal([]byte(config), &fields))
			_, err := fixture.services.device.Create(t.Context(), device.CreateDeviceRequest{ID: id, Name: id, Protocol: schema.ProtocolMC3E, ConnectionConfig: fields})
			require.NoError(t, err)
			require.NoError(t, fixture.services.device.ProbeAndActivate(t.Context(), id))
			_, err = fixture.services.workspace.AttachDevice(t.Context(), id)
			require.NoError(t, err)

		}
		load[i] = acceptanceScheduledDevice{id: id, peer: peer, want: map[string]int16{}}
	}
	// Prepare the canonical persisted graph without invoking hot workspace
	// re-projection for every Point. Measurement still uses the original runtime
	// Start/full bootstrap, repositories, profiles, and installation CAS.
	revisions := modbusshare.NewSQLWorkspaceRevisionStore(fixture.db)
	hydrateStartupShare(t.Context(), fixture.services, revisions, modbusshare.NewReconciler(fixture.services.modbusShare, revisions))
	seed := fixture.services
	groupsRepo := pollinggroup.NewSQLRepository(fixture.db)
	seed.point = point.NewService(point.NewSQLRepository(fixture.db), groupsRepo)
	seed.pollingGroup = pollinggroup.NewService(groupsRepo)
	seed.sourceRule = sourcerule.NewService(sourcerule.NewSQLRepository(fixture.db), seed.device, seed.point, nil)
	seed.sourceRule.SetTagMappingServices(seed.tag, seed.mapping)
	seed.sourceRule.SetCandidateScope(seed.workspace, seed.modbusShare)
	for i := range load {
		d := &load[i]
		id := d.id
		starts, counts := []int{0}, []int{points}
		if points == 3 {
			starts, counts = []int{0, 10, 20}, []int{1, 1, 1}
		}
		for j, start := range starts {
			rule, err := seed.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{DeviceID: id, StartAddress: fmt.Sprintf("D%d", start), Count: counts[j], DataType: schema.DataTypeInt16, NamingPrefix: fmt.Sprintf("load-%d-%d", i, j), Enabled: true})
			require.NoError(t, err)
			applyAcquisitionTags(t, seed, rule)
			links, err := seed.sourceRule.ListLinks(t.Context(), rule.ID)
			require.NoError(t, err)
			for _, link := range links {
				var address int
				_, err := fmt.Sscanf(link.Address, "D%d", &address)
				require.NoError(t, err)
				d.want[link.PointID] = int16(address + 1)
			}
		}
	}
	groups, err := seed.pollingGroup.List(t.Context())
	require.NoError(t, err)
	// Explicit fixed cadence, identical across all load dimensions. Overrun and
	// due coalescing are observations, not hidden by waiting for each prior scan.
	cadence := 1000
	for _, group := range groups {
		_, err = seed.pollingGroup.Update(t.Context(), group.ID, pollinggroup.UpdateRequest{IntervalMs: &cadence})
		require.NoError(t, err)
	}
	projection, err := fixture.services.workspace.ReadRuntimeProjection(t.Context())
	require.NoError(t, err)
	for i := range load {
		graph, ok := projection.AcquisitionSourceGraph(load[i].id)
		require.True(t, ok)
		load[i].sourceHead = graph.IdentityRevision()
	}
	// Probe/setup may connect and read. Close those connections BEFORE enabling
	// counters/delay; the first scheduled scan then includes a genuine fresh dial.
	require.NoError(t, fixture.manager.CloseAll())
	for i := range load {
		load[i].metrics = load[i].peer.acceptanceMetrics(rtt)
	}
	// Only the separate bootstrap wiring proof opts into historical capture-all.
	// Every performance pilot/matrix uses genuine pipeline interest by default.
	if len(captureAll) == 0 || !captureAll[0] {
		fixture.capture.acceptanceInterest = startAcceptanceInterestCapture(t, fixture, load)
	}
	return fixture, load, projection.WorkspaceID
}

//go:build f_write_group_fixture

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestFixturePollReturnsSafeErrorWhenCaptureQueueIsFull(t *testing.T) {
	env := newOutageEnv(t)
	fixture := env.services.fixture.(*groupFixture)
	activeAt := time.Now().UTC()
	fixture.clock.Set(activeAt)
	env.services.writeGroups.WithClock(fixture.clock.Now)
	group := env.createAndApply(t, "A")
	fixture.deps = &env.services
	fixture.pipe = env.services.groupPipe
	require.NoError(t, fixture.pause(t.Context()))
	fixture.capture.mu.Lock()
	fixture.capture.capacity = 0
	fixture.capture.mu.Unlock()

	_, err := env.db.ExecContext(t.Context(), `UPDATE devices SET description = '', status = 'active' WHERE id = ?`, "device-1")
	require.NoError(t, err)
	env.services.scheduler.AddDevice(&schema.Device{ID: "device-1", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{}`})
	for _, member := range group.Members[:1] {
		_, err := env.db.ExecContext(t.Context(), `UPDATE points SET description = '' WHERE id = ?`, member.PointID)
		require.NoError(t, err)
		point, pointErr := env.services.point.GetByID(t.Context(), member.PointID)
		require.NoError(t, pointErr)
		env.services.scheduler.AddPoint(point)
	}
	resolved, err := env.services.writeGroups.ResolveAppliedAt(t.Context(), group.ID, activeAt.Add(time.Minute))
	require.NoError(t, err)
	base := resolved.EffectiveAt.Truncate(10 * time.Second).Add(10 * time.Second)
	fixture.clock.Set(base.Add(time.Second))
	require.NoError(t, fixture.pipe.Reconcile(t.Context()))
	projection, err := env.services.workspace.RuntimeProjection(t.Context())
	require.NoError(t, err)
	require.NoError(t, env.services.runtime.ApplyWorkspaceProjection(t.Context(), projection))
	env.services.scheduler.AddDevice(&schema.Device{ID: "device-1", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{}`})
	point, err := env.services.point.GetByID(t.Context(), group.Members[0].PointID)
	require.NoError(t, err)
	env.services.scheduler.AddPoint(point)

	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/poll", strings.NewReader(`{"group_id":"`+group.ID+`","point_ids":["`+group.Members[0].PointID+`"]}`))
	recorder := httptest.NewRecorder()
	fixture.handlePoll(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), `"error":"manual poll capture failed"`)
	require.NotContains(t, recorder.Body.String(), "queue full")
	require.Empty(t, fixture.capture.samples)
}

func TestFixtureFaultsEndpointRespondsWhileClosureTickIsHeld(t *testing.T) {
	env := newOutageEnv(t)
	fixture := env.services.fixture.(*groupFixture)
	activeAt := time.Now().UTC()
	fixture.clock.Set(activeAt)
	env.services.writeGroups.WithClock(fixture.clock.Now)
	group := env.createAndApply(t, "A")
	fixture.deps = &env.services
	fixture.pipe = env.services.groupPipe
	fixture.faults = newFixtureFaultController(env.db, &env.services, fixture.clock.Now)
	require.NoError(t, fixture.faults.install())
	t.Cleanup(func() { require.NoError(t, fixture.faults.close(context.Background())) })

	resolved, err := env.services.writeGroups.ResolveAppliedAt(t.Context(), group.ID, activeAt.Add(time.Minute))
	require.NoError(t, err)
	base := resolved.EffectiveAt.Truncate(10 * time.Second).Add(10 * time.Second)
	fixture.clock.Set(base.Add(time.Second))
	require.NoError(t, fixture.pipe.Reconcile(t.Context()))
	_, err = fixture.faults.configure(t.Context(), fixtureFaultRequest{
		Kind: fixtureFaultClosureHold, GroupID: group.ID, Enabled: true,
	})
	require.NoError(t, err)

	observed := base.Add(5 * time.Second)
	fixture.clock.Set(observed.Add(time.Second))
	require.NoError(t, fixture.pipe.AcceptSample(t.Context(), env.envelope(group, 0, "faults-http-temperature", observed, 21.5)))
	require.NoError(t, fixture.pipe.AcceptSample(t.Context(), env.envelope(group, 1, "faults-http-pressure", observed, int64(1013))))
	fixture.clock.Set(base.Add(11 * time.Second))
	tickDone := make(chan struct{})
	go func() {
		fixture.pipe.TickAll(t.Context())
		close(tickDone)
	}()
	require.Eventually(t, func() bool {
		for _, state := range fixture.faults.snapshot() {
			if state.GroupID == group.ID {
				return state.Holding
			}
		}
		return false
	}, time.Second, time.Millisecond)

	fixture.port = 0
	require.NoError(t, fixture.startHTTP(t.Context()))
	t.Cleanup(func() { fixture.stopHTTP(context.Background()) })
	address := fixture.http.listener.Addr().String()
	client := &http.Client{Timeout: time.Second}
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://"+address+"/faults", http.NoBody)
	require.NoError(t, err)
	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)
	var body struct {
		Faults []fixtureFaultState `json:"faults"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Len(t, body.Faults, 1)
	require.Equal(t, group.ID, body.Faults[0].GroupID)
	require.True(t, body.Faults[0].Holding)

	// The held SQLite transaction owns the single configuration connection, so
	// release the process-local gate before asking the controller to drop SQL.
	releaseFixtureClosureGate(group.ID)
	select {
	case <-tickDone:
	case <-time.After(time.Second):
		t.Fatal("closure tick did not finish after release")
	}
	_, err = fixture.faults.configure(t.Context(), fixtureFaultRequest{
		Kind: fixtureFaultClosureHold, GroupID: group.ID, Enabled: false,
	})
	require.NoError(t, err)
}

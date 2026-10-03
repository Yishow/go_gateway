//go:build f_write_group_fixture

package main

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	timepkg "time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestFixtureClosureHoldBlocksUntilRelease(t *testing.T) {
	gate := &fixtureClosureGate{groupID: "group-closure", release: make(chan struct{})}
	installFixtureClosureGate(gate)
	t.Cleanup(func() { removeFixtureClosureGate(gate.groupID, gate); gate.releaseGate() })

	done := make(chan error, 1)
	go func() {
		_, err := fixtureClosureScalar(nil, []driver.Value{"group-closure"})
		done <- err
	}()
	require.Eventually(t, gate.holding.Load, timepkg.Second, timepkg.Millisecond)
	select {
	case <-done:
		t.Fatal("closure hold returned before release")
	default:
	}
	gate.releaseGate()
	require.NoError(t, <-done)
	require.False(t, gate.holding.Load())
}

func TestFixtureClosureHoldTriggerBlocksCheckpointTransaction(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	_, err = db.ExecContext(t.Context(), `CREATE TABLE wg_delivery_checkpoints (group_id TEXT, next_close TEXT)`)
	require.NoError(t, err)
	gate := &fixtureClosureGate{groupID: "group-checkpoint", release: make(chan struct{})}
	installFixtureClosureGate(gate)
	t.Cleanup(func() { removeFixtureClosureGate(gate.groupID, gate); gate.releaseGate() })
	controller := &fixtureFaultController{db: db}
	trigger, err := controller.createClosureHoldTrigger(t.Context(), gate.groupID)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, controller.dropTrigger(context.Background(), trigger)) })
	done := make(chan error, 1)
	go func() {
		_, execErr := db.ExecContext(context.Background(), `INSERT INTO wg_delivery_checkpoints (group_id, next_close) VALUES ('group-checkpoint', '2026-01-01T00:00:00Z')`)
		done <- execErr
	}()
	require.Eventually(t, gate.holding.Load, timepkg.Second, timepkg.Millisecond)
	gate.releaseGate()
	require.NoError(t, <-done)
}

func TestFixtureReceiptFailureTriggerIsGroupScoped(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	_, err = db.ExecContext(t.Context(), `
		CREATE TABLE wg_delivery_outbox (effect_key TEXT PRIMARY KEY, group_id TEXT NOT NULL);
		CREATE TABLE wg_delivery_receipts (effect_key TEXT PRIMARY KEY, payload_digest TEXT NOT NULL, committed_at TEXT NOT NULL);`)
	require.NoError(t, err)
	controller := &fixtureFaultController{db: db}
	trigger, err := controller.createReceiptFailureTrigger(t.Context(), "group-a")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, controller.dropTrigger(context.Background(), trigger)) })
	_, err = db.ExecContext(t.Context(), `INSERT INTO wg_delivery_outbox (effect_key, group_id) VALUES ('effect-a', 'group-a'), ('effect-b', 'group-b')`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO wg_delivery_receipts VALUES ('effect-a', 'digest', '2026-01-01T00:00:00Z')`)
	require.Error(t, err, "the owned group receipt is intentionally rejected")
	_, err = db.ExecContext(t.Context(), `INSERT INTO wg_delivery_receipts VALUES ('effect-b', 'digest', '2026-01-01T00:00:00Z')`)
	require.NoError(t, err, "a different group's receipt is unaffected")
}

func TestFixtureFaultRejectsUnknownCanonicalGroup(t *testing.T) {
	env := newOutageEnv(t)
	fixture := env.services.fixture.(*groupFixture)
	fixture.faults = newFixtureFaultController(env.db, &env.services, fixture.clock.Now)
	_, err := fixture.faults.configure(t.Context(), fixtureFaultRequest{
		Kind: "local_receipt_failure", GroupID: "foreign-group", Enabled: true,
	})
	require.ErrorIs(t, err, errFixtureFaultGroupUnknown)
}

func TestFixtureTargetFaultUsesAppliedConnectorWhenDraftChanges(t *testing.T) {
	env := newOutageEnv(t)
	fixture := env.services.fixture.(*groupFixture)
	activeAt := timepkg.Now().UTC()
	fixture.clock.Set(activeAt)
	env.services.writeGroups.WithClock(fixture.clock.Now)
	group := env.createAndApply(t, "A")
	current, err := env.services.writeGroups.Get(t.Context(), group.ID)
	require.NoError(t, err)
	draft := *current.Group
	draft.Members = append([]workspace.WriteGroupMember(nil), current.Group.Members...)
	draft.Destination.ConnectorID = "connector-B"
	workspaceState, err := env.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	updated, err := env.services.writeGroups.Update(t.Context(), group.ID, workspace.WriteGroupMutation{
		WorkspaceID:               workspaceState.ID,
		ExpectedWorkspaceRevision: workspaceState.DatabaseSetupRevision,
		ExpectedGroupRevision:     current.Group.Revision,
		ExpectedConnectorRevision: "connector-1",
		Group:                     &draft,
	})
	require.NoError(t, err)
	require.Equal(t, "connector-B", updated.Group.Destination.ConnectorID)
	require.Equal(t, group.AppliedRevision, updated.Group.AppliedRevision)

	fixture.clock.Set(activeAt.Add(timepkg.Minute))
	fixture.faults = newFixtureFaultController(env.db, &env.services, fixture.clock.Now)
	require.NoError(t, fixture.faults.install())
	t.Cleanup(func() { require.NoError(t, fixture.faults.close(context.Background())) })
	state, err := fixture.faults.configure(t.Context(), fixtureFaultRequest{
		Kind: string(dbtarget.FixtureDestinationFaultTargetCommitResponseLost), GroupID: group.ID, Enabled: true,
	})
	require.NoError(t, err)
	require.Equal(t, "connector-A", state.ConnectorID, "fault must follow the active applied destination")
}

func TestFixtureTargetFaultRejectsWithoutActiveSnapshot(t *testing.T) {
	env := newOutageEnv(t)
	fixture := env.services.fixture.(*groupFixture)
	activeAt := timepkg.Now().UTC()
	fixture.clock.Set(activeAt)
	env.services.writeGroups.WithClock(fixture.clock.Now)
	group := env.createAndApply(t, "A")
	fixture.faults = newFixtureFaultController(env.db, &env.services, fixture.clock.Now)
	require.NoError(t, fixture.faults.install())
	t.Cleanup(func() { require.NoError(t, fixture.faults.close(context.Background())) })
	_, err := fixture.faults.configure(t.Context(), fixtureFaultRequest{
		Kind: string(dbtarget.FixtureDestinationFaultTargetCommitResponseLost), GroupID: group.ID, Enabled: true,
	})
	require.ErrorIs(t, err, errFixtureFaultGroupUnknown)
}

func TestFixtureFaultHTTPUsesStablePact(t *testing.T) {
	env := newOutageEnv(t)
	fixture := env.services.fixture.(*groupFixture)
	activeAt := timepkg.Now().UTC()
	fixture.clock.Set(activeAt)
	env.services.writeGroups.WithClock(fixture.clock.Now)
	group := env.createAndApply(t, "A")
	fixture.clock.Set(activeAt.Add(timepkg.Minute))
	fixture.faults = newFixtureFaultController(env.db, &env.services, fixture.clock.Now)
	require.NoError(t, fixture.faults.install())
	t.Cleanup(func() { require.NoError(t, fixture.faults.close(context.Background())) })

	body := []byte(`{"kind":"closure_hold","group_id":"` + group.ID + `","enabled":true}`)
	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/fault", bytes.NewReader(body))
	response := httptest.NewRecorder()
	fixture.handleFault(response, request)
	require.Equal(t, http.StatusOK, response.Code)
	var faultEnvelope struct {
		Fault fixtureFaultState `json:"fault"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&faultEnvelope))
	require.Equal(t, "closure_hold", faultEnvelope.Fault.Kind)
	require.Equal(t, group.ID, faultEnvelope.Fault.GroupID)
	require.Equal(t, "connector-A", faultEnvelope.Fault.ConnectorID)
	require.True(t, faultEnvelope.Fault.Enabled)

	stateResponse := httptest.NewRecorder()
	fixture.handleState(stateResponse, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/state", http.NoBody))
	require.Equal(t, http.StatusOK, stateResponse.Code)
	var state fixtureStateResponse
	require.NoError(t, json.NewDecoder(stateResponse.Body).Decode(&state))
	require.Len(t, state.Faults, 1)
	require.Equal(t, faultEnvelope.Fault, state.Faults[0])

	disableBody := []byte(`{"kind":"closure_hold","group_id":"` + group.ID + `","enabled":false}`)
	disableRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/fault", bytes.NewReader(disableBody))
	disableResponse := httptest.NewRecorder()
	fixture.handleFault(disableResponse, disableRequest)
	require.Equal(t, http.StatusOK, disableResponse.Code)
	var disabled struct {
		Fault fixtureFaultState `json:"fault"`
	}
	require.NoError(t, json.NewDecoder(disableResponse.Body).Decode(&disabled))
	require.False(t, disabled.Fault.Enabled)
}

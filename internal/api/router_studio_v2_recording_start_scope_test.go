package api

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/recordingplan"
	datalinksettings "go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// scopeRig is one workspace with two devices behind the real group, schema,
// readiness and recording-start services. Device A owns a mean long group and
// a new wide group; device B owns one long group in a chosen state. Only the
// device activation call is simulated, and it records what it was asked for.
type scopeRig struct {
	recordingStartFixture
	longID, wideID, otherID string
}

var (
	scopeLongSources = []meanSource{{name: "half", tagType: "int16", role: "gauge"}, {name: "flt", tagType: "float32", role: "gauge"}}
	scopeWideSources = []meanSource{{name: "wone", tagType: "int16", role: "gauge"}, {name: "wtwo", tagType: "float32", role: "gauge"}}
)

// newScopeRig saves A's two profile groups and B's group. prepareA confirms
// A's explicit schema preparation; otherPrepared confirms and applies B, which
// simulates an already applied device. No collector runs in this fixture, so it
// is not evidence of real collecting or of a real device.
func newScopeRig(t *testing.T, prepareA, otherPrepared bool) scopeRig {
	t.Helper()
	managed := newManagedDestinationFixture(t, "sqlite")
	seedMeanSources(t, managed, scopeLongSources)
	seedMeanSources(t, managed, scopeWideSources)
	seedOtherDevice(t, managed)

	rig := scopeRig{recordingStartFixture: newRecordingStartServiceFixture(t, managed)}
	rig.longID = rig.saveProfile(t, meanLongRequest(t, managed, scopeLongSources))
	wide := wideProfileRequest(t, managed, scopeWideSources)
	renameCircuit(wide, "wide-main")
	rig.wideID = rig.saveProfile(t, wide)
	other := meanLongRequest(t, managed, []meanSource{{name: "bsrc", tagType: "float32", role: "gauge"}})
	renameCircuit(other, "other-main")
	for _, member := range other["group"].(map[string]any)["members"].([]map[string]any) {
		member["device_id"] = "device-B"
	}
	rig.otherID = rig.saveProfile(t, other)

	if prepareA {
		rig.prepare(t, rig.longID)
		rig.prepare(t, rig.wideID)
	}
	if otherPrepared {
		rig.prepare(t, rig.otherID)
		applied := performJSONRequest(t, managed.router, http.MethodPost, writeGroupsPath+"/"+rig.otherID+"/apply", scopeBody(t, managed, rig.otherID))
		require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	}
	return rig
}

func seedOtherDevice(t *testing.T, f managedGroupFixture) {
	t.Helper()
	ctx := t.Context()
	for _, statement := range []string{
		`INSERT INTO devices (id,name,protocol,status,connection_config,readiness_status) VALUES ('device-B','PLC B','modbus_tcp','draft','{}','{}')`,
		`INSERT INTO points (id,device_id,name,address,function,data_type,mode,enabled) VALUES ('point-bsrc','device-B','bsrc','40111','FC03','float32','read',1)`,
		`INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels) VALUES ('tag-bsrc','bsrc','bsrc','bsrc','float32','active','{}')`,
		`INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled) VALUES ('mapping-bsrc','point-bsrc','tag-bsrc','[]','active',1)`,
	} {
		_, err := f.db.ExecContext(ctx, statement)
		require.NoError(t, err)
	}
	_, err := f.workspace.AttachDevice(ctx, "device-B")
	require.NoError(t, err)
}

func renameCircuit(request map[string]any, circuit string) {
	group := request["group"].(map[string]any)
	group["name"] = circuit + " instantaneous"
	for _, member := range group["members"].([]map[string]any) {
		member["circuit_id"], member["circuit_name"] = circuit, "Circuit "+circuit
	}
}

func (r scopeRig) saveProfile(t *testing.T, request map[string]any) string {
	t.Helper()
	response := performJSONRequest(t, r.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, response))
	return saved["id"].(string)
}

func (r scopeRig) prepare(t *testing.T, groupID string) {
	t.Helper()
	// A group on a shared managed table may need no statements; the confirmation is still explicit.
	token := previewLong(t, r.managedGroupFixture, groupID)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, r.managedGroupFixture, groupID, token).Status)
}

// request freezes the current workspace and group revisions of device A for
// the named groups, as the reviewed primary start would.
func (r scopeRig) request(t *testing.T, requestID string, groupIDs ...string) workspace.RecordingStartRequest {
	t.Helper()
	current, err := r.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	request := workspace.RecordingStartRequest{
		RequestID: requestID, WorkspaceID: current.ID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		DeviceIDs: []string{"device-A"}, Groups: []workspace.RecordingStartGroupIntent{},
	}
	for _, id := range groupIDs {
		saved, err := r.groups.Get(t.Context(), id)
		require.NoError(t, err)
		request.Groups = append(request.Groups, workspace.RecordingStartGroupIntent{
			GroupID: id, ExpectedGroupRevision: saved.Group.Revision, ExpectedConnectorRevision: "connector-1",
		})
	}
	return request
}

// groupFacts is everything about a group that a start of another device must not change.
type groupFacts struct {
	Revision, AppliedRevision, Status, SchemaRevision, SchemaDigest string
}

func (r scopeRig) facts(t *testing.T, groupID string) groupFacts {
	t.Helper()
	saved, err := r.groups.Get(t.Context(), groupID)
	require.NoError(t, err)
	return groupFacts{
		Revision: saved.Group.Revision, AppliedRevision: saved.Group.AppliedRevision, Status: string(saved.Group.Status),
		SchemaRevision: saved.Group.Destination.SchemaRevision, SchemaDigest: saved.Group.Destination.SchemaDigest,
	}
}

func (r scopeRig) count(t *testing.T, query string, args ...any) int {
	t.Helper()
	var n int
	require.NoError(t, r.db.QueryRowContext(t.Context(), query, args...).Scan(&n))
	return n
}

func (r scopeRig) targetTables(t *testing.T) int {
	t.Helper()
	target, err := sql.Open("sqlite", r.file)
	require.NoError(t, err)
	defer target.Close()
	var n int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT count(*) FROM sqlite_master WHERE type='table' AND name LIKE 'gw_%'`).Scan(&n))
	return n
}

func TestUnifiedRecordingStartScope(t *testing.T) {
	t.Run("A long and wide start only their frozen scope while B is incomplete", func(t *testing.T) {
		rig := newScopeRig(t, true, false)
		otherBefore := rig.facts(t, rig.otherID)
		require.Empty(t, otherBefore.AppliedRevision, "B stays unprepared and unapplied")
		tablesBefore := rig.targetTables(t)

		operation, err := rig.start.Start(t.Context(), rig.request(t, "start-A-long-wide", rig.longID, rig.wideID))
		require.NoError(t, err)
		require.Equal(t, recordingplan.SchemaOperationSucceeded, operation.Status, "%+v", operation)
		require.Equal(t, []string{"device-A"}, operation.DeviceIDs)
		require.Len(t, operation.Groups, 2)
		for _, progress := range operation.Groups {
			require.True(t, progress.Applied, "%+v", progress)
			require.NotEqual(t, rig.otherID, progress.GroupID)
		}
		scopes := rig.activation.Scopes()
		require.Len(t, scopes, 1)
		require.Equal(t, []string{"device-A"}, scopes[0].DeviceIDs)
		require.ElementsMatch(t, []string{rig.longID, rig.wideID}, scopes[0].GroupIDs)
		require.Equal(t, otherBefore, rig.facts(t, rig.otherID), "the incomplete device is neither prepared, applied nor rewritten")
		require.Equal(t, tablesBefore, rig.targetTables(t), "start itself creates no structure")
	})

	t.Run("A wide start leaves an applied B untouched", func(t *testing.T) {
		rig := newScopeRig(t, true, true)
		otherBefore := rig.facts(t, rig.otherID)
		require.NotEmpty(t, otherBefore.AppliedRevision, "B is applied before A starts")

		operation, err := rig.start.Start(t.Context(), rig.request(t, "start-A-wide", rig.wideID))
		require.NoError(t, err)
		require.Equal(t, recordingplan.SchemaOperationSucceeded, operation.Status, "%+v", operation)
		scopes := rig.activation.Scopes()
		require.Len(t, scopes, 1)
		require.Equal(t, recordingStartActivationScope{DeviceIDs: []string{"device-A"}, GroupIDs: []string{rig.wideID}}, scopes[0])
		require.Equal(t, otherBefore, rig.facts(t, rig.otherID), "an applied device keeps its revision, applied revision and status")
		require.Equal(t, rig.facts(t, rig.longID).AppliedRevision, "", "A's other group was not part of the reviewed scope")
	})

	for _, layout := range []string{"long", "wide"} {
		t.Run("missing confirmation for A "+layout+" directs preparation without DDL", func(t *testing.T) {
			rig := newScopeRig(t, false, true)
			groupID := map[string]string{"long": rig.longID, "wide": rig.wideID}[layout]
			otherBefore := rig.facts(t, rig.otherID)
			tablesBefore := rig.targetTables(t)

			operation, err := rig.start.Start(t.Context(), rig.request(t, "start-unprepared-"+layout, groupID))
			require.NoError(t, err)
			require.Equal(t, recordingplan.SchemaOperationFailed, operation.Status)
			require.Equal(t, "preparation_required", operation.Reason)
			require.Equal(t, "prepare_schema", operation.NextAction)
			require.Zero(t, rig.activation.Calls())
			require.Empty(t, rig.facts(t, groupID).AppliedRevision)
			require.Equal(t, tablesBefore, rig.targetTables(t), "start never creates the structure")
			require.Zero(t, rig.count(t, `SELECT count(*) FROM recording_plans`), "no recording plan is claimed without confirmation")
			require.Equal(t, otherBefore, rig.facts(t, rig.otherID))
		})
	}

	t.Run("share only uses the zero group branch and keeps the Share barrier", func(t *testing.T) {
		for _, ready := range []bool{false, true} {
			name := "share-not-ready"
			if ready {
				name = "share-ready"
			}
			t.Run(name, func(t *testing.T) {
				rig := newScopeRig(t, true, true)
				share := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
				// A durable settings seam keeps the enabled output from binding a listener before hydration.
				require.NoError(t, share.SetSettingsRepository(datalinksettings.NewMemoryRepository()))
				settings := modbusshare.DefaultSettings()
				settings.Enabled, settings.Port, settings.SettingsRevision = true, 5020, "settings-1"
				require.NoError(t, share.ApplySettings(t.Context(), settings))
				hydration := modbusshare.HydrationState{
					State: modbusshare.HydrationStateFailed, WorkspaceID: rig.record.ID, WorkspaceRevision: "persisted-share-revision",
					SettingsRevision: share.Settings().SettingsRevision, ReadinessToken: "share-proof",
				}
				if ready {
					hydration.State, hydration.Readiness = modbusshare.HydrationStateReady, true
				}
				share.SetHydrationState(hydration)
				barrierCalls := 0
				services := &DatalinkServices{
					Workspace: rig.workspace, ModbusShare: share,
					RecordingStartShareBarrier: func(context.Context, workspace.RecordingStartRequest) error {
						barrierCalls++
						return nil
					},
					ShareRestore: func(context.Context, handlers.ActivateWorkspaceRequest) error { return nil },
				}
				rig.start.WithActivationBarrier(recordingStartActivationBarrier(services), recordingStartActivationResumeBarrier(services))
				request := rig.request(t, "start-share-only-"+name)
				request.WorkspaceRevision, request.SettingsRevision, request.ReadinessToken = hydration.WorkspaceRevision, hydration.SettingsRevision, hydration.ReadinessToken
				otherBefore := rig.facts(t, rig.otherID)
				tablesBefore := rig.targetTables(t)

				operation, err := rig.start.Start(t.Context(), request)
				require.NoError(t, err)
				require.Empty(t, operation.Groups, "Share-only selects no database group")
				require.Equal(t, tablesBefore, rig.targetTables(t), "no database structure or writer is prepared")
				require.Equal(t, otherBefore, rig.facts(t, rig.otherID))
				require.Zero(t, rig.count(t, `SELECT count(*) FROM recording_plans`))
				if ready {
					require.Equal(t, recordingplan.SchemaOperationSucceeded, operation.Status, "%+v", operation)
					require.Equal(t, []recordingStartActivationScope{{DeviceIDs: []string{"device-A"}, GroupIDs: []string{}}}, rig.activation.Scopes())
					require.Positive(t, barrierCalls, "the existing selected-projection barrier still runs before readiness and activation")
					return
				}
				require.Equal(t, recordingplan.SchemaOperationFailed, operation.Status)
				require.Equal(t, "share_not_ready", operation.Reason)
				require.Equal(t, "refresh_share", operation.NextAction)
				require.Zero(t, rig.activation.Calls())
			})
		}
	})

	t.Run("a discarded success response replays the recorded operation after restart without a second activation", func(t *testing.T) {
		rig := newScopeRig(t, true, true)
		otherBefore := rig.facts(t, rig.otherID)
		request := rig.request(t, "start-A-discarded-response", rig.longID, rig.wideID)

		// The service completes A; the client then loses or discards that success response.
		first, err := rig.start.Start(t.Context(), request)
		require.NoError(t, err)
		require.Equal(t, recordingplan.SchemaOperationSucceeded, first.Status, "%+v", first)
		firstScopes := rig.activation.Scopes()
		require.Len(t, firstScopes, 1)
		require.Equal(t, []string{"device-A"}, firstScopes[0].DeviceIDs)
		require.ElementsMatch(t, []string{rig.longID, rig.wideID}, firstScopes[0].GroupIDs)
		applied := map[string]string{rig.longID: rig.facts(t, rig.longID).AppliedRevision, rig.wideID: rig.facts(t, rig.wideID).AppliedRevision}
		for id, revision := range applied {
			require.NotEmpty(t, revision, id)
		}

		// A restarted service reads the same SQL ledger and receives the identical request again.
		restarted := newRecordingStartServiceFixture(t, rig.managedGroupFixture)
		replayed, err := restarted.start.Start(t.Context(), request)
		require.NoError(t, err)
		require.Equal(t, recordingplan.SchemaOperationSucceeded, replayed.Status, "%+v", replayed)
		require.Equal(t, first.OperationID, replayed.OperationID, "the same request returns the same operation")
		require.Equal(t, first.Groups, replayed.Groups, "the recorded per-group progress and applied revisions are returned")
		require.Zero(t, restarted.activation.Calls(), "the replay activates nothing")
		require.Equal(t, firstScopes, rig.activation.Scopes(), "the original activator was asked for A and its groups exactly once")
		require.Equal(t, applied[rig.longID], rig.facts(t, rig.longID).AppliedRevision, "the applied group is not applied again")
		require.Equal(t, applied[rig.wideID], rig.facts(t, rig.wideID).AppliedRevision)
		require.Equal(t, 1, rig.count(t, `SELECT count(*) FROM managed_schema_operations WHERE action='recording_start'`), "one operation, no second writer claim")
		require.Equal(t, otherBefore, rig.facts(t, rig.otherID))
	})
}

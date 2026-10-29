package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
)

type firstStartAPIRig struct {
	*successorRuntimeRig
	targetNamespace string
	schemaToken     recordingplan.SchemaPreviewToken
}

// Create and remove an actual draft through public routes. No version/history
// is deleted to manufacture NoHistory; G1 has never been applied or reconciled.
func newFirstStartAPIRig(t *testing.T, kind string, sourceLong, targetLong, counter bool, mode workspace.WriteGroupStorageMode, open ...func(string) *sql.DB) *firstStartAPIRig {
	t.Helper()
	var f managedGroupFixture
	var target *sql.DB
	var namespace string
	switch kind {
	case "postgres":
		f, namespace = newPostgresSuccessorDestination(t, mode == "timescale")
		target = f.target
	case "mysql", "mariadb":
		f, target, namespace = newMySQLSuccessorDestination(t, kind, kind)
	default:
		if len(open) > 0 {
			f = newManagedDestinationFixtureOn(t, "sqlite", newWriteGroupRouterFixtureOn(t, open[0]), "uint64")
		} else {
			f = newManagedDestinationFixture(t, "sqlite", "uint64")
		}
		namespace = f.namespace
	}
	clock := &atomic.Int64{}
	clock.Store(time.Date(2026, 10, 28, 11, 59, 0, 0, time.UTC).UnixNano())
	now := func() time.Time { return time.Unix(0, clock.Load()).UTC() }
	f.groups.WithClock(now).WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{})
	var request map[string]any
	if counter {
		tagID := "tag-A"
		require.NoError(t, measurement.NewService(measurement.NewSQLRepository(f.db)).CreateMeasurement(t.Context(), &measurement.MeasurementDefinition{
			ID: "first-start-energy", WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: "point-A", TagID: &tagID, Name: "Energy", Quantity: "energy", Unit: "kWh", SemanticKind: measurement.SemanticKindCounter, NumericEncoding: "uint64",
		}))
		request = profileRequest(t, f.writeGroupRouterFixture, "cumulative", 1, "minute", "main", "")
		request["group"].(map[string]any)["members"].([]map[string]any)[0]["measurement_id"] = "first-start-energy"
	} else {
		sources := []meanSource{{name: "pressure", tagType: "float32", role: "gauge", required: true}}
		seedMeanSources(t, f, sources)
		request = meanLongRequest(t, f, sources)
	}
	group := request["group"].(map[string]any)
	group["layout_version"] = workspace.WriteGroupLayoutManagedWideProfileV1
	if sourceLong {
		group["layout_version"] = workspace.WriteGroupLayoutManagedLongV1
	}
	group["storage_mode"] = mode
	if !sourceLong {
		group["storage_mode"] = ""
	}
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, created))
	id := saved["id"].(string)
	removed := performJSONRequest(t, f.router, http.MethodDelete, writeGroupsPath+"/"+id, scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, removed.Code, removed.Body.String())
	old, err := f.groups.Get(t.Context(), id)
	require.NoError(t, err)
	require.Empty(t, old.Group.AppliedRevision)
	require.Equal(t, workspace.WriteGroupStatusDeleted, old.Group.Status)
	var candidate workspace.WriteGroup
	raw, err := json.Marshal(old.Group)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &candidate))
	candidate.Name = "Never-applied successor"
	candidate.Destination.TableName = ""
	candidate.Destination.SchemaRevision, candidate.Destination.SchemaDigest = "", ""
	candidate.LayoutVersion = workspace.WriteGroupLayoutManagedWideProfileV1
	if targetLong {
		candidate.LayoutVersion = workspace.WriteGroupLayoutManagedLongV1
	}
	candidate.StorageMode = mode
	if !targetLong {
		candidate.StorageMode = ""
		candidate.WritePolicy.RecreateMissing = nil
	}
	if counter && targetLong {
		for i := range candidate.Members {
			candidate.Members[i].SourceValueType = ""
			candidate.Members[i].ResultValueType = ""
		}
	}
	if kind != "sqlite" {
		candidate.Destination.ConnectorID = "connector-B"
		candidate.Destination.ConnectorRevision = "connector-2"
		candidate.Destination.TableSchema = namespace
		if kind == "mysql" || kind == "mariadb" {
			candidate.Destination.Database = namespace
		}
	}
	var lineage string
	var revision int64
	require.NoError(t, f.db.QueryRow(`SELECT lineage_id,revision FROM write_group_lineages WHERE current_group_id=?`, id).Scan(&lineage, &revision))
	boundary := time.Date(2026, 10, 28, 12, 0, 12, 500000000, time.UTC)
	rig := &wideRig{f: f, kind: kind, groupID: id, table: old.Group.Destination.TableName, clock: clock, columns: map[string]string{}, frozen: map[string]workspace.WriteGroupMember{}}
	r := &successorRuntimeRig{successorAPIFixture: successorAPIFixture{wideRig: rig}, counter: counter, boundary: boundary}
	r.request = workspace.WriteGroupSuccessorReviewRequest{WorkspaceID: f.record.ID, ExpectedWorkspaceRevision: old.WorkspaceRevision, LineageID: lineage, ExpectedLineageRevision: revision,
		ExpectedPredecessorRevision: old.Group.Revision, ExpectedPredecessorConnectorRevision: old.Group.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: candidate.Destination.ConnectorRevision,
		OperationID: "first-start-original", Boundary: boundary, Group: &candidate}
	r.preparationReview = r.review(t)
	r.prepared = r.reserve(t, r.preparationReview)
	token := successorAPIPreview(t, r.successorAPIFixture, r.preparationReview, r.prepared)
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/successor", successorAPIConfirmRequest(t, r.successorAPIFixture, r.preparationReview, r.prepared, token))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Equal(t, "succeeded", decodeJSONBody(t, response)["data"].(map[string]any)["status"])
	if kind == "sqlite" {
		target = openSQLiteTarget(t, f.file)
	}
	r.target = target
	r.targets = dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	r.store = groupdelivery.NewStore(f.db).WithClock(now).WithCanonicalIntake()
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{Groups: f.groups, Tags: tag.NewService(tag.NewSQLRepository(f.db)), Destinations: r.targets,
		Inspector: dbtarget.NewReadOnlyTableInspector(r.targets), Store: r.store, Ledger: recordingplan.NewService(recordingplan.NewSQLRepository(f.db)), Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db))},
		grouppipeline.Config{Now: now, MaxFutureSkew: time.Minute, OnError: func(err error) { r.reports = append(r.reports, err) }})
	t.Cleanup(func() { require.NoError(t, r.pipe.Stop(time.Second)) })
	return &firstStartAPIRig{successorRuntimeRig: r, targetNamespace: namespace, schemaToken: token}
}

func (r *firstStartAPIRig) activationAPI() successorActivationAPIFixture {
	return successorActivationAPIFixture{successorAPIFixture: r.successorAPIFixture, reviewResult: r.preparationReview, prepared: r.prepared, boundary: r.boundary}
}
func (r *firstStartAPIRig) confirmation(t *testing.T) map[string]any {
	t.Helper()
	f := r.activationAPI()
	request := f.activationConfirmation(t)
	require.Equal(t, "first_start", request["envelope"].(map[string]any)["activation_kind"])
	request["action"] = "first-start-confirm"
	return request
}
func (r *firstStartAPIRig) statusRequest() map[string]any {
	return map[string]any{"action": "operation-status", "workspace_id": r.f.record.ID, "lineage_id": r.request.LineageID, "operation_id": r.request.OperationID}
}
func (r *firstStartAPIRig) start(t *testing.T) map[string]any {
	t.Helper()
	result := r.activationAPI().postActivation(t, r.confirmation(t))
	snapshot, err := r.f.groups.ResolveAppliedAt(t.Context(), r.prepared.Group.ID, r.boundary)
	require.NoError(t, err)
	r.snapshot = snapshot
	return result
}

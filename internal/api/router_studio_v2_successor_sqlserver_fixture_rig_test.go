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

type sqlServerSuccessorRig struct {
	*successorRuntimeRig
	targetNamespace string
	targetDB        *sql.DB
	targetDatabase  string
	fixtures        [2]*SQLServerCircuitFixture
}

func newSQLServerSuccessorRig(t *testing.T, mode workspace.WriteGroupStorageMode, counter, prepare bool, extraSources ...meanSource) *sqlServerSuccessorRig {
	t.Helper()
	f, targetDB, fixtures := newSQLServerSuccessorDestination(t, mode)
	targetNamespace, targetDatabase := fixtures[1].Manifest.Schema, fixtures[1].Manifest.Database
	clock := &atomic.Int64{}
	start := time.Date(2026, 10, 28, 11, 59, 30, 0, time.UTC)
	if counter {
		start = start.Add(-90 * time.Second)
	}
	clock.Store(start.Add(-time.Second).UnixNano())
	now := func() time.Time { return time.Unix(0, clock.Load()).UTC() }
	f.groups.WithClock(now).WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{})
	var request map[string]any
	if counter {
		tagID := "tag-A"
		require.NoError(t, measurement.NewService(measurement.NewSQLRepository(f.db)).CreateMeasurement(t.Context(), &measurement.MeasurementDefinition{
			ID: "successor-energy", WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: "point-A", TagID: &tagID,
			Name: "Energy", Quantity: "energy", Unit: "kWh", SemanticKind: measurement.SemanticKindCounter, NumericEncoding: "uint64",
		}))
		request = profileRequest(t, f.writeGroupRouterFixture, "cumulative", 1, "minute", "main", "")
		request["group"].(map[string]any)["members"].([]map[string]any)[0]["measurement_id"] = "successor-energy"
	} else {
		sources := append([]meanSource{{name: "pressure", tagType: "float32", role: "gauge", required: true}}, extraSources...)
		seedMeanSources(t, f, sources)
		request = meanLongRequest(t, f, sources)
	}
	group := request["group"].(map[string]any)
	group["layout_version"] = workspace.WriteGroupLayoutManagedLongV1
	group["storage_mode"] = mode
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	group["destination"].(map[string]any)["database"] = fixtures[0].Manifest.Database
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, created))
	id := saved["id"].(string)
	sourceToken := previewLong(t, f, id)
	require.True(t, sqlServerCircuitPreviewCreatesGroup(sourceToken), "the source table must be created by this test")
	registerSQLServerSuccessorPreparationAudit(t, f, sourceToken)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, sourceToken).Status)
	registerSQLServerSuccessorSchemaCleanup(t, f, sourceToken)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	snapshot, err := f.groups.ResolveAppliedAt(t.Context(), id, start.Add(time.Second))
	require.NoError(t, err)
	require.Equal(t, start, snapshot.EffectiveAt)
	clock.Store(start.Add(time.Second).UnixNano())
	rig := &wideRig{f: f, kind: "sqlserver", groupID: id, table: snapshot.Group.Destination.TableName, snapshot: snapshot, clock: clock, columns: map[string]string{}, frozen: map[string]workspace.WriteGroupMember{}}
	for _, member := range snapshot.Group.Members {
		rig.frozen[member.TagID] = member
		rig.columns[member.TagID] = member.TargetColumn
	}
	r := &successorRuntimeRig{successorAPIFixture: successorAPIFixture{wideRig: rig}, counter: counter, boundary: time.Date(2026, 10, 28, 12, 0, 12, 500000000, time.UTC)}
	r.targets = dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.db))
	r.store = groupdelivery.NewStore(f.db).WithClock(now).WithCanonicalIntake()
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{Groups: f.groups, Tags: tag.NewService(tag.NewSQLRepository(f.db)), Destinations: r.targets,
		Inspector: dbtarget.NewReadOnlyTableInspector(r.targets), Store: r.store, Ledger: ledger, Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db))},
		grouppipeline.Config{Now: now, MaxFutureSkew: time.Minute, OnError: func(err error) { r.reports = append(r.reports, err) }})
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	frozenBytes, err := r.store.RuntimeVersion(t.Context(), groupdelivery.GroupKey{WorkspaceID: f.record.ID, GroupID: id, GroupRevision: snapshot.GroupRevision})
	require.NoError(t, err)
	var frozen workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(frozenBytes, &frozen))
	snapshot = &frozen
	rig.snapshot = snapshot
	require.NotNil(t, snapshot.RuntimeLayout)
	require.NoError(t, workspace.ValidateSQLServerRuntimeLayout(snapshot.Group, snapshot.RuntimeLayout))
	require.Equal(t, fixtures[0].Manifest.Database, snapshot.RuntimeLayout.SQLServer.Group.Evidence.Database)
	require.Equal(t, f.namespace, snapshot.RuntimeLayout.SQLServer.Group.Evidence.Schema)

	t.Cleanup(func() { require.NoError(t, r.pipe.Stop(time.Second)) })
	r.target = targetDB
	// Reserve and prepare the same mode in destination B without enabling it.
	current, err := f.groups.Get(t.Context(), id)
	require.NoError(t, err)
	var candidate workspace.WriteGroup
	raw, err := json.Marshal(current.Group)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &candidate))
	candidate.Name = "SQL Server successor"
	candidate.Destination.ConnectorID = "connector-B"
	candidate.Destination.ConnectorRevision = "connector-2"
	candidate.Destination.TableSchema = targetNamespace
	candidate.Destination.Database = targetDatabase
	candidate.Destination.TableName = ""
	candidate.Destination.SchemaRevision, candidate.Destination.SchemaDigest = "", ""
	candidate.LayoutVersion = workspace.WriteGroupLayoutManagedLongV1
	candidate.StorageMode = mode
	var lineage string
	var revision int64
	require.NoError(t, f.db.QueryRow(`SELECT lineage_id,revision FROM write_group_lineages WHERE current_group_id=?`, id).Scan(&lineage, &revision))
	r.request = workspace.WriteGroupSuccessorReviewRequest{WorkspaceID: f.record.ID, ExpectedWorkspaceRevision: current.WorkspaceRevision,
		LineageID: lineage, ExpectedLineageRevision: revision, ExpectedPredecessorRevision: current.Group.Revision,
		ExpectedPredecessorConnectorRevision: current.Group.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: "connector-2",
		OperationID: "runtime-transition", Boundary: r.boundary, Group: &candidate}
	if !prepare {
		return &sqlServerSuccessorRig{successorRuntimeRig: r, targetNamespace: targetNamespace, targetDB: targetDB, targetDatabase: targetDatabase, fixtures: fixtures}
	}
	r.preparationReview = r.review(t)
	r.prepared = r.reserve(t, r.preparationReview)
	token := successorAPIPreview(t, r.successorAPIFixture, r.preparationReview, r.prepared)
	targetFixture := f
	targetFixture.target, targetFixture.namespace = targetDB, targetNamespace
	registerSQLServerSuccessorPreparationAudit(t, targetFixture, token)
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/successor", successorAPIConfirmRequest(t, r.successorAPIFixture, r.preparationReview, r.prepared, token))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var outcome struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &outcome))
	require.Equal(t, recordingplan.SchemaOperationSucceeded, outcome.Data.Status)
	registerSQLServerSuccessorSchemaCleanup(t, targetFixture, token)
	return &sqlServerSuccessorRig{successorRuntimeRig: r, targetNamespace: targetNamespace, targetDB: targetDB, targetDatabase: targetDatabase, fixtures: fixtures}
}

func (r *sqlServerSuccessorRig) activationAPI() successorActivationAPIFixture {
	return successorActivationAPIFixture{successorAPIFixture: r.successorAPIFixture, reviewResult: r.preparationReview, prepared: r.prepared, boundary: r.boundary}
}

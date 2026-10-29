package api

import (
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

type postgresSuccessorRig struct {
	*successorRuntimeRig
	targetNamespace string
}

func newPostgresSuccessorRuntimeRig(t *testing.T, sourceLong, targetLong, counter, timescale bool) *postgresSuccessorRig {
	return newPostgresSuccessorRig(t, sourceLong, targetLong, counter, timescale, true)
}

func newPostgresSuccessorRig(t *testing.T, sourceLong, targetLong, counter, timescale, prepare bool, extraSources ...meanSource) *postgresSuccessorRig {
	t.Helper()
	f, targetNamespace := newPostgresSuccessorDestination(t, timescale)
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
	group["layout_version"] = workspace.WriteGroupLayoutManagedWideProfileV1
	if sourceLong {
		group["layout_version"] = workspace.WriteGroupLayoutManagedLongV1
	}
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	if timescale {
		group["storage_mode"] = "timescale"
	}
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, created))
	id := saved["id"].(string)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, previewLong(t, f, id)).Status)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	snapshot, err := f.groups.ResolveAppliedAt(t.Context(), id, start.Add(time.Second))
	require.NoError(t, err)
	require.Equal(t, start, snapshot.EffectiveAt)
	clock.Store(start.Add(time.Second).UnixNano())
	rig := &wideRig{f: f, kind: "postgres", groupID: id, table: snapshot.Group.Destination.TableName, snapshot: snapshot, clock: clock, columns: map[string]string{}, frozen: map[string]workspace.WriteGroupMember{}}
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
	t.Cleanup(func() { require.NoError(t, r.pipe.Stop(time.Second)) })
	r.target = f.target
	// Reserve and prepare the opposite layout without enabling it.
	current, err := f.groups.Get(t.Context(), id)
	require.NoError(t, err)
	var candidate workspace.WriteGroup
	raw, err := json.Marshal(current.Group)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &candidate))
	candidate.Name = "PostgreSQL successor"
	candidate.Destination.ConnectorID = "connector-B"
	candidate.Destination.ConnectorRevision = "connector-2"
	candidate.Destination.TableSchema = targetNamespace
	candidate.Destination.TableName = ""
	candidate.Destination.SchemaRevision, candidate.Destination.SchemaDigest = "", ""
	candidate.LayoutVersion = workspace.WriteGroupLayoutManagedLongV1
	candidate.StorageMode = workspace.WriteGroupStorageModePlain
	if counter && targetLong {
		// Canonical unscaled cumulative long stores raw typing in its counter proof,
		// not the wide-only member metadata. Meaning and measurement stay identical.
		for i := range candidate.Members {
			candidate.Members[i].SourceValueType = ""
			candidate.Members[i].ResultValueType = ""
		}
	}
	if !targetLong {
		candidate.LayoutVersion = workspace.WriteGroupLayoutManagedWideProfileV1
		candidate.StorageMode = ""
		candidate.WritePolicy.RecreateMissing = nil
	}
	if timescale {
		candidate.StorageMode = "timescale"
	}
	var lineage string
	var revision int64
	require.NoError(t, f.db.QueryRow(`SELECT lineage_id,revision FROM write_group_lineages WHERE current_group_id=?`, id).Scan(&lineage, &revision))
	r.request = workspace.WriteGroupSuccessorReviewRequest{WorkspaceID: f.record.ID, ExpectedWorkspaceRevision: current.WorkspaceRevision,
		LineageID: lineage, ExpectedLineageRevision: revision, ExpectedPredecessorRevision: current.Group.Revision,
		ExpectedPredecessorConnectorRevision: current.Group.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: "connector-2",
		OperationID: "runtime-transition", Boundary: r.boundary, Group: &candidate}
	if !prepare {
		return &postgresSuccessorRig{successorRuntimeRig: r, targetNamespace: targetNamespace}
	}
	r.preparationReview = r.review(t)
	r.prepared = r.reserve(t, r.preparationReview)
	token := successorAPIPreview(t, r.successorAPIFixture, r.preparationReview, r.prepared)
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/successor", successorAPIConfirmRequest(t, r.successorAPIFixture, r.preparationReview, r.prepared, token))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var outcome struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &outcome))
	require.Equal(t, recordingplan.SchemaOperationSucceeded, outcome.Data.Status)
	return &postgresSuccessorRig{successorRuntimeRig: r, targetNamespace: targetNamespace}
}

package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/grouptestwrite"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const mysqlNativeReconfirmConnectorRevision = "native-reconfirm-connector-v2"

func TestMySQLNativeExplicitReconfirmationAfterConnectorRevision(t *testing.T) {
	for _, family := range []string{mysqlCircuitFamilyMySQL, mysqlCircuitFamilyMariaDB} {
		t.Run(family, func(t *testing.T) {
			runMySQLNativeExplicitReconfirmation(t, family)
		})
	}
}

func runMySQLNativeExplicitReconfirmation(t *testing.T, family string) {
	f := newMySQLCircuitPlainFixture(t, mysqlCircuitPlainFixtureOptions{
		Family: family, SeedMeasurement: measurement.MeasurementDefinition{ExactRecordingScale: "0.1"},
		Inputs: mysqlCircuitDecimalInputs, StorageMode: workspace.WriteGroupStorageModeNativeMonthly,
	})
	scope := workspace.StorageScopeOf(f.instant.group)
	lineageBefore, err := recordingplan.ReadStorageLineage(t.Context(), f.f.db, scope)
	require.NoError(t, err)
	require.Len(t, lineageBefore.Layout.MySQL.NativeMonthly.Partitions, 6)
	oldCapabilityRevision := lineageBefore.Layout.MySQL.CapabilityRevision

	oldSnapshots := map[string][]byte{}
	for _, profile := range []mysqlCircuitProfile{f.instant, f.cumulative} {
		payload, err := f.store.RuntimeVersion(t.Context(), profile.key)
		require.NoError(t, err)
		oldSnapshots[profile.id] = payload
	}

	// Register both production-backed runtime writers before invoking their
	// partition maintenance, and retain actual accepted effects for the check.
	f.accept(0)
	f.closeAfter(10 * time.Minute)
	waitForMySQLCircuitOutboxes(t, f, 41)
	oldPayloads := mysqlPartitionPayloadDigests(t, f)

	anchor, err := time.Parse(time.RFC3339, lineageBefore.Layout.MySQL.NativeMonthly.Partitions[0].LessThan)
	require.NoError(t, err)
	f.clock.Store(anchor.AddDate(0, 1, 2).UnixNano())
	require.NoError(t, f.pipe.MaintainMySQLPartitions(t.Context()))
	stateBefore, err := recordingplan.ReadMySQLPartitionMaintenanceState(t.Context(), f.f.db, scope)
	require.NoError(t, err)
	require.NotNil(t, stateBefore)
	require.Len(t, stateBefore.Current.Partitions, 7)

	_, err = f.f.db.ExecContext(t.Context(), `UPDATE database_connectors SET identity_revision = ? WHERE id = ?`,
		mysqlNativeReconfirmConnectorRevision, f.instant.group.Destination.ConnectorID)
	require.NoError(t, err)
	for _, id := range []string{f.instant.id, f.cumulative.id} {
		updateNativeReconfirmationGroup(t, f.f, id, mysqlNativeReconfirmConnectorRevision)
	}

	oldLayout := *lineageBefore.Layout
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.f.db))
	inspector := f.f.groups.ManagedTableInspector(targets)
	capability, err := inspector.ProbeMySQLRecordingAtRevision(t.Context(), f.instant.group.Destination.ConnectorID, mysqlNativeReconfirmConnectorRevision)
	require.NoError(t, err)
	require.Equal(t, mysqlNativeReconfirmConnectorRevision, capability.ConnectorRevision)
	require.NotEqual(t, oldCapabilityRevision, capability.CapabilityRevision)

	for _, id := range []string{f.instant.id, f.cumulative.id} {
		token := previewMySQLAtRevision(t, f.f, id, mysqlNativeReconfirmConnectorRevision)
		require.NotNil(t, token.GroupLayout)
		require.Equal(t, capability.CapabilityRevision, token.GroupLayout.MySQL.CapabilityRevision)
		require.Empty(t, token.Statements, "explicit reconfirmation must not CREATE an existing table")
		op := confirmMySQLAtRevision(t, f.f, id, mysqlNativeReconfirmConnectorRevision, token)
		require.Equal(t, recordingplan.SchemaOperationSucceeded, op.Status, op.Reason)
	}

	lineageAfter, err := recordingplan.ReadStorageLineage(t.Context(), f.f.db, scope)
	require.NoError(t, err)
	require.Equal(t, lineageBefore.Generation, lineageAfter.Generation)
	require.Equal(t, capability.CapabilityRevision, lineageAfter.Layout.MySQL.CapabilityRevision)
	require.Equal(t, mysqlNativeReconfirmConnectorRevision, lineageAfter.ConnectorRevision)
	stateAfter, err := recordingplan.ReadMySQLPartitionMaintenanceState(t.Context(), f.f.db, scope)
	require.NoError(t, err)
	require.Equal(t, stateBefore.Current, stateAfter.Current, "the seven authorized partitions remain current")
	require.Equal(t, lineageAfter.Generation, stateAfter.Generation)
	require.Equal(t, digestOfNativeLayout(t, *lineageAfter.Layout), stateAfter.BaseLayoutDigest)
	require.Equal(t, mysqlNativeReconfirmConnectorRevision, stateAfter.ConnectorRevision)

	require.Equal(t, oldPayloads, mysqlPartitionPayloadDigests(t, f), "accepted outbox bytes remain unchanged")
	for _, profile := range []mysqlCircuitProfile{f.instant, f.cumulative} {
		payload, err := f.store.RuntimeVersion(t.Context(), profile.key)
		require.NoError(t, err)
		require.Equal(t, oldSnapshots[profile.id], payload, "frozen runtime snapshot bytes remain unchanged")
		ready, err := f.f.groups.Readiness(t.Context(), profile.id)
		require.NoError(t, err)
		require.True(t, ready.SchemaReady, "%+v", ready.Issues)
		require.True(t, ready.Ready, "verified destination is ready: %+v", ready.Issues)
	}

	// The same physical catalog must still be inspected with the new capability,
	// while the old automatic proof cannot authorize the new binding.
	physical, err := inspector.InspectMySQLManagedTableAtRevision(t.Context(), capability.ConnectorID, capability.ConnectorRevision,
		capability.Database, capability.CapabilityRevision, f.table)
	require.NoError(t, err)
	actual := mysqlPartitionTargetInspection(physical)
	_, err = recordingplan.NormalizeMySQLPartitionCatalog(t.Context(), f.f.db, scope, oldLayout, actual)
	require.ErrorIs(t, err, recordingplan.ErrMySQLPartitionMaintenanceStateMismatch)

	altered := mysqlPartitionTargetInspection(physical)
	altered.MySQL.Partitions.Partitions[0].Name = "p_unauthorized"
	_, err = recordingplan.NormalizeMySQLPartitionCatalogForExplicitReconfirmation(t.Context(), f.f.db, scope,
		*lineageAfter.Layout, altered, mysqlNativeReconfirmConnectorRevision)
	require.ErrorIs(t, err, recordingplan.ErrIncompatibleExistingTable)

	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.f.db))
	testWriter := grouptestwrite.New(grouptestwrite.Dependencies{
		Groups: f.f.groups, Tags: tag.NewService(tag.NewSQLRepository(f.f.db)),
		Destinations: targets, Inspector: dbtarget.NewReadOnlyTableInspector(targets), Ledger: ledger,
	}, grouptestwrite.Config{})
	router := NewRouter(&DatalinkServices{
		Workspace: f.f.workspace, WriteGroups: f.f.groups, RecordingPlan: ledger,
		WriteGroupTestWrite: testWriter, DBTarget: targets,
	})
	for _, profile := range []mysqlCircuitProfile{f.instant, f.cumulative} {
		resp := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/"+profile.id+"/test-write-preview", nil)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	}
}

func updateNativeReconfirmationGroup(t *testing.T, f managedGroupFixture, id, connectorRevision string) {
	t.Helper()
	data := getGroupData(t, f, id)
	group := getGroup(t, f, id)
	destinationOf(group)["connector_revision"] = connectorRevision
	request := groupEditRequest(f.writeGroupRouterFixture, data, group)
	request["expected_connector_revision"] = connectorRevision
	response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
}

func mysqlRevisionScopeBody(t *testing.T, f managedGroupFixture, id, connectorRevision string) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, response))
	request := groupEditRequest(f.writeGroupRouterFixture, data, group)
	request["expected_connector_revision"] = connectorRevision
	delete(request, "group")
	return request
}

func previewMySQLAtRevision(t *testing.T, f managedGroupFixture, id, connectorRevision string) recordingplan.SchemaPreviewToken {
	t.Helper()
	request := mysqlRevisionScopeBody(t, f, id, connectorRevision)
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/schema-preview", request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data recordingplan.SchemaPreviewToken `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	return body.Data
}

func confirmMySQLAtRevision(t *testing.T, f managedGroupFixture, id, connectorRevision string, token recordingplan.SchemaPreviewToken) recordingplan.SchemaOperation {
	t.Helper()
	request := mysqlRevisionScopeBody(t, f, id, connectorRevision)
	request["token"], request["operation_id"] = token.Token, token.OperationID
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/schema-apply", request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	return body.Data
}

func digestOfNativeLayout(t *testing.T, layout recordingplan.GroupSchemaLayout) string {
	t.Helper()
	encoded, err := json.Marshal(layout)
	require.NoError(t, err)
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

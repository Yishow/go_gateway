package main

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/api"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestProductionDefaultUIUint64ToSQL(t *testing.T) {
	productionDefaultUIValueToSQL(t, schema.DataTypeUint64, 1, 0, uint64(9007199254740993), "uint64", uint64(9007199254740993), "9007199254740993")
}

func TestProductionDefaultUIScaledInt16ToSQL(t *testing.T) {
	productionDefaultUIValueToSQL(t, schema.DataTypeInt16, 0.5, 10, int16(243), "float64", float64(131.5), "131.5")
}

func productionDefaultUIValueToSQL(t *testing.T, source schema.DataType, scale, offset float64, input any, targetType string, expected any, digitsExpected string) {
	t.Helper()
	target := newManagedPipelineTarget(t, "sqlite")
	env := target.env
	_, err := env.db.ExecContext(t.Context(), `UPDATE devices SET status='active', last_test_success=1, last_test_error='', readiness_status='{"probe_status":"success"}', description='', created_at='2026-10-05 00:00:00', updated_at='2026-10-05 00:00:00' WHERE id='device-1'`)
	require.NoError(t, err)
	rule, err := env.services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{ID: "default-u64", DeviceID: "device-1", StartAddress: "41001", Count: 4, DataType: source, ScaleMultiplier: &scale, ScaleOffset: &offset, NamingPrefix: "Counter_", Enabled: true})
	require.NoError(t, err)
	links, err := env.services.sourceRule.ListLinks(t.Context(), rule.ID)
	require.NoError(t, err)
	require.NotEmpty(t, links)
	link := links[0]
	point := map[string]any{"id": link.PointID, "device_id": "device-1", "rule_id": rule.ID, "name": "Counter", "address": link.Address, "data_type": string(source), "_rule_scale": scale, "_rule_offset": offset}
	raw, err := json.Marshal(point)
	require.NoError(t, err)
	// Execute the actual default UI builder, rather than copy its defaults in Go.
	module, err := filepath.Abs("../../frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts")
	require.NoError(t, err)
	cmd := exec.CommandContext(t.Context(), "node", "--experimental-strip-types", "--input-type=module", "-e", `const {buildDefaultMapping}=await import(process.argv[1]); process.stdout.write(JSON.stringify(buildDefaultMapping(JSON.parse(process.argv[2]),0)));`, module, string(raw))
	output, err := cmd.Output()
	require.NoError(t, err)
	var request map[string]any
	require.NoError(t, json.Unmarshal(output, &request))
	require.Equal(t, targetType, request["target_type"])
	router := api.NewRouter(&api.DatalinkServices{Workspace: env.services.workspace, Device: env.services.device, SourceRule: env.services.sourceRule, Point: env.services.point, Tag: env.services.tag, Mapping: env.services.mapping})
	status, body := httpJSON(t, router, http.MethodPost, testWriteBase+"/mappings", request)
	require.Contains(t, []int{http.StatusOK, http.StatusCreated}, status, "%+v", body)
	data := body["data"].(map[string]any)
	persisted, err := env.services.mapping.GetByID(t.Context(), data["id"].(string))
	require.NoError(t, err)
	require.Equal(t, link.PointID, persisted.PointID)
	require.True(t, persisted.Enabled)
	require.Equal(t, schema.MappingStatusActive, persisted.Status)
	current, err := env.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Contains(t, current.OrderedDeviceIDs, "device-1")
	converted, err := mapping.ExecutePipeline(input, persisted.TransformPipeline)
	require.NoError(t, err)
	require.Equal(t, expected, converted.CurrentValue)
	var declared string
	require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT data_type FROM tags WHERE id=?`, persisted.TagID).Scan(&declared))
	require.Equal(t, targetType, declared)
	revisionBefore := persisted.TransformPipeline
	for range 2 {
		require.NoError(t, env.services.sourceRule.SyncDerivedPointState(t.Context()))
		afterRestart, err := env.services.mapping.GetByID(t.Context(), persisted.ID)
		require.NoError(t, err)
		require.Equal(t, revisionBefore, afterRestart.TransformPipeline)
		require.Equal(t, persisted.LastAppliedSignature, afterRestart.LastAppliedSignature)
		require.Equal(t, persisted.ProposedSignature, afterRestart.ProposedSignature)
		require.Equal(t, schema.MappingStatusActive, afterRestart.Status)
		require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT data_type FROM tags WHERE id=?`, persisted.TagID).Scan(&declared))
		require.Equal(t, targetType, declared)
	}

	clock := &testClock{}
	base := time.Now().UTC().Truncate(10 * time.Second).Add(-time.Minute)
	clock.set(base.Add(-time.Second))
	env.services.writeGroups.WithClock(clock.now)
	group := target.prepare(t, []workspace.WriteGroupMember{{DeviceID: "device-1", PointID: link.PointID, TagID: persisted.TagID, Required: true}}, workspace.WriteGroupRowPolicy{IntervalSeconds: 10})
	pipe := env.pipeline(clock, "default-ui")
	clock.set(base.Add(time.Second))
	require.NoError(t, pipe.Reconcile(t.Context()))
	require.NoError(t, pipe.AcceptSample(t.Context(), env.envelope(group, 0, "default-u64-sample", base.Add(2*time.Second), converted.CurrentValue)))
	var journal string
	require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_samples WHERE sample_id = 'default-u64-sample'`).Scan(&journal))
	require.Contains(t, journal, digitsExpected)
	clock.set(base.Add(11 * time.Second))
	pipe.TickAll(t.Context())
	var effect, payload string
	require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT effect_key,payload FROM wg_delivery_outbox WHERE group_id=?`, group.ID).Scan(&effect, &payload))
	require.Contains(t, payload, digitsExpected)
	sender := groupdelivery.NewSender(groupdelivery.NewStore(env.db), managedAckResolver{service: env.services.dbTarget, db: target.db, kind: target.kind}, groupdelivery.SenderConfig{Owner: "default-ui"})
	result, err := sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateCommitted, result.State)
	var digits string
	require.NoError(t, target.db.QueryRowContext(t.Context(), `SELECT CAST("`+group.Members[0].TargetColumn+`" AS TEXT) FROM `+managedSQLTable(target, group.Destination.TableName)).Scan(&digits))
	require.Equal(t, digitsExpected, digits)
}

package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

// TestMySQLGroupProbeLifecycle tests the complete group test-write lifecycle for
// MySQL and MariaDB across plain and native_monthly storage modes.
// It verifies:
//  1. Bare requests and stale scope fail before any target mutations occur.
//  2. Real grouptestwrite.New service wires canonical workspace, groups, tags, connectors, and plans.
//  3. For native_monthly, MaintainMySQLPartitions expands physical partitions (e.g. 7 partitions)
//     and test-write utilizes the current authorized proof while immutable confirmed layout remains 6 partitions.
//  4. Writes, readback verification, and cleanup operate independently; exact test record IDs are derived
//     truthfully from preview or operation detail and proven cleaned up.
//  5. Existing production rows and receipts remain completely untouched, verified by full row column snapshots and SHA256 digests.
//  6. Replaying a completed test-write returns retained status without re-inserting rows.
func TestMySQLGroupProbeLifecycle(t *testing.T) {
	for _, family := range []string{mysqlCircuitFamilyMySQL, mysqlCircuitFamilyMariaDB} {
		for _, mode := range []workspace.WriteGroupStorageMode{
			workspace.WriteGroupStorageModePlain,
			workspace.WriteGroupStorageModeNativeMonthly,
		} {
			t.Run(family+"_"+string(mode), func(t *testing.T) {
				runMySQLGroupProbeLifecycle(t, family, mode)
			})
		}
	}
}

func runMySQLGroupProbeLifecycle(t *testing.T, family string, mode workspace.WriteGroupStorageMode) {
	f := newMySQLCircuitPlainFixture(t, mysqlCircuitPlainFixtureOptions{
		Family:          family,
		SeedMeasurement: measurement.MeasurementDefinition{ExactRecordingScale: "0.1"},
		Inputs:          mysqlCircuitDecimalInputs,
		StorageMode:     mode,
	})

	// 1. Ingest sample 0 and wait for durable commitment of 41 rows & 41 receipts.
	f.accept(0)
	f.closeAfter(10 * time.Minute)
	waitForMySQLCircuitOutboxes(t, f, 41)
	assertMySQLPartitionRemoteCounts(t, f, 41)

	readReceiptDigests := func() map[string]string {
		rows, err := f.f.target.QueryContext(t.Context(),
			fmt.Sprintf("SELECT effect_key, payload_digest FROM %#q ORDER BY effect_key", dbtarget.EffectReceiptTable))
		require.NoError(t, err)
		defer rows.Close()
		receipts := make(map[string]string, 41)
		for rows.Next() {
			var key, digest string
			require.NoError(t, rows.Scan(&key, &digest))
			receipts[key] = digest
		}
		require.NoError(t, rows.Err())
		return receipts
	}

	readProductionRowDigests := func() map[string]string {
		rows, err := f.f.target.QueryContext(t.Context(),
			fmt.Sprintf("SELECT * FROM %#q ORDER BY bucket_start, record_id", f.table))
		require.NoError(t, err)
		defer rows.Close()

		cols, err := rows.Columns()
		require.NoError(t, err)

		recordIdx := -1
		for i, col := range cols {
			if col == "record_id" {
				recordIdx = i
				break
			}
		}
		require.NotEqual(t, -1, recordIdx, "target table must have record_id column")

		result := make(map[string]string, 41)
		for rows.Next() {
			values := make([]any, len(cols))
			ptrs := make([]any, len(cols))
			for i := range values {
				ptrs[i] = &values[i]
			}
			require.NoError(t, rows.Scan(ptrs...))

			canonical := make([]any, len(cols))
			for i, val := range values {
				switch v := val.(type) {
				case []byte:
					canonical[i] = string(v)
				case time.Time:
					canonical[i] = v.UTC().Format(time.RFC3339Nano)
				default:
					canonical[i] = v
				}
			}

			rawJSON, err := json.Marshal(canonical)
			require.NoError(t, err)
			sum := sha256.Sum256(rawJSON)
			hash := hex.EncodeToString(sum[:])

			recID, ok := canonical[recordIdx].(string)
			require.True(t, ok)
			result[recID] = hash
		}
		require.NoError(t, rows.Err())
		return result
	}

	receiptDigestsBefore := readReceiptDigests()
	require.Len(t, receiptDigestsBefore, 41)
	productionRowsBefore := readProductionRowDigests()
	require.Len(t, productionRowsBefore, 41)

	// 2. For native_monthly, advance horizon via MaintainMySQLPartitions to test
	// that current physical partition proof (7 partitions) is used while confirmed layout (6) is untouched.
	if mode == workspace.WriteGroupStorageModeNativeMonthly {
		scope := workspace.StorageScopeOf(f.instant.group)
		lineageBefore, err := recordingplan.ReadStorageLineage(t.Context(), f.f.db, scope)
		require.NoError(t, err)
		base := lineageBefore.Layout.MySQL.NativeMonthly
		require.NotNil(t, base)
		require.Len(t, base.Partitions, 6)

		anchor, err := time.Parse(time.RFC3339, base.Partitions[0].LessThan)
		require.NoError(t, err)
		f.clock.Store(anchor.AddDate(0, 1, 2).UnixNano())
		require.NoError(t, f.pipe.MaintainMySQLPartitions(t.Context()))

		state, err := recordingplan.ReadMySQLPartitionMaintenanceState(t.Context(), f.f.db, scope)
		require.NoError(t, err)
		require.NotNil(t, state)
		require.Equal(t, "ready", state.Status)
		require.Len(t, state.Current.Partitions, 7)

		// Restore clock to baseline offset
		f.clock.Store(f.base.Add(10 * time.Minute).UnixNano())
	}

	// 3. Build real grouptestwrite.New service and wire into Router.
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.f.db))
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.f.db))
	tags := tag.NewService(tag.NewSQLRepository(f.f.db))
	inspector := dbtarget.NewReadOnlyTableInspector(targets)

	testWriter := grouptestwrite.New(grouptestwrite.Dependencies{
		Groups:       f.f.groups,
		Tags:         tags,
		Destinations: targets,
		Inspector:    inspector,
		Ledger:       ledger,
	}, grouptestwrite.Config{})

	router := NewRouter(&DatalinkServices{
		Workspace:           f.f.workspace,
		WriteGroups:         f.f.groups,
		RecordingPlan:       ledger,
		WriteGroupTestWrite: testWriter,
		DBTarget:            targets,
	})

	// 4. Negative tests: bare requests and stale tokens must fail before writes.
	// Bare confirm without token/operation_id
	resp := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/"+f.instant.id+"/test-write", nil)
	require.Equal(t, http.StatusBadRequest, resp.Code)

	// Confirm with foreign / non-existent token
	resp = performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/"+f.instant.id+"/test-write", map[string]any{
		"token":        "tok-stale-invalid",
		"operation_id": "op-stale-invalid",
	})
	require.Equal(t, http.StatusNotFound, resp.Code)

	// Legacy route bare confirm
	resp = performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write", map[string]any{
		"plan_id": "plan-unresolved",
	})
	require.Equal(t, http.StatusBadRequest, resp.Code)

	// Legacy route preview for unresolved plan
	resp = performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write-preview", map[string]any{
		"plan_id": "plan-unresolved",
	})
	require.Equal(t, http.StatusUnprocessableEntity, resp.Code)

	// Full immutable digest comparison: receipts and production rows remain untouched
	require.Equal(t, receiptDigestsBefore, readReceiptDigests())
	require.Equal(t, productionRowsBefore, readProductionRowDigests())

	// 5. Test lifecycle for both instant and cumulative profiles.
	for _, profile := range []mysqlCircuitProfile{f.instant, f.cumulative} {
		// Preview
		previewResp := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/"+profile.id+"/test-write-preview", nil)
		require.Equal(t, http.StatusOK, previewResp.Code, previewResp.Body.String())
		previewData := decodeJSONBody(t, previewResp)["data"].(map[string]any)
		token := previewData["token"].(string)
		operationID := previewData["operation_id"].(string)
		require.NotEmpty(t, token)
		require.NotEmpty(t, operationID)
		require.Equal(t, "test_write", previewData["action"])
		require.Equal(t, profile.id, previewData["group_id"])
		values := previewData["values"].([]any)
		require.NotEmpty(t, values)

		var recordIDs []string
		if rawIDs, ok := previewData["record_ids"].([]any); ok && len(rawIDs) > 0 {
			for _, raw := range rawIDs {
				if id, ok := raw.(string); ok && id != "" {
					recordIDs = append(recordIDs, id)
				}
			}
		}

		// Confirm
		confirmResp := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/"+profile.id+"/test-write", map[string]any{
			"token":        token,
			"operation_id": operationID,
		})
		require.Equal(t, http.StatusOK, confirmResp.Code, confirmResp.Body.String())
		opData := decodeJSONBody(t, confirmResp)["data"].(map[string]any)
		require.Equal(t, operationID, opData["operation_id"])
		require.Equal(t, "succeeded", opData["status"], confirmResp.Body.String())
		require.Equal(t, grouptestwrite.WriteVerified, opData["write_outcome"])
		require.Equal(t, grouptestwrite.CleanupCleaned, opData["cleanup_status"], confirmResp.Body.String())

		// Derive test record IDs from operation ledger detail if not present in preview
		if len(recordIDs) == 0 {
			if detailStr, ok := opData["detail"].(string); ok && detailStr != "" {
				var prog struct {
					RecordIDs []string `json:"record_ids"`
				}
				if err := json.Unmarshal([]byte(detailStr), &prog); err == nil {
					recordIDs = append(recordIDs, prog.RecordIDs...)
				}
			}
		}
		require.NotEmpty(t, recordIDs, "exact test record IDs must be present in preview or operation detail")

		// Exact test record assertion: every test record ID must be completely cleaned up
		for _, recordID := range recordIDs {
			var count int
			require.NoError(t, f.f.target.QueryRowContext(t.Context(),
				fmt.Sprintf("SELECT COUNT(*) FROM %#q WHERE `record_id` = ?", f.table), recordID).Scan(&count))
			require.Equal(t, 0, count, "test record %s must be removed by cleanup", recordID)
		}

		// Full immutable digest comparison: receipts and production rows remain identical
		require.Equal(t, receiptDigestsBefore, readReceiptDigests(), "destination receipts must remain identical in key and digest")
		require.Equal(t, productionRowsBefore, readProductionRowDigests(), "destination production rows must remain identical across all columns and values")

		// Replay completed confirmation: returns retained status and performs no re-insert
		replayResp := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/"+profile.id+"/test-write", map[string]any{
			"token":        token,
			"operation_id": operationID,
		})
		require.Equal(t, http.StatusOK, replayResp.Code, replayResp.Body.String())
		replayData := decodeJSONBody(t, replayResp)["data"].(map[string]any)
		require.Equal(t, "succeeded", replayData["status"])
		require.Equal(t, grouptestwrite.WriteVerified, replayData["write_outcome"])
		require.Equal(t, grouptestwrite.CleanupCleaned, replayData["cleanup_status"])

		require.Equal(t, receiptDigestsBefore, readReceiptDigests(), "replay must not mutate receipts")
		require.Equal(t, productionRowsBefore, readProductionRowDigests(), "replay must not mutate production rows")
	}

	// 6. For native_monthly, verify original confirmed layout remains 6 partitions
	if mode == workspace.WriteGroupStorageModeNativeMonthly {
		scope := workspace.StorageScopeOf(f.instant.group)
		lineageAfter, err := recordingplan.ReadStorageLineage(t.Context(), f.f.db, scope)
		require.NoError(t, err)
		require.Equal(t, 6, len(lineageAfter.Layout.MySQL.NativeMonthly.Partitions), "immutable original layout remains 6 partitions")
	}
}

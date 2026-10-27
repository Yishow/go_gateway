package api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// TestMySQLDestinationCutover exercises the existing same-group Update/Apply
// boundary. Accepted effects keep their original destination; the new version
// consumes later samples. The reviewed source-transfer API is a separate flow.
func TestMySQLDestinationCutover(t *testing.T) {
	for _, family := range []string{mysqlCircuitFamilyMySQL, mysqlCircuitFamilyMariaDB} {
		t.Run(family, func(t *testing.T) {
			f := newMySQLCircuitPlainFixture(t, mysqlCircuitPlainFixtureOptions{
				Family: family, CanonicalApply: true,
				SeedMeasurement: measurement.MeasurementDefinition{ExactRecordingScale: "0.1"}, Inputs: mysqlCircuitDecimalInputs,
			})
			f.faultTargets.failOpen.Store(1_000_000)
			f.accept(0)
			f.closeAfter(10 * time.Minute)
			oldPayloads := mysqlCutoverFrozenOutbox(t, f)
			require.Len(t, oldPayloads, 41)
			assertMySQLPartitionRemoteCounts(t, f, 0)

			// The opaque saved configuration stays inside this disposable local
			// database. A distinct connector ID creates a distinct canonical table.
			const nextConnector = "connector-cutover"
			const nextRevision = "cutover-connector-v1"
			_, err := f.f.db.ExecContext(t.Context(), `INSERT INTO database_connectors
				(id, name, kind, status, connection_config, identity_revision, enabled)
				SELECT ?, 'Cutover fixture', kind, status, connection_config, ?, enabled
				FROM database_connectors WHERE id = ?`, nextConnector, nextRevision, f.instant.group.Destination.ConnectorID)
			require.NoError(t, err)
			f.f.groups.WithClock(func() time.Time { return f.base.Add(10 * time.Minute) })
			data := getGroupData(t, f.f, f.instant.id)
			group := getGroup(t, f.f, f.instant.id)
			destinationOf(group)["connector_id"] = nextConnector
			destinationOf(group)["connector_revision"] = nextRevision
			destinationOf(group)["table_name"] = ""
			destinationOf(group)["schema_revision"] = ""
			destinationOf(group)["schema_digest"] = ""
			request := groupEditRequest(f.f.writeGroupRouterFixture, data, group)
			request["expected_connector_revision"] = nextRevision
			for _, field := range []string{"expected_workspace_revision", "expected_group_revision", "expected_connector_revision"} {
				fresh := request[field]
				request[field] = "stale-cutover-revision"
				resp := performJSONRequest(t, f.f.router, http.MethodPut, writeGroupsPath+"/"+f.instant.id, request)
				require.Equal(t, http.StatusConflict, resp.Code, resp.Body.String())
				request[field] = fresh
			}
			resp := performJSONRequest(t, f.f.router, http.MethodPut, writeGroupsPath+"/"+f.instant.id, request)
			require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
			_, saved := groupSaveData(t, decodeJSONBody(t, resp))
			newTable := destinationOf(saved)["table_name"].(string)
			require.Regexp(t, longTablePattern, newTable)
			require.NotEqual(t, f.table, newTable)
			t.Cleanup(func() {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_, err := f.f.target.ExecContext(ctx, "DROP TABLE IF EXISTS `"+newTable+"`")
				require.NoError(t, err)
			})
			token := previewMySQLAtRevision(t, f.f, f.instant.id, nextRevision)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmMySQLAtRevision(t, f.f, f.instant.id, nextRevision, token).Status)
			apply := mysqlRevisionScopeBody(t, f.f, f.instant.id, nextRevision)
			for _, field := range []string{"expected_workspace_revision", "expected_group_revision", "expected_connector_revision"} {
				fresh := apply[field]
				apply[field] = "stale-cutover-revision"
				resp = performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.instant.id+"/apply", apply)
				require.Equal(t, http.StatusConflict, resp.Code, resp.Body.String())
				apply[field] = fresh
			}
			resp = performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.instant.id+"/apply", apply)
			require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
			current, err := f.f.groups.ResolveAppliedAt(t.Context(), f.instant.id, f.base.Add(time.Hour))
			require.NoError(t, err)
			require.True(t, current.EffectiveAt.After(f.base.Add(10*time.Minute)))
			require.Equal(t, newTable, current.Group.Destination.TableName)
			require.NotEqual(t, f.instant.group.Revision, current.GroupRevision)

			// Drain the already accepted batch after saving the new destination.
			f.faultTargets.failOpen.Store(0)
			waitForMySQLCircuitOutboxes(t, f, 41)
			assertMySQLPartitionRemoteCounts(t, f, 41)
			var count int
			require.NoError(t, f.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+newTable+"`").Scan(&count))
			require.Zero(t, count)

			// Reconcile and deliver a real post-boundary envelope through canonical
			// version lookup. Cumulative retains its original destination.
			observed := current.EffectiveAt.Add(time.Second)
			f.clock.Store(observed.Add(time.Second).UnixNano())
			require.NoError(t, f.pipe.Reconcile(t.Context()))
			member := f.instant.group.Members[0]
			require.NoError(t, f.pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{
				SampleID: "destination-cutover-after", WorkspaceID: f.f.record.ID, DeviceID: member.DeviceID,
				PointID: member.PointID, TagID: member.TagID, SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
				MeasurementID: f.definition.ID, DefinitionRevision: f.definition.DefinitionRevision,
				SourceBindingRevision: f.definition.SourceBindingRevision, SeriesEpoch: f.definition.SeriesEpoch,
				AcquisitionID: "destination-cutover-after-acquisition", ValueType: string(measurement.ExactUint64),
				ObservedAt: observed, ReceivedAt: observed.Add(time.Millisecond), TimeOrigin: "gateway", Value: uint64(120004), Quality: schema.QualityGood,
			}))
			f.closeAfter(20 * time.Minute)
			require.Eventually(t, func() bool {
				var pending int
				return f.f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state <> 'sql_committed'`).Scan(&pending) == nil && pending == 0
			}, 20*time.Second, 25*time.Millisecond)
			require.NoError(t, f.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+newTable+"` WHERE group_id=? AND group_revision=? AND value_decimal=12000.4 AND observed_at=?",
				f.instant.id, current.GroupRevision, observed).Scan(&count))
			require.Equal(t, 1, count, "the post-boundary raw value reaches only the new version and table")
			require.NoError(t, f.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+f.table+"` WHERE group_id=? AND observed_at>=?", f.instant.id, current.EffectiveAt).Scan(&count))
			require.Zero(t, count, "the old instantaneous destination cannot consume post-boundary input")
			require.NoError(t, f.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+f.table+"` WHERE group_id=? AND value_decimal=12000.4 AND delta_value=0.4", f.cumulative.id).Scan(&count))
			require.Equal(t, 1, count, "the untouched cumulative profile continues to its own original destination")
			afterPayloads := mysqlCutoverFrozenOutbox(t, f)
			for key, digest := range oldPayloads {
				require.Equal(t, digest, afterPayloads[key], "accepted old effect payload remains immutable")
			}
		})
	}
}

// TestMySQLLongLegacyCoexistence verifies saved custom configurations and their
// CAS contract are unaffected by creating a separate MySQL managed profile.
// SQL writes and Modbus-only activation have separate integration tests.
func TestMySQLLongLegacyCoexistence(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newWriteGroupRouterFixture(t)
			if kind == "postgres" {
				_, err := f.db.ExecContext(t.Context(), `UPDATE database_connectors SET kind='postgres', connection_config='{"host":"127.0.0.1","database":"compatibility","schema":"public"}' WHERE id='connector-A'`)
				require.NoError(t, err)
			}
			resp := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
			require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())
			_, existing := groupSaveData(t, decodeJSONBody(t, resp))
			existingID := existing["id"].(string)
			_, err := f.db.ExecContext(t.Context(), `INSERT INTO database_connectors
				(id,name,kind,status,connection_config,identity_revision)
				VALUES ('mysql-compat','MySQL compatibility','mysql','ready','{"host":"127.0.0.1","port":62516,"user":"fixture-only","database":"compatibility"}','mysql-compat-v1')`)
			require.NoError(t, err)
			request := profileRequest(t, f, "instantaneous", 15, "second", "mysql-compat", "")
			request["expected_connector_revision"] = "mysql-compat-v1"
			group := request["group"].(map[string]any)
			group["layout_version"], group["storage_mode"] = "managed_long_v1", "plain"
			destinationOf(group)["connector_id"], destinationOf(group)["connector_revision"] = "mysql-compat", "mysql-compat-v1"
			resp = performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
			require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
			_, managed := groupSaveData(t, decodeJSONBody(t, resp))
			require.Regexp(t, longTablePattern, destinationOf(managed)["table_name"])
			wrapper := managedGroupFixture{writeGroupRouterFixture: f}
			require.Equal(t, existing, getGroup(t, wrapper, existingID))
			stale := groupEditRequest(f, getGroupData(t, wrapper, existingID), existing)
			stale["expected_workspace_revision"] = "stale-compatibility-revision"
			resp = performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+existingID, stale)
			require.Equal(t, http.StatusConflict, resp.Code, resp.Body.String())
			require.Equal(t, existing, getGroup(t, wrapper, existingID))
		})
	}
}

// mysqlCutoverFrozenOutbox captures immutable accepted bytes and destination
// fields, excluding retry state and attempt timestamps that legitimately change.
func mysqlCutoverFrozenOutbox(t *testing.T, f *mysqlCircuitPlainFixture) map[string][6]string {
	t.Helper()
	rows, err := f.f.db.QueryContext(t.Context(), `SELECT effect_key, payload, payload_digest,
		connector_id, connector_revision, table_schema, table_name FROM wg_delivery_outbox
		WHERE group_id IN (?, ?)`, f.instant.id, f.cumulative.id)
	require.NoError(t, err)
	defer rows.Close()
	result := map[string][6]string{}
	for rows.Next() {
		var key string
		var fields [6]string
		require.NoError(t, rows.Scan(&key, &fields[0], &fields[1], &fields[2], &fields[3], &fields[4], &fields[5]))
		result[key] = fields
	}
	require.NoError(t, rows.Err())
	return result
}

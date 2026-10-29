package api

import (
	"context"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStudioV2LocatorEntryProductionRouteReadsExactMappingWithoutBootstrap(t *testing.T) {
	r := newRuntimeSourceRig(t)
	response := locatorRequest(t, r, "3", url.Values{"device_id": {"device-A"}, "mapping_id": {"mapping-source"}})
	data := locatorData(t, response)
	entry, ok := data["locator_entry"].(map[string]any)
	require.True(t, ok, "scoped entry must return revalidated native selection, not legacy readiness bootstrap: %s", response.Body.String())
	require.Equal(t, []any{"mapping-source"}, entry["revalidated_mapping_ids"])
	require.Nil(t, data["readiness_summary"])
	require.NotContains(t, data, "modbus_share")
}

func TestStudioV2LocatorEntryExactStepCollectionsAndNoWrites(t *testing.T) {
	r := newRuntimeSourceRig(t)
	for _, tc := range []struct {
		step  string
		query url.Values
		keys  []string
	}{
		{"1", url.Values{"device_id": {"device-A"}}, []string{"devices"}},
		{"2", url.Values{"device_id": {"device-A"}, "rule_id": {"rule-R"}}, []string{"devices", "source_rules"}},
		{"3", url.Values{"device_id": {"device-A"}, "mapping_id": {"mapping-source"}}, []string{"devices", "source_rules", "mappings"}},
		{"4", url.Values{"group_id": {r.f.group["id"].(string)}, "group_revision": {"historical-hint"}}, []string{"write_groups"}},
		{"4", url.Values{"group_id": {r.f.group["id"].(string)}, "device_id": {"device-A"}}, []string{"devices", "write_groups"}},
	} {
		t.Run(tc.step+tc.query.Encode(), func(t *testing.T) {
			before := locatorTotal(t, r)
			entry := locatorEntryData(t, locatorRequest(t, r, tc.step, tc.query))
			require.Equal(t, before, locatorTotal(t, r))
			for _, key := range []string{"devices", "source_rules", "mappings", "write_groups"} {
				if slices.Contains(tc.keys, key) {
					require.Contains(t, entry, key)
				} else {
					require.NotContains(t, entry, key)
				}
			}
			require.Equal(t, false, entry["hydration_scope"].(map[string]any)["complete_workspace"])
			require.Contains(t, entry, "rule_heads")
			require.Contains(t, entry, "group_heads")
			require.Contains(t, entry, "revalidated_mapping_ids")
			require.NotContains(t, entry, "database_config")
		})
	}
}
func TestStudioV2LocatorEntryStrictQueryNeverFallsBack(t *testing.T) {
	r := newRuntimeSourceRig(t)
	for _, q := range []string{"workspace_id=W", "entry=locator", "entry=locator&workspace_id=%ZZ&step=1&device_id=%ZZ", "entry=locator;workspace_id=W", "entry=locator&workspace_id=W&step=1&device_id=D&device_id=%ZZ", "entry=locator&workspace_id=W&step=1&device_id=D&readiness=1", "entry=locator&workspace_id=W&step=3&device_id=D&mapping_id=M&rule_id=R", "entry=locator&workspace_id=W&step=4&group_id=G&unknown="} {
		t.Run(q, func(t *testing.T) {
			before := locatorTotal(t, r)
			response := performJSONRequest(t, r.router, "GET", studioLocatorPath+"?"+q, nil)
			require.Equal(t, 400, response.Code, response.Body.String())
			require.Equal(t, before, locatorTotal(t, r))
		})
	}
	legacy := performJSONRequest(t, r.router, "GET", studioLocatorPath, nil)
	data := locatorData(t, legacy)
	require.NotContains(t, data, "locator_entry")
	require.NotNil(t, data["readiness_summary"])
}
func TestStudioV2LocatorEntryForeignMissingAndUnrepresentable(t *testing.T) {
	for _, tc := range []struct {
		name, step string
		q          url.Values
		sql        string
		code       int
	}{
		{"foreign workspace", "1", url.Values{"workspace_id": {"foreign"}, "device_id": {"device-A"}}, "", 403},
		{"foreign device", "1", url.Values{"device_id": {"foreign"}}, "", 403},
		{"missing device", "1", url.Values{"device_id": {"device-A"}}, `DELETE FROM devices WHERE id='device-A'`, 409},
		{"missing rule", "2", url.Values{"device_id": {"device-A"}, "rule_id": {"gone"}}, "", 409},
		{"missing mapping", "3", url.Values{"device_id": {"device-A"}, "mapping_id": {"gone"}}, "", 409},
		{"missing group", "4", url.Values{"group_id": {"gone"}}, "", 409},
		{"link mismatch", "3", url.Values{"device_id": {"device-A"}, "mapping_id": {"mapping-source"}}, `UPDATE source_rule_links SET mapping_id='mapping-A' WHERE id='link-source'`, 409},
		{"unsupported pipeline", "3", url.Values{"device_id": {"device-A"}, "mapping_id": {"mapping-source"}}, `UPDATE mappings SET transform_pipeline='[{"type":"secret","params":{"password":"leak"}}]' WHERE id='mapping-source'`, 503},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newRuntimeSourceRig(t)
			if tc.sql != "" {
				locatorExec(t, r, tc.sql)
			}
			res := locatorRequest(t, r, tc.step, tc.q)
			require.Equal(t, tc.code, res.Code, res.Body.String())
			require.NotContains(t, res.Body.String(), "password")
		})
	}
}
func TestStudioV2LocatorEntryIndependentHeadsAndConfigurationRace(t *testing.T) {
	for _, change := range []string{`UPDATE source_rules SET revision_id='R3' WHERE id='rule-R'`, `UPDATE source_rule_links SET mapping_id='mapping-A' WHERE id='link-source'`, `UPDATE points SET address='40009' WHERE id='point-source'`, `UPDATE mappings SET enabled=0 WHERE id='mapping-source'`, `UPDATE tags SET display_name='changed' WHERE id='tag-source'`, `DELETE FROM system_settings WHERE key='studio_v2_workspace'`} {
		t.Run(change, func(t *testing.T) {
			r := newRuntimeSourceRig(t)
			readers := locatorReaders(r)
			calls := 0
			readers.Rules = locatorRuleReads{Service: r.services.SourceRule, get: func(ctx context.Context, id string) (*schema.SourceRule, error) {
				got, err := r.services.SourceRule.GetByID(ctx, id)
				calls++
				if calls == 1 {
					locatorExec(t, r, change)
				}
				return got, err
			}}
			useLocatorReaders(r, readers)
			res := locatorRequest(t, r, "3", url.Values{"device_id": {"device-A"}, "mapping_id": {"mapping-source"}})
			require.Equal(t, 409, res.Code, res.Body.String())
			if strings.Contains(change, "DELETE") {
				var n int
				require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&n))
				require.Zero(t, n)
			}
		})
	}
}
func TestStudioV2LocatorEntryProviderFailureRemainsUnavailable(t *testing.T) {
	r := newRuntimeSourceRig(t)
	readers := locatorReaders(r)
	calls := 0
	readers.Devices = locatorDeviceRead(func(ctx context.Context, id string) (*schema.Device, error) {
		calls++
		if calls == 2 {
			return nil, errors.New("password=secret")
		}
		return r.services.Device.GetByID(ctx, id)
	})
	useLocatorReaders(r, readers)
	res := locatorRequest(t, r, "1", url.Values{"device_id": {"device-A"}})
	require.Equal(t, 503, res.Code, res.Body.String())
	require.NotContains(t, res.Body.String(), "password=secret")
}
func TestStudioV2LocatorEntryEmptyRevisionAndRedactedNativeDiagnostics(t *testing.T) {
	r := newRuntimeSourceRig(t)
	locatorExec(t, r, `UPDATE system_settings SET value=json_set(value,'$.database_setup_revision','') WHERE key='studio_v2_workspace'`)
	locatorExec(t, r, `UPDATE devices SET last_test_error='password=secret' WHERE id='device-A'`)
	res := locatorRequest(t, r, "2", url.Values{"device_id": {"device-A"}, "rule_id": {"rule-R"}})
	entry := locatorEntryData(t, res)
	require.Equal(t, "", entry["workspace_revision"])
	require.Equal(t, []any{}, entry["source_rules"].([]any)[0].(map[string]any)["skipped_addresses"])
	require.NotContains(t, res.Body.String(), "password=secret")
}
func TestStudioV2LocatorEntryEditorRequiresFreshBindingsAndPureHydration(t *testing.T) {
	r := newRuntimeSourceRig(t)
	wireLocatorDatabase(r)
	q := locatorEditorQuery(t, r)
	before := locatorTotal(t, r)
	res := locatorRequest(t, r, "4", q)
	entry := locatorEntryData(t, res)
	require.Equal(t, true, entry["editor_hydrated"])
	require.Contains(t, entry, "database_config")
	require.Contains(t, entry, "database_targets")
	require.Contains(t, entry, "devices")
	require.Contains(t, entry, "source_rules")
	require.Contains(t, entry, "mappings")
	require.Equal(t, before, locatorTotal(t, r))
}

func TestStudioV2LocatorEntryRejectsMissingExactLinkIDs(t *testing.T) {
	for _, column := range []string{"mapping_id", "tag_id"} {
		t.Run(column, func(t *testing.T) {
			r := newRuntimeSourceRig(t)
			locatorExec(t, r, "UPDATE source_rule_links SET "+column+"=NULL WHERE id='link-source'")
			before := locatorTotal(t, r)
			response := locatorRequest(t, r, "3", url.Values{"device_id": {"device-A"}, "mapping_id": {"mapping-source"}})
			require.Equal(t, 409, response.Code, response.Body.String())
			require.Equal(t, before, locatorTotal(t, r))
		})
	}
}
func TestStudioV2LocatorEntryNeverCreatesMissingWorkspace(t *testing.T) {
	r := newRuntimeSourceRig(t)
	locatorExec(t, r, `DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
	before := locatorTotal(t, r)
	res := locatorRequest(t, r, "1", url.Values{"device_id": {"device-A"}})
	require.Equal(t, 409, res.Code, res.Body.String())
	require.Equal(t, before, locatorTotal(t, r))
}
func TestStudioV2LocatorEntryEditorBindingsAndKnownOrphan(t *testing.T) {
	r := newRuntimeSourceRig(t)
	wireLocatorDatabase(r)
	q := locatorEditorQuery(t, r)
	for _, key := range []string{"expected_workspace_revision", "expected_group_revision"} {
		old := q.Get(key)
		q.Del(key)
		require.Equal(t, 400, locatorRequest(t, r, "4", q).Code)
		q.Set(key, "stale")
		require.Equal(t, 409, locatorRequest(t, r, "4", q).Code)
		q.Set(key, old)
	}
	q.Set("group_revision", "historical-only")
	locatorExec(t, r, `DELETE FROM points WHERE id='point-source'`)
	res := locatorRequest(t, r, "4", q)
	entry := locatorEntryData(t, res)
	require.Empty(t, entry["mappings"])
	require.Empty(t, entry["source_rules"])
	groups := entry["write_groups"].(map[string]any)["groups"].([]any)
	require.Len(t, groups[0].(map[string]any)["members"], 1)
}
func TestStudioV2LocatorEntryTombstoneDoesNotUseSameNameReplacement(t *testing.T) {
	r := newRuntimeSourceRig(t)
	draft := r.f.writeGroupRouterFixture.createRequest()
	current, err := r.f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	draft["expected_workspace_revision"] = current.DatabaseSetupRevision
	createdDraft := performJSONRequest(t, r.f.router, "POST", writeGroupsPath, draft)
	require.Equal(t, 201, createdDraft.Code, createdDraft.Body.String())
	_, draftGroup := groupSaveData(t, decodeJSONBody(t, createdDraft))
	id := draftGroup["id"].(string)
	removed := performJSONRequest(t, r.f.router, "DELETE", writeGroupsPath+"/"+id, scopeBody(t, r.f, id))
	require.Equal(t, 200, removed.Code, removed.Body.String())
	res := locatorRequest(t, r, "4", url.Values{"group_id": {id}})
	entry := locatorEntryData(t, res)
	group := entry["write_groups"].(map[string]any)["groups"].([]any)[0].(map[string]any)
	require.Equal(t, id, group["id"])
	require.Equal(t, "deleted", group["status"])
	// A different persisted ID with the same display name is never substituted.
	request := r.f.writeGroupRouterFixture.createRequest()
	record, err := r.f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	request["expected_workspace_revision"] = record.DatabaseSetupRevision
	request["group"].(map[string]any)["name"] = group["name"]
	created := performJSONRequest(t, r.f.router, "POST", writeGroupsPath, request)
	require.Equal(t, 201, created.Code, created.Body.String())
	_, replacement := groupSaveData(t, decodeJSONBody(t, created))
	require.NotEqual(t, id, replacement["id"])
	res = locatorRequest(t, r, "4", url.Values{"group_id": {id}})
	selected := locatorEntryData(t, res)["write_groups"].(map[string]any)["groups"].([]any)[0].(map[string]any)
	require.Equal(t, id, selected["id"])
	require.Equal(t, "deleted", selected["status"])
}

func TestStudioV2LocatorEntryGroupRereadAndEditorDeletionRace(t *testing.T) {
	for _, editor := range []bool{false, true} {
		t.Run(fmt.Sprint(editor), func(t *testing.T) {
			r := newRuntimeSourceRig(t)
			wireLocatorDatabase(r)
			q := url.Values{"group_id": {r.f.group["id"].(string)}}
			if editor {
				q = locatorEditorQuery(t, r)
			}
			readers := locatorReaders(r)
			calls := 0
			readers.Groups = locatorGroupRead{WriteGroupService: r.services.WriteGroups, get: func(ctx context.Context, id string) (*workspace.WriteGroupSaveResult, error) {
				g, err := r.services.WriteGroups.Get(ctx, id)
				calls++
				if calls == 1 {
					if editor {
						locatorExec(t, r, `DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
					} else {
						locatorExec(t, r, `UPDATE write_groups SET revision='changed' WHERE id=?`, id)
					}
				}
				return g, err
			}}
			useLocatorReaders(r, readers)
			res := locatorRequest(t, r, "4", q)
			require.Equal(t, 409, res.Code, res.Body.String())
			if editor {
				var n int
				require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&n))
				require.Zero(t, n)
			}
		})
	}
}
func TestStudioV2LocatorEntryEditorConnectorRaceAndError(t *testing.T) {
	for _, failure := range []bool{false, true} {
		t.Run(fmt.Sprint(failure), func(t *testing.T) {
			r := newRuntimeSourceRig(t)
			wireLocatorDatabase(r)
			q := locatorEditorQuery(t, r)
			readers := locatorReaders(r)
			calls := 0
			readers.Connectors = locatorConnectorRead(func(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
				calls++
				if calls == 2 && failure {
					return nil, errors.New("password=secret")
				}
				c, err := r.services.DBTarget.GetByID(ctx, id)
				if calls == 1 && !failure {
					locatorExec(t, r, `UPDATE database_connectors SET identity_revision='changed' WHERE id=?`, id)
				}
				return c, err
			})
			useLocatorReaders(r, readers)
			res := locatorRequest(t, r, "4", q)
			code := 409
			if failure {
				code = 503
			}
			require.Equal(t, code, res.Code, res.Body.String())
			require.NotContains(t, res.Body.String(), "password=secret")
		})
	}
}
func TestStudioV2LocatorEntryEditorMissingProofReappearsIsStale(t *testing.T) {
	r := newRuntimeSourceRig(t)
	wireLocatorDatabase(r)
	q := locatorEditorQuery(t, r)
	locatorExec(t, r, `DELETE FROM source_rule_links WHERE id='link-source'`)
	readers := locatorReaders(r)
	calls := 0
	readers.Groups = locatorGroupRead{WriteGroupService: r.services.WriteGroups, get: func(ctx context.Context, id string) (*workspace.WriteGroupSaveResult, error) {
		g, err := r.services.WriteGroups.Get(ctx, id)
		calls++
		if calls == 2 {
			locatorExec(t, r, `INSERT INTO source_rule_links (id,rule_id,address,point_id,tag_id,mapping_id) VALUES ('restored-link','rule-R','40001','point-source','tag-source','mapping-source')`)
		}
		return g, err
	}}
	useLocatorReaders(r, readers)
	res := locatorRequest(t, r, "4", q)
	require.Equal(t, 409, res.Code, res.Body.String())
}

func TestStudioV2LocatorEntryEditorCanonicalTargetsDoNotResurrectRawIntent(t *testing.T) {
	r := newLocatorMigratedRig(t)
	id := r.f.group["id"].(string)
	disabled := performJSONRequest(t, r.router, "POST", writeGroupsPath+"/"+id+"/disable", scopeBody(t, r.f, id))
	require.Equal(t, 200, disabled.Code, disabled.Body.String())
	locatorExec(t, r, `INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels) VALUES ('tag-other','other','other','Other','float32','active','{}')`)
	locatorExec(t, r, `INSERT INTO database_target_mappings (id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,enabled) VALUES ('unrelated','tag-other','connector-A','main','other','unrelated_column','insert',1)`)
	before := locatorTotal(t, r)
	entry := locatorEntryData(t, locatorRequest(t, r, "4", locatorEditorQuery(t, r)))
	rows := entry["database_targets"].([]any)
	require.Len(t, rows, 1)
	target := rows[0].(map[string]any)
	require.Equal(t, "legacy-A", target["id"])
	require.Equal(t, "point-A", target["point_id"])
	require.Equal(t, id, target["row_group_id"])
	require.Equal(t, "reviewed_temperature", target["column_name"])
	require.Equal(t, false, target["enabled"])
	require.NotContains(t, target, "canonical_group")
	require.Equal(t, before, locatorTotal(t, r))
	raw, err := r.services.DBMapping.GetByID(t.Context(), "legacy-A")
	require.NoError(t, err)
	require.Equal(t, "temperature", raw.ColumnName)
	require.True(t, raw.Enabled)
	unrelated, err := r.services.DBMapping.GetByID(t.Context(), "unrelated")
	require.NoError(t, err)
	require.Equal(t, "unrelated_column", unrelated.ColumnName)
}
func TestStudioV2LocatorEntryEditorMissingCanonicalLegacyRowIsUnavailable(t *testing.T) {
	r := newLocatorMigratedRig(t)
	q := locatorEditorQuery(t, r)
	require.Equal(t, 200, locatorRequest(t, r, "4", q).Code, "valid canonical target proof must succeed before fault")
	locatorExec(t, r, `DELETE FROM database_target_mappings WHERE id='legacy-A'`)
	res := locatorRequest(t, r, "4", q)
	require.Equal(t, 503, res.Code, res.Body.String())
}

func TestStudioV2LocatorEntrySharedMembersNeverHydrateUnselectedDevices(t *testing.T) {
	delivery := newRuntimeDeliveryRig(t)
	delivery.sharedGroup(t)
	r := &runtimeSourceRig{f: delivery.f, router: delivery.router}
	q := url.Values{"group_id": {r.f.group["id"].(string)}}
	entry := locatorEntryData(t, locatorRequest(t, r, "4", q))
	require.NotContains(t, entry, "devices")
	require.NotContains(t, entry, "source_rules")
	require.Len(t, entry["write_groups"].(map[string]any)["groups"].([]any)[0].(map[string]any)["members"], 2)
	// The saved group retains a former member without exposing its native config.
	locatorExec(t, r, `UPDATE system_settings SET value=json_set(value,'$.ordered_device_ids',json('["device-A"]')) WHERE key='studio_v2_workspace'`)
	entry = locatorEntryData(t, locatorRequest(t, r, "4", q))
	require.Len(t, entry["write_groups"].(map[string]any)["groups"].([]any)[0].(map[string]any)["members"], 2)
	require.NotContains(t, entry, "devices")
	q.Set("device_id", "device-B")
	require.Equal(t, 403, locatorRequest(t, r, "4", q).Code)
}

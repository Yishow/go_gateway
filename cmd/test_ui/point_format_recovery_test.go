package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/api/handlers"
	appconfig "go-gateway/internal/config"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
)

func TestProductionFormatRecoveryWithoutRestart(t *testing.T) {
	for _, tc := range []struct {
		name             string
		enabled          bool
		inactive         schema.DeviceStatus
		changedCandidate bool
	}{{"share-disabled", false, "", false}, {"share-enabled", true, "", false}, {"disabled-device", false, schema.DeviceStatusDisabled, false}, {"draft-device", false, schema.DeviceStatusDraft, false}, {"replay-changed-candidate", true, "", true}} {
		enabled := tc.enabled
		t.Run(tc.name, func(t *testing.T) {
			f := newWordProductionFixture(t, schema.ProtocolModbusTCP)
			require.NoError(t, f.services.modbusShare.ApplySettings(t.Context(), modbusshare.Settings{Enabled: enabled, BindAddress: "127.0.0.1", Port: 15020, SlaveID: 1, CapacityRegisters: 64}))
			rule, p, _ := f.createRule(t, schema.ProtocolModbusTCP, 0, schema.DataTypeUint32, "ABCD")
			if tc.inactive != "" {
				response := httptest.NewRecorder()
				f.router.ServeHTTP(response, httptest.NewRequest("POST", appconfig.Get().API.BasePath+"/datalink/devices/word-device/disable", nil))
				require.Equal(t, 200, response.Code, response.Body.String())
				if tc.inactive == schema.DeviceStatusDraft {
					_, e := f.db.Exec(`UPDATE devices SET status='draft' WHERE id='word-device'`)
					require.NoError(t, e)
				}
			}
			revisions := modbusshare.NewSQLWorkspaceRevisionStore(f.db)
			hydrateStartupShare(t.Context(), f.services, revisions, modbusshare.NewReconciler(f.services.modbusShare, revisions))
			require.NoError(t, f.services.runtime.Start(t.Context()))
			if enabled {
				rule, p = enableRecoveryRuleShare(t, f, rule, p)
			}
			f.router = wordProductionRouter(f.services, f.db)
			record, err := f.services.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)
			base := appconfig.Get().API.BasePath + "/datalink"
			for i, format := range []string{"CDAB", ""} {
				get := httptest.NewRecorder()
				f.router.ServeHTTP(get, httptest.NewRequest("GET", base+"/points/"+p.ID+"/format?workspace_id="+record.ID+"&device_id="+p.DeviceID+"&rule_id="+rule.ID, nil))
				require.Equal(t, 200, get.Code, get.Body.String())
				var read handlers.PointFormatReadResponse
				require.NoError(t, json.Unmarshal(get.Body.Bytes(), &read))
				require.True(t, read.Data.Editable)
				mode := schema.PointFormatOverride
				if format == "" {
					mode = schema.PointFormatInherit
				}
				ids := []string{"recover-override", "recover-reset"}
				save := point.FormatMutationRequest{FormatScope: point.FormatScope{WorkspaceID: record.ID, DeviceID: p.DeviceID, PointID: p.ID, RuleID: rule.ID, ExpectedWorkspaceRevision: read.Data.WorkspaceRevision, ExpectedPointRevision: &read.Data.PointRevision, ExpectedRuleRevision: read.Data.RuleRevision, RequestID: ids[i]}, FormatMode: mode, ConfiguredFormat: format}
				raw, e := json.Marshal(save)
				require.NoError(t, e)
				res := httptest.NewRecorder()
				request := httptest.NewRequest("PUT", base+"/points/"+p.ID+"/format", bytes.NewReader(raw))
				request.Header.Set("Content-Type", "application/json")
				f.router.ServeHTTP(res, request)
				require.Equal(t, 200, res.Code, res.Body.String())
				var saved handlers.PointFormatMutationResponse
				require.NoError(t, json.Unmarshal(res.Body.Bytes(), &saved))
				require.Empty(t, f.services.scheduler.PollNowContext(t.Context(), []string{p.ID}))
				recovery := point.FormatRecoveryRequest{FormatScope: save.FormatScope, ExpectedConnectionRevision: read.Data.ConnectionRevision, ExpectedSettingsRevision: f.services.modbusShare.Settings().SettingsRevision}
				recovery.ExpectedPointRevision = &saved.Data.Point.RevisionID
				recovery.ExpectedWorkspaceRevision = saved.Data.WorkspaceRevision
				raw, e = json.Marshal(map[string]any{"format_recovery": recovery})
				require.NoError(t, e)
				var first handlers.FormatRecoveryResult
				var firstLease uint64
				for attempt := 0; attempt < 2; attempt++ {
					if attempt == 1 && tc.changedCandidate {
						_, e := f.db.Exec(`UPDATE source_rule_candidate_snapshots SET payload=json_set(payload,'$.candidates[0].register',10) WHERE source_rule_id=? AND candidate_type='local_modbus_outputs'`, rule.ID)
						require.NoError(t, e)
						desired, e := f.services.sourceRule.BuildDesiredShareMappingsForDevices(t.Context(), record.ID, f.services.modbusShare.Settings(), []string{p.DeviceID})
						require.NoError(t, e)
						require.Len(t, desired, 1)
						require.Equal(t, uint16(10), desired[0].ZeroBasedRegister)
						require.Equal(t, uint16(0), f.services.modbusShare.ListMappings()[0].ZeroBasedRegister)
					} else if attempt == 1 && tc.inactive == "" {
						poll := httptest.NewRecorder()
						f.router.ServeHTTP(poll, httptest.NewRequest("POST", base+"/points/"+p.ID+"/poll", nil))
						require.Equal(t, 200, poll.Code, poll.Body.String())
						telemetry, e := f.services.point.GetByID(t.Context(), p.ID)
						require.NoError(t, e)
						require.NotNil(t, telemetry.LastReadAt)
						require.NotNil(t, telemetry.LastValue)
					}

					res = httptest.NewRecorder()
					request = httptest.NewRequest("POST", base+"/modbus-share/reconcile", bytes.NewReader(raw))
					request.Header.Set("Content-Type", "application/json")
					f.router.ServeHTTP(res, request)
					if attempt == 1 && tc.changedCandidate {
						require.Equal(t, 409, res.Code, res.Body.String())
						return
					}
					require.Equal(t, 200, res.Code, res.Body.String())
					evidence, e := f.services.point.ReadFormatRecovery(t.Context(), recovery)
					require.NoError(t, e)
					require.Equal(t, rule.ID, evidence.Profile.IdentityView().RuleID)
					if tc.inactive == "" {
						canonical, e := f.services.runtime.ResolvePointReadProfile(t.Context(), &evidence.Point)
						require.NoError(t, e)
						require.Equal(t, canonical.IdentityView(), evidence.Profile.IdentityView(), "recovery must retain ordinary source identity")
					}
					var body struct {
						Data handlers.FormatRecoveryResult `json:"data"`
					}
					require.NoError(t, json.Unmarshal(res.Body.Bytes(), &body))
					require.Equal(t, enabled, body.Data.ShareApplied)
					if tc.inactive == "" {
						require.Equal(t, "installed", body.Data.RuntimeState)
					} else {
						require.Equal(t, "not_started", body.Data.RuntimeState)
						require.Zero(t, f.services.runtime.CapturePointLease(p.ID, saved.Data.Point.RevisionID))
					}
					if attempt == 0 {
						if enabled {
							require.Len(t, f.services.modbusShare.ListMappings(), 1, "whole canonical Share plan must be restored, not only the target point")
						}
						first = body.Data
						firstLease = f.services.runtime.CapturePointLease(p.ID, saved.Data.Point.RevisionID)
					} else {
						require.Equal(t, first.Hydration.ReadinessToken, body.Data.Hydration.ReadinessToken)
						require.Equal(t, firstLease, f.services.runtime.CapturePointLease(p.ID, saved.Data.Point.RevisionID), "replay must not reinstall or invent an inactive lease")
						require.Equal(t, "replayed", body.Data.Outcome)
					}
				}
				words := []uint16{0x3344, 0x1122}
				if i == 1 {
					words = []uint16{0x1122, 0x3344}
				}
				f.plc.set(0, words)
				values := f.services.scheduler.PollNowContext(t.Context(), []string{p.ID})
				if tc.inactive == "" {
					require.Len(t, values, 1)
					require.Equal(t, uint32(0x11223344), values[0].Value)
				} else {
					require.Empty(t, values)
					device, e := f.services.device.GetByID(t.Context(), p.DeviceID)
					require.NoError(t, e)
					require.Equal(t, tc.inactive, device.Status)
				}
			}
		})
	}
}

func enableRecoveryRuleShare(t *testing.T, f wordProductionFixture, rule *schema.SourceRule, p *schema.Point) (*schema.SourceRule, *schema.Point) {
	t.Helper()
	enabled, start, stride := true, 40001, 2
	var err error
	rule, err = f.services.sourceRule.Update(t.Context(), rule.ID, sourcerule.UpdateRuleRequest{ShareEnabled: &enabled, ShareStartRegister: &start, ShareStride: &stride})
	require.NoError(t, err)
	snapshots, err := f.services.sourceRule.ListCandidateSnapshots(t.Context(), rule.ID)
	require.NoError(t, err)
	var ids []string
	for _, snapshot := range snapshots {
		if snapshot.CandidateType == schema.SourceRuleCandidateTypeTags {
			var payload struct {
				Candidates []schema.SourceRuleTagCandidate `json:"candidates"`
			}
			require.NoError(t, json.Unmarshal([]byte(snapshot.Payload), &payload))
			for _, c := range payload.Candidates {
				ids = append(ids, c.ID)
			}
		}
	}
	require.NotEmpty(t, ids)
	hydration, err := f.services.modbusShare.CheckHydration(t.Context())
	require.NoError(t, err)
	applied, err := f.services.sourceRule.ApplyTagCandidates(t.Context(), rule.ID, sourcerule.ApplyTagCandidatesRequest{WorkspaceID: hydration.WorkspaceID, ExpectedWorkspaceRevision: hydration.WorkspaceRevision, RevisionID: rule.RevisionID, CandidateIDs: ids})
	require.NoError(t, err)
	require.Len(t, applied.Results, len(ids))
	for _, result := range applied.Results {
		require.NotEqual(t, "failed", result.Status, result.Error)
	}
	p, err = f.services.point.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	return rule, p
}

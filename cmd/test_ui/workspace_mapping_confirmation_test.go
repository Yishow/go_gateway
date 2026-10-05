package main

import (
	"go-gateway/internal/api"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProductionWorkspaceConfirmationOnlySuccessfulOwnedSave(t *testing.T) {
	for _, scenario := range []string{"failed-confirmation", "source-edit", "same-rule-other-mapping"} {
		t.Run(scenario, func(t *testing.T) {
			env := newOutageEnv(t)
			_, err := env.db.ExecContext(t.Context(), `UPDATE devices SET status='active',last_test_success=1,last_test_error='',readiness_status='{"probe_status":"success"}',description='',created_at='2026-10-05 00:00:00',updated_at='2026-10-05 00:00:00' WHERE id='device-1'`)
			require.NoError(t, err)
			rule, err := env.services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{ID: "accept-source", DeviceID: "device-1", StartAddress: "42001", Count: 2, DataType: schema.DataTypeInt16, NamingPrefix: "Accept_", Enabled: true})
			require.NoError(t, err)
			links, err := env.services.sourceRule.ListLinks(t.Context(), rule.ID)
			require.NoError(t, err)
			require.Len(t, links, 2)
			router := api.NewRouter(&api.DatalinkServices{Workspace: env.services.workspace, Device: env.services.device, SourceRule: env.services.sourceRule, Point: env.services.point, Tag: env.services.tag, Mapping: env.services.mapping})
			requests := make([]map[string]any, 0, len(links))
			records := make([]*schema.Mapping, 0, len(links))
			for i, link := range links {
				request := map[string]any{"point_id": link.PointID, "rule_id": rule.ID, "device_id": "device-1", "address": link.Address, "tag_key": "accept." + link.Address, "display_name": "Point", "unit": "", "target_type": "int16", "scale": 1, "offset": 0, "enabled": true}
				status, body := httpJSON(t, router, http.MethodPost, testWriteBase+"/mappings", request)
				require.Equal(t, http.StatusCreated, status, "%+v", body)
				record, e := env.services.mapping.GetByID(t.Context(), body["data"].(map[string]any)["id"].(string))
				require.NoError(t, e)
				require.Equal(t, schema.MappingStatusActive, record.Status)
				requests = append(requests, request)
				records = append(records, record)
				require.Equal(t, i+1, len(records))
			}
			first, second := records[0], records[1]
			requests[0]["target_type"] = "float64"
			requests[0]["scale"] = 0.5
			requests[0]["offset"] = 10
			switch scenario {
			case "failed-confirmation":
				_, err = env.db.ExecContext(t.Context(), `CREATE TRIGGER fail_accept BEFORE UPDATE OF last_applied_signature ON mappings WHEN NEW.id='`+first.ID+`' AND NEW.last_applied_signature<>OLD.last_applied_signature BEGIN SELECT RAISE(ABORT,'owned injected confirmation failure'); END`)
				require.NoError(t, err)
				status, _ := httpJSON(t, router, http.MethodPut, testWriteBase+"/mappings/"+first.ID, requests[0])
				require.Equal(t, http.StatusInternalServerError, status)
				after, e := env.services.mapping.GetByID(t.Context(), first.ID)
				require.NoError(t, e)
				require.Equal(t, first.LastAppliedSignature, after.LastAppliedSignature)
				require.False(t, after.Enabled, "failed save cannot publish the pending pipeline")
				require.Equal(t, schema.MappingStatusDraft, after.Status)
			case "source-edit":
				multiplier := 2.0
				target := schema.DataTypeFloat64
				_, err = env.services.sourceRule.Update(t.Context(), rule.ID, sourcerule.UpdateRuleRequest{ScaleMultiplier: &multiplier, ScaleMultiplierSet: true, TargetDataType: &target})
				require.NoError(t, err)
				status, _ := httpJSON(t, router, http.MethodPut, testWriteBase+"/mappings/"+first.ID, requests[0])
				require.Equal(t, http.StatusConflict, status)
				after, e := env.services.mapping.GetByID(t.Context(), first.ID)
				require.NoError(t, e)
				require.Equal(t, first.LastAppliedSignature, after.LastAppliedSignature)
				require.Equal(t, first.TransformPipeline, after.TransformPipeline)
				require.Equal(t, schema.MappingStatusOutOfSync, after.Status)
				require.NotEmpty(t, after.BlockingReason)
				require.ErrorIs(t, env.services.sourceRule.ConfirmWorkspaceMapping(t.Context(), rule.ID, rule.RevisionID, first, schema.DataTypeInt16), mapping.ErrWorkspaceConfirmationConflict)
			case "same-rule-other-mapping":
				_, err = env.services.mapping.Update(t.Context(), second.ID, mapping.UpdateMappingRequest{TransformPipeline: []schema.TransformStep{{Type: schema.TransformScale, Params: map[string]any{"scale": 3.0, "offset": 0.0}}}})
				require.NoError(t, err)
				status, body := httpJSON(t, router, http.MethodPut, testWriteBase+"/mappings/"+first.ID, requests[0])
				require.Equal(t, http.StatusOK, status, "%+v", body)
				updated, e := env.services.mapping.GetByID(t.Context(), first.ID)
				require.NoError(t, e)
				require.Equal(t, schema.MappingStatusActive, updated.Status)
				require.NotEqual(t, first.LastAppliedSignature, updated.LastAppliedSignature)
				other, e := env.services.mapping.GetByID(t.Context(), second.ID)
				require.NoError(t, e)
				require.Equal(t, schema.MappingStatusOutOfSync, other.Status)
				require.Equal(t, second.LastAppliedSignature, other.LastAppliedSignature)
				status, _ = httpJSON(t, router, http.MethodPut, testWriteBase+"/mappings/"+second.ID, requests[1])
				require.Equal(t, http.StatusConflict, status, "unconfirmed manual edit cannot bypass candidate review")
				require.NoError(t, env.services.sourceRule.SyncDerivedPointState(t.Context()))
				again, e := env.services.mapping.GetByID(t.Context(), first.ID)
				require.NoError(t, e)
				require.Equal(t, updated.LastAppliedSignature, again.LastAppliedSignature)
				require.Equal(t, schema.MappingStatusActive, again.Status)
			}
		})
	}
}

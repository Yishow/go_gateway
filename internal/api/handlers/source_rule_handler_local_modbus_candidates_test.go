package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestSourceRuleHTTP_ShareConfigAndTagApplyProduceReadyLocalModbusCandidate(t *testing.T) {
	fixture := setupSourceRuleCandidatesFixture(t)

	createBody := []byte(`{"id":"rule-http-share-ready","device_id":"device-1","start_address":"40001","count":2,"data_type":"int16","naming_prefix":"HTTP_SHARE","enabled":true}`)
	createReq := newHandlerTestRequest(http.MethodPost, "/datalink/source-rules", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(createResp, createReq)
	require.Equal(t, http.StatusCreated, createResp.Code, createResp.Body.String())

	var created struct {
		Data struct {
			ID         string `json:"id"`
			RevisionID string `json:"revision_id"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(createResp.Body.Bytes(), &created))
	require.NotEmpty(t, created.Data.ID)
	require.NotEmpty(t, created.Data.RevisionID)

	updateReq := newHandlerTestRequest(http.MethodPut, "/datalink/source-rules/"+created.Data.ID, bytes.NewBufferString(`{"share_enabled":true,"share_start_register":40001,"share_stride":2}`))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(updateResp, updateReq)
	require.Equal(t, http.StatusOK, updateResp.Code, updateResp.Body.String())
	var updated struct {
		Data struct {
			RevisionID string `json:"revision_id"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(updateResp.Body.Bytes(), &updated))
	require.NotEmpty(t, updated.Data.RevisionID)

	candidatesReq := newHandlerTestRequest(http.MethodGet, "/datalink/source-rules/"+created.Data.ID+"/candidates", http.NoBody)
	candidatesResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(candidatesResp, candidatesReq)
	require.Equal(t, http.StatusOK, candidatesResp.Code)
	var candidates struct {
		Data struct {
			Tags struct {
				Candidates []struct {
					ID string `json:"id"`
				} `json:"candidates"`
			} `json:"tags"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(candidatesResp.Body.Bytes(), &candidates))
	require.Len(t, candidates.Data.Tags.Candidates, 2)
	tagCandidateIDs := []string{candidates.Data.Tags.Candidates[0].ID, candidates.Data.Tags.Candidates[1].ID}

	applyBody, err := json.Marshal(map[string]any{
		"revision_id":   updated.Data.RevisionID,
		"candidate_ids": tagCandidateIDs,
	})
	require.NoError(t, err)
	applyReq := newHandlerTestRequest(http.MethodPost, "/datalink/source-rules/"+created.Data.ID+"/tags/apply", bytes.NewReader(applyBody))
	applyReq.Header.Set("Content-Type", "application/json")
	applyResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(applyResp, applyReq)
	require.Equal(t, http.StatusOK, applyResp.Code, applyResp.Body.String())

	refreshReq := newHandlerTestRequest(http.MethodGet, "/datalink/source-rules/"+created.Data.ID+"/candidates", http.NoBody)
	refreshResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(refreshResp, refreshReq)
	require.Equal(t, http.StatusOK, refreshResp.Code)
	var refreshed struct {
		Data struct {
			LocalModbusOutputs struct {
				Status     schema.SourceRuleCandidateStatus `json:"status"`
				Candidates []struct {
					Register *uint16                                  `json:"register"`
					Status   schema.SourceRuleLocalModbusOutputStatus `json:"status"`
				} `json:"candidates"`
			} `json:"local_modbus_outputs"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(refreshResp.Body.Bytes(), &refreshed))
	require.Equal(t, schema.SourceRuleCandidateStatusReady, refreshed.Data.LocalModbusOutputs.Status)
	require.Len(t, refreshed.Data.LocalModbusOutputs.Candidates, 2)
	for index, expectedRegister := range []uint16{0, 2} {
		require.NotNil(t, refreshed.Data.LocalModbusOutputs.Candidates[index].Register)
		require.Equal(t, expectedRegister, *refreshed.Data.LocalModbusOutputs.Candidates[index].Register)
		require.Equal(t, schema.SourceRuleLocalModbusOutputStatusReady, refreshed.Data.LocalModbusOutputs.Candidates[index].Status)
	}
}

func TestSourceRuleHTTP_MultipleRulesKeepServerOwnedRegisterGeometry(t *testing.T) {
	fixture := setupSourceRuleCandidatesFixture(t)

	first := createAndApplyHTTPShareRule(t, fixture, "rule-http-share-first", "40001", 2, 40001, 2)
	second := createAndApplyHTTPShareRule(t, fixture, "rule-http-share-second", "40001", 1, 40011, 3)

	require.Equal(t, schema.SourceRuleCandidateStatusReady, first.Status)
	require.Len(t, first.Candidates, 2)
	require.Equal(t, []uint16{0, 2}, candidateRegisters(t, first.Candidates))
	require.Equal(t, schema.SourceRuleCandidateStatusReady, second.Status)
	require.Len(t, second.Candidates, 1)
	require.Equal(t, []uint16{10}, candidateRegisters(t, second.Candidates))
}

type localModbusHTTPCandidate struct {
	Register *uint16                                  `json:"register"`
	Status   schema.SourceRuleLocalModbusOutputStatus `json:"status"`
}

type localModbusHTTPCandidates struct {
	Status     schema.SourceRuleCandidateStatus `json:"status"`
	Candidates []localModbusHTTPCandidate       `json:"candidates"`
}

func createAndApplyHTTPShareRule(
	t *testing.T,
	fixture *sourceRuleCandidatesFixture,
	id, startAddress string,
	count, shareStart, shareStride int,
) localModbusHTTPCandidates {
	t.Helper()
	createBody, err := json.Marshal(map[string]any{
		"id":                   id,
		"device_id":            "device-1",
		"start_address":        startAddress,
		"count":                count,
		"data_type":            "int16",
		"naming_prefix":        id,
		"enabled":              true,
		"share_enabled":        true,
		"share_start_register": shareStart,
		"share_stride":         shareStride,
	})
	require.NoError(t, err)
	createReq := newHandlerTestRequest(http.MethodPost, "/datalink/source-rules", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(createResp, createReq)
	require.Equal(t, http.StatusCreated, createResp.Code, createResp.Body.String())

	var created struct {
		Data struct {
			ID         string `json:"id"`
			RevisionID string `json:"revision_id"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(createResp.Body.Bytes(), &created))

	tagIDs := loadHTTPTagCandidateIDs(t, fixture, created.Data.ID)
	applyBody, err := json.Marshal(map[string]any{
		"revision_id":   created.Data.RevisionID,
		"candidate_ids": tagIDs,
	})
	require.NoError(t, err)
	applyReq := newHandlerTestRequest(http.MethodPost, "/datalink/source-rules/"+created.Data.ID+"/tags/apply", bytes.NewReader(applyBody))
	applyReq.Header.Set("Content-Type", "application/json")
	applyResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(applyResp, applyReq)
	require.Equal(t, http.StatusOK, applyResp.Code, applyResp.Body.String())

	refreshReq := newHandlerTestRequest(http.MethodGet, "/datalink/source-rules/"+created.Data.ID+"/candidates", http.NoBody)
	refreshResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(refreshResp, refreshReq)
	require.Equal(t, http.StatusOK, refreshResp.Code, refreshResp.Body.String())
	var refreshed struct {
		Data struct {
			LocalModbusOutputs localModbusHTTPCandidates `json:"local_modbus_outputs"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(refreshResp.Body.Bytes(), &refreshed))
	return refreshed.Data.LocalModbusOutputs
}

func loadHTTPTagCandidateIDs(t *testing.T, fixture *sourceRuleCandidatesFixture, ruleID string) []string {
	t.Helper()
	candidatesReq := newHandlerTestRequest(http.MethodGet, "/datalink/source-rules/"+ruleID+"/candidates", http.NoBody)
	candidatesResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(candidatesResp, candidatesReq)
	require.Equal(t, http.StatusOK, candidatesResp.Code, candidatesResp.Body.String())
	var candidates struct {
		Data struct {
			Tags struct {
				Candidates []struct {
					ID string `json:"id"`
				} `json:"candidates"`
			} `json:"tags"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(candidatesResp.Body.Bytes(), &candidates))
	ids := make([]string, 0, len(candidates.Data.Tags.Candidates))
	for _, candidate := range candidates.Data.Tags.Candidates {
		ids = append(ids, candidate.ID)
	}
	require.NotEmpty(t, ids)
	return ids
}

func candidateRegisters(t *testing.T, candidates []localModbusHTTPCandidate) []uint16 {
	t.Helper()
	registers := make([]uint16, 0, len(candidates))
	for _, candidate := range candidates {
		require.NotNil(t, candidate.Register)
		require.Equal(t, schema.SourceRuleLocalModbusOutputStatusReady, candidate.Status)
		registers = append(registers, *candidate.Register)
	}
	return registers
}

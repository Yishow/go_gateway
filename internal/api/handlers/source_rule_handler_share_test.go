package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sourceRuleSharePayload struct {
	Success bool `json:"success"`
	Data    struct {
		ShareEnabled       bool `json:"share_enabled"`
		ShareStartRegister *int `json:"share_start_register"`
		ShareStride        *int `json:"share_stride"`
	} `json:"data"`
}

func TestSourceRuleHandler_ShareConfigRoundTripsCreateGetListUpdate(t *testing.T) {
	router, _ := setupSourceRuleRouter(t)

	createReq := httptest.NewRequest(http.MethodPost, "/datalink/source-rules", bytes.NewBufferString(`{
		"id":"rule-share-handler",
		"device_id":"device-1",
		"start_address":"40001",
		"count":1,
		"data_type":"int16",
		"naming_prefix":"SRC",
		"enabled":true,
		"share_enabled":true,
		"share_start_register":40001,
		"share_stride":2
	}`))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	router.ServeHTTP(createResp, createReq)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var created sourceRuleSharePayload
	require.NoError(t, json.Unmarshal(createResp.Body.Bytes(), &created))
	require.True(t, created.Success)
	assert.True(t, created.Data.ShareEnabled)
	require.NotNil(t, created.Data.ShareStartRegister)
	assert.Equal(t, 40001, *created.Data.ShareStartRegister)
	require.NotNil(t, created.Data.ShareStride)
	assert.Equal(t, 2, *created.Data.ShareStride)

	getReq := httptest.NewRequest(http.MethodGet, "/datalink/source-rules/rule-share-handler", nil)
	getResp := httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)
	require.Equal(t, http.StatusOK, getResp.Code)
	var fetched sourceRuleSharePayload
	require.NoError(t, json.Unmarshal(getResp.Body.Bytes(), &fetched))
	assert.True(t, fetched.Data.ShareEnabled)
	require.NotNil(t, fetched.Data.ShareStartRegister)
	assert.Equal(t, 40001, *fetched.Data.ShareStartRegister)
	require.NotNil(t, fetched.Data.ShareStride)
	assert.Equal(t, 2, *fetched.Data.ShareStride)

	listReq := httptest.NewRequest(http.MethodGet, "/datalink/source-rules?device_id=device-1", nil)
	listResp := httptest.NewRecorder()
	router.ServeHTTP(listResp, listReq)
	require.Equal(t, http.StatusOK, listResp.Code)
	// The list data is an array of response objects rather than wrapped payloads.
	var listEnvelope struct {
		Data []struct {
			ShareEnabled       bool `json:"share_enabled"`
			ShareStartRegister *int `json:"share_start_register"`
			ShareStride        *int `json:"share_stride"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(listResp.Body.Bytes(), &listEnvelope))
	require.Len(t, listEnvelope.Data, 1)
	assert.True(t, listEnvelope.Data[0].ShareEnabled)
	require.NotNil(t, listEnvelope.Data[0].ShareStartRegister)
	assert.Equal(t, 40001, *listEnvelope.Data[0].ShareStartRegister)
	require.NotNil(t, listEnvelope.Data[0].ShareStride)
	assert.Equal(t, 2, *listEnvelope.Data[0].ShareStride)

	updateReq := httptest.NewRequest(http.MethodPut, "/datalink/source-rules/rule-share-handler", bytes.NewBufferString(`{
		"share_enabled":false,
		"share_start_register":null,
		"share_stride":null
	}`))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	router.ServeHTTP(updateResp, updateReq)
	require.Equal(t, http.StatusOK, updateResp.Code)
	var updated sourceRuleSharePayload
	require.NoError(t, json.Unmarshal(updateResp.Body.Bytes(), &updated))
	assert.False(t, updated.Data.ShareEnabled)
	assert.Nil(t, updated.Data.ShareStartRegister)
	assert.Nil(t, updated.Data.ShareStride)

	finalGetReq := httptest.NewRequest(http.MethodGet, "/datalink/source-rules/rule-share-handler", nil)
	finalGetResp := httptest.NewRecorder()
	router.ServeHTTP(finalGetResp, finalGetReq)
	require.Equal(t, http.StatusOK, finalGetResp.Code)
	var finalPayload sourceRuleSharePayload
	require.NoError(t, json.Unmarshal(finalGetResp.Body.Bytes(), &finalPayload))
	assert.False(t, finalPayload.Data.ShareEnabled)
	assert.Nil(t, finalPayload.Data.ShareStartRegister)
	assert.Nil(t, finalPayload.Data.ShareStride)
}

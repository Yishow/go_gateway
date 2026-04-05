package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourceRuleHandler_Update_ReturnsNewRevisionID(t *testing.T) {
	router, _ := setupSourceRuleRouter(t)

	createBody, err := json.Marshal(sourcerule.CreateRuleRequest{
		ID:           "rule-revision",
		DeviceID:     "device-1",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	createReq, err := http.NewRequest(http.MethodPost, "/datalink/source-rules", bytes.NewBuffer(createBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	router.ServeHTTP(createResp, createReq)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var createPayload map[string]any
	require.NoError(t, json.Unmarshal(createResp.Body.Bytes(), &createPayload))
	createData := createPayload["data"].(map[string]any)
	initialRevision, ok := createData["revision_id"].(string)
	require.True(t, ok)
	require.NotEmpty(t, initialRevision)

	updateReq, err := http.NewRequest(http.MethodPut, "/datalink/source-rules/rule-revision", bytes.NewBufferString(`{"naming_prefix":"UPDATED"}`))
	require.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	router.ServeHTTP(updateResp, updateReq)
	require.Equal(t, http.StatusOK, updateResp.Code)

	var updatePayload map[string]any
	require.NoError(t, json.Unmarshal(updateResp.Body.Bytes(), &updatePayload))
	updateData := updatePayload["data"].(map[string]any)
	nextRevision, ok := updateData["revision_id"].(string)
	require.True(t, ok)
	require.NotEmpty(t, nextRevision)
	assert.NotEqual(t, initialRevision, nextRevision)
}

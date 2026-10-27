package api

import (
	"net/http"
	"testing"

	"go-gateway/internal/datalink/grouptestwrite"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupConfirmCleanupRepairRoutePassesAction(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	writer := &fakeTestWriter{
		outcome: &grouptestwrite.Outcome{
			Claim:     recordingplan.ClaimCompleted,
			Operation: finishedOperation("written_verified", "cleaned"),
		},
	}
	router := f.routerWithTestWriter(writer)

	// POST /studio-v2/workspace/write-groups/:id/test-write with action: "cleanup_repair"
	resp := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/g-1/test-write", map[string]any{
		"token":        "tok-1",
		"operation_id": "op-1",
		"action":       "cleanup_repair",
	})
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, 1, writer.confirms)
	require.Equal(t, "tok-1", writer.gotConfirm.Token)
	require.Equal(t, "op-1", writer.gotConfirm.OperationID)
	require.Equal(t, "g-1", writer.gotConfirm.GroupID)
	require.Equal(t, "cleanup_repair", writer.gotConfirm.Action)
}

func TestWriteGroupConfirmRejectsUnsupportedAction(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	writer := &fakeTestWriter{
		err: &grouptestwrite.UnsupportedError{Reason: "unknown test write action: invalid_act"},
	}
	router := f.routerWithTestWriter(writer)

	resp := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/g-1/test-write", map[string]any{
		"token":        "tok-1",
		"operation_id": "op-1",
		"action":       "invalid_act",
	})
	require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	body := decodeJSONBody(t, resp)
	errEnvelope := body["error"].(map[string]any)
	require.Equal(t, "WRITE_GROUP_TEST_WRITE_UNSUPPORTED", errEnvelope["code"])
}

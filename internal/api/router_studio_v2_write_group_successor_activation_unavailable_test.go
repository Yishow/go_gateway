package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

// Never-applied generations now have an explicit first-start action. The old
// activation-confirm contract remains cutover-only and must still fail closed.
func TestWriteGroupSuccessorActivationNeverAppliedRequiresFirstStart(t *testing.T) {
	r := newFirstStartAPIRig(t, "sqlite", false, false, false, workspace.WriteGroupStorageModePlain)
	request := r.confirmation(t)
	require.Equal(t, "first-start-confirm", request["action"])
	request["action"] = "activation-confirm"
	before := successorAPIAllState(t, r.f.db)
	rejected := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
	require.Equal(t, http.StatusBadRequest, rejected.Code, rejected.Body.String())
	require.Equal(t, before, successorAPIAllState(t, r.f.db), "no fabricated old version, end, fence or activation event")
}

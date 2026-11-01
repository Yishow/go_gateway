package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessorIdentityAndScopeFixedWireGolden(t *testing.T) {
	for _, tc := range []struct{ name, workspace, request, operationWant, scopeWant string }{
		{"empty", "", "", "successor_reserve-0ac2cee4216c7d3a225765ffd1d6ee3048458e106adc209efe8a92248bb321f2", "successor_reserve-8f3d1fafcce4a1d58e20e13b6f2792162ce77407fb81167a6931f2bef07f301f"},
		{"unicode", "工作區🙂", "請求", "successor_reserve-d4a6887b512a7394603dec930bfa3ef3350a8078373897216b27cc20c1f687c6", "successor_reserve-aceb8615c675efbf19a480d9500ea9654583efa4f3ebfa0cde4ce27cc245781f"},
		{"escaped", "\"\\\n<>&", "\u0000\t\r", "successor_reserve-3dc0894d6efb4fd35cd3f1d6ba73f6e8e59d27bf66ffe79a54720ff94cffdb89", "successor_reserve-1010f2f1a4672cd38d6d544f5354082518c55f94eb25387b79356650f5353eab"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.operationWant, successorOperationID(SuccessorReserveAction+"/v1", tc.workspace, tc.request))
			detail := successorOperationDetail{Action: SuccessorReserveAction, Scope: successorStoredScope{WorkspaceID: tc.workspace, RequestID: tc.request, LineageID: "lineage", PredecessorGroupID: "old", SuccessorGroupID: "new"}}
			require.Equal(t, tc.scopeWant, detail.scopeKey())
		})
	}
}

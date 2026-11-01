package point

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionMutationPrecedesPointGroupAssignmentWrites(t *testing.T) {
	for _, operation := range []string{"assign", "remove"} {
		t.Run(operation, func(t *testing.T) {
			lifecycle := &pointAcquisitionLifecycle{t: t}
			repo := &pointAcquisitionRepository{MemoryRepository: NewMemoryRepository(), beforeWrite: func() {
				require.True(t, lifecycle.active, "group membership write must already be quarantined")
			}}
			groups := NewMemoryPollingGroupRepository()
			require.NoError(t, groups.Create(t.Context(), &schema.PollingGroup{ID: "group", Name: "group", Enabled: true}))
			service := NewService(repo, groups)
			service.SetAcquisitionMutationLifecycle(lifecycle)
			groupID := "group"
			p := &schema.Point{ID: "p", DeviceID: "device", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &groupID, Enabled: true}
			require.NoError(t, repo.MemoryRepository.Create(t.Context(), p))
			var err error
			if operation == "assign" {
				err = service.AssignToGroup(t.Context(), p.ID, groupID)
			} else {
				err = service.RemoveFromGroup(t.Context(), p.ID)
			}
			require.NoError(t, err)
			require.Equal(t, 1, lifecycle.begun)
			require.Equal(t, 1, lifecycle.finished)
		})
	}
}

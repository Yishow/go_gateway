package sourcerule

import (
	"errors"
	"testing"

	"go-gateway/internal/datalink/point"

	"github.com/stretchr/testify/require"
)

func TestPointRollbackContinuesAndSyncsOnlyPersistedRestorations(t *testing.T) {
	restoreErr := errors.New("point storage unavailable")
	pointRepo := &rollbackPointRepository{
		Repository:     point.NewMemoryRepository(),
		updateFailures: map[int]error{1: restoreErr},
	}
	syncer := &stubRuntimeSync{}
	svc, rule := newRollbackUpdateService(t, pointRepo, NewMemoryRepository(), syncer)
	links, err := svc.ListLinks(t.Context(), rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)
	first, err := svc.pointSvc.GetByID(t.Context(), links[0].PointID)
	require.NoError(t, err)
	second := *first
	second.ID = "rollback-second-point"
	second.Address = "40002"
	second.Name = "second-original"
	require.NoError(t, pointRepo.Create(t.Context(), &second))
	for _, original := range []pointUpdatePlan{{point: first}, {point: &second}} {
		changed := *original.point
		changed.Name = "changed-" + changed.ID
		changed.DataFormat = "CDAB"
		require.NoError(t, pointRepo.Repository.Update(t.Context(), &changed))
	}
	syncer.upserted = nil

	err = svc.rollbackUpdatedPoints(t.Context(), []pointUpdatePlan{{point: first}, {point: &second}})

	require.ErrorIs(t, err, restoreErr)
	require.Equal(t, []string{second.ID}, syncer.upserted)
	unrestored, err := pointRepo.GetByID(t.Context(), first.ID)
	require.NoError(t, err)
	require.Equal(t, "changed-"+first.ID, unrestored.Name)
	require.Equal(t, "CDAB", unrestored.DataFormat)
	restored, err := pointRepo.GetByID(t.Context(), second.ID)
	require.NoError(t, err)
	require.Equal(t, second.Name, restored.Name)
	require.Equal(t, "ABCD", restored.DataFormat)
}

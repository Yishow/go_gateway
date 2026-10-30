package point

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

type restoreWithoutCAS struct{ Repository }

func TestRestoreRevisionRequiresOwnedAtomicCAS(t *testing.T) {
	ctx := t.Context()
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)
	p, err := svc.Create(ctx, CreatePointRequest{DeviceID: "d", Name: "p", Address: "40001", DataType: schema.DataTypeInt32})
	require.NoError(t, err)
	require.ErrorIs(t, svc.RestoreRevision(ctx, nil, "revision"), ErrFormatUnavailable)
	require.ErrorIs(t, svc.RestoreRevision(ctx, p, ""), ErrFormatUnavailable)
	require.ErrorIs(t, NewService(restoreWithoutCAS{repo}, nil).RestoreRevision(ctx, p, p.RevisionID), ErrFormatUnavailable)
	require.ErrorIs(t, svc.RestoreRevision(ctx, p, "foreign-revision"), ErrFormatConflict)
	require.NoError(t, svc.Delete(ctx, p.ID))
	require.Error(t, svc.RestoreRevision(ctx, p, p.RevisionID))
	_, err = svc.GetByID(ctx, p.ID)
	require.ErrorIs(t, err, ErrPointNotFound)
}

func TestUpdateMaterializedRejectsFormatIntentAndStalePlan(t *testing.T) {
	ctx := t.Context()
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)
	p, err := svc.Create(ctx, CreatePointRequest{DeviceID: "d", Name: "p", Address: "40001", DataType: schema.DataTypeInt32})
	require.NoError(t, err)
	format := "CDAB"
	_, err = svc.UpdateMaterialized(ctx, p, UpdatePointRequest{DataFormat: &format})
	require.ErrorIs(t, err, ErrFormatInvalid)
	mode := schema.PointFormatOverride
	_, err = svc.UpdateMaterialized(ctx, p, UpdatePointRequest{FormatMode: &mode})
	require.ErrorIs(t, err, ErrFormatInvalid)
	_, err = svc.UpdateMaterialized(ctx, p, UpdatePointRequest{ConfiguredFormat: &format})
	require.ErrorIs(t, err, ErrFormatInvalid)
	_, err = svc.UpdateMaterialized(ctx, nil, UpdatePointRequest{})
	require.ErrorIs(t, err, ErrFormatInvalid)
	_, err = NewService(restoreWithoutCAS{repo}, nil).UpdateMaterialized(ctx, p, UpdatePointRequest{})
	require.ErrorIs(t, err, ErrFormatUnavailable)
	name := "ordinary"
	accepted, err := svc.Update(ctx, p.ID, UpdatePointRequest{Name: &name})
	require.NoError(t, err)
	_, err = svc.UpdateMaterialized(ctx, p, UpdatePointRequest{})
	require.ErrorIs(t, err, ErrFormatConflict)
	got, err := svc.GetByID(ctx, p.ID)
	require.NoError(t, err)
	require.Equal(t, accepted.RevisionID, got.RevisionID)
}

package runtime

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"path/filepath"
	"testing"
)

type pointStateCASWindowRepo struct {
	point.Repository
	hook func()
}

func (r *pointStateCASWindowRepo) UpdateReadResultAtRevision(ctx context.Context, id, deviceID, revision string, value any, message string) error {
	r.hook()
	return r.Repository.(point.ReadResultRevisionWriter).UpdateReadResultAtRevision(ctx, id, deviceID, revision, value, message)
}

func TestIngestorTelemetryCASCannotOverwriteChangedSQLRevision(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	d := f.projection.Devices[0]
	p.RevisionID = "read-r1"
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "state.db"))
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	require.NoError(t, device.NewSQLRepository(db).Create(t.Context(), d))
	repo := point.NewSQLRepository(db)
	require.NoError(t, repo.Create(t.Context(), p))
	wrapper := &pointStateCASWindowRepo{Repository: repo, hook: func() {
		_, err := db.ExecContext(t.Context(), `UPDATE points SET revision_id='new-r2',format_mode='inherit',format_version=1 WHERE id=?`, p.ID)
		require.NoError(t, err)
	}}
	f.service.pointSvc = point.NewService(wrapper, nil)
	f.service.config.UpdatePointState = true
	f.service.writer = &mockWriter{}
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	values := f.service.scheduler.PollNowContext(t.Context(), []string{p.ID})
	require.Len(t, values, 1)
	require.Equal(t, "read-r1", values[0].PointRevision)
	f.service.handleCollectedValue(t.Context(), values[0])
	current, err := repo.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Equal(t, "new-r2", current.RevisionID)
	require.Nil(t, current.LastValue)
	require.Nil(t, current.LastReadAt)
}

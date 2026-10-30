package handlers

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"testing"
)

type failedFormatShare struct {
	calls  int
	ctxErr error
}

func (s *failedFormatShare) CheckHydration(context.Context) (modbusshare.HydrationState, error) {
	return modbusshare.HydrationState{}, nil
}
func (s *failedFormatShare) InvalidatePointFormatRevision(ctx context.Context, _, _, _ string) error {
	s.calls++
	s.ctxErr = ctx.Err()
	return errors.New("secret://share-failure")
}

type formatRuntimeStopSpy struct{ calls int }

func (s *formatRuntimeStopSpy) SuspendPointRevision(_, _ string) bool { s.calls++; return true }

func TestPointFormatLifecycleStopsOldPointEvenWhenShareFails(t *testing.T) {
	for _, cancelled := range []bool{false, true} {
		t.Run(map[bool]string{false: "share-error", true: "cancel-after-commit"}[cancelled], func(t *testing.T) {
			repo := point.NewMemoryRepository()
			svc := point.NewService(repo, nil)
			require.NoError(t, repo.Create(t.Context(), &schema.Point{ID: "point", RevisionID: "new"}))
			share, live := &failedFormatShare{}, &formatRuntimeStopSpy{}
			lifecycle := &pointFormatLifecycle{points: svc, share: share, runtime: live}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if cancelled {
				cancel()
			}
			err := lifecycle.InvalidatePointFormat(ctx, point.FormatInvalidation{PointID: "point", ExpectedPointRevision: "old", PointRevision: "new", WorkspaceID: "ws", ExpectedWorkspaceRevision: "old-ws", WorkspaceRevision: "new-ws"})
			require.ErrorIs(t, err, point.ErrFormatUnavailable)
			require.Equal(t, 1, live.calls)
			require.Equal(t, 1, share.calls)
			require.NoError(t, share.ctxErr, "cleanup cannot inherit a cancelled HTTP lifetime")
		})
	}
}

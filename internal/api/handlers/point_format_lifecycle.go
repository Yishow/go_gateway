package handlers

import (
	"context"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"time"
)

type pointFormatShareInvalidator interface {
	CheckHydration(context.Context) (modbusshare.HydrationState, error)
	InvalidatePointFormatRevision(context.Context, string, string, string) error
}
type pointFormatRuntimeInvalidator interface{ SuspendPointRevision(string, string) bool }

type pointFormatLifecycle struct {
	points  *point.Service
	share   pointFormatShareInvalidator
	runtime pointFormatRuntimeInvalidator
}

// ConfigurePointFormatLifecycle is startup wiring shared by both HTTP entrypoints.
func ConfigurePointFormatLifecycle(points *point.Service, share pointFormatShareInvalidator, runtime pointFormatRuntimeInvalidator) {
	if points == nil || share == nil || runtime == nil {
		return
	}
	points.SetFormatInvalidator(&pointFormatLifecycle{points: points, share: share, runtime: runtime})
}

func (l *pointFormatLifecycle) InvalidatePointFormat(ctx context.Context, change point.FormatInvalidation) error {
	// A committed configuration outlives the HTTP request. Each matching old
	// runtime boundary is independently closed even if another cleanup fails.
	cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	stopped := l.runtime.SuspendPointRevision(change.PointID, change.ExpectedPointRevision)
	shareErr := l.share.InvalidatePointFormatRevision(cleanup, change.WorkspaceID, change.ExpectedWorkspaceRevision, change.WorkspaceRevision)
	current, err := l.points.GetByID(cleanup, change.PointID)
	if !stopped || shareErr != nil || err != nil || current.RevisionID != change.PointRevision {
		return point.ErrFormatUnavailable
	}

	return nil
}

// CheckPointFormatScope preserves the existing Share bootstrap gate for fresh
// writes. Durable replay is exempt: it must not alter or require newer readiness.
func (l *pointFormatLifecycle) CheckPointFormatScope(ctx context.Context, scope point.FormatScope) error {
	state, err := l.share.CheckHydration(ctx)
	if err != nil || !state.Readiness || state.State != modbusshare.HydrationStateReady {
		return point.ErrFormatUnavailable
	}
	if state.WorkspaceID != scope.WorkspaceID {
		return point.ErrFormatScope
	}
	if state.WorkspaceRevision != scope.ExpectedWorkspaceRevision {
		return point.ErrFormatConflict
	}
	return nil
}

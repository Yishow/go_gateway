package handlers

import (
	"errors"
	"go-gateway/internal/datalink/dbtarget"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
)

// This optional pair scopes existing pure native editor reads. It never creates
// a workspace or grants permission for a later native mutation.
type studioLocatorNativeScope struct {
	WorkspaceID, Revision string
	Before                *workspace.Record
}

func parseStudioLocatorNativeScope(raw string) (*studioLocatorNativeScope, error) {
	q, err := url.ParseQuery(raw)
	_, hasID := q["locator_workspace_id"]
	_, hasRevision := q["locator_workspace_revision"]
	if err == nil && !hasID && !hasRevision {
		return nil, nil
	}
	if err != nil || !hasID || !hasRevision {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	for _, values := range q {
		if len(values) != 1 {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
	}
	if !runtimeDeliveryIDs(q.Get("locator_workspace_id")) || (q.Get("locator_workspace_revision") != "" && !runtimeDeliveryIDs(q.Get("locator_workspace_revision"))) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	return &studioLocatorNativeScope{WorkspaceID: q.Get("locator_workspace_id"), Revision: q.Get("locator_workspace_revision")}, nil
}
func studioLocatorNativeRequested(c *gin.Context) bool {
	scope, err := parseStudioLocatorNativeScope(c.Request.URL.RawQuery)
	return scope != nil || err != nil
}
func validateStudioLocatorNativeQuery(c *gin.Context) bool {
	_, err := parseStudioLocatorNativeScope(c.Request.URL.RawQuery)
	if err != nil {
		renderRuntimeHandlerError(c, http.StatusBadRequest, string(runtimeview.ReasonUnavailable), false)
		return false
	}
	return true
}
func beginStudioLocatorNativeRead(c *gin.Context, svc *workspace.Service, legacyFailure func()) (*workspace.Record, *studioLocatorNativeScope, bool) {
	scope, err := parseStudioLocatorNativeScope(c.Request.URL.RawQuery)
	if err != nil {
		renderRuntimeHandlerError(c, http.StatusBadRequest, string(runtimeview.ReasonUnavailable), false)
		return nil, nil, false
	}
	var record *workspace.Record
	if scope == nil {
		record, err = svc.GetOrCreate(c.Request.Context())
		if err != nil {
			legacyFailure()
			return nil, nil, false
		}
		return record, nil, true
	}
	record, err = svc.ReadSnapshot(c.Request.Context())
	if err != nil || record == nil {
		renderRuntimeEvidenceError(c, studioLocatorReadError(err))
		return nil, scope, false
	}
	if record.ID != scope.WorkspaceID {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonForeign))
		return nil, scope, false
	}
	if record.DatabaseSetupRevision != scope.Revision {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return nil, scope, false
	}
	scope.Before = record
	return record, scope, true
}
func finishStudioLocatorNativeRead(c *gin.Context, svc *workspace.Service, scope *studioLocatorNativeScope) bool {
	if scope == nil {
		return true
	}
	after, err := svc.ReadSnapshot(c.Request.Context())
	if err != nil || after == nil {
		renderRuntimeEvidenceError(c, studioLocatorReadError(err))
		return false
	}
	if !sameStudioLocatorRecord(scope.Before, after) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return false
	}
	return true
}

// The explicit scoped branch keeps known CAS failures and sanitizes provider
// failures; the ordinary native response contract remains unchanged.
func renderStudioLocatorNativeProviderError(c *gin.Context, scope *studioLocatorNativeScope, err error) bool {
	if scope == nil {
		return false
	}
	if errors.Is(err, dbtarget.ErrConnectorDisabled) || errors.Is(err, dbtarget.ErrValidation) || errors.Is(err, workspace.ErrValidation) {
		return false
	}
	reason := runtimeview.ReasonUnavailable
	if errors.Is(err, workspace.ErrWriteGroupNotFound) || errors.Is(err, workspace.ErrNotFound) || errors.Is(err, workspace.ErrSetupRevisionConflict) || errors.Is(err, workspace.ErrWriteGroupRevisionConflict) || errors.Is(err, dbtarget.ErrConnectorNotFound) || errors.Is(err, dbtarget.ErrConnectorRevisionConflict) {
		reason = runtimeview.ReasonStale
	}
	renderRuntimeEvidenceError(c, runtimeEvidenceError(reason))
	return true
}

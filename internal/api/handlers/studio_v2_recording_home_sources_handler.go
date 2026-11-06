package handlers

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"net/url"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
)

// RecordingHomeWorkspaceReader reads only an existing singleton.
type RecordingHomeWorkspaceReader interface {
	ReadSnapshot(context.Context) (*workspace.Record, error)
}

// RecordingHomeDeviceReader reads saved device presentation.
type RecordingHomeDeviceReader interface {
	GetByID(context.Context, string) (*schema.Device, error)
}

// RecordingHomeRuleReader reads scoped rules and their canonical links.
type RecordingHomeRuleReader interface {
	ListByDeviceIDs(context.Context, []string) ([]*schema.SourceRule, error)
	ListLinks(context.Context, string) ([]*schema.SourceRuleLink, error)
}

// RecordingHomePointReader reads the complete point catalog without a page limit.
type RecordingHomePointReader interface {
	List(context.Context, point.ListFilter) ([]*schema.Point, error)
}

// RecordingHomeMappingReader reads direct identities and complete point-scoped mappings.
type RecordingHomeMappingReader interface {
	GetByID(context.Context, string) (*schema.Mapping, error)
	List(context.Context, mapping.ListFilter) ([]*schema.Mapping, error)
}

// RecordingHomeTagReader reads saved tag presentation and ownership.
type RecordingHomeTagReader interface {
	GetByID(context.Context, string) (*schema.Tag, error)
}

// StudioV2RecordingHomeSourcesHandler serves readonly home observations.
type StudioV2RecordingHomeSourcesHandler struct {
	workspace RecordingHomeWorkspaceReader
	devices   RecordingHomeDeviceReader
	rules     RecordingHomeRuleReader
	points    RecordingHomePointReader
	mappings  RecordingHomeMappingReader
	tags      RecordingHomeTagReader
}

// NewStudioV2RecordingHomeSourcesHandler accepts only read capabilities.
func NewStudioV2RecordingHomeSourcesHandler(w RecordingHomeWorkspaceReader, d RecordingHomeDeviceReader, r RecordingHomeRuleReader, p RecordingHomePointReader, m RecordingHomeMappingReader, t RecordingHomeTagReader) *StudioV2RecordingHomeSourcesHandler {
	return &StudioV2RecordingHomeSourcesHandler{w, d, r, p, m, t}
}

// RecordingHomeSourcesEnvelope wraps a readonly presentation observation.
type RecordingHomeSourcesEnvelope struct {
	Success bool                                      `json:"success"`
	Data    workspace.RecordingHomeSourcesObservation `json:"data"`
}

// RecordingHomeSourcesError exposes a stable code, never a repository error.
type RecordingHomeSourcesError struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
}

var recordingHomeScopePattern = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9_.-]{0,127}$`)

// Get observes existing sources without initializing, repairing or reconciling.
// @Summary Observe recording home sources
// @Description Readonly two-pass presentation evidence, not an atomic snapshot or command authorization. Unknown contains no partial summary.
// @Tags Studio V2 Recording Home
// @Produce json
// @Param workspace_id query string true "Existing workspace ID"
// @Success 200 {object} RecordingHomeSourcesEnvelope
// @Failure 400 {object} RecordingHomeSourcesError
// @Failure 404 {object} RecordingHomeSourcesError
// @Failure 503 {object} RecordingHomeSourcesError
// @Router /v1/datalink/studio-v2/recording-home/sources [get]
func (h *StudioV2RecordingHomeSourcesHandler) Get(c *gin.Context) {
	query, err := url.ParseQuery(c.Request.URL.RawQuery)
	ids := query["workspace_id"]
	if err != nil || len(query) != 1 || len(ids) != 1 || !recordingHomeScopePattern.MatchString(ids[0]) {
		recordingHomeError(c, http.StatusBadRequest, "invalid-workspace-scope")
		return
	}
	if h == nil || homeNil(h.workspace) || homeNil(h.devices) || homeNil(h.rules) || homeNil(h.points) || homeNil(h.mappings) || homeNil(h.tags) {
		recordingHomeError(c, http.StatusServiceUnavailable, "source-unavailable")
		return
	}
	observed, err := workspace.ReadRecordingHomeSources(c.Request.Context(), ids[0], h.readSources)
	if err != nil {
		if errors.Is(err, workspace.ErrRecordingHomeWorkspaceNotFound) {
			recordingHomeError(c, http.StatusNotFound, "workspace-not-found")
		} else {
			recordingHomeError(c, http.StatusServiceUnavailable, "source-unavailable")
		}
		return
	}
	c.JSON(http.StatusOK, RecordingHomeSourcesEnvelope{Success: true, Data: observed})
}
func recordingHomeError(c *gin.Context, status int, code string) {
	c.JSON(status, RecordingHomeSourcesError{Code: code})
}
func homeNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	}
	return false
}

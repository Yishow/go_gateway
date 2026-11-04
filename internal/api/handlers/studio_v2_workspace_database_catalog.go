package handlers

import (
	"errors"
	"net/http"
	"net/url"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type studioV2WorkspaceDatabaseCatalogQuery struct {
	workspaceID           string
	expectedWorkspaceRev  string
	connectorID           string
	expectedConnectorRev  string
	groupID               string
	expectedGroupRevision string
}

type studioV2WorkspaceDatabaseCatalogResponse struct {
	WorkspaceID       string                       `json:"workspace_id"`
	WorkspaceRevision string                       `json:"workspace_revision"`
	ConnectorID       string                       `json:"connector_id"`
	ConnectorRevision string                       `json:"connector_revision"`
	Kind              schema.DatabaseConnectorKind `json:"kind"`
	DefaultSchema     string                       `json:"default_schema"`
	Schemas           []dbtarget.CatalogSchema     `json:"schemas"`
}

const (
	studioV2DatabaseCatalogConnectorIDKey               = "connector_id"
	studioV2DatabaseCatalogExpectedConnectorRevisionKey = "expected_connector_revision"
)

func parseStudioV2WorkspaceDatabaseCatalogQuery(raw string) (studioV2WorkspaceDatabaseCatalogQuery, error) {
	query, err := url.ParseQuery(raw)
	if err != nil {
		return studioV2WorkspaceDatabaseCatalogQuery{}, errors.New("database catalog query is invalid")
	}
	allowed := map[string]bool{
		workspaceIDKey: true, expectedWorkspaceRevisionField: true,
		studioV2DatabaseCatalogConnectorIDKey: true, studioV2DatabaseCatalogExpectedConnectorRevisionKey: true,
		"group_id": true, "expected_group_revision": true,
		"locator_workspace_id": true, "locator_workspace_revision": true,
	}
	for key, values := range query {
		if !allowed[key] || len(values) != 1 {
			return studioV2WorkspaceDatabaseCatalogQuery{}, errors.New("database catalog query is invalid")
		}
		value := values[0]
		if key == expectedWorkspaceRevisionField || key == "locator_workspace_revision" {
			if value != "" && !runtimeDeliveryIDs(value) {
				return studioV2WorkspaceDatabaseCatalogQuery{}, errors.New("database catalog query is invalid")
			}
			continue
		}
		if !runtimeDeliveryIDs(value) {
			return studioV2WorkspaceDatabaseCatalogQuery{}, errors.New("database catalog query is invalid")
		}
	}
	for _, key := range []string{workspaceIDKey, expectedWorkspaceRevisionField, studioV2DatabaseCatalogConnectorIDKey, studioV2DatabaseCatalogExpectedConnectorRevisionKey} {
		if _, ok := query[key]; !ok {
			return studioV2WorkspaceDatabaseCatalogQuery{}, errors.New("database catalog query is incomplete")
		}
	}
	if _, hasGroupID := query["group_id"]; hasGroupID != (query.Get("group_id") != "") {
		return studioV2WorkspaceDatabaseCatalogQuery{}, errors.New("database catalog group_id is invalid")
	}
	if query.Get("group_id") != "" {
		if _, ok := query["expected_group_revision"]; !ok || query.Get("expected_group_revision") == "" {
			return studioV2WorkspaceDatabaseCatalogQuery{}, errors.New("expected_group_revision is required")
		}
	} else if _, ok := query["expected_group_revision"]; ok {
		return studioV2WorkspaceDatabaseCatalogQuery{}, errors.New("group_id is required")
	}
	_, hasLocatorID := query["locator_workspace_id"]
	_, hasLocatorRevision := query["locator_workspace_revision"]
	if hasLocatorID != hasLocatorRevision {
		return studioV2WorkspaceDatabaseCatalogQuery{}, errors.New("locator workspace query is incomplete")
	}
	return studioV2WorkspaceDatabaseCatalogQuery{
		workspaceID:           query.Get(workspaceIDKey),
		expectedWorkspaceRev:  query.Get(expectedWorkspaceRevisionField),
		connectorID:           query.Get(studioV2DatabaseCatalogConnectorIDKey),
		expectedConnectorRev:  query.Get(studioV2DatabaseCatalogExpectedConnectorRevisionKey),
		groupID:               query.Get("group_id"),
		expectedGroupRevision: query.Get("expected_group_revision"),
	}, nil
}

// GetCatalog returns the user-visible schemas and tables of an enabled saved
// connector. Workspace, connector, and optional write-group revisions are
// checked before and after the external read so a changed destination never
// gets presented as a coherent catalog.
// @Summary Read the saved database catalog
// @Tags Studio V2 Recording Plans
// @Produce json
// @Param workspace_id query string true "Exact owning workspace"
// @Param expected_workspace_revision query string true "Exact saved database setup revision; empty is valid only for legacy-empty workspaces"
// @Param connector_id query string true "Enabled saved connector to inspect"
// @Param expected_connector_revision query string true "Current saved connector identity revision"
// @Param group_id query string false "Exact saved write group; requires expected_group_revision"
// @Param expected_group_revision query string false "Current saved write-group revision"
// @Param locator_workspace_id query string false "Locator-origin workspace identity; requires revision pair"
// @Param locator_workspace_revision query string false "Exact observed locator workspace revision"
// @Success 200 {object} studioV2WorkspaceDatabaseCatalogResponse
// @Failure 400 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 409 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Failure 503 {object} APIErrorResponse
// @Router /v1/datalink/studio-v2/workspace/database-catalog [get]
func (h *StudioV2WorkspaceDatabaseHandler) GetCatalog(c *gin.Context) {
	query, err := parseStudioV2WorkspaceDatabaseCatalogQuery(c.Request.URL.RawQuery)
	if err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}
	record, locatorScope, ok := beginStudioLocatorNativeRead(c, h.workspaceSvc, func() { renderStudioV2WorkspaceBootstrapError(c) })
	if !ok {
		return
	}
	if record.ID != query.workspaceID {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonForeign))
		return
	}
	if record.DatabaseSetupRevision != query.expectedWorkspaceRev {
		if locatorScope != nil {
			renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		} else {
			renderStudioV2WorkspaceDatabaseError(c, workspace.ErrSetupRevisionConflict)
		}
		return
	}

	var group *workspace.WriteGroup
	if query.groupID != "" {
		group, err = h.metadataGroup(c, record.ID, query.groupID, query.expectedConnectorRev)
		if err != nil {
			renderCatalogReadError(c, locatorScope, err)
			return
		}
		if group.Destination.ConnectorID != query.connectorID {
			renderCatalogReadError(c, locatorScope, workspace.ErrWriteGroupRevisionConflict)
			return
		}
	}
	if err := h.validateCatalogDestination(c, query.connectorID, query.expectedConnectorRev); err != nil {
		renderCatalogReadError(c, locatorScope, err)
		return
	}

	catalog, err := h.connectorSvc.ReadCatalog(c.Request.Context(), query.connectorID, query.expectedConnectorRev)
	if err != nil {
		renderCatalogReadError(c, locatorScope, err)
		return
	}

	// The external read is complete. Re-read all local authorities before
	// returning it so a connector, workspace, or group change cannot be hidden
	// behind a successful catalog query.
	if _, err := h.connectorSvc.ResolveSavedTarget(c.Request.Context(), query.connectorID, query.expectedConnectorRev); err != nil {
		renderCatalogReadError(c, locatorScope, err)
		return
	}
	if group != nil {
		if _, err := h.metadataGroup(c, record.ID, query.groupID, query.expectedConnectorRev); err != nil {
			renderCatalogReadError(c, locatorScope, err)
			return
		}
	}
	if locatorScope == nil {
		current, readErr := h.workspaceSvc.ReadSnapshot(c.Request.Context())
		if readErr != nil || current == nil {
			renderStudioV2WorkspaceBootstrapError(c)
			return
		}
		if current.ID != record.ID || current.DatabaseSetupRevision != record.DatabaseSetupRevision {
			renderStudioV2WorkspaceDatabaseError(c, workspace.ErrSetupRevisionConflict)
			return
		}
	}
	if !finishStudioLocatorNativeRead(c, h.workspaceSvc, locatorScope) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey: studioV2WorkspaceDatabaseCatalogResponse{
			WorkspaceID:       record.ID,
			WorkspaceRevision: record.DatabaseSetupRevision,
			ConnectorID:       query.connectorID,
			ConnectorRevision: query.expectedConnectorRev,
			Kind:              catalog.Kind,
			DefaultSchema:     catalog.DefaultSchema,
			Schemas:           catalog.Schemas,
		},
	})
}

// validateCatalogDestination reuses the managed destination path guard when
// the production write-group service is wired. The guard is deliberately
// applied before the external read, so an internal SQLite database can never
// be opened as a user-selected catalog.
func (h *StudioV2WorkspaceDatabaseHandler) validateCatalogDestination(c *gin.Context, connectorID, expectedRevision string) error {
	connector, err := h.connectorSvc.ResolveSavedTarget(c.Request.Context(), connectorID, expectedRevision)
	if err != nil {
		return err
	}
	if connector == nil || connector.Kind != schema.DatabaseConnectorKindSQLite {
		return nil
	}
	groups, ok := h.writeGroups.(*workspace.WriteGroupService)
	if !ok || groups == nil {
		return workspace.ErrWriteGroupServiceUnavailable
	}
	inspector := groups.ManagedTableInspector(h.connectorSvc)
	if inspector == nil {
		return workspace.ErrWriteGroupServiceUnavailable
	}
	_, err = inspector.ValidateDestination(c.Request.Context(), connectorID, expectedRevision)
	return err
}

func renderCatalogReadError(c *gin.Context, scope *studioLocatorNativeScope, err error) {
	if scope != nil && renderStudioLocatorNativeProviderError(c, scope, err) {
		return
	}
	if errors.Is(err, dbtarget.ErrManagedDestinationInternal) {
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "WRITE_GROUP_INTERNAL_DATABASE", "recording requires a distinct application-owned destination", false, "select_recording_destination")
		return
	}
	if errors.Is(err, dbtarget.ErrManagedDestinationInvalid) {
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "WRITE_GROUP_DESTINATION_INVALID", "the saved recording destination cannot be verified", false, "select_recording_destination")
		return
	}
	if errors.Is(err, workspace.ErrWriteGroupNotFound) || errors.Is(err, workspace.ErrWriteGroupRevisionConflict) || errors.Is(err, workspace.ErrWriteGroupValidation) || errors.Is(err, workspace.ErrWriteGroupServiceUnavailable) {
		renderWriteGroupError(c, err)
		return
	}
	renderStudioV2WorkspaceDatabaseError(c, err)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

type DatabaseTargetHandler struct {
	connectorSvc *dbtarget.ConnectorService
	mappingSvc   *dbtarget.MappingService
}

func NewDatabaseTargetHandler(connectorSvc *dbtarget.ConnectorService, mappingSvc *dbtarget.MappingService) *DatabaseTargetHandler {
	return &DatabaseTargetHandler{
		connectorSvc: connectorSvc,
		mappingSvc:   mappingSvc,
	}
}

type databaseConnectorResponse struct {
	ID               string                         `json:"id"`
	Name             string                         `json:"name"`
	Kind             schema.DatabaseConnectorKind   `json:"kind"`
	ConnectionConfig map[string]any                 `json:"connection_config"`
	IdentityRevision string                         `json:"identity_revision"`
	Status           schema.DatabaseConnectorStatus `json:"status"`
	LastCheckAt      any                            `json:"last_check_at,omitempty"`
	LastCheckError   string                         `json:"last_check_error,omitempty"`
	Enabled          bool                           `json:"enabled"`
	CreatedAt        any                            `json:"created_at"`
	UpdatedAt        any                            `json:"updated_at"`
}

func (h *DatabaseTargetHandler) ListConnectors(c *gin.Context) {
	filter := dbtarget.ConnectorListFilter{
		Enabled: parseBoolQuery(c, "enabled"),
	}

	connectors, err := h.connectorSvc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	data := make([]databaseConnectorResponse, 0, len(connectors))
	for _, connector := range connectors {
		data = append(data, toDatabaseConnectorResponse(connector))
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: data})
}

func (h *DatabaseTargetHandler) CreateConnector(c *gin.Context) {
	var req dbtarget.CreateConnectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	connector, err := h.connectorSvc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: toDatabaseConnectorResponse(connector)})
}

func (h *DatabaseTargetHandler) GetConnector(c *gin.Context) {
	connector, err := h.connectorSvc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: toDatabaseConnectorResponse(connector)})
}

func (h *DatabaseTargetHandler) UpdateConnector(c *gin.Context) {
	var req dbtarget.UpdateConnectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	connector, err := h.connectorSvc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		status := http.StatusUnprocessableEntity
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: toDatabaseConnectorResponse(connector)})
}

func (h *DatabaseTargetHandler) DeleteConnector(c *gin.Context) {
	if err := h.connectorSvc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		status := http.StatusInternalServerError
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true})
}

func (h *DatabaseTargetHandler) TestConnector(c *gin.Context) {
	connector, err := h.connectorSvc.TestConnection(c.Request.Context(), c.Param("id"))
	if err != nil {
		status := http.StatusInternalServerError
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: toDatabaseConnectorResponse(connector)})
}

func (h *DatabaseTargetHandler) ListTables(c *gin.Context) {
	tables, err := h.connectorSvc.ListTables(c.Request.Context(), c.Param("id"))
	if err != nil {
		status := http.StatusUnprocessableEntity
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: tables})
}

func (h *DatabaseTargetHandler) ValidateConnector(c *gin.Context) {
	result, err := h.mappingSvc.Validate(c.Request.Context(), c.Param("id"))
	if err != nil {
		status := http.StatusInternalServerError
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

func (h *DatabaseTargetHandler) ListMappings(c *gin.Context) {
	filter := dbtarget.TargetMappingListFilter{
		ConnectorID: optionalQuery(c, "connector_id"),
		TagID:       optionalQuery(c, "tag_id"),
		Enabled:     parseBoolQuery(c, "enabled"),
	}

	mappings, err := h.mappingSvc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: mappings})
}

func (h *DatabaseTargetHandler) CreateMapping(c *gin.Context) {
	var req dbtarget.CreateTargetMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	mapping, err := h.mappingSvc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: mapping})
}

func (h *DatabaseTargetHandler) GetMapping(c *gin.Context) {
	mapping, err := h.mappingSvc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: mapping})
}

func (h *DatabaseTargetHandler) UpdateMapping(c *gin.Context) {
	var req dbtarget.UpdateTargetMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	mapping, err := h.mappingSvc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		status := http.StatusUnprocessableEntity
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: mapping})
}

func (h *DatabaseTargetHandler) DeleteMapping(c *gin.Context) {
	if err := h.mappingSvc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		status := http.StatusInternalServerError
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true})
}

func toDatabaseConnectorResponse(connector *schema.DatabaseConnector) databaseConnectorResponse {
	response := databaseConnectorResponse{
		ID:               connector.ID,
		Name:             connector.Name,
		Kind:             connector.Kind,
		ConnectionConfig: map[string]any{},
		IdentityRevision: connector.IdentityRevision,
		Status:           connector.Status,
		LastCheckError:   connector.LastCheckError,
		Enabled:          connector.Enabled,
		CreatedAt:        connector.CreatedAt,
		UpdatedAt:        connector.UpdatedAt,
	}
	if connector.LastCheckAt != nil {
		response.LastCheckAt = connector.LastCheckAt
	}
	if err := json.Unmarshal([]byte(connector.ConnectionConfig), &response.ConnectionConfig); err != nil {
		response.ConnectionConfig = map[string]any{}
	}
	redactConnectionSecrets(response.ConnectionConfig)
	return response
}

func parseBoolQuery(c *gin.Context, key string) *bool {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil
	}
	return &value
}

func optionalQuery(c *gin.Context, key string) *string {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return nil
	}
	return &value
}

func isNotFoundError(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "不存在") || strings.Contains(message, "not found")
}

func redactConnectionSecrets(config map[string]any) {
	for _, key := range []string{"password", "passwd", "secret", "token", "api_key"} {
		if _, ok := config[key]; ok {
			config[key] = ""
		}
	}
}

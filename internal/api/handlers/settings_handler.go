package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/settings"

	"github.com/gin-gonic/gin"
)

// SettingsHandler 設定 API Handler
type SettingsHandler struct {
	svc         *settings.Service
	modbusShare *modbusshare.Service
}

// NewSettingsHandler 建立新的設定 Handler
func NewSettingsHandler(svc *settings.Service) *SettingsHandler {
	return &SettingsHandler{svc: svc}
}

// WithModbusShare wires the durable settings update to the listener lifecycle.
func (h *SettingsHandler) WithModbusShare(svc *modbusshare.Service) *SettingsHandler {
	h.modbusShare = svc
	return h
}

// List 列出所有設定
// GET /datalink/settings
// @Summary List durable settings
// @Tags datalink
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} APIErrorResponse
// @Router /datalink/settings [get]
func (h *SettingsHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		renderSafeError(c, http.StatusInternalServerError, ErrCodeSettingsUnavailable, true)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: items})
}

// UpdateRequest 更新設定請求
type UpdateSettingRequest struct {
	Value                    interface{} `json:"value"`
	ExpectedSettingsRevision string      `json:"expected_settings_revision,omitempty"`
}

// Update 更新設定
// PUT /datalink/settings/:key
// @Summary Update durable settings
// @Tags datalink
// @Accept json
// @Produce json
// @Param key path string true "Setting key"
// @Param request body UpdateSettingRequest true "Settings value and expected revision"
// @Success 200 {object} map[string]interface{}
// @Failure 409 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Router /datalink/settings/{key} [put]
func (h *SettingsHandler) Update(c *gin.Context) {
	key := c.Param("key")

	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderSafeError(c, http.StatusBadRequest, ErrCodeSettingsInvalid, false)
		return
	}

	if key == settings.KeyModbusShare && h.modbusShare != nil {
		payload, err := json.Marshal(req.Value)
		if err != nil {
			renderSafeError(c, http.StatusBadRequest, ErrCodeSettingsInvalid, false)
			return
		}
		var shareSettings modbusshare.Settings
		if err := json.Unmarshal(payload, &shareSettings); err != nil {
			renderSafeError(c, http.StatusBadRequest, ErrCodeSettingsInvalid, false)
			return
		}
		expected := req.ExpectedSettingsRevision
		if expected == "" {
			expected = shareSettings.SettingsRevision
		}
		if err := h.modbusShare.ApplySettingsCAS(c.Request.Context(), shareSettings, expected); err != nil {
			var shareErr *modbusshare.Error
			if errors.As(err, &shareErr) {
				status := http.StatusInternalServerError
				switch shareErr.Code {
				case modbusshare.ErrCodeRevisionConflict:
					status = http.StatusConflict
				case modbusshare.ErrCodeCapacityExceeded, modbusshare.ErrCodeInvalidGeometry, modbusshare.ErrCodeListenerBindFailed:
					status = http.StatusUnprocessableEntity
				}
				renderModbusShareAPIError(c, status, shareErr, ErrCodeSettingsUpdateFailed, true)
				return
			}
			renderSafeError(c, http.StatusInternalServerError, ErrCodeSettingsUpdateFailed, true)
			return
		}
	} else if err := h.svc.Set(c.Request.Context(), key, req.Value); err != nil {
		renderSafeError(c, http.StatusInternalServerError, ErrCodeSettingsUpdateFailed, true)
		return
	}

	// 取得更新後的設定
	item, err := h.svc.Get(c.Request.Context(), key)
	if err != nil {
		renderSafeError(c, http.StatusInternalServerError, ErrCodeSettingsUnavailable, true)
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: item})
}

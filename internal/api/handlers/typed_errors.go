package handlers

import (
	"errors"
	"unicode"

	"go-gateway/internal/datalink/modbusshare"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ErrCodePreviewInvalidRequest      = "preview_invalid_request"
	ErrCodePreviewUnavailable         = "preview_unavailable"
	ErrCodePreviewStreamClosed        = "preview_stream_closed"
	ErrCodeRuntimeDeviceNotFound      = "runtime_device_not_found"
	ErrCodeRuntimeSnapshotUnavailable = "runtime_snapshot_unavailable"
	ErrCodeRuntimeStreamUnavailable   = "runtime_stream_unavailable"
	ErrCodeWorkspaceNotReady          = "workspace_not_ready"
	ErrCodeActivationFailed           = "activation_failed"
	ErrCodeSettingsUnavailable        = "settings_unavailable"
	ErrCodeSettingsUpdateFailed       = "settings_update_failed"
	ErrCodeSettingsInvalid            = "settings_invalid"
	ErrCodeActivationRequestInvalid   = "activation_request_invalid"
	ErrCodeReadinessBlocked           = "readiness_blocked"
)

// TypedAPIErrorEnvelope represents a stable, safe typed error response.
type TypedAPIErrorEnvelope struct {
	Code              string `json:"code"`
	Message           string `json:"message"`
	Retryable         bool   `json:"retryable"`
	RequestID         string `json:"request_id"`
	Action            string `json:"action,omitempty"`
	WorkspaceRevision string `json:"workspace_revision,omitempty"`
	SettingsRevision  string `json:"settings_revision,omitempty"`
	DirtyState        string `json:"dirty_state,omitempty"`
}

func renderModbusShareAPIError(c *gin.Context, status int, err error, fallbackCode string, fallbackRetryable bool) {
	var shareErr *modbusshare.Error
	if errors.As(err, &shareErr) {
		action := shareErr.Action
		if action == "" {
			action = "Review the request and retry"
		}
		payload := gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: TypedAPIErrorEnvelope{
			Code: shareErr.Code, Message: typedAPIErrorMessage(shareErr.Code), Retryable: shareErr.Retryable,
			RequestID: getOrGenerateRequestID(c), Action: action,
			WorkspaceRevision: shareErr.WorkspaceRevision, SettingsRevision: shareErr.SettingsRevision,
			DirtyState: shareErr.DirtyState,
		}}
		c.JSON(status, payload)
		return
	}
	renderSafeError(c, status, fallbackCode, fallbackRetryable)
}

// APIErrorResponse is the stable top-level shape returned for typed failures.
type APIErrorResponse struct {
	Success bool                  `json:"success" example:"false"`
	Error   TypedAPIErrorEnvelope `json:"error"`
}

func getOrGenerateRequestID(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return "req-" + uuid.NewString()[:8]
	}
	reqID := c.GetHeader("X-Request-ID")
	if !isOpaqueRequestID(reqID) {
		reqID = "req-" + uuid.NewString()[:8]
	}
	return reqID
}

func isOpaqueRequestID(value string) bool {
	if value == "" || len(value) > 128 {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func typedAPIErrorMessage(code string) string {
	switch code {
	case ErrCodePreviewInvalidRequest:
		return "preview request is invalid"
	case ErrCodePreviewUnavailable:
		return "preview is unavailable"
	case ErrCodePreviewStreamClosed:
		return "preview stream closed unexpectedly"
	case ErrCodeRuntimeDeviceNotFound:
		return "runtime device was not found"
	case ErrCodeRuntimeSnapshotUnavailable:
		return "runtime snapshot is unavailable"
	case ErrCodeRuntimeStreamUnavailable:
		return "runtime stream is unavailable"
	case ErrCodeWorkspaceNotReady:
		return "workspace is not ready"
	case ErrCodeActivationFailed:
		return "workspace activation failed"
	case ErrCodeSettingsUnavailable:
		return "settings are temporarily unavailable"
	case ErrCodeSettingsUpdateFailed:
		return "settings could not be saved"
	case ErrCodeSettingsInvalid:
		return "settings request is invalid"
	case ErrCodeActivationRequestInvalid:
		return "activation request is invalid"
	case ErrCodeReadinessBlocked:
		return "workspace activation is blocked by readiness issues"
	default:
		return "request could not be completed"
	}
}

func renderSafeError(c *gin.Context, status int, code string, retryable bool) {
	renderTypedAPIError(c, status, code, retryable)
}

func renderTypedAPIError(c *gin.Context, status int, code string, retryable bool) {
	renderTypedAPIErrorWithData(c, status, code, retryable, nil)
}

func renderTypedAPIErrorWithData(c *gin.Context, status int, code string, retryable bool, data any) {
	payload := gin.H{
		apiResponseSuccessKey: false,
		apiResponseErrorKey: TypedAPIErrorEnvelope{
			Code:      code,
			Message:   typedAPIErrorMessage(code),
			Retryable: retryable,
			RequestID: getOrGenerateRequestID(c),
			Action:    "Review the request and retry",
		},
	}
	if data != nil {
		payload["data"] = data
	}
	c.JSON(status, payload)
}

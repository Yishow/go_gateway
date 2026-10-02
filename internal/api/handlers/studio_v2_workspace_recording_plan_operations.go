package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderRecordingOperationNotImplemented(c *gin.Context, code, message string) {
	c.JSON(http.StatusNotImplemented, gin.H{
		apiResponseSuccessKey: false,
		apiResponseErrorKey: TypedAPIErrorEnvelope{
			Code:      code,
			Message:   message,
			Retryable: false,
			RequestID: getOrGenerateRequestID(c),
			Action:    "wait_for_supported_operation",
		},
	})
}

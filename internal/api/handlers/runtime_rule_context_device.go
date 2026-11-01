package handlers

import (
	"encoding/json"
	"strings"

	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

// Legacy devices with no explicit availability field retain their existing semantics.
// Explicit malformed or unknown values must not become available by default.
func runtimeSavedDeviceAvailability(record *schema.Device) string {
	if record == nil {
		return runtimeEvidenceUnavailable
	}
	if strings.TrimSpace(record.ReadinessStatus) == "" {
		return availabilityStatusOf(record)
	}
	var readiness struct {
		AvailabilityStatus json.RawMessage `json:"availability_status"`
	}
	if json.Unmarshal([]byte(record.ReadinessStatus), &readiness) != nil {
		return runtimeEvidenceUnavailable
	}
	if len(readiness.AvailabilityStatus) == 0 {
		return availabilityStatusOf(record)
	}
	var status string
	if json.Unmarshal(readiness.AvailabilityStatus, &status) != nil {
		return runtimeEvidenceUnavailable
	}
	return runtimeSafeStatusChoice(status, runtimeEvidenceAvailable, runtimeEvidenceUnavailable)
}

// This normalization is deliberately limited to runtime requests; shared authoring
// and preview error envelopes retain their existing ownership and behavior.
func ensureRuntimeRequestID(c *gin.Context) {
	if c == nil || c.Request == nil {
		return
	}
	if runtimeSafeID(c.GetHeader("X-Request-ID")) == "" {
		c.Request.Header.Del("X-Request-ID")
	}
	c.Request.Header.Set("X-Request-ID", getOrGenerateRequestID(c))
}

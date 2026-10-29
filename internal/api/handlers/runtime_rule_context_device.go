package handlers

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/schema"
)

// Legacy devices with no explicit availability field retain their existing semantics.
// Explicit malformed or unknown values must not become available by default.
func runtimeSavedDeviceAvailability(record *schema.Device) string {
	if record == nil {
		return "unavailable"
	}
	if strings.TrimSpace(record.ReadinessStatus) == "" {
		return availabilityStatusOf(record)
	}
	var readiness struct {
		AvailabilityStatus json.RawMessage `json:"availability_status"`
	}
	if json.Unmarshal([]byte(record.ReadinessStatus), &readiness) != nil {
		return "unavailable"
	}
	if len(readiness.AvailabilityStatus) == 0 {
		return availabilityStatusOf(record)
	}
	var status string
	if json.Unmarshal(readiness.AvailabilityStatus, &status) != nil {
		return "unavailable"
	}
	return runtimeSafeStatusChoice(status, "available", "unavailable")
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

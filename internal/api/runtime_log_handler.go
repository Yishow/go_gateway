package api

import (
	"errors"
	"net/http"
	"time"

	"go-gateway/internal/diagnostics"

	"github.com/gin-gonic/gin"
)

const runtimeLogHeartbeat = 15 * time.Second

type runtimeLogHandler struct {
	broker    *diagnostics.Broker
	heartbeat time.Duration
}

func registerRuntimeLogRoutes(group *gin.RouterGroup, broker *diagnostics.Broker) {
	h := &runtimeLogHandler{broker: broker, heartbeat: runtimeLogHeartbeat}
	logs := group.Group("/system/logs", localLogGuard())
	logs.GET("", h.snapshot)
	logs.GET("/stream", h.stream)
}

func (h *runtimeLogHandler) snapshot(c *gin.Context) {
	query, _, err := parseRuntimeLogQuery(c.Request, false)
	if err != nil {
		runtimeLogError(c, err)
		return
	}
	if h.broker == nil {
		runtimeLogError(c, diagnostics.ErrClosed)
		return
	}
	snapshot, err := h.broker.Snapshot(query)
	if err != nil {
		runtimeLogError(c, err)
		return
	}
	c.JSON(http.StatusOK, snapshot)
}

func runtimeLogError(c *gin.Context, err error) {
	status, code := http.StatusServiceUnavailable, "logs_unavailable"
	if errors.Is(err, diagnostics.ErrInvalidQuery) || errors.Is(err, diagnostics.ErrInvalidCursor) {
		status, code = http.StatusBadRequest, "logs_invalid_query"
	} else if errors.Is(err, diagnostics.ErrCapacity) {
		code = "logs_capacity"
	}
	c.AbortWithStatusJSON(status, gin.H{"code": code})
}

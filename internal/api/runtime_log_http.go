package api

import (
	"net/http"
	"time"

	"go-gateway/internal/diagnostics"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const managedRequestID = "managed_diagnostic_request_id"

func managedHTTPMiddleware(broker *diagnostics.Broker, base string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if broker == nil || runtimeLogPath(c.Request.URL.Path, base) {
			c.Next()
			return
		}
		start := time.Now()
		c.Set(managedRequestID, uuid.NewString())
		c.Next()
		fields := managedHTTPFields(c)
		fields["duration_ms"] = float64(time.Since(start)) / float64(time.Millisecond)
		broker.Emit(diagnostics.Input{Code: "http.access", Fields: fields})
	}
}

func managedRecoveryMiddleware(broker *diagnostics.Broker) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.AbortWithStatus(http.StatusInternalServerError)
		if broker != nil {
			broker.Emit(diagnostics.Input{Code: "http.recovered", Fields: managedHTTPFields(c)})
		}
	})
}

func managedHTTPFields(c *gin.Context) map[string]any {
	route := c.FullPath()
	if route == "" {
		route = "unmatched"
	}
	return map[string]any{"method": c.Request.Method, "route": route, "status": c.Writer.Status(), "request_id": c.GetString(managedRequestID)}
}

func registerManagedHTTPRoutes(router *gin.Engine, broker *diagnostics.Broker) {
	if broker == nil {
		return
	}
	routes := router.Routes()
	paths := make([]string, 0, len(routes))
	for _, route := range routes {
		paths = append(paths, route.Path)
	}
	broker.RegisterRoutes(paths)
}

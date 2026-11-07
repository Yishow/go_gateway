package handlers

import (
	"errors"
	"net/http"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// WriteGroupDeleteData reports the permanently removed local generations.
type WriteGroupDeleteData struct {
	Deleted            bool   `json:"deleted"`
	GroupID            string `json:"group_id"`
	LineageID          string `json:"lineage_id"`
	DeletedGenerations int    `json:"deleted_generations"`
}

// WriteGroupDeleteResponse is the permanent-delete success envelope.
type WriteGroupDeleteResponse struct {
	Success bool                 `json:"success"`
	Data    WriteGroupDeleteData `json:"data"`
}

func renderPermanentDeleteError(c *gin.Context, err error) bool {
	var refusal *workspace.WriteGroupDeleteError
	if !errors.As(err, &refusal) {
		return false
	}
	action, message := "open_recording", "recording cannot be deleted yet"
	switch refusal.Code {
	case "WRITE_GROUP_DELETE_BACKLOG":
		action, message = "stop_and_wait", "recording has delivery data still waiting"
	case "WRITE_GROUP_DELETE_ACTIVE":
		action, message = "stop_first", "recording is still accepting data"
	case "WRITE_GROUP_DELETE_OPERATION_PENDING":
		message = "recording has an unfinished operation"
	case "WRITE_GROUP_DELETE_HAS_SUCCESSOR":
		message = "another recording still depends on this recording"
	default:
		return false
	}
	envelope := gin.H{"code": refusal.Code, "message": message, "retryable": false, "request_id": getOrGenerateRequestID(c), "action": action}
	if refusal.Details != nil {
		envelope["details"] = refusal.Details
	}
	c.JSON(http.StatusConflict, gin.H{"success": false, "error": envelope})
	return true
}

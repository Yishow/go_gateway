package handlers

import (
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

type studioV2WorkspaceMappingRequest struct {
	RuleID      string          `json:"rule_id"`
	Address     string          `json:"address"`
	TagKey      string          `json:"tag_key"`
	DisplayName string          `json:"display_name"`
	Unit        string          `json:"unit"`
	TargetType  schema.DataType `json:"target_type"`
	Scale       float64         `json:"scale"`
	Offset      float64         `json:"offset"`
	Enabled     bool            `json:"enabled"`
}

type studioV2WorkspaceMappingResponse struct {
	ID                  string                     `json:"id"`
	WorkspaceID         string                     `json:"workspace_id"`
	PointID             string                     `json:"point_id"`
	RuleID              string                     `json:"rule_id"`
	DeviceID            string                     `json:"device_id"`
	Address             string                     `json:"address"`
	TagID               string                     `json:"tag_id"`
	TagKey              string                     `json:"tag_key"`
	DisplayName         string                     `json:"display_name"`
	Unit                string                     `json:"unit"`
	TargetType          string                     `json:"target_type"`
	Scale               float64                    `json:"scale"`
	Offset              float64                    `json:"offset"`
	Enabled             bool                       `json:"enabled"`
	SaveState           string                     `json:"save_state"`
	RuntimeApplyStatus  string                     `json:"runtime_apply_status,omitempty"`
	RuntimeApplyMessage string                     `json:"runtime_apply_message,omitempty"`
	RuntimeApplyIssues  []workspace.ReadinessIssue `json:"runtime_apply_issues,omitempty"`
	CreatedAt           string                     `json:"created_at"`
	UpdatedAt           string                     `json:"updated_at"`
}

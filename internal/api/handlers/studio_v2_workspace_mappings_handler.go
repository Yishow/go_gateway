package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

const (
	ruleManagedTagLabelSource  = "source"
	ruleManagedTagLabelRuleID  = "source_rule_id"
	ruleManagedTagLabelAddress = "source_rule_address"
	ruleManagedTagLabelValue   = "source-rule"
)

type StudioV2WorkspaceMappingsHandler struct {
	workspaceSvc *workspace.Service
	deviceSvc    *device.Service
	ruleSvc      *sourcerule.Service
	pointSvc     *point.Service
	tagSvc       *tag.Service
	mappingSvc   *mapping.Service
}

func NewStudioV2WorkspaceMappingsHandler(workspaceSvc *workspace.Service, deviceSvc *device.Service, ruleSvc *sourcerule.Service, pointSvc *point.Service, tagSvc *tag.Service, mappingSvc *mapping.Service) *StudioV2WorkspaceMappingsHandler {
	if ruleSvc != nil && tagSvc != nil && mappingSvc != nil {
		ruleSvc.SetTagMappingServices(tagSvc, mappingSvc)
	}
	return &StudioV2WorkspaceMappingsHandler{
		workspaceSvc: workspaceSvc,
		deviceSvc:    deviceSvc,
		ruleSvc:      ruleSvc,
		pointSvc:     pointSvc,
		tagSvc:       tagSvc,
		mappingSvc:   mappingSvc,
	}
}

func (h *StudioV2WorkspaceMappingsHandler) List(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	rules, err := h.ruleSvc.ListByDeviceIDs(c.Request.Context(), record.OrderedDeviceIDs)
	if err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}

	payload := make([]studioV2WorkspaceMappingResponse, 0)
	for _, rule := range rules {
		items, err := h.listWorkspaceRuleMappings(c.Request.Context(), record.ID, rule.ID)
		if err != nil {
			renderStudioV2WorkspaceMappingError(c, err)
			return
		}
		payload = append(payload, items...)
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

func (h *StudioV2WorkspaceMappingsHandler) Create(c *gin.Context) {
	row, ok := h.resolveRequestRow(c)
	if !ok {
		return
	}

	defer row.release()
	record, rule, links, link, req := row.record, row.rule, row.links, row.link, row.request
	if mappingRecord, found, _, err := h.recoverWorkspaceLinkMappingForMutation(c.Request.Context(), link); err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	} else if found {
		var cleanupStatus string
		mappingRecord, cleanupStatus, err = h.saveExistingWorkspaceMapping(c.Request.Context(), rule, links, link, mappingRecord, req)
		if err != nil {
			renderStudioV2WorkspaceMappingError(c, err)
			return
		}
		payload, err := h.buildResponse(c.Request.Context(), record.ID, rule, link, mappingRecord.ID, mappingRecord.TagID)
		if err != nil {
			renderStudioV2WorkspaceMappingError(c, err)
			return
		}
		payload.CleanupStatus = cleanupStatus
		row.release()
		applyOutcome := resolveStudioV2ScopedRuntimeApplyOutcome(c.Request.Context(), h.workspaceSvc, h.deviceSvc, []string{rule.DeviceID}, []string{rule.DeviceID, link.PointID})
		payload.RuntimeApplyStatus = applyOutcome.Status
		payload.RuntimeApplyMessage = applyOutcome.Message
		payload.RuntimeApplyIssues = applyOutcome.Issues
		c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
		return
	}

	pointRecord, err := h.pointSvc.GetByID(c.Request.Context(), link.PointID)
	if err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}
	tagRecord, err := h.upsertTag(c.Request.Context(), rule.ID, req, link.Address)
	if err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}
	mappingRecord, err := h.mappingSvc.Create(c.Request.Context(), mapping.CreateMappingRequest{
		PointID:           pointRecord.ID,
		TagID:             tagRecord.ID,
		Enabled:           common.Ptr(false),
		TransformPipeline: buildWorkspaceMappingPipeline(pointRecord.DataType, req.TargetType, req.Scale, req.Offset),
	})
	if err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}

	link.TagID = common.Ptr(tagRecord.ID)
	link.MappingID = common.Ptr(mappingRecord.ID)
	link.UpdatedAt = time.Now()
	if err := h.ruleSvc.ReplaceLinks(c.Request.Context(), rule.ID, links); err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}
	if err := h.ruleSvc.ConfirmWorkspaceMapping(c.Request.Context(), rule.ID, rule.RevisionID, mappingRecord, req.TargetType); err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}

	payload, err := h.buildResponse(c.Request.Context(), record.ID, rule, link, mappingRecord.ID, tagRecord.ID)
	if err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}
	row.release()
	applyOutcome := resolveStudioV2ScopedRuntimeApplyOutcome(c.Request.Context(), h.workspaceSvc, h.deviceSvc, []string{rule.DeviceID}, []string{rule.DeviceID, link.PointID})
	payload.RuntimeApplyStatus = applyOutcome.Status
	payload.RuntimeApplyMessage = applyOutcome.Message
	payload.RuntimeApplyIssues = applyOutcome.Issues
	c.JSON(http.StatusCreated, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

func (h *StudioV2WorkspaceMappingsHandler) Update(c *gin.Context) {
	row, ok := h.requireWorkspaceMapping(c)
	if !ok {
		return
	}

	defer row.release()
	record, rule, links, link, mappingRecord, req := row.record, row.rule, row.links, row.link, row.mapping, row.request
	if req.RuleID != "" && req.RuleID != rule.ID {
		renderStudioV2WorkspaceValidationError(c, errors.New("workspace mapping ownership mismatch"))
		return
	}
	if req.Address != "" && normalizeWorkspaceAddress(req.Address) != normalizeWorkspaceAddress(link.Address) {
		renderStudioV2WorkspaceValidationError(c, errors.New("workspace mapping ownership mismatch"))
		return
	}

	mappingRecord, cleanupStatus, err := h.saveExistingWorkspaceMapping(c.Request.Context(), rule, links, link, mappingRecord, req)
	if err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}

	payload, err := h.buildResponse(c.Request.Context(), record.ID, rule, link, mappingRecord.ID, mappingRecord.TagID)
	if err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}
	payload.CleanupStatus = cleanupStatus
	row.release()
	applyOutcome := resolveStudioV2ScopedRuntimeApplyOutcome(c.Request.Context(), h.workspaceSvc, h.deviceSvc, []string{rule.DeviceID}, []string{rule.DeviceID, link.PointID})
	payload.RuntimeApplyStatus = applyOutcome.Status
	payload.RuntimeApplyMessage = applyOutcome.Message
	payload.RuntimeApplyIssues = applyOutcome.Issues
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

func (h *StudioV2WorkspaceMappingsHandler) buildResponse(ctx context.Context, workspaceID string, rule *schema.SourceRule, link *schema.SourceRuleLink, mappingID, tagID string) (studioV2WorkspaceMappingResponse, error) {
	mappingRecord, err := h.mappingSvc.GetByID(ctx, mappingID)
	if err != nil {
		return studioV2WorkspaceMappingResponse{}, err
	}
	tagRecord, err := h.tagSvc.GetByID(ctx, tagID)
	if err != nil {
		return studioV2WorkspaceMappingResponse{}, err
	}
	pointRecord, err := h.pointSvc.GetByID(ctx, link.PointID)
	if err != nil {
		return studioV2WorkspaceMappingResponse{}, err
	}
	targetType, scale, offset, err := decodeWorkspaceMappingPipeline(pointRecord.DataType, mappingRecord.TransformPipeline)
	if err != nil {
		return studioV2WorkspaceMappingResponse{}, err
	}
	return studioV2WorkspaceMappingResponse{
		ID:          mappingRecord.ID,
		WorkspaceID: workspaceID,
		PointID:     pointRecord.ID,
		RuleID:      rule.ID,
		DeviceID:    rule.DeviceID,
		Address:     link.Address,
		TagID:       tagRecord.ID,
		TagKey:      tagRecord.Key,
		DisplayName: tagRecord.DisplayName,
		Unit:        tagRecord.Unit,
		TargetType:  string(targetType),
		Scale:       scale,
		Offset:      offset,
		Enabled:     mappingRecord.Enabled,
		SaveState:   workspaceSaveStateSaved,
		CreatedAt:   mappingRecord.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:   mappingRecord.UpdatedAt.Format(time.RFC3339Nano),
	}, nil
}

func (h *StudioV2WorkspaceMappingsHandler) upsertTag(ctx context.Context, ruleID string, req studioV2WorkspaceMappingRequest, address string) (*schema.Tag, error) {
	key := tag.NormalizeTagKey(req.TagKey)
	tagRecord, err := h.tagSvc.GetByKey(ctx, key)
	if errors.Is(err, tag.ErrTagNotFound) {
		return h.tagSvc.Create(ctx, tag.CreateTagRequest{
			Key:         key,
			DisplayName: req.DisplayName,
			Unit:        req.Unit,
			DataType:    req.TargetType,
			Labels: map[string]string{
				ruleManagedTagLabelSource:  ruleManagedTagLabelValue,
				ruleManagedTagLabelRuleID:  ruleID,
				ruleManagedTagLabelAddress: normalizeWorkspaceAddress(address),
			},
		})
	}
	if err != nil {
		return nil, err
	}
	if !isRuleManagedWorkspaceTagOwnedBy(tagRecord, ruleID, address) {
		return nil, fmt.Errorf("%w: workspace mapping tag ownership mismatch", mapping.ErrValidation)
	}
	return h.tagSvc.Update(ctx, tagRecord.ID, tag.UpdateTagRequest{
		DisplayName: common.Ptr(req.DisplayName),
		Unit:        common.Ptr(req.Unit),
		DataType:    &req.TargetType,
	})
}

func (h *StudioV2WorkspaceMappingsHandler) deleteOrphanRuleManagedTag(ctx context.Context, tagID string) error {
	tagRecord, err := h.tagSvc.GetByID(ctx, tagID)
	if errors.Is(err, tag.ErrTagNotFound) {
		return nil
	}
	if err != nil || !isRuleManagedWorkspaceTag(tagRecord) {
		return err
	}
	mappings, err := h.mappingSvc.List(ctx, mapping.ListFilter{TagID: &tagID})
	if err != nil || len(mappings) > 0 {
		return err
	}
	if err := h.tagSvc.Delete(ctx, tagID); err != nil && !errors.Is(err, tag.ErrTagNotFound) {
		return err
	}
	return nil
}

func validateWorkspaceMappingRequest(req studioV2WorkspaceMappingRequest) error {
	switch {
	case strings.TrimSpace(req.RuleID) == "":
		return errors.New("rule_id is required")
	case strings.TrimSpace(req.Address) == "":
		return errors.New("address is required")
	case strings.TrimSpace(req.TagKey) == "":
		return errors.New("tag_key is required")
	case !req.TargetType.IsValid():
		return errors.New("target_type is invalid")
	default:
		return nil
	}
}

func buildWorkspaceMappingPipeline(pointType, targetType schema.DataType, scale, offset float64) []schema.TransformStep {
	steps := make([]schema.TransformStep, 0, 2)
	scaled := scale != 1 || offset != 0
	if scaled {
		steps = append(steps, schema.TransformStep{Type: schema.TransformScale, Order: len(steps), Params: map[string]any{"scale": scale, "offset": offset}})
	}
	// Only newly authored pipelines use this order. ExecutePipeline retains the
	// stored order of existing pipelines, including historical cast-before-scale.
	if targetType != pointType || scaled {
		steps = append(steps, schema.TransformStep{Type: schema.TransformCast, Order: len(steps), Params: map[string]any{"target_type": string(targetType)}})
	}
	return steps
}

func decodeWorkspaceMappingPipeline(pointType schema.DataType, raw string) (decodedType schema.DataType, decodedScale, decodedOffset float64, decodeErr error) {
	targetType := pointType
	scale, offset := 1.0, 0.0
	if strings.TrimSpace(raw) == "" {
		return targetType, scale, offset, nil
	}
	var steps []schema.TransformStep
	if err := json.Unmarshal([]byte(raw), &steps); err != nil {
		return "", 0, 0, err
	}
	for _, step := range steps {
		switch step.Type {
		case schema.TransformCast:
			if value, ok := step.Params["target_type"].(string); ok && schema.DataType(value).IsValid() {
				targetType = schema.DataType(value)
			}
		case schema.TransformScale:
			if value, ok := step.Params["scale"].(float64); ok {
				scale = value
			}
			if value, ok := step.Params["offset"].(float64); ok {
				offset = value
			}
		}
	}
	return targetType, scale, offset, nil
}

func normalizeWorkspaceAddress(address string) string {
	return strings.ToUpper(strings.TrimSpace(address))
}

func isRuleManagedWorkspaceTag(tagRecord *schema.Tag) bool {
	labels, ok := decodeRuleManagedWorkspaceTagLabels(tagRecord)
	if !ok {
		return false
	}
	return labels[ruleManagedTagLabelSource] == ruleManagedTagLabelValue
}

func isRuleManagedWorkspaceTagOwnedBy(tagRecord *schema.Tag, ruleID, address string) bool {
	labels, ok := decodeRuleManagedWorkspaceTagLabels(tagRecord)
	if !ok {
		return false
	}
	if labels[ruleManagedTagLabelSource] != ruleManagedTagLabelValue {
		return false
	}
	if labels[ruleManagedTagLabelRuleID] != strings.TrimSpace(ruleID) {
		return false
	}
	return normalizeWorkspaceAddress(labels[ruleManagedTagLabelAddress]) == normalizeWorkspaceAddress(address)
}

func isLegacyDirectWorkspaceBinding(link *schema.SourceRuleLink, mappingRecord *schema.Mapping, tagRecord *schema.Tag) bool {
	if link == nil || mappingRecord == nil || tagRecord == nil {
		return false
	}
	if strings.TrimSpace(tagRecord.Labels) != "" {
		return false
	}
	if link.MappingID == nil || link.TagID == nil {
		return false
	}
	return strings.TrimSpace(*link.MappingID) == strings.TrimSpace(mappingRecord.ID) &&
		strings.TrimSpace(*link.TagID) == strings.TrimSpace(tagRecord.ID) &&
		strings.TrimSpace(mappingRecord.PointID) == strings.TrimSpace(link.PointID) &&
		strings.TrimSpace(mappingRecord.TagID) == strings.TrimSpace(tagRecord.ID)
}

func decodeRuleManagedWorkspaceTagLabels(tagRecord *schema.Tag) (map[string]string, bool) {
	if tagRecord == nil || strings.TrimSpace(tagRecord.Labels) == "" {
		return nil, false
	}
	var labels map[string]string
	if err := json.Unmarshal([]byte(tagRecord.Labels), &labels); err != nil {
		return nil, false
	}
	return labels, true
}

package sourcerule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

type pointUpdatePlan struct {
	appliedRevision string
	name            string
	point           *schema.Point
}

type mappingRollbackState struct {
	enabled              bool
	transformPipeline    []schema.TransformStep
	status               schema.MappingStatus
	ruleCandidateID      string
	proposedSignature    string
	lastAppliedSignature string
	blockingReason       string
}

type tagRollbackState struct {
	dataType schema.DataType
}

type tagMappingSyncResult struct {
	createdTagIDs     []string
	createdMappingIDs []string
	updatedMappings   map[string]mappingRollbackState
	updatedTags       map[string]tagRollbackState
}

const (
	ruleManagedTagLabelSource  = "source"
	ruleManagedTagLabelRuleID  = "source_rule_id"
	ruleManagedTagLabelAddress = "source_rule_address"
	ruleManagedTagLabelValue   = "source-rule"
	ruleOriginManual           = "manual"
)

type CreateRuleRequest struct {
	ID               string          `json:"id,omitempty"`
	DeviceID         string          `json:"device_id"`
	StartAddress     string          `json:"start_address"`
	Count            int             `json:"count"`
	DataType         schema.DataType `json:"data_type"`
	NamingPrefix     string          `json:"naming_prefix"`
	Enabled          bool            `json:"enabled"`
	Locked           bool            `json:"locked"`
	Origin           string          `json:"origin,omitempty"`
	TemplateName     string          `json:"template_name,omitempty"`
	SkippedAddresses []string        `json:"skipped_addresses,omitempty"`
	// TargetDataType 目標資料型態（可空）。預設與 DataType 相同。
	TargetDataType *schema.DataType `json:"target_data_type,omitempty"`
	// ScaleMultiplier 縮放倍率（可空）。
	ScaleMultiplier *float64 `json:"scale_multiplier,omitempty"`
	// ScaleOffset 偏移量（可空）。
	ScaleOffset *float64 `json:"scale_offset,omitempty"`
	// DataFormat 字節序格式（可空）。空字串表示使用設備連線預設。
	DataFormat string `json:"data_format,omitempty"`
	// ShareEnabled controls whether this rule is exposed through the local Modbus Share output.
	ShareEnabled bool `json:"share_enabled"`
	// ShareStartRegister is the optional local Modbus holding-register start address.
	ShareStartRegister *int `json:"share_start_register"`
	// ShareStride is the optional number of local Modbus registers allocated per point.
	ShareStride *int `json:"share_stride"`
}

type UpdateRuleRequest struct {
	StartAddress     *string          `json:"start_address,omitempty"`
	Count            *int             `json:"count,omitempty"`
	DataType         *schema.DataType `json:"data_type,omitempty"`
	NamingPrefix     *string          `json:"naming_prefix,omitempty"`
	Enabled          *bool            `json:"enabled,omitempty"`
	Locked           *bool            `json:"locked,omitempty"`
	TemplateName     *string          `json:"template_name,omitempty"`
	SkippedAddresses *[]string        `json:"skipped_addresses,omitempty"`
	// TargetDataType 目標資料型態（可空）。
	TargetDataType    *schema.DataType `json:"target_data_type,omitempty"`
	TargetDataTypeSet bool             `json:"-"`
	// ScaleMultiplier 縮放倍率（可空）。
	ScaleMultiplier    *float64 `json:"scale_multiplier,omitempty"`
	ScaleMultiplierSet bool     `json:"-"`
	// ScaleOffset 偏移量（可空）。
	ScaleOffset    *float64 `json:"scale_offset,omitempty"`
	ScaleOffsetSet bool     `json:"-"`
	// DataFormat 字節序格式（可空）；傳 null 並搭配 DataFormatSet 可清空為連線預設。
	DataFormat    *string `json:"data_format,omitempty"`
	DataFormatSet bool    `json:"-"`
	// ShareEnabled controls whether this rule is exposed through the local Modbus Share output.
	ShareEnabled *bool `json:"share_enabled,omitempty"`
	// ShareStartRegister is the optional local Modbus holding-register start address.
	ShareStartRegister    *int `json:"share_start_register,omitempty"`
	ShareStartRegisterSet bool `json:"-"`
	// ShareStride is the optional number of local Modbus registers allocated per point.
	ShareStride    *int `json:"share_stride,omitempty"`
	ShareStrideSet bool `json:"-"`
}

func NewService(repo Repository, deviceSvc *device.Service, pointSvc *point.Service, runtimeSync RuntimeSyncer) *Service {
	return &Service{
		repo:        repo,
		deviceSvc:   deviceSvc,
		pointSvc:    pointSvc,
		runtimeSync: runtimeSync,
	}
}
func (s *Service) SetTagMappingServices(tagSvc *tag.Service, mappingSvc *mapping.Service) {
	s.tagSvc = tagSvc
	s.mappingSvc = mappingSvc
}

func (s *Service) Create(ctx context.Context, req CreateRuleRequest) (result *schema.SourceRule, err error) {
	if err := validateCreateRequest(req); err != nil {
		return nil, contextualRuleValidationError(req.ID, req.StartAddress, "", err)
	}

	deviceRecord, err := s.deviceSvc.GetByID(ctx, req.DeviceID)
	if err != nil {
		return nil, wrapSourceRuleDeviceError(req.DeviceID, err)
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		id, err = common.NewUUID()
		if err != nil {
			return nil, fmt.Errorf("建立來源規則 ID 失敗: %w", err)
		}
	}

	ctx, release := s.AcquireRuleMutation(ctx, id)
	defer release()
	ctx, finishAcquisition := s.beginSourceAcquisitionMutation(ctx, id, req.DeviceID)
	defer finishAcquisition()
	defer func(finalizeCtx context.Context) { err = s.finalizePointWrites(finalizeCtx, err) }(ctx)

	plannedAddresses, err := buildPlannedPointAddresses(strings.TrimSpace(strings.ToUpper(req.StartAddress)), req.Count, req.DataType, deviceRecord.Protocol)
	if err != nil {
		return nil, contextualRuleValidationError(id, req.StartAddress, deviceRecord.Protocol, err)
	}
	if err := s.preflightPointAddresses(ctx, req.DeviceID, nil, plannedAddresses, req.SkippedAddresses); err != nil {
		return nil, err
	}
	enabled, err := s.initialRuleEnabled(ctx, req.Enabled, req.DeviceID)
	if err != nil {
		return nil, contextualRuleValidationError(id, req.StartAddress, deviceRecord.Protocol, err)
	}
	skippedJSON, err := marshalSkippedAddresses(req.SkippedAddresses)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	rule := &schema.SourceRule{
		ID:                 id,
		DeviceID:           req.DeviceID,
		StartAddress:       strings.TrimSpace(strings.ToUpper(req.StartAddress)),
		Count:              req.Count,
		DataType:           req.DataType,
		NamingPrefix:       normalizeNamingPrefix(req.NamingPrefix),
		Enabled:            enabled,
		Locked:             req.Locked,
		Origin:             normalizeOrigin(req.Origin),
		TemplateName:       strings.TrimSpace(req.TemplateName),
		SkippedAddresses:   skippedJSON,
		TargetDataType:     req.TargetDataType,
		ScaleMultiplier:    req.ScaleMultiplier,
		ScaleOffset:        req.ScaleOffset,
		DataFormat:         normalizeDataFormat(req.DataFormat),
		ShareEnabled:       req.ShareEnabled,
		ShareStartRegister: cloneIntPtr(req.ShareStartRegister),
		ShareStride:        cloneIntPtr(req.ShareStride),
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := assignRuleRevisionID(rule); err != nil {
		return nil, err
	}
	s.noteRuleBaseline(ctx, rule)
	if err := s.repo.Create(ctx, rule); err != nil {
		return nil, s.dirtyUnknown(ctx)
	}
	createdPointIDs := make([]string, 0, req.Count)
	links := make([]*schema.SourceRuleLink, 0, req.Count)
	createdPoints := make([]*schema.Point, 0, req.Count)
	defaultPollingGroupID, err := s.resolveDerivedPointPollingGroupID(ctx)
	if err != nil {
		if rollbackErr := s.rollbackCreatedRule(ctx, rule.ID, rule.RevisionID, nil, tagMappingSyncResult{}); rollbackErr != nil {
			return nil, s.dirtyUnknown(ctx)
		}
		return nil, err
	}

	for _, address := range plannedAddresses {
		if containsAddress(req.SkippedAddresses, address) {
			continue
		}

		pointRecord, createErr := s.createRulePoint(ctx, point.CreatePointRequest{
			DeviceID:       rule.DeviceID,
			Name:           buildPointName(rule.NamingPrefix, address),
			Address:        address,
			DataType:       rule.DataType,
			DataFormat:     rule.DataFormat,
			Mode:           schema.PointModeReadOnly,
			Function:       "",
			Description:    "",
			PollingGroupID: defaultPollingGroupID,
		})
		if createErr != nil {
			if rollbackErr := s.rollbackCreatedRule(ctx, rule.ID, rule.RevisionID, createdPointIDs, tagMappingSyncResult{}); rollbackErr != nil {
				return nil, s.dirtyUnknown(ctx)
			}
			return nil, fmt.Errorf("建立規則衍生點位失敗: %w", createErr)
		}

		if !enabled {
			createdPointID := pointRecord.ID
			pointRecord, createErr = s.updateRulePoint(ctx, pointRecord, point.UpdatePointRequest{Enabled: common.Ptr(false)})
			if createErr != nil {
				if rollbackErr := s.rollbackCreatedRule(ctx, rule.ID, rule.RevisionID, append(createdPointIDs, createdPointID), tagMappingSyncResult{}); rollbackErr != nil {
					return nil, s.dirtyUnknown(ctx)
				}
				return nil, fmt.Errorf("設定規則衍生點位狀態失敗: %w", createErr)
			}
		}

		linkID, linkErr := common.NewUUID()
		if linkErr != nil {
			if rollbackErr := s.rollbackCreatedRule(ctx, rule.ID, rule.RevisionID, append(createdPointIDs, pointRecord.ID), tagMappingSyncResult{}); rollbackErr != nil {
				return nil, s.dirtyUnknown(ctx)
			}
			return nil, fmt.Errorf("建立來源規則連結 ID 失敗: %w", linkErr)
		}

		createdPointIDs = append(createdPointIDs, pointRecord.ID)
		createdPoints = append(createdPoints, pointRecord)
		links = append(links, &schema.SourceRuleLink{
			ID:        linkID,
			RuleID:    rule.ID,
			Address:   address,
			PointID:   pointRecord.ID,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	syncResult, err := s.syncRuleTagMappings(ctx, nil, rule, links, rule.Enabled)
	if err != nil {
		s.noteTransactionFailure(ctx, err)
		if rollbackErr := s.rollbackCreatedRule(ctx, rule.ID, rule.RevisionID, createdPointIDs, syncResult); rollbackErr != nil {
			return nil, s.dirtyUnknown(ctx)
		}
		return nil, fmt.Errorf("同步來源規則標籤映射失敗: %w", err)
	}

	if len(links) > 0 {
		if err := s.repo.CreateLinks(ctx, links); err != nil {
			s.noteTransactionFailure(ctx, err)
			if rollbackErr := s.rollbackCreatedRule(ctx, rule.ID, rule.RevisionID, createdPointIDs, syncResult); rollbackErr != nil {
				return nil, s.dirtyUnknown(ctx)
			}
			return nil, fmt.Errorf("建立來源規則連結失敗: %w", err)
		}
	}
	s.syncPoints(ctx, createdPoints)
	if err := s.persistCandidateSnapshots(ctx, rule, links); err != nil {
		if rollbackErr := s.rollbackCreatedRule(ctx, rule.ID, rule.RevisionID, createdPointIDs, syncResult); rollbackErr != nil {
			return nil, s.dirtyUnknown(ctx)
		}
		return nil, fmt.Errorf("持久化來源規則候選快照失敗: %w", err)
	}
	return rule, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*schema.SourceRule, error) {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	return rule, nil
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]*schema.SourceRule, error) {
	rules, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則失敗: %w", err)
	}
	return rules, nil
}

func (s *Service) ListByDeviceIDs(ctx context.Context, deviceIDs []string) ([]*schema.SourceRule, error) {
	result := make([]*schema.SourceRule, 0)

	for _, deviceID := range deviceIDs {
		rules, err := s.List(ctx, ListFilter{DeviceID: &deviceID})
		if err != nil {
			return nil, err
		}
		result = append(result, rules...)
	}

	return result, nil
}

func (s *Service) ListLinks(ctx context.Context, ruleID string) ([]*schema.SourceRuleLink, error) {
	links, err := s.repo.ListLinks(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則連結失敗: %w", err)
	}
	return links, nil
}

func (s *Service) Delete(ctx context.Context, id string) (err error) {
	ctx, release := s.AcquireRuleMutation(ctx, id)
	defer release()
	defer func(finalizeCtx context.Context) { err = s.finalizePointWrites(finalizeCtx, err) }(ctx)
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得來源規則失敗: %w", err)
	}
	ctx, finishAcquisition := s.beginSourceAcquisitionMutation(ctx, id, rule.DeviceID)
	defer finishAcquisition()
	links, err := s.repo.ListLinks(ctx, id)
	if err != nil {
		return fmt.Errorf("取得來源規則連結失敗: %w", err)
	}
	if err := s.deleteOwnedRule(ctx, rule); err != nil {
		return s.finishRollbackErrors(ctx, err)
	}
	for _, link := range links {
		if err := s.cleanupRuleLinkResources(ctx, rule, link); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Disable(ctx context.Context, id string) error {
	return s.setEnabled(ctx, id, false)
}

func (s *Service) Enable(ctx context.Context, id string) error {
	return s.setEnabled(ctx, id, true)
}

func (s *Service) syncRuleLinksEnabled(ctx context.Context, ruleID string, enabled bool) ([]pointUpdatePlan, error) {
	links, err := s.repo.ListLinks(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則連結失敗: %w", err)
	}
	defaultPollingGroupID, err := s.resolveDerivedPointPollingGroupID(ctx)
	if err != nil {
		return nil, err
	}

	appliedPlans := make([]pointUpdatePlan, 0, len(links))
	for _, link := range links {
		pointRecord, getErr := s.pointSvc.GetByID(ctx, link.PointID)
		if getErr != nil {
			return nil, s.finishRollbackErrors(ctx, fmt.Errorf("取得衍生點位失敗: %w", getErr), s.rollbackUpdatedPoints(ctx, appliedPlans))
		}
		original := pointRecord
		updateReq := point.UpdatePointRequest{
			Enabled: &enabled,
		}
		if enabled && pointRecord.PollingGroupID == nil && defaultPollingGroupID != nil {
			updateReq.PollingGroupID = defaultPollingGroupID
		}
		if s.syncUnchangedRulePoint(ctx, original, updateReq) {
			continue
		}
		pointRecord, getErr = s.updateRulePoint(ctx, original, updateReq)
		if getErr != nil {
			return nil, s.finishRollbackErrors(ctx, fmt.Errorf("更新衍生點位狀態失敗: %w", getErr), s.rollbackUpdatedPoints(ctx, appliedPlans))
		}
		appliedPlans = append(appliedPlans, pointUpdatePlan{
			appliedRevision: pointRecord.RevisionID,
			name:            original.Name,
			point:           original,
		})
		s.syncPoints(ctx, []*schema.Point{pointRecord})
	}
	return appliedPlans, nil
}

func (s *Service) setEnabled(ctx context.Context, id string, enabled bool) (err error) {
	ctx, release := s.AcquireRuleMutation(ctx, id)
	defer release()
	defer func(finalizeCtx context.Context) { err = s.finalizePointWrites(finalizeCtx, err) }(ctx)
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得來源規則失敗: %w", err)
	}
	ctx, finishAcquisition := s.beginSourceAcquisitionMutation(ctx, id, rule.DeviceID)
	defer finishAcquisition()

	previousRule := *rule
	previousLinks, err := s.repo.ListLinks(ctx, id)
	if err != nil {
		return fmt.Errorf("取得來源規則連結失敗: %w", err)
	}
	if err := s.captureLinkBaselines(ctx, previousLinks); err != nil {
		return err
	}
	if err := s.ensureRuleActivationAllowed(ctx, rule.DeviceID, enabled); err != nil {
		return err
	}
	rule.Enabled = enabled
	rule.UpdatedAt = time.Now()
	if err := s.updateOwnedRule(ctx, &previousRule, rule); err != nil {
		return fmt.Errorf("更新來源規則失敗: %w", err)
	}
	appliedPointPlans, err := s.syncRuleLinksEnabled(ctx, id, enabled)
	if err != nil {
		return s.finishRollbackErrors(ctx, err, s.restoreCurrentOwnedRule(ctx, &previousRule))
	}

	nextLinks := cloneSourceRuleLinks(previousLinks)
	syncResult, syncErr := s.syncRuleTagMappings(ctx, rule, rule, nextLinks, enabled)
	if syncErr != nil {
		if s.noteTransactionFailure(ctx, syncErr) {
			return s.dirtyUnknown(ctx)
		}
		return s.finishRollbackErrors(ctx, fmt.Errorf("同步來源規則標籤映射失敗: %w", syncErr), s.rollbackTagMappingSync(ctx, syncResult), s.rollbackUpdatedPoints(ctx, appliedPointPlans), s.restoreCurrentOwnedRule(ctx, &previousRule))
	}

	if err := s.replaceRuleLinks(ctx, id, previousLinks, nextLinks); err != nil {
		return s.finishRollbackErrors(ctx, err, s.rollbackTagMappingSync(ctx, syncResult), s.rollbackUpdatedPoints(ctx, appliedPointPlans), s.restoreCurrentOwnedRule(ctx, &previousRule))
	}

	return nil
}

// resolveDerivedPointPollingGroupID 取得來源規則衍生點位應使用的預設輪詢群組。
func (s *Service) resolveDerivedPointPollingGroupID(ctx context.Context) (*string, error) {
	if s.pointSvc == nil {
		return nil, nil
	}

	group, err := s.pointSvc.EnsureDefaultPollingGroup(ctx)
	if err != nil {
		return nil, fmt.Errorf("確保預設輪詢群組失敗: %w", err)
	}
	if group == nil {
		return nil, nil
	}

	groupID := group.ID
	return &groupID, nil
}

func (s *Service) syncRuleManagedTagDataType(
	ctx context.Context,
	rule *schema.SourceRule,
	link *schema.SourceRuleLink,
	tagRecord *schema.Tag,
	targetDataType schema.DataType,
	result *tagMappingSyncResult,
) (*schema.Tag, error) {
	if tagRecord.DataType == targetDataType {
		return tagRecord, nil
	}
	if !isRuleManagedTag(tagRecord, rule.ID, link.Address) {
		return nil, fmt.Errorf("來源規則衍生標籤資料型別衝突: tag=%s 目標=%s", tagRecord.DataType, targetDataType)
	}
	if result != nil {
		if result.updatedTags == nil {
			result.updatedTags = make(map[string]tagRollbackState)
		}
		if _, exists := result.updatedTags[tagRecord.ID]; !exists {
			result.updatedTags[tagRecord.ID] = tagRollbackState{dataType: tagRecord.DataType}
		}
	}

	updated, err := s.tagSvc.Update(ctx, tagRecord.ID, tag.UpdateTagRequest{
		DataType: &targetDataType,
	})
	if err != nil {
		return nil, fmt.Errorf("同步來源規則標籤資料型別失敗: %w", err)
	}
	return updated, nil
}

func (s *Service) validateTagAvailability(ctx context.Context, tagRecord *schema.Tag, pointID string) error {
	if tagRecord == nil {
		return fmt.Errorf("來源規則標籤不存在")
	}
	tagID := tagRecord.ID
	tagMappings, err := s.mappingSvc.List(ctx, mapping.ListFilter{TagID: &tagID})
	if err != nil {
		return fmt.Errorf("查詢 tag 既有映射失敗: %w", err)
	}
	if len(tagMappings) > 1 {
		return fmt.Errorf("tag %s 存在多條映射，無法同步來源規則", tagRecord.Key)
	}
	if len(tagMappings) == 1 && tagMappings[0].PointID != pointID {
		return fmt.Errorf("tag %s 已綁定其他 point", tagRecord.Key)
	}
	return nil
}

func marshalSkippedAddresses(addresses []string) (string, error) {
	if len(addresses) == 0 {
		return "[]", nil
	}

	normalized := make([]string, 0, len(addresses))
	for _, address := range addresses {
		address = strings.ToUpper(strings.TrimSpace(address))
		if address == "" {
			continue
		}
		normalized = append(normalized, address)
	}
	slices.Sort(normalized)
	payload, err := json.Marshal(normalized)
	if err != nil {
		return "", fmt.Errorf("序列化 skipped addresses 失敗: %w", err)
	}
	return string(payload), nil
}

func parseSkippedAddresses(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	var addresses []string
	if err := json.Unmarshal([]byte(raw), &addresses); err != nil {
		return nil, fmt.Errorf("解析 skipped addresses 失敗: %w", err)
	}
	return addresses, nil
}

func normalizeOrigin(origin string) string {
	switch strings.TrimSpace(origin) {
	case "", ruleOriginManual:
		return ruleOriginManual
	case "template":
		return "template"
	default:
		return ruleOriginManual
	}
}

func normalizeNamingPrefix(prefix string) string {
	compact := strings.ToUpper(strings.TrimSpace(prefix))
	if compact == "" {
		return "SRC"
	}

	var builder strings.Builder
	lastHyphen := false
	for _, ch := range compact {
		isAllowed := (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-'
		if isAllowed {
			builder.WriteRune(ch)
			lastHyphen = ch == '-'
			continue
		}
		if !lastHyphen {
			builder.WriteByte('-')
			lastHyphen = true
		}
	}

	normalized := strings.Trim(builder.String(), "-")
	if normalized == "" {
		return "SRC"
	}
	return normalized
}

func buildPointName(prefix, address string) string {
	return fmt.Sprintf("%s_%s", normalizeNamingPrefix(prefix), strings.ToUpper(strings.TrimSpace(address)))
}

func containsAddress(addresses []string, target string) bool {
	target = strings.ToUpper(strings.TrimSpace(target))
	return slices.ContainsFunc(addresses, func(address string) bool {
		return strings.ToUpper(strings.TrimSpace(address)) == target
	})
}

func normalizeAddressKey(address string) string {
	return strings.ToUpper(strings.TrimSpace(address))
}

func cloneSourceRuleLinks(links []*schema.SourceRuleLink) []*schema.SourceRuleLink {
	if len(links) == 0 {
		return nil
	}
	cloned := make([]*schema.SourceRuleLink, 0, len(links))
	for _, link := range links {
		if link == nil {
			continue
		}
		copyLink := *link
		if link.TagID != nil {
			copyLink.TagID = common.Ptr(*link.TagID)
		}
		if link.MappingID != nil {
			copyLink.MappingID = common.Ptr(*link.MappingID)
		}
		cloned = append(cloned, &copyLink)
	}
	return cloned
}

func (s *Service) cleanupRuleLinkResources(ctx context.Context, rule *schema.SourceRule, link *schema.SourceRuleLink) error {
	if s.mappingSvc != nil && link.MappingID != nil {
		if err := s.mappingSvc.Delete(ctx, *link.MappingID); err != nil && !errors.Is(err, mapping.ErrMappingNotFound) {
			return fmt.Errorf("刪除來源規則映射失敗: %w", err)
		}
	}
	if err := s.deleteRulePoint(ctx, link.PointID); err != nil {
		return err
	}
	if s.runtimeSync != nil {
		if !s.removeOwnedPointRuntime(ctx, link.PointID) {
			return s.dirtyUnknown(ctx)
		}
	}
	if s.tagSvc == nil || s.mappingSvc == nil || link.TagID == nil {
		return nil
	}
	return s.deleteRuleManagedTagIfOrphan(ctx, rule, link)
}

func (s *Service) deleteRuleManagedTagIfOrphan(ctx context.Context, rule *schema.SourceRule, link *schema.SourceRuleLink) error {
	tagRecord, err := s.tagSvc.GetByID(ctx, *link.TagID)
	if errors.Is(err, tag.ErrTagNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("取得來源規則標籤失敗: %w", err)
	}
	if !isRuleManagedTag(tagRecord, rule.ID, link.Address) {
		return nil
	}

	tagID := tagRecord.ID
	mappings, err := s.mappingSvc.List(ctx, mapping.ListFilter{TagID: &tagID})
	if err != nil {
		return fmt.Errorf("查詢來源規則標籤映射失敗: %w", err)
	}
	if len(mappings) > 0 {
		return nil
	}
	if err := s.tagSvc.Delete(ctx, tagID); err != nil && !errors.Is(err, tag.ErrTagNotFound) {
		return fmt.Errorf("刪除來源規則標籤失敗: %w", err)
	}
	return nil
}

func isRuleManagedTag(tagRecord *schema.Tag, ruleID, address string) bool {
	if tagRecord == nil || strings.TrimSpace(tagRecord.Labels) == "" {
		return false
	}
	var labels map[string]string
	if err := json.Unmarshal([]byte(tagRecord.Labels), &labels); err != nil {
		return false
	}
	return labels[ruleManagedTagLabelSource] == ruleManagedTagLabelValue &&
		labels[ruleManagedTagLabelRuleID] == ruleID &&
		labels[ruleManagedTagLabelAddress] == normalizeAddressKey(address)
}

func desiredRuleTargetDataType(rule *schema.SourceRule, pointRecord *schema.Point) schema.DataType {
	targetDataType := pointRecord.DataType
	if rule.TargetDataType != nil {
		targetDataType = *rule.TargetDataType
	}
	return targetDataType
}

func cloneDataTypePtr(value *schema.DataType) *schema.DataType {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneFloat64Ptr(value *float64) *float64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneIntPtr(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func encodeTransformPipeline(steps []schema.TransformStep) (string, error) {
	if steps == nil {
		steps = []schema.TransformStep{}
	}
	data, err := json.Marshal(steps)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func decodeTransformPipeline(raw string) ([]schema.TransformStep, error) {
	if strings.TrimSpace(raw) == "" {
		return []schema.TransformStep{}, nil
	}
	var steps []schema.TransformStep
	if err := json.Unmarshal([]byte(raw), &steps); err != nil {
		return nil, err
	}
	if steps == nil {
		return []schema.TransformStep{}, nil
	}
	return steps, nil
}

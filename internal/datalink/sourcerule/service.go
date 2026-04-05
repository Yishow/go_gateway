package sourcerule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

type Repository interface {
	Create(ctx context.Context, rule *schema.SourceRule) error
	Update(ctx context.Context, rule *schema.SourceRule) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*schema.SourceRule, error)
	List(ctx context.Context, filter ListFilter) ([]*schema.SourceRule, error)
	CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) error
	ListLinks(ctx context.Context, ruleID string) ([]*schema.SourceRuleLink, error)
	DeleteLinks(ctx context.Context, ruleID string) error
}

type RuntimeSyncer interface {
	UpsertPoint(point *schema.Point)
	RemovePoint(pointID string)
}

type ListFilter struct {
	DeviceID *string
	Enabled  *bool
}

type pointUpdatePlan struct {
	name  string
	point *schema.Point
}

type mappingRollbackState struct {
	enabled           bool
	transformPipeline []schema.TransformStep
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
}

type Service struct {
	repo        Repository
	deviceSvc   *device.Service
	pointSvc    *point.Service
	tagSvc      *tag.Service
	mappingSvc  *mapping.Service
	runtimeSync RuntimeSyncer
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

func (s *Service) Create(ctx context.Context, req CreateRuleRequest) (*schema.SourceRule, error) {
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	deviceRecord, err := s.deviceSvc.GetByID(ctx, req.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		id, err = common.NewUUID()
		if err != nil {
			return nil, fmt.Errorf("建立來源規則 ID 失敗: %w", err)
		}
	}

	enabled := req.Enabled && deviceRecord.Status == schema.DeviceStatusActive
	skippedJSON, err := marshalSkippedAddresses(req.SkippedAddresses)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	rule := &schema.SourceRule{
		ID:               id,
		DeviceID:         req.DeviceID,
		StartAddress:     strings.TrimSpace(strings.ToUpper(req.StartAddress)),
		Count:            req.Count,
		DataType:         req.DataType,
		NamingPrefix:     normalizeNamingPrefix(req.NamingPrefix),
		Enabled:          enabled,
		Locked:           req.Locked,
		Origin:           normalizeOrigin(req.Origin),
		TemplateName:     strings.TrimSpace(req.TemplateName),
		SkippedAddresses: skippedJSON,
		TargetDataType:   req.TargetDataType,
		ScaleMultiplier:  req.ScaleMultiplier,
		ScaleOffset:      req.ScaleOffset,
		DataFormat:       normalizeDataFormat(req.DataFormat),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := assignRuleRevisionID(rule); err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, rule); err != nil {
		return nil, fmt.Errorf("建立來源規則失敗: %w", err)
	}
	createdPointIDs := make([]string, 0, req.Count)
	links := make([]*schema.SourceRuleLink, 0, req.Count)
	createdPoints := make([]*schema.Point, 0, req.Count)
	defaultPollingGroupID, err := s.resolveDerivedPointPollingGroupID(ctx)
	if err != nil {
		return nil, err
	}

	for _, address := range buildPlannedPointAddresses(rule.StartAddress, rule.Count, rule.DataType, deviceRecord.Protocol) {
		if containsAddress(req.SkippedAddresses, address) {
			continue
		}

		pointRecord, createErr := s.pointSvc.Create(ctx, point.CreatePointRequest{
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
			s.rollbackCreatedPoints(ctx, createdPointIDs)
			_ = s.repo.Delete(ctx, rule.ID)
			return nil, fmt.Errorf("建立規則衍生點位失敗: %w", createErr)
		}

		if !enabled {
			disabled := false
			pointRecord, createErr = s.pointSvc.Update(ctx, pointRecord.ID, point.UpdatePointRequest{Enabled: &disabled})
			if createErr != nil {
				s.rollbackCreatedPoints(ctx, append(createdPointIDs, pointRecord.ID))
				_ = s.repo.Delete(ctx, rule.ID)
				return nil, fmt.Errorf("設定規則衍生點位狀態失敗: %w", createErr)
			}
		}

		linkID, linkErr := common.NewUUID()
		if linkErr != nil {
			s.rollbackCreatedPoints(ctx, append(createdPointIDs, pointRecord.ID))
			_ = s.repo.Delete(ctx, rule.ID)
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
		s.rollbackTagMappingSync(ctx, syncResult)
		s.rollbackCreatedPoints(ctx, createdPointIDs)
		_ = s.repo.Delete(ctx, rule.ID)
		return nil, fmt.Errorf("同步來源規則標籤映射失敗: %w", err)
	}

	if len(links) > 0 {
		if err := s.repo.CreateLinks(ctx, links); err != nil {
			s.rollbackTagMappingSync(ctx, syncResult)
			s.rollbackCreatedPoints(ctx, createdPointIDs)
			_ = s.repo.Delete(ctx, rule.ID)
			return nil, fmt.Errorf("建立來源規則連結失敗: %w", err)
		}
	}

	s.syncPoints(createdPoints)
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

func (s *Service) ListLinks(ctx context.Context, ruleID string) ([]*schema.SourceRuleLink, error) {
	links, err := s.repo.ListLinks(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則連結失敗: %w", err)
	}
	return links, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateRuleRequest) (*schema.SourceRule, error) {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}

	deviceRecord, err := s.deviceSvc.GetByID(ctx, rule.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}

	next := *rule
	if req.StartAddress != nil {
		next.StartAddress = strings.TrimSpace(strings.ToUpper(*req.StartAddress))
	}
	if req.Count != nil {
		next.Count = *req.Count
	}
	if req.DataType != nil {
		next.DataType = *req.DataType
	}
	if req.NamingPrefix != nil {
		next.NamingPrefix = normalizeNamingPrefix(*req.NamingPrefix)
	}
	if req.Locked != nil {
		next.Locked = *req.Locked
	}
	if req.TemplateName != nil {
		next.TemplateName = strings.TrimSpace(*req.TemplateName)
	}
	if req.SkippedAddresses != nil {
		skippedJSON, marshalErr := marshalSkippedAddresses(*req.SkippedAddresses)
		if marshalErr != nil {
			return nil, marshalErr
		}
		next.SkippedAddresses = skippedJSON
	}
	if req.TargetDataType != nil || req.TargetDataTypeSet {
		next.TargetDataType = cloneDataTypePtr(req.TargetDataType)
	}
	if req.ScaleMultiplier != nil || req.ScaleMultiplierSet {
		next.ScaleMultiplier = cloneFloat64Ptr(req.ScaleMultiplier)
	}
	if req.ScaleOffset != nil || req.ScaleOffsetSet {
		next.ScaleOffset = cloneFloat64Ptr(req.ScaleOffset)
	}
	if req.DataFormat != nil || req.DataFormatSet {
		if req.DataFormat != nil {
			next.DataFormat = normalizeDataFormat(*req.DataFormat)
		} else {
			next.DataFormat = ""
		}
	}
	if req.Enabled != nil {
		if *req.Enabled && deviceRecord.Status != schema.DeviceStatusActive {
			return nil, fmt.Errorf("設備尚未通過 probe readiness，不能啟用來源規則")
		}
		next.Enabled = *req.Enabled
	}

	if err := validateCreateRequest(CreateRuleRequest{
		DeviceID:     next.DeviceID,
		StartAddress: next.StartAddress,
		Count:        next.Count,
		DataType:     next.DataType,
		NamingPrefix: next.NamingPrefix,
	}); err != nil {
		return nil, err
	}
	if err := validateRuleDataFormat(next.DataFormat); err != nil {
		return nil, err
	}

	links, err := s.repo.ListLinks(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則連結失敗: %w", err)
	}

	desiredAddresses := buildPlannedPointAddresses(next.StartAddress, next.Count, next.DataType, deviceRecord.Protocol)
	skipped, err := parseSkippedAddresses(next.SkippedAddresses)
	if err != nil {
		return nil, err
	}

	filteredAddresses := make([]string, 0, len(desiredAddresses))
	for _, address := range desiredAddresses {
		if containsAddress(skipped, address) {
			continue
		}
		filteredAddresses = append(filteredAddresses, address)
	}

	linkByAddress := make(map[string]*schema.SourceRuleLink, len(links))
	keepPointIDs := make(map[string]struct{}, len(links))
	for _, link := range links {
		linkByAddress[normalizeAddressKey(link.Address)] = link
		keepPointIDs[link.PointID] = struct{}{}
	}

	// 連結仍指向已手動刪除的 point 時，先清關聯並改走下方「新建衍生點」路徑，避免儲存規則時 GetByID 失敗。
	for _, address := range filteredAddresses {
		key := normalizeAddressKey(address)
		existingLink, ok := linkByAddress[key]
		if !ok {
			continue
		}
		_, getErr := s.pointSvc.GetByID(ctx, existingLink.PointID)
		if getErr == nil {
			continue
		}
		if !errors.Is(getErr, point.ErrPointNotFound) {
			return nil, fmt.Errorf("取得衍生點位失敗: %w", getErr)
		}
		if err := s.cleanupRuleLinkResources(ctx, rule, existingLink); err != nil {
			return nil, fmt.Errorf("清理孤立來源規則連結資源失敗: %w", err)
		}
		delete(linkByAddress, key)
		delete(keepPointIDs, existingLink.PointID)
	}

	addedAddresses := make([]string, 0)
	for _, address := range filteredAddresses {
		if _, exists := linkByAddress[normalizeAddressKey(address)]; !exists {
			addedAddresses = append(addedAddresses, address)
		}
	}
	if err := s.validatePointAddresses(ctx, next.DeviceID, keepPointIDs, addedAddresses); err != nil {
		return nil, err
	}

	effectiveEnabled := next.Enabled && deviceRecord.Status == schema.DeviceStatusActive
	newLinks := make([]*schema.SourceRuleLink, 0, len(filteredAddresses))
	updatePlans := make([]pointUpdatePlan, 0, len(filteredAddresses))
	type pointCreatePlan struct {
		linkID  string
		address string
	}
	createPlans := make([]pointCreatePlan, 0, len(addedAddresses))
	removedLinks := make([]*schema.SourceRuleLink, 0)

	for _, address := range filteredAddresses {
		existingLink, exists := linkByAddress[normalizeAddressKey(address)]
		if exists {
			pointRecord, getErr := s.pointSvc.GetByID(ctx, existingLink.PointID)
			if getErr != nil {
				return nil, fmt.Errorf("取得衍生點位失敗: %w", getErr)
			}

			name := buildPointName(next.NamingPrefix, address)
			updatePlans = append(updatePlans, pointUpdatePlan{
				name:  name,
				point: pointRecord,
			})
			existingLink.Address = address
			existingLink.UpdatedAt = time.Now()
			newLinks = append(newLinks, existingLink)
			continue
		}

		linkID, linkErr := common.NewUUID()
		if linkErr != nil {
			return nil, fmt.Errorf("建立來源規則連結 ID 失敗: %w", linkErr)
		}
		createPlans = append(createPlans, pointCreatePlan{
			linkID:  linkID,
			address: address,
		})
	}

	desiredSet := make(map[string]struct{}, len(filteredAddresses))
	for _, address := range filteredAddresses {
		desiredSet[normalizeAddressKey(address)] = struct{}{}
	}
	for _, link := range links {
		if _, keep := desiredSet[normalizeAddressKey(link.Address)]; keep {
			continue
		}
		removedLinks = append(removedLinks, link)
	}

	next.UpdatedAt = time.Now()
	if err := assignRuleRevisionID(&next); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, &next); err != nil {
		return nil, fmt.Errorf("更新來源規則失敗: %w", err)
	}
	appliedUpdatePlans := make([]pointUpdatePlan, 0, len(updatePlans))
	defaultPollingGroupID, err := s.resolveDerivedPointPollingGroupID(ctx)
	if err != nil {
		return nil, err
	}
	for _, plan := range updatePlans {
		dataType := next.DataType
		df := next.DataFormat
		updateReq := point.UpdatePointRequest{
			Name:       &plan.name,
			DataType:   &dataType,
			DataFormat: &df,
			Enabled:    &effectiveEnabled,
		}
		if plan.point.PollingGroupID == nil && defaultPollingGroupID != nil {
			updateReq.PollingGroupID = defaultPollingGroupID
		}
		pointRecord, updateErr := s.pointSvc.Update(ctx, plan.point.ID, updateReq)
		if updateErr != nil {
			s.rollbackUpdatedPoints(ctx, appliedUpdatePlans)
			s.rollbackRuleState(ctx, rule, links)
			return nil, fmt.Errorf("更新衍生點位失敗: %w", updateErr)
		}
		appliedUpdatePlans = append(appliedUpdatePlans, plan)
		s.syncPoints([]*schema.Point{pointRecord})
	}

	createdPointIDs := make([]string, 0, len(createPlans))
	createdLinkMap := make(map[string]*schema.SourceRuleLink, len(createPlans))
	for _, plan := range createPlans {
		pointRecord, createErr := s.pointSvc.Create(ctx, point.CreatePointRequest{
			DeviceID:       next.DeviceID,
			Name:           buildPointName(next.NamingPrefix, plan.address),
			Address:        plan.address,
			DataType:       next.DataType,
			DataFormat:     next.DataFormat,
			Mode:           schema.PointModeReadOnly,
			PollingGroupID: defaultPollingGroupID,
		})
		if createErr != nil {
			s.rollbackUpdatedPoints(ctx, appliedUpdatePlans)
			s.rollbackRuleState(ctx, rule, links)
			s.rollbackCreatedPoints(ctx, createdPointIDs)
			return nil, fmt.Errorf("建立新增衍生點位失敗: %w", createErr)
		}
		if !effectiveEnabled {
			pointRecord, createErr = s.pointSvc.Update(ctx, pointRecord.ID, point.UpdatePointRequest{Enabled: &effectiveEnabled})
			if createErr != nil {
				s.rollbackUpdatedPoints(ctx, appliedUpdatePlans)
				s.rollbackRuleState(ctx, rule, links)
				s.rollbackCreatedPoints(ctx, append(createdPointIDs, pointRecord.ID))
				return nil, fmt.Errorf("更新新增衍生點位狀態失敗: %w", createErr)
			}
		}
		createdPointIDs = append(createdPointIDs, pointRecord.ID)
		createdLinkMap[normalizeAddressKey(plan.address)] = &schema.SourceRuleLink{
			ID:        plan.linkID,
			RuleID:    next.ID,
			Address:   plan.address,
			PointID:   pointRecord.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		s.syncPoints([]*schema.Point{pointRecord})
	}

	for _, address := range filteredAddresses {
		if link, exists := createdLinkMap[normalizeAddressKey(address)]; exists {
			newLinks = append(newLinks, link)
		}
	}

	syncResult, err := s.syncRuleTagMappings(ctx, rule, &next, newLinks, effectiveEnabled)
	if err != nil {
		s.rollbackTagMappingSync(ctx, syncResult)
		s.rollbackUpdatedPoints(ctx, appliedUpdatePlans)
		s.rollbackRuleState(ctx, rule, links)
		s.rollbackCreatedPoints(ctx, createdPointIDs)
		return nil, fmt.Errorf("同步來源規則標籤映射失敗: %w", err)
	}

	if err := s.repo.DeleteLinks(ctx, next.ID); err != nil {
		s.rollbackTagMappingSync(ctx, syncResult)
		s.rollbackUpdatedPoints(ctx, appliedUpdatePlans)
		s.rollbackRuleState(ctx, rule, links)
		s.rollbackCreatedPoints(ctx, createdPointIDs)
		return nil, fmt.Errorf("清除舊來源規則連結失敗: %w", err)
	}
	if len(newLinks) > 0 {
		if err := s.repo.CreateLinks(ctx, newLinks); err != nil {
			s.rollbackTagMappingSync(ctx, syncResult)
			s.rollbackUpdatedPoints(ctx, appliedUpdatePlans)
			s.rollbackRuleState(ctx, rule, links)
			s.rollbackCreatedPoints(ctx, createdPointIDs)
			return nil, fmt.Errorf("重建來源規則連結失敗: %w", err)
		}
	}

	for _, link := range removedLinks {
		if err := s.cleanupRuleLinkResources(ctx, &next, link); err != nil {
			return nil, fmt.Errorf("刪除已移除衍生點位失敗: %w", err)
		}
	}

	return &next, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得來源規則失敗: %w", err)
	}
	links, err := s.repo.ListLinks(ctx, id)
	if err != nil {
		return fmt.Errorf("取得來源規則連結失敗: %w", err)
	}
	for _, link := range links {
		if err := s.cleanupRuleLinkResources(ctx, rule, link); err != nil {
			return fmt.Errorf("刪除衍生點位失敗: %w", err)
		}
	}
	if err := s.repo.DeleteLinks(ctx, id); err != nil {
		return fmt.Errorf("刪除來源規則連結失敗: %w", err)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除來源規則失敗: %w", err)
	}
	return nil
}

func (s *Service) Disable(ctx context.Context, id string) error {
	return s.setEnabled(ctx, id, false)
}

func (s *Service) Enable(ctx context.Context, id string) error {
	return s.setEnabled(ctx, id, true)
}

func (s *Service) SyncDerivedPointState(ctx context.Context) error {
	rules, err := s.repo.List(ctx, ListFilter{})
	if err != nil {
		return fmt.Errorf("列出來源規則失敗: %w", err)
	}

	for _, rule := range rules {
		deviceRecord, getErr := s.deviceSvc.GetByID(ctx, rule.DeviceID)
		if getErr != nil {
			return fmt.Errorf("取得來源規則設備失敗: %w", getErr)
		}
		enabled := rule.Enabled && deviceRecord.Status == schema.DeviceStatusActive
		if _, err := s.syncRuleLinksEnabled(ctx, rule.ID, enabled); err != nil {
			return err
		}
		currentLinks, linkErr := s.repo.ListLinks(ctx, rule.ID)
		if linkErr != nil {
			return fmt.Errorf("取得來源規則連結失敗: %w", linkErr)
		}
		nextLinks := cloneSourceRuleLinks(currentLinks)
		syncResult, syncErr := s.syncRuleTagMappings(ctx, rule, rule, nextLinks, enabled)
		if syncErr != nil {
			s.rollbackTagMappingSync(ctx, syncResult)
			return fmt.Errorf("同步來源規則標籤映射失敗: %w", syncErr)
		}
		if err := s.replaceRuleLinks(ctx, rule.ID, currentLinks, nextLinks); err != nil {
			s.rollbackTagMappingSync(ctx, syncResult)
			return err
		}
	}

	return nil
}

func (s *Service) validatePointAddresses(ctx context.Context, deviceID string, keepPointIDs map[string]struct{}, addresses []string) error {
	filter := point.ListFilter{DeviceID: &deviceID}
	points, err := s.pointSvc.List(ctx, filter)
	if err != nil {
		return fmt.Errorf("取得既有點位失敗: %w", err)
	}

	existing := make(map[string]string, len(points))
	for _, item := range points {
		if _, ok := keepPointIDs[item.ID]; ok {
			continue
		}
		existing[strings.ToUpper(strings.TrimSpace(item.Address))] = item.ID
	}

	for _, address := range addresses {
		if existing[strings.ToUpper(strings.TrimSpace(address))] != "" {
			return fmt.Errorf("位址 %s 已存在既有點位", address)
		}
	}
	return nil
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
			s.rollbackUpdatedPoints(ctx, appliedPlans)
			return nil, fmt.Errorf("取得衍生點位失敗: %w", getErr)
		}
		original := pointRecord
		updateReq := point.UpdatePointRequest{
			Enabled: &enabled,
		}
		if enabled && pointRecord.PollingGroupID == nil && defaultPollingGroupID != nil {
			updateReq.PollingGroupID = defaultPollingGroupID
		}
		pointRecord, getErr = s.pointSvc.Update(ctx, link.PointID, point.UpdatePointRequest{
			Enabled:             updateReq.Enabled,
			PollingGroupID:      updateReq.PollingGroupID,
			ReplacePollingGroup: updateReq.ReplacePollingGroup,
		})
		if getErr != nil {
			s.rollbackUpdatedPoints(ctx, appliedPlans)
			return nil, fmt.Errorf("更新衍生點位狀態失敗: %w", getErr)
		}
		appliedPlans = append(appliedPlans, pointUpdatePlan{
			name:  original.Name,
			point: original,
		})
		s.syncPoints([]*schema.Point{pointRecord})
	}
	return appliedPlans, nil
}

func (s *Service) setEnabled(ctx context.Context, id string, enabled bool) error {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得來源規則失敗: %w", err)
	}

	deviceRecord, err := s.deviceSvc.GetByID(ctx, rule.DeviceID)
	if err != nil {
		return fmt.Errorf("取得設備失敗: %w", err)
	}
	previousRule := *rule
	previousLinks, err := s.repo.ListLinks(ctx, id)
	if err != nil {
		return fmt.Errorf("取得來源規則連結失敗: %w", err)
	}
	if enabled && deviceRecord.Status != schema.DeviceStatusActive {
		return fmt.Errorf("設備尚未通過 probe readiness，不能啟用來源規則")
	}
	rule.Enabled = enabled
	rule.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, rule); err != nil {
		return fmt.Errorf("更新來源規則失敗: %w", err)
	}
	appliedPointPlans, err := s.syncRuleLinksEnabled(ctx, id, enabled)
	if err != nil {
		_ = s.repo.Update(ctx, &previousRule)
		return err
	}

	nextLinks := cloneSourceRuleLinks(previousLinks)
	syncResult, syncErr := s.syncRuleTagMappings(ctx, rule, rule, nextLinks, enabled)
	if syncErr != nil {
		s.rollbackTagMappingSync(ctx, syncResult)
		s.rollbackUpdatedPoints(ctx, appliedPointPlans)
		_ = s.repo.Update(ctx, &previousRule)
		return fmt.Errorf("同步來源規則標籤映射失敗: %w", syncErr)
	}

	if err := s.replaceRuleLinks(ctx, id, previousLinks, nextLinks); err != nil {
		s.rollbackTagMappingSync(ctx, syncResult)
		s.rollbackUpdatedPoints(ctx, appliedPointPlans)
		_ = s.repo.Update(ctx, &previousRule)
		return err
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

func (s *Service) replaceRuleLinks(ctx context.Context, ruleID string, previousLinks, nextLinks []*schema.SourceRuleLink) error {
	if err := s.repo.DeleteLinks(ctx, ruleID); err != nil {
		return fmt.Errorf("清除來源規則連結失敗: %w", err)
	}
	if len(nextLinks) == 0 {
		return nil
	}
	if err := s.repo.CreateLinks(ctx, nextLinks); err != nil {
		_ = s.repo.DeleteLinks(ctx, ruleID)
		if len(previousLinks) > 0 {
			_ = s.repo.CreateLinks(ctx, previousLinks)
		}
		return fmt.Errorf("重建來源規則連結失敗: %w", err)
	}
	return nil
}

func (s *Service) syncRuleTagMappings(ctx context.Context, oldRule *schema.SourceRule, rule *schema.SourceRule, links []*schema.SourceRuleLink, enabled bool) (tagMappingSyncResult, error) {
	result := tagMappingSyncResult{
		updatedMappings: make(map[string]mappingRollbackState),
		updatedTags:     make(map[string]tagRollbackState),
	}
	if s.tagSvc == nil || s.mappingSvc == nil || len(links) == 0 {
		return result, nil
	}

	for _, link := range links {
		pointRecord, err := s.pointSvc.GetByID(ctx, link.PointID)
		if err != nil {
			return result, fmt.Errorf("取得衍生點位失敗: %w", err)
		}
		tagRecord, mappingRecord, err := s.ensureRuleTagMapping(ctx, oldRule, rule, pointRecord, link, enabled, &result)
		if err != nil {
			return result, err
		}
		link.TagID = stringPtr(tagRecord.ID)
		link.MappingID = stringPtr(mappingRecord.ID)
		link.UpdatedAt = time.Now()
	}

	return result, nil
}

func (s *Service) ensureRuleTagMapping(
	ctx context.Context,
	oldRule *schema.SourceRule,
	rule *schema.SourceRule,
	pointRecord *schema.Point,
	link *schema.SourceRuleLink,
	enabled bool,
	result *tagMappingSyncResult,
) (*schema.Tag, *schema.Mapping, error) {
	transformPipeline := s.buildRuleTransformPipeline(rule, pointRecord)
	var oldTransformPipeline []schema.TransformStep
	if oldRule != nil {
		oldTransformPipeline = s.buildRuleTransformPipeline(oldRule, &schema.Point{DataType: oldRule.DataType})
	} else {
		oldTransformPipeline = []schema.TransformStep{}
	}
	targetDataType := desiredRuleTargetDataType(rule, pointRecord)

	pointID := pointRecord.ID
	pointMappings, err := s.mappingSvc.List(ctx, mapping.ListFilter{PointID: &pointID})
	if err != nil {
		return nil, nil, fmt.Errorf("查詢 point 既有映射失敗: %w", err)
	}
	if len(pointMappings) > 1 {
		return nil, nil, fmt.Errorf("point %s 存在多條映射，無法同步來源規則", pointRecord.ID)
	}
	if len(pointMappings) == 1 {
		tagRecord, getErr := s.tagSvc.GetByID(ctx, pointMappings[0].TagID)
		if getErr != nil {
			return nil, nil, fmt.Errorf("取得既有映射標籤失敗: %w", getErr)
		}
		tagRecord, syncErr := s.syncRuleManagedTagDataType(ctx, rule, link, tagRecord, targetDataType, result)
		if syncErr != nil {
			return nil, nil, syncErr
		}
		mappingRecord, syncErr := s.syncRuleMapping(ctx, pointMappings[0], enabled, oldTransformPipeline, transformPipeline, result)
		if syncErr != nil {
			return nil, nil, syncErr
		}
		return tagRecord, mappingRecord, nil
	}

	tagRecord, createdTag, err := s.resolveRuleTag(ctx, rule, pointRecord, link, result)
	if err != nil {
		return nil, nil, err
	}
	if createdTag {
		result.createdTagIDs = append(result.createdTagIDs, tagRecord.ID)
	}
	if err := s.validateTagAvailability(ctx, tagRecord, pointRecord.ID); err != nil {
		return nil, nil, err
	}

	tagID := tagRecord.ID
	tagMappings, err := s.mappingSvc.List(ctx, mapping.ListFilter{TagID: &tagID})
	if err != nil {
		return nil, nil, fmt.Errorf("查詢 tag 既有映射失敗: %w", err)
	}
	if len(tagMappings) > 1 {
		return nil, nil, fmt.Errorf("tag %s 存在多條映射，無法同步來源規則", tagRecord.Key)
	}
	if len(tagMappings) == 1 {
		if tagMappings[0].PointID != pointRecord.ID {
			return nil, nil, fmt.Errorf("tag %s 已綁定其他 point", tagRecord.Key)
		}
		mappingRecord, syncErr := s.syncRuleMapping(ctx, tagMappings[0], enabled, oldTransformPipeline, transformPipeline, result)
		if syncErr != nil {
			return nil, nil, syncErr
		}
		return tagRecord, mappingRecord, nil
	}

	mappingEnabled := enabled
	mappingRecord, err := s.mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID:           pointRecord.ID,
		TagID:             tagRecord.ID,
		Enabled:           &mappingEnabled,
		TransformPipeline: transformPipeline,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("建立來源規則映射失敗: %w", err)
	}
	result.createdMappingIDs = append(result.createdMappingIDs, mappingRecord.ID)
	return tagRecord, mappingRecord, nil
}

func (s *Service) resolveRuleTag(
	ctx context.Context,
	rule *schema.SourceRule,
	pointRecord *schema.Point,
	link *schema.SourceRuleLink,
	result *tagMappingSyncResult,
) (*schema.Tag, bool, error) {
	targetDataType := desiredRuleTargetDataType(rule, pointRecord)

	if link.TagID != nil && strings.TrimSpace(*link.TagID) != "" {
		tagRecord, err := s.tagSvc.GetByID(ctx, *link.TagID)
		if err == nil {
			tagRecord, err = s.syncRuleManagedTagDataType(ctx, rule, link, tagRecord, targetDataType, result)
			if err != nil {
				return nil, false, err
			}
			return tagRecord, false, nil
		}
		if !errors.Is(err, tag.ErrTagNotFound) {
			return nil, false, fmt.Errorf("取得來源規則既有標籤失敗: %w", err)
		}
	}

	expectedKey := buildPointName(rule.NamingPrefix, link.Address)
	tagRecord, err := s.tagSvc.GetByKey(ctx, expectedKey)
	if err == nil {
		if tagRecord.Status == schema.TagStatusRetired {
			return nil, false, fmt.Errorf("自動產生的標籤鍵已被 retired tag 佔用: %s", expectedKey)
		}
		tagRecord, err = s.syncRuleManagedTagDataType(ctx, rule, link, tagRecord, targetDataType, result)
		if err != nil {
			return nil, false, err
		}
		return tagRecord, false, nil
	}
	if !errors.Is(err, tag.ErrTagNotFound) {
		return nil, false, fmt.Errorf("查詢來源規則標籤失敗: %w", err)
	}

	// 建立新 Tag 時使用目標型別
	created, err := s.tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         expectedKey,
		DisplayName: pointRecord.Name,
		DataType:    targetDataType,
		Labels: map[string]string{
			ruleManagedTagLabelSource:  ruleManagedTagLabelValue,
			ruleManagedTagLabelRuleID:  rule.ID,
			ruleManagedTagLabelAddress: normalizeAddressKey(link.Address),
		},
	})
	if err != nil {
		return nil, false, fmt.Errorf("建立來源規則標籤失敗: %w", err)
	}
	return created, true, nil
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

func (s *Service) syncRuleMapping(
	ctx context.Context,
	mappingRecord *schema.Mapping,
	enabled bool,
	oldTransformPipeline []schema.TransformStep,
	transformPipeline []schema.TransformStep,
	result *tagMappingSyncResult,
) (*schema.Mapping, error) {
	nextPipelineJSON, err := encodeTransformPipeline(transformPipeline)
	if err != nil {
		return nil, fmt.Errorf("序列化來源規則映射轉換管線失敗: %w", err)
	}
	if mappingRecord.Enabled == enabled && mappingRecord.TransformPipeline == nextPipelineJSON {
		return mappingRecord, nil
	}
	if result.updatedMappings == nil {
		result.updatedMappings = make(map[string]mappingRollbackState)
	}
	if _, exists := result.updatedMappings[mappingRecord.ID]; !exists {
		previousPipeline, err := decodeTransformPipeline(mappingRecord.TransformPipeline)
		if err != nil {
			return nil, fmt.Errorf("解析既有來源規則映射轉換管線失敗: %w", err)
		}
		result.updatedMappings[mappingRecord.ID] = mappingRollbackState{
			enabled:           mappingRecord.Enabled,
			transformPipeline: previousPipeline,
		}
	}

	updateReq := mapping.UpdateMappingRequest{}
	if mappingRecord.Enabled != enabled {
		updateReq.Enabled = &enabled
	}

	if mappingRecord.TransformPipeline != nextPipelineJSON {
		oldPipelineJSON, err := encodeTransformPipeline(oldTransformPipeline)
		if err != nil {
			return nil, fmt.Errorf("序列化既有來源規則映射轉換管線失敗: %w", err)
		}
		if mappingRecord.TransformPipeline != oldPipelineJSON {
			return nil, fmt.Errorf("映射存在手動編輯的轉換管線，無法自動覆蓋")
		}
		updateReq.TransformPipeline = transformPipeline
	}

	updated, err := s.mappingSvc.Update(ctx, mappingRecord.ID, updateReq)
	if err != nil {
		return nil, fmt.Errorf("同步來源規則映射設定失敗: %w", err)
	}
	return updated, nil
}

func validateCreateRequest(req CreateRuleRequest) error {
	if strings.TrimSpace(req.DeviceID) == "" {
		return fmt.Errorf("device_id is required")
	}
	if strings.TrimSpace(req.StartAddress) == "" {
		return fmt.Errorf("start_address is required")
	}
	if req.Count <= 0 {
		return fmt.Errorf("count must be greater than zero")
	}
	if err := validateRuleDataType(req.DataType); err != nil {
		return err
	}
	if err := validateRuleDataFormat(req.DataFormat); err != nil {
		return err
	}
	return nil
}

func normalizeDataFormat(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func validateRuleDataFormat(dataFormat string) error {
	normalized := normalizeDataFormat(dataFormat)
	if normalized == "" {
		return nil
	}
	switch normalized {
	case "ABCD", "BADC", "CDAB", "DCBA":
		return nil
	default:
		return fmt.Errorf("unsupported data_format: %s", dataFormat)
	}
}

func validateRuleDataType(dataType schema.DataType) error {
	if dataType == "" {
		return fmt.Errorf("data_type is required")
	}
	switch dataType {
	case schema.DataTypeBool,
		schema.DataTypeInt16,
		schema.DataTypeInt32,
		schema.DataTypeInt64,
		schema.DataTypeUint16,
		schema.DataTypeUint32,
		schema.DataTypeUint64,
		schema.DataTypeFloat32,
		schema.DataTypeFloat64,
		schema.DataTypeString:
		return nil
	default:
		return fmt.Errorf("unsupported data_type: %s", dataType)
	}
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
	sort.Strings(normalized)
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
	case "", "manual":
		return "manual"
	case "template":
		return "template"
	default:
		return "manual"
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

func buildPointName(prefix string, address string) string {
	return fmt.Sprintf("%s_%s", normalizeNamingPrefix(prefix), strings.ToUpper(strings.TrimSpace(address)))
}

func containsAddress(addresses []string, target string) bool {
	target = strings.ToUpper(strings.TrimSpace(target))
	for _, address := range addresses {
		if strings.ToUpper(strings.TrimSpace(address)) == target {
			return true
		}
	}
	return false
}

func normalizeAddressKey(address string) string {
	return strings.ToUpper(strings.TrimSpace(address))
}

func buildPlannedPointAddresses(startAddress string, count int, dataType schema.DataType, protocol schema.ProtocolType) []string {
	span := getDataTypeCellSpan(dataType)
	addresses := make([]string, 0, count)
	for index := 0; index < count; index++ {
		addresses = append(addresses, offsetAddress(startAddress, index*span, protocol))
	}
	return addresses
}

func getDataTypeCellSpan(dataType schema.DataType) int {
	switch dataType {
	case schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32:
		return 2
	case schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64:
		return 4
	default:
		return 1
	}
}

func offsetAddress(address string, delta int, protocol schema.ProtocolType) string {
	address = strings.ToUpper(strings.TrimSpace(address))
	if address == "" || delta == 0 {
		return address
	}

	switch protocol {
	case schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP:
		if len(address) < 2 {
			return address
		}
		prefix := address[:1]
		number := address[1:]
		value, ok := parseDecimal(number)
		if !ok {
			return address
		}
		next := value + delta
		if next < 1 {
			next = 1
		}
		return fmt.Sprintf("%s%04d", prefix, next)
	case schema.ProtocolFatekFBs, schema.ProtocolMC3E:
		area, number := splitAlphaNumeric(address)
		if area == "" || number == "" {
			return address
		}
		value, ok := parseDecimal(number)
		if !ok {
			return address
		}
		next := value + delta
		if next < 0 {
			next = 0
		}
		return fmt.Sprintf("%s%d", area, next)
	default:
		return address
	}
}

func parseDecimal(raw string) (int, bool) {
	value := 0
	for _, ch := range raw {
		if ch < '0' || ch > '9' {
			return 0, false
		}
		value = value*10 + int(ch-'0')
	}
	return value, true
}

func splitAlphaNumeric(input string) (string, string) {
	index := 0
	for ; index < len(input); index++ {
		ch := input[index]
		if ch >= '0' && ch <= '9' {
			break
		}
	}
	if index == 0 || index >= len(input) {
		return "", ""
	}
	return input[:index], input[index:]
}

func (s *Service) syncPoints(points []*schema.Point) {
	if s.runtimeSync == nil {
		return
	}
	for _, pointRecord := range points {
		s.runtimeSync.UpsertPoint(pointRecord)
	}
}

func (s *Service) rollbackCreatedPoints(ctx context.Context, pointIDs []string) {
	for _, pointID := range pointIDs {
		_ = s.pointSvc.Delete(ctx, pointID)
		if s.runtimeSync != nil {
			s.runtimeSync.RemovePoint(pointID)
		}
	}
}

func (s *Service) rollbackUpdatedPoints(ctx context.Context, plans []pointUpdatePlan) {
	for _, plan := range plans {
		name := plan.point.Name
		dataType := plan.point.DataType
		enabled := plan.point.Enabled
		_, _ = s.pointSvc.Update(ctx, plan.point.ID, point.UpdatePointRequest{
			Name:                &name,
			DataType:            &dataType,
			PollingGroupID:      plan.point.PollingGroupID,
			ReplacePollingGroup: true,
			Enabled:             &enabled,
		})
		s.syncPoints([]*schema.Point{plan.point})
	}
}

func (s *Service) rollbackRuleState(ctx context.Context, rule *schema.SourceRule, links []*schema.SourceRuleLink) {
	if rule == nil {
		return
	}
	_ = s.repo.Update(ctx, rule)
	_ = s.repo.DeleteLinks(ctx, rule.ID)
	if len(links) > 0 {
		_ = s.repo.CreateLinks(ctx, links)
	}
}

func (s *Service) rollbackTagMappingSync(ctx context.Context, result tagMappingSyncResult) {
	if s.mappingSvc != nil {
		for mappingID, state := range result.updatedMappings {
			enabled := state.enabled
			_, _ = s.mappingSvc.Update(ctx, mappingID, mapping.UpdateMappingRequest{
				Enabled:           &enabled,
				TransformPipeline: state.transformPipeline,
			})
		}
		for _, mappingID := range result.createdMappingIDs {
			_ = s.mappingSvc.Delete(ctx, mappingID)
		}
	}
	if s.tagSvc != nil {
		for tagID, state := range result.updatedTags {
			dataType := state.dataType
			_, _ = s.tagSvc.Update(ctx, tagID, tag.UpdateTagRequest{DataType: &dataType})
		}
		for _, tagID := range result.createdTagIDs {
			_ = s.tagSvc.Delete(ctx, tagID)
		}
	}
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
			copyLink.TagID = stringPtr(*link.TagID)
		}
		if link.MappingID != nil {
			copyLink.MappingID = stringPtr(*link.MappingID)
		}
		cloned = append(cloned, &copyLink)
	}
	return cloned
}

func stringPtr(value string) *string {
	copyValue := value
	return &copyValue
}

func (s *Service) cleanupRuleLinkResources(ctx context.Context, rule *schema.SourceRule, link *schema.SourceRuleLink) error {
	if s.mappingSvc != nil && link.MappingID != nil {
		if err := s.mappingSvc.Delete(ctx, *link.MappingID); err != nil && !errors.Is(err, mapping.ErrMappingNotFound) {
			return fmt.Errorf("刪除來源規則映射失敗: %w", err)
		}
	}
	if err := s.pointSvc.Delete(ctx, link.PointID); err != nil {
		return err
	}
	if s.runtimeSync != nil {
		s.runtimeSync.RemovePoint(link.PointID)
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

// buildRuleTransformPipeline 根據來源規則建構 Mapping 的轉換管線。
// 順序：先 cast（若目標型別與 Point 型別不同），再 scale（若有設定）。
func (s *Service) buildRuleTransformPipeline(rule *schema.SourceRule, pointRecord *schema.Point) []schema.TransformStep {
	steps := make([]schema.TransformStep, 0)
	order := 0

	targetDataType := desiredRuleTargetDataType(rule, pointRecord)

	// 若目標型別與 Point 型別不同，加入 cast 步驟
	if targetDataType != pointRecord.DataType {
		if !targetDataType.IsValid() {
			// 如果無效，回退為與 Point 相同型別
			targetDataType = pointRecord.DataType
		} else {
			steps = append(steps, schema.TransformStep{
				Type:  schema.TransformCast,
				Order: order,
				Params: map[string]interface{}{
					"target_type": string(targetDataType),
				},
			})
			order++
		}
	}

	// 若有 scale 設定，加入 scale 步驟
	if rule.ScaleMultiplier != nil || rule.ScaleOffset != nil {
		scale := 1.0
		if rule.ScaleMultiplier != nil {
			scale = *rule.ScaleMultiplier
		}
		offset := 0.0
		if rule.ScaleOffset != nil {
			offset = *rule.ScaleOffset
		}

		steps = append(steps, schema.TransformStep{
			Type:  schema.TransformScale,
			Order: order,
			Params: map[string]interface{}{
				"scale":  scale,
				"offset": offset,
			},
		})
	}

	return steps
}

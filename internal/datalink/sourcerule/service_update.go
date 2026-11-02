package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) Update(ctx context.Context, id string, req UpdateRuleRequest) (result *schema.SourceRule, err error) {
	if err := validateRuleScale(req.ScaleMultiplier, req.ScaleOffset); err != nil {
		return nil, err
	}
	ctx, release := s.AcquireRuleMutation(ctx, id)
	defer release()
	defer func(finalizeCtx context.Context) { err = s.finalizePointWrites(finalizeCtx, err) }(ctx)
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	if _, valid := ruleStateKey(rule); !valid {
		return nil, ErrRuleRevisionConflict
	}
	ctx, finishAcquisition := s.beginSourceAcquisitionMutation(ctx, id, rule.DeviceID)
	defer finishAcquisition()

	deviceRecord, err := s.deviceSvc.GetByID(ctx, rule.DeviceID)
	if err != nil {
		return nil, wrapSourceRuleDeviceError(rule.DeviceID, err)
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
	if req.ShareEnabled != nil {
		next.ShareEnabled = *req.ShareEnabled
	}
	if req.ShareStartRegister != nil || req.ShareStartRegisterSet {
		next.ShareStartRegister = cloneIntPtr(req.ShareStartRegister)
	}
	if req.ShareStride != nil || req.ShareStrideSet {
		next.ShareStride = cloneIntPtr(req.ShareStride)
	}
	if req.Enabled != nil {
		next.Enabled = *req.Enabled
	}

	if err := validateCreateRequest(CreateRuleRequest{
		DeviceID:           next.DeviceID,
		StartAddress:       next.StartAddress,
		Count:              next.Count,
		DataType:           next.DataType,
		NamingPrefix:       next.NamingPrefix,
		Enabled:            next.Enabled,
		ShareEnabled:       next.ShareEnabled,
		ShareStartRegister: next.ShareStartRegister,
		ShareStride:        next.ShareStride,
	}); err != nil {
		return nil, contextualRuleValidationError(id, next.StartAddress, deviceRecord.Protocol, err)
	}
	if err := validateRuleDataFormat(next.DataFormat); err != nil {
		return nil, contextualRuleValidationError(id, next.StartAddress, deviceRecord.Protocol, err)
	}

	links, err := s.listAndValidatePersistedMCLinks(ctx, rule, deviceRecord.Protocol)
	if err != nil {
		return nil, err
	}
	desiredAddresses, err := buildPlannedPointAddresses(next.StartAddress, next.Count, next.DataType, deviceRecord.Protocol)
	if err != nil {
		return nil, contextualRuleValidationError(id, next.StartAddress, deviceRecord.Protocol, err)
	}
	if err := s.validateRulePointFormats(ctx, rule, &next, deviceRecord.Protocol, links); err != nil {
		return nil, err
	}
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
	if err := s.preflightPointAddresses(ctx, next.DeviceID, keepPointIDs, filteredAddresses, nil); err != nil {
		return nil, err
	}
	// Readiness checks persist device state, so run them after read-only preflight.
	if req.Enabled != nil {
		if err := s.ensureRuleActivationAllowed(ctx, rule.DeviceID, *req.Enabled); err != nil {
			return nil, contextualRuleValidationError(id, next.StartAddress, deviceRecord.Protocol, err)
		}
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
	survivingLinks := make([]*schema.SourceRuleLink, 0, len(links))
	for _, link := range links {
		if _, keep := keepPointIDs[link.PointID]; keep {
			survivingLinks = append(survivingLinks, link)
		}
	}
	links = survivingLinks
	rollbackLinks := cloneSourceRuleLinks(links)

	// Points follow the device's activation readiness, like rule creation and
	// enabling do, not its status: the runtime projection already collects only
	// from active devices. Tying this to status "active" would switch derived
	// points and Tag mappings off for a device that is ready but not yet
	// activated, and nothing re-enables them when it is.
	effectiveEnabled, err := s.runtimeRuleEnabled(ctx, &next)
	if err != nil {
		return nil, err
	}
	newLinks := make([]*schema.SourceRuleLink, 0, len(filteredAddresses))
	updatePlans := make([]pointUpdatePlan, 0, len(filteredAddresses))
	type pointCreatePlan struct {
		linkID  string
		address string
	}
	createPlans := make([]pointCreatePlan, 0)
	removedLinks := make([]*schema.SourceRuleLink, 0)

	for _, address := range filteredAddresses {
		existingLink, exists := linkByAddress[normalizeAddressKey(address)]
		if exists {
			pointRecord, getErr := s.pointSvc.GetByID(ctx, existingLink.PointID)
			if getErr != nil {
				return nil, fmt.Errorf("取得衍生點位失敗: %w", getErr)
			}

			s.capturePointBaseline(ctx, pointRecord)
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
	if err := s.updateOwnedRule(ctx, rule, &next); err != nil {
		return nil, fmt.Errorf("更新來源規則失敗: %w", err)
	}
	appliedUpdatePlans := make([]pointUpdatePlan, 0, len(updatePlans))
	defaultPollingGroupID, err := s.resolveDerivedPointPollingGroupID(ctx)
	if err != nil {
		return nil, s.finishUpdateRollback(ctx, err, tagMappingSyncResult{}, appliedUpdatePlans, rule, rollbackLinks, nil)
	}
	for _, plan := range updatePlans {
		// Rule defaults are resolver inputs, never edits to a surviving point's intent
		// or legacy DataFormat policy. New identities are created separately below.
		dataType := next.DataType
		updateReq := point.UpdatePointRequest{
			Name:     &plan.name,
			DataType: &dataType,
			Enabled:  &effectiveEnabled,
		}
		if plan.point.PollingGroupID == nil && defaultPollingGroupID != nil {
			updateReq.PollingGroupID = defaultPollingGroupID
		}
		pointRecord, updateErr := s.pointSvc.UpdateMaterializedObserved(ctx, plan.point, updateReq, s.observePointWrite(ctx))
		if updateErr != nil {
			return nil, s.finishUpdateRollback(ctx, fmt.Errorf("更新衍生點位失敗: %w", updateErr), tagMappingSyncResult{}, appliedUpdatePlans, rule, rollbackLinks, nil)
		}
		plan.appliedRevision = pointRecord.RevisionID
		appliedUpdatePlans = append(appliedUpdatePlans, plan)
		s.syncPoints(ctx, []*schema.Point{pointRecord})
	}

	createdPointIDs := make([]string, 0, len(createPlans))
	createdRuntimePoints := []*schema.Point{}
	createdLinkMap := make(map[string]*schema.SourceRuleLink, len(createPlans))
	for _, plan := range createPlans {
		pointRecord, createErr := s.createRulePoint(ctx, point.CreatePointRequest{
			DeviceID:       next.DeviceID,
			Name:           buildPointName(next.NamingPrefix, plan.address),
			Address:        plan.address,
			DataType:       next.DataType,
			DataFormat:     next.DataFormat,
			Mode:           schema.PointModeReadOnly,
			PollingGroupID: defaultPollingGroupID,
		})
		if createErr != nil {
			return nil, s.finishUpdateRollback(ctx, fmt.Errorf("建立新增衍生點位失敗: %w", createErr), tagMappingSyncResult{}, appliedUpdatePlans, rule, rollbackLinks, createdPointIDs)
		}
		if !effectiveEnabled {
			createdPointID := pointRecord.ID
			pointRecord, createErr = s.updateRulePoint(ctx, pointRecord, point.UpdatePointRequest{Enabled: &effectiveEnabled})
			if createErr != nil {
				return nil, s.finishUpdateRollback(ctx, fmt.Errorf("更新新增衍生點位狀態失敗: %w", createErr), tagMappingSyncResult{}, appliedUpdatePlans, rule, rollbackLinks, append(createdPointIDs, createdPointID))
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
		createdRuntimePoints = append(createdRuntimePoints, pointRecord)
	}

	for _, address := range filteredAddresses {
		if link, exists := createdLinkMap[normalizeAddressKey(address)]; exists {
			newLinks = append(newLinks, link)
		}
	}

	syncResult, err := s.syncRuleTagMappings(ctx, rule, &next, newLinks, effectiveEnabled)
	if err != nil {
		return nil, s.finishUpdateRollback(ctx, fmt.Errorf("同步來源規則標籤映射失敗: %w", err), syncResult, appliedUpdatePlans, rule, rollbackLinks, createdPointIDs)
	}

	if err := s.repo.ReplaceLinks(ctx, next.ID, newLinks); err != nil {
		return nil, s.finishUpdateRollback(ctx, fmt.Errorf("重建來源規則連結失敗: %w", err), syncResult, appliedUpdatePlans, rule, rollbackLinks, createdPointIDs)
	}
	if err := s.persistCandidateSnapshots(ctx, &next, newLinks); err != nil {
		return nil, s.finishUpdateRollback(ctx, fmt.Errorf("持久化來源規則候選快照失敗: %w", err), syncResult, appliedUpdatePlans, rule, rollbackLinks, createdPointIDs)
	}

	for _, link := range removedLinks {
		if err := s.cleanupRuleLinkResources(ctx, &next, link); err != nil {
			return nil, s.finishRollbackErrors(ctx, fmt.Errorf("刪除已移除衍生點位失敗: %w", err), s.repo.DeleteCandidateSnapshots(ctx, next.ID, next.RevisionID))
		}
	}

	s.syncPoints(ctx, createdRuntimePoints)

	return &next, nil
}

package sourcerule

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
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
}

type Service struct {
	repo        Repository
	deviceSvc   *device.Service
	pointSvc    *point.Service
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
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.Create(ctx, rule); err != nil {
		return nil, fmt.Errorf("建立來源規則失敗: %w", err)
	}

	createdPointIDs := make([]string, 0, req.Count)
	links := make([]*schema.SourceRuleLink, 0, req.Count)
	createdPoints := make([]*schema.Point, 0, req.Count)

	for _, address := range buildPlannedPointAddresses(rule.StartAddress, rule.Count, rule.DataType, deviceRecord.Protocol) {
		if containsAddress(req.SkippedAddresses, address) {
			continue
		}

		pointRecord, createErr := s.pointSvc.Create(ctx, point.CreatePointRequest{
			DeviceID:    rule.DeviceID,
			Name:        buildPointName(rule.NamingPrefix, address),
			Address:     address,
			DataType:    rule.DataType,
			Mode:        schema.PointModeReadOnly,
			Function:    "",
			Description: "",
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

	if len(links) > 0 {
		if err := s.repo.CreateLinks(ctx, links); err != nil {
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
	if err := s.repo.Update(ctx, &next); err != nil {
		return nil, fmt.Errorf("更新來源規則失敗: %w", err)
	}

	appliedUpdatePlans := make([]pointUpdatePlan, 0, len(updatePlans))
	for _, plan := range updatePlans {
		dataType := next.DataType
		pointRecord, updateErr := s.pointSvc.Update(ctx, plan.point.ID, point.UpdatePointRequest{
			Name:     &plan.name,
			DataType: &dataType,
			Enabled:  &effectiveEnabled,
		})
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
			DeviceID: next.DeviceID,
			Name:     buildPointName(next.NamingPrefix, plan.address),
			Address:  plan.address,
			DataType: next.DataType,
			Mode:     schema.PointModeReadOnly,
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

	if err := s.repo.DeleteLinks(ctx, next.ID); err != nil {
		s.rollbackUpdatedPoints(ctx, appliedUpdatePlans)
		s.rollbackRuleState(ctx, rule, links)
		s.rollbackCreatedPoints(ctx, createdPointIDs)
		return nil, fmt.Errorf("清除舊來源規則連結失敗: %w", err)
	}
	if len(newLinks) > 0 {
		if err := s.repo.CreateLinks(ctx, newLinks); err != nil {
			s.rollbackUpdatedPoints(ctx, appliedUpdatePlans)
			s.rollbackRuleState(ctx, rule, links)
			s.rollbackCreatedPoints(ctx, createdPointIDs)
			return nil, fmt.Errorf("重建來源規則連結失敗: %w", err)
		}
	}

	for _, link := range removedLinks {
		if err := s.pointSvc.Delete(ctx, link.PointID); err != nil {
			return nil, fmt.Errorf("刪除已移除衍生點位失敗: %w", err)
		}
		if s.runtimeSync != nil {
			s.runtimeSync.RemovePoint(link.PointID)
		}
	}

	return &next, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	links, err := s.repo.ListLinks(ctx, id)
	if err != nil {
		return fmt.Errorf("取得來源規則連結失敗: %w", err)
	}
	for _, link := range links {
		if err := s.pointSvc.Delete(ctx, link.PointID); err != nil {
			return fmt.Errorf("刪除衍生點位失敗: %w", err)
		}
		if s.runtimeSync != nil {
			s.runtimeSync.RemovePoint(link.PointID)
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
		if err := s.syncRuleLinksEnabled(ctx, rule.ID, enabled); err != nil {
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

func (s *Service) syncRuleLinksEnabled(ctx context.Context, ruleID string, enabled bool) error {
	links, err := s.repo.ListLinks(ctx, ruleID)
	if err != nil {
		return fmt.Errorf("取得來源規則連結失敗: %w", err)
	}

	for _, link := range links {
		pointRecord, getErr := s.pointSvc.GetByID(ctx, link.PointID)
		if getErr != nil {
			return fmt.Errorf("取得衍生點位失敗: %w", getErr)
		}
		pointRecord, getErr = s.pointSvc.Update(ctx, link.PointID, point.UpdatePointRequest{
			Enabled: &enabled,
		})
		if getErr != nil {
			return fmt.Errorf("更新衍生點位狀態失敗: %w", getErr)
		}
		s.syncPoints([]*schema.Point{pointRecord})
	}
	return nil
}

func (s *Service) setEnabled(ctx context.Context, id string, enabled bool) error {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得來源規則失敗: %w", err)
	}

	if enabled {
		deviceRecord, getErr := s.deviceSvc.GetByID(ctx, rule.DeviceID)
		if getErr != nil {
			return fmt.Errorf("取得設備失敗: %w", getErr)
		}
		if deviceRecord.Status != schema.DeviceStatusActive {
			return fmt.Errorf("設備尚未通過 probe readiness，不能啟用來源規則")
		}
	}

	rule.Enabled = enabled
	rule.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, rule); err != nil {
		return fmt.Errorf("更新來源規則失敗: %w", err)
	}
	return s.syncRuleLinksEnabled(ctx, id, enabled)
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
	return nil
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
			Name:     &name,
			DataType: &dataType,
			Enabled:  &enabled,
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

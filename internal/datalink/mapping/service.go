// Package mapping 提供點位到標籤的映射管線功能。
//
// 本套件實作映射的 CRUD 操作、轉換管線執行和預覽功能。
package mapping

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// Repository 介面定義
// =============================================================================

// Repository 映射資料存取介面
type Repository interface {
	Create(ctx context.Context, mapping *schema.Mapping) error
	Update(ctx context.Context, mapping *schema.Mapping) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*schema.Mapping, error)
	GetByPointID(ctx context.Context, pointID string) ([]*schema.Mapping, error)
	GetByTagID(ctx context.Context, tagID string) ([]*schema.Mapping, error)
	List(ctx context.Context, filter ListFilter) ([]*schema.Mapping, error)
}

// ListFilter 映射列表篩選條件
type ListFilter struct {
	PointID *string
	TagID   *string
	Enabled *bool
	Limit   int
	Offset  int
}

// =============================================================================
// Service 服務層
// =============================================================================

// Service 映射管理服務
type Service struct {
	repo Repository
	mu   sync.RWMutex
}

// NewService 建立新的映射服務
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// =============================================================================
// CRUD 操作
// =============================================================================

// Create 建立新映射
func (s *Service) Create(ctx context.Context, req CreateMappingRequest) (*schema.Mapping, error) {
	if err := ValidateTransformPipeline(req.TransformPipeline); err != nil {
		return nil, fmt.Errorf("轉換管線無效: %w", err)
	}

	pipelineJSON, err := json.Marshal(req.TransformPipeline)
	if err != nil {
		return nil, fmt.Errorf("序列化轉換管線失敗: %w", err)
	}

	mapping := &schema.Mapping{
		ID:                generateUUID(),
		PointID:           req.PointID,
		TagID:             req.TagID,
		TransformPipeline: string(pipelineJSON),
		Enabled:           true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.repo.Create(ctx, mapping); err != nil {
		return nil, fmt.Errorf("建立映射失敗: %w", err)
	}

	return mapping, nil
}

// CreateMappingRequest 建立映射請求
type CreateMappingRequest struct {
	PointID           string                 `json:"point_id"`
	TagID             string                 `json:"tag_id"`
	TransformPipeline []schema.TransformStep `json:"transform_pipeline"`
}

// Update 更新映射
func (s *Service) Update(ctx context.Context, id string, req UpdateMappingRequest) (*schema.Mapping, error) {
	mapping, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得映射失敗: %w", err)
	}

	if req.TransformPipeline != nil {
		if err := ValidateTransformPipeline(req.TransformPipeline); err != nil {
			return nil, fmt.Errorf("轉換管線無效: %w", err)
		}
		pipelineJSON, err := json.Marshal(req.TransformPipeline)
		if err != nil {
			return nil, fmt.Errorf("序列化轉換管線失敗: %w", err)
		}
		mapping.TransformPipeline = string(pipelineJSON)
	}

	if req.Enabled != nil {
		mapping.Enabled = *req.Enabled
	}

	mapping.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, mapping); err != nil {
		return nil, fmt.Errorf("更新映射失敗: %w", err)
	}

	return mapping, nil
}

// UpdateMappingRequest 更新映射請求
type UpdateMappingRequest struct {
	TransformPipeline []schema.TransformStep `json:"transform_pipeline,omitempty"`
	Enabled           *bool                  `json:"enabled,omitempty"`
}

// Delete 刪除映射
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除映射失敗: %w", err)
	}
	return nil
}

// GetByID 根據 ID 取得映射
func (s *Service) GetByID(ctx context.Context, id string) (*schema.Mapping, error) {
	mapping, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得映射失敗: %w", err)
	}
	return mapping, nil
}

// List 列出映射
func (s *Service) List(ctx context.Context, filter ListFilter) ([]*schema.Mapping, error) {
	mappings, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出映射失敗: %w", err)
	}
	return mappings, nil
}

// =============================================================================
// 轉換管線驗證
// =============================================================================

// ValidateTransformPipeline 驗證轉換管線
func ValidateTransformPipeline(steps []schema.TransformStep) error {
	validTypes := map[schema.TransformType]bool{
		schema.TransformDecode:      true,
		schema.TransformCast:        true,
		schema.TransformScale:       true,
		schema.TransformLookup:      true,
		schema.TransformConditional: true,
		schema.TransformFormula:     true,
	}

	for i, step := range steps {
		if step.Order < 0 {
			return fmt.Errorf("步驟 %d: order 不能為負數", i+1)
		}
		if !validTypes[step.Type] {
			return fmt.Errorf("步驟 %d: 不支援的轉換類型 '%s'", i+1, step.Type)
		}

		switch step.Type {
		case schema.TransformScale:
			if step.Params == nil {
				return fmt.Errorf("步驟 %d: scale 需要參數", i+1)
			}
		case schema.TransformCast:
			if step.Params == nil {
				return fmt.Errorf("步驟 %d: cast 需要參數", i+1)
			}
		case schema.TransformFormula:
			if step.Params == nil {
				return fmt.Errorf("步驟 %d: formula 需要參數", i+1)
			}
		case schema.TransformConditional:
			if step.Params == nil {
				return fmt.Errorf("步驟 %d: conditional 需要參數", i+1)
			}
			if _, exists := step.Params["true_value"]; !exists {
				return fmt.Errorf("步驟 %d: conditional 缺少 true_value", i+1)
			}
			if _, exists := step.Params["false_value"]; !exists {
				return fmt.Errorf("步驟 %d: conditional 缺少 false_value", i+1)
			}
			if _, _, err := parseConditionalParams(step.Params); err != nil {
				return fmt.Errorf("步驟 %d: conditional 參數無效: %w", i+1, err)
			}
		}
	}

	return nil
}

// =============================================================================
// 轉換管線執行
// =============================================================================

// TransformContext 轉換上下文
type TransformContext struct {
	RawValue     interface{}
	CurrentValue interface{}
	StepResults  []StepResult
	Error        error
}

// StepResult 步驟結果
type StepResult struct {
	StepIndex int         `json:"step_index"`
	StepType  string      `json:"step_type"`
	Input     interface{} `json:"input"`
	Output    interface{} `json:"output"`
	Error     string      `json:"error,omitempty"`
}

// ExecutePipeline 執行轉換管線
func ExecutePipeline(raw interface{}, pipelineJSON string) (*TransformContext, error) {
	ctx := &TransformContext{
		RawValue:     raw,
		CurrentValue: raw,
		StepResults:  make([]StepResult, 0),
	}

	var steps []schema.TransformStep
	if err := json.Unmarshal([]byte(pipelineJSON), &steps); err != nil {
		ctx.Error = fmt.Errorf("解析轉換管線失敗: %w", err)
		return ctx, ctx.Error
	}

	orderedSteps := normalizeTransformSteps(steps)
	for i, step := range orderedSteps {
		result := StepResult{
			StepIndex: i,
			StepType:  string(step.Type),
			Input:     ctx.CurrentValue,
		}

		output, err := executeStep(ctx.CurrentValue, step)
		if err != nil {
			result.Error = err.Error()
			ctx.StepResults = append(ctx.StepResults, result)
			ctx.Error = fmt.Errorf("步驟 %d (%s) 失敗: %w", i+1, step.Type, err)
			return ctx, ctx.Error
		}

		result.Output = output
		ctx.StepResults = append(ctx.StepResults, result)
		ctx.CurrentValue = output
	}

	return ctx, nil
}

// executeStep 執行單一轉換步驟
func executeStep(input interface{}, step schema.TransformStep) (interface{}, error) {
	switch step.Type {
	case schema.TransformDecode:
		return executeDecode(input, step.Params)
	case schema.TransformCast:
		return executeCast(input, step.Params)
	case schema.TransformScale:
		return executeScale(input, step.Params)
	case schema.TransformLookup:
		return executeLookup(input, step.Params)
	case schema.TransformConditional:
		return executeConditional(input, step.Params)
	case schema.TransformFormula:
		return executeFormula(input, step.Params)
	default:
		return nil, fmt.Errorf("不支援的轉換類型: %s", step.Type)
	}
}

// =============================================================================
// 轉換步驟實作 (使用 map[string]interface{} 參數)
// =============================================================================

// executeDecode 解碼轉換
func executeDecode(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	format, _ := params["format"].(string)
	switch format {
	case "swap16":
		if v, ok := toUint16Value(input); ok {
			return ((v & 0xFF) << 8) | ((v >> 8) & 0xFF), nil
		}
	case "swap32":
		if v, ok := toUint32Value(input); ok {
			return ((v&0xFF)<<24 | ((v>>8)&0xFF)<<16 | ((v>>16)&0xFF)<<8 | (v >> 24)), nil
		}
	}

	return input, nil
}

// executeCast 型別轉換
func executeCast(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	toType, _ := params["to_type"].(string)
	if toType == "" {
		toType, _ = params["target_type"].(string)
	}
	targetType := schema.DataType(toType)

	switch targetType {
	case schema.DataTypeBool:
		return toBoolValue(input), nil
	case schema.DataTypeInt16:
		v, _ := toFloat64Value(input)
		return int16(v), nil
	case schema.DataTypeUint16:
		v, _ := toFloat64Value(input)
		return uint16(v), nil
	case schema.DataTypeInt32:
		v, _ := toFloat64Value(input)
		return int32(v), nil
	case schema.DataTypeUint32:
		v, _ := toFloat64Value(input)
		return uint32(v), nil
	case schema.DataTypeInt64:
		v, _ := toFloat64Value(input)
		return int64(v), nil
	case schema.DataTypeUint64:
		v, _ := toFloat64Value(input)
		return uint64(v), nil
	case schema.DataTypeFloat32:
		v, _ := toFloat64Value(input)
		return float32(v), nil
	case schema.DataTypeFloat64:
		v, _ := toFloat64Value(input)
		return v, nil
	case schema.DataTypeString:
		return fmt.Sprintf("%v", input), nil
	default:
		return input, nil
	}
}

// executeScale 縮放轉換
func executeScale(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	v, ok := toFloat64Value(input)
	if !ok {
		return input, fmt.Errorf("無法轉換為數值: %v", input)
	}

	multiplier := 1.0
	if m, ok := params["multiplier"].(float64); ok {
		multiplier = m
	} else if m, ok := params["scale"].(float64); ok {
		multiplier = m
	}

	divisor := 1.0
	if d, ok := params["divisor"].(float64); ok && d != 0 {
		divisor = d
	}

	offset := 0.0
	if o, ok := params["offset"].(float64); ok {
		offset = o
	}

	result := (v*multiplier)/divisor + offset

	if minVal, ok := params["min"].(float64); ok && result < minVal {
		result = minVal
	}
	if maxVal, ok := params["max"].(float64); ok && result > maxVal {
		result = maxVal
	}

	return result, nil
}

// executeLookup 查表轉換
func executeLookup(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	table, _ := params["table"].(map[string]interface{})
	if table == nil {
		return input, nil
	}

	key := fmt.Sprintf("%v", input)

	if value, exists := table[key]; exists {
		return value, nil
	}

	if defaultVal, exists := params["default"]; exists {
		return defaultVal, nil
	}

	return input, nil
}

// executeConditional 條件轉換
func executeConditional(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	v, _ := toFloat64Value(input)

	operator, threshold, err := parseConditionalParams(params)
	if err != nil {
		return input, err
	}
	trueValue := params["true_value"]
	falseValue := params["false_value"]

	result := evaluateCondition(v, operator, threshold)

	if result {
		return trueValue, nil
	}
	return falseValue, nil
}

// evaluateCondition 評估條件
func evaluateCondition(value float64, operator string, threshold float64) bool {
	switch operator {
	case "==", "eq":
		return value == threshold
	case "!=", "ne":
		return value != threshold
	case ">", "gt":
		return value > threshold
	case ">=", "gte":
		return value >= threshold
	case "<", "lt":
		return value < threshold
	case "<=", "lte":
		return value <= threshold
	default:
		return false
	}
}

func parseConditionalParams(params map[string]interface{}) (string, float64, error) {
	if params == nil {
		return "", 0, fmt.Errorf("缺少參數")
	}

	if condition, ok := params["condition"].(string); ok && strings.TrimSpace(condition) != "" {
		return parseConditionString(condition)
	}

	operator, _ := params["operator"].(string)
	operator = normalizeOperator(strings.TrimSpace(operator))
	if operator == "" {
		return "", 0, fmt.Errorf("缺少 operator 或 condition")
	}

	threshold, ok := toFloat64Value(params["threshold"])
	if !ok {
		return "", 0, fmt.Errorf("無效的 threshold")
	}

	return operator, threshold, nil
}

func normalizeOperator(operator string) string {
	switch operator {
	case "eq":
		return "=="
	case "ne":
		return "!="
	case "gt":
		return ">"
	case "gte":
		return ">="
	case "lt":
		return "<"
	case "lte":
		return "<="
	default:
		return operator
	}
}

func parseConditionString(condition string) (string, float64, error) {
	cond := strings.TrimSpace(condition)
	if strings.HasPrefix(cond, "value") {
		cond = strings.TrimSpace(strings.TrimPrefix(cond, "value"))
	}
	if strings.HasPrefix(cond, "x") {
		cond = strings.TrimSpace(strings.TrimPrefix(cond, "x"))
	}

	operators := []string{"==", "!=", ">=", "<=", ">", "<"}
	for _, op := range operators {
		idx := strings.Index(cond, op)
		if idx < 0 {
			continue
		}
		left := strings.TrimSpace(cond[:idx])
		right := strings.TrimSpace(cond[idx+len(op):])
		if left != "" && left != "value" && left != "x" {
			return "", 0, fmt.Errorf("無法解析 condition: %s", condition)
		}
		threshold, ok := toFloat64Value(right)
		if !ok {
			return "", 0, fmt.Errorf("無法解析門檻值: %s", right)
		}
		return op, threshold, nil
	}

	return "", 0, fmt.Errorf("無法解析 condition: %s", condition)
}

// executeFormula 公式轉換
func executeFormula(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	expression, _ := params["expression"].(string)
	if expression == "" {
		return input, nil
	}

	v, ok := toFloat64Value(input)
	if !ok {
		return input, fmt.Errorf("無法轉換為數值: %v", input)
	}

	expr := strings.TrimSpace(expression)
	expr = strings.ReplaceAll(expr, "value", strconv.FormatFloat(v, 'f', -1, 64))
	expr = strings.ReplaceAll(expr, "x", strconv.FormatFloat(v, 'f', -1, 64))

	result, err := evaluateSimpleExpression(expr, v)
	if err != nil {
		return input, err
	}

	return result, nil
}

// evaluateSimpleExpression 評估簡易表達式
func evaluateSimpleExpression(expr string, value float64) (float64, error) {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "abs(") {
		return math.Abs(value), nil
	}
	if strings.HasPrefix(expr, "sqrt(") {
		return math.Sqrt(value), nil
	}
	if strings.HasPrefix(expr, "round(") {
		return math.Round(value), nil
	}
	if strings.HasPrefix(expr, "floor(") {
		return math.Floor(value), nil
	}
	if strings.HasPrefix(expr, "ceil(") {
		return math.Ceil(value), nil
	}

	if v, err := strconv.ParseFloat(expr, 64); err == nil {
		return v, nil
	}

	return value, nil
}

// =============================================================================
// 輔助函數
// =============================================================================

func toFloat64Value(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case bool:
		if val {
			return 1, true
		}
		return 0, true
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func toUint16Value(v interface{}) (uint16, bool) {
	f, ok := toFloat64Value(v)
	if ok {
		return uint16(f), true
	}
	return 0, false
}

func toUint32Value(v interface{}) (uint32, bool) {
	f, ok := toFloat64Value(v)
	if ok {
		return uint32(f), true
	}
	return 0, false
}

func toBoolValue(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case int, int16, int32, int64, uint, uint16, uint32, uint64, float32, float64:
		f, _ := toFloat64Value(v)
		return f != 0
	case string:
		return val == "true" || val == "1" || val == "on" || val == "yes"
	}
	return false
}

func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func normalizeTransformSteps(steps []schema.TransformStep) []schema.TransformStep {
	if len(steps) == 0 {
		return steps
	}

	needsOrdering := false
	for _, step := range steps {
		if step.Order != 0 {
			needsOrdering = true
			break
		}
	}
	if !needsOrdering {
		return steps
	}

	ordered := make([]schema.TransformStep, len(steps))
	copy(ordered, steps)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Order < ordered[j].Order
	})
	return ordered
}

// =============================================================================
// 記憶體 Repository
// =============================================================================

// MemoryRepository 記憶體內映射儲存庫
type MemoryRepository struct {
	mu       sync.RWMutex
	mappings map[string]*schema.Mapping
}

// NewMemoryRepository 建立新的記憶體儲存庫
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		mappings: make(map[string]*schema.Mapping),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, mapping *schema.Mapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.mappings[mapping.ID]; exists {
		return fmt.Errorf("映射 ID 已存在: %s", mapping.ID)
	}

	mappingCopy := *mapping
	r.mappings[mapping.ID] = &mappingCopy
	return nil
}

func (r *MemoryRepository) Update(ctx context.Context, mapping *schema.Mapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.mappings[mapping.ID]; !exists {
		return fmt.Errorf("映射不存在: %s", mapping.ID)
	}

	mappingCopy := *mapping
	r.mappings[mapping.ID] = &mappingCopy
	return nil
}

func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.mappings[id]; !exists {
		return fmt.Errorf("映射不存在: %s", id)
	}

	delete(r.mappings, id)
	return nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*schema.Mapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mapping, exists := r.mappings[id]
	if !exists {
		return nil, fmt.Errorf("映射不存在: %s", id)
	}

	mappingCopy := *mapping
	return &mappingCopy, nil
}

func (r *MemoryRepository) GetByPointID(ctx context.Context, pointID string) ([]*schema.Mapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Mapping, 0)
	for _, mapping := range r.mappings {
		if mapping.PointID == pointID {
			mappingCopy := *mapping
			result = append(result, &mappingCopy)
		}
	}
	return result, nil
}

func (r *MemoryRepository) GetByTagID(ctx context.Context, tagID string) ([]*schema.Mapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Mapping, 0)
	for _, mapping := range r.mappings {
		if mapping.TagID == tagID {
			mappingCopy := *mapping
			result = append(result, &mappingCopy)
		}
	}
	return result, nil
}

func (r *MemoryRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Mapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Mapping, 0, len(r.mappings))

	for _, mapping := range r.mappings {
		if filter.PointID != nil && mapping.PointID != *filter.PointID {
			continue
		}
		if filter.TagID != nil && mapping.TagID != *filter.TagID {
			continue
		}
		if filter.Enabled != nil && mapping.Enabled != *filter.Enabled {
			continue
		}

		mappingCopy := *mapping
		result = append(result, &mappingCopy)
	}

	if filter.Offset > 0 {
		if filter.Offset >= len(result) {
			return []*schema.Mapping{}, nil
		}
		result = result[filter.Offset:]
	}
	if filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}

	return result, nil
}

// Count 計算映射數量
func (r *MemoryRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.mappings)
}

// Clear 清空所有映射
func (r *MemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mappings = make(map[string]*schema.Mapping)
}

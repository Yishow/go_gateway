package point

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/schema"
)

var (
	pointAddressPattern  = regexp.MustCompile(`^([A-Za-z]+\d+|\d+)$`)
	alphaFunctionPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

func (s *Service) validateBatchCreate(ctx context.Context, req BatchCreatePointsRequest) ([]string, error) {
	errors := make([]string, 0)

	if strings.TrimSpace(req.DeviceID) == "" {
		errors = append(errors, pointValidationError("device_id", req.DeviceID, "不能為空").Error())
	}
	if len(req.Points) == 0 {
		errors = append(errors, pointValidationError("points", len(req.Points), "不能為空").Error())
	}
	if err := validatePointDataType(req.DataType); err != nil {
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		return errors, nil
	}

	existingPoints, err := s.repo.ListByDevice(ctx, req.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("查詢既有點位失敗: %w", err)
	}

	existingByKey := make(map[string]struct{}, len(existingPoints))
	for _, pt := range existingPoints {
		existingByKey[buildBatchKey(pt.Address, pt.Function)] = struct{}{}
	}

	seenByKey := make(map[string]int, len(req.Points))
	for i, item := range req.Points {
		lineNo := i + 1
		name := strings.TrimSpace(item.Name)

		if name == "" {
			errors = append(errors, batchPointValidationError(lineNo, pointValidationError("name", item.Name, "不能為空")))
		}

		address, addrErr := validateAndNormalizeAddress(item.Address)
		if addrErr != nil {
			errors = append(errors, batchPointValidationError(lineNo, addrErr))
			continue
		}

		function, functionErr := validateAndNormalizeFunction(item.Function)
		if functionErr != nil {
			errors = append(errors, batchPointValidationError(lineNo, functionErr))
			continue
		}

		key := buildBatchKey(address, function)
		if first, exists := seenByKey[key]; exists {
			errors = append(errors, batchPointValidationError(lineNo, pointValidationError("address+function", key, fmt.Sprintf("與第 %d 筆資料重複", first))))
		} else {
			seenByKey[key] = lineNo
		}
		if _, exists := existingByKey[key]; exists {
			errors = append(errors, batchPointValidationError(lineNo, pointValidationError("address+function", key, "已存在於設備內")))
		}
	}

	return errors, nil
}

func buildBatchKey(address, function string) string {
	return normalizeAddress(address) + "|" + normalizeFunction(function)
}

func normalizeAddress(address string) string {
	return strings.ToUpper(strings.TrimSpace(address))
}

func validateAndNormalizeAddress(address string) (string, error) {
	normalized := normalizeAddress(address)
	if normalized == "" {
		return "", pointValidationError("address", address, "不能為空")
	}
	if !pointAddressPattern.MatchString(normalized) {
		return "", pointValidationError("address", address, "格式無效")
	}
	return normalized, nil
}

func normalizeFunction(function string) string {
	raw := strings.ToLower(strings.TrimSpace(function))
	switch raw {
	case "", "3", "03", "fc03", "holding", "read_holding":
		return "03"
	case "1", "01", "fc01", "coil":
		return "01"
	case "2", "02", "fc02", "discrete":
		return "02"
	case "4", "04", "fc04", "input", "read_input":
		return "04"
	default:
		return strings.ToUpper(strings.TrimSpace(function))
	}
}

func validateAndNormalizeFunction(function string) (string, error) {
	raw := strings.TrimSpace(function)
	if !isValidPointFunction(raw) {
		return "", pointValidationError("function", function, "格式無效")
	}
	return normalizeFunction(raw), nil
}

func isValidPointFunction(function string) bool {
	raw := strings.TrimSpace(function)
	if raw == "" {
		return true
	}

	normalized := strings.ToLower(raw)
	switch normalized {
	case "01", "02", "03", "04", "fc01", "fc02", "fc03", "fc04", "coil", "discrete", "holding", "input", "read_holding", "read_input":
		return true
	}

	if alphaFunctionPattern.MatchString(raw) {
		return true
	}
	if len(normalized) == 2 {
		if v, err := strconv.Atoi(normalized); err == nil {
			return v >= 1 && v <= 16
		}
	}
	return false
}

func validatePointDataType(dataType schema.DataType) error {
	if isValidDataType(dataType) {
		return nil
	}
	return pointValidationError("data_type", dataType, "不支援的資料型別")
}

func pointValidationError(field string, value interface{}, reason string) error {
	return fmt.Errorf("point 驗證失敗: field=%s value=%v reason=%s", field, value, reason)
}

func batchPointValidationError(index int, err error) string {
	return fmt.Sprintf("第 %d 筆資料 %s", index, err.Error())
}

package sourcerule

import (
	"fmt"
	"strings"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
)

const (
	candidateRevisionConflictMessage = "candidate revision conflict"
	candidateRetryAction             = "refresh workspace candidates and retry"
	localModbusStrideAction          = "increase share_stride to match the datatype register span"
)

func validateCreateRequest(req CreateRuleRequest) error {
	if strings.TrimSpace(req.DeviceID) == "" {
		return validationError("device_id is required")
	}
	if strings.TrimSpace(req.StartAddress) == "" {
		return validationError("start_address is required")
	}
	if req.Count <= 0 {
		return validationError("count must be greater than zero")
	}
	if err := validateRuleDataType(req.DataType); err != nil {
		return err
	}
	if err := validateRuleDataFormat(req.DataFormat); err != nil {
		return err
	}
	return validateShareConfig(req.ShareEnabled, req.ShareStartRegister, req.ShareStride, req.DataType)
}

func validateShareConfig(enabled bool, startRegister, stride *int, dataType schema.DataType) error {
	if !enabled {
		return nil
	}
	if startRegister != nil && *startRegister < 40001 {
		return validationError("share_start_register must be greater than or equal to 40001")
	}
	if stride != nil && *stride < 1 {
		return validationError("share_stride must be greater than zero")
	}
	if stride != nil && *stride < modbusshare.DataTypeSpan(dataType) {
		return &modbusshare.Error{
			Code:      modbusshare.ErrCodeInvalidGeometry,
			Message:   fmt.Sprintf("share_stride %d is smaller than datatype span %d", *stride, modbusshare.DataTypeSpan(dataType)),
			Retryable: false,
			Action:    localModbusStrideAction,
		}
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
		return validationError(fmt.Sprintf("unsupported data_format: %s", dataFormat))
	}
}

func validateRuleDataType(dataType schema.DataType) error {
	if dataType == "" {
		return validationError("data_type is required")
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
		return validationError(fmt.Sprintf("unsupported data_type: %s", dataType))
	}
}

func validationError(message string) error {
	return fmt.Errorf("%w: %s", ErrValidation, message)
}

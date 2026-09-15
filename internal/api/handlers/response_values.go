package handlers

const (
	apiResponseAppliedKey = "applied"

	apiValidationErrorCode = "validation"

	connectionNotFoundMessage               = "connection not found"
	connectionNotFoundOrNotConnectedMessage = "connection not found or not connected"
	connectionIDRequiredMessage             = "connection_id is required"
	notImplementedMessage                   = "not implemented"
	deviceNotFoundMessage                   = "Device not found"
	schedulerNotInitializedMessage          = "排程器未初始化"
	pointNotFoundMessage                    = "Point not found"
	studioV2MappingRowNotFoundMessage       = "Studio V2 mapping row not found"
	studioV2SourceRuleNotFoundMessage       = "Studio V2 source rule not found"

	runtimeStatusIdle            = "idle"
	runtimeApplyNotRunningStatus = "not_running"
	runtimeApplyAppliedStatus    = "applied"
	workspaceSaveStateSaved      = "saved"
)

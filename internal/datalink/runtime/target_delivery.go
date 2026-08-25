package runtime

import (
	"context"
	"time"
)

// TargetName identifies one production delivery target.
type TargetName string

const (
	// TargetDatabase is the durable database output target.
	TargetDatabase TargetName = "database"
	// TargetModbusShare is the local Modbus Share output target.
	TargetModbusShare TargetName = "modbus_share"
)

// TargetDeliveryOutcome records one target's independent delivery result.
// Err is nil when that target accepted the value.
type TargetDeliveryOutcome struct {
	Target TargetName
	Err    error
}

// TargetOutcomeWriter is an optional richer target seam. Implementations must
// attempt each configured target and return one outcome per target, even when
// another target fails.
type TargetOutcomeWriter interface {
	WriteTagValueOutcomes(context.Context, string, any, time.Time) []TargetDeliveryOutcome
}

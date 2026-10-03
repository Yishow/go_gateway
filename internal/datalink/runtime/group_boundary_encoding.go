package runtime

import (
	"errors"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/snapshot"
)

// Only a known value-representation refusal retains the existing quality-skip
// policy. Structural or unclassified errors leave durable closure/input open.
func structuralEncodingFailure(cause error) bool {
	var rowErr *dbtarget.GroupRowError
	return !errors.As(cause, &rowErr) || rowErr.Code != "sql-value-blocked"
}

func blockedOutcome(outcome snapshot.Outcome, cause error) snapshot.Outcome {
	code := "unknown"
	var rowErr *dbtarget.GroupRowError
	if errors.As(cause, &rowErr) {
		code = rowErr.Code
	}
	outcome.Kind = snapshot.OutcomeSkipped
	outcome.Reason = reasonEncodeBlocked + code
	outcome.EffectKey = ""
	outcome.Partial = false
	return outcome
}

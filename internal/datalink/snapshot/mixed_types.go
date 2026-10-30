package snapshot

import "time"

// Mixed modes dispatch the union of independently aligned member windows.
const (
	ModeMixedWide Mode = "mixed_wide"
	ModeMixedLong Mode = "mixed_long"
)

func (m Mode) mixed() bool { return m == ModeMixedWide || m == ModeMixedLong }

// Calculation is explicit mixed intent; it is never inferred from a source role.
type Calculation string

const (
	CalculationLatest Calculation = "latest"
	CalculationMean   Calculation = "mean"
	CalculationDelta  Calculation = "delta"
)

// MemberWindow retains the member's own aligned and actual bounds. A not-due
// result has only its intent and NextDue; it does not claim a computed window.
type MemberWindow struct {
	Calculation              Calculation
	Interval                 time.Duration
	AlignedStart, AlignedEnd time.Time
	ActualStart, ActualEnd   time.Time
	Truncated                string
	NextDue                  time.Time
}

// DispatchWatermark is mixed-only progress, not the legacy bucket-start cursor.
// ClosedThrough is the last committed union boundary, never a common window.
type DispatchWatermark struct {
	WorkspaceID, GroupID, GroupRevision string
	ClosedThrough                       time.Time
}

// MemberProgress is the next independent window, derived from dispatch progress.
// End is zero when an exact cutover leaves no future window for this member.
type MemberProgress struct{ Start, End time.Time }

const (
	MemberNotDue           MemberStatus = "not_due"
	ReasonNotDue                        = "not_due"
	ReasonWindowArithmetic              = "window-arithmetic"
)

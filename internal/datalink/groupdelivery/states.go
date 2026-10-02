package groupdelivery

// Outbox states. sql_committed is the only state backed by destination
// evidence; unknown and blocked never become success on their own.
const (
	StatePending     = "pending"
	StateSending     = "sending"
	StateRetrying    = "retrying"
	StateBlocked     = "blocked"
	StateUnknown     = "unknown"
	StateQuarantined = "quarantined"
	StateCommitted   = "sql_committed"
	// StateSkipped is an explicit operator decision to move past a quarantined
	// or blocked row; the payload is kept for audit and never delivered.
	StateSkipped = "operator_skipped"
)

// Dedupe capabilities of a destination.
const (
	DedupeReceipt   = "receipt"
	DedupeUniqueKey = "unique_key"
	DedupeNone      = "none"
)

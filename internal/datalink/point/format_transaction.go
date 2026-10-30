package point

// Keep database diagnostics private while retaining cleanup/unknown-outcome
// sentinels for callers deciding whether compensation is safe.
type formatTransactionError struct {
	cause     error
	attempted *FormatMutationResult
}

func (e *formatTransactionError) Error() string   { return ErrFormatUnavailable.Error() }
func (e *formatTransactionError) Unwrap() []error { return []error{ErrFormatUnavailable, e.cause} }
func formatTransactionFailure(cause error) error  { return &formatTransactionError{cause: cause} }

func formatAttemptFailure(cause error, attempted *FormatMutationResult) error {
	return &formatTransactionError{cause: cause, attempted: attempted}
}

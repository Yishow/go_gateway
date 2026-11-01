package datalink

// migrationCleanupError keeps cleanup causes inspectable without adding driver
// SQL, credentials or database paths to the migration's existing error text.
type migrationCleanupError struct{ cause error }

func (e migrationCleanupError) Error() string { return "migration cleanup failed" }
func (e migrationCleanupError) Unwrap() error { return e.cause }

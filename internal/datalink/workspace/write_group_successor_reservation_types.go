package workspace

// The reservation contract is private until operation-aware preparation and
// activation exist. A reserved result is neither reviewed nor schema-ready.
type writeGroupRecreateReservationRequest struct {
	WorkspaceID, ExpectedWorkspaceRevision string
	LineageID                              string
	ExpectedLineageRevision                int64
	PredecessorGroupID                     string
	ExpectedPredecessorRevision            string
	ExpectedPredecessorConnectorRevision   string
	ExpectedTargetConnectorRevision        string
	OperationID                            string
	Group                                  *WriteGroup
}

type writeGroupReservedDraft struct {
	WorkspaceRevision           string
	LineageID                   string
	LineageRevision, Generation int64
	PredecessorGroupID          string
	OperationID, IntentDigest   string
	Group                       *WriteGroup
}

type writeGroupReservationConflictError struct{ ReasonCode string }

func (e *writeGroupReservationConflictError) Error() string {
	return "write group successor reservation conflict"
}
func (e *writeGroupReservationConflictError) Unwrap() error { return ErrWriteGroupRevisionConflict }

func recreateReservationConflict(reason string) error {
	return &writeGroupReservationConflictError{ReasonCode: reason}
}

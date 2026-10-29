package workspace

// WriteGroupSuccessorOperationStatusRequest identifies the original reservation,
// not a current group head or the most recently issued preview.
type WriteGroupSuccessorOperationStatusRequest struct {
	WorkspaceID string `json:"workspace_id"`
	LineageID   string `json:"lineage_id"`
	OperationID string `json:"operation_id"`
}

// WriteGroupSuccessorOperationStatus contains historical operator-safe outcomes.
// A succeeded stage is not evidence of runtime running or committed target rows.
type WriteGroupSuccessorOperationStatus struct {
	WorkspaceID        string                              `json:"workspace_id"`
	LineageID          string                              `json:"lineage_id"`
	PredecessorGroupID string                              `json:"predecessor_group_id"`
	OperationID        string                              `json:"operation_id"`
	Reserve            WriteGroupSuccessorReserveStatus    `json:"reserve"`
	SchemaOperations   []WriteGroupSuccessorSchemaStatus   `json:"schema_operations"`
	Cutover            WriteGroupSuccessorCutoverStatus    `json:"cutover"`
	FirstStart         WriteGroupSuccessorFirstStartStatus `json:"first_start"`
}

// WriteGroupSuccessorReserveStatus preserves the original immutable result.
type WriteGroupSuccessorReserveStatus struct {
	Status string                          `json:"status"`
	Result *WriteGroupSuccessorPreparation `json:"result,omitempty"`
}

// WriteGroupSuccessorSchemaStatus omits token, SQL, private detail and raw errors.
type WriteGroupSuccessorSchemaStatus struct {
	OperationID string `json:"operation_id"`
	Status      string `json:"status"`
	ReasonCode  string `json:"reason_code,omitempty"`
	NextAction  string `json:"next_action,omitempty"`
}

// WriteGroupSuccessorCutoverStatus refers to the original validated event.
type WriteGroupSuccessorCutoverStatus struct {
	Status string             `json:"status"`
	Result *WriteGroupCutover `json:"result,omitempty"`
}

// WriteGroupSuccessorFirstStartStatus is separate from a historical cutover.
type WriteGroupSuccessorFirstStartStatus struct {
	Status string                               `json:"status"`
	Result *WriteGroupSuccessorFirstStartResult `json:"result,omitempty"`
}

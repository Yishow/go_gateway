package workspace

import "time"

// WriteGroupSuccessorFirstStartResult is one committed first activation, not a
// cutover from an invented predecessor version or proof of external SQL writes.
type WriteGroupSuccessorFirstStartResult struct {
	OperationID        string    `json:"operation_id"`
	RootOperationID    string    `json:"root_operation_id"`
	WorkspaceRevision  string    `json:"workspace_revision"`
	LineageID          string    `json:"lineage_id"`
	LineageRevision    int64     `json:"lineage_revision"`
	PredecessorGroupID string    `json:"predecessor_group_id"`
	GroupID            string    `json:"group_id"`
	AppliedRevision    string    `json:"applied_revision"`
	Boundary           time.Time `json:"boundary"`
	Generation         int64     `json:"generation"`
	Status             string    `json:"status"`
}

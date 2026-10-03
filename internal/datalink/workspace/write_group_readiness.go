package workspace

import (
	"context"
	"errors"

	"go-gateway/internal/datalink/dbtarget"
)

var ErrWriteGroupReadinessUnavailable = errors.New("write group readiness unavailable")

// WriteGroupTableInspector reads persisted destination table metadata without
// creating or altering target objects.
type WriteGroupTableInspector interface {
	InspectTable(ctx context.Context, connectorID, schema, table string) (*dbtarget.TableInspection, error)
}

// WriteGroupReadiness is the read-only configuration and schema gate for one
// canonical write group.
type WriteGroupReadiness struct {
	runtimeLayout     *WriteGroupRuntimeLayout
	WorkspaceID       string           `json:"workspace_id"`
	WorkspaceRevision string           `json:"workspace_revision"`
	GroupID           string           `json:"group_id"`
	GroupRevision     string           `json:"group_revision"`
	AppliedRevision   string           `json:"applied_revision"`
	ConfigReady       bool             `json:"config_ready"`
	SchemaReady       bool             `json:"schema_ready"`
	Ready             bool             `json:"ready"`
	SchemaDigest      string           `json:"schema_digest,omitempty"`
	Issues            []ReadinessIssue `json:"issues"`
}

// WithTableInspector wires an optional read-only destination metadata reader.
func (s *WriteGroupService) WithTableInspector(inspector WriteGroupTableInspector) *WriteGroupService {
	if s != nil {
		s.tableInspector = inspector
	}
	return s
}

// Readiness evaluates one canonical group without mutating workspace or
// destination state.
func (s *WriteGroupService) Readiness(ctx context.Context, id string) (*WriteGroupReadiness, error) {
	return s.readiness(ctx, id)
}

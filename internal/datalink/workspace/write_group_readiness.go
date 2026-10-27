package workspace

import (
	"context"
	"errors"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
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
	runtimeLayout   *WriteGroupRuntimeLayout
	counterBindings CounterBindingProofs
	// longRuntimeGate marks a long-layout group whose storage mode or profile the
	// runtime does not write; it adds a blocking issue without changing config or
	// schema state. counterGate marks a supported cumulative group that still has
	// a source without a saved measurement to confirm its counter semantics.
	longRuntimeGate bool
	counterGate     bool
	// reductionGate marks a group whose saved sample reduction the runtime
	// cannot serve yet: the profile stays savable, but no writer may start.
	reductionGate bool
	// counterIssues are the saved-measurement refusals of a cumulative group,
	// read from the database after the group was loaded.
	counterIssues        []ReadinessIssue
	WorkspaceID          string                                    `json:"workspace_id"`
	WorkspaceRevision    string                                    `json:"workspace_revision"`
	GroupID              string                                    `json:"group_id"`
	GroupRevision        string                                    `json:"group_revision"`
	AppliedRevision      string                                    `json:"applied_revision"`
	ConfigReady          bool                                      `json:"config_ready"`
	SchemaReady          bool                                      `json:"schema_ready"`
	Ready                bool                                      `json:"ready"`
	SchemaDigest         string                                    `json:"schema_digest,omitempty"`
	PartitionMaintenance *WriteGroupPartitionMaintenanceReadiness  `json:"partition_maintenance,omitempty"`
	SQLServerCurrent     *recordingplan.SQLServerPartitionEvidence `json:"sqlserver_current,omitempty"`
	Issues               []ReadinessIssue                          `json:"issues"`
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

func newWriteGroupReadiness(record *Record, group *WriteGroup) *WriteGroupReadiness {
	return &WriteGroupReadiness{
		WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision,
		GroupID: group.ID, GroupRevision: group.Revision, AppliedRevision: group.AppliedRevision,
		Issues: []ReadinessIssue{}, longRuntimeGate: group.isLongLayout() && !longRuntimeSupported(group),
		counterGate:   (longRuntimeSupported(group) || group.isWideProfile()) && cumulativeWithoutMeasurements(group),
		reductionGate: recordingReductionUnsupported(group, nil),
	}
}

func (s *WriteGroupService) newSnapshotReadiness(ctx context.Context, snapshot *writeGroupReadinessSnapshot) *WriteGroupReadiness {
	result := newWriteGroupReadiness(snapshot.record, snapshot.group)
	result.reductionGate = recordingReductionUnsupported(snapshot.group, snapshot.tagTypes)
	result.counterIssues = snapshot.counterIssues
	result.counterBindings = snapshot.counterBindings
	s.populateSQLServerPartitionMaintenanceReadiness(ctx, snapshot.connectorKind, snapshot.group, result)
	return result
}

func finalizeWriteGroupReadiness(readiness *WriteGroupReadiness) *WriteGroupReadiness {
	// Schema confirmation and runtime writability are separate facts: the gate is
	// added after ConfigReady and SchemaReady were decided, and only blocks Ready.
	if readiness.longRuntimeGate {
		readiness.Issues = append(readiness.Issues, longRuntimeUnavailableIssue(readiness.GroupID))
	}
	if readiness.counterGate {
		readiness.Issues = append(readiness.Issues, counterSemanticsRequiredIssue(readiness.GroupID))
	}
	if readiness.reductionGate {
		readiness.Issues = append(readiness.Issues, recordingReductionUnavailableIssue(readiness.GroupID))
	}
	readiness.Issues = append(readiness.Issues, readiness.counterIssues...)
	readiness.Ready = readiness.ConfigReady && readiness.SchemaReady && len(readiness.Issues) == 0
	return readiness
}

func writeGroupReadinessResult(readiness *WriteGroupReadiness) (*WriteGroupReadiness, error) {
	return readiness, nil
}

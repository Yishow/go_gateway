package grouptestwrite

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// Results of a test write. WriteOutcome and CleanupStatus are independent.
const (
	WriteVerified   = "written_verified"
	WriteUnverified = "written_unverified"
	WriteFailed     = "failed"
	WriteUnknown    = "unknown"

	CleanupNotAttempted = "not_attempted"
	CleanupCleaned      = "cleaned"
	CleanupFailed       = "failed"
	CleanupUnknown      = "unknown"
)

// Safe reason codes recorded with a result.
const (
	reasonCommitAmbiguous    = "commit-ambiguous"
	reasonRowRejected        = "destination-rejected-row"
	reasonDestinationDenied  = "destination-denied"
	reasonDestinationDown    = "destination-unavailable"
	reasonDestinationBlock   = "destination-blocked"
	reasonReadbackDenied     = "readback-denied"
	reasonReadbackFailed     = "readback-failed"
	reasonReadbackRows       = "readback-row-count"
	reasonReadbackMismatch   = "readback-mismatch"
	reasonReceiptMismatch    = "receipt-mismatch"
	reasonGroupChanged       = "group-changed-after-write"
	reasonGroupChangedBefore = "group-changed-before-write"
	reasonWriteUnconfirmed   = "write-unconfirmed"
	reasonCleanupDenied      = "cleanup-denied"
	reasonCleanupFailed      = "cleanup-failed"
	reasonCleanupAmbiguous   = "cleanup-commit-ambiguous"
	reasonCleanupRemaining   = "cleanup-rows-remain"
)

// ErrResultUnacknowledged means the target step finished but its result could
// not be recorded; the operation stays active until a retry reconciles it.
var ErrResultUnacknowledged = errors.New("test write result could not be recorded")

// Groups reads the saved canonical groups.
type Groups interface {
	Get(ctx context.Context, id string) (*workspace.WriteGroupSaveResult, error)
}

// Tags resolves persisted tag data types.
type Tags interface {
	GetByID(ctx context.Context, id string) (*schema.Tag, error)
}

// Destinations reads and opens saved connectors.
type Destinations interface {
	GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error)
	OpenDestination(ctx context.Context, connectorID, expectedRevision string) (*dbtarget.OpenedDestination, error)
}

// Inspector reads real table metadata without changing it.
type Inspector interface {
	InspectTable(ctx context.Context, connectorID, schemaName, tableName string) (*dbtarget.TableInspection, error)
}

// Dependencies wires the service to production services.
type Dependencies struct {
	Groups       Groups
	Tags         Tags
	Destinations Destinations
	Inspector    Inspector
	Ledger       *recordingplan.Service
}

// Config bounds each destination step; the zero value takes the default.
type Config struct {
	StepTimeout time.Duration
	Now         func() time.Time
}

const defaultStepTimeout = 20 * time.Second

// Service previews and confirms explicit test writes.
type Service struct {
	deps Dependencies
	cfg  Config
}

// New builds a Service.
func New(deps Dependencies, cfg Config) *Service {
	if cfg.StepTimeout <= 0 {
		cfg.StepTimeout = defaultStepTimeout
	}
	if cfg.Now == nil {
		cfg.Now = func() time.Time { return time.Now().UTC() }
	}
	return &Service{deps: deps, cfg: cfg}
}

// Workspace is the caller's server-resolved workspace identity.
type Workspace struct {
	ID       string
	Revision string
}

// PreviewValue is one value the confirmation will write.
type PreviewValue struct {
	Column string `json:"column"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}

// PreviewTarget names the saved destination of the test row.
type PreviewTarget struct {
	ConnectorID string `json:"connector_id"`
	Dialect     string `json:"dialect"`
	Database    string `json:"database"`
	Schema      string `json:"schema"`
	Table       string `json:"table"`
}

// Preview is the mutation-free description of a test write.
type Preview struct {
	Token         string         `json:"token"`
	OperationID   string         `json:"operation_id"`
	Action        string         `json:"action"`
	GroupID       string         `json:"group_id"`
	GroupRevision string         `json:"group_revision"`
	ExpiresAt     time.Time      `json:"expires_at"`
	Target        PreviewTarget  `json:"target"`
	OwnerColumn   string         `json:"owner_column"`
	OwnerValue    string         `json:"owner_value"`
	Dedupe        string         `json:"dedupe"`
	Values        []PreviewValue `json:"values"`
	// Cleanup describes the only rows the operation removes afterwards.
	Cleanup string `json:"cleanup"`
}

// Confirmation is a client's request to run a previewed test write.
type Confirmation struct {
	Token       string
	OperationID string
	// GroupID is the group the caller addressed. When set it must be the group
	// the preview was issued for; the token alone never retargets a request.
	GroupID string
}

// ErrGroupMismatch means the confirmation addresses a different group than the
// preview it carries.
var ErrGroupMismatch = errors.New("test write preview belongs to another group")

// Outcome is what a confirmation resolved to; Claim says whether this call ran
// it, found it running, found a retained result or found the scope busy.
type Outcome struct {
	Operation *recordingplan.SchemaOperation
	Claim     recordingplan.ClaimOutcome
}

func (s *Service) scopeFor(ws Workspace, p *plan) recordingplan.TestWriteScope {
	d := p.group.Destination
	return recordingplan.TestWriteScope{
		WorkspaceID: ws.ID, WorkspaceRevision: ws.Revision, GroupID: p.group.ID, GroupRevision: p.group.Revision,
		ConnectorID: d.ConnectorID, ConnectorRevision: d.ConnectorRevision, Dialect: string(p.dialect),
		Database: d.Database, Schema: d.TableSchema, Table: d.TableName,
	}
}

// Preview describes one test write for a saved group and persists its token.
// It reads the group and inspects the target but never writes to it.
func (s *Service) Preview(ctx context.Context, ws Workspace, groupID string) (*Preview, error) {
	p, err := s.buildPlan(ctx, strings.TrimSpace(groupID))
	if err != nil {
		return nil, err
	}
	token, err := s.deps.Ledger.PrepareTestWritePreview(ctx, s.scopeFor(ws, p), p.digest)
	if err != nil {
		return nil, fmt.Errorf("prepare test write preview: %w", err)
	}
	owner := ownerValue(token.OperationID)
	values := make([]PreviewValue, 0, len(p.members)+1)
	for _, m := range p.members {
		values = append(values, PreviewValue{Column: m.column, Type: string(m.kind), Value: fmt.Sprint(m.value.Value())})
	}
	values = append(values, PreviewValue{Column: p.ownerColumn, Type: "text", Value: owner})
	d := p.group.Destination
	return &Preview{
		Token: token.Token, OperationID: token.OperationID, Action: recordingplan.TestWriteAction,
		GroupID: p.group.ID, GroupRevision: p.group.Revision, ExpiresAt: token.ExpiresAt,
		Target:      PreviewTarget{ConnectorID: d.ConnectorID, Dialect: string(p.dialect), Database: d.Database, Schema: d.TableSchema, Table: d.TableName},
		OwnerColumn: p.ownerColumn, OwnerValue: owner, Dedupe: p.strategy, Values: values,
		Cleanup: fmt.Sprintf("remove only rows where %s = %q, and the receipt of this test when one is written", p.ownerColumn, owner),
	}, nil
}

// Confirm runs a previewed test write once. A repeat of a running operation
// reports it as running; a repeat of a finished one returns the retained result
// without touching the target, even after the token expired. An operation
// whose owner vanished is adopted and reconciled from target evidence rather
// than written again.
func (s *Service) Confirm(ctx context.Context, ws Workspace, req Confirmation) (*Outcome, error) {
	ledger := s.deps.Ledger
	op, claim, err := ledger.ResolveTestWriteReplay(ctx, ws.ID, req.Token, req.OperationID)
	if err != nil {
		return nil, err
	}
	if err := s.checkGroup(ctx, ws, req); err != nil {
		return nil, err
	}
	switch claim {
	case recordingplan.ClaimCompleted:
		return &Outcome{Operation: op, Claim: claim}, nil
	case recordingplan.ClaimInProgress:
		adopted, err := ledger.TakeOverTestWrite(ctx, op.OperationID)
		if err != nil {
			return nil, fmt.Errorf("adopt test write: %w", err)
		}
		if adopted == nil {
			return &Outcome{Operation: op, Claim: claim}, nil
		}
		return s.reconcile(ctx, ws, req, adopted)
	}
	return s.confirmNew(ctx, ws, req)
}

func (s *Service) confirmNew(ctx context.Context, ws Workspace, req Confirmation) (*Outcome, error) {
	token, err := s.deps.Ledger.PreviewTokenForWorkspace(ctx, ws.ID, req.Token)
	if err != nil {
		return nil, err
	}
	p, err := s.buildPlan(ctx, token.PlanID)
	if err != nil {
		return nil, err
	}
	validated, err := s.deps.Ledger.ValidateTestWriteToken(ctx, ws.ID, req.Token, s.scopeFor(ws, p), p.digest)
	if err != nil {
		return nil, err
	}
	op, claim, err := s.deps.Ledger.ClaimTestWrite(ctx, validated)
	if err != nil || claim != recordingplan.ClaimAcquired {
		return &Outcome{Operation: op, Claim: claim}, err
	}
	return s.run(ctx, validated, op, p, nil)
}

// reconcile finishes an adopted operation from its recorded progress and what
// the target holds. It never previews again and never re-inserts a row whose
// effect may already exist.
func (s *Service) reconcile(ctx context.Context, ws Workspace, req Confirmation, adopted *recordingplan.SchemaOperation) (*Outcome, error) {
	token, err := s.deps.Ledger.PreviewTokenForWorkspace(ctx, ws.ID, req.Token)
	if err != nil {
		return nil, err
	}
	resume := decodeProgress(adopted.Detail)
	p, err := s.buildPlan(ctx, token.PlanID)
	// A group edited since the preview can no longer reproduce the expected
	// row; reconciliation then relies on recorded progress and the owner marker.
	if err != nil || p.group.Revision != token.PlanRevision || p.digest != recordingplan.TestWriteContentDigest(token) {
		p = nil
	}
	return s.run(ctx, token, adopted, p, &resume)
}

// checkGroup refuses a confirmation whose addressed group is not the preview's.
func (s *Service) checkGroup(ctx context.Context, ws Workspace, req Confirmation) error {
	if strings.TrimSpace(req.GroupID) == "" {
		return nil
	}
	token, err := s.deps.Ledger.PreviewTokenForWorkspace(ctx, ws.ID, req.Token)
	if err != nil {
		return err
	}
	if token.PlanID != strings.TrimSpace(req.GroupID) {
		return ErrGroupMismatch
	}
	return nil
}

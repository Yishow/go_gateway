package grouptestwrite

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
)

// Steps recorded before their side effect so a restart knows what may have
// happened.
const (
	phaseClaimed  = "claimed"
	phaseWriting  = "writing"
	phaseCleaning = "cleaning"
)

// progress is the service-owned detail saved in the operation ledger. It holds
// only identifiers and safe codes.
type progress struct {
	Phase        string `json:"phase"`
	Schema       string `json:"schema"`
	OwnerColumn  string `json:"owner_column"`
	OwnerValue   string `json:"owner_value"`
	Table        string `json:"table"`
	Strategy     string `json:"strategy"`
	EffectKey    string `json:"effect_key"`
	WriteOutcome string `json:"write_outcome,omitempty"`
	WriteReason  string `json:"write_reason,omitempty"`
}

func decodeProgress(detail string) progress {
	var p progress
	if strings.TrimSpace(detail) == "" {
		return p
	}
	if err := json.Unmarshal([]byte(detail), &p); err != nil {
		// Unreadable progress may hide a started write, so it is treated as the
		// writing phase: the target decides, and a write is never repeated blindly.
		return progress{Phase: phaseWriting}
	}
	return p
}

func (p progress) encode() string {
	data, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(data)
}

// stepContext bounds one destination step independently of the request, so a
// disconnected client cannot leave a half-finished test row behind.
func (s *Service) stepContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), s.cfg.StepTimeout)
}

// run drives write, readback and cleanup for an operation this call owns. With
// resume set it continues an adopted operation from its recorded phase.
func (s *Service) run(ctx context.Context, token *recordingplan.SchemaPreviewToken, op *recordingplan.SchemaOperation, p *plan, resume *progress) (*Outcome, error) {
	state := progress{
		Phase: phaseClaimed, OwnerValue: ownerValue(op.OperationID), Schema: token.Schema, Table: token.TablePrefix,
	}
	if p != nil {
		state.OwnerColumn, state.Strategy = p.ownerColumn, p.strategy
	}
	if resume != nil && resume.Phase != "" {
		state = *resume
	}

	result := recordingplan.SchemaOperationResult{PayloadDigest: recordingplan.TestWriteContentDigest(token)}
	dest, err := s.openDestination(ctx, token)
	if err != nil {
		// Nothing was sent yet only when the write phase was never reached.
		return s.finishUnreachable(ctx, op, state, result, err)
	}
	if dest.Close != nil {
		defer func() { _ = dest.Close() }() //nolint:errcheck // releasing a finished connection; nothing to recover
	}

	if state.Phase != phaseCleaning {
		written, finished, err := s.writePhase(ctx, op, p, token.CreatedAt, dest, &state)
		if err != nil {
			return nil, err
		}
		if finished != nil {
			return s.finish(ctx, op, *finished)
		}
		state.WriteOutcome, state.WriteReason = written.outcome, written.reason
	}
	cleanup := s.cleanupPhase(ctx, op, p, dest, &state)
	result.WriteOutcome, result.Reason = state.WriteOutcome, state.WriteReason
	result.CleanupStatus, result.CleanupReason = cleanup.status, cleanup.reason
	result.Status = operationStatus(state.WriteOutcome)
	result.Detail = state.encode()
	return s.finish(ctx, op, result)
}

type stepResult struct {
	outcome string
	reason  string
}

type cleanupResult struct {
	status string
	reason string
}

func operationStatus(writeOutcome string) recordingplan.SchemaOperationStatus {
	switch writeOutcome {
	case WriteVerified:
		return recordingplan.SchemaOperationSucceeded
	case WriteUnverified:
		return recordingplan.SchemaOperationPartial
	case WriteFailed:
		return recordingplan.SchemaOperationFailed
	default:
		return recordingplan.SchemaOperationUnknown
	}
}

func (s *Service) openDestination(ctx context.Context, token *recordingplan.SchemaPreviewToken) (*dbtarget.OpenedDestination, error) {
	stepCtx, cancel := s.stepContext(ctx)
	defer cancel()
	return s.deps.Destinations.OpenDestination(stepCtx, token.ConnectorID, token.ConnectorRevision)
}

// finishUnreachable records the result when the destination could not be
// opened. Before the write phase nothing was sent, so that is a clean failure;
// at or after it the state of the target is unknown.
func (s *Service) finishUnreachable(ctx context.Context, op *recordingplan.SchemaOperation, state progress, result recordingplan.SchemaOperationResult, openErr error) (*Outcome, error) {
	reason := reasonDestinationDown
	if errors.Is(openErr, dbtarget.ErrDestinationBlocked) {
		reason = reasonDestinationBlock
	}
	result.Reason, result.Detail = reason, state.encode()
	switch state.Phase {
	case phaseClaimed:
		result.Status, result.WriteOutcome, result.CleanupStatus = recordingplan.SchemaOperationFailed, WriteFailed, CleanupNotAttempted
	case phaseCleaning:
		result.Status, result.WriteOutcome = operationStatus(state.WriteOutcome), state.WriteOutcome
		result.Reason = state.WriteReason
		result.CleanupStatus, result.CleanupReason = CleanupUnknown, reason
	default:
		result.Status, result.WriteOutcome, result.CleanupStatus = recordingplan.SchemaOperationUnknown, WriteUnknown, CleanupUnknown
		result.CleanupReason = reason
	}
	return s.finish(ctx, op, result)
}

func (s *Service) finish(ctx context.Context, op *recordingplan.SchemaOperation, result recordingplan.SchemaOperationResult) (*Outcome, error) {
	finishCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.cfg.StepTimeout)
	defer cancel()
	finished, err := s.deps.Ledger.FinishSchemaOperation(finishCtx, op.OperationID, op.Owner, result)
	if err != nil {
		return &Outcome{Operation: op, Claim: recordingplan.ClaimAcquired}, errors.Join(ErrResultUnacknowledged, err)
	}
	return &Outcome{Operation: finished, Claim: recordingplan.ClaimAcquired}, nil
}

func (s *Service) saveProgress(ctx context.Context, op *recordingplan.SchemaOperation, state progress) error {
	saveCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.cfg.StepTimeout)
	defer cancel()
	return s.deps.Ledger.SaveTestWriteProgress(saveCtx, op.OperationID, op.Owner, state.encode())
}

// writePhase writes the operation-owned row and verifies it by readback. It
// returns the write outcome, or a complete terminal result when the operation
// must end here (nothing was written, or the target state cannot be known).
func (s *Service) writePhase(ctx context.Context, op *recordingplan.SchemaOperation, p *plan, at time.Time, dest *dbtarget.OpenedDestination, state *progress) (stepResult, *recordingplan.SchemaOperationResult, error) {
	owner := state.OwnerValue
	if p == nil {
		// Degraded reconciliation: the group changed after the preview, so the
		// expected row cannot be reproduced; report what the target holds.
		return s.degradedWrite(ctx, op, dest, state)
	}
	row, err := p.row(owner, at)
	if err != nil {
		return stepResult{}, nil, err
	}
	state.EffectKey = row.EffectKey
	ref := p.ownedRef(owner)

	if state.Phase == phaseWriting {
		// A restart in the writing phase: the row may already be there.
		if existing, readErr := s.readOwned(ctx, dest, p, ref); readErr == nil && len(existing) > 0 {
			return s.verify(ctx, dest, p, row, ref), nil, nil
		}
		if p.strategy != dbtarget.GroupEffectReceipt {
			// Without destination dedupe a repeat could duplicate a committed row.
			return stepResult{}, &recordingplan.SchemaOperationResult{
				Status: recordingplan.SchemaOperationUnknown, WriteOutcome: WriteUnknown, Reason: reasonWriteUnconfirmed,
				CleanupStatus: CleanupNotAttempted, PayloadDigest: op.PayloadDigest, Detail: state.encode(),
			}, nil
		}
	}
	state.Phase = phaseWriting
	if err := s.saveProgress(ctx, op, *state); err != nil {
		return stepResult{}, nil, errors.Join(ErrResultUnacknowledged, err)
	}
	digest, err := groupdelivery.RowPayloadDigest(row)
	if err != nil {
		return stepResult{}, nil, unsupported(ReasonSQLValueBlocked)
	}
	stepCtx, cancel := s.stepContext(ctx)
	_, insertErr := dbtarget.InsertGroupRow(stepCtx, dest.DB, dbtarget.GroupInsertRequest{
		Kind: dest.Kind, SchemaName: p.group.Destination.TableSchema, TableName: p.group.Destination.TableName,
		Row: row, PayloadDigest: digest, Strategy: p.strategy, CommittedAt: s.cfg.Now().Format(time.RFC3339Nano),
	})
	cancel()
	if insertErr != nil {
		return s.afterFailedInsert(ctx, dest, p, row, ref, op, state, insertErr)
	}
	return s.verify(ctx, dest, p, row, ref), nil, nil
}

// afterFailedInsert decides the outcome of a failed insert: a pre-commit
// failure wrote nothing; a commit failure is checked against the target and is
// unknown when that cannot settle it.
func (s *Service) afterFailedInsert(ctx context.Context, dest *dbtarget.OpenedDestination, p *plan, row dbtarget.EncodedRow, ref dbtarget.OwnedRowRef, op *recordingplan.SchemaOperation, state *progress, insertErr error) (stepResult, *recordingplan.SchemaOperationResult, error) {
	var insertFailure *dbtarget.GroupInsertError
	ambiguous := errors.As(insertErr, &insertFailure) && insertFailure.Ambiguous()
	if ambiguous {
		if rows, err := s.readOwned(ctx, dest, p, ref); err == nil && len(rows) > 0 {
			return s.verify(ctx, dest, p, row, ref), nil, nil
		}
		return stepResult{}, &recordingplan.SchemaOperationResult{
			Status: recordingplan.SchemaOperationUnknown, WriteOutcome: WriteUnknown, Reason: reasonCommitAmbiguous,
			CleanupStatus: CleanupNotAttempted, PayloadDigest: op.PayloadDigest, Detail: state.encode(),
		}, nil
	}
	reason := reasonDestinationDown
	switch dbtarget.ClassifyInsertError(insertErr) {
	case dbtarget.InsertErrorRow:
		reason = reasonRowRejected
	case dbtarget.InsertErrorTarget:
		reason = reasonDestinationDenied
	}
	return stepResult{}, &recordingplan.SchemaOperationResult{
		Status: recordingplan.SchemaOperationFailed, WriteOutcome: WriteFailed, Reason: reason,
		CleanupStatus: CleanupNotAttempted, PayloadDigest: op.PayloadDigest, Detail: state.encode(),
	}, nil
}

func (s *Service) readOwned(ctx context.Context, dest *dbtarget.OpenedDestination, p *plan, ref dbtarget.OwnedRowRef) ([][]any, error) {
	stepCtx, cancel := s.stepContext(ctx)
	defer cancel()
	return dbtarget.ReadOwnedRows(stepCtx, dest.DB, ref, p.readColumns())
}

// readColumns lists the columns a readback compares, member values first.
func (p *plan) readColumns() []string {
	columns := make([]string, 0, len(p.members)+2)
	for _, m := range p.members {
		columns = append(columns, m.column)
	}
	columns = append(columns, p.ownerColumn)
	if column := strings.TrimSpace(p.group.RowPolicy.ProvenanceColumn); column != "" {
		columns = append(columns, column)
	}
	return columns
}

// verify reads the committed row back and compares typed values, scope
// identity and provenance. Anything short of a full match is unverified, with a
// safe reason; it is never reported as verified.
func (s *Service) verify(ctx context.Context, dest *dbtarget.OpenedDestination, p *plan, row dbtarget.EncodedRow, ref dbtarget.OwnedRowRef) stepResult {
	rows, err := s.readOwned(ctx, dest, p, ref)
	if err != nil {
		if dbtarget.ClassifyInsertError(err) == dbtarget.InsertErrorTarget {
			return stepResult{outcome: WriteUnverified, reason: reasonReadbackDenied}
		}
		return stepResult{outcome: WriteUnverified, reason: reasonReadbackFailed}
	}
	if len(rows) != 1 {
		return stepResult{outcome: WriteUnverified, reason: reasonReadbackRows}
	}
	if !s.rowMatches(p, row, rows[0]) {
		return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
	}
	if p.strategy == dbtarget.GroupEffectReceipt {
		digest, err := groupdelivery.RowPayloadDigest(row)
		if err != nil {
			return stepResult{outcome: WriteUnverified, reason: reasonReceiptMismatch}
		}
		stepCtx, cancel := s.stepContext(ctx)
		defer cancel()
		stored, found, err := dbtarget.ReadEffectReceiptDigest(stepCtx, dest.DB, dest.Kind, p.group.Destination.TableSchema, row.EffectKey)
		if err != nil {
			return stepResult{outcome: WriteUnverified, reason: reasonReadbackFailed}
		}
		if !found || stored != digest {
			return stepResult{outcome: WriteUnverified, reason: reasonReceiptMismatch}
		}
	}
	return stepResult{outcome: WriteVerified}
}

// rowMatches compares the row read back with what was written: typed member
// values by exact comparison, the ownership marker, and provenance content.
func (s *Service) rowMatches(p *plan, row dbtarget.EncodedRow, got []any) bool {
	columnType := make(map[string]string, len(p.columns))
	for _, column := range p.columns {
		columnType[strings.ToLower(column.Name)] = column.DataType
	}
	for i, m := range p.members {
		decoded, err := dbtarget.DecodeExactValue(p.dialect, columnType[strings.ToLower(m.column)], m.kind, got[i])
		if err != nil || !decoded.Equal(m.value) {
			return false
		}
	}
	owner, ok := textValue(got[len(p.members)])
	if !ok || owner != row.EntityKey {
		return false
	}
	if provenance := strings.TrimSpace(p.group.RowPolicy.ProvenanceColumn); provenance != "" {
		return provenanceMatches(row, provenance, got[len(p.members)+1])
	}
	return true
}

func textValue(raw any) (string, bool) {
	switch v := raw.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	}
	return "", false
}

func provenanceMatches(row dbtarget.EncodedRow, column string, got any) bool {
	var want any
	for _, cell := range row.Cells {
		if strings.EqualFold(cell.Column, column) {
			want = cell.Value
		}
	}
	wantText, ok := textValue(want)
	if !ok {
		return false
	}
	gotText, ok := textValue(got)
	if !ok {
		return false
	}
	var wantJSON, gotJSON any
	if json.Unmarshal([]byte(wantText), &wantJSON) != nil || json.Unmarshal([]byte(gotText), &gotJSON) != nil {
		return false
	}
	return reflect.DeepEqual(wantJSON, gotJSON)
}

// degradedWrite reconciles an adopted operation whose group changed since the
// preview. The expected row cannot be reproduced, so an owned row that exists
// is reported as written but unverified.
func (s *Service) degradedWrite(ctx context.Context, op *recordingplan.SchemaOperation, dest *dbtarget.OpenedDestination, state *progress) (stepResult, *recordingplan.SchemaOperationResult, error) {
	ref := dbtarget.OwnedRowRef{
		Kind: dest.Kind, SchemaName: state.Schema, TableName: state.Table, OwnerColumn: state.OwnerColumn, OwnerValue: state.OwnerValue,
	}
	if state.Phase == phaseClaimed {
		// The write never started, so nothing needs reconciling.
		return stepResult{}, &recordingplan.SchemaOperationResult{
			Status: recordingplan.SchemaOperationFailed, WriteOutcome: WriteFailed, Reason: reasonGroupChangedBefore,
			CleanupStatus: CleanupNotAttempted, PayloadDigest: op.PayloadDigest, Detail: state.encode(),
		}, nil
	}
	if state.OwnerColumn == "" {
		return stepResult{}, &recordingplan.SchemaOperationResult{
			Status: recordingplan.SchemaOperationUnknown, WriteOutcome: WriteUnknown, Reason: reasonWriteUnconfirmed,
			CleanupStatus: CleanupUnknown, CleanupReason: reasonGroupChanged, PayloadDigest: op.PayloadDigest, Detail: state.encode(),
		}, nil
	}
	stepCtx, cancel := s.stepContext(ctx)
	defer cancel()
	rows, err := dbtarget.ReadOwnedRows(stepCtx, dest.DB, ref, []string{state.OwnerColumn})
	if err == nil && len(rows) == 0 {
		return stepResult{}, &recordingplan.SchemaOperationResult{
			Status: recordingplan.SchemaOperationUnknown, WriteOutcome: WriteUnknown, Reason: reasonWriteUnconfirmed,
			CleanupStatus: CleanupNotAttempted, PayloadDigest: op.PayloadDigest, Detail: state.encode(),
		}, nil
	}
	return stepResult{outcome: WriteUnverified, reason: reasonGroupChanged}, nil, nil
}

// cleanupPhase removes only the rows this operation owns. The intent and the
// write outcome are saved first so a restart finishes the cleanup instead of
// writing again.
func (s *Service) cleanupPhase(ctx context.Context, op *recordingplan.SchemaOperation, p *plan, dest *dbtarget.OpenedDestination, state *progress) cleanupResult {
	state.Phase = phaseCleaning
	if err := s.saveProgress(ctx, op, *state); err != nil {
		return cleanupResult{status: CleanupUnknown, reason: reasonCleanupFailed}
	}
	ref := dbtarget.OwnedRowRef{
		Kind: dest.Kind, SchemaName: state.Schema, TableName: state.Table, OwnerColumn: state.OwnerColumn, OwnerValue: state.OwnerValue,
	}
	if p != nil {
		ref = p.ownedRef(state.OwnerValue)
	}
	effectKey := ""
	if state.Strategy == dbtarget.GroupEffectReceipt {
		effectKey = state.EffectKey
	}
	stepCtx, cancel := s.stepContext(ctx)
	_, err := dbtarget.RemoveOwnedRows(stepCtx, dest.DB, ref, effectKey)
	cancel()
	switch {
	case err == nil:
	case errors.Is(err, dbtarget.ErrOwnedCleanupAmbiguous):
		return cleanupResult{status: CleanupUnknown, reason: reasonCleanupAmbiguous}
	case dbtarget.ClassifyInsertError(err) == dbtarget.InsertErrorTarget:
		return cleanupResult{status: CleanupFailed, reason: reasonCleanupDenied}
	default:
		return cleanupResult{status: CleanupFailed, reason: reasonCleanupFailed}
	}
	checkCtx, cancelCheck := s.stepContext(ctx)
	defer cancelCheck()
	if remaining, readErr := dbtarget.ReadOwnedRows(checkCtx, dest.DB, ref, []string{state.OwnerColumn}); readErr == nil && len(remaining) > 0 {
		return cleanupResult{status: CleanupFailed, reason: reasonCleanupRemaining}
	}
	return cleanupResult{status: CleanupCleaned}
}

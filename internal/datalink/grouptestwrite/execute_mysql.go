package grouptestwrite

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

const (
	mysqlTestCommittedAtLayout = "2006-01-02T15:04:05.000000Z"
	mysqlTestRecordColumn      = "record_id"
)

func mysqlLongTestRowRef(kind schema.DatabaseConnectorKind, state progress) dbtarget.GroupLongTestRowRef {
	return dbtarget.GroupLongTestRowRef{
		Kind: kind, SchemaName: state.Schema, TableName: state.Table,
		OwnerColumn: state.OwnerColumn, OwnerValue: state.TableOwnerValue,
		ProvenanceColumn: state.ProvenanceColumn,
		WorkspaceID:      state.WorkspaceID, GroupID: state.GroupID, GroupRevision: state.GroupRevision,
		OperationOwner: state.OwnerValue,
		RecordIDs:      state.RecordIDs, SampleIDs: state.SampleIDs,
		EffectKeys: state.EffectKeys, EffectDigests: state.EffectDigests,
	}
}

func (s *Service) writeMySQLPhase(ctx context.Context, op *recordingplan.SchemaOperation, p *plan, at time.Time, dest *dbtarget.OpenedDestination, state *progress) (stepResult, *recordingplan.SchemaOperationResult, error) {
	rows, err := p.mysqlRows(state.OwnerValue, at)
	if err != nil {
		return stepResult{}, nil, err
	}
	state.Schema, state.Table = p.group.Destination.TableSchema, p.group.Destination.TableName
	state.OwnerColumn, state.TableOwnerValue = p.ownerColumn, "gateway"
	state.ProvenanceColumn = p.group.RowPolicy.ProvenanceColumn
	state.WorkspaceID, state.GroupID, state.GroupRevision = p.group.WorkspaceID, p.group.ID, p.group.Revision
	state.SampleIDs = make([]string, len(rows))
	state.RecordIDs = make([]string, len(rows))
	state.EffectKeys = make([]string, len(rows))
	state.EffectDigests = make(map[string]string, len(rows))
	for i, row := range rows {
		digest, err := groupdelivery.RowPayloadDigest(row)
		if err != nil {
			return stepResult{}, nil, unsupported(ReasonSQLValueBlocked)
		}
		state.RecordIDs[i] = row.RecordID
		state.SampleIDs[i] = fmt.Sprintf("%s-%d", state.OwnerValue, i)
		state.EffectKeys[i] = row.EffectKey
		state.EffectDigests[row.EffectKey] = digest
	}
	state.EffectKey = state.EffectKeys[0]
	if state.Phase == phaseWriting {
		ref := mysqlLongTestRowRef(dest.Kind, *state)
		existing, readErr := s.readLongTestRows(ctx, dest, ref)
		if readErr != nil {
			// A row we cannot prove is ours may already be committed. Do not
			// insert again; let bounded exact-ID cleanup and reconciliation
			// decide the outcome.
			return stepResult{outcome: WriteUnknown, reason: reasonReadbackFailed}, nil, nil //nolint:nilerr // Persist uncertain target evidence as an explicit operation outcome.
		}
		if len(existing) > 0 {
			return s.verifyMySQL(ctx, dest, p, *state, rows), nil, nil
		}
	}
	state.Phase = phaseWriting
	// IDs, effect keys and digests are durable before any row is written.
	if err := s.saveProgress(ctx, op, *state); err != nil {
		return stepResult{}, nil, errors.Join(ErrResultUnacknowledged, err)
	}
	for _, row := range rows {
		stepCtx, cancel := s.stepContext(ctx)
		_, insertErr := s.insertGroupRow(stepCtx, dest.DB, dbtarget.GroupInsertRequest{
			Kind: dest.Kind, SchemaName: p.group.Destination.TableSchema, TableName: p.group.Destination.TableName,
			Row: row, PayloadDigest: state.EffectDigests[row.EffectKey], Strategy: dbtarget.GroupEffectReceipt,
			CommittedAt: s.cfg.Now().UTC().Format(mysqlTestCommittedAtLayout), MySQL: p.mysql.proof,
		})
		cancel()
		if insertErr != nil {
			return s.afterFailedMySQLInsert(ctx, dest, p, rows, op, state, insertErr)
		}
	}
	return s.verifyMySQL(ctx, dest, p, *state, rows), nil, nil
}

func (s *Service) insertGroupRow(ctx context.Context, db *sql.DB, req dbtarget.GroupInsertRequest) (dbtarget.GroupInsertResult, error) {
	if s.insertRow != nil {
		return s.insertRow(ctx, db, req)
	}
	return dbtarget.InsertGroupRow(ctx, db, req)
}

func (s *Service) readLongTestRows(ctx context.Context, dest *dbtarget.OpenedDestination, ref dbtarget.GroupLongTestRowRef) ([][]any, error) {
	stepCtx, cancel := s.stepContext(ctx)
	defer cancel()
	return dbtarget.ReadGroupLongTestRows(stepCtx, dest.DB, ref, []string{mysqlTestRecordColumn})
}

func (s *Service) afterFailedMySQLInsert(ctx context.Context, dest *dbtarget.OpenedDestination, p *plan, rows []dbtarget.EncodedRow, op *recordingplan.SchemaOperation, state *progress, insertErr error) (stepResult, *recordingplan.SchemaOperationResult, error) {
	ref := mysqlLongTestRowRef(dest.Kind, *state)
	existing, readErr := s.readLongTestRows(ctx, dest, ref)
	if readErr == nil && len(existing) > 0 {
		result := s.verifyMySQLRows(ctx, dest, p, *state, rows, false)
		if result.outcome == WriteVerified || result.outcome == WriteUnverified {
			// A known committed subset is a partial write, never a clean failure.
			// run still executes bounded cleanup for exactly these record IDs.
			return result, nil, nil
		}
	}
	var insertFailure *dbtarget.GroupInsertError
	ambiguous := readErr != nil || (errors.As(insertErr, &insertFailure) && insertFailure.Ambiguous())
	reason := reasonDestinationDown
	switch dbtarget.ClassifyInsertError(insertErr) {
	case dbtarget.InsertErrorRow:
		reason = reasonRowRejected
	case dbtarget.InsertErrorTarget:
		reason = reasonDestinationDenied
	}
	if readErr == nil && len(existing) == 0 && !ambiguous {
		// A non-ambiguous pre-commit failure with a confirmed empty target wrote
		// nothing. This is a terminal failed/not-attempted result, not a cleanup.
		return stepResult{}, &recordingplan.SchemaOperationResult{
			Status: recordingplan.SchemaOperationFailed, WriteOutcome: WriteFailed, Reason: reason,
			CleanupStatus: CleanupNotAttempted, PayloadDigest: op.PayloadDigest, Detail: state.encode(),
		}, nil
	}
	if ambiguous {
		return stepResult{outcome: WriteUnknown, reason: reasonCommitAmbiguous}, nil, nil //nolint:nilerr // Persist ambiguous commit and continue the owned cleanup phase.
	}
	return stepResult{outcome: WriteFailed, reason: reason}, nil, nil
}

func (s *Service) verifyMySQL(ctx context.Context, dest *dbtarget.OpenedDestination, p *plan, state progress, expected []dbtarget.EncodedRow) stepResult {
	return s.verifyMySQLRows(ctx, dest, p, state, expected, true)
}

func (s *Service) verifyMySQLRows(ctx context.Context, dest *dbtarget.OpenedDestination, p *plan, state progress, expected []dbtarget.EncodedRow, exactCount bool) stepResult {
	ref := mysqlLongTestRowRef(dest.Kind, state)
	stepCtx, cancel := s.stepContext(ctx)
	columns := append([]string{mysqlTestRecordColumn}, p.readColumns()...)
	rows, err := dbtarget.ReadGroupLongTestRows(stepCtx, dest.DB, ref, columns)
	cancel()
	if err != nil {
		if dbtarget.ClassifyInsertError(err) == dbtarget.InsertErrorTarget {
			return stepResult{outcome: WriteUnverified, reason: reasonReadbackDenied}
		}
		return stepResult{outcome: WriteUnverified, reason: reasonReadbackFailed}
	}
	if len(rows) > len(expected) || (exactCount && len(rows) != len(expected)) {
		return stepResult{outcome: WriteUnverified, reason: reasonReadbackRows}
	}
	columnIndex := make(map[string]int, len(columns))
	columnType := make(map[string]string, len(p.columns))
	for i, name := range columns {
		columnIndex[strings.ToLower(name)] = i
	}
	for _, column := range p.columns {
		columnType[strings.ToLower(column.Name)] = column.DataType
	}
	recordIndex, ok := columnIndex[mysqlTestRecordColumn]
	if !ok {
		return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
	}
	gotByRecord := make(map[string][]any, len(rows))
	for _, row := range rows {
		recordID, ok := textValue(row[recordIndex])
		if !ok || recordID == "" {
			return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
		}
		gotByRecord[recordID] = row
	}
	matched := 0
	for i, want := range expected {
		got, ok := gotByRecord[want.RecordID]
		if !ok {
			if exactCount {
				return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
			}
			continue
		}
		if len(got) != len(columns) {
			return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
		}
		member := p.mysql.members[i]
		slot := mysqlValueSlot(member.kind)
		slotIndex, ok := columnIndex[slot]
		if !ok {
			return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
		}
		if !mysqlTestValueMatches(columnType[slot], member.kind, member.value, got[slotIndex]) {
			return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
		}
		for _, cell := range want.Cells {
			cellIndex, ok := columnIndex[strings.ToLower(cell.Column)]
			if !ok {
				continue
			}
			if strings.EqualFold(cell.Column, slot) {
				continue
			}
			if strings.EqualFold(cell.Column, p.group.RowPolicy.ProvenanceColumn) {
				if !mysqlTestJSONEqual(cell.Value, got[cellIndex]) {
					return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
				}
				continue
			}
			if !mysqlTestCellEqual(cell.Value, got[cellIndex]) {
				return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
			}
		}
		digest := state.EffectDigests[want.EffectKey]
		if digest == "" {
			return stepResult{outcome: WriteUnverified, reason: reasonReceiptMismatch}
		}
		receiptCtx, cancelReceipt := s.stepContext(ctx)
		stored, found, err := dbtarget.ReadEffectReceiptDigest(receiptCtx, dest.DB, dest.Kind, state.Schema, want.EffectKey)
		cancelReceipt()
		if err != nil || !found || stored != digest {
			return stepResult{outcome: WriteUnverified, reason: reasonReceiptMismatch}
		}
		matched++
	}
	if matched != len(rows) {
		return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
	}
	if !exactCount && len(rows) < len(expected) {
		return stepResult{outcome: WriteUnverified, reason: reasonReadbackRows}
	}
	return stepResult{outcome: WriteVerified}
}

func mysqlTestValueMatches(declared string, kind measurement.ExactType, expected measurement.ExactValue, raw any) bool {
	decoded, err := dbtarget.DecodeExactValue(dbtarget.SQLDialectMySQL, declared, kind, raw)
	return err == nil && decoded.Equal(expected)
}

func mysqlTestCellEqual(want, got any) bool {
	if want == nil {
		return got == nil
	}
	if raw, ok := got.([]byte); ok {
		switch expected := want.(type) {
		case string:
			return string(raw) == expected
		case time.Time:
			parsed, err := time.Parse(time.RFC3339Nano, string(raw))
			return err == nil && parsed.Equal(expected)
		}
	}
	if expected, ok := want.(time.Time); ok {
		actual, ok := got.(time.Time)
		return ok && actual.Equal(expected)
	}
	return reflect.DeepEqual(want, got)
}

// mysqlTestJSONEqual tolerates the server's JSON formatting while preserving
// every numeric digit; float64 decoding could hide a changed large integer.
func mysqlTestJSONEqual(want, got any) bool {
	left, leftOK := textValue(want)
	right, rightOK := textValue(got)
	if !leftOK || !rightOK || !json.Valid([]byte(left)) || !json.Valid([]byte(right)) {
		return false
	}
	decode := func(text string) (any, error) {
		decoder := json.NewDecoder(bytes.NewBufferString(text))
		decoder.UseNumber()
		var value any
		err := decoder.Decode(&value)
		return value, err
	}
	a, errA := decode(left)
	b, errB := decode(right)
	return errA == nil && errB == nil && reflect.DeepEqual(a, b)
}

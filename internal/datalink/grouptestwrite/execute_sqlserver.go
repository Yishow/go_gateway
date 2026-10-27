package grouptestwrite

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
)

const sqlServerTestCommittedAtLayout = "2006-01-02T15:04:05.000000Z"

// writeSQLServerPhase writes the operation-owned SQL Server long rows and
// verifies them by exact readback. It keeps the same progress-before-write,
// restart, partial-write and cleanup contracts as the MySQL phase.
func (s *Service) writeSQLServerPhase(ctx context.Context, op *recordingplan.SchemaOperation, p *plan, at time.Time, dest *dbtarget.OpenedDestination, state *progress) (stepResult, *recordingplan.SchemaOperationResult, error) {
	rows, err := p.sqlserverRows(state.OwnerValue, at)
	if err != nil {
		return stepResult{}, nil, err
	}
	state.Schema, state.Table = p.group.Destination.TableSchema, p.group.Destination.TableName
	state.OwnerColumn, state.TableOwnerValue = p.ownerColumn, sqlServerTestTableOwner
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
			return stepResult{outcome: WriteUnknown, reason: reasonReadbackFailed}, nil, nil //nolint:nilerr // persist uncertain target evidence
		}
		if len(existing) > 0 {
			return s.verifySQLServer(ctx, dest, p, *state, rows), nil, nil
		}
	}
	state.Phase = phaseWriting
	if err := s.saveProgress(ctx, op, *state); err != nil {
		return stepResult{}, nil, errors.Join(ErrResultUnacknowledged, err)
	}
	for _, row := range rows {
		stepCtx, cancel := s.stepContext(ctx)
		_, insertErr := s.insertGroupRow(stepCtx, dest.DB, dbtarget.GroupInsertRequest{
			Kind: dest.Kind, SchemaName: p.group.Destination.TableSchema, TableName: p.group.Destination.TableName,
			Row: row, PayloadDigest: state.EffectDigests[row.EffectKey], Strategy: dbtarget.GroupEffectReceipt,
			CommittedAt: s.cfg.Now().UTC().Format(sqlServerTestCommittedAtLayout), SQLServer: p.sqlserver.proof,
		})
		cancel()
		if insertErr != nil {
			return s.afterFailedSQLServerInsert(ctx, dest, p, rows, op, state, insertErr)
		}
	}
	return s.verifySQLServer(ctx, dest, p, *state, rows), nil, nil
}

func (s *Service) afterFailedSQLServerInsert(ctx context.Context, dest *dbtarget.OpenedDestination, p *plan, rows []dbtarget.EncodedRow, op *recordingplan.SchemaOperation, state *progress, insertErr error) (stepResult, *recordingplan.SchemaOperationResult, error) {
	ref := mysqlLongTestRowRef(dest.Kind, *state)
	existing, readErr := s.readLongTestRows(ctx, dest, ref)
	if readErr == nil && len(existing) > 0 {
		result := s.verifySQLServerRows(ctx, dest, p, *state, rows, false)
		if result.outcome == WriteVerified || result.outcome == WriteUnverified {
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
		return stepResult{}, &recordingplan.SchemaOperationResult{
			Status: recordingplan.SchemaOperationFailed, WriteOutcome: WriteFailed, Reason: reason,
			CleanupStatus: CleanupNotAttempted, PayloadDigest: op.PayloadDigest, Detail: state.encode(),
		}, nil
	}
	if ambiguous {
		return stepResult{outcome: WriteUnknown, reason: reasonCommitAmbiguous}, nil, nil
	}
	return stepResult{outcome: WriteFailed, reason: reason}, nil, nil
}

func (s *Service) verifySQLServer(ctx context.Context, dest *dbtarget.OpenedDestination, p *plan, state progress, expected []dbtarget.EncodedRow) stepResult {
	return s.verifySQLServerRows(ctx, dest, p, state, expected, true)
}

func (s *Service) verifySQLServerRows(ctx context.Context, dest *dbtarget.OpenedDestination, p *plan, state progress, expected []dbtarget.EncodedRow, exactCount bool) stepResult {
	if p == nil || p.sqlserver == nil {
		return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
	}
	ref := mysqlLongTestRowRef(dest.Kind, state)
	stepCtx, cancel := s.stepContext(ctx)
	columns := sqlServerReadbackColumns(p)
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
		member := p.sqlserver.members[i]
		slot := sqlServerValueSlot(member.kind)
		slotIndex, ok := columnIndex[slot]
		if !ok {
			return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
		}
		if !sqlServerTestValueMatches(columnType[slot], member.kind, member.value, got[slotIndex]) {
			return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
		}
		for _, cell := range want.Cells {
			cellIndex, ok := columnIndex[strings.ToLower(cell.Column)]
			if !ok {
				return stepResult{outcome: WriteUnverified, reason: reasonReadbackMismatch}
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
			if !sqlServerTestCellEqual(cell.Value, got[cellIndex]) {
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

// sqlServerReadbackColumns selects every canonical row cell, not only the
// member slots, so bucket/time/profile/quality/NULL cells cannot escape
// verification. record_id is requested once.
func sqlServerReadbackColumns(p *plan) []string {
	if p == nil {
		return nil
	}
	columns := make([]string, 0, len(p.columns)+1)
	seen := make(map[string]struct{}, len(p.columns)+1)
	add := func(column string) {
		key := strings.ToLower(column)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		columns = append(columns, column)
	}
	add(mysqlTestRecordColumn)
	for _, column := range p.columns {
		add(column.Name)
	}
	add(p.ownerColumn)
	if provenance := strings.TrimSpace(p.group.RowPolicy.ProvenanceColumn); provenance != "" {
		add(provenance)
	}
	return columns
}

func sqlServerTestValueMatches(declared string, kind measurement.ExactType, expected measurement.ExactValue, raw any) bool {
	decoded, err := dbtarget.DecodeExactValue(dbtarget.SQLDialectSQLServer, declared, kind, raw)
	return err == nil && decoded.Equal(expected)
}

func sqlServerTestCellEqual(want, got any) bool {
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
		if actual, ok := got.(time.Time); ok {
			return actual.Equal(expected)
		}
		if text, ok := textValue(got); ok {
			parsed, err := time.Parse(time.RFC3339Nano, text)
			if err != nil {
				parsed, err = time.Parse("2006-01-02T15:04:05.000000", text)
			}
			return err == nil && parsed.Equal(expected)
		}
	}
	if expected, ok := want.(bool); ok {
		if actual, ok := got.(bool); ok {
			return actual == expected
		}
		if actual, ok := got.(int64); ok {
			return (actual != 0) == expected
		}
	}
	return reflect.DeepEqual(want, got)
}

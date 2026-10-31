package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"go-gateway/internal/datalink/snapshot"
)

type mixedSourceOutcome struct{ id, due, kind, raw, digest string }

func (s *Store) mixedSourceRecords(ctx context.Context, scope SourceScope, ref SourceRef) (SourceRecords, error) {
	var out SourceRecords
	facts, err := s.mixedSourceFacts(ctx, scope, ref)
	if err != nil {
		return out, err
	}
	query := `SELECT record_id,due_at,kind,outcome,outcome_digest FROM wg_delivery_mixed_outcomes WHERE workspace_id=? AND group_id=? AND group_revision=?`
	args := []any{scope.WorkspaceID, scope.GroupID, scope.GroupRevision}
	if !scope.Wide {
		query += ` AND json_extract(outcome,'$.EntityKey')=?`
		args = append(args, SourceEntityKey(ref.DeviceID, ref.PointID, ref.TagID))
	}
	var newest mixedSourceOutcome
	err = s.db.QueryRowContext(ctx, query+` ORDER BY due_at DESC,record_id DESC LIMIT 1`, args...).Scan(&newest.id, &newest.due, &newest.kind, &newest.raw, &newest.digest)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return out, err
	}
	if err == nil {
		record, e := s.verifyMixedSource(ctx, scope, ref, facts, newest)
		if e != nil {
			return out, e
		}
		out.Latest = &record
	}
	entity := SourceEntityKey(ref.DeviceID, ref.PointID, ref.TagID)
	if scope.Wide {
		entity = ""
	}
	committed, err := s.sourceRows(ctx, scope, partitionKey(GroupKey{GroupID: scope.GroupID}, entity), StateCommitted, sourceCandidateLimit)
	if err != nil {
		return out, err
	}
	for _, row := range committed {
		var candidate mixedSourceOutcome
		err = s.db.QueryRowContext(ctx, `SELECT record_id,due_at,kind,outcome,outcome_digest FROM wg_delivery_mixed_outcomes WHERE workspace_id=? AND group_id=? AND group_revision=? AND record_id=?`, scope.WorkspaceID, scope.GroupID, scope.GroupRevision, row.recordID).Scan(&candidate.id, &candidate.due, &candidate.kind, &candidate.raw, &candidate.digest)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return out, err
		}
		record, e := s.verifyMixedSource(ctx, scope, ref, facts, candidate)
		if e != nil {
			return out, e
		}
		if record.Available && record.Delivery == SourceDeliveryCommitted {
			out.LastSQLCommitted = &record
			break
		}
	}
	return out, nil
}
func (s *Store) mixedSourceRow(ctx context.Context, scope SourceScope, id string) (sourceRow, error) {
	var row sourceRow
	err := s.db.QueryRowContext(ctx, `SELECT effect_key,record_id,workspace_id,group_id,group_revision,entity_key,connector_id,connector_revision,bucket_start,payload,payload_digest,state,committed_at,COALESCE(due_at,'') FROM wg_delivery_outbox WHERE workspace_id=? AND group_id=? AND group_revision=? AND record_id=?`, scope.WorkspaceID, scope.GroupID, scope.GroupRevision, id).Scan(&row.effectKey, &row.recordID, &row.workspaceID, &row.groupID, &row.groupRevision, &row.entityKey, &row.connectorID, &row.connectorRevision, &row.bucketStart, &row.payload, &row.digest, &row.state, &row.committedAt, &row.dueAt)
	return row, err
}
func mixedSourceIdentity(scope SourceScope, o snapshot.Outcome) bool {
	var id string
	var err error
	if scope.Wide {
		if !o.BucketStart.IsZero() || o.EntityKey != "" {
			return false
		}
		id, err = snapshot.MixedDispatchRecordID(scope.WorkspaceID, scope.GroupID, scope.GroupRevision, o.DueAt)
	} else {
		if len(o.Members) != 1 || o.EntityKey != o.Members[0].MemberKey {
			return false
		}
		id, err = snapshot.ProfileRecordID(scope.WorkspaceID, scope.GroupID, scope.GroupRevision, o.EntityKey, "mixed", o.BucketStart)
	}
	return err == nil && id == o.RecordID
}

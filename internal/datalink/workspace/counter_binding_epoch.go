package workspace

import (
	"context"
	"database/sql"
	"errors"
	"reflect"

	"go-gateway/internal/datalink/counterproof"
)

// sameAppliedCounterMeaningInTx permits display-only reuse only when original
// immutable evidence proves the same complete meaning. Rekeying is transient:
// no old descriptor, accepted sample or digest is normalized or rewritten.
func (s *WriteGroupService) sameAppliedCounterMeaningInTx(ctx context.Context, tx *sql.Tx, active *WriteGroupAppliedSnapshot, reviewed CounterBindingProofs) (bool, error) {
	if len(reviewed) == 0 {
		return true, nil
	}
	var raw string
	err := tx.QueryRowContext(ctx, s.repo.query(`SELECT payload FROM wg_runtime_versions WHERE workspace_id=$1 AND group_id=$2 AND group_revision=$3`), active.WorkspaceID, active.GroupID, active.AppliedRevision).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return false, &CounterBindingProofError{Reason: ErrCounterBindingProofMissing}
	}
	if err != nil {
		return false, err
	}
	original, err := counterproof.ResolveSnapshot([]byte(raw), counterproof.SnapshotScope{WorkspaceID: active.WorkspaceID, GroupID: active.GroupID, AppliedRevision: active.AppliedRevision})
	if err != nil {
		return false, err
	}
	comparable, err := rekeyCounterBindings(reviewed, active.AppliedRevision)
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(original, comparable), nil
}

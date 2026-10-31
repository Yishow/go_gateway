package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// verifyCounterBindingsInTx compares reviewed facts with the definitions read
// by validation in this same transaction. Neither map outlives its transaction
// as mutable repository state; readiness retains only sealed value proofs.
func (s *WriteGroupService) verifyCounterBindingsInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup, readiness *WriteGroupReadiness, definitions map[string]*measurement.MeasurementDefinition) (CounterBindingProofs, error) {
	var types map[string]schema.DataType
	if readiness.runtimeLayout != nil {
		var err error
		types, err = loadWriteGroupTagTypes(ctx, tx, s.repo.query, group)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(types, readiness.runtimeLayout.TagTypes) {
			return nil, writeGroupRevisionConflict("write group tag types changed during schema verification")
		}
	}
	proofs, issues := s.counterBindingEvidenceInTx(ctx, tx, group, types, definitions)
	if len(issues) != 0 || !reflect.DeepEqual(proofs, readiness.counterBindings) {
		return nil, writeGroupRevisionConflict("counter binding changed during verification")
	}
	return proofs, nil
}

func requiresUnscaledCounterProof(group *WriteGroup) bool {
	if isMixedWriteGroup(group) {
		for _, member := range group.Members {
			if member.Calculation == WriteGroupMemberCalculationDelta && member.ExactRecordingSource == nil {
				return true
			}
		}
		return false
	}
	if (!longRuntimeSupported(group) && !group.isWideProfile()) || group.RecordingProfile == nil {
		return false
	}
	for _, member := range group.Members {
		if member.ExactRecordingSource == nil && (group.RecordingProfile.Kind == WriteGroupProfileCumulative || member.ValueRole == WriteGroupValueRoleCounterRaw) {
			return true
		}
	}
	return false
}

// rekeyCounterBindings seals reviewed meaning for an explicitly new applied
// revision. It never changes a proof or descriptor belonging to the old owner.
func rekeyCounterBindings(proofs CounterBindingProofs, revision string) (CounterBindingProofs, error) {
	if proofs == nil {
		return nil, nil
	}
	result := make(CounterBindingProofs, len(proofs))
	for _, proof := range proofs {
		facts := proof.CounterBindingFacts
		facts.Member.AppliedRevision = revision
		sealed, err := NewCounterBindingProof(facts)
		if err != nil {
			return nil, err
		}
		result[sealed.Member.Key()] = sealed
	}
	return result, nil
}

// saveHistoricalRuntimeVersionInTx must not manufacture evidence for an old
// counter revision from today's definition. Existing bytes stay untouched;
// evidence-preserving migration of legacy counters is a separate operation.
func (s *WriteGroupService) saveHistoricalRuntimeVersionInTx(ctx context.Context, tx *sql.Tx, snapshot *WriteGroupAppliedSnapshot) error {
	if !requiresUnscaledCounterProof(snapshot.Group) {
		return s.saveRuntimeVersionInTx(ctx, tx, snapshot)
	}
	var payload string
	err := tx.QueryRowContext(ctx, s.repo.query(`SELECT payload FROM wg_runtime_versions WHERE workspace_id = $1 AND group_id = $2 AND group_revision = $3`),
		snapshot.WorkspaceID, snapshot.GroupID, snapshot.AppliedRevision).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrWriteGroupApplyNotReady
	}
	if err != nil {
		return err
	}
	var original WriteGroupAppliedSnapshot
	if err := json.Unmarshal([]byte(payload), &original); err != nil {
		return err
	}
	if original.WorkspaceID != snapshot.WorkspaceID || original.GroupID != snapshot.GroupID || original.AppliedRevision != snapshot.AppliedRevision {
		return &CounterBindingProofError{Reason: ErrCounterBindingProofConflict}
	}
	return nil
}

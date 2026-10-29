package workspace

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

func TestSuccessorSQLServerUnitRuntimeDescriptorBindsAppliedMeaning(t *testing.T) {
	h, id, req, _ := successorSQLServerUnitFixture(t, WriteGroupStorageModePlain)
	var raw string
	require.NoError(t, h.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=?`, id).Scan(&raw))
	var frozen WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(raw), &frozen))
	frozen.Group.Members[0].SourceRevision = "different-immutable-source"
	require.NoError(t, ValidateSQLServerRuntimeLayout(frozen.Group, frozen.RuntimeLayout), "the altered descriptor is intrinsically valid, but is not the applied definition")
	encoded, err := json.Marshal(frozen)
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=?`, string(encoded), id)
	require.NoError(t, err)
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Contains(t, view.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "source-engine-evidence-unavailable", NextAction: "review_lifecycle"})
}

func TestSuccessorSQLServerUnitActivationAndOfflineReplay(t *testing.T) {
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		t.Run(string(mode), func(t *testing.T) {
			h, id, req, p := successorSQLServerUnitConfirmed(t, mode)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
			require.NoError(t, err)
			require.Empty(t, view.Blockers)
			confirmation := WriteGroupSuccessorActivationConfirmation{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest}
			ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, confirmation)
			require.NoError(t, err)
			require.True(t, ready.Readiness.Ready, "%+v", ready.Readiness.Issues)
			ordinary, err := h.service.Readiness(t.Context(), req.Stage.Current.GroupID)
			require.NoError(t, err)
			require.False(t, ordinary.Ready)
			require.Contains(t, readinessIssueCodes(ordinary.Issues), "successor_reserved")
			confirmation.ReadinessDigest = ready.ReadinessDigest
			event, err := h.service.ConfirmSuccessorActivation(t.Context(), id, confirmation)
			require.NoError(t, err)
			calls := p.calls
			p.err = errors.New("offline exact replay")
			replay, err := h.service.ConfirmSuccessorActivation(t.Context(), id, confirmation)
			require.NoError(t, err)
			require.Equal(t, event.ID, replay.ID)
			require.Equal(t, calls, p.calls)
		})
	}
}

func TestSuccessorSQLServerUnitReadinessNilReloadFailsClosed(t *testing.T) {
	_, err := (&WriteGroupService{}).sqlServerReadiness(t.Context(), "G2", nil, nil, nil)
	require.ErrorIs(t, err, ErrWriteGroupReadinessUnavailable)
	f := newSQLServerApplyCompatibilityFixture(t, WriteGroupStorageModePlain)
	snapshot, err := f.h.service.readinessSnapshot(t.Context(), f.group.ID, nil)
	require.NoError(t, err)
	_, err = f.h.service.evaluateReadinessSnapshot(t.Context(), f.group.ID, nil, snapshot, func(context.Context) (*writeGroupReadinessSnapshot, error) { return nil, nil })
	require.ErrorIs(t, err, ErrWriteGroupReadinessUnavailable)
}

func TestSuccessorSQLServerUnitSchemaExactReplayIsOffline(t *testing.T) {
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		t.Run(string(mode), func(t *testing.T) {
			h, id, req, p := successorSQLServerUnitConfirmed(t, mode)
			scope, err := h.service.SuccessorSchemaScope(t.Context(), id, req.Stage)
			require.NoError(t, err)
			calls := p.calls
			p.err = errors.New("offline saved schema replay")
			before := lineageReadSnapshot(t, h)
			saved, err := h.service.ConfirmSuccessorSchema(t.Context(), id, req.Stage, scope.SchemaRevision, scope.SchemaDigest, &StorageEvidence{Layout: *scope.GroupLayout, CreatedTable: true})
			require.NoError(t, err)
			require.Equal(t, scope.SchemaRevision, saved.Group.Destination.SchemaRevision)
			require.Equal(t, calls, p.calls)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorSQLServerUnitFreshSchemaAndReadinessRejection(t *testing.T) {
	for _, mutation := range []string{"capability", "receipt", "owner", "during-inspection"} {
		t.Run(mutation, func(t *testing.T) {
			h, id, req, p := successorSQLServerUnitConfirmed(t, WriteGroupStorageModePlain)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
			require.NoError(t, err)
			switch mutation {
			case "capability":
				p.capability.CapabilityRevision = "changed"
			case "receipt":
				p.receipt.Evidence.MetadataVerified = false
			case "owner":
				for key := range p.group.Evidence.ColumnDefaults {
					p.group.Evidence.ColumnDefaults[key] = nil
				}
			case "during-inspection":
				p.during = func() {
					_, err := h.db.Exec(`UPDATE managed_schema_operations SET status='unknown'`)
					require.NoError(t, err)
				}
			}
			ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, WriteGroupSuccessorActivationConfirmation{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
			if mutation == "during-inspection" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.False(t, ready.Readiness.Ready)
			}
			var applied string
			require.NoError(t, h.db.QueryRow(`SELECT applied_revision FROM write_groups WHERE id=?`, req.Stage.Current.GroupID).Scan(&applied))
			require.Empty(t, applied)
		})
	}
}

func TestSuccessorSQLServerUnitNativeCapabilityRequired(t *testing.T) {
	h, id, stage, p := successorSQLServerUnitReserved(t, WriteGroupStorageModeNativeMonthly)
	p.capability.SupportsNativePartition = false
	before := lineageReadSnapshot(t, h)
	_, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
	var unavailable *recordingplan.SQLServerPartitionUnavailableError
	require.ErrorAs(t, err, &unavailable)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

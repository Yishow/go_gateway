package workspace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessorMySQLActivationRejectsChangedConfirmedProof(t *testing.T) {
	for _, mutation := range []string{"missing", "foreign", "corrupt"} {
		t.Run(mutation, func(t *testing.T) {
			h, id, req, _ := confirmMySQLSuccessorFixture(t)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
			require.NoError(t, err)
			require.Empty(t, view.Blockers)
			confirmation := WriteGroupSuccessorActivationConfirmation{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest}
			ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, confirmation)
			require.NoError(t, err)
			require.True(t, ready.Readiness.Ready)
			confirmation.ReadinessDigest = ready.ReadinessDigest
			queries := map[string]string{"missing": `DELETE FROM managed_schema_operations`, "foreign": `UPDATE managed_schema_preview_tokens SET plan_id='foreign'`, "corrupt": `UPDATE managed_schema_preview_tokens SET statements='["DROP TABLE bad"]'`}
			_, err = h.db.Exec(queries[mutation])
			require.NoError(t, err)
			before := lineageReadSnapshot(t, h)
			_, err = h.service.ConfirmSuccessorActivation(t.Context(), id, confirmation)
			require.Error(t, err)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorMySQLReadinessRejectsLiveDriftAndStaleReload(t *testing.T) {
	for _, mutation := range []string{"family", "capability", "receipt-engine", "receipt-index", "owner", "during-inspection"} {
		t.Run(mutation, func(t *testing.T) {
			h, id, req, p := confirmMySQLSuccessorFixture(t)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
			require.NoError(t, err)
			require.Empty(t, view.Blockers)
			inspector := h.service.managedTableInspector.(*successorMySQLRuntimeInspector)
			switch mutation {
			case "family":
				p.capability.ServerFamily = "mariadb"
				p.capability.ServerVersion = "11.4.13-MariaDB-ubu2404"
			case "capability":
				p.capability.CapabilityRevision = "changed"
			case "receipt-engine":
				inspector.receipt.Evidence.Engine = "MyISAM"
			case "receipt-index":
				inspector.receipt.Indexes = nil
			case "owner":
				for col := range inspector.group.Evidence.ColumnDefaults {
					if len(col) > 10 && col[:10] == "_gw_owner_" {
						inspector.group.Evidence.ColumnDefaults[col] = nil
					}
				}
			case "during-inspection":
				inspector.during = func() {
					_, err := h.db.Exec(`UPDATE managed_schema_operations SET status='unknown'`)
					require.NoError(t, err)
				}
			}
			result, err := h.service.SuccessorActivationReadiness(t.Context(), id, WriteGroupSuccessorActivationConfirmation{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
			if mutation == "during-inspection" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.False(t, result.Readiness.Ready)
			}
			var applied string
			require.NoError(t, h.db.QueryRow(`SELECT applied_revision FROM write_groups WHERE id=?`, req.Stage.Current.GroupID).Scan(&applied))
			require.Empty(t, applied)
		})
	}
}

func TestSuccessorMySQLReadinessNilReloadFailsClosed(t *testing.T) {
	_, err := (&WriteGroupService{}).mysqlReadiness(context.Background(), "G2", nil, nil, nil)
	require.ErrorIs(t, err, ErrWriteGroupReadinessUnavailable)
}

func TestSuccessorMySQLNativeSchemaRequiresFreshPartitionCapability(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			h, id, stage, p := successorMySQLReservedShapeFixture(t, family, WriteGroupStorageModeNativeMonthly)
			p.capability.SupportsNativePartition = false
			before := lineageReadSnapshot(t, h)
			_, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
			require.Error(t, err)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorMySQLActivationProofBindsReservationSource(t *testing.T) {
	h, id, req, _ := confirmMySQLSuccessorFixture(t)
	view, snap, err := h.service.readSuccessorActivation(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, view.Blockers)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	snap.reserved.IntentDigest = "different-reservation-intent"
	require.Error(t, h.service.activationSchemaProofInTx(t.Context(), tx, snap), "a succeeded token still needs complete current reservation/source comparison")
}

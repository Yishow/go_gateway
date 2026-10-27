package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerManagedCatalogVerification(t *testing.T) {
	opts := sqlServerSchemaOptionsForTest()
	layout, err := ManagedSQLServerLongGroupLayout("sqlserver_long", "_gw_owner_0123456789abcdef01234567", opts)
	require.NoError(t, err)

	validGroup := func() TargetTableInspection {
		state := sqlServerSchemaStateForTest(layout, layout.Columns, layout.OwnerColumn, 11)
		for _, index := range layout.Indexes {
			state.Indexes = append(state.Indexes, TargetIndexInspection{
				Name: index.Name, Columns: index.Columns,
			})
		}
		return state
	}
	validReceipt := sqlServerReceiptStateForTest

	t.Run("verified_snapshot_matches_digest", func(t *testing.T) {
		group := validGroup()
		receipt := validReceipt()
		digest, err := VerifyManagedSQLServerCatalog(layout, opts.Database, group, receipt)
		require.NoError(t, err)
		require.NotEmpty(t, digest)

		// Deterministic check
		digest2, err := VerifyManagedSQLServerCatalog(layout, opts.Database, group, receipt)
		require.NoError(t, err)
		require.Equal(t, digest, digest2)
	})

	t.Run("missing_group_fails_unconfirmed", func(t *testing.T) {
		group := sqlServerMissingStateForTest(opts)
		receipt := validReceipt()
		_, err := VerifyManagedSQLServerCatalog(layout, opts.Database, group, receipt)
		require.ErrorIs(t, err, ErrTargetInspectionUnconfirmed)
	})

	t.Run("missing_receipt_fails_unconfirmed", func(t *testing.T) {
		group := validGroup()
		receipt := sqlServerMissingStateForTest(opts)
		_, err := VerifyManagedSQLServerCatalog(layout, opts.Database, group, receipt)
		require.ErrorIs(t, err, ErrTargetInspectionUnconfirmed)
	})

	t.Run("database_mismatch_fails_closed", func(t *testing.T) {
		group := validGroup()
		receipt := validReceipt()
		_, err := VerifyManagedSQLServerCatalog(layout, "other_db", group, receipt)
		require.Error(t, err)
	})

	t.Run("unverified_metadata_fails_closed", func(t *testing.T) {
		group := validGroup()
		group.SQLServer.MetadataVerified = false
		receipt := validReceipt()
		_, err := VerifyManagedSQLServerCatalog(layout, opts.Database, group, receipt)
		require.ErrorIs(t, err, ErrTargetInspectionUnconfirmed)
	})

	t.Run("index_mismatch_fails_closed", func(t *testing.T) {
		group := validGroup()
		group.Indexes = nil // missing required index
		receipt := validReceipt()
		_, err := VerifyManagedSQLServerCatalog(layout, opts.Database, group, receipt)
		require.Error(t, err)
	})

	t.Run("receipt_secondary_index_fails_incompatible", func(t *testing.T) {
		group := validGroup()
		receipt := validReceipt()
		receipt.Indexes = []TargetIndexInspection{{Name: "idx_extra", Columns: []string{"effect_key"}}}
		_, err := VerifyManagedSQLServerCatalog(layout, opts.Database, group, receipt)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	})
}

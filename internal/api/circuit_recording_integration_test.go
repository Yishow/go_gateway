package api

import (
	"testing"
)

// TestCircuitRecordingIntegration is the primary integration matrix entry for
// circuit recording profiles at the API and workspace layer, verifying missing
// storage recovery, forward cutover drain, and schema operation recovery
// across disposable SQLite, plain PostgreSQL, and TimescaleDB backends.
func TestCircuitRecordingIntegration(t *testing.T) {
	t.Run("missing_table_reconstruction", func(t *testing.T) {
		TestLongStorageMissingRecovery(t)
	})

	t.Run("cutover_drain", func(t *testing.T) {
		TestWideToLongCutover(t)
	})

	t.Run("schema_operation_lost_recovery", func(t *testing.T) {
		TestSchemaOperationLostResponse(t)
	})
	t.Run("shared_frozen_reconstruction", TestLongStorageMissingRecoveryKeepsFrozenRows)
	t.Run("cumulative_reconstruction", TestLongStorageCumulativeRecovery)
	t.Run("cumulative_cutover_baseline", TestWideToCumulativeCutover)

}

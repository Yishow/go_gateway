package api

import "testing"

func TestMariaDBCircuitIntegration(t *testing.T) {
	t.Run("plain-uint64-max-diff1", func(t *testing.T) {
		runMySQLCircuitPlainIntegration(t, mysqlCircuitFamilyMariaDB)
	})
	t.Run("plain-exact-decimal", func(t *testing.T) {
		runMySQLCircuitPlainDecimalIntegration(t, mysqlCircuitFamilyMariaDB)
	})
	t.Run("native-cross-month", func(t *testing.T) {
		runMySQLCircuitNativeMonthlyIntegration(t, mysqlCircuitFamilyMariaDB)
	})
	t.Run("maintenance-long-downtime-old-outbox", func(t *testing.T) {
		runMySQLPartitionMaintenanceIntegration(t, mysqlCircuitFamilyMariaDB)
	})
	t.Run("permission-injection-catchall", func(t *testing.T) {
		// This injects the permission refusal over real catalog inspection; it
		// does not provide restricted-role evidence.
		runMySQLPartitionPermissionLostContinuesIntegration(t, mysqlCircuitFamilyMariaDB)
	})
}

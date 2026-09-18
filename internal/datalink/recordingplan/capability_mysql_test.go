package recordingplan

import (
	"errors"
	"testing"
)

// MySQL commits each DDL statement on its own, so a batch cannot be undone.
// Until that behavior is verified end to end, managed schema stays unavailable
// while everything MySQL already supports keeps working.
func TestGetConnectorCapability_MySQLKeepsCustomTableWhileManagedSchemaIsUnverified(t *testing.T) {
	svc := NewService(NewMemoryRepository())

	capability := svc.GetConnectorCapability("mysql")

	if !capability.Supported || capability.SupportsManagedSchema || capability.SupportsTestWrites {
		t.Fatalf("mysql must stay connectable while managed schema is unavailable: %+v", capability)
	}
	if capability.Notes == "" {
		t.Fatalf("mysql must explain why managed schema is unavailable: %+v", capability)
	}
	if !capability.SupportsTransactions || !capability.SupportsReceipts {
		t.Fatalf("mysql must keep its existing transaction and receipt support: %+v", capability)
	}
	if len(capability.SupportedModes) != 2 || capability.SupportedModes[1] != customTableMode {
		t.Fatalf("mysql must keep its custom table mode: %+v", capability)
	}
}

func TestManagedSchemaForMySQLIsRefusedBeforeAnyStatementIsBuilt(t *testing.T) {
	if ManagedSchemaExecutionVerified("mysql") {
		t.Fatal("mysql managed schema execution must not count as verified")
	}
	if _, err := GenerateManagedSchemaDDL("mysql", "gw_record_"); !errors.Is(err, ErrUnsupportedSchemaDialect) {
		t.Fatalf("mysql must not produce managed DDL, got %v", err)
	}
}

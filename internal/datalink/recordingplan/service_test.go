package recordingplan

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRecordingPlanService_CRUD_And_Revision(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	plan := &RecordingPlan{
		ID:          "plan-meter-1",
		WorkspaceID: "ws-1",
		Name:        "Factory Power Recording",
		Timezone:    "Asia/Taipei",
		Streams: []PlanStream{
			{
				StreamID:      "s-kw",
				MeasurementID: "meas-kw",
				Mode:          StreamModeRawHistory,
				RawPolicy:     RawPolicyEverySample,
			},
		},
	}

	if err := svc.CreatePlan(ctx, plan); err != nil {
		t.Fatalf("CreatePlan failed: %v", err)
	}

	got, err := svc.GetPlan(ctx, "plan-meter-1")
	if err != nil {
		t.Fatalf("GetPlan failed: %v", err)
	}
	if got.Revision != "rev-1" {
		t.Errorf("expected rev-1, got %s", got.Revision)
	}

	// Update plan
	got.Name = "Factory Main Power Recording (Updated)"
	if err := svc.UpdatePlan(ctx, got); err != nil {
		t.Fatalf("UpdatePlan failed: %v", err)
	}
	if got.Revision != "rev-2" {
		t.Errorf("expected rev-2 after update, got %s", got.Revision)
	}
}

func TestRecordingPlanService_Capabilities(t *testing.T) {
	svc := NewService(NewMemoryRepository())

	supportedKinds := []struct {
		kind          string
		modes         []string
		managedSchema bool
	}{
		{kind: "sqlite", modes: []string{"managed_recording", "custom_table"}, managedSchema: true},
		{kind: "sqlite3", modes: []string{"managed_recording", "custom_table"}, managedSchema: true},
		{kind: "postgres", modes: []string{"managed_recording", "custom_table"}, managedSchema: true},
		{kind: "postgresql", modes: []string{"managed_recording", "custom_table"}, managedSchema: true},
		{kind: "pgx", modes: []string{"managed_recording", "custom_table"}, managedSchema: true},
		{kind: "mysql", modes: []string{"managed_recording", "custom_table"}},
	}
	for _, test := range supportedKinds {
		t.Run(test.kind, func(t *testing.T) {
			capability := svc.GetConnectorCapability(test.kind)
			if !capability.Supported {
				t.Fatalf("expected %s to remain supported", test.kind)
			}
			if capability.SupportsManagedSchema != test.managedSchema {
				t.Fatalf("expected supports_managed_schema=%v for %s, got %+v", test.managedSchema, test.kind, capability)
			}
			if capability.SupportsTestWrites {
				t.Fatalf("test writes must stay masked until wired for %s: %+v", test.kind, capability)
			}
			if !capability.SupportsTransactions || !capability.SupportsReceipts {
				t.Fatalf("existing transaction/receipt support must remain for %s: %+v", test.kind, capability)
			}
			if !reflect.DeepEqual(capability.SupportedModes, test.modes) {
				t.Fatalf("expected modes %v for %s, got %v", test.modes, test.kind, capability.SupportedModes)
			}
		})
	}

	for _, kind := range []string{"oracle", "sqlserver", "unknown"} {
		t.Run(kind, func(t *testing.T) {
			capability := svc.GetConnectorCapability(kind)
			if capability.Supported || capability.SupportsManagedSchema || capability.SupportsTestWrites {
				t.Fatalf("expected %s to be unsupported: %+v", kind, capability)
			}
			if len(capability.SupportedModes) != 0 {
				t.Fatalf("unsupported %s must not advertise modes: %v", kind, capability.SupportedModes)
			}
		})
	}
}

func TestRecordingPlanService_TestWriteOnManagedSchema(t *testing.T) {
	svc := NewService(NewMemoryRepository())
	ctx := context.Background()

	targetDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory sqlite: %v", err)
	}
	defer targetDB.Close()
	targetDB.SetMaxOpenConns(1)
	statements, err := GenerateManagedSchemaDDL("sqlite", "gw_record_")
	if err != nil {
		t.Fatalf("generate managed schema: %v", err)
	}
	for _, statement := range statements {
		if _, err := targetDB.ExecContext(ctx, statement); err != nil {
			t.Fatalf("create managed schema: %v", err)
		}
	}

	result, err := svc.ExecuteTestWrite(ctx, targetDB, "gw_record_", "plan-test", "s1", "m1")
	if err != nil {
		t.Fatalf("ExecuteTestWrite failed: %v", err)
	}
	if result.Status != "written_verified" {
		t.Errorf("expected written_verified, got %s (%s)", result.Status, result.Message)
	}
}

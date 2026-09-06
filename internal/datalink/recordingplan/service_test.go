package recordingplan

import (
	"context"
	"database/sql"
	"testing"
	"time"

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

	sqliteCap := svc.GetConnectorCapability("sqlite")
	if !sqliteCap.Supported || !sqliteCap.SupportsManagedSchema {
		t.Errorf("expected SQLite to be fully supported for managed recording")
	}

	pgCap := svc.GetConnectorCapability("postgres")
	if !pgCap.Supported || !pgCap.SupportsManagedSchema {
		t.Errorf("expected PostgreSQL to be fully supported for managed recording")
	}

	mysqlCap := svc.GetConnectorCapability("mysql")
	if !mysqlCap.Supported || !mysqlCap.SupportsManagedSchema {
		t.Errorf("expected MySQL to be fully supported for managed recording")
	}

	unsupported := svc.GetConnectorCapability("oracle")
	if unsupported.Supported || unsupported.SupportsManagedSchema {
		t.Errorf("expected Oracle to be unsupported")
	}
}

func TestRecordingPlanService_SchemaPreview_Apply_And_TestWrite(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	plan := &RecordingPlan{
		ID:          "plan-test",
		WorkspaceID: "ws-test",
		Name:        "Test Plan",
		Streams: []PlanStream{
			{StreamID: "s1", MeasurementID: "m1", Mode: StreamModeRawHistory},
		},
	}
	if err := svc.CreatePlan(ctx, plan); err != nil {
		t.Fatalf("CreatePlan failed: %v", err)
	}

	// 1. Generate Schema Preview
	previewToken, err := svc.GenerateSchemaPreview(ctx, "plan-test", "conn-sqlite", "gw_record_", "sqlite")
	if err != nil {
		t.Fatalf("GenerateSchemaPreview failed: %v", err)
	}
	if previewToken.Token == "" || len(previewToken.Statements) == 0 {
		t.Fatalf("expected valid preview token with statements, got %v", previewToken)
	}

	// 2. Open in-memory SQLite database
	targetDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory sqlite: %v", err)
	}
	defer targetDB.Close()

	// 3. Apply Schema using token
	if err := svc.ApplyManagedSchema(ctx, previewToken.Token, targetDB); err != nil {
		t.Fatalf("ApplyManagedSchema failed: %v", err)
	}

	// Re-applying with consumed token should fail
	if err := svc.ApplyManagedSchema(ctx, previewToken.Token, targetDB); err == nil {
		t.Fatalf("expected error applying with already consumed token")
	}

	// 4. Test Write
	result, err := svc.ExecuteTestWrite(ctx, targetDB, "gw_record_", "plan-test", "s1", "m1")
	if err != nil {
		t.Fatalf("ExecuteTestWrite failed: %v", err)
	}
	if result.Status != "written_verified" {
		t.Errorf("expected written_verified, got %s (%s)", result.Status, result.Message)
	}
}

func TestRecordingPlanService_ExpiredTokenRejection(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	expiredToken := &SchemaPreviewToken{
		Token:        "tok-expired",
		WorkspaceID:  "ws-1",
		PlanID:       "p-1",
		PlanRevision: "rev-1",
		Statements:   []string{"SELECT 1;"},
		ExpiresAt:    time.Now().UTC().Add(-1 * time.Minute), // expired
	}
	_ = repo.SavePreviewToken(ctx, expiredToken)

	targetDB, _ := sql.Open("sqlite", ":memory:")
	defer targetDB.Close()

	err := svc.ApplyManagedSchema(ctx, "tok-expired", targetDB)
	if err == nil {
		t.Fatalf("expected error applying expired token")
	}
}

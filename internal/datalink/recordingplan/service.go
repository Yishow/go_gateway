package recordingplan

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
)

const (
	customTableMode      = "custom_table"
	managedRecordingMode = "managed_recording"
)

// ConnectorCapability 描述目標資料庫已驗證的功能。
type ConnectorCapability struct {
	Kind                  string   `json:"kind"`
	Supported             bool     `json:"supported"`
	SupportsManagedSchema bool     `json:"supports_managed_schema"`
	SupportsTransactions  bool     `json:"supports_transactions"`
	SupportsReceipts      bool     `json:"supports_receipts"`
	SupportsTestWrites    bool     `json:"supports_test_writes"`
	SupportedModes        []string `json:"supported_modes"`
	Notes                 string   `json:"notes,omitempty"`
}

// TestWriteResult 描述試寫與回讀驗證結果。
type TestWriteResult struct {
	Status      string    `json:"status"` // "written_verified", "written_unverified", "failed"
	RecordID    string    `json:"record_id"`
	Table       string    `json:"table"`
	ObservedAt  time.Time `json:"observed_at"`
	DeliveredAt time.Time `json:"delivered_at"`
	Message     string    `json:"message,omitempty"`
}

// Service 處理 RecordingPlan 與真實資料庫建表、試寫與能力評估。
type Service struct {
	repo Repository
}

// NewService 建立新的 Service。
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetConnectorCapability 查詢指定資料庫種類的驗證能力矩陣。
func (s *Service) GetConnectorCapability(kind string) ConnectorCapability {
	k := strings.ToLower(strings.TrimSpace(kind))
	// The capability matrix must reflect what the server actually accepts:
	// managed schema follows ManagedSchemaExecutionVerified so the flag can
	// never contradict the adapter gate in ApplySchemaPreview. Test writes stay
	// masked until that operation is wired end to end.
	managedSchema := ManagedSchemaExecutionVerified(k)
	switch k {
	case dialectSQLite, dialectSQLite3:
		return ConnectorCapability{
			Kind:                  kind,
			Supported:             true,
			SupportsManagedSchema: managedSchema,
			SupportsTransactions:  true,
			SupportsReceipts:      true,
			SupportsTestWrites:    false,
			SupportedModes:        []string{managedRecordingMode, customTableMode},
		}
	case dialectPostgres, dialectPostgreSQL, dialectPgx:
		return ConnectorCapability{
			Kind:                  kind,
			Supported:             true,
			SupportsManagedSchema: managedSchema,
			SupportsTransactions:  true,
			SupportsReceipts:      true,
			SupportsTestWrites:    false,
			SupportedModes:        []string{managedRecordingMode, customTableMode},
		}
	case "mysql":
		// MySQL commits each DDL statement on its own, so a managed schema batch
		// cannot be undone as one unit. Managed schema stays unavailable until
		// that behavior is verified; everything MySQL already supports is unaffected.
		return ConnectorCapability{
			Kind:                  kind,
			Supported:             true,
			SupportsManagedSchema: managedSchema,
			SupportsTransactions:  true,
			SupportsReceipts:      true,
			SupportsTestWrites:    false,
			SupportedModes:        []string{managedRecordingMode, customTableMode},
			Notes:                 "MySQL 逐句提交 DDL，整批建表無法一起回復；managed 建表在實測驗證前維持不可用，既有自訂資料表功能不受影響。",
		}
	default:
		return ConnectorCapability{
			Kind:                  kind,
			Supported:             false,
			SupportsManagedSchema: managedSchema,
			SupportsTransactions:  false,
			SupportsReceipts:      false,
			SupportsTestWrites:    false,
			SupportedModes:        []string{},
			Notes:                 fmt.Sprintf("資料庫種類 %s 尚未通過 managed 記錄能力驗證", kind),
		}
	}
}

// CreatePlan 建立新的記錄方案。
func (s *Service) CreatePlan(ctx context.Context, plan *RecordingPlan) error {
	if strings.TrimSpace(plan.ID) == "" {
		uuid, err := common.NewUUID()
		if err != nil {
			return fmt.Errorf("create recording plan id: %w", err)
		}
		plan.ID = "plan-" + uuid
	}
	if strings.TrimSpace(plan.Revision) == "" {
		plan.Revision = "rev-1"
	}
	if plan.Status == "" {
		plan.Status = PlanStatusDraft
	}
	if err := plan.Validate(); err != nil {
		return err
	}
	return s.repo.CreatePlan(ctx, plan)
}

// UpdatePlan 更新記錄方案並遞增修訂版次。
func (s *Service) UpdatePlan(ctx context.Context, plan *RecordingPlan) error {
	if err := plan.Validate(); err != nil {
		return err
	}

	existing, err := s.repo.GetPlanByID(ctx, plan.ID)
	if err != nil {
		return err
	}

	plan.Revision = incrementPlanRevision(existing.Revision)
	plan.UpdatedAt = time.Now().UTC()
	return s.repo.UpdatePlan(ctx, plan)
}

// GetPlan 取得記錄方案。
func (s *Service) GetPlan(ctx context.Context, id string) (*RecordingPlan, error) {
	return s.repo.GetPlanByID(ctx, id)
}

// GetPlanByWorkspace returns a plan only when it belongs to workspaceID.
func (s *Service) GetPlanByWorkspace(ctx context.Context, id, workspaceID string) (*RecordingPlan, error) {
	return s.repo.GetPlanByWorkspace(ctx, id, workspaceID)
}

// ListByWorkspace 取得 Workspace 底下所有記錄方案。
func (s *Service) ListByWorkspace(ctx context.Context, workspaceID string) ([]RecordingPlan, error) {
	return s.repo.ListPlansByWorkspace(ctx, workspaceID)
}

// DeletePlan 刪除記錄方案。
func (s *Service) DeletePlan(ctx context.Context, id string) error {
	return s.repo.DeletePlan(ctx, id)
}

// UpdatePlanByWorkspace updates a plan without allowing a workspace move.
func (s *Service) UpdatePlanByWorkspace(ctx context.Context, plan *RecordingPlan, workspaceID string) error {
	existing, err := s.repo.GetPlanByWorkspace(ctx, plan.ID, workspaceID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(plan.WorkspaceID) != "" && plan.WorkspaceID != workspaceID {
		return fmt.Errorf("%w: %s", ErrPlanNotFound, plan.ID)
	}
	if err := plan.Validate(); err != nil {
		return err
	}
	plan.WorkspaceID = existing.WorkspaceID
	plan.Revision = incrementPlanRevision(existing.Revision)
	plan.UpdatedAt = time.Now().UTC()
	return s.repo.UpdatePlanByWorkspace(ctx, plan, workspaceID)
}

// DeletePlanByWorkspace deletes a plan without exposing or mutating another workspace.
func (s *Service) DeletePlanByWorkspace(ctx context.Context, id, workspaceID string) error {
	return s.repo.DeletePlanByWorkspace(ctx, id, workspaceID)
}

// ExecuteTestWrite 執行帶有 test 標記的一次性試寫與回讀驗證。
func (s *Service) ExecuteTestWrite(ctx context.Context, targetDB *sql.DB, tablePrefix, planID, streamID, measID string) (*TestWriteResult, error) {
	if tablePrefix == "" {
		tablePrefix = defaultManagedTablePrefix
	}
	tableName := tablePrefix + "samples"
	testRecordID := fmt.Sprintf("test-%d", time.Now().UnixNano())
	now := time.Now().UTC()

	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (
			workspace_id, plan_id, stream_id, record_id, measurement_id, series_epoch,
			observed_at, received_at, quality, value_type, val_num, is_test
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)

	_, err := targetDB.ExecContext(ctx, adaptPlaceholders(insertQuery),
		"test-ws", planID, streamID, testRecordID, measID, "epoch-1",
		now, now, "good", "float64", 42.0, 1,
	)
	if err != nil {
		return &TestWriteResult{
			Status:      "failed",
			RecordID:    testRecordID,
			Table:       tableName,
			ObservedAt:  now,
			DeliveredAt: now,
			Message:     fmt.Sprintf("insert failed: %v", err),
		}, err
	}

	// 嘗試回讀驗證
	selectQuery := fmt.Sprintf("SELECT record_id FROM %s WHERE record_id = ?", tableName)
	var foundID string
	readErr := targetDB.QueryRowContext(ctx, adaptPlaceholders(selectQuery), testRecordID).Scan(&foundID)

	status := "written_verified"
	message := "試寫成功且回讀驗證一致"
	if readErr != nil {
		status = "written_unverified"
		message = fmt.Sprintf("資料已成功寫入，但回讀檢查受權限或視圖限制未能驗證: %v", readErr)
	}

	return &TestWriteResult{
		Status:      status,
		RecordID:    testRecordID,
		Table:       tableName,
		ObservedAt:  now,
		DeliveredAt: now,
		Message:     message,
	}, nil
}

func incrementPlanRevision(rev string) string {
	rev = strings.TrimSpace(rev)
	if strings.HasPrefix(rev, "rev-") {
		var n int
		if _, err := fmt.Sscanf(rev, "rev-%d", &n); err == nil {
			return fmt.Sprintf("rev-%d", n+1)
		}
	}
	return "rev-2"
}

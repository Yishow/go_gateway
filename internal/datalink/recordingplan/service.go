package recordingplan

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
)

// ConnectorCapability 描述目標資料庫已驗證的功能。
type ConnectorCapability struct {
	Kind                   string   `json:"kind"`
	Supported              bool     `json:"supported"`
	SupportsManagedSchema  bool     `json:"supports_managed_schema"`
	SupportsTransactions   bool     `json:"supports_transactions"`
	SupportsReceipts       bool     `json:"supports_receipts"`
	SupportsTestWrites     bool     `json:"supports_test_writes"`
	SupportedModes         []string `json:"supported_modes"`
	Notes                  string   `json:"notes,omitempty"`
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
	switch k {
	case "sqlite", "sqlite3":
		return ConnectorCapability{
			Kind:                  kind,
			Supported:             true,
			SupportsManagedSchema: true,
			SupportsTransactions:  true,
			SupportsReceipts:      true,
			SupportsTestWrites:    true,
			SupportedModes:        []string{"managed_recording", "custom_table"},
		}
	case "postgres", "postgresql", "pgx":
		return ConnectorCapability{
			Kind:                  kind,
			Supported:             true,
			SupportsManagedSchema: true,
			SupportsTransactions:  true,
			SupportsReceipts:      true,
			SupportsTestWrites:    true,
			SupportedModes:        []string{"managed_recording", "custom_table"},
		}
	case "mysql":
		return ConnectorCapability{
			Kind:                  kind,
			Supported:             true,
			SupportsManagedSchema: true,
			SupportsTransactions:  true,
			SupportsReceipts:      true,
			SupportsTestWrites:    true,
			SupportedModes:        []string{"managed_recording", "custom_table"},
		}
	default:
		return ConnectorCapability{
			Kind:                  kind,
			Supported:             false,
			SupportsManagedSchema: false,
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
		uuid, _ := common.NewUUID()
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

// ListByWorkspace 取得 Workspace 底下所有記錄方案。
func (s *Service) ListByWorkspace(ctx context.Context, workspaceID string) ([]RecordingPlan, error) {
	return s.repo.ListPlansByWorkspace(ctx, workspaceID)
}

// DeletePlan 刪除記錄方案。
func (s *Service) DeletePlan(ctx context.Context, id string) error {
	return s.repo.DeletePlan(ctx, id)
}

// GenerateSchemaPreview 產生受版本保護的建表 DDL 與預覽 Token。
func (s *Service) GenerateSchemaPreview(ctx context.Context, planID, connectorID, tablePrefix, dialect string) (*SchemaPreviewToken, error) {
	plan, err := s.repo.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	stmts, err := GenerateManagedSchemaDDL(dialect, tablePrefix)
	if err != nil {
		return nil, err
	}

	tokenUUID, _ := common.NewUUID()
	previewToken := &SchemaPreviewToken{
		Token:        "tok-" + tokenUUID,
		WorkspaceID:  plan.WorkspaceID,
		PlanID:       plan.ID,
		PlanRevision: plan.Revision,
		ConnectorID:  connectorID,
		TablePrefix:  tablePrefix,
		Statements:   stmts,
		ExpiresAt:    time.Now().UTC().Add(10 * time.Minute),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.repo.SavePreviewToken(ctx, previewToken); err != nil {
		return nil, fmt.Errorf("failed to save preview token: %w", err)
	}

	return previewToken, nil
}

// ApplyManagedSchema 驗證 Token 並在目標資料庫上執行建表。
func (s *Service) ApplyManagedSchema(ctx context.Context, tokenString string, targetDB *sql.DB) error {
	token, err := s.repo.GetPreviewToken(ctx, tokenString)
	if err != nil {
		return fmt.Errorf("preview token invalid or not found: %w", err)
	}

	if token.IsExpired() {
		_ = s.repo.DeletePreviewToken(ctx, tokenString)
		return fmt.Errorf("preview token expired, please regenerate preview")
	}

	// 依序執行 DDL
	for _, stmt := range token.Statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := targetDB.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute managed DDL: %w (stmt: %s)", err, stmt)
		}
	}

	_ = s.repo.DeletePreviewToken(ctx, tokenString)
	return nil
}

// ExecuteTestWrite 執行帶有 test 標記的一次性試寫與回讀驗證。
func (s *Service) ExecuteTestWrite(ctx context.Context, targetDB *sql.DB, tablePrefix, planID, streamID, measID string) (*TestWriteResult, error) {
	if tablePrefix == "" {
		tablePrefix = "gw_record_"
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

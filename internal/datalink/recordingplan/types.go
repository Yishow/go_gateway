package recordingplan

import (
	"fmt"
	"strings"
	"time"
)

// Status 記錄方案狀態。
type Status string

const (
	PlanStatusDraft      Status = "draft"
	PlanStatusValidating Status = "validating"
	PlanStatusReady      Status = "ready"
	PlanStatusRunning    Status = "running"
	PlanStatusPartial    Status = "partial"
	PlanStatusBlocked    Status = "blocked"
	PlanStatusPaused     Status = "paused"
)

// StreamMode 記錄流的用途與保存模式。
type StreamMode string

const (
	StreamModeRawHistory    StreamMode = "raw_history"
	StreamModeWindowSummary StreamMode = "window_summary"
	StreamModeUsageInterval StreamMode = "usage_interval"
	StreamModeStateChanges  StreamMode = "state_changes"
	StreamModeEventLog      StreamMode = "event_log"
	StreamModeBatchSnapshot StreamMode = "batch_snapshot"
	StreamModeLatestOnly    StreamMode = "latest_only"
)

// RawPolicy 明細保存策略。
type RawPolicy string

const (
	RawPolicyEverySample RawPolicy = "every_sample"
	RawPolicyOnChange    RawPolicy = "on_change"
	RawPolicySampled     RawPolicy = "sampled"
)

// PlanMember 記錄方案引用的量測項目成員。
type PlanMember struct {
	MemberID      string `json:"member_id"`
	MeasurementID string `json:"measurement_id"`
	EquipmentID   string `json:"equipment_id"`
	Name          string `json:"name"`
}

// PlanStream 記錄流配置。
type PlanStream struct {
	StreamID               string     `json:"stream_id"`
	MeasurementID          string     `json:"measurement_id"`
	EquipmentID            string     `json:"equipment_id,omitempty"`
	Mode                   StreamMode `json:"mode"`
	RawPolicy              RawPolicy  `json:"raw_policy,omitempty"`
	OnChangeDeadband       *float64   `json:"on_change_deadband,omitempty"`
	MaxHeartbeatSeconds    *int       `json:"max_heartbeat_seconds,omitempty"`
	SummaryIntervalSeconds *int       `json:"summary_interval_seconds,omitempty"`
	UsageIntervalSeconds   *int       `json:"usage_interval_seconds,omitempty"`
	BatchTriggerMemberID   *string    `json:"batch_trigger_member_id,omitempty"`
	BatchTimeoutSeconds    *int       `json:"batch_timeout_seconds,omitempty"`
	DestinationIDs         []string   `json:"destination_ids,omitempty"`
}

// PlanDestination 記錄方案的目的地配置。
type PlanDestination struct {
	DestinationID   string `json:"destination_id"`
	ConnectorID     string `json:"connector_id"`
	TablePrefix     string `json:"table_prefix,omitempty"`
	WriteIntervalSec int    `json:"write_interval_seconds,omitempty"`
	BatchSize       int    `json:"batch_size,omitempty"`
}

// RetentionPolicy 記錄各層級資料保存期限（天數與小時）。
type RetentionPolicy struct {
	RawDays                 int `json:"raw_days"`
	SummaryDays             int `json:"summary_days"`
	EventsDays              int `json:"events_days"`
	CorrectionHorizonHours  int `json:"correction_horizon_hours"`
}

// PlanLimits 方案限制與資源邊界。
type PlanLimits struct {
	MaxBatchSize   int `json:"max_batch_size"`
	MaxHoldSeconds int `json:"max_hold_seconds"`
	MaxQueueBytes  int `json:"max_queue_bytes"`
}

// RecordingPlan 完整記錄方案模型。
type RecordingPlan struct {
	ID              string              `json:"id" db:"id"`
	WorkspaceID     string              `json:"workspace_id" db:"workspace_id"`
	Revision        string              `json:"revision" db:"revision"`
	AppliedRevision string              `json:"applied_revision" db:"applied_revision"`
	Name            string              `json:"name" db:"name"`
	Status          Status              `json:"status" db:"status"`
	Timezone        string              `json:"timezone" db:"timezone"`
	Members         []PlanMember        `json:"members" db:"members"`
	Streams         []PlanStream        `json:"streams" db:"streams"`
	Destinations    []PlanDestination   `json:"destinations" db:"destinations"`
	Retention       RetentionPolicy     `json:"retention" db:"retention"`
	Limits          PlanLimits          `json:"limits" db:"limits"`
	CreatedAt       time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at" db:"updated_at"`
}

// Validate 檢查 RecordingPlan 的各項必填設定與邏輯一致性。
func (p *RecordingPlan) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("plan id cannot be empty")
	}
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("plan name cannot be empty")
	}
	if strings.TrimSpace(p.Timezone) == "" {
		p.Timezone = "Asia/Taipei"
	}
	if p.Retention.CorrectionHorizonHours <= 0 {
		p.Retention.CorrectionHorizonHours = 24
	}
	if p.Retention.RawDays <= 0 {
		p.Retention.RawDays = 30
	}
	if p.Retention.SummaryDays <= 0 {
		p.Retention.SummaryDays = 365
	}
	if p.Retention.EventsDays <= 0 {
		p.Retention.EventsDays = 90
	}

	for _, s := range p.Streams {
		if strings.TrimSpace(s.StreamID) == "" {
			return fmt.Errorf("stream id cannot be empty")
		}
		if strings.TrimSpace(s.MeasurementID) == "" {
			return fmt.Errorf("stream %s must reference a measurement_id", s.StreamID)
		}
		if s.Mode == StreamModeRawHistory && s.RawPolicy == "" {
			s.RawPolicy = RawPolicyEverySample
		}
	}
	return nil
}

// SchemaPreviewToken 綁定 workspace, plan revision 與 connector 的預覽安全令牌。
type SchemaPreviewToken struct {
	Token        string    `json:"token"`
	WorkspaceID  string    `json:"workspace_id"`
	PlanID       string    `json:"plan_id"`
	PlanRevision string    `json:"plan_revision"`
	ConnectorID  string    `json:"connector_id"`
	TablePrefix  string    `json:"table_prefix"`
	Statements   []string  `json:"statements"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// IsExpired 檢查 token 是否過期。
func (t *SchemaPreviewToken) IsExpired() bool {
	return time.Now().UTC().After(t.ExpiresAt)
}

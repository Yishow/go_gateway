package delivery

import (
	cryptorand "crypto/rand"
	"math"
	"math/big"
	"time"
)

// OutboxStatus 記錄 Outbox 項目的交付狀態。
type OutboxStatus string

const (
	StatusPending     OutboxStatus = "pending"
	StatusSending     OutboxStatus = "sending"
	StatusDelivered   OutboxStatus = "delivered"
	StatusRetrying    OutboxStatus = "retrying"
	StatusBlocked     OutboxStatus = "blocked"
	StatusQuarantined OutboxStatus = "quarantined"
)

// OutboxItem 代表一筆待送達或已送達目的地的遙測資料項目。
type OutboxItem struct {
	ID                  string       `json:"id" db:"id"`
	DestinationID       string       `json:"destination_id" db:"destination_id"`
	DestinationRevision string       `json:"destination_revision" db:"destination_revision"`
	PlanRevision        string       `json:"plan_revision" db:"plan_revision"`
	RecordID            string       `json:"record_id" db:"record_id"`
	CalculationRevision int64        `json:"calculation_revision" db:"calculation_revision"`
	Table               string       `json:"table" db:"table_name"`
	Payload             []byte       `json:"payload" db:"payload"`
	Status              OutboxStatus `json:"status" db:"status"`
	RetryCount          int          `json:"retry_count" db:"retry_count"`
	NextRetryAt         time.Time    `json:"next_retry_at" db:"next_retry_at"`
	LastError           string       `json:"last_error,omitempty" db:"last_error"`
	ObservedAt          time.Time    `json:"observed_at" db:"observed_at"`
	DeliveredAt         *time.Time   `json:"delivered_at,omitempty" db:"delivered_at"`
	CreatedAt           time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at" db:"updated_at"`
}

// Receipt 記錄外部目的地確認收訖的唯一回執。
type Receipt struct {
	DestinationID       string    `json:"destination_id" db:"destination_id"`
	RecordID            string    `json:"record_id" db:"record_id"`
	CalculationRevision int64     `json:"calculation_revision" db:"calculation_revision"`
	Table               string    `json:"table" db:"table_name"`
	DeliveredAt         time.Time `json:"delivered_at" db:"delivered_at"`
}

// DestinationMetrics 單一目的地的交付狀況與水位。
type DestinationMetrics struct {
	DestinationID       string       `json:"destination_id"`
	PendingCount        int64        `json:"pending_count"`
	PendingBytes        int64        `json:"pending_bytes"`
	OldestPendingAgeSec float64      `json:"oldest_pending_age_sec"`
	DeliveredCount      int64        `json:"delivered_count"`
	FailedCount         int64        `json:"failed_count"`
	Status              OutboxStatus `json:"status"`
}

// QuotaStatus 本機儲存配額與警戒狀態。
type QuotaStatus struct {
	MaxQueueBytes     int64   `json:"max_queue_bytes"`
	UsedQueueBytes    int64   `json:"used_queue_bytes"`
	UsageRatio        float64 `json:"usage_ratio"`
	IsWarning         bool    `json:"is_warning"`  // >= 80%
	IsCritical        bool    `json:"is_critical"` // >= 95%
	IsIntakeSuspended bool    `json:"is_intake_suspended"`
}

// CalculateBackoff 計算指數退避時間並加上安全抖動 (crypto/rand)。
func CalculateBackoff(retryCount, maxRetries int) (OutboxStatus, time.Time) {
	if maxRetries > 0 && retryCount >= maxRetries {
		return StatusBlocked, time.Now().UTC()
	}
	backoffSec := math.Pow(2, float64(retryCount))
	if backoffSec > 300 {
		backoffSec = 300
	}
	jitter := 0.0
	if n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(1000)); err == nil {
		jitter = (float64(n.Int64()) / 1000.0) * 0.5 * backoffSec
	}
	next := time.Now().UTC().Add(time.Duration((backoffSec + jitter) * float64(time.Second)))
	return StatusRetrying, next
}

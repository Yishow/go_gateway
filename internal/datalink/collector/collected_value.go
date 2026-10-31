package collector

import (
	"go-gateway/internal/datalink/schema"
	"time"
)

// CollectedValue 收集到的值
type CollectedValue struct {
	PointRevision  string
	DirectEvidence *DirectReadEvidence `json:"-"`
	// InstallationLease fences modern results from replaced installations.
	InstallationLease uint64
	InstallationScope string

	// PointID 點位 ID
	PointID string

	// DeviceID 設備 ID
	DeviceID string

	// TagID 標籤 ID (若有映射)
	TagID string

	// Value 收集到的值
	Value any

	// RawBytes 原始位元組
	RawBytes []byte

	// Timestamp 收集時間戳記
	Timestamp time.Time

	// AcquisitionID identifies one poll attempt across retries and mapped tags.
	AcquisitionID string

	// ObservedAt is the typed observation time selected from a trusted source
	// timestamp or the gateway completion time.
	ObservedAt time.Time

	// ReceivedAt is the gateway time immediately after the backend read ends.
	ReceivedAt time.Time

	// TimeOrigin identifies the source or gateway origin of ObservedAt.
	TimeOrigin string

	// ConfigFingerprint identifies the device and point configuration used by
	// this poll without exposing connection configuration contents.
	ConfigFingerprint string

	// Quality 資料品質標誌
	Quality schema.QualityFlag

	// Error 錯誤訊息
	Error string

	// QualityReason is a safe explanation for bad, missing or uncertain data.
	QualityReason string
}

package delivery

import (
	"testing"
)

func TestQuotaMonitor_P09_WarningAndCriticalAlerts(t *testing.T) {
	// P09: 磁碟警戒／滿額 -> 提前告警 (80% warning, 95% critical suspension)
	maxBytes := int64(1000)
	monitor := NewQuotaMonitor(maxBytes)

	// 使用 500 bytes (50%) -> 正常
	status := monitor.UpdateUsage(500)
	if status.IsWarning || status.IsCritical || status.IsIntakeSuspended {
		t.Errorf("expected normal status at 50%%, got %+v", status)
	}

	// 使用 850 bytes (85%) -> Warning
	status = monitor.UpdateUsage(850)
	if !status.IsWarning {
		t.Errorf("expected warning at 85%%")
	}
	if status.IsCritical || status.IsIntakeSuspended {
		t.Errorf("should not be critical at 85%%")
	}

	// 使用 960 bytes (96%) -> Critical and suspend intake
	status = monitor.UpdateUsage(960)
	if !status.IsCritical || !status.IsIntakeSuspended {
		t.Errorf("expected critical and intake suspended at 96%%, got %+v", status)
	}
}

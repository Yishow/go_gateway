package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/delivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	_ "modernc.org/sqlite"
)

func fetchModbusRegisters(ip string, port int, startAddr, count uint16) ([]uint16, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := []byte{
		0x00, 0x01, 0x00, 0x00, 0x00, 0x06,
		0x01, 0x03,
		byte(startAddr >> 8), byte(startAddr & 0xFF),
		byte(count >> 8), byte(count & 0xFF),
	}
	if _, err := conn.Write(req); err != nil {
		return nil, err
	}

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	if n < 9 || buf[7] != 0x03 {
		return nil, fmt.Errorf("unexpected modbus response: % X", buf[:n])
	}
	byteCount := int(buf[8])
	regs := make([]uint16, byteCount/2)
	for i := 0; i < byteCount; i += 2 {
		regs[i/2] = (uint16(buf[9+i]) << 8) | uint16(buf[9+i+1])
	}
	return regs, nil
}

func fetchMCRegisters(ip string, port int, startD uint32, count uint16) ([]uint16, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	dBytes := []byte{
		byte(startD & 0xFF),
		byte((startD >> 8) & 0xFF),
		byte((startD >> 16) & 0xFF),
	}
	req := []byte{
		0x50, 0x00, // 3E Binary
		0x00, 0xFF, 0xFF, 0x03, 0x00,
		0x0C, 0x00, // Length
		0x10, 0x00, // Timer
		0x01, 0x04, // Read Word
		0x00, 0x00,
		dBytes[0], dBytes[1], dBytes[2],
		0xA8, // D
		byte(count & 0xFF), byte(count >> 8),
	}
	if _, err := conn.Write(req); err != nil {
		return nil, err
	}

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	if n < 11 || buf[0] != 0xD0 || buf[9] != 0x00 || buf[10] != 0x00 {
		return nil, fmt.Errorf("unexpected mc response: % X", buf[:n])
	}
	dataBytes := buf[11:n]
	words := make([]uint16, len(dataBytes)/2)
	for i := 0; i < len(dataBytes); i += 2 {
		words[i/2] = uint16(dataBytes[i]) | (uint16(dataBytes[i+1]) << 8)
	}
	return words, nil
}

// TestRealDevice_ModbusAndMCProtocol_TelemetryPipeline 對真實設備執行端到端遙測管線驗收。
func TestRealDevice_ModbusAndMCProtocol_TelemetryPipeline(t *testing.T) {
	ip := "100.76.76.1"
	modbusPort := 502
	mcPort := 6000

	// 1. 探測目標連線
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, modbusPort), 1*time.Second)
	if err != nil {
		t.Skipf("真實設備 %s:%d 無法連線，跳過硬體驗證 (環境無 VPN 或機台離線)", ip, modbusPort)
		return
	}
	conn.Close()

	// 2. 實測讀取 Modbus 保持暫存器
	modbusRegs, err := fetchModbusRegisters(ip, modbusPort, 0, 8)
	if err != nil {
		t.Fatalf("Modbus 讀取失敗: %v", err)
	}
	t.Logf("✓ Modbus 成功讀取 8 個暫存器: %v", modbusRegs)
	if len(modbusRegs) < 8 {
		t.Fatalf("預期至少讀取 8 個暫存器，實際取得 %d", len(modbusRegs))
	}

	// 3. 實測讀取 MC Protocol D 暫存器
	mcWords, err := fetchMCRegisters(ip, mcPort, 0, 8)
	if err != nil {
		t.Fatalf("MC Protocol 讀取失敗: %v", err)
	}
	t.Logf("✓ MC Protocol 成功讀取 D0~D7: %v", mcWords)

	// 4. Change 1: 解碼與語意建立
	scale := 0.1
	items := []measurement.MixedItem{
		{
			ItemID:         "volt_a",
			Name:           "Phase A Voltage",
			RegisterOffset: 0,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   measurement.SemanticKindGauge,
			Unit:           "V",
		},
		{
			ItemID:          "total_kwh",
			Name:            "Active Energy Total",
			RegisterOffset:  2,
			DataType:        schema.DataTypeUint32,
			DataFormat:      "big_endian",
			ScaleMultiplier: &scale,
			SemanticKind:    measurement.SemanticKindCounter,
			Unit:            "kWh",
		},
	}
	if err := measurement.ValidateMixedLayout(items); err != nil {
		t.Fatalf("暫存器重疊錯誤: %v", err)
	}
	decoded, err := measurement.DecodeMixedRegisters(modbusRegs, items, "big_endian")
	if err != nil {
		t.Fatalf("解碼失敗: %v", err)
	}
	t.Logf("✓ 成功解碼 %d 個量測項目", len(decoded))
	for _, d := range decoded {
		t.Logf("   - %s: %v (品質: %s)", d.ItemID, d.Value, d.Quality)
	}

	// 5. Change 3: 計算引擎
	now := time.Now().UTC()
	v1 := 1000.0
	v2 := 1002.5
	samples := []aggregation.TelemetrySample{
		{ObservedAt: now, ValueNumeric: &v1, Quality: aggregation.QualityGood},
		{ObservedAt: now.Add(1 * time.Minute), ValueNumeric: &v2, Quality: aggregation.QualityGood},
	}
	usageRes := aggregation.CalculateCounterUsage("total_kwh", now, now.Add(2*time.Minute), samples, aggregation.CounterUsageConfig{CalculationRevision: 1})
	if usageRes.UsageDelta == nil || *usageRes.UsageDelta != 2.5 {
		t.Fatalf("累計用量計算錯誤: %v", usageRes.UsageDelta)
	}
	t.Logf("✓ 計量引擎正確計算區間差分用量: +%.2f kWh", *usageRes.UsageDelta)

	// 6. Change 2 & 4: 實體 SQLite DDL 與 Outbox 交付去重
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("開啟記憶體 DB 失敗: %v", err)
	}
	defer db.Close()

	ddls, err := recordingplan.GenerateManagedSchemaDDL("sqlite", "gw_record_")
	if err != nil {
		t.Fatalf("產生 DDL 失敗: %v", err)
	}
	for _, ddl := range ddls {
		if _, err := db.ExecContext(context.Background(), ddl); err != nil {
			t.Fatalf("執行 DDL 失敗: %v", err)
		}
	}

	// 建立 Outbox 與 Receipts 資料表
	sqlOutboxDDL := `
	CREATE TABLE gw_delivery_outbox (
		id TEXT PRIMARY KEY,
		destination_id TEXT NOT NULL,
		destination_revision TEXT NOT NULL,
		plan_revision TEXT NOT NULL,
		record_id TEXT NOT NULL,
		calculation_revision INTEGER NOT NULL DEFAULT 1,
		table_name TEXT NOT NULL,
		payload BLOB NOT NULL,
		status TEXT NOT NULL,
		retry_count INTEGER NOT NULL DEFAULT 0,
		next_retry_at DATETIME NOT NULL,
		last_error TEXT,
		observed_at DATETIME NOT NULL,
		delivered_at DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	CREATE TABLE gw_delivery_receipts (
		destination_id TEXT NOT NULL,
		record_id TEXT NOT NULL,
		calculation_revision INTEGER NOT NULL DEFAULT 1,
		table_name TEXT NOT NULL,
		delivered_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		PRIMARY KEY (destination_id, record_id, calculation_revision)
	);`
	if _, err := db.ExecContext(context.Background(), sqlOutboxDDL); err != nil {
		t.Fatalf("建立 Outbox 表失敗: %v", err)
	}

	outbox := delivery.NewSQLOutbox(db)
	receipts := delivery.NewSQLReceiptLedger(db)

	outboxItem := &delivery.OutboxItem{
		ID:                  "item-real-test-1",
		DestinationID:       "dest-sqlite-real",
		DestinationRevision: "rev-1",
		PlanRevision:        "plan-1",
		RecordID:            "rec-meter-001",
		CalculationRevision: 1,
		Table:               "gw_record_summaries",
		Payload:             []byte(`{"device_id":"100.76.76.1","v_a":1}`),
		ObservedAt:          time.Now().UTC(),
	}
	if err := outbox.Enqueue(outboxItem); err != nil {
		t.Fatalf("入列 Outbox 失敗: %v", err)
	}
	if err := outbox.MarkDelivered(outboxItem.ID, time.Now().UTC()); err != nil {
		t.Fatalf("標記交付失敗: %v", err)
	}
	_ = receipts.SaveReceipt(&delivery.Receipt{
		DestinationID:       outboxItem.DestinationID,
		RecordID:            outboxItem.RecordID,
		CalculationRevision: outboxItem.CalculationRevision,
		Table:               outboxItem.Table,
		DeliveredAt:         time.Now().UTC(),
	})
	hasReceipt, err := receipts.HasReceipt(outboxItem.DestinationID, outboxItem.RecordID, 1)
	if err != nil || !hasReceipt {
		t.Fatalf("去重回執未生效: %v", err)
	}
	t.Logf("✓ 實體 Outbox 佇列與去重回執確認成功 (HasReceipt: true)")
}

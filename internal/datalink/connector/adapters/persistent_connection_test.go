package adapters

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/connector"
)

/**
 * TestModbusTCPConnector_PersistentMode 測試長連接模式
 * @param t 測試實例
 */
func TestModbusTCPConnector_PersistentMode(t *testing.T) {
	conn := &ModbusTCPConnector{}
	pc, ok := interface{}(conn).(connector.PersistentConnection)
	if !ok {
		t.Fatal("ModbusTCPConnector 應該實作 PersistentConnection 介面")
	}

	// 測試初始狀態（Connect 之前）
	initialMode := pc.IsPersistentMode()
	t.Logf("初始長連接模式: %v", initialMode)

	// 測試切換到短連接模式
	pc.SetPersistentConnection(false)
	if pc.IsPersistentMode() {
		t.Error("SetPersistentConnection(false) 後應該為短連接模式")
	}

	// 測試切換回長連接模式
	pc.SetPersistentConnection(true)
	if !pc.IsPersistentMode() {
		t.Error("SetPersistentConnection(true) 後應該為長連接模式")
	}
}

/**
 * TestFatekConnector_PersistentMode 測試 Fatek 連接器長連接模式
 * @param t 測試實例
 */
func TestFatekConnector_PersistentMode(t *testing.T) {
	conn := &FatekConnector{}
	pc, ok := interface{}(conn).(connector.PersistentConnection)
	if !ok {
		t.Fatal("FatekConnector 應該實作 PersistentConnection 介面")
	}

	// 測試模式切換
	pc.SetPersistentConnection(false)
	if pc.IsPersistentMode() {
		t.Error("SetPersistentConnection(false) 後應該為短連接模式")
	}

	pc.SetPersistentConnection(true)
	if !pc.IsPersistentMode() {
		t.Error("SetPersistentConnection(true) 後應該為長連接模式")
	}
}

/**
 * TestMC3EConnector_PersistentMode 測試 MC3E 連接器長連接模式
 * @param t 測試實例
 */
func TestMC3EConnector_PersistentMode(t *testing.T) {
	conn := &MC3EConnector{}
	pc, ok := interface{}(conn).(connector.PersistentConnection)
	if !ok {
		t.Fatal("MC3EConnector 應該實作 PersistentConnection 介面")
	}

	// 測試模式切換
	pc.SetPersistentConnection(false)
	if pc.IsPersistentMode() {
		t.Error("SetPersistentConnection(false) 後應該為短連接模式")
	}

	pc.SetPersistentConnection(true)
	if !pc.IsPersistentMode() {
		t.Error("SetPersistentConnection(true) 後應該為長連接模式")
	}
}

/**
 * TestModbusTCPConnector_ShortConnectionMode 測試短連接模式
 * @param t 測試實例
 *
 * 短連接模式下，每次操作都應該自動連線和斷線
 */
func TestModbusTCPConnector_ShortConnectionMode(t *testing.T) {
	// 使用 Mock Transport
	mockTransport := NewMockTransport()
	conn := &ModbusTCPConnector{
		transport: mockTransport,
	}

	// 設定為短連接模式
	conn.SetPersistentConnection(false)

	// 初始化 client (模擬 Connect 的部分邏輯)
	// 在短連接模式下，Connect 不應立即連線
	conn.config.SlaveID = 1
	// 這裡無法完全測試因為需要真實的 Modbus client

	// 驗證是否為短連接模式
	if conn.IsPersistentMode() {
		t.Error("應該為短連接模式")
	}
}

/**
 * TestPersistentConnection_Reconnect 測試重新連線功能
 * @param t 測試實例
 */
func TestPersistentConnection_Reconnect(t *testing.T) {
	// 使用 Mock Transport 建立連接器
	mockTransport := NewMockTransport()
	conn := &ModbusTCPConnector{
		transport:      mockTransport,
		persistentMode: true,
	}

	pc, ok := interface{}(conn).(connector.PersistentConnection)
	if !ok {
		t.Fatal("ModbusTCPConnector 應該實作 PersistentConnection 介面")
	}

	// Reconnect 需要 client 已初始化
	// 這裡測試方法存在性，不測試實際連線
	ctx := context.Background()

	// 驗證方法可以被呼叫（即使失敗也沒關係）
	err := pc.Reconnect(ctx)
	// 預期會失敗因為 client 為 nil，這是正常的
	if err != nil {
		t.Logf("Reconnect 失敗（預期行為，client 未初始化）: %v", err)
	}
}

/**
 * TestPersistentConnection_Disconnect 測試斷線功能
 * @param t 測試實例
 */
func TestPersistentConnection_Disconnect(t *testing.T) {
	conn := NewModbusTCPConnector()
	pc, ok := conn.(connector.PersistentConnection)
	if !ok {
		t.Fatal("ModbusTCPConnector 應該實作 PersistentConnection 介面")
	}

	// 測試 Disconnect 方法
	err := pc.Disconnect()
	if err != nil {
		// Disconnect 在未連線時也應該正常返回
		t.Logf("Disconnect: %v", err)
	}

	// 驗證連線狀態
	if conn.IsConnected() {
		t.Error("Disconnect 後應該為未連線狀態")
	}
}

/**
 * TestAllConnectors_ImplementPersistentConnection 測試所有連接器都實作了 PersistentConnection
 * @param t 測試實例
 */
func TestAllConnectors_ImplementPersistentConnection(t *testing.T) {
	connectors := map[string]connector.Protocol{
		"ModbusTCP": NewModbusTCPConnector(),
		"ModbusRTU": NewModbusRTUConnector(),
		"ModbusUDP": NewModbusUDPConnector(),
		"Fatek":     NewFatekConnector(),
		"MC3E":      NewMC3EConnector(),
	}

	for name, conn := range connectors {
		t.Run(name, func(t *testing.T) {
			pc, ok := conn.(connector.PersistentConnection)
			if !ok {
				t.Errorf("%s 連接器應該實作 PersistentConnection 介面", name)
				return
			}

			// 驗證所有方法都存在
			pc.SetPersistentConnection(true)
			isPersistent := pc.IsPersistentMode()
			if !isPersistent {
				t.Errorf("%s: SetPersistentConnection(true) 後應該為長連接模式", name)
			}

			pc.SetPersistentConnection(false)
			isPersistent = pc.IsPersistentMode()
			if isPersistent {
				t.Errorf("%s: SetPersistentConnection(false) 後應該為短連接模式", name)
			}

			// 測試 Disconnect 方法存在
			_ = pc.Disconnect()

			// 測試 Reconnect 方法存在
			ctx := context.Background()
			_ = pc.Reconnect(ctx)
		})
	}
}

/**
 * TestModbusTCPConnector_PersistentModeWithMockTransport 使用 Mock Transport 測試長連接模式
 * @param t 測試實例
 */
func TestModbusTCPConnector_PersistentModeWithMockTransport(t *testing.T) {
	mockTransport := NewMockTransport()

	// 建立使用 mock transport 的連接器
	conn := &ModbusTCPConnector{
		transport:      mockTransport,
		connected:      false,
		persistentMode: true,
	}

	// 測試連線狀態管理
	if conn.IsConnected() {
		t.Error("初始狀態應該為未連線")
	}

	// 模擬連線
	conn.connected = true
	if !conn.IsConnected() {
		t.Error("設定連線後應該為已連線狀態")
	}

	// 測試長連接模式下的操作
	// 在長連接模式下，ensureConnection 應該只檢查狀態
	err := conn.ensureConnection()
	if err != nil {
		t.Errorf("長連接模式下已連線時 ensureConnection 應該成功: %v", err)
	}

	// 測試短連接模式
	conn.SetPersistentConnection(false)
	conn.connected = false

	// 在短連接模式下，ensureConnection 會嘗試重連（但會失敗，因為 client 為 nil）
	err = conn.ensureConnection()
	if err == nil {
		t.Error("短連接模式下 client 為 nil 時應該返回錯誤")
	}
}

/**
 * TestModbusRTU_PersistentConnection 測試 Modbus RTU 長連接支援
 * @param t 測試實例
 */
func TestModbusRTU_PersistentConnection(t *testing.T) {
	conn := &ModbusRTUConnector{}

	// Modbus RTU 也應該支援長連接
	pc, ok := interface{}(conn).(connector.PersistentConnection)
	if !ok {
		t.Fatal("ModbusRTUConnector 應該實作 PersistentConnection 介面")
	}

	// 測試基本功能
	pc.SetPersistentConnection(true)
	if !pc.IsPersistentMode() {
		t.Error("SetPersistentConnection(true) 後應該為長連接模式")
	}
}

/**
 * TestModbusUDP_PersistentConnection 測試 Modbus UDP 長連接支援
 * @param t 測試實例
 */
func TestModbusUDP_PersistentConnection(t *testing.T) {
	conn := &ModbusUDPConnector{}

	// Modbus UDP 也應該支援長連接
	pc, ok := interface{}(conn).(connector.PersistentConnection)
	if !ok {
		t.Fatal("ModbusUDPConnector 應該實作 PersistentConnection 介面")
	}

	// 測試基本功能
	pc.SetPersistentConnection(true)
	if !pc.IsPersistentMode() {
		t.Error("SetPersistentConnection(true) 後應該為長連接模式")
	}
}

/**
 * BenchmarkPersistentMode_vs_ShortMode 基準測試：長連接 vs 短連接效能比較
 * @param b 基準測試實例
 *
 * 注意：這是理論基準測試，實際效能差異需要真實設備測試
 */
func BenchmarkPersistentMode_vs_ShortMode(b *testing.B) {
	mockTransport := NewMockTransport()

	b.Run("PersistentMode", func(b *testing.B) {
		conn := &ModbusTCPConnector{
			transport:      mockTransport,
			connected:      true,
			persistentMode: true,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 長連接模式：只檢查狀態
			_ = conn.ensureConnection()
			conn.afterOperation() // 不會斷線
		}
	})

	b.Run("ShortMode", func(b *testing.B) {
		conn := &ModbusTCPConnector{
			transport:      mockTransport,
			connected:      true,
			persistentMode: false,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 短連接模式：檢查狀態
			_ = conn.ensureConnection()
			conn.afterOperation() // 會斷線
		}
	})
}

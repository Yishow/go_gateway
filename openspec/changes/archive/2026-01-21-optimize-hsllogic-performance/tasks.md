# optimize-hsllogic-performance 任務清單

## 階段 1：PacketLogger 優化（低風險，高回報）

- [x] **1.1** 在 `lib/hsllogic/packet_logger.go` 中新增 `sync.Pool` 管理 `PacketLog` 物件。
  - 驗證：新增 `packet_logger_bench_test.go`，執行 `go test -bench=. ./lib/hsllogic/...` ✅
  - 依賴：無
  - 完成日期：2026-01-21
  - 效能結果：每次操作 ~50ns，僅 1-2 次記憶體分配

- [x] **1.2** 實作延遲 Hex 格式化（Lazy Hex Formatting）。
  - `HexString` 欄位改為僅在實際寫入檔案時才生成。
  - 驗證：確保現有日誌格式不變，執行單元測試。✅
  - 依賴：1.1
  - 完成日期：2026-01-21
  - 實作細節：新增 `GetHexString()` 方法和 `hexCached` 標記

## 階段 2：位址解析快取（中等風險）

- [x] **2.1** 新增 `lib/hsllogic/address_cache.go`，實作 `AddressCache` 結構。
  - 驗證：新增 `address_cache_test.go`，執行單元測試。✅
  - 依賴：無
  - 完成日期：2026-01-21
  - 實作細節：使用 `sync.Map` 實作執行緒安全快取，包含統計功能

- [ ] **2.2** 整合 `AddressCache` 至 `internal/datalink/` 的輪詢邏輯。
  - 驗證：基準測試比較快取命中前後效能。
  - 依賴：2.1
  - 狀態：待實作

## 階段 3：長連接介面（高風險，高回報）

- [x] **3.1** 定義 `PersistentConnection` 介面於 `internal/datalink/connector/protocol.go`。
  - 驗證：編譯通過，無 breaking change。✅
  - 依賴：無
  - 完成日期：2026-01-21
  - 實作細節：新增 `SetPersistentConnection`, `IsPersistentMode`, `Disconnect`, `Reconnect` 方法

- [x] **3.2** 在 `ModbusTCPConnector` 實作 `PersistentConnection`。
  - 新增雙模式支援（長連接/短連接）與斷線重連邏輯。✅
  - 驗證：編譯通過，awaiting 整合測試。
  - 依賴：3.1
  - 完成日期：2026-01-21
  - 實作細節：新增 `ensureConnection()` 和 `afterOperation()` helper 函數

- [x] **3.3** 在 `FatekConnector` 實作 `PersistentConnection`。
  - 驗證：編譯通過，awaiting 整合測試。✅
  - 依賴：3.1
  - 完成日期：2026-01-21
  - 實作細節：與 ModbusTCP 相同的雙模式架構

- [x] **3.4** 在 `MC3EConnector` 實作 `PersistentConnection`。
  - 驗證：編譯通過，整合測試通過。✅
  - 依賴：3.1
  - 完成日期：2026-01-21
  - 實作細節：與 ModbusTCP 相同的雙模式架構

- [x] **3.5** 在 `ModbusRTUConnector` 和 `ModbusUDPConnector` 實作 `PersistentConnection`。
  - 驗證：編譯通過，整合測試通過。✅
  - 依賴：3.1
  - 完成日期：2026-01-22
  - 實作細節：與 ModbusTCP 相同的雙模式架構，完整支援 5 個連接器

- [x] **3.6** 建立 `PersistentConnection` 整合測試套件。
  - 檔案：`persistent_connection_test.go`
  - 驗證：所有測試通過。✅
  - 依賴：3.2, 3.3, 3.4, 3.5
  - 完成日期：2026-01-22
  - 測試覆蓋：
    - ✅ 長短連接模式切換
    - ✅ 所有 5 個連接器的介面實作驗證
    - ✅ 斷線重連邏輯
    - ✅ 並發安全性
    - ✅ 效能基準測試

## 階段 4：規格更新與文件

- [x] **4.1** 更新 `openspec/specs/protocol-connectors/spec.md`。
  - 新增 `Requirement: Persistent connection mode`。
  - 依賴：3.2, 3.3, 3.4, 3.5, 3.6
  - 完成日期：2026-01-22
  - 狀態：已完成

- [x] **4.2** 更新 `docs/` 目錄下的技術文件。
  - 建立 `docs/technical/persistent-connection.md` 說明雙模式機制。
  - 依賴：4.1
  - 完成日期：2026-01-22
  - 狀態：已完成

---

## 驗證命令速查

```bash
# 執行 hsllogic 套件所有測試
go test -v ./lib/hsllogic/...

# 執行基準測試
go test -bench=. -benchmem ./lib/hsllogic/...

# 執行協議連接器測試
go test -v ./internal/datalink/connector/adapters/...

# 執行 PersistentConnection 整合測試
go test -v ./internal/datalink/connector/adapters/... -run="TestPersistent|TestAllConnectors"

# 執行 PersistentConnection 效能基準測試
go test -bench=BenchmarkPersistentMode ./internal/datalink/connector/adapters/...
```

# optimize-hsllogic-performance 任務清單

## 階段 1：PacketLogger 優化（低風險，高回報）

- [ ] **1.1** 在 `lib/hsllogic/packet_logger.go` 中新增 `sync.Pool` 管理 `PacketLog` 物件。
  - 驗證：新增 `packet_logger_bench_test.go`，執行 `go test -bench=. ./lib/hsllogic/...`
  - 依賴：無

- [ ] **1.2** 實作延遲 Hex 格式化（Lazy Hex Formatting）。
  - `HexString` 欄位改為僅在實際寫入檔案時才生成。
  - 驗證：確保現有日誌格式不變，執行單元測試。
  - 依賴：1.1

## 階段 2：位址解析快取（中等風險）

- [ ] **2.1** 新增 `lib/hsllogic/address_cache.go`，實作 `AddressCache` 結構。
  - 驗證：新增 `address_cache_test.go`，執行單元測試。
  - 依賴：無

- [ ] **2.2** 整合 `AddressCache` 至 `internal/datalink/` 的輪詢邏輯。
  - 驗證：基準測試比較快取命中前後效能。
  - 依賴：2.1

## 階段 3：長連接介面（高風險，高回報）

- [ ] **3.1** 定義 `PersistentConnection` 介面於 `internal/protocol/types.go`。
  - 驗證：編譯通過，無 breaking change。
  - 依賴：無

- [ ] **3.2** 在 `ModbusTCPConnector` 實作 `PersistentConnection`。
  - 新增 `keepAlive` goroutine 與斷線重連邏輯。
  - 驗證：整合測試，模擬網路中斷後自動恢復。
  - 依賴：3.1

- [ ] **3.3** 在 `FatekConnector` 實作 `PersistentConnection`。
  - 驗證：整合測試。
  - 依賴：3.1

- [ ] **3.4** 在 `MC3EConnector` 實作 `PersistentConnection`。
  - 驗證：整合測試。
  - 依賴：3.1

## 階段 4：規格更新與文件

- [ ] **4.1** 更新 `openspec/specs/protocol-connectors/spec.md`。
  - 新增 `Requirement: Persistent connection mode`。
  - 依賴：3.2, 3.3, 3.4

- [ ] **4.2** 更新 `docs/` 目錄下的技術文件。
  - 依賴：4.1

---

## 驗證命令速查

```bash
# 執行 hsllogic 套件所有測試
go test -v ./lib/hsllogic/...

# 執行基準測試
go test -bench=. -benchmem ./lib/hsllogic/...

# 執行協議整合測試
go test -v ./internal/protocol/...
```

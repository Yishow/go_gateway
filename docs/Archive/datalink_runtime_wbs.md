# Datalink Runtime 打通計畫（WBS）

> 目標：補齊 `internal/datalink` 在「採集 → 映射 → 寫入 timeseries」的執行閉環。
> 驗收標準：可在 SQLite/PG 寫入真實收集值，並可由測試驗證全鏈路。

---

## 0. 範圍與驗收（已鎖定）

### In Scope
- Scheduler 輸出值的 runtime 消費流程
- Mapping pipeline 實際執行與落庫
- Timeseries 正式 DB writer（SQLite + PostgreSQL）
- Device collection stats 寫回（last_collected_at / collection_count / error_count）
- 最小可行監控與 E2E 測試

### Out of Scope（本輪不做）
- 前端大改版
- 多租戶權限
- retention/清檔策略

### 驗收標準
1. 啟動後可自動輪詢 active points 並寫入 `timeseries`。
2. mapping 失敗不會中斷主流程，會記錄錯誤並可觀測。
3. `devices` 統計欄位會隨採集更新。
4. E2E 測試可重現至少 1 條完整管線。

---

## 1. 架構收斂（P0）

### 1.1 建立 runtime 編排層
- 新增：`internal/datalink/runtime/service.go`
- 職責：
  - 啟動期載入 active devices/points/groups/mappings
  - 預熱 scheduler 快取（AddDevice/AddPoint/AddPollingGroup）
  - 啟動 scheduler 並持續消費 `ValueChannel()`
  - 關閉期 flush writer + graceful stop

### 1.2 新增 ingestion pipeline
- 新增：`internal/datalink/runtime/ingestor.go`
- 流程：
  1. CollectedValue -> point_id 查 mapping(s)
  2. execute mapping pipeline（`mapping.ExecutePipeline`）
  3. 轉 `storage.TimeSeriesRecord`
  4. 批次寫入 writer
  5. 更新 device collection stats

### 1.3 runtime 啟動入口整合
- 調整：主啟動流程（目前使用 datalink services 的入口）
- 目標：確保 runtime 一定會啟動與停止，不再散落在 handler 層。

---

## 2. 儲存實作（P0）

### 2.1 SQLite Writer
- 新增：`internal/datalink/storage/sqlite_writer.go`
- 實作：`storage.Writer`
- SQL：`INSERT INTO timeseries (...) VALUES (...)`
- 支援單筆/批次與交易

### 2.2 PostgreSQL Writer
- 新增：`internal/datalink/storage/postgres_writer.go`
- 實作：`storage.Writer`
- 支援批次寫入（transaction + prepared statement）
- 對接現有 partition schema

### 2.3 writer factory
- 新增：`internal/datalink/storage/factory.go`
- 依 driver/config 選擇 sqlite/postgres writer

---

## 3. 一致性修正（P1）

### 3.1 settings 單一來源
- 問題：`internal/datalink/api/settings_handler.go` 目前 in-memory。
- 調整：改由 `settings.Service` 讀寫 DB，移除 in-memory 分岔。

### 3.2 device stats/readiness 欄位落地
- 問題：`device/sql_repo.go` 尚未完整讀寫 migration 002 新欄位。
- 調整：補齊 Create/Get/List/Update scan 與 update 欄位。

### 3.3 queue 背壓策略可配置
- 目標：讓 `emitValue` 滿佇列時策略可選（預設 drop-oldest）並紀錄計數。

---

## 4. 可觀測與錯誤處理（P1）

### 4.1 runtime metrics（最小版）
- 新增計數：
  - collected_total
  - mapped_total
  - write_success_total
  - write_error_total
  - mapping_error_total
  - queue_drop_total

### 4.2 錯誤事件記錄
- 新增：`internal/datalink/runtime/error_sink.go`
- 先記錄到 log + 結構化事件（後續可擴 DB dead-letter table）

---

## 5. 測試與驗證（P0/P1）

### 5.1 單元測試
- runtime ingestor 正常/異常路徑
- writer sqlite/postgres 寫入行為
- settings 與 device repo 欄位一致性

### 5.2 E2E 測試（最小可用）
- 新增：`internal/datalink/runtime/e2e_pipeline_test.go`
- 模擬：
  - fake connector 回傳值
  - scheduler 輸出 -> mapping -> write
  - 驗證 timeseries 有資料 + quality 正確

### 5.3 回歸測試
- 現有 point poll / mapping preview 不受影響

---

## 6. 交付節奏

### Milestone A（先打通）
- runtime + sqlite writer + e2e（SQLite）
- 可在本機看到資料進 `timeseries`

### Milestone B（生產準備）
- postgres writer + stats/readiness 修正 + settings 單一來源

### Milestone C（穩定化）
- queue 策略 + metrics + 錯誤事件整理

---

## 7. 風險與緩解

1. **啟動流程分散**：runtime 可能未被啟動
   - 緩解：集中單一組裝點，並加啟動自檢 log。
2. **mapping 執行耗時拉高採集延遲**
   - 緩解：ingestor 與 scheduler 解耦，採批次寫入。
3. **DB 寫入失敗造成資料遺失**
   - 緩解：錯誤事件＋重試策略（下一輪可加 buffer spill）。

---

## 8. 檔案級變更清單（預計）

### 新增
- `internal/datalink/runtime/service.go`
- `internal/datalink/runtime/ingestor.go`
- `internal/datalink/runtime/error_sink.go`
- `internal/datalink/storage/sqlite_writer.go`
- `internal/datalink/storage/postgres_writer.go`
- `internal/datalink/storage/factory.go`
- `internal/datalink/runtime/e2e_pipeline_test.go`

### 修改
- `internal/datalink/device/sql_repo.go`
- `internal/datalink/api/settings_handler.go`
- 啟動組裝入口（依實際 main/router wiring）

---

## 9. 預設決策（已採用）
- 先 SQLite 打通，再接 PostgreSQL。
- queue 滿載預設：drop-oldest + 計數告警。
- 測試門檻：主鏈 E2E 1 條 + 相關測試覆蓋 80%+。

---

## 10. 進度（Milestone A）

### 已完成
- [x] 新增 `internal/datalink/runtime`，由 runtime service 消費 scheduler `ValueChannel()` 並執行 ingestion。
- [x] 新增 `internal/datalink/storage/sqlite_writer.go`，支援單筆/批次寫入 `timeseries`（SQLite）。
- [x] 新增 `internal/virtual/server/modbus/udp_server.go`，支援本機 Modbus UDP source。
- [x] 新增 `cmd/loadtest_modbus/main.go`，可建立 15 TCP + 15 UDP 本機伺服器並執行 10 分鐘壓測。
- [x] 壓測報告輸出包含：`Total Reads`、`Total Writes`、`Error Count`、`Mismatch Count`。
- [x] 新增關鍵路徑測試（runtime ingestion、sqlite writer、udp server）。

### 使用步驟（Milestone A）
1. 執行 10 分鐘本機壓測：
   ```bash
   go run ./cmd/loadtest_modbus \
     -duration=10m \
     -tcp=15 \
     -udp=15 \
     -interval-ms=500 \
     -db-path=tmp/datalink_runtime_loadtest.db
   ```
2. 結果判讀：
   - `Result: PASS` 表示本輪驗證通過。
   - 觀察報表中的 `Total Reads/Total Writes/Error Count/Mismatch Count`。
3. 驗證資料落庫：
   - `timeseries` 表應有資料，且 `Mismatch Count = 0`。

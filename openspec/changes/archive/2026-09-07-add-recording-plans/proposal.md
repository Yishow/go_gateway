# 以記錄方案取代逐欄填資料庫

## Why

同一來源需要曲線、用量、狀態與最新值，不能以一個 connector 寫入間隔或 insert/upsert 選項概括。
一般使用者應選用途與保存期限，由系統建立可擴充、可追查的資料結構。

## What Changes

- 新增版本化 RecordingPlan、逐項 stream、用途／頻率／保存策略與真實預覽。
- 預設 managed 長表，保留精確型別、設備身份、原讀取時間、品質與證據。
- 支援歷史、摘要、用量、狀態、事件、批次快照及只留最新值；不增加任意 SQL 執行入口。
- 建立資料庫能力、連線測試、真實表欄位查詢、安全建表與一次試寫 API。
- 既有自訂表、grouped row、單點寫入保留為進階相容模式，不自動遷移生產資料。

## Capabilities

### New Capabilities
- `recording-plan-storage`: 記錄模式、型別保存、stream 與 transaction 契約。
- `recording-database-setup`: 能力偵測、建表、試寫、憑證與目標驗證。

### Modified Capabilities
- `database-target-workbench`: 將 managed 方案與 legacy grouped 模式分清楚，改以已驗證能力顯示資料庫。

## Impact

依賴 add-measurement-semantics。重用 dbtarget service/repository、writer、schema 工具與 workspace ownership。
新增 plan/storage abstraction；新 managed runtime 功能必須等 change 4 可靠交付通過才能正式啟動。
既有 API 保持相容；新 API 與資料表詳見共享 contracts。不是新增另一套工作台。

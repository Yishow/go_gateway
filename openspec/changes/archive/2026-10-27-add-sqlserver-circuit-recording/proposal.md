## Why

SQL Server 現場需要與既有迴路記錄相同的瞬時／累積週期、原始值＋差值與可靠長表寫入，但畫面存在 SQL Server 選項不能代表實際具備此能力。需要獨立提案補齊 SQL Server 方言、schema、交付與分區驗收，讓各資料庫分批落地。

## What Changes

- 將 add-circuit-recording-profiles 的 canonical profiles、精確 counter 與 managed_long_v1 延伸到 SQL Server，不重作採集或差值引擎。
- 完成 sqlserver connector 的儲存、加密連線、真實 metadata、operation capabilities、明確 schema 預覽確認、Apply、production durable write 及 readback。
- 預設普通長表，可明選 SQL Server 原生每月時間分區，管理本案受管的 partition function／scheme 與對齊索引，分區持續新增但不自動刪除歷史。
- 對應 UTC datetime2(6)、精確 decimal、Unicode 及 bool；不能靜默四捨五入或把大整數經浮點轉換。
- 在同一交易寫資料列與 receipt，支援 commit ACK 遺失、重啟、離線重送、已確認受管缺表恢復與部分失敗修復。
- SQL Server 能力以可丟棄真實 instance 驗證後開放，保留既有 SQLite／PostgreSQL／MySQL 路徑。本輪依使用者 2026-10-27 指示使用已拉取的 SQL Server 2025／ARM 模擬環境執行完整功能驗收；如實記錄版本與架構，不以年份或原生 x86 作為執行門檻，不因此免除任一資料庫／瀏覽器驗收或自動開放正式能力。

- 依使用者追加授權，最小擴充共同來源：已保存、明確確認的十進位 scale 將整數 counter 精確轉成 decimal，再沿用原 counter／journal／outbox；保留原 decimal raw/delta 驗收。共同實作及證據歸 MySQL/MariaDB 批次，SQL Server 後續沿用。

## Capabilities

### New Capabilities

- `sqlserver-long-recording`: SQL Server 長表設定、精確值、schema scope、耐久交付、受管表恢復與可選原生時間分區。

### Modified Capabilities

無。繼承共同 recording-database-setup、database-output-delivery 與前置長表契約，新增資料庫專屬要求。

## Impact

- 技術前置：add-circuit-recording-profiles 已實作驗證；共同前置已封存於 35390b4。實作批次排在 add-mysql-mariadb-circuit-recording 之後，避免共用 capability/schema/UI 檔案重疊；SQL Server 不依賴 MySQL SQL 方言或分區實作。
- 本案唯一負責 SQL Server；不改寫 MySQL／MariaDB 提案，不擴充其他 SQL Server 雲端衍生服務、報表、歷史搬移或設備協議。
- 範圍：internal/datalink/measurement/、runtime/、grouppipeline/、dbtarget/、recordingplan/、workspace/、groupdelivery/、既有外部 SQL adapter 所在 internal/datalink/、必要 schema metadata migration、internal/api/handlers/、cmd/test_ui/service_wiring.go、go.mod／go.sum 中必要 driver，以及既有 Step 4/types/services/hooks/locales 與測試。
- 已查證 GetConnectorCapability 目前未將 sqlserver 視為支援；OpenDestination 只允許 SQLite／PostgreSQL。實作時核對並沿用實際 adapter，必要時使用 Microsoft go-mssqldb，不由 UI 或 dependency 名稱推定端到端能力。
- 使用者已授權按順序實作；真實資料庫驗收使用模擬來源與可丟棄 SQL Server，不接正式資料庫或執行部署。

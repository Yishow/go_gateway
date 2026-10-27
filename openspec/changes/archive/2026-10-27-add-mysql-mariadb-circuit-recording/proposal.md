## Why

迴路記錄目前提案僅涵蓋 SQLite／PostgreSQL；使用 MySQL 或 MariaDB 的現場也需要相同的瞬時值、累積原始值與差值，以及受管長表與缺表恢復。兩者可共用邏輯設計，但 DDL、分區及伺服器能力必須各別驗證，不能由相同 wire protocol 推定全部相容。

## What Changes

- 將既有 canonical circuit profiles 與 managed_long_v1 延伸到 MySQL／MariaDB，沿用瞬時 1–300 秒／預設 15 秒、累積 1–600 分鐘／預設 10 分鐘，以及原始值＋差值契約，不重作計算引擎。
- 沿用 mysql connector kind，實際探測並保存 MySQL／MariaDB family、版本與 operation capabilities；完成保存、連線、真實 metadata、schema 預覽確認、Apply、production durable delivery 與 readback。
- 預設普通 InnoDB 長表，可明確選原生每月時間分區；分區預建、catch-all 維護與受管表缺失恢復沿用確認、ownership、operation ledger，沒有自動刪除政策。
- 處理逐句 DDL 提交、部分成功與回應遺失；同一操作可查證／恢復，不能聲稱整批 rollback。
- 以精確型別、UTC microsecond、交易內 row＋receipt 及 payload digest 保證可靠重送；MySQL／MariaDB 各有真實可丟棄資料庫驗收，未驗證能力維持不可用。

- 依使用者追加授權，最小擴充共同來源：已保存、明確確認的十進位 scale 將整數 counter 精確轉成 decimal，再沿用原 counter／journal／outbox；保留原 decimal raw/delta 驗收。共同實作及證據歸 MySQL/MariaDB 批次，SQL Server 後續沿用。

## Capabilities

### New Capabilities

- `mysql-mariadb-long-recording`: MySQL／MariaDB 的已驗證長表設定、schema 操作、耐久交付、缺表恢復與可選原生時間分區。

### Modified Capabilities

無。沿用既有 recording-database-setup、database-output-delivery 及前置長表契約，以資料庫專屬能力補足行為，不重寫共同規格。

## Impact

- 必要前置：add-circuit-recording-profiles 的共同資料模型、counter、durable closure、schema scope 與 Basic UI 已實作驗證。共同前置已封存於 35390b4；本案已開始實作，各操作仍須依 tasks 驗證後才開放。
- 唯一負責 MySQL／MariaDB 新長表路徑；不移入既有 legacy writer WIP、歷史 custom-table 修復或 SQL Server 工作。add-sqlserver-circuit-recording 為另一批，安排在本案之後執行以避免共用 capability、schema 與 UI 檔案衝突；技術語意仍直接依賴共同前置。
- 範圍：internal/datalink/measurement/、runtime/、grouppipeline/、dbtarget/、recordingplan/、workspace/、groupdelivery/、必要 schema metadata migration、internal/api/handlers/、cmd/test_ui/service_wiring.go 及既有 Studio Step 4、types/services/hooks、en/zh-TW locales 和對應測試。
- 現有 capability_mysql_test.go 明確將 managed schema 視為未驗證；OpenDestination 目前只放行 SQLite／PostgreSQL。連線 driver 存在與舊 custom 模式不作為新能力完成證據。已封存 fix-activation-barrier-and-mysql-target-defects 是舊路徑背景，不能代替本案驗收。
- 本輪按本案 tasks 分批實作、使用模擬來源及可丟棄 MySQL／MariaDB。無正式資料遷移、安裝或部署。

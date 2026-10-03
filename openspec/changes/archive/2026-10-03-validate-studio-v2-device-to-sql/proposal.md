## Why

component tests、已歸檔規格與UI成功提示無法證明使用者的設定真的把設備值寫入正確SQL資料列。需要用production binary、simulated device、真實UI及可丟棄SQLite/PostgreSQL，把持久設定、交付和故障恢復串成可重跑證據。

## What Changes

- 新增production-entrypoint驗收harness，UI建立設備/點位/Tags/群組後直接查內部persisted IDs與外部SQL row，禁止以mock API作為最終證據。
- 真採集驗證既有int16、uint16、bool、uint64、float32、float64、string七型別；decimal僅列explicit codec／SQL round-trip證據，設備採集未支援，不新增或推測decimal Tag型別。
- 涵蓋多設備同位址、混合型別、缺值/bad/stale/late、DB outage、重啟、duplicate提交、unknown commit、設定編輯與concurrent revisions。
- 保存versioned machine-readable witness、執行環境、命令、SQL內容與UI截圖；每次測試前後確保可丟棄資源及自身資料清理。
- 分開已執行的Linux simulator證據與尚未執行的Windows/ARM/embedded browser/LAN/PLC/SCADA現場驗收，不用測試名稱宣稱跨平台通過。

## Capabilities

### New Capabilities

- `studio-v2-device-to-sql-acceptance`: 可重跑的production binary到SQL驗收與證據邊界

### Modified Capabilities

無。沿用既有相關契約，新增production範圍要求，不重寫其他能力。

## Impact

scripts/、frontend/tests/e2e/、cmd/test_ui/測試接縫、Go integration tests及docs驗收紀錄；前置：前五案全數驗證。只驗證本批basic grouped writing；不新增報表、retention產品或真設備部署。

已授權依A→F實作；1.1–1.3已有production witness，剩餘矩陣按tasks與validation的实际证据繼續，不以artifact完成代替產品驗收。來源、依賴與移交見 [總覽](../../../../docs/plans/studio-v2-write-groups/README.md) 及 [現況證據](../../../../docs/plans/studio-v2-write-groups/evidence.md)。

2026-10-03使用者確認維持現有型別：取代原F1.2及共用fixture的真採集decimal承諾；decimal codec契約與B既有完成證據保留，七型別真UI→SQL證據才是F1.2/F1.3驗收範圍。

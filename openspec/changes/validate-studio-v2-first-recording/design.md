## Context

既有 run.mjs 先以外部 SQL 建目標 readings，再操作 UI；兩台設備驗證是各自一群組。沿用現有 harness 與 production binary，新增缺少的情境，不否定或覆寫歷史證據。

## Goals / Non-Goals

證明從沒有採集表到持續可查資料，以及 A-C 的故障組合。不是 PLC 現場認證、效能保證或全平台部署。

## Decisions

1. 基本起點：全新 workspace、已啟動的 loopback simulator、已知點位表。SQLite 目的檔可尚不存在；PostgreSQL server/account/database/可用 schema 由環境提供，但沒有採集表。獨立 SQL 必須先證明目標為空。允許測試環境佈署，不允許預先建產品目標表。
2. 所有 setup mutation、schema 確認和 start 都透過真 UI；API 與獨立 SQL 只作觀察/斷言。fixture SQL 僅能建立非目標的隔離 neighbor 或故障控制，不可為待驗的錄製表建表或 INSERT 測試讀值。
3. 一個預設 60 秒區間的首筆測試保持產品預設；另對需要較快故障循環的案例，在 UI 明確選既有 10 秒 interval，不偷偷改 production default。成功路徑持續觀察至少三個正式 bucket。
4. 首次可用性目標：已知連線/點位、相同受控環境與準備好的 browser/backend 下，從首次打開空白 setup 到獨立 SELECT 看見第一列不超過 300 秒，SQLite／PostgreSQL 各量測三次並記錄每次結果與最慢值。這是產品驗收目標，不是現場保證；不得用預建表、直接 API 或少驗欄位換時間。
5. 固定矩陣：草稿只存不中斷；慢 sender 成功/失敗；DB 離線冷啟動新 ACK；ACK 後 Disable/Delete/supersede 重啟排空；同群組 A/B 同/不同欄；until race；延遲補送時間；開始記錄的 duplicate/lost response/stale/partial；custom 與 Share-only 不退化。
6. 七種既有採集型別與 uint64 9007199254740993 要走真設備模擬→UI→SQL。decimal 僅現有 codec/SQL 層，不新增設備採集。缺值/bad/stale/silent、不確定 commit、容量/poison/worker 等既有保護回歸不刪除。
7. 每個場景記錄期望與實際 source/group/revision、SQL、journal/outbox/checkpoint/receipt、UI stage、cleanup、命令、source/build 與執行環境。SQL committed、readback 與 cleanup 分開記錄；未跑與 blocked 不能算 PASS。

## Risks / Trade-offs

- 計時受環境影響：記錄初始化前置與環境、三次結果，不外推真設備；若未達標只修本範圍流程瓶頸，不開效能平台。
- 故障注入污染正式版：使用既有 build tag 和 loopback 隔離，驗證 normal binary 未啟用。
- 清理錯報：run-owned 資源清理後獨立查不存在，任何 cleanup/child exit failure 均不得發布 PASS。

## Migration Plan

不碰正式 DB 或部署。新增證據不覆蓋既有 evidence-f 的歷史 run；fixture 有唯一命名空間、可查清理與不變的 neighbor 證據。

## Verification

任何必驗 SQLite／PostgreSQL 情境未完成均不得稱 A-F 整體完成。各案 focused tests、正式最低檢查和本案 UI/SQL 矩陣都完成後，再按原流程 archive。Windows、原生硬體、LAN、PLC、SCADA、正式 DB、長時間 soak 等未驗範圍逐項保留。

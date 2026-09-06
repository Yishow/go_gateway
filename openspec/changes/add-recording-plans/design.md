## Context

目前 dbtarget 以 mapping 與 grouped bucket 寫出；前端仍有範例欄位來源。
新方案要支援多設備、多種資料，並重用既有 database connection 和安全契約。
歷史規格僅列 SQLite/PostgreSQL，但最新 main 有 MySQL inspection/schema 修正；須改成真實能力矩陣，
不能把 SQL Server 等 UI 選項直接視為完整可用。

## Goals / Non-Goals

目標：一般流程不用 SQL 與逐欄 mapping；資料結構能分離樣本、區間、事件與最新值。
不做：換掉既有 config DB、默默修改既有資料表、擴張所有 DB driver、以 upsert 取代歷史。

## Decisions

1. Plan 與 DB connection 分離。一個 plan 可以明確指定多個已驗證 destinations；
   同一 measurement 供多個 stream，採集與計算重用，但不同用途／目的地有獨立交付狀態。
2. 記錄模式及欄位依 contracts。raw every_sample 為預設；on_change 使用相較最後保留值的絕對死區，
   初筆、品質變更與最大心跳必留。sampled 模式標示會省略中間明細。
   非數值 state/text 只做相等性比較，不使用數值死區；counter/delta 用量計算不得從省略明細推算。
3. 摘要與計量入口在降採樣前。明細、計算證據與傳輸 batch 各自分離；
   flush 每 N 秒只改送出時機，不改觀測時間、樣本數、平均或用量。
4. managed 寫入長表加 definition ledger，關鍵索引和 unique receipt 由 schema plan 建立。
   精確十進位欄是權威值，原型別保留；近似 numeric 欄只供視覺／索引，不能代替累積運算。
   外部 schema 決定支援的數值長度，超出須 fail closed，不 truncate 或轉浮點。
5. 同一設備的一輪快照顯示完整度；非同時讀取保留 member timestamps 與最大允許偏差。
   batch_snapshot 使用明確 trigger measurement、edge/sequence、batch_id、freshness 與 timeout；
   缺項預設輸出 partial 記錄，不從舊 cache 補齊。要求完整的進階策略可拒絕該批但必須保留失敗事件。
   同一 trigger token 只產生一個 batch；重複輪詢 held-high 不重複產生。
6. latest_only 只更新比目前更晚的可信樣本，明示沒有歷史；使用者可另開歷史 stream。
   歷史追加不可覆寫；摘要修訂帶 calculation_revision，不能被晚送舊版蓋回。
7. 連線測試獨立於欄位 mapping。真實 introspection 提供 schema/table/columns、型別、
   NOT NULL/default、PK/unique、生成欄、權限；未知與查無資料必須分開，不退回範例。
8. 新建 managed 表先 preview schema diff，確認後才建立；預設不建 database。
   既有表要求確認欄位相容與必要列 key；無權限可匯出建表指令交管理者。
   managed 前綴碰到既有非本系統表時阻擋，不拿 CREATE IF NOT EXISTS 當相容證明。
9. 憑證使用後端 reference；改 endpoint/user/DB identity 即清除預覽並要求重新驗證。
   外連驗證沿用現有 TLS 與允許端點／路徑邊界；本機 SQLite 路徑位於 Gateway，不是瀏覽器電腦。
   不允許任意檔案覆寫或無確認的資料表結構更改。
10. 試寫：先顯示真實讀值、目的地、筆數與影響；確認後實際寫入帶穩定 test identity 的記錄。
    同一試寫 token 重按不得重複寫。managed 測試紀錄標記 test 並排除正式統計；
    自訂表若不能安全區別測試資料，預設不允許試寫，需專用表或明確同意副作用。
    不宣稱 rollback 能撤回 trigger/外部副作用；有查詢權限才查回驗證，否則顯示未驗證。

## Risks / Trade-offs

長表筆數較多，但混合型別、擴充項目與跨設備查詢不必反覆改表；wide view 可供現有 BI。
每個 adapter 必須證明精度／索引／unique/transaction；只測連線不代表全能力通過。
多個目的地不做分散式原子提交，狀態分別顯示，可靠性由 change 4 實作。

## Migration Plan

加法式建立 plan/stream/definition schema；現有 dbtarget mapping 保持原模式。
使用者預覽後才選擇加入新方案；不自動雙寫到同一用途導致重複。
回退停用 plan，保留外部資料與自訂表；不得執行自動 DROP TABLE。

## Open Questions

實際表名前綴、資料保留與頻率為部署設定；明確在啟用前確認。
SQLite/PostgreSQL/MySQL 是本組驗收矩陣；其他 kind 只有通過同等能力測試後才可選。

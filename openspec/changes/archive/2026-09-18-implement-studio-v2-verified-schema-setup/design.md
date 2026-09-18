# 設計：先看清楚要改什麼，再真的動手

## Context and Scope
In scope：真實 metadata、版本綁定 preview、明確 schema apply、operation 查詢、partial/unknown/retry、所有相關建表入口防旁路。
Out of scope：任意 SQL 編輯器、DROP/rename/破壞性 ALTER、正式資料搬移、資料庫管理權限自動提升。
來源定位為本案 proposal；一般 `dbTargetAPI.listTables` 已回傳欄位，優先重用而非再複製一套查詢。舊模擬欄位可留測試與新表建議，但不得進 production metadata。

## Decisions
### 1. 查詢結果不能只有一個成功燈
結果至少包含 connector_id、connector_revision、schema、table、inspection_status 與 columns。inspection_status 的伺服器契約為 `exists|missing|forbidden|failed`；`not_checked|checking` 是客戶端在未查詢與進行中的本地狀態，不會出現在回應 envelope。exists 時才把 columns 當成已查得欄位。missing 不等於 forbidden，空清單不自動代表表不存在。沒有來源權限時顯示受限，不 fallback 到範例。

查詢 key 含 workspace、connector/revision、database、schema、table。metadata envelope 必須包含 `workspace_id`、`connector_id`、`connector_revision`、`database`、`schema`、`table`、`inspection_status` 與 `columns`。目前 `listTables` 只回傳 array，並非完整 inspection 契約；需由 service／type／React Query hook 承接 scope envelope，同時相容舊 array 呼叫。改任一身分立即令舊資料失效；缺欄、無權與查詢失敗不得轉成 missing。取消要求或忽略舊答案不能宣稱已取消伺服器副作用。儲存後讀真正 backend scope，不對未保存草稿發起 target mutation。

### 2. Preview 契約
擬定要求：plan_id 或明確 custom-table 設定參照、connector_id、expected_workspace_revision、expected_connector_revision、expected_plan_revision（plan 模式必要）、table scope。preview/apply 都呼叫 workflow B 的 server target resolver，重新解析並核對 workspace、plan、connector、database、schema、table 與各 revision；dialect 由 connector 決定，若客戶端仍送 dialect 且不一致，必須拒絕。

回覆：token、operation_id、operation scope、各 revisions、生成內容 digest、statements、expires_at、可讀的新增表／欄位摘要。token 為不透明隨機值；伺服器持久保存其 digest、身分、版本、允許動作與狀態，不接受瀏覽器送回修改後的 statements。有效期預設 10 分鐘是本案沿用的設計值，不是效能測量結果。

預覽只讀，不為了探測而建表、插資料或更動 readiness。表已存在且完全相容，可回空 statements 與明確 no-change reason。

### 3. Apply 契約及入口收斂
要求包含 token、operation_id、expected revisions；提交的 token 必須屬於目前工作區且與已保存 scope 完全一致。missing fields=400、foreign/not found=404、stale/expired/另一操作佔用同範圍=409、incompatible/permission validation=422、未實作 adapter=501。資料庫連線不可用=503；真正 server exception 使用安全 500。狀態矩陣需由實作與 generated API docs 同步驗證，不能單靠文件宣稱。

operation_id 為後端建立、與同一 token 綁定的一次明確確認穩定編號，preview 回覆時提供；客戶端不得任意換號繞過同一 token 的唯一操作。重複確認不得重複執行；同 token／operation 處理中回 202 與同一 operation_id，終態回 200 及保存的結果；只有另一 operation 佔用同一範圍才回 409，可透過擬新增的 `GET /studio-v2/workspace/database-operations/:operation_id` 查詢。讀取也做工作區歸屬檢查。此路徑相對於現有 datalink API base。

operation ledger 由 recordingplan 的 SQLRepository、MemoryRepository 與 types 共同定義，並配合本地 migration。MemoryRepository 只驗證介面語意；多程序與 restart 的不重複 mutation 必須用 durable SQL 實測。status handler 對 foreign／unknown operation 回安全 404，running／partial／unknown 均可查詢。

`StudioV2WorkspaceDatabaseHandler.GenerateSchema`、`DatabaseTargetHandler.GenerateSchema`（`internal/api/handlers/dbtarget_handler_tooling.go`）及其他公開的非預覽建表入口也必須接相同確認 gate，不能保留 dry_run=false 的旁路。preview 舊入口可相容；無 token 的 mutation 客戶端得到明確升級錯誤。`EnsureWorkspaceSchema` 轉成僅核對／準備狀態；activation 不暗中建立或修改表。Local Modbus-only 不要求 DB schema。

### 4. 操作狀態與資料庫差異
operation 狀態為 `pending|running|succeeded|partial|failed|unknown`。結果包含 executed statements 數、查核的 schema digest、安全原因及下一步；succeeded 需要實際執行後核對，no-op 也需查核既有表相容。

對能真正保證本批 DDL 一起完成的 adapter 使用目標端交易。不能保證時，執行前須有部分完成政策；逐項記錄已執行與已核對範圍，禁止假稱全部還原。此案 release 最低實測 SQLite、PostgreSQL；PostgreSQL 在本案指外部 target adapter，本地設定仍依既有 repository，不新增未證明的 PostgreSQL local config repository。MySQL 必須獨立驗證其 DDL 行為才能啟用 managed schema。未驗證不等於刪除原有 custom-table 功能。SQL Server/其他種類不靠下拉選單即宣稱支援。

資料庫已變更但本機 acknowledgement 未存妥，回 unknown，靠持久 operation 與真實 metadata 核對，不重新盲跑。無法證明剩餘語句安全時停在 partial/unknown，由新預覽確認修復，不刪除使用者資料。不得將外部多資料庫行為包裝成單一交易。

### 5. 並行與安全
同一 connector/scope 的 schema mutation 與身分修改協調順序；採已固定的連線快照，操作前再次比對版本。操作 ledger 必須能阻止兩個程序對同 token 重複取得執行權，不能只靠前端 disabled 或單程序記憶體旗標。不要持有 Modbus 全域 projection lock 等待外部 DB。

表名、schema、prefix 以各 adapter 的識別字驗證及引用規則處理；值使用參數，不執行客戶端 SQL。預覽與錯誤不包含密碼、完整 DSN 或憑證。既有表若撞名但不相容，無副作用拒絕，給「換表名或人工檢查」動作。

## Implementation Contract

進入 3.1 前，主代理必須核對 `fix-studio-v2-database-result-truthfulness` 1.1–1.8 及 `fix-studio-v2-database-workflow` 1.2、1.3、2.1、2.2、2.4 的完成與測試證據；不可用未發布 C2 代號或 frontend 草稿替代已存 scope。該階段不依賴 workflow 2.3 的 metadata 消費者。

3.1–3.2 交付帶 connector revision／table scope 的真實 metadata，exists/missing/forbidden/failed 與新表建議分開；3.3–3.4 交付持久 preview／operation 與原子 claim，拒絕缺欄、過期、跨範圍與舊版，重複同操作保持同 identity；3.5–3.7 依 adapter 以真表及故障／重啟測試證明 apply 結果。

3.8 必須同時覆蓋 workspace handler、generic DatabaseTargetHandler 及 activation 的 EnsureWorkspaceSchema；任一公開 dry_run=false 路徑缺確認都不得 mutation。3.9 同步 Step4Database／SchemaSetupSection 的 token、revision、operation 查詢、readonly 與 stale response；3.10 完成 router 狀態矩陣、generated docs、完整基準及隔離 DB 證據。測試定位以 tasks 與 delta scenarios 為準。

範圍限 schema 與安全 metadata；試寫保持 `fix-studio-v2-database-result-truthfulness` 501，直到 `fix-studio-v2-database-workflow` 3.3／3.4 另行驗證。新增 ledger／preview migration 先採相容欄位與讀取；舊 token 缺 scope／revision／digest／期限時回安全 stale 結果要求重建，不能補猜有效性。觀察到 partial/unknown 不刪表、不重放 SQL，也不宣稱跨資料庫原子回滾。

## Rollout and Recovery
先增加 metadata/operation 結構與讀取，再接 preview、最後啟用 gated apply。未完成功能保留 `fix-studio-v2-database-result-truthfulness` 的 501。rollback 前停止新的 schema mutation，保留 operation 狀態並核對未知結果；不刪已建立的表作為程式回退。已進行的 target operation 無法被瀏覽器離頁保證取消。

## Verification
以隔離 SQLite 檔及 PostgreSQL 測試庫驗證表的實際存在、欄位、型態、相容空變更；禁止只檢查 HTTP 200。覆蓋 expired、stale、跨工作區、同 token 雙擊／跨程序並行、途中斷線、重啟、執行一半失敗、拒絕權限、撞名，以及一般建表入口旁路與 activation 不建表。所有測試資料只在專用測試 scope，不碰正式資料。

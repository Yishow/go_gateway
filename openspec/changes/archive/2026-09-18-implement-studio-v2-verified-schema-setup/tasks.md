# Implementation Tasks

前置：`fix-studio-v2-database-result-truthfulness` 1.1–1.8 及 `fix-studio-v2-database-workflow` 1.2、1.3、2.1、2.2、2.4；核對其完成與測試證據後才能開始 3.1。本案完成後才解鎖 workflow 2.3 與試寫階段；不依賴 workflow 整案完成。

## 3. 真實欄位與受確認保護的建表

全部待執行；先以會失敗的合約／整合測試固定行為，再實作，每批結束保持可建置。

- [x] 3.1 為真實欄位查詢補 exists/missing/forbidden/failed 測試；檢查 `frontend/src/services/datalink.ts` 的 listTables 與後端 adapter，查詢失敗不可回示範欄位。 對照：Capability-backed database preparation；設計「查詢結果不能只有一個成功燈」。
- [x] 3.2 [after: 3.1] 實作包含 `workspace_id`／database 的所選 scope metadata envelope 與版本資訊；目前 `listTables` 只回 array，並非完整 inspection 契約，需同步 service／type／React Query hook 並相容舊 array 呼叫；缺欄、無權與失敗不得轉成 missing。在 Step 4 以該資料取代 production sample schema，切換目標時舊回覆失效。
- [x] 3.3 [after: 3.2] 在 `internal/datalink/recordingplan/` 新增完整 preview scope、版本、digest、到期與持久 token；preview/apply 呼叫 workflow B 的 server target resolver，重驗 workspace／plan／connector／database／schema／table／revision，client dialect 不一致即拒絕；測試預覽零副作用、相容空變更及舊 token 缺保護欄位時要求重建；migration 先加相容欄位，不將舊資料補猜成有效。 對照：Revision-bound schema preview and explicit creation；設計「Preview 契約」。
- [x] 3.4 [after: 3.3] 在 `internal/datalink/recordingplan/` 的 SQLRepository／MemoryRepository／types 與本地 migrations 實作 operation ledger、工作區讀取及單次取得執行權；更新既有 `internal/api/router_studio_v2_recording_routes.go`，新增 status handler 與 `GET /studio-v2/workspace/database-operations/:operation_id`，foreign／unknown 回安全 404，running／partial／unknown 可查；以後端發給同 token 的固定 operation_id，測同 token 多請求／換號、同操作 running=202、終態=200、別操作佔用=409。MemoryRepository 只驗介面語意，多程序／restart 不重複 mutation 用 durable SQL 實測。 對照：Durable schema operation identity and honest recovery；設計「並行與安全」。
- [x] 3.5 [after: 3.4] 在 SQLite 路徑接真正 apply 與執行後 schema 核對；先測 stale/expired/foreign，再開放合法建立。
- [x] 3.6 [after: 3.5] 在 PostgreSQL 路徑接實際 apply，測試權限不足、既有表相容與不相容，依 adapter 保證記錄 transaction 結果。 對照：設計「操作狀態與資料庫差異」。
- [x] 3.7 [after: 3.6] 補 MySQL 的能力判定與 partial/unknown 案例；未有實測證據時 managed schema 維持 unavailable，不影響原有 custom-table 路徑。
- [x] 3.8 [after: 3.7] 將 `internal/api/handlers/studio_v2_workspace_database_handler.go` 的一般建表／EnsureWorkspaceSchema 與 `internal/api/handlers/dbtarget_handler_tooling.go` 的 DatabaseTargetHandler.GenerateSchema 收斂到相同確認／只讀檢查政策；以真正 router 分別測 workspace／generic dry_run=false 無 token 均拒絕、有有效確認才能建表、activation 不建表且 Local Modbus-only 不要求 DB；保留已有安全查詢功能。 對照：設計「Apply 契約及入口收斂」。
- [x] 3.9 [after: 3.8] 在 `frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx` 接 scope-bound preview/apply 與 operation 查詢，顯示 partial/unknown 並維持同 operation_id 重試。
- [x] 3.10 [after: 3.9] 同步前後端型別、API 文件與語系，跑隔離 DB 整合、重啟／故障注入及完整驗證；具體覆蓋第二句 DDL 前後失敗、target 成功但 local ack 失敗、程序重啟與 response lost，核對各 adapter 實際 schema、operation 狀態及無重放；記錄 adapter 分別支援範圍，`fix-studio-v2-database-workflow` 3.3／3.4 完成前不解鎖試寫。

# Implementation Tasks

前置：C1、C2。全部待執行；先以會失敗的合約／整合測試固定行為，再實作，每批結束保持可建置。

- [ ] 3.1 為真實欄位查詢補 exists/missing/forbidden/failed 測試；檢查 `frontend/src/services/datalink.ts` 的 listTables 與後端 adapter，查詢失敗不可回示範欄位。
- [ ] 3.2 [after: 3.1] 實作所選 scope 的 metadata 回覆與版本資訊；在 Step 4 以該資料取代 production sample schema，切換目標時舊回覆失效。
- [ ] 3.3 [after: 3.2] 在 `internal/datalink/recordingplan/` 新增完整 preview scope、版本、digest、到期與持久 token；測試預覽零副作用與空變更。
- [ ] 3.4 [after: 3.3] 實作 operation ledger、工作區讀取及單次取得執行權；用同 token 多請求與程序重啟測試證明不重複 mutation。
- [ ] 3.5 [after: 3.4] 在 SQLite 路徑接真正 apply 與執行後 schema 核對；先測 stale/expired/foreign，再開放合法建立。
- [ ] 3.6 [after: 3.5] 在 PostgreSQL 路徑接實際 apply，測試權限不足、既有表相容與不相容，依 adapter 保證記錄 transaction 結果。
- [ ] 3.7 [after: 3.6] 補 MySQL 的能力判定與 partial/unknown 案例；未有實測證據時 managed schema 維持 unavailable，不影響原有 custom-table 路徑。
- [ ] 3.8 [after: 3.7] 將 `internal/api/handlers/studio_v2_workspace_database_handler.go` 的一般建表與 EnsureWorkspaceSchema 收斂到相同確認／只讀檢查政策；測試 tokenless mutation 旁路與 activation 暗中建表均被阻止。
- [ ] 3.9 [after: 3.8] 在 `frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx` 接 scope-bound preview/apply 與 operation 查詢，顯示 partial/unknown 並維持同 operation_id 重試。
- [ ] 3.10 [after: 3.9] 同步前後端型別、API 文件與語系，跑隔離 DB 整合、重啟／故障注入及完整驗證；記錄 adapter 分別支援範圍，C4 完成前不解鎖試寫。

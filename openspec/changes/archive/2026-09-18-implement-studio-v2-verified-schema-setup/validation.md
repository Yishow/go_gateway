# 驗證紀錄：implement-studio-v2-verified-schema-setup

## 2026-09-16 3.1–3.10 實作與驗證

### 各任務落地重點

- **3.1–3.2**：`dbtarget.InspectTable` 回報 exists／missing／forbidden／failed，只有 exists 帶實際欄位；新增 `GET /studio-v2/workspace/database-metadata`，未綁定回 404、缺 `expected_connector_revision` 回 400、版本過期回 409、連線停用回 422。Step 4 只用實查欄位配對，不再以示範欄位猜測。
- **3.3**：預覽綁定工作區、方案、連線、資料庫、schema、表名前綴與各版本，並帶 operation 編號、內容摘要與 10 分鐘期限，持久保存（migration 022）。預覽只讀；表已存在且相容回空語句與 `schema_already_compatible`；表名前綴以識別字規則驗證（舊版會把 `gw;drop table x;` 直接拼進 DDL）。舊 token 缺保護欄位一律要求重新預覽。
- **3.4**：operation 帳本（migration 023）以 operation_id 主鍵、token 唯一索引與「pending/running 範圍唯一」的部分索引，用單一條件式寫入取得執行權。重複確認：處理中 202、終態 200 回保留結果、別的 operation 佔用同範圍 409。新增 `GET /studio-v2/workspace/database-operations/:operation_id`，外來與不存在同為安全 404。
- **3.5**：SQLite 真正套用。確認必須帶 `operation_id` 與三個版本，缺一 400；伺服器重驗版本、期限、工作區歸屬與內容摘要，並確認目標自預覽後未變動，才在單一交易內執行；執行後以實際 schema 判定 succeeded／failed／partial／unknown。結果寫不回本地時 operation 保留 running，重試不重跑。移除舊的 `ApplyManagedSchema` 旁路。
- **3.6**：PostgreSQL 真正套用。執行前把 `search_path` 綁到查核用的 schema，DDL 不會落到 `public`；整批在單一交易內，失敗即回復；42501 歸類為權限不足並記為 failed。
- **3.7**：MySQL 逐句提交 DDL，無法整批回復，因此 managed 建表在實測驗證前維持不可用，能力查詢附上原因；既有 custom-table 功能與旗標不受影響。
- **3.8**：workspace 與 generic 兩個公開建表入口的 `dry_run=false` 一律回 409 `SCHEMA_CONFIRMATION_REQUIRED`，`dry_run=true` 只讀計畫保留；`EnsureWorkspaceSchema` 改為只核對，缺表時啟用回 422 `WORKSPACE_SCHEMA_PREPARATION_REQUIRED` 且不建表；未綁定資料庫（Local Modbus-only）不受影響。
- **3.9**：前端套用改送完整確認（token、operation_id、三個版本），回傳改為 operation 並顯示 partial／unknown 與下一步，重試沿用同一個 operation_id；新增 operation 狀態查詢；`SchemaSetupSection` 不再送出未確認的建表，改為顯示可行動的拒絕訊息；13 個後端錯誤碼補進白名單與 en／zh-TW 兩個語系。
- **3.10**：移除前端已死的舊 apply 型別、解析器與服務方法；新增「結果未存妥 → 重啟後不重跑」的持久化測試；重新產生 Swagger 文件。

### Adapter 支援範圍

| Adapter | managed 建表 | 實測依據 |
| --- | --- | --- |
| SQLite | 已開放 | 預覽零副作用、相容即空變更、單一交易回復、權限與目標無法查核的分流、重啟不重跑（本機檔案） |
| PostgreSQL | 已開放 | 真實伺服器 17.6：建至指定 schema 且不外洩 public、整批回復、權限不足（USAGE 無 CREATE 角色）、既有表相容與不相容 |
| MySQL | 維持不可用 | 無實測證據；逐句提交無法整批回復，能力查詢附原因，既有 custom-table 不受影響 |

試寫維持 501，待 `fix-studio-v2-database-workflow` 3.3／3.4 另行驗證後才解鎖。

### 驗證結果

| 結果 | 檢查 |
| --- | --- |
| PASS | `gofmt`、`go vet ./internal/...`、`golangci-lint run ./internal/...`（0 issues） |
| PASS | `go test ./...`（未設 POSTGRES_DSN）：44 packages ok、6 無測試、0 失敗 |
| PASS | 帶 `POSTGRES_DSN` 序列執行 `go test -p 1 -count=1 ./internal/datalink/dbtarget/ ./internal/api/`：PostgreSQL 執行器 3 項與路由 3 項實際執行並通過 |
| PASS | 前端 `tsc -b`、`npm run lint`、`npx vitest run`（168 files／962 tests） |
| PASS | `make check-lines`、`git diff --check` |
| PASS | 37 項 mutation 檢查全數被測試擋下，檔案逐一還原為逐位元相同 |
| PASS | `make gen-docs` 重新產生 Swagger：新增 `database-operations/{operation_id}` 與本案錯誤碼／確認欄位，差異未夾帶無關項目 |
| NOT RUN | 真實 MySQL；瀏覽器／E2E；跨程序並行建表的實機驗證；現場 PLC |
| NOT RUN | 未 commit、未部署 |

### 已知事項

- 帶 `POSTGRES_DSN` 全平行執行 `go test ./...` 時，本案的 PostgreSQL 測試會與 repo 既有的 `TestPostgresPartitionMigration` 爭用同一個資料庫而逾時；該既有測試單獨執行亦失敗，且 `internal/datalink/storage/` 與初始 migration 未被本案修改。PostgreSQL 驗證以序列執行為準。
- 本案的 PostgreSQL 測試改為在清理時一併刪除可能落在 `public` 的探測表：先前 3.6 的 mutation 檢查曾讓 DDL 落到 `public` 並留下 `gw_record_a`，導致後續執行失敗；該殘留已清除。
- `swag` 由 `go install ...@v1.16.6` 取得，但執行檔自報 v1.16.4；`go.mod`／`go.sum` 未被改動。
- 能力查詢的 `supports_managed_schema` 原對 SQLite／PostgreSQL 硬編 false；2026-09-16 修正為由 `ManagedSchemaExecutionVerified` 推導，與 `ApplySchemaPreview` 的 adapter gate 同源，capability 矩陣不再與實際接受的執行矛盾（SQLite／PostgreSQL 顯示 true，MySQL 維持 false 並附原因）。`supports_test_writes` 維持 false，直到試寫端對端接線。
- 2026-09-16 修補（誠實回復）：ack 寫入失敗或執行中斷線後停在 running 的 schema operation，重試原會永久回 202 並鎖住 scope；現以 `SchemaOperationLease`（10 分鐘設計值，對齊預覽期限）判定 claim 過期，過期後由目標實際 schema 證據收斂為 succeeded（`recovered_from_target_evidence`＋digest）或 unknown（`verification_unavailable`＋下一步），不重新盲跑。證據：`schema_apply_recovery_test.go`、`router_studio_v2_recording_schema_apply_recovery_test.go`。
- 2026-09-16 修補（憑證不外洩）：SQLite 目標以 `dsn` 鍵保存連線字串時，preview scope 與 metadata 的 `database` 欄位改經 `sanitizeDSNValue` 清除 userinfo 與 `password=` 片段，符合「預覽與錯誤不包含密碼、完整 DSN」。證據：`studio_v2_workspace_recording_target_database_test.go`。

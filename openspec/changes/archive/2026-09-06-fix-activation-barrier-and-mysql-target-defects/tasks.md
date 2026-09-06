## 1. Studio V2 設定存檔狀態收斂（前端）

- [x] 1.1 [P] 依 design 決策「D1 存檔中狀態僅由實際發出的存檔請求設定」，讓 `onUpdateModbusShare` 與 `onReset` 只更新設定資料與清除既有 `save_error`，不再寫入 `save_state: 'saving'`；`saving` 僅由 `onSave` 排入存檔佇列時設定，並在請求結束後收斂為 `saved` 或 `save-error`。落地 spec 需求「Settings save state reflects only real save requests」。先在 `frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx` 寫出紅燈測試：編輯任一 Modbus Share 欄位後 `state.settings.modbus_share.save_state` 不為 `'saving'`、按下重置後亦然、按下儲存後依序經過 `'saving'` 與 `'saved'`。驗證：`cd frontend && npx vitest run tests/unit/workbench-v2/settings-save-state-convergence.test.tsx` 由紅轉綠。
- [x] 1.2 依 design 決策「D3 存檔擁有權失效不得吞掉終態」，把 ownership 檢查的作用範圍限縮為「是否回填伺服器回傳的 `settings_revision` 欄位資料」；`save_state` 與 `save_error` 的終態寫入、`operationError` 設定與重試動作註冊一律不受 ownership 影響。落地 spec 需求「Settings save always converges to a terminal state」。先在 `frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx` 補上紅燈測試：存檔請求在飛行中時再編輯一個一般設定欄位，該請求失敗後 `save_state` 為 `'save-error'`、`operationError` 有值且提供重試；同情境下請求成功時 `save_state` 為 `'saved'`；被取代的回應不得把舊的 `settings_revision` 寫回較新的編輯。驗證：同一測試檔由紅轉綠。
- [x] 1.3 依 design 決策「D4 連接器成功清單以生命週期而非目前狀態修剪」，讓刪除連接器時把該 id 自 `successfulAdds` 移除，並在樂觀寫回列表時排除已刪除的 id，使「新增 A → 刪除 A → 新增 B」後列表只含 B。落地 spec 需求「Deleted connectors do not reappear after later additions」。先在 `frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx` 寫出紅燈測試斷言該三步後列表長度為 1 且不含 A。驗證：`cd frontend && npx vitest run tests/unit/workbench-v2/settings-operations-race.test.tsx` 由紅轉綠。

## 2. 啟用 autosave barrier 有界化（前端）

- [x] 2.1 [P] 依 design 決策「D2 autosave settlement 等待改為有界」，讓 `useStudioV2AutosaveSettlement` 的等待具備 10 秒上限（以具名常數表示、可由參數注入供測試覆寫），逾時時以當下狀態的 barrier 快照解決承諾而非拒絕，等待者解決後自清單移除、元件卸載時清除計時器，使殘留 `saving` 情境下 `activateWorkspace` 走既有 `modbus_share_save_incomplete` 失敗路徑且 Step 4 離開 `activating`。落地 spec 需求「Activation autosave settlement is time-bounded」。先在新增的 `frontend/tests/unit/workbench-v2/autosave-settlement-timeout.test.ts` 寫出紅燈測試：無存檔進行中時立即解決；注入短逾時且存在永不收斂的 `saving` 時，承諾在逾時後解決且 `pending_saves >= 1`；逾時解決後狀態再變更不得二次解決。並於 `frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx` 補上斷言：該情境下啟用結果碼為 `modbus_share_save_incomplete` 且 `phase` 回到非 `activating` 的終態。驗證：`cd frontend && npx vitest run tests/unit/workbench-v2/autosave-settlement-timeout.test.ts tests/unit/workbench-v2/step4-first-activation.test.tsx` 由紅轉綠。

## 3. 裝置站號預設值正規化（前端）

- [x] 3.1 [P] 依 design 決策「D5 站號預設值在載入正規化階段寫入設定」，讓裝置載入正規化流程沿用 `protocols.ts` 既有協議預設值把缺漏的 `station_no` 寫入裝置設定，並讓 `ConnectionConfigForm` 的站號欄位直接反映設定值、缺值時顯示為空，使畫面顯示與 `isStudioV2DeviceValid` 判定一致。落地 spec 需求「Displayed connection defaults are backed by device configuration」。先在 `frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx` 寫出紅燈測試：載入一台 `protocol: 'mc_3e'`、設定僅含 host/port/timeout 的裝置後，其 `config.station_no` 為 0、裝置有效性為有效、且載入正規化不對未編輯裝置發出 autosave 請求。驗證：`cd frontend && npx vitest run tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx` 由紅轉綠。

## 4. MySQL schema 產生正確性（後端）

- [x] 4.1 [P] 依 design 決策「D6 MySQL 唯一索引改用方言合法語法」，讓 `buildEnsureUniqueIndexStatement` 依資料庫類型輸出語句——MySQL 不含 `IF NOT EXISTS`，SQLite 與 PostgreSQL 維持既有含 `IF NOT EXISTS` 的形式——回傳值形狀不變，既有的「欄位已是主鍵或唯一鍵則不產生語句」與同批次索引名稱去重機制維持不變。落地 spec 需求「Generated schema statements use dialect-valid syntax」。先在新增的 `internal/datalink/dbtarget/tooling_service_mysql_schema_test.go` 寫出紅燈測試，逐一斷言三種 kind 的輸出語句字串。驗證：`go test ./internal/datalink/dbtarget/ -run MySQLSchema` 由紅轉綠。
- [x] 4.2 依 design 決策「D7 schema 產生器改以連接器解析預設 schema」，把 `buildSchemaGenerateStatements` 的參數由資料庫類型改為 `*schema.DatabaseConnector`，空白 `table_schema` 一律以 `defaultSchemaForConnector` 解析並同步調整呼叫端，使 MySQL 連接器的映射不再產生指向 `main` 的語句、且能與 `inspectMySQLTables` 依真實資料庫名稱建立的資料表清單成功比對。落地 spec 需求「Schema generation resolves the default schema from the connector」。在同一測試檔補上紅燈測試：連線設定 `database=gateway_metrics` 且映射 `table_schema` 為空時，產生語句以 `gateway_metrics` 限定資料表；該表已存在於檢查結果時不產生建表語句。驗證：`go test ./internal/datalink/dbtarget/` 由紅轉綠。
- [x] 4.3 完成 4.1 與 4.2 後補上迴歸測試，證明 spec 需求「Upsert mappings require a guaranteed unique key on MySQL」成立：啟用中的 MySQL upsert 映射對「時間戳欄位既非主鍵也非唯一鍵」的既有資料表執行 schema 產生時，輸出必定包含建立該欄位唯一索引的語句；時間戳欄位已為主鍵或唯一鍵時不輸出索引語句；並斷言語句執行失敗時 `GenerateSchema` 記為失敗且保留原始資料庫錯誤訊息。驗證：`go test ./internal/datalink/dbtarget/` 全綠。

## 5. MySQL 連線安全與探測路徑（後端）

- [x] 5.1 [P] 依 design 決策「D8 MySQL cleartext 認證改為失敗關閉」，讓 MySQL 連線描述組裝在未指定時不啟用 cleartext 認證，並依「明確 `tls` 值優先 → 要求 `use_tls` 時採加密且不退回明文 → 明確啟用 cleartext 時採加密且不退回明文 → 其餘不強制 TLS」四條規則決定 TLS 模式，且 `AllowFallbackToPlaintext` 不再無條件為 true；正式連線與探測用管理連線共用同一組解析，不各自複製。落地 spec 需求「MySQL connections fail closed rather than sending cleartext passwords」。同步改寫 `internal/datalink/dbtarget/service_mysql_test.go` 中斷言預設啟用 cleartext 與 `preferred` 的既有案例，並補上四條規則的 DSN 解析斷言與「兩處組裝結果一致」的斷言。驗證：`go test ./internal/datalink/dbtarget/` 全綠，且以 `mysqldriver.ParseDSN` 解析的結果符合上表四種組合。
- [x] 5.2 依 design 決策「D10 移除探測路徑的重複自動建庫分支」，移除 `probeConnector` 中針對 MySQL 的第二段建庫重試分支（該路徑已由 `openExternalDBManagerWithMySQLDatabaseEnsure` 內建處理），PostgreSQL 分支維持不變，使一次探測對 MySQL 最多嘗試建庫一次。落地 spec 需求「Connector probing attempts database creation once per kind」。以既有的建庫函式替身計數斷言：資料庫不存在時建庫被呼叫一次、連線重試一次；非「資料庫不存在」的失敗不觸發建庫且回報原始失敗分類。驗證：`go test ./internal/datalink/dbtarget/` 全綠。

## 6. Share 啟用前置檢查範圍（後端）

- [x] 6.1 [P] 依 design 決策「D9 Share 停用時保留工作區層級前置檢查」，把 barrier validator 改為三分支：Share 啟用時維持完整檢查順序；Share 停用且 hydration 就緒時仍檢查 readiness token 與 workspace revision、略過 settings revision；Share 停用且 hydration 未就緒時放行。落地 spec 需求「Disabled Share preserves workspace-level activation gates」。在 `internal/api/router_modbus_share_test.go` 中為三個分支各補一個案例（停用且就緒時錯誤 token 被拒為 `modbus_share_save_incomplete`、停用且就緒時錯誤 workspace revision 被拒為 `modbus_share_revision_conflict`、停用且未就緒時放行且不呼叫 ShareRestore），並確認既有的 `TestConfigureModbusShareActivation_DisabledShareSkipsRestoreAndActivatesWorkspace` 維持 hydration 未就緒的前提。驗證：`go test ./internal/api/ -run ModbusShareActivation` 全綠。
- [x] 6.2 依 spec 需求「Share settings read failure is not reported as disabled」，讓 `GetSettings` 回傳錯誤時 barrier validator 回報可重試的內部前置檢查失敗，不再套用 Share 已停用的錯誤碼，也不從未填值的 settings 讀取 revision。在 `internal/api/router_modbus_share_test.go` 以會回錯的 Share 設定來源斷言回應錯誤碼不是停用碼、且標記為可重試。驗證：`go test ./internal/api/` 全綠。

## 7. 整體驗證與回報

- [x] 7.1 執行前端完整驗證基準（依 `AGENTS.md` 所列命令，至少包含 `cd frontend && npx vitest run` 與型別檢查），確認本次新增與修改的測試全綠且無既有測試迴歸。驗證：命令結束碼為 0，輸出無 failed 案例。
- [x] 7.2 執行後端完整驗證基準（依 `AGENTS.md` 所列命令，至少包含 `go build ./...`、`go vet ./...`、`go test ./internal/datalink/dbtarget/... ./internal/api/...`），確認建置、靜態檢查與測試全數通過。驗證：三個命令結束碼皆為 0。
- [x] 7.3 執行 `git diff --check` 確認無空白與格式異常；確認本次改動未觸及 `docs/technical/studio-surface-inventory/`，若觸及則以 `go run ./cmd/studio_inventory_changelog` 補上含 `summary`、`surface`、`files`、`reason` 的 changelog 紀錄。最後產出完成回報，內容必含修改檔案、驗證結果、風險、後續建議四項，並明列未能於當前環境執行的項目與對應風險。驗證：`git diff --check` 無輸出，回報四項齊備。

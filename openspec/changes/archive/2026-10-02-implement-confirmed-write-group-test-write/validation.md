# 實作驗證紀錄

2026-10-02：依總覽 A→F 順序實作。前置 A、B、C 已驗證並歸檔。本文件只記錄 D 實際執行的命令與結果；D 完成**不代表** E 的 UI 整合或 F 的 production UI→SQL 驗收已完成，Step 4 的試寫按鈕仍由 `supports_test_writes=false` 關閉。

## 設計決定與偏離（先於實作寫明）

- **共用帳本**：試寫是既有 `managed_schema_operations` 中 `action=test_write` 的 operation，token 存在既有 `managed_schema_preview_tokens`（action、operation_id、revisions、digest 都用既有欄位；試寫內容 digest 放在 `statements[0]`）。沒有新增 repository。`026_operation_test_write_sqlite.up.sql` 只 additive 新增 5 欄，舊 row 與建表流程不變（既有建表 token 對試寫 route 回 `ErrPreviewTokenKind`；試寫 token 對建表 apply 回 legacy／stale）。
- **擁有權決定**：試寫 row 的 entity key 欄位值為 `gw-test-<operation uuid>`，讀回與清理只用這個完全相同的值（`dbtarget.OwnedRowRef` 會拒絕不在 `gw-test-` 命名空間內的值）。group 沒有 `entity_key_column` 時無法證明哪一列是試寫 row，preview 回 422，不做任何寫入或廣泛 DELETE（design 要求「無法安全標識就拒絕」）。**偏離設計文字的一點**：design 寫「共用 production codec/sender」。實作共用 production **row layout 與 insert path**（`GroupRowLayout.EncodeRow`、`InsertGroupRow`、`ClassifyInsertError`、dedupe strategy、`RowPayloadDigest`），但**不經 delivery outbox／Sender**：試寫要等 target commit 才有結果，且不應進 group backlog／quota／partition 順序。
- **舊 plan 路由**：解析只用 group 的 migration provenance（`SourceKind=recording-plan`）。recording plan 的 migration review 目前仍是 blocked，沒有 production 流程寫入這種 provenance，所以今天舊的 `recording-plans/test-write-preview`／`test-write` 實際上都回 422 `RECORDING_TEST_WRITE_PLAN_UNRESOLVED`。這是誠實的結果，不是已涵蓋舊 plan。
- **capability**：`supports_test_writes` 刻意維持 false（它把關 Step 4 的舊按鈕，那個按鈕送裸 plan_id，現在會得到 400）；新增 `supports_group_test_writes`（SQLite／PostgreSQL true、MySQL／其他 false）只表示後端能力。UI 整合屬 E。
- 501 `RECORDING_TEST_WRITE_NOT_IMPLEMENTED` 已隨完整測試通過移除；原本釘住 501 的 6 個測試／斷言（recording plans 的 contract／target witness／capability／handler 測試、swagger 測試）依 spec MODIFIED 的新契約**有意改寫**為：裸 `plan_id` 回 400 validation、swagger 記載 200／202、capability 新旗標。這是契約改變，不是為了讓測試通過而放寬預期；舊測試的「無 repo mutation」斷言保留。

## 1.1 TestWritePreviewIsReadOnlyAndScoped

- `recordingplan.Service.PrepareTestWritePreview`／`ValidateTestWriteToken`／`ResolveTestWriteReplay`／`ClaimTestWrite`：token 綁 action、workspace／group／connector revisions、dialect／database／schema／table 與內容 digest。錯 kind → `ErrPreviewTokenKind`（422）、foreign／unknown → not found（404）、stale（任一 revision／表／內容變更）與首次 claim 的 expired → 409。
- `grouptestwrite.Service.Preview`：讀已保存 group、connector（revision 必須吻合）、真實 inspected 表與 tag 型別；以 `GroupRowLayout` 預先編碼一次，欄位型別放不下測試值（`test-value-blocked`）在 preview 就被拒；回傳每欄值、擁有者欄位與值、dedupe、清理描述。**不寫 target、不 claim operation**（測試斷言：未開任何 destination 連線、表內 row 數不變、帳本沒有 operation）。
- 不支援原因（皆為安全代碼）：`test-ownership-unsupported`、`table-missing`、`table-unavailable`、`layout-blocked`（附 layout issue code）、`receipt-table-missing`、`tag-unavailable`、`tag-type-unsupported`、`test-value-blocked`、`destination-*`。
- 舊 plan adapter：`workspace.WriteGroupService.ResolveRecordingPlanGroup` 只看 provenance，0 或多個候選都回 `ErrRecordingPlanGroupUnresolved`（不猜第一筆）。

### 實際執行（1.1）

| 命令／證據 | 結果 |
| --- | --- |
| RED：`TestWritePreviewIsBoundToActionScopeAndContent` 等 5 個 token 測試先於實作 | 編譯失敗（`undefined: TestWriteScope`）；實作後 PASS |
| `go test ./internal/datalink/recordingplan` | PASS（含 `TestWriteTokensNeverCrossKinds`：建表 token 不能試寫、試寫 token 不能建表 apply） |
| `go test ./internal/datalink/grouptestwrite -run 'TestWritePreview'` | PASS（唯讀、不支援的擁有權／缺表／缺欄／缺 receipt 表／connector revision 已變／unknown group） |
| 路由層 `TestWritePreviewRouteIsScopedToTheWorkspaceAndGroup`、`TestWriteRoutesMapServiceErrorsToTypedStatuses`、`TestWriteLegacyPlanRoutesResolveToTheSameGroupOrRefuse` | PASS；workspace 身分只來自伺服器，driver 錯誤文字不出現在回應 |

## 1.2 AtomicTestWriteClaim

- 帳本新增 `TakeOverSchemaOperation`（owner 換手、僅限 lease 已過期且仍 active）與 `SaveSchemaOperationProgress`（owner 才能存進度，同時續約 lease）；SQL 與 memory 實作行為相同並有同一組測試。claim 沿用既有 unique token／operation／active scope（同表）索引：同 operation 執行中重送 → 202；已保存結果 → 200（即使 token 已過期，不再執行）；同表另一個 operation → 409 並回傳佔用者的 operation_id；stale／expired 首次 claim → 409。
- 與 schema operation 互斥：scope key 以表名取代前綴，所以只在**同一張表的試寫之間**互斥；試寫與建表（前綴範圍）不共 scope，這是限制，不是已涵蓋。

### 實際執行（1.2）

| 命令／證據 | 結果 |
| --- | --- |
| RED：`TestOperationLedgerKeepsTestWriteFactsForSQLAndMemory`、`TestOperationLedgerTakeOverOnlyAfterTheLeaseRanOut` | 編譯失敗（缺 `TakeOverSchemaOperation`）；實作後 PASS |
| `TestWriteClaimRunsOnceAndReplaysRunningAndSavedResults`、`TestWriteTakeOverOnlyAfterTheLeaseAndOnlyOnce` | PASS |
| `TestWriteRunningDuplicateAndBusyScope`（第一個確認卡在開連線，期間重送同 operation 與另一個 operation） | 202 與 409，只寫一次 |
| `TestWriteSafeRejectionsHappenBeforeAnyTargetMutation`：schema token、stale group／workspace revision、內容變更、expired、foreign／unknown、他人的 operation id | 全部在開任何 destination 連線前被拒 |

## 2.1 IdempotentProductionTestWrite

- 寫入前先保存進度 `writing`；以 production layout＋`InsertGroupRow`（策略依 group 的 dedupe：`receipt` 或單純 insert）寫入；等 target commit 後才有結果。commit 階段失敗（回應遺失）先以讀回判定：row 在 → 繼續驗證；不在 → `unknown`/`commit-ambiguous`（不重插）。commit 前失敗 → `failed`（`destination-rejected-row`／`destination-denied`／`destination-unavailable`）。
- 重啟（lease 過期後重送同 operation，由新 owner 接手）：`writing` 階段 target 已有擁有者 row → 直接讀回與清理；沒有 → 只有 `receipt` dedupe 以同一 effect key 重試，其他 `unknown`；`cleaning` 階段只完成清理（row 已消失也**不會**重寫）；`claimed`（從未寫入）從頭執行一次。

### 實際執行（2.1）

| 命令／證據 | 結果 |
| --- | --- |
| 真實 SQLite 目標＋ SQLite 帳本檔：`TestOwnedCleanupAndRestartResumeFromRecordedProgress`（9 個情境：commit 後死亡、清理中死亡、清理後死亡、寫入前死亡、無 dedupe 且無 row、receipt dedupe 重試、receipt 存在、group 試寫後被編輯、lease 未過期不接手） | PASS；以 `Faults.Reject` 計數 INSERT，需要「不再寫」的情境 INSERT 為 0，receipt 重試恰為 1 |
| `TestWriteOutcomesWhenTheTargetRefusesOrIsAmbiguous`：destination 以 trigger 拒絕、pre-commit 失敗（注入 reset）、離線、commit 回應遺失（`lostcommit`，row 在 → 驗證通過且只寫 1 次；row 被另一連線刪掉 → unknown）、preview 後 connector 被編輯 | PASS |
| 真實行程 `kill -9`：`TestOwnedCleanupAndRestartSurvivesARealKill`（子行程跑真實 `Confirm`，在清理的 DELETE 處卡住並寫出 ready 檔；父行程確認 target 已有測試 row 後 `SIGKILL`，將 lease 調舊，同檔案重啟確認） | PASS（無 dedupe 與 receipt dedupe 各一次）：重啟後 verified＋cleaned，INSERT 為 0，只剩 production row，再重送回 200 的保存結果 |
| 真實 PostgreSQL 16（disposable container `POSTGRES_DSN`）：`TestPostgresTestWrite*` | PASS：型別化值（含 uint64 `9007199254740993`、bool、double、text、JSONB provenance）驗證、receipt 同交易寫入與清理、兩個 commit 回應都遺失（寫入由讀回與 receipt 判定一次，清理 `unknown` 但 row 實際已刪） |

## 2.2 TypedReadbackEvidence

- `verify`：只在**恰好一列**、每個 member 欄位以 `DecodeExactValue` 還原後與寫入的 `ExactValue` 逐值相等、擁有者欄位相符、provenance JSON 結構相等，且（`receipt`）destination receipt 的 digest 吻合時才是 `written_verified`。讀取被拒／失敗、列數不是 1、值不符、receipt 不符 → `written_unverified` 與安全 reason（`readback-denied`／`readback-failed`／`readback-row-count`／`readback-mismatch`／`receipt-mismatch`）。table／time 不由前端或請求填入，來自 token 與帳本。

### 實際執行（2.2）

| 命令／證據 | 結果 |
| --- | --- |
| `TestTypedReadbackEvidenceRequiresTheWholePayloadToMatch`：同 owner 但值被 trigger 改成不同值、整數 off-by-one、provenance 不同、SELECT 被拒 | 全部 `written_unverified`＋原因；owned row 仍被清理 |
| 變異檢查：暫時把值比對改成永遠通過 | 3 個 mismatch 子測試失敗（驗證測試確實在檢查值） |
| 真實 PostgreSQL：以 `GRANT INSERT`-only 的角色執行 | `written_unverified`／`readback-denied`；清理 `failed`／`cleanup-denied`，row 仍在且如實回報 |

## 2.3 OwnedCleanupAndRestart

- 清理只在 `gw-test-` 命名空間內以 owner 完全相等刪除，並在同一交易刪除該 effect key 的 receipt；receipt 刪除失敗時整筆回滾（row 不會先消失）。刪除後再讀一次確認不再有 owned row。結果：`cleaned`、`failed`（`cleanup-denied`／`cleanup-failed`／`cleanup-rows-remain`）、`unknown`（清理 commit 回應遺失，`cleanup-commit-ambiguous`）、`not_attempted`（沒有寫入）。清理失敗不否認已寫入，也不顯示已清乾淨；retry／重啟只回保存的結果。

### 實際執行（2.3）

| 命令／證據 | 結果 |
| --- | --- |
| `dbtarget` 單元測試：只刪自己的 row 與 receipt、拒絕 `line-a`／空值／`%`／`gw-test-`／`gw-test`、receipt 表缺失時回滾 | PASS |
| `TestOwnedCleanupStatusIsIndependentOfWriteVerification`：DELETE 被拒（row 仍在、回 verified＋cleanup failed，再重送回保存結果且不再寫入）、清理 commit 回應遺失（INSERT 恰 1 次） | PASS |
| `TestOwnedCleanupNeverReachesProductionNeighborsAfterRestart`：另一個 operation 的 `gw-test-someone-else` row 與 production row | 皆不受影響 |
| 變異檢查：暫時讓清理錯誤一律視為 cleaned | `delete_not_permitted` 子測試失敗 |

## 3.1 完整契約開放

- 路由：`POST write-groups/:id/test-write-preview`、`/test-write`；舊 `recording-plans/test-write-preview`、`/test-write` 轉入同一服務；`GET database-operations/:operation_id` 沿用。未接線的服務回 503 `WRITE_GROUP_TEST_WRITE_UNAVAILABLE`（不是假成功）。結果 JSON 帶 `write_outcome`、`cleanup_status`、`reason`、`cleanup_reason`，不暴露 owner 與 detail。
- Go／Swagger（`swag init` 重新產生）／TypeScript 型別、嚴格 parser（拒絕互相矛盾：running 卻有結果、status 與 outcome 不一致、failed 卻有清理、verified 卻說沒嘗試清理、preview 的擁有者值不在列出的值中）、service（`testWritePreview`／`testWrite`／`testWriteOperation`）、hooks（preview／confirm 不重試／operation 查詢）、typed error 白名單與 en／zh-TW 字串。
- 生產接線：`wireGatewayServices` 建立共用帳本與 `grouptestwrite.Service`；`DatalinkServices.WriteGroupTestWrite` 為空時路由只回 503。

### 實際執行（3.1）

| 命令 | 結果 |
| --- | --- |
| `go test ./internal/api/...` | 全 PASS，除 `TestStudioV2RecordingRoutes_PostgresApplyCreatesAndVerifiesTheManagedSchema`（有 `POSTGRES_DSN` 時 `verification_unavailable`；C 階段已在乾淨 HEAD worktree 確認同樣失敗，是既有的環境問題，與 D 無關；無 `POSTGRES_DSN` 時該測試跳過） |
| `npm run lint && npm test -- --run && npm run build`（frontend） | PASS：180 files／1054 tests；新增 `studioV2WriteGroupTestWrite.test.ts` 9 項；locale 檔維持 511 行 |

## 3.2 API／operation fault matrix 與實際 DB 證據

- `cmd/test_ui` 以**production 接線**（真實 configuration DB、真實儲存的 group 與 connector、真實 destination）經 HTTP 走完整流程，SQLite（無 dedupe、receipt）與 PostgreSQL（無 dedupe、receipt）各一次：preview 零 mutation → confirm 200 `written_verified`＋`cleaned` → 共用 operation 狀態端點 → 重送 200 不再寫入；另有不支援擁有權 422、preview 後編輯 group 409 stale、真實建表 token 422 wrong kind、他人／未知 token 404。

### 實際執行（3.2）

| 命令 | 結果 |
| --- | --- |
| `go test ./cmd/test_ui -run TestProductionTestWrite`（`POSTGRES_DSN` 指向 disposable PostgreSQL 16） | PASS（SQLite 2 個 dedupe 情境、unsafe／stale／wrong-kind、PostgreSQL 2 個 dedupe 情境） |
| `go test -race`（grouptestwrite、recordingplan、dbtarget、groupdelivery、grouppipeline、api/handlers、cmd/test_ui，含 live PostgreSQL） | PASS；測試未留下 PostgreSQL schema／role |
| `go test -p 1 -parallel 1 ./...`（無 `POSTGRES_DSN`） | 48 個 package，零失敗 |
| `golangci-lint run ./...` | 僅剩既有的 `service_probe.go:201` |
| `git diff --check`、`make check-lines` | PASS |

已知限制與風險（不是已完成）：

1. 試寫 row 會短暫出現在 production 表（entity 欄位為 `gw-test-…`）；若表有下游消費者，這一列可能被看到。preview 已揭露目標與清理方式，但沒有量測。
2. 擁有權完全依賴 entity key 欄位的值唯一；若表有 unique／check 約束使 `gw-test-…` 被拒，會以 `destination-rejected-row` 失敗並清楚回報，不會改用別的寫法。
3. `unique_key` dedupe 的 group 試寫改用無 dedupe 的 insert（A 的 row policy 沒有 record key 欄位）；commit 回應遺失時以讀回判定，判定不了就是 `unknown`。
4. 舊 plan 路由目前一律 422（沒有 provenance 的 producer）；Step 4 舊按鈕仍關閉，需 E 改用 group 流程。
5. 試寫與建表 operation 不互斥（scope 不同）；只有同表試寫之間互斥。
6. 沒有執行：Windows／ARM、MySQL、真實 PLC／LAN、完整 production binary 的 `kill -9`（只做了 service 層的真實行程終止）、UI 流程（E／F）。

## spectra-verify（2026-10-02）

無 Critical；3 個 Warning：

- **W1（已修）**：spec／design 寫「共用 production sender」，實作是共用 row layout 與 insert path 而不經 outbox（理由見上）。已把 delta spec 與 design 的文字改成與實作一致，而不是改實作。
- **W2（接受，已揭露）**：舊 plan 路由目前一律 422（沒有 provenance 的 producer）；成功解析的路徑只在單元／路由測試中以「0 個候選 → 拒絕」驗證。後續 change 需補上 provenance 來源才有意義。
- **W3（已由我方補證）**：verifier 環境沒有 `POSTGRES_DSN`，PostgreSQL 項目在它那邊是 not-run。同一份程式碼在有 disposable PostgreSQL 16 的環境下通過（見 3.2），驗證者本身未重跑，不當作它的確認。
- Suggestion：scope 限制已寫入 design；`cmd/test_ui/test_ui`（我 `go build` 產生的 74MB 二進位）已刪除；`supports_test_writes` 與 `supports_group_test_writes` 的區分已在文件說明。

## spectra-review（2026-10-02）與處置

審查範圍為近似歸屬（無驗證 base，與前面未提交的 A／B／C 工作交錯）；只審了 test-write 的後端核心，前端、`cmd/test_ui`、Swagger 與測試檔未審。1 Warning、2 Suggestion：

1. **Warning（已修）**：確認路由忽略 `:id`，用 group B 的網址帶 group A 的 token 會寫到 A 的表。現在 `Confirmation.GroupID` 帶入網址的 group（舊 plan 路由帶解析出的 group；未帶 plan_id 時以 token 為準），與 token 綁定的 group 不同回 409 `WRITE_GROUP_TEST_WRITE_PREVIEW_STALE`，連已保存結果也不經由別的 group 位址回傳，且在碰 target 前拒絕。測試：service 層（不開連線、row 數不變、保存結果也拒）、路由層（對應 409、`:id` 傳進 service）。
2. **Suggestion（已修）**：進度欄位損毀時原本當作 `claimed` 從頭寫，無 dedupe 時可能重複寫。現在視為 `writing`，由 target 證據決定，無 dedupe 且找不到 row → `unknown`。測試：損毀進度＋無 dedupe，INSERT 為 0。
3. **Suggestion（不處理）**：單次確認對 token 與 tag 的重複讀取；成本小、操作罕見，未量測也不最佳化。

重跑：`go test ./internal/datalink/grouptestwrite ./cmd/test_ui ./internal/api`（含 live PostgreSQL）PASS；對受影響套件跑 `golangci-lint` 無新增問題。review 建議的 `--base` 重跑沒有可驗證的 base，未執行。

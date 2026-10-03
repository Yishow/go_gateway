# Review 後續修補證據

日期：2026-10-04。檢視的文件 commit 為 `bf463e3884ec1c3620303a8be6cd3abcd6e37ef9`；本次程式修補基準為 `6676b2b0877a4532ea8a0f57d33824c164a3b56f`。

## 授權與交付

使用者確認「先完成 bug 與規格修補」。本次完成以下三處修補，並修正 A/D/E/F 的設計與驗收契約。D/E/F 的受管建表、開始記錄協調及空白資料庫端到端功能仍待實作；六案 50 個產品 tasks 均未勾選、未 archive。本報告記錄提交前的修補證據；提交狀態以 Git 歷史為準，本輪沒有推送或部署。

原先草稿影響 applied snapshot、離線與歷史 journal 恢復、settlement deadline、entity layout 的修補已存在於基準；本次不重做。原始 review 保留於 [review.md](review.md)，不能作為此次測試結果。

## 程式修補

| 問題與影響 | 修法 | 此次證據 |
| --- | --- | --- |
| 同語意 Disable→Apply 沿用舊 applied revision，可能沿用已退休的 intake boundary，回傳接受卻未保存新樣本 | Disabled 群組走既有新版本 Apply 路徑；一般顯示名稱修改仍保留原 applied revision | 正式 Create→Apply→Disable→Apply regression；兩 revision journal、切換、重啟與 closure regression；workspace/grouppipeline race |
| connector 預設表 A，但保存群組指向 B 時，編輯器仍查 A 的欄位，可能誤判或誤配 | metadata query 帶 group ID、group revision；後端解析目前 workspace 的持久群組與 connector revision，再 inspection 保存的表。未保存表名變更不冒充已查證欄位 | disposable SQLite 真實 A/B 表 inspection；unknown/deleted/foreign/stale/missing revision 回歸；service、hook/query cache 與未保存 scope 測試 |
| uint64 被提案為 BIGINT；完整 unsigned 範圍超過 signed 64-bit，安全的 TEXT／NUMERIC 反而可能被 UI 擋住 | SQLite／portable 提案 TEXT，PostgreSQL 提案 NUMERIC(20,0)；實際 editor、建議與欄位判定都傳入 dialect。SQLite affinity 優先判 INT，避免 CHARINT 等誤判 | proposal、compatibility、member/editor regression；PostgreSQL NUMERIC 接受、BIGINT 拒絕與 SQLite TEXT 的瀏覽器操作 |

沒有新增 pipeline 狀態、queue、token 或 ledger。SQL codec 與 backend readiness 仍是最終寫入判定；UI 相容判定不能代替實際 SQL。

### 修改檔案定位

- lifecycle：`internal/datalink/workspace/write_group_lifecycle.go`；回歸於同目錄 `write_group_reenable_test.go` 及 `internal/datalink/grouppipeline/reconcile_revision_test.go`。
- metadata：`internal/api/handlers/studio_v2_workspace_database_metadata.go` 與 `studio_v2_workspace_group_metadata_test.go`；前端 `services/studioV2WorkspaceDatabase.ts`、`hooks/datalink/useStudioV2WorkspaceDatabase.ts`、`keys.ts`、`steps/step4/useStep4TargetColumns.ts`、`writeGroup/GroupEditor.tsx`。
- uint64：`frontend/src/features/datalink/workbench-v2/state/writeGroup/columns.ts`、`proposal.ts`；同 feature 的 `steps/step4/writeGroup/GroupColumnProposal.tsx`、`GroupMemberTable.tsx`、`GroupEditor.tsx`。
- 前端回歸：`frontend/tests/unit/services/studioV2WorkspaceDatabaseMetadata.test.ts`、`hooks/useStudioV2GroupMetadata.test.tsx`；`workbench-v2/writeGroupColumns.test.ts`、`writeGroupMembers.test.tsx`、`writeGroupProposal.test.tsx`。
- 文件：本目錄 README／本報告／截圖；A/D/E/F 的 design、specs、tasks，以及 D/E proposal。B/C artifacts 保持原狀。

## 規格修補與後續流程

| 案 | 已補的契約 | 保持待實作的部分 |
| --- | --- | --- |
| A | 重新啟用包含退休完成前的快速情境，舊 accepted input 保留原身分 | A 的完整產品驗收與交付 |
| D | canonical group 映射既有 preview token 的 PlanID/PlanRevision；schema/source/layout stale、ownership 與 operation replay 保護；group-scoped metadata、完整 uint64 範圍 | group-scoped schema API、明確 UI 確認後建表、metadata/receipt 真 SQL |
| E | 持久 `(workspace_id, device_id, canonical managed role)` 重用；換目的地走 CAS 與新確認；start/recovery 固定原選取 scope，不用整個 workspace activation 冒充 | 開始協調入口、真正 scoped readiness/activation、基本四步流程收斂 |
| F | 缺檔 SQLite 只 stat、已有檔 read-only；每次 fresh namespace；setup／桶等待／delivery 分段；自動化速度與操作者觀察分開 | 真 embedded UI 空白目的地→持續 SQL、兩 DB／三次計時與完整故障矩陣 |

後續先完成 A/B 必要保護，接 D/E 的單設備 SQLite 垂直流程，立即觀察首次設定與首列 SQL。C 的多 entity 及 PostgreSQL／完整故障矩陣仍為正式 release/archive gate；初步 SQLite 回饋不必等待全部 50 tasks，也不能代替完整驗收。沒有建立另一種快速模式或另一套設定模型。

完整 uint64 驗收明列 `9007199254740993`、`9223372036854775808`、`18446744073709551615`，禁止經 JS Number／浮點中轉。本次修 UI 契約，尚未宣稱完成三值的真 PostgreSQL 讀回。

## 實際檢查

環境：macOS darwin/arm64；模組 Go 1.25.5，執行器 Go 1.27.1。Go 使用可丟棄 cache，資料庫測試使用既有 fixture／可丟棄 SQLite；沒有操作真 PLC 或正式 DB。

| 命令／檢查 | 結果與限制 |
| --- | --- |
| `go test ./...` | 重跑 PASS；首次僅 `modbusshare.TestService_Status_AfterStartAndStop` 因 5020 listener bind 失敗，單獨重跑及全套重跑通過，未修改 Modbus 程式 |
| `go vet ./...` | PASS |
| `golangci-lint run ./...` | PASS，0 issues |
| `go test -race ./internal/datalink/workspace`、`go test -race ./internal/datalink/grouppipeline` | PASS |
| lifecycle/runtime cutoff 與 closure focused suite | PASS；正式 service 與兩 revision pipeline regression 通過 |
| `go test ./internal/api/handlers ./internal/datalink/dbtarget -count=1` | PASS；包含真 SQLite metadata inspection，PostgreSQL live fixture 的 skip 不算實測通過 |
| `npm run lint` | PASS；最後 compatibility 調整另跑 affected ESLint PASS |
| `npm test -- --run` | 182 files／1105 tests PASS；最後 SQLite affinity 調整另跑 columns／proposal／members 29 tests PASS |
| `npm run build` | PASS，最後前端來源再建置 PASS；保留既有 Browserslist 與 chunk 大小提示 |
| `openspec validate --changes --strict --no-interactive` | 六案全數 PASS |
| `spectra validate <A/D/E/F> --json` | 四案 PASS，沒有 errors/warnings |
| `spectra analyze <A/D/E/F> --json` | 沒有 Critical／Warning；仍有 concrete-example 等 Suggestion，未當作產品驗收通過 |
| 最後唯讀 source／artifact review | 沒有本次引入的 P1/P2；先前缺證據檔與 E 依賴矛盾兩處已修正 |
| `git diff --check`、`make check-lines` | PASS；GroupEditor、write_group_lifecycle 為 300 行以上 warning，沒有 500 行以上 blocker |

metadata 回歸曾在修補前重現：B 得到 A 欄位，無效 group scope 仍回 200。前端 group query 參數、uint64 proposal／compatibility、實際 editor dialect 與 SQLite affinity 皆先看見失敗，再修補通過。正式 service regression 另以基準 lifecycle source 的 Go overlay 重現沿用舊 revision 的失敗，現行修補來源通過；不修改工作樹來跑 RED。舊 same-revision pipeline fake reproducer 與正式 service regression 分開；正式修法讓停用後 Apply 產生新 revision，不保留 fake 狀態作為新增 runtime fallback 的理由。

### 瀏覽器截圖

agent-browser 載入實際 `GroupEditor`、React Query hook 與 locale；metadata 使用明示的測試資料。實際操作接受建議、加入缺欄位 Tag，以及選擇不相容 BIGINT。臨時 fixture page 已移除，瀏覽器與 dev server 已關閉；未執行 save、Apply、建表或 SQL 寫入。

- [PostgreSQL：NUMERIC(20,0) 可指定，缺欄位提案同型別](evidence/review-repairs-uint64-postgres.png)
- [SQLite：TEXT 可指定，缺欄位提案 TEXT](evidence/review-repairs-uint64-sqlite.png)
- [BIGINT：uint64 顯示不相容](evidence/review-repairs-uint64-reject-bigint.png)

這些截圖是元件 UI 證據，不是 F 的空白資料庫端到端驗收，也不是現場驗收。

## 未執行與範圍外發現

- `POSTGRES_DSN` 未設定，兩個 exact codec live PostgreSQL tests 明確 SKIP；PG declaration／rounding gate 測試有通過，不能替代 live readback。
- 未執行真 embedded UI→空白資料庫、完整 Playwright 產品 e2e、三個持續桶／兩 DB 計時矩陣、Windows／ARM、PLC/LAN 或部署驗收。這些屬後續 D/E/F 或現場證據。
- 額外 source finding：現有 int16/int32/uint16 的 PostgreSQL proposal 仍為 INTEGER，而 exact type 將小整數歸入 ExactInt64/ExactUint64，PG codec 不接受 INTEGER。這不是本次 uint64 修改引入；留給 D task 2.2 的完整型別一致性修補與回歸，未擅自擴大本輪修法。

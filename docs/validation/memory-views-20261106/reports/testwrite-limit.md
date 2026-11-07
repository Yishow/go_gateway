> 歷史報告閱覽副本：正規化本機絕對路徑與附件連結；保留當時結論與未驗項目。原始位元組及完整附件見 `local-only/original-evidence.tar.gz`，不要將歷史 PASS 當成新版重驗。

# 試寫 422 唯讀分析報告

## 結論

- 本案是目前試寫能力的設計限制：SQLite mixed 寬表群組沒有 `row_policy.entity_key_column`，現行試寫計畫因此在產生預覽前安全拒絕。422 本身符合「無法證明測試資料 ownership 時拒絕」的規格；已查規格**沒有明文規定所有 mixed 群組一律禁止試寫**，也沒有找到本案 mixed 寬表必須支援試寫的明文要求。因此不把這個 unsupported 自動定為必修 bug。
- 正常 mixed 記錄與試寫是不同通道。程式有 mixed 寬表的正式 runtime boundary 和 row codec；UI 的套用條件沒有把試寫成功列為條件。這次 422 發生在建立試寫 token 之前。**但現有四份擷取只能證明 schema apply 成功與試寫預覽被拒，不能證明該群組的正常記錄已啟動、SQL 已交付或讀回成功**；正常主線仍以 t144 的實際 E2E 證據為準。
- 現行簡化 UI 新建的 mixed 群組沒有可填的實體鍵設定，不能建議靠補欄位或切換長表繞過本案。進階 UI 存在非 mixed 的自訂資料表／實體鍵設定，但需另建不同群組及相容的既有資料表；這是不同記錄路徑，本次未實際操作或驗收。
- 最小建議是**修正 unsupported 的 operator 文案**，不要暗示一定是資料表本身無法辨識測試列或要求 mixed 使用者尋找不存在的欄位。是否新增 mixed 試寫能力是另一項產品決策；目前不建議為測試工具要求而新增。

## 已觀察的同一群組證據

唯讀擷取位於 `$TMP/go-gateway-full-recording-e2e-20261106/terra/fresh-fixed/`：

- `sqlite-group-draft.json`：workspace `7bb2a3c5-7a0f-45c1-927a-3966c6546624`；group `f20b3395-7867-4c1b-94b2-15cb89a56af3`，revision `fc681a45-3016-4129-ad04-ef9605829b5c`；SQLite `main.sensor_readings`，`managed_wide_profile_v1`、mixed `member-calculation-v1`，6 個 int16 來源各採 latest／15 秒；row policy 有 `due_at`、`record_id`、`group_id`、`provenance`，沒有 entity key column。
- `sqlite-connector.json`：同 workspace、connector `646ebdb4-7463-4fdd-af8b-07f49dbfd0ba` 與 revision `ea17c370-9e6d-493f-8805-57244b396caf`。
- `sqlite-schema-apply-request.json`：同 group/revision 的 schema apply HTTP 200、operation `succeeded`、執行 2 statements。這是 schema apply 證據，不等於正常記錄或試寫成功。
- `sqlite-test-preview-422.json`：同 group 的 `test-write-preview` HTTP 422，`WRITE_GROUP_TEST_WRITE_UNSUPPORTED / test-ownership-unsupported`。後端 action 建議選擇有 entity key 的資料表，對目前 mixed 編輯旅程並非可直接照做的指示。

## 呼叫鏈與判斷

1. Production wiring 將 `grouptestwrite.New` 接上 saved group、tag、destination、read-only inspector 與 operation ledger（`cmd/test_ui/service_wiring.go:198`；`cmd/test_ui/main.go:202`），route 指向 preview handler（`internal/api/router_studio_v2_write_group_routes.go:45`）。
2. Handler 的 `Preview → preview → testWriter.Preview`（`internal/api/handlers/studio_v2_workspace_write_group_test_write_handler.go:115`）呼叫 `buildPlan`；有 unsupported error 時轉成上述 422（同檔 :264）。`internal/datalink/grouptestwrite/service.go:187` 在 buildPlan 成功後才建立 preview token。
3. SQLite/PG 的 `buildPlan` 先檢查 connector revision 與實際表存在，然後把 `group.RowPolicy.EntityKeyColumn` 作為測試 ownership 欄位；空值即回 `ReasonOwnershipUnsafe`（`internal/datalink/grouptestwrite/plan.go:126`、:157、:166）。這是本案 422 的直接原因。
4. `internal/datalink/grouptestwrite/plan.go:193` 目前使用舊式 entity-keyed row layout，並未建立 mixed 的 result shape／due-at／成員計算視窗；`row` 也產生舊式 entity/bucket outcome（:306）。mixed 的正式 boundary 另外建立 versioned mixed codec（`internal/datalink/runtime/group_boundary_mixed.go:78`），其契約明確拒絕 entity-keyed mixed-wide 與不符合的 identity／window（`internal/datalink/dbtarget/group_mixed_contract.go:69`、:124、:199）。因此即使只補 entity key、略過 guard，也無法證明能正確試寫。
5. 正常記錄路徑採 mixed runtime boundary；編輯器的 `canApply` 看 saved revision、readiness 與 activation capability（`frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx:334`），試寫是旁邊獨立面板（同檔 :463）。這只證明前述程式關係，不代替實際 runtime/SQL 驗收；其他 readiness、revision 或未決操作仍可能阻擋操作。

## 規格與歷史限制

- `openspec/specs/recording-database-setup/spec.md:3825`：試寫需共用正式 row layout/codec/insert，並分別呈現 write/readback/cleanup；:3852 明文規定無法安全辨識 operation-owned 資料時，預覽應拒絕，不能擴大刪除。
- `openspec/changes/archive/2026-10-02-implement-confirmed-write-group-test-write/design.md:17`：test key 不得與正式 row 共用，無法證明 ownership 就 unsupported；這是安全設計，不是放寬 guard 的理由。
- `openspec/specs/per-member-recording-calculation/spec.md:9`、:26、:57 描述 mixed 記錄正式資料語意與 wide/long 支援；本次檢索未找到它對 mixed 試寫能力的一般保證或全面禁止。
- `docs/plans/studio-v2-flow-completion/evidence-f/fresh-recovery-manifest-20261004.json:36` 曾將無 entity key 的 basic managed group 安全拒絕列為 negative regression。該舊案例不是本案 mixed 成功／失敗的現場驗證。

## UI 可用設定與最小重現

- 簡化記錄新建時，`GroupEditor.tsx:159` 使用 `newMixedDraft`；`recording-editor/recordingMixedDraft.ts:13` 預設 wide、entity 欄空白。因 mixed 是 profile，`GroupPolicyFields.tsx:22`、:35 將「實體鍵欄位」藏在非 profile 分支。即使未儲存時可選 long（`RecordingMixedEditor.tsx:40`），SQLite/PG 現行 `buildPlan` 仍先要求 entity 欄，不能宣稱這是解法。
- 非 mixed 的可選路徑：進 Step 4「進階設定與工作區狀態」→ 新建一般寫入群組 → 選「既有資料表」→ 在「進階：列與寫入規則」填實際的「實體鍵欄位」，並配置與表內實際型別吻合的成員欄位。前端存在這些欄位（`Step4Database.tsx:291`、`GroupEditor.tsx:159`、`GroupIdentityFields.tsx:36`、`GroupPolicyFields.tsx:35`）；既有 SQLite 試寫 harness 使用此類 custom group 與 entity 欄（`internal/datalink/grouptestwrite/fixture_test.go:142`、:186）。前提是已有獲授權、相容且能以 gateway 自有測試識別安全讀寫刪除的表，並通過實際 preview／confirm；本次未用 UI 操作或驗證此路徑。不可把原 mixed 群組改稱為此路徑或直接拿其受管表冒用。
- 原案可重現序列已由既有擷取見證：建立 SQLite connector → 兩設備各 3 點 → 用簡化 UI 新建 mixed wide、latest/15 秒 → 儲存 → 明確 schema apply → 按試寫預覽 → 422。從 `buildPlan` 的分支順序推論，最小程式條件是已保存 mixed SQLite 群組、目標表存在、entity 欄空白；成員數與 15 秒非觸發此 guard 所必需，但單成員縮減**未實跑**。

## 最小建議與 ownership

此階段沒有產品檔案修改。若主代理只處理 operator UX，建議核准的最小檔案 ownership：

1. `frontend/src/i18n/locales/zh-TW/workbench-v2.json` 和 `frontend/src/i18n/locales/en/workbench-v2.json`：修正 `errors.WRITE_GROUP_TEST_WRITE_UNSUPPORTED`，以「此群組目前不支援安全試寫；請依正常記錄的準備與交付狀態確認」表達能力邊界，避免把群組驗證限制武斷歸因為表本身。現行面板透過 `getSafeErrorMessage` 顯示該本地化訊息（`GroupTestWritePanel.tsx:260`）。
2. 如主代理要求行為驗證，可在 `frontend/tests/unit/workbench-v2/recording-test-write-phases.test.tsx` 補一個 422 的顯示測試，確認不顯示成功、文案不宣稱表必有缺陷；該測試須先在原版 RED，再修改文案驗證 GREEN。這是建議，非本輪已執行。
3. API 消費者另會看到 `internal/api/handlers/studio_v2_workspace_write_group_test_write_handler.go:274` 的一般 action「選有 entity key 的表」；若本次也要修正 API 指示，需另納入此 handler 與 `internal/api/router_studio_v2_write_group_test_write_test.go` 的精確 ownership。單純修 operator DOM 不需要改後端。

若產品決策明確要求 mixed 試寫，才另協調 `internal/datalink/grouptestwrite/plan.go`、`service.go`、`execute.go` 及相鄰測試（必要時新增同目錄 mixed 測試／實作檔）的 ownership；若正式 codec 本身需改，另先協調 `internal/datalink/dbtarget/` 的具體檔案。此不是目前建議的必要修復或已獲授權實作。評估時須先用真實 mixed 寬表與既有正式 codec 寫 RED，覆蓋 preview 無外部寫入、確認後 typed readback、僅移除本次資料及 receipt、鄰近 production row/receipt 不變、重送／中斷復原與 stale scope。受管 schema 確有 `_gw_owner_…` 欄及預設 `gateway`（`internal/datalink/workspace/write_group_schema.go:169`；`internal/datalink/recordingplan/schema_group.go:325`），但這是表 ownership 線索，**不是已證明可直接拿來當 operation ownership 的方案**。mixed record identity 由 workspace/group/revision/due-at 決定（`internal/datalink/snapshot/mixed_identity.go:8`），須先證明測試 record/receipt 與正常記錄不衝突；不能只放寬 guard、直接改表或用 DDL/INSERT 繞過。

## 驗證與限制

- 已讀四份指定擷取、適用文件、相關 source/spec/test，並核對 repo 為 `7a01e7d5`、與本機 upstream tracking ref 0/0；11 個既有 WIP 保留。codebase-memory 對所查路徑回報 no recorded issue／metadata match，屬 best effort，並已讀取關鍵 source。
- **NOT RUN**：本輪沒有跑產品測試、建立新的 RED loop、操作任何瀏覽器／服務、對目標資料庫執行 SQL，或做任何產品修正。實際 422 是既有 E2E 擷取證據，縮減重現、非 mixed UI、正常記錄的啟動／SQL readback 都沒有由本輪驗證。
- 修改只有本報告；沒有 stage、commit、push、stash 或啟停服務。


## 附錄：父代理 SQLite verifier 的取時假設修正

### 結論與修改

本次修正的是父代理 harness，與產品行為無關：原 verifier 在 UI 停用後取得目前群組，錯誤要求 current saved revision、applied revision 與歷史 running revision 必須相同。實際 saved revision 已變為 `a0cbd1e7-0607-419d-bd87-46574724f578`，而三個歷史 bucket 的 frozen runtime/applied revision 是 `fc681a45-3016-4129-ad04-ef9605829b5c`。停用後的目前 saved revision 不是歷史 SQL 交付的版本權威。

僅修改 `$TMP/go-gateway-full-recording-e2e-20261106/parent-checks/verify_sqlite.py`：移除對即時 API 的依賴，以唯讀 `wg_runtime_versions` 和 `write_group_versions` 核對歷史 scope、frozen group 與 effective time，再核對當時 catalog、實際 SQLite schema、outbox 歷史 revision、目的地與內部兩份 receipt、每個 bucket 的 6 個值／provenance 和完整視窗。保留至少 3 個不同 bucket 的要求；時間以實際時間比較，避免 outbox 的奈秒格式與目標列字串格式不同而誤判。Python 最佳化模式會拒絕執行，避免忽略 assertions 而假報 PASS。

### RED → GREEN 與負向檢查

- **RED：** 對未修改的 verifier 執行記憶體測試，以設定庫的目前群組狀態替代即時 API，且攔截結果檔寫入；原第 12 行因 current saved revision 不等於歷史 applied revision 而 `AssertionError`。此替代只重現取時假設，不是重送服務請求。
- **GREEN：** 修正後以真實唯讀 SQLite 目標、設定庫與擷取檔執行，結果檔寫入攔截在記憶體；3 個完整 bucket、18 個值與雙 receipt 均通過。
- **父代理直接驗證：** 使用者回報父代理已直接執行修後 verifier，exit 0，3 個完整 bucket／18 個值／雙 receipt／frozen version 與 catalog 交叉通過；此處記錄為父代理回報，非本代理重跑服務。
- **負向：** 在記憶體中分別模擬缺少 frozen version、outbox revision 不符、內部 receipt 不符、目的地 receipt 不符、SQL 值不符及重複 bucket，全部拒絕，且未寫出 PASS；最佳化模式亦拒絕。這些模擬僅修改讀回資料，不修改資料庫。

此階段只修改上述 verifier 與本報告；未修改產品、未操作 t144 的瀏覽器或服務，也未以此次 SQLite 結果宣稱試寫通過。

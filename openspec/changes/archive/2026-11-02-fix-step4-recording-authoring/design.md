## Context

已核對 `CircuitProfileConfigFields` 的 numeric role options、`defaultCircuitProfileLocalDraft` 的 cumulative keys、`CircuitProfileBasicPanel` 的 cumulative-only error 顯示，以及 `GroupTestWritePanel` 的 preview／scope reset。現有外部 review reproduction 的四個案例是：numeric state 無 latest 選項、instant-only counter 缺名稱卻無必填訊息、啟用累積默選溫度／壓力／文字、instant-only counter 查詢錯誤沒有 retry。實作時將案例轉為 repo-relative 正式測試，不把外部臨時檔當持久依賴。整合交付順序為 backend frozen replay → existing cutover boundary → `manage-recording-lifecycle-safely` → 本 change → `guide-step4-recording-workflow`；既有前端契約上的垂直切片可先實作；saved-role／live-operation 接線與最終驗收仍須等 lifecycle interface，詳細 ownership 見「執行分期與交接」。

## Goals / Non-Goals

**Goals:**

- In scope：上述四項 authoring 缺陷、目的地布局能力早期 gate、test-write 未決操作復原；保留當前草稿、已保存語意與使用者範圍。
- 任一 disabled dependent action 都可解釋缺什麼，以及如何修正；正常讀取與純 navigation 不被不相關錯誤封鎖。

**Non-Goals:**

- Out of scope：全面 Step 4 重排、新 layout 支援、counter replay 演算法、事件觸發、資料庫→PLC、SQL 歷史清理、批次換所有設備設定。
- 不把 type／name／unit 當作 counter confirmation；不默改既有 role，不讓 UI 直接修改已接受資料的語意。

## Decisions

### 角色與計數名稱前置提示

對可表示數值的 source 提供明確 `state`（latest）、`gauge`（mean）與需語意確認的 `counter_raw`；非數值只顯示合法 role。saved role 以已保存值顯示，不因新增選項或 reload 改寫。任何未來 saved-role 變更仍須經明確支援該語意變更的 lifecycle revision／review，依 `freeze-counter-replay-semantics` 保留歷史 binding；目前已發布的布局／目的地 review 不提供 role 變更能力，依下列鎖定決策處理。

將 `requiresCounterConfirmation` 作為名稱要求與 loading/error gate 的共同條件，不再只看 cumulative_enabled。名稱被 confirmation 所需時，先顯示 inline 必填與原因並連結焦點；一般非計數的瞬時來源仍可無迴路名稱。名稱僅為 label，既有 circuit/source IDs 不因填入或更名改變。

拒絕把 numeric machine state 的 1→2 平均成 1.5，也不因名稱像「電量」就當 counter。

### 已保存角色鎖定與支援的 lifecycle 導覽

**Supersedes**: fix-step4-recording-authoring / 角色與計數名稱前置提示

本決策只取代上述決策及 pending 4.1 對「受審核 saved-role callback 已提供角色變更」的歧義假設；numeric role、名稱提示及未來語意變更必須 review／保護 frozen binding 的要求不變。這是既有 unavailable 時鎖定 fallback 的明確化，不新增 scope、不取消使用者要求，也不改已完成 task 的描述／provenance。

已核對發布基準 `2d99b9ae`：`SavedRecordingActionHost.onAction` 接受 `SavedRecordingActionIntent` 的 `edit | stop | remove | recreate | transition`，沒有 role-changing action。`Step4Database` 將 host 接至 Basic／Advanced；`RecordingChangeController.requestAction` 重新讀取並比對 lifecycle scope、revision、server action availability 與 `successor_seed`。`RecordingChangeDialog` 僅編輯名稱、布局及目的地；`changeTargetLayout` 不重建來源語意。`WriteGroupService.Lifecycle`／`lifecycleSuccessorSeed` 只提供符合 server guards 的操作資訊，並非 mutation authority。

`ValidateWriteGroupTransitionStructure`／`sameTransitionMemberMeaning` 保留 source identity、role、measurement、source/result types 與 exact recording source；`TestWriteGroupTransitionMatrixRejectsIdentityAndMeaningChanges` 的 role-change 案例回傳 `recording-meaning-change`。因此 4.1 只整合支援的 lifecycle 導覽，saved-role select 維持鎖定並說明目前不支援角色變更；布局／目的地 review 不得標示成角色編輯。不新增 role mutation、ensure/update fallback 或新的 successor 規格。缺少安全角色變更能力不是解鎖理由。

驗收保留 `TestRecordingMemberMetadataPersistence` 的不同 explicit role conflict、`TestCounterBindingChangedMeaningCannotReuseTheSameSavedRevision` 的 frozen snapshot 不變，以及 `TestCounterBindingCutoverReviewedDefinitionEditIsStale`；Basic／Advanced、reload 與 review 開關不得改 saved role 或既有 frozen binding。4.2 的正式 browser cases 與驗收範圍仍保留；實際執行／截圖只依本案「現有雲端瀏覽器驗收的條件式交付例外」處理，非沿用其他 change 的 waiver。

### 累積來源明確選擇與查詢重試

新 cumulative profile 預設 `cumulative_keys=[]`，啟用時展開／指向來源 picker，摘要顯示尚未選擇計數器。保存過的 explicit cumulative selection 原樣還原；關閉再開啟不能重設其選擇。來源總勾選不能暗中把 disabled profile 的全部來源預選。

凡本 scope 有 selected cumulative source 或 instantaneous counter_raw 且 measurement query 失敗，都顯示安全 error＋retry。retry 僅 refetch 原 workspace/device scope 的 measurement definitions，保留草稿，不 ensure group、schema apply 或 recording start；只封鎖依賴該查詢的確認／開始。

### 布局能力與不支援設定

以已保存 connector ID/revision 的能力物件產生 `layout_options` 與 reason，不把 supports_managed_long 解讀為 wide。SQLite／PostgreSQL 使用現有實作能力；MySQL／SQL Server 新建只提供經證明可用的 long，wide 明確 unavailable。刷新、error、unknown、scope mismatch 時不放行寫入；既存不支援設定仍顯示原值及修復連結，不偷偷改成 long。

`CircuitPurposeFields` 與 request builder 共用 capability gate，並保留 backend validation。UI 不以 DB driver 存在當作功能支援；storage options 亦綁原 connector revision。

### 試寫 operation 與 editor 分離

新增 `recordingTestWriteRecovery`，只保存非秘密 `workspace_id`, `group_id`, `group_revision`, `operation_id`, `operation_kind`, `created_at` 與必要可查詢的 target fingerprint；operation/token truth 留在 server ledger。多個群組的 unresolved references 分開保存，切群組不覆寫；重新載入先查 ledger，不憑 browser cache 宣稱成功。

新 preview 僅在同 group 沒有 unresolved write／cleanup 時可用；pending/running、untyped lost response、unknown write、cleanup failed/unknown 不能被 preview 清掉。允許瀏覽其他群組，但持續有「待釐清試寫」入口。若瀏覽器 reference 不可用，靠 workspace-scoped server operation listing 找回；此唯讀查詢缺口是本 change 的 backend interface 範圍。

write、readback、cleanup 分別顯示；cleanup repair 僅走已支援且 server-authorized 的同 operation repair，不重新插入測試列。absence/error 不等於未寫入；未證明可安全重送時只提供查詢與人工處理指引。

### 現有雲端瀏覽器驗收的條件式交付例外

**Supersedes**: fix-step4-recording-authoring / 已保存角色鎖定與支援的 lifecycle 導覽

只取代上述決策及本案最終驗收中「無法實際執行 browser／截圖就一律停止交付」的門檻，不取代其產品行為或其他前置。授權來源是 2026-10-28 18:12:29 UTC 明確詢問剩餘四案同一雲端限制下是否略過實際瀏覽器／截圖、明列未驗證、其餘測試照跑且不新增環境；使用者於 18:19:30 UTC 回覆「繼續」，18:19:36 UTC 已確認此解讀。適用範圍僅 `fix-step4-recording-authoring`、`guide-step4-recording-workflow`、`expose-runtime-recording-rules`、`make-runtime-logs-actionable`；不是由第三案例外自動延伸。

現有 cloud 的 browser launch、官方 bundled browser download 及 managed-browser loopback blockers 已由本輪交接確認未變。若本案驗收時實際執行仍不可用，可條件式略過實際 browser／截圖交付 gate；保留既有環境證據與可重跑命令，不新增環境、不更改安全限制、不公開服務或繞過拒絕。不因同一未變阻擋反覆重試已拒絕操作；若支援路徑已可用則依原案例執行並如實記錄結果。

正式 browser test cases、discovery、E2E 專用 typecheck、lint 及 spec scenario 對照仍須完成；逐案例及截圖欄位記錄 `NOT RUN`／`SKIP`／「未驗證」和原因，禁止標 `PASS`。discovery／編譯／jsdom／mock fixture 通過不等於真 browser、幾何／焦點或截圖驗收；不得用合成影像補位。`Browser verification exceptions remain explicit` 固定此報告契約。本案對應 pending task 4.2 只有在這些準備檢查及其餘原驗收完成後，才可依條件式交付門檻承接後續工作，並持續揭露未驗證風險；本次 artifact 更新不勾選任何任務。

產品邏輯、依賴順序、逐檔 ownership、backend DB／frozen proof、其他測試與 review 全部保留；未執行其他平台／現場仍未驗證。authoring active CLI 10 baseline 與所有已完成描述／provenance 不改；runtime CLI 4 已完成，1.5 仍須等 guided 最終交接，不能因本例外提前啟動。本例外只處理實際 browser／截圖，不能用於略過 DB、API 或安全契約。

## Implementation Contract

- `CircuitProfileConfigFields`、`CircuitProfileBasicPanel` 與 draft builders 接受同一 role/selection/required-name 決策；selected counter 的 blocked reason 不得藏在 cumulative-only 分支。
- 以測試驗證整數 state 1→2 保留 2，gauge 使用現有 mean，counter_raw 缺 definition 不送出；此處不重寫後端計算演算法。
- 4.1 的本機 SQLite acceptance 使用 `internal/api/router_studio_v2_mean_long_runtime_test.go` 的 `TestMeanLongRecording`、`newManagedDestinationFixture(t, "sqlite")`、`seedMeanSources`、`meanLongRequest`、`newMeanPipeline` 與 `readMeanRows`：由 public circuit-profile request 新建 numeric `value_role=state`，經 schema preparation／Apply 的 frozen proof、production pipeline 與 SQL readback 驗證。現有 state 10→30 latest 30 對照不是本案指定 1→2 的完成證據；以 focused case 補 state 1→2 得 integer 2、非 decimal 1.5，gauge 同窗 1→2 得 mean 1.5，並核對 state provenance。既有 focused selector 為 `go test ./internal/api -run '^TestMeanLongRecording$' -count=1`；新增獨立案例時列出其精確 selector，不能以零測試代替。
- 上述 SQLite 使用既有 test-local disposable files，不要求真設備、外部 DB 或新環境；保留既有 frozen-role rejection tests，禁止為驗收改寫歷史 proof。這不是 4.2 browser gate 的替代，也不宣稱未執行現場驗收通過。
- 新 cumulative 的溫度／文字候選不自動 selected；已保存明確來源不因修正預設而丟失。操作其他設備不改本設備草稿。
- capability gate 覆蓋 UI、builder、late response 與 connector revision change；只查看不支援舊配置不會保存任何改動。
- query failure、stale response、重複 retry、lost test-write response、reload／switch revision 都須保留對應 identity 與恢復入口。
- 正式測試使用 `frontend/tests/unit/workbench-v2/`、`frontend/tests/integration/`、`frontend/tests/e2e/`；新增 backend operation listing 時驗證 workspace ownership、safe diagnostics 與生成 API 文件。

## Risks / Trade-offs

- [新累積預設不選需要多一步] → picker 與已選摘要給明確指引，避免使用者被迫替溫度或文字確認計數語意。
- [多筆 pending operation 讓畫面雜亂] → 主面板只顯示本 scope，獨立復原清單保留其他 scope；不能為簡潔隱藏未決寫入。
- [只用 session storage 無法跨 reload/session 復原] → server ledger 為真相，local reference 僅加速定位；權限／查詢失敗顯示未確認。
- [saved role 修改與 replay 衝突] → 先完成後端 frozen binding 契約，沒有可用安全 review 時只顯示修復入口、不靜默解鎖。

## Migration Plan

Expand：加入 pure decision helpers、capability model 與 operation recovery references，保留舊已保存值。Migrate：依 role/name、selection/retry、layout gate、operation recovery 分批換呼叫端，每批維持可建置。Contract：確認所有 Basic／Advanced 入口共用 guard 後移除舊 cumulative-only gate／單一 preview-owned operation state；不清除 server ledger。

## 執行分期與交接

本次來源為已核准的「依賴先行、可獨立部分並行」，只變更執行順序，不新增產品需求、change、環境、部署或權限。保留七個 change ID 與各 spec scenario；spec.md 不變，因可觀察產品契約沒有改變。

**Supersedes:** 原先 whole-change serial execution 僅替換為 phased preparation/activation gates；不是取消 backend correctness 或最終整合前置。共同決策見 `guide-step4-recording-workflow / 分期準備與啟用門檻`。

- 最終 UI／定位／整合驗收仍依 `freeze-counter-replay-semantics` → `fix-recording-cutover-boundaries` → `manage-recording-lifecycle-safely` → `fix-step4-recording-authoring` → `guide-step4-recording-workflow` → `expose-runtime-recording-rules` → `make-runtime-logs-actionable`。唯一前移切片是 runtime 1.3–1.4 安全讀取／redaction／effective provider：在 freeze→cutover→lifecycle 交接後完成，作為 guided 2.2 的額外前置；authoring 仍為 guided 掛載前置。runtime dashboard、locator／URL 與最終驗收仍等 guided 完成；不得用準備測試取代任何整合門檻。
- 準備成果必須是可編譯、可測且範圍明確的程式；fixture／parser／未掛載元件通過只證明該準備切片，不代表 production wiring、資料庫正確性或產品驗收。每個產品切片各自 RED→GREEN→Refactor，結束時受影響測試與建置保持通過，不先提交一整套已知失敗測試。
- 每次啟動先核對目前工作樹及 owner；同一檔案只交一位 worker。需要其他 owner 檔案時先交接，不能以不同 function／相鄰行為理由並寫。共享 locale、types、routing、fixtures 也適用；不新增測試環境。
- 現行 freeze worker 保留 `internal/datalink/{workspace,runtime,counterproof,groupdelivery,snapshot,grouppipeline}/` 的修改權；後續實作僅讀其既有契約，不同時改同一檔案。backend 交接後重讀 frozen proof／revision/reason 的實際結果再接線。
- `add-windows-tray-runtime-logs` 維持獨立 active change 的 13/23 及原完成證據；broker、file sink、SSE、local-only guard 與 Windows lifecycle 不因本次分期移交或放寬。所有未執行產品／平台／現場驗收仍未完成。

## 本案檔案分工與開工條件

- 可先行：tasks 1.1、1.2、2.1、2.2、2.3，僅針對既有 `WriteGroupValueRole`、已保存 role、measurement query、connector ID/revision capabilities 做前端修正；不依賴正在改動的 counter proof 寫入。saved role 必須維持 locked，不能把新增 numeric state 選項變成 saved-role mutation。
- 前端 owner：`frontend/src/features/datalink/workbench-v2/steps/step4/{CircuitProfileConfigFields,CircuitProfileBasicPanel,CircuitPurposeFields}.tsx`、`steps/step4/circuitProfileActions.ts`、`state/circuitProfileBasic.ts`；必要新 pure guards 只放 `state/recordingAuthoring*.ts`、必要新 role 子元件只放 `steps/step4/CircuitProfileRoleOptions.tsx`。正式任務各自再縮小寫入清單。
- 測試 owner：本案新增 `frontend/tests/unit/workbench-v2/circuit-profile-authoring-*.test.tsx`／`.test.ts`；可延伸現有 `circuit-profile-actions-reduction.test.ts`、`circuit-profile-mixed-sources.test.tsx`、`device-recording-purpose*.test.tsx`、`mysql-mariadb-recording.test.tsx`、`sqlserver-recording.test.tsx`。其餘共用 fixtures 先唯讀，確需改時由 parent 獨占移交。
- authoring 完成前，guided 不得修改上述檔案、BasicRecordingPanel、Step4Database 或共享 `frontend/src/i18n/locales/{en,zh-TW}/workbench-v2.json`。authoring 需要新增文案時獨占該 locale；既有超長檔不得增行，沿既有分檔方式維持 line gate，不能藉此全域重排。guided 準備元件接受 typed display labels，不同時加同一 locale。
- tasks 3.1–3.3、4.1–4.3 的啟動門檻：freeze、cutover、lifecycle 的相關最終驗收與穩定接口已交接，parent 確認沒有共享檔案 worker。server ledger listing 由 parent 另派 backend owner，精確分派 `internal/api/handlers/studio_v2_workspace_write_group_test_write*.go`、對應 router registration 與 `internal/datalink/recordingplan/` 的必要讀取擴充；先核對實際 symbol，再固定檔案清單。前端 worker 不自行修改這些模組。
- 若 server listing 需要目前 freeze owner 的 workspace/handler shared 檔案，等待交接，不繞過權限另建假 listing。3.2 必須完成真 workspace ownership query 與 API 測試，local reference 或 fixture 不滿足本項。
- 4.1 依「已保存角色鎖定與支援的 lifecycle 導覽」消費現有 host 的支援操作；沒有 role-changing review 能力時維持鎖定及 truthful unsupported 說明，不新增 ensure/update fallback mutation。最終 browser／SQL／scenario gates 不因前端先行而縮減。

## 首個可執行任務封包：1.1

- 行為：新 numeric instantaneous mean profile 可明確選 `state` 並送出既有 `value_role=state`（每來源 latest）；`gauge` 仍 mean、`counter_raw` 仍需既有 confirmation。non-numeric 只提供合法 role，saved-role select 仍 disabled，reload/navigation 不重寫既有 role。
- 允許寫入：`frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields.tsx`、必要的新 `CircuitProfileRoleOptions.tsx`、`frontend/tests/unit/workbench-v2/circuit-profile-authoring-numeric-state.test.tsx`、`frontend/tests/unit/workbench-v2/circuit-profile-actions-reduction.test.ts`。讀取既有 `circuitProfileBasic.ts`、`circuitProfileActions.ts`、type、fixtures；若其實需要寫入，先由 parent 擴定此任務清單，不讓 worker 默默跨 scope。
- 契約：不新增 role、API 或 sample-reduction 演算法；既有 mean profile 的 state member 是 latest，payload 不能把 state 轉 gauge。UI 不計算假的 1→2 結果；前端測真 authored payload，最後整合另用 backend 實際結果核對 latest 2 而非 1.5。
- 情境：int16/float64 numeric state 選項及 authored role；gauge／counter 對照；text/bool 合法 role；saved numeric state、gauge、counter 的選項／disabled／reload 保留；選取或切 scope 不呼叫 ensure/update/start/DDL。現有 action-reduction 測試也必須通過，不改 saved proof fixtures 來配合未完成 backend。
- 方法：先新增本切片具名回歸測試並記錄實際 RED，再完成 UI 修正與 GREEN／Refactor；不預先啟用名稱／累積／retry 的其他已知失敗案例。確認所有新測試實際被執行，不能以 skip/todo/fixture-only 當完成。
- focused 命令（frontend）：`npm test -- --run tests/unit/workbench-v2/circuit-profile-authoring-numeric-state.test.tsx tests/unit/workbench-v2/circuit-profile-actions-reduction.test.ts`；再跑受影響的 `circuit-profile-mixed-sources.test.tsx`、`circuit-profile-basic.test.tsx`；對改動 source/tests 跑 `npx eslint`、`npx tsc --noEmit` 與 `npm run build`。repo root 跑 `git diff --check`、`make check-lines`。這是前端切片證據，不是新 live transition 或整案驗收。

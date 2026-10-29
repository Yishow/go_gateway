## Context

`Step4Database` 顯示 Basic、Advanced 群組及 workspace overview，卻明確不接 workspace-wide activation；`GroupLifecycleBar` 把 managed profile 生效導回 Basic scoped start。`CircuitProfileSchemaPreparation` 分開保存／schema preview-confirm，這些安全階段必須保留，但不必同時攤開所有操作與技術欄位。

本 change 是既有 React route 內的引導與資訊層，不新增設定 authority；內容依據 source 現況，而非較舊 inventory 的未接線敘述。沒有外部視覺稿，故不宣稱像素比對。

## Goals / Non-Goals

**Goals:**

- In scope：一個設備／迴路 scope、一個當前主要下一步、可讀規則摘要、直接可發現的已設定記錄管理、可展開但完整的專業細節。
- 新手能看懂「哪台設備、哪些資料、多久、寫到哪、現在缺什麼」，專業操作者能查原資料形狀、政策、revision 與交付證據。

**Non-Goals:**

- Out of scope：新建 V3／legacy Studio route、workspace-wide 一鍵全啟動、改資料庫支援矩陣、增 SQL history deletion、報表／保留政策、event-trigger、DB→PLC。
- 不以簡化移除 CAS、readiness、DDL 確認、operation ledger、source ownership、Modbus Share fail-closed，亦不另寫 counter／cutover 演算法。

## Decisions

### 單一 scope 與下一步

維持父層 canonical selected device/circuit；view 選擇不改設定、不啟動，不因該設備來源為空而自動切另一台。以 `RecordingNextAction` 的 derived state 決定 primary action：選資料、補必要語意、保存／準備、審核 DDL、開始、查詢未決操作、前往監看。未決操作優先於新開始，錯誤顯示安全原因與可行動 retry。

步驟名稱是引導，不是新的 persisted state machine；truth 仍來自既有 server queries/revisions。每個 scope 只有一個突出的 primary action，次要修復／詳情仍可用。將「保存」「建立資料表」「開始記錄」串成可讀進度但不合併成無確認的副作用。

### Owner-local 操作觀察與唯讀導覽分界

**Supersedes**: guide-step4-recording-workflow / 單一 scope 與下一步

只修正該決策在已完成隔離 caller 準備中的兩個假設：`RecordingWorkflowInput.operation` 的 `{ kind: 'none' }` 不是實際存在的全域無操作 query proof，`recordingNextAction` 不能要求不存在的 lifecycle `review_scoped_start` action。原產品需求、已完成任務描述／provenance、server revision/readiness/ownership/confirmation guards、依賴與 browser 條件式例外全部保留；本決策由既有 pending 2.3 掛載工作補正，不新增 task、backend DTO、authority、persistence format 或 retry rewrite。

真實 `WriteGroupLifecycleActionName` 只有 `edit | stop | remove | recreate | transition`。`Step4Database.reviewScopedStart` 是返回／聚焦指定 scope review 入口的 UI intent，不是 start permission。普通 Basic／Circuit 現有按鈕直接呼叫 start owner，尚無獨立確認 dialog；2.3 以重用既有 summary 的小型 presentation-only scope-review／explicit-confirm 區補足審核，不新增 start API、path 或 authority。`GroupSchemaPanel` 的 preview 是會建立 ledger/token 的 POST owner command，不是唯讀 NAV；DDL confirmation 與最終 start confirmation 仍分別交由既有 owner handler，在明確 consent 後執行。

僅提供 ephemeral、token-free 的 scoped owner observation：owner 為 `schema | start | test_write`；狀態為 `unresolved`（可有原 operationId；pending/running/unknown/lost_ack/query_failed/cleanup_pending）、`observed`（該 owner retained none 或 resolved；只帶 owner-known terminal status、opaque operation ID 與 typed reason，不複製既有 operation DTO）或 `unavailable`（loading/unreadable/invalid_record/scope_mismatch/revision_mismatch）。`observed/none` 只表示成功檢視該精確 owner 的 memory/storage 未留操作，不表示 server 全域無操作，亦不授權 mutation。實際結果與後續 readiness 仍由既有 server proof 判斷，不能持久化第二份 phase/readiness。

| Owner 觀察／條件 | Primary 與安全 NAV | Owner command 界線 |
| --- | --- | --- |
| 已知 unresolved 且有 O1 | 優先透過原 owner 查 O1；open/focus/details 可用 | 不得建立 save/preview/DDL/new start；不更換 operation ID；既有同請求 recovery 僅依下述 fresh proof 作次要明確操作 |
| unresolved 無 ID／lost ACK | 聚焦原 retained request 的查明／人工指引，不編造 ID | 不自動 resend/retry，不開新 operation |
| unavailable／loading／stale／foreign | 安全 recovery/inspection 指引；不揭露 foreign 資料；本地 NAV 保留 | 命令 unavailable，不把缺證據當 idle |
| observed owner-local none | 可選下一個 prerequisite 或開 review；不是全域 absence | 明確 command 仍須 owner 最新 scope/revision、health、readiness/capability 與 consent |
| observed resolved | 顯示原精確結果，依目前 server proof 選下一步 | succeeded 不代表 SQL committed/readback/cleanup；不可跳過既有 guards |
| offline | 本地 open/focus/repair navigation 保留，O1 不變 | network query/preview/DDL/start disabled |

保留既有同請求 retry，不把 primary 改為 query 解讀成刪除原 recovery 能力。`RecordingStartService.finishRecordingStart` 在 `internal/datalink/workspace/recording_start.go` 的既有結算路徑可回 typed `failed`／`partial` 與 `next_action: retry_same_request`；`BasicRecordingPanel.handleRetry` 重用 stored request，`useCircuitProfileStart.start(true)` 已先 refetch。2.3 僅 guarded delegate 這些既有 owner handler：每次使用者明確選擇次要 recovery 時，先取得 fresh、精確匹配 O1／owner／workspace／device／group 的 server proof，且必須為 typed failed/partial + retry_same_request，才可重送完全不變的原 request_id、payload 與 revisions；primary 仍是 query。這是原 O1 recovery，不是 NEW operation，不重建 draft／intent、不換 ID、不自動 retry。缺 ID、unknown、pending、untyped、lookup failure 或其他 next_action 均不得 POST，僅 query／人工查明；其他 owner 的 schema/test-write unresolved／unreadable guards 仍適用。原 request 與 tokens 留在 owner，不複製到 view model；不改 backend 或 retry 演算法。驗收保留 Database 路徑的既有同請求 retry 測試及精確 request／call-count assertions，不可移到 Share-only 以掩蓋功能流失。

已知 live／retained request 優先於空或錯誤 storage。`loadBasicRecordingIntent` 目前將 missing、JSON/schema/scope 無效及讀取失敗都回 undefined；`readSchemaOperation` 優先 memory fence、其次接受非空 opaque stored ID，missing group/empty storage/exception 都回 undefined，並未提供 malformed-record 驗證。若掛載需區別，僅對這兩個既有 helper 增加 additive readonly inspection，保留原 writer、key/version、格式及相容 wrapper；無法讀取／無效不能冒充 observed none，已有 O1 不可被清空。schema 的任何非空 opaque stored string 仍按原 lookup identity 保留為 unresolved，即使外觀可疑也不新增 UUID grammar／JSON parser 或改成 none；安全顯示與原查詢 identity 分離。

Presentation bridge 若需要，僅傳當前 scope 的 token-free status、opaque operation ID、typed reason 與既有 guarded callbacks，不傳 raw persisted request/draft/token/error，不複製 operation DTO、preview proof 或授權旗標；完整 DTO 留在 owner 與既有 result presenter。每次呼叫都由 owner 重驗當前 scope/epoch/recovery/preview 與必要 proof；stale report／closure 無副作用。NAV 僅 open/focus，不暗送 preview POST。`confirm_schema` 只在 review/consent 後 delegate 現有 schema handler；review scoped start 只打開 presentation review，最終確認才 delegate 同一既有 start handler。scope-review 重用 2.2 summary；open/cancel/reopen 均不送 POST。workspace/device/group/revision 或任何顯示中的原生 draft config、selected sources、profiles、interval、destination、semantics 在確認前改變，都必須使 review consent 失效（draft 變動不必等 revision 更新）。只用 owner 當前 epoch／consent fingerprint，不複製持久設定或新增 start authority。standalone／Advanced `GroupSchemaPanel` 預設行為不變。

2.3 驗收須先固定：O1 優先且無新操作、無 ID 不虛構 request、readable-empty 不是全域 absence、unreadable/malformed 不變 idle、offline NAV 保留且 network 禁用、單純開 review 不 preview/DDL/start、stale callback 重新檢查後不 dispatch；再將 owner status／callbacks 接入單一 primary，保留全部原 guard／call-count assertions。此修正不擴張未來 3.x 管理／導航或 runtime locator 範圍。

### 摘要先行與進階詳情

顯示「選設備與資料 → 記錄方式 → 確認儲存規則 → 準備並開始」的 Step 4 內部階段。`RecordingRuleSummary` 清楚列來源與型別、layout、profile、mean/latest、週期及單位、缺值／遲到策略、目的表、安全 owner 摘要。saved draft 與 applied 差異並列標註，不拿新 draft 當正在執行。

一般路徑不要求看 UUID、revision、SQL column mapping；詳細 drawer 仍包含 source binding、row policy、counter policy／epoch、schema proof、operation 與 durable delivery stages。已保存、已套用、已接受、SQL committed、readback、cleanup 必須分開，不用綠色總狀態混為成功。

### 已設定記錄與工作區狀態

Basic 與 Advanced 讀同一 WriteGroup。已設定卡片直接提供「修改」「停止新記錄」「移除設定」與伺服器允許的「變更布局／目的地」「重新建立」，不自行猜 canDelete/canApply。disabled action 顯示理由；有 pending backlog 時顯示仍將補送原目的地的影響。

Advanced 標題改「進階設定與工作區狀態」，說明這是所有已保存群組概況，不代表即將全啟用。managed profile 仍只有 reviewed scoped start 生效入口；legacy/custom 保留原能力但明示範圍，不自動轉成 Basic profile。runtime handoff 攜帶真正 activated device IDs，不擴大為所有設備。

### 標準 scope 選擇介面

提供 typed `StudioSelectionRequest`：`workspaceId`, `step`, 選用 `deviceId`, `ruleId`, `mappingId`, `groupId`, `groupRevisionHint`。接收者以現有 backend ownership/revisions 驗證後才選擇目標；hint 僅顯示歷史差異，不能授權 mutation。缺失、跨 workspace、衝突或已刪除時顯示安全 unavailable／唯讀歷史，不選另一群組，選擇操作沒有 save/apply/start 副作用。URL query 名稱、parsing、resolver links 與路由接線唯一由後續 `expose-runtime-recording-rules` 擁有；本 change 僅提供無 URL 耦合的 interface，依賴單向。

### 相依責任與漸進替換

原先整個 change 循序的執行安排已由「分期準備與啟用門檻」取代；產品整合仍要求 `fix-step4-recording-authoring` 先建立 role/name/selection/capability/unknown-operation 正確行為；`manage-recording-lifecycle-safely` 提供 action availability 與 successor；本 change 只重排與重用。新 live transition 仍等待 `freeze-counter-replay-semantics`、`fix-recording-cutover-boundaries` 完成，不用 UI fallback 解鎖。

先增加隔離且未掛載的 presentation-only components，等待前置驗收及檔案交接後才在既有 surface 共用，逐塊遷移後才移除重複操作；不另建獨立 Basic 設定資料模型。拒絕純改標題卻保留多個互相矛盾的 primary start，也拒絕把所有細節永遠藏起來。

### 分期準備與啟用門檻

**Supersedes**: guide-step4-recording-workflow / 相依責任與漸進替換

僅取代該決策中 whole-change serial execution 的排程限制，保留其所有產品 guards 與正式整合順序。前端 authoring 垂直切片、未掛載 guided 呈現、runtime contract、依該契約的 logs 安全擴充可依下列 ownership 分期準備；既有相依責任、後端兩案循序、lifecycle 前置與最終驗收均不變。拒絕七案直接同時修改共享檔案，也拒絕 fixture-only 綠燈視為產品完成。

### 現有雲端瀏覽器驗收的條件式交付例外

**Supersedes**: guide-step4-recording-workflow / 分期準備與啟用門檻

只取代上述決策及本案最終驗收中「無法實際執行 browser／截圖就一律停止交付」的門檻，不取代其產品行為或其他前置。授權來源是 2026-10-28 18:12:29 UTC 明確詢問剩餘四案同一雲端限制下是否略過實際瀏覽器／截圖、明列未驗證、其餘測試照跑且不新增環境；使用者於 18:19:30 UTC 回覆「繼續」，18:19:36 UTC 已確認此解讀。適用範圍僅 `fix-step4-recording-authoring`、`guide-step4-recording-workflow`、`expose-runtime-recording-rules`、`make-runtime-logs-actionable`；不是由第三案例外自動延伸。

現有 cloud 的 browser launch、官方 bundled browser download 及 managed-browser loopback blockers 已由本輪交接確認未變。若本案驗收時實際執行仍不可用，可條件式略過實際 browser／截圖交付 gate；保留既有環境證據與可重跑命令，不新增環境、不更改安全限制、不公開服務或繞過拒絕。不因同一未變阻擋反覆重試已拒絕操作；若支援路徑已可用則依原案例執行並如實記錄結果。

正式 browser test cases、discovery、E2E 專用 typecheck、lint 及 spec scenario 對照仍須完成；逐案例及截圖欄位記錄 `NOT RUN`／`SKIP`／「未驗證」和原因，禁止標 `PASS`。discovery／編譯／jsdom／mock fixture 通過不等於真 browser、幾何／焦點或截圖驗收；不得用合成影像補位。`Browser verification exceptions remain explicit` 固定此報告契約。本案對應 pending task 4.1 只有在這些準備檢查及其餘原驗收完成後，才可依條件式交付門檻承接後續工作，並持續揭露未驗證風險；本次 artifact 更新不勾選任何任務。

產品邏輯、依賴順序、逐檔 ownership、backend DB／frozen proof、其他測試與 review 全部保留；未執行其他平台／現場仍未驗證。authoring active CLI 10 baseline 與所有已完成描述／provenance 不改；runtime CLI 4 已完成，1.5 仍須等 guided 最終交接，不能因本例外提前啟動。本例外只處理實際 browser／截圖，不能用於略過 DB、API 或安全契約。

## Implementation Contract

- 使用者可以在同 scope 看懂每個 disabled 原因，點擊修復入口後焦點落在對應欄位；錯誤／進度由適當 alert/status 宣告。
- 一般非計數瞬時來源可走完保存→schema review→start，不強迫計量語意、報表或進階群組編輯；計數來源保留明確 confirmation。
- Basic/Advanced 互切、details 展開、Back/Cancel、切設備、重整不產生 save/apply/start 副作用，不丟失待決操作或覆寫尚未保存草稿。
- 以 390、768、1280 CSS px 視窗寬作為設計驗收條件（非效能量測），主要動作與修復原因可操作且無頁面水平溢出；新元件沿用現有 theme，不做全站換色。長來源清單可捲動、鍵盤選取與螢幕閱讀標籤完整，不以縮小字體掩蓋溢出。
- 真正開始前摘要列本次設備、profile、來源、destination 與影響；只啟動使用者審核的範圍，部分成功不報全部完成。
- 驗收用正式 component/integration/browser 測試記錄新設定、重開既有 wide/long/custom、failed probe、schema conflict、pending operation、部分成功與 runtime handoff；截圖附上同一受驗版本；keyboard Tab 可到達 primary action，修復按鈕將焦點帶到其欄位，drawer 關閉回到原觸發器。一般非計數路徑不要求額外進入 Advanced 才能完成，沒有隱藏的第二個 start。

## Risks / Trade-offs

- [多層摘要讓細節難找] → 每張卡有清楚「查看完整規則／交付證據」，保持同 scope 的返回位置與焦點。
- [UI derived phase 與 server reality 漂移] → phase 只根據帶 revision 的既有 queries，不持久化第二份 readiness；loading/error/unknown 分開。
- [進階仍被誤解為整個 workspace start] → 更名、明列只讀概況與當前 scope，測試確認沒有新增 workspace start handler。
- [相依 change 同改元件] → 準備檔案隔離，掛載按 backend frozen replay → existing cutover boundary → lifecycle → authoring → guided 的順序交接；同一 shared file 不並寫，重讀最後 source 後保留 guards。

## Migration Plan

Expand：新增 derived next-action／rule-summary 元件但不改 API。Migrate：依 destination/source、preparation/start、saved-recording/advanced 順序替換；各批單獨可建置且保留底層 guard。Contract：只移除重複 presentation；回退只恢復舊呈現，不撤銷 server operations、歷史或已保存設定。

## 執行分期與交接

本次來源為已核准的「依賴先行、可獨立部分並行」，只變更執行順序，不新增產品需求、change、環境、部署或權限。保留七個 change ID 與各 spec scenario；spec.md 不變，因可觀察產品契約沒有改變。

**Supersedes:** 原先 whole-change serial execution 僅替換為 phased preparation/activation gates；不是取消 backend correctness 或最終整合前置。共同決策見 `guide-step4-recording-workflow / 分期準備與啟用門檻`。

- 最終 UI／定位／整合驗收仍依 `freeze-counter-replay-semantics` → `fix-recording-cutover-boundaries` → `manage-recording-lifecycle-safely` → `fix-step4-recording-authoring` → `guide-step4-recording-workflow` → `expose-runtime-recording-rules` → `make-runtime-logs-actionable`。唯一前移切片是 runtime 1.3–1.4 安全讀取／redaction／effective provider：在 freeze→cutover→lifecycle 交接後完成，作為 guided 2.2 的額外前置；authoring 仍為 guided 掛載前置。runtime dashboard、locator／URL 與最終驗收仍等 guided 完成；不得用準備測試取代任何整合門檻。
- 準備成果必須是可編譯、可測且範圍明確的程式；fixture／parser／未掛載元件通過只證明該準備切片，不代表 production wiring、資料庫正確性或產品驗收。每個產品切片各自 RED→GREEN→Refactor，結束時受影響測試與建置保持通過，不先提交一整套已知失敗測試。
- 每次啟動先核對目前工作樹及 owner；同一檔案只交一位 worker。需要其他 owner 檔案時先交接，不能以不同 function／相鄰行為理由並寫。共享 locale、types、routing、fixtures 也適用；不新增測試環境。
- 現行 freeze worker 保留 `internal/datalink/{workspace,runtime,counterproof,groupdelivery,snapshot,grouppipeline}/` 的修改權；後續實作僅讀其既有契約，不同時改同一檔案。backend 交接後重讀 frozen proof／revision/reason 的實際結果再接線。
- `add-windows-tray-runtime-logs` 維持獨立 active change 的 13/23 及原完成證據；broker、file sink、SSE、local-only guard 與 Windows lifecycle 不因本次分期移交或放寬。所有未執行產品／平台／現場驗收仍未完成。

## 分期 ownership 矩陣

| 階段／owner | 可修改的獨占檔案或模組 | 不可跳過的接線門檻 |
| --- | --- | --- |
| authoring 準備 | `steps/step4/{CircuitProfileConfigFields,CircuitProfileBasicPanel,CircuitPurposeFields}.tsx`、`steps/step4/circuitProfileActions.ts`、`state/circuitProfileBasic.ts` 及該案 tests | saved role 及 operation recovery 的 backend／lifecycle 交接 |
| guided 1.1–2.1 準備 | 新 `state/recordingWorkflowView.ts`、`state/studioSelection.ts`、新 `steps/step4/{RecordingNextAction,RecordingRuleSummary}.tsx`、新 `frontend/tests/unit/workbench-v2/recording-guided-*.test.ts(x)` | 不掛載現有 Step4，不修改 authoring 檔案、root wiring、共享 fixture 或 locale |
| runtime 1.1 契約 | 新 `internal/datalink/runtimeview/contract*.go`、新 `frontend/src/types/runtimeRuleEvidence.ts`、新 `frontend/src/services/runtimeRuleEvidenceParser.ts` 及專屬 tests | 精確契約審定才可產生 parser/consumer；此階段不接實際 reader |
| logs 1.1–1.3 準備 | 新 `internal/diagnostics/domain_*.go` 及專屬 tests | 等 runtime 1.1 契約交接；不改 broker／sink／guard／writer |
| guided 2.2 起掛載 | 本案 Impact 列出的既有 Step4 檔案及 `frontend/src/i18n/locales/{en,zh-TW}/workbench-v2.json` | freeze→cutover→lifecycle→authoring 及 runtime 1.3–1.4 真 provider／redaction 證據與逐檔獨占移交 |
| runtime 1.3–1.4 安全 provider | runtime-context handlers、安全 types/parser、共用 services/hooks；不掛載 runtime-dashboard 或 Step 1–4 route binding | freeze→cutover→lifecycle、backend proof revisions 與逐檔移交；只限安全 projection/provider/service/hook，1.5 起 dashboard／定位仍等 guided 最終證據 |
| logs 2.1 起 emitter／viewer | producer 所需 backend/wiring、既有 diagnostics projection/query 與 runtime-logs UI | runtime 最終契約／讀取驗收與 producer owner 交接；Windows owner 的 broker/sink/guard 不移交 |

表內 `state/`、`steps/` 相對 `frontend/src/features/datalink/workbench-v2/`；明確新檔不是現有 production 契約的佔位替身。若新檔已被別人建立，先停該檔並由 parent 決定 owner，不覆寫。

## 本案準備與掛載契約

- 1.1／1.2 完成歷史保留；1.4 補償後的 pure view calculation／摘要改消費 runtime 安全 facts、scoped selection、schema/operation evidence，對未知 lifecycle availability 回 unavailable，不先猜可開始／可刪除；沒有 API、DDL、save/start 或修改 session draft 的副作用。
- `StudioSelectionRequest` 由本案唯一宣告；選取回傳經驗證的同 workspace identity 或 typed unavailable。準備時用明確的 read-only evidence provider 參數，沒有 backend 證據就拒絕；未掛載的 adapter 測試不是完整導航驗收。
- labels 由 typed props 傳入供隔離元件測試；locale 註冊與真翻譯在 2.2 後由已取得 owner 的 worker 完成，不提前碰共享 JSON，也不以未翻譯 placeholder 掛到產品。
- 2.2 掛載前重讀 authoring/lifecycle 最終來源及 runtime 1.3–1.4 的安全 provider／redaction 實測證據；逐塊接入真 query/callback 並跑整合測試，不替換已驗收的 role、capability、unknown-operation 或 frozen-binding guards。4.1–4.2 的 browser/a11y/viewport 與最終 gates 原樣保留。
- 本案不擁有 URL 名稱、query parser、EntityLocator 或 resolver。runtime 唯一擁有這些，再把已驗證 locator 映射成 StudioSelectionRequest；logs 只消費 runtime resolver。導航依賴為 guided selection → runtime locator binding → logs；安全讀取切片另由 runtime provider → guided summary，分期而非整案相互等待。

## 安全生效證據交接

**Supersedes:** 本案「執行分期與交接」及「契約準備與 production 讀取分界」（guided 對應「本案準備與掛載契約」）中把全部 runtime reader 延後到 guided 最終驗收的門檻；只前移安全 provider，不改產品需求、七案範圍或最終驗收。

- 已查證 `StudioV2WorkspaceWriteGroupsHandler.Get` → `WriteGroupService.Get` 回 editable saved group，不能當 immutable effective reader。`setAppliedInTx` 可先將 saved 的 applied_revision 改 B；`resolveVersionInTx` 按 effective_at <= as-of 選版本，因此未到 boundary 時 `ResolveAppliedAt` 仍為 A。`ResolveAppliedAt` 的 defensive Group／EffectiveAt 是 backend authority，不等於已有 public reader；frozen runtime layout 目前由 grouppipeline recovery 私有鏈持有，交接後僅透過安全讀取 adapter 消費，不改 writer/recovery authority。
- 固定語意為 saved revision、latest requested applied revision（可 scheduled）、effective revision at server as-of、各 snapshot effective_at；不得把 saved.applied_revision 當當下生效版本。回應對齊 workspace/entity/revision/as-of，切版競態回 stale-reference，缺 effective/layout proof 回 unavailable；時間來自 server 證據，不由 browser interval 倒推。
- retained ended（包括 EndsAt 已到期）／disabled／tombstone snapshot 可供歷史與 backlog 解讀，不代表允許 intake 或 start；intake 狀態及 lifecycle action availability 另依對應 revision 的 authority。既有操作未決優先、CAS、readiness、confirmation 不變。
- runtime 唯一擁有安全 projection／DTO/parser 及共用 service/hook；guided 消費安全 facts 與 evidence status，不將 `GroupFacts` cast 為完整 `WriteGroup`，不補造被刻意省略的 destination.database（可能是私人 SQLite path），不使用 editable GET 冒充 applied snapshot。可見摘要若需額外安全欄位，仍由同一契約 owner 檢證投影，不能復原 raw connector。
- runtime 1.2 先補償已完成 1.1 並交接 guide 1.4、通知 logs 維持相容；1.2 僅隔離契約與受影響 logs consumer tests，可先行；1.3 才在 backend/lifecycle 交接後執行所有 legacy/new projection redaction，必須在首次 production 曝光前同批通過。1.4 前移原 2.1 provider 並交付 API/service/hook 證據。guided 2.2 同時等待 authoring、lifecycle、guide 1.4／2.1 及 runtime 1.3–1.4；runtime 1.5 起 dashboard、locator/URL、browser／最終驗收仍等待 guided 最終交接，故無循環。
- 完成紀錄不回寫：guide 1.1／1.2 的隔離契約與 next-action 測試不證明真 effective reader；runtime 1.1 的 parser 證據不證明 production redaction。補償均保持 pending；不得以本次 artifacts 修正聲稱 source、browser、SQL 或現場驗收完成。

## Context

來源基準為雲端 checkout main 的 913cead9；本案只從 source 查證，未連線使用者私人 host。AGENTS.md／CLAUDE.md 優先於舊 inventory 與 openspec/config.yaml 的歷史技術描述。已讀 studio inventory 的 START_HERE.md、context.json、CURRENT_STATE.md。

- RuntimeSetupContextPanel 以 source_rules.slice 與 mappings.slice 展示摘要，RuntimeSetupDatabase 只呈現 connector／host／database／table 與 target count；useRuntimeDashboardState 把 workspace-wide setup 直接傳入。
- runtime_workspace_setup_context.go 的 mapRuntimeSetupSourceRule 未回 scale／offset／data format；mapRuntimeSetupMapping 未回 transformation pipeline；runtime-context 沒有 canonical write groups。
- schema_source_rule_models.go 的 SourceRule、schema_tag_mapping_models.go 的 Mapping 已有轉換資料；Point.PollingGroupID 對應 PollingGroup.IntervalMs，不能以寫入週期推論採集週期。
- WriteGroup 是目的地與記錄政策 authority；WriteGroupService.ResolveAppliedAt 提供 immutable applied snapshot；grouppipeline.DeliveryForRevision、sourceEvidence 及 groupdelivery.LastCommittedEffect 是 durable 證據。
- RealtimeLogsPanel 無條件 LIVE，useRuntimeStream 的 logs 是 browser transport／stale 事件，不是 SQL receipt。RuntimeDashboardPage 將歷史報表設為底部分頁預設。
- 相關 specs 為 studio-runtime-workspace-view、runtime-dashboard-truthfulness、runtime-write-group-delivery；本案不覆寫其既有要求。active add-windows-tray-runtime-logs 的 local-only 日誌權限不變。

## Goals / Non-Goals

**Goals:**
- In scope：唯讀規則鏈、有效設定與草稿差異、durable delivery drilldown、typed unknown、安全 identity、精準設定定位、focused monitor 的資訊順序及測試。
- 讓操作者直接判斷哪個規則採集、怎麼轉換、多久寫哪張表，以及尚未寫入的原因。

**Non-Goals:**
- Out of scope：新 writer、修改聚合／counter 算法、schema DDL、刪表、變更 cutover、啟停或重試 mutation、新報表平台、歷史資料查詢擴張、外部部署。
- 不開放遠端程序日誌，不代理或放寬 /system/logs，不建立另一份 durable receipt store。
- 不自動測 PLC／正式 DB，不以 UI 優化移除 revision、readiness、ownership 或確認門檻。

## Decisions

### 共用安全讀取模型

本案唯一擁有 RuleEvidenceView 與 EntityLocator 契約；make-runtime-logs-actionable 只引用此 identity／revision／stage／action，不建立第二套映射。採 backend read projection，而非前端把新 draft 與舊 receipt 自行 join。

RuleEvidenceView 由既有 workspace/runtime-context 讀取鏈提供 summary；展開詳細資料走同 workspace 的唯讀分頁。沿同一 workspace revision 收集資料；遇到讀取期間 revision 改變回 typed stale-reference，不拼出混合結果。過渡期只保留安全的舊欄位；所有 operator-accessible API projection 立即移除或遮蔽 unsafe endpoint／credential／private-path 舊欄位，舊前端以缺值相容處理，不透過另一 legacy 欄位繼續洩漏。

來源項目包含 workspace/device/rule/point/tag/mapping identity、saved revision、effective runtime revision、proof state、enabled、protocol function、來源 register/address、read/result datatype、有效 polling interval 與其來源、safe transform steps。byte order／scale／offset／cast 採已驗證的型別與參數；未知或不安全 transform 顯示 typed unsupported detail，不 stringify 任意 params。來源 register/address 不是網路 endpoint。

若 runtime 只有 projection version 但沒有足以還原 effective source facts 的證據，回 unknown，不拿目前 SourceRule 填入 active 區。API 必須明確標識 saved-only、effective 或 unavailable。

### 生效版本與交付證據

群組項目包含 group ID、draft/applied revisions、effective_at、members、connector alias/kind/ID/revision、schema/table/column、layout、storage mode、interval、latest／mean／cumulative delta、write/incomplete policy。active 欄位由 ResolveAppliedAt 及 frozen runtime layout 取得，不用 editable destination。

DeliveryView 為唯一交付讀取來源：現行 revision stages 與 historical backlog 分列；receipt、source evidence、window、oldest pending、quota、storage repair 不混為一個 ready。SQL committed 不等於 readback verified 或 cleanup completed。測試寫入操作有獨立 provenance，不列作 runtime receipt。

next_due_at／next_retry_at／attempt_duration_ms 為 nullable evidence；只有實際 scheduler／outbox／attempt 讀值才展示，沒有就寫「尚無證據」。不從 interval 製造精準倒數。舊表 backlog 維持 frozen connector revision/table，active 群組更新不能重新標記它。

### 精準定位與設備範圍

EntityLocator 使用 step、workspace_id、device_id、rule_id、mapping_id、group_id 的適用欄位；runtime detail 的 group_revision 僅作歷史選取，絕不繞過 mutation CAS。所有 identity 都經 workspace ownership 檢查，不以名稱比對。

guide-step4-recording-workflow 可在隔離準備期先提供 canonical scope-selection interface，不依賴本案；本案的 locator binding 仍待該介面實際交接與 guided 最終驗收。其後本案擁有 query／resolver／deeplink、Step 1–3 選取與 Step 4 interface binding；guided 只消費先行安全 provider，不消費本案 locator／URL query，分期避免循環依賴。manage-recording-lifecycle-safely 提供 tombstone 行為；刪除、foreign workspace、stale reference 顯示安全原因並回該 workspace 的來源清單，不選另一同名項目。

預設只看 selected device 的來源／群組成員；多設備共用群組明示共有範圍與 group-level totals，不能歸屬為該設備每個 Tag 的成功數。提供搜尋／分頁與詳情，取代不可展開的 +N。

### Focused monitor 與降級

主視線依序為來源規則、資料庫寫入規則、最近交付證據；即時值與 collector health 保留，歷史報表維持次要且不增功能。transport 事件明確標「連線事件」，LIVE 只能在 stream ready 且資料 freshness 有證據時成立。

設備切換清除舊選取與其請求，遲到結果不覆蓋新設備。SSE 故障保留最後成功 snapshot、顯示時間／stale 並使用既有 polling recovery；錯誤提供 typed action，不顯示 raw backend exception。

### 現有雲端瀏覽器驗收的條件式交付例外

**Supersedes**: expose-runtime-recording-rules / Focused monitor 與降級

只取代上述決策及本案最終驗收中「無法實際執行 browser／截圖就一律停止交付」的門檻，不取代其產品行為或其他前置。授權來源是 2026-10-28 18:12:29 UTC 明確詢問剩餘四案同一雲端限制下是否略過實際瀏覽器／截圖、明列未驗證、其餘測試照跑且不新增環境；使用者於 18:19:30 UTC 回覆「繼續」，18:19:36 UTC 已確認此解讀。適用範圍僅 `fix-step4-recording-authoring`、`guide-step4-recording-workflow`、`expose-runtime-recording-rules`、`make-runtime-logs-actionable`；不是由第三案例外自動延伸。

現有 cloud 的 browser launch、官方 bundled browser download 及 managed-browser loopback blockers 已由本輪交接確認未變。若本案驗收時實際執行仍不可用，可條件式略過實際 browser／截圖交付 gate；保留既有環境證據與可重跑命令，不新增環境、不更改安全限制、不公開服務或繞過拒絕。不因同一未變阻擋反覆重試已拒絕操作；若支援路徑已可用則依原案例執行並如實記錄結果。

正式 browser test cases、discovery、E2E 專用 typecheck、lint 及 spec scenario 對照仍須完成；逐案例及截圖欄位記錄 `NOT RUN`／`SKIP`／「未驗證」和原因，禁止標 `PASS`。discovery／編譯／jsdom／mock fixture 通過不等於真 browser、幾何／焦點或截圖驗收；不得用合成影像補位。`Browser verification exceptions remain explicit` 固定此報告契約。本案對應 pending task 4.2 只有在這些準備檢查及其餘原驗收完成後，才可依條件式交付門檻承接後續工作，並持續揭露未驗證風險；本次 artifact 更新不勾選任何任務。

產品邏輯、依賴順序、逐檔 ownership、backend DB／frozen proof、其他測試與 review 全部保留；未執行其他平台／現場仍未驗證。authoring active CLI 10 baseline 與所有已完成描述／provenance 不改；runtime CLI 4 已完成，1.5 仍須等 guided 最終交接，不能因本例外提前啟動。本例外只處理實際 browser／截圖，不能用於略過 DB、API 或安全契約。

## Implementation Contract

- 跨案 reason 消費：freeze-counter-replay-semantics 擁有 counter-replay-proof-missing／counter-replay-proof-conflict／counter-replay-proof-unsupported；fix-recording-cutover-boundaries 擁有 recording-ownership-unavailable／recording-ownership-conflict／recording-handoff-proof-conflict／recording-cutover-finalized-overlap。本案按 group/revision、pending count 與 typed action 呈現，finalized output 與 replay quarantine 分開；只導向恢復原 immutable evidence 的指引，不重新貼 epoch/unit、不修復來源資料或重送 finalized output。
- API：擴充 runtime-context 的安全 summary，細節採同範圍 read-only query；list 預設 50、上限 100 筆為設計限制，不是量測結果，server cursor 與 workspace revision 綁定。無權限／不存在／stale 分別為安全 typed error，不默默切換設備。禁止 N 個群組無上限 fan-out。
- UI：規則完整列表、選取詳細內容、active/draft 對照、現行／歷史交付、可用／缺失證據、精準設定連結；讀取失敗可重試且不改設定。
- 安全：response、DOM、URL、analytics 不含 host/port、username、DSN、密碼、token、私人 SQLite path、raw SQL／driver text；目的地顯示安全 alias/kind/schema/table，不輸出原始 connector DTO。
- 接口 ownership：fix-step4-recording-authoring 擁有記錄 policy 正確性；manage-recording-lifecycle-safely 擁有停用／移除；本案依賴 guide-step4-recording-workflow 的 scope-selection interface 並加入 EntityLocator binding；make-runtime-logs-actionable 消費本案只讀契約。缺少依賴時保留 typed unavailable，不冒充完成。
- 驗收：API tests 對 draft/applied 不同目的地、讀取中 revision 切換、foreign scope、unknown timing、receipt mismatch；UI tests 對超過摘要數量、A→B 遲到、Back/Forward、共享群組、disabled backlog、no-data、delta unknown、mean partial、local-only logs 說明。
- 展示原則：採集成功、accepted、queued、SQL committed、readback、cleanup 分別有證據；空資料不是 0，legacy/custom 無 per-source proof 不推算。
- 舊 caller 的安全欄位維持相容；unsafe 欄位不因過渡而保留，前端 parser 缺少新欄位時顯示 unavailable，不 fallback 到 unsafe DTO。

## Risks / Trade-offs

- [跨 API revision 不一致] → 同 revision projection 與 typed stale-reference；重試只讀取，不修改／套用。
- [把總群組成功當每來源成功] → group totals 與 verified source evidence 分開；absence 永不補值。
- [規則多造成高負載] → 分頁與展開才讀詳細資料，測試請求數／有界回應；本案不宣稱尚未量測的效能改善。
- [舊註解與能力漂移] → readiness／frozen proof 為依據，不能以 write_group_profile.go 的舊 unsupported 註解決定 UI。
- [安全名義下隱藏必要診斷] → 保留 source address、typed transform、safe table identity 與 action，拒絕任意 connector dump。

## Migration Plan

先 expand 新 typed projection 與 parser，舊 UI 仍可建置；再分別遷移來源、群組、交付及定位；最後收斂過時摘要和無條件 LIVE。每批執行受影響測試，無 schema migration。rollback 只回復讀取 UI，不更動 runtime、outbox 或設定。部署與現場驗證另需授權。

## Open Questions

無待決產品問題。實作前依當時 source 重核 API envelope 與 capability；若無法證明 effective source/timing，依本案固定回 unavailable，不能延後安全判斷。

## 執行分期與交接

本次來源為已核准的「依賴先行、可獨立部分並行」，只變更執行順序，不新增產品需求、change、環境、部署或權限。保留七個 change ID 與各 spec scenario；spec.md 不變，因可觀察產品契約沒有改變。

**Supersedes:** 原先 whole-change serial execution 僅替換為 phased preparation/activation gates；不是取消 backend correctness 或最終整合前置。共同決策見 `guide-step4-recording-workflow / 分期準備與啟用門檻`。

- 最終 UI／定位／整合驗收仍依 `freeze-counter-replay-semantics` → `fix-recording-cutover-boundaries` → `manage-recording-lifecycle-safely` → `fix-step4-recording-authoring` → `guide-step4-recording-workflow` → `expose-runtime-recording-rules` → `make-runtime-logs-actionable`。唯一前移切片是 runtime 1.3–1.4 安全讀取／redaction／effective provider：在 freeze→cutover→lifecycle 交接後完成，作為 guided 2.2 的額外前置；authoring 仍為 guided 掛載前置。runtime dashboard、locator／URL 與最終驗收仍等 guided 完成；不得用準備測試取代任何整合門檻。
- 準備成果必須是可編譯、可測且範圍明確的程式；fixture／parser／未掛載元件通過只證明該準備切片，不代表 production wiring、資料庫正確性或產品驗收。每個產品切片各自 RED→GREEN→Refactor，結束時受影響測試與建置保持通過，不先提交一整套已知失敗測試。
- 每次啟動先核對目前工作樹及 owner；同一檔案只交一位 worker。需要其他 owner 檔案時先交接，不能以不同 function／相鄰行為理由並寫。共享 locale、types、routing、fixtures 也適用；不新增測試環境。
- 現行 freeze worker 保留 `internal/datalink/{workspace,runtime,counterproof,groupdelivery,snapshot,grouppipeline}/` 的修改權；後續實作僅讀其既有契約，不同時改同一檔案。backend 交接後重讀 frozen proof／revision/reason 的實際結果再接線。
- `add-windows-tray-runtime-logs` 維持獨立 active change 的 13/23 及原完成證據；broker、file sink、SSE、local-only guard 與 Windows lifecycle 不因本次分期移交或放寬。所有未執行產品／平台／現場驗收仍未完成。

## 契約準備與 production 讀取分界

- 1.1 可先行的成果只有精確契約與未掛載 parser：先把 Go／TypeScript 欄位、enum、nullable、revision binding、typed unavailable/stale-reference、reason/action 對照寫成可審核契約，再實作 parser。欄位未定不得讓另一案猜測或自行設計平行 DTO；1.1 不通過就沒有下游 logs 準備權。
- 契約唯一 owner 是本案。獨占新 `internal/datalink/runtimeview/contract.go`、`contract_test.go`、`frontend/src/types/runtimeRuleEvidence.ts`、`frontend/src/services/runtimeRuleEvidenceParser.ts`、`frontend/tests/unit/services/runtime-rule-evidence-parser.test.ts`。`runtimeview` 只存安全唯讀 DTO／enum，不能查 DB、變更 writer 或維護新 receipt store；供 diagnostics 引用時不 import API handlers，避免 import cycle。
- `RuleEvidenceView` 契約必須完整區分 workspace/device/entity identity、saved/effective revision/proof state、來源 polling/transform、群組 frozen destination/layout/policy、DeliveryView current/backlog/source receipt、nullable scheduler timing。`EntityLocator` 保持前述 step/workspace/device/rule/mapping/group/revision-hint；stage/action 僅定義觀察階段及唯讀下一步，不能變成 mutation capability。輸入來源與 unavailable 規則逐欄對應既有 SourceRule/Mapping/Point/PollingGroup、ResolveAppliedAt、DeliveryForRevision，不填入編造值。
- 1.1 的 serializer/parser tests 覆蓋 missing/unknown、cross revision、unsafe fields；fixtures 明示為測試資料，未掛接 runtime API。準備完成不宣稱 operator API 已 redacted，也不宣稱真讀取完成。
- 1.3–1.4 production 安全讀取僅等待 freeze、cutover、lifecycle 的正式交接；1.5 起 dashboard／locator／最終驗收仍等 authoring、guided 交接，parent 確認沒有共用檔案 worker。先在第一個 production 批次移除所有 legacy/new operator projection 的 unsafe 欄位，確認舊 UI 缺值相容，再接新 reader；不得先上新 reader 而把舊洩漏延後。
- 真 reader 只讀前置完成後的 immutable proof/revision，counter-replay-proof-*／recording-ownership-* 等 reason 以 owner 最終定義為準。若 1.1 契約由 pending 1.2 補償，先由本案維持同一相容契約並通知 logs；不得把未知 proof 結構硬編成 effective。source proof 不足仍 saved-only/unavailable，不以 browser join 或假造資料通過。
- 既有 handler/read-model 整合 owner：`internal/api/handlers/runtime_workspace_setup_context.go`、`runtime_workspace_setup_types.go`、`runtime_handler.go` 及必要新 `runtime_rule_*.go`；`internal/api/router.go` 與 Step 1–4 route/root wiring 要逐檔移交。workspace/grouppipeline 只消費其已交接 API，不同時編輯 backend worker 檔案。
- frontend 真 reader owner：`frontend/src/features/datalink/runtime-dashboard/`、`types/studioV2RuntimeContext.ts`、`services/studioV2RuntimeContext.ts`、`hooks/datalink/useStudioV2RuntimeContext.ts` 與專屬 tests；共享 locales 由 parent 排他移交。1.1 沒有權先改這些 production 檔案。
- guide 唯一提供無 URL 的 StudioSelectionRequest；本案唯一擁有 EntityLocator、URL/query parser、resolver、Step 1–4 binding。可在契約期對該介面唯讀檢視，但真 binding 等 guide 交接，logs 不反向定義 locator。最終所有 spec scenarios、redaction、browser、real evidence gates 保留。

## 首個可執行契約任務封包：1.1

- 交付：列出每一安全欄位的 JSON key、Go/TS type、允許 enum、缺值／unknown／stale 結果與 authority；定稿後以同一 contract 實作 serializer/parser。沒有 production import、handler、HTTP route、query hook 或 UI 掛載；不聲稱已有資料或目前 API 已安全化。
- 寫入清單僅五檔：`internal/datalink/runtimeview/contract.go`、`internal/datalink/runtimeview/contract_test.go`、`frontend/src/types/runtimeRuleEvidence.ts`、`frontend/src/services/runtimeRuleEvidenceParser.ts`、`frontend/tests/unit/services/runtime-rule-evidence-parser.test.ts`。若 line gate 要求拆檔，先讓 parent 排他擴定清單。Go contract 只能 import standard library；前端 parser 只 import 新 type，無副作用或既有 API client import。
- 既有穩定來源可唯讀：SourceRule/Mapping/Point/PollingGroup models、WriteGroup/WriteGroupAppliedSnapshot 的已存在識別欄位、DeliveryView 的 stage/backlog/source/receipt 與 timing shape。正在修改的 counter proof 內容不得拷貝成另一份 authority；只預留明確 unavailable 的證據狀態，不偽造尚未完成的 proof。1.2 真 adapter 再消費 backend 最終契約。
- 必測：有效 saved-only/effective 可區分且 evidence/revision 齊全才 effective；missing/unknown enum 安全降級或拒絕；mixed workspace/revision 拒絕；current/backlog 與 test-write/runtime/readback/cleanup 不互換；nullable timing 不合成；unsafe connector/host/path/raw-error/nested params 不被輸出；locator 歷史 hint 不變成 mutation permission；Go JSON 與 TS parser 以同一組具名契約案例對照，不靠任意 JSON blob 或 any。
- 驗證：`go test ./internal/datalink/runtimeview`、`go vet ./internal/datalink/runtimeview`；frontend 的 `npm test -- --run tests/unit/services/runtime-rule-evidence-parser.test.ts`、新 TS 檔的 `npx eslint`、`npx tsc --noEmit`、`npm run build`；root `git diff --check`、`make check-lines`。用 import/reference 檢查確認 production 尚無 consumer。此組 PASS 後 parent 才可交接 logs 1.1，不能因此解鎖 runtime 1.2。

## 安全生效證據交接

**Supersedes:** 本案「執行分期與交接」及「契約準備與 production 讀取分界」（guided 對應「本案準備與掛載契約」）中把全部 runtime reader 延後到 guided 最終驗收的門檻；只前移安全 provider，不改產品需求、七案範圍或最終驗收。

- 已查證 `StudioV2WorkspaceWriteGroupsHandler.Get` → `WriteGroupService.Get` 回 editable saved group，不能當 immutable effective reader。`setAppliedInTx` 可先將 saved 的 applied_revision 改 B；`resolveVersionInTx` 按 effective_at <= as-of 選版本，因此未到 boundary 時 `ResolveAppliedAt` 仍為 A。`ResolveAppliedAt` 的 defensive Group／EffectiveAt 是 backend authority，不等於已有 public reader；frozen runtime layout 目前由 grouppipeline recovery 私有鏈持有，交接後僅透過安全讀取 adapter 消費，不改 writer/recovery authority。
- 固定語意為 saved revision、latest requested applied revision（可 scheduled）、effective revision at server as-of、各 snapshot effective_at；不得把 saved.applied_revision 當當下生效版本。回應對齊 workspace/entity/revision/as-of，切版競態回 stale-reference，缺 effective/layout proof 回 unavailable；時間來自 server 證據，不由 browser interval 倒推。
- retained ended（包括 EndsAt 已到期）／disabled／tombstone snapshot 可供歷史與 backlog 解讀，不代表允許 intake 或 start；intake 狀態及 lifecycle action availability 另依對應 revision 的 authority。既有操作未決優先、CAS、readiness、confirmation 不變。
- runtime 唯一擁有安全 projection／DTO/parser 及共用 service/hook；guided 消費安全 facts 與 evidence status，不將 `GroupFacts` cast 為完整 `WriteGroup`，不補造被刻意省略的 destination.database（可能是私人 SQLite path），不使用 editable GET 冒充 applied snapshot。可見摘要若需額外安全欄位，仍由同一契約 owner 檢證投影，不能復原 raw connector。
- runtime 1.2 先補償已完成 1.1 並交接 guide 1.4、通知 logs 維持相容；1.2 僅隔離契約與受影響 logs consumer tests，可先行；1.3 才在 backend/lifecycle 交接後執行所有 legacy/new projection redaction，必須在首次 production 曝光前同批通過。1.4 前移原 2.1 provider 並交付 API/service/hook 證據。guided 2.2 同時等待 authoring、lifecycle、guide 1.4／2.1 及 runtime 1.3–1.4；runtime 1.5 起 dashboard、locator/URL、browser／最終驗收仍等待 guided 最終交接，故無循環。
- 完成紀錄不回寫：guide 1.1／1.2 的隔離契約與 next-action 測試不證明真 effective reader；runtime 1.1 的 parser 證據不證明 production redaction。補償均保持 pending；不得以本次 artifacts 修正聲稱 source、browser、SQL 或現場驗收完成。

- 已查證 `grouppipeline.DeliveryView.CounterReplayQuarantines` 與 `counter_quarantine.go` 的 accepted-input scope/revision、pending_samples、reason/action；目前 `runtimeview.DeliveryEvidence` 僅 current/backlog/test-write，未承載此獨立陣列。pending 1.2 依既有 backend 權威補償安全 DTO/parser 與 consumer tests，不將 input quarantine 計入 finalized SQL output，不修改 backend producer 或另造 receipt authority；production adapter 與卡片仍各由 1.4／2.2 交接。

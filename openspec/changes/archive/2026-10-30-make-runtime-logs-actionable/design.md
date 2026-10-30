## Context

查證來源為雲端 checkout main 的 913cead9，未連線使用者私人 host。internal/diagnostics/projection.go 的 templates 僅有程序／HTTP／泛用 runtime 事件；project 只允許 HTTP 與 shutdown 的固定欄位。cmd/test_ui/service_wiring.go 的 pipeline OnError 發 runtime.group_failed，沒有 group identity。RuntimeLogsPage 顯示 code/message 與 fields JSON；另一路 runtime dashboard 的 useRuntimeStream logs 是 browser connection/stale 事件。

既有 active add-windows-tray-runtime-logs 已擁有 safety projection、bounded broker、file sink、snapshot/replay/SSE、localLogGuard 與 Windows lifecycle。其已完成證據與未完成工作全部保留。本案不從歷史勾選推定當前平台驗收通過，也不重新定義 log viewer transport。

已讀來源：internal/diagnostics/event.go、projection.go，internal/api/runtime_log_guard.go、runtime_log_handler.go，cmd/test_ui/service_wiring.go，frontend/src/features/runtime-logs/RuntimeLogsPage.tsx、runtimeLogBuffer.ts、runtimeLogSession.ts，以及 durable delivery_view.go／committed_effect.go。

## Goals / Non-Goals

**Goals:**
- In scope：安全 domain-event producer、typed event projection、有界故障合併、entity/stage/actionable 篩選、可定位事件詳細內容、production wiring 與安全回歸。
- 操作者能知道哪個來源／寫入群組、哪個生效版本出錯，何時開始、是否恢復、下一步去哪裡修正。

**Non-Goals:**
- Out of scope：任意 raw log capture、逐 sample／逐列永久日誌、另一套 report/log database、完整 audit 平台、任意 SQL 展示、遠端 process-log access、替代 durable receipt。
- 不修改採集／記錄／cutover／retry/skip 的行為，不自動執行修復；按 action 只導向原有需確認的介面。
- 不重新實作 Windows tray、file rotation、transport 或放寬 localhost guard；不部署或接觸正式設備與資料庫。

## Decisions

### 事件只引用共用真相

準備依賴 expose-runtime-recording-rules 的 1.1 完成並審定 RuleEvidenceView／EntityLocator／stage／action 精確契約；production 接線及最終驗收硬依賴該 change 完整交接。該案擁有唯讀真相與定位，本案只引用 identity／revision，不複製群組、來源、目的地 authority。

producer 從實際執行結果建立 DomainEventRef：workspace_id、device_id、rule_id、mapping_id、point_id、tag_id、group_id 的適用集合；group_revision／connector_revision、stage、typed code、action、opaque request/operation/effect reference。缺 scope 證據時只發不帶錯誤歸屬的安全 unavailable 事件，不從名稱或最近操作猜測。

SQL committed 事件必須在匹配 receipt 的 durable transition 後發出，不能從 buffer accept、OnError 為 nil 或 HTTP 200 推斷。發事件失敗不改變 delivery outcome；重啟不把歷史 receipt 重發成新 commit。log gap 後由 read model 查當前真相，不能補造完整事件歷史。

### 安全白名單與本機門禁

在 diagnostics project 先驗證 code、opaque IDs、enum、合法時間與有界 scalar，再送到 memory/file/SSE；禁止任意 error、params、connector DTO、payload 或 nested map。訊息由 fixed code 對應 en／zh-TW 文案，fields 不再直接 dump 成主要 operator UI。

允許安全階段、時間、發生次數、revision 與 action 參照；不帶 host/port、username、DSN、credential、私人 path、raw SQL／driver text、sample value。目的地與規則細節透過共用安全 read model 開啟，不塞每條事件。

localLogGuard 與 /api/v1/system/logs、/stream 路徑保持同一限制：loopback peer/local address/Host、同源、拒絕 proxy headers。遠端 /studio/runtime 仍能用其既有安全 workspace APIs 看規則／交付，不能 tunnel 此 feed；遠端 /studio/logs 403 後停止自動重試並說明需在 gateway 本機開啟。使用者私人網段網址不是門禁變更授權。

### 有界合併與隔離

同 workspace/entity/applied revision/stage/code 的重複故障合併為一組；首次立即送出，持續故障以既有 runtimeLogHeartbeat 間隔作摘要窗口，記 first_seen_at、last_seen_at、occurrences；恢復獨立事件並終止舊 fault 狀態，revision 改變不共用狀態。

coalescer 同時受既有 AdmissionRecords／AdmissionBytes 硬界線限制；高 cardinality、飽和、關閉或失敗只走既有非阻塞 admission／drop accounting，不無界增加 map 或阻塞 SQL。coalesced occurrences 與 dropped events 不混算；counter 用經驗證的十進位字串避免瀏覽器精度損失。容量不足無法保留的 detail 以既有 gap/drop 語意呈現，不宣稱已完整捕捉。

既有實作常數為 8 KiB/event、ring 2000 records/4 MiB、replay 500 records/1 MiB、live 256 records/512 KiB、最多 8 subscribers；browser 1000 records/2 MiB。這些是 source 查證的現有設計限制，不是吞吐量量測。不得為新欄位放大現有 caps；首筆成功、故障恢復與窗口成功摘要取代每 sample 或每 SQL row spam。

### 可行動篩選與歷史定位

擴充既有 typed filter：workspace_id/device_id/rule_id/group_id、stage、actionable；與 level/source/q 一起作用於 snapshot/replay/live，filter identity 改變取消舊 stream、cursor 綁定相同 filter，保留既有 bounded replay semantics。

事件列顯示時間、設備／群組安全 identity、階段、問題、次數及 action；展開顯示原 applied revision、首次／最近時間、receipt/reference。修復只導向 EntityLocator，不自動 retry／skip／DDL。歷史事件不因今日 draft 改名或刪除而被重標成當前狀態；resolver 不可用時展示 immutable ID 與安全 fallback。

active fault 只由同 scope/revision 的 typed lifecycle evidence 判斷；有限 retention 的「沒有事件」不是沒有故障，掉落與重啟標示不消失。「只看需處理」含 unresolved/unknown 與缺乏 recovery 證據的事件，不以搜尋字串猜測。

### 現有雲端瀏覽器驗收的條件式交付例外

**Supersedes**: make-runtime-logs-actionable / 可行動篩選與歷史定位

只取代上述決策及本案最終驗收中「無法實際執行 browser／截圖就一律停止交付」的門檻，不取代其產品行為或其他前置。授權來源是 2026-10-28 18:12:29 UTC 明確詢問剩餘四案同一雲端限制下是否略過實際瀏覽器／截圖、明列未驗證、其餘測試照跑且不新增環境；使用者於 18:19:30 UTC 回覆「繼續」，18:19:36 UTC 已確認此解讀。適用範圍僅 `fix-step4-recording-authoring`、`guide-step4-recording-workflow`、`expose-runtime-recording-rules`、`make-runtime-logs-actionable`；不是由第三案例外自動延伸。

現有 cloud 的 browser launch、官方 bundled browser download 及 managed-browser loopback blockers 已由本輪交接確認未變。若本案驗收時實際執行仍不可用，可條件式略過實際 browser／截圖交付 gate；保留既有環境證據與可重跑命令，不新增環境、不更改安全限制、不公開服務或繞過拒絕。不因同一未變阻擋反覆重試已拒絕操作；若支援路徑已可用則依原案例執行並如實記錄結果。

正式 browser test cases、discovery、E2E 專用 typecheck、lint 及 spec scenario 對照仍須完成；逐案例及截圖欄位記錄 `NOT RUN`／`SKIP`／「未驗證」和原因，禁止標 `PASS`。discovery／編譯／jsdom／mock fixture 通過不等於真 browser、幾何／焦點或截圖驗收；不得用合成影像補位。`Browser verification exceptions remain explicit` 固定此報告契約。本案對應 pending task 4.2 只有在這些準備檢查及其餘原驗收完成後，才可依條件式交付門檻承接後續工作，並持續揭露未驗證風險；本次 artifact 更新不勾選任何任務。

產品邏輯、依賴順序、逐檔 ownership、backend DB／frozen proof、其他測試與 review 全部保留；未執行其他平台／現場仍未驗證。authoring active CLI 10 baseline 與所有已完成描述／provenance 不改；runtime CLI 4 已完成，1.5 仍須等 guided 最終交接，不能因本例外提前啟動。本例外只處理實際 browser／截圖，不能用於略過 DB、API 或安全契約。

## Implementation Contract

- 跨案 reason 消費：freeze-counter-replay-semantics 擁有 counter-replay-proof-missing／counter-replay-proof-conflict／counter-replay-proof-unsupported；fix-recording-cutover-boundaries 擁有 recording-ownership-unavailable／recording-ownership-conflict／recording-handoff-proof-conflict／recording-cutover-finalized-overlap。本案按 group/revision、pending count 與 typed action 呈現，finalized output 與 replay quarantine 分開；只導向恢復原 immutable evidence 的指引，不重新貼 epoch/unit、不修復來源資料或重送 finalized output。
- Producer matrix：source apply/reject；collector/mapping failure/recovery；group intake applied/paused/retiring；delivery retry/blocked/quarantined/unknown/recovered；quota refusal；SQL first-confirmed/recovery summary；storage repair/cutover result。每個 producer 需寫明 authority、emit boundary、identity 與 evidence reference。
- 領域狀態 mutation 仍由 fix-step4-recording-authoring、manage-recording-lifecycle-safely 與既有 cutover owner 處理；事件只觀察已確定結果。guide-step4-recording-workflow 的選取介面由 expose-runtime-recording-rules 綁定，本案只呼叫其 resolver，沒有反向依賴。
- DTO：沿既有 RuntimeLogRecord 的 fixed code/message 與安全 scalar fields 擴充；保留 instance_id、字串 sequence、timestamp、level/source、truncated 與既有 envelope。新舊 event version 同時可讀；未知 code 安全 suppressed，不原樣印出。
- Failure：投影失敗 suppressed；producer 缺 identity unavailable；admission/file/subscriber loss 獨立顯示；403 停 retry；局部 resolver 失敗保留事件，不清空整頁；logging 不回傳新的 delivery error。
- Acceptance：真 production collector／group pipeline 使用 simulator／可丟棄 DB，證明失敗、recovery、首筆 receipt 有正確 identity；故障 burst/high-cardinality 不超容量且不阻塞 writer；queued/unknown/receipt mismatch 不假報成功；跨 revision／刪除／重啟不改歷史。
- Security：secret fixture、hostile labels、raw exception、超長欄位與 XSS 覆蓋 API／memory／file／SSE／DOM；IPv4/IPv6 local 正向與 LAN/origin/proxy 負向測試沿用既有門禁。
- Transport：snapshot→replay→live、pause/resume/clear、filter A→B、future cursor、gap/reset、overflow、slow subscriber、cancel/shutdown 不回歸；不增加 log endpoint 自我放大。

## Risks / Trade-offs

- [事件與交易成功順序競態] → receipt／ledger transition 後才 emit，不以 event 發送成功影響交易。
- [高 cardinality 壓垮合併記憶體] → 使用現有 count/byte cap 與 drop accounting，壓力測試以實際配置界線斷言，不宣稱未量測的效能數字。
- [有限日誌不是完整 audit] → UI 明示 retention／gap 與當前 evidence link，禁止 no event = healthy。
- [任意 field 造成洩漏] → 逐 code 白名單、正規化 typed scalars、所有 sink 共用一個投影器。
- [與既有 active change 重疊] → broker／transport／Windows 保留原 owner，新增領域事件在本案驗證，不重置舊任務。

## Migration Plan

先 expand safe event refs、filter parser 與 renderer，舊 lifecycle events 保持可讀；逐組加入 collector、delivery、lifecycle producers，每批可建置及測試；consumer 遷移後才撤除同一結果的重複泛用事件。rollback 只停新 producer／顯示，不更動 receipts、outbox、file ownership、guard 或 runtime。

## Open Questions

無待決產品問題。若某 producer 尚無足夠 identity 或 receipt evidence，依明確 unavailable／suppressed 契約呈現；不得先記 raw error 或假造關聯。

## 執行分期與交接

本次來源為已核准的「依賴先行、可獨立部分並行」，只變更執行順序，不新增產品需求、change、環境、部署或權限。保留七個 change ID 與各 spec scenario；spec.md 不變，因可觀察產品契約沒有改變。

**Supersedes:** 原先 whole-change serial execution 僅替換為 phased preparation/activation gates；不是取消 backend correctness 或最終整合前置。共同決策見 `guide-step4-recording-workflow / 分期準備與啟用門檻`。

- 最終 UI／定位／整合驗收仍依 `freeze-counter-replay-semantics` → `fix-recording-cutover-boundaries` → `manage-recording-lifecycle-safely` → `fix-step4-recording-authoring` → `guide-step4-recording-workflow` → `expose-runtime-recording-rules` → `make-runtime-logs-actionable`。唯一前移切片是 runtime 1.3–1.4 安全讀取／redaction／effective provider：在 freeze→cutover→lifecycle 交接後完成，作為 guided 2.2 的額外前置；authoring 仍為 guided 掛載前置。runtime dashboard、locator／URL 與最終驗收仍等 guided 完成；不得用準備測試取代任何整合門檻。
- 準備成果必須是可編譯、可測且範圍明確的程式；fixture／parser／未掛載元件通過只證明該準備切片，不代表 production wiring、資料庫正確性或產品驗收。每個產品切片各自 RED→GREEN→Refactor，結束時受影響測試與建置保持通過，不先提交一整套已知失敗測試。
- 每次啟動先核對目前工作樹及 owner；同一檔案只交一位 worker。需要其他 owner 檔案時先交接，不能以不同 function／相鄰行為理由並寫。共享 locale、types、routing、fixtures 也適用；不新增測試環境。
- 現行 freeze worker 保留 `internal/datalink/{workspace,runtime,counterproof,groupdelivery,snapshot,grouppipeline}/` 的修改權；後續實作僅讀其既有契約，不同時改同一檔案。backend 交接後重讀 frozen proof／revision/reason 的實際結果再接線。
- `add-windows-tray-runtime-logs` 維持獨立 active change 的 13/23 及原完成證據；broker、file sink、SSE、local-only guard 與 Windows lifecycle 不因本次分期移交或放寬。所有未執行產品／平台／現場驗收仍未完成。

## 獨立安全擴充與 emitter 接線

- 1.1–1.3 準備門檻只有 runtime 1.1 的精確 `RuleEvidenceView`／`EntityLocator`／stage/action 契約已完成審定；producer matrix 可唯讀對照現有來源，但 backend emit boundary 最終仍待其 owner 交接確認。沒有精確契約時不自行宣告同名平行 DTO。
- 準備 owner 只新增 `internal/diagnostics/domain_event.go`、`domain_projection.go`、`domain_coalescer.go` 與同前綴 `_test.go`。引用 runtime owner 的 `internal/datalink/runtimeview` 型別，不改該 package；不碰既有 `projection.go`、`broker.go`、`event.go`、`query.go`、file sink、SSE、guard 或任何 writer/wiring。若需要其他檔案才能保證編譯，交 parent 重新劃分，不能自行越界。
- 新 projection 是實際 fixed-code/scalar allowlist 函式，不是保留秘密的 placeholder；新 coalescer 是有現有 AdmissionRecords/AdmissionBytes 上限、注入 clock/非阻塞 admission callback 的實際模組。準備測試驗證 first/repeat/recovery/revision isolation 與容量；輸出尚未接 memory/file/SSE，也不能宣稱 production domain events 已可見。
- 舊 broker pipeline 在準備期間維持既有拒絕未知 domain code 的安全 suppression；不能為展示 mock event 而提早放開 project。memory/file/SSE/DOM 的 domain field 全鏈安全證據在 2.1 與後續整合補齊；1.2/1.3 單元通過不是該整體驗收完成。
- 2.1 起啟動門檻為 freeze→cutover→lifecycle→authoring→guided→runtime 的正式整合證據及 parent 逐檔移交。接真 producer 前先接固定 domain projector 到既有安全 admission，實測所有 sink redaction/drop accounting；如果這需要 Windows owner 正在修改的檔案，就等待 owner，不另建 broker 或 bypass sink。
- backend producer 檔案由 parent 在 gate 後按 matrix 逐個獨占分派（`cmd/test_ui/service_wiring.go`、runtime、workspace、grouppipeline、groupdelivery、對應 handlers）；不能在 freeze worker 進行中改其中任何檔案。emit 只能位於既有 authority 已確定 transition 後，logging failure 不回寫 delivery 結果。
- viewer/filter 檔案由本案在 runtime contract 最終交接後擁有 `internal/api/runtime_log_query.go`、必要既有 diagnostics query、`frontend/src/features/runtime-logs/`、`frontend/src/types/runtimeLogs.ts`、`frontend/src/services/runtimeLogs.ts` 及專屬測試；`runtime_log_guard.go` 仍不改。共享 locale 的寫入另行排他移交。
- 所有 local-only、IPv4/IPv6/origin/proxy、caps、string cursor、snapshot/replay/live、403 stop-retry 與既有 Windows 13/23 保留；最終真 producer＋可丟棄 DB＋browser proof 不可由 fixture 代替。runtime owns locator，guide owns selection，本案只讀兩者的交接結果，不產生反向依賴。

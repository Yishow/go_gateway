## Context

useUpdateStudioV2DatabaseConfigMutation 的 onSuccess 依序等待 databaseConfig、databaseTargets、bootstrap invalidation。useStudioV2DatabaseAutosave 的 runSave 於 mutateAsync 之後送 acknowledgeDbConnectorSave，已有 submitted snapshot、newer edit、workspace/mounted fences 與 unknownConnectorRef recovery。

查到 rule-centric-studio-refactor / 3. Introduce a candidate/review/apply layer as the missing orchestration boundary，理由是保存 inputs 不等於下游已 apply；guide-step4-recording-workflow / Owner-local 操作觀察與唯讀導覽分界 不允許 lost ACK 自動重送。本案保留兩項決策。

## Goals / Non-Goals

**Goals:** 使用者在合法保存 ACK 後即知道該次設定已保存，刷新中的資料另標狀態。
**Non-Goals:** 不全面修改 mutation semantics、不取消 target/write-group/schema/start/lifecycle/test-write 的 await、readback、CAS 或 operation ledger，不 optimistic 假成功、不改後端 API。瀏覽器驗收與 archive 留待全部 change 實作後於 Mac 合併執行。

## Decisions

### 保存與刷新分開結算

針對 database config save，mutation promise 只在有效 server ACK 的 scope、connector/setup revisions 與 submitted intent 可被現有 reducer 接受後結算。該 ACK 只能證明 configuration saved；不包含 test、running、SQL committed、readback verified 或 cleanup。若 response 不提供所需證據，保留確認讀取後才結算，不從 HTTP 200 推論完整保存。

refresh status 掛在現有 query state 的 scoped projection，狀態 refreshing/fresh/stale-error；不建第二份 persisted config。已確認保存後 refresh 失敗顯示「已保存，狀態更新失敗」及唯讀 retry，不把保存改成失敗；也不啟用需 fresh proof 的動作。尚在 pending/unknown 的保存仍保留原恢復與 dirty barrier。

### 精確失效與舊回覆隔離

在接受 ACK 同步標記 databaseConfig/databaseTargets/bootstrap 與 affected home summary invalidated，再以獨立 promise 並行刷新；所有 rejection 必須被 query/error UI 接住。相同 key in-flight 合併，跨 scope 不合併。開始新 mutation 時先 cancel／fence 舊讀取；older ACK/read response 不覆寫更新 local intent 或新 revision。

使用者在刷新中再編輯，saved 只指前次 submitted snapshot，新欄位仍 dirty/queued；Close/Cancel、Back/Forward 依原未送／已送／unknown 契約，不把關閉視為取消 server write。timeout/lost ACK 不自動 replay。

### 受限遷移與等待契約

先增 scoped refresh coordinator 與可觀察狀態，維持原 API；僅遷移 useUpdateStudioV2DatabaseConfigMutation 及其 autosave caller。其他 useInvalidateWriteGroups／target mutation 保留既有等待，直到另有完整 readback／caller 證據才可另案提出。

依賴 reuse-recording-home-snapshots 的 query keys 與 freshness 定義完成；沒有它時不得自製第二種 home invalidation。合成 deferred-promise 測試：三個 200 ms refresh 同時開始，config ACK 可先結算，refresh pending 仍明示；不是宣稱實際保存必達 1 ms。

### 最終 Mac 驗收與同範圍保留

**Supersedes**: shorten-recording-save-feedback / Goals / Non-Goals

依最新使用者指示，非瀏覽器實作、verify/review 與修復先完成；browser 改於全部 change 完成後一起上 Mac，archive 保留到實際驗收之後。這替換原 task 3.1 的雲端 browser／不使用 Mac 條件與 3.2 的 browser 前置，不取消任何驗收。

AdvancedWorkbenchV2Page 只保留已 hydrate、相同 workspace／API context 且有已確認 config save 的刷新中／錯誤畫面；初始或 malformed bootstrap 仍拒絕。Step4SupportPanels 的折疊摘要顯示刷新狀態，既有 RecordingDraftReadState 與 readonly guards 保持 consequential operation unavailable。Bootstrap read 的 API context key 必須保留原 canonical prefix，並更新 B 的 workspace fallback lookup，使 generic mutation owners 在 await 前仍捕捉正確 owner；不建另一份 home authority。

## Implementation Contract

- 保留 mutation request/response types、retry:false、CAS/revision 與既有 reducer acknowledgement；不得改寫 target/schema/lifecycle helpers 的語意。
- 接受 ACK 前驗 scope/revision；將 mutation completion 與 refresh completion 分開返回／呈現，不洩漏 raw backend exception、DSN 或 credentials。
- database-save-feedback.test.tsx 以可控制 deferred promises 覆蓋 ACK 先於刷新、refresh fail、409、新er edit、scope switch/unmount、late response、timeout/lost ACK 與 repeated click。
- 現有 database-autosave、database-save-queue、workspace activation barrier、reload recovery 與 scoped database save/readback tests 必須回歸；未涉及 SQL 的成功測試不能冒稱 DB receipt。
- useStudioV2WorkspaceQuery、database config/target read hooks 與 services 消費 AbortSignal，捕捉 backend／完整 baseURL／workspace／generation；被取消但晚到的 GET 不得寫入新快取。B 的 captureRecordingHomeSourceMutation／completeRecordingHomeSourceMutation 必須沿用，必要 bootstrap fallback adapter 調整須有 regression。
- en/zh-TW workbench-v2 文案、實際保存 indicator 與 disabled reason 必須一致；若只修改 hook 而 UI 看不到 stale/error，不算完成。
- In scope 是 configuration ACK/refresh；Out of scope 是所有 consequential operations 的成功語意、採集/DB runtime、全域 optimistic caching。
- 每批先寫失敗 regression，再小範圍遷移，最後移除該 configuration flow 的串行 refresh await；每批 build/type checks 保持可執行。完成一次完整 frontend lint/test/build，browser 保存／返回流程另列。

## Risks / Trade-offs

- [使用者把已保存看成已啟用] → state 文案分離與可見證據測試，readiness 未確認時仍 disable。
- [背景舊回覆覆寫新 ACK] → cancel plus generation/scope/revision fence，並測 newer edit preserved。
- [共用 helper 影響未知 caller] → 只遷移 named config caller，其他等待不動。

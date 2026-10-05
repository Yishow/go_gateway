## Context

基準 d9ba336e。診斷已用真 SQL one-connection pool 重現並行 409/500 以及受控 DELETE gap 404；原現場 3222/5173 持續運作，禁止寫現場資料。隔離 worktree /tmp/go-gateway-mapping-save-concurrency 保護 HMR。

## Goals / Non-Goals

**Goals:** 修好同規則八列並行儲存，維持不同規則進展、原子連結替換及真 stale source 拒絕；前端排隊、保留最新修改、可重試、安全具體錯誤。

**Non-Goals:** 不新增資料模型、daemon、跨程序鎖服務、driver、schema 變更、報表或盲目重試。不得放寬 accepted hash/CAS，不能因並行方便而自動 reapply 真正變更的 source。禁止修改現場 DB/設備/服務。

既有 mapping `enabled` 目前由 source rule 的 lifecycle 決定；手動將單列 false 持久化需要另外定義 operator intent，本次保持既有語意。批次驗收使用使用者本次操作「全部啟用」與 transform；前端仍忠實提交 false/true，不把排隊成功宣稱為單列停用重啟政策已修好。WriteGroup 停用／重啟的既有保護另以回歸確認。

## Decisions

### 規則協調涵蓋 production mutation 與 save-confirm

同一 sourcerule.Service 以 rule 身份協調 public create/update/delete、link mutation、candidate apply、derived sync 與 workspace mapping save-confirm；nested call 重用明確 context ownership，不能重入死鎖。檢查現有 workspace/runtime lock 順序，避免在 source guard 內取得會再次呼叫 source 的 workspace lock；跨多 rule 一次只持有一個，不能不必要地阻塞其他 rule。SQL CAS 仍為持久驗證，不把 process mutex 冒充跨程序保護。

### 連結替換是 repository 原子操作

SQL 在同一 transaction DELETE+INSERT，失敗完整 rollback，reader 看舊版或新版而非中間空集合；Memory 用單一 lock swap。遷移 service 的兩條 replace path，保留原有錯誤與 rollback 契約，不留下 silent fallback。實作如需要擴大策略先回報。

既有 PUT 保存若在 link replace／confirm 後段失敗，以 exact last-written mapping/tag 欄位 CAS 補償回上一個 accepted draft；直接修改過的資料不得覆蓋。補上 SQL mapping Update 原本漏寫、但 DTO 與 Memory 已支援的 tag_id，否則改 Tag Key 的保存會確認失敗。DELETE 先替換連結再刪 mapping，刪除失敗安全恢復連結；新 POST 若後段失敗，可保留唯一 owned disabled unconfirmed draft，再以同身份明確重試恢復。不新增 schema 或 accepted intent 模型，也不改保存 request DTO。

PUT／POST 已確認或 DELETE 已成功、但 orphan tag 清理失敗時，以成功結果加安全 cleanup_status（DELETE 另有固定 cleanup_message）區分，不能引導重試已完成的操作。前端只辨識 allowlisted status，以既有本地 state 的暫態通知呈現固定安全文案；不顯示 raw cleanup_message，不宣稱清理完成。

### 前端 rule queue 與 typed failure

同 rule 最多一個 request，其他 rule 獨立；排隊工作讀最新 state，已送出 response 不覆蓋新修改，保持 point pending/failed metadata；失敗停止自動重送，顯示明確 retry 並保留草稿。保留 HTTP status/code/request ID；映射 stale conflict、not found、save failed 安全文字進既有 en/zh-TW，禁止 SQL、credential 或 raw exception 進 operator DOM。

## Implementation Contract

後端 task: production API 八列並行於 pool 1/4 全成功，link/tag/point/pipeline/accepted hash 可讀回；替換唯一鍵/trigger失敗回滾；無 DELETE gap；save 與 source edit/sync/reapply/delete 協調且真正 stale source 409；不同 rule 不被長時間持鎖擋住。覆蓋服務 runtime wrapper、舊 direct endpoints 與 workspace routes 的讀改連結流程，維持 lock order。

前端 task: 批次 enabled/transform queue、快改同列保留最後草稿、各 rule 獨立、錯誤保留與明確 retry，hydration 不偽稱最新值已存。安全 code 使用 workspace_mapping_conflict / workspace_mapping_not_found / workspace_mapping_save_failed，status 保留 409/404/500。不改 DTO 的基本 mapping fields。

驗證 task: 先失敗再修；focused tests、完整 Go test/vet/lint、frontend lint/test/build、line/spec gates。隔離 mock PLC+owned DB+新 binary/UI 真實批次保存、preview與SQL readback，證據對應最終 source；不操作現場服務。review→fix→rereview，本地 commit、兩樹狀態回查。

## Risks / Trade-offs

單 gateway process 的 rule guard 不提供跨程序協調；持久 CAS 與原子替換不放寬。前端 HMR 風險以 worktree 隔離；Spectra 此版本會解析主 repo root，artifact CLI 暫存於主 repo 的本任務 docs，再複製至 worktree，原樹最終只清理此 owned change。source review 以隔離 worktree 的實際 Git diff/HEAD 為準，不冒稱 CLI 原樹 snapshot 代表新 source。

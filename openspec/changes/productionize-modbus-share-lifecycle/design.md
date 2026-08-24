## Context

/studio/v2 已經有 source-rule 的 durable Share 欄位、candidate snapshot、apply/restore seam，以及 internal/datalink/modbusshare 的 listener 與 memory projection；目前這些片段仍容許 browser session、固定 handler default、非 workspace-scoped mapping 與非原子 runtime mutation 互相混用。結果是 hydration、全域設定、workspace ownership、register geometry、activation、restart 與 listener lifecycle 沒有共同的 revision 邊界。

本 design 以既有 source-rule-runtime、runtime-workspace-reconciliation、studio-v2-live-config-apply、datalink-workbench-v2-settings、datalink-workbench-v2-step4-database、local-modbus-memory-workbench、protocol-servers 為約束，並延續 active change fix-protocol-address-adaptation-v2 已定義的 share_enabled、share_start_register、share_stride durable fields。新 change 只收斂 Share 的正式生命週期，不另建平行 mapping model，也不改 connector/protocol transport。

## Goals / Non-Goals

**Goals:**

- 建立 /studio/v2 bootstrap hydration → persisted workspace/settings → backend Share projection 的單向資料流。
- 讓 global settings.modbus_share.enabled、workspace ownership、source-rule/tag relationship、register geometry 與 workspace revision 成為所有 Share API 與 activation 的共同 gate。
- 在完整 prevalidation 後，以 serialized、atomic、idempotent、CAS-checked reconcile 更新 runtime projection；失敗保留舊狀態或顯式標記 dirty/unknown 並可 recovery。
- 讓 activation 等待同 revision 的 durable saves，並讓 restart 從 persisted desired mappings 恢復 listener 與 register projection。
- 讓 listener lifecycle 由 durable address/port/slave-id settings 控制，disabled-by-default 且不再由 handler 無條件 fallback 到 5020。
- 讓 delete/move/retype 清理舊 datatype span 或產生明確 invalidated state，並提供 typed error codes 與 operator-safe diagnostics。
- 將 configure→save→activate→Modbus read→restart→restore/read 納入可重跑的 HTTP、runtime、EXE acceptance matrix；實體 PLC/SCADA field acceptance 另列。

**Non-Goals:**

- 不修改 legacy /studio、/gateway/* 或新增繞過 /studio/v2 的產品入口；C 的 legacy retirement 是後續 change。
- 不替代 A 的 release correctness；A 未完成前本 change 不宣稱可發佈。
- 不重新設計既有 source-rule durable fields、candidate snapshot、tag/mapping repository 或 workspace readiness model。
- 不改 Modbus/FATEK/MC connector、frame、transport、polling 或資料型別讀取語意。
- 不把 full i18n presentation、真實 PLC/SCADA、交換器與現場硬體驗收納入 repo-only completion gate。

## Decisions

### Backend owns the truth and browser hydration gates every Share operation

Persisted source rules、persisted workspace mappings、persisted global settings 與 backend runtime status 是唯一真實來源。前端 display defaults 只可在未載入資料時呈現 skeleton/unknown，不得轉成 summary、list、delete、upsert 或 activation input。

/studio/v2 bootstrap 必須回傳 hydration_state、workspace_id、workspace_revision、settings_revision、readiness 與 Share status；ready 以外的狀態封鎖 Share mutation。hydration_failed 需保留可重試的 typed diagnostic，不能使用 default settings 繼續啟動。global enabled=false 時，summary 不渲染；list/delete/upsert/activation 不進入 mapping repository 或 runtime service，回傳 modbus_share_disabled 與目前設定 revision。

選擇 backend truth 而非 browser truth，是為了讓兩個 tab、直接 HTTP caller 與 process restart 對同一份 persisted state 得到一致結果；代價是 bootstrap 未完成時 UI 必須顯示 unknown/blocked，而不能提供看似可操作的暫存摘要。

### Workspace ownership is proven by persisted relationships

每一筆可被 Share 操作的 mapping 必須能由 persisted workspace mapping record、source rule revision 與 persisted tag relationship 三者證明 workspace_id scope。ownership query 必須以 repository/service 的 workspace scope 執行，且讀取到的 rule、tag、mapping 必須仍然存在且互相關聯。

tag_id 非空不是 ownership proof；browser 新建但尚未 durable save 的 tag、手工填入的 tag id、另一 workspace 的 tag 或沒有 source-rule 關聯的 mapping 都不得被刪除、移動或覆寫。無法證明 ownership 時，API 回傳 modbus_share_workspace_scope，保留既有 runtime state 並提供 recovery/diagnostic，而非猜測刪除。

選擇關聯證據而非 tag id allowlist，是為了避免一個 workspace 的 reconcile 清掉另一個 workspace 或外部建立的 register。代價是 legacy orphan mapping 必須進入 unknown/repair 流程，不能自動清除。

### Register geometry is canonical and backend-validated

Share output 的單位固定為兩 bytes per holding register。durable/user contract 的 share_start_register 使用 human holding-register domain（例如 40001）；runtime projection 的 register 使用 zero-based index（40001 → 0）。資料型別 span 固定如下：

- bool、int16、uint16: 1 register / 2 bytes
- int32、uint32、float32: 2 registers / 4 bytes
- int64、uint64、float64: 4 registers / 8 bytes

每個 desired mapping 的 stride_registers 必須大於或等於該 datatype 的 span；範圍採半開區間 [zero_based_register, zero_based_register + span_registers)。backend 依 listener 的 persisted capacity_registers 驗證下界、上界、整數溢位與 pairwise range collision；任何一筆失敗都使整批 prevalidation 失敗。frontend 可以提前提示，但不能放寬或取代 backend validator。API response 同時提供 human start、zero-based start、span、stride、capacity 與 datatype，避免 UI/SCADA 誤讀單位。

選擇明確 geometry 而非只存一個 register，是為了讓多 register numeric value、stride 與 stale cleanup 使用同一計算；代價是既有只支援單 word 的 mapping 需要 migration/validation，不能靜默截斷。

### Reconcile is prevalidated, atomic, serialized, idempotent, and revision-checked

Share reconcile 的輸入是 workspace_id、expected_workspace_revision、global settings revision 與完整 desired mapping set，不接受只 patch 單一 tag 的 browser projection。service 先在 immutable snapshot 上完成 hydration、global gate、ownership、datatype span、capacity、range collision、source-rule revision 與 listener settings prevalidation，全部通過後才取得 workspace lock 並以 CAS 確認 revision 未變。

mutation 以 workspace scope 序列化；同一 revision 的重送回傳相同 outcome，不重複建立 listener 或 mapping。durable repository mutation 與 runtime projection swap 必須以 transaction/staged-copy + lock 的方式完成，禁止先刪舊 mapping 再逐筆 upsert。成功後以新 revision 回應 desired/applied/invalidated summary。

若 mutation 或 runtime swap 失敗，系統先 rollback durable/runtime snapshot；rollback 也失敗時，狀態明確為 dirty_unknown，回應 modbus_share_reconcile_failed、dirty scope、observed revision 與可重跑的 recovery action，且不得回報 success。stale expected_workspace_revision 或 concurrent request 回傳 modbus_share_revision_conflict，讓 caller 重新 hydration；不能以 last-write-wins 蓋掉另一個 operator 的配置。

選擇完整 desired-set transaction 而非逐筆 patch，是為了在 collision、partial failure 與 double submit 下保持 SCADA 看到舊狀態或新狀態之一；代價是大批 mapping 更新需一次通過完整 validation，並需要 bounded lock 與 deterministic ordering。

### Activation is an autosave barrier

前端 activation controller 先 await 所有必要 durable saves：global Share settings、source-rule Share fields、tag/mapping changes 與任何 required candidate/apply decision。任一 save error、pending save、stale response 或 hydration/readiness 非 ready 都禁止呼叫 activation；UI 顯示可重試的 save diagnostic，不能靠 local reducer state 送出。

activation request 必須攜帶同一 workspace_revision、settings_revision 與 readiness token。backend 重新讀取 persisted state 並確認 revision、ownership、candidate snapshot 與 connector/device readiness；版本不一致回傳 modbus_share_revision_conflict 或 modbus_share_save_incomplete。只有通過 barrier 的完整 desired set 才交給 reconcile，並回傳 projection outcome、listener state 與 revision。

選擇 await barrier 而非讓 activation 自行猜測最近一次 autosave，是為了消除「UI 顯示已保存但 backend 仍是舊值」的 race；代價是 network save latency 會明確延後 activation。

### Candidate/apply/restore remains the only projection seam

Share candidate 必須由既有 source-rule candidate generation/validation 產生，apply 只接受 persisted candidate snapshot、source-rule revision 與 workspace scope；share_enabled 只決定該 rule 是否進入 desired set，不能越過 global disabled gate。runtime mapping API 不得接受 direct tag/register write 作為正式 activation；若保留 diagnostic endpoint，必須標記為 read-only 或回傳 modbus_share_projection_required。

service restart 先載入 persisted global settings 與 workspace/rule snapshots，再由同一 restore/reconcile service 產生 desired mapping set。global disabled 時 restore 不建立 listener、不寫入 runtime mapping；enabled 且資料無效時維持 listener stopped/degraded 並回報 typed error，不從瀏覽器 session 或 memory cache 恢復。

選擇重用 candidate/apply/restore seam，是為了讓既有 rule-level durable fields 與多協議 point/tag contract 只有一條 projection path；代價是 direct legacy callers 必須遷移到 workspace-scoped API，不能保留繞過 Share service 的 mutation。

### Global settings drive listener lifecycle

settings.modbus_share durable object 是 listener lifecycle 的 source of truth，至少包含 enabled、bind_address、port、slave_id 與 capacity_registers，每次更新產生 settings_revision。migration 對既有安裝顯式 materialize 一組相容值；5020 只能是資料庫中已保存或 operator 明確指定的值，不能是 handler 的無條件 fallback。

lifecycle state 使用 disabled、starting、running、stopping、failed；enabled=false 時 process startup、bootstrap、activation 與 settings save 後都不得持有 listener，且 Share mapping mutation 不執行。enable/start、disable/stop、bind address/port/slave id 更新與 reconcile 共用 serialized lifecycle lock；bind conflict、invalid address、port range、slave id range 或 capacity failure 均保留 durable setting、停止新 listener、保留舊 runtime projection並回傳 modbus_share_listener_bind_failed 或對應 typed code。

status/diagnostics 只回傳 bind state、address、port、slave id、mapping counts、error code、retryable 與 action hint，不回傳 credentials、DSN 或 raw secrets。這讓 operator 能知道「未啟用」「尚未 hydration」「bind 失敗」「需要重試」的差異，也避免固定 port 假設。

### Datatype changes invalidate old spans

每個 runtime mapping identity 包含 workspace、source rule、tag、human/zero-based register、datatype、span、stride 與 source revision。desired set 變更時，reconcile 以 identity 比對 delete/move/retype：舊 span 先在 staged memory copy 清除，再將新 span 寫入；commit 後 response 明確列出 removed、invalidated、upserted 與 observed revision。

若新範圍與舊範圍碰撞，整批維持舊 projection；若舊 mapping 的 ownership 無法證明，系統不得清除它，改回傳 modbus_share_workspace_scope、invalidated_unknown 與 recovery action。任何 invalidated/unknown 狀態都不得繼續回報 ready 或讓 activation 假裝完成，避免 SCADA 讀到未告知的 stale value。

選擇 span-aware diff 而非只更新 mapping row，是因為 retype 後舊的第二至第四 register 不會自動消失；代價是 reconcile 需要記錄 clear/invalidate outcome，並在 dirty state 時提供 recovery。

## Implementation Contract

**Observable behavior**

- /studio/v2 bootstrap 未完成或 failed 時，Share summary 顯示 blocked/unknown；不會以 display default 顯示 ready，也不能 list/delete/upsert/activate。
- global settings.modbus_share.enabled=false 時 summary 不出現、listener 不存在、list/delete/upsert/activation 不造成 durable 或 runtime mutation；rule share_enabled=true 不得越權。
- owned persisted mapping 可被同 workspace scope 的 candidate/apply/reconcile 操作；只有非空或未持久化 tag_id 的 mapping 不可授權刪除。
- 40001 轉換為 runtime register 0，每個 register 為 2 bytes，int32/int64/float span 分別為 2/4 registers，stride 不小於 span，capacity/range collision 在 backend 完整驗證。
- activate 只能在同一 workspace/settings revision 的必要 saves 全部成功後進行；同一 request 重送是 idempotent，舊 revision 會被拒絕。
- restart 由 persisted desired mappings 恢復相同 register projection；direct runtime mapping mutation 不得改變正式 Share projection。
- listener 由 durable settings 的 address/port/slave id 啟停；disabled-by-default 且 failure status 包含 typed code、retryable 與 action hint。
- move/retype/delete 不留下可被誤讀的舊 span；成功回報 removed/invalidated，無法證明 ownership 則回報 unknown 並保留待 recovery。

**Interface / data shape**

- Bootstrap/status 需攜帶 hydration_state、workspace_id、workspace_revision、settings_revision、readiness、enabled、listener_state、bind_address、port、slave_id、capacity_registers。
- Desired mapping 需攜帶 workspace_id、source_rule_id、source_rule_revision、tag_id、data_type、share_start_register（human 40001 domain）、zero_based_register、span_registers、stride_registers。
- Reconcile/activation request 需攜帶 expected_workspace_revision、expected_settings_revision 與完整 desired set；response 需攜帶 outcome（aligned、applied、disabled、failed、dirty_unknown、invalidated_unknown）、new revision、changed sets 與 diagnostics。
- Error envelope 使用 typed code、operator-safe message、retryable、workspace_revision、settings_revision、dirty_state 與 action；不回傳 secrets 或任意 raw database/transport payload。

**Failure modes**

- hydration/readiness/save barrier 不完整：拒絕 mutation/activation，維持舊 projection並回傳 modbus_share_hydration_required 或 modbus_share_save_incomplete。
- global disabled：不讀取/寫入 mapping mutation path，回傳 modbus_share_disabled。
- invalid ownership、register span、capacity、collision、datatype、settings 或 revision：整批 prevalidation fail，舊 projection不變。
- listener bind/start/stop、repository commit 或 rollback 失敗：回傳 typed failure；rollback 不確定時狀態為 dirty_unknown，recovery 可重跑 persisted desired set。
- restart 無法 restore：不從 browser/memory 猜測，listener 保持 stopped/failed，status 顯示 actionable diagnostic。

**Acceptance criteria**

- focused backend tests cover hydration/global gate, ownership proof, bytes/register conversion, multi-register span, stride/capacity/collision, atomic failure, CAS/double submit, barrier, restart restore, listener lifecycle and stale invalidation.
- frontend tests cover blocked hydration, hidden disabled summary, save-error barrier, revision conflict and operator-safe error rendering; frontend validation remains a hint only.
- HTTP integration and EXE smoke run configure→save→activate→real Modbus client read→restart→restore/read; disabled mode proves no listener and no mutation.
- spectra analyze productionize-modbus-share-lifecycle --json has no Critical/Warning findings and spectra validate productionize-modbus-share-lifecycle passes.
- field hardware acceptance is a separate gate requiring a real SCADA/Modbus client, target network/bind policy, EXE identity/hash and operator sign-off; repo tests do not close it.

**Scope boundaries**

- In scope: /studio/v2 Share API/state, existing source-rule candidate/apply/restore seam, workspace/settings/runtime projection, local Modbus listener lifecycle, schema/migration fields needed for these contracts, and focused tests. Approved implementation seams are internal/datalink/modbusshare/, internal/virtual/memory/, and internal/virtual/server/modbus/, together with the explicitly listed API/workspace/rule/runtime/frontend paths in the proposal.
- Out of scope: legacy /studio retirement, C's full i18n/UI presentation, internal/protocol/ and internal/datalink/connector/ connector/transport/frame/polling behavior, unrelated outputs, release packaging work owned by A, and field hardware acceptance.

## Risks / Trade-offs

- [Atomic reconcile adds locking and transaction coordination] → Use bounded per-workspace/lifecycle locks, deterministic desired ordering, explicit CAS conflict and idempotency tests.
- [Legacy orphan mappings cannot be safely auto-deleted] → Fail closed with invalidated_unknown, preserve state, and expose a recovery action instead of guessing ownership.
- [Persisted listener settings can conflict with existing deployments] → Materialize migration values explicitly, show bind/action diagnostics, and require operator confirmation in deployment acceptance.
- [Large numeric spans increase collision frequency] → Display span/stride/capacity in the summary, validate the complete set before mutation, and retain old projection on failure.
- [Bootstrap blocking can look like a UI regression] → Show hydration/readiness state and retry action; never present defaults as confirmed runtime truth.

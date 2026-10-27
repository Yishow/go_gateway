## Context

本 change 依基準 `913cead9` 與已執行 review fixture 規劃；測得數字只描述固定測試輸入／輸出，不是產線損失估計。`buildFrozen` 對無 `ExactRecordingSource` 的 member 讀取 current `CounterSources`；`samplePayload` 未攜帶 counter binding；`assembleBoundary` 先建立當下 binding 再 `Replay`。現有 `TestFrozenExactRecordingSourceDoesNotRefetchDefinition` 已表達已縮放 source 不重查 editable definition 的正確方向，本 change 延伸到未縮放 cumulative 與 instantaneous counter_raw。

相關既有規格：`durable-recording-delivery` 的 revision-safe apply、evidence retention，以及 `circuit-recording-profiles` 的同 epoch 比較／原子 baseline。新增要求只補足 acceptance proof 及舊 journal 無法證明時的明確處置，不重寫既有演算法。CLI context 的舊 Go／Viper 描述不是版本來源；實作遵守 AGENTS.md、go.mod 與 production wiring。

## Goals / Non-Goals

**Goals:**

- In scope：所有既有支援 layout 中 counter member 的 apply-time semantic proof、acceptance／restart／reconcile／closure／retry 相容性與隔離診斷。
- 已接受歷史從原 immutable proof 決定意義；新 epoch 由明確新 revision 生效，不影響已接受資料。
- 未能證明舊 binding 的資料完整保留並可查明範圍；合法已 finalized effects 不因缺少重算證據而被阻擋。

**Non-Goals:**

- Out of scope：修改 counter 算法、精度規則、量測編輯 UX、歷史修補工具、對外資料庫 migration／刪表／清史。
- 新 layout／profile 轉換矩陣與 successor generation 由 `manage-recording-lifecycle-safely` 負責；cutover ownership barrier 由 `fix-recording-cutover-boundaries` 負責。
- 不將既有未縮放來源硬改成 scaled source；不新增平行 writer 或以拒絕所有量測編輯迴避 replay 問題。

## Decisions

### 凍結完整 counter proof

在 `WriteGroupAppliedSnapshot` 的 runtime proof 增加版本化 counter binding map，鍵為完整 member identity，不只 Tag label。內容涵蓋 workspace/group/applied revision、measurement identity、definition/source-binding revisions、series epoch、counter kind、unit、source/mapping revisions、exact scale（含未縮放明確值）、source/result type、rollover modulus 與 max-rate policy。canonical serialized proof 具有 version 與 digest；保存到 `wg_runtime_versions` 的同一 apply transaction，驗證量測 revisions 仍匹配確認資料，失敗則 apply 不成功。

`ExactRecordingSource` 使用既有資料轉成同一 proof 語意；不得因此重新查詢 definition。其他 counter 在 apply 時取得凍結 facts，之後 `buildFrozen`／replay 不再呼叫 current definition provider。替代方案「每次 restart 比較目前單位」無法回復原意；「只存 epoch」遺漏單位與 rollover policy，均拒絕。

### Acceptance 與 proof identity 原子關聯

每個新 accepted journal entry 在原 group/revision key 下攜帶 proof version/digest reference；reference 指向先前持久化且 immutable 的 snapshot。`AppendSample` 成功才回 ACK，缺少 proof 不得 ACK。新 codec 以明確版本將 proof reference 納入新 identity digest；舊 codec digest 保持可驗證，不能因新增欄位讓相同舊 sample 重送被誤判內容衝突。

sample 自帶 measurement/epoch claims 必須與凍結 proof 一致；空 claims 不等於可重新取 current definition，僅能由已 applied member proof 補足既定意義。與新 epoch 不符的即時輸入回傳 scoped typed refusal，提示重新檢查及 apply，不能 silently label 成舊或新 epoch。選用完整 snapshot reference 而非每筆重複所有 policy，以避免兩套可變 truth。

### Legacy proof 分類遷移與隔離

Expand：新 reader 同時辨識舊 codec 及新 proof format，先部署安全讀取與缺證據阻擋。Migrate：只有原 immutable snapshot／既有 exact proof 能證明該 accepted set 的全部 facts 時，建立附加 proof reference／migration audit；保留舊 payload 與 digest。單憑 current definition、相同名稱／單位或不完整 baseline 不構成證明。Contract：所有未 finalized counter replay 都必須通過 proof gate，移除 current-definition fallback。

無法證明、proof 矛盾、版本不支援，分別以 `counter-replay-proof-missing`、`counter-replay-proof-conflict`、`counter-replay-proof-unsupported` 隔離 group/revision；資料、checkpoint、pending count 與 frozen destination 不變。安全 status 包含 affected group/revision、pending count、reason/action；action 指示保留資料、恢復該 revision 的可信證據或聯絡維護者，不提供一鍵「套用目前設定」或自動刪除。已 finalized outbox 保持原 bytes/destination，與 replay quarantine 分開。

### 重播與新 epoch 的驗收隔離

把 replay invariant 與新 revision 的 baseline reset 分開測試。固定回歸輸入 100、105 在 epoch-1/kWh 下，編輯為 epoch-2/MWh 後仍輸出舊 delta 5 kWh；另以新明確 applied revision 接受新 epoch sample，首筆 NULL delta／baseline。相同 proof 的 crash/retry 不重算 finalized payload；不同 proof 不接續舊 baseline。

## Implementation Contract

- Observable：量測編輯不改變已接受 journal／output 意義；恢復不依賴 editable repository 在線；無證據的舊資料顯示安全阻擋而非成功寫入或消失。
- Data shape：版本化 counter proof 屬 applied runtime snapshot；journal 只 reference immutable proof，decode 明確辨識新舊 format；對外不暴露 DSN、raw exception 或樣本數值。
- Failure：proof persistence／revision CAS／digest 驗證失敗時不接受新 sample、不消耗舊 input、不推進 checkpoint；保留既有 output identity。definition 暫時不可讀對已凍結 replay 不構成失敗。
- Scope：既有 managed long、managed wide profile 的 cumulative／counter_raw，及已縮放／未縮放、整數／浮點／exact decimal 的已支援組合；不新增原本不支援的組合。
- Acceptance：在 `counter_binding_replay_test.go` 覆蓋 edit/delete/outage、claims mismatch、crash 前後；在 `counter_binding_proof_test.go` 覆蓋新舊 digest、tamper/unknown version/quarantine、retry；公開 measurement PUT→restart→SQLite readback 覆蓋原失敗路徑，斷言 unit、epoch、delta、baseline 及 immutable snapshot。測試命令選到實際案例且 skip 不算通過。
- Cross-change ordering：本 change 無同批 prerequisite；`fix-recording-cutover-boundaries` 依賴本 change，避免 runtime/recovery 檔案平行修改；`manage-recording-lifecycle-safely` 再依賴兩者。`expose-runtime-recording-rules`／`make-runtime-logs-actionable` 只能讀本 change 提供的安全 reason/action，不負責修復資料。
- Rollback：新格式一旦寫入，不允許舊 binary 以無版本 fallback 重播；回退需停止 affected intake、保留 durable store，使用能辨識 proof format 的版本。不得回退 schema 而刪除證據。

## Risks / Trade-offs

- [舊 accepted data 缺完整證據而阻擋] → 明確列為安全相容性影響；保留原 journal、scope 診斷及人工修復指引，不承諾可由 current definition 自動救回。
- [新增 proof 讓 digest／重送相容性改變] → codec 版本化、舊 payload byte-preserving fixtures、新舊 resend identity 測試。
- [apply 與 measurement edit 競爭] → 同 transaction pinned revision 驗證；測試只有一個一致的 binding 能被 applied，無半套 proof。
- [已縮放來源重複轉換] → existing exact source proof adapter 與未縮放 proof 分流建構，golden tests 鎖住 raw/result types 與 digits。
- [影響其他群組／目的地] → quarantine 只阻擋 affected group/revision 的未 finalized replay；healthy groups 與 finalized payload delivery 各有回歸。

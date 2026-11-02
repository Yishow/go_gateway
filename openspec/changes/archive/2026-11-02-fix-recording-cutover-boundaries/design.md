## Context

基準 `913cead9` 的 `applyCutoverInTx` 已保存 cutover／revision end，並檢查 future boundary；`Pipeline.AcceptSample` 則向目前記憶體 boundaries 廣播，`reconcileGroup` 延後才套用 `SetEndsAt`。`Assembler.SetEndsAt` 只改 config，既有 selector 與 replay 中的跨界 observation 仍可能影響輸出。已重現的固定 fixture：邊界 22.5 秒、22.499 秒值 6、23 秒值 7，舊 SQLite row 錯寫 7；另 long row end 12.5 秒卻選到 13 秒值 200。這些是測試觀察值，不是產線資料估計。

既有 `long-table-group-recording` 已承諾不可雙寫、half-open ownership、partial boundary、pending rows frozen destination；本 change 補齊 intake／closure 競爭與未 finalized journal 交接的可測契約。`window-average-recording` 的 actual window bounds 與 circuit profile counter baseline 契約繼續有效。

## Goals / Non-Goals

**Goals:**

- In scope：目前已允許的 forward cutover 與 rollback，含 source 的既有 legacy/custom wide、managed wide profile，及既有 `managedCutoverTarget` 已允許的 managed long／managed wide profile targets；不開放任何新轉換組合。
- 按 observed_at 決定唯一 owner，並使 durable acceptance、cutover commit、closure 看見同一持久化 owner epoch。
- 已接受未 finalized samples 可驗證移交、不丟失、不雙計；latest、mean、state/counter_raw、cumulative 與 replay 都只使用 actual bounds 的 eligible set。

**Non-Goals:**

- Out of scope：新增 wide↔long／profile 轉換矩陣、successor generation、一般刪除／停用產品流程；全部由 `manage-recording-lifecycle-safely` 定義。
- 不回填、不改寫已 finalized row／receipt、不刪除外部表與歷史、不新增 retention、排程器或平行 writer。
- 不以縮短 reconcile interval、只標 partial、只改 metadata 或濾掉錯歸樣本代替正確交接。

## Decisions

### 既有切換矩陣界線澄清

原 Goals 的 target long 單一列舉不完整；本段只澄清「目前已允許」的原有範圍，不新增轉換能力。已讀 `write_group_profile.go:managedCutoverTarget`，其既有 predicate 同時允許 managed long 與 managed wide profile；`reviewForwardTargets` 沿用該 predicate，`TestReviewedWideMeanTransition` 已有 wide-profile target／rollback 驗收。實作須保留這些現有能力；新的 layout/profile 轉換矩陣仍屬後續 lifecycle change。本澄清不改任何已完成任務或產品 spec 語意。

### 持久化 ownership fence 統一線性化點

沿用本機設定／delivery store 的 transaction authority，增加 owner-fence revision 與 cutover intent 的明確關聯；`applyCutoverInTx`、`AppendSample` 與 `CommitClosure` 在同一持久化 authority 驗證 owner epoch。不能僅靠 process-local mutex 或 reconcile cache。cutover commit 包含 source end、target start、受影響來源集合與 fence revision；acceptance／closure 取得過期 fence 必須重新解析／重試或安全拒絕，不得以舊 cache commit。

cutover 仍要求 boundary 在現在之後，但另檢查已 finalized 範圍：若先前 accepted future observations 已導致與候選 boundary 衝突的 finalized effects，拒絕該 cutover 並要求重新 review 新邊界，不改寫已 finalized 結果。成功 commit 即使 runtime 尚未 reconcile，持久化 authority 已阻止舊 owner 接收或關閉跨界 input。新 owner 尚未建好時，持久化 pending handoff 及其原 sample identity 才能 ACK；否則回傳可重試的未接受，不宣稱已記錄。

### 並行階段與唯一寫入責任

本段調整未完成工作的排程與責任，不改 spec 可觀察要求或已完成 1.1／其 provenance。1.2 與 1.3 都依賴 1.1，彼此可並行；2.1 必須等待 1.2、1.3 都完成。3.1 仍依賴 2.4，所有 pending handoff／closure guards、4.3 整合驗收及 lifecycle 前置 gate 不變。

1.3 僅負責 workspace cutover-side finalized-overlap transaction gate：在 `applyCutoverInTx` 的同一 SQL transaction，檢查 `plan.ended` 的所有結束 revisions（含既有 rollback／retained successor），覆蓋已 finalized future observation 與 closure 先取得寫入 authority 的 closure-wins 次序，衝突回傳安全 `recording-cutover-finalized-overlap`／fresh-review，失敗不更改 source revision、baseline 或既有 effects。可讀既有 checkpoint／effect authority，無需改 Store／helper API。`groupdelivery.Store.inTx` 先取得 SQLite write authority；`Store.Reclaim` 可回收 consumed samples、完成 outbox 與 bucket history，但 `wg_delivery_checkpoints.next_close` 保留 first-unclosed cursor。必測 post-reclaim 與合法 boundary equality／nonoverlap，不能以缺少 rows 推論無 finalized history，也不能拒絕所有 checkpoint。

1.3 不交付完整 intake／closure barrier。反方向 cutover-wins 的 stale `CommitClosure` 拒絕／重試屬 3.1：cutover 先 commit 時，不得推進 checkpoint、consumed set、baseline 或 outbox；3.1 同時保留 pending handoff 與 eligible-input guards。完整雙次序 deterministic races 必須在 3.1／4.3 驗證，不以 1.3 局部成功替代整合 gate。

唯一寫入責任：active 1.2 lead 與其指定 helper 擁有 groupdelivery Store／新 ownership modules、grouppipeline 與 runtime acceptance；1.3 worker 僅擁有 `internal/datalink/workspace/write_group_cutover_apply.go`、新增 `write_group_cutover_finalized.go`／focused test，以及 `internal/api/handlers/studio_v2_workspace_write_group_errors.go`／新增 focused API test。每個共享欄位／檔案由一位指定 owner 寫入；workspace error declaration 若需額外檔案，先交由協調者確認唯一 owner。1.3 不改 schema、新 endpoint 或 Closure type，不擴充切換矩陣；其他 workers 不同時改其檔案，整合時再由 lead 驗證雙側契約。

2.1 的唯一寫入 owner 為 transfer lead 及其指定 helper；其負責 Store transfer／verified receipt、runtime replay 授權、transfer-specific closure guard 與 callable adapter，其他 workers 不同時修改共享介面。2.2 接手 live handoff／Reconcile 啟用及 source／target eligible-set 重建；每個共享欄位／檔案仍由協調者指定唯一 owner。此為原任務內 expand → migrate 的責任澄清，保留已完成 1.1／1.2／1.3、active 2.1 的原描述與 provenance、1.2／1.3 並行及 2.1 join 依賴，不新增產品範圍或 task。

### 半開窗 filter 與原子 durable handoff

ownership 以 workspace／source／destination-ownership scope／profile kind 決定；既有獨立 instantaneous/cumulative profiles 不互相排斥。cutover 的 moved source 在 observed_at < boundary 屬舊 owner，等於或大於 boundary 屬新 owner；未搬移 source 保留既有 reviewed successor 路徑。新舊 group revision identity 不等於 sample identity，不能靠複製 row 清除衝突。

已在舊 journal 的未 finalized 跨界 sample，以同一 local transaction 保存 target acceptance 或 pending handoff、舊端 transferred disposition、cutover identity、原 sample ID/digest 與兩端 owner refs。原 payload／原 proof 保留，不刪記錄；retry 依穩定 transfer identity 不重複移交。轉交僅改 owner，不改 accepted measurement 意義；新 owner immutable proof 不相容時隔離 handoff，禁止重新貼 epoch／單位。`freeze-counter-replay-semantics` 是此 proof 的先決條件。

2.1 是安全 expand，不是 pending-only：提供可呼叫的原子 transfer API／adapter，在同一 transaction 保存 source disposition、target exact-byte acceptance 或 pending／quarantined receipt、stable transfer identity；保留原 payload、digest、proof 及 `created_at`。相容者須經 verified `Store.Restore` 與 runtime replay 授權，才能由新 worker／restart 消費；不相容或 unknown evidence 只能 durable blocked。不得經普通 target `AppendSample` 繞過原 group-bound proof，不重讀 current definition、不改貼 target proof；pending receipt 不是 target accepted，更不代表外部 SQL committed。

target acceptance 採同一 journal 的明確邏輯 namespace，而非 magic prefix 或另一套 writer：`acceptance_kind` 區分 native／handoff，`sample_id` 是內部 physical address，`origin_sample_id` 及 origin scope 保存原 logical identity；歷史 native rows 保留所有舊欄位／bytes，空 origin metadata 依原 sample_id 解讀。native logical identity 與 handoff origin 各有獨立唯一性限制，任意 native ID 與內部 physical address 相同時可配置另一 storage address，但不拒絕該 logical ID、不改原 payload／digest／created_at／seq，也不放寬普通 native duplicate 或跨 workspace 拒絕。所有 reader、association、predecessor 比較使用經驗證的 logical identity。

Snapshot 的 `AcceptanceOrigin` 只作處理器記憶體中的 tagged composite seen key（origin scope＋原 SampleID），不進 codec、hash 或外部 JSON，也不是持久化授權。普通 `AppendSample` 拒絕非零 annotation；只有 verified `Restore`／recognition 可衍生它，私有授權同時綁 origin tuple、原 digest 及 target proof。必測 A 的 s-1@23 移入 B 後，B native s-1@24／35 仍是不同觀測，原觀測重送仍去重；transitive origin、I/C scopes、counter readings、GC／restart 及任意 physical-ID 碰撞不得改變原資料。

2.1 同時提供 opaque `HandoffView`／receipt-disposition guard：input set、exclusion set 及 receipt view 必須在同一 read transaction 捕捉並綁定 assembler；closure 只有在該 view 仍有效時才能前進，source 或 target assembler stale 均安全拒絕，checkpoint／consumed set／baseline／outbox 不前進。這是 transfer 特定保護，不替代 3.1 的一般 owner epoch fence 或完整雙次序 closure races。2.1 的驗證包含 atomicity／idempotency、相容 exact-byte acceptance 的新 worker／restart replay、blocked receipt 及雙側 stale view 安全拒絕；正常完成後 pins／`Reclaim` 沿既有政策釋放／回收，不永久保存所有 payload。

分期原因：`groupdelivery.Store.Restore` 會驗證 sample digest 與 group-scoped proof；runtime `assembleBoundary`／`validateCounterReplay` 再核對 immutable proof。source-scope counter proof 不能靠 target copy 變成有效 target input，必須 verified receipt 授權。只改 journal 而不重建兩端 assembler，可能使舊端仍含 transferred input，或新端漏 input 而產生錯誤／empty closure。因此 2.1 不自動掛進 live `Reconcile`；2.2 才把 live handoff／Reconcile 與 source／target eligible-set 重建、receipt view 刷新一起啟用，兩端完成一致重建前不可解除 guard。2.3 保留全部 reducer 矩陣，2.4 保留完整 immutable proof 衝突矩陣，3.1 及後續整合驗收不縮減。

舊 selector 必須由符合 actual bounds 的未 finalized accepted set 重建；不能只刪最新選中值而遺失更早有效候選。相同過濾套用 latest、mean 的 sum/count、state/counter_raw，以及 cumulative 全部 pre-decimation observation；不把下一 owner 的 observation 用作舊 baseline。設計選擇保留 journal 與 transfer audit，而非兩邊重新 append 無關樣本，避免 drop／double count。

### Closure 與 recovery 不跨 ownership 進度

closure transaction 同時檢查 fence、未處理 handoff 與 consumed-input set；過期 plan 不 commit、不推進 checkpoint。只在正確 owner 已接收或 pending handoff 有耐久憑證後，舊 owner 才能結算自己的 eligible set。新 owner 不得先關閉仍有按時已接受但待移交的 bucket。replay 根據持久化 disposition 與 proof 重建 owner，不再按當前 cache 猜測。

實際 window bounds 同時控制選樣與 row/provenance；邊界截窗標 partial，新 cumulative owner 首筆為 baseline／NULL delta。原 accepted timestamp 決定 handoff 是否在 closure deadline 前被接受，reconcile／restart 的晚到不能將已接受 sample 降級丟棄；真正首次接收已超過 deadline 的輸入仍按既有 late-after-close 契約診斷，不重開 finalized window。已 finalized 舊 payload 正常 drain 到 frozen destination。

**3.3 期限範例時序澄清（取代原範例的混合敘述）**：本決策僅 supersedes `Accepted before deadline fixture` 原先把第 14 秒 acceptance 與第 20 秒 pending handoff 混為一談的 example wording；原 requirement／scenario、已完成 3.3 與其證據保持不變，不取消要求、不新增 task、轉換能力或 late policy。

採用下列兩條合法時序，而非讓 source 在第 14 秒接受後再回溯套用第 10 秒 boundary：`planBoundary` 要求 future boundary，commit 後的 `validateCutoverAcceptance` 也會拒絕跨界舊 owner。固定測試條件與證據見 `.task-evidence/cutover-3.3-lead/README.md`。

- 原第 14→20 秒範例保留：`internal/datalink/grouppipeline/cutover_reconcile_deadline_test.go:TestCutoverDeadlineOriginalAccepted14RemainsEligibleAfterReopen20` 先在第 9 秒 commit 第 10 秒 boundary；successor 原生接受 observed_at 11、ReceivedAt 13、Store `created_at` 14 的值 7。window [10,15)、deadline 16 的一般 recovery／reconciliation 在 physical reopen 20 後保留原 acceptance 資格與值，不宣稱此路徑有 pending handoff。
- 真正 delayed handoff 另列範例：`internal/datalink/runtime/group_boundary_cutover_deadline_test.go:TestGroupBoundaryCutoverDeadlineDelayedHandoffRemainsEligible` 在第 9 秒由 source 接受 observed_at 11 的值 200，第 9.5 秒 commit future boundary 10；deadline 16 的 closure 因未完成 handoff 被阻擋，physical reopen／handoff 20 後保留 accepted_at 9 與原值，不重新解讀為 late。首次到達 17 的獨立 late fixture 仍以實際 finalize 16 為前提；不因 wall clock 已過 deadline 就拒絕尚未關閉 window 的輸入。

### 以 deterministic 競爭 fixture 驗證 cutover

測試使用受控 clock／barrier，不靠 sleep 猜測 race。固定 cases 包含先接受 future sample 後 cutover commit、commit 後 reconcile 前 sample、同時 closure 與 commit、邊界前／等於／後、replay／rollback。驗證 journal 與 transfer audit、consumed ownership、outbox effect identity、SQLite 實際 row 的 observed_at／value／delta／bounds，不只看 metadata 或 HTTP 成功。

## Implementation Contract

- Behavior：successful cutover 的唯一歸屬在 persisted boundary 即成立，不等待 reconcile；失敗 cutover 不改 source owner、baseline 或已接受資料。
- Data shape：owner-fence revision 綁定既有 cutover/revision-end；journal disposition 支援 active、transferred、pending handoff／quarantined 的明確等價表示，記錄 source/target scope、sample ID/digest、proof reference 與 transfer ID；不得以刪 source entry 表示完成。
- Failure：owner authority 不可讀或 fence stale 時 typed `recording-ownership-unavailable`／`recording-ownership-conflict`；proof 不符為 `recording-handoff-proof-conflict`；finalized overlap 為 `recording-cutover-finalized-overlap`。一般 operator response 只顯示 scoped reason、pending 狀態與 retry/review action，不含 DSN 或 raw value。
- In scope：已支援切換矩陣內的所有 profile kinds、saved latest/mean 與 completeness policy；不為通過測試開放新矩陣。
- Acceptance：固定 fixture 舊值是 6／100，新值是 7／200；恰好 boundary 的 sample 只歸新 owner。測試 old/new journal＋SQL union 等於 eligible accepted set，交集為空；counter/new epoch／baseline 與 proof 保持相容。
- Verification：`cutover_ownership_test.go` 與 `group_boundary_cutover_replay_test.go` 覆蓋 race/crash/handoff；`cutover_eligibility_test.go` 覆蓋各 profile 的選樣與 actual bounds；公開 cutover-review/apply→延遲 reconcile→SQLite readback fixture 證明實際 wiring。新增測試先失敗、再最小修復、最後重構與 race gate。
- Cross-change ordering：必須先完成 `freeze-counter-replay-semantics` 的 proof／legacy gate，再開始本 change；`manage-recording-lifecycle-safely` 依賴本 change 且只擴充新的 lifecycle intents，不另造第二套 boundary engine。
- Migration：expand fence/disposition reader＋安全 gate，再遷移可證明的 open journal，最後關閉無 fence 的舊 acceptance／closure fallback；每一批保持可建置。未知歸屬保留並隔離，已 finalized data 絕不重新分配。舊 binary 不得略過新 disposition，rollback 先停止 affected intake 並保留 store。

**4.2：既有 rollback／獨立 profiles 的同 origin 匯流契約**

本段只細化未完成 4.2 的 scope propagation／rollback 驗收，不新增產品 requirement、task、轉換組合或 owner authority，也不改已完成 history／provenance。既有 `HandoffView`／`OwnershipView` 繼續是 transfer／ownership 的授權依據，沒有第二套 writer／owner model；不同 I/C output scopes 仍各自獨立。

- 已觀察失敗：`.task-evidence/cutover-4.2-lead/02-convergence.log` 與 `08-convergence-exact-red.log` 的公開 fixture 從首次 apply 即有 immutable measurement／epoch／unit proof 的 managed-wide instantaneous `counter_raw` source A/R1 開始，source／target 使用不同 canonical circuit 及真實 SQLite connectors，屬既有支援的獨立 ownership scopes。固定測試條件是第 8 秒接受 observed_at 第 25 秒的觀測，第 9 秒 commit 第 10 秒 forward boundary 至 I/C；兩支都攜帶同一 A/R1 origin。第 19 秒 commit 第 20 秒 rollback boundary 至 A/R3，第一支產生 target seq 4，第二支報 `ErrOwnershipConflict`。這是既有「Rollback and independent profiles」要求的缺口；部分搬移與 native I/C 的既有通過結果不代表此匯流已通過。
- Canonical journal：完整 `(target workspace/group/revision scope, origin workspace/group/revision/originSeq, original sample identity)` 只對應一個 target journal；多條各自驗證的 transfer edges 能引用同一 canonical target seq。原 sample identity 包含原 SampleID、member 與 observed_at；不同 origins 即使 SampleID 相同也不得合併。重用既有完整 origin/native uniqueness，不以放寬 journal identity、刪 row、複製新 ID 或重寫 proof 消除衝突。
- 逐支驗證：所有共享 target seq 的 peers 都必須核對完整 origin、原 identity、原 payload bytes／digest／raw digest、原 proof／`created_at`、target scope／bucket，以及 original/target proof references 與 origin/target runtime digests；任何不符都 fail closed。各 edge 的 Source／SourceSeq、Fence、SourceBucket、Targets 是合法且獨立的 branch facts，須依各自 persisted authority 驗證，不要求跨 branch 相等。只有所有相關 peer 都成立才授權共享 target，不能因先找到相容支線就忽略另一條 corrupt 或未授權 peer。
- Audit shape：保留 canonical target row 的原 payload、digest、proof、`created_at`、seq 與不可變 `handoff_id`。其他 alias edges 在既有 receipt JSON 增加 optional immutable `CanonicalTransferID`，直接指向既有 canonical transfer audit；禁止指向 alias、循環或不存在的 audit。canonical edge 留空此欄位，舊 v1 receipt 以省略空欄位保持原 IDs、serialized bytes 與 hash 相容；新增 reference 只能在獨立驗證該 edge 與 canonical 完整相符後保存，既有 edge 不重編 ID 或改其歷史 bytes。
- Reader authority：live canonical 必須由 target row 的 `handoff_id` 選出，再驗證該 canonical、所有 peers 與各 branch lineage，絕不依 SQL 第一列或任意排序選擇。`Store.Restore`、`RecognizeTransferredSample`、transitive replay、closure guard 及正常 reclaim 使用同一驗證結果；不得只驗 canonical 就漏掉壞 alias，亦不得把未授權 peer 忽略為無關 receipt。
- GC／failure：正常 reclaim 後共享 target raw row 已不在時，使用 compact immutable meaning、每條 terminal edge 及持久化 canonical reference 驗證；缺 canonical、reference 不符或 corrupt peer 一律 fail closed，不以缺 raw row 推論成功。第二支失敗須原子回滾該支，保留已成功第一支、全部仍需的資料／pins，阻擋受影響 closure，修復後 retry 能前進；不能偽造 consumed／finalized、永久拒絕合法觀測，或以永久保留 payload 迴避既有 GC。
- Acceptance（固定 fixture 的預期值，非已通過宣稱）：公開 `TestCutoverPublicCompatibleFanoutRollbackConvergesOneOrigin` 完成 rollback reconcile 後應有 4 個 journals／4 條 receipts，其中兩條 rollback edges 共用一個 canonical target seq；原始、I、C、A/R3 的原 bytes／proof／acceptance time 不變，只有 A/R3 消費一次並產生既有 actual-bounds SQL 結果。另以 mixed branches 的 proof 衝突、不同 origin 同 SampleID 不合併、第二支 transaction 故障後 retry、physical close/reopen／正常 GC 後重試不增 effects、alias／canonical 缺失或損壞後 Restore／Recognize 安全拒絕驗證，不能只驗 metadata 或 pending 狀態。
- Expand → migrate → contract：先擴充 optional reference 及兼容 v1 的 reader／peer 驗證；再以 anticipated `046_write_group_handoff_convergence_sqlite.up.sql` 僅將 receipt 的 `idx_wg_handoff_target_seq` unique index 改為 nonunique，保留所有 target full-origin／native unique indexes、原 rows／IDs／payload／proof；最後啟用經驗證的 coalescing acceptance。各批保持可建置、相關測試通過，未受影響的舊 records 完全不變；migration 要驗證冪等 upgrade、physical reopen，以及注入 migration fault 時原子 rollback 後可安全重試，不靠刪資料或鬆綁其他 constraints。
- Scope boundary：這是既有 SQLite local receipt 表及公開已支援切換路徑的實作澄清，不新增外部 DB／測試環境、DB/layout/protocol conversion、retention policy 或 current-definition relabel。保留既有 deployment rollback 的 affected-intake 停止／store 保全要求；舊 reader 不得略過新 alias authorization。

## Risks / Trade-offs

- [ownership 查核增加 local transaction 成本] → scope-indexed fence、只重建受影響 bucket，維持既有 bounded backlog／per-tick 工作限制；量測基準與 contention 後再調整，不預先宣稱效能提升。
- [跨 group handoff 一半成功] → source disposition、target acceptance／pending intent 與 transfer identity 同 transaction，逐個 crash point 測試 rollback／replay。
- [修正只覆蓋新樣本、漏既有 buffer] → mandatory seeded legacy/open-journal cases，同時斷言選樣結果與 old/new ledger union／intersection。
- [邊界與 future-skew／lateness 交互] → observed_at 決定 owner，原 durable acceptance time 決定 deadline；finalized overlap 拒絕 cutover，不更改已交付歷史。
- [跨 profile 誤判重複 writer] → owner scope 保留 profile kind，測試 instantaneous/cumulative 的合法雙 profile 與同 kind 唯一 owner。

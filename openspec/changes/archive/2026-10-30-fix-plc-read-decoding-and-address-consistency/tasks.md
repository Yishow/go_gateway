## 1. 共用解碼

- [x] 1.1 完成「共用 byte order 解碼保留既有 16 位元與預設值」及 `Multiword numeric decoding follows all configured byte orders`：先在 `lib/hsllogic/plc_byte_order_regression_test.go` 用獨立 known-bit vectors 重現 BADC/DCBA 的 32／64 位元失敗，再讓 register-to-value 重用既有 byte transform；驗證四種 order 的 signed/unsigned/float、完整 uint64、ABCD/CDAB 與單 word 結果，保存改前失敗及改後通過證據。
- [x] 1.2 [after: 1.1] 在 `internal/datalink/connector/adapters/plc_read_decoding_regression_test.go` 用獨立本機 wire fixtures 驗證 Modbus 與 MC 讀取均遵守共用解碼，含 float32(1)、float64(1)、帶符號值及完整 64 位元；執行 adapters focused tests 並回歸 point override／connection default 與 direct MC typed batch，不能只呼叫 converter。

## 2. FATEK 完整值與失敗回應

- [x] 2.1 [after: 1.2] 完成「FATEK typed scalar 與原生元件 count 分開處理」中的完整 typed 解碼，對應 `FATEK typed register decoding preserves complete values`：先重現 D0／DR0 的 float32 與 64-bit 首值 fallback，再使 `convertFatekValue` 回 `(any, error)` 並由 `FatekConnector.Read` 傳遞失敗；在 adapters regression tests 驗證 high-word-first bits、Go 回傳型別、int64/uint64/float64 與不足資料 bad quality，不新增 FATEK endian 功能。
- [x] 2.2 [after: 2.1] 完成 `Scalar point reads preserve protocol count units`：在 `point_read.go` 增加 address-aware count helper，遷移 `pollPointWithAcquisition`、`PollNowContext`、`ReadPointNow`；先用 `point_count_contract_test.go` 的 FATEK 真實 adapter fixture 證明原生 32-bit 多讀，再驗證 D0 float32/64-bit count 2/4、DR0 count 1/2、無效位址不送 wire，以及 Modbus、MC、FATEK direct component Count 不變。
- [x] 2.3 [after: 2.1] 完成「FATEK Cmd 44 嚴格驗證後才形成成功值」及 `FATEK discrete payloads are validated before success`：先在 `internal/protocol/fatek/status_payload_test.go` 重現空／短／長／非 binary payload，再驗證 `ReadStatus` 與 adapter 回 error、bad quality、無部分成功值；涵蓋合法 scalar、batch、count 1/255 及拒絕 0/256，保持 Cmd46 長度與 LRC 回歸通過。此項與 2.1 共用 adapter regression 檔，須依序修改。

## 3. MC 位址與既有規則保護

- [x] 3.1 完成「MC 舊規則更新在任何 mutation 前阻擋隱含重建」及 `MC persisted source identities survive radix correction`：先在 `internal/datalink/sourcerule/mc_radix_compatibility_test.go` 建立已儲存 W0…W16、完整 point/tag/mapping 的 fixture，驗證守門邏輯能辨識 corrected stored-range mismatch；於 `Service.Update` 的任何 cleanup／write 前回傳 `mc_address_plan_conflict`，以 repository call spy 與前後 snapshot 證明 name/scale/range patch 皆不改 IDs、位址、關聯及 revisions，並驗證 startup/reload 不重新展開。planner 尚未切換時以測試中的獨立 corrected-plan 向量驗證 guard，不先修改 production parser 來製造失敗。
- [x] 3.2 [after: 3.1] 完成「MC 新規則採現行 runtime 可確定解析的 radix」、`Protocol-aware point address derivation` 及 `Protocol-radix source address validation and expansion` 的 backend 垂直修復：先重現 W0→WA、WA/BFF start、X0+16 的錯誤，再同步 backend planner 與 point 詞形驗證；在 sourcerule planner/create/update 測試驗證 W0 count17 可完整建立、X20/Y10/B100、invalid X8/YA/WG/Z0 無寫入；point 測試只驗證 WA/BFF 等合法詞形可儲存，不擴充 generic point API 的協議語意。用不變的 runtime parser 核對實際 offset。不得新增 MC area 或把任意字串放行。
- [x] 3.3 [after: 3.2] 讓前端 `AddressParser` 的 parse/format/offset/expand 與 backend 新規則契約相同：先在 `frontend/tests/unit/utils/addressParser.test.ts` 更新獨立規格預期取得改前失敗，再驗證 X/Y octal、W/B hex、D/M decimal、word/bit 類別與 all-letter hex suffix，保留 Modbus/FATEK 行為。舊 hex X/Y 測試須依本 change 的 MODIFIED delta 改寫，不為使測試通過而改 runtime grammar。
- [x] 3.4 [after: 3.3] 完成「MC 已儲存點位以 scoped snapshot 作為前端事實來源」及 `Persisted MC points remain authoritative in Step 2` 的 hydration/autosave：用新 `mcPersistedPoints` adapter 與 `useStudioV2PersistedMCPoints` 沿用既有 candidates GET，傳遞 workspace/rule revisions，新 hook 自有包含 rule ID、workspace_id、expected_workspace_revision、revision_id 的 key，不沿用現有 unscoped `useSourceRuleCandidatesQuery`；保留 GET 既有 database-output snapshot refresh，不發 POST recompute/apply；在 `useStudioV2AutosaveState`、rule/mapping autosave 導入 persisted point identity 與 checking/conflict/unavailable gate。以 `mc-persisted-hydration.test.tsx` 驗證舊 W0…W16 不生成新位址、不重綁 mapping、403/409/503／缺資料／late old-scope response 都不授權 mutation，且未受影響草稿保留。
- [x] 3.5 [after: 3.4] 讓 Step2 selectors、point list、`RuleEditor`／`RangeSummary` 與 occupied/conflict 檢查使用同一 persisted MC point 來源；以 `mc-persisted-view.test.tsx`、`step2-rule.test.tsx` 及 occupied-range 相關回歸證明 legacy mismatch 不顯示假連續範圍／無衝突，新草稿與 compatible rule 仍可操作。必要阻擋訊息進既有 en／zh-TW locales，原始錯誤與 endpoint 不出現在 operator DOM，不新增轉換 UI。

## 4. Modbus 明確位址優先序

- [x] 4.1 [after: 2.3] 完成「Modbus raw numeric 依明確 Function 使用 zero-based offset」及 `Explicit raw Modbus addresses preserve requested function`：先在 adapter regression tests 捕捉 Address100+FC03 被送成 FC01/offset99 的失敗，再修 `parseModbusAddress` 的 raw／傳統五位／alias 分類；驗證 TCP/RTU/UDP 共用路徑、read/write caller、raw 0/100/65535、overflow/負數／壞字串、傳統與 alias prefix 優先，以及未指定 Function 的歷史 fallback。不自動重寫 persisted point；交接須列出短數字+stored03 的 before/after 實際目標。

## 5. 整合、必要收斂與交付

- [x] 5.1 [after: 1.2, 2.2, 2.3, 3.5, 4.1] 完成 expand→migrate→contract：查證 shared decoder、FATEK converter、count helper 與 Step2 selector 的實際 callers，確認全部必要路徑已遷移後，才移除本次取代的重複 order／首值 fallback 分支；仍被使用的舊 API 保留。以 caller 檢查與受影響 focused tests 證明沒有新 dead branch／編譯錯誤，超長檔不增行，禁止順手清理其他模組。
- [x] 5.2 [after: 5.1] 執行 backend 交付驗證：`go test` 涵蓋 `./lib/hsllogic`、`./internal/protocol/fatek`、`./internal/datalink/connector/...`、collector、device、point、sourcerule 及相關 API source-rule 選測；以真實本機 adapter fixture 穿過 periodic/manual/read-now 驗證 malformed Cmd44 與 typed-count 錯誤沒有 quality-good 新樣本，以既有 lifecycle 測試確認映射及舊規則保護。對受影響 packages 執行 vet/lint；記錄實際命令與結果，零測試或 skip 不算通過。
- [x] 5.3 [after: 5.1] 執行前端交付驗證：addressParser、MC persisted hydration/view、Step2／occupied-range／rule與mapping autosave 相關選測，加上型別檢查與修改檔案 ESLint；以真實瀏覽器和本機 simulator 驗證新 W0/X0 規則建立、儲存重載，以及舊 mismatch rule 的列表／blocked 狀態與無 mutation request，保留截圖與網路請求證據。真 PLC、正式 DB、Windows／現場未驗須列 NOT RUN。
- [x] 5.4 [after: 5.2, 5.3] 對最終差異執行 `git diff --check`、`make check-lines`，核對 proposal/spec/design/tasks 與已執行證據；交接五類修復、兩項 breaking grammar/precedence、回復限制、範圍外項目與 NOT RUN。確認 unrelated frontend WIP、持久化資料及 point/tag/mapping 身分未被本次工作擅自修改；未獲授權不 commit／push／部署。

## 6. Review 發現的原契約缺口

以下承接本次 review 的 R1–R5，補齊既有契約，保留上方已執行的實作與驗證紀錄；不增加協議、area、API、schema 或現場驗收範圍。

- [x] 6.1 修復 R1：MC scope 進入 checking／unavailable 時，沒有 server record 的新 mapping／database target 草稿也須保留原 point ID 與使用者輸入，不被 default 值取代；先重現丟失，再驗證恢復 compatible 後仍使用該草稿，既有 persisted 及其他設備行為不變。
- [x] 6.2 修復 R2：cold MC locator／editor 以實際 mapping record 的 point_id／address 建立單點 geometry，沒有 Step2 暫存 snapshot 也不丟失已有映射；先重現回傳空值，再驗證 MC 與非 MC 入口，不能為此放寬一般 Step2 的 scoped evidence／mutation gate。
- [x] 6.3 修復 R3：新增 MC 規則的 nextFreeStartAddress 與現行 parser 使用相同 radix，X0/count8 之後為 X10，W0/count17 之後為 W11，WA／BFF 等 suffix 可延續；先重現預填 X8，再測試規則新增 caller 與非 MC fallback 回歸。
- [x] 6.4 修復 R4：MC 新 plan 在 frontend／backend 拒絕超過現有 MC3E client 0xFFFFFF 範圍的位址與 expansion overflow，不能保存 runtime 無法讀取的 WFFFFFFFF 或回繞為 W0；保留舊 direct runtime grammar，邊界輸入不得讓摘要拋例外。先重現，再測最大合法值、超界、位移與拒絕前無 persistence。瀏覽器補驗同案缺口：WA／WF／BFF 正確區分 area 與 suffix，compatible evidence 不得略過當前草稿的 start／完整範圍驗證與 Continue 阻擋。
- [x] 6.5 修復 R5：compatible MC rule 的 X8 patch 經 UpdateWithRuntimeReconcile 於 mutation 前 validation 拒絕時，不執行 snapshot restore writes；保留已 mutation 後失敗所需 rollback。先以 repository spy 重現 Update／ReplaceLinks，再驗證 write-free rejection 與既有 rollback。
- [x] 6.6 [after: 6.1, 6.2, 6.3, 6.4, 6.5] 對修復後快照重跑 review／verify、受影響測試及必要 browser 回歸，更新 validation／scenario 對照與 scope 清單；確認四個既有 WIP／index 未被誤納，準備完整 archive／commit 計畫供確認，不 push。

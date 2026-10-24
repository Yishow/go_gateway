## 1. 擴充契約並保持舊資料可讀

- [x] 1.1 落實「迴路型態與 canonical 群組」與 Stable circuit membership and independent profiles：在 WriteGroup／成員保存穩定 circuit identity、profile kind 與 canonical create key；新增 nullable metadata migration，舊群組缺欄位仍讀為 wide，重載／重送／改名不產生重複群組。定位 internal/datalink/workspace/write_group.go 與 internal/datalink/schema/migrations/；先寫 TestCircuitProfileCanonicalIdentity、TestLegacyGroupDecode 失敗案例，再實作、整理並跑 workspace focused tests。
- [x] 1.2 [after: 1.1] 實現 Bounded recording intervals independent of acquisition：新增 profile interval/unit 的 API、TS types、儲存與 hydration 契約，瞬時 1–300 秒／預設 15 秒、累積 1–600 分鐘／預設 10 分鐘，拒絕越界、小數及錯單位，保留舊值與採集頻率。以新增 TestRecordingProfileIntervalBounds、API round-trip 測試及前端 profile-boundaries.test.ts 驗證；檢查 UI 和後端不能以不同單位解讀同一值。
- [x] 1.3 [after: 1.2] 實現「長表布局與 TimescaleDB」的 plain 部分、Versioned shared typed long storage、Scoped confirmation and truthful delivery evidence 與 Stable history access without automatic deletion：新增固定 typed 長表布局、workspace/destination storage ownership、複合主鍵與 source/time index，保留群組個別確認與 stale revision 拒絕。定位 managedSchemaLayout／ManagedSchemaScope 及 internal/datalink/recordingplan/；先補 TestLongLayoutExactTypes、TestSharedStorageConfirmationScope、TestSchemaOperationLostResponse，驗證 SQLite schema/readback、新點免加欄、空 DDL 仍確認及未建 retention policy。

## 2. 建立可恢復的記錄與計量切片

- [x] 2.1 [after: 1.3] 實現「瞬時選樣與累積精確計算」中的 Deterministic per-source instantaneous records：沿用 runtime／grouppipeline 的接受與閉窗路徑，依 UTC 半開窗挑選每個來源的合格值並展開長列，missing/bad/stale 以 NULL 與原因寫出，保留真實觀察時間及晚到診斷。先新增 TestInstantaneousProfileWindowEdges、TestLongProfileIndependentMissing、TestProfileLateSample，確認失敗後實作，再跑舊 snapshot completeness 回歸；不可改動 wide 的 required/skip_row。
- [x] 2.2 [after: 2.1] 實現 Cumulative raw values and exact differences：新增不經 float64 強轉的整數／decimal 計算路徑，raw、delta、區間起訖與 provenance 一起形成 typed record，首窗建立 baseline 並輸出 NULL delta。定位 internal/datalink/aggregation/ 與 internal/datalink/measurement/；先寫 TestProfileCounterExactDelta，覆蓋 9007199254740992→9007199254740993、12000.0→12000.4→12000.7 及 float 型別標示，再實作並跑精確編解碼測試。
- [x] 2.3 [after: 2.2] 實現 Counter discontinuities and gaps remain explicit：檢查窗內 pre-decimation 序列，重置／同時刻衝突不計差值，signed 需確認，rollover 需單圈證據，缺窗保留 baseline 並標示實際 span_total，epoch／binding／profile 修訂重建基準。先寫 TestProfileCounterResetWithinWindow、TestProfileSignedAndRollover、TestProfileCounterGapAndEpoch 的固定預期值，再實作；驗證 missing 並非零與未知差值不被絶對值修補。
- [x] 2.4 [after: 2.3] 落實「耐久閉窗與缺表恢復」中的 Atomic baseline and delivery progress：將 long records、counter baseline、consumed checkpoint 與 outbox intent 接進現有同一閉窗交易及 production SampleSink，重啟與 SQL retry 使用凍結結果。定位 internal/datalink/groupdelivery/types.go 的 Closure、Pipeline.AcceptSample 及 cmd/test_ui/service_wiring.go；先寫 TestProfileClosureCrashRecovery、TestProfileOfflineRestart、TestProfileCommitAckLost，於 journal ACK、閉窗前後、SQL ACK 遺失注入故障，確認每個 effect 一列、基準只前移一次，並回歸 legacy ownership suppression。

## 3. 擴充目的地並提供明確切換

- [x] 3.1 [after: 2.4] 完成「長表布局與 TimescaleDB」的 Optional explicit TimescaleDB mode：唯讀探測 extension/function/權限，明選後經 preview/confirm 建立空 hypertable，主鍵含 bucket_start，預覽並凍結 7 日 chunk interval；無能力回 timescale_unavailable，不能裝 extension 或自動降級。先寫 TestTimescaleCapabilityGate 與 TestTimescaleConfirmedLayout，再在可丟棄 Timescale fixture 驗證 catalog／跨 chunk／重送；記錄實測版本及 image digest。
- [x] 3.2 [after: 3.1] 完成「耐久閉窗與缺表恢復」中的 Restore only previously confirmed missing managed storage：確認時保存 recreate_missing 政策，catalog 證實缺表後依已確認 generation/layout 協調重建與續送，保留 committed receipts 並顯示歷史未保證復原。先寫 TestLongStorageMissingRecovery、TestLongStorageRecoveryDenied，覆蓋 concurrent plain／hypertable 恢復、異質同名表、舊 revision、缺 owner／權限與 SQLite 檔案消失；不改造 legacy writer 修復。
- [x] 3.3 [after: 3.2] 實現「基本操作與前向切換」的 Explicit forward cutover preserves wide history：新增 reviewed/CAS cutover 與新邊界 rollback，邊界截斷窗標 partial、來源只屬一個 active revision，新累積窗重新 baseline，舊 pending rows 送原目的地。先寫 TestWideToLongCutover、TestProfileCutoverFailure、TestProfileRollbackDrain；確認舊 wide/custom 布局、歷史、間隔及 ownership 不變，失敗不雙寫。

## 4. 提供可操作的基本流程

- [x] 4.1 [after: 3.3] 實現 Four-step intent-led setup 與 Complete basic preparation with one effective configuration：在既有 BasicRecordingPanel 及 Step4Database 顯示目的地、迴路來源、雙型態／單位、schema 預覽確認與開始記錄，plain/Timescale 選擇接同一 canonical group，原始 Tag 不強迫 measurement，cumulative 缺語意才阻擋；en/zh-TW 文字與鍵盤操作同步。先新增前端 circuit-profile-basic.test.tsx 的新建／reload／未知語意／離線草稿／確認失敗案例，看到失敗後實作並跑四步流程相關回歸。
- [x] 4.2 [after: 4.1] 實現 Simple circuit choices and scoped batch edits、Explain advanced groups and workspace responsibilities 及 Basic missing-value policy reveals required gating：批次設定只影響勾選迴路，先預覽覆寫；進階顯示相同 profile、工作區只管 scoped readiness／activation；wide required 規則及 long NULL 原因分開呈現，合併相同 blocker 並提供修復入口。先寫 circuit-profile-batch.test.tsx 與 circuit-profile-status.test.tsx，驗證覆寫確認、未選保持、baseline 未知、SQL commit 與 local durable 分別顯示。
- [x] 4.3 [after: 4.2] 完成 contract 階段：新建 Basic 只保留 profile 的有效週期控制，既有 wide/custom 繼續原編輯／執行，Basic／Advanced 切換與 forward cutover 不新增第二份設定。以新增 circuit-profile-compatibility.test.tsx 及後端 TestProfileLegacyCoexistence 驗證既有 60 秒群組、custom table、Modbus-only、shared destination、新增點位空 DDL 與 stale revision；更新操作說明解釋進階群組／工作區用途，內容 review 對照本案 specs。

## 5. 驗證交付範圍與證據

- [x] 5.1 [after: 4.3] 執行可丟棄 SQLite、plain PostgreSQL、TimescaleDB 的整合矩陣：固定模擬輸入驗證 exact raw/delta、first baseline、gap/reset、跨 chunk、SQL ACK 遺失、缺表重建與 cutover/drain，對照 schema/catalog/readback 及 effect identity，確認三模式邏輯結果相同且沒有 retention/compression job。新增的整合測試名稱以 TestCircuitRecordingIntegration 為入口，執行對應 Go package 的該測試並記錄實際選中案例、版本、image digest 與 cleanup 結果；缺環境或 skip 標 NOT RUN，不勾完成。
- [x] 5.2 [after: 5.1] 使用模擬來源完成實際瀏覽器操作：選迴路、雙型態、批次覆寫、reload、越界錯誤、確認／開始、修復及切換，留下 zh-TW 截圖與 SQL readback 對照；依 docs/testing.md 執行受影響 Go package tests/vet/lint、前端 lint/tests/build、git diff --check 與 make check-lines，將 specs 每項 scenario 對應到實跑證據。交付記錄區分 PASS／FAIL／NOT RUN，不以此替代正式設備、Windows 或現場驗收，不自動 commit／archive／部署。

上述 15 項已完成；逐項 scenario、真實資料庫與瀏覽器證據見 docs/technical/evidence/circuit-recording-scenario-coverage.md、docs/technical/studio-v2-circuit-recording-db-verification.md 及 docs/technical/studio-v2-circuit-recording-browser-verification.md。提案階段的 analyze／validate 不作為產品完成證據。

每項估計 30–60 分鐘；實際起迄與 focused verification 用於驗證估計。若單項超過 60 分鐘且無可驗收切片，先停在可建置狀態並拆細該項剩餘範圍，不擴增需求。涉及演算法、交易、恢復及 UI 的任務採 red → green → refactor；每個切片結束需可建置並通過所列 focused tests。

上述 Test* 與新前端測試檔名稱是提案時指定的驗證入口；實際選中名稱、PASS／SKIP／NOT RUN 以交付證據為準。序列依賴同時涵蓋共享 WriteGroup、schema、runtime 與 Studio 檔案，避免重疊編輯。

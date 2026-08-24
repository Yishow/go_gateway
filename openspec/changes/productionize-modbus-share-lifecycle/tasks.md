## 1. hydration、全域 gate 與 ownership contract

- [ ] 1.1 [TDD-Red] 為 B1 Hydration-gated Share truth、B2 Global Share setting is an absolute gate、B3 Durable workspace ownership is required 建立 failing backend/API tests：bootstrap pending/failed、global disabled、rule flag 越權、未持久化 tag_id、cross-workspace delete/list/upsert 均必須產生指定 typed error 且無 mutation；以 internal/api/handlers/*modbus_share*_test.go、internal/datalink/sourcerule/*_test.go 的 focused test names 與 mutation assertions 作為驗證目標。
- [ ] 1.2 [P] 實作 Backend owns the truth and browser hydration gates every Share operation：bootstrap/status 帶 hydration_state、workspace/settings revision/readiness，typed error envelope 與 global disabled gate 在所有 Share list/delete/upsert/activation/restore 路徑生效；以 handler route tests 驗證 display defaults 不會進入 repository/runtime，且 B1 Hydration-gated Share truth 與 B2 Global Share setting is an absolute gate scenarios 全部通過。
- [ ] 1.3 [P] 實作 Workspace ownership is proven by persisted relationships：沿用 workspace mapping、source-rule revision、tag 關聯查詢建立 ownership proof，禁止只以非空 tag_id 授權刪除並保留未知 orphan mapping；以 cross-workspace、unpersisted-tag、owned-mapping focused tests 驗證 B3 Durable workspace ownership is required。

## 2. register geometry 與 backend validation

- [ ] 2.1 [TDD-Red] 先建立 B4 Register units, datatype spans, and capacity are backend contracts 的 failing tests，覆蓋 40001→zero-based 0、每 register 2 bytes、int32/float32 span=2、int64/float64 span=4、stride 小於 span、capacity overflow、range collision 與 frontend hint 不可放寬 backend；以 internal/datalink/modbusshare/*validation*_test.go 或同等 focused test target 驗證紅燈案例具體失敗。
- [ ] 2.2 實作 Register geometry is canonical and backend-validated：建立可重用 datatype span/stride/capacity/range validator、response 同時回傳 human/zero-based/span/stride/capacity，並把完整 prevalidation 接到 Share API；以 register contract tests 驗證 B4 Register units, datatype spans, and capacity are backend contracts，且 invalid set 在 mutation 前不改 durable/runtime state。

## 3. atomic reconcile、CAS 與 failure recovery

- [ ] 3.1 [TDD-Red] 為 B5 Reconcile is atomic, serialized, idempotent, and revision-checked 建立 failing tests：完整 desired-set collision prevalidation、repository/runtime partial failure、rollback success、rollback uncertainty dirty_unknown、同 revision double submit 與 stale revision CAS；以 deterministic old/new projection assertions 驗證不能先刪後半套 upsert。
- [ ] 3.2 實作 Reconcile is prevalidated, atomic, serialized, idempotent, and revision-checked：以 workspace lock、完整 desired snapshot、durable transaction 或 staged-copy swap、revision CAS 與 recovery outcome 串起 backend reconcile；以 3.1 tests、race/concurrency focused run 及 response 的 applied/failed/dirty_unknown/new revision assertions 驗證 B5 Reconcile is atomic, serialized, idempotent, and revision-checked。

## 4. activation barrier、candidate seam 與 restart

- [ ] 4.1 [P] 實作 Activation is an autosave barrier：frontend activation controller await global settings、source-rule fields、mapping/candidate saves 並攜帶同一 workspace/settings revision/readiness token，backend 對 pending/error/stale save 回傳 modbus_share_save_incomplete 或 revision conflict；以 frontend save-error/pending tests 與 internal/api/handlers/studio_v2_workspace_activation_handler*_test.go 驗證 B6 Activation waits for a durable autosave barrier。
- [ ] 4.2 [P] 實作 Candidate/apply/restore remains the only projection seam：將既有 source-rule candidate generation/validation/apply/restore 接到 Share reconcile，拒絕沒有 workspace candidate/revision 的 direct runtime mapping mutation，並讓 restart 從 persisted desired set 恢復；以 source-rule restore、direct-API rejection、process-restart projection equality tests 驗證 B7 Existing candidate/apply/restore seam is authoritative。

## 5. global settings 與 listener lifecycle

- [ ] 5.1 實作 Global settings drive listener lifecycle：將 enabled、bind_address、port、slave_id、capacity_registers 與 settings_revision 持久化並作為 startup/enable/disable/stop/bind/reconcile 的唯一設定來源，移除無條件 5020 fallback，保留 disabled/starting/running/stopping/failed 狀態與 operator-safe typed diagnostics；以 listener lifecycle、bind conflict、disabled-no-listener、settings migration tests 驗證 B8 Global settings control the listener lifecycle 與 Local Modbus server sink on port 5020。

## 6. stale span 與 workbench surface

- [ ] 6.1 實作 Datatype changes invalidate old spans：以 mapping identity 與 span-aware diff 在 delete/move/retype 前清除 staged 舊範圍或標記 invalidated，ownership 不明時回傳 invalidated_unknown 並阻擋 ready；以 int64 delete、float32→int64 retype、move collision、unknown ownership focused tests 驗證 B9 Mapping changes invalidate obsolete datatype spans。
- [ ] 6.2 [P] 將 frontend Local Modbus surface 收斂為 Hydrated workspace-scoped Local Modbus operations、Backend-authoritative Local Modbus projection、Span-aware review and recovery feedback：hydration 前 blocked、每次 request 帶 workspace/revision、candidate snapshot 為 truth、direct runtime bind 不顯示 applied、完整 span/capacity/collision/dirty diagnostics 可見；以 frontend/tests/unit/pages/datalink/workbench-local-modbus-review-surface.test.tsx、frontend/tests/unit/workbench-v2/step4-share*.test.tsx 與 local-modbus-memory-workbench contract assertions 驗證。

## 7. frontend barrier 與 operator-safe diagnostics

- [ ] 7.1 [P] 補齊 B1 Hydration-gated Share truth、B2 Global Share setting is an absolute gate、B6 Activation waits for a durable autosave barrier 的 /studio/v2 UI 行為：summary 在 disabled 時隱藏、defaults 不冒充 ready、save error/pending/revision conflict 顯示可行動且不含 secrets 的錯誤；以 frontend Vitest focused run、i18n key coverage 與 DOM assertions 驗證。

## 8. end-to-end acceptance 與 closeout verification

- [ ] 8.1 實作 B10 End-to-end acceptance and field boundary 的可重跑 acceptance harness：HTTP configure→durable save→activate→real Modbus client read→EXE restart→restore/read，並執行 disabled/no-listener/no-mutation、capacity/collision、partial failure、ownership、double-submit negative cases；以 backend integration、frontend integration 與 target EXE smoke logs/hash/read evidence 驗證，另把 field PLC/SCADA sign-off 保留為獨立 gate。
- [ ] 8.2 完成 implementation handoff verification：逐一回讀 B1-B10、Local Modbus operation requirements 與 design decisions 的 observable behavior，使用父代理指定的 authorized fixed-point 分別檢查 approved application seams（internal/api/、internal/datalink/modbusshare/、internal/datalink/sourcerule/、internal/datalink/workspace/、internal/datalink/runtime/、internal/datalink/schema/migrations/、internal/datalink/mapping/、internal/virtual/memory/、internal/virtual/server/modbus/、frontend/src/）與本 change artifact paths 的實際 diff；分開盤點 git diff、git diff --cached、git status --short 與 git ls-files --others --exclude-standard 的 staged/unstaged/untracked inventory，執行 fixed-point/application/artifact 對應的 git diff --check，再執行 spectra status --change productionize-modbus-share-lifecycle --json、spectra analyze productionize-modbus-share-lifecycle --json、spectra validate productionize-modbus-share-lifecycle 及風險相稱的 focused Go/Frontend tests；不得只檢查 parked change directory，並回報未執行的 live hardware checks，不把 build/test 視為 field acceptance。

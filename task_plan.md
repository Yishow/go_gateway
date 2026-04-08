# UI/UX 改善執行計畫

**開始時間**：2026-03-22  
**目標**：依 findings.md 改善專案 UI/UX 缺失  
**策略**：優先處理 P0，再逐步處理 P1、P2

---

## 臨時任務：2026-03-23 規範文件對齊 ✅

- [x] 深度盤點 `README.md`、`Makefile`、`frontend/package.json`、`scripts/build.ps1`、`frontend/tsconfig.json`、`frontend/eslint.config.js`、`.golangci.yml`
- [x] 對齊 `AGENTS.md` 的專案總覽、建置與測試指令、程式碼樣式、測試規範、安全考量與 repo 特定規則
- [x] 對齊 `CLAUDE.md` 的專案總覽、常用命令、樣式 / 測試 / 安全摘要與工作流程
- [x] 修正兩份文件對規範優先順序的描述，使其一致為：`AGENTS.md -> Agent 專屬文件 -> .github/instructions/`
- [x] 以 `git diff --check -- AGENTS.md CLAUDE.md` 驗證文件 patch 無 whitespace / 格式錯誤

**備註**：
- 本次僅更新文件與 planning records，未碰其他進行中的程式碼變更。
- 由於屬文件對齊任務，未額外執行 build / test；驗證以來源盤點與 diff 檢查為主。

---

## 臨時任務：2026-04-07 Workbench UX Phase 2–5 rollout planning（本輪）

- [x] 完成 `docs/superpowers/specs/2026-04-07-workbench-phase-2-to-5-rollout-design.md`
- [x] 完成 reviewer loop，取得雙 reviewer `✅ Approved`
- [x] 取得使用者 spec 核准
- [x] 依新 spec 修正 `openspec/changes/workbench-ux-operator-efficiency/{proposal,design,tasks}.md`
- [x] 依修正後 OpenSpec 建立 Phase 2 implementation plan（shared -> baseline -> v2 -> v1 -> v3 -> compare）
- [x] 同步更新 SQL todo DAG 與本輪 planning records

**備註**：
- 最新 spec commit：`abc2059`
- 目前 blocker 已從 brainstorming gate 移除，但在 OpenSpec 修正完成前，仍不可進入 Phase 2 實作。

### 本輪執行骨架
1. OpenSpec correction：更新 `proposal.md`、`design.md`、`tasks.md` 為 winner-led rollout model
2. Phase 2 shared contract：補齊 Source shared acceptance、scenario matrix、done conditions
3. Phase 2 baseline：凍結並收集 baseline evidence
4. Phase 2 `v2`：先完成 canonical Source flow
5. Phase 2 `v1`：在相同 flow 上完成高質感 full-flow surface
6. Phase 2 `v3`：只完成 minimum-obligation compare surface
7. Phase 2 compare：輸出 compare rubric、推薦結論、checkpoint update

---

## Phase 1: P0 高優先級改善（影響可用性）

### 1.1 統一載入狀態與錯誤處理 ✅
- [x] 建立統一 Spinner 元件（`frontend/src/components/ui/spinner.tsx`）
- [x] 建立統一 Skeleton 元件（`frontend/src/components/ui/skeleton.tsx`）
- [x] 建立 ErrorBoundary 元件（`frontend/src/components/ErrorBoundary.tsx`）
- [x] 更新 App.tsx 加入 ErrorBoundary
- [x] 更新 WorkbenchDeviceStep 使用 Spinner
- [x] 補充測試

### 1.2 補齊 ARIA 屬性（螢幕閱讀器支援）
- [ ] WorkbenchStepRail 補充 `aria-current`、`aria-label`
- [ ] MemoryGrid 補充 `role="grid"`、`aria-label`
- [ ] TagBindingStudio 批量操作按鈕補充 `aria-describedby`
- [ ] SourceCanvasSection 地址選擇器補充 `aria-selected`
- [ ] 補充焦點管理（步驟切換後自動聚焦）
- [ ] 補充測試

### 1.3 補充 Studio 主流程 E2E 測試
- [ ] 建立 `frontend/tests/e2e/studio/device-to-output.spec.ts`
- [ ] 建立 `frontend/tests/e2e/studio/cross-step-traceability.spec.ts`
- [ ] 執行測試驗證

---

## Phase 2: P1 中優先級改善（影響體驗）

### 2.1 補充跨步驟追蹤視覺反饋
- [ ] 為 focused 項目補充視覺高亮樣式
- [ ] 在 Inspector Panel 顯示追蹤路徑
- [ ] 補充測試

### 2.2 補充移動端響應式設計
- [ ] WorkbenchFrame 補充響應式斷點
- [ ] 小螢幕下 Inspector Panel 改用 Drawer
- [ ] 移動端下 StepRail 改為頂部 Tab
- [ ] 補充移動端 E2E 測試

### 2.3 補充 Inspector Panel 與 BottomSummaryBar 互動設計
- [ ] BottomSummaryBar 就緒指示器補充點擊跳轉
- [ ] Inspector Panel 補充 "返回列表" 按鈕
- [ ] 補充互動設計文件

---

## Phase 3: P2 低優先級改善（改進體驗）

### 3.1 補充鍵盤導航支援
- [ ] 補充鍵盤快捷鍵文件
- [ ] 為主流程補充快捷鍵（Ctrl+1/2/3/4）
- [ ] 為 Grid/Canvas 補充方向鍵導航

### 3.2 補充無障礙性測試
- [ ] 補充 a11y 單元測試（jest-axe）
- [ ] 補充 a11y E2E 測試（@axe-core/playwright）

### 3.3 補充元件測試覆蓋
- [ ] 補齊核心元件測試（DeviceForm、MemoryGrid、MappingCanvas）
- [ ] 補充 Wizard 子元件測試

---

## 當前狀態

**Phase**: 2026-04-08 Workbench UX reopened Phase 3 compare complete
**進度**: `3.1` 全部完成；compare gate 已正式判定 `v2` 繼續維持 Tag canonical owner，`v1` / `v3` 則作為可吸收亮點的 distinct variants 保留 evidence 與 branch checkpoint
**下一步**: 若你核准，進入 `4.1 共用基礎`（Output shared foundation）；若要先 review，本輪 compare 輸出與證據也都已備齊

---

## 臨時任務：2026-04-06 staged code review / bugfix

- [x] 盤點所有未提交變更與 staged diff
- [x] 執行前端 targeted tests、Go targeted tests、frontend build
- [x] 修正 review findings：
  - 新 SQLite WAL sidecar 檔未被 `.gitignore` 忽略
  - staged change 新增的 MQTT topic-based source planner 位址支援與現有 canvas / backend point model 不相容
- [x] 補寫 planning records
- [ ] 建立繁中詳細 commit

**備註**：
- 本輪以 review staged 變更為主，未擴大處理 repo 內既有、與本次 diff 無關的 MQTT 資料模型舊債。
- Source planner 的 MQTT topic-based 位址支援先回收，避免前端接受 topic 但後端與畫布模型無法正確處理。

---

## 臨時任務：2026-04-06 uncommitted review / bugfix（本輪）

- [x] 重新盤點目前 staged 變更，確認實際範圍為 `point`、`sourcerule` 與 `sourceCanvasModel`
- [x] 審查並確認一個前端顯示 regression：
  - `formatSourceValue()` 新增 object/array/JSON 支援後，numeric string 不再套用 `hex` / `binary` / `float` 格式
- [x] 修正 formatter，讓 `point.last_value` / JSON payload 內的 numeric string 仍遵守數值顯示模式
- [x] 補上 regression test case（numeric string 與 JSON 內包 numeric string）
- [x] 執行可在目前環境完成的驗證：
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npx eslint src/pages/datalink/workbench/sourceCanvasModel.ts src/pages/datalink/workbench/__tests__/sourceCanvasModel.test.ts`
  - `GOCACHE/TEMP/TMP` 指向 workspace 後執行 `go test ./internal/datalink/point ./internal/datalink/sourcerule`
- [ ] 建立繁中詳細 commit

**備註**：
- `vitest` 在目前 sandbox 仍會因 `esbuild` 啟動子行程遭 `spawn EPERM` 阻擋，因此本輪以前端型別檢查、ESLint 與新增 regression test 內容作為最小可驗證證據。

---

## 臨時任務：2026-04-06 README 規範整併（本輪）

- [x] 深度盤點 `AGENTS.md`、`CLAUDE.md`、`.github/instructions/*`、`Makefile`、`go.mod`、`frontend/package.json`、`cmd/test_ui/main.go`、`internal/api/router.go`、`internal/web/embed.go`
- [x] 對齊並重寫 `README.md`，補齊：
  - 專案總覽
  - 建置與測試指令
  - 程式碼樣式
  - 測試指引與測試要求
  - 安全考量
  - 命名規則
  - error handling pattern
  - 禁止事項
  - repo 特定規則（規範優先序、OpenSpec、planning-with-files、路由主線）
- [x] 同步更新 `task_plan.md`、`findings.md`、`progress.md`
- [ ] 建立繁中詳細 commit

**備註**：
- 本輪為文件整併與規範對齊，不涉及程式行為變更。
- 路由主線與 legacy redirect 描述以目前 `frontend/src/App.tsx` 現況為準（`/studio` 主線、`/test` 工具、`/datalink/*` compat 收斂）。

---

## 臨時任務：2026-04-06 AGENTS / CLAUDE / README 規範整併（本輪）

- [x] 深度比對 `AGENTS.md`、`CLAUDE.md`、`README.md` 與 repo 現況
- [x] 在 `AGENTS.md` 補齊與強化：
  - 命名規則
  - Error Handling Pattern
  - 測試指引與測試要求
  - 禁止事項
  - build/test 指令補全（gatev11、points down、longtask-smoke）
- [x] 在 `CLAUDE.md` 補齊與強化：
  - 命名規則
  - Error Handling Pattern
  - 測試要求（提交前）
  - 禁止事項
  - 其他 repo 特定規則
  - 常用命令補全（gate/migration/smoke）
- [x] 確認 `README.md` 仍與上述規範維持一致（章節結構與內容對齊）
- [x] 同步更新 `task_plan.md`、`findings.md`、`progress.md`
- [ ] 建立繁中詳細 commit

**備註**：
- 本輪以規範文件一致化為目標，未涉及程式邏輯變更。

---

## 臨時任務：2026-04-06 檔案行數規範強制落地（本輪）

- [x] 新增行數檢查腳本 `scripts/check_file_lines.sh`（300 警告 / 500 阻擋）
- [x] 新增 ignore 清單 `.line-limit-ignore`
- [x] 新增 CI workflow `.github/workflows/file-line-limit.yml`
- [x] 新增本機 hook `.githooks/pre-commit`
- [x] 新增 `make check-lines` 命令
- [x] 更新 `.github/pull_request_template.md`，加入行數檢查與超過 300 行說明欄位
- [x] 更新 `AGENTS.md`、`CLAUDE.md`、`README.md`：
  - 強制聲明 `AGENTS.md` 與 `CLAUDE.md` 必須互相參考，不可只讀一份
  - 補上行數規範與執行方式
- [ ] 建立繁中詳細 commit

**備註**：
- 腳本預設以「本次變更檔案」為檢查範圍，避免一次性阻擋歷史超長檔；歷史超長檔若本次變更未增加行數，允許通過並保留逐步收斂空間。

---

## 臨時任務：2026-04-06 未提交變更 code review / bugfix / commit（本輪）

- [x] review 所有未提交變更（排除使用者指定忽略檔案）
- [x] 修正一個行數檢查腳本 bug：
  - 本地模式未涵蓋 untracked 新檔
  - 改為同時檢查 staged / unstaged / untracked
- [x] 重新執行驗證：
  - `bash scripts/check_file_lines.sh`
  - `make check-lines`
  - `git diff --check`
- [x] 建立繁中詳細 commit

---

## 臨時任務：2026-04-06 Workbench UX Phase -1 API 實作（本輪）

- [x] 補齊 SourceRule output apply API（無 mock）：
  - `POST /api/v1/datalink/source-rules/:id/database-outputs/apply`
  - `POST /api/v1/datalink/source-rules/:id/local-modbus/apply`
- [x] 補齊 DB target 工具 API（無 mock）：
  - `POST /api/v1/datalink/db-targets/connectors/:id/schema/generate`
  - `POST /api/v1/datalink/db-targets/connectors/:id/mappings/dry-run`
  - `GET /api/v1/datalink/db-targets/connectors/:id/write-history`
- [x] 補 service + handler 測試：
  - revision mismatch
  - per-item partial success
  - schema_missing / connector_unavailable
  - local-modbus deferred apply contract
- [x] 更新 OpenSpec tasks：
  - `-1.2`、`-1.3`、`-1.4` 標記完成
- [x] 執行本輪驗證（targeted）：
  - `go test ./internal/datalink/sourcerule -run 'TestService_Apply(DatabaseOutputCandidates|LocalModbusOutputCandidates)'`
  - `go test ./internal/datalink/dbtarget -run 'Test(MappingService_DryRun|ConnectorService_GenerateSchema|ConnectorService_ListWriteHistory)'`
  - `go test ./internal/api/handlers -run 'TestSourceRuleHandler_Apply(DatabaseOutputs|LocalModbusOutputs)|TestDatabaseTargetHandler_(DryRunMappings|GenerateSchema|ListWriteHistory)'`
  - `bash scripts/check_file_lines.sh`
  - `git diff --check`

**備註**：
- `internal/api/handlers` 全量測試在目前 sandbox 仍會因 `listen tcp :0` 權限限制失敗；本輪改用 targeted tests 驗證新增端點與契約。

---

## 臨時任務：2026-04-07 Workbench UX Phase 0 啟動（本輪）

- [x] 重新讀取 `proposal.md`、`design.md`、`tasks.md`、`specs/datalink-api/spec.md`
- [x] 建立 fleet SQL todo DAG：
  - `phaseN-shared -> phaseN-baseline -> phaseN-(v1|v2|v3) -> phaseN-compare`
- [x] 啟用 `frontend-design` 並補讀 UX / color design context
- [x] 建立正式 worktree：
  - `.worktrees/woe-base-current-ui`
  - `.worktrees/woe-v1-radix`
  - `.worktrees/woe-v2-mui`
  - `.worktrees/woe-v3-antd`
- [x] 落地 Phase 0 shared foundation 檔案：
  - `frontend/src/styles/workbench-experiment-tokens.ts`
  - `frontend/src/pages/datalink/workbench/workbenchExperimentContract.ts`
- [x] 對齊 OpenSpec worktree 路徑到 repo 既有 `.worktrees/` 規則
- [x] 將 shared foundation commit 套用到 baseline / v1 / v2 / v3 worktree
- [x] 在 baseline branch 建立 freeze checkpoint 與證據
- [x] 平行派發 v1 / v2 / v3 subagents，只做 Phase 0
- [x] 完成 Phase 0 compare gate，確認四套表面共同起跑線

**備註**：
- 本輪保留既有 `.worktrees/workbench-v1|v2|v3` 舊實驗工作樹，不納入本次正式 matrix。
- `/studio` route contract 目前仍由 `frontend/src/App.tsx` 明確守住，Phase 0 shared foundation 只補 contract 與 token source，不直接改 baseline UI 表面。
- Phase 0 compare 結論：
  - baseline / v1 / v2 / v3 皆可在 `4173/4174/4175/4176` 開啟同一個 `/studio`
  - 四版透過 frontend proxy 命中同一組真實 API（`/api/v1/datalink/devices` response hash 一致）
  - 推薦版本（僅 Phase 0）：`v2 / MUI`，原因是 step shell 語意最清楚、版本辨識最穩定；但維護成本目前也是四版中最高
  - Phase 1 以前暫不推進，等待使用者確認下一輪

---

## 臨時任務：2026-04-07 Workbench UX Phase 1 Device 啟動（本輪）

- [x] 將 `phase1-shared` 轉為結構化 contract：
  - `WORKBENCH_DEVICE_COMPARE_SCENARIOS`
  - `WORKBENCH_DEVICE_COMPARE_ACCEPTANCE`
  - `WORKBENCH_DEVICE_COMPARE_CRITICAL_TASK`
- [x] 明確鎖定 Device compare 的六個動作：
  - create
  - edit
  - clone
  - connect
  - probe
  - diagnostics
- [ ] 以 baseline `/studio` 實際操作 Device step，建立 baseline 證據
- [x] 以 baseline `/studio` 實際操作 Device step，建立 baseline 證據
- [x] 完成 `phase1-baseline`
- [ ] 平行派發 `phase1-v1` / `phase1-v2` / `phase1-v3`
- [x] 平行派發 `phase1-v1` / `phase1-v2` / `phase1-v3`
- [x] 完成 `phase1-compare`

**備註**：
- baseline 首次 agent-browser 探查已確認：create / edit / clone / test-connection 都存在於 Device surface。
- baseline 已完成 Device 基線量測，證據包含：
  - `phase1-baseline-device-create.png`
  - `phase1-baseline-device-detail.png`
  - `phase1-baseline-device-edit.png`
  - `phase1-baseline-device-test-current-draft.png`
- baseline 的 connect / probe 語意分層目前是：
  - 已選設備 detail 露出 `Test connection`
  - 編輯器露出 `Test current draft settings`
  - diagnostics 區內才會出現 `Connect phase` 與 `Protocol probe`
- 這代表 baseline 並非缺少 probe 能力，而是 **probe 被包在診斷結果裡，不是首屏可辨識的一級操作**；這將作為 Phase 1 compare 的主觀察項。
- v2 / MUI 已完成 `1.1.v2`：
  - commit：`16d27b3`
  - 1741 行 monolith 被拆成 6 個 MUI 元件，單檔最高 273 行
  - connect / probe diagnostics 已提升為一級 panel 與 stage 狀態
  - 證據：
    - `phase1-v2-device-list.png`
    - `phase1-v2-device-diagnostics.png`
    - `phase1-v2-device-create-editor.png`
  - 已知風險：
    - `/api/v1/datalink/mappings` 500 為 pre-existing backend issue
    - 未使用 `@mui/lab`，改以 `Button + CircularProgress` 實作 loading 狀態
- v1 / shadcn-Radix 已完成 `1.1.v1`：
  - commits：`b1a6b9c`、`4896fc8`
  - `WorkbenchDeviceStep` 拆成 8 個角色檔案；主檔 462 行，其餘多數 <250 行
  - connect / probe diagnostics 已提升為明確分層區塊
  - 證據：
    - `phase1-v1-device-list.png`
    - `phase1-v1-device-diagnostics.png`
  - timing：
    - list -> detail：`459ms`
    - detail -> editor：`397ms`
    - editor -> diagnostics：`381ms`
  - 已知風險：
    - 主協調檔 `WorkbenchDeviceStep.tsx` 仍有 462 行，雖未超過 500 行硬上限，但維護風險高於 v2
- v3 / Ant Design 已完成 `1.1.v3`：
  - commits：`dc496b9`、`6ba09f8`
  - `WorkbenchDeviceStep` 主檔 258 行，但 `DeviceFormDrawer.tsx` 有 422 行
  - connect / probe diagnostics 已提升為雙階段區塊
  - 證據：
    - `phase1-v3-device-list.png`
    - `phase1-v3-device-diagnostics.png`
  - timing：
    - list -> detail：`486ms`
    - detail -> editor：`467ms`
    - editor -> diagnostics：`462ms`
  - 已知風險：
    - foundation tests 28 failures 為 Phase 0 Ant Design shell 既有問題，非本次 Device step regression
    - `DeviceFormDrawer.tsx` 已超過 300 行軟上限，維護風險偏高
- Phase 1 compare 結論：
  - **推薦版本：v2 / MUI**
  - 理由：
    - connect / probe diagnostics 最清楚，兩階段狀態在 UI 上最穩定、最容易被 operator 立即辨識
    - create / edit mode 語意最直接，與 Device detail / diagnostics / editor 的責任分層最完整
    - 6 個主要檔案全部維持在 300 行以下，維護面優於 v1 / v3
  - compare 摘要：
    - `baseline`：最快、風險最低，但 probe 語意藏在 diagnostics 裡，不適合作為新 Device 表面
    - `v1`：速度最接近 baseline，且把 probe 語意前移；是低風險備選
    - `v2`：邏輯清晰度與系統完整性最佳，雖互動略慢於 baseline / v1，但最平衡
  - `v3`：資訊密度與工作台感最強，但互動較慢，且維護 / foundation 風險最高

---

## 臨時任務：2026-04-07 Workbench UX Phase 1R Device 重做（使用者直接打回）

- [x] 明確重新定義 round 方向：
  - 三版都要是不同 interaction model + 不同視覺語言
  - 三版都允許徹底重做 Device interaction skeleton
  - 三版都要明確對齊 `design-md`
- [x] 與使用者確認三 archetypes 策略：
  - `v1` = Linear Control Room
  - `v2` = Sentry Incident Desk
  - `v3` = ClickHouse Data Cockpit
- [x] 將 Phase 1R 三 archetypes 與共通 domain model 回寫 OpenSpec
- [x] 平行派發 `phase1r-v1` / `phase1r-v2` / `phase1r-v3`
- [x] 發現 shared token contract 過窄，依 blocked 規則重開 shared foundation
- [x] 擴充 `frontend/src/styles/workbench-experiment-tokens.ts`：
  - archetype-specific palette semantics（linear / sentry / clickhouse）
  - typography cues（510 / 700 / 900 權重與 tracking）
  - depth / treatment semantics（glass / inset / neon glow）
- [x] 補 shared token contract unit test：
  - `frontend/tests/unit/features/datalink/workbench-experiment-tokens.test.ts`
- [ ] 讓 `phase1r-v1` / `phase1r-v2` / `phase1r-v3` 依新 shared contract 對齊後續作
- [ ] 完成 Phase 1R compare gate

**備註**：
- 直接觸發點是 v2 開始長 branch-local `sentryVisualTokens.ts`；雖然該檔多半是 style preset 而非完整新 palette，但它暴露出 shared token contract 仍不足以支撐三 archetypes 的真正視覺語言。
- 本輪 shared repair 已把「不可在 branch 內新增 version-only token file」回寫到 OpenSpec design / tasks，避免後續再次 drift。
- 使用者明確否定上一輪 Device compare，理由是：
  - 三版 UI 差異不足
  - interaction model 幾乎沒變
  - 視覺語言沒有真正對齊 `design-md`
  - 看起來仍太像 baseline

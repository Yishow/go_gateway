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

**Phase**: 1.2 補齊 ARIA 屬性  
**進度**: 0/6 完成  
**下一步**: WorkbenchStepRail 補充 aria-current、aria-label

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

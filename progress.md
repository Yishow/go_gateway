# UI/UX 改善執行進度

**開始時間**：2026-03-22  
**當前 Phase**：2026-04-08 Workbench UX Phase 2 source compare ready

---

## 執行記錄

### 2026-04-08

#### Session 9: Workbench UX Phase 2–5 rollout design review / planning handoff（本輪）
- ✅ 完成 `docs/superpowers/specs/2026-04-07-workbench-phase-2-to-5-rollout-design.md` 修補：
  - baseline / baseline check 定義
  - `v3` 最低交付義務
  - shared scenario matrix
  - OpenSpec amendment outputs 與 planning handoff criteria
- ✅ reviewer loop 通過：
  - `spec-reviewer-rollout-fast`：`✅ Approved`
  - `spec-reviewer-rollout`：`✅ Approved`
- ✅ 使用者核准 spec，可進入 implementation planning
- ✅ 完成 OpenSpec correction：
  - 更新 `openspec/changes/workbench-ux-operator-efficiency/proposal.md`
  - 更新 `openspec/changes/workbench-ux-operator-efficiency/design.md`
  - 更新 `openspec/changes/workbench-ux-operator-efficiency/tasks.md`
  - `openspec validate --changes workbench-ux-operator-efficiency` 通過
- ✅ 完成 Phase 2 shared contract、2.1.baseline 與 `2.1.v1` / `2.1.v2` / `2.1.v3` Source rollout
- ✅ compare evidence：baseline 與三版均已具備 overview / focused / handoff；v1 / v2 / v3 已完成 browser verification 與 test/build
- ⏳ 下一步：開始 2.1.compare Source 比較與推薦

**驗證結果**：
```bash
git --no-pager diff --check -- docs/superpowers/specs/2026-04-07-workbench-phase-2-to-5-rollout-design.md
```
- 結果：通過

**修改檔案**：
- docs/superpowers/specs/2026-04-07-workbench-phase-2-to-5-rollout-design.md
- task_plan.md
- findings.md
- progress.md

### 2026-03-23

#### Session 2: AGENTS / CLAUDE 規範文件對齊 ✅
- ✅ 深度盤點 `README.md`、`Makefile`、`frontend/package.json`、`scripts/build.ps1`、`frontend/tsconfig.json`、`frontend/eslint.config.js`、`.golangci.yml`
- ✅ 更新 `AGENTS.md`，補齊專案總覽、完整建置 / 測試命令、程式碼樣式、測試規範、安全考量與 repo 特定規則
- ✅ 更新 `CLAUDE.md`，對齊 `AGENTS.md` 並補齊命令、樣式 / 測試 / 安全摘要與 agent 工作流程
- ✅ 收斂兩份文件的規範優先順序描述
- ✅ 補寫 planning records（`task_plan.md`、`findings.md`、`progress.md`）

**驗證 / 修改摘要**：`git diff --check -- AGENTS.md CLAUDE.md` 通過；修改 `AGENTS.md`、`CLAUDE.md` 與三份 planning records。

### 2026-03-22

#### Session 1: Phase 1.1 統一載入狀態與錯誤處理 ✅
- ✅ 建立統一 Spinner 元件（`frontend/src/components/ui/spinner.tsx`）
- ✅ 建立統一 Skeleton 元件（`frontend/src/components/ui/skeleton.tsx`）
- ✅ 建立 ErrorBoundary 元件（`frontend/src/components/ErrorBoundary.tsx`）
- ✅ 更新 App.tsx 加入 ErrorBoundary
- ✅ 更新 WorkbenchDeviceStep 使用 Spinner
- ✅ 補充測試（Spinner、Skeleton、ErrorBoundary）
- ✅ 測試通過（21 tests passed）

**測試 / 修改摘要**：21 個前端單元測試通過；新增 `spinner`、`skeleton`、`ErrorBoundary` 與對應測試，並更新 `App.tsx`、`WorkbenchDeviceStep.tsx`。

#### Session 1: Phase 1.2 補齊 ARIA 屬性（進行中）
- ✅ WorkbenchStepRail 補充 `aria-label`（包含就緒狀態）
- ⏳ MemoryGrid 補充 `role="grid"`、`aria-label`（待續）
- ⏳ TagBindingStudio 批量操作按鈕補充 `aria-describedby`（待續）
- ⏳ SourceCanvasSection 地址選擇器補充 `aria-selected`（待續）
- ⏳ 補充焦點管理（待續）
- ⏳ 補充測試（待續）

**下一步**：繼續 Phase 1.2 - MemoryGrid ARIA 屬性

### 2026-04-06

#### Session 3: staged code review / bugfix
- ? 盤點 staged 變更 10 檔，確認主要集中在 Source planner / address parser / SQLite DSN
- ? 執行 review 驗證：
  - `cd frontend && npm run test -- tests/unit/utils/addressParser.test.ts tests/unit/features/datalink/sourcePlannerContract.test.ts tests/unit/pages/datalink/workbench-source-canvas-model.test.ts tests/unit/pages/datalink/workbench-source-step.test.tsx`
  - `cd frontend && npm run build`
  - `go test ./internal/datalink ./cmd/test_ui`
- ? 修正兩個 review findings：
  - `.gitignore` 補上 `datalink.db-wal` / `datalink.db-shm`
  - 回收 staged diff 中不完整的 MQTT topic-based source planner 位址支援，避免前端接受 topic 但現有 canvas / backend point model 無法正確處理
- ? 同步更新 `frontend/tests/unit/utils/addressParser.test.ts`，改回驗證目前正式支援的 planner 預設位址 contract

**驗證 / 修改摘要**：前端 targeted tests 83 個通過、build 成功、Go target packages 通過；修正 `.gitignore` 與 address parser contract，並同步 planning records。

#### Session 4: uncommitted review / bugfix（本輪）
- ? 重新確認目前 staged 變更實際集中於：
  - `internal/datalink/point/*`
  - `internal/datalink/sourcerule/*`
  - `frontend/src/pages/datalink/workbench/sourceCanvasModel.ts`
- ? 發現並修正前端 formatter regression：
  - `formatSourceValue()` 對 numeric string 未套用 `hex` / `binary` / `float` 格式
  - 補強後會先將可解析字串轉成 number，再沿用既有格式化邏輯
- ? 補上 regression tests：
  - numeric string 直接輸入
  - JSON payload 內包 numeric string
- ? 執行驗證：
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npx eslint src/pages/datalink/workbench/sourceCanvasModel.ts src/pages/datalink/workbench/__tests__/sourceCanvasModel.test.ts`
  - `$env:GOCACHE='C:\\AIProject\\go_gateway\\.gocache'; $env:TEMP='C:\\AIProject\\go_gateway\\.gocache\\tmp'; $env:TMP='C:\\AIProject\\go_gateway\\.gocache\\tmp'; go test ./internal/datalink/point ./internal/datalink/sourcerule`

**驗證 / 修改摘要**：`tsc --noEmit`、targeted ESLint、Go targeted tests 通過；修正 `sourceCanvasModel` formatter regression 並補回歸測試。Vitest 仍受 sandbox `spawn EPERM` 限制。

#### Session 5: README 規範整併（本輪）
- ✅ 深度盤點規範與實作來源：
  - `AGENTS.md`
  - `CLAUDE.md`
  - `.github/instructions/go.instructions.md`
  - `.github/instructions/reactjs.instructions.md`
  - `.github/instructions/typescript-5-es2022.instructions.md`
  - `Makefile`
  - `go.mod`
  - `frontend/package.json`
  - `scripts/build.ps1`
  - `cmd/test_ui/main.go`
  - `internal/api/router.go`
  - `internal/web/embed.go`
  - `frontend/src/App.tsx`
- ✅ 重寫 `README.md`，補齊：
  - 專案總覽與模組結構
  - 完整建置/測試命令（含 gate/migration）
  - 程式碼樣式與命名規則
  - error handling pattern
  - 測試指引與測試要求
  - 安全考量與禁止事項
  - repo 特定規則（規範優先序、OpenSpec、planning-with-files）
- ✅ 同步更新 planning records：`task_plan.md`、`findings.md`、`progress.md`

**驗證 / 修改摘要**：`git diff --check -- README.md task_plan.md findings.md progress.md` 通過；重寫 `README.md` 並同步三份 planning records。

#### Session 6: AGENTS / CLAUDE / README 規範整併（本輪）
- ✅ 重新盤點 `AGENTS.md`、`CLAUDE.md`、`README.md` 與實際 repo 流程，確認缺口集中在：
  - 命名規則（需獨立章節）
  - Error Handling Pattern（需獨立章節）
  - 禁止事項（需可執行條列）
  - 測試要求（提交前 gate）
- ✅ 更新 `AGENTS.md`：
  - 補上 `命名規則`、`Error Handling Pattern`、`測試指引與測試要求`、`禁止事項`
  - 補全命令清單（`gatev11`、`points-precheck-down`、`points-migrate-down`、`longtask-smoke`）
- ✅ 更新 `CLAUDE.md`：
  - 補上 `命名規則`、`Error Handling Pattern`、`測試要求（提交前）`、`禁止事項`、`其他 Repo 特定規則`
  - 補全命令清單（gate/migration/smoke）
- ✅ 微調 `README.md`，補齊 `make longtask-smoke`，使命令清單與 `AGENTS.md` / `CLAUDE.md` 完全一致
- ✅ 同步更新 planning records：`task_plan.md`、`findings.md`、`progress.md`

**驗證 / 修改摘要**：`git diff --check -- AGENTS.md CLAUDE.md README.md task_plan.md findings.md progress.md` 通過；更新三份規範文件並同步 planning records。

#### Session 7: 檔案行數規範強制落地（本輪）
- ✅ 新增 `scripts/check_file_lines.sh`：
  - 預設門檻：`>300` 警告、`>500` 阻擋
  - 預設檢查範圍：本次變更檔案（支援 CI base SHA / staged / working tree）
  - legacy guard：歷史上已 >500 的檔案僅允許不增加行數
- ✅ 新增 `.line-limit-ignore`（lock/build/imported docs 等噪音排除）
- ✅ 新增 CI workflow：`.github/workflows/file-line-limit.yml`
- ✅ 新增本機 hook：`.githooks/pre-commit`
- ✅ 新增 `make check-lines` 入口
- ✅ 更新 `.github/pull_request_template.md`，要求：
  - 行數檢查通過
  - 超過 300 行需補 rationale + split plan
- ✅ 更新 `AGENTS.md` / `CLAUDE.md` / `README.md`：
  - 強制聲明 `AGENTS.md` 與 `CLAUDE.md` 必須互相參考，不可只讀一份
  - 補上行數規範與執行方式

**驗證 / 修改摘要**：`check_file_lines.sh` 與 `diff --check` 通過；新增腳本、CI、hook、ignore 與 PR template，並同步規範與 planning records。

#### Session 8: 未提交變更 code review / bugfix / commit（本輪）
- ✅ 完成未提交變更 code review（以目前工作樹為範圍）
- ✅ 發現並修正 bug：
  - `scripts/check_file_lines.sh` 本地 fallback 模式漏檢 untracked 新檔
  - 已改為 staged / unstaged / untracked 三路合併檢查
- ✅ 驗證：
  - `bash scripts/check_file_lines.sh`
  - `make check-lines`
  - `git --no-pager diff --check -- .github/pull_request_template.md .github/workflows/file-line-limit.yml .line-limit-ignore .githooks/pre-commit Makefile AGENTS.md CLAUDE.md README.md scripts/check_file_lines.sh task_plan.md findings.md progress.md`
- ✅ 建立繁中詳細 commit

#### Session 9: Workbench UX Phase -1 API 實作（本輪）
- ✅ 完成 SourceRule output apply API（無 mock）：
  - `POST /datalink/source-rules/:id/database-outputs/apply`
  - `POST /datalink/source-rules/:id/local-modbus/apply`
- ✅ 完成 DB target tooling API（無 mock）：
  - `POST /datalink/db-targets/connectors/:id/schema/generate`
  - `POST /datalink/db-targets/connectors/:id/mappings/dry-run`
  - `GET /datalink/db-targets/connectors/:id/write-history`
- ✅ 完成測試（service + handler）：
  - revision mismatch
  - per-item partial success
  - schema_missing / connector_unavailable
  - local-modbus deferred apply
- ✅ 更新 OpenSpec：`workbench-ux-operator-efficiency/tasks.md` 之 `-1.2`、`-1.3`、`-1.4` 已標記完成

**驗證結果**：
```bash
GOCACHE=$(pwd)/.gocache GOTMPDIR=$(pwd)/.gotmp TMPDIR=$(pwd)/.gotmp \
go test ./internal/datalink/sourcerule -run 'TestService_Apply(DatabaseOutputCandidates|LocalModbusOutputCandidates)'

GOCACHE=$(pwd)/.gocache GOTMPDIR=$(pwd)/.gotmp TMPDIR=$(pwd)/.gotmp \
go test ./internal/datalink/dbtarget -run 'Test(MappingService_DryRun|ConnectorService_GenerateSchema|ConnectorService_ListWriteHistory)'

GOCACHE=$(pwd)/.gocache GOTMPDIR=$(pwd)/.gotmp TMPDIR=$(pwd)/.gotmp \
go test ./internal/api/handlers -run 'TestSourceRuleHandler_Apply(DatabaseOutputs|LocalModbusOutputs)|TestDatabaseTargetHandler_(DryRunMappings|GenerateSchema|ListWriteHistory)'

bash scripts/check_file_lines.sh
git diff --check
```
- 結果：以上命令通過。

**限制**：
- `go test ./internal/api/handlers` 全量在 sandbox 會因既有測試需要 bind TCP 埠（`listen tcp :0`）而失敗；本輪改用 targeted tests 驗證新增 API 契約。

#### Session 10: Phase -1 完成後可平行 worktree 的就緒檢查（本輪）
- ✅ 確認 OpenSpec `tasks.md` Phase -1 全數完成（`-1.1 ~ -1.4`）
- ✅ 確認新 API 路由已掛載於 router：
  - `source-rules/:id/database-outputs/apply`
  - `source-rules/:id/local-modbus/apply`
  - `db-targets/connectors/:id/schema/generate`
  - `db-targets/connectors/:id/mappings/dry-run`
  - `db-targets/connectors/:id/write-history`
- ✅ 再次執行 targeted 驗證（全部通過）：
  - `go test ./internal/datalink/sourcerule -run 'TestService_Apply(DatabaseOutputCandidates|LocalModbusOutputCandidates)'`
  - `go test ./internal/datalink/dbtarget -run 'Test(MappingService_DryRun|ConnectorService_GenerateSchema|ConnectorService_ListWriteHistory)'`
  - `go test ./internal/api/handlers -run 'TestSourceRuleHandler_Apply(DatabaseOutputs|LocalModbusOutputs)|TestDatabaseTargetHandler_(DryRunMappings|GenerateSchema|ListWriteHistory)'`
  - `bash scripts/check_file_lines.sh`
  - `git diff --check`
- ✅ 盤點 worktree 狀態：
  - `.worktrees/workbench-v1|v2|v3` 皆已存在，且 `.worktrees/` 已由 `.gitignore` 忽略
  - 三者目前 HEAD 皆為 `b8688e5`（舊基線），主工作樹為 `main@3a18025` 且含本次未提交變更

**補充驗證**：
- 嘗試執行全量 `go test ./...`：因 sandbox 不允許 bind TCP/UDP（`listen tcp :0`、`listen udp 127.0.0.1:0`）失敗，屬環境限制，非本輪 API 變更引入。

#### Session 11: Workbench UX Phase 0 啟動（本輪）
- ✅ 讀取 OpenSpec context：
  - `proposal.md`
  - `design.md`
  - `tasks.md`
  - `specs/datalink-api/spec.md`
- ✅ 建立 fleet SQL todo graph：
  - `phase0-shared`
  - `phase0-baseline`
  - `phase0-v1|v2|v3`
  - `phase0-compare`
  - 並串起 Phase 1 ~ 6 的 shared / baseline / variants / compare 依賴
- ✅ 啟用 `frontend-design` 並讀取 `ux-psychology.md`、`color-system.md`
- ✅ 建立正式 Phase 0 worktree：
  - `.worktrees/woe-base-current-ui`
  - `.worktrees/woe-v1-radix`
  - `.worktrees/woe-v2-mui`
  - `.worktrees/woe-v3-antd`
- ✅ 為四個 worktree 執行基礎 setup：
  - `go mod download`
  - `cd frontend && npm install --silent`
- ✅ 落地 shared foundation 檔案：
  - `frontend/src/styles/workbench-experiment-tokens.ts`
  - `frontend/src/pages/datalink/workbench/workbenchExperimentContract.ts`
- ✅ 對齊 OpenSpec worktree path：
  - `openspec/changes/workbench-ux-operator-efficiency/design.md`
  - `openspec/changes/workbench-ux-operator-efficiency/tasks.md`
- ✅ 執行本輪 shared-foundation 驗證：
  - `openspec validate --changes "workbench-ux-operator-efficiency"`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npx eslint src/styles/workbench-experiment-tokens.ts src/pages/datalink/workbench/workbenchExperimentContract.ts`
  - `bash scripts/check_file_lines.sh`
  - `git diff --check`

**驗證結果**：
```bash
openspec validate --changes "workbench-ux-operator-efficiency"
cd frontend && npx tsc --noEmit
cd frontend && npx eslint src/styles/workbench-experiment-tokens.ts src/pages/datalink/workbench/workbenchExperimentContract.ts
bash scripts/check_file_lines.sh
git diff --check
```
- 結果：
  - `openspec validate` 通過
  - TypeScript / ESLint 通過
  - line-limit 通過（僅 `design.md`、`tasks.md` 既有長檔 warning）
  - `git diff --check` 通過

**下一步**：
- 將 shared foundation commit 套用到 baseline / v1 / v2 / v3
- baseline branch 建立 freeze checkpoint
- 平行派發 v1 / v2 / v3 subagents 只做 Phase 0

#### Session 12: Workbench UX Phase 0 compare gate（本輪）
- ✅ baseline / v1 / v2 / v3 全部完成 Phase 0 checkpoint：
  - baseline：`e0928d8`
  - v1：`53e4fca`
  - v2：`e3a245f`
  - v3：`e1c3fac`
- ✅ 使用 `agent-browser` 收集四版 `/studio` 證據：
  - `files/phase0-baseline-compare.png`
  - `files/phase0-v1-compare.png`
  - `files/phase0-v2-compare.png`
  - `files/phase0-v3-compare.png`
- ✅ 使用 `agent-browser` 驗證四版 `?step=source` deep-link：
  - baseline `4173`
  - v1 `4174`
  - v2 `4175`
  - v3 `4176`
- ✅ 以 frontend proxy 驗證 shared real API：
  - `GET /api/v1/datalink/devices`
  - 四版 response bytes = `3076`
  - 四版 sha256 = `5cafa581345cfa9225eb221c3c821671ddae03523385cb464fb6669fe6aa9ad4`
- ✅ 量測基本 `/studio` response timing（3 次平均）：
  - baseline = `2.2ms`
  - v1 = `1.8ms`
  - v2 = `2.0ms`
  - v3 = `1.8ms`
- ✅ 完成 Phase 0 compare 結論：
  - 操作順暢度：四版都能穩定載入與切到 Source；差異主要在 shell/step rail 語意
  - 邏輯清晰度：`v2` 最佳，`v3` 次之，`v1` 接近 baseline
  - 對系統的完整性：四版皆維持同 `/studio` + 同 API + 同資料
  - 首屏資訊密度：baseline/v1 最熟悉，v2 最平衡，v3 最有工作台感但略厚重
  - 關鍵操作時間：四版差異極小，均在 `1.8ms ~ 2.2ms`
  - 實作 / 維護風險：baseline < v1 < v3 < v2
  - 推薦版本（Phase 0 only）：`v2 / MUI`

**驗證結果**：
```bash
curl http://127.0.0.1:4173/4174/4175/4176/studio
agent-browser open /studio
agent-browser open /studio?step=source
curl http://127.0.0.1:{4173,4174,4175,4176}/api/v1/datalink/devices
```
- 結果：
  - 四版 `/studio` 可達
  - 四版 deep-link 可達
  - 四版 API proxy hash 一致

#### Session 13: Workbench UX Phase 1 shared criteria（本輪）
- ✅ 在 `workbenchExperimentContract.ts` 補齊 Device phase shared contract：
  - `WORKBENCH_DEVICE_COMPARE_SCENARIOS`
  - `WORKBENCH_DEVICE_COMPARE_ACCEPTANCE`
  - `WORKBENCH_DEVICE_COMPARE_CRITICAL_TASK`
- ✅ 鎖定 Device compare 六個固定動作：
  - create
  - edit
  - clone
  - connect
  - probe
  - diagnostics
- ✅ 使用 `agent-browser` 先探查 baseline Device step：
  - `/studio` Device list 首屏
  - 開啟 `Create device` editor
  - 開啟既有 `Browser Smoke PLC` detail
- ✅ baseline 初步觀察：
  - create / edit / clone 路徑存在
  - 已選設備 detail 有 `Test connection`
  - `probe` 不在 detail 首屏按鈕，而是在 diagnostics 區結果內揭露

#### Session 14: Workbench UX Phase 1 baseline（Current Studio）
- ✅ 使用 `agent-browser` 完成 baseline Device 實際操作：
  - 開啟 `Create device`
  - 開啟既有 `Browser Smoke PLC` detail
  - 進入 `Edit device`
  - 觸發 `Test current draft settings`
- ✅ 保留證據：
  - `phase1-baseline-device-create.png`
  - `phase1-baseline-device-detail.png`
  - `phase1-baseline-device-edit.png`
  - `phase1-baseline-device-test-current-draft.png`
- ✅ baseline 診斷文字證據重點：
  - `Test current draft settings`
  - `Failed`
  - `Connect phase`
  - `Protocol probe`
  - `Skipped`
- ✅ baseline 結論：
  - create / edit / clone 主路徑是完整的
  - connect / probe 其實都有，但拆層不夠前置
  - operator 需要先進入 editor 並執行 diagnostics，才看得到 `Protocol probe`
  - 這使 baseline 在「系統完整性」上仍成立，但在「操作順暢度」與「邏輯清晰度」上留有優化空間

**下一步**：
- 完成 `phase1-baseline`
- 平行派發 `phase1-v1` / `phase1-v2` / `phase1-v3`

#### Session 15: Workbench UX Phase 1 v2 / MUI 完成
- ✅ v2 子代理完成 `1.1.v2`
- ✅ branch / commit：
  - `woe-v2-mui`
  - `16d27b3`
- ✅ 結構變更：
  - `WorkbenchDeviceStep` 拆成 6 個 MUI 元件
  - 單檔最高 273 行，符合 line-limit 目標
- ✅ 驗證：
  - `tsc --noEmit`
  - ESLint
  - Vitest form-model tests `8/8`
  - `agent-browser` 實際操作驗證
- ✅ 證據：
  - `phase1-v2-device-list.png`
  - `phase1-v2-device-diagnostics.png`
  - `phase1-v2-device-create-editor.png`
- ✅ 重要結論：
  - connect / probe diagnostics 已提升為一級顯示
  - create / edit mode 的辨識度提升
  - `/mappings` 500 為 pre-existing backend issue，非 v2 Device surface regression

#### Session 16: Workbench UX Phase 1 v1 / shadcn-Radix 完成
- ✅ v1 branch / commit：
  - `woe-v1-radix`
  - `b1a6b9c`
  - `4896fc8`（補齊 diagnostics i18n 文案）
- ✅ 結構變更：
  - `WorkbenchDeviceStep` 拆成 8 個角色檔案
  - `WorkbenchDeviceDiagnosticsPanel` 明確分層 `Connect` / `Protocol Probe`
- ✅ 驗證：
  - `workbench-device-form-model` + `workbench-device-step-editor` 共 11 個測試通過
  - 新增檔案 lint 通過
  - `agent-browser` 實際操作驗證
- ✅ 證據：
  - `phase1-v1-device-list.png`
  - `phase1-v1-device-diagnostics.png`
- ✅ timing：
  - list -> detail：`459ms`
  - detail -> editor：`397ms`
  - editor -> diagnostics：`381ms`
- ✅ 重要結論：
  - v1 在維持 baseline 節奏的前提下，把 connect / probe 語意前移
  - 主協調檔 462 行，維護風險低於 monolith，但仍高於 v2 的拆分品質

#### Session 17: Workbench UX Phase 1 v3 / Ant Design 完成
- ✅ v3 branch / commit：
  - `woe-v3-antd`
  - `dc496b9`
- ✅ 結構變更：
  - `WorkbenchDeviceStep` 主檔降到 258 行
  - 新增 `DeviceFormDrawer`、`DeviceListTable`、`DeviceDetailCard`、`DeviceDiagnosticsPanel`、`deviceConnectionFieldDefs`
- ✅ 驗證：
  - `tsc --noEmit`
  - `npm run lint`
  - `npm run build`
  - editor tests `3/3`
  - `agent-browser` 實際操作驗證
- ✅ 證據：
  - `phase1-v3-device-list.png`
  - `phase1-v3-device-diagnostics.png`
- ✅ timing：
  - list -> detail：`486ms`
  - detail -> editor：`467ms`
  - editor -> diagnostics：`462ms`
- ✅ 重要結論：
  - v3 具備明確的雙階段 diagnostics 與高資訊密度
  - foundation tests 28 failures 為 Phase 0 shell 既有問題，不是 Device refactor regression
  - `DeviceFormDrawer.tsx` 422 行，形成明確維護風險

#### Session 18: Workbench UX Phase 1 compare gate
- ✅ compare 範圍：
  - baseline (`4173`)
  - v1 / shadcn-Radix (`4174`)
  - v2 / MUI (`4175`)
  - v3 / Ant Design (`4176`)
- ✅ compare 依據：
  - `agent-browser` 實際操作證據
  - 真實 API / 同一路由 `/studio`
  - key operation timings
  - code split / line count / test 狀態
- ✅ 關鍵操作時間（ms）：
  - baseline：`458 / 381 / 377`
  - v1：`459 / 397 / 381`
  - v2：`502 / 493 / 386`
  - v3：`486 / 467 / 462`
  - 順序：`list -> detail / detail -> editor / editor -> diagnostics`
- ✅ 比較結論：
  - **操作順暢度**：`v1` ≈ `baseline` > `v3` > `v2`
  - **邏輯清晰度**：`v2` > `v3` > `v1` > `baseline`
  - **系統完整性**：`v2` > `v1` ≈ `v3` > `baseline`
  - **首屏資訊密度**：`v3` > `v2` > `v1` > `baseline`
  - **實作 / 維護風險**：`baseline` < `v2` < `v1` < `v3`
- ✅ 推薦版本：
  - **v2 / MUI**
  - 理由：connect / probe diagnostics 最清楚、editor mode 語意最直接、主要檔案都壓在 300 行內，是 Device phase 下最平衡的版本

#### Session 19: 使用者直接打回 Phase 1 Device round
- ❌ 使用者明確否定上一輪 Device compare 結果
- ❌ 原因：
  - 三版看起來差異不足
  - 互動模型沒有真正重做
  - 視覺語言沒有對齊 `/Users/yishow/prj/awesome-design-md/design-md`
  - 整體仍太像 baseline
- ✅ 新方向已確認：
  - 採 `Three Archetypes` 重開 Device round
  - `v1` = Linear Control Room
  - `v2` = Sentry Incident Desk
  - `v3` = ClickHouse Data Cockpit
  - 三版都允許完全重做 Device interaction skeleton，只保留同 route / 同 API / 同 token / 同 domain contract

# UI/UX 改善執行進度

**開始時間**：2026-03-22  
**當前 Phase**：1.2 補齊 ARIA 屬性

---

## 執行記錄

### 2026-03-23

#### Session 2: AGENTS / CLAUDE 規範文件對齊 ✅
- ✅ 深度盤點 `README.md`、`Makefile`、`frontend/package.json`、`scripts/build.ps1`、`frontend/tsconfig.json`、`frontend/eslint.config.js`、`.golangci.yml`
- ✅ 更新 `AGENTS.md`，補齊專案總覽、完整建置 / 測試命令、程式碼樣式、測試規範、安全考量與 repo 特定規則
- ✅ 更新 `CLAUDE.md`，對齊 `AGENTS.md` 並補齊命令、樣式 / 測試 / 安全摘要與 agent 工作流程
- ✅ 收斂兩份文件的規範優先順序描述
- ✅ 補寫 planning records（`task_plan.md`、`findings.md`、`progress.md`）

**驗證結果**：
```bash
git --no-pager diff --check -- AGENTS.md CLAUDE.md
```
- 結果：通過，無 whitespace / patch 格式錯誤

**修改檔案**：
- AGENTS.md
- CLAUDE.md
- task_plan.md
- findings.md
- progress.md

### 2026-03-22

#### Session 1: Phase 1.1 統一載入狀態與錯誤處理 ✅
- ✅ 建立統一 Spinner 元件（`frontend/src/components/ui/spinner.tsx`）
- ✅ 建立統一 Skeleton 元件（`frontend/src/components/ui/skeleton.tsx`）
- ✅ 建立 ErrorBoundary 元件（`frontend/src/components/ErrorBoundary.tsx`）
- ✅ 更新 App.tsx 加入 ErrorBoundary
- ✅ 更新 WorkbenchDeviceStep 使用 Spinner
- ✅ 補充測試（Spinner、Skeleton、ErrorBoundary）
- ✅ 測試通過（21 tests passed）

**測試結果**：
```
✓ tests/unit/components/ErrorBoundary.test.tsx (4 tests)
✓ tests/unit/components/ui/spinner.test.tsx (6 tests)
✓ tests/unit/components/ui/skeleton.test.tsx (11 tests)
```

**修改檔案**：
- frontend/src/components/ui/spinner.tsx（新增）
- frontend/src/components/ui/skeleton.tsx（新增）
- frontend/src/components/ErrorBoundary.tsx（新增）
- frontend/src/App.tsx（加入 ErrorBoundary）
- frontend/src/pages/datalink/workbench/WorkbenchDeviceStep.tsx（使用 Spinner）
- frontend/tests/unit/components/ErrorBoundary.test.tsx（新增）
- frontend/tests/unit/components/ui/spinner.test.tsx（新增）
- frontend/tests/unit/components/ui/skeleton.test.tsx（新增）

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

**驗證結果**：
```bash
cd frontend && npm run test -- tests/unit/utils/addressParser.test.ts tests/unit/features/datalink/sourcePlannerContract.test.ts tests/unit/pages/datalink/workbench-source-canvas-model.test.ts tests/unit/pages/datalink/workbench-source-step.test.tsx
cd frontend && npm run build
go test ./internal/datalink ./cmd/test_ui
```
- 結果：通過（frontend 83 tests passed；build 成功；Go target packages 通過）

**修改檔案**：
- .gitignore
- frontend/src/utils/addressParser.ts
- frontend/tests/unit/utils/addressParser.test.ts
- task_plan.md
- findings.md
- progress.md

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

**驗證結果**：
```bash
cd frontend && npx tsc --noEmit
cd frontend && npx eslint src/pages/datalink/workbench/sourceCanvasModel.ts src/pages/datalink/workbench/__tests__/sourceCanvasModel.test.ts
$env:GOCACHE='C:\\AIProject\\go_gateway\\.gocache'; $env:TEMP='C:\\AIProject\\go_gateway\\.gocache\\tmp'; $env:TMP='C:\\AIProject\\go_gateway\\.gocache\\tmp'; go test ./internal/datalink/point ./internal/datalink/sourcerule
```
- 結果：通過

**額外限制**：
- `cd frontend && npm run test -- src/pages/datalink/workbench/__tests__/sourceCanvasModel.test.ts` 在目前 sandbox 會因 `esbuild` 啟動子行程失敗而回報 `spawn EPERM`，因此本輪未能取得 Vitest 執行證據。

**修改檔案**：
- frontend/src/pages/datalink/workbench/sourceCanvasModel.ts
- frontend/src/pages/datalink/workbench/__tests__/sourceCanvasModel.test.ts
- task_plan.md
- findings.md
- progress.md

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

**驗證結果**：
```bash
git --no-pager diff --check -- README.md task_plan.md findings.md progress.md
```
- 結果：通過，無 whitespace / patch 格式錯誤

**修改檔案**：
- README.md
- task_plan.md
- findings.md
- progress.md

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

**驗證結果**：
```bash
git --no-pager diff --check -- AGENTS.md CLAUDE.md README.md task_plan.md findings.md progress.md
```
- 結果：通過，無 whitespace / patch 格式錯誤

**修改檔案**：
- AGENTS.md
- CLAUDE.md
- README.md
- task_plan.md
- findings.md
- progress.md

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

**驗證結果**：
```bash
bash scripts/check_file_lines.sh
git --no-pager diff --check -- .github/pull_request_template.md .github/workflows/file-line-limit.yml .line-limit-ignore .githooks/pre-commit Makefile AGENTS.md CLAUDE.md README.md scripts/check_file_lines.sh task_plan.md findings.md progress.md
```
- 結果：`check_file_lines.sh` 通過（含 warning）；`diff --check` 通過（無 whitespace / patch 格式錯誤）

**修改檔案**：
- scripts/check_file_lines.sh（新增）
- .line-limit-ignore（新增）
- .github/workflows/file-line-limit.yml（新增）
- .githooks/pre-commit（新增）
- Makefile
- .github/pull_request_template.md
- AGENTS.md
- CLAUDE.md
- README.md
- task_plan.md
- findings.md
- progress.md

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

<!-- SPECTRA:START v1.0.2 -->

# Spectra Instructions

This project uses Spectra for Spec-Driven Development(SDD). Specs live in `openspec/specs/`, change proposals in `openspec/changes/`.

## Use `/spectra-*` skills when:

- A discussion needs structure before coding → `/spectra-discuss`
- User wants to plan, propose, or design a change → `/spectra-propose`
- Tasks are ready to implement → `/spectra-apply`
- There's an in-progress change to continue → `/spectra-ingest`
- User asks about specs or how something works → `/spectra-ask`
- Implementation is done → `/spectra-archive`
- Commit only files related to a specific change → `/spectra-commit`

## Workflow

discuss? → propose → apply ⇄ ingest → archive

- `discuss` is optional — skip if requirements are clear
- Requirements change mid-work? Plan mode → `ingest` → resume `apply`

## Parked Changes

Changes can be parked（暫存）— temporarily moved out of `openspec/changes/`. Parked changes won't appear in `spectra list` but can be found with `spectra list --parked`. To restore: `spectra unpark <name>`. The `/spectra-apply` and `/spectra-ingest` skills handle parked changes automatically.

<!-- SPECTRA:END -->

# CLAUDE.md

本檔是 Claude 在本 repo 的執行補充。共通規範（專案總覽、目錄結構、全部建置/測試/gate 命令、程式碼樣式、命名、錯誤處理、測試要求、安全、禁止事項、行數規範、commit/PR 慣例）的**單一來源是 `AGENTS.md`**——本檔不重複那些內容；兩檔重疊或衝突時，一律以 `AGENTS.md` 為準。

## 必讀路由（動手前）

1. `AGENTS.md`：共通規範入口。所有命令（`make build`、gate/migration/smoke、前後端驗證基準）都在那裡，不要憑記憶執行。
2. 依修改檔案類型讀 `.github/instructions/` 對應規範：
   - `*.go`、`go.mod`、`go.sum` → `go.instructions.md`
   - `*.tsx`、`*.jsx`、`*.js`、`*.css`、`*.scss` → `reactjs.instructions.md`
   - `*.ts` → `reactjs.instructions.md` ＋ `typescript-5-es2022.instructions.md`
3. 本檔其餘章節：只放 repo 特定規則與架構脈絡。

## 系統資料流（高階）

1. 協議連線：`connector/manager` 依設備配置建立 / 管理連線。
2. 資料採集：`collector/scheduler` 依 polling 設定輪詢 Tag。
3. 映射處理：`mapping` 與 transform 將來源值轉為 Point / Tag / Output 目標格式。
4. 儲存與對外：寫入 storage，並透過 API / SSE 提供前端。

## Repo 特定規則

- 前端主線入口為 `/studio/v2`；`/studio` 保留為既有完整工作台 fallback；工程測試工具集中於 `/test`。前端 build 產物嵌入 `cmd/test_ui/static`。（此行與 AGENTS.md 重複，為下兩條 inventory 規則提供上下文，屬刻意錨點，維護時勿刪。）
- Commit 訊息用繁中祈使句主旨（例：`補強...`、`完成...`，見 AGENTS.md），**不用**全域規範的 `<type>: <description>` 格式——專案慣例優先。
- 任務涉及 `/studio`、`/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*` 或 `docs/technical/studio-surface-inventory/` 時：先讀 `docs/technical/studio-surface-inventory/` 下的 `START_HERE.md`、`context.json`、`CURRENT_STATE.md`，不足以回答時才展開該目錄完整 md/html 文件。
- 修改 `docs/technical/studio-surface-inventory/` 內任何文件，必須同步寫入該目錄的 `changelog.sqlite`；changelog 一律用 `go run ./cmd/studio_inventory_changelog ...` 管理，至少留下 `summary`、`surface`、`files`、`reason`。

## 執行與回報要求（AGENTS.md 之外的補充）

- 實作不可只交付最小可動版本；錯誤處理與邊界條件屬任務範圍，不是後續工作。
- 每次任務的完成回報必含四項：修改檔案、驗證結果、風險、後續建議。
- 文件任務至少執行 `git diff --check`，確認無格式 / 空白異常。
- 未能在當前環境執行完整測試時，回報必須明列「已執行」與「未執行」項目及對應風險。

## 文件維護

- 共通規範調整 → 改 `AGENTS.md`；Claude 執行脈絡或 repo 補充脈絡 → 改本檔。
- 本檔於 2026-07-03 精簡為路由（原 197 行完整版見 git history，另備份於 `~/.claude/backups/2026-07-03-fable-institution/go_gateway/CLAUDE.md`）。維護原則：與 `AGENTS.md` 重複的規範內容一律刪除、不回加（標註「刻意錨點」者除外）；新增內容前先確認 `AGENTS.md` 沒有涵蓋。

# 2026-04-09 Frontend Redesign Workflow Skill Spec

## 目標

將本次「盤點現有前後端、重定義 IA、生成 AI 設計 prompt、規劃接回 backend」的流程 skill 化，讓未來可以快速重新評估不同前端版本。

## Skill 名稱建議

通用版：

`frontend-redesign-integration-review`

專案版：

`studio-frontend-redesign-review`

## 觸發情境

以下情境應觸發此 skill：

- 重新設計目前專案前端
- 想比較不同前端 IA / layout 方向
- 想用 Google Stitch / Pencil MCP 先做 screen exploration
- 想知道 AI 產版後怎麼接現有 backend
- 想盤點目前 API 與前端的耦合位置
- 想同時規劃 `/studio` 與 `/test`

## Skill 輸入

### 必要輸入

- repo root
- 目前主產品入口 route
- 目前工程工具入口 route
- redesign 目標

### 選填輸入

- 指定要重構的主線物件
- 指定是否含 `/test`
- 指定 AI 工具：`stitch` / `pencil` / both
- 指定是否要產 OpenSpec change

## Skill 輸出

skill 執行後應產出：

1. `master plan`
2. `OpenSpec change draft`
3. `Google Stitch prompts`
4. `Pencil MCP planning doc`
5. `backend integration plan`
6. 每一個納入 scope 的頁面，都必須有可重用的 `baseline / comparison` 文件
7. 如果 repo 中原本沒有該頁，必須用 accepted Stitch page 產出一份 `Stitch-planned baseline`

## Workflow

### Step 1: Repo context check

必讀：

- `AGENTS.md`
- agent-specific file
- `.github/instructions/*`

再讀：

- app routes
- page entrypoints
- frontend services / hooks
- backend router
- current OpenSpec

### Step 2: Current frontend map

盤點：

- 目前 route map
- 主要 page
- 主產品頁與工具頁
- 現有 screen ownership
- 若存在獨立工具頁，例如 `/test`，必須額外記錄 page shell、component composition、state ownership

輸出：

- current route inventory
- current page inventory
- page-by-page baseline / comparison document

### Step 3: Backend API inventory

盤點：

- `internal/api/router.go`
- 對應 handlers 與 service layer
- 目前 `/studio` 用哪些 `/api/v1/datalink/*`
- 目前 `/test` 用哪些 `/api/v1/test/*` / `/debug/*`

輸出：

- API groups
- screen-to-endpoint mapping
- 目前缺口與多餘 legacy

### Step 4: Frontend-backend integration map

盤點：

- `frontend/src/services/*.ts`
- `frontend/src/hooks/datalink/*.ts`
- 頁面如何經由 hooks 與 services 取數據 / mutate
- 對於獨立工具頁，還要盤點 page-level state 與 component-level state 的 ownership

輸出：

- `route -> page -> hook -> service -> API` integration chain
- 哪些是 stable contract
- 哪些需要 adapter
- 若為工具頁，補 `page shell -> workspace -> diagnostics` 的 ownership baseline
- 若為主產品頁，補 `route -> shell -> step/workspace -> state owner -> service/API` 的 ownership baseline

### Step 5: IA redesign

輸出：

- new object model
- new workflow
- page map
- workspace skeleton
- readiness model

### Step 6: AI generation package

如果目標包含 `stitch`：

- 產 `plain language, screen-by-screen` prompt package

如果目標包含 `pencil`：

- 產 screen structure + operation planning doc

### Step 7: Backend reconnection plan

AI 產版後，不直接重寫 backend。

skill 必須產出：

- 哪些畫面直接重用現有 API
- 哪些畫面需要 view-model adapter
- 哪些地方真的需要 backend amendment
- 接線順序

## Backend reconnection protocol

### 原則

- preserve current contract first
- adapt in frontend before changing backend
- use services and hooks, not direct fetch in screen components
- create adapters for new IA

### 建議 adapter 類型

- `device-capability-adapter`
- `semantic-refinement-adapter`
- `delivery-group-adapter`
- `database-row-planner-adapter`
- `modbus-block-planner-adapter`
- `mqtt-message-planner-adapter`
- `test-console-adapter`

### 接回步驟

1. AI 生成 layout / skeleton
2. 建 route shell
3. 建 adapters
4. adapters 接現有 hooks / services
5. 補 mutation wiring
6. 補 error / retry / dry-run handling
7. 跑 targeted tests / build

## 對 `/studio` 與 `/test` 的要求

### `/studio`

- workflow-first
- domain object clear
- product-oriented
- destination model clear

### `/test`

- tool-first
- field engineer oriented
- packet/log/monitor ready
- independent route, not appended to studio flow
- 必須先產出一份「current architecture baseline」，避免每次重做 `/test` redesign 都重新探索現況

### 所有頁面的 baseline 規範

- 不只 `/test`
- `/studio` 與其 destination family pages 也必須有逐頁 baseline
- 若該頁已存在於 code，baseline 要寫 current route、page owner、state owner、service/API basis
- 若該頁不存在於 code，baseline 要寫 accepted Stitch page id、規劃用途、預期 route 與目前缺口

## Skill 建議結構

```text
frontend-redesign-integration-review/
├── SKILL.md
├── references/
│   ├── api-inventory-template.md
│   ├── integration-map-template.md
│   ├── reconnection-checklist.md
│   └── iteration-rules.md
```

專案版 `studio-frontend-redesign-review` 再額外補 repo 專屬 domain、route、API 與 workflow 規則。

## 本 repo 的 Stitch 模型預設

- 對這個 repo 的 Stitch generation / edit 工作，預設使用 `GEMINI_3_1_PRO`
- 只有在使用者明確指定其他模型時才覆蓋

## SKILL.md 應包含的核心程序

- 先盤點 repo 與 OpenSpec
- 再盤點 route / page / API / integration
- 先建立 / 更新所有頁面的 baseline / comparison 文件
- 再輸出新 IA
- 再輸出 Stitch / Pencil package
- 再輸出 backend reconnection plan
- 若 scope 夠大，必須維護 `task_plan.md` / `findings.md` / `progress.md`

## AI 產版後的迭代規則

先分問題層級，再選手段：

1. 結構錯誤：改 prompt
2. 局部 layout 問題：修生成結果
3. 方向已穩：才進前端實作與 backend 對接

判斷原則：

- IA、workflow、主次層級錯了，不要硬 patch 產出
- screen map 穩了但局部不順，可保留產出做 refinement
- route、adapter、hooks、services、verify 才屬於最後落地階段

## 成功標準

skill 完成後，至少要讓使用者得到：

- 一份可討論的新主線規格
- 一份 API inventory
- 一份 frontend-backend integration map
- 一套可逐頁比對的 baseline 文件
- 一組可拿去 AI 產版的 prompt
- 一份 AI 產版後接回 backend 的執行策略
- 一套不滿意版面時可重複執行的迭代規則

## 是否要再跑一次專案版 skill

不需要為了「讓已存在的文件生效」再跑一次。

目前文件已經產出。

之後只有在以下情況建議再跑：

- repo 結構已明顯變動
- backend API inventory 已改
- 你要重新評估另一版 IA
- 你要重產 Stitch / Pencil 套件
- 你要重新生成 integration map 與 reconnection plan

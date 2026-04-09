# 2026-04-09 Studio/Test Frontend Redesign Master Plan

本文件改為索引頁，避免單一 markdown 超過 repo 行數規範。

## 內容索引

### 1. 流程與 IA 主規格

見 [2026-04-09-studio-workflow-and-ia.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-studio-workflow-and-ia.md)

內容包含：

- `/studio` 新主線
- `/test` 定位
- `Point` / `Tag` / `Tag Group` / `Delivery Group`
- `Destination Hub` / `Share Hub`
- `Database` / `Local Modbus` / `MQTT` workspace 骨架
- readiness 與 handoff 規則

### 2. API Inventory 與前後端整合主規格

見 [2026-04-09-api-inventory-and-integration.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-api-inventory-and-integration.md)

內容包含：

- 現有 backend API inventory
- 現有 frontend service / hook integration
- `/studio` 與 `/test` 對應哪些 API
- AI 產版後如何接回現有 backend
- 建議 adapter 層

### 3. OpenSpec change

見 [proposal.md](/Users/yishow/prj/go_gateway/openspec/changes/reframe-studio-destination-flow/proposal.md) 、[design.md](/Users/yishow/prj/go_gateway/openspec/changes/reframe-studio-destination-flow/design.md) 、[tasks.md](/Users/yishow/prj/go_gateway/openspec/changes/reframe-studio-destination-flow/tasks.md)

### 4. Google Stitch prompts

見 [2026-04-09-google-stitch-prompts.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-google-stitch-prompts.md)

`/test` 更細的 page-level prompts 見 [2026-04-09-test-console-stitch-prompts.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-test-console-stitch-prompts.md)

### 5. Pencil MCP plan

見 [2026-04-09-pencil-mcp-plan.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-pencil-mcp-plan.md)

`/test` 現況分析與新 IA 見 [2026-04-09-test-console-analysis-and-ia.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-test-console-analysis-and-ia.md)

### 6. Skill 定義

設計說明見 [2026-04-09-redesign-workflow-skill-spec.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-redesign-workflow-skill-spec.md)

通用版 skill：

- [frontend-redesign-integration-review/SKILL.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/skill-package-installable/frontend-redesign-integration-review/SKILL.md)

正式 skill package 草案見：

- [SKILL.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/skill-package-installable/studio-frontend-redesign-review/SKILL.md)
- [api-inventory-template.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/skill-package-installable/studio-frontend-redesign-review/references/api-inventory-template.md)
- [integration-map-template.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/skill-package-installable/studio-frontend-redesign-review/references/integration-map-template.md)
- [reconnection-checklist.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/skill-package-installable/studio-frontend-redesign-review/references/reconnection-checklist.md)

### 7. AI 版面迭代手冊

見 [2026-04-09-ai-layout-iteration-playbook.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-ai-layout-iteration-playbook.md)

內容包含：

- 何時該改 prompts
- 何時該沿用成品做局部 refinement
- 何時才該進入前端實作修補
- 如何避免一邊修 UI 一邊把前後端契約修壞

## 現在的正式方向

- `/studio` 主線：`Device -> Source -> Tag -> Destination`
- `/test`：`Field Engineer Debug Console`
- `Destination` 內任一合法目的地完成即可
- `Database` 與 `Share / Publish` 為平行去向
- `Share / Publish` 下再分 `Local Modbus` 與 `MQTT`
- skill 的終點不是產出 mockup，而是 **AI 產前端後完成與現有 backend 對接與驗證**
- skill 採雙層：
  - 通用版：可跨 repo 重用
  - 專案版：保留本 repo 的 `/studio` / `/test` domain 規則

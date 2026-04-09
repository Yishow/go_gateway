# Studio 主線重定義進度

**開始時間**：2026-04-09
**當前 Phase**：IA / flow definition

---

## 執行記錄

### 2026-04-09

#### Session 1：重置 planning records 並建立新主線
- ✅ 已讀取現有 `/studio` route、Output target switch、相關 OpenSpec 與 planning records
- ✅ 已確認目前程式與 spec 仍把 Output 視為 `Local Modbus + Database` 的 unified workspace
- ✅ 已與使用者確認新的產品主線：
  - 初版先收斂到 `Device -> Source -> Tag -> Database`
  - `Local Modbus` / `MQTT` 為 `Share / Publish` 附屬能力
- ✅ 已與使用者確認 IA 方向：
  - `Share / Publish` 採獨立子路由
  - 路徑暫定為 `/studio/share`
- ✅ 後續補充已修正前一版假設：
  - `Database` 不一定是唯一終點
  - 有些情境會在整理資料後直接走 `Share / Publish`
  - 因此主線抽象已改寫為 `Device -> Source -> Tag -> Destination`
- ✅ 已與使用者確認最後一步命名採 `Destination`
- ✅ 已與使用者確認 `Destination` 採 `Destination Hub -> destination-specific workspace` 的 IA
- ✅ 已與使用者確認主線完成條件：
  - `Destination` 內任一合法目的地完成即可
  - 不要求 `Database` 與 `Share / Publish` 同時完成
- ✅ 已與使用者確認 `Destination Hub` 可做情境化推薦，但不自動跳轉
- ✅ 已與使用者確認 `Share Hub` 只作能力入口頁，不承擔配置 / 監控總覽
- ✅ 已與使用者確認新的 destination UX 原則：
  - 各目的地共用上游資料上下文
  - 但 `Database`、`Local Modbus`、`MQTT` 都要有各自的批次規劃便利操作
  - 不可退化成逐筆綁定
- ✅ 已與使用者確認 `Delivery Group` 作為三種目的地共用的 grouping 抽象，但不升成新的主線 step
- ✅ 已與使用者確認 `Delivery Group` 採全域主物件 + destination projection / override
- ✅ 已與使用者確認 destination 頁面的主列表單位採 `Delivery Group`
- ✅ 已與使用者確認 `Database Workspace` 採三區式骨架：
  - 左 `Delivery Groups`
  - 中 `Row Planner`
  - 右 `Schema / Apply Preview / Validation`
- ✅ 已與使用者確認 `Database` 主操作模型採 `Template-assisted row composition`
- ✅ 已與使用者確認 `Local Modbus` 主操作模型採 `Block-first register composition`
- ✅ 已與使用者確認 `MQTT` 主操作模型採 `Template-assisted message composition`
- ✅ 已與使用者確認 `Point` 在新主線中定位為 internal acquisition unit，而非 primary operator object
- ✅ 已與使用者確認 `Point` 在 UI 中只保留於 Source / Inspector / Traceability / Exception 四個位置
- ✅ 已與使用者確認 `Tag` 為新主線中的 primary semantic data object
- ✅ 已與使用者確認 `Tag step` 採 `Semantic Refinement Board` 模型
- ✅ 已與使用者確認 `Tag Workspace` 採三區式骨架與第一版便利功能
- ✅ 已與使用者確認 `Tag Group -> Delivery Group suggestions` 的 handoff 模型
- ✅ 已補充新範圍：
  - 新 OpenSpec change 不只是 backend integration note，而是完整前端主線重定義
  - `/test` 工程工具頁也納入本輪 redesign / API inventory / integration package
- ✅ 已與使用者確認 `/test` 定位為 `Field Engineer Debug Console`
- ✅ 已與使用者確認 `/test` 與 `/studio` 共用設計語言，但 `/test` 明顯更工具化、更高密度
- ✅ 已將舊的 `task_plan.md`、`findings.md`、`progress.md` 歸檔至：
  - `docs/planning-archive/2026-04-09-studio-mainline-reset/`
- ✅ 已重置根目錄三份 planning files，只保留 2026-04-09 之後的新路線
- ✅ 已完成第一輪 5 份輸出：
  - `docs/frontend-redesign/2026-04-09-studio-test-redesign-master-plan.md`
  - `docs/frontend-redesign/2026-04-09-google-stitch-prompts.md`
  - `docs/frontend-redesign/2026-04-09-pencil-mcp-plan.md`
  - `docs/frontend-redesign/2026-04-09-redesign-workflow-skill-spec.md`
  - `openspec/changes/reframe-studio-destination-flow/` change package
- ✅ 已將過長的 master 文檔拆成：
  - `docs/frontend-redesign/2026-04-09-studio-workflow-and-ia.md`
  - `docs/frontend-redesign/2026-04-09-api-inventory-and-integration.md`
  - `docs/frontend-redesign/2026-04-09-studio-test-redesign-master-plan.md`（索引頁）
- ✅ 已建立正式 skill package 草案：
  - `docs/frontend-redesign/skill-package-installable/studio-frontend-redesign-review/SKILL.md`
  - 對應 `references/` 模板檔
- ✅ 已將現有 backend API inventory、frontend service/hook integration、AI 產版後接回 backend 策略納入文件
- ✅ 已補正式 installable skill 的 `reconnection-checklist.md`
- ✅ 已補 AI 產版後的 `prompt -> generated refinement -> implementation` 迭代手冊
- ✅ 已完成 `/test` 現況診斷與新 IA 定稿：
  - `docs/frontend-redesign/2026-04-09-test-console-analysis-and-ia.md`
- ✅ 已完成 `/test` page-level Stitch prompts：
  - `docs/frontend-redesign/2026-04-09-test-console-stitch-prompts.md`
- ✅ 已同步修正 master plan、Google Stitch index、Pencil MCP plan 中 `/test` 的結構描述
- ✅ 已新增通用版 skill package：
  - `docs/frontend-redesign/skill-package-installable/frontend-redesign-integration-review/SKILL.md`
  - 對應 `references/` 通用模板檔
- ✅ 已執行 `git diff --check`，新增文件目前無格式錯誤

## 下一步

- 補 `/test` 的 Pencil page-level 結構
- 補 OpenSpec spec delta 深度
- 視需要正式安裝 / 驗證 skill packages

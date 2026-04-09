# Studio 主線重定義執行計畫

**開始時間**：2026-04-09
**目標**：重定義 `/studio` 為以 Database 為主線終點的產品流程，並將資料分享能力拆為附屬工作區，後續可直接餵給 Google Stitch / Pencil MCP，並沉澱成可重複使用的 skill。

---

## 已確認決策

- [x] `Local Modbus` 不再屬於主線 step 完成條件
- [x] `MQTT` 與 `Local Modbus` 同屬資料分享能力
- [x] 資料分享能力統一命名為：`Share / Publish`
- [x] `Share / Publish` 採獨立子路由，而不是繼續塞在主線 step 內
- [x] IA 方向選定：`/studio/share`
- [x] 最終主線命名定為：`Device -> Source -> Tag -> Destination`
- [x] `Destination` 底下至少包含兩條目的地：
  - `Database`
  - `Share / Publish`
- [x] `Destination` 先進入 `Destination Hub`，再分流到各自 workspace
- [x] 主線完成條件：
  - 在 `Destination` 內完成任一合法目的地即可
  - 不要求 `Database` 與 `Share / Publish` 同時完成
- [x] `Destination Hub` 可提供情境化推薦，但不可自動跳轉
- [x] `Destination` 內引入共用抽象：`Delivery Group`
  - 供 `Database` / `Local Modbus` / `MQTT` 共用
  - 不升成新的主線 step
- [x] `Delivery Group` 採共用主物件 + destination projection / override
- [x] `Destination` 內的主列表單位為 `Delivery Group`，不是單一 tag

---

## 當前主線 IA 草案

### 主線
1. `Device`
2. `Source`
3. `Tag`
4. `Destination`

### Destination 結構
- `Destination Hub`
  - `Database`
  - `Share / Publish`
  - 可依目前資料形態 / operator intent 顯示推薦，但保留人工選擇

### 共用抽象
- `Delivery Group`
  - 一組有語意的 tag / point 集合
  - 可被不同目的地重用
  - 由系統先建議、操作員再調整
  - 作為全域主物件存在
  - 各目的地只在此基礎上做 projection / override
  - 在 destination 相關頁面中作為主列表 / 主操作單位
  - 在各目的地內分別投影為：
    - `Database row / column set`
    - `Local Modbus register block`
    - `MQTT topic + payload template`

### 附屬能力
- `/studio/share`
  - `Local Modbus`
  - `MQTT`
  - `Share Hub` 只作能力入口頁，不承擔配置與監控總覽

### 目前產品敘事
- 主線負責：從 PLC 取資料、規劃來源、完成 Tag review，並決定資料送往哪個目的地
- 目的地至少包含：
  - `Database`
  - `Share / Publish`
- `Destination` 不直接落到某個子頁，而是先經過一個 hub 讓操作員選擇目標去向
- 只要 `Destination` 內至少一條去向可用，主線即可視為完成
- `Share / Publish` 仍包含：
  - `Local Modbus`
  - `MQTT`

---

## 待完成項目

### Phase 1：資訊架構與操作流定稿
- [in_progress] 定義 `Destination` step 的完成條件與主要操作節奏
- [pending] 定義 `/studio/share` 的入口位置、進入時機、返回主線方式
- [pending] 定義 `Local Modbus` / `MQTT` 在 share workspace 的並列或分層關係
- [in_progress] 定義各目的地的批次規劃便利操作（Database row planner、Local Modbus register block planner、MQTT publish template）
- [pending] 釐清哪些狀態屬主線 readiness，哪些只屬 share readiness

### 已定義的 workspace 骨架
- [x] `Tag Workspace`
  - 左：`Candidate / Tag Groups`
  - 中：`Semantic Editor`
  - 右：`Traceability / Preview / Readiness`
  - 主操作模型：`Semantic Refinement Board`
  - 便利功能：
    - `Auto naming suggestion`
    - `Rename with pattern`
    - `Data type confirm / convert hint`
    - `Unit / meaning annotation`
    - `Keep / drop toggle`
    - `Merge / split candidate`
    - `Key tag / critical tag marker`
    - `Destination-readiness preview`
- [x] `Database Workspace`
  - 左：`Delivery Groups`
  - 中：`Row Planner`
  - 右：`Schema / Apply Preview / Validation`
  - 主操作模型：`Template-assisted row composition`
  - 便利功能：
    - `Auto column suggestion`
    - `Rename all with pattern`
    - `Column order drag sort`
    - `Drop / keep toggle`
    - `Data type hint / coercion warning`
    - `Row preview`
- [x] `Local Modbus Workspace`
  - 左：`Delivery Groups`
  - 中：`Register Block Planner`
  - 右：`Register Map Preview / Conflict / Apply`
  - 主操作模型：`Block-first register composition`
  - 便利功能：
    - `Sequential auto-pack`
    - `Gap-aware auto-pack`
    - `Width-aware placement`
    - `Base register quick set`
    - `Contiguous block preview`
    - `Conflict heatmap`
    - `Dry run before apply`
- [x] `MQTT Workspace`
  - 左：`Delivery Groups`
  - 中：`Payload / Topic Planner`
  - 右：`Topic Preview / Payload Preview / Validation / Apply`
  - 主操作模型：`Template-assisted message composition`
  - 便利功能：
    - `Topic pattern generator`
    - `Payload schema suggestion`
    - `Flat vs nested payload switch`
    - `Rename keys with pattern`
    - `Metadata include toggle`
    - `Sample message preview`
    - `Dry run / publish validation`

### Phase 2：頁面與 layout 重構描述
- [x] 重新拆分頁面清單，只保留新 IA 下的有效頁面
- [x] 為每個頁面撰寫 Google Stitch 版 prompt
- [x] 為每個頁面撰寫 Pencil MCP 版 layout prompt / structure hints
- [x] 為主線與附屬能力各自定義更順暢的操作 choreography
- [x] 重定義 `Point` / `Tag` 在新主線中的角色與 UI 所屬層級
- [x] 定義 `Tag` step 的 `Semantic Refinement Board` 骨架與便利功能
- [x] 補 `/test` 現況分析與 session-centered 新 IA
- [x] 補 `/test` 更細的 page-level Stitch prompts

### Phase 3：規格與實作對齊
- [x] 明確列出現行 OpenSpec 與新 IA 的衝突點
- [x] 擬定新的 OpenSpec change，涵蓋：
  - `/studio` 新主線與 destination model
  - `Point` / `Tag` / `Delivery Group` 新角色
  - `Destination Hub` / `Share Hub`
  - `/test` 工程工具頁在新 IA 中的定位
  - 新前端如何對接現有 backend API contract
- [x] 規劃前端 route、state ownership、shell 導覽調整方案

### Phase 4：Skill 化
- [x] 定義這套「重新評估前端主線 / IA / layout」流程的 skill scope
- [x] 拆出固定輸入、固定輸出、決策節點與 prompt 模板
- [x] 設計下次可快速比較不同前端版本的 workflow
- [x] 補正式 installable skill 參考檔與 backend reconnection checklist
- [x] 補 AI 產版後的 prompt / refinement / implementation 迭代規則
- [x] 新增通用版 `frontend-redesign-integration-review` skill

---

## 重要限制

- `/studio` 仍是主產品入口
- `Local Modbus` / `MQTT` 不能再假裝只是 Database 的附屬後處理；它們本身也可能是有效目的地
- `/test` 仍保留為到工廠 debug 用的工程工具頁，不併入 `/studio` 主產品主線
- 後續若正式實作，需回補 OpenSpec，因現行 spec 仍把 Output 定義為 unified workspace

---

## 當前狀態

**Phase**：Design package drafted
**下一步**：審閱 5 份輸出，確認是否要：
- 補 `/test` 的 Pencil page-level 結構
- 補 OpenSpec spec delta 深度
- 正式安裝 / 驗證 skill packages

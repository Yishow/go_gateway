# 2026-04-09 Studio Workflow and IA

## 目標

定義 `/studio` 與 `/test` 的新前端資訊架構與資料抽象。

## 產品分流

### `/studio`

主產品入口，主線為：

`Device -> Source -> Tag -> Destination`

### `/test`

獨立工程工具頁，定位為 `Field Engineer Debug Console`。

## 核心資料抽象

### Point

`Point` 是 internal acquisition unit。

- machine-facing
- 主要出現在 source 展開結果、inspector、traceability、exception tools
- 不再作為主產品主流程的主要操作物件

### Tag

`Tag` 是 primary semantic data object。

- operator-facing
- 承載命名、型別、單位、資料意義
- 是後續 `Destination` 與 `Delivery Group` 的穩定輸入單位

### Tag Group

`Tag Group` 是上游語意分組。

- 用於 Tag step 內做語意整理
- 代表這些 tag 在語意上應被一起理解

### Delivery Group

`Delivery Group` 是下游交付分組。

- 系統先建議，操作員再調整
- 是全域共用主物件
- 各目的地保存自己的 projection / override
- 不升成新的主線 step

關係：

`Source -> Point -> Tag -> Tag Group -> Delivery Group -> Destination`

handoff 模型：

`Tag Group -> Delivery Group suggestions`

## `/studio` 新 IA

### Route map

- `/studio`
- `/studio?step=destination`
- `/studio/database`
- `/studio/share`
- `/studio/share/local-modbus`
- `/studio/share/mqtt`

### Step rail

- `Device`
- `Source`
- `Tag`
- `Destination`

### Destination Hub

`Destination` 先進入 `Destination Hub`，再分流到具體目的地。

目的地卡片：

- `Database`
- `Share / Publish`

規則：

- 允許依目前 tag / data shape / operator intent 顯示推薦
- 不可自動跳轉
- 只要完成任一合法目的地，`Destination` 即可視為完成

### Share Hub

`Share Hub` 只作能力入口頁。

能力卡片：

- `Local Modbus`
- `MQTT`

不承擔：

- 全域分享 KPI dashboard
- 跨協定設定總覽
- 所有配置直接堆在 hub

## 各步驟設計

### Device

目的：建立設備能力上下文。

主要任務：

- 建立 / 編輯 / clone device
- connect / probe 分階段診斷
- 暴露 protocol traits、address model、word order、unit id

### Source

目的：規劃如何從設備讀資料。

主要任務：

- source rule planning
- address lattice / register planning canvas
- live values / link mode
- 產生 point 與 coverage context

### Tag

Tag step 重定義為 `Semantic Refinement Board`。

目的：把來源資料整理成可交付的語意資料物件。

骨架：

- 左：`Candidate / Tag Groups`
- 中：`Semantic Editor`
- 右：`Traceability / Preview / Readiness`

第一版便利功能：

- `Auto naming suggestion`
- `Rename with pattern`
- `Data type confirm / convert hint`
- `Unit / meaning annotation`
- `Keep / drop toggle`
- `Merge / split candidate`
- `Key tag / critical tag marker`
- `Destination-readiness preview`

### Destination

目的：將整理後的資料送往至少一個合法目的地。

不是要求所有 output 一次配完，而是讓資料先有有效去向。

## Destination workspaces

### Database Workspace

骨架：

- 左：`Delivery Groups`
- 中：`Row Planner`
- 右：`Schema / Apply Preview / Validation`

主操作模型：

`Template-assisted row composition`

第一版便利功能：

- `Auto column suggestion`
- `Rename all with pattern`
- `Column order drag sort`
- `Drop / keep toggle`
- `Data type hint / coercion warning`
- `Row preview`

適用例：

電力儀表 10 個 tag 便利地組成同一筆 row。

### Local Modbus Workspace

骨架：

- 左：`Delivery Groups`
- 中：`Register Block Planner`
- 右：`Register Map Preview / Conflict / Apply`

主操作模型：

`Block-first register composition`

第一版便利功能：

- `Sequential auto-pack`
- `Gap-aware auto-pack`
- `Width-aware placement`
- `Base register quick set`
- `Contiguous block preview`
- `Conflict heatmap`
- `Dry run before apply`

### MQTT Workspace

骨架：

- 左：`Delivery Groups`
- 中：`Payload / Topic Planner`
- 右：`Topic Preview / Payload Preview / Validation / Apply`

主操作模型：

`Template-assisted message composition`

第一版便利功能：

- `Topic pattern generator`
- `Payload schema suggestion`
- `Flat vs nested payload switch`
- `Rename keys with pattern`
- `Metadata include toggle`
- `Sample message preview`
- `Dry run / publish validation`

## `/test` 定位

`/test` 定位為 `Field Engineer Debug Console`。

角色：

- 現場快速連線測試
- raw read/write / batch / monitor
- packet / log / debug
- profile / protocol / connection config

設計原則：

- 與 `/studio` 共用品牌 / 設計語言
- 但操作語氣更工具化、更高密度、更直接
- 允許更多 monospace、debug panel、instant status

## readiness 規則

主線 readiness：

- `Device`
- `Source`
- `Tag`
- `Destination`

Destination readiness：

- `Destination ready` 只代表至少一條合法目的地完成
- 不代表所有目的地都完成

目的地內部 readiness 可分開顯示：

- `Database readiness`
- `Share readiness`
- `Local Modbus readiness`
- `MQTT readiness`

## 與現行規格的主要衝突

- 現行仍把 `Output` 視為單一步驟內的 unified workspace
- 現行把 `Local Modbus` 和 `Database` 並列在同一個 step 中處理
- 現行沒有 `Destination Hub` / `Share Hub`
- 現行沒有 `Delivery Group` 這個共用抽象
- 現行 Point / Tag 心智仍偏舊的 rule-review / mapping 模型

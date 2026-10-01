# 2026-04-09 Studio Final Design Checklist

## 目的

依照目前最新決策，先把 `/studio` 的 redesign UI 全部確認完畢，再進入任何 frontend-backend reconnection。

這份清單只處理：

- `/studio` 的設計收斂
- Stitch 方向確認
- 哪些畫面還需要 focused refinement
- 哪些畫面需要補成完整 screen family

這份清單暫時不處理：

- `/test` 的最終接線
- backend reconnection
- adapter implementation
- typecheck / build / tests

## 最新優先順序

1. 先把 `/studio` 補好
2. `/test` 先保留 accepted 狀態
3. 之後再回頭針對目前 repo 內現有 `/test` 實作做 review
4. 等 `/studio` 與 `/test` 的 UI 全部確認完畢，再一次做接線與驗證

## `/studio` 目前狀態總覽

### 已有 accepted direction

- `Device`
- `Source`
- `Tag`
- `Destination`

### 目前仍需補完的地方

- `Tag` 的語意層與 readiness 還不夠完整
- `Destination` 目前更像 workspace，不夠像完整的 `Destination Hub`
- `Database`、`Local Modbus`、`MQTT` 需要作為完整 screen family 被定義清楚
- `MQTT Workspace` 必須完成設計，但暫時不要求可直接對接現有 backend

## 畫面層級完成標準

### A. 主線 step 完成標準

對每個主線 step，都要同時滿足：

- 有清楚的主操作物件
- 有清楚的左中右或主次區域層級
- 有明確的 readiness / next action 語意
- 不再殘留舊 datalink 頁面心智
- 可被描述成之後要實作的 screen shell

### B. Destination family 完成標準

除了主線 step 外，`Destination` 相關頁面還要額外滿足：

- `Destination Hub` 存在且語意清楚
- `Database`、`Local Modbus`、`MQTT` 各自有獨立 workspace identity
- 三者看得出來共用上游 `Delivery Group`，但操作模型不同
- 不再退化成逐筆欄位綁定畫面

## `/studio` 最後一輪建議處理順序

### 1. Tag focused refinement

目的：

- 把目前 accepted `Tag` 從「審核列表」提升成真正的 `Semantic Refinement Board`

必須補強：

- `Tag Group`
- `Delivery Group suggestion`
- `Destination-readiness preview`
- traceability 與 preview 的存在感
- 讓右側欄不只是附屬摘要，而是決策依據

不建議重生整頁。

建議作法：

- 對既有 accepted `Tag` screen 做一次 focused `edit`
- 只改結構與資訊層級，不改整體風格家族

### 2. Destination Hub

目的：

- 補出 `Destination` 這一步真正的入口層

必須具備：

- `Database`
- `Share / Publish`
- 推薦去向，但不可自動跳轉
- 清楚說明主線完成條件是「任一合法目的地完成即可」
- 顯示 `Delivery Group` readiness 或資料形態摘要

這一張很重要，因為它決定：

- `Destination` 是一個真正的 step
- 不是直接掉進某個目的地 workspace

### 3. Database Workspace 確認

目前方向：

- 左 `Delivery Groups`
- 中 `Row Planner`
- 右 `Schema / Apply Preview / Validation`

最後確認重點：

- row composition 是否夠明顯
- 同一筆 row 收多個 tag 的便利性是否有被視覺化
- 是否看得出批次命名、排序、drop/keep、preview

### 4. Local Modbus Workspace 確認

目前方向：

- 左 `Delivery Groups`
- 中 `Register Block Planner`
- 右 `Register Map Preview / Conflict / Apply`

最後確認重點：

- contiguous block 規劃是否夠明顯
- sequential / gap-aware / width-aware 的便利性是否有被表現
- 是否看得出這不是一格一格手填 register 的舊 UI

### 5. MQTT Workspace 設計完成

目的：

- 先把 `MQTT` 作為 `Destination` 目的地的 UI 設計完成
- 目前不要求直接對接 backend

必要原則：

- 它是 `Destination` family 的正式成員
- 共用 `Delivery Group` 上游抽象
- 但主操作模型必須明確是 `Topic / Payload Planner`

必須具備：

- 左 `Delivery Groups`
- 中 `Payload / Topic Planner`
- 右 `Topic Preview / Payload Preview / Validation / Apply`

必須看得出來的便利功能：

- topic pattern generator
- payload schema suggestion
- flat vs nested switch
- metadata include toggle
- sample payload preview

### 6. Share Hub 是否需要補一張

這張目前可列為選配。

如果在 Stitch 專案中發現：

- `Destination Hub -> Share / Publish -> Local Modbus / MQTT`

這個節奏不夠清楚，則補一張 `Share Hub`。

如果 `Destination Hub` 已能清楚承接到 `Local Modbus / MQTT`，則可先不補。

## `/test` 目前策略

目前不做大方向 redesign。

### 保留原則

- accepted `Command Lab`
- accepted `Live Monitor`
- accepted `Device Scan`
- accepted `RTU Polling Jobs`

### 暫不做的事

- 不再大幅重生 `/test`
- 不急著把 `/test` 接上 backend
- 不急著做 `/test` adapter implementation

### 之後要做的事

等 `/studio` UI 全部確認完後，再做：

- 以目前 repo 內既有 `/test` 實作為基準的回頭 review
- 確認哪些能力已存在
- 確認 accepted screens 與現有元件落差
- 再決定 `/test` 的接線順序

## 現階段禁止混做的事

在 `/studio` 設計尚未完全確認前，先不要：

- 開始大規模前端實作
- 開始大量改 hooks / services
- 因為某個 workspace 想接資料，就提早改 backend API
- 把 `MQTT Workspace` 誤當成現有 contract 已 ready

## 這一輪完成後的下一步

當以下條件同時成立時，才進下一階段：

- `/studio` 的主線與 destination family 都已確認
- `MQTT Workspace` 已完成設計
- `/test` accepted screens 保持穩定

之後才進：

1. `/test` 現有實作 review
2. accepted screens reconnection plan refresh
3. 一次性進入前端接線與驗證

## 建議結論

目前最合理的策略是：

- 把這輪當成 **UI 設計定稿階段**
- 不是 integration 階段

因此接下來的工作應該是：

- 完成 `/studio`
- 補齊 `MQTT`
- 收斂 `Destination family`
- 延後 `/test` 接線
- 等整體 UI 確認完再一次接回 backend

# Proposal: Workbench UX Operator Efficiency Experiment

## Why

目前系統尚未正式上線，因此這次 change 的目標不是「重構並取代現在的 UI」，而是：

1. **保留目前 `/studio` 當 baseline**
2. **額外做 3 個平行新版本**
3. 在**同一組真實 API、同一個 `/studio` route、同一組 design tokens** 下，比較不同 UI kit 帶來的操作差異

要比較的重點不是單純視覺好不好看，而是：

- 操作是否順暢
- 邏輯是否清楚
- 對系統主流程的承接是否完整

因此本 change 改寫為 **baseline + 3 個 UI kit 版本** 的正式對照實驗，而不是單一路線的新 UI 重構。

## Experiment Model

| 版本 | 角色 | UI kit / 表面 | 說明 |
| --- | --- | --- | --- |
| `baseline` | 現況基線 | Current `/studio` UI | 凍結目前 UI，作為正式比較對象 |
| `v1` | 新版本 1 | shadcn/Radix | 使用 token-first 主題，允許微調互動邏輯 |
| `v2` | 新版本 2 | MUI | 使用 token-first 主題，允許微調互動邏輯 |
| `v3` | 新版本 3 | Ant Design | 使用 token-first 主題，允許微調互動邏輯 |

### Version Strategy

- 四個版本有**同一個大目標**：提升資深操作員在 `/studio` 主線中的效率。
- `v1` / `v2` / `v3` 的差異主要來自：
  - UI kit
  - kit 所擅長的元件語彙
  - 在不偏離主流程前提下的**微調互動邏輯**
- 不做四套完全不同產品，不做四套不同後端契約。

## Shared Design Tokens

三個新版本必須先建立一套**共享的 design-token 主題層**，再映射到不同 UI kit。

### 主題來源

- 可從 `/Users/yishow/prj/awesome-design-md/design-md/` 提取視覺元素與語意線索，整理成主題來源庫。
- 本 change 目前選定的主題混合來源為：`linear.app + sentry + clickhouse`

### Token 範圍

至少包含：

- color roles
- typography
- spacing
- radius
- surface / elevation
- border / divider
- status / severity
- focus / hover / active states

### Token 抽取原則

| 來源 | 優先提取 | 不直接搬用 |
| --- | --- | --- |
| `linear.app` | 版面節奏、資訊層級、精簡但高密度的表面語彙 | 品牌識別本身與特定產品 copy |
| `sentry` | severity / diagnostics / issue 狀態表面 | 與錯誤監控產品綁定的專屬語意 |
| `clickhouse` | 表格、資料密度、指標面板的結構感 | 特定資料庫產品的品牌裝飾 |

抽取後必須先轉成**語意 token**，再決定是否映射到 kit 元件；不得直接把外部設計稿當成成品樣式搬進來。

### Token 規則

1. 三個新版本共享同一組**語意 token**。
2. 差異發生在 kit component mapping 與互動細節，不是發生在三套完全不同主題。
3. baseline 不強制回補 token 化，但必須能被同一組比較 rubric 評估。

## Shared Constraints

1. **baseline + 三個新版本共用同一組真實 API 契約**，不得使用 mock API、假資料流程或版本專屬後端語意。
2. **baseline + 三個新版本都使用同一個 `/studio` route**；版本隔離只靠 worktree / branch / port，不得額外開 `/studio-v1`、`/studio-v2` 之類平行產品路由。
3. **不影響現行系統**：目前 UI 先被凍結為 baseline，比較期間主線與既有 `/studio` 行為保持穩定。
4. **相同主流程範圍**：四套表面都必須完整覆蓋 Device / Source / Tag / Output。
5. **比較 gate 強制**：每個 phase 都必須完成 `共用基礎 -> baseline -> v1 -> v2 -> v3 -> compare`，才能進入下一個 phase。
6. **可微調互動邏輯，但不可偏離同一個大目標**：三個新版本可以因 kit 特性微調操作表面，但不能變成三套不同產品。

## Comparison Outputs

每個 `compare` 子任務都必須產出固定格式的比較結果，至少包含：

| 指標 | 說明 |
| --- | --- |
| 操作順暢度 | 點擊 / 輸入 / 確認步驟是否精簡，操作是否卡頓或繞路 |
| 邏輯清晰度 | 使用者是否容易理解畫面結構、狀態關係與下一步 |
| 對系統的完整性 | 是否完整承接 Device / Source / Tag / Output 主線，不留下明顯缺口 |
| 首屏資訊密度 | 首屏能否看見必要狀態與上下文 |
| 關鍵操作時間 | 完成指定任務的耗時 |
| 實作 / 維護風險 | 元件複雜度、檔案膨脹風險、後續維護成本 |
| 推薦結論 | 本 phase 建議保留的版本與理由 |

## Stop Conditions

若出現以下任一情況，必須先停下來補 spec / contract，不得直接硬做：

1. 真實 API 契約不足，無法讓四套表面在同條件下比較。
2. baseline 未被凍結，導致比較基線可能漂移。
3. shared design-token 主題層未定義，導致三個新版本的比較失去共同視覺語意。
4. 需要新增與 repo 現況衝突的產品路由或平行主入口。
5. 某一版需要版本專屬後端語意，導致四套表面不再可比。

## What Changes

本 change 不再描述「單一路線的新 Studio 重構」，而是改為：

1. 建立 **baseline + 3 個 UI kit 版本** 的 `/studio` 對照實驗 OpenSpec。
2. 明確定義 baseline 的量測角色，以及 `v1=shadcn/Radix`、`v2=MUI`、`v3=Ant Design`。
3. 明確定義共享 design-token 主題層與外部主題來源。
4. 明確定義後續若新增 `v4+` 版本時的擴充規則。
5. 將 UI 工作拆成矩陣式 tasks：`共用基礎 -> baseline -> v1 -> v2 -> v3 -> compare`。
6. 以比較證據與推薦版本作為主要交付，而不是直接把某一版併入主線。

## Future Variant Expansion

若之後還要增加更多版本比較，遵循以下規則：

1. 新版本必須沿用同一組 shared API 與同一組 design-token 主題。
2. 新版本必須以新的 worktree / branch / port 隔離，不新增新產品 route。
3. 新版本必須補進同一套任務矩陣，例如新增 `0.2.v4`、`1.1.v4`、`2.1.v4`。
4. 新版本必須接受與 baseline / v1 / v2 / v3 相同的 compare rubric。
5. 新版本加入前，必須先補 branch / worktree / port 命名與 shared token mapping。
6. 未補齊矩陣與 compare 欄位前，不得把新版本加入正式推薦結論。

## Non-Goals

1. 本 change 不在實驗階段引入新的永久產品路由。
2. 本 change 不以 mock 流程製作視覺展示稿。
3. 本 change 不讓三個新版本各自演化出不同後端契約。
4. 本 change 不把三個新版本做成三套完全不同產品概念。
5. 本 change 不直接宣布哪一版成為正式主線；正式收斂應在比較完成後再做決策。

## Impact

### OpenSpec

- 重寫 `openspec/changes/workbench-ux-operator-efficiency/proposal.md`
- 重寫 `openspec/changes/workbench-ux-operator-efficiency/design.md`
- 重寫 `openspec/changes/workbench-ux-operator-efficiency/tasks.md`
- 保留已完成的 `specs/datalink-api/spec.md` 作為 baseline + 三個新版本共用 API 基線

### Implementation Planning

- 後續實作將以一個 baseline worktree 與三個新版本 worktree / branch / port 並行比較
- 三個新版本先共用一套 design-token 主題，再各自映射到 shadcn/Radix、MUI、Ant Design
- 每個 worktree 內部可自由調整模組結構，但對外都必須維持 `/studio` 與同一組 API 契約
- 主線 `/studio` 在推薦版本決定前不應被實驗性 UI 直接覆蓋

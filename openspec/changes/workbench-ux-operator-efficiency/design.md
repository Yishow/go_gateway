# Design: Workbench UX Operator Efficiency Experiment

## 設計目標

本 change 的目標不是立刻交付單一新 UI，而是先**凍結目前 UI 當 baseline**，再建立**三個 UI kit 版本**，在相同 API、相同主流程與相同 design-token 語意下，驗證哪一種表面最能提升資深操作員效率。

設計必須同時滿足五個前提：

1. **主線穩定**：目前 canonical `/studio` 先被凍結為 baseline，不被實驗直接取代。
2. **版本可比**：baseline 與三個新版本都只能改 UI 表面，不得各自擁有不同的 API 語意。
3. **路由收斂**：baseline 與三個新版本都維持 `/studio`，不創造平行產品入口。
4. **token 共享**：三個新版本先共用同一套 design-token 主題，再映射到不同 kit。
5. **證據導向**：每個 phase 都產出固定格式的比較結論，再決定是否前進。

## 實驗拓樸

```text
                             shared API-ready base
                                     │
        ┌───────────────┬────────────┼────────────┬───────────────┐
        │               │            │            │               │
        ▼               ▼            ▼            ▼               ▼
   baseline worktree   worktree v1  worktree v2  worktree v3   mainline
   branch: woe-base    branch: woe-radix branch: woe-mui branch: woe-antd current
   port:   4173        port:   4174      port:   4175    port:   4176     stable
   route:  /studio     route:  /studio   route:  /studio route:  /studio
        │               │            │            │
        └───────────────┴──── same real backend API ───────────────┘

baseline is a frozen snapshot of the current UI for comparison;
mainline stays stable until a winner is chosen
```

### 共用基線

- baseline 與三個新版本都從 **Phase -1 API 已完成** 的共同基線分出。
- baseline 必須在任何新版本實驗前先凍結，作為可重複量測的 current UI 快照。
- 若實驗中發現 API 契約不足，只能補到**共同基線**，再同步到 baseline 與三個新版本，不得只修某一版。
- 版本識別依靠 branch / worktree / port，不依靠 route path。

## Worktree / Branch / Port 命名

### 建議命名

| 版本 | branch | worktree 目錄名 | frontend port | 備註 |
| --- | --- | --- | --- | --- |
| `baseline` | `woe-base-current-ui` | `.worktrees/woe-base-current-ui` | `4173` | 凍結目前 `/studio` 作比較基線 |
| `v1` | `woe-v1-radix` | `.worktrees/woe-v1-radix` | `4174` | shared tokens + shadcn/Radix |
| `v2` | `woe-v2-mui` | `.worktrees/woe-v2-mui` | `4175` | shared tokens + MUI |
| `v3` | `woe-v3-antd` | `.worktrees/woe-v3-antd` | `4176` | shared tokens + Ant Design |

### 執行順序

1. 先切 `baseline` worktree，凍結 current UI。
2. 從 `linear.app + sentry + clickhouse` 抽出 shared design-token 語意層。
3. 再建立 `v1` / `v2` / `v3` worktree，分別接上各自 UI kit。
4. 每個 phase 都依序完成 `baseline -> v1 -> v2 -> v3 -> compare`。
5. compare 結論完成前，不進入下一個 phase。

### Worktree 操作稿

```bash
# 1. 建立四個 worktree
git worktree add .worktrees/woe-base-current-ui -b woe-base-current-ui
git worktree add .worktrees/woe-v1-radix -b woe-v1-radix
git worktree add .worktrees/woe-v2-mui -b woe-v2-mui
git worktree add .worktrees/woe-v3-antd -b woe-v3-antd

# 2. 啟動共享 backend（任一選定工作目錄即可，但 API 契約必須一致）
./start.sh --dev-mode --port 8080

# 3. 在各 worktree 啟動前端 dev server
cd .worktrees/woe-base-current-ui/frontend && VITE_DEV_PORT=4173 VITE_API_PROXY_TARGET=http://127.0.0.1:8080 npm run dev -- --host 0.0.0.0
cd .worktrees/woe-v1-radix/frontend        && VITE_DEV_PORT=4174 VITE_API_PROXY_TARGET=http://127.0.0.1:8080 npm run dev -- --host 0.0.0.0
cd .worktrees/woe-v2-mui/frontend          && VITE_DEV_PORT=4175 VITE_API_PROXY_TARGET=http://127.0.0.1:8080 npm run dev -- --host 0.0.0.0
cd .worktrees/woe-v3-antd/frontend         && VITE_DEV_PORT=4176 VITE_API_PROXY_TARGET=http://127.0.0.1:8080 npm run dev -- --host 0.0.0.0
```

#### 操作注意

- backend 共享 `8080` 只是預設建議；若改用其他 port，四個前端都必須同步更新 `VITE_API_PROXY_TARGET`。
- compare 時要固定使用同一組測試資料與同一組 backend 狀態。
- 若新增 `v4+`，沿用同樣格式分配新的 branch / worktree / port。

### Shared Backend Rule

- 優先使用同一個共享 backend API instance 作為三個新版本與 baseline 的比較來源。
- 若開發期間因 worktree 隔離需要各自啟動 backend，也必須保證：
  - API 契約完全一致
  - 測試資料來源一致
  - compare 時的量測條件一致

## Shared Design Tokens

### 外部主題來源

- 主題元素來源：`/Users/yishow/prj/awesome-design-md/design-md/`
- 此來源是**靈感與元素庫**，不是直接照抄單一品牌頁面。
- 本 change 目前選定的 shared token 參考混合為：`linear.app + sentry + clickhouse`

### 主題抽取分工

| 來源 | 主要吸收的元素 | 轉成 token 時的責任 |
| --- | --- | --- |
| `linear.app` | 版面節奏、留白比例、精簡工具感 | 轉成 spacing / layout / neutral surface token |
| `sentry` | severity、issue、diagnostic callout | 轉成 status / severity / diagnostic emphasis token |
| `clickhouse` | 表格密度、資料面板、指標分組 | 轉成 data-surface / table / density token |

### 抽取步驟

1. 先從三個參考來源收集候選元素。
2. 將候選元素分類為 `preserve`、`adapt`、`reject`。
3. 只把 `preserve` / `adapt` 類元素轉成語意 token。
4. 完成 shared token 表後，再做 kit mapping。
5. 若某個 kit 無法自然承接某個 token，優先調整 mapping，不直接改 shared token 語意。

### token 層責任

三個新版本必須先整理出共享的語意 token，至少涵蓋：

- color roles
- typography scale
- spacing scale
- radius
- elevation / surface
- border / divider
- severity / status
- interaction states（hover / active / focus / selected）

### mapping 規則

| 層級 | 責任 |
| --- | --- |
| shared design tokens | 定義共同主題語意與視覺一致性 |
| UI kit mapping | 把 token 映射到 shadcn/Radix、MUI、Ant Design 的元件 props / class / theme API |
| variant interaction | 根據 kit 特性微調局部互動邏輯與元件組合 |

### 建議語意命名

| token family | 建議命名 | 用途 |
| --- | --- | --- |
| `surface.*` | `surface.app`, `surface.panel`, `surface.elevated`, `surface.overlay` | 控制主背景、卡片、浮層、覆蓋層 |
| `text.*` | `text.primary`, `text.secondary`, `text.muted`, `text.inverse` | 控制主要文字、次要說明、弱化文字與反白文字 |
| `accent.*` | `accent.primary`, `accent.subtle`, `accent.strong` | 控制主強調色與其強弱層次 |
| `status.*` | `status.success`, `status.warning`, `status.error`, `status.info` | 控制狀態與 severity 呈現 |
| `border.*` | `border.default`, `border.muted`, `border.strong` | 控制常規分隔、弱分隔與強分隔 |
| `focus.*` | `focus.ring`, `focus.offset`, `focus.invalid` | 控制鍵盤 focus、錯誤 focus 與可及性高亮 |
| `density.*` | `density.compact`, `density.default`, `density.relaxed` | 控制表格、表單、列表的資訊密度節奏 |
| `data.*` | `data.rowHover`, `data.rowSelected`, `data.metricAccent` | 控制資料表面互動與指標強調 |

### 重要限制

1. 三個新版本必須共用同一組語意 token。
2. 差異可以來自 kit 表達能力與微調互動邏輯，不可以來自三套完全不同主題。
3. baseline 不要求回補 token 化，但必須用同一組比較 rubric 評估。

## 後續擴充版本規則

若未來要再加入 `v4+` 版本，遵循以下設計規則：

1. **先擴矩陣，再做 UI**：先補 proposal / design / tasks 的版本欄位與 compare 格式，再開始實作。
2. **沿用 shared token 主題**：新增版本只能映射既有 token 語意，不可自創另一套主題宇宙。
3. **沿用相同 route / API**：新增版本仍只能使用 `/studio` 與同一組真實 API。
4. **先分配新命名**：新增版本前，先補新的 branch / worktree / port 命名。
5. **沿用相同比較標準**：新增版本必須補 baseline 對照與完整 compare 證據，才能納入推薦。

### `v4+` 操作手順

1. 在 proposal 補上一列新版本定義（例如 `v4 = Mantine`）。
2. 在 design 補新版本的 kit mapping、branch / worktree / port。
3. 在 tasks 為每個 phase 補上對應的 `v4` 子任務。
4. 重新執行 `openspec validate ...`。
5. 建立新 worktree 並接上 shared tokens。
6. 依同樣的 compare rubric 取得 baseline 對照證據。
7. 完整 compare 前，不得把 `v4+` 放進推薦結論。

## baseline 與三個 kit 版本

| 版本 | UI kit / 表面 | 允許的差異 | 預期優勢 | 主要風險 |
| --- | --- | --- | --- | --- |
| `baseline` | current `/studio` | 不新增新設計，只量測現況 | 能提供真正的現況參考值 | 若 baseline 漂移，後續比較全部失真 |
| `v1` | shadcn/Radix | 可依 primitives 自由組合，微調高頻互動表面 | 彈性最高，容易做精準工程感表面 | 需要較多組裝與一致性約束 |
| `v2` | MUI | 可利用成熟表單 / data-display 元件，微調操作節奏 | 表單與資料結構表現穩定 | 若控制不好，視覺可能過於制式 |
| `v3` | Ant Design | 可利用資料輸入、表格、工作流元件，微調明確操作回饋 | 對狀態、表格、批次操作很有優勢 | 若過度使用，畫面可能顯得厚重 |

### 版本策略

- 四個版本有**同一個大目標**：提升 `/studio` 主流程的操作效率。
- `v1` / `v2` / `v3` 可以因 kit 能力差異而微調互動邏輯，例如：
  - 表單拆分方式
  - action 位置
  - 狀態提示的呈現
  - table / card / drawer / modal 的組合
- 一般情況下不允許演化成三套完全不同的資訊架構或產品流程。
- **例外：Phase 1 Device 重做 round** 已被直接 review 指出「三版差異不足且過度接近 baseline」，因此在這一個 phase 內，`v1` / `v2` / `v3` 被明確允許演化成三種不同的 interaction skeleton，但仍必須共用：
  - 同一個 `/studio` route
  - 同一組真實 API / hook / mutation
  - 同一份 workbench domain state
  - 同一組 shared design-token 語意層

## 共同不變項

以下契約對 baseline 與三個新版本全部強制成立：

1. **相同四步主線**
   - Device
   - Source
   - Tag
   - Output
2. **相同真實後端 API**
   - 使用已完成的 datalink API delta spec 與後端端點
   - 禁止 mock-only 流程
3. **相同 route contract**
   - 對外都使用 `/studio`
   - 不新增 `/studio-v1`、`/studio-v2`、`/tools/modbus` 等實驗專用產品入口
4. **相同驗收場景**
   - 建立 / 編輯設備並做 connect + probe
   - 建立或套用 source rule 並檢視 plan / live
   - review tag candidates 並執行 batch action
   - 設定 output 綁定並處理阻塞 / 診斷
5. **相同比較格式**
   - 操作錄影 / 截圖
   - 指定任務步數
   - 關鍵操作耗時
   - 邏輯清晰度筆記
   - 對系統完整性的評估
   - 實作 / 維護風險判斷

## 比較方法

每個 compare checkpoint 都要回收 baseline 與三個新版本的下列資料：

| 類別 | 證據 |
| --- | --- |
| 操作順暢度 | 完成場景的步數、切換次數、額外確認次數、是否繞路 |
| 邏輯清晰度 | 畫面結構是否好理解、狀態與下一步是否明確 |
| 對系統的完整性 | 是否完整承接四步主線，是否留下明顯缺口 |
| 首屏資訊密度 | 首屏是否能看見必要狀態、上下文與下一步 |
| 關鍵操作時間 | 完成指定任務的耗時 |
| 實作成本 | 元件複雜度、是否引入超長檔風險、測試成本 |
| 維護風險 | 共享邏輯可否抽離、後續擴充是否清楚 |

### 比較輸出格式

每次 compare 至少輸出：

```text
phase:
scenario:

baseline:
  操作順暢度:
  邏輯清晰度:
  對系統的完整性:
  首屏資訊密度:
  關鍵操作時間:
  實作 / 維護風險:

v1 (shadcn/Radix):
  ...

v2 (MUI):
  ...

v3 (Ant Design):
  ...

recommendation:
reason:
```

## Phase 切分設計

### Phase 0：baseline 與 token 基礎

- 先建立 baseline worktree 並凍結 current UI。
- 從外部主題來源整理共享 design-token 語意層。
- 固定 baseline / v1 / v2 / v3 的 branch / worktree / port 命名。
- 固定共同驗收場景、量測表與 compare 輸出格式。
- 確認 baseline 與三個新版本都能在各自環境以 `/studio` 啟動並連到同一組 API。

### Phase 1：Device

- 比較設備列表、編輯、clone、connect / probe 的最佳表面。
- 不允許改變底層 connect / probe API 語意。
- 可因 kit 能力不同微調表單、結果區、診斷提示的編排。

#### Phase 1R：Three Archetypes reboot

- 本輪不是「同一個 Device 畫面換不同 kit」；本輪明確要求三版各自成為**不同的工作台 archetype**。
- 觸發原因：上一輪 Device compare 被直接 review 判定為：
  - 三版差異不足
  - interaction model 沒有翻新
  - 視覺語言沒有真正對齊 `design-md`
  - 整體仍過度接近 baseline
- 因此本輪 Device phase 允許**完全重做 interaction skeleton**，但不允許改 route / API / domain contract。

#### Phase 1R 共通 domain model

- 三個 archetype 都必須把 Device 拆成同一套 domain layers：
  1. `Selection layer`：找設備、篩選、選中設備
  2. `Task layer`：目前正在執行的 create / edit / clone / connect / probe 任務
  3. `Diagnostics layer`：connect 與 probe 必須是兩個獨立可辨識的 stage
  4. `System layer`：capability、recent tests、readiness、tag / output 摘要
- 本輪不再接受 baseline 那種「先進 editor，之後才從結果裡看懂 probe」的互動模型。
- `create / edit / clone` 必須有明確 mode 語意，不得只是同一張表單套不同資料。

#### Phase 1R 三個 archetype

| 版本 | Archetype | design source mapping | interaction model | 核心目標 |
| --- | --- | --- | --- | --- |
| `v1` | Linear Control Room | 骨架 = `linear.app`；資訊區 = `linear.app` + 少量 `sentry` 診斷徽章；數據視覺 = 精簡 `clickhouse` stats strip | 左欄 device rail / 中央 task canvas / 右欄 live inspector | 低噪音、最快上手、最像產品主線 |
| `v2` | Sentry Incident Desk | 骨架 = `sentry`；資訊區 = `sentry` 主導；數據視覺 = `clickhouse` 用於 diagnostics / throughput / recent tests 摘要 | 任務導向 command center，diagnostics 是第一公民 | 最強 troubleshooting / connect-probe flow |
| `v3` | ClickHouse Data Cockpit | 骨架 = `clickhouse`；資訊區 = `linear.app` 結構 + `sentry` 告警語意；數據視覺 = `clickhouse` 主導 | 高密度 cockpit，首屏就是控制台 | 最像工業 ops console |

#### Phase 1R 視覺語言原則

- `v1 / Linear Control Room`
  - 近黑、低噪音、精密、極簡
  - Inter Variable / 細 border / indigo active state
- `v2 / Sentry Incident Desk`
  - 深紫黑、診斷導向、事件處理感
  - Rubik + uppercase labels + lime / coral / purple 狀態對比
- `v3 / ClickHouse Data Cockpit`
  - 純黑 + neon、高對比、高密度
  - heavy Inter + KPI tiles + 高可見 health / throughput signals
- 三版都不能再呈現「同一個灰黑 Tailwind 表面換 kit」的感覺。

### Phase 2：Source

- 比較 source rule 編修、模板套用與畫面模式切換的最佳表面。
- 著重在重複輸入、邏輯清晰度與資訊承接。

### Phase 3：Tag

- 比較 tag review、批次操作、diff preview 與 apply 回饋。
- 著重在 batch decision 成本、狀態回饋與系統承接完整性。

### Phase 4：Output

- 比較 output 綁定、狀態表面與阻塞診斷。
- 著重在映射狀態可見性、步驟連續性與資訊完整性。

### Phase 5：跨步驟 shell / diagnostics

- 比較 readiness、global blockers 與跨步驟診斷提示的最佳表面。
- 著重在系統主線是否被完整承接，而不是額外開 route。

### Phase 6：總結與推薦

- 在相同端對端流程下重跑 baseline 與三個新版本。
- 整理每 phase 的 compare 結論。
- 產出推薦版本、保留理由、淘汰理由與後續建議。

## Stop Conditions 與升級規則

遇到以下狀況，必須先回到 OpenSpec 補契約：

1. 某一版需要專屬 API 欄位或回應結構。
2. baseline 未被凍結，無法提供穩定比較基線。
3. shared design-token 層未定義，無法確保三個新版本比較公平。
4. 某一版需要新的永久產品 route 才能成立。
5. baseline 與三個新版本無法使用相同驗收場景比較。
6. compare 結果格式不足以支持推薦結論。
7. 新增版本未先補 OpenSpec 矩陣就直接開工。

## 實作邊界

- 本 design **不強制**三個新版本採相同內部檔案結構；版本可以在各自 worktree 內做最合適的拆分。
- 但任何內部做法都必須遵守：
  - repo 的 `/studio` 主線規範
  - 檔案行數限制
  - 同一組 API 契約
  - 同一組 shared design tokens
  - 同一組 compare gate

這樣可讓實驗真正聚焦在**同主題、同 API、不同 UI kit 表面**的比較，而不是預設某一種單一路徑。
在 Phase 1R 中，這個比較會進一步提升為**同主題、同 API、不同 archetype + 不同 kit + 不同視覺語言**的比較。

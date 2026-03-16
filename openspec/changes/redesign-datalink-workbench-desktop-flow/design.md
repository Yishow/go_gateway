## Context

目前 `/datalink/workbench` 已經能承接 `來源可視化 -> Tag 綁定 -> Local Modbus / Database` 的基本主線，但 desktop 資訊架構仍停留在過渡期狀態：

- page shell 仍保留偏 interim 的 summary dock，導致 1920×1080 下主工作區可用寬度不足
- Step 1 在主欄內再切一個 inline inspector，讓 desktop 版面雖寬實窄
- Step 2 有基本 toolbar 與內容，但缺少足夠完整的規則管理、coverage/gap、值/轉換呈現
- Step 3 / 4 已具備部分功能，但主區資訊密度與 operational feedback 還不足以支撐真實作業

本變更不是要回頭修 SmartDashboard，而是要把已存在的 `/datalink/workbench` 收斂成可正式規劃、可測試、可持續實作的新主工作台。

限制條件：

- 不重用舊 presentational components
- 保留既有 datalink hooks/services/types 作為 domain/data 底座
- 不在這一輪新增後端 template API
- 需維持 Local Modbus 與 Database 同一條主流程的產品方向
- 需優先滿足 1920×1080 desktop 穩定性與資訊密度

## Goals / Non-Goals

**Goals:**
- 將 workbench shell 收斂為單一 desktop 骨架：`StepRail + ContextBar + PrimaryWorkArea + InspectorPanel + BottomSummaryBar`
- 讓 Step 1~4 都成為新的 workbench-first workspace，而不是 legacy UI 的包裝層
- 讓 Step 2 成為真正的來源規畫 + 即時驗證工作台
- 讓 Step 3 / 4 具備足夠的資訊密度、差異預覽、readiness 與 operational feedback
- 定義可實作的 OpenSpec 契約，讓後續 planning/tasks 可直接落地

**Non-Goals:**
- 這一輪不重寫 datalink domain model 或 API 契約（除非後續 spec delta 明確要求）
- 這一輪不新增雲端同步或多人協作能力
- 這一輪不把 template 持久化擴充成後端資源
- 這一輪不刪除 legacy 路由或舊頁面，只重新定義 workbench 主線

## Decisions

### 1. 以五區 desktop shell 取代過渡期 ActionDock 版面

**Decision:**
workbench desktop shell 統一改為 `StepRail + ContextBar + PrimaryWorkArea + InspectorPanel + BottomSummaryBar`，不再保留 page-level `ActionDock` 作為常駐結構。

**Why:**
目前 layout 問題不是單一 CSS 細節，而是 shell IA drift。只要 summary dock 還是常駐外殼，Step 1/2 就會持續爭奪主欄寬度。

**Alternatives considered:**
- 保留 ActionDock，只微調寬度：無法根治 Step 1/2 被壓縮的問題。
- 把 inspector 繼續做在各 step 內：會重複產生不同 step 的 inline side panel，持續造成版面不一致。

### 2. 以 workbench-first 新元件重做每一個 step

**Decision:**
Step 1~4 均以新的 workspace/component boundary 實作，只重用 hooks/services/types/helper，不重組 legacy UI component tree。

**Why:**
使用者已明確要求不可沿用舊組件。現有 UI 的主要問題是資訊架構與畫面 composition，而不是單純樣式。

**Alternatives considered:**
- 以舊 `SmartDashboard` / `LocalModbusWorkbenchPage` 拆件重組：短期較快，但會把既有結構問題帶入新工作台。
- 只做 skinning：無法達成資訊密度與主流程收斂目標。

### 3. Step 2 採單一 `AddressCanvas` 幾何 + 多資訊視圖

**Decision:**
Step 2 使用一套固定 16-bit lattice 幾何，透過 `Plan / Live / Link` 切換資訊層，而不是切成多套獨立畫面。

**Why:**
使用者同時需要規畫、值驗證、鏈路理解。如果每個需求都變成獨立畫面，會重新製造 context switch。

**Alternatives considered:**
- 規畫頁 / live 頁 / link 頁分開：資訊清楚但心智模型破碎。
- 只保留 grid/table：無法滿足 live validation 與 downstream linkage 的可視化需求。

### 4. Step 2 template 能力限定為 local persistence

**Decision:**
Source rule template 僅保存在 browser-local storage，不新增後端 template API。

**Why:**
使用者需要實用的 rule/配方重用，但本輪主要問題是 desktop UX，不適合把範圍擴大到新的後端資源模型。

**Alternatives considered:**
- 新增 template API：功能更完整，但會把前端 redesign 擴大為 API / migration 變更。
- 不做 template：無法滿足使用者選定的「規則模板/配方」需求。

### 5. Step 3 與 Step 4 以資訊密度優先，而非按鈕優先

**Decision:**
Step 3 主區優先呈現 source→tag 資訊、命名規則、value/transform、merge/status；Step 4 主區優先呈現 register map / schema snapshot / readiness，而不是把操作按鈕堆成表單。

**Why:**
使用者反覆指出現在 UI 的問題不是功能按鈕不存在，而是看不懂、無法快速判斷、資訊太散。

**Alternatives considered:**
- 延續表單/表格主導的佈局：功能可做，但不會讓頁面進入可維運狀態。
- 把大部分資訊留給 Inspector：主區會再次失去可掃視性。

### 6. 將 Local Modbus 與 Database 視為同一條 Output workflow 的兩個 target

**Decision:**
Step 4 保留單一 `OutputWorkspace`，用 target switcher 切換 Local Modbus / Database，但共享 output candidate board、readiness 模型與 inspector trace。

**Why:**
產品目標已明確：tag 下一步是綁到本地 Modbus 或寫入資料庫。兩者若仍維持兩套獨立頁面，主流程會再次裂開。

**Alternatives considered:**
- 保持 Local Modbus 專頁：沿用成本較低，但與主工作流斷裂。
- 先只做 Local Modbus、Database 繼續 disabled：會再次讓輸出主線不完整。

### 7. Readiness 採 UI-derived cross-step model

**Decision:**
workbench 使用 UI 衍生的 `draft / ready / partial / blocked / applied` readiness 模型，不直接要求後端 enum 重構。

**Why:**
目前 device/tag/output 已存在不同 vocabularies。若要求後端先統一 enum，會拖慢 UI redesign；但若沒有共享 readiness，跨 step summary 與 output candidate 狀態就會失真。

**Alternatives considered:**
- 直接沿用各 domain enum：主畫面狀態不一致、難以比較。
- 新增後端通用 enum：太重，超出本輪範圍。

## Risks / Trade-offs

- **[Desktop shell 改動面廣]** → 以 shell/frame、step workspace、inspector、bottom bar 分層拆解，避免一次性巨改。
- **[Step 2 scope 容易膨脹]** → 嚴格限定一套幾何 + 三種資訊視圖，不再新增第四種主模式。
- **[Template 功能被誤做成大平台]** → 明確限制為 local persistence，不進後端。
- **[Output step 同時處理兩個 target 易變複雜]** → 共享 candidate/readiness/inspector，但 target-specific studio 分離。
- **[Legacy route 與新 workbench 並存造成導流混亂]** → 先保留相容入口，後續由 tasks 規劃 fallback/redirect 策略，不在本 artifact 直接刪除舊頁。

## Migration Plan

1. 先以新 shell/frame 取代現有 workbench desktop 結構。
2. 分步重做 Step 1 -> Step 2 -> Step 3 -> Step 4 的主區與 inspector surface。
3. 保留現有 data hooks / services / API contracts，先讓 UI 邏輯切換成功。
4. 補齊 i18n、keyboard/a11y、responsive 與 1920×1080 regression 驗證。
5. 在新 workbench 穩定後，再決定 legacy redirect / fallback 收斂策略。

**Rollback strategy:**
- 保持 `/datalink/workbench` 的 step-level feature flag / compatibility fallback 在切換期間可回退到上一版 route shell。
- 不在同一批 change 直接刪除 SmartDashboard 或獨立 Local Modbus route。

## Open Questions

- Local Modbus 獨立頁在這一輪實作完成後，是保留為 fallback/debug tool，還是改成導向 workbench Output step 的相容入口？
- Step 2 的 auxiliary audit surface 最終採右側次級 drawer、底部 ledger，或依 breakpoint 自適應切換？

## Why

`datalink-workbench-v2-shell` 與 `datalink-workbench-v2-step1-device` 兩個 change 已落地 shell 與 Step 1。本 change 把 Step 2「接入規則 + 記憶體網格 + Modbus Share」的完整靜態 UI 移植進 v2 容器，取代 `Step2RulePlaceholder`。

Step 2 是整個 datalink workbench 中互動最密集的一頁，扮演核心角色：

- **多規則 × 多裝置矩陣**：每個 device 可有多條 source rule，每條 rule 展開出一組點位（point）。原型 `step2-rule.jsx`（676 行）已支援多規則 tab、規則歸屬 device、自動推斷 Modbus function code、跨規則位址衝突偵測。
- **點位記憶體網格**：以 4/6/8 col grid 視覺呈現一條規則的位址範圍，每格代表一個 register；支援單擊 toggle 略過、Shift+click 範圍切換、Ctrl/⌘+click 多選、批次工具列、衝突格紅 ring + ! 徽章。這個視覺隱喻是整個 Dark Industrial Telemetry 風格的核心識別。
- **Local Modbus Share 配置**：每條規則可決定是否對外發布到 Local Modbus（重新成為 SCADA 可訂閱來源），且可手動指定 share register 起點或自動接續配置；網格格子右下角即時顯示對應 share 位址。
- **合併點位表**：跨規則的點位合併展示，可一眼看出未來會 `POST /points × N`、`POST /source-rules × M`，以及哪些位址有衝突。

本 change 不接後端 API；規則新增/刪除/skip/share 計算都在 reducer + selector 完成。後續 backend-wiring 會把 `state.rules` 與 `state.points` 接到 `POST /source-rules` 與衍生 points 的 API。

## What Changes

- 新增 Step 2 完整 UI 於 `frontend/src/features/datalink/workbench-v2/steps/step2/`，三層版面：
  - **上層** Rule tab 列：水平捲動 rule card；多裝置情境下顯示「所屬裝置」彩色點 + 名稱第三行；hover 顯示啟用 toggle + 刪除 icon。
  - **中層** 12-col grid：左 col 5（規則編輯器，含所屬裝置 select、2×2 grid 基本參數、範圍摘要與 Modbus function code chip、線性轉換 details、Local Modbus Share details）；右 col 7（當前規則點位網格 + 批次工具列）。
  - **下層** 合併點位表：跨規則合併、衝突高亮、Share 位址顯示、表尾「繼續到映射」按鈕（disabled 條件：總啟用點位 === 0 或位址衝突數 > 0）。
- 在 `state/sourceRule.ts` 落地：
  - `derivePoints(rule, deviceId, skippedSet)` 與 `deriveAllPoints(rules, fallbackDeviceId)` 對齊原型；
  - `computeShareLayout(rules, baseRegister)` 計算每條規則的 share start/stride/end，支援手動 / 自動配置；
  - `detectAddressConflicts(allPoints)` 回傳衝突位址 set；
  - `fnFromAddr(addr)` 推斷 Modbus function code（4xxxx → holding_register，3xxxx → input_register，2xxxx → discrete_input，0xxxx/1xxxx → coil）；
  - `dataTypeWidth(type)` 與 `formatAddr(n)` 格式化常數。
- 擴充 `useWorkbenchV2State` reducer 加入 rule actions（`addRule`、`removeRule`、`updateRule`、`renameRule`、`toggleRuleEnabled`、`updateRuleSkipped`、`toggleRuleShareEnabled`、`updateRuleShareStart`、`updateRuleShareStride`）。
- 衍生 state：在 selector 層計算 `selectAllPoints`、`selectConflicts`、`selectShareLayout`，避免每次 render 重算。
- **點位網格互動**：單擊 toggle skip；Shift+click 範圍切換（基於 `lastClickedIdx`）；Ctrl/⌘+click 加入/移出 `gridSelection` 多選 set；批次工具列「全部啟用 / 全部略過 / 反轉啟用 / 略過選取 / 啟用選取 / 清除選取」；網格底紋使用 24px `.bg-grid`；衝突格紅 ring + 右上 ! 徽章；多選格藍 ring；啟用且 share 開啟時格子右下浮動 mono 標籤顯示對應 share 位址。
- **Local Modbus Share 視覺整合**：規則編輯器 details 區塊顯示「Share 起始 Register」「Stride」雙欄 + 「透過 Modbus Share 對外發布此規則的點位」toggle + 「改回自動分配」reset 連結；計算後的 share start/end chip 顯示在 details summary。
- **跨規則裝置色彩**：rule card 使用 6 色循環（與 device 共用色票但獨立 index）；多裝置情境下 rule card 第三行渲染 device 顏色點 + name。
- 修改 `WorkbenchV2Shell` 在 `current === 2` 時改 render `Step2Rule`；shell spec MODIFIED placeholder requirement 從名單移除 Step 2。
- 擴充 i18n `workbench-v2` namespace 加入 Step 2 microcopy（rule tab、grid hint、批次按鈕、share details、合併表 column 名）。
- 不引入新依賴；沿用 shell + step1 共用元件與 `useDeviceColor` hook。

## Non-Goals (optional)

- 不接 `POST /source-rules` 與 bulk-create points API；保留 reducer 形狀供 backend-wiring 後續換 mutation。
- 不在本 change 處理 Modbus Share 的 server-side 設定（綁定位址、port、slave_id）；那屬於 Settings 頁，由 `datalink-workbench-v2-settings` change 落地。本 change 讀 `state.settings.modbus_share` 決定是否渲染 share 對應 UI。
- 不處理規則匯入 / 匯出（JSON 上下載）；原型也未涵蓋。
- 不支援跨裝置「複製規則」操作；新增規則僅 makeDefaultRule。
- 不處理規則拖拉排序；tab 列保持新增順序。
- 不變更 Step 1/3/4/Settings 任何行為。

## Capabilities

### New Capabilities

- `datalink-workbench-v2-step2-rule`：定義 v2 Step 2 多規則工作區的需求、點位網格互動契約、Modbus Share 位址配置規則、跨規則位址衝突偵測、Modbus function code 推斷、線性轉換預設值繼承機制、合併點位表結構。

### Modified Capabilities

- `datalink-workbench-v2-shell`：縮窄 `Placeholder step and settings surfaces` 需求範圍，把 Step 2 從 placeholder 名單移除（剩餘 Step 3/4 與 Settings 仍走 placeholder）。

## Impact

- Affected specs:
  - 新增 `openspec/specs/datalink-workbench-v2-step2-rule/spec.md`
  - 修改 `openspec/specs/datalink-workbench-v2-shell/spec.md`
- Affected code:
  - New:
    - `frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step2/index.ts`
    - `frontend/src/features/datalink/workbench-v2/state/sourceRule.ts`
    - `frontend/src/features/datalink/workbench-v2/state/selectors.ts`
    - `frontend/tests/workbench-v2/step2-rule.test.tsx`
    - `frontend/tests/workbench-v2/step2-grid.test.tsx`
    - `frontend/tests/workbench-v2/step2-share.test.tsx`
    - `frontend/tests/workbench-v2/sourceRule.test.ts`
  - Modified:
    - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts`（新增 rule actions）
    - `frontend/src/features/datalink/workbench-v2/state/types.ts`（補完 `Rule`、`Point`、`ShareLayout` 型別）
    - `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx`（current=2 切到 Step2Rule）
    - `frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx`（標記為 rollback only）
    - `frontend/src/i18n/locales/zh-TW/workbench-v2.json`、`frontend/src/i18n/locales/en/workbench-v2.json`
- 不變更後端任何檔案。

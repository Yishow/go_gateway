## 1. State 層擴充

- [ ] 1.1 在 `state/types.ts` 補完 `Rule`、`Point`、`ShareLayout` 型別（含 10 種 data_type union、function code union、四種 data_format 字串聯集），落地設計決策「拆檔策略：以「資料區塊 / 互動行為」拆 10 個元件 + 2 個 state module」的型別契約。**行為**：所有 step2 子元件可正確 import；`Rule['share_start_register']` 與 `Rule['share_stride']` 可為 `number | null`。**驗證**：`tsc --noEmit` 通過；新增 `frontend/tests/workbench-v2/types-step2.test-d.ts` type-level 斷言。
- [ ] 1.2 建立 `state/sourceRule.ts`，匯出 `derivePoints`、`deriveAllPoints`、`computeShareLayout`、`detectAddressConflicts`、`fnFromAddr`、`dataTypeWidth`、`formatAddr` 7 個純函式，落地設計決策「衍生 state：純函式 + selector hook」並實作需求 **Modbus function code inference**、**Modbus Share layout**、**Cross-rule address conflict detection**、**Rule-level scale inheritance to points**。**行為**：每個函式為純函式，輸入相同輸出相同；`computeShareLayout` cursor 用 `max` 防止手動指定點落後。**驗證**：新增 `frontend/tests/workbench-v2/sourceRule.test.ts`，7 個函式各 2–3 個 case；覆蓋 `fnFromAddr` 對 4 種 prefix 的對應、`dataTypeWidth` 對 10 種 type 的回傳值、`computeShareLayout` 的自動 / 手動 / disabled 三種混合場景、`detectAddressConflicts` 的跨規則同址 + skipped 不衝突。
- [ ] 1.3 建立 `state/selectors.ts`，匯出 `useAllPoints`、`useShareLayout`、`useConflictAddrs`、`useRulePoints` 4 個 hook，內部用 `useMemo` cache。**行為**：rule 變更後 hook 回傳值更新；其他 state 變更（如 device test status）不觸發重算（useMemo 依賴限縮）。**驗證**：新增 `frontend/tests/workbench-v2/selectors.test.tsx` 用 `renderHook` 對 4 個 hook 各別測試 cache 行為。
- [ ] 1.4 擴充 `state/useWorkbenchV2State.ts` reducer，新增 10 個 rule actions（`addRule`、`removeRule`、`updateRule`、`renameRule`、`toggleRuleEnabled`、`updateRuleSkipped`、`toggleRuleSkippedAddress`、`toggleRuleShareEnabled`、`updateRuleShareStart`、`updateRuleShareStride`），落地設計決策「改變參數時清空 skipped_addresses」並實作需求 **Rule editor with linked reset**。**行為**：`updateRule` 收到 patch 含 `start_address`/`count`/`data_type` 任一時自動 reset `skipped_addresses=[]`；`updateRule` 收到 count 時夾在 1–64；`removeRule` 在 selected 是被刪 rule 時自動切到剩餘第一個（如有）。**驗證**：新增 `frontend/tests/workbench-v2/reducer-step2.test.ts`，10 個 action 各 1–2 個 case；額外覆蓋 reset skipped、count clamp、remove cascade selection 切換。

## 2. UI 子元件

- [ ] 2.1 建立 `steps/step2/RuleTabRail.tsx`，渲染水平捲動 rule tab + 「+ 新增規則」按鈕，落地需求 **Multi-rule tab management**。**行為**：每個 tab 顯示顏色點 + inline 改名 + start_address/data_type + 點位計數；hover 顯示啟用 toggle 與 close icon；`state.devices.length >= 2` 時 tab 第三行顯示裝置色點 + 裝置名（用 `useDeviceColor`）。**驗證**：`frontend/tests/workbench-v2/step2-rule.test.tsx` 中 `describe('RuleTabRail')` 測試 add rule、close icon visibility、inline rename、multi-device device row。
- [ ] 2.2 建立 `steps/step2/RangeSummary.tsx`，顯示 `{startAddr} → {endAddr}` chip + Modbus function code chip + 「共 N register」label，落地需求 **Modbus function code inference**。**行為**：function code 依 `start_address` 第一位數字推斷；4 種 prefix 對應不同 tone（amber/slate/cyan/blue）。**驗證**：`step2-rule.test.tsx` 中 `describe('RangeSummary')` 四個 it 對應 4 種 prefix。
- [ ] 2.3 建立 `steps/step2/ScaleSection.tsx`，details 區塊含 multiplier / offset / byte order，落地設計決策「線性轉換預設值繼承到 Step 3」與需求 **Rule editor with linked reset** 的線性轉換部分。**行為**：輸入 onChange dispatch `updateRule`；不影響 skipped。**驗證**：`step2-rule.test.tsx` 中測試 multiplier 改值反映 state。
- [ ] 2.4 建立 `steps/step2/ShareSection.tsx`，details 區塊含 share_start_register / share_stride / share_enabled toggle / 「改回自動分配」reset link，落地設計決策「Modbus Share 位址布局：自動接續 + 手動覆寫」並實作需求 **Modbus Share layout** 的 UI 部分。**行為**：share_start_register 留空時 placeholder `(自動)`、reducer 接到空字串轉 null；toggle share_enabled 不影響 share_start_register；reset link 只在手動值存在時顯示。**驗證**：新增 `frontend/tests/workbench-v2/step2-share.test.tsx`，測試自動 / 手動切換、reset 行為、disabled 時 details summary 顯示「全域未啟用」chip。
- [ ] 2.5 建立 `steps/step2/RuleEditor.tsx`，組合 device select + 2×2 基本參數 grid + RangeSummary + ScaleSection + ShareSection + 底部「啟用此規則」toggle，落地需求 **Rule editor with linked reset**。**行為**：device select 在 `devices.length === 1` 時 disabled；count input min=1/max=64；改 start_address/count/data_type 觸發 skipped reset。**驗證**：`step2-rule.test.tsx` 中 `describe('RuleEditor')` 測 count clamp、data_type 切換 stride 更新、device select disabled 條件。
- [ ] 2.6 建立 `steps/step2/PointGridToolbar.tsx`，渲染 6 個批次按鈕（全部啟用 / 全部略過 / 反轉啟用 / 略過選取 / 啟用選取 / 清除選取）+ 右側 mono summary `total {count} · stride {width}`，落地需求 **Point grid batch toolbar**。**行為**：略過選取 / 啟用選取在 gridSelection.size === 0 時 disabled；清除選取只在 size > 0 時 render。**驗證**：`frontend/tests/workbench-v2/step2-grid.test.tsx` 中 `describe('PointGridToolbar')` 測 6 個按鈕行為與 disabled 條件、反轉啟用 swap 邏輯。
- [ ] 2.7 建立 `steps/step2/PointGrid.tsx`，4/6/8 col 自適應 grid，每個 cell 渲染位址 + 名稱 + 衝突徽章 + share 浮動 mono 標籤，含 3 種點擊修飾鍵 handler，落地設計決策「點位網格互動：點擊 + Shift + Ctrl/⌘ 三模式」與「衝突偵測：在 `state.points` 層去重」，並實作需求 **Point grid with modifier-key interactions**。**行為**：plain click toggle skip；Shift+click 對 `lastClickedIdx → idx` 範圍以「當前 cell skip 狀態」為基準批次切換；Cmd/Ctrl+click toggle `gridSelection` 成員；switch rule 重置 selection 與 lastClickedIdx（useEffect 依賴 `selectedRuleId`）；share 開啟時 cell 右下浮 `→{shareAddr}` mono 小字。**驗證**：`step2-grid.test.tsx` 中 `describe('PointGrid')` 測 3 種 click 模式、switch rule reset selection、share 浮動標籤、衝突紅 ring。
- [ ] 2.8 建立 `steps/step2/MergedPointTable.tsx`，渲染合併點位表（8 欄）+ sticky header + 表尾「繼續到映射」，落地需求 **Merged point table summary** 與 **Cross-rule address conflict detection** 的合併表部分。**行為**：每列依 device color + rule color 渲染；衝突列地址紅字 + ⚠ + `衝突` chip；skipped 列 opacity 40% + `跳過` chip；繼續按鈕在 `totalEnabled === 0 || conflictAddrs.size > 0` 時 disabled。**驗證**：`step2-rule.test.tsx` 中 `describe('MergedPointTable')` 測列數 = total points、衝突列高亮、繼續按鈕 disabled 條件。

## 3. 容器與整合

- [ ] 3.1 建立 `steps/step2/Step2Rule.tsx`，組裝 RuleTabRail + 12-col grid（RuleEditor + PointGrid + PointGridToolbar）+ MergedPointTable + 底部 footer，落地設計決策「拆檔策略：以「資料區塊 / 互動行為」拆 10 個元件 + 2 個 state module」的容器。**行為**：state 全部由 selector / reducer 提供；onContinue 由 footer 觸發。**驗證**：`step2-rule.test.tsx` 主測試 case `renders full Step 2 with default state`。
- [ ] 3.2 建立 `steps/step2/index.ts` barrel re-export `Step2Rule` 與必要型別。**行為**：外部 `from 'steps/step2'` 可取得 Step2Rule。**驗證**：`tsc --noEmit` 通過。
- [ ] 3.3 修改 `shell/WorkbenchV2Shell.tsx` 在 `current === 2 && view === 'flow'` 時改 render `<Step2Rule onContinue={...} />`，落地 modified 需求 **Placeholder step and settings surfaces**（Step 2 從 placeholder 名單移除）。**行為**：current=2 看到完整 Step 2；current=3/4 仍看 placeholder；settings 仍看 SettingsPlaceholder。**驗證**：`shell.test.tsx` 既有 placeholder 測試更新；新增 it 斷言 current=2 顯示 Rule tab。
- [ ] 3.4 在 `steps/Step2RulePlaceholder.tsx` 加 rollback only 註解，不刪除檔案。**行為**：檔案保留供緊急 rollback；不被 import。**驗證**：`grep -rn "Step2RulePlaceholder" frontend/src` 只在自身或註解出現。

## 4. i18n

- [ ] 4.1 擴充 `frontend/src/i18n/locales/zh-TW/workbench-v2.json`，新增 `step2.*` keys：rule tab title/subtitle、grid hint（單擊 / Shift / Ctrl）、批次按鈕、Modbus function code label、Byte order label、Share details label、合併表 column 名、繼續按鈕、`全域未啟用` chip、`改回自動分配` link。**行為**：i18n 完整覆蓋。**驗證**：`step2-rule.test.tsx` 在 zh-TW 斷言 `接入規則` 與 `當前規則點位網格` 出現。
- [ ] 4.2 擴充 `frontend/src/i18n/locales/en/workbench-v2.json` 對應英文翻譯。**行為**：對等翻譯。**驗證**：`step2-rule.test.tsx` 在 en 斷言 `Source Rules` 出現。

## 5. 驗證

- [ ] 5.1 跑 `cd frontend && npm run lint`，零錯誤。
- [ ] 5.2 跑 `cd frontend && npm run typecheck`，零錯誤。
- [ ] 5.3 跑 `cd frontend && npm run test -- --run workbench-v2`，全綠。
- [ ] 5.4 跑 `cd frontend && npm run build`，產出乾淨。
- [ ] 5.5 跑 `bash scripts/check_file_lines.sh`，確認 step2 所有新檔 ≤ 300 行。
- [ ] 5.6 手動 smoke：完整跑「新增第二條規則 → 改 count → Shift+click skip 範圍 → Ctrl+click 多選 → 批次略過選取 → 開 share → 手動改 share 起點 → 看合併表衝突高亮 → 點繼續到 Step 3 placeholder」。

## 1. State 層擴充

- [ ] 1.1 在 `state/types.ts` 補完 `Settings`、`SettingsConnector`、`TimeseriesSettings`、`SchedulerSettings`、`ModbusShareSettings`、`GeneralSettings` 全部子型別，落地設計決策「拆檔策略：12 個元件 + 1 個 state module」型別契約。**行為**：所有 settings 子元件 import 型別解析正確；General 與 Connector 與 Step 4 共用的子集對齊。**驗證**：`tsc --noEmit` 通過；新增 `types-settings.test-d.ts` 對 union 型別 assertion。
- [ ] 1.2 建立 `state/settingsDefaults.ts`，匯出 `DEFAULT_SETTINGS`（含預設 `TimeSeries Prod` connector + 4 個 section 預設值）、`makeDefaultConnector(idx)`，落地設計決策「Settings 預設值集中管理」並實作需求 **Connector pool CRUD** 的 `Default connector` Scenario 與五個 section 各自的 default。**行為**：DEFAULT_SETTINGS 的 modbus_share.enabled === true、connectors 含 1 個 ready connector；`makeDefaultConnector(2)` 回傳 name `新連線 2` + kind postgres + status unknown。**驗證**：新增 `frontend/tests/workbench-v2/settingsDefaults.test.ts` 覆蓋 DEFAULT_SETTINGS 結構與 makeDefaultConnector 命名規則。
- [ ] 1.3 擴充 `useWorkbenchV2State` reducer 加 8 個 settings actions（`updateSettings`、`updateSettingsSection`、`addConnector`、`updateConnector`、`removeConnector`、`startConnectorTest`、`completeConnectorTest`、`resetSettingsToDefaults`），落地設計決策「Connector 池：以陣列管理 + per-row status」並實作需求 **Connector pool CRUD** 的全部 Scenarios。**行為**：`updateConnector` 收到 patch 含 kind/host/port 任一時自動把 status 回 unknown；`addConnector` 不帶 payload 由 reducer 用 `makeDefaultConnector` 生；`resetSettingsToDefaults` 把整個 `state.settings` 換 DEFAULT_SETTINGS。**驗證**：新增 `frontend/tests/workbench-v2/reducer-settings.test.ts` 覆蓋 8 個 action 與 kind/host/port 觸發 status reset 邏輯。

## 2. UI 子元件

- [ ] 2.1 建立 `settings/SettingsHeader.tsx`，渲染漸層標題卡（11×11 sliders icon + 標題 + 一行說明）。**行為**：純展示元件。**驗證**：新增 `frontend/tests/workbench-v2/settings.test.tsx` 中 `describe('SettingsHeader')` 測 icon 與標題渲染。
- [ ] 2.2 建立 `settings/TimeseriesSection.tsx`，渲染 4 個 Field（write_precision / partition_interval / batch_size / retention_days），落地需求 **Timeseries strategy fields**。**行為**：select 與 number onChange 即時 dispatch `updateSettingsSection`。**驗證**：`settings.test.tsx` 中 `describe('TimeseriesSection')` 覆蓋 4 個欄位預設值與更新。
- [ ] 2.3 建立 `settings/SchedulerSection.tsx`，渲染 5 個欄位（4 個 number + 1 個 toggle auto_start），落地需求 **Scheduler defaults fields**。**行為**：dispatch `updateSettingsSection`。**驗證**：`settings.test.tsx` 中 `describe('SchedulerSection')` 測 auto_start toggle 與 4 個 number 更新。
- [ ] 2.4 建立 `settings/ModbusShareSection.tsx`，總開關放 aside；enabled 時顯示 4-col grid（bind_address / port / slave_id / base_register），停用時顯示說明文字，落地設計決策「連動效應」與需求 **Local Modbus Share configuration**。**行為**：toggle 立即切換顯隱；base_register 變更立即生效（透過 selector 流到 Step 2）。**驗證**：`settings.test.tsx` 中 `describe('ModbusShareSection')` 覆蓋 toggle 顯隱與 base_register 更新；額外加 integration test 驗證 `useShareLayout` 在 base_register 變更後回傳新位址。
- [ ] 2.5 建立 `settings/UiSection.tsx`、`settings/ApiSection.tsx`、`settings/DiagnosticsSection.tsx`，各 3–4 個 Field，落地需求 **UI / API / Diagnostics options**。**行為**：theme select 寫入 state 但不實際換主題；其他欄位即時 dispatch。**驗證**：`settings.test.tsx` 中三個 describe 各別覆蓋。
- [ ] 2.6 建立 `settings/ConnectorRow.tsx`，渲染單一 connector：類型 emoji icon + 12-col grid 表單 + 跨欄狀態列（status chip + last_check_at + 測試 / 刪除 button），落地需求 **Connector pool CRUD** 的 row 部分。**行為**：所有欄位 onChange dispatch `updateConnector`；改 kind/host/port 自動 reset status；status chip 4 種 tone（ready=success/testing=info/unreachable=error/unknown=draft）；last_check_at 只在 status==='ready' 時顯示。**驗證**：新增 `frontend/tests/workbench-v2/settings-connectors.test.tsx` 中 `describe('ConnectorRow')` 覆蓋欄位更新、status chip 切換、kind 切換 reset。
- [ ] 2.7 建立 `settings/ConnectorPoolSection.tsx`，組合 SectionCard + 「+ 新增連接器」aside button + 空狀態 + connector 列表（divide-y），含 mock 測試 timer，落地需求 **Connector pool CRUD** 與 **Connector mock test**。**行為**：addConnector dispatch；test 點擊 startConnectorTest → setTimeout 900ms → completeConnectorTest（85/15 by Math.random）。**驗證**：`settings-connectors.test.tsx` 中 `describe('ConnectorPoolSection')` 覆蓋 addConnector、mock test 成功 / 失敗 path（用 `vi.useFakeTimers()` + `vi.spyOn(Math, 'random')`）、removeConnector。
- [ ] 2.8 建立 `settings/SaveBar.tsx`，渲染 sticky 底部列（info text + 重設為預設 ghost + 儲存所有設定 success），落地設計決策「「重設為預設」與「儲存所有設定」按鈕」與「Sticky 儲存列定位」並實作需求 **Save bar actions**。**行為**：`重設為預設` 點擊先 confirm 再 dispatch resetSettingsToDefaults；`儲存所有設定` 點擊 console.warn。**驗證**：`settings.test.tsx` 中 `describe('SaveBar')` 覆蓋 confirm true/false、save no-op + console.warn spy。

## 3. 容器與整合

- [ ] 3.1 建立 `settings/SettingsPage.tsx`，組裝 SettingsHeader + ConnectorPoolSection + （TimeseriesSection + SchedulerSection 並排）+ ModbusShareSection + （UiSection + ApiSection + DiagnosticsSection 並排）+ SaveBar，落地需求 **Settings page layout**。**行為**：5 個 SectionCard 順序與並排 layout 正確；save bar sticky bottom-4 z-10；scroll 時 save bar 保持可見。**驗證**：`settings.test.tsx` 主測試 `renders all sections in order with sticky save bar`。
- [ ] 3.2 建立 `settings/index.ts` barrel re-export。**行為**：`from 'settings'` 可取得 SettingsPage。**驗證**：`tsc --noEmit` 通過。
- [ ] 3.3 修改 `shell/WorkbenchV2Shell.tsx` 在 `view === 'settings'` 時改 render `<SettingsPage />`，移除對 `SettingsPlaceholder` 的 import 與條件分支，落地 modified 需求 **Coexisting v2 workbench route** 與 REMOVED 需求 **Placeholder step and settings surfaces**。**行為**：view='settings' 顯示完整 settings；shell 不再有任何 placeholder import。**驗證**：`shell.test.tsx` 既有 settings 測試更新；新增 it 斷言 `import.*Placeholder` 不再出現於 shell tree。
- [ ] 3.4 在 `settings/SettingsPlaceholder.tsx`（與其他 4 個 step placeholder）加 rollback only 註解。**行為**：保留供 rollback；未被 import。**驗證**：`grep -rn "SettingsPlaceholder\|Step.DevicePlaceholder\|Step.RulePlaceholder\|Step.MappingPlaceholder\|Step.DatabasePlaceholder" frontend/src` 全部僅在自身或註解出現。

## 4. i18n

- [ ] 4.1 擴充 `frontend/src/i18n/locales/zh-TW/workbench-v2.json`，新增 `settings.*` keys：標題列文案（系統設定 / 資料庫連接器池、時序儲存策略、排程器、Local Modbus Share、介面偏好）、Connector 池（新增連接器、空狀態、測試 / 刪除、status 四種文字）、Timeseries（時間精度 / 分區間隔 / 批次寫入大小 / 保留天數 + hints）、Scheduler（5 個 label + 開機自動啟動 collector）、Modbus Share（綁定位址 / Port / Slave ID / 起始 Register + 全域未啟用提示文字）、UI / API / 診斷三區所有 label 與 hint、SaveBar（設定會立即套用… / 重設為預設 / 儲存所有設定）。**行為**：覆蓋。**驗證**：`settings.test.tsx` 在 zh-TW 斷言 `資料庫連接器池` 與 `Local Modbus Share` 出現。
- [ ] 4.2 擴充 `frontend/src/i18n/locales/en/workbench-v2.json` 對應翻譯（Connector Pool / Timeseries Strategy / Scheduler Defaults / Local Modbus Share / Interface / API Endpoint / Diagnostics / Reset to Defaults / Save All Settings 等）。**行為**：對等。**驗證**：`settings.test.tsx` 在 en 斷言 `Connector Pool` 出現。

## 5. 驗證

- [ ] 5.1 跑 `cd frontend && npm run lint`，零錯誤。
- [ ] 5.2 跑 `cd frontend && npm run typecheck`，零錯誤。
- [ ] 5.3 跑 `cd frontend && npm run test -- --run workbench-v2`，全綠。
- [ ] 5.4 跑 `cd frontend && npm run build`，產出乾淨。
- [ ] 5.5 跑 `bash scripts/check_file_lines.sh`，確認 settings 所有新檔 ≤ 300 行。
- [ ] 5.6 手動 smoke：完整跑「進入 settings → 新增 connector → 測試 connector 看 mock 動畫成功與失敗 → 改 Modbus Share base_register → 切到 Step 2 看 share 位址重算 → 切回 settings 點重設為預設 confirm → 全部回到初值 → 點儲存所有設定看 console.warn」。

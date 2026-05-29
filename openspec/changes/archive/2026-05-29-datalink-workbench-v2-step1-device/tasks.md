## 1. State 層擴充

- [x] 1.1 在 `frontend/src/features/datalink/workbench-v2/state/types.ts` 補完 `Device`、`DeviceTest`、`ProtocolId`、`ReadinessStage` 型別，落地設計決策「拆檔策略：以「資料區塊 / 互動行為」拆分而非單一巨檔」所需的型別契約。**行為**：`import type { Device, DeviceTest, ProtocolId, ReadinessStage } from '.../state/types'` 可在所有 step1 子元件正確解析；`DeviceTest['stages']` 為 `Record<string, { status; latency_ms? }>`。**驗證**：在 `frontend/tests/workbench-v2/types.test-d.ts` 加 type-level assertion 斷言型別 shape；`tsc --noEmit` 通過。
- [x] 1.2 建立 `frontend/src/features/datalink/workbench-v2/state/protocols.ts`，匯出 `PROTOCOLS`（6 個 protocol 元資料）、`getStagesForProtocol(protocol)` 回傳含 `group: 'connect' | 'probe'` 的 stage 序列、`getDefaultConfig(protocol)`，落地設計決策「Readiness 群組：以 `connect` / `probe` 語義標籤」。**行為**：呼叫 `getStagesForProtocol('mqtt')` 回傳 `[{id:'resolve',group:'connect'},{id:'connect',group:'connect'},{id:'subscribe',group:'probe'}]`。**驗證**：新增 `frontend/tests/workbench-v2/protocols.test.ts` 對 6 個 protocol 各自斷言序列與 default config（對齊 design Decisions 內的 defaults table）。
- [x] 1.3 建立 `frontend/src/features/datalink/workbench-v2/state/deviceColors.ts`，匯出 `DEVICE_COLORS` 常數與 `useDeviceColor(deviceId)` hook，落地設計決策「設備色彩：6 色循環常數 + hook 共用」並實作需求 **Device color cycle**。**行為**：`useDeviceColor(devices[0].id).name === 'blue'`；包到 6 個後 wrap-around 回 `blue`。**驗證**：新增 `frontend/tests/workbench-v2/deviceColors.test.tsx` 測試 6 個位置 + wrap-around；用 `renderHook` 包 provider。
- [x] 1.4 在 `state/useWorkbenchV2State.ts` 擴充 reducer 加入 9 個 device actions（`addDevice`、`removeDevice`、`updateDevice`、`updateDeviceConfig`、`renameDevice`、`changeDeviceProtocol`、`startDeviceTest`、`advanceDeviceTest`、`completeDeviceTest`、`failDeviceTest`），落地設計決策「Mock 測試動畫：使用 setTimeout + reducer action」與「Protocol 切換：清空 test 結果與 config」。**行為**：每個 action 為純函式 immutable 更新；`changeDeviceProtocol` 同時把 `test=null` `status='draft'` `config=getDefaultConfig(next)`。**驗證**：新增 `frontend/tests/workbench-v2/reducer-step1.test.ts`，9 個 action 各 1–2 個 case。
- [x] 1.5 在 reducer 抽出純函式 `cascadeRemoveDevice(state, deviceId)`，落地設計決策「刪除 device 連動清理 rules：在 reducer 層處理」並實作需求 **Multi-device tab management** 的 `Cascade clean-up on device removal` Scenario。**行為**：呼叫後同時清掉 devices/rules/points/mappings/db.targets 內所有與該 deviceId 相關項；若被刪 device 為當前 selected 則自動切到剩餘第一個。**驗證**：`reducer-step1.test.ts` 加 cascade 測試（GIVEN dev-A + 2 rules + 4 points + mappings + targets，刪除後全部清空）。

## 2. UI 子元件

- [x] 2.1 建立 `steps/step1/ProtocolSelector.tsx`，渲染 6 張 protocol 卡片（icon + 名稱 + 描述），選中時藍底 + 右上 check icon，落地需求 **Protocol switching with reset**。**行為**：點卡片觸發 `onChange(protocol)`，父元件呼叫 reducer `changeDeviceProtocol` 完成 test/config 重置。**驗證**：`step1.test.tsx` 中 `describe('ProtocolSelector')` 斷言渲染 6 張卡片、點擊 mqtt 卡片觸發 onChange、active 卡片有 highlight class。
- [x] 2.2 建立 `steps/step1/ConnectionConfigForm.tsx`，依 `protocol` prop 切換三種佈局（TCP 系 / RTU / MQTT），落地需求 **Connection config form per protocol**。**行為**：TCP 顯示 Host/Port/Slave ID/Timeout；RTU 顯示 Port/Baud Rate/Parity(N/E/O select)/Slave ID；MQTT 顯示 Broker/Username/Client ID；每欄 onChange 觸發 reducer `updateDeviceConfig`。**驗證**：`step1.test.tsx` 中 `describe('ConnectionConfigForm')` 三個 it 各別測 TCP/RTU/MQTT 欄位渲染。
- [x] 2.3 建立 `steps/step1/ReadinessStages.tsx`，接受 `stages: ReadinessStage[]` 與 `testStatus`，依 `group` 分成兩段（連線通道 / 協議握手），每段帶標題分隔線，落地設計決策「Readiness 群組：以 `connect` / `probe` 語義標籤」並實作需求 **Readiness check with explicit connect / probe separation**。**行為**：每個 stage 渲染圓圈（pending=灰 / success=emerald 勾 / running=旋轉）+ stage id + latency；分隔線顯示 group label。**驗證**：新增 `frontend/tests/workbench-v2/step1-readiness.test.tsx`，三個 it 各別測 modbus_tcp/modbus_rtu/mqtt 的群組分組正確。
- [x] 2.4 建立 `steps/step1/ConnectionTestPanel.tsx`，組合 send payload 預覽（`<pre>` mono JSON）+ `<ReadinessStages>` + 「執行測試」按鈕 + success readiness panel，落地需求 **Send payload preview block**。**行為**：payload 內容 reactive 反映 `device.protocol` 與 `device.config`；點按鈕觸發父元件 `onRunTest(device.id)`；success 後渲染綠色 readiness card 列出 Connection / Protocol Probe / 可進行下一步。**驗證**：`step1.test.tsx` 中 `describe('ConnectionTestPanel')` 測試 payload 反映、按鈕 disabled 條件、success card 渲染；新增 spy 斷言「payload 渲染但不發任何 fetch」。
- [x] 2.5 建立 `steps/step1/DeviceTabRail.tsx`，渲染水平捲動 tab 列 + 「+ 新增設備」按鈕，落地需求 **Multi-device tab management**。**行為**：每個 tab 顯示顏色點（用 `useDeviceColor`）、inline 改名 input、test status chip；hover 顯示刪除 icon（只在 `devices.length >= 2`）；點 tab 切 selected；點 + 新增。**驗證**：`step1.test.tsx` 中 `describe('DeviceTabRail')` 測試 add device、inline rename、刪除 confirm + cascade（用 mock `window.confirm`）。
- [x] 2.6 建立 `steps/step1/DeviceEditor.tsx`，組合 SectionCard（左 col 7）含設備名稱 + 描述 + `<ProtocolSelector>` + `<ConnectionConfigForm>`。**行為**：name/description 雙向綁定 reducer；切 protocol 自動透過 reducer 重置 test 與 config。**驗證**：`step1.test.tsx` 中 `describe('DeviceEditor')` 測 name input 改值反映 state、protocol 切換後 ConnectionConfigForm 換佈局。

## 3. 容器與整合

- [x] 3.1 建立 `steps/step1/Step1Device.tsx`，組裝 DeviceTabRail + DeviceEditor + ConnectionTestPanel + 底部 footer，並實作 mock 測試動畫 timer（每 380ms 推進）與卸載 cleanup，落地設計決策「Mock 測試動畫：使用 setTimeout + reducer action」並實作需求 **Readiness check with explicit connect / probe separation** 的 `Stage progression during test` Scenario。**行為**：`runDeviceTest(deviceId)` 啟動 chain；元件 unmount 時 clearTimeout 取消 pending tick；按鈕 disabled 期間不可重複觸發。**驗證**：`step1.test.tsx` 用 `vi.useFakeTimers()` 推 3 × 380ms 斷言 stage 依序變 success、最後 status=tested、latency 在 42–60ms；額外測 unmount cleanup（spy clearTimeout）。
- [x] 3.2 建立 `steps/step1/index.ts` barrel，re-export `Step1Device`、`DEVICE_COLORS`、`useDeviceColor`。**行為**：外部可從 `steps/step1` import 主元件與顏色 helper。**驗證**：`tsc --noEmit` 通過。
- [x] 3.3 修改 `shell/WorkbenchV2Shell.tsx` 在 `current === 1 && view === 'flow'` 時改 render `<Step1Device onContinue={...} />` 而非 `<Step1DevicePlaceholder />`，落地 modified 需求 **Placeholder step and settings surfaces**（Step 1 從 placeholder 名單移除）。**行為**：進入 step 1 看到完整 Step 1；切到 step 2 看到 Step2Placeholder（仍存在）；切到 settings 看到 SettingsPlaceholder。**驗證**：`shell.test.tsx` 既有測試需更新；新增 it 斷言 `current=1` 不再渲染「完整內容由後續 change 交付」訊息。
- [x] 3.4 在 `steps/Step1DevicePlaceholder.tsx` 加註解標示「Replaced by `steps/step1/Step1Device.tsx`；保留以供 rollback」，不刪除檔案。**行為**：檔案存在但未被 import。**驗證**：`grep -rn "Step1DevicePlaceholder" frontend/src` 只在註解或 `Step1DevicePlaceholder.tsx` 自身出現。

## 4. 底部 footer 與 continue gate

- [x] 4.1 在 `Step1Device.tsx` 內 render 底部 footer：顯示 `{N} 個設備，{M} 已通過測試`、warning chip（未全測時）、「全部建立並繼續」按鈕，落地需求 **Continue gate**。**行為**：所有 device test.status === 'success' 時按鈕啟用；點擊呼叫 `onContinue`，shell 將 1 加入 completed、current 切到 2。**驗證**：`step1.test.tsx` 中 `it('disables continue when not all tested')` 與 `it('enables and triggers onContinue when all tested')` 兩個 case。

## 5. i18n

- [x] 5.1 擴充 `frontend/src/i18n/locales/zh-TW/workbench-v2.json`，新增 `step1.*` keys：protocol 名稱（Modbus TCP / Modbus RTU / Modbus UDP / Fatek FBS / MC 3E / MQTT）+ 描述、stage 名稱（resolve / connect / probe / open_port / handshake / subscribe）、群組標題（連線通道 / 協議握手）、按鈕（執行測試 / 全部建立並繼續 / 新增設備）、readiness 訊息、警告 chip。**行為**：i18n key 完整覆蓋 step1 microcopy。**驗證**：`step1.test.tsx` 在 zh-TW 渲染斷言 `連線通道` 出現。
- [x] 5.2 擴充 `frontend/src/i18n/locales/en/workbench-v2.json` 對應英文翻譯（Connection Channel / Protocol Probe / Run Test / Create All and Continue / Add Device / …）。**行為**：對等翻譯。**驗證**：`step1.test.tsx` 在 en 渲染斷言 `Connection Channel` 出現。

## 6. 驗證

- [ ] 6.1 跑 `cd frontend && npm run lint`，零錯誤。
- [ ] 6.2 跑 `cd frontend && npm run typecheck`，零錯誤。
- [ ] 6.3 跑 `cd frontend && npm run test -- --run workbench-v2`，所有 `workbench-v2` 測試綠燈。
- [ ] 6.4 跑 `cd frontend && npm run build`，產出乾淨。
- [ ] 6.5 跑 `bash scripts/check_file_lines.sh`，確認 step1 所有新檔 ≤ 300 行。
- [ ] 6.6 手動 smoke：完整跑「預設 PLC-生產線-01 → 新增設備 2 (mqtt) → 切回設備 1 → 執行測試 → 切到設備 2 → 執行測試 → 點繼續到 Step 2 placeholder」；刪除設備時確認 confirm dialog 出現、刪除後 SummaryRail device 數量正確。

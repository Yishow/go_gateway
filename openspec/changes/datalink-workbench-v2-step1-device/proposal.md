## Why

`datalink-workbench-v2-shell` change 已落地 `/studio/v2` shell 與 4 個 step placeholder。本 change 是 phase 1 的第二棒，把 Step 1（多裝置設定）的完整靜態 UI 移植進 v2 容器，取代 `Step1DevicePlaceholder`。

Step 1 是整個 4 步驟流程的入口，必須做到：

- **多裝置同步管理**：原型 `step1-device.jsx` 已支援 tab 列管理多個 device、各自配置不同 protocol；現行 `/studio` 仍偏向單裝置心智模型，無法表達多生產線情境。
- **probe / connect 分離診斷**：依 `CLAUDE.md` 與 `AGENTS.md`，readiness 必須明確區分 connect（網路 / 序列埠 / TCP 通道）與 probe（協議層握手 / 訂閱），讓使用者知道是哪一階段失敗。
- **設備色彩語言**：Step 1 為每個 device 指派固定色（blue/emerald/amber/fuchsia/cyan/rose 循環），後續 Step 2 規則、Step 3 點位、Step 4 摘要欄都會延用同一個色票，建立跨 step 視覺一致性。

本 change 不接後端 API；測試行為以 mock 動畫（每階段 380ms 推進）撐起，留給 phase 2 backend-wiring 換成真實 `POST /devices/:id/test`。

## What Changes

- 新增 Step 1 完整 UI 於 `frontend/src/features/datalink/workbench-v2/steps/step1/`，三層版面：
  - **上層** Device tab 列（`SectionCard`）：水平捲動的 device card，每張可 inline 改名、刪除、顯示測試狀態 chip（`✓ 38ms` / `測試中…` / `未測試`），尾端「+ 新增設備」按鈕。
  - **中層** 12-col grid 編輯區：左 col 7（設備名稱 + 描述 + 6 個 protocol 卡片選擇器 + 協議專屬連線參數），右 col 5（連線測試 panel，含 send payload 預覽、分階段狀態、執行測試按鈕、success readiness panel）。
  - **下層** 摘要列：`{N} 個設備，{M} 已通過測試`，右側「全部建立並繼續」按鈕（disabled until 全部測試通過）。
- 支援 6 種協議：`modbus_tcp` / `modbus_rtu` / `modbus_udp` / `fatek_fbs` / `mc_3e` / `mqtt`。每種協議連線參數欄位不同（TCP 系：Host/Port/Slave ID/Timeout；RTU：Port/Baud Rate/Parity/Slave ID；MQTT：Broker/Username/Client ID）。
- **probe / connect 分離**：依協議產生不同 readiness stages：
  - TCP/UDP/Fatek/MC：`resolve` → `connect` → `probe`
  - RTU：`open_port` → `handshake` → `probe`
  - MQTT：`resolve` → `connect` → `subscribe`
  分兩個語義組別：`connect` 群組（網路 / 通道層）與 `probe` 群組（協議層握手 / 訂閱），UI 上以小標分隔，避免使用者把通道層失敗誤判為協議錯誤。
- **設備色彩語言**：在 `state/deviceColors.ts` 落地 6 色循環常數，並提供 `useDeviceColor(deviceId)` hook 供後續 step change 共用。
- **互動行為**：
  - 點 tab 切換 selected device，inline 改名直接 onChange。
  - 切換 protocol 自動清空 `test` 與 `status='draft'`，避免顯示過期測試結果。
  - 點刪除前 confirm（最後一個 device 不可刪）；刪除時連動清掉 `state.rules` 內屬於該 device 的規則。
  - 「執行測試」mock：每 380ms 推進一個 stage，最後 success 延遲取 42–60ms 隨機，並把 `status` 更新為 `tested`。
- 擴充 `useWorkbenchV2State` reducer actions：`addDevice` / `removeDevice` / `updateDevice` / `updateDeviceConfig` / `renameDevice` / `runDeviceTest`。
- 修改 `WorkbenchV2Shell` 在 `current === 1` 且 `view === 'flow'` 時改 render `Step1Device` 而非 `Step1DevicePlaceholder`。
- 擴充 i18n `workbench-v2` namespace 加入 Step 1 microcopy（protocol 卡片描述、stage 名稱、按鈕、readiness 訊息）。
- 不引入新依賴；沿用 shell change 的共用元件（`SectionCard`、`Field`、`Input`、`Select`、`Button`、`StatusChip`、`Toggle`、`Icon`）。

## Non-Goals (optional)

- 不接 `POST /api/v1/datalink/devices/test` 真實後端；測試流程仍為 mock 動畫。`send payload 預覽`區塊顯示的 JSON 是文字示意，不會實際 POST。
- 不在本 change 處理「儲存草稿」（頂列按鈕仍是 stub），留待 backend-wiring。
- 不導入 SSE 即時連線狀態；現有 `status: 'draft' | 'tested' | 'active'` 為純前端狀態，`active` 由後續 commit 流程設置。
- 不變更 Step 2/3/4/Settings 任何行為（這些仍是 placeholder）。
- 不涉及 mobile responsive；維持與 shell change 一致只支援 ≥ 1024px viewport。

## Capabilities

### New Capabilities

- `datalink-workbench-v2-step1-device`：定義 v2 Step 1 多裝置工作區的需求、6 種 protocol 配置形式、probe/connect 分離 readiness、設備色彩循環、tab 切換與連動清理規則。

### Modified Capabilities

- `datalink-workbench-v2-shell`：縮窄 `Placeholder step and settings surfaces` 需求範圍，把 Step 1 從 placeholder 名單移除（剩餘 Step 2/3/4 與 Settings 仍走 placeholder 直到對應 change 落地）。

## Impact

- Affected specs:
  - 新增 `openspec/specs/datalink-workbench-v2-step1-device/spec.md`
  - 修改 `openspec/specs/datalink-workbench-v2-shell/spec.md`（縮窄 placeholder 範圍）
- Affected code:
  - New:
    - `frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step1/index.ts`
    - `frontend/src/features/datalink/workbench-v2/state/deviceColors.ts`
    - `frontend/src/features/datalink/workbench-v2/state/protocols.ts`
    - `frontend/tests/workbench-v2/step1.test.tsx`
    - `frontend/tests/workbench-v2/step1-readiness.test.tsx`
  - Modified:
    - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts`（新增 device actions）
    - `frontend/src/features/datalink/workbench-v2/state/types.ts`（補完 `DeviceTest`、`ReadinessStage` 型別）
    - `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx`（current=1 切到 Step1Device）
    - `frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx`（移除或標記為 dead code 註解，由 Step1Device 取代）
    - `frontend/src/i18n/locales/zh-TW/workbench-v2.json`、`frontend/src/i18n/locales/en/workbench-v2.json`
- 不變更後端任何檔案。

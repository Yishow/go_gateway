## Context

`datalink-workbench-v2-shell` 已建立 `/studio/v2` 容器、`useWorkbenchV2State` reducer 與共用元件。本 change 在這個容器內把原型 `step1-device.jsx`（357 行）移植成 TypeScript，且必須滿足 `CLAUDE.md` 的「Step 1 要明確區分 connect / probe 診斷」硬性要求。

原型 Step 1 涵蓋了多裝置 tab、6 種協議切換、3 階段 readiness、刪除連動 rules 清理等多項互動，總行數 357 行；TypeScript 拆檔後預計 7 個檔案，每檔 100–250 行區間以符合 repo 300/500 行硬上限。

## Goals / Non-Goals

**Goals:**

- 1:1 視覺與互動復刻原型 Step 1，使用 shell change 提供的共用元件。
- 在 readiness 流程中以結構化方式區分 `connect` 與 `probe` 階段群組。
- 設備色彩 6 色循環，對後續 step change 暴露 `useDeviceColor` hook 與常數，避免色彩規則散落。
- 補完 state types 與 reducer actions，使後續 Step 2 可直接從 state 讀 devices 並衍生 rule.device_id。

**Non-Goals:**

- 不接 `POST /devices/test` 真實 API。
- 不處理 SSE 即時狀態。
- 不變更舊 `/studio` 的 Step 1（既有 `device-registry` spec 不動）。
- 不處理裝置列表分頁、搜尋、過濾（單一閘道情境下 device 數量 ≤ 10）。

## Decisions

### 拆檔策略：以「資料區塊 / 互動行為」拆分而非單一巨檔

原型把所有邏輯塞在一個 357 行的 `Step1Device` 函式內；移植時拆為 7 個檔案：

- `Step1Device.tsx`（容器，~200 行）：state derivation、業務邏輯、底部 footer、子元件組合。
- `DeviceTabRail.tsx`（~120 行）：上層 tab 列。
- `DeviceEditor.tsx`（~80 行）：中層左側 SectionCard 容器。
- `ProtocolSelector.tsx`（~80 行）：6 個 protocol 卡片選擇器。
- `ConnectionConfigForm.tsx`（~150 行）：協議專屬連線參數三種佈局（TCP 系 / RTU / MQTT）。
- `ConnectionTestPanel.tsx`（~150 行）：中層右側 SectionCard，含 send payload 區、stage list、執行按鈕、success readiness。
- `ReadinessStages.tsx`（~100 行）：分階段狀態列表元件，內含 connect/probe 群組分隔。

**Alternatives considered**：(A) 把所有 UI 塞在一個 `Step1Device.tsx`（複製原型形狀）— 超過 300 行警戒線，違反 repo 規則；(B) 過度拆分成 10+ 檔 — 增加跨檔追蹤成本，且每個檔案不到 50 行 props 多於邏輯。

### Readiness 群組：以 `connect` / `probe` 語義標籤

`ReadinessStage` TypeScript 型別新增 `group: 'connect' | 'probe'` 欄位：

```ts
type ReadinessStage = {
  id: string;        // 'resolve' | 'tcp' | 'modbus_probe' | 'open_port' | ...
  label: string;    // i18n label
  group: 'connect' | 'probe';
  status: 'pending' | 'running' | 'success' | 'failed';
  latency_ms?: number;
  message?: string;
};
```

依協議產生 stage 序列的 lookup table 集中在 `state/protocols.ts`：

| Protocol     | Stage 1            | Stage 2            | Stage 3              |
| ------------ | ------------------ | ------------------ | -------------------- |
| modbus_tcp   | resolve (connect)  | connect (connect)  | probe (probe)        |
| modbus_udp   | resolve (connect)  | connect (connect)  | probe (probe)        |
| fatek_fbs    | resolve (connect)  | connect (connect)  | probe (probe)        |
| mc_3e        | resolve (connect)  | connect (connect)  | probe (probe)        |
| modbus_rtu   | open_port (connect)| handshake (connect)| probe (probe)        |
| mqtt         | resolve (connect)  | connect (connect)  | subscribe (probe)    |

UI 上分兩段呈現：「連線通道（connect）」分隔線 + 兩個 stage，「協議握手（probe）」分隔線 + 一個 stage。每段都有合計延遲顯示。

**Alternatives considered**：(A) 不分群組，照原型平鋪三 stage — 違反 `CLAUDE.md`「明確區分」要求；(B) 三段（resolve / connect / probe）— resolve 在語義上仍屬於 connect 通道層，分太細沒意義。

### 設備色彩：6 色循環常數 + hook 共用

`state/deviceColors.ts` 匯出：

```ts
export const DEVICE_COLORS = ['blue', 'emerald', 'amber', 'fuchsia', 'cyan', 'rose'] as const;
export type DeviceColor = typeof DEVICE_COLORS[number];

export function useDeviceColor(deviceId: string): {
  name: DeviceColor;
  bg: string;     // tailwind class fragment, e.g. 'bg-blue-500/15'
  border: string;
  text: string;
  solid: string;
};
```

實作 inside hook：用 `useWorkbenchV2State` 拿 `devices`，`devices.findIndex(d => d.id === deviceId) % 6` 對應到 `DEVICE_COLORS`。Step 2/3/4 都會 import 此 hook 取得統一顏色，避免每個 step 自己維護一份 color cycling。

**Alternatives considered**：(A) 用 hash function（如 `hashCode(deviceId) % 6`）對應顏色 — 加裝置時不會 reshuffle 既有 device 顏色，但失去原型「依新增順序循環」的可預期性；(B) 把顏色直接寫進 `Device` 型別 — state 多一個欄位，且改名/排序時要重算。

### Mock 測試動畫：使用 setTimeout + reducer action

`runDeviceTest(deviceId)` action 由 component 在 useEffect 內啟動：

```ts
function runTest(deviceId: string) {
  dispatch({ type: 'startDeviceTest', deviceId });
  const stages = getStagesForProtocol(device.protocol);
  let idx = 0;
  const tick = () => {
    if (idx >= stages.length) {
      dispatch({ type: 'completeDeviceTest', deviceId, latency: 42 + rand(18) });
      return;
    }
    dispatch({
      type: 'advanceDeviceTest',
      deviceId,
      stageId: stages[idx].id,
      stageLatency: 10 + rand(18),
    });
    idx++;
    setTimeout(tick, 380);
  };
  setTimeout(tick, 220);
}
```

把 timer 管理放在 component，reducer 只處理純函式 state 更新。卸載元件時 component cleanup 應 cancel 進行中的 timer。

**Alternatives considered**：(A) 在 reducer 用 thunk middleware — `useReducer` 不支援 middleware，會引入 reduce-thunk 依賴；(B) 用 React Query mutation — 還沒接後端，此階段引入 React Query 違反「phase 2 才接 API」邊界。

### 刪除 device 連動清理 rules：在 reducer 層處理

`removeDevice(deviceId)` reducer action 同時：

1. 過濾 `state.devices` 移除該 device。
2. 過濾 `state.rules` 移除所有 `device_id === deviceId` 的規則。
3. 過濾 `state.points` 移除衍生自上述 rules 的 points。
4. 過濾 `state.mappings` 移除對應 point 的 mapping。
5. 過濾 `state.db.targets` 移除對應 point 的 target。
6. 若被刪 device 是當前 selected → 把 selected 切到第一個剩下的 device。

把這個連動清理放在 reducer 內，確保不會有「孤兒 rule」或「指向已刪 device 的 point」狀態，比每個 step 自己驗證更安全。

**Alternatives considered**：(A) component 內呼叫多個 dispatch — race condition 風險；(B) 不連動，留給後續 step 自己處理孤兒 — 違反「設備刪除應同時移除所屬規則」的明確需求。

### Protocol 切換：清空 test 結果與 config

切換 `protocol` 時：

- 把 `test` 設為 `null`，`status` 改回 `'draft'`。
- 把 `config` 替換為對應協議的預設值（TCP: `{ host: '192.168.1.100', port: 502, slave_id: 1, timeout: 5 }`，RTU: `{ port: '/dev/ttyUSB0', baud: 9600, parity: 'N', slave_id: 1 }`，MQTT: `{ broker: 'mqtts://broker.local:8883', username: '', client_id: 'gw-01' }`）。

避免使用者切換協議後看到上一個協議的 host/port 殘留而誤判。

## Implementation Contract

#### Behavior

- 進入 `/studio/v2`（current=1）：顯示一張預設 device card（`PLC-生產線-01` / `modbus_tcp` / `192.168.1.100:502`）+ 「+ 新增設備」按鈕。
- 點「+ 新增設備」：tab 列尾端新增 `設備 2`，自動 selected，host 自動帶 `192.168.1.101`。
- 點 device card：成為 selected，中層編輯區換成該 device 內容。
- 在 tab 內 inline 改名：`devices[].name` 即時更新，tab 文字同步、breadcrumb 不變（breadcrumb 顯示 step 名而非 device 名）。
- 切換 protocol（點下方 6 張 protocol 卡片其一）：選中卡片右上出現 check icon、藍底 highlight；連線參數區塊布局換成該協議專屬欄位；該 device `test` 清空、status 回 `'draft'`。
- 點「執行測試」：按鈕 disabled、icon 換 refresh 旋轉；每 380ms 推進一個 stage、stage 圓圈變 emerald 並顯示 latency；全部完成後 success readiness panel 出現，列出「Connection Configuration ✓」「Protocol Probe ✓ ({latency}ms)」「可進行下一步：設定接入規則」。
- 全部 device 都 `test.status === 'success'` 時，底部「全部建立並繼續」按鈕啟用；點擊呼叫 `onContinue`，shell `completeStep(1)` + `setCurrent(2)`。
- 在 ≥ 2 個 device 的情況下，hover device card 出現 close icon；點擊跳 confirm dialog（瀏覽器原生 `window.confirm`）；確認後該 device 與所屬規則同步移除。
- 切到 Step 1 後再切回 Step 1（透過 StepRail）：先前測試過的 device 仍保留 `test.status === 'success'`，不重新驗證。

#### Interface / Data Shape

- `state.devices[]` 元素型別：
  ```ts
  type Device = {
    id: string;                 // dev-xxx
    name: string;
    description: string;
    protocol: ProtocolId;
    config: Record<string, string | number>;
    status: 'draft' | 'tested' | 'active';
    test: DeviceTest | null;
  };

  type DeviceTest = {
    status: 'running' | 'success' | 'failed';
    latency_ms?: number;
    stages: Record<string, { status: 'pending' | 'running' | 'success' | 'failed'; latency_ms?: number }>;
    tested_at?: string; // ISO timestamp
  };

  type ProtocolId = 'modbus_tcp' | 'modbus_rtu' | 'modbus_udp' | 'fatek_fbs' | 'mc_3e' | 'mqtt';
  ```
- Reducer actions：
  - `{ type: 'addDevice'; device: Device }`
  - `{ type: 'removeDevice'; deviceId: string }`
  - `{ type: 'updateDevice'; deviceId: string; patch: Partial<Device> }`
  - `{ type: 'updateDeviceConfig'; deviceId: string; patch: Partial<Device['config']> }`
  - `{ type: 'renameDevice'; deviceId: string; name: string }`
  - `{ type: 'changeDeviceProtocol'; deviceId: string; protocol: ProtocolId }`
  - `{ type: 'startDeviceTest'; deviceId: string }`
  - `{ type: 'advanceDeviceTest'; deviceId: string; stageId: string; stageLatency: number }`
  - `{ type: 'completeDeviceTest'; deviceId: string; totalLatency: number }`
  - `{ type: 'failDeviceTest'; deviceId: string; stageId: string; message: string }`
- `protocols.ts` exports：
  - `PROTOCOLS`: array of `{ id: ProtocolId; nameKey: string; descKey: string }`
  - `getStagesForProtocol(protocol: ProtocolId): ReadinessStage[]`
  - `getDefaultConfig(protocol: ProtocolId): Device['config']`

#### Failure Modes

- Mock 測試完成前 unmount 元件：useEffect cleanup 取消 pending setTimeout，state 留在最後一個 reducer 已套用的階段；下次回到該 step 時 device card 顯示「測試中…」可手動再按一次「執行測試」（會重置 `test.stages` 重跑）。
- 使用者連續點「執行測試」：按鈕 disabled 期間不可重複觸發；若某種原因觸發了第二次，第二次的 `startDeviceTest` reducer 會清空既有 `stages` 重新開始。
- Confirm dialog 取消：state 完全不變，device 與 rules 都保留。
- Protocol 切換後 config 欄位驗證失敗（例如 port 留空）：UI 顯示紅框錯誤訊息，「執行測試」按鈕保持 disabled 直到欄位填妥。

#### Acceptance Criteria

- `frontend/tests/workbench-v2/step1.test.tsx` 覆蓋：
  - render 預設一張 device card。
  - 點「新增設備」後 devices 變兩筆且第二筆 selected。
  - 在 tab 內輸入新名稱：state.devices[i].name 更新。
  - 切到 `mqtt` 協議：右側連線參數欄位顯示 Broker / Username / Client ID；`test` 與 `status` 重置。
  - 「執行測試」mock：use fake timers 推進 3 × 380ms，斷言三個 stage 依序變 success，最後 device.status === `tested`、test.latency_ms 在 42–60 之間。
  - 全部 device 通過測試後「全部建立並繼續」啟用，點擊 onContinue 被呼叫且 shell completed set 含 1。
  - 刪除 device 同步移除 state.rules 內 `device_id === deletedDeviceId` 的規則。
- `frontend/tests/workbench-v2/step1-readiness.test.tsx` 覆蓋：
  - 每個 protocol 的 stage 序列順序與 group 標籤正確（用 `getStagesForProtocol` 直接斷言）。
  - ReadinessStages 元件渲染時，stage 依 `group` 分成兩段，且每段都有標題。
  - mqtt：connect 群組含 resolve + connect，probe 群組含 subscribe。
  - modbus_rtu：connect 群組含 open_port + handshake，probe 群組含 probe。
- `cd frontend && npm run lint && npm run typecheck && npm run test -- --run workbench-v2 && npm run build` 全綠。
- 手動 smoke：在 `/studio/v2` 完整跑一輪「新增第二個 device → 切到 mqtt → 執行測試 → 切回第一個 device → 全部測試完 → 點繼續到 Step 2 placeholder」。

#### Scope Boundaries

**In scope:**

- Step 1 完整靜態 UI（7 個檔案）。
- `useWorkbenchV2State` 新增 9 個 device actions。
- `state/types.ts` 補完 `Device`、`DeviceTest`、`ReadinessStage` 型別。
- `state/deviceColors.ts`、`state/protocols.ts`。
- shell 換掉 Step1DevicePlaceholder，shell spec 縮窄 placeholder 範圍。
- i18n `workbench-v2` 加入 Step 1 microcopy。
- 上述項目的 unit test。

**Out of scope:**

- Step 2/3/4 與 Settings 行為（仍是 placeholder）。
- 任何後端 API 呼叫（mock 動畫）。
- E2E 測試。
- 「儲存草稿」按鈕的真實行為。
- 失敗測試的視覺對應（只實作 success path；fail tone 留 type-level 但 UI 不展示）。

## Risks / Trade-offs

- [Mock 動畫 latency 與真實後端差距] → 後續 backend-wiring change 會把 `runDeviceTest` 改成真實 API call，stage progress 改由 SSE 推送；保持 reducer actions 形狀不變，只換 component handler。
- [刪除連動清理涉及多個 state slice，reducer 變胖] → reducer 寫成 `case 'removeDevice': return cascadeRemoveDevice(state, action.deviceId);`，把連動邏輯抽到純函式 `cascadeRemoveDevice` 並單元測試。
- [`useDeviceColor` 依賴 device 順序，新增 device 時顏色循環不會打亂既有顏色] → 預期行為；若未來改為 hash 對應再改 hook 實作即可，外部 API 不變。
- [切換 protocol 清空 config 可能讓使用者已輸入的 host/port 一鍵清空] → 預期行為，避免 host:port 跟新協議不相容；同時在 UI 上加 hint「切換協議會重置連線參數」。

## Migration Plan

1. 在 `state/types.ts` 補完 `Device`、`DeviceTest`、`ProtocolId`、`ReadinessStage` 型別。
2. 在 `state/protocols.ts` 落地 `PROTOCOLS`、`getStagesForProtocol`、`getDefaultConfig`，先寫 protocols.test.ts。
3. 在 `state/deviceColors.ts` 落地 `DEVICE_COLORS`、`useDeviceColor`，先寫 colors.test.ts。
4. 擴充 `useWorkbenchV2State` reducer，TDD：先測 addDevice / removeDevice cascade / startDeviceTest / advanceDeviceTest / completeDeviceTest。
5. 從子元件實作起：`ProtocolSelector` → `ConnectionConfigForm` → `ReadinessStages` → `ConnectionTestPanel` → `DeviceEditor` → `DeviceTabRail`。
6. 組裝 `Step1Device.tsx`，覆寫 shell 內 `current === 1` 的渲染分支。
7. 補 i18n 翻譯 key。
8. 跑 lint / test / build / check-lines。
9. 手動 smoke 完整流程。
10. PR review 重點：reducer cascade 邏輯、useEffect cleanup、ReadinessStages 群組分隔。

**Rollback strategy**：本 change 不刪除任何 shell change 的檔案；rollback 只需把 `WorkbenchV2Shell` 內 `current === 1` 的渲染改回 `Step1DevicePlaceholder`、移除 step1 目錄與相關 state 變更。reducer 新增的 actions 與 placeholder 並存，rollback 後 actions 變成 dead code 但不影響運作。

## Open Questions

- 失敗測試（mock 隨機 fail）是否在本 change 落地？傾向 **不**，原型也沒展示 fail 視覺，UI 在 fail 時行為留給 backend-wiring change 處理；本 change reducer `failDeviceTest` action 保留但 UI 不會自然觸發。
- 「描述」欄位是否必填？傾向 **否**，與原型一致；UI 上以「選填說明」hint 標示。

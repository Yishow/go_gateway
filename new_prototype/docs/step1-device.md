# Step 1 · 新增裝置

## 目的

讓使用者一次新增「一個或多個」工業設備，逐一設定協議與連線參數，並各自完成 readiness 連線測試。所有設備測試通過才能進入下一步。

對應後端：`POST /devices` → `POST /devices/:id/test` → `POST /devices/:id/activate`

---

## 完整 Prompt (可直接生成此頁)

> 為 go_gateway 系統做一個「新增裝置」頁，支援多設備。版面分三層：上方裝置 tab 列、中間選中裝置的編輯器、底部摘要+繼續按鈕。
>
> **上層 - 裝置 tab 列 (SectionCard)**
> - 標題「設備列表」，副標「可同時設定多個設備，每個設備有各自的協議與連線參數」
> - icon 用 `device` 圖示
> - 右側 aside chip：「N 個設備 · M 已測試」
> - 內容是水平捲動的 tab 列，每個 tab 是 `min-w-[220px]` 的卡片：
>   - 左上：對應顏色點 (從 blue/emerald/amber/fuchsia/cyan/rose 循環) + inline 可編輯的裝置名稱
>   - 右上：hover 才顯示的刪除 X 鈕（≥2 個時才可刪）
>   - 下方：`{protocol} · {host}` + 測試狀態 chip (`✓ 38ms` / `測試中…` / `未測試`)
>   - 選中時邊框與背景使用該設備顏色
> - 最右邊：虛線邊框的「+ 新增設備」按鈕；按下產生新設備，配 `192.168.1.10N` 的 host，並自動切到該設備
>
> **中層 - 選中設備的編輯器 (12-col grid)**
>
> 左 (col 7) `SectionCard "編輯：<name>"`：
> - 兩欄表單：「設備名稱」(required) + 「描述」(選填)
> - 「通訊協議」：6 個卡片式選項 (2-col on mobile / 3-col on md+) — Modbus TCP/RTU/UDP、Fatek FBS、MC 3E、MQTT；每張顯示名稱 + 一句說明；選中時藍底，右上有勾選圖示
> - 連線參數區塊 (rounded-xl, 內嵌 mono 字體)：
>   - **TCP/UDP/Fatek/MC 系列**：Host (mono) / Port / Slave ID / Timeout(s)
>   - **modbus_rtu**：Port (/dev/ttyUSB0) / Baud Rate / Parity (N/E/O) / Slave ID
>   - **mqtt**：Broker / Username / Client ID
> - 切換協議時清空 test 並重置為 draft 狀態
>
> 右 (col 5) `SectionCard "連線測試"`：
> - aside chip：未測試 (draft) / 測試中… (info) / 已通過 (success)
> - 內容：
>   - **送出 payload 預覽**：`POST /api/v1/datalink/devices/test` + 真實 JSON payload (反映目前表單)，用 mono font
>   - **分階段狀態**：依協議產生不同階段名 — TCP: `resolve` `tcp` `modbus_probe`；RTU: `open_port` `handshake` `probe`；MQTT: `resolve` `connect` `subscribe`
>     - 未跑：灰底圓圈帶階段編號
>     - 已通過：emerald 圓圈帶勾，右側顯示 `{N}ms`
>     - 跑到該階段時的等候態：「等候…」
>   - 「執行測試」按鈕 (secondary)：點擊後逐一推進每階段 (380ms 間隔)，全部成功後 totals 在 42-60ms 區間隨機
>   - 成功後顯示 emerald readiness panel：「Readiness 檢查通過 · Connection Configuration ✓ · Protocol Probe ✓ · 可進行下一步」
>
> **底部 - 摘要與繼續**
> - `<div className="rounded-xl border bg-slate-900/40 px-4 py-3 flex items-center justify-between">`
> - 左：「{N} 個設備，{M} 已通過測試」+ 必要時 warning chip「尚有 X 個設備未通過測試」
> - 右：`<Button>` 「全部建立並繼續」(disabled 直到全部測試通過)
>
> **互動細節**
> - 點 tab 切到該設備；點 + 新增；點 ✕ 刪除（會 confirm，同時移除該設備所屬的 rules）
> - 設備名稱可直接 inline 編輯（不需進編輯器）
> - 切換協議自動清掉測試結果
> - 全部設備必須完成測試才能繼續
>
> **配色**
> - 裝置顏色循環：blue / emerald / amber / fuchsia / cyan / rose
> - 同一個設備在後續所有 step（Step 2 rule tab、合併點位表、SummaryRail）都用同一顏色
>
> 用 React 18 + Tailwind，dark slate 主題；mono 字體用 JetBrains Mono；所有狀態變化用 `transition-all duration-200`；新增設備時用 sweep-in 動畫 (240ms ease-out)。

---

## State 介面

```ts
Props: {
  state: { devices: Device[], rules: Rule[], ... };
  setState: (updater: (prev: State) => State) => void;
  onContinue: () => void;
}

Local state: {
  selectedDeviceId: string | null;
  testing: boolean;  // 是否正在跑測試
}
```

---

## 主要函式

```js
const addDevice = () => {
  const newDev = makeDefaultDevice({ name: `設備 ${devices.length + 1}`, ... });
  setState(s => ({ ...s, devices: [...s.devices, newDev] }));
  setSelectedDeviceId(newDev.id);
};

const removeDevice = (id) => {
  if (!confirm('刪除此設備將同時移除所屬規則...')) return;
  setState(s => ({
    ...s,
    devices: s.devices.filter(d => d.id !== id),
    rules: s.rules.filter(r => r.device_id !== id),  // 連動清理
  }));
};

const runTest = () => {
  // 依 protocol 決定階段名稱
  // 每 380ms 推進一個階段，更新 updateDevice({ test: {...} })
  // 全部完成後 status = 'tested', test.status = 'success'
};
```

---

## 驗收標準

- [ ] 初始載入顯示 1 個預設設備 (PLC-生產線-01, Modbus TCP)
- [ ] 點「+ 新增設備」立即新增第 2 個設備，自動切過去
- [ ] 切換協議時，連線參數區塊根據協議顯示不同欄位
- [ ] 「執行測試」按鈕按下後，3 個階段依序變綠 (約 1.5 秒完成)
- [ ] 任一設備未測試時，底部「全部建立並繼續」disabled
- [ ] 刪除設備時，已關聯的 rules 同步被清掉
- [ ] tab 列在窄螢幕可水平捲動，不會擠壓內容

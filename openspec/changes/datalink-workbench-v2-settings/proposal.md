## Why

`datalink-workbench-v2-shell` 與 step1/2/3/4 五個 change 已落地 shell 與 4 步驟主流程。本 change 是 phase 1 的最後一棒，把獨立的「設定」頁完整靜態 UI 移植進 v2 容器，取代 `SettingsPlaceholder`。

設定頁是 4 步驟流程之外的系統級管理頁，涵蓋 5 個彼此獨立但相互關聯的區塊：

- **資料庫連接器池**：可註冊多個目標資料庫供 Step 4 挑選；每個 connector 含類型、host、port、database、username、預設寫入間隔、啟用狀態、最近一次測試結果與時間戳。
- **時序儲存策略**：時間精度（秒 / 毫秒）、分區間隔（每日 / 每週 / 每月）、批次寫入大小、保留天數。
- **排程預設**：預設輪詢間隔、重試次數、重試延遲、斷路器門檻、開機自動啟動 collector。
- **Local Modbus Share**：總開關 + 綁定位址 + Port + Slave ID + 起始 Register；Step 2 的 share register 計算會讀這裡。
- **介面 / API / 診斷**：主題 / 語言 / 位址顯示格式、Base URL / API 版本 / 請求逾時、日誌等級 / SSE 心跳 / Debug 面板 / 操作審計記錄。

這 5 個區塊與主 4 步驟流程連動的關鍵：

- `settings.modbus_share.enabled` 影響 Step 2 share 位址計算（已在 step2 change 對接）。
- `settings.modbus_share.base_register` 是 Step 2 自動配置 share register 起點。
- `settings.scheduler.default_interval_ms` 是新建立 polling group 的預設值（commit 時送出）。
- `settings.connectors` 是 Step 4 的 connector 來源（後續 phase 2 會把 Step 4 內建單一 connector 改為從這裡選）。

本 change 不接後端；所有設定儲存於 `state.settings`。Connector 測試是 mock 動畫（900ms 隨機 success/fail）。底部 sticky 儲存列的「儲存所有設定」按鈕在本 change 仍為 noop（後續 backend-wiring 才會接到 `PATCH /settings` 之類的端點）。

## What Changes

- 新增 Settings 完整 UI 於 `frontend/src/features/datalink/workbench-v2/settings/`，整頁長卷軸由 5 個 SectionCard 區塊 + 1 個標題列 + 1 個 sticky 底部儲存列組成：
  - **標題列**（漸層卡）：11×11 圓角藍色 sliders icon + 「系統設定」標題 + 一行說明。
  - **資料庫連接器池**：標題 + aside「+ 新增連接器」按鈕；空狀態 hint；列表 divide-y，每行：類型 emoji icon + 12-col grid 表單（連線名稱 / 類型 / Host / Port / Database / 使用者 / 預設寫入間隔 / 啟用 toggle）+ 跨欄狀態列（status chip + last_check_at 時間戳 + 測試 / 刪除 button）。
  - **時序儲存策略 + 排程預設**（並排 col 6+6）：兩個 SectionCard。
  - **Local Modbus Share**（總開關在 aside）：啟用時顯示 4-col grid（綁定位址 / Port / Slave ID / 起始 Register）；停用時僅顯示說明文字。
  - **介面 / API / 診斷**（並排 col 4+4+4）：三個 SectionCard，各含 3–4 個 Field + Select + Toggle。
  - **底部 sticky 儲存列**：rounded-xl backdrop-blur + 左 info「設定會立即套用，並於下次重啟後生效」+ 右「重設為預設」ghost button + 「儲存所有設定」success button。
- **Connector 測試**：mock 900ms 後 85% 成功、15% 失敗；失敗時 `last_check_error: 'connection refused'`、status `unreachable`。
- **連動效應**：本 change 落地 settings 寫入；Step 2 share 計算與 Step 4 connector 既有實作會即時讀到變更（已由前面 change 的 selector / useEffect 處理）。
- 擴充 `useWorkbenchV2State` reducer 加 settings actions（`updateSettings`、`updateSettingsSection`、`addConnector`、`updateConnector`、`removeConnector`、`startConnectorTest`、`completeConnectorTest`、`resetSettingsToDefaults`）。
- 修改 `WorkbenchV2Shell` 在 `view === 'settings'` 時改 render `SettingsPage`；shell spec MODIFIED placeholder 完全清空（不再有 placeholder）。
- 擴充 i18n 加入 Settings microcopy（5 個區塊的標題 / 副標 / Field labels / hints / Select options / Toggle labels / 按鈕 / 空狀態文案 / 連接器測試狀態文字）。
- 不引入新依賴。

## Non-Goals (optional)

- 不接 `PATCH /settings` / `POST /db-connectors/test` 真實後端；儲存列的「儲存所有設定」按鈕為 noop（後續 backend-wiring 才接）。
- 不變更 `database-target-workbench` / `local-modbus-memory-workbench` 等既有 spec。
- 不處理設定匯入 / 匯出（JSON 上下載）。
- 不處理「主題切換」實作 light mode（v2 仍 dark-only；主題 select 只是存設定但 UI 不換）。
- 不處理「語言切換」即時 reload；切了之後須重新整理頁面（後續 backend-wiring 會處理 useTranslation 動態切換）。
- 不變更 Step 1/2/3/4 行為。

## Capabilities

### New Capabilities

- `datalink-workbench-v2-settings`：定義 v2 Settings 頁的 5 個區塊需求、connector 池 CRUD 與測試契約、時序儲存策略欄位、排程預設欄位、Local Modbus Share 配置、介面 / API / 診斷選項、sticky 底部儲存列行為、與主 4 步驟流程的連動效應規範。

### Modified Capabilities

- `datalink-workbench-v2-shell`：縮窄 `Placeholder step and settings surfaces` 需求，把 Settings 從 placeholder 名單移除（至此 phase 1 完成，無 placeholder 殘留）。

## Impact

- Affected specs:
  - 新增 `openspec/specs/datalink-workbench-v2-settings/spec.md`
  - 修改 `openspec/specs/datalink-workbench-v2-shell/spec.md`
- Affected code:
  - New:
    - `frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/index.ts`
    - `frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts`
    - `frontend/tests/workbench-v2/settings.test.tsx`
    - `frontend/tests/workbench-v2/settings-connectors.test.tsx`
  - Modified:
    - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts`（新增 settings actions）
    - `frontend/src/features/datalink/workbench-v2/state/types.ts`（補完 `Settings`、`SettingsConnector`、`TimeseriesSettings`、`SchedulerSettings`、`ModbusShareSettings`、`GeneralSettings` 型別）
    - `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx`（view='settings' 切到 SettingsPage）
    - `frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx`（標 rollback only）
    - `frontend/src/i18n/locales/zh-TW/workbench-v2.json`、`frontend/src/i18n/locales/en/workbench-v2.json`
- 不變更後端任何檔案。

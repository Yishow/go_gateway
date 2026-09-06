## Why

`fix-studio-v2-setup-and-db-flow` 已完成實作但尚未提交，程式碼審查在該差異上找出 12 項缺陷，其中三類會直接讓使用者無法完成設定流程或造成資安風險：

1. **Step 4「啟用」永久卡死**：Modbus Share 設定欄位一經編輯就把 `save_state` 標成 `saving` 卻無人清除，而 autosave barrier 的等待承諾沒有逾時也不會被拒絕，使用者按下啟用後會停在 `activating`，既無錯誤訊息也無法重試或重置。
2. **MySQL 資料庫目標新支援路徑不可用**：schema 產生器對 MySQL 送出非法索引語法、預設 schema 名稱與連線資料庫名稱不一致，導致 upsert 映射的 schema 產生直接失敗，或退化成每個時間桶重複寫入的 append-only 行為。
3. **MySQL 密碼可能以明文上線路**：cleartext 認證預設開啟並同時允許退回未加密連線，目標資料庫密碼會在伺服器未提供 TLS 時無聲地以明文傳送。

其餘缺陷（設定存檔錯誤被靜默吞掉、已刪除連接器復活、啟用前置檢查被整段略過、死碼與錯誤分類錯置）雖不致當機，但都會讓使用者看到與實際狀態不符的畫面或錯誤原因。這些缺陷全部位於尚未提交的差異中，現在修復的成本遠低於進入主線之後。

## What Changes

**Studio V2 設定與啟用流程（前端）**

- Modbus Share 欄位編輯與重置不再擅自寫入 `save_state: 'saving'`；`saving` 僅由真正發出的存檔請求設定，並由該請求的成功或失敗路徑收斂為 `saved` 或 `save-error`。
- autosave settlement 等待加入逾時上限；逾時的等待者會以可辨識的 barrier 結果收斂（而非永久擱置），讓 Step 4 能顯示錯誤並提供重試。
- 存檔進行中若使用者再編輯任一設定欄位，存檔的失敗結果不再被靜默吞掉：`save_state` 與錯誤訊息必須收斂到終態，使用者看得到失敗並可重試。
- 新增連接器的成功清單改以「已成功建立且尚未刪除」為準，已刪除的連接器不會在後續新增時重新出現在列表。
- Step 1 連線表單不再顯示未寫入設定的站號預設值；站號預設值改在裝置載入正規化階段寫入設定，讓畫面顯示與 `isStudioV2DeviceValid` 判定一致。

**MySQL 資料庫目標（後端）**

- schema 產生器針對 MySQL 輸出合法的唯一索引語法（MySQL 不支援索引的 `IF NOT EXISTS`），SQLite 與 PostgreSQL 維持既有語法。
- schema 產生器改以連接器設定解析預設 schema 名稱，與映射建立/更新時採用的預設值一致，不再對 MySQL 產生指向 `main` 的語句。
- upsert 映射在 MySQL 上的唯一鍵保證納入 schema 產生的驗收條件，避免無聲退化為重複寫入。
- **BREAKING**：MySQL 連接器的 `allow_cleartext_passwords` 預設值由 true 改為 false；明確啟用 cleartext 認證時不再允許退回未加密連線，連線失敗會回報可辨識的錯誤，而非以明文送出密碼。
- 移除連線探測中重複執行資料庫自動建立的死碼分支。

**Modbus Share 啟用前置檢查（後端）**

- Share 停用時仍保留工作區層級的前置檢查（readiness token 與 workspace revision），僅略過 Share 專屬檢查；Share 從未完成 hydration 時維持「不阻擋工作區啟用」的既有行為。
- 讀取 Share 設定失敗時回報內部讀取失敗，不再誤報為「Share 已停用」。

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-settings`: 設定頁存檔狀態機必須收斂到終態；編輯欄位不得偽造存檔中狀態；連接器列表不得復活已刪除項目。
- `modbus-share-lifecycle`: 啟用 autosave barrier 必須有界（逾時可收斂）；Share 停用時的啟用前置檢查範圍與設定讀取失敗的錯誤分類需明確界定。
- `database-output-delivery`: MySQL 目標的 schema 產生必須輸出方言合法且 schema 名稱正確的語句，並保證 upsert 所需唯一鍵；MySQL 連線的密碼傳輸必須失敗關閉而非明文退回。
- `studio-v2-device-validity`: 裝置設定表單顯示的必填值必須與裝置有效性判定使用同一份設定資料。

## Impact

- Affected specs:
  - `datalink-workbench-v2-settings`
  - `modbus-share-lifecycle`
  - `database-output-delivery`
  - `studio-v2-device-validity`
- Affected code:
  - Modified:
    - `frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts`
    - `frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts`
    - `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx`
    - `frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts`
    - `frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx`
    - `frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx`
    - `frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx`
    - `internal/api/router_modbus_share.go`
    - `internal/api/router_modbus_share_test.go`
    - `internal/datalink/dbtarget/tooling_service.go`
    - `internal/datalink/dbtarget/tooling_service_helpers.go`
    - `internal/datalink/dbtarget/service.go`
    - `internal/datalink/dbtarget/service_probe.go`
    - `internal/datalink/dbtarget/service_mysql_test.go`
  - New:
    - `frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx`
    - `frontend/tests/unit/workbench-v2/autosave-settlement-timeout.test.ts`
    - `internal/datalink/dbtarget/tooling_service_mysql_schema_test.go`
  - Removed: (none)
- Affected behaviour: 既有 MySQL 連接器若依賴 cleartext 認證且伺服器未啟用 TLS，升級後連線會失敗並回報明確錯誤，需改用 TLS 或改用非 cleartext 認證外掛。

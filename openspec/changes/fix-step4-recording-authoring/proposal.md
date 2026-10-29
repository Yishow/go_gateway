## Why

Step 4 的來源角色、名稱、批次選取與錯誤復原不夠一致；MySQL／SQL Server 能力僅允許長表時仍可看到無效寬表選擇。試寫結果未知時，新預覽或切換版本可覆蓋原 operation，使使用者失去查明寫入與清理結果的入口。

### Root Cause

`CircuitProfileConfigFields` 同時暴露來源總勾選及各 profile 勾選，numeric role 選項不完整；`circuitProfileBasic`／`circuitProfileActions` 分開處理顯示名稱、保存比對與來源 metadata。`CircuitPurposeFields` 沒有接收目的地布局能力，`GroupTestWritePanel` 新預覽及 scope reset 會清除未決狀態；安全 gate 與畫面狀態未以同一個已保存身分串起。

## What Changes

- 校正 numeric state 的明確 latest 選擇、counter confirmation 所需名稱的即時提示、新累積 profile 不默選全部來源，以及任何已選 counter 的 measurement 查詢失敗重試；不依名稱推斷語意，也不覆寫已保存 role。
- 以 connector revision 綁定的能力限制布局與儲存選項；未知、更新中、讀取失敗均不冒充支援，既存不支援設定保留可讀且給出原因。
- 將 test-write 非終態 operation 與當前 editor preview 分離保存；重整、切群組、改版後仍可查原 operation，結果未釐清不得用新 preview 隱藏它。
- 計數 binding 的 backend authority 依賴 `freeze-counter-replay-semantics`；本 change 不自行變更 replay 或 live transition，轉換交由 `manage-recording-lifecycle-safely`；整合交付在 backend correctness 與 lifecycle 之後、guided workflow 之前。

### Success Criteria

- numeric state 保留最新整數狀態、gauge 才平均；counter_raw 必須確認語意。只要 counter confirmation 需要迴路名稱，即使只啟用瞬時記錄也先顯示必填理由；一般非計數瞬時來源名稱仍選填。
- 新累積 profile 預設不選來源，明示請選計數器；已保存的明確選擇不重設。任何已選 counter 的查詢失敗都顯示安全原因與重試，不只 cumulative_enabled 才顯示。
- MySQL／SQL Server 不允許新增或送出 managed wide；SQLite／PostgreSQL 仍按實際能力顯示。
- lost response、pending／running、unknown write、cleanup failed／unknown 的原 operation 都可找回；write、readback、cleanup 各自呈現，沒有假成功。

## Capabilities

### New Capabilities

- `step4-recording-authoring`: 保留來源角色、名稱、選取、能力限制及未決試寫復原契約。

### Modified Capabilities

無；既有安全與 backend 契約不改。

## Impact

- Affected specs: `step4-recording-authoring`（新增的可觀察 authoring／recovery 契約）。
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields.tsx`, `frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel.tsx`, `frontend/src/features/datalink/workbench-v2/state/circuitProfileBasic.ts`, `frontend/src/features/datalink/workbench-v2/steps/step4/circuitProfileActions.ts`, `frontend/src/features/datalink/workbench-v2/steps/step4/CircuitPurposeFields.tsx`, `frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel.tsx`, `frontend/tests/unit/workbench-v2/`, `frontend/tests/e2e/`。
  - New: `frontend/src/features/datalink/workbench-v2/state/recordingTestWriteRecovery.ts`（非秘密的 scoped operation reference）；若現有 operation query 不足，擴充 `internal/api/handlers/` 的唯讀復原契約與測試。
  - Removed: 無。

## 執行順序調整

**Supersedes:** 原先整個 change 必須等待前案全部完成的執行安排，僅改為可隔離準備先行、正式啟用與整合驗收維持原依賴順序。authoring 的角色／名稱／累積選擇／讀取重試／布局 guard 可先使用既有 backend 契約完成；saved-role mutation、server operation listing、lifecycle 接線與整合驗收仍有明確後置門檻。 七個 change 與本案需求／驗收不變；本次只更新 artifacts，不代表產品已實作。詳見 design.md「執行分期與交接」。

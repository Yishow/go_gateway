## Why

在既有 `/studio/v2` Step 4 改為設備範圍清楚、每個階段只有一個主要下一步的記錄引導，讓一般使用者不必理解群組、revision 或 operation ID 才能完成；專業細節仍可展開。

### Motivation

現在 Basic 的來源、型態、資料表準備與開始並列，Advanced 再列群組及工作區總覽，容易讓人不知道下一步或誤認工作區會全部啟用。`Step4Database` 已不接 workspace-wide start，因此「進階群組與工作區啟用」也不符合實際行為。

## What Changes

- 以「選設備與資料 → 記錄方式 → 確認儲存規則 → 準備並開始」呈現 Step 4 內的進度；只突出當前安全可執行動作，阻擋時改為具體修復入口。
- 將現有 Basic／Advanced 共用的 canonical group 狀態整理為「已設定的記錄」卡片，直接顯示修改、停止新記錄、移除設定及能力允許的轉換入口。
- Advanced 改為「進階設定與工作區狀態」，保留 row policy、型別／角色、欄位、revisions、schema、計數政策與交付證據，不新增 workspace-wide start。
- 預覽已儲存及目前生效規則的差異；開始前明列本次設備／profile／目的地範圍，完成後導向同 scope 的 runtime monitor。
- 提供不依賴 URL 的 typed scope selection interface，驗證 workspace/device/group ownership 與最新 revisions；歷史提示不作 mutation 授權。URL parsing／resolver links 交由後續 `expose-runtime-recording-rules`。
- 重用 `fix-step4-recording-authoring` 與 `manage-recording-lifecycle-safely`，不在 UI 放寬 capability、ownership、confirmation 或 unknown-operation 保護；新 live transition 仍須後端正確性前置通過。

## Capabilities

### New Capabilities

- `guided-step4-recording`: 同設備範圍的下一步引導、摘要、管理與安全選取介面。

### Modified Capabilities

無；保留既有 guided-recording-workflow 契約。

## Impact

- Affected specs: `guided-step4-recording`（新增的 Step 4 呈現契約，補充既有 `guided-recording-workflow`）。
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx`, `frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx`, `frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel.tsx`, `frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileSchemaPreparation.tsx`, `frontend/src/features/datalink/workbench-v2/steps/step4/WorkspaceScopeNote.tsx`, `frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx`, `frontend/src/i18n/locales/en/workbench-v2.json`, `frontend/src/i18n/locales/zh-TW/workbench-v2.json`, `frontend/tests/unit/workbench-v2/`, `frontend/tests/e2e/`。
  - New: `frontend/src/features/datalink/workbench-v2/steps/step4/RecordingNextAction.tsx`, `frontend/src/features/datalink/workbench-v2/steps/step4/RecordingRuleSummary.tsx`（沿既有 state／hooks 的呈現元件）。
  - Removed: 無獨立 route、API 或產品能力；只在呼叫端遷移完成後移除重複操作區。

## 執行順序調整

**Supersedes:** 原先整個 change 必須等待前案全部完成的執行安排，僅改為可隔離準備先行、正式啟用與整合驗收維持原依賴順序。只新增隔離的 presentation／typed selection 元件可先行；掛載到 Step 4、接 lifecycle actions 及最終 browser 驗收，仍等待 authoring 與 lifecycle 交接。 七個 change 與本案需求／驗收不變；本次只更新 artifacts，不代表產品已實作。詳見 design.md「執行分期與交接」。

**Supersedes（安全讀取順序補正）:** 上述整案等待僅保留於 runtime dashboard／locator／最終驗收；runtime 1.3–1.4 的單一安全 projection、全舊新 redaction、applied-config-provider 與共用 service/hook 在 freeze→cutover→lifecycle 後先交付，guided 2.2 與 authoring 前置一併消費。區分 saved、scheduled/latest requested applied 與 as-of effective revision／effective_at，retained snapshot 不代表 intake permission；詳見 design「安全生效證據交接」。已完成歷史不變，以 pending compensation 修正 caller 契約。

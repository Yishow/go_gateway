## Why

目前 Modbus Share 主要把 browser session 中的暫存 mapping 當成可操作狀態，導致 workspace hydration 尚未完成時 UI 可能顯示 display default、全域設定停用時 rule flag 仍可能越權、重啟後 listener 與 mapping 無法可靠恢復，且錯誤可能留下半套或 stale register。這對 SCADA 而言不是單純 UI 缺陷，而是會讀到錯誤資料、誤刪其他 workspace mapping 或在未持久化配置上啟動服務的生命週期風險；現在需要在既有 source-rule candidate/apply/restore seam 上建立 backend-authoritative、可恢復的正式契約。

## What Changes

- 將 /studio/v2 bootstrap hydration、global settings 與 workspace revision/readiness 納入 Share 操作前置條件；hydration 完成前 summary、list、upsert、delete、activation 一律 fail-closed，display defaults 只作顯示而不作 truth。
- 以 persisted workspace mapping 與可由 source-rule/tag 關聯證明的 durable ownership 建立 workspace scope；非空或未持久化的 tag_id 不取得刪除權，所有 API 都拒絕越界操作。
- 建立明確的 bytes/register register contract：human 40001 與 backend zero-based index 互轉，int32/int64/float 的多 register span、stride 最小寬度、capacity、上界與 range collision 由 backend 完整驗證。
- 將 reconcile 設計為完整 prevalidation 後的 atomic、serialized、idempotent mutation；以 workspace revision/CAS 防 double submit 與 concurrency，失敗保留舊狀態或回報可恢復的 dirty/unknown，不採先刪後半套 upsert。
- activation 必須等待必要 durable saves 成功，沿用同一 workspace revision/readiness；save error 或未完成 autosave barrier 會阻擋 activation。
- 接入既有 source-rule candidate/apply/restore seam；backend activation 是唯一 projection authority，direct API 不得繞過 Share projection，process restart 依 persisted desired mappings 自動恢復。
- 讓 settings.modbus_share.enabled 真正控制 server lifecycle：安全預設停用，enable/disable/start/stop、bind address、port、slave id 都從 durable settings 讀取；不再無條件固定 5020，失敗狀態提供 typed error code 與 operator-safe 診斷。
- delete、move、retype mapping 時清除舊 datatype span 或明確標記 invalidated，避免 SCADA 讀到 stale register；完整 acceptance 覆蓋 configure→save→activate→Modbus read→EXE restart→restore/read，以及 disabled、collision、capacity、partial failure、ownership、double-submit。
- 將 legacy /studio retirement 排除在本 change 外；A release-correctness 必須先完成，C legacy 退場後再由後續工作處理；不泛化 connector/protocol transport，也不把 field hardware acceptance 偽裝成 repo 測試結果。

## Non-Goals

- 不修改 legacy /studio、/gateway/* 或建立平行產品入口；本 change 只服務 /studio/v2，並依賴 C 後續退場決策。
- 不重寫既有 rule-level durable fields 或另建平行 mapping model；沿用已存在的 share_enabled、share_start_register、share_stride 與 source-rule candidate/apply/restore seam。
- 不改造 Modbus、FATEK、MC 或其他 connector/protocol transport；本 change 只定義 Share output projection 與 listener lifecycle。
- 不修改 internal/protocol/ 或 internal/datalink/connector/ 的 connector、transport、frame、polling 或 protocol behavior。
- 不把真實 PLC/SCADA、交換器、現場硬體 acceptance 宣稱為自動化測試；field hardware acceptance 另列為部署 gate。
- 不負責完整 i18n presentation；backend contract 只提供可行動且不洩漏敏感資訊的 diagnostics，前端文字由既有 i18n/C 相關工作承接。

## Capabilities

### New Capabilities

- modbus-share-lifecycle: 定義 /studio/v2 Modbus Share 的 hydration、global gate、durable ownership、register contract、atomic reconcile、activation barrier、restart recovery、listener lifecycle、stale-span invalidation 與 acceptance contract。

### Modified Capabilities

- protocol-servers: 將 Local Modbus Share 從固定 5020 listener 改為由 durable global settings 控制、disabled-by-default 且可報告 lifecycle failure 的 server contract。
- local-modbus-memory-workbench: 將 rule-scoped candidate/apply 與 register review 明確收斂到 hydrated、workspace-scoped、backend-authoritative projection；保留 UI 作為 review/diagnostic surface，不以 browser session 取得 ownership 或 runtime truth。

## Impact

- Affected specs:
  - New: openspec/changes/productionize-modbus-share-lifecycle/specs/modbus-share-lifecycle/spec.md
  - Modified: openspec/changes/productionize-modbus-share-lifecycle/specs/protocol-servers/spec.md
  - Modified: openspec/changes/productionize-modbus-share-lifecycle/specs/local-modbus-memory-workbench/spec.md
- Affected code and approved Share/server seams:
  - Backend API/workspace/rule/runtime: internal/api/, internal/datalink/sourcerule/, internal/datalink/workspace/, internal/datalink/runtime/, internal/datalink/schema/migrations/, internal/datalink/mapping/
  - Share service seam: internal/datalink/modbusshare/
  - Virtual memory/server seams: internal/virtual/memory/ and internal/virtual/server/modbus/
  - Frontend: frontend/src/pages/datalink/workbench-v2/, frontend/src/features/datalink/workbench-v2/, frontend/src/services/, frontend/src/types/, frontend/src/i18n/
  - Tests: corresponding internal/api/, internal/datalink/sourcerule/, internal/datalink/workspace/, internal/datalink/runtime/, internal/datalink/mapping/, internal/datalink/modbusshare/, internal/datalink/schema/, internal/virtual/, frontend/tests/unit/, frontend/tests/integration/, and frontend/tests/e2e/ contracts
- Explicitly excluded from this change: internal/protocol/ and internal/datalink/connector/ connector/transport implementation; release packaging owned by A; legacy retirement owned by C.
- Release prerequisite: A release-correctness change must be complete before implementation; C legacy retirement remains a later dependency.
- External acceptance: a documented field test must verify a real Modbus client/SCADA read and restart recovery on the target EXE; it is not substituted by unit, integration, build, or Spectra validation.

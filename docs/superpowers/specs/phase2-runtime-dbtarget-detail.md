# Phase 2 細化規劃：Runtime / Live Value + Database Target

- 基於：`docs/superpowers/specs/2026-03-15-datalink-workbench-design.md` §14–§15
- 初版日期：2026-03-16
- **本次修訂：2026-07-25 — 全面配合「UI 全新重設計」方向**
- 修訂要旨：底層 hooks / services / types 沿用；**所有 UI 組件（含 MemoryGrid 及 workbench 頁面元件）均列為「待重新設計」**，本文不預設任何舊元件名稱或佈局，改以功能區域（functional zone）與 UI 契約（prop / state shape）描述。

---

## 0. 沿用 vs 重做判定表

> **核心原則**：Domain 層與 API 層穩定沿用，UI 呈現層全部重做。

### 0.1 ✅ 沿用層（Domain / API / State）

| 層級 | 資產 | 說明 |
|------|------|------|
| **Go 後端 API** | 所有 §1 列出的 REST / SSE endpoint | 契約不變，後端無需因 UI 重設計而修改 |
| **Go 後端 Domain** | `runtime.Service`、`dbtarget.Service`、`dbtarget.MappingService`、migration 表結構 | 業務邏輯與 DB schema 不受 UI 影響 |
| **前端 Services** | `frontend/src/services/datalink.ts`（`deviceAPI`, `pointAPI`, `tagAPI`, `mappingAPI`, `modbusShareAPI`, `pollingGroupAPI`）| HTTP 呼叫層穩定，新 UI 繼續透過這些 service 存取後端 |
| **前端 React Query Hooks** | `useDevicesQuery`, `usePointsQuery`, `useTagsQuery`, `useMappingsQuery`, `usePollingGroupsQuery`, `useTestConnectionMutation` 等 | Query / mutation hooks 為 domain-bound，與 UI 解耦 |
| **前端 Types** | `frontend/src/types/datalink.ts`（`Device`, `Point`, `Tag`, `Mapping`, `PollingGroup`, `DataType`, `ProtocolType`, `TransformStep` 等）| Type 定義反映 domain model，不隨 UI 變 |
| **前端 Feature Helpers** | `frontend/src/features/datalink/`（`typedOccupancy.ts`, `batchNaming.ts`, `commitAudit.ts`, `commitLifecycle.ts`, `allocationStrategy.ts`, `tagEditImpact.ts`, `validationFlow.ts` 等）| 演算法與規則邏輯，UI-agnostic |
| **Phase 2 新 Hooks（待建）** | `useRuntimeStream`, `useRuntimeStatus`（§2.5 定義的 interface）| Hook 的 input/output shape 不變；內部 SSE / polling 邏輯不變 |
| **Phase 2 新 Services（待建）** | `dbTargetAPI`（connector CRUD / test / tables / mapping / validate）| 純 HTTP 呼叫層 |
| **Phase 2 新 Types（待建）** | `RuntimeStreamEvent`, `RuntimeStatus`, `LivePointValue`, `DatabaseConnector`, `DatabaseTargetMapping`, `DbTargetValidationIssue` 等 | Type 定義反映 domain model |

### 0.2 🔄 重做層（UI 組件）

| 舊資產 | 處置 | 原因 |
|--------|------|------|
| `SmartDashboardPage.tsx`（1,498 行）| **不沿用** | 肥大頁面元件，融合了 workflow / modal / data mapping / action routing，不符合新架構 |
| `LocalModbusWorkbenchPage.tsx` | **不沿用** | 將合併為 Output 步驟的一部分，不再是獨立頁面 |
| `MemoryGrid.tsx`（382 行）| **不沿用舊元件** | UI 全面重設計；但其「位址空間視覺化」的核心概念與 cell 狀態模型可作為新元件的設計參考 |
| `smart-dashboard/` 下所有子元件 | **不沿用** | Header / Workspace / ModbusPanel / CommitFlow 等均為舊 layout 專屬 |
| `useSmartDashboard*` hooks | **不沿用** | 這些是 UI 編排 hooks（非 domain hooks），與舊佈局耦合 |
| Phase 1 spec 提到的 `WorkbenchStepRail`, `WorkbenchContextBar`, `SourceCanvasSection`, `SourceInspectorPanel`, `BatchTagBinder`, `OutputTargetsPanel`, `WorkbenchBottomBar` | **不預設沿用** | 這些是舊 spec 的元件名，新設計將從功能區域出發重新定義元件邊界 |

> **注意**：「不沿用」不代表「刪除」。舊路由 `/datalink` (SmartDashboard) 保留作為 rollback 通道（見 §5.2），但新 workbench 不從舊元件 fork，而是從零構建。

---

## 1. API 契約（不受 UI 重設計影響，完整保留）

### 1.1 Runtime Status（REST snapshot）

**`GET /api/v1/datalink/runtime/status`**

> 回傳 collector 整體與各設備的即時健康摘要。

Request：無 body。可選 query `?device_id=xxx` 過濾單台設備。

Response：
```
{
  "running": bool,
  "uptime_seconds": number,
  "collectors": [
    {
      "device_id": string,
      "device_name": string,
      "protocol": ProtocolType,
      "status": "idle" | "running" | "warning" | "error",
      "points_total": number,
      "points_healthy": number,
      "points_stale": number,
      "points_error": number,
      "last_read_at": string | null,
      "last_error": string | null,
      "breaker_state": "closed" | "open" | "half-open"
    }
  ]
}
```

設計原則：
- 與現有 `GET /dashboard/device-statuses` 的差異是「runtime 維度」——回傳的是 scheduler 看到的即時統計，不是 device 靜態狀態。
- `breaker_state` 直接暴露 circuit breaker 狀態。
- 回傳欄位全部可從現有 `Scheduler`、`ConnectionManager`、`health.Tracker` 組裝，無需新 DB 表。

### 1.2 Runtime Stream（SSE）

**`GET /api/v1/datalink/runtime/stream`**

> 持久 SSE 連線，推送即時 point 值變化。

Query params：
- `device_id`（必填）：只訂閱一台設備的 point 變化。
- `point_ids`（選填，逗號分隔）：進一步過濾。

Event shapes：

**`value` 事件（主資料）**
```
event: value
data: {
  "device_id": string,
  "point_id": string,
  "address": string,
  "raw_value": string | number | boolean | null,
  "transformed_value": string | number | boolean | null,
  "quality": "good" | "bad" | "uncertain",
  "stale": bool,
  "timestamp": string
}
```

**`status` 事件（設備狀態變化）**
```
event: status
data: {
  "device_id": string,
  "status": "idle" | "running" | "warning" | "error",
  "breaker_state": "closed" | "open" | "half-open",
  "last_error": string | null
}
```

**`heartbeat` 事件**
```
event: heartbeat
data: { "ts": string }
```

設計原則：
- 必須 device-scoped。全量推送在 100+ point 場景會壓垮瀏覽器 EventSource。
- `stale` 由後端計算（比較 last_read_at 與 polling interval），前端不做時間比較。
- quality 三值對齊現有 `QualityFlag`（good / bad / uncertain），SSE 用 JSON string。
- 後端實作：在 `runtime.Service` 的 `handleCollectedValue` 流程內，廣播到 channel hub，SSE handler 訂閱。複用現有 `SSEHandler` 的 client map + Flusher 模式。

### 1.3 Point Poll（單次觸發）

**`POST /api/v1/datalink/points/:id/poll`**（已存在）

Phase 2 確認 response 須包含：
```
{
  "point_id": string,
  "value": unknown,
  "transformed_value": unknown | null,
  "quality": "good" | "bad" | "uncertain",
  "stale": bool,
  "timestamp": string,
  "error": string
}
```

- `transformed_value`：查 mapping → 執行 `ExecutePipeline` → 回傳。
- `stale`：比較 timestamp 與該 point 所屬 polling group 的 interval。
- 批次版 `POST /points/poll` 同理，陣列回傳。

### 1.4 Database Target 契約

#### 1.4.1 連接器管理

**`GET /api/v1/datalink/db-targets/connectors`**
```
Response: {
  "connectors": [
    {
      "id": string,
      "name": string,
      "kind": "sqlite" | "postgres" | "mysql" | "sqlserver",
      "status": "ready" | "unreachable" | "auth_failed" | "error",
      "last_check_at": string | null,
      "last_error": string | null
    }
  ]
}
```

**`POST /api/v1/datalink/db-targets/connectors`**
```
Request: {
  "name": string,
  "kind": "sqlite" | "postgres" | "mysql" | "sqlserver",
  "connection_string": string
}
```

**`POST /api/v1/datalink/db-targets/connectors/:id/test`**
```
Response: {
  "reachable": bool,
  "latency_ms": number,
  "error": string | null
}
```

#### 1.4.2 目標表探索

**`GET /api/v1/datalink/db-targets/connectors/:id/tables`**
```
Response: {
  "tables": [
    {
      "schema": string,
      "name": string,
      "columns": [
        {
          "name": string,
          "type": string,
          "nullable": bool,
          "is_primary_key": bool
        }
      ]
    }
  ]
}
```

#### 1.4.3 Tag → Column 映射

**`GET /api/v1/datalink/db-targets/mappings`**
Query：`?connector_id=xxx`

**`PUT /api/v1/datalink/db-targets/mappings/:tagId`**
```
Request: {
  "connector_id": string,
  "table_schema": string,
  "table_name": string,
  "column_name": string,
  "write_mode": "insert" | "upsert",
  "timestamp_column": string | null
}
```

**`POST /api/v1/datalink/db-targets/mappings/validate`**
```
Request: { "connector_id": string }
Response: {
  "valid": bool,
  "issues": [
    {
      "tag_id": string,
      "tag_key": string,
      "issue": "type_mismatch" | "column_not_found" | "table_not_found" | "connector_unreachable" | "missing_mapping",
      "detail": string
    }
  ]
}
```

設計原則：
- 路由前綴 `/db-targets/` 與現有 `/modbus-share/` 平行。
- `write_mode` 只有 insert / upsert 兩種，不做 delete/truncate。
- validate 回傳 issue 陣列而非 boolean，讓 UI 可精準對應。
- connector 的 `connection_string` 只在 create/update 時傳入，GET 不回傳。

---

## 2. 前端 UI 呈現策略（全新 Workbench Component System）

> **本節不指定元件名稱、佈局骨架或元件樹**。UI 設計將以獨立的 design phase 決定。本節定義的是：各功能區域（functional zone）需要消費的資料、需要驅動的行為、以及 quality 標準。

### 2.1 功能區域定義

Phase 2 前端需支援以下功能區域。每個區域的具體元件邊界、視覺設計、佈局位置由 UI 設計階段決定。

#### Zone A：即時值顯示（Live Value Display）

**用途**：在使用者操作位址空間或點位清單時，同步顯示 point 的即時採集值。

**UI 契約（需要的資料）**：
- 每個 point 的 `{ raw_value, transformed_value, quality, stale, timestamp }`
- 資料來源：`useRuntimeStream` hook 的 `getLiveValue(pointId)` selector

**行為需求**：
- 值變化時有微動畫回饋（pulse / highlight），持續 ≤ 300ms
- quality 三態 + stale 須有視覺區分（色彩 + 圖示，具體設計不預設）
- 不可在 SSE 斷線時顯示假值——凍結最後已知值並標示斷線

**不預設的事**：
- 不預設這是 grid cell、table row、card、或任何特定元件
- 不預設 live value 的位置（右上角、行內、overlay 等）

#### Zone B：設備健康摘要（Device Health Summary）

**用途**：讓使用者一眼了解目前所選設備的 runtime 狀態與整體 point 健康度。

**UI 契約**：
- `{ running, status, breaker_state, points_total, points_healthy, points_stale, points_error, last_read_at }`
- 資料來源：`useRuntimeStatus` hook（REST 10s 輪詢）+ `useRuntimeStream` 的 `status` 事件（即時更新）

**行為需求**：
- 三種主要狀態需視覺區分：Running / Warning / Stopped
- Breaker state（closed / open / half-open）需可檢閱（展開、tooltip、或其他互動）
- 提供「啟動 / 停止 Collector」的操作進入點

#### Zone C：資料庫輸出設定（Database Target Configuration）

**用途**：完整的 DB connector 管理 → 表探索 → tag-to-column 映射 → 驗證 → 啟用 流程。

**UI 契約**：
- Connector 列表 + CRUD：`useDbConnectors` hook
- 連線測試結果：`{ reachable, latency_ms, error }`
- Table/column schema 樹：`useDbTables(connectorId)` hook
- Mapping 列表與 CRUD：`useDbTargetMappings(connectorId)` hook
- Validation issues 陣列：`useDbTargetValidate(connectorId)` hook

**行為需求**：
- 設定流程是線性引導（select connector → explore tables → create mappings → validate → enable），但使用者可在步驟間自由切換
- Type mismatch 等 issues 需與對應 tag 視覺關聯（不預設是 inline badge、toast、或 panel）
- Blocked 條件（§3.4 列舉的 7 種）必須 proactively 顯示，不能等使用者點「啟用」才報錯

#### Zone D：Validation 狀態總覽（Readiness Summary）

**用途**：跨所有 output target（Modbus + Database）的整體就緒度。

**UI 契約**：
- Modbus 映射計數 + 狀態
- Database 映射計數 + validation 狀態
- Runtime 是否 running
- 資料來源：組合 `useRuntimeStatus` + `useDbTargetValidate` + 現有 `useMappingsQuery`

**行為需求**：
- 一目了然的 overall readiness indicator（非多處分散）
- 可展開查看各 target 的細項問題

### 2.2 Quality 視覺語意（設計約束，非色彩指定）

| Quality 狀態 | 語意層級 | 視覺區分需求 |
|-------------|---------|-------------|
| `good` | Positive / success | 主要（正常態），不應過度搶眼 |
| `uncertain` | Caution / degraded | 需與 good 明確區分，提示值可能過時 |
| `bad` | Error / failure | 最高對比，立即引起注意 |
| `stale` | Inactive / outdated | 弱化顯示，區別於 bad（stale 有舊值，bad 無值） |

> 具體色彩、圖示、動畫由 UI design phase 決定。上表僅約束語意層級的相對關係。

### 2.3 Stale / Error / Disconnected 處理規則

這些規則是行為需求，任何 UI 設計都必須遵守：

- **Stale**：由後端 SSE event 的 `stale: true` 決定，前端不做本地計時。值保留但視覺弱化。
- **Error（quality=bad）**：顯示錯誤標示 + 可查看 `last_error` 全文。提供建議動作入口（「檢查連線」「確認地址」）。
- **斷線（SSE 連線中斷）**：
  - 在適當範圍顯示斷線提示「即時連線已中斷，正在重連…」
  - 值凍結為最後已知值 + stale 視覺
  - 超過 3 次重連失敗後，提供手動重連操作
  - 自動降級為 REST 輪詢（見 §5.2 Fallback）

### 2.4 前端 State 管理（Hook 層，UI 無關）

以下 hooks 定義的是 interface 與行為契約，可被任何 UI 元件消費：

**`useRuntimeStream(deviceId: string)`**：
- 維護 `Map<pointId, LivePointValue>` in-memory store
- 提供 `getLiveValue(pointId): LivePointValue | undefined` selector
- 暴露 `connectionState: 'connecting' | 'connected' | 'reconnecting' | 'disconnected' | 'fallback'`
- 自動 reconnect，exponential backoff，max 5 attempts
- 超過 max attempts → 切換 `connectionState` 為 `'fallback'`

**`useRuntimeStatus(options?: { deviceId?: string })`**：
- 每 10 秒 `GET /runtime/status`
- 回傳 `RuntimeStatus` 型別（含 `running`, `uptime_seconds`, `collectors[]`）
- 提供 device-level 篩選

**`useDbConnectors()`** / **`useDbTables(connectorId)`** / **`useDbTargetMappings(connectorId)`** / **`useDbTargetValidate(connectorId)`**：
- 標準 React Query hooks，CRUD + query
- 型別定義在 `frontend/src/types/datalink.ts` 擴充

### 2.5 新 UI 元件系統的設計約束

雖然不預設元件名稱與佈局，以下約束是新 workbench component system 必須遵守的設計邊界：

1. **頁面元件不超過 300 行**：避免重蹈 SmartDashboardPage（1,498 行）的覆轍。超過則拆分。
2. **狀態由 Provider/Context 統管**：workbench 級共享狀態（selected device、current step、live stream connection）由統一的 context provider 管理，不在頁面元件中堆積 useState。
3. **Modal 最少化**：主流程不使用全螢幕 modal。優先使用 inline 展開、drawer、或 step transition。
4. **Step 間狀態保持**：使用者在步驟間切換時，前一步的選擇與編輯狀態不丟失。
5. **漸進式揭露**：Phase 2 功能（live value、database target）在後端未就緒時以 disabled / coming-soon 狀態出現，不隱藏。
6. **Live value 是非阻塞疊加**：live value 的呈現不能改變底層操作流程（選擇、編輯、拖拽等）。它是「疊加層」——即使 SSE 斷線，所有靜態操作仍可正常進行。
7. **響應式但桌面優先**：workbench 以 ≥1280px 寬度為主要設計目標。窄螢幕可隱藏次要區域但不破壞核心流程。
8. **元件可獨立測試**：每個功能區域對應的元件必須可以用 mock hook 獨立渲染測試，不依賴整個 workbench context。

---

## 3. Database Target Domain 物件與流程

> 本節為 domain 層，不受 UI 重設計影響，完整保留。

### 3.1 需要的 Domain 物件（後端）

**`DatabaseConnector`**（新 DB 表 `database_connectors`）
- `id` (UUID)
- `name` (string, unique)
- `kind` (enum: sqlite / postgres / mysql / sqlserver)
- `connection_config` (encrypted JSON，含 host/port/user/password/dbname)
- `status` (enum: ready / unreachable / auth_failed / error)
- `last_check_at` (timestamp)
- `last_check_error` (string)
- `enabled` (bool)
- `created_at`, `updated_at`

**`DatabaseTargetMapping`**（新 DB 表 `database_target_mappings`）
- `id` (UUID)
- `tag_id` (FK → tags)
- `connector_id` (FK → database_connectors)
- `table_schema` (string)
- `table_name` (string)
- `column_name` (string)
- `write_mode` (enum: insert / upsert)
- `timestamp_column` (string, nullable)
- `enabled` (bool)
- `created_at`, `updated_at`
- UNIQUE constraint: `(tag_id, connector_id)`

### 3.2 設定流程（使用者操作順序）

1. **選擇 / 建立 Connector**
   - 選擇已有 connector 或建立新的
   - 填入連線參數（host / port / user / password / dbname）
   - 觸發連線測試 → 呼叫 `POST .../test`
   - 測試通過 → connector status = `ready`
   - 測試失敗 → 顯示具體錯誤，不允許繼續

2. **探索目標表**
   - Connector ready 後，呼叫 `GET .../tables`
   - 呈現 schema → table → columns 的層級結構（具體 UI 元件不預設）
   - 使用者選擇目標 table

3. **建立 Tag → Column 映射**
   - 列出已 active 且有 point mapping 的 tag
   - 列出目標 table 的 columns
   - 提供對應操作（具體互動方式不預設——可為拖拽、下拉、或其他）
   - 系統 proactively 提示 type mismatch

4. **驗證完整性**
   - 呼叫 `POST .../validate`
   - 呈現所有 issues（與對應 tag 關聯）
   - 全部解決 → 狀態 `valid`
   - 有問題 → 狀態 `warning` 或 `blocked`

5. **啟用寫入**
   - 驗證通過後，啟用 mapping
   - Runtime Service 的 ingestor 在寫入 timeseries 的同時，寫入目標 DB

### 3.3 Validation 狀態機

| 狀態 | 條件 | 語意 |
|------|------|------|
| `draft` | 有 connector 但無 mapping | 初始態 |
| `incomplete` | 有部分 mapping 但 validate 有 issues | 需人工介入 |
| `valid` | validate 回傳 `valid: true` | 可啟用 |
| `active` | mapping 已啟用寫入 | 寫入進行中 |
| `error` | connector unreachable 或寫入失敗 | 需排障 |
| `blocked` | 見 §3.4 blocked 條件 | 結構性問題，需修正前置條件 |

> UI 如何呈現這些狀態（色彩、圖示、位置）由 UI design phase 決定。本規格只約束每個狀態必須有獨立視覺區分且使用者可辨識。

### 3.4 Blocked 條件（不允許啟用或操作）

1. **Connector unreachable**：測試連線失敗 → 不允許建立 mapping
2. **Tag 未 active**：tag 狀態為 draft / retired → 不允許建 mapping，提示「請先啟用 Tag」
3. **Tag 無 point mapping**：tag 沒有上游 point → mapping 無意義，提示「此 Tag 無資料來源」
4. **Type 不相容且無轉換**：tag dataType 與 column type 無法隱式轉換 → blocked，建議加 transform
5. **Column 不存在**：目標表結構變更 → 映射失效，需重新選擇
6. **重複映射衝突**：同 tag 已映射到同 connector 的另一 table/column → 顯示衝突
7. **Collector 未啟動**：runtime 未 running → 可設定但標示「寫入功能需 Collector 運行」

---

## 4. Release 邊界：骨架 vs 不可假做

### 4.1 可以先建骨架（在 UI 重設計的 Phase 1 交付中）

| 項目 | 骨架做法 | 理由 |
|------|---------|------|
| Database Target 功能區域 | 渲染佔位區塊，顯示「資料庫輸出（即將推出）」+ disabled 互動 | 讓使用者知道功能規劃存在，不預設元件形式 |
| Output Target 的 tab / section 結構 | Local Modbus 可操作，Database 區域 disabled | 資訊架構預佈線，不開放未就緒功能 |
| Readiness Summary 的 database 項 | 顯示 `—` 或 `N/A` | 佔位，數值來自 API |
| `useRuntimeStream` hook | 定義完整 interface，內部回傳 `connectionState: 'disconnected'` + 空 Map | 讓新 UI 元件提前消費 hook，Phase 2 啟用時只需 feature flag 切換 |
| Live value 顯示區域 | 預留 value display 位置，Phase 1 顯示 `last_value` 靜態值（來自 point query） | Phase 2 接入 SSE 時替換資料源，不改 UI 結構 |

### 4.2 不可假做（必須等後端就緒）

| 項目 | 原因 |
|------|------|
| 假的 live value 動畫或隨機值 | 使用者無法區分真假資料，會誤判系統狀態 |
| 假的 database connector 列表 | 使用者嘗試設定卻無法儲存，喪失信任 |
| 假的 validate 結果 | 必須反映真實 DB schema，假 result 比沒有更糟 |
| SSE 連線的假 heartbeat | 前端誤判連線狀態，隱藏真實問題 |
| Write mode 設定 UI（insert/upsert）| 需後端實際支援才能承諾行為 |

### 4.3 灰色地帶的處理原則

- **Point poll**：現有 `POST /points/:id/poll` 已實作，但缺 `transformed_value`。Phase 1 可用現有 response；Phase 2 擴充 response 欄位（向後相容，加欄位不刪欄位）。
- **Runtime status**：Phase 1 可用現有 `GET /dashboard/device-statuses` 作為降級方案。Phase 2 新增 `/runtime/status` 後切換。
- **MemoryGrid 概念**：即使舊 `MemoryGrid.tsx` 元件不沿用，「位址空間視覺化」的概念（cell state model、span-by-datatype 計算、address range navigation）在新 UI 中仍是核心功能需求。新元件的設計文件應參考 `MemoryGrid` 的 cell state 定義（available / planned / used / linked / conflict / selected）作為功能基線。

---

## 5. Rollout / Fallback / Migration 建議

### 5.1 Rollout 策略

**Phase 2a — Runtime / Live Value（先做）**
1. 後端：在 `runtime.Service` 新增 value broadcast hub → SSE handler 訂閱
2. 後端：新增 `GET /runtime/status` endpoint
3. 後端：擴充 `POST /points/:id/poll` response 加入 `transformed_value` 與 `stale`
4. 前端：實作 `useRuntimeStream` + `useRuntimeStatus` hooks（已有 interface 定義）
5. 前端：在新 workbench 的 Zone A（Live Value Display）與 Zone B（Device Health Summary）接入 hooks
6. 驗收：選一台有 polling 的設備，確認值即時顯示且 quality 正確

**Phase 2b — Database Target（後做）**
1. 後端：新增 `database_connectors` + `database_target_mappings` 表（migration 004）
2. 後端：實作 connector CRUD + test + table exploration
3. 後端：實作 mapping CRUD + validate
4. 後端：在 ingestor 流程加入 target write 分支
5. 前端：在新 workbench 的 Zone C（Database Target Configuration）啟用完整功能
6. 前端：Zone D（Readiness Summary）接入 database validation 數據
7. 驗收：建立 connector → 探索表 → 映射 → validate → 啟用寫入 → 確認資料寫入

### 5.2 Fallback 機制

| 場景 | Fallback |
|------|---------|
| SSE 連線不穩定 | 前端 `useRuntimeStream` 自動降級為 REST 輪詢（`connectionState: 'fallback'`）；live value 顯示靜態值 + 降級提示 |
| Runtime status endpoint 不可用 | 使用現有 `GET /dashboard/device-statuses` + `GET /dashboard/stats` 組合 |
| Database connector 連線失敗 | mapping 可設定但 status = error；不阻擋其他 output target；寫入失敗計入 error 計數但不中斷採集 |
| 新 workbench UI 重大問題 | 舊 SmartDashboard 路由 `/datalink` 保留不刪除，使用者可隨時切回。兩套 UI 共用相同 services / hooks / types，資料一致 |
| Feature flag 關閉 | `runtime_stream_enabled` = false → live value zone 顯示靜態值；`db_target_enabled` = false → database zone 顯示 coming-soon |

### 5.3 Migration 注意事項

1. **DB Migration 004**：只加新表、不改現有表。`database_connectors` 與 `database_target_mappings` 是獨立表，不需 FK 到 `devices` 或 `mappings`。
2. **API 版本**：新 endpoint 全部掛在 `/api/v1/datalink/` 下，不需 v2。現有 endpoint 的 response 擴充是向後相容（加欄位不刪欄位）。
3. **SSE 路由**：新增 `/runtime/stream` 與現有 `/preview/stream` 並存，不合併。兩者職責不同（preview 是單一 mapping 的 transform 預覽；runtime 是設備級即時值）。
4. **前端路由**：新 workbench 的路由由 UI design phase 決定（可能是新路由或替換現有路由）。無論如何，舊 `/datalink` (SmartDashboard) 路由在過渡期保留。
5. **Feature flag**：在 `settings` 表新增 `runtime_stream_enabled` 與 `db_target_enabled` 兩個開關。前端讀取 settings → 決定各功能區域的啟用狀態。
6. **前端程式碼共存**：由於底層 services / hooks / types 共用，新舊 UI 不存在資料層衝突。舊元件檔案在新 workbench 穩定後才移除。

---

## 6. 子代理任務切法

> 以下根據「沿用 vs 重做」判定，重新劃分工作單元。後端 agent 不變，前端 agent 配合 UI 全新設計調整。

### Agent 1：Runtime SSE Backend（不變）
- 範圍：`internal/datalink/runtime/` + `internal/api/handlers/`
- 任務：
  - 在 `runtime.Service` 實作 value broadcast hub（fan-out channel）
  - 新增 `RuntimeSSEHandler`，訂閱 hub，推送 `value` / `status` / `heartbeat` 事件
  - 新增 `GET /runtime/status` handler
  - 擴充 poll handler response 加入 `transformed_value` 與 `stale`
  - 路由掛載
  - 單元測試 + 整合測試
- 依賴：無，可最先啟動
- 產出驗收：`go test ./internal/datalink/runtime/... ./internal/api/handlers/...`

### Agent 2：Runtime Frontend Hooks + Types（不變）
- 範圍：`frontend/src/hooks/datalink/` + `frontend/src/types/` + `frontend/src/services/`
- 任務：
  - 實作 `useRuntimeStream` hook（SSE 連線、reconnect、fallback、LivePointValue map）
  - 實作 `useRuntimeStatus` hook（REST 輪詢）
  - 新增 type 定義：`RuntimeStreamEvent`, `RuntimeStatus`, `LivePointValue`, `RuntimeConnectionState`
  - 單元測試
- 依賴：Agent 1 的 API shape（可用 type 先行）
- 產出驗收：`cd frontend && npm run test && npm run lint`

### Agent 3：Live Value UI — 全新設計（重大變更）
- 範圍：新 workbench 元件目錄（路徑由 UI design phase 決定）
- 任務：
  - **不是接入舊元件**，而是在全新 workbench component system 中設計並實作 Zone A + Zone B
  - Zone A（Live Value Display）：設計新元件消費 `useRuntimeStream`，在位址空間/點位清單中呈現即時值
  - Zone B（Device Health Summary）：設計新元件消費 `useRuntimeStatus`，呈現 runtime 摘要與操作入口
  - 實作 §2.3 的 stale / error / disconnected 行為規則
  - 遵守 §2.5 的設計約束（300 行上限、context provider、minimal modal 等）
- 依賴：Agent 2 完成 + UI 設計方案確定
- 產出驗收：新元件可獨立 mount + 視覺驗收 + `npm run test`
- **⚠ 此 Agent 須等 UI 設計方案確定後才能啟動**

### Agent 4：Database Target Backend（不變）
- 範圍：`internal/datalink/dbtarget/`（新 package）+ `internal/api/handlers/` + migration
- 任務：
  - 新增 migration 004（`database_connectors` + `database_target_mappings` 表）
  - 實作 `dbtarget.Service`（connector CRUD + test + table discovery）
  - 實作 `dbtarget.MappingService`（mapping CRUD + validate）
  - 實作 target writer（在 ingestor 流程中分支寫入）
  - 新增 API handlers + 路由掛載
  - 單元測試 + 整合測試
- 依賴：無，可與 Agent 1 平行
- 產出驗收：`go test ./internal/datalink/dbtarget/... ./internal/api/handlers/...`

### Agent 5：Database Target Frontend — 全新設計（重大變更）
- 範圍：`frontend/src/services/` + `frontend/src/hooks/datalink/` + `frontend/src/types/` + 新 workbench 元件目錄
- 任務：
  - 新增 `dbTargetAPI` service（沿用 service 模式）
  - 新增 React Query hooks（`useDbConnectors`, `useDbTables`, `useDbTargetMappings`, `useDbTargetValidate`）
  - 新增 type 定義
  - **不是實作舊 `DatabaseTargetCard`**，而是在全新 workbench 中設計 Zone C + Zone D
  - Zone C（Database Target Configuration）：設計新元件承載 connector 管理 → 表探索 → mapping → validate → enable 完整流程
  - Zone D（Readiness Summary）：設計新元件整合所有 output target 的 validation 狀態
  - 實作 §3.4 的 blocked 條件判斷與回饋
  - 遵守 §2.5 的設計約束
- 依賴：Agent 4 的 API shape + UI 設計方案確定
- 產出驗收：`cd frontend && npm run test && npm run lint && npm run build`
- **⚠ 此 Agent 須等 UI 設計方案確定後才能啟動**

### Agent 6：Feature Flag + Integration 驗收（調整）
- 範圍：跨後端 settings + 前端 feature detection
- 任務：
  - settings 表新增 `runtime_stream_enabled` / `db_target_enabled`
  - 前端讀取 settings 控制 Zone A–D 的啟用狀態
  - SSE fallback 邏輯（連線失敗降級為輪詢）
  - 確保舊 SmartDashboard 路由仍可正常運作（rollback 通道驗證）
  - E2E 測試：完整 Device → Source → Tag → Output 流程
  - 文件更新：API spec / CLAUDE.md 補充 Phase 2 路由
- 依賴：Agent 1–5 全部完成
- 產出驗收：全量 `go test ./...` + `cd frontend && npm run test && npm run build`

### 平行性分析（更新）

```
Timeline:
───────────────────────────────────────────────────────
Agent 1 (Runtime SSE BE)  ████████░░░░░░░░░░░░░░░░░░░
Agent 2 (Runtime FE Hook) ░░░░████████░░░░░░░░░░░░░░░  (等 Agent 1 type)
                          ┌── UI Design Phase ────────┐
Agent 3 (Live Value UI)   ░░░░░░░░░░░░░░████████░░░░░  (等 Agent 2 + UI 設計)
Agent 4 (DB Target BE)    ████████████░░░░░░░░░░░░░░░  (與 Agent 1 平行)
Agent 5 (DB Target FE)    ░░░░░░░░░░░░░░░░████████░░░  (等 Agent 4 + UI 設計)
Agent 6 (Integration)     ░░░░░░░░░░░░░░░░░░░░████████ (等全部)
───────────────────────────────────────────────────────
```

關鍵變化：
- Agent 3 和 Agent 5 新增「UI 設計方案確定」前置依賴。
- Agent 1、2、4 不受影響，可立即啟動。
- 建議在 Agent 1 + 4 開工的同時，平行進行 UI 設計探索，縮短關鍵路徑。

---

## 附錄 A：與現有架構的對齊確認

| 現有元件 | Phase 2 如何使用 | 受 UI 重設計影響？ |
|---------|-----------------|-------------------|
| `SSEHandler` (generic) | 複用 client map + Flusher 模式，新建 `RuntimeSSEHandler` | ❌ 後端 |
| `Scheduler.valueChan` | runtime.Service 已消費；新增 fan-out 到 SSE hub | ❌ 後端 |
| `ConnectionManager.ConnectionStatus` | 組裝 `/runtime/status` 的 device-level 資料 | ❌ 後端 |
| `health.Tracker` + `health.Breaker` | 提供 breaker_state 給 runtime status | ❌ 後端 |
| `QualityFlag` (good/bad/uncertain) | SSE event 直接使用相同三值 | ❌ 後端 |
| `ExecutePipeline()` | poll response 新增 `transformed_value` 時呼叫 | ❌ 後端 |
| `timeseries` 表 | Database target 寫入是額外分支，不影響現有 | ❌ 後端 |
| settings 表 | 新增 feature flag key | ❌ 後端 |
| `usePreviewStream` hook | `useRuntimeStream` 採用相同 EventSource + reconnect 模式 | ❌ Hook 層 |
| `datalink.ts` services | 新增 `dbTargetAPI`，現有 service 沿用 | ❌ Service 層 |
| `datalink.ts` types | 擴充新 type，現有 type 沿用 | ❌ Type 層 |
| `features/datalink/` helpers | 演算法邏輯沿用（occupancy、naming、validation 等） | ❌ 邏輯層 |
| `SmartDashboardPage.tsx` | 保留作 rollback 路由，不在新 workbench 中使用 | ✅ 不沿用 |
| `MemoryGrid.tsx` | 概念參考（cell state model），元件不沿用 | ✅ 不沿用 |
| `LocalModbusWorkbenchPage.tsx` | 功能合併至新 workbench 的 Output zone | ✅ 不沿用 |
| `smart-dashboard/` 子元件 | 不在新 workbench 中使用 | ✅ 不沿用 |
| `useSmartDashboard*` hooks | UI 編排 hooks，與舊佈局耦合，不沿用 | ✅ 不沿用 |

## 附錄 B：MemoryGrid 概念遷移指南

雖然 `MemoryGrid.tsx` 元件不沿用，其核心設計概念在新 workbench 中仍是必要功能。以下列出需遷移的概念（非程式碼）：

| 概念 | 原始位置 | 在新 UI 中的需求 |
|------|---------|-----------------|
| Cell State Model（available / planned / used / linked / conflict / selected） | `MemoryGrid.tsx` props | 任何「位址空間視覺化」元件都需要這套狀態模型 |
| Span-by-DataType 計算 | `features/datalink/typedOccupancy.ts` | ✅ 邏輯沿用，不需重做 |
| Address Range Navigation（centerAddress + range） | `MemoryGrid.tsx` props | 需提供位址空間導航能力（具體 UI 不預設） |
| Multi-select with contiguous range | `MemoryGrid.tsx` onSelect | 需支援範圍選擇操作（具體互動不預設） |
| Conflict detection visualization | `MemoryGrid.tsx` cell rendering | 衝突必須視覺可辨（具體呈現不預設） |
| Live value overlay | Phase 2 新增 | 在位址空間中疊加即時值（Phase 2 核心功能） |

> 新 UI 設計團隊應將 `MemoryGrid.tsx` 作為「功能需求參考」閱讀，而非「程式碼 fork 來源」。

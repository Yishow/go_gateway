## Context

V2 settings UI 已有完整畫面與本地 reducer，但它目前不是 backend truth。使用者已明確要求 `/studio/v2` 全站對接，因此 settings 必須從「只影響本地 state 的 demo 面」升級為正式產品 surface。

## Goals / Non-Goals

**Goals:**

- settings page 啟動時讀取真實 backend state
- connector pool CRUD / test 使用真實 backend APIs
- `儲存所有設定` 實際寫入 backend，而不是 noop
- 移除 operator-visible mock / noop 行為

**Non-Goals:**

- 不重做 settings 視覺設計
- 不改寫 Step 1-4 autosave 契約
- 不在此 change 內重新設計 runtime dashboard

## Decisions

### Backend becomes settings source of truth

Settings page 進頁時應以 backend state 為真相來源，前端 reducer 僅作為編輯中的工作狀態，不再把初始 defaults 視為產品完成狀態。

### Reuse existing backend endpoints

既有 `/settings` 與 `/db-targets/connectors/*` APIs 已存在，應優先重用並補齊 V2 對映，而不是新增一套只給 V2 的平行 settings contract。

### Save bar remains explicit, but no longer noop

現有 UI 已有 `儲存所有設定` 的操作心智。最小變更是保留 explicit save bar，但讓它真正寫入 backend 並回報失敗，而不是改成另一套 autosave 心智。

## Implementation Contract

- Behavior:
  - settings page 啟動時讀取 backend settings 與 connector data
  - connector CRUD / test 由真實 API 驅動
  - `儲存所有設定` 會將變更寫入 backend，失敗時回報 actionable error
  - operator-visible mock timer、Math.random、warning-only noop 必須移除
- Interface / data shape:
  - 重用 `GET /api/v1/datalink/settings`、`PUT /api/v1/datalink/settings/:key`
  - 重用 `GET/POST/PUT/DELETE /api/v1/datalink/db-targets/connectors`
  - 重用 `POST /api/v1/datalink/db-targets/connectors/:id/test`
- Failure modes:
  - backend 讀取失敗時不得假裝以 defaults 成功啟動
  - save 失敗時不得只留 console warning
  - connector test 失敗不得由隨機值主導
- Acceptance criteria:
  - backend tests 覆蓋 settings read/write 與 connector CRUD/test
  - frontend tests 覆蓋 boot、save、connector test success/failure
- Scope boundaries:
  - In scope: settings boot/persist、connector pool backend wiring
  - Out of scope: layout redesign、runtime page rewrite
- Not complete if:
  - connector test 仍依賴 Math.random 或 mock timer
  - save bar 仍是 noop
  - boot 失敗時仍默默退回本地 defaults 冒充 backend state

## Risks / Trade-offs

- [Risk] 既有 `/settings` key/value shape 與 V2 reducer state 並非 1:1 → Mitigation: 明確建立 mapping layer，避免在 UI 直接耦合 backend storage key
- [Risk] connector pool 既有 backend model 與 V2 顯示欄位差異 → Mitigation: 以 service/hook 做對映，不新增第二套 connector persistence

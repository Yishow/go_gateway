## Context

V2 要成為預設入口時，不能再只依賴 `useWorkbenchV2State()` 的前端本地資料。後續 autosave、runtime 工作區視角、Step 4 第一次啟動都需要一個後端可以識別的工作區。這個 change 先只建立唯一 V2 工作區與 bootstrap API，不把 Step 1~4 的實際資料流一併塞進來。

## Goals / Non-Goals

**Goals:**

- 提供唯一 V2 工作區的後端資料邊界
- 打開 `/studio/v2` 時，若 workspace 不存在就自動建立
- workspace 在服務重啟後仍可讀回
- workspace 與 `/studio` legacy 資料預設隔離

**Non-Goals:**

- 不支援多 workspace 建立與切換
- 不匯入 `/studio` 既有 device / rule / output 資料
- 不在這個 change 內完成 Step 1~4 autosave

## Decisions

### Singleton workspace model

V2 先只允許一個 workspace。這和使用者已經決定的「先只做一個工作區」一致，也能避免在 autosave 之前先陷入 workspace switch 複雜度。

### Auto-bootstrap on read

以讀取 workspace 的 API 直接承擔 bootstrap 行為。前端打開 `/studio/v2` 時只需做一次讀取，不需要多一步建立 workspace 的按鈕或精靈。

### Legacy isolation

workspace 只管理 V2 自己的資料，不會自動把 `/studio` 既有資料搬進來。這可避免在入口切換前，就把 legacy 與 V2 資料模型混成一團。

### Workspace metadata contract

bootstrap response 先只回 workspace metadata 與排序骨架，不在這一步偷塞 Step 1~4 snapshot。後續 stepwise autosave change 再逐段擴充。

## Implementation Contract

- Behavior:
  - 使用者開啟 `/studio/v2` 時，系統可取得唯一 workspace；若不存在則自動建立
  - workspace 在服務重啟後仍可再次讀到
  - V2 workspace 不會自動帶入 `/studio` legacy 主線資料
- Interface / data shape:
  - `GET /api/v1/datalink/studio-v2/workspace`
  - success payload 至少包含：
    - `id`
    - `kind` = `single`
    - `status` = `empty` | `ready`
    - `ordered_device_ids`
    - `created_at`
    - `updated_at`
- Failure modes:
  - 若資料庫初始化失敗，API 回傳可行動錯誤，不回傳半成品 workspace
  - 若 workspace metadata 壞掉，不自動吃 `/studio` 資料當 fallback
- Acceptance criteria:
  - 後端 integration test 覆蓋首次讀取自動建立與重啟後重讀
  - 前端 route boot test 驗證 V2 頁面依 workspace API 啟動，不再只靠本地預設 state
- Scope boundaries:
  - In scope: singleton workspace、bootstrap API、frontend boot hook
  - Out of scope: multi-workspace、legacy import、step autosave
- Not complete if:
  - `/studio/v2` 還要靠前端假資料才能開
  - workspace 重啟後消失
  - bootstrap 失敗時偷偷改吃 `/studio` 舊資料

## Risks / Trade-offs

- [Risk] 只建立 workspace metadata，短期內前端仍需和本地 state 並存 → Mitigation: 後續 stepwise autosave changes 逐步替換各 step source of truth
- [Risk] singleton 模式未來可能要擴成 multi-workspace → Mitigation: API shape 保留 `kind` 與 `ordered_device_ids`，避免之後完全打掉

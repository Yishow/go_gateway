## Context

`/studio/v2` 已經被 inventory 與產品決策定義為主線，但目前 router 仍以 `/studio` 為首頁與 fallback。這個 change 只處理入口矩陣，不碰 V2 內部資料流、Step 4 啟動語意或 runtime 觀察邏輯。

## Goals / Non-Goals

**Goals:**

- 讓使用者從 `/` 與通用入口進站時直接落到 `/studio/v2`
- 保留 `/studio` 直接可達，避免把既有桌面主線整包強制改寫
- 補齊 route regression tests，避免後續 change 再把首頁拉回 `/studio`

**Non-Goals:**

- 不改 `/studio/v2` 內部的自動存、工作區或 runtime handoff
- 不重寫 task-specific legacy deep links 的既有去向
- 不處理 `/studio` 與 `/studio/v2` 的資料互通

## Decisions

### Default route matrix

`/` 與 catch-all fallback 直接導到 `/studio/v2`。這是唯一和產品決策一致、也最容易被測試固定下來的入口規則。

### Legacy redirect boundary

只改「通用 landing 類」legacy 路由，例如 `/datalink`、`/datalink/workbench`。仍帶有功能語意的 legacy deep links 維持原本導向，避免這個 change 偷偷改壞其它尚未整理的工作流。

### Preserve direct `/studio`

`/studio` 仍保留直接開啟 legacy 主線的能力。這個 change 不是表態 legacy surface 已下線，而是先把新的預設入口切正。

## Implementation Contract

- Behavior:
  - 開啟 `/` 時，系統直接進入 `/studio/v2`
  - 開啟未知路由時，系統直接進入 `/studio/v2`
  - 開啟通用 legacy landing 路由時，系統直接進入 `/studio/v2`
  - 直接開啟 `/studio` 時，系統仍進入 legacy `DatalinkWorkbenchPage`
- Interface / data shape:
  - Router contract 仍維持 `react-router` route tree
  - 不新增新 query 參數，不新增新 API
- Failure modes:
  - 若 routing 測試缺漏，最常見回歸會是 fallback 又被改回 `/studio`
  - 若把 task-specific legacy deep links 一起改掉，會造成 scope 外功能漂移
- Acceptance criteria:
  - `frontend/tests/unit/workbench-v2/routing.test.tsx` 覆蓋 `/`、`*`、`/datalink`、`/datalink/workbench` 與 `/studio`
  - 手動確認 `/studio` 不會被自動轉去 `/studio/v2`
- Scope boundaries:
  - In scope: route defaults、fallback、通用 legacy landing redirect
  - Out of scope: V2 autosave、workspace、runtime 啟動
- Not complete if:
  - `/` 還會先進 `/studio`
  - 未知路由還會先進 `/studio`
  - `/studio` 被強制轉成 `/studio/v2`

## Risks / Trade-offs

- [Risk] 只切入口，V2 內部主線仍未完成 → Mitigation: 其它 changes 會分別補 workspace、autosave、activation、runtime
- [Risk] 誤改 legacy deep links → Mitigation: 把 change scope 限在 landing 類路由，並用 routing tests 固定

## Why

目前產品決策已經明確把 `/studio/v2` 定為主入口，但實作上的 `/`、未知路由與通用 legacy 入口仍會把使用者帶到 `/studio`。這讓使用者一進站就落到暫停維護的 surface，也讓後續 V2 主線規劃無法先從入口層收斂。

## What Changes

- 將 `/`、未知路由與通用 legacy 入口改為落到 `/studio/v2`
- 保留直接開啟 `/studio` 的能力，不把 legacy 主線強制轉成 V2
- 補齊 route contract 測試，固定入口矩陣與 fallback 行為

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-shell`: 新增 `/studio/v2` 作為預設入口與通用 legacy 落點的要求

## Impact

- Affected specs: `datalink-workbench-v2-shell`
- Affected code:
  - Modified: `frontend/src/App.tsx`, `frontend/src/features/datalink/legacyRoutes.ts`, `frontend/tests/unit/workbench-v2/routing.test.tsx`
  - Modified: `docs/technical/studio-surface-inventory/context.json`, `docs/technical/studio-surface-inventory/studio-v2-runtime.md`
  - Removed: none

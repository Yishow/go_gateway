# Studio Surface Inventory: Current State

Last updated: `2026-06-02`

## One-screen Summary

- `/studio`
  - 既有完整版工作台。
  - 目前保留為 fallback。
  - 不要把近期工作重新拉回 `/studio` 首頁重整，除非使用者明確要求。
- `/studio/v2`
  - 目前主重點。
  - 目標是面向使用者、簡單、有指引、好觀察。
  - `/`、unknown routes 與 generic `/datalink` landing 已收斂到 `/studio/v2`。
  - 下一步重點是把多台設備送出真正接上後端。
- `/studio/runtime`
  - 追隨 `/studio/v2`。
  - 角色是 post-setup focused monitor，不是獨立 fleet-first dashboard。
  - 還需要承接 V2 的多台設備送出結果。
- `/test`
  - 獨立工程測試工具入口。
  - legacy `/templates`、`/history`、`/compare`、`/analyzer` 已收斂到 `/test`。
- `/gateway/*`
  - experimental surfaces。
  - 目前不是產品主線。

## What Was Added Recently

1. `studio-surface-inventory` 現在有固定 onboarding 入口：
   - [START_HERE.md](./START_HERE.md)
   - [context.json](./context.json)
   - [CURRENT_STATE.md](./CURRENT_STATE.md)
2. `studio-surface-inventory` 現在有 SQLite changelog：
   - DB: `docs/technical/studio-surface-inventory/changelog.sqlite`
   - Tool: `go run ./cmd/studio_inventory_changelog ...`
3. `AGENTS.md` 與 `CLAUDE.md` 已明確要求：
   - 任務若涉及 `studio` surfaces / `studio-surface-inventory`
   - 先讀 onboarding 入口
   - 不要一開始就掃完整 md/html
4. 目前實作已確認：
   - `frontend/src/App.tsx` 的 `/` 會 redirect 到 `/studio/v2`
   - unknown routes fallback 也會回到 `/studio/v2`
   - `/test` 維持獨立頁而非掛回產品主流程

## Read This First Next Time

如果新對話任務涉及以下任一項：

- `/studio`
- `/studio/v2`
- `/studio/runtime`
- `/test`
- `/gateway/*`
- `docs/technical/studio-surface-inventory/`

先讀：

1. `AGENTS.md`
2. `CLAUDE.md`
3. [START_HERE.md](./START_HERE.md)
4. [context.json](./context.json)
5. [CURRENT_STATE.md](./CURRENT_STATE.md)

只有在這五個入口不足以回答問題時，才展開：

- [studio-mainline.md](./studio-mainline.md)
- [studio-v2-runtime.md](./studio-v2-runtime.md)
- [test-tooling.md](./test-tooling.md)
- [gateway-experiments.md](./gateway-experiments.md)
- [backend-api-registry.md](./backend-api-registry.md)
- [gap-roadmap.md](./gap-roadmap.md)
- [index.html](./index.html)

## What Not To Re-Do

- 不要每次重新盤點全部 surface 與全部 API。
- 不要把 `/studio` 或 `/gateway/*` 當成這一輪優先實作線，除非使用者改變方向。
- 不要把 `/test` 誤當成產品主流程；它是獨立工程測試工具。
- 不要只改 inventory 文件而不寫 changelog。
- 不要直接跳進完整 HTML/長文檔，先用 onboarding 入口縮小範圍。

## Likely Next Work

最可能的下一批工作仍是：

1. `/studio/v2` 多台設備送出契約真正持久化到後端
2. `/studio/runtime` 多台設備交接與 lifecycle contract
3. inventory 文件持續和實作同步

## Fast Commands

```bash
go run ./cmd/studio_inventory_changelog list -limit 10
go run ./cmd/studio_inventory_changelog add \
  -surface "overview" \
  -summary "..." \
  -files "file1,file2" \
  -reason "..."
```

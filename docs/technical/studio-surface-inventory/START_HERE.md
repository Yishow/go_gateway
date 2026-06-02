# Studio Surface Inventory: Start Here

這份文件是 `studio-surface-inventory` 的第一入口，目的是讓接手的 AI 或人類維護者不用先掃完整包文件，就能在幾分鐘內掌握重點。

## 這組文件在做什麼

- 記錄 `go_gateway` 前端 surface 與後端 API 的對照關係。
- 區分哪些頁面是主產品、哪些先暫停、哪些只是工程工具或 prototype。
- 記錄前端實際行為、資料真相、handoff、缺口與後續 roadmap。

## 目前產品決策

1. `/studio`
   - 既有完整版工作台。
   - 需求面應可全面覆蓋。
   - 目前保留為 fallback，不列為近期施作重點。
2. `/studio/v2`
   - 面向使用者的簡化入口。
   - 要簡單、有指引、容易觀察。
   - 這是近期重點施作線，也是目前服務預設入口。
3. `/studio/runtime`
   - 跟著 `/studio/v2` 一起看。
   - 定位是 post-setup focused monitor，不是 fleet-first dashboard。
4. `/test`
   - 獨立工程測試工具入口。
   - 不承載產品主流程，但要維持可用與完整記錄。
5. `/gateway/*`
   - experimental / prototype surfaces。
   - 保留紀錄，不列入主產品優先線。

## 先讀哪些檔案

### 最小必要

1. [README.md](./README.md)
2. [context.json](./context.json)
3. [CURRENT_STATE.md](./CURRENT_STATE.md)
4. [studio-v2-runtime.md](./studio-v2-runtime.md)

### 需要完整脈絡時再讀

1. [studio-mainline.md](./studio-mainline.md)
2. [test-tooling.md](./test-tooling.md)
3. [gateway-experiments.md](./gateway-experiments.md)
4. [backend-api-registry.md](./backend-api-registry.md)
5. [gap-roadmap.md](./gap-roadmap.md)
6. [index.html](./index.html)

## 常見任務對應入口

| 任務 | 先讀 |
| --- | --- |
| 調整 `/studio/v2` 主流程 | `context.json`、`studio-v2-runtime.md`、`backend-api-registry.md` |
| 調整 `/studio/runtime` 監看頁 | `context.json`、`studio-v2-runtime.md`、`backend-api-registry.md` |
| 盤點 `/studio` 與 `/studio/v2` 差異 | `README.md`、`studio-mainline.md`、`studio-v2-runtime.md` |
| 判斷某 API 是否已接前端 | `backend-api-registry.md`、對應 surface 文件 |
| 規劃下一個 change | `gap-roadmap.md`、對應 surface 文件、最新 changelog |

## 接手 AI 的建議流程

1. 先讀 `AGENTS.md` 與 `CLAUDE.md`。
2. 若任務涉及 `studio-surface-inventory` 或 `studio` surfaces，先讀本檔、`context.json`、`CURRENT_STATE.md`。
3. 只有在本檔不足以回答問題時，才展開對應 surface 文件。
4. 若修改這個目錄下的文件，務必同步寫入 `changelog.sqlite`。

## Changelog 入口

- DB 路徑：`docs/technical/studio-surface-inventory/changelog.sqlite`
- 寫入工具：`go run ./cmd/studio_inventory_changelog ...`
- 常用命令：

```bash
go run ./cmd/studio_inventory_changelog list -limit 10
go run ./cmd/studio_inventory_changelog add \
  -surface "studio-v2-runtime" \
  -summary "更新 V2/runtime onboarding context" \
  -files "docs/technical/studio-surface-inventory/START_HERE.md,docs/technical/studio-surface-inventory/context.json" \
  -reason "讓接手 AI 不必先掃完整文件"
```

## 目前最該記住的事

- `/studio/v2` 是目前預設入口與重點線。
- `/test` 是獨立工程測試工具，不是產品主流程。
- `/studio` 保留為 fallback，不是近期主線。
- `studio-surface-inventory` 不只是 route 清單，而是行為、API、狀態與缺口的維護台帳。
- 接手時不要先掃全量 md；先看 `START_HERE.md`、`context.json`、`CURRENT_STATE.md`。

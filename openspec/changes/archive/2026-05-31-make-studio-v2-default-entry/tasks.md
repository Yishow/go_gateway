## 1. 入口矩陣

- [x] 1.1 交付 `Default v2 entry route`：讓 `/` 與 catch-all fallback 進站時都直接落到 `/studio/v2`，並以 `frontend/tests/unit/workbench-v2/routing.test.tsx` 驗證首頁與未知路由的 redirect contract。
- [x] 1.2 交付 `Legacy landing paths converge on /studio/v2`：讓 `/datalink` 與 `/datalink/workbench` 這類通用入口收斂到 `/studio/v2`，並以 routing tests 驗證 generic legacy landing 不再落到 `/studio`。

## 2. preserve direct `/studio` 與回歸

- [x] 2.1 交付 `Preserve direct \`/studio\`` 邊界：明確覆蓋 design decision `preserve direct \`/studio\``，依照 `### Preserve direct \`/studio\`` 與 `### Default route matrix` 保留直接開啟 `/studio` 仍可進入 legacy 主線的行為，並以 routing tests 驗證 `/studio` 不會被自動轉去 `/studio/v2`。
- [x] 2.2 依照 `### Legacy redirect boundary` 固定 task-specific legacy deep links 維持原本導向，不讓這個 change 偷改 scope 外 surface，並以 route matrix review 與 `git diff --check` 驗證變更只落在入口層。

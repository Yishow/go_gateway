## Why

`/studio` 已經不再使用，繼續保留它會讓 legacy frontend surface、測試與文件引用持續被誤認為產品契約。使用者已明確授權立即刪除，並接受以 Git revert 復原；本 change 將刪除限定為一個可獨立驗收的 bounded batch，保留 `/studio/v2`、`/studio/runtime`、`/test` 與 `/gateway/*`。

## What Changes

- 建立刪除前 inventory，證明 `/studio` route registration、legacy-only imports、lazy chunks、assets、測試與文件引用的範圍，以及不屬於 legacy 的 preserved route/assets。
- 立即移除專用 `/studio` route registration 與所有已證明只供 legacy 使用的 frontend imports/chunks/assets、obsolete tests、obsolete docs references；不建立 `/studio` tombstone 或 special redirect。
- 刪除後 `/studio` 必須和任意 unknown route 完全走既有 generic unknown-route policy；不得殘留 legacy-specific redirect/fallback、legacy mount 或 legacy chunk。
- 以獨立 Git commit 作為 rollback anchor；若需要復原，revert 該 commit、重新 build/redeploy，並明確記錄 Git 只能恢復 tracked source/docs，不能恢復 deployment/runtime data 或 browser bookmarks。
- 保留 C 類 release polish：typed backend/runtime errors、safe operator copy、Step 3 preview/SSE truthfulness、Share accessibility/single-flight、CommitProgress failure presentation 與正式 API/inventory/release evidence。這些任務不依賴 `/studio` 退場，也不修改 B 的 Share lifecycle/CAS/concurrency。
- 本 change 不新增 telemetry table/API/config/migration，不寫入或修改資料庫資料，不等待 14-day usage、deprecation window、redirect window 或 A/B evidence 才能刪除 `/studio`。

## Capabilities

### New Capabilities

- `legacy-studio-retirement`: Define the immediately executable, bounded deletion of legacy `/studio`, including pre-delete inventory, generic unknown-route equivalence, Git rollback, and non-target route boundaries.

### Modified Capabilities

- `datalink-workbench-v2-shell`: Make v2 the only product setup entry and remove the dedicated legacy `/studio` route while preserving generic unknown-route behavior and i18n/accessibility polish.
- `datalink-ui`: Keep the primary guided workflow in `/studio/v2` and make live preview failure/reconnect and accessible controls truthful.
- `post-setup-runtime-dashboard`: Require typed runtime errors and visible reconnect/degraded/error states without silent success.
- `embedded-frontend-delivery`: Remove the legacy `/studio` asset graph immediately while preserving `/studio/v2`, `/studio/runtime`, `/test`, and `/gateway/*` identities/assets.
- `api-docs`: Extend generated/documented API contracts for typed preview/runtime errors and the post-deletion generic unknown-route behavior.

## Impact

- Affected specs: `legacy-studio-retirement`, `datalink-workbench-v2-shell`, `datalink-ui`, `post-setup-runtime-dashboard`, `embedded-frontend-delivery`, `api-docs`.
- Affected code and artifacts:
  - Frontend route registry, legacy-only workbench imports/chunks/assets, obsolete tests, and obsolete docs references.
  - Existing generic unknown-route policy and preserved `/studio/v2`, `/studio/runtime`, `/test`, `/gateway/*` route/assets.
  - Git rollback/rebuild/redeploy runbook and pre-delete inventory; no runtime data or database mutation.
  - Typed error/runtime/preview UI, generated Swagger, API registry, release evidence, and studio-surface inventory for the independent C polish batch. Inventory edits require `studio_inventory_changelog` records.
- The `/studio` deletion is an immediate breaking change. Git revert restores tracked source/docs only; it does not restore deployment/runtime data or browser bookmarks.
